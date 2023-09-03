package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/bjarke-xyz/uber-clone-backend/internal/cfg"
	"github.com/bjarke-xyz/uber-clone-backend/internal/domain"
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

	userRepo    domain.UserRepository
	vehicleRepo domain.VehicleRepository
	rideRepo    domain.RideRepository

	osrClient *service.OpenRouteServiceClient

	broker *broker
}

func NewAPI(ctx context.Context, logger *slog.Logger, cfg *cfg.Cfg, pool *pgxpool.Pool, osrClient *service.OpenRouteServiceClient) *api {
	userRepo := repository.NewPostgresUser(pool)
	vehicleRepo := repository.NewPostgresVehicle(pool)
	rideRepo := repository.NewPostgresRide(pool)

	broker := &broker{
		Notifier:       make(chan []byte, 1),
		newClients:     make(chan chan []byte),
		closingClients: make(chan chan []byte),
		clients:        make(map[chan []byte]bool),
	}
	go broker.listen(logger)

	return &api{
		logger:      logger,
		cfg:         cfg,
		userRepo:    userRepo,
		vehicleRepo: vehicleRepo,
		rideRepo:    rideRepo,
		osrClient:   osrClient,
		broker:      broker,
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

	r.Get("/v1/health", a.healthCheckHandler)

	r.Route("/v1/vehicles", func(r chi.Router) {
		r.Use(a.firebaseJwtVerifier)
		r.Get("/", a.getVehiclesHandler)
		r.Put("/{vehicleID}/position", a.updateVehiclePositionHandler)
	})
	r.Get("/v1/sim/events", a.handleEvents)
	r.Get("/v1/sim-vehicles", a.handleGetSimulatedVehicles)
	r.Get("/v1/sim/logs", a.handleGetRecentLogs)

	r.Route("/v1/rides", func(r chi.Router) {
		r.Use(a.firebaseJwtVerifier)
		r.Get("/mine", a.handleGetMyRideRequests)
		r.Get("/available", a.handleGetAvailableRideRequests)
		r.Post("/", a.handleCreateRideRequest)
		r.Put("/{rideRequestID}/claim", a.handleClaimRideRequest)
		r.Post("/{rideRequestID}/directions", a.handleGetRideDirections)
		r.Put("/{rideRequestID}/finish", a.handleFinishRide)
	})
	r.Get("/v1/sim-rides", a.handleGetSimulatedRides)

	r.Route("/v1/me", func(r chi.Router) {
		r.Use(a.firebaseJwtVerifier)
		r.Get("/user", a.handleGetMyUser)
		r.Post("/log", a.handlePostUserLog)
	})
	r.Get("/v1/sim-users", a.handleGetSimUsers)

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

func (a *api) respond(w http.ResponseWriter, r *http.Request, data any) {
	a.respondStatus(w, r, http.StatusOK, data)
}

func (a *api) respondStatus(w http.ResponseWriter, r *http.Request, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			a.logger.Error("error sending response", "error", err)
		}
	}
}

func (a *api) errorResponse(w http.ResponseWriter, _ *http.Request, status int, err error) {
	a.logger.Error("error", "error", err)
	if errors.Is(err, domain.ErrNotFound) {
		status = http.StatusNotFound
	}
	http.Error(w, err.Error(), status)
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
