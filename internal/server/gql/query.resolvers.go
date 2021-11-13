package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/eveisesi/skillz"
	"github.com/eveisesi/skillz/internal/server/gql/generated"
	"github.com/gofrs/uuid"
)

func (r *authAttemptResolver) Status(ctx context.Context, obj *skillz.AuthAttempt) (string, error) {
	return skillz.StatusMap[obj.Status], nil
}

func (r *authAttemptResolver) URL(ctx context.Context, obj *skillz.AuthAttempt) (string, error) {
	return r.auth.AuthorizationURI(ctx, obj.State), nil
}

func (r *characterResolver) Corporation(ctx context.Context, obj *skillz.Character) (*skillz.Corporation, error) {
	return r.dataloaders.CorporationLoader().Load(ctx, obj.CorporationID)
}

func (r *corporationResolver) Alliance(ctx context.Context, obj *skillz.Corporation) (*skillz.Alliance, error) {
	if !obj.AllianceID.Valid {
		return nil, nil
	}

	return r.dataloaders.AllianceLoader().Load(ctx, obj.AllianceID.Uint)
}

func (r *queryResolver) Auth(ctx context.Context) (*skillz.AuthAttempt, error) {
	return r.auth.InitializeAttempt(ctx)
}

func (r *queryResolver) User(ctx context.Context, id uuid.UUID) (*skillz.User, error) {
	return r.user.User(ctx, id)
}

func (r *userResolver) Scopes(ctx context.Context, obj *skillz.User) ([]string, error) {
	var out = make([]string, 0, len(obj.Scopes))
	for _, scope := range obj.Scopes {
		out = append(out, scope.String())
	}
	return out, nil
}

func (r *userResolver) Character(ctx context.Context, obj *skillz.User) (*skillz.Character, error) {
	return r.dataloaders.CharacterLoader().Load(ctx, obj.CharacterID)
}

// AuthAttempt returns generated.AuthAttemptResolver implementation.
func (r *Resolver) AuthAttempt() generated.AuthAttemptResolver { return &authAttemptResolver{r} }

// Character returns generated.CharacterResolver implementation.
func (r *Resolver) Character() generated.CharacterResolver { return &characterResolver{r} }

// Corporation returns generated.CorporationResolver implementation.
func (r *Resolver) Corporation() generated.CorporationResolver { return &corporationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// User returns generated.UserResolver implementation.
func (r *Resolver) User() generated.UserResolver { return &userResolver{r} }

type authAttemptResolver struct{ *Resolver }
type characterResolver struct{ *Resolver }
type corporationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type userResolver struct{ *Resolver }
