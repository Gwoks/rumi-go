package event_report

import (
	"context"
	"fmt"
	"time"

	"rumi-go/internal/model"
)

// CreateEventReport creates a new event report
func (u *EventReportUsecase) CreateEventReport(ctx context.Context, req model.EventReportRequest) (*model.EventReportResponse, error) {
	// Create event report
	report := &model.EventReport{
		UserID:    req.UserID,
		TeacherID: req.TeacherID,
		EventID:   req.EventID,
		Notes:     req.Notes,
		Score:     req.Score,
		CreatedAt: time.Now().In(u.location),
		UpdatedAt: time.Now().In(u.location),
	}

	if err := u.database.EventReportStore().Create(ctx, report); err != nil {
		return nil, fmt.Errorf("failed to create event report: %w", err)
	}

	return u.toEventReportResponse(report), nil
}
