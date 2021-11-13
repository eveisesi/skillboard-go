package skillz

import (
	"context"
	"time"

	"github.com/volatiletech/null"
)

type UniverseRepository interface {
	Bloodline(ctx context.Context, bloodlineID uint) (*Bloodline, error)
	Bloodlines(ctx context.Context) ([]*Bloodline, error)
	CreateBloodline(ctx context.Context, bloodline *Bloodline) error
	UpdateBloodline(ctx context.Context, bloodline *Bloodline) error

	Category(ctx context.Context, categoryID uint) (*Category, error)
	Categories(ctx context.Context) ([]*Category, error)
	CreateCategory(ctx context.Context, category *Category) error
	UpdateCategory(ctx context.Context, category *Category) error

	Constellation(ctx context.Context, constellationID uint) (*Constellation, error)
	Constellations(ctx context.Context) ([]*Constellation, error)
	CreateConstellation(ctx context.Context, constellation *Constellation) error
	UpdateConstellation(ctx context.Context, constellation *Constellation) error

	Faction(ctx context.Context, factionID uint) (*Faction, error)
	Factions(ctx context.Context) ([]*Faction, error)
	CreateFaction(ctx context.Context, faction *Faction) error
	UpdateFaction(ctx context.Context, faction *Faction) error

	Group(ctx context.Context, groupID uint) (*Group, error)
	Groups(ctx context.Context, operators ...*Operator) ([]*Group, error)
	CreateGroup(ctx context.Context, group *Group) error
	UpdateGroup(ctx context.Context, group *Group) error

	Race(ctx context.Context, raceID uint) (*Race, error)
	Races(ctx context.Context, operators ...*Operator) ([]*Race, error)
	CreateRace(ctx context.Context, race *Race) error
	UpdateRace(ctx context.Context, race *Race) error

	Region(ctx context.Context, regionID uint) (*Region, error)
	Regions(ctx context.Context) ([]*Region, error)
	CreateRegion(ctx context.Context, region *Region) error
	UpdateRegion(ctx context.Context, region *Region) error

	SolarSystem(ctx context.Context, solarSystemID uint) (*SolarSystem, error)
	SolarSystems(ctx context.Context, operators ...*Operator) ([]*SolarSystem, error)
	CreateSolarSystem(ctx context.Context, solarSystem *SolarSystem) error
	UpdateSolarSystem(ctx context.Context, solarSystem *SolarSystem) error

	Station(ctx context.Context, stationID uint) (*Station, error)
	Stations(ctx context.Context, operators ...*Operator) ([]*Station, error)
	CreateStation(ctx context.Context, station *Station) (*Station, error)
	UpdateStation(ctx context.Context, id uint, station *Station) (*Station, error)

	Structure(ctx context.Context, structureID uint64) (*Structure, error)
	Structures(ctx context.Context) ([]*Structure, error)
	CreateStructure(ctx context.Context, structure *Structure) (*Structure, error)
	UpdateStructure(ctx context.Context, id uint64, structure *Structure) (*Structure, error)

	Type(ctx context.Context, typeID uint) (*Type, error)
	Types(ctx context.Context) ([]*Type, error)
	CreateType(ctx context.Context, item *Type) error
	UpdateType(ctx context.Context, item *Type) error
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
