package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/bjarke-xyz/uber-clone-backend/internal/core"
	"github.com/bjarke-xyz/uber-clone-backend/internal/core/rides"
	"github.com/samber/lo"
)

type postgresRideRepository struct {
	conn Connection
}

const rideRequestColumns = `id, rider_id, driver_id, from_lat, from_lng, from_name,
			to_lat, to_lng, to_name, state, directions_json_version, directions_json, price, currency, created_at, updated_at`

func (p *postgresRideRepository) fetch(ctx context.Context, query string, args ...interface{}) ([]rides.RideRequest, error) {
	rows, err := p.conn.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rr := make([]rides.RideRequest, 0)
	for rows.Next() {
		var r rides.RideRequest
		if err := rows.Scan(
			&r.ID,
			&r.RiderID,
			&r.DriverID,
			&r.FromLat,
			&r.FromLng,
			&r.FromName,
			&r.ToLat,
			&r.ToLng,
			&r.ToName,
			&r.State,
			&r.DirectionsJsonVersion,
			&r.DirectionsJson,
			&r.Price,
			&r.Currency,
			&r.CreatedAt,
			&r.UpdatedAt,
		); err != nil {
			return nil, err
		}
		if r.DirectionsJsonVersion != nil && r.DirectionsJson != nil && len(*r.DirectionsJson) > 0 {
			switch *r.DirectionsJsonVersion {
			case 1:
				{
					directionsObj := &rides.ORSDirections{}
					err = json.Unmarshal([]byte(*r.DirectionsJson), directionsObj)
					r.DirectionsV1 = directionsObj
				}
			}
			if err != nil {
				return nil, fmt.Errorf("failed to unmarshal v%v directions json: %w", r.DirectionsJsonVersion, err)
			}
		}
		rr = append(rr, r)
	}
	return rr, nil
}

// CreateRequest implements rides.RideRepository.
func (p *postgresRideRepository) CreateRequest(ctx context.Context, ride *rides.RideRequest) error {
	if err := ride.Validate(); err != nil {
		return err
	}
	sql := `INSERT INTO ride_requests (rider_id, driver_id, from_lat, from_lng, from_name,
										to_lat, to_lng, to_name, state, created_at, updated_at) VALUES
			($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id`
	return p.conn.QueryRow(ctx, sql, ride.RiderID, ride.DriverID, ride.FromLat, ride.FromLng, ride.FromName,
		ride.ToLat, ride.ToLng, ride.ToName, ride.State, ride.CreatedAt, ride.UpdatedAt).Scan(&ride.ID)
}

// GetRequests implements rides.RideRepository.
func (p *postgresRideRepository) GetRequests(ctx context.Context, state rides.RideRequestState) ([]rides.RideRequest, error) {
	sql := fmt.Sprintf("SELECT %v FROM ride_requests WHERE state = $1", rideRequestColumns)
	return p.fetch(ctx, sql, state)
}

// GetByID implements rides.RideRepository.
func (p *postgresRideRepository) GetByID(ctx context.Context, id int64) (rides.RideRequest, error) {
	sql := fmt.Sprintf("SELECT %v FROM ride_requests WHERE id = $1", rideRequestColumns)
	ridesList, err := p.fetch(ctx, sql, id)
	if err != nil {
		return rides.RideRequest{}, err
	}
	if len(ridesList) == 0 {
		return rides.RideRequest{}, core.Errorf(core.ENOTFOUND, "ride with id %v not found", id)
	}
	return ridesList[0], nil
}

// GetByUserID implements rides.RideRepository.
func (p *postgresRideRepository) GetByUserID(ctx context.Context, userId int64) ([]rides.RideRequest, error) {
	sql := fmt.Sprintf("SELECT %v FROM ride_requests WHERE rider_id = $1 OR driver_id = $1", rideRequestColumns)
	return p.fetch(ctx, sql, userId)
}

// GetByUserIDs implements rides.RideRepository.
func (p *postgresRideRepository) GetByUserIDs(ctx context.Context, userIds []int64, states []rides.RideRequestState) ([]rides.RideRequest, error) {
	ridesList := make([]rides.RideRequest, 0)
	if len(userIds) == 0 {
		return ridesList, nil
	}
	userIdsStrs := lo.Map(userIds, func(x int64, index int) string {
		return strconv.FormatInt(x, 10)
	})
	userIdsStr := strings.Join(userIdsStrs, ",")
	statesWhere := ""
	if len(states) > 0 {
		statesStrs := lo.Map(states, func(x rides.RideRequestState, index int) string { return fmt.Sprintf("%v", x) })
		stateStr := strings.Join(statesStrs, ",")
		statesWhere = fmt.Sprintf(" AND state IN (%v)", stateStr)
	}
	sql := fmt.Sprintf("SELECT %v FROM ride_requests WHERE (rider_id IN (%v) OR driver_id IN (%v)) %v ORDER BY created_at DESC", rideRequestColumns, userIdsStr, userIdsStr, statesWhere)
	return p.fetch(ctx, sql)
}

// UpdateRequestState implements rides.RideRepository.
func (p *postgresRideRepository) UpdateRequestState(ctx context.Context, rideId int64, state rides.RideRequestState) error {
	sql := "UPDATE ride_requests SET state = $2, updated_at = $3  WHERE id = $1"
	_, err := p.conn.Exec(ctx, sql, rideId, state, time.Now().UTC())
	return err
}

// ClaimRequest implements rides.RideRepository.
func (p *postgresRideRepository) ClaimRequest(ctx context.Context, requestId int64, driverID int64) error {
	sql := "UPDATE ride_requests SET state = $2, updated_at = $3, driver_id = $4 WHERE id = $1"
	_, err := p.conn.Exec(ctx, sql, requestId, rides.RiderRequestStateAccepted, time.Now().UTC(), driverID)
	return err
}

// UpdateRideDirections implements rides.RideRepository.
func (p *postgresRideRepository) UpdateRideDirections(ctx context.Context, requestId int64, directionsVersion int, directions *rides.ORSDirections, price int) error {
	sql := "UPDATE ride_requests SET directions_json_version = $2, directions_json = $3, price = $4 WHERE id = $1"
	directionsBytes, err := json.Marshal(directions)
	if err != nil {
		return fmt.Errorf("failed to marshal directions: %w", err)
	}
	directionsStr := string(directionsBytes)
	_, err = p.conn.Exec(ctx, sql, requestId, directionsVersion, directionsStr, price)
	return err
}

func NewPostgresRide(conn Connection) rides.RideRepository {
	return &postgresRideRepository{conn: conn}
}
