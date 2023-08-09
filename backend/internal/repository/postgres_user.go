package repository

import (
	"context"
	"fmt"

	"github.com/bjarke-xyz/uber-clone-backend/internal/domain"
)

type postgresUserRepository struct {
	conn Connection
}

func NewPostgresUser(conn Connection) domain.UserRepository {
	return &postgresUserRepository{conn: conn}
}

func (p *postgresUserRepository) fetch(ctx context.Context, query string, args ...interface{}) ([]domain.User, error) {
	rows, err := p.conn.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	uu := make([]domain.User, 0)
	for rows.Next() {
		var u domain.User
		if err := rows.Scan(
			&u.ID,
			&u.UserID,
			&u.Name,
			&u.Simulated,
		); err != nil {
			return nil, err
		}
		uu = append(uu, u)
	}
	return uu, nil
}

// CreateOrUpdate implements domain.UserRepository.
func (p *postgresUserRepository) CreateOrUpdate(ctx context.Context, user *domain.User) error {
	if err := user.Validate(); err != nil {
		return err
	}
	sql := `INSERT INTO users (user_uid, name, simulated) VALUES ($1, $2, false)
			ON CONFLICT(user_uid) DO NOTHING
			RETURNING id`
	return p.conn.QueryRow(ctx, sql, user.UserID, user.Name).Scan(&user.ID)
}

// Delete implements domain.UserRepository.
func (p *postgresUserRepository) Delete(ctx context.Context, id int64) error {
	sql := "DELETE FROM users WHERE id = $1"
	_, err := p.conn.Exec(ctx, sql, id)
	return err
}

// GetByID implements domain.UserRepository.
func (p *postgresUserRepository) GetByID(ctx context.Context, id int64) (domain.User, error) {
	sql := "SELECT id, user_uid, name, simulated FROM users WHERE id = $1"
	users, err := p.fetch(ctx, sql, id)
	if err != nil {
		return domain.User{}, err
	}
	if len(users) == 0 {
		return domain.User{}, domain.ErrNotFound
	}
	return users[0], nil
}

// GetByUserID implements domain.UserRepository.
func (p *postgresUserRepository) GetByUserID(ctx context.Context, userID string) (domain.User, error) {
	sql := "SELECT id, user_uid, name, simulated FROM users WHERE user_uid = $1"
	users, err := p.fetch(ctx, sql, userID)
	if err != nil {
		return domain.User{}, err
	}
	if len(users) == 0 {
		return domain.User{}, domain.ErrNotFound
	}
	return users[0], nil
}

// GetSimulatedUsers implements domain.UserRepository.
func (p *postgresUserRepository) GetSimulatedUsers(ctx context.Context) ([]domain.User, error) {
	sql := "SELECT id, user_uid, name, simulated FROM users WHERE simulated = true"
	users, err := p.fetch(ctx, sql)
	if err != nil {
		return users, fmt.Errorf("failed to get simulated users: %w", err)
	}
	return users, nil

}
