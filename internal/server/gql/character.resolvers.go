package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/eveisesi/skillz"
	"github.com/eveisesi/skillz/internal/server/gql/generated"
)

func (r *characterImplantResolver) Character(ctx context.Context, obj *skillz.CharacterImplant) (*skillz.Character, error) {
	return r.dataloaders.CharacterLoader().Load(ctx, obj.CharacterID)
}

func (r *userResolver) Character(ctx context.Context, obj *skillz.User) (*skillz.Character, error) {
	return r.dataloaders.CharacterLoader().Load(ctx, obj.CharacterID)
}

// Character returns generated.CharacterResolver implementation.
func (r *Resolver) Character() generated.CharacterResolver { return &characterResolver{r} }

type characterResolver struct{ *Resolver }
