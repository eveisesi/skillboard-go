package clone

import (
	"context"
	"database/sql"
	"time"

	"github.com/davecgh/go-spew/spew"
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
	Implants(ctx context.Context, user *skillz.User) ([]*skillz.CharacterImplant, error)
}

type Service struct {
	cache cache.CloneAPI
	etag  etag.API
	esi   esi.CloneAPI

	clones skillz.CloneRepository

	scopes []skillz.Scope
}

var _ API = (*Service)(nil)

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

	var err error
	var funcs = []func(context.Context, *skillz.User) error{s.updateClones, s.updateImplants}

	for _, f := range funcs {
		err = f(ctx, user)
		if err != nil {
			break
		}
	}

	return err

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

	clones, err = s.clones.CharacterCloneMeta(ctx, user.CharacterID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, errors.Wrap(err, "failed to fetch character clones from data store")
	}

	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	homeLocation, homeLocationErr := s.clones.CharacterDeathClone(ctx, user.CharacterID)
	jumpClones, jumpClonesErr := s.clones.CharacterJumpClones(ctx, user.CharacterID)

	if homeLocationErr == nil && jumpClonesErr == nil {
		clones.HomeLocation = homeLocation
		clones.JumpClones = jumpClones
	}

	return clones, s.cache.SetCharacterClones(ctx, user.CharacterID, clones, time.Hour)

}

func (s *Service) updateClones(ctx context.Context, user *skillz.User) error {

	etagID, etag, err := s.esi.Etag(ctx, esi.GetCharacterClones, &esi.Params{CharacterID: null.Uint64From(user.CharacterID)})
	if err != nil {
		return errors.Wrap(err, "failed to fetch tag for expiry check")
	}

	if etag != nil && etag.CachedUntil.Unix() > time.Now().Unix() {
		return nil
	}

	mods := s.esi.BaseCharacterModifiers(ctx, user, etagID, etag)

	updatedClones, err := s.esi.GetCharacterClones(ctx, user.CharacterID, mods...)
	if err != nil {
		return errors.Wrap(err, "failed to fetch character clones from ESI")
	}

	spew.Dump(updatedClones)

	if updatedClones != nil {
		err = s.insertCharacterClones(ctx, updatedClones)
		if err != nil {
			return errors.Wrap(err, "failed to process character clones")
		}

		err = s.cache.SetCharacterClones(ctx, user.CharacterID, updatedClones, time.Hour)
		if err != nil {
			return errors.Wrap(err, "failed to cache character clones")
		}

	}

	return nil

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

	err = s.clones.DeleteCharacterJumpClones(ctx, clones.CharacterID)
	if err != nil {
		return err
	}

	if len(clones.JumpClones) > 0 {
		err = s.clones.CreateCharacterJumpClones(ctx, clones.JumpClones)
		if err != nil {
			return err
		}
	}

	return nil

}

func (s *Service) Implants(ctx context.Context, user *skillz.User) ([]*skillz.CharacterImplant, error) {

	implants, err := s.cache.CharacterImplants(ctx, user.CharacterID)
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}

	if implants != nil {
		return implants, nil
	}

	implants, err = s.clones.CharacterImplants(ctx, user.CharacterID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, errors.Wrap(err, "failed to fetch character implants from data store")
	}

	return implants, nil

}

func (s *Service) updateImplants(ctx context.Context, user *skillz.User) error {

	etagID, etag, err := s.esi.Etag(ctx, esi.GetCharacterImplants, &esi.Params{CharacterID: null.Uint64From(user.CharacterID)})
	if err != nil {
		return errors.Wrap(err, "failed to fetch etag for expiry check")
	}

	if etag != nil && etag.CachedUntil.Unix() > time.Now().Unix() {
		return nil
	}

	mods := s.esi.BaseCharacterModifiers(ctx, user, etagID, etag)

	implantsOk, err := s.esi.GetCharacterImplants(ctx, user.CharacterID, mods...)
	if err != nil {
		return errors.Wrap(err, "failed to fetch character implants from ESI")
	}

	if implantsOk.Updated {
		implants := implantsOk.Implants

		err = s.clones.DeleteCharacterImplants(ctx, user.CharacterID)
		if err != nil {
			return errors.Wrap(err, "failed to update character implants")
		}

		if len(implants) > 0 {
			err = s.clones.CreateCharacterImplants(ctx, implants)
			if err != nil {
				return errors.Wrap(err, "failed to update character implants")
			}

			err = s.cache.SetCharacterImplants(ctx, user.CharacterID, implants, time.Hour)
			if err != nil {
				return errors.Wrap(err, "failed to cache character implants")
			}
		}

	}

	return nil

}
