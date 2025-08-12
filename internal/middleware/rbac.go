package middleware

import (
	"net/http"
)

// RoleConfig defines which roles can access an endpoint
type RoleConfig struct {
	AllowedRoles []string
	GuestAllowed bool // Allow unauthenticated access
}

// RBACMiddleware creates middleware for role-based access control
func RBACMiddleware(config RoleConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// If guest access is allowed, skip authentication
			if config.GuestAllowed {
				next.ServeHTTP(w, r)
				return
			}

			// Get user from context (set by AuthMiddleware)
			user, ok := r.Context().Value("user").(*User)
			if !ok {
				http.Error(w, "Authentication required", http.StatusUnauthorized)
				return
			}

			// Check if user's role is allowed
			if !isRoleAllowed(user.Role, config.AllowedRoles) {
				http.Error(w, "Insufficient permissions", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// isRoleAllowed checks if the user's role is in the allowed roles list
func isRoleAllowed(userRole string, allowedRoles []string) bool {
	for _, role := range allowedRoles {
		if role == userRole {
			return true
		}
	}
	return false
}

// Predefined role configurations for common use cases
var (
	// AdminOnly - only admin users can access
	AdminOnly = RoleConfig{
		AllowedRoles: []string{"admin"},
		GuestAllowed: false,
	}

	// AdminAndUser - admin and user roles can access
	AdminAndUser = RoleConfig{
		AllowedRoles: []string{"admin", "user"},
		GuestAllowed: false,
	}

	// UserOnly - only authenticated users (any role) can access
	UserOnly = RoleConfig{
		AllowedRoles: []string{"admin", "user", "manager"},
		GuestAllowed: false,
	}

	// GuestOnly - only unauthenticated users can access
	GuestOnly = RoleConfig{
		AllowedRoles: []string{},
		GuestAllowed: true,
	}

	// Public - anyone can access (no restrictions)
	Public = RoleConfig{
		AllowedRoles: []string{"admin", "user", "manager"},
		GuestAllowed: true,
	}
)

// Helper functions for common role combinations
func AdminAndManager() RoleConfig {
	return RoleConfig{
		AllowedRoles: []string{"admin", "manager"},
		GuestAllowed: false,
	}
}

func UserAndAbove() RoleConfig {
	return RoleConfig{
		AllowedRoles: []string{"admin", "user", "manager"},
		GuestAllowed: false,
	}
}

func CustomRoles(roles ...string) RoleConfig {
	return RoleConfig{
		AllowedRoles: roles,
		GuestAllowed: false,
	}
}
