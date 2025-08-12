package utils

import (
	"errors"
	"fmt"
	"net/url"

	"rumi-go/internal/database"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	sqltrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/database/sql"
	sqlxtrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/jmoiron/sqlx"
)

func InitSqlxDB(dbConfig database.Config) (*sqlx.DB, error) {
	driverName := dbConfig.Driver
	spanServiceName := fmt.Sprintf("%s.%s", driverName, dbConfig.Host)
	switch driverName {
	case "postgres":
		sqltrace.Register(driverName, &pq.Driver{},
			sqltrace.WithServiceName(spanServiceName),
			sqltrace.WithCustomTag("db.name", dbConfig.Name),
		)
		dataSourceName := fmt.Sprintf("%s://%s:%s@%s:%d/%s?sslmode=disable",
			driverName,
			dbConfig.User,
			url.QueryEscape(dbConfig.Password),
			dbConfig.Host,
			dbConfig.Port,
			dbConfig.Name)
		db, err := sqlxtrace.Connect("postgres", dataSourceName)
		if err != nil {
			return nil, err
		}
		db.SetMaxOpenConns(dbConfig.MaxOpen)
		db.SetMaxIdleConns(dbConfig.MaxIdle)
		db.SetConnMaxLifetime(dbConfig.MaxLifetime)
		db.SetConnMaxIdleTime(dbConfig.MaxIdleTime)
		return db, nil
	case "mysql":
		sqltrace.Register(driverName, &mysql.MySQLDriver{},
			sqltrace.WithServiceName(spanServiceName),
			sqltrace.WithCustomTag("db.name", dbConfig.Name),
		)
		dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
			dbConfig.User,
			dbConfig.Password,
			dbConfig.Host,
			dbConfig.Port,
			dbConfig.Name)
		db, err := sqlxtrace.Connect("mysql", dataSourceName)
		if err != nil {
			return nil, err
		}
		db.SetMaxOpenConns(dbConfig.MaxOpen)
		db.SetMaxIdleConns(dbConfig.MaxIdle)
		db.SetConnMaxLifetime(dbConfig.MaxLifetime)
		db.SetConnMaxIdleTime(dbConfig.MaxIdleTime)
		return db, nil
	default:
		return nil, errors.New("Unsupported driver: " + driverName)
	}
}
