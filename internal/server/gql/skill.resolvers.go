package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/eveisesi/skillz"
	"github.com/eveisesi/skillz/internal/server/gql/generated"
)

func (r *characterSkillResolver) Info(ctx context.Context, obj *skillz.CharacterSkill) (*skillz.Type, error) {
	return r.dataloaders.TypeLoader().Load(ctx, obj.SkillID)
}

func (r *characterSkillQueueResolver) Info(ctx context.Context, obj *skillz.CharacterSkillQueue) (*skillz.Type, error) {
	return r.dataloaders.TypeLoader().Load(ctx, obj.SkillID)
}

// CharacterSkill returns generated.CharacterSkillResolver implementation.
func (r *Resolver) CharacterSkill() generated.CharacterSkillResolver {
	return &characterSkillResolver{r}
}

// CharacterSkillQueue returns generated.CharacterSkillQueueResolver implementation.
func (r *Resolver) CharacterSkillQueue() generated.CharacterSkillQueueResolver {
	return &characterSkillQueueResolver{r}
}

type characterSkillResolver struct{ *Resolver }
type characterSkillQueueResolver struct{ *Resolver }
