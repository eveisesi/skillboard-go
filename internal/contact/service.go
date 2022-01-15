package contact

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/eveisesi/skillz"
	"github.com/eveisesi/skillz/internal/alliance"
	"github.com/eveisesi/skillz/internal/cache"
	"github.com/eveisesi/skillz/internal/character"
	"github.com/eveisesi/skillz/internal/corporation"
	"github.com/eveisesi/skillz/internal/esi"
	"github.com/eveisesi/skillz/internal/etag"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/volatiletech/null"
)

type API interface {
	skillz.Processor
	Contacts(ctx context.Context, characterID uint64) ([]*skillz.CharacterContact, error)
}

type Service struct {
	logger      *logrus.Logger
	cache       cache.ContactAPI
	etag        etag.API
	esi         esi.ContactAPI
	character   character.API
	corporation corporation.API
	alliance    alliance.API

	contacts skillz.ContactRepository

	scopes []skillz.Scope
}

var _ API = (*Service)(nil)

func New(logger *logrus.Logger, cache cache.ContactAPI, etag etag.API, esi esi.ContactAPI, character character.API, corporation corporation.API, alliance alliance.API, contacts skillz.ContactRepository) *Service {
	return &Service{
		logger:      logger,
		cache:       cache,
		etag:        etag,
		esi:         esi,
		character:   character,
		corporation: corporation,
		alliance:    alliance,
		contacts:    contacts,
		scopes:      []skillz.Scope{skillz.ReadContactsV1},
	}
}

func (s *Service) Process(ctx context.Context, user *skillz.User) error {

	var err error
	var funcs = []func(context.Context, *skillz.User) error{s.updateContacts}

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

func (s *Service) Contacts(ctx context.Context, characterID uint64) ([]*skillz.CharacterContact, error) {

	contacts, err := s.cache.CharacterContacts(ctx, characterID)
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}

	if len(contacts) > 0 {
		return contacts, nil
	}

	contacts, err = s.contacts.CharacterContacts(ctx, characterID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, errors.Wrap(err, "failed to fetch character contacts from data store")
	}

	for _, contact := range contacts {

		switch contact.ContactType {
		case skillz.AllianceContactType:
			contact.Alliance, err = s.alliance.Alliance(ctx, contact.ContactID)
		case skillz.CorporationContactType:
			contact.Corporation, err = s.corporation.Corporation(ctx, contact.ContactID)
		}

		if err != nil {
			return nil, err
		}

	}

	defer func() {
		err := s.cache.SetCharacterContacts(ctx, characterID, contacts, time.Hour)
		if err != nil {
			s.logger.WithError(err).Error("failed to cache character standings")
		}
	}()

	return contacts, nil

}

func (s *Service) updateContacts(ctx context.Context, user *skillz.User) error {

	s.logger.WithFields(logrus.Fields{
		"service": "contact",
		"userID":  user.ID.String(),
	}).Info("updating contacts")

	etagID, etag, err := s.esi.Etag(ctx, esi.GetCharacterContacts, &esi.Params{CharacterID: null.Uint64From(user.CharacterID)})
	if err != nil {
		return errors.Wrap(err, "failed to fetch tag for expiry check")
	}

	if etag != nil && etag.CachedUntil.Unix() > time.Now().Unix() {
		return nil
	}

	mods := s.esi.BaseCharacterModifiers(ctx, user, etagID, etag)

	updatedContacts, err := s.esi.GetCharacterContacts(ctx, user.CharacterID, mods...)
	if err != nil {
		return errors.Wrap(err, "failed to fetch character clones from ESI")
	}

	for _, contact := range updatedContacts {
		switch contact.ContactType {
		case skillz.AllianceContactType:
			_, err = s.alliance.Alliance(ctx, contact.ContactID)
			if err != nil {
				fmt.Println(err)
			}
		case skillz.CorporationContactType:
			_, err = s.corporation.Corporation(ctx, contact.ContactID)
			if err != nil {
				fmt.Println(err)
			}

		case skillz.CharacterContactType:
			_, err = s.character.Character(ctx, uint64(contact.ContactID))
			if err != nil {
				fmt.Println(err)
			}
		}
	}

	if updatedContacts != nil {
		err = s.contacts.DeleteCharacterContacts(ctx, user.CharacterID)
		if err != nil {
			return errors.Wrap(err, "failed to remove older contacts")
		}

		err = s.contacts.CreateCharacterContacts(ctx, updatedContacts)
		if err != nil {
			return errors.Wrap(err, "failed to add contacts to data store")
		}

		err = s.cache.SetCharacterContacts(ctx, user.CharacterID, updatedContacts, time.Hour)
		if err != nil {
			return errors.Wrap(err, "failed to cache character clones")
		}

	}

	return nil

}
