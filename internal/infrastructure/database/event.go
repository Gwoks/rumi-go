package database

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"rumi-go/internal/model"

	"github.com/jmoiron/sqlx"
)

type EventStore struct {
	db      *sqlx.DB
	timeout time.Duration
}

func NewEventStore(db *sqlx.DB, timeout time.Duration) *EventStore {
	return &EventStore{
		db:      db,
		timeout: timeout,
	}
}

var (
	getEventQuery    = `SELECT * FROM events WHERE id = ?`
	createEventQuery = `INSERT INTO events (event_name, event_description, event_date, created_at, updated_at) VALUES (?, ?, ?, ?, ?)`
	updateEventQuery = `UPDATE events SET event_name = ?, event_description = ?, event_date = ?, updated_at = ? WHERE id = ?`
	deleteEventQuery = `DELETE FROM events WHERE id = ?`
	listEventsQuery  = `SELECT * FROM events ORDER BY created_at DESC LIMIT ? OFFSET ?`
	countEventsQuery = `SELECT COUNT(*) FROM events`
)

func (s *EventStore) Get(ctx context.Context, eventID int64) (*model.Event, error) {
	ctxWithDeadline, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	var event model.Event
	err := s.db.GetContext(ctxWithDeadline, &event, getEventQuery, eventID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("event not found")
		}
		return nil, err
	}
	return &event, nil
}

func (s *EventStore) Create(ctx context.Context, event *model.Event) error {
	ctxWithDeadline, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	result, err := s.db.ExecContext(ctxWithDeadline, createEventQuery,
		event.Event, event.Description, event.Date, event.CreatedAt, event.UpdatedAt)
	if err != nil {
		return err
	}

	// Get the inserted ID
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	event.ID = id
	return nil
}

func (s *EventStore) Update(ctx context.Context, event *model.Event) error {
	ctxWithDeadline, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	result, err := s.db.ExecContext(ctxWithDeadline, updateEventQuery,
		event.Event, event.Description, event.Date, event.UpdatedAt, event.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("event not found")
	}

	return nil
}

func (s *EventStore) Delete(ctx context.Context, eventID int64) error {
	ctxWithDeadline, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	result, err := s.db.ExecContext(ctxWithDeadline, deleteEventQuery, eventID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("event not found")
	}

	return nil
}

func (s *EventStore) List(ctx context.Context, limit, offset int) ([]*model.Event, error) {
	ctxWithDeadline, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	var events []*model.Event
	err := s.db.SelectContext(ctxWithDeadline, &events, listEventsQuery, limit, offset)
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (s *EventStore) Count(ctx context.Context) (int64, error) {
	ctxWithDeadline, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	var count int64
	err := s.db.GetContext(ctxWithDeadline, &count, countEventsQuery)
	return count, err
}
