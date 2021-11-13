package esi

import (
	"context"
	"fmt"
	"net/http"

	"github.com/eveisesi/skillz"
	"github.com/pkg/errors"
)

type AllianceAPI interface {
	alliance
	modifiers
}

type alliance interface {
	GetAlliance(ctx context.Context, allianceID uint, mods ...ModifierFunc) (*skillz.Alliance, error)
}

func (s *Service) GetAlliance(ctx context.Context, allianceID uint, mods ...ModifierFunc) (*skillz.Alliance, error) {

	var alliance = new(skillz.Alliance)
	var out = &out{Data: alliance}
	endpoint := fmt.Sprintf(endpoints[GetAlliance], allianceID)
	err := s.request(ctx, http.MethodGet, endpoint, nil, http.StatusOK, out, mods...)
	if err != nil {
		return alliance, errors.Wrap(err, "failed to execute request to ESI for Character data")
	}

	if out.Status == http.StatusNotModified {
		return nil, nil
	}

	if alliance.ID == 0 {
		alliance.ID = allianceID
	}

	return alliance, nil

}
