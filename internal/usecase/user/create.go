package user

import (
	"context"
	"errors"
	"fmt"

	"rumi-go/internal/model"

	"golang.org/x/crypto/bcrypt"
)

// CreateUser creates a new user
func (u *UserUsecase) CreateUser(ctx context.Context, req model.UserRequest) (*model.UserResponse, error) {
	// Check if user already exists
	existingUser, err := u.database.UserStore().GetByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := &model.User{
		Email:    req.Email,
		Password: string(hashedPassword),
		Name:     req.Name,
		Phone:    stringToPtr(req.Phone),
		Role:     req.Role,
		Address:  stringToPtr(req.Address),
	}

	if err := u.database.UserStore().Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return u.toUserResponse(user), nil
}
