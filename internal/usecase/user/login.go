package user

import (
	"context"
	"fmt"
	"time"

	"rumi-go/internal/config"
	"rumi-go/internal/model"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// LoginUser authenticates a user and returns user info with JWT token
func (u *UserUsecase) LoginUser(ctx context.Context, req model.LoginRequest) (*model.LoginResponse, error) {
	// Get user by email
	user, err := u.database.UserStore().GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, err
	}

	// Generate JWT token
	token, err := u.generateJWTToken(user)
	if err != nil {
		return nil, err
	}

	return &model.LoginResponse{
		Token: token,
	}, nil
}

// generateJWTToken generates a JWT token for the user
func (u *UserUsecase) generateJWTToken(user *model.User) (string, error) {
	// Get JWT configuration
	jwtConfig := config.Get().JWT

	// Create claims
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Duration(jwtConfig.Expiration) * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
		"iss":     "rumi-go",
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token with secret
	tokenString, err := token.SignedString([]byte(jwtConfig.Secret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}
