package rides

import (
	"context"
	"math"
	"time"

	"github.com/bjarke-xyz/uber-clone-backend/internal/core"
	"github.com/bjarke-xyz/uber-clone-backend/internal/core/payments"
	"github.com/bjarke-xyz/uber-clone-backend/internal/core/users"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/samber/lo"
)

type RideService struct {
	rideRepo           RideRepository
	userRepo           users.UserRepository
	routeServiceClient RouteServiceClient
	paymentsService    *payments.PaymentsService
}

func NewService(rideRepo RideRepository, userRepo users.UserRepository, routeServiceClient RouteServiceClient, paymentsService *payments.PaymentsService) *RideService {
	return &RideService{
		rideRepo:           rideRepo,
		userRepo:           userRepo,
		routeServiceClient: routeServiceClient,
		paymentsService:    paymentsService,
	}
}

func (r *RideService) GetRideRequestsByUserID(ctx context.Context, userID string) ([]RideRequest, error) {
	user, err := r.userRepo.GetByUserID(ctx, userID)
	if err != nil {
		return []RideRequest{}, core.Errorw(core.EINTERNAL, err)
	}
	rideRequests, err := r.rideRepo.GetByUserID(ctx, user.ID)
	if err != nil {
		return []RideRequest{}, core.Errorw(core.EINTERNAL, err)
	}
	return rideRequests, nil
}

func (r *RideService) GetAvailableRideRequests(ctx context.Context) ([]RideRequest, error) {
	return r.rideRepo.GetRequests(ctx, RiderRequestStateAvailable)
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
func (r *RideService) CreateRideRequest(ctx context.Context, userID string, input *CreateRideInput) (RideRequest, error) {
	if err := input.Validate(); err != nil {
		return RideRequest{}, core.Errorw(core.EINTERNAL, err)
	}
	user, err := r.userRepo.GetByUserID(ctx, userID)
	if err != nil {
		return RideRequest{}, core.Errorw(core.EINTERNAL, err)
	}

	rideRequest := &RideRequest{
		RiderID:   user.ID,
		FromLat:   input.FromLat,
		FromLng:   input.FromLng,
		FromName:  input.FromName,
		ToLat:     input.ToLat,
		ToLng:     input.ToLng,
		ToName:    input.ToName,
		State:     RiderRequestStateAvailable,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	err = r.rideRepo.CreateRequest(ctx, rideRequest)
	if err != nil {
		return RideRequest{}, core.Errorw(core.EINTERNAL, err)
	}
	rideReq, err := r.rideRepo.GetByID(ctx, rideRequest.ID)
	if err != nil {
		return RideRequest{}, core.Errorw(core.EINTERNAL, err)
	}
	return rideReq, nil
}

func (r *RideService) ClaimRideRequest(ctx context.Context, userID string, rideRequestId int64) error {
	user, err := r.userRepo.GetByUserID(ctx, userID)
	if err != nil {
		return core.Errorw(core.EINTERNAL, err)
	}

	rideReq, err := r.rideRepo.GetByID(ctx, rideRequestId)
	if err != nil {
		return core.Errorw(core.EINTERNAL, err)
	}

	if rideReq.State != RiderRequestStateAvailable {
		return core.Errorf(core.EINVALID, "cannot claim non-available ride")
	}

	err = r.rideRepo.ClaimRequest(ctx, rideReq.ID, user.ID)
	if err != nil {
		return core.Errorw(core.EINTERNAL, err)
	}
	return nil
}

func (r *RideService) GetRideDirections(ctx context.Context, rideRequestId int64, optionalStartLat float64, optionalStartLng float64) (*ORSDirections, error) {
	// user, err := a.userRepo.GetByUserID(r.Context(), token.Subject)
	// if err != nil {
	// 	a.errorResponse(w, r, http.StatusInternalServerError, err)
	// 	return
	// }

	rideReq, err := r.rideRepo.GetByID(ctx, rideRequestId)
	if err != nil {
		return nil, core.Errorw(core.EINTERNAL, err)
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
		directions, err = r.routeServiceClient.GetDirections(locations)
		if err != nil {
			return nil, core.Errorw(core.EINTERNAL, err)
		}
		totalDistance := 0
		for _, route := range directions.Routes {
			totalDistance += int(math.Ceil(route.Summary.Distance))
		}
		price := r.paymentsService.CalculatePrice(totalDistance)
		err = r.rideRepo.UpdateRideDirections(ctx, rideRequestId, 1, directions, price)
		if err != nil {
			return nil, core.Errorw(core.EINTERNAL, err)
		}
	}
	return directions, nil
}

func (r *RideService) FinishRide(ctx context.Context, userID string, rideRequestId int64) error {
	user, err := r.userRepo.GetByUserID(ctx, userID)
	if err != nil {
		return core.Errorw(core.EINTERNAL, err)
	}

	rideReq, err := r.rideRepo.GetByID(ctx, rideRequestId)
	if err != nil {
		return core.Errorw(core.EINTERNAL, err)
	}
	if rideReq.DriverID == nil || *rideReq.DriverID != user.ID {
		return core.Errorf(core.EINVALID, "cannot change ride")
	}

	err = r.rideRepo.UpdateRequestState(ctx, rideReq.ID, RiderRequestStateFinished)
	if err != nil {
		return core.Errorw(core.EINTERNAL, err)
	}
	return nil
}

func (r *RideService) GetSimulatedRides(ctx context.Context) ([]RideRequest, error) {
	simUsers, err := r.userRepo.GetSimulatedUsers(ctx)
	if err != nil {
		return []RideRequest{}, core.Errorw(core.EINTERNAL, err)
	}
	simUserIds := lo.Map(simUsers, func(item users.User, index int) int64 { return item.ID })
	states := []RideRequestState{RiderRequestStateAvailable, RiderRequestStateAccepted, RiderRequestStateInProgress}
	rides, err := r.rideRepo.GetByUserIDs(ctx, simUserIds, states)
	if err != nil {
		return []RideRequest{}, core.Errorw(core.EINTERNAL, err)
	}
	return rides, nil
}
