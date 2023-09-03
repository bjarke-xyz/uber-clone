package api

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/bjarke-xyz/uber-clone-backend/internal/domain"
	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/samber/lo"
)

func (a *api) handleGetMyRideRequests(w http.ResponseWriter, r *http.Request) {
	token, _ := TokenFromContext(r.Context())
	user, err := a.userRepo.GetByUserID(r.Context(), token.Subject)
	if err != nil {
		a.errorResponse(w, r, http.StatusInternalServerError, err)
		return
	}
	rideRequests, err := a.rideRepo.GetByUserID(r.Context(), user.ID)
	if err != nil {
		a.errorResponse(w, r, http.StatusInternalServerError, err)
		return
	}
	a.respond(w, r, rideRequests)
}

func (a *api) handleGetAvailableRideRequests(w http.ResponseWriter, r *http.Request) {
	rideRequests, err := a.rideRepo.GetRequests(r.Context(), domain.RiderRequestStateAvailable)
	if err != nil {
		a.errorResponse(w, r, http.StatusInternalServerError, err)
		return
	}
	a.respond(w, r, rideRequests)
}

type CreateRideInput struct {
	FromLat  float64 `json:"fromLat"`
	FromLng  float64 `json:"fromLng"`
	FromName string  `json:"fromName"`
	ToLat    float64 `json:"toLat"`
	ToLng    float64 `json:"toLng"`
	ToName   string  `json:"toName"`
}

func (c *CreateRideInput) Validate() error {
	return validation.ValidateStruct(c,
		validation.Field(&c.FromLat, validation.Required),
		validation.Field(&c.FromLng, validation.Required),
		validation.Field(&c.FromName, validation.Required),
		validation.Field(&c.ToLat, validation.Required),
		validation.Field(&c.ToLng, validation.Required),
		validation.Field(&c.ToName, validation.Required),
	)
}

func (a *api) handleCreateRideRequest(w http.ResponseWriter, r *http.Request) {
	token, _ := TokenFromContext(r.Context())
	input := &CreateRideInput{}
	if err := json.NewDecoder(r.Body).Decode(input); err != nil {
		a.errorResponse(w, r, http.StatusInternalServerError, fmt.Errorf("failed to decode body: %w", err))
		return
	}
	if err := input.Validate(); err != nil {
		a.errorResponse(w, r, http.StatusInternalServerError, err)
		return
	}
	user, err := a.userRepo.GetByUserID(r.Context(), token.Subject)
	if err != nil {
		a.errorResponse(w, r, http.StatusInternalServerError, err)
		return
	}

	rideRequest := &domain.RideRequest{
		RiderID:   user.ID,
		FromLat:   input.FromLat,
		FromLng:   input.FromLng,
		FromName:  input.FromName,
		ToLat:     input.ToLat,
		ToLng:     input.ToLng,
		ToName:    input.ToName,
		State:     domain.RiderRequestStateAvailable,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	err = a.rideRepo.CreateRequest(r.Context(), rideRequest)
	if err != nil {
		a.errorResponse(w, r, http.StatusInternalServerError, err)
		return
	}
	rideReq, err := a.rideRepo.GetByID(r.Context(), rideRequest.ID)
	if err != nil {
		a.errorResponse(w, r, http.StatusInternalServerError, err)
		return
	}
	a.respond(w, r, rideReq)
}

func (a *api) handleClaimRideRequest(w http.ResponseWriter, r *http.Request) {
	token, _ := TokenFromContext(r.Context())
	rideRequestIdStr := chi.URLParam(r, "rideRequestID")
	rideRequestId, err := strconv.ParseInt(rideRequestIdStr, 10, 64)
	if err != nil {
		a.errorResponse(w, r, http.StatusInternalServerError, err)
		return
	}

	user, err := a.userRepo.GetByUserID(r.Context(), token.Subject)
	if err != nil {
		a.errorResponse(w, r, http.StatusInternalServerError, err)
		return
	}

	rideReq, err := a.rideRepo.GetByID(r.Context(), rideRequestId)
	if err != nil {
		a.errorResponse(w, r, http.StatusInternalServerError, err)
		return
	}

	if rideReq.State != domain.RiderRequestStateAvailable {
		a.errorResponse(w, r, http.StatusBadRequest, domain.ErrCannotClaimNonAvailableRide)
		return
	}

	err = a.rideRepo.ClaimRequest(r.Context(), rideReq.ID, user.ID)
	if err != nil {
		a.errorResponse(w, r, http.StatusInternalServerError, err)
		return
	}
}

func (a *api) handleGetRideDirections(w http.ResponseWriter, r *http.Request) {
	// token, _ := TokenFromContext(r.Context())
	rideRequestIdStr := chi.URLParam(r, "rideRequestID")
	rideRequestId, err := strconv.ParseInt(rideRequestIdStr, 10, 64)
	if err != nil {
		a.errorResponse(w, r, http.StatusInternalServerError, err)
		return
	}

	var optionalStartLat float64 = 0
	var optionalStartLng float64 = 0
	query := r.URL.Query()
	optionalStartLatStr := query.Get("startLat")
	optionalStartLngStr := query.Get("startLng")
	if optionalStartLatStr != "" && optionalStartLngStr != "" {
		optionalStartLat, err = strconv.ParseFloat(optionalStartLatStr, 64)
		if err != nil {
			a.errorResponse(w, r, http.StatusBadRequest, err)
			return
		}
		optionalStartLng, err = strconv.ParseFloat(optionalStartLngStr, 64)
		if err != nil {
			a.errorResponse(w, r, http.StatusBadRequest, err)
			return
		}
	}

	// user, err := a.userRepo.GetByUserID(r.Context(), token.Subject)
	// if err != nil {
	// 	a.errorResponse(w, r, http.StatusInternalServerError, err)
	// 	return
	// }

	rideReq, err := a.rideRepo.GetByID(r.Context(), rideRequestId)
	if err != nil {
		a.errorResponse(w, r, http.StatusInternalServerError, err)
		return
	}
	// if rideReq.DriverID == nil || *rideReq.DriverID != user.ID || rideReq.RiderID != user.ID {
	// 	err = fmt.Errorf("cannot get directions of ride you have not claimed or are not a rider of")
	// 	a.errorResponse(w, r, http.StatusBadRequest, err)
	// 	return
	// }

	directions := rideReq.DirectionsV1
	if directions == nil {
		locations := make([][]float64, 0)
		if optionalStartLat > 0 && optionalStartLng > 0 {
			locations = append(locations, []float64{optionalStartLng, optionalStartLat})
		}
		locations = append(locations, []float64{rideReq.FromLng, rideReq.FromLat}, []float64{rideReq.ToLng, rideReq.ToLat})
		directions, err = a.osrClient.GetDirections(locations)
		if err != nil {
			a.errorResponse(w, r, http.StatusInternalServerError, err)
			return
		}
		totalDistance := 0
		for _, route := range directions.Routes {
			totalDistance += int(math.Ceil(route.Summary.Distance))
		}
		price := domain.CalculatePrice(totalDistance)
		err = a.rideRepo.UpdateRideDirections(r.Context(), rideRequestId, 1, directions, price)
		if err != nil {
			a.errorResponse(w, r, http.StatusInternalServerError, err)
			return
		}
	}

	a.respond(w, r, directions)
}

func (a *api) handleFinishRide(w http.ResponseWriter, r *http.Request) {
	token, _ := TokenFromContext(r.Context())
	rideRequestIdStr := chi.URLParam(r, "rideRequestID")
	rideRequestId, err := strconv.ParseInt(rideRequestIdStr, 10, 64)
	if err != nil {
		a.errorResponse(w, r, http.StatusInternalServerError, err)
		return
	}

	user, err := a.userRepo.GetByUserID(r.Context(), token.Subject)
	if err != nil {
		a.errorResponse(w, r, http.StatusInternalServerError, err)
		return
	}

	rideReq, err := a.rideRepo.GetByID(r.Context(), rideRequestId)
	if err != nil {
		a.errorResponse(w, r, http.StatusInternalServerError, err)
		return
	}
	if rideReq.DriverID == nil || *rideReq.DriverID != user.ID {
		err = fmt.Errorf("cannot change ride")
		a.errorResponse(w, r, http.StatusBadRequest, err)
		return
	}

	a.rideRepo.UpdateRequestState(r.Context(), rideReq.ID, domain.RiderRequestStateFinished)
	w.WriteHeader(http.StatusOK)
}

func (a *api) handleGetSimulatedRides(w http.ResponseWriter, r *http.Request) {
	simUsers, err := a.userRepo.GetSimulatedUsers(r.Context())
	if err != nil {
		a.errorResponse(w, r, http.StatusInternalServerError, err)
		return
	}
	simUserIds := lo.Map(simUsers, func(item domain.User, index int) int64 { return item.ID })
	states := []domain.RideRequestState{domain.RiderRequestStateAvailable, domain.RiderRequestStateAccepted, domain.RiderRequestStateInProgress}
	rides, err := a.rideRepo.GetByUserIDs(r.Context(), simUserIds, states)
	if err != nil {
		a.errorResponse(w, r, http.StatusInternalServerError, err)
		return
	}
	a.respond(w, r, rides)
}
