package dataloaders

import (
	"time"

	"github.com/eveisesi/skillz/internal/alliance"
	"github.com/eveisesi/skillz/internal/character"
	"github.com/eveisesi/skillz/internal/clone"
	"github.com/eveisesi/skillz/internal/corporation"
	"github.com/eveisesi/skillz/internal/server/gql/dataloaders/generated"
	"github.com/eveisesi/skillz/internal/universe"
)

type API interface {
	AllianceLoader() *generated.AllianceLoader
	CharacterLoader() *generated.CharacterLoader
	CorporationLoader() *generated.CorporationLoader
}

type Service struct {
	wait  time.Duration
	batch int

	character   character.API
	corporation corporation.API
	alliance    alliance.API
	clone       clone.API
	universe    universe.API
}

func New(wait time.Duration, batch int, character character.API, corporation corporation.API, alliance alliance.API, clone clone.API, universe universe.API) *Service {
	return &Service{
		wait, batch, character, corporation, alliance, clone, universe,
	}
}
