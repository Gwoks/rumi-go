package event

import (
	"context"
	"fmt"
	"time"

	"rumi-go/internal/model"
)

// CreateEvent creates a new event
func (u *EventUsecase) CreateEvent(ctx context.Context, req model.EventRequest) (*model.EventResponse, error) {
	// Create event
	event := &model.Event{
		Event:       req.Event,
		Description: req.Description,
		Date:        req.Date,
		CreatedAt:   time.Now().In(u.location),
		UpdatedAt:   time.Now().In(u.location),
	}

	if err := u.database.EventStore().Create(ctx, event); err != nil {
		return nil, fmt.Errorf("failed to create event: %w", err)
	}

	return u.toEventResponse(event), nil
}
