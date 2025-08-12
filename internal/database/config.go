package database

import (
	"fmt"
	"time"
)

// Config holds database configuration parameters
type Config struct {
	Driver           string        `mapstructure:"driver"`
	Name             string        `mapstructure:"name"`
	Host             string        `mapstructure:"host"`
	Port             int           `mapstructure:"port"`
	User             string        `mapstructure:"user"`
	Password         string        `mapstructure:"password"`
	MaxOpen          int           `mapstructure:"max_open"`
	MaxIdle          int           `mapstructure:"max_idle"`
	MaxLifetime      time.Duration `mapstructure:"max_life_time"`
	MaxIdleTime      time.Duration `mapstructure:"max_idle_time"`
	StatementTimeout time.Duration `mapstructure:"statement_timeout"`
}

// DSN returns the database connection string for MySQL
func (c Config) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.User, c.Password, c.Host, c.Port, c.Name)
}
