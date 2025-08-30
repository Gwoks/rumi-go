package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/afghan/rumi-backend/internal/database"
	"github.com/afghan/rumi-backend/internal/domain/entity"
	"github.com/afghan/rumi-backend/internal/domain/repository"
)

// userActivityRepositoryImpl implements UserActivityRepository interface
type userActivityRepositoryImpl struct {
	db *database.DB
}

// NewUserActivityRepository creates a new user activity repository
func NewUserActivityRepository(db *database.DB) repository.UserActivityRepository {
	return &userActivityRepositoryImpl{
		db: db,
	}
}

// Create creates a new user activity record
func (r *userActivityRepositoryImpl) Create(ctx context.Context, activity *entity.UserActivity) (*entity.UserActivity, error) {
	query := `
		INSERT INTO user_activities (user_id, activity_type, description, ip_address, user_agent, metadata, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)`

	now := time.Now()
	activity.CreatedAt = now

	// Convert metadata to JSON string
	var metadataJSON []byte
	var err error
	if activity.Metadata != nil {
		metadataJSON, err = activity.MarshalMetadata()
		if err != nil {
			return nil, fmt.Errorf("failed to marshal metadata: %w", err)
		}
	}

	result, err := r.db.ExecContext(ctx, query,
		activity.UserID,
		activity.ActivityType,
		activity.Description,
		activity.IPAddress,
		activity.UserAgent,
		metadataJSON,
		activity.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create user activity: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}

	activity.ID = id
	return activity, nil
}

// GetByID retrieves an activity by ID
func (r *userActivityRepositoryImpl) GetByID(ctx context.Context, id int64) (*entity.UserActivity, error) {
	query := `
		SELECT ua.id, ua.user_id, ua.activity_type, ua.description, ua.ip_address, ua.user_agent, ua.metadata, ua.created_at,
			   u.email, u.name, u.phone, u.role, u.is_active, u.created_at as user_created_at, u.updated_at as user_updated_at
		FROM user_activities ua
		LEFT JOIN users u ON ua.user_id = u.id
		WHERE ua.id = ?`

	var activity entity.UserActivity
	var user entity.PublicUser
	var metadataJSON []byte

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&activity.ID,
		&activity.UserID,
		&activity.ActivityType,
		&activity.Description,
		&activity.IPAddress,
		&activity.UserAgent,
		&metadataJSON,
		&activity.CreatedAt,
		&user.Email,
		&user.Name,
		&user.Phone,
		&user.Role,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("activity not found")
		}
		return nil, fmt.Errorf("failed to get activity by id: %w", err)
	}

	// Unmarshal metadata
	if len(metadataJSON) > 0 {
		if err := activity.UnmarshalMetadata(metadataJSON); err != nil {
			return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
		}
	}

	user.ID = activity.UserID
	activity.User = &user

	return &activity, nil
}

// GetByUserID retrieves activities by user ID with pagination
func (r *userActivityRepositoryImpl) GetByUserID(ctx context.Context, userID int64, limit, offset int) ([]*entity.UserActivity, error) {
	query := `
		SELECT id, user_id, activity_type, description, ip_address, user_agent, metadata, created_at
		FROM user_activities
		WHERE user_id = ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?`

	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get activities by user id: %w", err)
	}
	defer rows.Close()

	var activities []*entity.UserActivity
	for rows.Next() {
		var activity entity.UserActivity
		var metadataJSON []byte

		err := rows.Scan(
			&activity.ID,
			&activity.UserID,
			&activity.ActivityType,
			&activity.Description,
			&activity.IPAddress,
			&activity.UserAgent,
			&metadataJSON,
			&activity.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan activity: %w", err)
		}

		// Unmarshal metadata
		if len(metadataJSON) > 0 {
			if err := activity.UnmarshalMetadata(metadataJSON); err != nil {
				return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
			}
		}

		activities = append(activities, &activity)
	}

	return activities, nil
}

// List retrieves activities with pagination and joins user info
func (r *userActivityRepositoryImpl) List(ctx context.Context, limit, offset int) ([]*entity.UserActivity, error) {
	query := `
		SELECT ua.id, ua.user_id, ua.activity_type, ua.description, ua.ip_address, ua.user_agent, ua.metadata, ua.created_at,
			   u.email, u.name, u.phone, u.role, u.is_active, u.created_at as user_created_at, u.updated_at as user_updated_at
		FROM user_activities ua
		LEFT JOIN users u ON ua.user_id = u.id
		ORDER BY ua.created_at DESC
		LIMIT ? OFFSET ?`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list activities: %w", err)
	}
	defer rows.Close()

	var activities []*entity.UserActivity
	for rows.Next() {
		var activity entity.UserActivity
		var user entity.PublicUser
		var metadataJSON []byte

		err := rows.Scan(
			&activity.ID,
			&activity.UserID,
			&activity.ActivityType,
			&activity.Description,
			&activity.IPAddress,
			&activity.UserAgent,
			&metadataJSON,
			&activity.CreatedAt,
			&user.Email,
			&user.Name,
			&user.Phone,
			&user.Role,
			&user.IsActive,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan activity: %w", err)
		}

		// Unmarshal metadata
		if len(metadataJSON) > 0 {
			if err := activity.UnmarshalMetadata(metadataJSON); err != nil {
				return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
			}
		}

		user.ID = activity.UserID
		activity.User = &user
		activities = append(activities, &activity)
	}

	return activities, nil
}

// ListByType retrieves activities by type with pagination
func (r *userActivityRepositoryImpl) ListByType(ctx context.Context, activityType entity.ActivityType, limit, offset int) ([]*entity.UserActivity, error) {
	query := `
		SELECT ua.id, ua.user_id, ua.activity_type, ua.description, ua.ip_address, ua.user_agent, ua.metadata, ua.created_at,
			   u.email, u.name, u.phone, u.role, u.is_active, u.created_at as user_created_at, u.updated_at as user_updated_at
		FROM user_activities ua
		LEFT JOIN users u ON ua.user_id = u.id
		WHERE ua.activity_type = ?
		ORDER BY ua.created_at DESC
		LIMIT ? OFFSET ?`

	rows, err := r.db.QueryContext(ctx, query, activityType, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list activities by type: %w", err)
	}
	defer rows.Close()

	var activities []*entity.UserActivity
	for rows.Next() {
		var activity entity.UserActivity
		var user entity.PublicUser
		var metadataJSON []byte

		err := rows.Scan(
			&activity.ID,
			&activity.UserID,
			&activity.ActivityType,
			&activity.Description,
			&activity.IPAddress,
			&activity.UserAgent,
			&metadataJSON,
			&activity.CreatedAt,
			&user.Email,
			&user.Name,
			&user.Phone,
			&user.Role,
			&user.IsActive,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan activity: %w", err)
		}

		// Unmarshal metadata
		if len(metadataJSON) > 0 {
			if err := activity.UnmarshalMetadata(metadataJSON); err != nil {
				return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
			}
		}

		user.ID = activity.UserID
		activity.User = &user
		activities = append(activities, &activity)
	}

	return activities, nil
}

// Count returns total number of activities
func (r *userActivityRepositoryImpl) Count(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM user_activities`

	var count int64
	err := r.db.GetContext(ctx, &count, query)
	if err != nil {
		return 0, fmt.Errorf("failed to count activities: %w", err)
	}

	return count, nil
}

// CountByUserID returns total number of activities for a user
func (r *userActivityRepositoryImpl) CountByUserID(ctx context.Context, userID int64) (int64, error) {
	query := `SELECT COUNT(*) FROM user_activities WHERE user_id = ?`

	var count int64
	err := r.db.GetContext(ctx, &count, query, userID)
	if err != nil {
		return 0, fmt.Errorf("failed to count activities by user id: %w", err)
	}

	return count, nil
}

// CountByType returns total number of activities by type
func (r *userActivityRepositoryImpl) CountByType(ctx context.Context, activityType entity.ActivityType) (int64, error) {
	query := `SELECT COUNT(*) FROM user_activities WHERE activity_type = ?`

	var count int64
	err := r.db.GetContext(ctx, &count, query, activityType)
	if err != nil {
		return 0, fmt.Errorf("failed to count activities by type: %w", err)
	}

	return count, nil
}

// Delete deletes an activity by ID
func (r *userActivityRepositoryImpl) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM user_activities WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete activity: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("activity not found")
	}

	return nil
}

// GetRecentActivities gets recent activities across all users
func (r *userActivityRepositoryImpl) GetRecentActivities(ctx context.Context, limit int) ([]*entity.UserActivity, error) {
	query := `
		SELECT ua.id, ua.user_id, ua.activity_type, ua.description, ua.ip_address, ua.user_agent, ua.metadata, ua.created_at,
			   u.email, u.name, u.phone, u.role, u.is_active, u.created_at as user_created_at, u.updated_at as user_updated_at
		FROM user_activities ua
		LEFT JOIN users u ON ua.user_id = u.id
		ORDER BY ua.created_at DESC
		LIMIT ?`

	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent activities: %w", err)
	}
	defer rows.Close()

	var activities []*entity.UserActivity
	for rows.Next() {
		var activity entity.UserActivity
		var user entity.PublicUser
		var metadataJSON []byte

		err := rows.Scan(
			&activity.ID,
			&activity.UserID,
			&activity.ActivityType,
			&activity.Description,
			&activity.IPAddress,
			&activity.UserAgent,
			&metadataJSON,
			&activity.CreatedAt,
			&user.Email,
			&user.Name,
			&user.Phone,
			&user.Role,
			&user.IsActive,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan activity: %w", err)
		}

		// Unmarshal metadata
		if len(metadataJSON) > 0 {
			if err := activity.UnmarshalMetadata(metadataJSON); err != nil {
				return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
			}
		}

		user.ID = activity.UserID
		activity.User = &user
		activities = append(activities, &activity)
	}

	return activities, nil
}
