package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/eveisesi/skillz"
	"github.com/eveisesi/skillz/internal"
	"github.com/eveisesi/skillz/internal/server/gql/generated"
)

func (r *queryResolver) Auth(ctx context.Context, state *string) (*skillz.AuthAttempt, error) {
	if state != nil {
		return r.auth.AuthAttempt(ctx, *state)
	}
	return r.auth.InitializeAttempt(ctx)
}

func (r *queryResolver) User(ctx context.Context) (*skillz.User, error) {
	return internal.UserFromContext(ctx), nil
}

func (r *queryResolver) Character(ctx context.Context, id uint64) (*skillz.Character, error) {
	return r.character.Character(ctx, id)
}

func (r *userResolver) Scopes(ctx context.Context, obj *skillz.User) ([]string, error) {
	var out = make([]string, 0, len(obj.Scopes))
	for _, scope := range obj.Scopes {
		out = append(out, scope.String())
	}
	return out, nil
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// User returns generated.UserResolver implementation.
func (r *Resolver) User() generated.UserResolver { return &userResolver{r} }

type queryResolver struct{ *Resolver }
type userResolver struct{ *Resolver }
