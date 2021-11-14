package esi

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/volatiletech/null"
)

type (
	Params struct {
		AllianceID      null.Uint
		CategoryID      null.Uint
		CharacterID     null.Uint64
		ConstellationID null.Uint
		ContractID      null.Uint
		CorporationID   null.Uint
		From            null.Uint64
		GroupID         null.Uint
		MailID          null.Uint
		ItemID          null.Uint
		LastMailID      null.Uint64
		Page            null.Uint
		RegionID        null.Uint
		StationID       null.Uint
		SolarSystemID   null.Uint
		StructureID     null.Uint64
	}

	EndpointID        uint
	resolverFunc      func(params *Params) (string, error)
	endpointResolvers map[EndpointID]resolverFunc
	endpointMap       map[EndpointID]string
)

const (
	GetAlliance EndpointID = iota
	GetCharacter
	GetCharacterCorporationHistory
	GetCharacterAttributes
	GetCharacterAssets
	GetCharacterSkills
	GetCharacterSkillQueue
	GetCharacterClones
	GetCharacterImplants
	GetCharacterContacts
	GetCharacterContactLabels
	GetCharacterContracts
	GetCharacterContractItems
	GetCharacterContractBids
	GetCharacterFittings
	GetCharacterLocation
	GetCharacterMailHeaders
	GetCharacterMailHeader
	GetCharacterMailLabels
	GetCharacterMailLists
	GetCharacterOnline
	GetCharacterShip
	GetCharacterWalletBalance
	GetCharacterWalletTransactions
	GetCharacterWalletJournal
	GetCorporation
	GetCorporationAllianceHistory
	GetAncestries
	GetAsteroidBelt
	GetBloodlines
	GetCategories
	GetCategory
	GetConstellation
	GetConstellations
	GetFactions
	GetGroups
	GetGroup
	GetMoon
	GetPlanet
	GetRaces
	GetRegions
	GetRegion
	GetSolarSystem
	GetStation
	GetStructure
	GetType
	PostUniverseNames
)

var Resolvers = endpointResolvers{
	GetAlliance: resolverFuncs["allianceID"](GetAlliance),

	GetCharacter:                   resolverFuncs["characterID"](GetCharacter),
	GetCharacterCorporationHistory: resolverFuncs["characterID"](GetCharacterCorporationHistory),
	GetCharacterClones:             resolverFuncs["characterID"](GetCharacterClones),
	GetCharacterImplants:           resolverFuncs["characterID"](GetCharacterImplants),
	GetCharacterSkills:             resolverFuncs["characterID"](GetCharacterSkills),
	GetCharacterSkillQueue:         resolverFuncs["characterID"](GetCharacterSkillQueue),
	GetCharacterAttributes:         resolverFuncs["characterID"](GetCharacterAttributes),

	GetCorporation:                resolverFuncs["corporationID"](GetCorporation),
	GetCorporationAllianceHistory: resolverFuncs["corporationID"](GetCorporationAllianceHistory),

	GetRegion:        resolverFuncs["regionID"](GetRegion),
	GetConstellation: resolverFuncs["constellationID"](GetConstellation),
	GetSolarSystem:   resolverFuncs["solarSystemID"](GetSolarSystem),
	GetStation:       resolverFuncs["stationID"](GetStation),
	GetStructure:     resolverFuncs["structureID"](GetStructure),

	GetCategory: resolverFuncs["categoryID"](GetCategory),
	GetGroup:    resolverFuncs["groupID"](GetGroup),
	GetType:     resolverFuncs["typeID"](GetType),
}

var ErrNilParams = errors.New("received nil for params")

type ErrInvalidParameter struct {
	Parameter string
}

func (e ErrInvalidParameter) Error() string {
	return fmt.Sprintf("%s required", e.Parameter)
}

var endpoints = endpointMap{
	GetAlliance:                    "/v4/alliances/%d/",
	GetCharacter:                   "/v5/characters/%d/",
	GetCharacterClones:             "/v4/characters/%d/clones/",
	GetCharacterImplants:           "/v2/characters/%d/implants/",
	GetCharacterSkills:             "/v2/characters/%d/skills/",
	GetCharacterSkillQueue:         "/v2/characters/%d/skillqueue/",
	GetCharacterAttributes:         "/v2/characters/%d/attributes/",
	GetCharacterCorporationHistory: "/v1/characters/%d/corporationhistory/",
	GetCorporation:                 "/v5/corporations/%d/",
	GetCorporationAllianceHistory:  "/v3/corporations/%d/alliancehistory/",

	GetBloodlines: "/v1/universe/bloodlines/",
	GetRaces:      "/v1/universe/races/",

	GetRegions:        "/v1/universe/regions/",
	GetRegion:         "/v1/universe/regions/%d/",
	GetConstellations: "/v1/universe/constellations/",
	GetConstellation:  "/v1/universe/constellations/%d/",
	GetSolarSystem:    "/v4/universe/systems/%d/",
	GetStation:        "/v2/universe/stations/%d/",
	GetStructure:      "/v2/universe/structures/%d/",

	GetCategories: "/v1/universe/categories/",
	GetCategory:   "/v1/universe/categories/%d/",
	GetGroups:     "/v1/universe/groups/",
	GetGroup:      "/v1/universe/groups/%d/",
	GetType:       "/v3/universe/types/%d/",
}

var resolverFuncs = map[string]func(endpoint EndpointID) resolverFunc{
	"characterID": func(endpoint EndpointID) resolverFunc {
		return func(params *Params) (string, error) {
			if params == nil {
				return "", ErrNilParams
			}
			if !params.CharacterID.Valid {
				return "", ErrInvalidParameter{"characterID"}
			}

			path := endpoints[endpoint]

			return hash(fmt.Sprintf(path, params.CharacterID.Uint64)), nil
		}
	},
	"corporationID": func(endpoint EndpointID) resolverFunc {
		return func(params *Params) (string, error) {
			if params == nil {
				return "", ErrNilParams
			}
			if !params.CorporationID.Valid {
				return "", ErrInvalidParameter{"corporationID"}
			}

			path := endpoints[endpoint]

			return hash(fmt.Sprintf(path, params.CorporationID.Uint)), nil

		}
	},
	"allianceID": func(endpoint EndpointID) resolverFunc {
		return func(params *Params) (string, error) {
			if params == nil {
				return "", ErrNilParams
			}
			if !params.AllianceID.Valid {
				return "", ErrInvalidParameter{"allianceID"}
			}

			path := endpoints[endpoint]

			return hash(fmt.Sprintf(path, params.AllianceID.Uint)), nil

		}
	},
	"categoryID": func(endpoint EndpointID) resolverFunc {
		return func(params *Params) (string, error) {
			if params == nil {
				return "", ErrNilParams
			}
			if !params.CategoryID.Valid {
				return "", ErrInvalidParameter{"categoryID"}
			}

			path := endpoints[endpoint]

			return hash(fmt.Sprintf(path, params.CategoryID.Uint)), nil

		}
	},
	"groupID": func(endpoint EndpointID) resolverFunc {
		return func(params *Params) (string, error) {
			if params == nil {
				return "", ErrNilParams
			}
			if !params.GroupID.Valid {
				return "", ErrInvalidParameter{"groupID"}
			}

			path := endpoints[endpoint]

			return hash(fmt.Sprintf(path, params.GroupID.Uint)), nil

		}
	},
	"regionID": func(endpoint EndpointID) resolverFunc {
		return func(params *Params) (string, error) {
			if params == nil {
				return "", ErrNilParams
			}
			if !params.RegionID.Valid {
				return "", ErrInvalidParameter{"regionID"}
			}

			path := endpoints[endpoint]

			return hash(fmt.Sprintf(path, params.RegionID.Uint)), nil

		}
	},
	"constellationID": func(endpoint EndpointID) resolverFunc {
		return func(params *Params) (string, error) {
			if params == nil {
				return "", ErrNilParams
			}
			if !params.ConstellationID.Valid {
				return "", ErrInvalidParameter{"constellationID"}
			}

			path := endpoints[endpoint]

			return hash(fmt.Sprintf(path, params.ConstellationID.Uint)), nil

		}
	},
	"solarSystemID": func(endpoint EndpointID) resolverFunc {
		return func(params *Params) (string, error) {
			if params == nil {
				return "", ErrNilParams
			}
			if !params.SolarSystemID.Valid {
				return "", ErrInvalidParameter{"solarSystemID"}
			}

			path := endpoints[endpoint]

			return hash(fmt.Sprintf(path, params.SolarSystemID.Uint)), nil

		}
	},
	"structureID": func(endpoint EndpointID) resolverFunc {
		return func(params *Params) (string, error) {
			if params == nil {
				return "", ErrNilParams
			}
			if !params.StructureID.Valid {
				return "", ErrInvalidParameter{"structureID"}
			}

			path := endpoints[endpoint]

			return hash(fmt.Sprintf(path, params.StructureID.Uint64)), nil

		}
	},
	"stationID": func(endpoint EndpointID) resolverFunc {
		return func(params *Params) (string, error) {
			if params == nil {
				return "", ErrNilParams
			}
			if !params.StationID.Valid {
				return "", ErrInvalidParameter{"stationID"}
			}

			path := endpoints[endpoint]

			return hash(fmt.Sprintf(path, params.StationID.Uint)), nil

		}
	},
	"typeID": func(endpoint EndpointID) resolverFunc {
		return func(params *Params) (string, error) {
			if params == nil {
				return "", ErrNilParams
			}
			if !params.ItemID.Valid {
				return "", ErrInvalidParameter{"typeID"}
			}

			path := endpoints[endpoint]

			return hash(fmt.Sprintf(path, params.ItemID.Uint)), nil

		}
	},
}
