package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/eveisesi/skillz"
	"github.com/eveisesi/skillz/internal"
	"github.com/eveisesi/skillz/internal/server/gql/generated"
	"github.com/eveisesi/skillz/internal/server/gql/model"
)

func (r *characterResolver) Skills(ctx context.Context, obj *skillz.Character) (*model.CharacterSkills, error) {
	meta, err := r.skill.Skills(ctx, internal.UserFromContext(ctx))
	if err != nil {
		return nil, err
	}

	return &model.CharacterSkills{
		Meta:   meta,
		Skills: meta.Skills,
	}, nil
}

func (r *characterResolver) Queue(ctx context.Context, obj *skillz.Character) ([]*skillz.CharacterSkillQueue, error) {
	return r.skill.SkillQueue(ctx, internal.UserFromContext(ctx))
}

func (r *characterResolver) Attributes(ctx context.Context, obj *skillz.Character) (*skillz.CharacterAttributes, error) {
	return r.skill.Attributes(ctx, internal.UserFromContext(ctx))
}

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
