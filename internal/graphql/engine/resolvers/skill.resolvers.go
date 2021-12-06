package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/eveisesi/skillz"
	"github.com/eveisesi/skillz/internal/graphql/engine"
)

func (r *characterFlyableShipResolver) Info(ctx context.Context, obj *skillz.CharacterFlyableShip) (*skillz.Type, error) {
	return r.dataloaders.TypeLoader().Load(ctx, obj.ShipTypeID)
}

func (r *characterSkillResolver) Info(ctx context.Context, obj *skillz.CharacterSkill) (*skillz.Type, error) {
	return r.dataloaders.TypeLoader().Load(ctx, obj.SkillID)
}

func (r *characterSkillQueueResolver) Info(ctx context.Context, obj *skillz.CharacterSkillQueue) (*skillz.Type, error) {
	return r.dataloaders.TypeLoader().Load(ctx, obj.SkillID)
}

// CharacterFlyableShip returns engine.CharacterFlyableShipResolver implementation.
func (r *Resolver) CharacterFlyableShip() engine.CharacterFlyableShipResolver {
	return &characterFlyableShipResolver{r}
}

// CharacterSkill returns engine.CharacterSkillResolver implementation.
func (r *Resolver) CharacterSkill() engine.CharacterSkillResolver { return &characterSkillResolver{r} }

// CharacterSkillQueue returns engine.CharacterSkillQueueResolver implementation.
func (r *Resolver) CharacterSkillQueue() engine.CharacterSkillQueueResolver {
	return &characterSkillQueueResolver{r}
}

type characterFlyableShipResolver struct{ *Resolver }
type characterSkillResolver struct{ *Resolver }
type characterSkillQueueResolver struct{ *Resolver }
