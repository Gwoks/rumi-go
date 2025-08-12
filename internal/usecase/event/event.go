package event

import (
	"time"

	"rumi-go/internal/infrastructure/database"
	"rumi-go/internal/model"
)

type EventUsecase struct {
	location *time.Location
	database database.Store
}

func NewEventUsecase(database database.Store) *EventUsecase {
	loc, _ := time.LoadLocation("Asia/Jakarta")

	return &EventUsecase{
		location: loc,
		database: database,
	}
}

// toEventResponse converts Event model to EventResponse
func (u *EventUsecase) toEventResponse(event *model.Event) *model.EventResponse {
	return &model.EventResponse{
		ID:          event.ID,
		Event:       event.Event,
		Description: event.Description,
		Date:        event.Date,
		CreatedAt:   event.CreatedAt,
		UpdatedAt:   event.UpdatedAt,
	}
}
