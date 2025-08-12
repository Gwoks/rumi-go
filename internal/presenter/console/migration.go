package console

import (
	"fmt"

	"rumi-go/internal/config"
	"rumi-go/internal/migration"

	"gopkg.in/ukautz/clif.v1"
)

const (
	ERR_NO_CHANGE = "no change"
)

func (c *Console) StartMigration(migrate migration.Migrations) *clif.Command {
	return clif.NewCommand("migrate", "migrate the database", func(o *clif.Command, in clif.Input, out clif.Output) error {
		out.Printf("Migration commands available:\n")
		out.Printf("Use the following commands:\n")
		out.Printf("  migrate:create --filename <name>  - create a new migration file\n")
		out.Printf("  migrate:run                        - run pending migrations\n")
		out.Printf("  migrate:rollback                      - rollback the last migration\n")
		out.Printf("  migrate:mysql:up                   - run MySQL migrations\n")
		out.Printf("  migrate:mysql:down                 - rollback MySQL migrations\n")
		out.Printf("  migrate:mysql:status               - check MySQL migration status\n")
		return nil
	})
}

func (c *Console) StartMigrationCreate(migrate migration.Migrations) *clif.Command {
	cmd := clif.NewCommand("migrate:create", "create the migration file", func(o *clif.Command, in clif.Input, out clif.Output) error {
		filename := o.Option("filename").String()
		if filename == "" {
			return fmt.Errorf("filename is required. Use --filename <name>")
		}
		if err := migrate.Create(filename); err != nil {
			return fmt.Errorf("failed to create migration file: %w", err)
		}
		out.Printf("Created migration file: %s\n", filename)
		return nil
	})
	cmd.AddOption(clif.NewOption("filename", "f", "Name of migration file", "", true, false))
	return cmd
}

func (c *Console) StartMigrationUp(migrate migration.Migrations) *clif.Command {
	return clif.NewCommand("migrate:run", "run the migration files", func(o *clif.Command, in clif.Input, out clif.Output) error {
		if err := migrate.Up(); err != nil && err.Error() != ERR_NO_CHANGE {
			return fmt.Errorf("migration up failed: %w", err)
		}
		out.Printf("Migrations applied successfully\n")
		return nil
	})
}

func (c *Console) StartMigrationDown(migrate migration.Migrations) *clif.Command {
	return clif.NewCommand("migrate:rollback", "rollback the migration", func(o *clif.Command, in clif.Input, out clif.Output) error {
		if err := migrate.Down(); err != nil && err.Error() != ERR_NO_CHANGE {
			return fmt.Errorf("migration down failed: %w", err)
		}
		out.Printf("Migrations rolled back successfully\n")
		return nil
	})
}

// New MySQL migration commands
func (c *Console) StartMySQLMigrationUp() *clif.Command {
	return clif.NewCommand("migrate:mysql:up", "run MySQL migrations", func(o *clif.Command, in clif.Input, out clif.Output) error {
		// Load configuration
		if err := config.Load("."); err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}
		conf := config.Get()

		// Run migrations
		if err := migration.RunMigrations(conf.Database, "./migrations"); err != nil {
			return fmt.Errorf("MySQL migration failed: %w", err)
		}

		out.Printf("MySQL migrations completed successfully\n")
		return nil
	})
}

func (c *Console) StartMySQLMigrationDown() *clif.Command {
	return clif.NewCommand("migrate:mysql:down", "rollback MySQL migrations", func(o *clif.Command, in clif.Input, out clif.Output) error {
		// Load configuration
		if err := config.Load("."); err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}
		conf := config.Get()

		// Rollback migrations
		if err := migration.RollbackMigrations(conf.Database, "./migrations"); err != nil {
			return fmt.Errorf("MySQL migration rollback failed: %w", err)
		}

		out.Printf("MySQL migrations rolled back successfully\n")
		return nil
	})
}

func (c *Console) StartMySQLMigrationStatus() *clif.Command {
	return clif.NewCommand("migrate:mysql:status", "check MySQL migration status", func(o *clif.Command, in clif.Input, out clif.Output) error {
		out.Printf("MySQL migration status:\n")
		out.Printf("  - Migrations directory: ./migrations\n")
		out.Printf("  - Available migrations:\n")
		out.Printf("    * 001_create_users_table.up.sql - Create users table\n")
		out.Printf("    * 002_insert_admin_user.up.sql - Insert admin user\n")
		out.Printf("  - To run migrations: make migrate-mysql-up\n")
		out.Printf("  - To rollback: make migrate-mysql-down\n")
		return nil
	})
}
