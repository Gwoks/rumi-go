package usecase

import (
	"context"
	"fmt"
	"strings"

	"github.com/afghan/rumi-backend/internal/domain/entity"
	"github.com/afghan/rumi-backend/internal/domain/repository"
	"github.com/afghan/rumi-backend/internal/domain/usecase"
	"github.com/afghan/rumi-backend/pkg/auth"
)

// userManagementUsecaseImpl implements UserManagementUsecase interface
type userManagementUsecaseImpl struct {
	userRepo        repository.UserRepository
	activityRepo    repository.UserActivityRepository
	passwordService *auth.PasswordService
}

// NewUserManagementUsecase creates a new user management usecase
func NewUserManagementUsecase(
	userRepo repository.UserRepository,
	activityRepo repository.UserActivityRepository,
	passwordService *auth.PasswordService,
) usecase.UserManagementUsecase {
	return &userManagementUsecaseImpl{
		userRepo:        userRepo,
		activityRepo:    activityRepo,
		passwordService: passwordService,
	}
}

// GetAllUsers retrieves all users with pagination (admin only)
func (u *userManagementUsecaseImpl) GetAllUsers(ctx context.Context, limit, offset int) ([]*entity.PublicUser, int64, error) {
	// Validate pagination
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	users, err := u.userRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get users: %w", err)
	}

	total, err := u.userRepo.Count(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	// Convert to public users
	publicUsers := make([]*entity.PublicUser, len(users))
	for i, user := range users {
		publicUsers[i] = user.ToPublicUser()
	}

	return publicUsers, total, nil
}

// GetUserByID retrieves user by ID (admin only)
func (u *userManagementUsecaseImpl) GetUserByID(ctx context.Context, id int64) (*entity.PublicUser, error) {
	user, err := u.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user.ToPublicUser(), nil
}

// UpdateUser updates user information (admin only)
func (u *userManagementUsecaseImpl) UpdateUser(ctx context.Context, userID int64, req *entity.CreateUserRequest) (*entity.PublicUser, error) {
	// Get existing user
	user, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Validate request
	if req.Name == "" {
		return nil, fmt.Errorf("name is required")
	}
	if req.Phone == "" {
		return nil, fmt.Errorf("phone is required")
	}
	if req.Email == "" {
		return nil, fmt.Errorf("email is required")
	}

	// Check if email is being changed and if it already exists
	if user.Email != strings.ToLower(req.Email) {
		exists, err := u.userRepo.EmailExists(ctx, req.Email)
		if err != nil {
			return nil, fmt.Errorf("failed to check email existence: %w", err)
		}
		if exists {
			return nil, fmt.Errorf("email already exists")
		}
		user.Email = strings.ToLower(req.Email)
	}

	// Validate role
	if req.Role != "" && !req.Role.IsValid() {
		return nil, fmt.Errorf("invalid role")
	}

	// Update user fields
	user.Name = req.Name
	user.Phone = req.Phone
	if req.Role != "" {
		user.Role = req.Role
	}

	// Hash password if provided
	if req.Password != "" {
		err := u.passwordService.IsPasswordValid(req.Password)
		if err != nil {
			return nil, err
		}

		hashedPassword, err := u.passwordService.HashPassword(req.Password)
		if err != nil {
			return nil, fmt.Errorf("failed to hash password: %w", err)
		}
		user.Password = hashedPassword
	}

	// Update user
	updatedUser, err := u.userRepo.Update(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	// Log activity
	_ = u.logActivity(ctx, userID, entity.ActivityUserUpdated, fmt.Sprintf("User %s updated", user.Email))

	return updatedUser.ToPublicUser(), nil
}

// DeleteUser soft deletes a user (admin only)
func (u *userManagementUsecaseImpl) DeleteUser(ctx context.Context, userID int64) error {
	// Get user to log activity
	user, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Soft delete user
	err = u.userRepo.Delete(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	// Log activity
	_ = u.logActivity(ctx, userID, entity.ActivityUserDeleted, fmt.Sprintf("User %s deleted", user.Email))

	return nil
}

// SetUserActive sets user active status (admin only)
func (u *userManagementUsecaseImpl) SetUserActive(ctx context.Context, userID int64, isActive bool) error {
	// Get user to log activity
	user, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Set active status
	err = u.userRepo.SetActive(ctx, userID, isActive)
	if err != nil {
		return fmt.Errorf("failed to set user active status: %w", err)
	}

	// Log activity
	activityType := entity.ActivityUserActivated
	description := fmt.Sprintf("User %s activated", user.Email)
	if !isActive {
		activityType = entity.ActivityUserDeactivated
		description = fmt.Sprintf("User %s deactivated", user.Email)
	}
	_ = u.logActivity(ctx, userID, activityType, description)

	return nil
}

// UpdateProfile updates user's own profile (name, phone only)
func (u *userManagementUsecaseImpl) UpdateProfile(ctx context.Context, userID int64, req *entity.UpdateProfileRequest) (*entity.PublicUser, error) {
	// Get existing user
	user, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Update only allowed fields
	user.Name = req.Name
	user.Phone = req.Phone

	// Update user
	updatedUser, err := u.userRepo.Update(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to update profile: %w", err)
	}

	// Log activity
	_ = u.logActivity(ctx, userID, entity.ActivityProfileUpdate, "Profile updated")

	return updatedUser.ToPublicUser(), nil
}

// ChangePassword changes user's password
func (u *userManagementUsecaseImpl) ChangePassword(ctx context.Context, userID int64, req *entity.ChangePasswordRequest) error {
	// Get user
	user, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Verify current password
	err = u.passwordService.VerifyPassword(user.Password, req.CurrentPassword)
	if err != nil {
		return fmt.Errorf("current password is incorrect")
	}

	// Validate new password
	err = u.passwordService.IsPasswordValid(req.NewPassword)
	if err != nil {
		return err
	}

	// Hash new password
	hashedPassword, err := u.passwordService.HashPassword(req.NewPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Update password
	user.Password = hashedPassword
	_, err = u.userRepo.Update(ctx, user)
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	// Log activity
	_ = u.logActivity(ctx, userID, entity.ActivityPasswordChange, "Password changed")

	return nil
}

// logActivity is a helper method to log user activities
func (u *userManagementUsecaseImpl) logActivity(ctx context.Context, userID int64, activityType entity.ActivityType, description string) error {
	activity := &entity.UserActivity{
		UserID:       userID,
		ActivityType: activityType,
		Description:  description,
		IPAddress:    "127.0.0.1", // Could be extracted from context
		UserAgent:    "System",    // Could be extracted from context
	}

	_, err := u.activityRepo.Create(ctx, activity)
	return err
}
