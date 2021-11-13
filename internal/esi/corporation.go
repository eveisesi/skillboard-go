package esi

import (
	"context"
	"fmt"
	"net/http"

	"github.com/eveisesi/skillz"
	"github.com/pkg/errors"
)

type CorporationAPI interface {
	corporations
	modifiers
}

type corporations interface {
	GetCorporation(ctx context.Context, corporationID uint, mods ...ModifierFunc) (*skillz.Corporation, error)
}

func (s *Service) GetCorporation(ctx context.Context, corporationID uint, mods ...ModifierFunc) (*skillz.Corporation, error) {

	var corporation = new(skillz.Corporation)
	var out = new(out)
	out.Data = corporation
	endpoint := fmt.Sprintf(endpoints[GetCorporation], corporationID)
	err := s.request(ctx, http.MethodGet, endpoint, nil, http.StatusOK, out, mods...)
	if err != nil {
		return corporation, errors.Wrap(err, "failed to execute request to ESI for Corporation Data")
	}

	if out.Status == http.StatusNotModified {
		return nil, nil
	}

	if corporation.ID == 0 {
		corporation.ID = corporationID
	}

	return corporation, nil

}

func (s *Service) GetCorporationAllianceHistory(ctx context.Context, corporationID uint, mods ...ModifierFunc) ([]*skillz.CorporationAllianceHistory, error) {

	var history = make([]*skillz.CorporationAllianceHistory, 0, 128)
	var out = &out{Data: history}
	endpoint := fmt.Sprintf(endpoints[GetCorporationAllianceHistory], corporationID)
	err := s.request(ctx, http.MethodGet, endpoint, nil, http.StatusOK, out, mods...)
	if err != nil {
		return history, errors.Wrap(err, "failed to execute request to ESI for Corporation Data")
	}

	if out.Status == http.StatusNotModified {
		return nil, nil
	}

	return history, nil

}
