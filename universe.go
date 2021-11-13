package skillz

import (
	"context"
	"time"

	"github.com/volatiletech/null"
)

type UniverseRepository interface {
	bloodlineRepository
	categoryRepository
	constellationRepository
	factionRepository
	groupRepository
	raceRepository
	regionRepository
	solarSystemRepository
	stationRepository
	structureRepository
	typeRepository
}

type bloodlineRepository interface {
	Bloodline(ctx context.Context, id uint) (*Bloodline, error)
	Bloodlines(ctx context.Context, operators ...*Operator) ([]*Bloodline, error)
	CreateBloodline(ctx context.Context, bloodline *Bloodline) (*Bloodline, error)
	UpdateBloodline(ctx context.Context, id uint, bloodline *Bloodline) (*Bloodline, error)
	DeleteBloodline(ctx context.Context, id uint) (bool, error)
}

type categoryRepository interface {
	Category(ctx context.Context, id uint) (*Category, error)
	Categories(ctx context.Context, operators ...*Operator) ([]*Category, error)
	CreateCategory(ctx context.Context, group *Category) (*Category, error)
	UpdateCategory(ctx context.Context, id uint, category *Category) (*Category, error)
	DeleteCategory(ctx context.Context, id uint) (bool, error)
}

type constellationRepository interface {
	Constellation(ctx context.Context, id uint) (*Constellation, error)
	Constellations(ctx context.Context, operators ...*Operator) ([]*Constellation, error)
	CreateConstellation(ctx context.Context, constellation *Constellation) (*Constellation, error)
	UpdateConstellation(ctx context.Context, id uint, constellation *Constellation) (*Constellation, error)
	DeleteConstellation(ctx context.Context, id uint) (bool, error)
}

type factionRepository interface {
	Faction(ctx context.Context, id uint) (*Faction, error)
	Factions(ctx context.Context, operators ...*Operator) ([]*Faction, error)
	CreateFaction(ctx context.Context, faction *Faction) (*Faction, error)
	UpdateFaction(ctx context.Context, id uint, faction *Faction) (*Faction, error)
	DeleteFaction(ctx context.Context, id uint) (bool, error)
}

type groupRepository interface {
	Group(ctx context.Context, id uint) (*Group, error)
	Groups(ctx context.Context, operators ...*Operator) ([]*Group, error)
	CreateGroup(ctx context.Context, group *Group) (*Group, error)
	UpdateGroup(ctx context.Context, id uint, group *Group) (*Group, error)
	DeleteGroup(ctx context.Context, id uint) (bool, error)
}

type raceRepository interface {
	Race(ctx context.Context, id uint) (*Race, error)
	Races(ctx context.Context, operators ...*Operator) ([]*Race, error)
	CreateRace(ctx context.Context, race *Race) (*Race, error)
	UpdateRace(ctx context.Context, id uint, race *Race) (*Race, error)
	DeleteRace(ctx context.Context, id uint) (bool, error)
}

type regionRepository interface {
	Region(ctx context.Context, id uint) (*Region, error)
	Regions(ctx context.Context, operators ...*Operator) ([]*Region, error)
	CreateRegion(ctx context.Context, region *Region) (*Region, error)
	UpdateRegion(ctx context.Context, id uint, region *Region) (*Region, error)
	DeleteRegion(ctx context.Context, id uint) (bool, error)
}

type solarSystemRepository interface {
	SolarSystem(ctx context.Context, id uint) (*SolarSystem, error)
	SolarSystems(ctx context.Context, operators ...*Operator) ([]*SolarSystem, error)
	CreateSolarSystem(ctx context.Context, solarSystem *SolarSystem) (*SolarSystem, error)
	UpdateSolarSystem(ctx context.Context, id uint, solarSystem *SolarSystem) (*SolarSystem, error)
	DeleteSolarSystem(ctx context.Context, id uint) (bool, error)
}

type stationRepository interface {
	Station(ctx context.Context, id uint) (*Station, error)
	Stations(ctx context.Context, operators ...*Operator) ([]*Station, error)
	CreateStation(ctx context.Context, station *Station) (*Station, error)
	UpdateStation(ctx context.Context, id uint, solarSystem *Station) (*Station, error)
	DeleteStation(ctx context.Context, id uint) (bool, error)
}

type structureRepository interface {
	Structure(ctx context.Context, id uint64) (*Structure, error)
	Structures(ctx context.Context, operators ...*Operator) ([]*Structure, error)
	CreateStructure(ctx context.Context, solarSystem *Structure) (*Structure, error)
	UpdateStructure(ctx context.Context, id uint64, struture *Structure) (*Structure, error)
	DeleteStructure(ctx context.Context, id uint64) (bool, error)
}

type typeRepository interface {
	Type(ctx context.Context, id uint) (*Type, error)
	Types(ctx context.Context, operators ...*Operator) ([]*Type, error)
	CreateType(ctx context.Context, item *Type) (*Type, error)
	UpdateType(ctx context.Context, id uint, item *Type) (*Type, error)
	DeleteType(ctx context.Context, id uint) (bool, error)
}

type Bloodline struct {
	ID            uint      `db:"id" json:"bloodline_id"`
	Name          string    `db:"name" json:"name"`
	RaceID        uint      `db:"race_id" json:"race_id"`
	CorporationID uint      `db:"corporation_id" json:"corporation_id"`
	ShipTypeID    uint      `db:"ship_type_id" json:"ship_type_id"`
	Charisma      uint      `db:"charisma" json:"charisma"`
	Intelligence  uint      `db:"intelligence" json:"intelligence"`
	Memory        uint      `db:"memory" json:"memory"`
	Perception    uint      `db:"perception" json:"perception"`
	Willpower     uint      `db:"willpower" json:"willpower"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`
}

type Category struct {
	ID        uint      `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	Published bool      `db:"published" json:"published"`
	Groups    []uint    `db:"-" json:"groups,omitempty"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type Constellation struct {
	ID        uint      `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	RegionID  uint      `db:"region_id" json:"region_id"`
	SystemIDs []uint    `db:"-" json:"systems,omitempty"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type Faction struct {
	ID                   uint      `db:"id" json:"faction_id"`
	Name                 string    `db:"name" json:"name"`
	IsUnique             bool      `db:"is_unique" json:"is_unique"`
	SizeFactor           float64   `db:"size_factor" json:"size_factor"`
	StationCount         uint      `db:"station_count" json:"station_count"`
	StationSystemCount   uint      `db:"station_system_count" json:"station_system_count"`
	CorporationID        null.Uint `db:"corporation_id,omitempty" json:"corporation_id,omitempty"`
	MilitiaCorporationID null.Uint `db:"militia_corporation_id,omitempty" json:"militia_corporation_id,omitempty"`
	SolarSystemID        null.Uint `db:"solar_system_id,omitempty" json:"solar_system_id,omitempty"`
	CreatedAt            time.Time `db:"created_at" json:"created_at"`
	UpdatedAt            time.Time `db:"updated_at" json:"updated_at"`
}

func (Faction) IsContactInfo() {}

type Group struct {
	ID         uint      `db:"id" json:"id"`
	Name       string    `db:"name" json:"name"`
	Published  bool      `db:"published" json:"published"`
	CategoryID uint      `db:"category_id" json:"category_id"`
	Types      []uint    `db:"-" json:"types,omitempty"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time `db:"updated_at" json:"updated_at"`
}

type Race struct {
	ID        uint      `db:"id" json:"race_id"`
	Name      string    `db:"name" json:"name"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type Region struct {
	ID               uint      `db:"id" json:"id"`
	Name             string    `db:"name" json:"name"`
	ConstellationIDs []uint    `db:"-" json:"constellations,omitempty"`
	CreatedAt        time.Time `db:"created_at" json:"created_at"`
	UpdatedAt        time.Time `db:"updated_at" json:"updated_at"`
}

type SolarSystem struct {
	ID              uint        `db:"id" json:"id"`
	Name            string      `db:"name" json:"name"`
	ConstellationID uint        `db:"constellation_id" json:"constellation_id"`
	SecurityStatus  float64     `db:"security_status" json:"security_status"`
	StarID          null.Uint   `db:"star_id,omitempty" json:"star_id,omitempty"`
	SecurityClass   null.String `db:"security_class,omitempty" json:"security_class,omitempty"`
	CreatedAt       time.Time   `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time   `db:"updated_at" json:"updated_at"`
}

type Station struct {
	ID                       uint      `db:"id" json:"id"`
	Name                     string    `db:"name" json:"name"`
	SystemID                 uint      `db:"system_id" json:"system_id"`
	TypeID                   uint      `db:"type_id" json:"type_id"`
	RaceID                   null.Uint `db:"race_id" json:"race_id"`
	OwnerCorporationID       null.Uint `db:"owner_corporation_id" json:"owner"`
	MaxDockableShipVolume    float64   `db:"max_dockable_ship_volume" json:"max_dockable_ship_volume"`
	OfficeRentalCost         float64   `db:"office_rental_cost" json:"office_rental_cost"`
	ReprocessingEfficiency   float64   `db:"reprocessing_efficiency" json:"reprocessing_efficiency"`
	ReprocessingStationsTake float64   `db:"reprocessing_stations_take" json:"reprocessing_stations_take"`
	CreatedAt                time.Time `db:"created_at" json:"created_at"`
	UpdatedAt                time.Time `db:"updated_at" json:"updated_at"`
}

func (Station) IsCloneLocationInfo() {}

type Structure struct {
	ID            uint64    `db:"id" json:"id"`
	Name          string    `db:"name" json:"name"`
	OwnerID       uint      `db:"owner_id" json:"owner_id"`
	SolarSystemID uint      `db:"solar_system_id" json:"solar_system_id"`
	TypeID        null.Uint `db:"type_id" json:"type_id"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`
}

func (Structure) IsCloneLocationInfo() {}

type Type struct {
	ID             uint         `db:"id" json:"id"`
	Name           string       `db:"name" json:"name"`
	GroupID        uint         `db:"group_id" json:"group_id"`
	Published      bool         `db:"published" json:"published"`
	Capacity       float64      `db:"capacity,omitempty" json:"capacity,omitempty"`
	MarketGroupID  null.Uint    `db:"market_group_id,omitempty" json:"market_group_id,omitempty"`
	Mass           null.Float64 `db:"mass,omitempty" json:"mass,omitempty"`
	PackagedVolume null.Float64 `db:"packaged_volume,omitempty" json:"packaged_volume,omitempty"`
	PortionSize    null.Uint    `db:"portion_size,omitempty" json:"portion_size,omitempty"`
	Radius         null.Float64 `db:"radius,omitempty" json:"radius,omitempty"`
	Volume         float64      `db:"volume,omitempty" json:"volume,omitempty"`
	CreatedAt      time.Time    `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time    `db:"updated_at" json:"updated_at"`
}
