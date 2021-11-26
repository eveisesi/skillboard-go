package etag

import (
	"context"
	"database/sql"
	"time"

	"github.com/eveisesi/skillz"
	"github.com/eveisesi/skillz/internal/cache"
	"github.com/pkg/errors"
)

type API interface {
	Etag(ctx context.Context, path string) (*skillz.Etag, error)
	InsertEtag(ctx context.Context, etag *skillz.Etag) error
}

type Service struct {
	cache cache.EtagAPI

	etag skillz.EtagRepository
}

func New(cache cache.EtagAPI, etag skillz.EtagRepository) *Service {
	return &Service{
		cache: cache,
		etag:  etag,
	}
}

func (s *Service) Etag(ctx context.Context, path string) (*skillz.Etag, error) {

	etag, err := s.cache.EtagByPath(ctx, path)
	if err != nil {
		return nil, err
	}

	if etag != nil {
		return etag, nil
	}

	etag, err = s.etag.Etag(ctx, path)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, errors.Wrap(err, "failed to fetch etag from data store")
	}

	return etag, s.cache.SetEtag(ctx, etag, time.Minute*5)

}

func (s *Service) InsertEtag(ctx context.Context, etag *skillz.Etag) error {

	err := s.etag.InsertEtag(ctx, etag)
	if err != nil {
		return errors.Wrap(err, "failed to persist etag to datastore")
	}

	return s.cache.SetEtag(ctx, etag, time.Since(etag.CachedUntil))

}
