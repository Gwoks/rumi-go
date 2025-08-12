package event_report

import (
	"context"
)

// DeleteEventReport deletes an event report by ID
func (u *EventReportUsecase) DeleteEventReport(ctx context.Context, reportID int64) error {
	return u.database.EventReportStore().Delete(ctx, reportID)
}
