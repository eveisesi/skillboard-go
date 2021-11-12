package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/eveisesi/skillz"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

type AllianceAPI interface {
	Alliance(ctx context.Context, allianceID uint) (*skillz.Alliance, error)
	SetAlliance(ctx context.Context, alliance *skillz.Alliance, expires time.Duration) error
}

const (
	allianceKeyPrefix = "alliances"
)

func (s *Service) Alliance(ctx context.Context, allianceID uint) (*skillz.Alliance, error) {

	key := generateKey(allianceKeyPrefix, hashUint(allianceID))

	result, err := s.redis.Get(ctx, key).Bytes()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}

	if errors.Is(err, redis.Nil) {
		return nil, nil
	}

	var alliance = new(skillz.Alliance)
	err = json.Unmarshal(result, alliance)
	return alliance, errors.Wrapf(err, "failed to decode cached alliance to structure")

}

func (s *Service) SetAlliance(ctx context.Context, alliance *skillz.Alliance, expires time.Duration) error {

	data, err := json.Marshal(alliance)
	if err != nil {
		return errors.Wrap(err, "failed to encode alliance to json")
	}

	key := generateKey(allianceKeyPrefix, hashUint(alliance.ID))

	err = s.redis.Set(ctx, key, data, expires).Err()
	return errors.Wrap(err, "failed to cache alliance")

}
