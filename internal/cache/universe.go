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

type UniverseAPI interface {
	Bloodline(ctx context.Context, id uint) (*skillz.Bloodline, error)
	SetBloodline(ctx context.Context, bloodline *skillz.Bloodline) error
	Race(ctx context.Context, id uint) (*skillz.Race, error)
	SetRace(ctx context.Context, race *skillz.Race) error
	Faction(ctx context.Context, id uint) (*skillz.Faction, error)
	SetFaction(ctx context.Context, faction *skillz.Faction) error
	Region(ctx context.Context, id uint) (*skillz.Region, error)
	SetRegion(ctx context.Context, region *skillz.Region) error
	Constellation(ctx context.Context, id uint) (*skillz.Constellation, error)
	SetConstellation(ctx context.Context, constellation *skillz.Constellation) error
	SolarSystem(ctx context.Context, id uint) (*skillz.SolarSystem, error)
	SetSolarSystem(ctx context.Context, solarSystem *skillz.SolarSystem) error
	Station(ctx context.Context, id uint) (*skillz.Station, error)
	SetStation(ctx context.Context, station *skillz.Station) error
	Structure(ctx context.Context, id uint64) (*skillz.Structure, error)
	SetStructure(ctx context.Context, structure *skillz.Structure) error
	Category(ctx context.Context, id uint) (*skillz.Category, error)
	SetCategory(ctx context.Context, category *skillz.Category) error
	Group(ctx context.Context, id uint) (*skillz.Group, error)
	SetGroup(ctx context.Context, group *skillz.Group) error
	GroupsByCategoryID(ctx context.Context, id uint) ([]*skillz.Group, error)
	SetGroupsByCategoryID(ctx context.Context, categoryID uint, groups []*skillz.Group) error
	Type(ctx context.Context, id uint) (*skillz.Type, error)
	SetType(ctx context.Context, item *skillz.Type) error
	SkillTypes(ctx context.Context) ([]*skillz.Type, error)
	SetSkillTypes(ctx context.Context, items []*skillz.Type, expires time.Duration) error
	SkillGroups(ctx context.Context) ([]*skillz.Group, error)
	SetSkillGroups(ctx context.Context, items []*skillz.Group, expires time.Duration) error
	ShipTypes(ctx context.Context) ([]*skillz.Type, error)
	SetShipTypes(ctx context.Context, items []*skillz.Type, expires time.Duration) error
	TypeAttributes(ctx context.Context, id uint) ([]*skillz.TypeDogmaAttribute, error)
	SetTypeAttributes(ctx context.Context, id uint, attributes []*skillz.TypeDogmaAttribute) error
	TypesByGroupID(ctx context.Context, id uint) ([]*skillz.Type, error)
	SetTypesByGroupID(ctx context.Context, groupID uint, types []*skillz.Type) error
}

const (
	keyBloodline       = "bloodline"
	keyRace            = "race"
	keyFaction         = "faction"
	keyCategory        = "category"
	keyGroup           = "group"
	keyGroupByCategory = "group-by-category"
	keyType            = "type"
	keySkillTypes      = "types-skills"
	keySkillGroups     = "groups-skills"
	keyShipsTypes      = "types-ships"
	keyTypeAttributes  = "type-attributes"
	keyTypesByGroup    = "types-by-group"
	keyRegion          = "region"
	keyConstellation   = "constellation"
	keySolarSystem     = "solar-system"
	keyStation         = "station"
	keyStructure       = "structure"
)

func (s *Service) Bloodline(ctx context.Context, bloodlineID uint) (*skillz.Bloodline, error) {

	key := generateKey(keyBloodline, strconv.FormatUint(uint64(bloodlineID), 10))
	result, err := s.redis.Get(ctx, key).Bytes()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, errors.Wrapf(err, errorFFormat, universeAPI, "Bloodline", "failed to fetch results from cache")
	}

	if errors.Is(err, redis.Nil) {
		return nil, nil
	}

	var bloodline = new(skillz.Bloodline)
	err = json.Unmarshal(result, bloodline)
	return bloodline, errors.Wrapf(err, errorFFormat, universeAPI, "Bloodline", "failed to decode json to structure")

}

func (s *Service) SetBloodline(ctx context.Context, bloodline *skillz.Bloodline) error {

	key := generateKey(keyBloodline, strconv.FormatUint(uint64(bloodline.ID), 10))
	data, err := json.Marshal(bloodline)
	if err != nil {
		return errors.Wrapf(err, errorFFormat, universeAPI, "SetBloodline", "failed to encode struct as json")
	}

	err = s.redis.Set(ctx, key, data, time.Hour).Err()
	return errors.Wrapf(err, errorFFormat, universeAPI, "SetBloodline", "failed to write cache")

}

func (s *Service) Race(ctx context.Context, id uint) (*skillz.Race, error) {

	key := generateKey(keyRace, strconv.FormatUint(uint64(id), 10))
	result, err := s.redis.Get(ctx, key).Bytes()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, errors.Wrapf(err, errorFFormat, universeAPI, "Race", "failed to fetch results from cache")
	}

	if errors.Is(err, redis.Nil) {
		return nil, nil
	}

	var race = new(skillz.Race)
	err = json.Unmarshal(result, race)
	return race, errors.Wrapf(err, errorFFormat, universeAPI, "Race", "failed to decode json to structure")

}

func (s *Service) SetRace(ctx context.Context, race *skillz.Race) error {

	key := generateKey(keyRace, strconv.FormatUint(uint64(race.ID), 10))
	data, err := json.Marshal(race)
	if err != nil {
		return errors.Wrapf(err, errorFFormat, universeAPI, "SetRace", "failed to encode structure as json")
	}

	err = s.redis.Set(ctx, key, data, time.Hour).Err()
	return errors.Wrapf(err, errorFFormat, universeAPI, "SetRace", "failed to write cache")

}

func (s *Service) Faction(ctx context.Context, id uint) (*skillz.Faction, error) {

	key := generateKey(keyFaction, strconv.FormatUint(uint64(id), 10))
	result, err := s.redis.Get(ctx, key).Bytes()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, errors.Wrapf(err, errorFFormat, universeAPI, "Faction", "failed to fetch results from cache")
	}

	if errors.Is(err, redis.Nil) {
		return nil, nil
	}

	var faction = new(skillz.Faction)
	err = json.Unmarshal(result, faction)
	return faction, errors.Wrapf(err, errorFFormat, universeAPI, "Faction", "failed to decode json to structure")

}

func (s *Service) SetFaction(ctx context.Context, faction *skillz.Faction) error {

	key := generateKey(keyFaction, strconv.FormatUint(uint64(faction.ID), 10))
	data, err := json.Marshal(faction)
	if err != nil {
		return errors.Wrapf(err, errorFFormat, universeAPI, "SetFaction", "failed to encode structure as json")
	}

	err = s.redis.Set(ctx, key, data, time.Hour).Err()
	return errors.Wrapf(err, errorFFormat, universeAPI, "SetFaction", "failed to write cache")

}

func (s *Service) Region(ctx context.Context, id uint) (*skillz.Region, error) {

	key := generateKey(keyRegion, strconv.FormatUint(uint64(id), 10))
	result, err := s.redis.Get(ctx, key).Bytes()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, errors.Wrapf(err, errorFFormat, universeAPI, "Region", "failed to fetch results from cache")
	}

	if errors.Is(err, redis.Nil) {
		return nil, nil
	}

	var region = new(skillz.Region)
	err = json.Unmarshal(result, region)
	return region, errors.Wrapf(err, errorFFormat, universeAPI, "Region", "failed to decode json to structure")

}

func (s *Service) SetRegion(ctx context.Context, region *skillz.Region) error {

	key := generateKey(keyRegion, strconv.FormatUint(uint64(region.ID), 10))
	data, err := json.Marshal(region)
	if err != nil {
		return errors.Wrapf(err, errorFFormat, universeAPI, "SetRegion", "failed to encode structure as json")
	}

	err = s.redis.Set(ctx, key, data, time.Hour).Err()
	return errors.Wrapf(err, errorFFormat, universeAPI, "SetRegion", "failed to write cache")

}

func (s *Service) Constellation(ctx context.Context, id uint) (*skillz.Constellation, error) {

	key := generateKey(keyConstellation, strconv.FormatUint(uint64(id), 10))
	result, err := s.redis.Get(ctx, key).Bytes()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, errors.Wrapf(err, errorFFormat, universeAPI, "Constellation", "failed to fetch results from cache")
	}

	if errors.Is(err, redis.Nil) {
		return nil, nil
	}

	var constellation = new(skillz.Constellation)
	err = json.Unmarshal(result, constellation)
	return constellation, errors.Wrapf(err, errorFFormat, universeAPI, "Constellation", "failed to decode json to structure")

}

func (s *Service) SetConstellation(ctx context.Context, constellation *skillz.Constellation) error {

	key := generateKey(keyConstellation, strconv.FormatUint(uint64(constellation.ID), 10))
	data, err := json.Marshal(constellation)
	if err != nil {
		return errors.Wrapf(err, errorFFormat, universeAPI, "SetConstellation", "failed to encode structure as json")
	}

	err = s.redis.Set(ctx, key, data, time.Hour).Err()
	return errors.Wrapf(err, errorFFormat, universeAPI, "SetConstellation", "failed to write cache")

}

func (s *Service) SolarSystem(ctx context.Context, id uint) (*skillz.SolarSystem, error) {

	key := generateKey(keySolarSystem, strconv.FormatUint(uint64(id), 10))
	result, err := s.redis.Get(ctx, key).Bytes()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, errors.Wrapf(err, errorFFormat, universeAPI, "SolarSystem", "failed to fetch results from cache")
	}

	if errors.Is(err, redis.Nil) {
		return nil, nil
	}

	var solarSystem = new(skillz.SolarSystem)
	err = json.Unmarshal(result, solarSystem)
	return solarSystem, errors.Wrapf(err, errorFFormat, universeAPI, "SolarSystem", "failed to decode json to structure")

}

func (s *Service) SetSolarSystem(ctx context.Context, solarSystem *skillz.SolarSystem) error {

	key := generateKey(keySolarSystem, strconv.FormatUint(uint64(solarSystem.ID), 10))
	data, err := json.Marshal(solarSystem)
	if err != nil {
		return errors.Wrapf(err, errorFFormat, universeAPI, "SetSolarSystem", "failed to encode structure as json")
	}

	err = s.redis.Set(ctx, key, data, time.Hour).Err()
	return errors.Wrapf(err, errorFFormat, universeAPI, "SetSolarSystem", "failed to write cache")

}

func (s *Service) Station(ctx context.Context, id uint) (*skillz.Station, error) {

	key := generateKey(keyStation, strconv.FormatUint(uint64(id), 10))
	result, err := s.redis.Get(ctx, key).Bytes()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, errors.Wrapf(err, errorFFormat, universeAPI, "Station", "failed to fetch results from cache")
	}

	if errors.Is(err, redis.Nil) {
		return nil, nil
	}

	var station = new(skillz.Station)
	err = json.Unmarshal(result, station)
	return station, errors.Wrapf(err, errorFFormat, universeAPI, "Station", "failed to decode json to structure")

}

func (s *Service) SetStation(ctx context.Context, station *skillz.Station) error {

	key := generateKey(keyStation, strconv.FormatUint(uint64(station.ID), 10))
	data, err := json.Marshal(station)
	if err != nil {
		return errors.Wrapf(err, errorFFormat, universeAPI, "SetStation", "failed to encode structure as json")
	}

	err = s.redis.Set(ctx, key, data, time.Hour).Err()
	return errors.Wrapf(err, errorFFormat, universeAPI, "SetStation", "failed to write cache")

}

func (s *Service) Structure(ctx context.Context, id uint64) (*skillz.Structure, error) {

	key := generateKey(keyStructure, strconv.FormatUint(uint64(id), 10))
	result, err := s.redis.Get(ctx, key).Bytes()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, errors.Wrapf(err, errorFFormat, universeAPI, "Structure", "failed to fetch results from cache")
	}

	if errors.Is(err, redis.Nil) {
		return nil, nil
	}

	var structure = new(skillz.Structure)
	err = json.Unmarshal(result, structure)
	return structure, errors.Wrapf(err, errorFFormat, universeAPI, "Structure", "failed to decode json to structure")

}

func (s *Service) SetStructure(ctx context.Context, structure *skillz.Structure) error {

	key := generateKey(keyStructure, strconv.FormatUint(uint64(structure.ID), 10))
	data, err := json.Marshal(structure)
	if err != nil {
		return errors.Wrapf(err, errorFFormat, universeAPI, "SetStructure", "failed to encode structure as json")
	}

	err = s.redis.Set(ctx, key, data, time.Hour).Err()
	return errors.Wrapf(err, errorFFormat, universeAPI, "SetStructure", "failed to write cache")

}

func (s *Service) Category(ctx context.Context, id uint) (*skillz.Category, error) {

	key := generateKey(keyCategory, strconv.FormatUint(uint64(id), 10))
	result, err := s.redis.Get(ctx, key).Bytes()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, errors.Wrapf(err, errorFFormat, universeAPI, "Category", "failed to fetch results from cache")
	}

	if errors.Is(err, redis.Nil) {
		return nil, nil
	}

	var category = new(skillz.Category)
	err = json.Unmarshal(result, category)
	return category, errors.Wrapf(err, errorFFormat, universeAPI, "Category", "failed to decode json to structure")

}

func (s *Service) SetCategory(ctx context.Context, category *skillz.Category) error {

	key := generateKey(keyCategory, strconv.FormatUint(uint64(category.ID), 10))
	data, err := json.Marshal(category)
	if err != nil {
		return errors.Wrapf(err, errorFFormat, universeAPI, "SetCategory", "failed to encode structure as json")
	}

	err = s.redis.Set(ctx, key, data, time.Hour).Err()
	return errors.Wrapf(err, errorFFormat, universeAPI, "SetCategory", "failed to write cache")

}

func (s *Service) Group(ctx context.Context, id uint) (*skillz.Group, error) {

	key := generateKey(keyGroup, strconv.FormatUint(uint64(id), 10))
	result, err := s.redis.Get(ctx, key).Bytes()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, errors.Wrapf(err, errorFFormat, universeAPI, "Group", "failed to fetch results from cache")
	}

	if errors.Is(err, redis.Nil) {
		return nil, nil
	}

	var group = new(skillz.Group)
	err = json.Unmarshal(result, group)
	return group, errors.Wrapf(err, errorFFormat, universeAPI, "Group", "failed to decode json to structure")

}

func (s *Service) SetGroup(ctx context.Context, group *skillz.Group) error {

	key := generateKey(keyGroup, strconv.FormatUint(uint64(group.ID), 10))
	data, err := json.Marshal(group)
	if err != nil {
		return errors.Wrapf(err, errorFFormat, universeAPI, "SetGroup", "failed to encode structure as json")
	}

	err = s.redis.Set(ctx, key, data, time.Hour).Err()
	return errors.Wrapf(err, errorFFormat, universeAPI, "SetGroup", "failed to write cache")

}

func (s *Service) GroupsByCategoryID(ctx context.Context, id uint) ([]*skillz.Group, error) {

	var groups = make([]*skillz.Group, 0)

	key := generateKey(keyGroupByCategory, strconv.FormatUint(uint64(id), 10))
	result, err := s.redis.Get(ctx, key).Bytes()
	if err != nil && !errors.Is(err, redis.Nil) {
		return groups, errors.Wrapf(err, errorFFormat, universeAPI, "GroupsByCategoryID", "failed to fetch results from cache")
	}

	if errors.Is(err, redis.Nil) {
		return groups, nil
	}

	err = json.Unmarshal(result, &groups)
	return groups, errors.Wrapf(err, errorFFormat, universeAPI, "GroupsByCategoryID", "failed to decode json to structure")

}

func (s *Service) SetGroupsByCategoryID(ctx context.Context, categoryID uint, groups []*skillz.Group) error {

	key := generateKey(keyGroupByCategory, strconv.FormatUint(uint64(categoryID), 10))
	data, err := json.Marshal(groups)
	if err != nil {
		return errors.Wrapf(err, errorFFormat, universeAPI, "SetGroupsByCategoryID", "failed to encode structure as json")
	}

	err = s.redis.Set(ctx, key, data, time.Hour).Err()
	return errors.Wrapf(err, errorFFormat, universeAPI, "SetGroupsByCategoryID", "failed to write cache")

}

func (s *Service) Type(ctx context.Context, id uint) (*skillz.Type, error) {

	key := generateKey(keyType, strconv.FormatUint(uint64(id), 10))
	result, err := s.redis.Get(ctx, key).Bytes()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, errors.Wrapf(err, errorFFormat, universeAPI, "Type", "failed to fetch results from cache")
	}

	if errors.Is(err, redis.Nil) {
		return nil, nil
	}

	var item = new(skillz.Type)
	err = json.Unmarshal(result, item)
	return item, errors.Wrapf(err, errorFFormat, universeAPI, "Type", "failed to decode json to structure")

}

func (s *Service) SetType(ctx context.Context, item *skillz.Type) error {

	key := generateKey(keyType, strconv.FormatUint(uint64(item.ID), 10))
	data, err := json.Marshal(item)
	if err != nil {
		return errors.Wrapf(err, errorFFormat, universeAPI, "SetType", "failed to encode structure as json")
	}

	err = s.redis.Set(ctx, key, data, time.Hour).Err()
	return errors.Wrapf(err, errorFFormat, universeAPI, "SetType", "failed to write cache")

}

func (s *Service) decodeTypes(ctx context.Context, results []string) ([]*skillz.Type, error) {

	var out = make([]*skillz.Type, 0, len(results))

	for _, result := range results {
		var item = new(skillz.Type)
		err := json.Unmarshal([]byte(result), item)
		if err != nil {
			return nil, errors.Wrapf(err, errorFFormat, universeAPI, "decodeTypes", "failed to decode json to structure")
		}

		out = append(out, item)
	}

	return out, nil
}

func (s *Service) ShipTypes(ctx context.Context) ([]*skillz.Type, error) {

	key := generateKey(keyShipsTypes)
	results, err := s.redis.SMembers(ctx, key).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, errors.Wrapf(err, errorFFormat, universeAPI, "ShipTypes", "failed to fetch results from cache")
	}

	if errors.Is(err, redis.Nil) || len(results) == 0 {
		return nil, nil
	}

	items, err := s.decodeTypes(ctx, results)
	if err != nil {
		return nil, err
	}

	return items, nil

}

func (s *Service) SkillTypes(ctx context.Context) ([]*skillz.Type, error) {

	key := generateKey(keySkillTypes)
	results, err := s.redis.SMembers(ctx, key).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, errors.Wrapf(err, errorFFormat, universeAPI, "SkillTypes", "failed to fetch results from cache")
	}

	if errors.Is(err, redis.Nil) || len(results) == 0 {
		return nil, nil
	}

	items, err := s.decodeTypes(ctx, results)
	if err != nil {
		return nil, err
	}

	return items, nil

}

func (s *Service) formatSliceSkillTypes(ctx context.Context, items []*skillz.Type) ([]interface{}, error) {

	members := make([]interface{}, 0, len(items))
	for _, item := range items {
		data, err := json.Marshal(item)
		if err != nil {
			return nil, errors.Wrapf(err, errorFFormat, universeAPI, "formatSliceSkillTypes", "failed to encode struct as json")
		}

		members = append(members, string(data))
	}

	return members, nil
}

func (s *Service) SetSkillTypes(ctx context.Context, items []*skillz.Type, expires time.Duration) error {

	members, err := s.formatSliceSkillTypes(ctx, items)
	if err != nil {
		return err
	}

	key := generateKey(keySkillTypes)
	err = s.redis.SAdd(ctx, key, members...).Err()
	if err != nil {
		return errors.Wrapf(err, errorFFormat, universeAPI, "SetSkillTypes", "failed to write cache")
	}

	if expires > 0 {
		err = s.redis.Expire(ctx, key, expires).Err()
	}

	return errors.Wrapf(err, errorFFormat, universeAPI, "SetSkillTypes", "failed to set expiry on set")

}

func (s *Service) SetShipTypes(ctx context.Context, items []*skillz.Type, expires time.Duration) error {

	members, err := s.formatSliceSkillTypes(ctx, items)
	if err != nil {
		return err
	}

	key := generateKey(keyShipsTypes)
	err = s.redis.SAdd(ctx, key, members...).Err()
	if err != nil {
		return errors.Wrapf(err, errorFFormat, universeAPI, "SetShipTypes", "failed to write cache")
	}

	if expires > 0 {
		err = s.redis.Expire(ctx, key, expires).Err()
	}

	return errors.Wrapf(err, errorFFormat, universeAPI, "SetShipTypes", "failed to set expiry on set")

}

func (s *Service) SkillGroups(ctx context.Context) ([]*skillz.Group, error) {

	key := generateKey(keySkillGroups)
	results, err := s.redis.SMembers(ctx, key).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, errors.Wrapf(err, errorFFormat, universeAPI, "SkillGroups", "failed to fetch results from cache")
	}

	if errors.Is(err, redis.Nil) || len(results) == 0 {
		return nil, nil
	}

	var out = make([]*skillz.Group, 0, len(results))

	for _, result := range results {
		var item = new(skillz.Group)
		err := json.Unmarshal([]byte(result), item)
		if err != nil {
			return nil, errors.Wrapf(err, errorFFormat, universeAPI, "SkillGroups", "failed to decode json to structure")
		}

		out = append(out, item)
	}

	return out, nil

}

func (s *Service) SetSkillGroups(ctx context.Context, groups []*skillz.Group, expires time.Duration) error {

	members := make([]interface{}, 0, len(groups))
	for _, item := range groups {
		data, err := json.Marshal(item)
		if err != nil {
			return errors.Wrapf(err, errorFFormat, universeAPI, "SetSkillGroups", "failed to encode struct as json")
		}

		members = append(members, string(data))
	}

	if len(members) == 0 {
		return nil
	}

	key := generateKey(keySkillGroups)
	err := s.redis.SAdd(ctx, key, members...).Err()
	if err != nil {
		return errors.Wrapf(err, errorFFormat, universeAPI, "SetSkillGroups", "failed to write cache")
	}

	if expires > 0 {
		err = s.redis.Expire(ctx, key, expires).Err()
	}

	return errors.Wrapf(err, errorFFormat, universeAPI, "SetSkillGroups", "failed to set expiry on set")

}

func (s *Service) TypeAttributes(ctx context.Context, id uint) ([]*skillz.TypeDogmaAttribute, error) {

	var attributes = make([]*skillz.TypeDogmaAttribute, 0)

	key := generateKey(keyTypeAttributes, strconv.FormatUint(uint64(id), 10))
	result, err := s.redis.Get(ctx, key).Bytes()
	if err != nil && !errors.Is(err, redis.Nil) {
		return attributes, errors.Wrapf(err, errorFFormat, universeAPI, "TypeAttributes", "failed to fetch results from cache")
	}

	if errors.Is(err, redis.Nil) {
		return attributes, nil
	}

	err = json.Unmarshal(result, &attributes)
	return attributes, errors.Wrapf(err, errorFFormat, universeAPI, "TypeAttributes", "failed to decode json to structure")

}

func (s *Service) SetTypeAttributes(ctx context.Context, id uint, attributes []*skillz.TypeDogmaAttribute) error {

	key := generateKey(keyTypeAttributes, strconv.FormatUint(uint64(id), 10))
	data, err := json.Marshal(attributes)
	if err != nil {
		return errors.Wrapf(err, errorFFormat, universeAPI, "SetTypeAttributes", "failed to encode structure as json")
	}

	err = s.redis.Set(ctx, key, data, time.Hour).Err()
	return errors.Wrapf(err, errorFFormat, universeAPI, "SetTypeAttributes", "failed to write cache")

}

func (s *Service) TypesByGroupID(ctx context.Context, id uint) ([]*skillz.Type, error) {

	var types = make([]*skillz.Type, 0)

	key := generateKey(keyTypesByGroup, strconv.FormatUint(uint64(id), 10))
	result, err := s.redis.Get(ctx, key).Bytes()
	if err != nil && !errors.Is(err, redis.Nil) {
		return types, errors.Wrapf(err, errorFFormat, universeAPI, "TypesByGroupID", "failed to fetch results from cache")
	}

	if errors.Is(err, redis.Nil) {
		return types, nil
	}

	err = json.Unmarshal(result, &types)
	return types, errors.Wrapf(err, errorFFormat, universeAPI, "TypesByGroupID", "failed to decode json to structure")

}

func (s *Service) SetTypesByGroupID(ctx context.Context, groupID uint, types []*skillz.Type) error {

	key := generateKey(keyTypesByGroup, strconv.FormatUint(uint64(groupID), 10))
	data, err := json.Marshal(types)
	if err != nil {
		return errors.Wrapf(err, errorFFormat, universeAPI, "SetTypesByGroupID", "failed to encode structure as json")
	}

	err = s.redis.Set(ctx, key, data, time.Hour).Err()
	return errors.Wrapf(err, errorFFormat, universeAPI, "SetTypesByGroupID", "failed to write cache")

}
