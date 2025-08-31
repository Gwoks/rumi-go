package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/afghan/rumi-backend/internal/domain/entity"
	"github.com/afghan/rumi-backend/internal/domain/repository"
	"github.com/afghan/rumi-backend/internal/domain/usecase"
)

// childUsecaseImpl implements ChildUsecase interface
type childUsecaseImpl struct {
	childRepo repository.ChildRepository
	userRepo  repository.UserRepository
}

// NewChildUsecase creates a new child usecase
func NewChildUsecase(
	childRepo repository.ChildRepository,
	userRepo repository.UserRepository,
) usecase.ChildUsecase {
	return &childUsecaseImpl{
		childRepo: childRepo,
		userRepo:  userRepo,
	}
}

// CreateChild creates a new child for a user
func (u *childUsecaseImpl) CreateChild(ctx context.Context, userID int64, req *entity.CreateChildRequest) (*entity.Child, error) {
	// Validate request
	if req.Name == "" {
		return nil, fmt.Errorf("name is required")
	}
	if req.BirthDate == "" {
		return nil, fmt.Errorf("birth date is required")
	}

	// Parse birth date
	birthDate, err := time.Parse("2006-01-02", req.BirthDate)
	if err != nil {
		return nil, fmt.Errorf("invalid birth date format, expected YYYY-MM-DD")
	}

	// Check if birth date is not in the future
	if birthDate.After(time.Now()) {
		return nil, fmt.Errorf("birth date cannot be in the future")
	}

	// Verify user exists
	_, err = u.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Create child
	child := &entity.Child{
		UserID:    userID,
		Name:      req.Name,
		NickName:  req.NickName,
		BirthDate: birthDate,
		IsActive:  true,
	}

	createdChild, err := u.childRepo.Create(ctx, child)
	if err != nil {
		return nil, fmt.Errorf("failed to create child: %w", err)
	}

	return createdChild, nil
}

// GetChildByID retrieves a child by ID (with ownership check for non-admin users)
func (u *childUsecaseImpl) GetChildByID(ctx context.Context, childID int64, userID int64, userRole entity.UserRole) (*entity.Child, error) {
	// Admin can access any child
	if userRole == entity.RoleAdmin {
		child, err := u.childRepo.GetByID(ctx, childID)
		if err != nil {
			return nil, fmt.Errorf("child not found: %w", err)
		}
		return child, nil
	}

	// Regular users can only access their own children
	child, err := u.childRepo.GetByIDAndUserID(ctx, childID, userID)
	if err != nil {
		return nil, fmt.Errorf("child not found: %w", err)
	}

	return child, nil
}

// GetUserChildren retrieves children for a specific user
func (u *childUsecaseImpl) GetUserChildren(ctx context.Context, userID int64, limit, offset int) ([]*entity.Child, int64, error) {
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

	children, err := u.childRepo.GetByUserID(ctx, userID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get user children: %w", err)
	}

	total, err := u.childRepo.CountByUserID(ctx, userID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count user children: %w", err)
	}

	return children, total, nil
}

// GetAllChildren retrieves all children (admin only)
func (u *childUsecaseImpl) GetAllChildren(ctx context.Context, limit, offset int) ([]*entity.Child, int64, error) {
	// Validate pagination
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	children, err := u.childRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get all children: %w", err)
	}

	total, err := u.childRepo.Count(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count all children: %w", err)
	}

	return children, total, nil
}

// UpdateChild updates a child (with ownership check for non-admin users)
func (u *childUsecaseImpl) UpdateChild(ctx context.Context, childID int64, userID int64, userRole entity.UserRole, req *entity.UpdateChildRequest) (*entity.Child, error) {
	// Validate request
	if req.Name == "" {
		return nil, fmt.Errorf("name is required")
	}
	if req.BirthDate == "" {
		return nil, fmt.Errorf("birth date is required")
	}

	// Parse birth date
	birthDate, err := time.Parse("2006-01-02", req.BirthDate)
	if err != nil {
		return nil, fmt.Errorf("invalid birth date format, expected YYYY-MM-DD")
	}

	// Check if birth date is not in the future
	if birthDate.After(time.Now()) {
		return nil, fmt.Errorf("birth date cannot be in the future")
	}

	// Get existing child with ownership check
	var child *entity.Child
	if userRole == entity.RoleAdmin {
		child, err = u.childRepo.GetByID(ctx, childID)
	} else {
		child, err = u.childRepo.GetByIDAndUserID(ctx, childID, userID)
	}
	if err != nil {
		return nil, fmt.Errorf("child not found: %w", err)
	}

	// Update child fields
	child.Name = req.Name
	child.NickName = req.NickName
	child.BirthDate = birthDate

	// Only admin can change active status
	if req.IsActive != nil && userRole == entity.RoleAdmin {
		child.IsActive = *req.IsActive
	}

	// Update child
	updatedChild, err := u.childRepo.Update(ctx, child)
	if err != nil {
		return nil, fmt.Errorf("failed to update child: %w", err)
	}

	return updatedChild, nil
}

// DeleteChild soft deletes a child (with ownership check for non-admin users)
func (u *childUsecaseImpl) DeleteChild(ctx context.Context, childID int64, userID int64, userRole entity.UserRole) error {
	// Get existing child with ownership check
	var err error
	if userRole == entity.RoleAdmin {
		_, err = u.childRepo.GetByID(ctx, childID)
	} else {
		_, err = u.childRepo.GetByIDAndUserID(ctx, childID, userID)
	}
	if err != nil {
		return fmt.Errorf("child not found: %w", err)
	}

	// Soft delete child
	err = u.childRepo.Delete(ctx, childID)
	if err != nil {
		return fmt.Errorf("failed to delete child: %w", err)
	}

	return nil
}

// SetChildActive sets child active status (admin only)
func (u *childUsecaseImpl) SetChildActive(ctx context.Context, childID int64, isActive bool) error {
	// Verify child exists
	_, err := u.childRepo.GetByID(ctx, childID)
	if err != nil {
		return fmt.Errorf("child not found: %w", err)
	}

	// Set active status
	err = u.childRepo.SetActive(ctx, childID, isActive)
	if err != nil {
		return fmt.Errorf("failed to set child active status: %w", err)
	}

	return nil
}
