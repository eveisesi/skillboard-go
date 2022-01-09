package esi

import (
	"context"
	"fmt"
	"net/http"

	"github.com/eveisesi/skillz"
	"github.com/pkg/errors"
	"github.com/volatiletech/null"
)

type CloneAPI interface {
	clones
	etags
	modifiers
}

type clones interface {
	GetCharacterClones(ctx context.Context, characterID uint64, mods ...ModifierFunc) (*CharacterClonesOK, error)
	GetCharacterImplants(ctx context.Context, characterID uint64, mods ...ModifierFunc) (*CharacterImplantsOk, error)
}

type CharacterClonesOK struct {
	HomeLocation struct {
		LocationID   uint64 `json:"location_id"`
		LocationType string `json:"location_type"`
	} `json:"home_location"`
	JumpClones []struct {
		JumpCloneID  uint   `json:"jump_clone_id"`
		LocationID   uint64 `json:"location_id"`
		LocationType string `json:"location_type"`
		Implants     []uint `json:"implants"`
	} `json:"jump_clones"`
	LastCloneJumpDate     null.Time `json:"last_clone_jump_date"`
	LastStationChangeDate null.Time `json:"last_station_change_date"`
}

func (s *Service) GetCharacterClones(ctx context.Context, characterID uint64, mods ...ModifierFunc) (*CharacterClonesOK, error) {

	var clones = new(CharacterClonesOK)
	var out = new(out)
	out.Data = clones
	endpoint := fmt.Sprintf(endpoints[GetCharacterClones], characterID)
	err := s.request(ctx, http.MethodGet, endpoint, nil, http.StatusOK, out, mods...)

	if out.Status == http.StatusNotModified {
		return nil, nil
	}

	return clones, errors.Wrap(err, "failed to execute request to ESI for Character data")

}

type CharacterImplantsOk struct {
	Implants []*skillz.CharacterImplant
	Updated  bool
}

func (s *Service) GetCharacterImplants(ctx context.Context, characterID uint64, mods ...ModifierFunc) (*CharacterImplantsOk, error) {

	// var implants = make([]*skillz.CharacterImplant, 0, 10)
	var implantIDs = make([]uint, 0, 10)
	var out = new(out)
	out.Data = &implantIDs
	endpoint := fmt.Sprintf(endpoints[GetCharacterImplants], characterID)
	err := s.request(ctx, http.MethodGet, endpoint, nil, http.StatusOK, out, mods...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to exec request to ESI for character implants")
	}

	var implants = &CharacterImplantsOk{
		Implants: make([]*skillz.CharacterImplant, 0, 10),
		Updated:  out.Status == http.StatusOK,
	}

	if out.Status == http.StatusNotModified {
		return implants, nil
	}

	for _, id := range implantIDs {
		implants.Implants = append(implants.Implants, &skillz.CharacterImplant{
			CharacterID: characterID,
			ImplantID:   id,
		})
	}

	return implants, nil

}
