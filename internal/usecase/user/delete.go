package user

import (
	"context"
)

// DeleteUser deletes a user
func (u *UserUsecase) DeleteUser(ctx context.Context, userID int64) error {
	return u.database.UserStore().Delete(ctx, userID)
}
