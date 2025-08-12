package user

import (
	"context"
	"fmt"

	"rumi-go/internal/model"
)

// GetUserInfo retrieves user information by ID
func (u *UserUsecase) GetUserInfo(ctx context.Context, userID int64) (*model.UserResponse, error) {
	user, err := u.database.UserStore().Get(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	return u.toUserResponse(user), nil
}

// GetUserByEmail retrieves user information by email
func (u *UserUsecase) GetUserByEmail(ctx context.Context, email string) (*model.UserResponse, error) {
	user, err := u.database.UserStore().GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	return u.toUserResponse(user), nil
}

// ListUsers retrieves a list of users with pagination
func (u *UserUsecase) ListUsers(ctx context.Context, limit, offset int) ([]*model.UserResponse, error) {
	users, err := u.database.UserStore().List(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	var responses []*model.UserResponse
	for _, user := range users {
		responses = append(responses, u.toUserResponse(user))
	}

	return responses, nil
}
