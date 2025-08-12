package config

import (
	"fmt"
	"log"
	"strings"

	"rumi-go/internal/database"

	"github.com/spf13/viper"
)

var c Config

type Config struct {
	App         AppConfig       `mapstructure:"app"`
	Server      ServerConfig    `mapstructure:"server"`
	Redis       RedisConfig     `mapstructure:"redis"`
	JWT         JWTConfig       `mapstructure:"jwt"`
	Database    database.Config `mapstructure:"database"`
	LogDatabase database.Config `mapstructure:"log_database"`
}

type AppConfig struct {
	Name      string `mapstructure:"name"`
	Env       string `mapstructure:"env"`
	ApiPrefix string `mapstructure:"api_prefix"`
}

type ServerConfig struct {
	Port int `mapstructure:"port"`
}

type RedisConfig struct {
	Address  string `mapstructure:"address"`
	Password string `mapstructure:"password"`
	Port     string `mapstructure:"port"`
	DB       int    `mapstructure:"db"`
}

type JWTConfig struct {
	Secret     string `mapstructure:"secret"`
	Expiration int    `mapstructure:"expiration"` // in hours
}

// Load loads configuration from file and environment variables using Viper
func Load(configPath string) error {
	v := viper.New()

	// Set default config file path
	if configPath == "" {
		configPath = "."
	}

	// Set config file name
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(configPath)

	// Enable reading from environment variables
	v.AutomaticEnv()
	v.SetEnvPrefix("RUMI")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Read config file
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return fmt.Errorf("error reading config file: %w", err)
		}
		// Config file not found, continue with defaults and env vars
	}

	// Unmarshal configuration into struct
	if err := v.Unmarshal(&c); err != nil {
		return fmt.Errorf("unable to decode config into struct: %w", err)
	}

	return validate()
}

// validate validates the configuration
func validate() error {
	// Validate required fields
	if c.App.Name == "" {
		return fmt.Errorf("app.name is required")
	}

	if c.Database.Name == "" {
		return fmt.Errorf("database.name is required")
	}

	if c.Database.Host == "" {
		return fmt.Errorf("database.host is required")
	}

	// Validate environment
	validEnvs := []string{"local", "development", "staging", "production"}
	isValidEnv := false
	for _, env := range validEnvs {
		if c.App.Env == env {
			isValidEnv = true
			break
		}
	}
	if !isValidEnv {
		return fmt.Errorf("app.env must be one of: %s", strings.Join(validEnvs, ", "))
	}

	// Validate port ranges
	if c.Server.Port < 1 || c.Server.Port > 65535 {
		return fmt.Errorf("server.port must be between 1 and 65535")
	}

	if c.Database.Port < 1 || c.Database.Port > 65535 {
		return fmt.Errorf("database.port must be between 1 and 65535")
	}

	return nil
}

// Get returns the loaded configuration
func Get() *Config {
	return &c
}

// GetApp returns app configuration
func GetApp() AppConfig {
	return c.App
}

// GetServer returns server configuration
func GetServer() ServerConfig {
	return c.Server
}

// GetRedis returns redis configuration
func GetRedis() RedisConfig {
	return c.Redis
}

// GetJWT returns JWT configuration
func GetJWT() JWTConfig {
	return c.JWT
}

// GetDatabase returns database configuration
func GetDatabase() database.Config {
	return c.Database
}

// GetLogDatabase returns log database configuration
func GetLogDatabase() database.Config {
	return c.LogDatabase
}

// IsDevelopment returns true if running in development mode
func IsDevelopment() bool {
	return c.App.Env == "development" || c.App.Env == "local"
}

// IsProduction returns true if running in production mode
func IsProduction() bool {
	return c.App.Env == "production"
}

// PrintConfig prints the current configuration (without sensitive data)
func PrintConfig() {
	log.Printf("=== Configuration ===")
	log.Printf("App: %s (env: %s)", c.App.Name, c.App.Env)
	log.Printf("Server: Port %d", c.Server.Port)
	log.Printf("Redis: %s:%s (db: %d)", c.Redis.Address, c.Redis.Port, c.Redis.DB)
	log.Printf("Database: %s@%s:%d/%s", c.Database.User, c.Database.Host, c.Database.Port, c.Database.Name)
	log.Printf("Log DB: %s@%s:%d/%s", c.LogDatabase.User, c.LogDatabase.Host, c.LogDatabase.Port, c.LogDatabase.Name)
	log.Printf("====================")
}
