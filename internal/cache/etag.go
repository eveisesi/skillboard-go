package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/eveisesi/skillz"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

type EtagAPI interface {
	EtagByPath(ctx context.Context, path string) (*skillz.Etag, error)
	SetEtag(ctx context.Context, etag *skillz.Etag, expires time.Duration) error
}

const (
	etagKeyPrefix = "etag"
)

func (s *Service) EtagByPath(ctx context.Context, path string) (*skillz.Etag, error) {
	if s.disabled {
		return nil, nil
	}
	key := generateKey(etagKeyPrefix, hash(path))

	result, err := s.redis.Get(ctx, key).Bytes()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, errors.Wrapf(err, errorFFormat, etagAPI, "EtagByPath", "failed to fetch results from cache")
	}

	if errors.Is(err, redis.Nil) {
		return nil, nil
	}

	etag := new(skillz.Etag)
	err = json.Unmarshal(result, etag)
	return etag, errors.Wrapf(err, errorFFormat, etagAPI, "EtagByPath", "failed to decode json to structure")

}

func (s *Service) SetEtag(ctx context.Context, etag *skillz.Etag, expires time.Duration) error {
	if s.disabled {
		return nil
	}
	data, err := json.Marshal(etag)
	if err != nil {
		return errors.Wrapf(err, errorFFormat, universeAPI, "SetBloodline", "failed to encode struct as json")
	}

	key := generateKey(etagKeyPrefix, hash(etag.Path))
	err = s.redis.Set(ctx, key, data, expires).Err()
	return errors.Wrapf(err, errorFFormat, etagAPI, "SetEtag", "failed to write cache")

}
