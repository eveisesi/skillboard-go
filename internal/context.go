package internal

import (
	"context"

	"github.com/lestrrat-go/jwx/jwt"
)

type ctxKey int

const (
	ctxToken ctxKey = iota
)

func ContextWithToken(ctx context.Context, token jwt.Token) context.Context {
	return context.WithValue(ctx, ctxToken, token)
}

func TokenFromContext(ctx context.Context) jwt.Token {
	if token, ok := ctx.Value(ctxToken).(jwt.Token); ok {
		return token
	}
	return nil
}
