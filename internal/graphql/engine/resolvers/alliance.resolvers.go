package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/eveisesi/skillz"
)

func (r *corporationResolver) Alliance(ctx context.Context, obj *skillz.Corporation) (*skillz.Alliance, error) {
	if !obj.AllianceID.Valid {
		return nil, nil
	}

	return r.dataloaders.AllianceLoader().Load(ctx, obj.AllianceID.Uint)
}
