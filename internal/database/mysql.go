package database

import (
	"context"
	"fmt"
	"time"

	"github.com/afghan/rumi-backend/internal/config"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// DB wraps sqlx.DB with additional functionality
type DB struct {
	*sqlx.DB
}

// NewMySQLConnection creates a new MySQL database connection
func NewMySQLConnection(cfg *config.Config) (*DB, error) {
	dsn := cfg.GetDSN()

	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(1 * time.Minute)

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &DB{DB: db}, nil
}

// Close closes the database connection
func (d *DB) Close() error {
	return d.DB.Close()
}

// BeginTx starts a new transaction with context
func (d *DB) BeginTx(ctx context.Context) (*sqlx.Tx, error) {
	return d.DB.BeginTxx(ctx, nil)
}

// WithTransaction executes a function within a transaction
func (d *DB) WithTransaction(ctx context.Context, fn func(*sqlx.Tx) error) error {
	tx, err := d.BeginTx(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	if err := fn(tx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("transaction failed: %v, rollback failed: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}
