package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/eveisesi/skillz"
	"github.com/eveisesi/skillz/internal/graphql/engine"
)

func (r *characterResolver) Corporation(ctx context.Context, obj *skillz.Character) (*skillz.Corporation, error) {
	return r.dataloaders.CorporationLoader().Load(ctx, obj.CorporationID)
}

// Corporation returns engine.CorporationResolver implementation.
func (r *Resolver) Corporation() engine.CorporationResolver { return &corporationResolver{r} }

type corporationResolver struct{ *Resolver }