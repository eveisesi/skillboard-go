package esi

import (
	"context"
	"fmt"
	"net/http"

	"github.com/eveisesi/skillz"
	"github.com/pkg/errors"
)

type CharacterAPI interface {
	characters
	modifiers
}

type characters interface {
	GetCharacter(ctx context.Context, characterID uint64, mods ...ModifierFunc) (*skillz.Character, error)
	GetCharacterHistory(ctx context.Context, characterID uint64) ([]*skillz.CharacterCorporationHistory, error)
}

func isCharacterValid(character *skillz.Character) bool {
	return character.Name != "" && character.CorporationID > 0
}

func (s *Service) GetCharacter(ctx context.Context, characterID uint64, mods ...ModifierFunc) (*skillz.Character, error) {

	var character = new(skillz.Character)
	var out = new(out)
	out.Data = character
	endpoint := fmt.Sprintf(endpoints[GetCharacter], characterID)
	err := s.request(ctx, http.MethodGet, endpoint, nil, http.StatusOK, out, mods...)
	if err != nil {
		return character, errors.Wrap(err, "failed to execute request to ESI for Character data")
	}

	if !isCharacterValid(character) {
		return nil, errors.New("invalid character returned from ESI")
	}

	if character.ID == 0 {
		character.ID = characterID
	}

	return character, nil

}

func (s *Service) GetCharacterHistory(ctx context.Context, characterID uint64) ([]*skillz.CharacterCorporationHistory, error) {

	var out = new(out)
	var history = make([]*skillz.CharacterCorporationHistory, 0, 256)
	out.Data = &history
	endpoint := fmt.Sprintf("/v2/characters/%d/corporationhistory/", characterID)
	err := s.request(ctx, http.MethodGet, endpoint, nil, http.StatusOK, out)

	for _, record := range history {
		record.CharacterID = characterID
	}

	return history, errors.Wrap(
		err,
		"failed to execute request to ESI for Character Corporation History data",
	)

}
