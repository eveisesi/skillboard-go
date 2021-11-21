package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/eveisesi/skillz"
	"github.com/eveisesi/skillz/internal/server/gql/generated"
	"github.com/eveisesi/skillz/internal/server/gql/model"
)

func (r *characterCloneResolver) Jump(ctx context.Context, obj *skillz.CharacterCloneMeta) ([]*skillz.CharacterJumpClone, error) {
	return obj.JumpClones, nil
}

func (r *characterCloneResolver) Death(ctx context.Context, obj *skillz.CharacterCloneMeta) (*skillz.CharacterDeathClone, error) {
	return obj.HomeLocation, nil
}

func (r *characterDeathCloneResolver) LocationType(ctx context.Context, obj *skillz.CharacterDeathClone) (string, error) {
	return obj.LocationType.String(), nil
}

func (r *characterDeathCloneResolver) Location(ctx context.Context, obj *skillz.CharacterDeathClone) (model.LocationInfo, error) {
	switch obj.LocationType {
	case skillz.CloneLocationTypeStation:
		return r.dataloaders.StationLoader().Load(ctx, uint(obj.LocationID))
	case skillz.CloneLocationTypeStructure:
		return r.dataloaders.StructureLoader().Load(ctx, obj.LocationID)
	default:
		return nil, fmt.Errorf("%s is not a resolvable location type", obj.LocationType.String())
	}
}

func (r *characterJumpCloneResolver) LocationType(ctx context.Context, obj *skillz.CharacterJumpClone) (string, error) {
	return obj.LocationType.String(), nil
}

func (r *characterJumpCloneResolver) Location(ctx context.Context, obj *skillz.CharacterJumpClone) (model.LocationInfo, error) {
	switch obj.LocationType {
	case skillz.CloneLocationTypeStation:
		return r.dataloaders.StationLoader().Load(ctx, uint(obj.LocationID))
	case skillz.CloneLocationTypeStructure:
		return r.dataloaders.StructureLoader().Load(ctx, obj.LocationID)
	default:
		return nil, fmt.Errorf("%s is not a resolvable location type", obj.LocationType.String())
	}
}

// CharacterClone returns generated.CharacterCloneResolver implementation.
func (r *Resolver) CharacterClone() generated.CharacterCloneResolver {
	return &characterCloneResolver{r}
}

// CharacterDeathClone returns generated.CharacterDeathCloneResolver implementation.
func (r *Resolver) CharacterDeathClone() generated.CharacterDeathCloneResolver {
	return &characterDeathCloneResolver{r}
}

// CharacterImplant returns generated.CharacterImplantResolver implementation.
func (r *Resolver) CharacterImplant() generated.CharacterImplantResolver {
	return &characterImplantResolver{r}
}

// CharacterJumpClone returns generated.CharacterJumpCloneResolver implementation.
func (r *Resolver) CharacterJumpClone() generated.CharacterJumpCloneResolver {
	return &characterJumpCloneResolver{r}
}

type characterCloneResolver struct{ *Resolver }
type characterDeathCloneResolver struct{ *Resolver }
type characterImplantResolver struct{ *Resolver }
type characterJumpCloneResolver struct{ *Resolver }
