package main

import (
	"log"
	"os"

	"rumi-go/internal/config"
	"rumi-go/internal/migration"
	"rumi-go/internal/presenter/console"

	_ "github.com/go-sql-driver/mysql" // MySQL driver
	"gopkg.in/ukautz/clif.v1"
)

func main() {
	if len(os.Args) == 1 {
		return
	}

	// Load configuration first
	err := config.Load(".")
	if err != nil {
		log.Fatalln("Failed to load configuration:", err)
	}

	migrateCommand := migration.NewDBMigrations(config.Get().Database, "migrations/")

	cli := clif.New("rumi-go", "0.0.1", "")
	cmd := console.Console{}

	cli.Add(cmd.StartServer())
	cli.Add(cmd.StartMigration(migrateCommand))
	cli.Add(cmd.StartMigrationCreate(migrateCommand))
	cli.Add(cmd.StartMigrationUp(migrateCommand))
	cli.Add(cmd.StartMigrationDown(migrateCommand))

	// Add new MySQL migration commands
	cli.Add(cmd.StartMySQLMigrationUp())
	cli.Add(cmd.StartMySQLMigrationDown())
	cli.Add(cmd.StartMySQLMigrationStatus())

	// Add configuration command
	cli.Add(cmd.StartConfigShow())

	cli.Run()
}
