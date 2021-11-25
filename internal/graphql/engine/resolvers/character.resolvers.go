package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"github.com/eveisesi/skillz/internal/graphql/engine"
)

// Character returns engine.CharacterResolver implementation.
func (r *Resolver) Character() engine.CharacterResolver { return &characterResolver{r} }

type characterResolver struct{ *Resolver }
