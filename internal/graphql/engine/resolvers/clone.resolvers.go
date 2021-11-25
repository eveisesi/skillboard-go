package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/eveisesi/skillz"
	"github.com/eveisesi/skillz/internal/graphql/engine"
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

func (r *characterDeathCloneResolver) Location(ctx context.Context, obj *skillz.CharacterDeathClone) (engine.LocationInfo, error) {
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

func (r *characterJumpCloneResolver) Location(ctx context.Context, obj *skillz.CharacterJumpClone) (engine.LocationInfo, error) {
	switch obj.LocationType {
	case skillz.CloneLocationTypeStation:
		return r.dataloaders.StationLoader().Load(ctx, uint(obj.LocationID))
	case skillz.CloneLocationTypeStructure:
		return r.dataloaders.StructureLoader().Load(ctx, obj.LocationID)
	default:
		return nil, fmt.Errorf("%s is not a resolvable location type", obj.LocationType.String())
	}
}

// CharacterClone returns engine.CharacterCloneResolver implementation.
func (r *Resolver) CharacterClone() engine.CharacterCloneResolver { return &characterCloneResolver{r} }

// CharacterDeathClone returns engine.CharacterDeathCloneResolver implementation.
func (r *Resolver) CharacterDeathClone() engine.CharacterDeathCloneResolver {
	return &characterDeathCloneResolver{r}
}

// CharacterImplant returns engine.CharacterImplantResolver implementation.
func (r *Resolver) CharacterImplant() engine.CharacterImplantResolver {
	return &characterImplantResolver{r}
}

// CharacterJumpClone returns engine.CharacterJumpCloneResolver implementation.
func (r *Resolver) CharacterJumpClone() engine.CharacterJumpCloneResolver {
	return &characterJumpCloneResolver{r}
}

type characterCloneResolver struct{ *Resolver }
type characterDeathCloneResolver struct{ *Resolver }
type characterImplantResolver struct{ *Resolver }
type characterJumpCloneResolver struct{ *Resolver }
