package repository

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/bjarke-xyz/uber-clone-backend/internal/domain"
	"github.com/samber/lo"
)

type postgresVehicleRepository struct {
	conn Connection
}

func NewPostgresVehicle(conn Connection) domain.VehicleRepository {
	return &postgresVehicleRepository{conn: conn}
}

func (p *postgresVehicleRepository) fetch(ctx context.Context, query string, args ...interface{}) ([]domain.Vehicle, error) {
	rows, err := p.conn.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	vv := make([]domain.Vehicle, 0)
	for rows.Next() {
		var v domain.Vehicle
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

// GetByIdAndOwnerId implements domain.VehicleRepository.
func (p *postgresVehicleRepository) GetByIdAndOwnerId(ctx context.Context, vehicleId int64, userId int64) (domain.Vehicle, error) {
	sql := "SELECT id, registration_country, registration_number, owner_id, icon FROM vehicles WHERE id = $1 AND owner_id = $2"
	vehicles, err := p.fetch(ctx, sql, vehicleId, userId)
	if err != nil {
		return domain.Vehicle{}, err
	}
	if len(vehicles) == 0 {
		return domain.Vehicle{}, domain.ErrNotFound
	}
	return vehicles[0], nil
}

// CreateOrUpdate implements domain.VehicleRepository.
func (p *postgresVehicleRepository) CreateOrUpdate(ctx context.Context, v *domain.Vehicle) error {
	panic("unimplemented")
}

// Delete implements domain.VehicleRepository.
func (p *postgresVehicleRepository) Delete(ctx context.Context, id int64) error {
	panic("unimplemented")
}

// GetByID implements domain.VehicleRepository.
func (p *postgresVehicleRepository) GetByID(ctx context.Context, id int64) (domain.Vehicle, error) {
	sql := "SELECT id, registration_country, registration_number, owner_id, icon FROM vehicles WHERE id = $1"
	vehicles, err := p.fetch(ctx, sql, id)
	if err != nil {
		return domain.Vehicle{}, err
	}
	if len(vehicles) == 0 {
		return domain.Vehicle{}, domain.ErrNotFound
	}
	return vehicles[0], nil
}

// GetByOwnerId implements domain.VehicleRepository.
func (p *postgresVehicleRepository) GetByOwnerId(ctx context.Context, userId int64) ([]domain.Vehicle, error) {
	sql := "SELECT id, registration_country, registration_number, owner_id, icon FROM vehicles WHERE owner_id = $1"
	vehicles, err := p.fetch(ctx, sql, userId)
	if err != nil {
		return vehicles, err
	}
	return vehicles, nil
}

// GetSimulatedVehicles implements domain.VehicleRepository.
func (p *postgresVehicleRepository) GetSimulatedVehicles(ctx context.Context) ([]domain.Vehicle, error) {
	sql := `SELECT v.id, v.registration_country, v.registration_number, v.owner_id, v.icon FROM vehicles v
			WHERE v.owner_id IN (SELECT id FROM users WHERE simulated = true)`
	vehicles, err := p.fetch(ctx, sql)
	if err != nil {
		return vehicles, err
	}
	return vehicles, nil
}

// GetVehiclePositions implements domain.VehicleRepository.
func (p *postgresVehicleRepository) GetVehiclePositions(ctx context.Context, vehicleIds []int64) ([]domain.VehiclePosition, error) {
	positions := make([]domain.VehiclePosition, 0)
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
		var p domain.VehiclePosition
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

// UpdatePosition implements domain.VehicleRepository.
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
