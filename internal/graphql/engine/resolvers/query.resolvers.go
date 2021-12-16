package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/eveisesi/skillz"
	"github.com/eveisesi/skillz/internal/graphql/engine"
)

func (r *queryResolver) InitializeAuth(ctx context.Context) (string, error) {
	attempt, err := r.auth.InitializeAttempt(ctx)
	if err != nil {
		return "", err
	}

	return r.auth.AuthorizationURI(ctx, attempt.State), nil
}

func (r *queryResolver) FinalizeAuth(ctx context.Context, code string, state string) (*skillz.User, error) {
	return r.user.Login(ctx, code, state)
}

func (r *queryResolver) SearchUser(ctx context.Context, term string) ([]*skillz.User, error) {
	return r.user.SearchUsers(ctx, term)
}

func (r *queryResolver) User(ctx context.Context, id uint64) (*skillz.User, error) {
	return r.user.UserByCharacterID(ctx, id)
}

func (r *queryResolver) Clones(ctx context.Context, id uint64) (*skillz.CharacterCloneMeta, error) {
	user, err := r.user.UserByCharacterID(ctx, id)
	if err != nil {
		return nil, err
	}

	return r.clone.Clones(ctx, user.CharacterID)
}

func (r *queryResolver) Implants(ctx context.Context, id uint64) ([]*skillz.CharacterImplant, error) {
	user, err := r.user.UserByCharacterID(ctx, id)
	if err != nil {
		return nil, err
	}

	return r.clone.Implants(ctx, user.CharacterID)
}

func (r *queryResolver) SkillMeta(ctx context.Context, id uint64) (*skillz.CharacterSkillMeta, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Skills(ctx context.Context, id uint64) ([]*skillz.CharacterSkillGroup, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Queue(ctx context.Context, id uint64) ([]*skillz.CharacterSkillQueue, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Attributes(ctx context.Context, id uint64) (*skillz.CharacterAttributes, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Flyable(ctx context.Context, id uint64) ([]*skillz.CharacterFlyableShip, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Contacts(ctx context.Context, id uint64) ([]*skillz.CharacterContact, error) {
	panic(fmt.Errorf("not implemented"))
}

// Query returns engine.QueryResolver implementation.
func (r *Resolver) Query() engine.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
func (r *queryResolver) Character(ctx context.Context, id uint64) (*skillz.Character, error) {
	return r.character.Character(ctx, id)
}
