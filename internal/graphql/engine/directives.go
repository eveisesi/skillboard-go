package resolvers

import (
	"context"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
	"github.com/eveisesi/skillz/internal"
)

func IsAuthed(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {

	user := internal.UserFromContext(ctx)
	if user == nil {
		return nil, fmt.Errorf("invalidate authentication detected, please supply a valid token")
	}

	return next(ctx)

}
