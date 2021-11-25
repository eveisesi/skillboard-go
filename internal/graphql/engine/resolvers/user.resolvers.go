package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/eveisesi/skillz"
	"github.com/eveisesi/skillz/internal/graphql/engine"
)

func (r *userResolver) Scopes(ctx context.Context, obj *skillz.User) ([]string, error) {
	var out = make([]string, 0, len(obj.Scopes))
	for _, scope := range obj.Scopes {
		out = append(out, scope.String())
	}
	return out, nil
}

func (r *userResolver) Character(ctx context.Context, obj *skillz.User) (*skillz.Character, error) {
	return r.character.Character(ctx, obj.CharacterID)
}

func (r *userResolver) Clone(ctx context.Context, obj *skillz.User) (*skillz.CharacterCloneMeta, error) {
	return r.clone.Clones(ctx, obj)
}

func (r *userResolver) Implants(ctx context.Context, obj *skillz.User) ([]*skillz.CharacterImplant, error) {
	return r.clone.Implants(ctx, obj)
}

func (r *userResolver) Skills(ctx context.Context, obj *skillz.User) (*engine.CharacterSkills, error) {
	meta, err := r.skill.Meta(ctx, obj.CharacterID)
	if err != nil {
		return nil, err
	}

	skills, err := r.skill.Skillz(ctx, obj.CharacterID)
	if err != nil {
		return nil, err
	}

	return &engine.CharacterSkills{
		Meta:   meta,
		Skills: skills,
	}, nil
}

func (r *userResolver) Queue(ctx context.Context, obj *skillz.User) ([]*skillz.CharacterSkillQueue, error) {
	return r.skill.SkillQueue(ctx, obj.CharacterID)
}

func (r *userResolver) Attributes(ctx context.Context, obj *skillz.User) (*skillz.CharacterAttributes, error) {
	return r.skill.Attributes(ctx, obj.CharacterID)
}

// User returns engine.UserResolver implementation.
func (r *Resolver) User() engine.UserResolver { return &userResolver{r} }

type userResolver struct{ *Resolver }
