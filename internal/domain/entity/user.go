package entity

import (
	"time"
)

// UserRole represents user role type
type UserRole string

const (
	RoleAdmin UserRole = "admin"
	RoleUser  UserRole = "user"
)

// IsValid checks if the role is valid
func (r UserRole) IsValid() bool {
	switch r {
	case RoleAdmin, RoleUser:
		return true
	default:
		return false
	}
}

// User represents a user in the system
type User struct {
	ID        int64     `json:"id" db:"id"`
	Email     string    `json:"email" db:"email"`
	Name      string    `json:"name" db:"name"`
	Phone     string    `json:"phone" db:"phone"`
	Password  string    `json:"-" db:"password"` // Never include in JSON responses
	Role      UserRole  `json:"role" db:"role"`
	IsActive  bool      `json:"is_active" db:"is_active"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// IsAdmin returns true if user has admin role
func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

// IsUser returns true if user has user role
func (u *User) IsUser() bool {
	return u.Role == RoleUser
}

// ToPublicUser returns user data safe for public consumption
func (u *User) ToPublicUser() *PublicUser {
	return &PublicUser{
		ID:        u.ID,
		Email:     u.Email,
		Name:      u.Name,
		Phone:     u.Phone,
		Role:      u.Role,
		IsActive:  u.IsActive,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// PublicUser represents user data safe for public consumption (without sensitive data)
type PublicUser struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Phone     string    `json:"phone"`
	Role      UserRole  `json:"role"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateUserRequest represents request to create a new user
type CreateUserRequest struct {
	Email    string   `json:"email" binding:"required,email"`
	Name     string   `json:"name" binding:"required,min=2,max=100"`
	Phone    string   `json:"phone" binding:"required,min=10,max=15"`
	Password string   `json:"password" binding:"required,min=6"`
	Role     UserRole `json:"role,omitempty"`
}

// LoginRequest represents login request
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// AuthResponse represents authentication response
type AuthResponse struct {
	User  *PublicUser `json:"user"`
	Token string      `json:"token"`
}

// Session represents user session data stored in JWT
type Session struct {
	UserID int64    `json:"user_id"`
	Email  string   `json:"email"`
	Role   UserRole `json:"role"`
}
