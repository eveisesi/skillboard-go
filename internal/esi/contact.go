package esi

import (
	"context"
	"fmt"
	"net/http"

	"github.com/eveisesi/skillz"
	"github.com/pkg/errors"
)

type ContactAPI interface {
	contacts
	etags
	modifiers
}

type contacts interface {
	GetCharacterContacts(ctx context.Context, characterID uint64, mods ...ModifierFunc) ([]*skillz.CharacterContact, error)
}

func (s *Service) GetCharacterContacts(ctx context.Context, characterID uint64, mods ...ModifierFunc) ([]*skillz.CharacterContact, error) {

	var contacts = make([]*skillz.CharacterContact, 0)
	var out = new(out)
	out.Data = &contacts
	endpoint := fmt.Sprintf(endpoints[GetCharacterContacts], characterID)
	err := s.request(ctx, http.MethodGet, endpoint, nil, http.StatusOK, out, mods...)

	if err == nil {
		for _, contact := range contacts {
			contact.CharacterID = characterID
		}
	}

	if out.Status == http.StatusNotModified {
		return nil, nil
	}

	return contacts, errors.Wrap(err, "failed to execute request to ESI for Character data")

}
