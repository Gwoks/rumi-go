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

// userRepositoryImpl implements UserRepository interface using sqlx with raw queries
type userRepositoryImpl struct {
	db *database.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *database.DB) repository.UserRepository {
	return &userRepositoryImpl{
		db: db,
	}
}

// Create creates a new user
func (r *userRepositoryImpl) Create(ctx context.Context, user *entity.User) (*entity.User, error) {
	query := `
		INSERT INTO users (email, name, phone, password, role, is_active, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	result, err := r.db.ExecContext(ctx, query,
		user.Email,
		user.Name,
		user.Phone,
		user.Password,
		user.Role,
		user.IsActive,
		user.CreatedAt,
		user.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}

	user.ID = id
	return user, nil
}

// GetByID retrieves a user by ID
func (r *userRepositoryImpl) GetByID(ctx context.Context, id int64) (*entity.User, error) {
	query := `
		SELECT id, email, name, phone, password, role, is_active, created_at, updated_at
		FROM users 
		WHERE id = ? AND is_active = TRUE`

	var user entity.User
	err := r.db.GetContext(ctx, &user, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	return &user, nil
}

// GetByEmail retrieves a user by email
func (r *userRepositoryImpl) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	query := `
		SELECT id, email, name, phone, password, role, is_active, created_at, updated_at
		FROM users 
		WHERE email = ?`

	var user entity.User
	err := r.db.GetContext(ctx, &user, query, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return &user, nil
}

// Update updates an existing user
func (r *userRepositoryImpl) Update(ctx context.Context, user *entity.User) (*entity.User, error) {
	query := `
		UPDATE users 
		SET name = ?, phone = ?, role = ?, is_active = ?, updated_at = ?
		WHERE id = ?`

	user.UpdatedAt = time.Now()

	result, err := r.db.ExecContext(ctx, query,
		user.Name,
		user.Phone,
		user.Role,
		user.IsActive,
		user.UpdatedAt,
		user.ID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return nil, fmt.Errorf("user not found or no changes made")
	}

	return user, nil
}

// Delete soft deletes a user by ID
func (r *userRepositoryImpl) Delete(ctx context.Context, id int64) error {
	query := `UPDATE users SET is_active = FALSE, updated_at = ? WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// List retrieves users with pagination
func (r *userRepositoryImpl) List(ctx context.Context, limit, offset int) ([]*entity.User, error) {
	query := `
		SELECT id, email, name, phone, password, role, is_active, created_at, updated_at
		FROM users 
		WHERE is_active = TRUE
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?`

	var users []*entity.User
	err := r.db.SelectContext(ctx, &users, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	return users, nil
}

// Count returns total number of users
func (r *userRepositoryImpl) Count(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM users WHERE is_active = TRUE`

	var count int64
	err := r.db.GetContext(ctx, &count, query)
	if err != nil {
		return 0, fmt.Errorf("failed to count users: %w", err)
	}

	return count, nil
}

// GetByRole retrieves users by role
func (r *userRepositoryImpl) GetByRole(ctx context.Context, role entity.UserRole, limit, offset int) ([]*entity.User, error) {
	query := `
		SELECT id, email, name, phone, password, role, is_active, created_at, updated_at
		FROM users 
		WHERE role = ? AND is_active = TRUE
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?`

	var users []*entity.User
	err := r.db.SelectContext(ctx, &users, query, role, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get users by role: %w", err)
	}

	return users, nil
}

// UpdateLastLogin updates user's last login time
func (r *userRepositoryImpl) UpdateLastLogin(ctx context.Context, userID int64) error {
	query := `UPDATE users SET updated_at = ? WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query, time.Now(), userID)
	if err != nil {
		return fmt.Errorf("failed to update last login: %w", err)
	}

	return nil
}

// SetActive sets user active status
func (r *userRepositoryImpl) SetActive(ctx context.Context, userID int64, isActive bool) error {
	query := `UPDATE users SET is_active = ?, updated_at = ? WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query, isActive, time.Now(), userID)
	if err != nil {
		return fmt.Errorf("failed to set user active status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// EmailExists checks if email already exists
func (r *userRepositoryImpl) EmailExists(ctx context.Context, email string) (bool, error) {
	query := `SELECT COUNT(*) FROM users WHERE email = ?`

	var count int64
	err := r.db.GetContext(ctx, &count, query, email)
	if err != nil {
		return false, fmt.Errorf("failed to check email existence: %w", err)
	}

	return count > 0, nil
}
