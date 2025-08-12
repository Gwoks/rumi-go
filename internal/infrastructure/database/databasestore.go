package database

import (
	"context"
	"fmt"
	"log/slog"

	"rumi-go/internal/database"
	"rumi-go/internal/model"
	"rumi-go/internal/utils"

	"github.com/jmoiron/sqlx"
)

type Store interface {
	Close() error
	DB() *sqlx.DB
	UserStore() model.UserSQLStore
	EventStore() model.EventSQLStore
	EventReportStore() model.EventReportSQLStore
}

type SQLStore struct {
	db          *sqlx.DB
	user        model.UserSQLStore
	event       model.EventSQLStore
	eventReport model.EventReportSQLStore
}

func NewStore(ctx context.Context, dbConfig database.Config) (Store, error) {
	slog.InfoContext(ctx, fmt.Sprintf("monolithsqlstore database config, max open %d, max idle %d, max lifetime %d",
		dbConfig.MaxOpen, dbConfig.MaxIdle, dbConfig.MaxLifetime))

	db, err := utils.InitSqlxDB(dbConfig)
	if err != nil {
		return nil, err
	}

	return &SQLStore{
		db:          db,
		user:        NewUserStore(db, dbConfig.StatementTimeout),
		event:       NewEventStore(db, dbConfig.StatementTimeout),
		eventReport: NewEventReportStore(db, dbConfig.StatementTimeout),
	}, nil
}

func (s *SQLStore) Close() error {
	return s.db.Close()
}

func (s *SQLStore) UserStore() model.UserSQLStore {
	return s.user
}

func (s *SQLStore) EventStore() model.EventSQLStore {
	return s.event
}

func (s *SQLStore) EventReportStore() model.EventReportSQLStore {
	return s.eventReport
}

func (s *SQLStore) DB() *sqlx.DB {
	return s.db
}
