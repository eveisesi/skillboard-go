package clone

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/eveisesi/skillz"
	"github.com/eveisesi/skillz/internal/cache"
	"github.com/eveisesi/skillz/internal/esi"
	"github.com/eveisesi/skillz/internal/etag"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"github.com/volatiletech/null"
)

type Service struct {
	cache cache.CloneAPI
	etag  etag.API
	esi   esi.CloneAPI

	clones skillz.CloneRepository

	scopes []skillz.Scope
}

var _ skillz.Processor = new(Service)

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

	fmt.Println("received user ::", user.CharacterID)
	return nil

}

func (s *Service) Scopes() []skillz.Scope {
	return s.scopes
}

func (s *Service) Clones(ctx context.Context, characterID uint64) (*skillz.CharacterCloneMeta, error) {

	clones, err := s.cache.CharacterClones(ctx, characterID)
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}

	if clones != nil {
		return clones, nil
	}

	etagID, etag, err := s.esi.Etag(ctx, esi.GetCharacterClones, &esi.Params{CharacterID: null.Uint64From(characterID)})
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch tag for expiry check")
	}

	clones, err = s.clones.CharacterCloneMeta(ctx, characterID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, errors.Wrap(err, "failed to fetch character clones from data store")
	}

	exists := err == nil

	if exists && etag != nil && etag.CachedUntil.Unix() > time.Now().Add(time.Minute).Unix() {

		homeLocation, homeLocationErr := s.clones.CharacterDeathClone(ctx, characterID)
		jumpClones, jumpClonesErr := s.clones.CharacterJumpClones(ctx, characterID)

		if homeLocationErr == nil && jumpClonesErr == nil {
			clones.HomeLocation = homeLocation
			clones.JumpClones = jumpClones
			return clones, nil
		}

	}

	mods := append(make([]esi.ModifierFunc, 0, 2), s.esi.CacheEtag(ctx, etagID))
	if etag != nil && etag.Etag != "" {
		mods = append(mods, s.esi.AddIfNoneMatchHeader(ctx, etag.Etag))
	}

	updatedClones, err := s.esi.GetCharacterClones(ctx, characterID, mods...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch character clones from ESI")
	}

	if updatedClones == nil {
		return clones, nil
	}

	switch exists {
	case true:
		err = s.updateCharacterClones(ctx, updatedClones)
	case false:
		err = s.insertCharacterClones(ctx, updatedClones)
	}

	if err != nil {
		return nil, errors.Wrap(err, "failed to process character clones")
	}

	return updatedClones, s.cache.SetCharacterClones(ctx, characterID, updatedClones, time.Hour)

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

func (s *Service) Implants(ctx context.Context, characterID uint64) ([]*skillz.CharacterImplant, error) {

	implants, err := s.cache.CharacterImplants(ctx, characterID)
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}

	if implants != nil {
		return implants, nil
	}

	etagID, etag, err := s.esi.Etag(ctx, esi.GetCharacterImplants, &esi.Params{CharacterID: null.Uint64From(characterID)})
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetche tag for expiry check")
	}

	implants, err = s.clones.CharacterImplants(ctx, characterID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, errors.Wrap(err, "failed to fetch character implants from data store")
	}

	exists := err == nil

	if exists && etag != nil && etag.CachedUntil.Unix() > time.Now().Add(time.Minute).Unix() {
		return implants, nil
	}

	mods := append(make([]esi.ModifierFunc, 0, 2), s.esi.CacheEtag(ctx, etagID))
	if etag != nil && etag.Etag != "" {
		mods = append(mods, s.esi.AddIfNoneMatchHeader(ctx, etag.Etag))
	}

	implants, err = s.esi.GetCharacterImplants(ctx, characterID, mods...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch character implants from ESI")
	}

	err = s.clones.DeleteCharacterImplants(ctx, characterID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update character implants")
	}

	err = s.clones.CreateCharacterImplants(ctx, implants)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update character implants")
	}

	err = s.cache.SetCharacterImplants(ctx, characterID, implants, time.Hour)

	return implants, err

}
