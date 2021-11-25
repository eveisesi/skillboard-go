package cache

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

type PageAPI interface {
	Page(ctx context.Context, path string) (string, error)
	SetPage(ctx context.Context, path, contents string, duration time.Duration) error
}

const (
	pageKeyPrefix = "page"
)

func (s *Service) Page(ctx context.Context, path string) (string, error) {

	key := generateKey(pageKeyPrefix, hash(path))

	result, err := s.redis.Get(ctx, key).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return "", errors.Wrapf(err, errorFFormat, pageAPI, "Page", "failed to fetch results from cache")
	}

	return result, nil

}

func (s *Service) SetPage(ctx context.Context, path, content string, expires time.Duration) error {

	data := []byte(content)

	key := generateKey(pageKeyPrefix, hash(path))
	err := s.redis.Set(ctx, key, data, expires).Err()
	return errors.Wrapf(err, errorFFormat, pageAPI, "SetPage", "failed to write cache")

}
