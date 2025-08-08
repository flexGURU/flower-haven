package repository

import (
	"context"
	"time"

	"github.com/flexGURU/flower-haven/backend/pkg"
)

type User struct {
	ID           uint32    `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	Address      *string   `json:"address,omitempty"`
	PhoneNumber  string    `json:"phone_number"`
	RefreshToken *string   `json:"refresh_token,omitempty"`
	Password     *string   `json:"password,omitempty"`
	IsAdmin      bool      `json:"is_admin"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
}

type UpdateUser struct {
	ID           uint32  `json:"id"`
	Name         *string `json:"name"`
	Address      *string `json:"address"`
	PhoneNumber  *string `json:"phone_number"`
	Password     *string `json:"password"`
	IsAdmin      *bool   `json:"is_admin"`
	IsActive     *bool   `json:"is_active"`
	RefreshToken *string `json:"refresh_token"`
}

type UserFilter struct {
	Pagination *pkg.Pagination
	Search     *string
	IsAdmin    *bool
	IsActive   *bool
}

type UserRepository interface {
	CreateUser(ctx context.Context, user *User) (*User, error)
	GetUserByID(ctx context.Context, id int64) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	UpdateUser(ctx context.Context, user *UpdateUser) (*User, error)
	ListUsers(ctx context.Context, filter *UserFilter) ([]*User, *pkg.Pagination, error)
}
