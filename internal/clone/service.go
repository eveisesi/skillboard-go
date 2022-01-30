package clone

import (
	"context"
	"database/sql"
	"time"

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
	Implants(ctx context.Context, characterID uint64) ([]*skillz.CharacterImplant, error)
}

type Service struct {
	logger   *logrus.Logger
	cache    cache.CloneAPI
	etag     etag.API
	esi      esi.CloneAPI
	universe universe.API

	clones skillz.CloneRepository
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
	}
}

func (s *Service) Process(ctx context.Context, user *skillz.User) error {

	var err error
	var funcs = []func(context.Context, *skillz.User) error{}
	for _, scope := range user.Scopes {
		switch scope {
		case skillz.ReadImplantsV1:
			funcs = append(funcs, s.updateImplants)
		}
	}

	for _, f := range funcs {
		err = f(ctx, user)
		if err != nil {
			s.logger.WithError(err).Error("processor func returned an error")
		}
	}

	return err

}

func (s *Service) Implants(ctx context.Context, characterID uint64) ([]*skillz.CharacterImplant, error) {

	implants, err := s.cache.CharacterImplants(ctx, characterID)
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}

	if len(implants) > 0 {
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

	s.logger.WithFields(logrus.Fields{
		"service": "clone",
		"userID":  user.ID.String(),
	}).Info("updating implants")

	etagID, _, err := s.esi.Etag(ctx, esi.GetCharacterImplants, &esi.Params{CharacterID: null.Uint64From(user.CharacterID)})
	if err != nil {
		return errors.Wrap(err, "failed to fetch etag for expiry check")
	}

	// if etag != nil && etag.CachedUntil.Unix() > time.Now().Unix() {
	// 	return nil
	// }

	mods := s.esi.BaseCharacterModifiers(ctx, user, etagID, nil)

	implants, err := s.esi.GetCharacterImplants(ctx, user.CharacterID, mods...)
	if err != nil {
		return errors.Wrap(err, "failed to fetch character implants from ESI")
	}

	if implants != nil {
		err = s.clones.DeleteCharacterImplants(ctx, user.CharacterID)
		if err != nil {
			return errors.Wrap(err, "failed to update character implants")
		}

		if len(implants) > 0 {

			for _, implant := range implants {

				implantType, err := s.universe.Type(ctx, implant.ImplantID)
				if err != nil {
					return err
				}

				attribute := implantType.GetAttribute(skillz.ImplantSlotAttributeID)
				if attribute == nil {
					return errors.Wrap(err, "failed to fetch implant slot attribute from implant type")
				}

				implant.Slot = uint(attribute.Value)

			}

			err = s.clones.CreateCharacterImplants(ctx, implants)
			if err != nil {
				return errors.Wrap(err, "failed to update character implants")
			}

		}

	}

	return nil

}
