package vehicles

import (
	"context"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
)

type Vehicle struct {
	ID int64

	RegistrationCountry string
	RegistrationNumber  string
	OwnerID             int64
	Icon                string

	LastRecordedPosition *VehiclePosition `json:"lastRecordedPosition"`
}

type VehiclePosition struct {
	ID         int64     `json:"id"`
	VehicleID  int64     `json:"vehicleId"`
	Lat        float64   `json:"lat"`
	Lng        float64   `json:"lng"`
	RecordedAt time.Time `json:"recordedAt"`
	Bearing    float32   `json:"bearing"`
	Speed      float32   `json:"speed"`
}

func (v *Vehicle) Validate() error {
	return validation.ValidateStruct(v,
		validation.Field(&v.RegistrationCountry, validation.Required),
		validation.Field(&v.RegistrationNumber, validation.Required),
		validation.Field(&v.Icon, validation.Required),
	)
}

type VehicleRepository interface {
	GetByID(context.Context, int64) (Vehicle, error)
	GetByOwnerId(context.Context, int64) ([]Vehicle, error)
	GetByIdAndOwnerId(ctx context.Context, vehicleId int64, userId int64) (Vehicle, error)
	GetSimulatedVehicles(ctx context.Context) ([]Vehicle, error)

	CreateOrUpdate(context.Context, *Vehicle) error
	Delete(context.Context, int64) error

	GetVehiclePositions(ctx context.Context, vehicleIds []int64) ([]VehiclePosition, error)
	UpdatePosition(ctx context.Context, vehicleId int64, lat float64, lng float64) error
}
