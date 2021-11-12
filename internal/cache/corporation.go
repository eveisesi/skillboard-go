package cache

import (
	"context"
	"encoding/json"
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

	key := generateKey(corporationKeyPrefix, hashUint(corporationID))

	result, err := s.redis.Get(ctx, key).Bytes()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}

	if errors.Is(err, redis.Nil) {
		return nil, nil
	}

	var corporation = new(skillz.Corporation)
	err = json.Unmarshal(result, corporation)
	return corporation, errors.Wrapf(err, "failed to decode cached corporation to structure")

}

func (s *Service) SetCorporation(ctx context.Context, corporation *skillz.Corporation, expires time.Duration) error {

	data, err := json.Marshal(corporation)
	if err != nil {
		return errors.Wrap(err, "failed to encode corporation to json")
	}

	key := generateKey(corporationKeyPrefix, hashUint(corporation.ID))

	err = s.redis.Set(ctx, key, data, expires).Err()
	return errors.Wrap(err, "failed to cache corporation")

}
