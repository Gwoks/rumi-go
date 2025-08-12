package user

import (
	"time"

	"rumi-go/internal/infrastructure/database"
	"rumi-go/internal/model"
)

const (
	// Transfer success message
	TransferSuccessMessage = "balance added, journal entry and deposit history created"
)

type UserUsecase struct {
	location *time.Location
	database database.Store
}

func NewUserUsecase(database database.Store) *UserUsecase {
	loc, _ := time.LoadLocation("Asia/Jakarta")

	return &UserUsecase{
		location: loc,
		database: database,
	}
}

// toUserResponse converts User model to UserResponse
func (u *UserUsecase) toUserResponse(user *model.User) *model.UserResponse {
	return &model.UserResponse{
		ID:      user.ID,
		Email:   user.Email,
		Name:    user.Name,
		Phone:   user.Phone,
		Role:    user.Role,
		Address: user.Address,
	}
}

// stringToPtr converts a string to a pointer, returns nil if empty
func stringToPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
