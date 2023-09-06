package repository

import (
	"context"

	"github.com/bjarke-xyz/uber-clone-backend/internal/core"
	"github.com/bjarke-xyz/uber-clone-backend/internal/core/users"
)

type postgresUserRepository struct {
	conn Connection
}

func NewPostgresUser(conn Connection) users.UserRepository {
	return &postgresUserRepository{conn: conn}
}

func (p *postgresUserRepository) fetch(ctx context.Context, query string, args ...interface{}) ([]users.User, error) {
	rows, err := p.conn.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	uu := make([]users.User, 0)
	for rows.Next() {
		var u users.User
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

// CreateOrUpdate implements users.UserRepository.
func (p *postgresUserRepository) CreateOrUpdate(ctx context.Context, user *users.User) error {
	if err := user.Validate(); err != nil {
		return err
	}
	sql := `INSERT INTO users (user_uid, name, simulated) VALUES ($1, $2, false)
			ON CONFLICT(user_uid) DO NOTHING
			RETURNING id`
	return p.conn.QueryRow(ctx, sql, user.UserID, user.Name).Scan(&user.ID)
}

// Delete implements users.UserRepository.
func (p *postgresUserRepository) Delete(ctx context.Context, id int64) error {
	sql := "DELETE FROM users WHERE id = $1"
	_, err := p.conn.Exec(ctx, sql, id)
	return err
}

// GetByID implements users.UserRepository.
func (p *postgresUserRepository) GetByID(ctx context.Context, id int64) (users.User, error) {
	sql := "SELECT id, user_uid, name, simulated FROM users WHERE id = $1"
	userList, err := p.fetch(ctx, sql, id)
	if err != nil {
		return users.User{}, err
	}
	if len(userList) == 0 {
		return users.User{}, core.Errorf(core.ENOTFOUND, "user with id %v not found", id)
	}
	return userList[0], nil
}

// GetByUserID implements users.UserRepository.
func (p *postgresUserRepository) GetByUserID(ctx context.Context, userID string) (users.User, error) {
	sql := "SELECT id, user_uid, name, simulated FROM users WHERE user_uid = $1"
	userList, err := p.fetch(ctx, sql, userID)
	if err != nil {
		return users.User{}, err
	}
	if len(userList) == 0 {
		return users.User{}, core.Errorf(core.ENOTFOUND, "user with id %v not found", userID)
	}
	return userList[0], nil
}

// GetSimulatedUsers implements users.UserRepository.
func (p *postgresUserRepository) GetSimulatedUsers(ctx context.Context) ([]users.User, error) {
	sql := "SELECT id, user_uid, name, simulated FROM users WHERE simulated = true"
	userList, err := p.fetch(ctx, sql)
	if err != nil {
		return userList, err
	}
	return userList, nil

}
