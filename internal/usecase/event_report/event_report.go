package event_report

import (
	"time"

	"rumi-go/internal/infrastructure/database"
	"rumi-go/internal/model"
)

type EventReportUsecase struct {
	location *time.Location
	database database.Store
}

func NewEventReportUsecase(database database.Store) *EventReportUsecase {
	loc, _ := time.LoadLocation("Asia/Jakarta")

	return &EventReportUsecase{
		location: loc,
		database: database,
	}
}

// toEventReportResponse converts EventReport model to EventReportResponse
func (u *EventReportUsecase) toEventReportResponse(report *model.EventReport) *model.EventReportResponse {
	return &model.EventReportResponse{
		ID:        report.ID,
		UserID:    report.UserID,
		TeacherID: report.TeacherID,
		EventID:   report.EventID,
		Notes:     report.Notes,
		Score:     report.Score,
		CreatedAt: report.CreatedAt,
		UpdatedAt: report.UpdatedAt,
	}
}
