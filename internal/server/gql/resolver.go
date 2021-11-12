package resolvers

import (
	"github.com/eveisesi/skillz/internal/auth"
	"github.com/eveisesi/skillz/internal/character"
	"github.com/eveisesi/skillz/internal/user"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	auth      auth.API
	character character.API
	user      user.API
}

func New(
	auth auth.API,
	character character.API,
	user user.API,
) *Resolver {
	return &Resolver{
		auth:      auth,
		character: character,
		user:      user,
	}
}
