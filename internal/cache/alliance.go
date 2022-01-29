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

type AllianceAPI interface {
	Alliance(ctx context.Context, allianceID uint) (*skillz.Alliance, error)
	SetAlliance(ctx context.Context, alliance *skillz.Alliance, expires time.Duration) error
}

const (
	allianceKeyPrefix = "alliances"
)

func (s *Service) Alliance(ctx context.Context, allianceID uint) (*skillz.Alliance, error) {
	if s.disabled {
		return nil, nil
	}

	key := generateKey(allianceKeyPrefix, strconv.Itoa(int(allianceID)))

	result, err := s.redis.Get(ctx, key).Bytes()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, errors.Wrapf(err, errorFFormat, allianceAPI, "Alliance", "failed to fetch results from cache")
	}

	if errors.Is(err, redis.Nil) {
		return nil, nil
	}

	var alliance = new(skillz.Alliance)
	err = json.Unmarshal(result, alliance)
	return alliance, errors.Wrapf(err, errorFFormat, allianceAPI, "Alliance", "failed to decode json to structure")

}

func (s *Service) SetAlliance(ctx context.Context, alliance *skillz.Alliance, expires time.Duration) error {

	data, err := json.Marshal(alliance)
	if err != nil {
		return errors.Wrapf(err, errorFFormat, allianceAPI, "SetAlliance", "failed to encode struct as json")
	}

	key := generateKey(allianceKeyPrefix, strconv.Itoa(int(alliance.ID)))
	err = s.redis.Set(ctx, key, data, expires).Err()
	return errors.Wrapf(err, errorFFormat, allianceAPI, "SetAlliance", "failed to write cache")

}
