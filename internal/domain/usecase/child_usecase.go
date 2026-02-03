package usecase

import (
	"context"

	"github.com/afghan/rumi-backend/internal/domain/entity"
)

// ChildUsecase defines the interface for child management business logic
type ChildUsecase interface {
	// CreateChild creates a new child for a user
	CreateChild(ctx context.Context, userID int64, req *entity.CreateChildRequest) (*entity.Child, error)

	// GetChildByID retrieves a child by ID (with ownership check for non-admin users)
	GetChildByID(ctx context.Context, childID int64, userID int64, userRole entity.UserRole) (*entity.Child, error)

	// GetUserChildren retrieves children for a specific user
	GetUserChildren(ctx context.Context, userID int64, limit, offset int) ([]*entity.Child, int64, error)

	// GetAllChildren retrieves all children (admin only)
	GetAllChildren(ctx context.Context, limit, offset int) ([]*entity.Child, int64, error)

	// UpdateChild updates a child (with ownership check for non-admin users)
	UpdateChild(ctx context.Context, childID int64, userID int64, userRole entity.UserRole, req *entity.UpdateChildRequest) (*entity.Child, error)

	// DeleteChild soft deletes a child (with ownership check for non-admin users)
	DeleteChild(ctx context.Context, childID int64, userID int64, userRole entity.UserRole) error

	// SetChildActive sets child active status (admin only)
	SetChildActive(ctx context.Context, childID int64, isActive bool) error
}
