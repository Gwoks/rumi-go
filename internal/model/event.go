package model

import (
	"context"
	"time"
)

type Event struct {
	ID          int64     `db:"id"`
	Event       string    `db:"event_name"`
	Description string    `db:"event_description"`
	Date        string    `db:"event_date"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

type EventRequest struct {
	Event       string `json:"event_name" validate:"required"`
	Description string `json:"event_description"`
	Date        string `json:"event_date"`
}

type EventResponse struct {
	ID          int64     `json:"id"`
	Event       string    `json:"event_name"`
	Description string    `json:"event_description"`
	Date        string    `json:"event_date"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type EventSQLStore interface {
	Get(ctx context.Context, eventID int64) (*Event, error)
	Create(ctx context.Context, event *Event) error
	Update(ctx context.Context, event *Event) error
	Delete(ctx context.Context, eventID int64) error
	List(ctx context.Context, limit, offset int) ([]*Event, error)
	Count(ctx context.Context) (int64, error)
}
