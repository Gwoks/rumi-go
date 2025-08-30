package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/afghan/rumi-backend/internal/config"
	"github.com/afghan/rumi-backend/internal/database"
	"github.com/afghan/rumi-backend/internal/handler"
	"github.com/afghan/rumi-backend/internal/infrastructure/repository"
	"github.com/afghan/rumi-backend/internal/infrastructure/usecase"
	"github.com/afghan/rumi-backend/internal/middleware"
	"github.com/afghan/rumi-backend/pkg/auth"
	"github.com/afghan/rumi-backend/pkg/response"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Set Gin mode based on environment
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	// Connect to database
	db, err := database.NewMySQLConnection(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize services
	jwtService := auth.NewJWTService(cfg)
	passwordService := auth.NewPasswordService(cfg)

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	userActivityRepo := repository.NewUserActivityRepository(db)

	// Initialize use cases
	authUsecase := usecase.NewAuthUsecase(userRepo, userActivityRepo, jwtService, passwordService)
	userManagementUsecase := usecase.NewUserManagementUsecase(userRepo, userActivityRepo, passwordService)
	userActivityUsecase := usecase.NewUserActivityUsecase(userActivityRepo, userRepo)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authUsecase)
	userHandler := handler.NewUserHandler(userManagementUsecase, userActivityUsecase)

	// Initialize middlewares
	authMiddleware := middleware.NewAuthMiddleware(authUsecase)

	// Setup router
	router := setupRouter(cfg, authHandler, userHandler, authMiddleware)

	// Setup server
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("🚀 Server starting on %s:%d", cfg.Server.Host, cfg.Server.Port)
		log.Printf("📝 Environment: %s", cfg.Server.Environment)
		log.Printf("🗄️ Database: %s@%s:%d/%s", cfg.Database.User, cfg.Database.Host, cfg.Database.Port, cfg.Database.Name)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Server shutting down...")

	// Give outstanding requests 30 seconds to complete
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("✅ Server gracefully stopped")
}

// setupRouter configures the Gin router with all routes and middleware
func setupRouter(cfg *config.Config, authHandler *handler.AuthHandler, userHandler *handler.UserHandler, authMiddleware *middleware.AuthMiddleware) *gin.Engine {
	router := gin.New()

	// Global middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.SetupCORS(cfg))

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		response.Success(c, "Server is healthy", gin.H{
			"status":      "ok",
			"environment": cfg.Server.Environment,
			"timestamp":   time.Now(),
		})
	})

	// API routes
	api := router.Group("/api/v1")
	{
		// Auth routes (matching Angular API service)
		auth := api.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/signup", authHandler.Signup)
			auth.POST("/logout", authHandler.Logout)
			auth.POST("/refresh", authHandler.RefreshToken)
			auth.POST("/validate", authHandler.ValidateToken)

			// Protected routes
			auth.Use(authMiddleware.RequireAuth())
			auth.GET("/profile", authHandler.GetProfile)
		}

		// Protected routes
		protected := api.Group("")
		protected.Use(authMiddleware.RequireAuth())
		{
			// User profile routes
			protected.GET("/profile", authHandler.GetProfile)
			protected.PUT("/profile", userHandler.UpdateProfile)
			protected.PUT("/profile/password", userHandler.ChangePassword)

			// User activities (own activities)
			protected.GET("/my-activities", func(c *gin.Context) {
				userID, _ := middleware.GetUserIDFromContext(c)
				c.Params = append(c.Params, gin.Param{Key: "id", Value: fmt.Sprintf("%d", userID)})
				userHandler.GetUserActivities(c)
			})

			// Admin only routes
			admin := protected.Group("")
			admin.Use(authMiddleware.RequireAdmin())
			{
				// User management
				admin.GET("/users", userHandler.GetAllUsers)
				admin.GET("/users/:id", userHandler.GetUser)
				admin.PUT("/users/:id", userHandler.UpdateUser)
				admin.DELETE("/users/:id", userHandler.DeleteUser)
				admin.PUT("/users/:id/active", userHandler.SetUserActive)

				// User activities management
				admin.GET("/users/:id/activities", userHandler.GetUserActivities)
				admin.GET("/activities", userHandler.GetAllActivities)
			}
		}
	}

	// Catch all unmatched routes
	router.NoRoute(func(c *gin.Context) {
		response.NotFound(c, "Endpoint not found")
	})

	return router
}
