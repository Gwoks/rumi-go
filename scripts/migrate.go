package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/afghan/rumi-backend/internal/config"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to database
	db, err := sqlx.Connect("mysql", cfg.GetDSN())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Run migrations
	if err := runMigrations(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	fmt.Println("✅ Migrations completed successfully!")
}

func runMigrations(db *sqlx.DB) error {
	// Create migrations table if it doesn't exist
	createMigrationsTable := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version VARCHAR(255) PRIMARY KEY,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`

	if _, err := db.Exec(createMigrationsTable); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Get list of migration files
	migrationFiles, err := filepath.Glob("migrations/*.up.sql")
	if err != nil {
		return fmt.Errorf("failed to read migration files: %w", err)
	}

	for _, file := range migrationFiles {
		version := filepath.Base(file)
		version = version[:len(version)-7] // Remove ".up.sql"

		// Check if migration already applied
		var count int
		err := db.Get(&count, "SELECT COUNT(*) FROM schema_migrations WHERE version = ?", version)
		if err != nil {
			return fmt.Errorf("failed to check migration status: %w", err)
		}

		if count > 0 {
			fmt.Printf("⏭️  Skipping %s (already applied)\n", version)
			continue
		}

		// Read and execute migration
		content, err := ioutil.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", file, err)
		}

		if _, err := db.Exec(string(content)); err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", version, err)
		}

		// Record migration as applied
		if _, err := db.Exec("INSERT INTO schema_migrations (version) VALUES (?)", version); err != nil {
			return fmt.Errorf("failed to record migration %s: %w", version, err)
		}

		fmt.Printf("✅ Applied migration: %s\n", version)
	}

	return nil
}
