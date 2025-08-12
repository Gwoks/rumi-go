package event_report

import (
	"context"

	"rumi-go/internal/model"
)

// GetEventReport retrieves an event report by ID
func (u *EventReportUsecase) GetEventReport(ctx context.Context, reportID int64) (*model.EventReportResponse, error) {
	report, err := u.database.EventReportStore().Get(ctx, reportID)
	if err != nil {
		return nil, err
	}

	return u.toEventReportResponse(report), nil
}

// ListEventReports retrieves a list of event reports with pagination
func (u *EventReportUsecase) ListEventReports(ctx context.Context, limit, offset int) ([]*model.EventReportResponse, error) {
	reports, err := u.database.EventReportStore().List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	responses := make([]*model.EventReportResponse, len(reports))
	for i, report := range reports {
		responses[i] = u.toEventReportResponse(report)
	}

	return responses, nil
}

// GetEventReportsByEventID retrieves event reports by event ID
func (u *EventReportUsecase) GetEventReportsByEventID(ctx context.Context, eventID int64) ([]*model.EventReportResponse, error) {
	reports, err := u.database.EventReportStore().GetByEventID(ctx, eventID)
	if err != nil {
		return nil, err
	}

	responses := make([]*model.EventReportResponse, len(reports))
	for i, report := range reports {
		responses[i] = u.toEventReportResponse(report)
	}

	return responses, nil
}

// GetEventReportsByUserID retrieves event reports by user ID
func (u *EventReportUsecase) GetEventReportsByUserID(ctx context.Context, userID int64) ([]*model.EventReportResponse, error) {
	reports, err := u.database.EventReportStore().GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	responses := make([]*model.EventReportResponse, len(reports))
	for i, report := range reports {
		responses[i] = u.toEventReportResponse(report)
	}

	return responses, nil
}
