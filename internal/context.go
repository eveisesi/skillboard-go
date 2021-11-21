package internal

import (
	"context"

	"github.com/eveisesi/skillz"
	"github.com/lestrrat-go/jwx/jwt"
)

type ctxKey int

const (
	ctxToken ctxKey = iota
	ctxUser
	ctxUpdateable
)

func UpdateableFromContext(ctx context.Context) bool {
	if updateable, ok := ctx.Value(ctxUpdateable).(bool); ok {
		return updateable
	}

	return false
}

func ContextWithUpdateable(ctx context.Context, updateable bool) context.Context {
	return context.WithValue(ctx, ctxUpdateable, updateable)
}

func ContextWithToken(ctx context.Context, token jwt.Token) context.Context {
	return context.WithValue(ctx, ctxToken, token)
}

func TokenFromContext(ctx context.Context) jwt.Token {
	if token, ok := ctx.Value(ctxToken).(jwt.Token); ok {
		return token
	}
	return nil
}

func ContextWithUser(ctx context.Context, user *skillz.User) context.Context {
	return context.WithValue(ctx, ctxUser, user)
}

func UserFromContext(ctx context.Context) *skillz.User {
	if user, ok := ctx.Value(ctxUser).(*skillz.User); ok {
		return user
	}

	return nil
}
