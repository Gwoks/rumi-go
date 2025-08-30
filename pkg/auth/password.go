package auth

import (
	"fmt"

	"github.com/afghan/rumi-backend/internal/config"
	"golang.org/x/crypto/bcrypt"
)

// PasswordService provides password hashing operations
type PasswordService struct {
	cost int
}

// NewPasswordService creates a new password service
func NewPasswordService(cfg *config.Config) *PasswordService {
	return &PasswordService{
		cost: cfg.Security.BCryptCost,
	}
}

// HashPassword hashes a plain text password
func (p *PasswordService) HashPassword(password string) (string, error) {
	if password == "" {
		return "", fmt.Errorf("password cannot be empty")
	}

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), p.cost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	return string(hashedBytes), nil
}

// VerifyPassword verifies a plain text password against a hashed password
func (p *PasswordService) VerifyPassword(hashedPassword, plainPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return fmt.Errorf("invalid password")
		}
		return fmt.Errorf("failed to verify password: %w", err)
	}

	return nil
}

// IsPasswordValid checks if password meets minimum requirements
func (p *PasswordService) IsPasswordValid(password string) error {
	if len(password) < 6 {
		return fmt.Errorf("password must be at least 6 characters long")
	}

	// Add more password validation rules here if needed
	// e.g., must contain uppercase, lowercase, numbers, special characters

	return nil
}
