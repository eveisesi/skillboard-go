package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/eveisesi/skillz"
	"github.com/eveisesi/skillz/internal"
	"github.com/eveisesi/skillz/internal/server/gql/generated"
)

func (r *queryResolver) InitializeAuth(ctx context.Context) (string, error) {
	attempt, err := r.auth.InitializeAttempt(ctx)
	if err != nil {
		return "", err
	}

	return r.auth.AuthorizationURI(ctx, attempt.State), nil
}

func (r *queryResolver) FinalizeAuth(ctx context.Context, code string, state string) (*skillz.AuthAttempt, error) {
	return r.user.Login(ctx, code, state)
}

func (r *queryResolver) User(ctx context.Context) (*skillz.User, error) {
	return internal.UserFromContext(ctx), nil
}

func (r *queryResolver) UserByID(ctx context.Context) (*skillz.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Character(ctx context.Context, id uint64) (*skillz.Character, error) {
	return r.character.Character(ctx, id)
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
