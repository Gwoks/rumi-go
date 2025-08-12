package user

import (
	"context"
	"fmt"

	"rumi-go/internal/model"

	"golang.org/x/crypto/bcrypt"
)

// UpdateUser updates user information
func (u *UserUsecase) UpdateUser(ctx context.Context, userID int64, req model.UserRequest) (*model.UserResponse, error) {
	// Get existing user
	existingUser, err := u.database.UserStore().Get(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Update fields
	existingUser.Name = req.Name
	existingUser.Phone = stringToPtr(req.Phone)
	existingUser.Role = req.Role
	existingUser.Address = stringToPtr(req.Address)

	// Update password if provided
	if req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, fmt.Errorf("failed to hash password: %w", err)
		}
		existingUser.Password = string(hashedPassword)
	}

	if err := u.database.UserStore().Update(ctx, existingUser); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return u.toUserResponse(existingUser), nil
}
