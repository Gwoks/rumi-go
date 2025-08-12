package model

import (
	"context"
	"encoding/json"
	"time"
)

type EventReport struct {
	ID        int64           `db:"id"`
	UserID    int64           `db:"user_id"`
	TeacherID int64           `db:"teacher_id"`
	EventID   int64           `db:"event_id"`
	Notes     string          `db:"notes"`
	Score     json.RawMessage `db:"score"`
	CreatedAt time.Time       `db:"created_at"`
	UpdatedAt time.Time       `db:"updated_at"`
}

type EventReportRequest struct {
	UserID    int64           `json:"user_id" validate:"required"`
	TeacherID int64           `json:"teacher_id" validate:"required"`
	EventID   int64           `json:"event_id" validate:"required"`
	Notes     string          `json:"notes"`
	Score     json.RawMessage `json:"score"`
}

type EventReportResponse struct {
	ID        int64           `json:"id"`
	UserID    int64           `json:"user_id"`
	TeacherID int64           `json:"teacher_id"`
	EventID   int64           `json:"event_id"`
	Notes     string          `json:"notes"`
	Score     json.RawMessage `json:"score"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

type EventReportSQLStore interface {
	Get(ctx context.Context, reportID int64) (*EventReport, error)
	Create(ctx context.Context, report *EventReport) error
	Update(ctx context.Context, report *EventReport) error
	Delete(ctx context.Context, reportID int64) error
	List(ctx context.Context, limit, offset int) ([]*EventReport, error)
	Count(ctx context.Context) (int64, error)
	GetByEventID(ctx context.Context, eventID int64) ([]*EventReport, error)
	GetByUserID(ctx context.Context, userID int64) ([]*EventReport, error)
}
