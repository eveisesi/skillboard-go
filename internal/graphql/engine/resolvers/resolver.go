package resolvers

import (
	"github.com/eveisesi/skillz/internal/alliance"
	"github.com/eveisesi/skillz/internal/auth"
	"github.com/eveisesi/skillz/internal/character"
	"github.com/eveisesi/skillz/internal/clone"
	"github.com/eveisesi/skillz/internal/corporation"
	"github.com/eveisesi/skillz/internal/graphql/dataloaders"
	"github.com/eveisesi/skillz/internal/skill"
	"github.com/eveisesi/skillz/internal/user"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	auth        auth.API
	alliance    alliance.API
	character   character.API
	clone       clone.API
	corporation corporation.API
	dataloaders dataloaders.API
	skill       skill.API
	user        user.API
}

func New(
	alliance alliance.API,
	auth auth.API,
	character character.API,
	clone clone.API,
	corporation corporation.API,
	dataloaders dataloaders.API,
	skill skill.API,
	user user.API,
) *Resolver {
	return &Resolver{
		alliance:    alliance,
		auth:        auth,
		character:   character,
		clone:       clone,
		corporation: corporation,
		dataloaders: dataloaders,
		skill:       skill,
		user:        user,
	}
}
