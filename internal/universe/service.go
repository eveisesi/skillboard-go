package universe

import (
	"context"
	"database/sql"
	"time"

	"github.com/eveisesi/skillz"
	"github.com/eveisesi/skillz/internal"
	"github.com/eveisesi/skillz/internal/cache"
	"github.com/eveisesi/skillz/internal/esi"
	"github.com/pkg/errors"
	"github.com/volatiletech/null"
)

type API interface {
	Bloodline(ctx context.Context, bloodlineID uint) (*skillz.Bloodline, error)
	Category(ctx context.Context, categoryID uint) (*skillz.Category, error)
	Constellation(ctx context.Context, constellationID uint) (*skillz.Constellation, error)
	Faction(ctx context.Context, id uint) (*skillz.Faction, error)
	Group(ctx context.Context, groupID uint) (*skillz.Group, error)
	Race(ctx context.Context, id uint) (*skillz.Race, error)
	Region(ctx context.Context, regionID uint) (*skillz.Region, error)
	SolarSystem(ctx context.Context, solarSystemID uint) (*skillz.SolarSystem, error)
	Station(ctx context.Context, stationID uint) (*skillz.Station, error)
	Structure(ctx context.Context, structureID uint64) (*skillz.Structure, error)
	Type(ctx context.Context, itemID uint) (*skillz.Type, error)
}

type Service struct {
	cache cache.UniverseAPI
	esi   esi.UniverseAPI

	universe skillz.UniverseRepository
}

func New(cache cache.UniverseAPI, esi esi.UniverseAPI, universe skillz.UniverseRepository) *Service {
	return &Service{

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
	if err != nil {
		return nil, err
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

	etagID, etag, err := s.esi.Etag(ctx, esi.GetCategory, &esi.Params{CategoryID: null.UintFrom(categoryID)})
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch etag for expiry check")
	}

	category, err = s.universe.Category(ctx, categoryID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, errors.Wrap(err, "failed to fetch category from data store")
	}

	exists := err == nil
	if !exists || etag == nil || etag.CachedUntil.Unix() < time.Now().Add(-1*time.Minute).Unix() {
		mods := append(make([]esi.ModifierFunc, 0, 2), s.esi.CacheEtag(ctx, etagID))
		if etag != nil && etag.Etag != "" {
			mods = append(mods, s.esi.AddIfNoneMatchHeader(ctx, etag.Etag))
		}

		updatedCategory, err := s.esi.GetCategory(ctx, categoryID, mods...)
		if err != nil {
			return nil, err
		}

		if updatedCategory != nil {
			switch exists {
			case true:
				err = s.universe.UpdateCategory(ctx, updatedCategory)
				if err != nil {
					return nil, errors.Wrap(err, "failed to save character to data store")
				}
			case false:
				err = s.universe.CreateCategory(ctx, updatedCategory)
				if err != nil {
					return nil, errors.Wrap(err, "failed to save character to data store")
				}

			}

			category = updatedCategory
		}

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

	etagID, etag, err := s.esi.Etag(ctx, esi.GetConstellation, &esi.Params{ConstellationID: null.UintFrom(constellationID)})
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch etag for expiry check")
	}

	constellation, err = s.universe.Constellation(ctx, constellationID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, errors.Wrap(err, "failed to fetch constellation from data store")
	}

	exists := err == nil
	if !exists || etag == nil || etag.CachedUntil.Unix() < time.Now().Add(-1*time.Minute).Unix() {
		mods := append(make([]esi.ModifierFunc, 0, 2), s.esi.CacheEtag(ctx, etagID))
		if etag != nil && etag.Etag != "" {
			mods = append(mods, s.esi.AddIfNoneMatchHeader(ctx, etag.Etag))
		}

		updatedConstellation, err := s.esi.GetConstellation(ctx, constellationID, mods...)
		if err != nil {
			return nil, err
		}

		if updatedConstellation != nil {
			switch exists {
			case true:
				err = s.universe.UpdateConstellation(ctx, updatedConstellation)
				if err != nil {
					return nil, errors.Wrap(err, "failed to save character to data store")
				}
			case false:
				err = s.universe.CreateConstellation(ctx, updatedConstellation)
				if err != nil {
					return nil, errors.Wrap(err, "failed to save character to data store")
				}

			}

			constellation = updatedConstellation
		}

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
	if err != nil {
		return nil, err
	}

	err = s.cache.SetFaction(ctx, faction)

	return faction, err
}

func (s *Service) Group(ctx context.Context, groupID uint) (*skillz.Group, error) {

	group, err := s.cache.Group(ctx, groupID)
	if err != nil {
		return nil, err
	}

	if group != nil {
		return group, nil
	}

	etagID, etag, err := s.esi.Etag(ctx, esi.GetGroup, &esi.Params{GroupID: null.UintFrom(groupID)})
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch etag for expiry check")
	}

	group, err = s.universe.Group(ctx, groupID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, errors.Wrap(err, "failed to fetch group from data store")
	}

	exists := err == nil
	if !exists || etag == nil || etag.CachedUntil.Unix() < time.Now().Add(-1*time.Minute).Unix() {
		mods := append(make([]esi.ModifierFunc, 0, 2), s.esi.CacheEtag(ctx, etagID))
		if etag != nil && etag.Etag != "" {
			mods = append(mods, s.esi.AddIfNoneMatchHeader(ctx, etag.Etag))
		}

		updatedGroup, err := s.esi.GetGroup(ctx, groupID, mods...)
		if err != nil {
			return nil, err
		}

		if updatedGroup != nil {
			switch exists {
			case true:
				err = s.universe.UpdateGroup(ctx, updatedGroup)
				if err != nil {
					return nil, errors.Wrap(err, "failed to save character to data store")
				}
			case false:
				err = s.universe.CreateGroup(ctx, updatedGroup)
				if err != nil {
					return nil, errors.Wrap(err, "failed to save character to data store")
				}

			}

			group = updatedGroup
		}

	}

	return group, s.cache.SetGroup(ctx, group)

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
	if err != nil {
		return nil, err
	}

	err = s.cache.SetRace(ctx, race)

	return race, err

}

func (s *Service) Region(ctx context.Context, regionID uint) (*skillz.Region, error) {

	region, err := s.cache.Region(ctx, regionID)
	if err != nil {
		return nil, err
	}

	if region != nil {
		return region, nil
	}

	etagID, etag, err := s.esi.Etag(ctx, esi.GetRegion, &esi.Params{RegionID: null.UintFrom(regionID)})
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch etag for expiry check")
	}

	region, err = s.universe.Region(ctx, regionID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, errors.Wrap(err, "failed to fetch region from data store")
	}

	exists := err == nil
	if !exists || etag == nil || etag.CachedUntil.Unix() < time.Now().Add(-1*time.Minute).Unix() {
		mods := append(make([]esi.ModifierFunc, 0, 2), s.esi.CacheEtag(ctx, etagID))
		if etag != nil && etag.Etag != "" {
			mods = append(mods, s.esi.AddIfNoneMatchHeader(ctx, etag.Etag))
		}

		updatedRegion, err := s.esi.GetRegion(ctx, regionID, mods...)
		if err != nil {
			return nil, err
		}

		if updatedRegion != nil {
			switch exists {
			case true:
				err = s.universe.UpdateRegion(ctx, updatedRegion)
				if err != nil {
					return nil, errors.Wrap(err, "failed to save character to data store")
				}
			case false:
				err = s.universe.CreateRegion(ctx, updatedRegion)
				if err != nil {
					return nil, errors.Wrap(err, "failed to save character to data store")
				}

			}

			region = updatedRegion
		}

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

	etagID, etag, err := s.esi.Etag(ctx, esi.GetSolarSystem, &esi.Params{SolarSystemID: null.UintFrom(solarSystemID)})
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch etag for expiry check")
	}

	solarSystem, err = s.universe.SolarSystem(ctx, solarSystemID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, errors.Wrap(err, "failed to fetch solarSystem from data store")
	}

	exists := err == nil
	if !exists || etag == nil || etag.CachedUntil.Unix() < time.Now().Add(-1*time.Minute).Unix() {
		mods := append(make([]esi.ModifierFunc, 0, 2), s.esi.CacheEtag(ctx, etagID))
		if etag != nil && etag.Etag != "" {
			mods = append(mods, s.esi.AddIfNoneMatchHeader(ctx, etag.Etag))
		}

		updatedSolarSystem, err := s.esi.GetSolarSystem(ctx, solarSystemID, mods...)
		if err != nil {
			return nil, err
		}

		if updatedSolarSystem != nil {
			switch exists {
			case true:
				err = s.universe.UpdateSolarSystem(ctx, updatedSolarSystem)
				if err != nil {
					return nil, errors.Wrap(err, "failed to save character to data store")
				}
			case false:
				err = s.universe.CreateSolarSystem(ctx, updatedSolarSystem)
				if err != nil {
					return nil, errors.Wrap(err, "failed to save character to data store")
				}

			}

			solarSystem = updatedSolarSystem
		}

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

	etagID, etag, err := s.esi.Etag(ctx, esi.GetStation, &esi.Params{StationID: null.UintFrom(stationID)})
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch etag for expiry check")
	}

	station, err = s.universe.Station(ctx, stationID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, errors.Wrap(err, "failed to fetch station from data store")
	}

	exists := err == nil
	if !exists || etag == nil || etag.CachedUntil.Unix() < time.Now().Add(-1*time.Minute).Unix() {
		mods := append(make([]esi.ModifierFunc, 0, 2), s.esi.CacheEtag(ctx, etagID))
		if etag != nil && etag.Etag != "" {
			mods = append(mods, s.esi.AddIfNoneMatchHeader(ctx, etag.Etag))
		}

		updatedStation, err := s.esi.GetStation(ctx, stationID, mods...)
		if err != nil {
			return nil, err
		}

		if updatedStation != nil {
			switch exists {
			case true:
				err = s.universe.UpdateStation(ctx, updatedStation)
				if err != nil {
					return nil, errors.Wrap(err, "failed to save character to data store")
				}
			case false:
				err = s.universe.CreateStation(ctx, updatedStation)
				if err != nil {
					return nil, errors.Wrap(err, "failed to save character to data store")
				}

			}

			station = updatedStation
		}

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
			return nil, errors.Wrap(err, "failed to resolve structure id, no user available to query esi with")
		}

		mods := append(make([]esi.ModifierFunc, 0, 2), s.esi.CacheEtag(ctx, etagID), s.esi.AddAuthorizationHeader(ctx, user.AccessToken))
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

func (s *Service) Type(ctx context.Context, itemID uint) (*skillz.Type, error) {

	item, err := s.cache.Type(ctx, itemID)
	if err != nil {
		return nil, err
	}

	if item != nil {
		return item, nil
	}

	etagID, etag, err := s.esi.Etag(ctx, esi.GetType, &esi.Params{ItemID: null.UintFrom(itemID)})
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch etag for expiry check")
	}

	item, err = s.universe.Type(ctx, itemID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, errors.Wrap(err, "failed to fetch item from data store")
	}

	exists := err == nil
	if !exists || etag == nil || etag.CachedUntil.Unix() < time.Now().Add(-1*time.Minute).Unix() {
		mods := append(make([]esi.ModifierFunc, 0, 2), s.esi.CacheEtag(ctx, etagID))
		if etag != nil && etag.Etag != "" {
			mods = append(mods, s.esi.AddIfNoneMatchHeader(ctx, etag.Etag))
		}

		updatedType, err := s.esi.GetType(ctx, itemID, mods...)
		if err != nil {
			return nil, err
		}

		if updatedType != nil {
			switch exists {
			case true:
				err = s.universe.UpdateType(ctx, updatedType)
				if err != nil {
					return nil, errors.Wrap(err, "failed to save character to data store")
				}
			case false:
				err = s.universe.CreateType(ctx, updatedType)
				if err != nil {
					return nil, errors.Wrap(err, "failed to save character to data store")
				}

			}

			item = updatedType
		}

	}

	return item, s.cache.SetType(ctx, item)

}
