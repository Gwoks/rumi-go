package usecase

import (
	"context"
	"fmt"

	"github.com/afghan/rumi-backend/internal/domain/entity"
	"github.com/afghan/rumi-backend/internal/domain/repository"
	"github.com/afghan/rumi-backend/internal/domain/usecase"
	"github.com/afghan/rumi-backend/pkg/auth"
)

// authUsecaseImpl implements AuthUsecase interface
type authUsecaseImpl struct {
	userRepo        repository.UserRepository
	activityRepo    repository.UserActivityRepository
	jwtService      *auth.JWTService
	passwordService *auth.PasswordService
}

// NewAuthUsecase creates a new auth usecase
func NewAuthUsecase(
	userRepo repository.UserRepository,
	activityRepo repository.UserActivityRepository,
	jwtService *auth.JWTService,
	passwordService *auth.PasswordService,
) usecase.AuthUsecase {
	return &authUsecaseImpl{
		userRepo:        userRepo,
		activityRepo:    activityRepo,
		jwtService:      jwtService,
		passwordService: passwordService,
	}
}

// Login authenticates user and returns JWT token
func (a *authUsecaseImpl) Login(ctx context.Context, req *entity.LoginRequest) (*entity.AuthResponse, error) {
	// Validate request
	if req.Email == "" {
		return nil, fmt.Errorf("email is required")
	}
	if req.Password == "" {
		return nil, fmt.Errorf("password is required")
	}

	// Get user by email
	user, err := a.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Check if user is active
	if !user.IsActive {
		return nil, fmt.Errorf("account is deactivated")
	}

	// Verify password
	err = a.passwordService.VerifyPassword(user.Password, req.Password)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Generate JWT token
	token, err := a.jwtService.GenerateToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Update last login
	err = a.userRepo.UpdateLastLogin(ctx, user.ID)
	if err != nil {
		// Log but don't fail the login
		fmt.Printf("failed to update last login for user %d: %v\n", user.ID, err)
	}

	// Log activity
	_ = a.logActivity(ctx, user.ID, entity.ActivityLogin, "User logged in")

	return &entity.AuthResponse{
		User:  user.ToPublicUser(),
		Token: token,
	}, nil
}

// Signup creates a new user account and returns JWT token
func (a *authUsecaseImpl) Signup(ctx context.Context, req *entity.CreateUserRequest) (*entity.AuthResponse, error) {
	// Validate request
	if req.Email == "" {
		return nil, fmt.Errorf("email is required")
	}
	if req.Name == "" {
		return nil, fmt.Errorf("name is required")
	}
	if req.Phone == "" {
		return nil, fmt.Errorf("phone is required")
	}
	if req.Password == "" {
		return nil, fmt.Errorf("password is required")
	}

	// Validate password requirements
	err := a.passwordService.IsPasswordValid(req.Password)
	if err != nil {
		return nil, err
	}

	// Check if email already exists
	exists, err := a.userRepo.EmailExists(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check email: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("email already exists")
	}

	// Hash password
	hashedPassword, err := a.passwordService.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Set default role if not specified
	role := req.Role
	if role == "" {
		role = entity.RoleUser
	}

	// Validate role
	if !role.IsValid() {
		return nil, fmt.Errorf("invalid role")
	}

	// Create user
	user := &entity.User{
		Email:    req.Email,
		Name:     req.Name,
		Phone:    req.Phone,
		Password: hashedPassword,
		Role:     role,
		IsActive: true,
	}

	createdUser, err := a.userRepo.Create(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Generate JWT token
	token, err := a.jwtService.GenerateToken(createdUser)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Log activity
	_ = a.logActivity(ctx, createdUser.ID, entity.ActivityAccountCreated, "Account created")

	return &entity.AuthResponse{
		User:  createdUser.ToPublicUser(),
		Token: token,
	}, nil
}

// Logout invalidates user session
func (a *authUsecaseImpl) Logout(ctx context.Context, token string) error {
	// For JWT tokens, we can't easily invalidate them without maintaining a blacklist
	// This is a simple implementation that just validates the token
	// In production, you might want to implement token blacklisting with Redis

	_, err := a.jwtService.ValidateToken(token)
	if err != nil {
		return fmt.Errorf("invalid token")
	}

	// Here you could add logic to blacklist the token
	// For now, we'll just return success
	return nil
}

// RefreshToken generates a new JWT token from existing valid token
func (a *authUsecaseImpl) RefreshToken(ctx context.Context, token string) (*entity.AuthResponse, error) {
	// Get session from token
	session, err := a.jwtService.GetSessionFromToken(token)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	// Get user to ensure they still exist and are active
	user, err := a.userRepo.GetByID(ctx, session.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	if !user.IsActive {
		return nil, fmt.Errorf("account is deactivated")
	}

	// Generate new token
	newToken, err := a.jwtService.GenerateToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &entity.AuthResponse{
		User:  user.ToPublicUser(),
		Token: newToken,
	}, nil
}

// ValidateToken validates JWT token and returns user session data
func (a *authUsecaseImpl) ValidateToken(ctx context.Context, token string) (*entity.Session, error) {
	session, err := a.jwtService.GetSessionFromToken(token)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	// Optionally verify user still exists and is active
	user, err := a.userRepo.GetByID(ctx, session.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	if !user.IsActive {
		return nil, fmt.Errorf("account is deactivated")
	}

	return session, nil
}

// GetProfile retrieves user profile by ID
func (a *authUsecaseImpl) GetProfile(ctx context.Context, userID int64) (*entity.PublicUser, error) {
	user, err := a.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	return user.ToPublicUser(), nil
}

// logActivity is a helper method to log user activities
func (a *authUsecaseImpl) logActivity(ctx context.Context, userID int64, activityType entity.ActivityType, description string) error {
	activity := &entity.UserActivity{
		UserID:       userID,
		ActivityType: activityType,
		Description:  description,
		IPAddress:    "127.0.0.1", // Could be extracted from context
		UserAgent:    "System",    // Could be extracted from context
	}

	_, err := a.activityRepo.Create(ctx, activity)
	return err
}
