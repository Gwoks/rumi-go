package event_report

import (
	"rumi-go/internal/usecase/event_report"
)

type EventReportHandler struct {
	eventReportUsecase *event_report.EventReportUsecase
}

func NewEventReportHandler(eventReportUsecase *event_report.EventReportUsecase) *EventReportHandler {
	return &EventReportHandler{
		eventReportUsecase: eventReportUsecase,
	}
}
