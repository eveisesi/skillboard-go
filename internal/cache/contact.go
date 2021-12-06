package cache

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/eveisesi/skillz"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

type ContactAPI interface {
	CharacterContacts(ctx context.Context, characterID uint64) ([]*skillz.CharacterContact, error)
	SetCharacterContacts(ctx context.Context, characterID uint64, contacts []*skillz.CharacterContact, expires time.Duration) error
}

const (
	characterContactsKeyPrefix = "character::contacts"
)

func (s *Service) CharacterContacts(ctx context.Context, characterID uint64) ([]*skillz.CharacterContact, error) {

	key := generateKey(characterContactsKeyPrefix, strconv.FormatUint(characterID, 10))
	results, err := s.redis.SMembers(ctx, key).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, errors.Wrapf(err, errorFFormat, contactAPI, "CharacterContact", "failed to fetch results from cache")
	}

	if errors.Is(err, redis.Nil) || len(results) == 0 {
		return nil, nil
	}

	contacts := make([]*skillz.CharacterContact, 0, len(results))
	for _, result := range results {
		var contact = new(skillz.CharacterContact)
		err = json.Unmarshal([]byte(result), contact)
		if err != nil {
			return contacts, errors.Wrapf(err, errorFFormat, contactAPI, "CharacterContact", "failed to decode json to structure")
		}

		contacts = append(contacts, contact)
	}

	return contacts, nil

}

func (s *Service) SetCharacterContacts(ctx context.Context, characterID uint64, contacts []*skillz.CharacterContact, expires time.Duration) error {

	members := make([]interface{}, 0, len(contacts))
	for _, contact := range contacts {
		data, err := json.Marshal(contact)
		if err != nil {
			return errors.Wrapf(err, errorFFormat, contactAPI, "SetCharacterContacts", "failed to encode struct as json")
		}

		members = append(members, string(data))
	}

	if len(members) == 0 {
		return nil
	}
	key := generateKey(characterContactsKeyPrefix, strconv.FormatUint(characterID, 10))
	err := s.redis.SAdd(ctx, key, members...).Err()
	if err != nil {
		return errors.Wrapf(err, errorFFormat, contactAPI, "SetCharacterContacts", "failed to write cache")
	}

	err = s.redis.Expire(ctx, key, expires).Err()

	return errors.Wrapf(err, errorFFormat, contactAPI, "SetCharacterContacts", "failed to set expiry on set")

}
