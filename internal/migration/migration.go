package migration

import (
	"fmt"
	"os"
	"time"

	"rumi-go/internal/database"
	"rumi-go/internal/utils"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/pkg/errors"
)

type Migrations interface {
	Up() error
	Down() error
	Create(title string) error
}

type Migrator interface {
	Up() error
	Steps(n int) error
}

func NewDBMigrations(config database.Config, sourceFile string) Migrations {
	return &DBMigrations{
		sourceFile: sourceFile,
		config:     config,
	}
}

type DBMigrations struct {
	migrate    Migrator
	sourceFile string
	config     database.Config
}

func (m *DBMigrations) init() error {
	if m.migrate != nil {
		return nil
	}

	sqlconn, err := utils.InitSqlxDB(m.config)
	if err != nil {
		return err
	}

	defer sqlconn.Close()

	sourceFile := fmt.Sprintf("file://%s", m.sourceFile)
	driver, err := mysql.WithInstance(sqlconn.DB, &mysql.Config{})

	if err != nil {
		return err
	}

	mi, err := migrate.NewWithDatabaseInstance(sourceFile, m.config.Name, driver)
	if err != nil {
		return err
	}

	m.migrate = mi

	return nil
}

func (m *DBMigrations) Up() error {
	if err := m.init(); err != nil {
		return err
	}
	err := m.migrate.Up()
	if err != nil {
		return err
	}

	return nil
}

func (m *DBMigrations) Down() error {
	if err := m.init(); err != nil {
		return err
	}

	err := m.migrate.Steps(-1)
	if err != nil {
		return err
	}

	return nil
}

func (m *DBMigrations) Create(title string) error {
	if title == "" {
		return errors.New("Title can't be empty")
	}
	fileNameUp, fileNameDown := m.generateFileName(title)

	if _, err := os.Create(fileNameUp); err != nil {
		return err
	}

	if _, err := os.Create(fileNameDown); err != nil {
		_ = os.Remove(fileNameUp)
		return err
	}

	return nil
}

func (m *DBMigrations) generateFileName(title string) (fileNameUp string, fileNameDown string) {
	now := time.Now()
	unixTime := now.Unix()

	fileNameUp = fmt.Sprintf("%s/%d_%s.up.sql", m.sourceFile, unixTime, title)
	fileNameDown = fmt.Sprintf("%s/%d_%s.down.sql", m.sourceFile, unixTime, title)

	return fileNameUp, fileNameDown
}
