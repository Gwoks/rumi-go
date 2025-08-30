package usecase

import (
	"context"

	"github.com/afghan/rumi-backend/internal/domain/entity"
)

// AuthUsecase defines the interface for authentication business logic
type AuthUsecase interface {
	// Login authenticates user and returns JWT token
	Login(ctx context.Context, req *entity.LoginRequest) (*entity.AuthResponse, error)

	// Signup creates a new user account and returns JWT token
	Signup(ctx context.Context, req *entity.CreateUserRequest) (*entity.AuthResponse, error)

	// Logout invalidates user session (could be used for token blacklisting)
	Logout(ctx context.Context, token string) error

	// RefreshToken generates a new JWT token from existing valid token
	RefreshToken(ctx context.Context, token string) (*entity.AuthResponse, error)

	// ValidateToken validates JWT token and returns user session data
	ValidateToken(ctx context.Context, token string) (*entity.Session, error)

	// GetProfile retrieves user profile by ID
	GetProfile(ctx context.Context, userID int64) (*entity.PublicUser, error)
}

// UserUsecase defines the interface for user management business logic
type UserUsecase interface {
	// Create creates a new user (admin only)
	Create(ctx context.Context, req *entity.CreateUserRequest) (*entity.PublicUser, error)

	// GetByID retrieves user by ID
	GetByID(ctx context.Context, id int64) (*entity.PublicUser, error)

	// Update updates user profile
	Update(ctx context.Context, user *entity.User) (*entity.PublicUser, error)

	// Delete soft deletes a user
	Delete(ctx context.Context, id int64) error

	// List retrieves users with pagination
	List(ctx context.Context, limit, offset int) ([]*entity.PublicUser, error)

	// SetActive sets user active status
	SetActive(ctx context.Context, userID int64, isActive bool) error
}
