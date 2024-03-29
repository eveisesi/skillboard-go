package alliance

import (
	"context"
	"database/sql"
	"time"

	"github.com/eveisesi/skillz"
	"github.com/eveisesi/skillz/internal/cache"
	"github.com/eveisesi/skillz/internal/esi"
	"github.com/eveisesi/skillz/internal/etag"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/volatiletech/null"
)

type API interface {
	Alliance(ctx context.Context, alliance uint) (*skillz.Alliance, error)
}

type Service struct {
	logger *logrus.Logger
	cache  cache.AllianceAPI
	esi    esi.AllianceAPI
	etag   etag.API

	alliance skillz.AllianceRepository
}

var _ API = new(Service)

func New(logger *logrus.Logger, cache cache.AllianceAPI, esi esi.AllianceAPI, etag etag.API, alliance skillz.AllianceRepository) *Service {
	return &Service{
		logger:   logger,
		cache:    cache,
		esi:      esi,
		etag:     etag,
		alliance: alliance,
	}
}

func (s *Service) Alliance(ctx context.Context, allianceID uint) (*skillz.Alliance, error) {

	alliance, err := s.cache.Alliance(ctx, allianceID)
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}

	if alliance != nil {
		return alliance, nil
	}

	etagID, err := esi.Resolvers[esi.GetAlliance](&esi.Params{AllianceID: null.UintFrom(allianceID)})
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate etag ID")
	}

	etag, err := s.etag.Etag(ctx, etagID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch etag for expiry check")
	}

	alliance, err = s.alliance.Alliance(ctx, allianceID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, errors.Wrap(err, "failed to fetch alliance record from data store")
	}

	exists := err == nil
	now := time.Now()
	if !exists || etag == nil || etag.CachedUntil.Unix() < now.Add(-1*time.Minute).Unix() {
		mods := append(make([]esi.ModifierFunc, 0, 2), s.esi.CacheEtag(ctx, etagID, null.TimeFrom(now.AddDate(0, 0, 14)).Ptr()))
		if etag != nil && etag.Etag != "" {
			mods = append(mods, s.esi.AddIfNoneMatchHeader(ctx, etag.Etag))
		}

		updatedAlliance, err := s.esi.GetAlliance(ctx, allianceID, mods...)
		if err != nil {
			return nil, errors.Wrap(err, "failed to fetch alliance from ESI")
		}

		if updatedAlliance != nil {
			switch exists {
			case true:
				err = s.alliance.UpdateAlliance(ctx, updatedAlliance)
				if err != nil {
					return nil, errors.Wrap(err, "failed to save alliance to data store")
				}
			case false:
				err = s.alliance.CreateAlliance(ctx, updatedAlliance)
				if err != nil {
					return nil, errors.Wrap(err, "failed to save alliance to data store")
				}

			}
			alliance = updatedAlliance
		}

	}

	defer func(x context.Context, alliance *skillz.Alliance) {
		err = s.cache.SetAlliance(ctx, alliance, time.Hour)
		if err != nil {
			s.logger.WithError(err).Error("failed to cache alliance")
		}
	}(context.Background(), alliance)

	return alliance, nil

}
