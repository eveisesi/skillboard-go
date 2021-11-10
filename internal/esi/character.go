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
}

type characters interface {
	GetCharacter(ctx context.Context, character *skillz.Character) (*skillz.Character, error)
	GetCharacterHistory(ctx context.Context, character *skillz.Character) ([]*skillz.CharacterCorporationHistory, error)
}

func isCharacterValid(character *skillz.Character) bool {
	return character.Name != "" && character.CorporationID > 0
}

func (s *Service) GetCharacter(ctx context.Context, character *skillz.Character) (*skillz.Character, error) {

	var out = new(out)
	out.Data = character
	endpoint := fmt.Sprintf("/v5/characters/%d/", character.ID)
	err := s.request(ctx, http.MethodGet, endpoint, nil, http.StatusOK, out)
	if err != nil {
		return character, errors.Wrap(err, "failed to execute request to ESI for Character data")
	}

	return character, nil

}

func (s *Service) GetCharacterHistory(ctx context.Context, character *skillz.Character) ([]*skillz.CharacterCorporationHistory, error) {

	var out = new(out)
	var history = make([]*skillz.CharacterCorporationHistory, 0, 256)
	out.Data = &history
	endpoint := fmt.Sprintf("/v2/characters/%d/corporationhistory/", character.ID)

	return history, errors.Wrap(
		s.request(ctx, http.MethodGet, endpoint, nil, http.StatusOK, out),
		"failed to execute request to ESI for Character Corporation History data",
	)

}
