package model

import (
	"context"
	"database/sql"
)

type User struct {
	ID       int64   `db:"id"`
	Email    string  `db:"email"`
	Password string  `db:"password"`
	Name     string  `db:"name"`
	Phone    *string `db:"phone"`
	Role     string  `db:"role"`
	Address  *string `db:"address"`
}

// UserRequest represents the request for creating a new user
type UserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Name     string `json:"name" validate:"required"`
	Phone    string `json:"phone"`
	Role     string `json:"role" validate:"required,oneof=admin user manager"`
	Address  string `json:"address"`
}

// UserResponse represents the response for user operations
type UserResponse struct {
	ID        int64        `json:"id"`
	Email     string       `json:"email"`
	Name      string       `json:"name"`
	Phone     *string      `json:"phone"`
	Role      string       `json:"role"`
	Address   *string      `json:"address"`
	CreatedAt sql.NullTime `json:"created_at"`
	UpdatedAt sql.NullTime `json:"updated_at"`
}

// LoginRequest represents the request for user login
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// LoginResponse represents the response for user login
type LoginResponse struct {
	Token string       `json:"token"`
}

// UserSQLStore interface for database operations
type UserSQLStore interface {
	Get(ctx context.Context, userID int64) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, userID int64) error
	List(ctx context.Context, limit, offset int) ([]*User, error)
	Count(ctx context.Context) (int64, error)
}
