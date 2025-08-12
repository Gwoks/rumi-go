package database

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"rumi-go/internal/model"

	"github.com/jmoiron/sqlx"
)

type EventReportStore struct {
	db      *sqlx.DB
	timeout time.Duration
}

func NewEventReportStore(db *sqlx.DB, timeout time.Duration) *EventReportStore {
	return &EventReportStore{
		db:      db,
		timeout: timeout,
	}
}

var (
	getEventReportQuery        = `SELECT * FROM event_report WHERE id = ?`
	createEventReportQuery     = `INSERT INTO event_report (user_id, teacher_id, event_id, notes, score, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)`
	updateEventReportQuery     = `UPDATE event_report SET user_id = ?, teacher_id = ?, event_id = ?, notes = ?, score = ?, updated_at = ? WHERE id = ?`
	deleteEventReportQuery     = `DELETE FROM event_report WHERE id = ?`
	listEventReportsQuery      = `SELECT * FROM event_report ORDER BY created_at DESC LIMIT ? OFFSET ?`
	countEventReportsQuery     = `SELECT COUNT(*) FROM event_report`
	getByEventIDQuery          = `SELECT * FROM event_report WHERE event_id = ? ORDER BY created_at DESC`
	getByUserIDQuery           = `SELECT * FROM event_report WHERE user_id = ? ORDER BY created_at DESC`
)

func (s *EventReportStore) Get(ctx context.Context, reportID int64) (*model.EventReport, error) {
	ctxWithDeadline, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	var report model.EventReport
	err := s.db.GetContext(ctxWithDeadline, &report, getEventReportQuery, reportID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("event report not found")
		}
		return nil, err
	}
	return &report, nil
}

func (s *EventReportStore) Create(ctx context.Context, report *model.EventReport) error {
	ctxWithDeadline, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	result, err := s.db.ExecContext(ctxWithDeadline, createEventReportQuery,
		report.UserID, report.TeacherID, report.EventID, report.Notes, report.Score, report.CreatedAt, report.UpdatedAt)
	if err != nil {
		return err
	}

	// Get the inserted ID
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	report.ID = id
	return nil
}

func (s *EventReportStore) Update(ctx context.Context, report *model.EventReport) error {
	ctxWithDeadline, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	result, err := s.db.ExecContext(ctxWithDeadline, updateEventReportQuery,
		report.UserID, report.TeacherID, report.EventID, report.Notes, report.Score, report.UpdatedAt, report.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("event report not found")
	}

	return nil
}

func (s *EventReportStore) Delete(ctx context.Context, reportID int64) error {
	ctxWithDeadline, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	result, err := s.db.ExecContext(ctxWithDeadline, deleteEventReportQuery, reportID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("event report not found")
	}

	return nil
}

func (s *EventReportStore) List(ctx context.Context, limit, offset int) ([]*model.EventReport, error) {
	ctxWithDeadline, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	var reports []*model.EventReport
	err := s.db.SelectContext(ctxWithDeadline, &reports, listEventReportsQuery, limit, offset)
	if err != nil {
		return nil, err
	}

	return reports, nil
}

func (s *EventReportStore) Count(ctx context.Context) (int64, error) {
	ctxWithDeadline, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	var count int64
	err := s.db.GetContext(ctxWithDeadline, &count, countEventReportsQuery)
	return count, err
}

func (s *EventReportStore) GetByEventID(ctx context.Context, eventID int64) ([]*model.EventReport, error) {
	ctxWithDeadline, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	var reports []*model.EventReport
	err := s.db.SelectContext(ctxWithDeadline, &reports, getByEventIDQuery, eventID)
	if err != nil {
		return nil, err
	}

	return reports, nil
}

func (s *EventReportStore) GetByUserID(ctx context.Context, userID int64) ([]*model.EventReport, error) {
	ctxWithDeadline, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	var reports []*model.EventReport
	err := s.db.SelectContext(ctxWithDeadline, &reports, getByUserIDQuery, userID)
	if err != nil {
		return nil, err
	}

	return reports, nil
}
