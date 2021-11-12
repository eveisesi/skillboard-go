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
	"github.com/volatiletech/null"
)

type API interface {
	Character(ctx context.Context, characterID uint64) (*skillz.Character, error)
}

type Service struct {
	cache cache.CharacterAPI
	etag  etag.API
	esi   esi.CharacterAPI

	character skillz.CharacterRepository
}

var _ API = new(Service)

func New(cache cache.CharacterAPI, esi esi.CharacterAPI, etag etag.API, character skillz.CharacterRepository) *Service {
	return &Service{
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

	etagID, err := esi.Resolvers[esi.GetCharacter](&esi.Params{CharacterID: null.Uint64From(characterID)})
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate etag ID")
	}

	etag, err := s.etag.Etag(ctx, etagID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetche tag for expiry check")
	}

	character, err = s.character.Character(ctx, characterID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, errors.Wrap(err, "failed to fetch character record from data store")
	}

	exists := err == nil

	if err == nil && etag != nil && etag.CachedUntil.Unix() > time.Now().Add(time.Minute).Unix() {
		return character, nil
	}

	mods := append(make([]esi.ModifierFunc, 0, 2), s.esi.CacheEtag(ctx, etagID))
	if etag != nil && etag.Etag != "" {
		mods = append(mods, s.esi.AddIfNoneMatchHeader(ctx, etag.Etag))
	}

	character, err = s.esi.GetCharacter(ctx, characterID, mods...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch character from ESI")
	}

	switch exists {
	case true:
		err = s.character.UpdateCharacter(ctx, character)
		if err != nil {
			return nil, errors.Wrap(err, "failed to save character to data store")
		}
	case false:
		err = s.character.CreateCharacter(ctx, character)
		if err != nil {
			return nil, errors.Wrap(err, "failed to save character to data store")
		}

	}

	return character, s.cache.SetCharacter(ctx, character, time.Hour)

}
