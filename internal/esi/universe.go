package esi

import (
	"context"
	"fmt"
	"net/http"

	"github.com/eveisesi/skillz"
	"github.com/pkg/errors"
)

type UniverseAPI interface {
	universe
	modifiers
	etags
}

type universe interface {
	GetBloodlines(ctx context.Context, mods ...ModifierFunc) ([]*skillz.Bloodline, error)
	GetRaces(ctx context.Context, mods ...ModifierFunc) ([]*skillz.Race, error)

	GetCategories(ctx context.Context, mods ...ModifierFunc) ([]uint, error)
	GetCategory(ctx context.Context, categoryID uint, mods ...ModifierFunc) (*skillz.Category, error)
	GetGroups(ctx context.Context, mods ...ModifierFunc) ([]uint, error)
	GetGroup(ctx context.Context, groupID uint, mods ...ModifierFunc) (*skillz.Group, error)
	GetType(ctx context.Context, typeID uint, mods ...ModifierFunc) (*skillz.Type, error)

	// Map
	GetRegions(ctx context.Context, mods ...ModifierFunc) ([]uint, error)
	GetRegion(ctx context.Context, regionID uint, mods ...ModifierFunc) (*skillz.Region, error)
	GetConstellations(ctx context.Context, mods ...ModifierFunc) ([]uint, error)
	GetConstellation(ctx context.Context, constellationID uint, mods ...ModifierFunc) (*skillz.Constellation, error)
	GetSolarSystem(ctx context.Context, solarSystemID uint, mods ...ModifierFunc) (*skillz.SolarSystem, error)
	GetStructure(ctx context.Context, structureID uint64, mods ...ModifierFunc) (*skillz.Structure, error)
	GetStation(ctx context.Context, stationID uint, mods ...ModifierFunc) (*skillz.Station, error)
}

func (s *Service) GetBloodlines(ctx context.Context, mods ...ModifierFunc) ([]*skillz.Bloodline, error) {

	var bloodlines = make([]*skillz.Bloodline, 0, 32)
	var out = &out{Data: &bloodlines}
	endpoint := endpoints[GetBloodlines]
	err := s.request(ctx, http.MethodGet, endpoint, nil, http.StatusOK, out, mods...)
	if err != nil {
		return bloodlines, errors.Wrap(err, "encounted error requesting bloodlines from ESI")
	}

	if out.Status == http.StatusNotModified {
		return nil, nil
	}

	return bloodlines, nil

}

func (s *Service) GetRaces(ctx context.Context, mods ...ModifierFunc) ([]*skillz.Race, error) {

	var races = make([]*skillz.Race, 0, 32)
	var out = &out{Data: &races}
	endpoint := endpoints[GetRaces]
	err := s.request(ctx, http.MethodGet, endpoint, nil, http.StatusOK, out, mods...)
	if err != nil {
		return races, errors.Wrap(err, "encounted error requesting races from ESI")
	}

	if out.Status == http.StatusNotModified {
		return nil, nil
	}

	return races, nil

}

func (s *Service) GetCategories(ctx context.Context, mods ...ModifierFunc) ([]uint, error) {

	var categories = make([]uint, 0)
	var out = &out{Data: &categories}
	endpoint := endpoints[GetCategories]
	err := s.request(ctx, http.MethodGet, endpoint, nil, http.StatusOK, out, mods...)
	if err != nil {
		return categories, errors.Wrap(err, "encounted error requesting categories from ESI")
	}

	if out.Status == http.StatusNotModified {
		return nil, nil
	}

	return categories, nil

}

func (s *Service) GetCategory(ctx context.Context, categoryID uint, mods ...ModifierFunc) (*skillz.Category, error) {

	var category = new(skillz.Category)
	var out = &out{Data: category}
	endpoint := fmt.Sprintf(endpoints[GetCategory], categoryID)
	err := s.request(ctx, http.MethodGet, endpoint, nil, http.StatusOK, out, mods...)
	if err != nil {
		return category, errors.Wrap(err, "encounted error requesting category from ESI")
	}

	if out.Status == http.StatusNotModified {
		return nil, nil
	}

	category.ID = categoryID

	return category, nil

}

func (s *Service) GetGroups(ctx context.Context, mods ...ModifierFunc) ([]uint, error) {

	var groups = make([]uint, 0)
	var out = &out{Data: &groups}
	endpoint := endpoints[GetGroup]
	err := s.request(ctx, http.MethodGet, endpoint, nil, http.StatusOK, out, mods...)
	if err != nil {
		return groups, errors.Wrap(err, "encounted error requesting groups from ESI")
	}

	if out.Status == http.StatusNotModified {
		return nil, nil
	}

	return groups, nil

}

func (s *Service) GetGroup(ctx context.Context, groupID uint, mods ...ModifierFunc) (*skillz.Group, error) {

	var group = new(skillz.Group)
	var out = &out{Data: group}
	endpoint := fmt.Sprintf(endpoints[GetGroup], groupID)
	err := s.request(ctx, http.MethodGet, endpoint, nil, http.StatusOK, out, mods...)
	if err != nil {
		return group, errors.Wrap(err, "encounted error requesting group from ESI")
	}

	if out.Status == http.StatusNotModified {
		return nil, nil
	}

	group.ID = groupID

	return group, nil

}

func (s *Service) GetType(ctx context.Context, typeID uint, mods ...ModifierFunc) (*skillz.Type, error) {

	var item = new(skillz.Type)
	var out = &out{Data: item}
	endpoint := fmt.Sprintf(endpoints[GetType], typeID)
	err := s.request(ctx, http.MethodGet, endpoint, nil, http.StatusOK, out, mods...)
	if err != nil {
		return item, errors.Wrap(err, "encounted error requesting item from ESI")
	}

	if out.Status == http.StatusNotModified {
		return nil, nil
	}

	item.ID = typeID
	for _, attribute := range item.Attributes {
		attribute.TypeID = typeID
	}

	return item, nil

}

// #### Map Data ####

func (s *Service) GetRegions(ctx context.Context, mods ...ModifierFunc) ([]uint, error) {

	var regions = make([]uint, 0)
	var out = &out{Data: &regions}
	endpoint := endpoints[GetRegion]
	err := s.request(ctx, http.MethodGet, endpoint, nil, http.StatusOK, out, mods...)
	if err != nil {
		return regions, errors.Wrap(err, "encounted error requesting regions from ESI")
	}

	if out.Status == http.StatusNotModified {
		return nil, nil
	}

	return regions, nil

}

func (s *Service) GetRegion(ctx context.Context, regionID uint, mods ...ModifierFunc) (*skillz.Region, error) {

	var region = new(skillz.Region)
	var out = &out{Data: region}
	endpoint := fmt.Sprintf(endpoints[GetRegion], regionID)
	err := s.request(ctx, http.MethodGet, endpoint, nil, http.StatusOK, out, mods...)
	if err != nil {
		return region, errors.Wrap(err, "encounted error requesting region from ESI")
	}

	if out.Status == http.StatusNotModified {
		return nil, nil
	}

	region.ID = regionID

	return region, nil

}

func (s *Service) GetConstellations(ctx context.Context, mods ...ModifierFunc) ([]uint, error) {

	var constellations = make([]uint, 0)
	var out = &out{Data: &constellations}
	endpoint := endpoints[GetConstellation]
	err := s.request(ctx, http.MethodGet, endpoint, nil, http.StatusOK, out, mods...)
	if err != nil {
		return constellations, errors.Wrap(err, "encounted error requesting constellations from ESI")
	}

	if out.Status == http.StatusNotModified {
		return nil, nil
	}

	return constellations, nil

}

func (s *Service) GetConstellation(ctx context.Context, constellationID uint, mods ...ModifierFunc) (*skillz.Constellation, error) {

	var constellation = new(skillz.Constellation)
	var out = &out{Data: constellation}
	endpoint := fmt.Sprintf(endpoints[GetConstellation], constellationID)
	err := s.request(ctx, http.MethodGet, endpoint, nil, http.StatusOK, out, mods...)
	if err != nil {
		return constellation, errors.Wrap(err, "encounted error requesting constellation from ESI")
	}

	if out.Status == http.StatusNotModified {
		return nil, nil
	}

	constellation.ID = constellationID

	return constellation, nil

}

func (s *Service) GetSolarSystem(ctx context.Context, solarSystemID uint, mods ...ModifierFunc) (*skillz.SolarSystem, error) {

	var system = new(skillz.SolarSystem)
	var out = &out{Data: system}
	endpoint := fmt.Sprintf(endpoints[GetSolarSystem], solarSystemID)
	err := s.request(ctx, http.MethodGet, endpoint, nil, http.StatusOK, out, mods...)
	if err != nil {
		return system, errors.Wrap(err, "encounted error requesting solar system from ESI")
	}

	if out.Status == http.StatusNotModified {
		return nil, nil
	}

	system.ID = solarSystemID

	return system, nil

}

func (s *Service) GetStructure(ctx context.Context, structureID uint64, mods ...ModifierFunc) (*skillz.Structure, error) {

	var structure = new(skillz.Structure)
	var out = &out{Data: structure}
	endpoint := fmt.Sprintf(endpoints[GetStructure], structureID)
	err := s.request(ctx, http.MethodGet, endpoint, nil, http.StatusOK, out, mods...)
	if err != nil {
		return structure, errors.Wrap(err, "encounted error requesting structure from ESI")
	}

	if out.Status == http.StatusNotModified {
		return nil, nil
	}

	structure.ID = structureID

	return structure, nil

}

func (s *Service) GetStation(ctx context.Context, stationID uint, mods ...ModifierFunc) (*skillz.Station, error) {

	var station = new(skillz.Station)
	var out = &out{Data: station}
	endpoint := fmt.Sprintf(endpoints[GetStation], stationID)
	err := s.request(ctx, http.MethodGet, endpoint, nil, http.StatusOK, out, mods...)
	if err != nil {
		return station, errors.Wrap(err, "encounted error requesting station from ESI")
	}

	if out.Status == http.StatusNotModified {
		return nil, nil
	}

	station.ID = stationID

	return station, nil

}
