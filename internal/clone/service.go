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
	"github.com/eveisesi/skillz/internal/universe"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/volatiletech/null"
)

type API interface {
	skillz.Processor
	Clones(ctx context.Context, characterID uint64) (*skillz.CharacterCloneMeta, error)
	Implants(ctx context.Context, characterID uint64) ([]*skillz.CharacterImplant, error)
}

type Service struct {
	logger   *logrus.Logger
	cache    cache.CloneAPI
	etag     etag.API
	esi      esi.CloneAPI
	universe universe.API

	clones skillz.CloneRepository

	scopes []skillz.Scope
}

var _ API = (*Service)(nil)

func New(logger *logrus.Logger, cache cache.CloneAPI, etag etag.API, esi esi.CloneAPI, universe universe.API, clones skillz.CloneRepository) *Service {
	return &Service{
		logger:   logger,
		cache:    cache,
		etag:     etag,
		esi:      esi,
		universe: universe,

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

		time.Sleep(time.Second)
	}

	return err

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

	clones, err = s.clones.CharacterCloneMeta(ctx, characterID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, errors.Wrap(err, "failed to fetch character clones from data store")
	}

	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	homeLocation, homeLocationErr := s.DeathClone(ctx, characterID)
	jumpClones, jumpClonesErr := s.JumpClone(ctx, characterID)

	if homeLocationErr == nil {
		clones.HomeLocation = homeLocation
	} else {
		s.logger.WithError(jumpClonesErr).Error("failed to fetch death clone")
	}

	if jumpClonesErr == nil {
		clones.JumpClones = jumpClones
	} else {
		s.logger.WithError(jumpClonesErr).Error("failed to fetch jump clones")
	}

	return clones, s.cache.SetCharacterClones(ctx, characterID, clones, time.Hour)

}

func (s *Service) DeathClone(ctx context.Context, characterID uint64) (*skillz.CharacterDeathClone, error) {

	death, err := s.clones.CharacterDeathClone(ctx, characterID)
	if err != nil {
		return nil, err
	}

	switch death.LocationType {
	case skillz.CloneLocationTypeStation:
		death.Station, err = s.universe.Station(ctx, uint(death.LocationID))
	case skillz.CloneLocationTypeStructure:
		death.Structure, err = s.universe.Structure(ctx, death.LocationID)
	}

	return death, err

}

func (s *Service) JumpClone(ctx context.Context, characterID uint64) ([]*skillz.CharacterJumpClone, error) {

	clones, err := s.clones.CharacterJumpClones(ctx, characterID)
	if err != nil {
		return nil, err
	}

	spew.Dump(clones)

	for _, clone := range clones {
		switch clone.LocationType {
		case skillz.CloneLocationTypeStation:
			spew.Dump(clone.LocationID)
			clone.Station, err = s.universe.Station(ctx, uint(clone.LocationID))
			spew.Dump(clone.Station, err)
		case skillz.CloneLocationTypeStructure:
			clone.Structure, err = s.universe.Structure(ctx, clone.LocationID)
		}

		if err != nil {
			return nil, err
		}

		clone.Implants = make([]*skillz.Type, 0, len(clone.ImplantIDs))
		for _, implantID := range clone.ImplantIDs {
			implant, err := s.universe.Type(ctx, uint(implantID))
			if err != nil {
				return nil, err
			}

			clone.Implants = append(clone.Implants, implant)
		}

		clone.ImplantIDs = nil

	}

	return clones, nil
}

func (s *Service) updateClones(ctx context.Context, user *skillz.User) error {

	s.logger.WithFields(logrus.Fields{
		"service": "clone",
		"userID":  user.ID.String(),
	}).Info("updating clones")

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

	if updatedClones != nil {
		clones, err := s.insertCharacterClones(ctx, user.CharacterID, updatedClones)
		if err != nil {
			return errors.Wrap(err, "failed to process character clones")
		}

		err = s.cache.SetCharacterClones(ctx, user.CharacterID, clones, time.Hour)
		if err != nil {
			return errors.Wrap(err, "failed to cache character clones")
		}

		homeClone := clones.HomeLocation
		switch homeClone.LocationType {
		case skillz.CloneLocationTypeStation:
			_, err = s.universe.Station(ctx, uint(homeClone.LocationID))
		case skillz.CloneLocationTypeStructure:
			_, err = s.universe.Structure(ctx, homeClone.LocationID)
		}
		if err != nil {
			return errors.Wrap(err, "failed to resolve home location id")
		}

		jumpClones := clones.JumpClones
		for _, clone := range jumpClones {
			switch clone.LocationType {
			case skillz.CloneLocationTypeStation:
				_, err = s.universe.Station(ctx, uint(clone.LocationID))
			case skillz.CloneLocationTypeStructure:
				_, err = s.universe.Structure(ctx, clone.LocationID)
			}

			if err != nil {
				return errors.Wrap(err, "failed to resolve clone location id")
			}

		}

	}

	return nil

}

func (s *Service) insertCharacterClones(ctx context.Context, characterID uint64, esiClones *esi.CharacterClonesOK) (*skillz.CharacterCloneMeta, error) {

	clones := &skillz.CharacterCloneMeta{
		CharacterID:           characterID,
		LastCloneJumpDate:     esiClones.LastCloneJumpDate,
		LastStationChangeDate: esiClones.LastStationChangeDate,
		HomeLocation: &skillz.CharacterDeathClone{
			CharacterID:  characterID,
			LocationID:   esiClones.HomeLocation.LocationID,
			LocationType: skillz.CloneLocationType(esiClones.HomeLocation.LocationType),
		},
		JumpClones: make([]*skillz.CharacterJumpClone, 0, len(esiClones.JumpClones)),
	}
	for _, esiJC := range esiClones.JumpClones {
		jc := &skillz.CharacterJumpClone{
			CharacterID:  characterID,
			JumpCloneID:  esiJC.JumpCloneID,
			LocationID:   esiJC.LocationID,
			LocationType: skillz.CloneLocationType(esiJC.LocationType),
			Implants:     make([]*skillz.Type, 0, len(esiJC.Implants)),
		}

		for _, implantID := range esiJC.Implants {
			_, err := s.universe.Type(ctx, implantID)
			if err != nil {
				s.logger.WithError(err).Error("failed to fetch implant type")
			}
			jc.ImplantIDs = append(jc.ImplantIDs, uint64(implantID))
		}
		clones.JumpClones = append(clones.JumpClones, jc)
	}

	spew.Dump(clones)

	err := s.clones.CreateCharacterCloneMeta(ctx, clones)
	if err != nil {
		return clones, err
	}

	err = s.clones.CreateCharacterDeathClone(ctx, clones.HomeLocation)
	if err != nil {
		return clones, err
	}

	err = s.clones.DeleteCharacterJumpClones(ctx, clones.CharacterID)
	if err != nil {
		return clones, err
	}

	if len(clones.JumpClones) > 0 {
		err = s.clones.CreateCharacterJumpClones(ctx, clones.JumpClones)
		if err != nil {
			return clones, err
		}
	}

	return clones, nil

}

func (s *Service) Implants(ctx context.Context, characterID uint64) ([]*skillz.CharacterImplant, error) {

	implants, err := s.cache.CharacterImplants(ctx, characterID)
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}

	if implants != nil {
		return implants, nil
	}

	implants, err = s.clones.CharacterImplants(ctx, characterID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, errors.Wrap(err, "failed to fetch character implants from data store")
	}

	for _, implant := range implants {

		implantType, err := s.universe.Type(ctx, implant.ImplantID)
		if err != nil {
			return nil, err
		}

		implant.Type = implantType

	}

	defer func() {
		err = s.cache.SetCharacterImplants(ctx, characterID, implants, time.Hour)
		if err != nil {
			s.logger.WithError(err).Error("failed to cache character implants")
		}
	}()

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

			for _, implant := range implants {
				t, err := s.universe.Type(ctx, implant.ImplantID)
				if err != nil {
					s.logger.WithError(err).Error("failed to fetch implantType")
				}

				implant.Type = t
			}

			err = s.cache.SetCharacterImplants(ctx, user.CharacterID, implants, time.Hour)
			if err != nil {
				return errors.Wrap(err, "failed to cache character implants")
			}
		}

	}

	return nil

}
