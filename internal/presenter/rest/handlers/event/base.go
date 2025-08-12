package event

import (
	"rumi-go/internal/usecase/event"
)

type EventHandler struct {
	eventUsecase *event.EventUsecase
}

func NewEventHandler(eventUsecase *event.EventUsecase) *EventHandler {
	return &EventHandler{
		eventUsecase: eventUsecase,
	}
}
