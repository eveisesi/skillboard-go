package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/eveisesi/skillz"
	"github.com/go-redis/redis"
	"github.com/pkg/errors"
)

type CharacterAPI interface {
	Character(ctx context.Context, characterID uint64) (*skillz.Character, error)
	SetCharacter(ctx context.Context, characterID uint64, character *skillz.Character, expires time.Duration) error
}

const (
	characterKeyPrefix = "characters"
)

func (s *Service) Character(ctx context.Context, characterID uint64) (*skillz.Character, error) {

	key := generateKey(characterKeyPrefix, hashUint64(characterID))

	result, err := s.redis.Get(ctx, key).Bytes()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, errors.Wrap(err, "failed to query cahce for character")
	}

	if errors.Is(err, redis.Nil) {
		return nil, nil
	}

	character := new(skillz.Character)
	return character, errors.Wrap(json.Unmarshal(result, character), "failed to decode cached character to structure")

}

func (s *Service) SetCharacter(ctx context.Context, characterID uint64, character *skillz.Character, expires time.Duration) error {

	data, err := json.Marshal(character)
	if err != nil {
		return errors.Wrap(err, "failed to encode character to json")
	}

	return errors.Wrap(
		s.redis.Set(ctx, generateKey(characterKeyPrefix, hashUint64(characterID)), data, expires).Err(),
		"failed to cache character",
	)

}
