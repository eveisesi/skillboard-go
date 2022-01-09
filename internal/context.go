package internal

import (
	"context"

	"github.com/eveisesi/skillz"
	"github.com/gofrs/uuid"
)

type ctxKey int

const (
	ctxSessionID ctxKey = iota
	ctxUser
	ctxIsAuthed
)

func ContextWithSessionID(ctx context.Context, sessionID uuid.UUID) context.Context {
	return context.WithValue(ctx, ctxSessionID, sessionID)
}

func SessionIDFromContext(ctx context.Context) uuid.UUID {
	if sessionID, ok := ctx.Value(ctxSessionID).(uuid.UUID); ok {
		return sessionID
	}
	return uuid.Nil
}

func ContextWithUser(ctx context.Context, user *skillz.User) context.Context {

	ctx = context.WithValue(ctx, ctxUser, user)
	ctx = context.WithValue(ctx, ctxIsAuthed, true)

	return ctx
}

func IsAuthedFromContext(ctx context.Context) bool {

	_, ok := ctx.Value(ctxIsAuthed).(bool)

	return ok

}

func UserFromContext(ctx context.Context) *skillz.User {
	if user, ok := ctx.Value(ctxUser).(*skillz.User); ok {
		return user
	}

	return nil
}
