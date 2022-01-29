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

type CorporationAPI interface {
	Corporation(ctx context.Context, corporationID uint) (*skillz.Corporation, error)
	SetCorporation(ctx context.Context, corporation *skillz.Corporation, expires time.Duration) error
}

const (
	corporationKeyPrefix = "corporations"
)

func (s *Service) Corporation(ctx context.Context, corporationID uint) (*skillz.Corporation, error) {
	if s.disabled {
		return nil, nil
	}
	key := generateKey(corporationKeyPrefix, strconv.FormatUint(uint64(corporationID), 10))

	result, err := s.redis.Get(ctx, key).Bytes()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, errors.Wrapf(err, errorFFormat, corporationAPI, "Corporation", "failed to fetch results from cache")
	}

	if errors.Is(err, redis.Nil) {
		return nil, nil
	}

	var corporation = new(skillz.Corporation)
	err = json.Unmarshal(result, corporation)
	return corporation, errors.Wrapf(err, errorFFormat, corporationAPI, "Corporation", "failed to decode json to structure")

}

func (s *Service) SetCorporation(ctx context.Context, corporation *skillz.Corporation, expires time.Duration) error {
	if s.disabled {
		return nil
	}
	data, err := json.Marshal(corporation)
	if err != nil {
		return errors.Wrapf(err, errorFFormat, corporationAPI, "SetCorporation", "failed to encode struct as json")
	}

	key := generateKey(corporationKeyPrefix, strconv.FormatUint(uint64(corporation.ID), 10))
	err = s.redis.Set(ctx, key, data, expires).Err()
	return errors.Wrapf(err, errorFFormat, corporationAPI, "SetCorporation", "failed to write cache")

}
