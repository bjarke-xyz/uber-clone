package domain

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation"
)

type User struct {
	ID int64 `json:"id"`

	Name      string `json:"name"`
	Simulated bool   `json:"simulated"`
	// Firebase auth info
	UserID string `json:"userId"`
}

func (u *User) Validate() error {
	return validation.ValidateStruct(u,
		validation.Field(&u.Name, validation.Required, validation.Length(2, 200)),
	)
}

type UserRepository interface {
	GetByID(context.Context, int64) (User, error)
	GetByUserID(context.Context, string) (User, error)
	GetSimulatedUsers(context.Context) ([]User, error)

	CreateOrUpdate(context.Context, *User) error
	Delete(context.Context, int64) error
}
