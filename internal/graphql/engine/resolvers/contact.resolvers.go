package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/eveisesi/skillz"
	"github.com/eveisesi/skillz/internal/graphql/engine"
)

func (r *characterContactResolver) ContactType(ctx context.Context, obj *skillz.CharacterContact) (string, error) {
	return string(obj.ContactType), nil
}

// CharacterContact returns engine.CharacterContactResolver implementation.
func (r *Resolver) CharacterContact() engine.CharacterContactResolver {
	return &characterContactResolver{r}
}

type characterContactResolver struct{ *Resolver }
