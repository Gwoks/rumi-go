package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Database Database `json:"database"`
	Server   Server   `json:"server"`
	JWT      JWT      `json:"jwt"`
	Security Security `json:"security"`
	CORS     CORS     `json:"cors"`
}

type Database struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type Server struct {
	Host        string `json:"host"`
	Port        int    `json:"port"`
	Environment string `json:"environment"`
}

type JWT struct {
	Secret      string `json:"secret"`
	ExpiryHours int    `json:"expiry_hours"`
}

type Security struct {
	BCryptCost int `json:"bcrypt_cost"`
}

type CORS struct {
	AllowedOrigins []string `json:"allowed_origins"`
	AllowedMethods []string `json:"allowed_methods"`
	AllowedHeaders []string `json:"allowed_headers"`
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Load .env file if exists
	_ = godotenv.Load()

	config := &Config{
		Database: Database{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvAsInt("DB_PORT", 3306),
			User:     getEnv("DB_USER", "root"),
			Password: getEnv("DB_PASSWORD", ""),
			Name:     getEnv("DB_NAME", "rumi_db"),
		},
		Server: Server{
			Host:        getEnv("SERVER_HOST", "0.0.0.0"),
			Port:        getEnvAsInt("SERVER_PORT", 8080),
			Environment: getEnv("ENVIRONMENT", "development"),
		},
		JWT: JWT{
			Secret:      getEnv("JWT_SECRET", "default-secret-change-me"),
			ExpiryHours: getEnvAsInt("JWT_EXPIRY_HOURS", 24),
		},
		Security: Security{
			BCryptCost: getEnvAsInt("BCRYPT_COST", 12),
		},
		CORS: CORS{
			AllowedOrigins: getEnvAsSlice("CORS_ALLOWED_ORIGINS", []string{"*"}),
			AllowedMethods: getEnvAsSlice("CORS_ALLOWED_METHODS", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
			AllowedHeaders: getEnvAsSlice("CORS_ALLOWED_HEADERS", []string{"Content-Type", "Authorization"}),
		},
	}

	if err := config.validate(); err != nil {
		return nil, err
	}

	return config, nil
}

// validate checks if required configuration values are set
func (c *Config) validate() error {
	if c.Database.Password == "" {
		return fmt.Errorf("DB_PASSWORD is required")
	}
	if c.JWT.Secret == "default-secret-change-me" && c.Server.Environment == "production" {
		return fmt.Errorf("JWT_SECRET must be set in production")
	}
	return nil
}

// GetDSN returns database connection string
func (c *Config) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.Database.User,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.Name,
	)
}

// IsDevelopment returns true if environment is development
func (c *Config) IsDevelopment() bool {
	return c.Server.Environment == "development"
}

// IsProduction returns true if environment is production
func (c *Config) IsProduction() bool {
	return c.Server.Environment == "production"
}

// Helper functions
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return fallback
}

func getEnvAsSlice(key string, fallback []string) []string {
	if value := os.Getenv(key); value != "" {
		return strings.Split(value, ",")
	}
	return fallback
}
