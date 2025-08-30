package repository

import (
	"context"

	"github.com/afghan/rumi-backend/internal/domain/entity"
)

// UserActivityRepository defines the interface for user activity data operations
type UserActivityRepository interface {
	// Create creates a new user activity record
	Create(ctx context.Context, activity *entity.UserActivity) (*entity.UserActivity, error)

	// GetByID retrieves an activity by ID
	GetByID(ctx context.Context, id int64) (*entity.UserActivity, error)

	// GetByUserID retrieves activities by user ID with pagination
	GetByUserID(ctx context.Context, userID int64, limit, offset int) ([]*entity.UserActivity, error)

	// List retrieves activities with pagination and optional filtering
	List(ctx context.Context, limit, offset int) ([]*entity.UserActivity, error)

	// ListByType retrieves activities by type with pagination
	ListByType(ctx context.Context, activityType entity.ActivityType, limit, offset int) ([]*entity.UserActivity, error)

	// Count returns total number of activities
	Count(ctx context.Context) (int64, error)

	// CountByUserID returns total number of activities for a user
	CountByUserID(ctx context.Context, userID int64) (int64, error)

	// CountByType returns total number of activities by type
	CountByType(ctx context.Context, activityType entity.ActivityType) (int64, error)

	// Delete deletes an activity by ID (soft delete not needed for logs)
	Delete(ctx context.Context, id int64) error

	// GetRecentActivities gets recent activities across all users
	GetRecentActivities(ctx context.Context, limit int) ([]*entity.UserActivity, error)
}
