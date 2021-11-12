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
	GetFactions
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
	GetCharacter:                   resolverFuncs["characterID"](GetCharacter),
	GetCharacterCorporationHistory: resolverFuncs["characterID"](GetCharacterCorporationHistory),
	GetCorporation:                 resolverFuncs["corporationID"](GetCorporation),
	GetCorporationAllianceHistory:  resolverFuncs["corporationID"](GetCorporationAllianceHistory),
	GetAlliance:                    resolverFuncs["corporationID"](GetAlliance),
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
	GetCharacterCorporationHistory: "/v1/characters/%d/corporationhistory/",
	GetCorporation:                 "/v5/corporations/%d/",
	GetCorporationAllianceHistory:  "/v3/corporations/%d/alliancehistory/",
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
}
