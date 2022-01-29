package cache

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/eveisesi/skillz"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

type CharacterAPI interface {
	Character(ctx context.Context, characterID uint64) (*skillz.Character, error)
	SetCharacter(ctx context.Context, character *skillz.Character, expires time.Duration) error
}

const (
	characterKeyPrefix = "characters"
)

func (s *Service) Character(ctx context.Context, characterID uint64) (*skillz.Character, error) {
	if s.disabled {
		return nil, nil
	}
	key := generateKey(characterKeyPrefix, strconv.FormatUint(characterID, 10))
	result, err := s.redis.Get(ctx, key).Bytes()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, errors.Wrapf(err, errorFFormat, characterAPI, "Character", "failed to fetch results from cache")
	}

	if errors.Is(err, redis.Nil) {
		return nil, nil
	}

	var character = new(skillz.Character)
	err = json.Unmarshal(result, character)
	return character, errors.Wrapf(err, errorFFormat, characterAPI, "Character", "failed to decode json to structure")

}

func (s *Service) SetCharacter(ctx context.Context, character *skillz.Character, expires time.Duration) error {
	if s.disabled {
		return nil
	}
	data, err := json.Marshal(character)
	if err != nil {
		return errors.Wrapf(err, errorFFormat, characterAPI, "SetCharacter", "failed to encode struct as json")
	}

	key := generateKey(characterKeyPrefix, strconv.FormatUint(character.ID, 10))
	err = s.redis.Set(ctx, key, data, expires).Err()
	return errors.Wrapf(err, errorFFormat, characterAPI, "SetCharacter", "failed to write cache")

}
