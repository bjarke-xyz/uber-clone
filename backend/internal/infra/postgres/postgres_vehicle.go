package postgres

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/bjarke-xyz/uber-clone-backend/internal/core"
	"github.com/bjarke-xyz/uber-clone-backend/internal/core/vehicles"
	"github.com/samber/lo"
)

type postgresVehicleRepository struct {
	conn Connection
}

func NewPostgresVehicle(conn Connection) vehicles.VehicleRepository {
	return &postgresVehicleRepository{conn: conn}
}

func (p *postgresVehicleRepository) fetch(ctx context.Context, query string, args ...interface{}) ([]vehicles.Vehicle, error) {
	rows, err := p.conn.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	vv := make([]vehicles.Vehicle, 0)
	for rows.Next() {
		var v vehicles.Vehicle
		if err := rows.Scan(
			&v.ID,
			&v.RegistrationCountry,
			&v.RegistrationNumber,
			&v.OwnerID,
			&v.Icon,
		); err != nil {
			return nil, err
		}
		vv = append(vv, v)
	}
	return vv, nil
}

// GetByIdAndOwnerId implements vehicles.VehicleRepository.
func (p *postgresVehicleRepository) GetByIdAndOwnerId(ctx context.Context, vehicleId int64, userId int64) (vehicles.Vehicle, error) {
	sql := "SELECT id, registration_country, registration_number, owner_id, icon FROM vehicles WHERE id = $1 AND owner_id = $2"
	vehicleList, err := p.fetch(ctx, sql, vehicleId, userId)
	if err != nil {
		return vehicles.Vehicle{}, err
	}
	if len(vehicleList) == 0 {
		return vehicles.Vehicle{}, core.Errorf(core.ENOTFOUND, "vehicle not found")
	}
	return vehicleList[0], nil
}

// CreateOrUpdate implements vehicles.VehicleRepository.
func (p *postgresVehicleRepository) CreateOrUpdate(ctx context.Context, v *vehicles.Vehicle) error {
	panic("unimplemented")
}

// Delete implements vehicles.VehicleRepository.
func (p *postgresVehicleRepository) Delete(ctx context.Context, id int64) error {
	panic("unimplemented")
}

// GetByID implements vehicles.VehicleRepository.
func (p *postgresVehicleRepository) GetByID(ctx context.Context, id int64) (vehicles.Vehicle, error) {
	sql := "SELECT id, registration_country, registration_number, owner_id, icon FROM vehicles WHERE id = $1"
	vehicleList, err := p.fetch(ctx, sql, id)
	if err != nil {
		return vehicles.Vehicle{}, err
	}
	if len(vehicleList) == 0 {
		return vehicles.Vehicle{}, core.Errorf(core.ENOTFOUND, "vehicle with id %v not found", id)
	}
	return vehicleList[0], nil
}

// GetByOwnerId implements vehicles.VehicleRepository.
func (p *postgresVehicleRepository) GetByOwnerId(ctx context.Context, userId int64) ([]vehicles.Vehicle, error) {
	sql := "SELECT id, registration_country, registration_number, owner_id, icon FROM vehicles WHERE owner_id = $1"
	vehicleList, err := p.fetch(ctx, sql, userId)
	if err != nil {
		return vehicleList, err
	}
	return vehicleList, nil
}

// GetSimulatedVehicles implements vehicles.VehicleRepository.
func (p *postgresVehicleRepository) GetSimulatedVehicles(ctx context.Context) ([]vehicles.Vehicle, error) {
	sql := `SELECT v.id, v.registration_country, v.registration_number, v.owner_id, v.icon FROM vehicles v
			WHERE v.owner_id IN (SELECT id FROM users WHERE simulated = true)`
	vehicleList, err := p.fetch(ctx, sql)
	if err != nil {
		return vehicleList, err
	}
	return vehicleList, nil
}

// GetVehiclePositions implements vehicles.VehicleRepository.
func (p *postgresVehicleRepository) GetVehiclePositions(ctx context.Context, vehicleIds []int64) ([]vehicles.VehiclePosition, error) {
	positions := make([]vehicles.VehiclePosition, 0)
	if len(vehicleIds) == 0 {
		return positions, nil
	}
	vehicleIdsStrs := lo.Map(vehicleIds, func(x int64, index int) string {
		return strconv.FormatInt(x, 10)
	})
	vehicleIdsStr := strings.Join(vehicleIdsStrs, ",")
	sql := fmt.Sprintf("SELECT id, vehicle_id, lat, lng, recorded_at FROM vehicle_positions WHERE vehicle_id IN (%v)", vehicleIdsStr)
	rows, err := p.conn.Query(ctx, sql)
	if err != nil {
		return positions, err
	}
	defer rows.Close()
	for rows.Next() {
		var p vehicles.VehiclePosition
		if err := rows.Scan(
			&p.ID,
			&p.VehicleID,
			&p.Lat,
			&p.Lng,
			&p.RecordedAt,
		); err != nil {
			return positions, err
		}
		positions = append(positions, p)
	}
	return positions, nil

}

// UpdatePosition implements vehicles.VehicleRepository.
func (p *postgresVehicleRepository) UpdatePosition(ctx context.Context, vehicleId int64, lat float64, lng float64) error {
	countSql := "SELECT COUNT(*) FROM vehicle_positions WHERE vehicle_id = $1"

	var count int
	err := p.conn.QueryRow(ctx, countSql, vehicleId).Scan(&count)
	if err != nil {
		return err
	}
	var sql string
	if count > 0 {
		sql = `UPDATE vehicle_positions SET lat = $2, lng = $3, recorded_at = $4 WHERE vehicle_id = $1`
	} else {
		sql = `INSERT INTO vehicle_positions (vehicle_id, lat, lng, recorded_at) VALUES ($1, $2, $3, $4)`
	}
	_, err = p.conn.Exec(ctx, sql, vehicleId, lat, lng, time.Now().UTC())
	return err
}
