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
	return r.clone.Clones(ctx, obj.CharacterID)
}

func (r *userResolver) Implants(ctx context.Context, obj *skillz.User) ([]*skillz.CharacterImplant, error) {
	return r.clone.Implants(ctx, obj.CharacterID)
}

func (r *userResolver) SkillMeta(ctx context.Context, obj *skillz.User) (*skillz.CharacterSkillMeta, error) {
	return r.skill.Meta(ctx, obj.CharacterID)
}

func (r *userResolver) Skills(ctx context.Context, obj *skillz.User) ([]*skillz.CharacterSkillGroup, error) {
	return r.skill.SkillsGrouped(ctx, obj.CharacterID)
}

func (r *userResolver) Queue(ctx context.Context, obj *skillz.User) ([]*skillz.CharacterSkillQueue, error) {
	return r.skill.SkillQueue(ctx, obj.CharacterID)
}

func (r *userResolver) Attributes(ctx context.Context, obj *skillz.User) (*skillz.CharacterAttributes, error) {
	return r.skill.Attributes(ctx, obj.CharacterID)
}

func (r *userResolver) Flyable(ctx context.Context, obj *skillz.User) ([]*skillz.CharacterFlyableShip, error) {
	return r.skill.Flyable(ctx, obj.CharacterID)
}

func (r *userResolver) Contacts(ctx context.Context, obj *skillz.User) ([]*skillz.CharacterContact, error) {
	return r.contact.Contacts(ctx, obj.CharacterID)
}

// User returns engine.UserResolver implementation.
func (r *Resolver) User() engine.UserResolver { return &userResolver{r} }

type userResolver struct{ *Resolver }
