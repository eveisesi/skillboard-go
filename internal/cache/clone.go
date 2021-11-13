package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/eveisesi/skillz"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

type CloneAPI interface {
	CharacterClones(ctx context.Context, characterID uint64) (*skillz.CharacterCloneMeta, error)
	SetCharacterClones(ctx context.Context, characterID uint64, clones *skillz.CharacterCloneMeta, expires time.Duration) error
	CharacterImplants(ctx context.Context, characterID uint64) ([]*skillz.CharacterImplant, error)
	SetCharacterImplants(ctx context.Context, characterID uint64, implants []*skillz.CharacterImplant, expires time.Duration) error
}

const (
	characterClonesKeyPrefix   = "character::clones"
	characterImplantsKeyPrefix = "character::implants"
)

func (s *Service) CharacterClones(ctx context.Context, characterID uint64) (*skillz.CharacterCloneMeta, error) {

	key := generateKey(characterClonesKeyPrefix, hashUint64(characterID))

	result, err := s.redis.Get(ctx, key).Bytes()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}

	if errors.Is(err, redis.Nil) {
		return nil, nil
	}

	var clones = new(skillz.CharacterCloneMeta)
	err = json.Unmarshal(result, clones)
	return clones, errors.Wrapf(err, "failed to decode cached clones to structure")
}

func (s *Service) SetCharacterClones(ctx context.Context, characterID uint64, clones *skillz.CharacterCloneMeta, expires time.Duration) error {

	data, err := json.Marshal(clones)
	if err != nil {
		return errors.Wrap(err, "failed to encode character clones to json")
	}

	key := generateKey(characterClonesKeyPrefix, hashUint64(characterID))

	err = s.redis.Set(ctx, key, data, expires).Err()
	return errors.Wrap(err, "failed to cache character clones")

}

func (s *Service) CharacterImplants(ctx context.Context, characterID uint64) ([]*skillz.CharacterImplant, error) {

	key := generateKey(characterImplantsKeyPrefix, hashUint64(characterID))

	result, err := s.redis.Get(ctx, key).Bytes()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}

	if errors.Is(err, redis.Nil) {
		return nil, nil
	}

	var implants = make([]*skillz.CharacterImplant, 0, 10)
	err = json.Unmarshal(result, &implants)
	return implants, errors.Wrapf(err, "failed to decode cached implants to structure")

}

func (s *Service) SetCharacterImplants(ctx context.Context, characterID uint64, implants []*skillz.CharacterImplant, expires time.Duration) error {

	data, err := json.Marshal(implants)
	if err != nil {
		return errors.Wrap(err, "failed to encode character implants to json")
	}

	key := generateKey(characterImplantsKeyPrefix, hashUint64(characterID))

	err = s.redis.Set(ctx, key, data, expires).Err()
	return errors.Wrap(err, "failed to cache character implants")

}
