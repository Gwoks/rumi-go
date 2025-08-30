package usecase

import (
	"context"
	"fmt"

	"github.com/afghan/rumi-backend/internal/domain/entity"
	"github.com/afghan/rumi-backend/internal/domain/repository"
	"github.com/afghan/rumi-backend/internal/domain/usecase"
)

// userActivityUsecaseImpl implements UserActivityUsecase interface
type userActivityUsecaseImpl struct {
	activityRepo repository.UserActivityRepository
	userRepo     repository.UserRepository
}

// NewUserActivityUsecase creates a new user activity usecase
func NewUserActivityUsecase(
	activityRepo repository.UserActivityRepository,
	userRepo repository.UserRepository,
) usecase.UserActivityUsecase {
	return &userActivityUsecaseImpl{
		activityRepo: activityRepo,
		userRepo:     userRepo,
	}
}

// LogActivity logs a user activity
func (u *userActivityUsecaseImpl) LogActivity(ctx context.Context, req *entity.CreateActivityRequest) error {
	// Validate request
	if req.UserID <= 0 {
		return fmt.Errorf("user ID is required")
	}
	if req.ActivityType == "" {
		return fmt.Errorf("activity type is required")
	}
	if req.Description == "" {
		return fmt.Errorf("description is required")
	}

	// Verify user exists
	_, err := u.userRepo.GetByID(ctx, req.UserID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Create activity
	activity := &entity.UserActivity{
		UserID:       req.UserID,
		ActivityType: req.ActivityType,
		Description:  req.Description,
		IPAddress:    req.IPAddress,
		UserAgent:    req.UserAgent,
		Metadata:     req.Metadata,
	}

	_, err = u.activityRepo.Create(ctx, activity)
	if err != nil {
		return fmt.Errorf("failed to log activity: %w", err)
	}

	return nil
}

// GetUserActivities retrieves activities for a specific user
func (u *userActivityUsecaseImpl) GetUserActivities(ctx context.Context, userID int64, limit, offset int) ([]*entity.UserActivity, int64, error) {
	// Validate pagination
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	// Verify user exists
	_, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, 0, fmt.Errorf("user not found: %w", err)
	}

	activities, err := u.activityRepo.GetByUserID(ctx, userID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get user activities: %w", err)
	}

	total, err := u.activityRepo.CountByUserID(ctx, userID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count user activities: %w", err)
	}

	return activities, total, nil
}

// GetAllActivities retrieves all activities (admin only)
func (u *userActivityUsecaseImpl) GetAllActivities(ctx context.Context, limit, offset int) ([]*entity.UserActivity, int64, error) {
	// Validate pagination
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	activities, err := u.activityRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get all activities: %w", err)
	}

	total, err := u.activityRepo.Count(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count all activities: %w", err)
	}

	return activities, total, nil
}

// GetRecentActivities gets recent activities
func (u *userActivityUsecaseImpl) GetRecentActivities(ctx context.Context, limit int) ([]*entity.UserActivity, error) {
	if limit <= 0 || limit > 50 {
		limit = 10
	}

	activities, err := u.activityRepo.GetRecentActivities(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent activities: %w", err)
	}

	return activities, nil
}

// GetActivitiesByType retrieves activities by type
func (u *userActivityUsecaseImpl) GetActivitiesByType(ctx context.Context, activityType entity.ActivityType, limit, offset int) ([]*entity.UserActivity, int64, error) {
	// Validate pagination
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	activities, err := u.activityRepo.ListByType(ctx, activityType, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get activities by type: %w", err)
	}

	total, err := u.activityRepo.CountByType(ctx, activityType)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count activities by type: %w", err)
	}

	return activities, total, nil
}
