package clone

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
	skillz.Processor
	Clones(ctx context.Context, user *skillz.User) (*skillz.CharacterCloneMeta, error)
}

type Service struct {
	cache cache.CloneAPI
	etag  etag.API
	esi   esi.CloneAPI

	clones skillz.CloneRepository

	scopes []skillz.Scope
}

var _ API = new(Service)

func New(cache cache.CloneAPI, etag etag.API, esi esi.CloneAPI, clones skillz.CloneRepository) *Service {
	return &Service{
		cache:  cache,
		etag:   etag,
		esi:    esi,
		clones: clones,
		scopes: []skillz.Scope{skillz.ReadImplantsV1, skillz.ReadClonesV1},
	}
}

func (s *Service) Process(ctx context.Context, user *skillz.User) error {

	_, err := s.Clones(ctx, user)
	if err != nil {
		return err
	}

	_, err = s.Implants(ctx, user)
	if err != nil {
		return err
	}

	return nil

}

func (s *Service) Scopes() []skillz.Scope {
	return s.scopes
}

func (s *Service) Clones(ctx context.Context, user *skillz.User) (*skillz.CharacterCloneMeta, error) {

	clones, err := s.cache.CharacterClones(ctx, user.CharacterID)
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}

	if clones != nil {
		return clones, nil
	}

	etagID, etag, err := s.esi.Etag(ctx, esi.GetCharacterClones, &esi.Params{CharacterID: null.Uint64From(user.CharacterID)})
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch tag for expiry check")
	}

	clones, err = s.clones.CharacterCloneMeta(ctx, user.CharacterID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, errors.Wrap(err, "failed to fetch character clones from data store")
	}

	exists := err == nil

	if !exists || etag == nil || etag.CachedUntil.Unix() < time.Now().Add(-1*time.Minute).Unix() {

		mods := append(make([]esi.ModifierFunc, 0, 3), s.esi.CacheEtag(ctx, etagID), s.esi.AddAuthorizationHeader(ctx, user.AccessToken))
		if etag != nil && etag.Etag != "" {
			mods = append(mods, s.esi.AddIfNoneMatchHeader(ctx, etag.Etag))
		}

		updatedClones, err := s.esi.GetCharacterClones(ctx, user.CharacterID, mods...)
		if err != nil {
			return nil, errors.Wrap(err, "failed to fetch character clones from ESI")
		}

		if updatedClones != nil {
			switch exists {
			case true:
				err = s.updateCharacterClones(ctx, updatedClones)
			case false:
				err = s.insertCharacterClones(ctx, updatedClones)
			}

			if err != nil {
				return nil, errors.Wrap(err, "failed to process character clones")
			}

			return updatedClones, s.cache.SetCharacterClones(ctx, user.CharacterID, updatedClones, time.Hour)
		}

	}

	homeLocation, homeLocationErr := s.clones.CharacterDeathClone(ctx, user.CharacterID)
	jumpClones, jumpClonesErr := s.clones.CharacterJumpClones(ctx, user.CharacterID)

	if homeLocationErr == nil && jumpClonesErr == nil {
		clones.HomeLocation = homeLocation
		clones.JumpClones = jumpClones
	}

	return clones, s.cache.SetCharacterClones(ctx, user.CharacterID, clones, time.Hour)

}

func (s *Service) insertCharacterClones(ctx context.Context, clones *skillz.CharacterCloneMeta) error {

	err := s.clones.CreateCharacterCloneMeta(ctx, clones)
	if err != nil {
		return err
	}

	err = s.clones.CreateCharacterDeathClone(ctx, clones.HomeLocation)
	if err != nil {
		return err
	}

	return s.clones.CreateCharacterJumpClones(ctx, clones.JumpClones)

}

func (s *Service) updateCharacterClones(ctx context.Context, clones *skillz.CharacterCloneMeta) error {

	err := s.clones.UpdateCharacterCloneMeta(ctx, clones)
	if err != nil {
		return err
	}

	err = s.clones.UpdateCharacterDeathClone(ctx, clones.HomeLocation)
	if err != nil {
		return err
	}

	err = s.clones.DeleteCharacterJumpClones(ctx, clones.CharacterID)
	if err != nil {
		return err
	}

	return s.clones.CreateCharacterJumpClones(ctx, clones.JumpClones)

}

func (s *Service) Implants(ctx context.Context, user *skillz.User) ([]*skillz.CharacterImplant, error) {

	implants, err := s.cache.CharacterImplants(ctx, user.CharacterID)
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}

	if implants != nil {
		return implants, nil
	}

	etagID, etag, err := s.esi.Etag(ctx, esi.GetCharacterImplants, &esi.Params{CharacterID: null.Uint64From(user.CharacterID)})
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetche tag for expiry check")
	}

	implants, err = s.clones.CharacterImplants(ctx, user.CharacterID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, errors.Wrap(err, "failed to fetch character implants from data store")
	}

	exists := err == nil

	if !exists || etag == nil || etag.CachedUntil.Unix() < time.Now().Add(-1*time.Minute).Unix() {
		mods := append(make([]esi.ModifierFunc, 0, 3), s.esi.CacheEtag(ctx, etagID), s.esi.AddAuthorizationHeader(ctx, user.AccessToken))
		if etag != nil && etag.Etag != "" {
			mods = append(mods, s.esi.AddIfNoneMatchHeader(ctx, etag.Etag))
		}

		implantsOk, err := s.esi.GetCharacterImplants(ctx, user.CharacterID, mods...)
		if err != nil {
			return nil, errors.Wrap(err, "failed to fetch character implants from ESI")
		}

		if implantsOk.Updated {
			implants = implantsOk.Implants

			err = s.clones.DeleteCharacterImplants(ctx, user.CharacterID)
			if err != nil {
				return nil, errors.Wrap(err, "failed to update character implants")
			}

			err = s.clones.CreateCharacterImplants(ctx, implants)
			if err != nil {
				return nil, errors.Wrap(err, "failed to update character implants")
			}
		}
	}

	err = s.cache.SetCharacterImplants(ctx, user.CharacterID, implants, time.Hour)

	return implants, err

}
