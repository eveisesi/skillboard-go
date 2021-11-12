package cache

import (
	"context"
	"encoding/json"
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

	key := generateKey(characterKeyPrefix, hashUint64(characterID))

	result, err := s.redis.Get(ctx, key).Bytes()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}

	if errors.Is(err, redis.Nil) {
		return nil, nil
	}

	var character = new(skillz.Character)
	err = json.Unmarshal(result, character)
	return character, errors.Wrapf(err, "failed to decode cached character to structure")

}

func (s *Service) SetCharacter(ctx context.Context, character *skillz.Character, expires time.Duration) error {

	data, err := json.Marshal(character)
	if err != nil {
		return errors.Wrap(err, "failed to encode character to json")
	}

	key := generateKey(characterKeyPrefix, hashUint64(character.ID))

	err = s.redis.Set(ctx, key, data, expires).Err()
	return errors.Wrap(err, "failed to cache character")

}
