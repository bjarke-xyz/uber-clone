package api

import (
	"context"
	"net/http"

	"github.com/bjarke-xyz/uber-clone-backend/internal/core"
	"github.com/bjarke-xyz/uber-clone-backend/internal/core/rides"
)

func (a *api) handleGetMyRideRequests(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	token, _ := TokenFromContext(r.Context())
	rideRequests, err := a.rideService.GetRideRequestsByUserID(ctx, token.Subject)
	if err != nil {
		return err
	}
	return a.respond(w, r, rideRequests)
}

func (a *api) handleGetAvailableRideRequests(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	rideRequests, err := a.rideService.GetAvailableRideRequests(ctx)
	if err != nil {
		return core.Errorw(core.EINTERNAL, err)
	}
	return a.respond(w, r, rideRequests)
}

func (a *api) handleCreateRideRequest(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	token, _ := TokenFromContext(r.Context())
	input := &rides.CreateRideInput{}
	if err := decodeBody(r.Body, input); err != nil {
		return err
	}
	rideReq, err := a.rideService.CreateRideRequest(ctx, token.Subject, input)
	if err != nil {
		return err
	}
	return a.respond(w, r, rideReq)
}

func (a *api) handleClaimRideRequest(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	token, _ := TokenFromContext(r.Context())
	rideRequestId, err := urlParamInt(r, "rideRequestID")
	if err != nil {
		return err
	}
	err = a.rideService.ClaimRideRequest(ctx, token.Subject, rideRequestId)
	if err != nil {
		return err
	}
	return a.respondStatus(w, r, http.StatusNoContent, nil)
}

func (a *api) handleGetRideDirections(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	rideRequestId, err := urlParamInt(r, "rideRequestID")
	if err != nil {
		return err
	}
	optionalStartLat, startLatOk, err := queryParamFloat(r, "startLat")
	if startLatOk && err != nil {
		return err
	}
	optionalStartLng, startLngOk, err := queryParamFloat(r, "startLng")
	if startLngOk && err != nil {
		return err
	}
	directions, err := a.rideService.GetRideDirections(ctx, rideRequestId, optionalStartLat, optionalStartLng)
	if err != nil {
		return err
	}
	return a.respond(w, r, directions)
}

func (a *api) handleFinishRide(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	token, _ := TokenFromContext(r.Context())
	rideRequestId, err := urlParamInt(r, "rideRequestID")
	if err != nil {
		return err
	}
	err = a.rideService.FinishRide(ctx, token.Subject, rideRequestId)
	if err != nil {
		return err
	}
	return a.respondStatus(w, r, http.StatusNoContent, nil)
}

func (a *api) handleGetSimulatedRides(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	rides, err := a.rideService.GetSimulatedRides(ctx)
	if err != nil {
		return err
	}
	return a.respond(w, r, rides)
}
