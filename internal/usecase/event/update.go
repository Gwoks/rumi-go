package event

import (
	"context"
	"fmt"
	"time"

	"rumi-go/internal/model"
)

// UpdateEvent updates an existing event
func (u *EventUsecase) UpdateEvent(ctx context.Context, eventID int64, req model.EventRequest) (*model.EventResponse, error) {
	// Get existing event
	existingEvent, err := u.database.EventStore().Get(ctx, eventID)
	if err != nil {
		return nil, err
	}

	// Update fields
	existingEvent.Event = req.Event
	existingEvent.Description = req.Description
	existingEvent.Date = req.Date
	existingEvent.UpdatedAt = time.Now().In(u.location)

	if err := u.database.EventStore().Update(ctx, existingEvent); err != nil {
		return nil, fmt.Errorf("failed to update event: %w", err)
	}

	return u.toEventResponse(existingEvent), nil
}
