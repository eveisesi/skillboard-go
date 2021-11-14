package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/eveisesi/skillz"
	"github.com/eveisesi/skillz/internal/server/gql/generated"
)

func (r *authAttemptResolver) Status(ctx context.Context, obj *skillz.AuthAttempt) (string, error) {
	return skillz.StatusMap[obj.Status], nil
}

func (r *authAttemptResolver) URL(ctx context.Context, obj *skillz.AuthAttempt) (string, error) {
	return r.auth.AuthorizationURI(ctx, obj.State), nil
}

// AuthAttempt returns generated.AuthAttemptResolver implementation.
func (r *Resolver) AuthAttempt() generated.AuthAttemptResolver { return &authAttemptResolver{r} }

type authAttemptResolver struct{ *Resolver }
