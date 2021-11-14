package mysql

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/eveisesi/skillz"
	"github.com/pkg/errors"
)

type UniverseRepository struct {
	db QueryExecContext

	bloodlines,
	races, factions tableConf

	regions, constellations, solarSystems tableConf

	stations, structures tableConf

	categories, groups, types tableConf
}

const (
	BloodlineID            string = "id"
	BloodlineName          string = "name"
	BloodlineRaceID        string = "race_id"
	BloodlineCorporationID string = "corporation_id"
	BloodlineShipTypeID    string = "ship_type_id"
	BloodlineCharisma      string = "charisma"
	BloodlineIntelligence  string = "intelligence"
	BloodlineMemory        string = "memory"
	BloodlinePerception    string = "perception"
	BloodlineWillpower     string = "willpower"

	CategoryID        string = "id"
	CategoryName      string = "name"
	CategoryPublished string = "published"

	ConstellationID       string = "id"
	ConstellationName     string = "name"
	ConstellationRegionID string = "region_id"

	FactionID                   string = "id"
	FactionName                 string = "name"
	FactionIsUnique             string = "is_unique"
	FactionSizeFactor           string = "size_factor"
	FactionStationCount         string = "station_count"
	FactionStationSystemCount   string = "station_system_count"
	FactionCorporationID        string = "corporation_id"
	FactionMilitiaCorporationID string = "militia_corporation_id"
	FactionSolarSystemID        string = "solar_system_id"

	GroupID         string = "id"
	GroupName       string = "name"
	GroupPublished  string = "published"
	GroupCategoryID string = "category_id"

	RaceID   string = "id"
	RaceName string = "name"

	RegionID   string = "id"
	RegionName string = "name"

	SolarSystemID              string = "id"
	SolarSystemName            string = "name"
	SolarSystemConstellationID string = "constellation_id"
	SolarSystemSecurityStatus  string = "security_status"
	SolarSystemStarID          string = "star_id"
	SolarSystemSecurityClass   string = "security_class"

	StationID                       string = "id"
	StationName                     string = "name"
	StationSystemID                 string = "system_id"
	StationTypeID                   string = "type_id"
	StationRaceID                   string = "race_id"
	StationOwnerCorporationID       string = "owner_corporation_id"
	StationMaxDockableShipVolume    string = "max_dockable_ship_volume"
	StationOfficeRentalCost         string = "office_rental_cost"
	StationReprocessingEfficiency   string = "reprocessing_efficiency"
	StationReprocessingStationsTake string = "reprocessing_stations_take"

	StructureID            string = "id"
	StructureName          string = "name"
	StructureOwnerID       string = "owner_id"
	StructureSolarSystemID string = "solar_system_id"
	StructureTypeID        string = "type_id"

	TypesID             string = "id"
	TypesName           string = "name"
	TypesGroupID        string = "group_id"
	TypesPublished      string = "published"
	TypesCapacity       string = "capacity"
	TypesMarketGroupID  string = "market_group_id"
	TypesMass           string = "mass"
	TypesPackagedVolume string = "packaged_volume"
	TypesPortionSize    string = "portion_size"
	TypesRadius         string = "radius"
	TypesVolume         string = "volume"
)

var _ skillz.UniverseRepository = (*UniverseRepository)(nil)

func NewUniverseRepository(db QueryExecContext) *UniverseRepository {
	return &UniverseRepository{
		db: db,
		bloodlines: tableConf{
			table: "bloodlines",
			columns: []string{
				BloodlineID,
				BloodlineName,
				BloodlineRaceID,
				BloodlineCorporationID,
				BloodlineShipTypeID,
				BloodlineCharisma,
				BloodlineIntelligence,
				BloodlineMemory,
				BloodlinePerception,
				BloodlineWillpower,
				ColumnCreatedAt,
				ColumnUpdatedAt,
			},
		},
		factions: tableConf{
			table: "factions",
			columns: []string{
				FactionID,
				FactionName,
				FactionIsUnique,
				FactionSizeFactor,
				FactionStationCount,
				FactionStationSystemCount,
				FactionCorporationID,
				FactionMilitiaCorporationID,
				FactionSolarSystemID,
				ColumnCreatedAt,
				ColumnUpdatedAt,
			},
		},
		races: tableConf{
			table: "races",
			columns: []string{
				RaceID, RaceName,
				ColumnCreatedAt, ColumnUpdatedAt,
			},
		},

		regions: tableConf{
			table: "map_regions",
			columns: []string{
				RegionID, RegionName,
				ColumnCreatedAt, ColumnUpdatedAt,
			},
		},
		constellations: tableConf{
			table: "map_constellations",
			columns: []string{ConstellationID,
				ConstellationName,
				ConstellationRegionID,
				ColumnCreatedAt,
				ColumnUpdatedAt,
			},
		},
		solarSystems: tableConf{
			table: "map_solar_systems",
			columns: []string{
				SolarSystemID,
				SolarSystemName,
				SolarSystemConstellationID,
				SolarSystemSecurityStatus,
				SolarSystemStarID,
				SolarSystemSecurityClass,
				ColumnCreatedAt, ColumnUpdatedAt,
			},
		},

		stations: tableConf{
			table: "map_stations",
			columns: []string{
				StationID,
				StationName,
				StationSystemID,
				StationTypeID,
				StationRaceID,
				StationOwnerCorporationID,
				StationMaxDockableShipVolume,
				StationOfficeRentalCost,
				StationReprocessingEfficiency,
				StationReprocessingStationsTake,
				ColumnCreatedAt,
				ColumnUpdatedAt,
			},
		},
		structures: tableConf{
			table: "map_structures",
			columns: []string{
				StructureID,
				StructureName,
				StructureOwnerID,
				StructureSolarSystemID,
				StructureTypeID,
				ColumnCreatedAt,
				ColumnUpdatedAt,
			},
		},

		categories: tableConf{
			table: "type_categories",
			columns: []string{
				CategoryID,
				CategoryName,
				CategoryPublished,
				ColumnCreatedAt,
				ColumnUpdatedAt,
			},
		},
		groups: tableConf{
			table: "type_groups",
			columns: []string{
				GroupID,
				GroupName,
				GroupPublished,
				GroupCategoryID,
				ColumnCreatedAt,
				ColumnUpdatedAt,
			},
		},
		types: tableConf{
			table: "types",
			columns: []string{
				TypesID,
				TypesName,
				TypesGroupID,
				TypesPublished,
				TypesCapacity,
				TypesMarketGroupID,
				TypesMass,
				TypesPackagedVolume,
				TypesPortionSize,
				TypesRadius,
				TypesVolume,
				ColumnCreatedAt,
				ColumnUpdatedAt,
			},
		},
	}
}

func (r *UniverseRepository) Bloodline(ctx context.Context, bloodlineID uint) (*skillz.Bloodline, error) {

	query, args, err := sq.Select(r.bloodlines.columns...).
		From(r.bloodlines.table).
		Where(sq.Eq{BloodlineID: bloodlineID}).
		Limit(1).
		ToSql()
	if err != nil {
		return nil, errors.Wrapf(err, errorFFormat, universeRepository, "Bloodline", "failed to generate sql")
	}

	var bloodline = new(skillz.Bloodline)
	err = r.db.GetContext(ctx, bloodline, query, args...)
	return bloodline, errors.Wrapf(err, prefixFormat, universeRepository, "Bloodline")

}

func (r *UniverseRepository) Bloodlines(ctx context.Context) ([]*skillz.Bloodline, error) {

	query, args, err := sq.Select(r.bloodlines.columns...).From(r.bloodlines.table).ToSql()
	if err != nil {
		return nil, errors.Wrapf(err, errorFFormat, universeRepository, "Bloodlines", "failed to generate sql")
	}

	var bloodlines = make([]*skillz.Bloodline, 0)
	err = r.db.SelectContext(ctx, &bloodlines, query, args...)
	return bloodlines, errors.Wrapf(err, prefixFormat, universeRepository, "Bloodlines")

}

func (r *UniverseRepository) CreateBloodline(ctx context.Context, bloodline *skillz.Bloodline) error {

	now := time.Now()
	bloodline.CreatedAt = now
	bloodline.UpdatedAt = now

	query, args, err := sq.Insert(r.bloodlines.table).
		SetMap(map[string]interface{}{
			BloodlineID:            bloodline.ID,
			BloodlineName:          bloodline.Name,
			BloodlineRaceID:        bloodline.RaceID,
			BloodlineCorporationID: bloodline.CorporationID,
			BloodlineShipTypeID:    bloodline.ShipTypeID,
			BloodlineCharisma:      bloodline.Charisma,
			BloodlineIntelligence:  bloodline.Intelligence,
			BloodlineMemory:        bloodline.Memory,
			BloodlinePerception:    bloodline.Perception,
			BloodlineWillpower:     bloodline.Willpower,
			ColumnCreatedAt:        bloodline.CreatedAt,
			ColumnUpdatedAt:        bloodline.UpdatedAt,
		}).
		ToSql()
	if err != nil {
		return errors.Wrapf(err, errorFFormat, universeRepository, "CreateBloodline", "failed to generate sql")
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	return errors.Wrapf(err, prefixFormat, universeRepository, "CreateBloodline")

}

func (r *UniverseRepository) UpdateBloodline(ctx context.Context, bloodline *skillz.Bloodline) error {

	bloodline.UpdatedAt = time.Now()

	query, args, err := sq.Update(r.bloodlines.table).
		SetMap(map[string]interface{}{
			BloodlineName:          bloodline.Name,
			BloodlineRaceID:        bloodline.RaceID,
			BloodlineCorporationID: bloodline.CorporationID,
			BloodlineShipTypeID:    bloodline.ShipTypeID,
			BloodlineCharisma:      bloodline.Charisma,
			BloodlineIntelligence:  bloodline.Intelligence,
			BloodlineMemory:        bloodline.Memory,
			BloodlinePerception:    bloodline.Perception,
			BloodlineWillpower:     bloodline.Willpower,
			ColumnUpdatedAt:        bloodline.UpdatedAt,
		}).
		Where(sq.Eq{BloodlineID: bloodline.ID}).ToSql()
	if err != nil {
		return errors.Wrapf(err, errorFFormat, universeRepository, "UpdateBloodline", "failed to generate sql")
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	return errors.Wrapf(err, prefixFormat, universeRepository, "CreateBloodline")

}

func (r *UniverseRepository) Category(ctx context.Context, categoryID uint) (*skillz.Category, error) {

	query, args, err := sq.Select(r.categories.columns...).
		From(r.categories.table).
		Where(sq.Eq{CategoryID: categoryID}).
		Limit(1).
		ToSql()
	if err != nil {
		return nil, errors.Wrapf(err, errorFFormat, universeRepository, "Category", "failed to generate sql")
	}

	var category = new(skillz.Category)
	err = r.db.GetContext(ctx, category, query, args...)
	return category, errors.Wrapf(err, prefixFormat, universeRepository, "Category")

}

func (r *UniverseRepository) Categories(ctx context.Context) ([]*skillz.Category, error) {

	query, args, err := sq.Select(r.categories.columns...).
		From(r.categories.table).
		ToSql()
	if err != nil {
		return nil, errors.Wrapf(err, errorFFormat, universeRepository, "Categories", "failed to generate sql")
	}

	var categories = make([]*skillz.Category, 0)
	err = r.db.SelectContext(ctx, &categories, query, args...)

	return categories, errors.Wrapf(err, prefixFormat, universeRepository, "Categories")

}

func (r *UniverseRepository) CreateCategory(ctx context.Context, category *skillz.Category) error {

	now := time.Now()
	category.CreatedAt = now
	category.UpdatedAt = now

	query, args, err := sq.Insert(r.categories.table).
		SetMap(map[string]interface{}{
			CategoryID:        category.ID,
			CategoryName:      category.Name,
			CategoryPublished: category.Published,
			ColumnCreatedAt:   category.CreatedAt,
			ColumnUpdatedAt:   category.UpdatedAt,
		}).
		ToSql()
	if err != nil {
		return errors.Wrapf(err, errorFFormat, universeRepository, "CreateCategory", "failed to generate sql")
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	return errors.Wrapf(err, prefixFormat, universeRepository, "CreateCategory")

}

func (r *UniverseRepository) UpdateCategory(ctx context.Context, category *skillz.Category) error {

	category.UpdatedAt = time.Now()

	query, args, err := sq.Update(r.categories.table).
		SetMap(map[string]interface{}{
			CategoryName:      category.Name,
			CategoryPublished: category.Published,
			ColumnUpdatedAt:   category.UpdatedAt,
		}).
		Where(sq.Eq{CategoryID: category.ID}).
		ToSql()
	if err != nil {
		return errors.Wrapf(err, errorFFormat, universeRepository, "UpdateCategory", "failed to generate sql")
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	return errors.Wrapf(err, prefixFormat, universeRepository, "UpdateCategory")

}

func (r *UniverseRepository) Constellation(ctx context.Context, constellationID uint) (*skillz.Constellation, error) {

	query, args, err := sq.Select(r.constellations.columns...).
		From(r.constellations.table).
		Where(sq.Eq{ConstellationID: constellationID}).
		Limit(1).
		ToSql()
	if err != nil {
		return nil, errors.Wrapf(err, errorFFormat, universeRepository, "Constellation", "failed to generate sql")
	}

	var constellation = new(skillz.Constellation)
	err = r.db.GetContext(ctx, constellation, query, args...)
	return constellation, errors.Wrapf(err, prefixFormat, universeRepository, "Constellation")

}

func (r *UniverseRepository) Constellations(ctx context.Context) ([]*skillz.Constellation, error) {

	query, args, err := sq.Select(r.constellations.columns...).
		From(r.constellations.table).
		ToSql()
	if err != nil {
		return nil, errors.Wrapf(err, errorFFormat, universeRepository, "Constellations", "failed to generate sql")
	}

	var constellations = make([]*skillz.Constellation, 0)
	err = r.db.GetContext(ctx, &constellations, query, args...)
	return constellations, errors.Wrapf(err, prefixFormat, universeRepository, "Constellations")

}

func (r *UniverseRepository) CreateConstellation(ctx context.Context, constellation *skillz.Constellation) error {

	now := time.Now()
	constellation.CreatedAt = now
	constellation.UpdatedAt = now

	query, args, err := sq.Insert(r.constellations.table).
		SetMap(map[string]interface{}{
			ConstellationID:       constellation.ID,
			ConstellationName:     constellation.Name,
			ConstellationRegionID: constellation.RegionID,
			ColumnCreatedAt:       constellation.CreatedAt,
			ColumnUpdatedAt:       constellation.UpdatedAt,
		}).
		ToSql()
	if err != nil {
		return errors.Wrapf(err, errorFFormat, universeRepository, "CreateConstellation", "failed to generate sql")
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	return errors.Wrapf(err, prefixFormat, universeRepository, "CreateConstellation")

}

func (r *UniverseRepository) UpdateConstellation(ctx context.Context, constellation *skillz.Constellation) error {

	constellation.UpdatedAt = time.Now()

	query, args, err := sq.Update(r.constellations.table).
		SetMap(map[string]interface{}{
			ConstellationID:       constellation.ID,
			ConstellationName:     constellation.Name,
			ConstellationRegionID: constellation.RegionID,
			ColumnUpdatedAt:       constellation.UpdatedAt,
		}).
		Where(sq.Eq{ConstellationID: constellation.ID}).
		ToSql()
	if err != nil {
		return errors.Wrapf(err, errorFFormat, universeRepository, "UpdateConstellation", "failed to generate sql")
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	return errors.Wrapf(err, prefixFormat, universeRepository, "UpdateConstellation")

}

func (r *UniverseRepository) Faction(ctx context.Context, factionID uint) (*skillz.Faction, error) {

	query, args, err := sq.Select(r.factions.columns...).
		From(r.factions.table).
		Where(sq.Eq{FactionID: factionID}).
		Limit(1).
		ToSql()
	if err != nil {
		return nil, errors.Wrapf(err, errorFFormat, universeRepository, "Faction", "failed to generate sql")
	}

	var faction = new(skillz.Faction)
	err = r.db.GetContext(ctx, faction, query, args...)
	return faction, errors.Wrapf(err, prefixFormat, universeRepository, "Faction")

}

func (r *UniverseRepository) Factions(ctx context.Context) ([]*skillz.Faction, error) {

	query, args, err := sq.Select(r.factions.columns...).
		From(r.factions.table).
		ToSql()
	if err != nil {
		return nil, errors.Wrapf(err, errorFFormat, universeRepository, "Factions", "failed to generate sql")
	}

	var factions = make([]*skillz.Faction, 0)
	err = r.db.SelectContext(ctx, &factions, query, args...)
	return factions, errors.Wrapf(err, prefixFormat, universeRepository, "Factions")

}

func (r *UniverseRepository) CreateFaction(ctx context.Context, faction *skillz.Faction) error {

	now := time.Now()
	faction.CreatedAt = now
	faction.UpdatedAt = now

	query, args, err := sq.Insert(r.factions.table).
		SetMap(map[string]interface{}{
			FactionID:                   faction.ID,
			FactionName:                 faction.Name,
			FactionIsUnique:             faction.IsUnique,
			FactionSizeFactor:           faction.SizeFactor,
			FactionStationCount:         faction.StationCount,
			FactionStationSystemCount:   faction.StationSystemCount,
			FactionCorporationID:        faction.CorporationID,
			FactionMilitiaCorporationID: faction.MilitiaCorporationID,
			FactionSolarSystemID:        faction.SolarSystemID,
			ColumnCreatedAt:             faction.CreatedAt,
			ColumnUpdatedAt:             faction.UpdatedAt,
		}).
		ToSql()
	if err != nil {
		return errors.Wrapf(err, errorFFormat, universeRepository, "CreateFaction", "failed to generate sql")
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	return errors.Wrapf(err, prefixFormat, universeRepository, "CreateFaction")

}

func (r *UniverseRepository) UpdateFaction(ctx context.Context, faction *skillz.Faction) error {

	faction.UpdatedAt = time.Now()

	query, args, err := sq.Update(r.factions.table).
		SetMap(map[string]interface{}{
			FactionName:                 faction.Name,
			FactionIsUnique:             faction.IsUnique,
			FactionSizeFactor:           faction.SizeFactor,
			FactionStationCount:         faction.StationCount,
			FactionStationSystemCount:   faction.StationSystemCount,
			FactionCorporationID:        faction.CorporationID,
			FactionMilitiaCorporationID: faction.MilitiaCorporationID,
			FactionSolarSystemID:        faction.SolarSystemID,
			ColumnUpdatedAt:             faction.UpdatedAt,
		}).
		Where(sq.Eq{FactionID: faction.ID}).ToSql()
	if err != nil {
		return errors.Wrapf(err, errorFFormat, universeRepository, "UpdateFaction", "failed to generate sql")
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	return errors.Wrapf(err, prefixFormat, universeRepository, "UpdateFaction")

}

func (r *UniverseRepository) Group(ctx context.Context, groupID uint) (*skillz.Group, error) {

	query, args, err := sq.Select(r.groups.columns...).
		From(r.groups.table).
		Where(sq.Eq{GroupID: groupID}).
		Limit(1).
		ToSql()
	if err != nil {
		return nil, errors.Wrapf(err, errorFFormat, universeRepository, "Group", "failed to generate sql")
	}

	var group = new(skillz.Group)
	err = r.db.GetContext(ctx, group, query, args...)

	return group, errors.Wrapf(err, prefixFormat, universeRepository, "Group")

}

func (r *UniverseRepository) Groups(ctx context.Context, operators ...*skillz.Operator) ([]*skillz.Group, error) {

	query, args, err := sq.Select(r.groups.columns...).
		From(r.groups.table).
		ToSql()
	if err != nil {
		return nil, errors.Wrapf(err, errorFFormat, universeRepository, "Groups", "failed to generate sql")
	}

	var groups = make([]*skillz.Group, 0)
	err = r.db.SelectContext(ctx, &groups, query, args...)

	return groups, errors.Wrapf(err, prefixFormat, universeRepository, "Groups")

}

func (r *UniverseRepository) CreateGroup(ctx context.Context, group *skillz.Group) error {

	now := time.Now()
	group.CreatedAt = now
	group.UpdatedAt = now

	query, args, err := sq.Insert(r.groups.table).
		SetMap(map[string]interface{}{
			GroupID:         group.ID,
			GroupName:       group.Name,
			GroupPublished:  group.Published,
			GroupCategoryID: group.CategoryID,
			ColumnCreatedAt: group.CreatedAt,
			ColumnUpdatedAt: group.UpdatedAt,
		}).
		ToSql()
	if err != nil {
		return errors.Wrapf(err, errorFFormat, universeRepository, "CreateGroup", "failed to generate sql")
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	return errors.Wrapf(err, prefixFormat, universeRepository, "CreateGroup")

}

func (r *UniverseRepository) UpdateGroup(ctx context.Context, group *skillz.Group) error {

	group.UpdatedAt = time.Now()

	query, args, err := sq.Update(r.groups.table).
		SetMap(map[string]interface{}{
			GroupName:       group.Name,
			GroupPublished:  group.Published,
			GroupCategoryID: group.CategoryID,
			ColumnUpdatedAt: group.UpdatedAt,
		}).
		Where(sq.Eq{GroupID: group.ID}).ToSql()
	if err != nil {
		return errors.Wrapf(err, errorFFormat, universeRepository, "UpdateGroup", "failed to generate sql")
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	return errors.Wrapf(err, prefixFormat, universeRepository, "UpdateGroup")

}

func (r *UniverseRepository) Race(ctx context.Context, raceID uint) (*skillz.Race, error) {

	query, args, err := sq.Select(r.races.columns...).
		From(r.races.table).
		Where(sq.Eq{RaceID: raceID}).
		Limit(1).
		ToSql()
	if err != nil {
		return nil, errors.Wrapf(err, errorFFormat, universeRepository, "Race", "failed to generate sql")
	}

	var race = new(skillz.Race)
	err = r.db.GetContext(ctx, race, query, args...)
	return race, errors.Wrapf(err, prefixFormat, universeRepository, "Race")

}

func (r *UniverseRepository) Races(ctx context.Context, operators ...*skillz.Operator) ([]*skillz.Race, error) {

	query, args, err := sq.Select(r.races.columns...).
		From(r.races.table).
		ToSql()
	if err != nil {
		return nil, errors.Wrapf(err, errorFFormat, universeRepository, "Races", "failed to generate sql")
	}

	var races = make([]*skillz.Race, 0)
	err = r.db.SelectContext(ctx, &races, query, args...)
	return races, errors.Wrapf(err, prefixFormat, universeRepository, "Races")

}

func (r *UniverseRepository) CreateRace(ctx context.Context, race *skillz.Race) error {

	now := time.Now()
	race.CreatedAt = now
	race.UpdatedAt = now

	query, args, err := sq.Insert(r.races.table).
		SetMap(map[string]interface{}{
			RaceID:          race.ID,
			RaceName:        race.Name,
			ColumnCreatedAt: race.CreatedAt,
			ColumnUpdatedAt: race.UpdatedAt,
		}).
		ToSql()
	if err != nil {
		return errors.Wrapf(err, errorFFormat, universeRepository, "CreateRace", "failed to generate sql")
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	return errors.Wrapf(err, prefixFormat, universeRepository, "CreateRace")

}

func (r *UniverseRepository) UpdateRace(ctx context.Context, race *skillz.Race) error {

	race.UpdatedAt = time.Now()

	query, args, err := sq.Update(r.races.table).
		SetMap(map[string]interface{}{
			RaceName:        race.Name,
			ColumnUpdatedAt: race.UpdatedAt,
		}).
		Where(sq.Eq{RaceID: race.ID}).ToSql()
	if err != nil {
		return errors.Wrapf(err, errorFFormat, universeRepository, "UpdateRace", "failed to generate sql")
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	return errors.Wrapf(err, prefixFormat, universeRepository, "UpdateRace")

}

func (r *UniverseRepository) Region(ctx context.Context, regionID uint) (*skillz.Region, error) {

	query, args, err := sq.Select(r.regions.columns...).
		From(r.regions.table).
		Where(sq.Eq{RegionID: regionID}).
		Limit(1).
		ToSql()
	if err != nil {
		return nil, errors.Wrapf(err, errorFFormat, universeRepository, "Region", "failed to generate sql")
	}

	var region = new(skillz.Region)
	err = r.db.GetContext(ctx, region, query, args...)
	return region, errors.Wrapf(err, prefixFormat, universeRepository, "Region")

}

func (r *UniverseRepository) Regions(ctx context.Context) ([]*skillz.Region, error) {

	query, args, err := sq.Select(r.regions.columns...).
		From(r.regions.table).
		ToSql()
	if err != nil {
		return nil, errors.Wrapf(err, errorFFormat, universeRepository, "Regions", "failed to generate sql")
	}

	var regions = make([]*skillz.Region, 0)
	err = r.db.SelectContext(ctx, &regions, query, args...)
	return regions, errors.Wrapf(err, prefixFormat, universeRepository, "Regions")

}

func (r *UniverseRepository) CreateRegion(ctx context.Context, region *skillz.Region) error {

	now := time.Now()
	region.CreatedAt = now
	region.UpdatedAt = now

	query, args, err := sq.Insert(r.regions.table).
		SetMap(map[string]interface{}{
			RegionID:        region.ID,
			RegionName:      region.Name,
			ColumnCreatedAt: region.CreatedAt,
			ColumnUpdatedAt: region.UpdatedAt,
		}).
		ToSql()
	if err != nil {
		return errors.Wrapf(err, errorFFormat, universeRepository, "CreateRegion", "failed to generate sql")
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	return errors.Wrapf(err, prefixFormat, universeRepository, "CreateRegion")

}

func (r *UniverseRepository) UpdateRegion(ctx context.Context, region *skillz.Region) error {

	region.UpdatedAt = time.Now()

	query, args, err := sq.Update(r.regions.table).
		SetMap(map[string]interface{}{
			RegionName:      region.Name,
			ColumnUpdatedAt: region.UpdatedAt,
		}).
		Where(sq.Eq{RegionID: region.ID}).ToSql()
	if err != nil {
		return errors.Wrapf(err, errorFFormat, universeRepository, "UpdateRegion", "failed to generate sql")
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	return errors.Wrapf(err, prefixFormat, universeRepository, "UpdateRegion")

}

func (r *UniverseRepository) SolarSystem(ctx context.Context, solarSystemID uint) (*skillz.SolarSystem, error) {

	query, args, err := sq.Select(r.solarSystems.columns...).
		From(r.solarSystems.table).
		Where(sq.Eq{SolarSystemID: solarSystemID}).
		Limit(1).
		ToSql()
	if err != nil {
		return nil, errors.Wrapf(err, errorFFormat, universeRepository, "SolarSystem", "failed to generate sql")
	}

	var solarSystem = new(skillz.SolarSystem)
	err = r.db.GetContext(ctx, solarSystem, query, args...)
	return solarSystem, errors.Wrapf(err, prefixFormat, universeRepository, "SolarSystem")

}

func (r *UniverseRepository) SolarSystems(ctx context.Context, operators ...*skillz.Operator) ([]*skillz.SolarSystem, error) {

	query, args, err := sq.Select(r.solarSystems.columns...).
		From(r.solarSystems.table).
		ToSql()
	if err != nil {
		return nil, errors.Wrapf(err, errorFFormat, universeRepository, "SolarSystems", "failed to generate sql")
	}

	var solarSystems = make([]*skillz.SolarSystem, 0)
	err = r.db.SelectContext(ctx, &solarSystems, query, args...)
	return solarSystems, errors.Wrapf(err, prefixFormat, universeRepository, "SolarSystems")

}

func (r *UniverseRepository) CreateSolarSystem(ctx context.Context, solarSystem *skillz.SolarSystem) error {

	now := time.Now()
	solarSystem.CreatedAt = now
	solarSystem.UpdatedAt = now

	query, args, err := sq.Insert(r.solarSystems.table).SetMap(map[string]interface{}{
		SolarSystemID:              solarSystem.ID,
		SolarSystemName:            solarSystem.Name,
		SolarSystemConstellationID: solarSystem.ConstellationID,
		SolarSystemSecurityStatus:  solarSystem.SecurityStatus,
		SolarSystemStarID:          solarSystem.StarID,
		SolarSystemSecurityClass:   solarSystem.SecurityClass,
		ColumnCreatedAt:            solarSystem.CreatedAt,
		ColumnUpdatedAt:            solarSystem.UpdatedAt,
	}).ToSql()
	if err != nil {
		return errors.Wrapf(err, errorFFormat, universeRepository, "CreateSolarSystem", "failed to generate sql")
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	return errors.Wrapf(err, prefixFormat, universeRepository, "CreateSolarSystem")

}

func (r *UniverseRepository) UpdateSolarSystem(ctx context.Context, solarSystem *skillz.SolarSystem) error {

	now := time.Now()
	solarSystem.UpdatedAt = now

	query, args, err := sq.Update(r.solarSystems.table).SetMap(map[string]interface{}{
		SolarSystemID:              solarSystem.ID,
		SolarSystemName:            solarSystem.Name,
		SolarSystemConstellationID: solarSystem.ConstellationID,
		SolarSystemSecurityStatus:  solarSystem.SecurityStatus,
		SolarSystemStarID:          solarSystem.StarID,
		SolarSystemSecurityClass:   solarSystem.SecurityClass,
		ColumnUpdatedAt:            solarSystem.UpdatedAt,
	}).Where(sq.Eq{SolarSystemID: solarSystem.ID}).ToSql()
	if err != nil {
		return errors.Wrapf(err, errorFFormat, universeRepository, "UpdateSolarSystem", "failed to generate sql")
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	return errors.Wrapf(err, prefixFormat, universeRepository, "UpdateSolarSystem")

}

func (r *UniverseRepository) Station(ctx context.Context, stationID uint) (*skillz.Station, error) {

	query, args, err := sq.Select(r.stations.columns...).
		From(r.stations.table).
		Where(sq.Eq{StationID: stationID}).
		Limit(1).
		ToSql()
	if err != nil {
		return nil, errors.Wrapf(err, errorFFormat, universeRepository, "Station", "failed to generate sql")
	}

	var station = new(skillz.Station)
	err = r.db.GetContext(ctx, station, query, args...)
	return station, errors.Wrapf(err, prefixFormat, universeRepository, "Station")

}

func (r *UniverseRepository) Stations(ctx context.Context, operators ...*skillz.Operator) ([]*skillz.Station, error) {

	query, args, err := sq.Select(r.stations.columns...).
		From(r.stations.table).
		ToSql()
	if err != nil {
		return nil, errors.Wrapf(err, errorFFormat, universeRepository, "Stations", "failed to generate sql")
	}

	var stations = make([]*skillz.Station, 0)
	err = r.db.SelectContext(ctx, &stations, query, args...)
	return stations, errors.Wrapf(err, prefixFormat, universeRepository, "Stations")

}

func (r *UniverseRepository) CreateStation(ctx context.Context, station *skillz.Station) error {

	now := time.Now()
	station.CreatedAt = now
	station.UpdatedAt = now

	query, args, err := sq.Insert(r.stations.table).
		SetMap(map[string]interface{}{
			StationID:                       station.ID,
			StationName:                     station.Name,
			StationSystemID:                 station.SystemID,
			StationTypeID:                   station.TypeID,
			StationRaceID:                   station.RaceID,
			StationOwnerCorporationID:       station.OwnerCorporationID,
			StationMaxDockableShipVolume:    station.MaxDockableShipVolume,
			StationOfficeRentalCost:         station.OfficeRentalCost,
			StationReprocessingEfficiency:   station.ReprocessingEfficiency,
			StationReprocessingStationsTake: station.ReprocessingStationsTake,
			ColumnCreatedAt:                 station.CreatedAt,
			ColumnUpdatedAt:                 station.UpdatedAt,
		}).ToSql()
	if err != nil {
		return errors.Wrapf(err, errorFFormat, universeRepository, "CreateStation", "failed to generate sql")
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	return errors.Wrapf(err, prefixFormat, universeRepository, "CreateStation")

}

func (r *UniverseRepository) UpdateStation(ctx context.Context, station *skillz.Station) error {

	station.UpdatedAt = time.Now()

	query, args, err := sq.Update(r.stations.table).
		SetMap(map[string]interface{}{
			StationName:                     station.Name,
			StationSystemID:                 station.SystemID,
			StationTypeID:                   station.TypeID,
			StationRaceID:                   station.RaceID,
			StationOwnerCorporationID:       station.OwnerCorporationID,
			StationMaxDockableShipVolume:    station.MaxDockableShipVolume,
			StationOfficeRentalCost:         station.OfficeRentalCost,
			StationReprocessingEfficiency:   station.ReprocessingEfficiency,
			StationReprocessingStationsTake: station.ReprocessingStationsTake,
			ColumnUpdatedAt:                 station.UpdatedAt,
		}).
		Where(sq.Eq{StationID: station.ID}).ToSql()
	if err != nil {
		return errors.Wrapf(err, errorFFormat, universeRepository, "UpdateStation", "failed to generate sql")
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	return errors.Wrapf(err, prefixFormat, universeRepository, "UpdateStation")

}

func (r *UniverseRepository) Structure(ctx context.Context, structureID uint64) (*skillz.Structure, error) {

	query, args, err := sq.Select(r.structures.columns...).
		From(r.structures.table).
		Where(sq.Eq{StructureID: structureID}).
		Limit(1).
		ToSql()
	if err != nil {
		return nil, errors.Wrapf(err, errorFFormat, universeRepository, "Structure", "failed to generate sql")
	}

	var structure = new(skillz.Structure)
	err = r.db.GetContext(ctx, structure, query, args...)
	return structure, errors.Wrapf(err, prefixFormat, universeRepository, "Structure")

}

func (r *UniverseRepository) Structures(ctx context.Context) ([]*skillz.Structure, error) {

	query, args, err := sq.Select(r.structures.columns...).
		From(r.structures.table).
		ToSql()
	if err != nil {
		return nil, errors.Wrapf(err, errorFFormat, universeRepository, "Structures", "failed to generate sql")
	}

	var structures = make([]*skillz.Structure, 0)
	err = r.db.SelectContext(ctx, &structures, query, args...)
	return structures, errors.Wrapf(err, prefixFormat, universeRepository, "Structures")

}

func (r *UniverseRepository) CreateStructure(ctx context.Context, structure *skillz.Structure) error {

	now := time.Now()
	structure.CreatedAt = now
	structure.UpdatedAt = now

	query, args, err := sq.Insert(r.structures.table).
		SetMap(map[string]interface{}{
			StructureID:            structure.ID,
			StructureName:          structure.Name,
			StructureOwnerID:       structure.OwnerID,
			StructureSolarSystemID: structure.SolarSystemID,
			StructureTypeID:        structure.TypeID,
			ColumnCreatedAt:        structure.CreatedAt,
			ColumnUpdatedAt:        structure.UpdatedAt,
		}).ToSql()
	if err != nil {
		return errors.Wrapf(err, errorFFormat, universeRepository, "CreateStructure", "failed to generate sql")
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	return errors.Wrapf(err, prefixFormat, universeRepository, "CreateStructure")

}

func (r *UniverseRepository) UpdateStructure(ctx context.Context, structure *skillz.Structure) error {

	structure.UpdatedAt = time.Now()

	query, args, err := sq.Update(r.structures.table).
		SetMap(map[string]interface{}{
			StructureName:          structure.Name,
			StructureOwnerID:       structure.OwnerID,
			StructureSolarSystemID: structure.SolarSystemID,
			StructureTypeID:        structure.TypeID,
			ColumnUpdatedAt:        structure.UpdatedAt,
		}).Where(sq.Eq{StructureID: structure.ID}).ToSql()
	if err != nil {
		return errors.Wrapf(err, errorFFormat, universeRepository, "UpdateStructure", "failed to generate sql")
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	return errors.Wrapf(err, prefixFormat, universeRepository, "UpdateStructure")

}

func (r *UniverseRepository) Type(ctx context.Context, typeID uint) (*skillz.Type, error) {

	query, args, err := sq.Select(r.types.columns...).
		From(r.types.table).
		Where(sq.Eq{TypesID: typeID}).
		Limit(1).
		ToSql()
	if err != nil {
		return nil, errors.Wrapf(err, errorFFormat, universeRepository, "Type", "failed to generate sql")
	}

	var item = new(skillz.Type)
	err = r.db.GetContext(ctx, item, query, args...)
	return item, errors.Wrapf(err, prefixFormat, universeRepository, "Type")

}

func (r *UniverseRepository) Types(ctx context.Context) ([]*skillz.Type, error) {

	query, args, err := sq.Select(r.types.columns...).
		From(r.types.table).
		ToSql()
	if err != nil {
		return nil, errors.Wrapf(err, errorFFormat, universeRepository, "Types", "failed to generate sql")
	}

	var types = make([]*skillz.Type, 0)
	err = r.db.SelectContext(ctx, &types, query, args...)
	return types, errors.Wrapf(err, prefixFormat, universeRepository, "Types")

}

func (r *UniverseRepository) CreateType(ctx context.Context, item *skillz.Type) error {

	now := time.Now()
	item.CreatedAt = now
	item.UpdatedAt = now

	query, args, err := sq.Insert(r.types.table).
		SetMap(map[string]interface{}{
			TypesID:             item.ID,
			TypesName:           item.Name,
			TypesGroupID:        item.GroupID,
			TypesPublished:      item.Published,
			TypesCapacity:       item.Capacity,
			TypesMarketGroupID:  item.MarketGroupID,
			TypesMass:           item.Mass,
			TypesPackagedVolume: item.PackagedVolume,
			TypesPortionSize:    item.PortionSize,
			TypesRadius:         item.Radius,
			TypesVolume:         item.Volume,
			ColumnCreatedAt:     item.CreatedAt,
			ColumnUpdatedAt:     item.UpdatedAt,
		}).ToSql()
	if err != nil {
		return errors.Wrapf(err, errorFFormat, universeRepository, "CreateType", "failed to generate sql")
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	return errors.Wrapf(err, prefixFormat, universeRepository, "CreateType")

}

func (r *UniverseRepository) UpdateType(ctx context.Context, item *skillz.Type) error {

	item.UpdatedAt = time.Now()

	query, args, err := sq.Update(r.types.table).
		SetMap(map[string]interface{}{
			TypesName:           item.Name,
			TypesGroupID:        item.GroupID,
			TypesPublished:      item.Published,
			TypesCapacity:       item.Capacity,
			TypesMarketGroupID:  item.MarketGroupID,
			TypesMass:           item.Mass,
			TypesPackagedVolume: item.PackagedVolume,
			TypesPortionSize:    item.PortionSize,
			TypesRadius:         item.Radius,
			TypesVolume:         item.Volume,
			ColumnUpdatedAt:     item.UpdatedAt,
		}).
		Where(sq.Eq{TypesID: item.ID}).ToSql()
	if err != nil {
		return errors.Wrapf(err, errorFFormat, universeRepository, "UpdateType", "failed to generate sql")
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	return errors.Wrapf(err, prefixFormat, universeRepository, "UpdateType")

}
