package character

import (
	"context"
	"database/sql"
	"time"

	"github.com/eveisesi/skillz"
	"github.com/eveisesi/skillz/internal/cache"
	"github.com/eveisesi/skillz/internal/esi"
	"github.com/pkg/errors"
)

type API interface {
	Character(ctx context.Context, character *skillz.Character) (*skillz.Character, error)
}

type Service struct {
	cache cache.CharacterAPI
	esi   esi.CharacterAPI

	character skillz.CharacterRepository
}

func New(cache cache.CharacterAPI, esi esi.CharacterAPI, charachter skillz.CharacterRepository) *Service {
	return &Service{
		cache:     cache,
		esi:       esi,
		character: charachter,
	}
}

func (s *Service) Character(ctx context.Context, character *skillz.Character) (*skillz.Character, error) {

	character, err := s.cache.Character(ctx, character.ID)
	if err != nil {
		return nil, err
	}

	if character != nil {
		return character, nil
	}

	character, err = s.cache.Character(ctx, character.ID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, errors.Wrap(err, "failed to fetch character record from data store")
	}

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	return character, s.cache.SetCharacter(ctx, character.ID, character, time.Hour)

}
