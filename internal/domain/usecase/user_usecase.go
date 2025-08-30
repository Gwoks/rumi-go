package usecase

import (
	"context"

	"github.com/afghan/rumi-backend/internal/domain/entity"
)

// UserManagementUsecase defines the interface for user management business logic
type UserManagementUsecase interface {
	// GetAllUsers retrieves all users with pagination (admin only)
	GetAllUsers(ctx context.Context, limit, offset int) ([]*entity.PublicUser, int64, error)

	// GetUserByID retrieves user by ID (admin only)
	GetUserByID(ctx context.Context, id int64) (*entity.PublicUser, error)

	// UpdateUser updates user information (admin only)
	UpdateUser(ctx context.Context, userID int64, req *entity.CreateUserRequest) (*entity.PublicUser, error)

	// DeleteUser soft deletes a user (admin only)
	DeleteUser(ctx context.Context, userID int64) error

	// SetUserActive sets user active status (admin only)
	SetUserActive(ctx context.Context, userID int64, isActive bool) error

	// UpdateProfile updates user's own profile (name, phone only)
	UpdateProfile(ctx context.Context, userID int64, req *entity.UpdateProfileRequest) (*entity.PublicUser, error)

	// ChangePassword changes user's password
	ChangePassword(ctx context.Context, userID int64, req *entity.ChangePasswordRequest) error
}

// UserActivityUsecase defines the interface for user activity business logic
type UserActivityUsecase interface {
	// LogActivity logs a user activity
	LogActivity(ctx context.Context, req *entity.CreateActivityRequest) error

	// GetUserActivities retrieves activities for a specific user
	GetUserActivities(ctx context.Context, userID int64, limit, offset int) ([]*entity.UserActivity, int64, error)

	// GetAllActivities retrieves all activities (admin only)
	GetAllActivities(ctx context.Context, limit, offset int) ([]*entity.UserActivity, int64, error)

	// GetRecentActivities gets recent activities
	GetRecentActivities(ctx context.Context, limit int) ([]*entity.UserActivity, error)

	// GetActivitiesByType retrieves activities by type
	GetActivitiesByType(ctx context.Context, activityType entity.ActivityType, limit, offset int) ([]*entity.UserActivity, int64, error)
}
