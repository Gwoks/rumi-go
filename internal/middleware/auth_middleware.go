package middleware

import (
	"strings"

	"github.com/afghan/rumi-backend/internal/domain/entity"
	"github.com/afghan/rumi-backend/internal/domain/usecase"
	"github.com/afghan/rumi-backend/pkg/response"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware provides JWT authentication middleware
type AuthMiddleware struct {
	authUsecase usecase.AuthUsecase
}

// NewAuthMiddleware creates a new auth middleware
func NewAuthMiddleware(authUsecase usecase.AuthUsecase) *AuthMiddleware {
	return &AuthMiddleware{
		authUsecase: authUsecase,
	}
}

// RequireAuth validates JWT token and sets user context
func (a *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractTokenFromHeader(c)
		if token == "" {
			response.Unauthorized(c, "Authorization token required")
			c.Abort()
			return
		}

		session, err := a.authUsecase.ValidateToken(c.Request.Context(), token)
		if err != nil {
			response.Unauthorized(c, "Invalid or expired token", err.Error())
			c.Abort()
			return
		}

		// Set user session in context
		c.Set("user_id", session.UserID)
		c.Set("user_email", session.Email)
		c.Set("user_role", session.Role)
		c.Set("session", session)

		c.Next()
	}
}

// RequireRole validates user role
func (a *AuthMiddleware) RequireRole(roles ...entity.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user role from context (set by RequireAuth middleware)
		userRole, exists := c.Get("user_role")
		if !exists {
			response.Unauthorized(c, "Authentication required")
			c.Abort()
			return
		}

		role, ok := userRole.(entity.UserRole)
		if !ok {
			response.InternalServerError(c, "Invalid user role type")
			c.Abort()
			return
		}

		// Check if user has required role
		hasRole := false
		for _, requiredRole := range roles {
			if role == requiredRole {
				hasRole = true
				break
			}
		}

		if !hasRole {
			response.Forbidden(c, "Insufficient permissions")
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireAdmin validates admin role
func (a *AuthMiddleware) RequireAdmin() gin.HandlerFunc {
	return a.RequireRole(entity.RoleAdmin)
}

// OptionalAuth validates JWT token if present but doesn't require it
func (a *AuthMiddleware) OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractTokenFromHeader(c)
		if token == "" {
			// No token provided, continue without auth
			c.Next()
			return
		}

		session, err := a.authUsecase.ValidateToken(c.Request.Context(), token)
		if err != nil {
			// Invalid token, continue without auth
			c.Next()
			return
		}

		// Set user session in context
		c.Set("user_id", session.UserID)
		c.Set("user_email", session.Email)
		c.Set("user_role", session.Role)
		c.Set("session", session)

		c.Next()
	}
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

// GetUserIDFromContext gets user ID from gin context
func GetUserIDFromContext(c *gin.Context) (int64, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}

	id, ok := userID.(int64)
	return id, ok
}

// GetUserEmailFromContext gets user email from gin context
func GetUserEmailFromContext(c *gin.Context) (string, bool) {
	email, exists := c.Get("user_email")
	if !exists {
		return "", false
	}

	userEmail, ok := email.(string)
	return userEmail, ok
}

// GetUserRoleFromContext gets user role from gin context
func GetUserRoleFromContext(c *gin.Context) (entity.UserRole, bool) {
	role, exists := c.Get("user_role")
	if !exists {
		return "", false
	}

	userRole, ok := role.(entity.UserRole)
	return userRole, ok
}

// GetSessionFromContext gets user session from gin context
func GetSessionFromContext(c *gin.Context) (*entity.Session, bool) {
	session, exists := c.Get("session")
	if !exists {
		return nil, false
	}

	userSession, ok := session.(*entity.Session)
	return userSession, ok
}
