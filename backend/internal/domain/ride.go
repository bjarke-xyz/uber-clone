package domain

import (
	"context"
	"errors"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
)

var ErrCannotClaimNonAvailableRide = errors.New("cannot claim non-avaiable ride")

type RideRequestState int

const (
	RiderRequestStateAvailable RideRequestState = iota
	RiderRequestStateAccepted
	RiderRequestStateInProgress
	RiderRequestStateFinished
)

type ORSDirections struct {
	Bbox   []float64 `json:"bbox"`
	Routes []struct {
		Summary struct {
			Distance float64 `json:"distance"`
			Duration float64 `json:"duration"`
		} `json:"summary"`
		Segments []struct {
			Distance float64 `json:"distance"`
			Duration float64 `json:"duration"`
			Steps    []struct {
				Distance    float64 `json:"distance"`
				Duration    float64 `json:"duration"`
				Type        int     `json:"type"`
				Instruction string  `json:"instruction"`
				Name        string  `json:"name"`
				WayPoints   []int   `json:"way_points"`
				Maneuver    struct {
					Location      []float64 `json:"location"`
					BearingBefore int       `json:"bearing_before"`
					BearingAfter  int       `json:"bearing_after"`
				} `json:"maneuver"`
			} `json:"steps"`
		} `json:"segments"`
		Bbox      []float64 `json:"bbox"`
		Geometry  string    `json:"geometry"`
		WayPoints []int     `json:"way_points"`
		Legs      []any     `json:"legs"`
	} `json:"routes"`
	Metadata struct {
		Attribution string `json:"attribution"`
		Service     string `json:"service"`
		Timestamp   int64  `json:"timestamp"`
		Query       struct {
			Coordinates [][]float64 `json:"coordinates"`
			Profile     string      `json:"profile"`
			Format      string      `json:"format"`
		} `json:"query"`
		Engine struct {
			Version   string    `json:"version"`
			BuildDate time.Time `json:"build_date"`
			GraphDate time.Time `json:"graph_date"`
		} `json:"engine"`
	} `json:"metadata"`
}

type RideRequest struct {
	ID int64 `json:"id"`

	RiderID  int64  `json:"riderId"`
	DriverID *int64 `json:"driverId"`

	FromLat  float64 `json:"fromLat"`
	FromLng  float64 `json:"fromLng"`
	FromName string  `json:"fromName"`

	ToLat  float64 `json:"toLat"`
	ToLng  float64 `json:"toLng"`
	ToName string  `json:"toName"`

	State RideRequestState `json:"state"`

	DirectionsJsonVersion *int           `json:"directionsVersion"`
	DirectionsJson        *string        `json:"-"`
	DirectionsV1          *ORSDirections `json:"directions"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (r *RideRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.FromLat, validation.Required),
		validation.Field(&r.FromLng, validation.Required),
		validation.Field(&r.FromName, validation.Required),
		validation.Field(&r.ToLat, validation.Required),
		validation.Field(&r.ToLng, validation.Required),
		validation.Field(&r.ToName, validation.Required),
	)
}

type RideRepository interface {
	GetRequests(context.Context, RideRequestState) ([]RideRequest, error)
	GetByID(context.Context, int64) (RideRequest, error)
	GetByUserID(context.Context, int64) ([]RideRequest, error)
	GetByUserIDs(ctx context.Context, userIds []int64, states []RideRequestState) ([]RideRequest, error)
	CreateRequest(context.Context, *RideRequest) error
	UpdateRequestState(context.Context, int64, RideRequestState) error
	ClaimRequest(ctx context.Context, requestID int64, driverID int64) error
	UpdateRideDirections(ctx context.Context, requestId int64, directionsVersion int, directions *ORSDirections) error
}
