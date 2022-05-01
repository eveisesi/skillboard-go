package universe

import (
	"context"
	"database/sql"

	"github.com/eveisesi/skillz"
	"github.com/eveisesi/skillz/internal"
	"github.com/eveisesi/skillz/internal/cache"
	"github.com/eveisesi/skillz/internal/esi"
	"github.com/eveisesi/skillz/internal/mysql"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/volatiletech/null"
)

type API interface {
	Bloodline(ctx context.Context, bloodlineID uint) (*skillz.Bloodline, error)
	Category(ctx context.Context, categoryID uint) (*skillz.Category, error)
	Constellation(ctx context.Context, constellationID uint) (*skillz.Constellation, error)
	Faction(ctx context.Context, id uint) (*skillz.Faction, error)
	Group(ctx context.Context, groupID uint) (*skillz.Group, error)
	GroupsByCategory(ctx context.Context, categoryID uint) ([]*skillz.Group, error)
	Race(ctx context.Context, id uint) (*skillz.Race, error)
	Region(ctx context.Context, regionID uint) (*skillz.Region, error)
	SolarSystem(ctx context.Context, solarSystemID uint) (*skillz.SolarSystem, error)
	Station(ctx context.Context, stationID uint) (*skillz.Station, error)
	Structure(ctx context.Context, structureID uint64) (*skillz.Structure, error)
	Type(ctx context.Context, itemID uint) (*skillz.Type, error)
	SkillTypesHydrated(ctx context.Context) ([]*skillz.Type, error)
	TypeGroupsHydrated(ctx context.Context, categoryID uint) ([]*skillz.Group, error)
	TypeAttributes(ctx context.Context, id uint) ([]*skillz.TypeDogmaAttribute, error)
	TypesByGroup(ctx context.Context, groupID uint) ([]*skillz.Type, error)
}

type Service struct {
	logger *logrus.Logger
	cache  cache.UniverseAPI
	esi    esi.UniverseAPI

	universe skillz.UniverseRepository
}

func New(logger *logrus.Logger, cache cache.UniverseAPI, esi esi.UniverseAPI, universe skillz.UniverseRepository) *Service {
	return &Service{
		logger:   logger,
		cache:    cache,
		esi:      esi,
		universe: universe,
	}
}

func (s *Service) Bloodline(ctx context.Context, bloodlineID uint) (*skillz.Bloodline, error) {

	bloodline, err := s.cache.Bloodline(ctx, bloodlineID)
	if err != nil {
		return nil, err
	}

	if bloodline != nil {
		return bloodline, nil
	}

	bloodline, err = s.universe.Bloodline(ctx, bloodlineID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, errors.Wrap(err, "failed to fetch bloodline from data store")
	}

	err = s.cache.SetBloodline(ctx, bloodline)

	return bloodline, err

}

func (s *Service) Category(ctx context.Context, categoryID uint) (*skillz.Category, error) {

	category, err := s.cache.Category(ctx, categoryID)
	if err != nil {
		return nil, err
	}

	if category != nil {
		return category, nil
	}

	category, err = s.universe.Category(ctx, categoryID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, errors.Wrap(err, "failed to fetch category from data store")
	}

	return category, s.cache.SetCategory(ctx, category)

}

func (s *Service) Constellation(ctx context.Context, constellationID uint) (*skillz.Constellation, error) {

	constellation, err := s.cache.Constellation(ctx, constellationID)
	if err != nil {
		return nil, err
	}

	if constellation != nil {
		return constellation, nil
	}

	constellation, err = s.universe.Constellation(ctx, constellationID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, errors.Wrap(err, "failed to fetch constellation from data store")
	}

	return constellation, s.cache.SetConstellation(ctx, constellation)

}

func (s *Service) Faction(ctx context.Context, id uint) (*skillz.Faction, error) {
	faction, err := s.cache.Faction(ctx, id)
	if err != nil {
		return nil, err
	}

	if faction != nil {
		return faction, nil
	}

	faction, err = s.universe.Faction(ctx, id)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, errors.Wrap(err, "failed to fetch faction from data store")
	}

	return faction, s.cache.SetFaction(ctx, faction)
}

func (s *Service) Group(ctx context.Context, groupID uint) (*skillz.Group, error) {

	group, err := s.cache.Group(ctx, groupID)
	if err != nil {
		return nil, err
	}

	if group != nil {
		return group, nil
	}

	group, err = s.universe.Group(ctx, groupID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, errors.Wrap(err, "failed to fetch group from data store")
	}

	return group, s.cache.SetGroup(ctx, group)

}

func (s *Service) GroupsByCategory(ctx context.Context, categoryID uint) ([]*skillz.Group, error) {

	groups, err := s.cache.GroupsByCategoryID(ctx, categoryID)
	if err != nil {
		return nil, err
	}

	if len(groups) > 0 {
		return groups, nil
	}

	groups, err = s.universe.Groups(ctx, skillz.NewEqualOperator(mysql.GroupCategoryID, categoryID))
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, errors.Wrap(err, "failed to fetch groups from data store")
	}

	defer func(categoryID uint, groups []*skillz.Group) {
		err = s.cache.SetGroupsByCategoryID(ctx, categoryID, groups)
		if err != nil {
			s.logger.WithError(err).Error("failed to cache groups by category id")
		}
	}(categoryID, groups)

	return groups, nil

}

func (s *Service) Race(ctx context.Context, id uint) (*skillz.Race, error) {

	race, err := s.cache.Race(ctx, id)
	if err != nil {
		return nil, err
	}

	if race != nil {
		return race, nil
	}

	race, err = s.universe.Race(ctx, id)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, errors.Wrap(err, "failed to fetch race from data store")
	}

	return race, s.cache.SetRace(ctx, race)

}

func (s *Service) Region(ctx context.Context, regionID uint) (*skillz.Region, error) {

	region, err := s.cache.Region(ctx, regionID)
	if err != nil {
		return nil, err
	}

	if region != nil {
		return region, nil
	}

	region, err = s.universe.Region(ctx, regionID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, errors.Wrap(err, "failed to fetch region from data store")
	}

	return region, s.cache.SetRegion(ctx, region)

}

func (s *Service) SolarSystem(ctx context.Context, solarSystemID uint) (*skillz.SolarSystem, error) {

	solarSystem, err := s.cache.SolarSystem(ctx, solarSystemID)
	if err != nil {
		return nil, err
	}

	if solarSystem != nil {
		return solarSystem, nil
	}

	solarSystem, err = s.universe.SolarSystem(ctx, solarSystemID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, errors.Wrap(err, "failed to fetch solarSystem from data store")
	}

	return solarSystem, s.cache.SetSolarSystem(ctx, solarSystem)

}

func (s *Service) Station(ctx context.Context, stationID uint) (*skillz.Station, error) {

	station, err := s.cache.Station(ctx, stationID)
	if err != nil {
		return nil, err
	}

	if station != nil {
		return station, nil
	}

	station, err = s.universe.Station(ctx, stationID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, errors.Wrap(err, "failed to fetch station from data store")
	}

	return station, s.cache.SetStation(ctx, station)

}

// TODO: Need Auth Token from Currently Authenticated Character
func (s *Service) Structure(ctx context.Context, structureID uint64) (*skillz.Structure, error) {

	structure, err := s.cache.Structure(ctx, structureID)
	if err != nil {
		return nil, err
	}

	if structure != nil {
		return structure, nil
	}

	etagID, etag, err := s.esi.Etag(ctx, esi.GetStructure, &esi.Params{StructureID: null.Uint64From(structureID)})
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch etag for expiry check")
	}

	structure, err = s.universe.Structure(ctx, structureID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, errors.Wrap(err, "failed to fetch structure from data store")
	}

	exists := err == nil
	if !exists {
		user := internal.UserFromContext(ctx)
		if user == nil {
			return nil, nil
		}

		mods := append(make([]esi.ModifierFunc, 0, 2), s.esi.CacheEtag(ctx, etagID, nil), s.esi.AddAuthorizationHeader(ctx, user.AccessToken))
		if etag != nil && etag.Etag != "" {
			mods = append(mods, s.esi.AddIfNoneMatchHeader(ctx, etag.Etag))
		}

		updatedStructure, err := s.esi.GetStructure(ctx, structureID, mods...)
		if err != nil {
			return nil, err
		}

		if updatedStructure != nil {
			switch exists {
			case true:
				err = s.universe.UpdateStructure(ctx, updatedStructure)
				if err != nil {
					return nil, errors.Wrap(err, "failed to save character to data store")
				}
			case false:
				err = s.universe.CreateStructure(ctx, updatedStructure)
				if err != nil {
					return nil, errors.Wrap(err, "failed to save character to data store")
				}

			}

			structure = updatedStructure
		}

	}

	return structure, s.cache.SetStructure(ctx, structure)

}

func (s *Service) SkillTypesHydrated(ctx context.Context) ([]*skillz.Type, error) {

	types, err := s.cache.SkillTypes(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch skill types from cache")
	}

	if len(types) > 0 {
		return types, err
	}

	groups, err := s.GroupsByCategory(ctx, 16)
	if err != nil {
		return nil, err
	}

	mapGroups := make(map[uint]*skillz.Group)
	groupIDs := make([]interface{}, 0, len(groups))
	for _, group := range groups {
		mapGroups[group.ID] = group
		groupIDs = append(groupIDs, group.ID)
	}

	skillTypes, err := s.universe.Types(ctx, skillz.NewInOperator("group_id", groupIDs))
	if err != nil {
		return nil, err
	}

	for _, t := range skillTypes {
		skillGroup, ok := mapGroups[t.GroupID]
		if !ok {
			continue
		}

		t.Group = skillGroup

	}

	defer func(types []*skillz.Type) {
		err := s.cache.SetSkillTypes(ctx, skillTypes, 0)
		if err != nil {
			s.logger.WithError(err).Error("failed to cache skill types")
		}
	}(skillTypes)

	return skillTypes, nil

}

const (
	CategoryShips    = uint(6)
	CategorySkills   = uint(16)
	CategoryImplants = uint(20)
)

func (s *Service) TypeGroupsHydrated(ctx context.Context, categoryID uint) ([]*skillz.Group, error) {

	groups, err := s.cache.GroupsByCategoryID(ctx, categoryID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch skill groups from cache")
	}

	if len(groups) > 0 {
		return groups, err
	}

	groups, err = s.GroupsByCategory(ctx, categoryID)
	if err != nil {
		return nil, err
	}

	mapGroups := make(map[uint]*skillz.Group)
	groupIDs := make([]interface{}, 0, len(groups))
	for _, group := range groups {
		mapGroups[group.ID] = group
		groupIDs = append(groupIDs, group.ID)
	}

	skillTypes, err := s.universe.Types(ctx, skillz.NewInOperator(mysql.TypesGroupID, groupIDs), skillz.NewEqualOperator(mysql.TypesPublished, 1))
	if err != nil {
		return nil, err
	}

	skillIDs := make([]uint, 0, len(skillTypes))
	for _, skill := range skillTypes {
		skillIDs = append(skillIDs, skill.ID)
	}

	skillDogmaAttributes, err := s.universe.TypeDogmaAttributesBulk(ctx, skillIDs)
	if err != nil {
		return nil, err
	}

	var mapTypeDogmaAttributes = make(map[uint][]*skillz.TypeDogmaAttribute)
	for _, attribute := range skillDogmaAttributes {
		if _, ok := mapTypeDogmaAttributes[attribute.TypeID]; !ok {
			mapTypeDogmaAttributes[attribute.TypeID] = make([]*skillz.TypeDogmaAttribute, 0, 50)
		}

		mapTypeDogmaAttributes[attribute.TypeID] = append(mapTypeDogmaAttributes[attribute.TypeID], attribute)
	}

	for _, t := range skillTypes {
		group := mapGroups[t.GroupID]
		if group == nil {
			continue
		}

		t.Attributes = mapTypeDogmaAttributes[t.ID]

		if group.Types == nil {
			group.Types = make([]*skillz.Type, 0, 15)
		}

		group.Types = append(group.Types, t)

	}

	defer func(groups []*skillz.Group) {
		err := s.cache.SetGroupsByCategoryID(ctx, categoryID, groups)
		if err != nil {
			s.logger.WithError(err).Error("failed to cache skill groups")
		}
	}(groups)

	return groups, nil

}

func (s *Service) ShipTypesHydrated(ctx context.Context) ([]*skillz.Type, error) {

	types, err := s.cache.ShipTypes(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch ship types from cache")
	}

	if len(types) > 0 {
		return types, err
	}

	groups, err := s.GroupsByCategory(ctx, 6)
	if err != nil {
		return nil, err
	}

	mapGroups := make(map[uint]*skillz.Group)
	groupIDs := make([]interface{}, 0, len(groups))
	for _, group := range groups {
		mapGroups[group.ID] = group
		groupIDs = append(groupIDs, group.ID)
	}

	shipTypes, err := s.universe.Types(ctx, skillz.NewInOperator("group_id", groupIDs))
	if err != nil {
		return nil, err
	}

	tIDs := make([]uint, 0, len(shipTypes))
	for _, t := range shipTypes {
		tIDs = append(tIDs, t.ID)
	}

	attributes, err := s.universe.TypeDogmaAttributesBulk(ctx, tIDs)
	if err != nil {
		return nil, err
	}

	mapAttributes := make(map[uint][]*skillz.TypeDogmaAttribute)
	for _, attribute := range attributes {
		if _, ok := mapAttributes[attribute.TypeID]; !ok {
			mapAttributes[attribute.TypeID] = make([]*skillz.TypeDogmaAttribute, 0)
		}

		mapAttributes[attribute.TypeID] = append(mapAttributes[attribute.TypeID], attribute)
	}

	for _, t := range shipTypes {
		group, ok := mapGroups[t.GroupID]
		if !ok {
			continue
		}
		t.Group = group
		attriutes, ok := mapAttributes[t.ID]
		if !ok {
			continue
		}
		t.Attributes = attriutes
	}

	defer func(shipTypes []*skillz.Type) {
		err = s.cache.SetShipTypes(ctx, shipTypes, 0)
		if err != nil {
			s.logger.WithError(err).Error("failed to cache ship types")
		}
	}(shipTypes)

	return shipTypes, nil

}

func (s *Service) Type(ctx context.Context, itemID uint) (*skillz.Type, error) {

	item, err := s.cache.Type(ctx, itemID)
	if err != nil {
		return nil, err
	}

	if item != nil {
		return item, nil
	}

	item, err = s.universe.Type(ctx, itemID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, errors.Wrap(err, "failed to fetch item from data store")
	}

	attributes, err := s.universe.TypeDogmaAttributes(ctx, item.ID)
	if err != nil {
		return nil, err
	}

	item.Attributes = attributes

	defer func() {
		err := s.cache.SetType(ctx, item)
		if err != nil {
			s.logger.WithError(err).Error("failed to cache type")
		}
	}()

	return item, nil

}

func (s *Service) TypeAttributes(ctx context.Context, id uint) ([]*skillz.TypeDogmaAttribute, error) {

	attributes, err := s.cache.TypeAttributes(ctx, id)
	if err != nil {
		return nil, err
	}

	if len(attributes) > 0 {
		return attributes, nil
	}

	attributes, err = s.universe.TypeDogmaAttributes(ctx, id)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, errors.Wrap(err, "failed to fetch type attributes from data store")
	}

	return attributes, s.cache.SetTypeAttributes(ctx, id, attributes)

}

func (s *Service) TypesByGroup(ctx context.Context, groupID uint) ([]*skillz.Type, error) {

	types, err := s.cache.TypesByGroupID(ctx, groupID)
	if err != nil {
		return nil, err
	}

	if len(types) > 0 {
		return types, nil
	}

	types, err = s.universe.Types(ctx, skillz.NewEqualOperator(mysql.TypesGroupID, groupID))
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	for _, t := range types {
		t.Attributes, err = s.TypeAttributes(ctx, t.ID)
		if err != nil {
			return nil, err
		}
	}

	return types, s.cache.SetTypesByGroupID(ctx, groupID, types)

}
