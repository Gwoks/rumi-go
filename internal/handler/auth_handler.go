package handler

import (
	"strings"

	"github.com/afghan/rumi-backend/internal/domain/entity"
	"github.com/afghan/rumi-backend/internal/domain/usecase"
	"github.com/afghan/rumi-backend/internal/middleware"
	"github.com/afghan/rumi-backend/pkg/response"
	"github.com/gin-gonic/gin"
)

// AuthHandler handles authentication related requests
type AuthHandler struct {
	authUsecase usecase.AuthUsecase
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authUsecase usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{
		authUsecase: authUsecase,
	}
}

// Login handles user login - POST /api/v1/auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req entity.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", err.Error())
		return
	}

	// Normalize email
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	authResponse, err := h.authUsecase.Login(c.Request.Context(), &req)
	if err != nil {
		if strings.Contains(err.Error(), "invalid credentials") {
			response.Unauthorized(c, "Invalid email or password")
			return
		}
		if strings.Contains(err.Error(), "deactivated") {
			response.Forbidden(c, "Account is deactivated")
			return
		}
		response.InternalServerError(c, "Login failed", err.Error())
		return
	}

	response.Success(c, "Login successful", authResponse)
}

// Signup handles user registration - POST /api/v1/auth/signup
func (h *AuthHandler) Signup(c *gin.Context) {
	var req entity.CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", err.Error())
		return
	}

	// Normalize email
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	// Trim whitespace from other fields
	req.Name = strings.TrimSpace(req.Name)
	req.Phone = strings.TrimSpace(req.Phone)

	// Set default role to user for signup
	req.Role = entity.RoleUser

	authResponse, err := h.authUsecase.Signup(c.Request.Context(), &req)
	if err != nil {
		if strings.Contains(err.Error(), "email already exists") {
			response.Conflict(c, "Email address is already registered")
			return
		}
		if strings.Contains(err.Error(), "password must be") {
			response.BadRequest(c, err.Error())
			return
		}
		response.InternalServerError(c, "Registration failed", err.Error())
		return
	}

	response.Created(c, "Account created successfully", authResponse)
}

// Logout handles user logout - POST /api/v1/auth/logout
func (h *AuthHandler) Logout(c *gin.Context) {
	// Extract token from Authorization header
	token := extractTokenFromHeader(c)
	if token == "" {
		response.BadRequest(c, "Authorization token required")
		return
	}

	err := h.authUsecase.Logout(c.Request.Context(), token)
	if err != nil {
		response.BadRequest(c, "Invalid token", err.Error())
		return
	}

	response.Success(c, "Logout successful", nil)
}

// RefreshToken handles token refresh - POST /api/v1/auth/refresh
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	// Extract token from Authorization header
	token := extractTokenFromHeader(c)
	if token == "" {
		response.BadRequest(c, "Authorization token required")
		return
	}

	authResponse, err := h.authUsecase.RefreshToken(c.Request.Context(), token)
	if err != nil {
		if strings.Contains(err.Error(), "invalid token") || strings.Contains(err.Error(), "expired") {
			response.Unauthorized(c, "Invalid or expired token")
			return
		}
		if strings.Contains(err.Error(), "deactivated") {
			response.Forbidden(c, "Account is deactivated")
			return
		}
		response.InternalServerError(c, "Token refresh failed", err.Error())
		return
	}

	response.Success(c, "Token refreshed successfully", authResponse)
}

// GetProfile handles getting user profile - GET /api/v1/auth/profile
func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	profile, err := h.authUsecase.GetProfile(c.Request.Context(), userID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			response.NotFound(c, "User not found")
			return
		}
		response.InternalServerError(c, "Failed to get profile", err.Error())
		return
	}

	response.Success(c, "Profile retrieved successfully", profile)
}

// ValidateToken handles token validation - POST /api/v1/auth/validate
func (h *AuthHandler) ValidateToken(c *gin.Context) {
	// Extract token from Authorization header
	token := extractTokenFromHeader(c)
	if token == "" {
		response.BadRequest(c, "Authorization token required")
		return
	}

	session, err := h.authUsecase.ValidateToken(c.Request.Context(), token)
	if err != nil {
		response.Unauthorized(c, "Invalid or expired token", err.Error())
		return
	}

	response.Success(c, "Token is valid", session)
}

// extractTokenFromHeader extracts JWT token from Authorization header
func extractTokenFromHeader(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return ""
	}

	// Expected format: "Bearer <token>"
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}

	return parts[1]
}
