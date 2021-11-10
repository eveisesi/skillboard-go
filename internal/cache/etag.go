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
	EtagByPath(ctx context.Context, pathHashed string) (*skillz.Etag, error)
	SetEtag(ctx context.Context, etagID string, etag *skillz.Etag, expires time.Duration) error
}

const (
	etagKeyPrefix = "etag"
)

func (s *Service) EtagByPath(ctx context.Context, path string) (*skillz.Etag, error) {

	key := generateKey(etagKeyPrefix, hash(path))

	result, err := s.redis.Get(ctx, key).Bytes()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, errors.Wrap(err, "failed to query cache for record")
	}

	if errors.Is(err, redis.Nil) {
		return nil, nil
	}

	etag := new(skillz.Etag)
	err = json.Unmarshal(result, etag)

	return etag, errors.Wrap(err, "failed to decode cached etag to structure")

}

func (s *Service) SetEtag(ctx context.Context, path string, etag *skillz.Etag, expires time.Duration) error {

	data, err := json.Marshal(etag)
	if err != nil {
		return errors.Wrap(err, "failed to encode etag to json")
	}

	return errors.Wrap(
		s.redis.Set(ctx, generateKey(etagKeyPrefix, hash(path)), data, expires).Err(),
		"failed to cache etag",
	)

}
