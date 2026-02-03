package repository

import (
	"context"

	"github.com/afghan/rumi-backend/internal/domain/entity"
)

// ChildRepository defines the interface for child data operations
type ChildRepository interface {
	// Create creates a new child
	Create(ctx context.Context, child *entity.Child) (*entity.Child, error)

	// GetByID retrieves a child by ID
	GetByID(ctx context.Context, id int64) (*entity.Child, error)

	// GetByUserID retrieves children by user ID with pagination
	GetByUserID(ctx context.Context, userID int64, limit, offset int) ([]*entity.Child, error)

	// Update updates an existing child
	Update(ctx context.Context, child *entity.Child) (*entity.Child, error)

	// Delete soft deletes a child by ID
	Delete(ctx context.Context, id int64) error

	// List retrieves children with pagination (admin only)
	List(ctx context.Context, limit, offset int) ([]*entity.Child, error)

	// Count returns total number of children
	Count(ctx context.Context) (int64, error)

	// CountByUserID returns total number of children for a user
	CountByUserID(ctx context.Context, userID int64) (int64, error)

	// SetActive sets child active status
	SetActive(ctx context.Context, childID int64, isActive bool) error

	// GetByIDAndUserID retrieves a child by ID and user ID (for ownership check)
	GetByIDAndUserID(ctx context.Context, id, userID int64) (*entity.Child, error)
}
