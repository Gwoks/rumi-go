package event

import (
	"context"
)

// DeleteEvent deletes an event by ID
func (u *EventUsecase) DeleteEvent(ctx context.Context, eventID int64) error {
	return u.database.EventStore().Delete(ctx, eventID)
}
