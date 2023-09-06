package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/bjarke-xyz/uber-clone-backend/internal/auth"
	"github.com/bjarke-xyz/uber-clone-backend/internal/cfg"
	"github.com/bjarke-xyz/uber-clone-backend/internal/core"
	"github.com/bjarke-xyz/uber-clone-backend/internal/core/payments"
	"github.com/bjarke-xyz/uber-clone-backend/internal/core/rides"
	"github.com/bjarke-xyz/uber-clone-backend/internal/core/users"
	"github.com/bjarke-xyz/uber-clone-backend/internal/core/vehicles"
	"github.com/bjarke-xyz/uber-clone-backend/internal/repository"
	"github.com/bjarke-xyz/uber-clone-backend/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	sseClientGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "uberclone_sse_clients",
		Help: "The total number of sse clients",
	})
)

var (
	TokenCtxKey = &contextKey{"Token"}
	ErrorCtxKey = &contextKey{"Error"}
)

type api struct {
	logger *slog.Logger
	cfg    *cfg.Cfg

	authClient auth.AuthClient

	paymentsService *payments.PaymentsService
	rideService     *rides.RideService
	userService     *users.UserService
	vehicleService  *vehicles.VehicleService

	userRepo    users.UserRepository
	vehicleRepo vehicles.VehicleRepository
	rideRepo    rides.RideRepository

	broker *broker
}

func NewAPI(ctx context.Context, logger *slog.Logger, cfg *cfg.Cfg, authClient auth.AuthClient, pool *pgxpool.Pool, osrClient *service.OpenRouteServiceClient) *api {
	userRepo := repository.NewPostgresUser(pool)
	vehicleRepo := repository.NewPostgresVehicle(pool)
	rideRepo := repository.NewPostgresRide(pool)

	paymentsService := payments.NewService()
	rideService := rides.NewService(rideRepo, userRepo, osrClient, paymentsService)
	userService := users.NewService(userRepo)
	vehicleService := vehicles.NewService(vehicleRepo, userRepo)

	broker := &broker{
		Notifier:       make(chan []byte, 1),
		newClients:     make(chan chan []byte),
		closingClients: make(chan chan []byte),
		clients:        make(map[chan []byte]bool),
	}
	go broker.listen(logger)

	return &api{
		logger:          logger,
		cfg:             cfg,
		authClient:      authClient,
		paymentsService: paymentsService,
		rideService:     rideService,
		userService:     userService,
		vehicleService:  vehicleService,
		userRepo:        userRepo,
		vehicleRepo:     vehicleRepo,
		rideRepo:        rideRepo,
		broker:          broker,
	}
}

func (a *api) Server(port int) *http.Server {
	return &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: a.Routes(),
	}
}

func (a *api) Routes() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.AllowAll().Handler)

	r.Get("/v1/health", a.requestWrapper(a.healthCheckHandler))

	r.Route("/v1/vehicles", func(r chi.Router) {
		r.Use(a.firebaseJwtVerifier)
		r.Get("/", a.requestWrapper(a.getVehiclesHandler))
		r.Put("/{vehicleID}/position", a.requestWrapper(a.updateVehiclePositionHandler))
	})
	r.Get("/v1/sim/events", a.handleEvents)
	r.Get("/v1/sim-vehicles", a.requestWrapper(a.handleGetSimulatedVehicles))
	r.Get("/v1/sim/logs", a.requestWrapper(a.handleGetRecentLogs))

	r.Route("/v1/rides", func(r chi.Router) {
		r.Use(a.firebaseJwtVerifier)
		r.Get("/mine", a.requestWrapper(a.handleGetMyRideRequests))
		r.Get("/available", a.requestWrapper(a.handleGetAvailableRideRequests))
		r.Post("/", a.requestWrapper(a.handleCreateRideRequest))
		r.Put("/{rideRequestID}/claim", a.requestWrapper(a.handleClaimRideRequest))
		r.Post("/{rideRequestID}/directions", a.requestWrapper(a.handleGetRideDirections))
		r.Put("/{rideRequestID}/finish", a.requestWrapper(a.handleFinishRide))
	})
	r.Get("/v1/sim-rides", a.requestWrapper(a.handleGetSimulatedRides))

	r.Route("/v1/me", func(r chi.Router) {
		r.Use(a.firebaseJwtVerifier)
		r.Get("/user", a.requestWrapper(a.handleGetMyUser))
		r.Post("/log", a.requestWrapper(a.handlePostUserLog))
	})
	r.Get("/v1/sim-users", a.requestWrapper(a.handleGetSimUsers))

	r.Route("/v1/payments", func(r chi.Router) {
		r.Get("/currencies", a.requestWrapper(a.handleGetCurrencies))
	})

	return r
}

type broker struct {

	// Events are pushed to this channel by the main events-gathering routine
	Notifier chan []byte

	// New client connections
	newClients chan chan []byte

	// Closed client connections
	closingClients chan chan []byte

	// Client connections registry
	clients map[chan []byte]bool
}

func (a *api) respond(w http.ResponseWriter, r *http.Request, data any) error {
	return a.respondStatus(w, r, http.StatusOK, data)
}

func (a *api) respondStatus(w http.ResponseWriter, r *http.Request, status int, data any) error {
	if data != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			a.logger.Error("error sending response", "error", err)
		}
		return err
	} else {
		w.WriteHeader(status)
	}
	return nil
}

func (a *api) requestWrapper(handler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := handler(r.Context(), w, r); err != nil {
			a.logger.Error("error", "error", err)
			status := http.StatusInternalServerError
			var e *core.Error
			if errors.As(err, &e) {
				status = mapToHttpStatus(e)
			}
			http.Error(w, err.Error(), status)
		}
	}
}

func mapToHttpStatus(err *core.Error) int {
	switch err.Code {
	case core.ECONFLICT:
		return http.StatusConflict
	case core.EINTERNAL:
		return http.StatusInternalServerError
	case core.EINVALID:
		return http.StatusBadRequest
	case core.ENOTFOUND:
		return http.StatusNotFound
	case core.ENOTIMPLEMENTED:
		return http.StatusNotImplemented
	case core.EUNAUTHORIZED:
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}

func decodeBody(r io.Reader, input any) error {
	if err := json.NewDecoder(r).Decode(input); err != nil {
		return core.Errorf(core.EINVALID, "failed to decode body: %v", err)
	}
	return nil
}

func urlParamInt(r *http.Request, key string) (int64, error) {
	valueStr := chi.URLParam(r, key)
	valueInt, err := strconv.ParseInt(valueStr, 10, 64)
	if err != nil {
		return 0, core.Errorw(core.EINVALID, err)
	}
	return valueInt, nil
}

func queryParamFloat(r *http.Request, key string) (float64, bool, error) {
	query := r.URL.Query()
	valueStr := query.Get(key)
	if valueStr == "" {
		return 0, false, nil
	}
	valueFloat, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		return 0, true, core.Errorw(core.EINVALID, err)
	}
	return valueFloat, true, nil
}

func (broker *broker) listen(logger *slog.Logger) {
	for {
		select {
		case s := <-broker.newClients:

			// A new client has connected.
			// Register their message channel
			broker.clients[s] = true
			logger.Info("Client added", "clients", len(broker.clients))
			sseClientGauge.Inc()
		case s := <-broker.closingClients:

			// A client has dettached and we want to
			// stop sending them messages.
			delete(broker.clients, s)
			logger.Info("Removed client", "clients", len(broker.clients))
			sseClientGauge.Dec()
		case event := <-broker.Notifier:

			// We got a new event from the outside!
			// Send event to all connected clients
			for clientMessageChan := range broker.clients {
				clientMessageChan <- event
			}
		}
	}

}
