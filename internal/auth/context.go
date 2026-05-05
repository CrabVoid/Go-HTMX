package auth

import (
	"context"

	"github.com/google/uuid"
)

type contextKey string

const UserIDKey contextKey = "userID"

func GetUserID(ctx context.Context) uuid.UUID {
	userID, ok := ctx.Value(UserIDKey).(uuid.UUID)
	if !ok {
		return uuid.Nil
	}
	return userID
}
