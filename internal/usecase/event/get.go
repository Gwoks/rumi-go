package event

import (
	"context"

	"rumi-go/internal/model"
)

// GetEvent retrieves an event by ID
func (u *EventUsecase) GetEvent(ctx context.Context, eventID int64) (*model.EventResponse, error) {
	event, err := u.database.EventStore().Get(ctx, eventID)
	if err != nil {
		return nil, err
	}

	return u.toEventResponse(event), nil
}

// ListEvents retrieves a list of events with pagination
func (u *EventUsecase) ListEvents(ctx context.Context, limit, offset int) ([]*model.EventResponse, error) {
	events, err := u.database.EventStore().List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	responses := make([]*model.EventResponse, len(events))
	for i, event := range events {
		responses[i] = u.toEventResponse(event)
	}

	return responses, nil
}
