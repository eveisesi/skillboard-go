package character

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
	Character(ctx context.Context, characterID uint64) (*skillz.Character, error)
}

type Service struct {
	logger *logrus.Logger
	cache  cache.CharacterAPI
	etag   etag.API
	esi    esi.CharacterAPI

	character skillz.CharacterRepository
}

var _ API = new(Service)

func New(logger *logrus.Logger, cache cache.CharacterAPI, esi esi.CharacterAPI, etag etag.API, character skillz.CharacterRepository) *Service {
	return &Service{
		logger:    logger,
		cache:     cache,
		esi:       esi,
		etag:      etag,
		character: character,
	}
}

func (s *Service) Character(ctx context.Context, characterID uint64) (*skillz.Character, error) {

	character, err := s.cache.Character(ctx, characterID)
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}

	if character != nil {
		return character, nil
	}

	etagID, etag, err := s.esi.Etag(ctx, esi.GetCharacter, &esi.Params{CharacterID: null.Uint64From(characterID)})
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch etag for expiry check")
	}

	character, err = s.character.Character(ctx, characterID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, errors.Wrap(err, "failed to fetch character record from data store")
	}

	exists := err == nil
	now := time.Now()
	if !exists || etag == nil || etag.CachedUntil.Unix() < now.Add(-1*time.Minute).Unix() {
		mods := append(make([]esi.ModifierFunc, 0, 2), s.esi.CacheEtag(ctx, etagID, null.TimeFrom(now.AddDate(0, 0, 3)).Ptr()))
		if etag != nil && etag.Etag != "" {
			mods = append(mods, s.esi.AddIfNoneMatchHeader(ctx, etag.Etag))
		}

		updatedCharacter, err := s.esi.GetCharacter(ctx, characterID, mods...)
		if err != nil {
			return nil, errors.Wrap(err, "failed to fetch character from ESI")
		}

		if updatedCharacter != nil {
			switch exists {
			case true:
				err = s.character.UpdateCharacter(ctx, updatedCharacter)
				if err != nil {
					return nil, errors.Wrap(err, "failed to save character to data store")
				}
			case false:
				err = s.character.CreateCharacter(ctx, updatedCharacter)
				if err != nil {
					return nil, errors.Wrap(err, "failed to save character to data store")
				}

			}

			character = updatedCharacter
		}

	}

	defer func(ctx context.Context, character *skillz.Character) {
		err = s.cache.SetCharacter(ctx, character, time.Hour)
		if err != nil {
			s.logger.WithError(err).Error("failed to cache character")
		}
	}(context.Background(), character)

	return character, nil

}
