package middleware

import (
	"time"

	"github.com/afghan/rumi-backend/internal/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SetupCORS configures CORS middleware
func SetupCORS(cfg *config.Config) gin.HandlerFunc {
	corsConfig := cors.Config{
		AllowOrigins:     cfg.CORS.AllowedOrigins,
		AllowMethods:     cfg.CORS.AllowedMethods,
		AllowHeaders:     cfg.CORS.AllowedHeaders,
		ExposeHeaders:    []string{"Content-Length", "X-Request-ID"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	// If development environment and origins contain "*", allow all origins
	if cfg.IsDevelopment() && containsWildcard(cfg.CORS.AllowedOrigins) {
		corsConfig.AllowAllOrigins = true
		corsConfig.AllowOrigins = nil // Must be nil when AllowAllOrigins is true
	}

	return cors.New(corsConfig)
}

// containsWildcard checks if slice contains "*"
func containsWildcard(slice []string) bool {
	for _, item := range slice {
		if item == "*" {
			return true
		}
	}
	return false
}
