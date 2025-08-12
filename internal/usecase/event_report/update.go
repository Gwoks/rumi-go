package event_report

import (
	"context"
	"fmt"
	"time"

	"rumi-go/internal/model"
)

// UpdateEventReport updates an existing event report
func (u *EventReportUsecase) UpdateEventReport(ctx context.Context, reportID int64, req model.EventReportRequest) (*model.EventReportResponse, error) {
	// Get existing event report
	existingReport, err := u.database.EventReportStore().Get(ctx, reportID)
	if err != nil {
		return nil, err
	}

	// Update fields
	existingReport.UserID = req.UserID
	existingReport.TeacherID = req.TeacherID
	existingReport.EventID = req.EventID
	existingReport.Notes = req.Notes
	existingReport.Score = req.Score
	existingReport.UpdatedAt = time.Now().In(u.location)

	if err := u.database.EventReportStore().Update(ctx, existingReport); err != nil {
		return nil, fmt.Errorf("failed to update event report: %w", err)
	}

	return u.toEventReportResponse(existingReport), nil
}
