package vehicles

import (
	"context"

	"github.com/bjarke-xyz/uber-clone-backend/internal/core"
	"github.com/bjarke-xyz/uber-clone-backend/internal/core/users"
	validation "github.com/go-ozzo/ozzo-validation"
)

type VehicleService struct {
	vehicleRepo VehicleRepository
	userRepo    users.UserRepository
}

func NewService(vehicleRepo VehicleRepository, userRepo users.UserRepository) *VehicleService {
	return &VehicleService{
		vehicleRepo: vehicleRepo,
		userRepo:    userRepo,
	}
}

func (s *VehicleService) enrichWithPositions(ctx context.Context, vehiclesList []Vehicle) error {
	vehicleIds := make([]int64, len(vehiclesList))
	for i, v := range vehiclesList {
		vehicleIds[i] = v.ID
	}
	vehiclePositions, err := s.vehicleRepo.GetVehiclePositions(ctx, vehicleIds)
	if err != nil {
		return err
	}
	vehiclePosMap := make(map[int64]VehiclePosition)
	for _, p := range vehiclePositions {
		vehiclePosMap[p.VehicleID] = p
	}
	for i, vehicle := range vehiclesList {
		pos, ok := vehiclePosMap[vehicle.ID]
		if ok {
			vehiclesList[i].LastRecordedPosition = &pos
		}
	}
	return nil
}

func (s *VehicleService) GetVehicles(ctx context.Context, userID string) ([]Vehicle, error) {
	user, err := s.userRepo.GetByUserID(ctx, userID)
	if err != nil {
		return []Vehicle{}, core.Errorw(core.EINTERNAL, err)
	}
	vehiclesList, err := s.vehicleRepo.GetByOwnerId(ctx, user.ID)
	if err != nil {
		return []Vehicle{}, core.Errorw(core.EINTERNAL, err)
	}
	err = s.enrichWithPositions(ctx, vehiclesList)
	if err != nil {
		return []Vehicle{}, core.Errorw(core.EINTERNAL, err)
	}
	return vehiclesList, nil
}

func (s *VehicleService) GetSimulatedVehicles(ctx context.Context) ([]Vehicle, error) {
	vehicleList, err := s.vehicleRepo.GetSimulatedVehicles(ctx)
	if err != nil {
		return []Vehicle{}, core.Errorw(core.EINTERNAL, err)
	}
	err = s.enrichWithPositions(ctx, vehicleList)
	if err != nil {
		return []Vehicle{}, core.Errorw(core.EINTERNAL, err)
	}
	return vehicleList, nil
}

type UpdateVehiclePositionInput struct {
	Lat     float64 `json:"lat"`
	Lng     float64 `json:"lng"`
	Bearing float32 `json:"bearing"`
	Speed   float32 `json:"speed"`
}

func (i *UpdateVehiclePositionInput) Validate() error {
	return validation.ValidateStruct(i,
		validation.Field(&i.Lat, validation.Required),
		validation.Field(&i.Lng, validation.Required),
	)
}

func (a *VehicleService) UpdateVehiclePosition(ctx context.Context, userID string, vehicleId int64, input *UpdateVehiclePositionInput) error {
	err := input.Validate()
	if err != nil {
		return core.Errorf(core.EINVALID, "invalid input: %v", err)
	}
	user, err := a.userRepo.GetByUserID(ctx, userID)
	if err != nil {
		return core.WrapErr(err)
	}
	vehicle, err := a.vehicleRepo.GetByIdAndOwnerId(ctx, int64(vehicleId), user.ID)
	if err != nil {
		return core.WrapErr(err)
	}
	err = a.vehicleRepo.UpdatePosition(ctx, vehicle.ID, input.Lat, input.Lng)
	if err != nil {
		return core.WrapErr(err)
	}
	// TODO: emit event that position was updated
	return nil
}
