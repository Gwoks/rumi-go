package repository

import (
	"context"

	"github.com/afghan/rumi-backend/internal/domain/entity"
)

// UserRepository defines the interface for user data operations
type UserRepository interface {
	// Create creates a new user
	Create(ctx context.Context, user *entity.User) (*entity.User, error)

	// GetByID retrieves a user by ID
	GetByID(ctx context.Context, id int64) (*entity.User, error)

	// GetByEmail retrieves a user by email
	GetByEmail(ctx context.Context, email string) (*entity.User, error)

	// Update updates an existing user
	Update(ctx context.Context, user *entity.User) (*entity.User, error)

	// Delete soft deletes a user by ID
	Delete(ctx context.Context, id int64) error

	// List retrieves users with pagination
	List(ctx context.Context, limit, offset int) ([]*entity.User, error)

	// Count returns total number of users
	Count(ctx context.Context) (int64, error)

	// GetByRole retrieves users by role
	GetByRole(ctx context.Context, role entity.UserRole, limit, offset int) ([]*entity.User, error)

	// UpdateLastLogin updates user's last login time
	UpdateLastLogin(ctx context.Context, userID int64) error

	// SetActive sets user active status
	SetActive(ctx context.Context, userID int64, isActive bool) error

	// EmailExists checks if email already exists
	EmailExists(ctx context.Context, email string) (bool, error)
}
