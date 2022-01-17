package cache

import (
	"context"
	"time"

	"github.com/pkg/errors"
)

type PageAPI interface {
	Page(ctx context.Context, page string) (string, error)
	SetPage(ctx context.Context, page string, data []byte, expires time.Duration) error
}

const (
	pagePrefix = "page"
)

func (s *Service) Page(ctx context.Context, page string) (string, error) {

	key := generateKey(pagePrefix, hash(page))

	result, err := s.redis.Get(ctx, key).Result()
	return result, errors.Wrapf(err, errorFFormat, pageAPI, "Page", "failed to fetch results from cache")

}

func (s *Service) SetPage(ctx context.Context, page string, data []byte, expires time.Duration) error {

	key := generateKey(pagePrefix, hash(page))

	err := s.redis.Set(ctx, key, data, expires).Err()
	return errors.Wrapf(err, errorFFormat, pageAPI, "SetPage", "failed to write cache")

}
