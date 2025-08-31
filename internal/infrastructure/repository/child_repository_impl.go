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

// childRepositoryImpl implements ChildRepository interface
type childRepositoryImpl struct {
	db *database.DB
}

// NewChildRepository creates a new child repository
func NewChildRepository(db *database.DB) repository.ChildRepository {
	return &childRepositoryImpl{
		db: db,
	}
}

// Create creates a new child
func (r *childRepositoryImpl) Create(ctx context.Context, child *entity.Child) (*entity.Child, error) {
	query := `
		INSERT INTO children (user_id, name, nick_name, birth_date, is_active, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)`

	now := time.Now()
	child.CreatedAt = now
	child.UpdatedAt = now

	result, err := r.db.ExecContext(ctx, query,
		child.UserID,
		child.Name,
		child.NickName,
		child.BirthDate,
		child.IsActive,
		child.CreatedAt,
		child.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create child: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}

	child.ID = id
	return child, nil
}

// GetByID retrieves a child by ID
func (r *childRepositoryImpl) GetByID(ctx context.Context, id int64) (*entity.Child, error) {
	query := `
		SELECT c.id, c.user_id, c.name, c.nick_name, c.birth_date, c.is_active, c.created_at, c.updated_at,
			   u.email, u.name as user_name, u.phone, u.role, u.is_active as user_is_active, 
			   u.created_at as user_created_at, u.updated_at as user_updated_at
		FROM children c
		LEFT JOIN users u ON c.user_id = u.id
		WHERE c.id = ? AND c.is_active = TRUE`

	var child entity.Child
	var user entity.PublicUser

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&child.ID,
		&child.UserID,
		&child.Name,
		&child.NickName,
		&child.BirthDate,
		&child.IsActive,
		&child.CreatedAt,
		&child.UpdatedAt,
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
			return nil, fmt.Errorf("child not found")
		}
		return nil, fmt.Errorf("failed to get child by id: %w", err)
	}

	user.ID = child.UserID
	child.User = &user

	return &child, nil
}

// GetByUserID retrieves children by user ID with pagination
func (r *childRepositoryImpl) GetByUserID(ctx context.Context, userID int64, limit, offset int) ([]*entity.Child, error) {
	query := `
		SELECT id, user_id, name, nick_name, birth_date, is_active, created_at, updated_at
		FROM children
		WHERE user_id = ? AND is_active = TRUE
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?`

	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get children by user id: %w", err)
	}
	defer rows.Close()

	var children []*entity.Child
	for rows.Next() {
		var child entity.Child

		err := rows.Scan(
			&child.ID,
			&child.UserID,
			&child.Name,
			&child.NickName,
			&child.BirthDate,
			&child.IsActive,
			&child.CreatedAt,
			&child.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan child: %w", err)
		}

		children = append(children, &child)
	}

	return children, nil
}

// Update updates an existing child
func (r *childRepositoryImpl) Update(ctx context.Context, child *entity.Child) (*entity.Child, error) {
	query := `
		UPDATE children 
		SET name = ?, nick_name = ?, birth_date = ?, is_active = ?, updated_at = ?
		WHERE id = ?`

	child.UpdatedAt = time.Now()

	result, err := r.db.ExecContext(ctx, query,
		child.Name,
		child.NickName,
		child.BirthDate,
		child.IsActive,
		child.UpdatedAt,
		child.ID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update child: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return nil, fmt.Errorf("child not found or no changes made")
	}

	return child, nil
}

// Delete soft deletes a child by ID
func (r *childRepositoryImpl) Delete(ctx context.Context, id int64) error {
	query := `UPDATE children SET is_active = FALSE, updated_at = ? WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to delete child: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("child not found")
	}

	return nil
}

// List retrieves children with pagination and joins user info (admin only)
func (r *childRepositoryImpl) List(ctx context.Context, limit, offset int) ([]*entity.Child, error) {
	query := `
		SELECT c.id, c.user_id, c.name, c.nick_name, c.birth_date, c.is_active, c.created_at, c.updated_at,
			   u.email, u.name as user_name, u.phone, u.role, u.is_active as user_is_active, 
			   u.created_at as user_created_at, u.updated_at as user_updated_at
		FROM children c
		LEFT JOIN users u ON c.user_id = u.id
		WHERE c.is_active = TRUE
		ORDER BY c.created_at DESC
		LIMIT ? OFFSET ?`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list children: %w", err)
	}
	defer rows.Close()

	var children []*entity.Child
	for rows.Next() {
		var child entity.Child
		var user entity.PublicUser

		err := rows.Scan(
			&child.ID,
			&child.UserID,
			&child.Name,
			&child.NickName,
			&child.BirthDate,
			&child.IsActive,
			&child.CreatedAt,
			&child.UpdatedAt,
			&user.Email,
			&user.Name,
			&user.Phone,
			&user.Role,
			&user.IsActive,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan child: %w", err)
		}

		user.ID = child.UserID
		child.User = &user
		children = append(children, &child)
	}

	return children, nil
}

// Count returns total number of children
func (r *childRepositoryImpl) Count(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM children WHERE is_active = TRUE`

	var count int64
	err := r.db.GetContext(ctx, &count, query)
	if err != nil {
		return 0, fmt.Errorf("failed to count children: %w", err)
	}

	return count, nil
}

// CountByUserID returns total number of children for a user
func (r *childRepositoryImpl) CountByUserID(ctx context.Context, userID int64) (int64, error) {
	query := `SELECT COUNT(*) FROM children WHERE user_id = ? AND is_active = TRUE`

	var count int64
	err := r.db.GetContext(ctx, &count, query, userID)
	if err != nil {
		return 0, fmt.Errorf("failed to count children by user id: %w", err)
	}

	return count, nil
}

// SetActive sets child active status
func (r *childRepositoryImpl) SetActive(ctx context.Context, childID int64, isActive bool) error {
	query := `UPDATE children SET is_active = ?, updated_at = ? WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query, isActive, time.Now(), childID)
	if err != nil {
		return fmt.Errorf("failed to set child active status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("child not found")
	}

	return nil
}

// GetByIDAndUserID retrieves a child by ID and user ID (for ownership check)
func (r *childRepositoryImpl) GetByIDAndUserID(ctx context.Context, id, userID int64) (*entity.Child, error) {
	query := `
		SELECT id, user_id, name, nick_name, birth_date, is_active, created_at, updated_at
		FROM children
		WHERE id = ? AND user_id = ? AND is_active = TRUE`

	var child entity.Child
	err := r.db.GetContext(ctx, &child, query, id, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("child not found")
		}
		return nil, fmt.Errorf("failed to get child by id and user id: %w", err)
	}

	return &child, nil
}
