package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/eveisesi/skillz"
	"github.com/eveisesi/skillz/internal/graphql/engine"
)

func (r *queryResolver) InitializeAuth(ctx context.Context) (string, error) {
	attempt, err := r.auth.InitializeAttempt(ctx)
	if err != nil {
		return "", err
	}

	return r.auth.AuthorizationURI(ctx, attempt.State), nil
}

func (r *queryResolver) FinalizeAuth(ctx context.Context, code string, state string) (*skillz.User, error) {
	return r.user.Login(ctx, code, state)
}

func (r *queryResolver) User(ctx context.Context, id uint64) (*skillz.User, error) {
	return r.user.UserByCharacterID(ctx, id)
}

func (r *queryResolver) Character(ctx context.Context, id uint64) (*skillz.Character, error) {
	return r.character.Character(ctx, id)
}

// Query returns engine.QueryResolver implementation.
func (r *Resolver) Query() engine.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
