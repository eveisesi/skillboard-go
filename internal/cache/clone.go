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

func (s *Service) CharacterImplants(ctx context.Context, characterID uint64) ([]*skillz.CharacterImplant, error) {
	if s.disabled {
		return nil, nil
	}
	var implants = make([]*skillz.CharacterImplant, 0, 10)

	key := generateKey(characterImplantsKeyPrefix, strconv.FormatUint(characterID, 10))
	result, err := s.redis.Get(ctx, key).Bytes()
	if err != nil && !errors.Is(err, redis.Nil) {
		return implants, errors.Wrapf(err, errorFFormat, cloneAPI, "CharacterImplants", "failed to fetch results from cache")
	}

	if errors.Is(err, redis.Nil) {
		return implants, nil
	}

	err = json.Unmarshal(result, &implants)
	return implants, errors.Wrapf(err, errorFFormat, cloneAPI, "CharacterImplants", "failed to decode json to structure")

}

func (s *Service) SetCharacterImplants(ctx context.Context, characterID uint64, implants []*skillz.CharacterImplant, expires time.Duration) error {
	if s.disabled {
		return nil
	}
	data, err := json.Marshal(implants)
	if err != nil {
		return errors.Wrapf(err, errorFFormat, cloneAPI, "SetCharacterImplants", "failed to encode struct as json")
	}

	key := generateKey(characterImplantsKeyPrefix, strconv.FormatUint(characterID, 10))
	err = s.redis.Set(ctx, key, data, expires).Err()
	return errors.Wrapf(err, errorFFormat, cloneAPI, "SetCharacterImplants", "failed to write cache")

}
