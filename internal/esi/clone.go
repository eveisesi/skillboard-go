package esi

import (
	"context"
	"fmt"
	"net/http"

	"github.com/eveisesi/skillz"
	"github.com/pkg/errors"
)

type CloneAPI interface {
	clones
	etags
	modifiers
}

type clones interface {
	GetCharacterClones(ctx context.Context, characterID uint64, mods ...ModifierFunc) (*skillz.CharacterCloneMeta, error)
	GetCharacterImplants(ctx context.Context, characterID uint64, mods ...ModifierFunc) ([]*skillz.CharacterImplant, error)
}

func (s *Service) GetCharacterClones(ctx context.Context, characterID uint64, mods ...ModifierFunc) (*skillz.CharacterCloneMeta, error) {

	var clones = new(skillz.CharacterCloneMeta)
	var out = new(out)
	out.Data = clones
	endpoint := fmt.Sprintf(endpoints[GetCharacterClones], characterID)
	err := s.request(ctx, http.MethodGet, endpoint, nil, http.StatusOK, out, mods...)

	if clones != nil {
		clones.CharacterID = characterID
		clones.HomeLocation.CharacterID = characterID
		for _, clone := range clones.JumpClones {
			clone.CharacterID = characterID
		}
	}

	if out.Status == http.StatusNotModified {
		return nil, nil
	}

	return clones, errors.Wrap(err, "failed to execute request to ESI for Character data")

}

func (s *Service) GetCharacterImplants(ctx context.Context, characterID uint64, mods ...ModifierFunc) ([]*skillz.CharacterImplant, error) {

	// var implants = make([]*skillz.CharacterImplant, 0, 10)
	var implantIDs = make([]uint, 0, 10)
	var out = new(out)
	out.Data = &implantIDs
	endpoint := fmt.Sprintf(endpoints[GetCharacterClones], characterID)
	err := s.request(ctx, http.MethodGet, endpoint, nil, http.StatusOK, out, mods...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to exec request to ESI for character implants")
	}

	if out.Status == http.StatusNotModified {
		return nil, nil
	}

	var implants = make([]*skillz.CharacterImplant, 0, 10)

	for _, id := range implantIDs {
		implants = append(implants, &skillz.CharacterImplant{
			CharacterID: characterID,
			ImplantID:   id,
		})
	}

	return implants, nil

}
