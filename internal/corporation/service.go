package corporation

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
	"github.com/volatiletech/null"
)

type API interface {
	Corporation(ctx context.Context, corporation uint) (*skillz.Corporation, error)
}

type Service struct {
	cache cache.CorporationAPI
	esi   esi.CorporationAPI
	etag  etag.API

	corporation skillz.CorporationRepository
}

var _ API = new(Service)

func New(cache cache.CorporationAPI, esi esi.CorporationAPI, etag etag.API, corporation skillz.CorporationRepository) *Service {
	return &Service{
		cache:       cache,
		esi:         esi,
		etag:        etag,
		corporation: corporation,
	}
}

func (s *Service) Corporation(ctx context.Context, corporationID uint) (*skillz.Corporation, error) {

	corporation, err := s.cache.Corporation(ctx, corporationID)
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}

	if corporation != nil {
		return corporation, nil
	}

	etagID, err := esi.Resolvers[esi.GetCorporation](&esi.Params{CorporationID: null.UintFrom(corporationID)})
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate etag ID")
	}

	etag, err := s.etag.Etag(ctx, etagID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch etag for expiry check")
	}

	corporation, err = s.corporation.Corporation(ctx, corporationID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, errors.Wrap(err, "failed to fetch corporation record from data store")
	}

	exists := err == nil
	var n = time.Now()
	if !exists || etag == nil || etag.CachedUntil.Unix() < n.Add(-1*time.Minute).Unix() {
		mods := append(make([]esi.ModifierFunc, 0, 2), s.esi.CacheEtag(ctx, etagID, null.TimeFrom(n.AddDate(0, 0, 14)).Ptr()))
		if etag != nil && etag.Etag != "" {
			mods = append(mods, s.esi.AddIfNoneMatchHeader(ctx, etag.Etag))
		}

		updateCorporation, err := s.esi.GetCorporation(ctx, corporationID, mods...)
		if err != nil {
			return nil, errors.Wrap(err, "failed to fetch corporation from ESI")
		}

		if updateCorporation != nil {
			switch exists {
			case true:
				err = s.corporation.UpdateCorporation(ctx, updateCorporation)
				if err != nil {
					return nil, errors.Wrap(err, "failed to save corporation to data store")
				}
			case false:
				err = s.corporation.CreateCorporation(ctx, updateCorporation)
				if err != nil {
					return nil, errors.Wrap(err, "failed to save corporation to data store")
				}

			}
			corporation = updateCorporation
		}

	}

	return corporation, s.cache.SetCorporation(ctx, corporation, time.Hour)

}
