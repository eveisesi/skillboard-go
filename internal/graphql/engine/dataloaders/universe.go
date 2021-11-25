package dataloaders

import (
	"context"

	"github.com/eveisesi/skillz"
	"github.com/eveisesi/skillz/internal/graphql/engine/dataloaders/generated"
)

func (s *Service) StationLoader() *generated.StationLoader {
	return generated.NewStationLoader(generated.StationLoaderConfig{
		MaxBatch: s.batch,
		Wait:     s.wait,
		Fetch: func(ctx context.Context, keys []uint) ([]*skillz.Station, []error) {
			var errors = make([]error, len(keys))
			var results = make([]*skillz.Station, len(keys))

			for i, k := range keys {
				result, err := s.universe.Station(ctx, k)
				if err != nil {
					errors[i] = err
					continue
				}

				results[i] = result
			}

			return results, errors

		},
	})
}

func (s *Service) StructureLoader() *generated.StructureLoader {
	return generated.NewStructureLoader(generated.StructureLoaderConfig{
		MaxBatch: s.batch,
		Wait:     s.wait,
		Fetch: func(ctx context.Context, keys []uint64) ([]*skillz.Structure, []error) {
			var errors = make([]error, len(keys))
			var results = make([]*skillz.Structure, len(keys))

			for i, k := range keys {
				result, err := s.universe.Structure(ctx, k)
				if err != nil {
					errors[i] = err
					continue
				}

				results[i] = result
			}

			return results, errors
		},
	})
}

func (s *Service) RegionLoader() *generated.RegionLoader {
	return generated.NewRegionLoader(generated.RegionLoaderConfig{
		MaxBatch: s.batch,
		Wait:     s.wait,
		Fetch: func(ctx context.Context, keys []uint) ([]*skillz.Region, []error) {
			var errors = make([]error, len(keys))
			var results = make([]*skillz.Region, len(keys))

			for i, k := range keys {
				result, err := s.universe.Region(ctx, k)
				if err != nil {
					errors[i] = err
					continue
				}

				results[i] = result
			}

			return results, errors
		},
	})
}

func (s *Service) ConstellationLoader() *generated.ConstellationLoader {
	return generated.NewConstellationLoader(generated.ConstellationLoaderConfig{
		MaxBatch: s.batch,
		Wait:     s.wait,
		Fetch: func(ctx context.Context, keys []uint) ([]*skillz.Constellation, []error) {
			var errors = make([]error, len(keys))
			var results = make([]*skillz.Constellation, len(keys))

			for i, k := range keys {
				result, err := s.universe.Constellation(ctx, k)
				if err != nil {
					errors[i] = err
					continue
				}

				results[i] = result
			}

			return results, errors
		},
	})
}

func (s *Service) SolarSystemLoader() *generated.SolarSystemLoader {
	return generated.NewSolarSystemLoader(generated.SolarSystemLoaderConfig{
		MaxBatch: s.batch,
		Wait:     s.wait,
		Fetch: func(ctx context.Context, keys []uint) ([]*skillz.SolarSystem, []error) {
			var errors = make([]error, len(keys))
			var results = make([]*skillz.SolarSystem, len(keys))

			for i, k := range keys {
				result, err := s.universe.SolarSystem(ctx, k)
				if err != nil {
					errors[i] = err
					continue
				}

				results[i] = result
			}

			return results, errors
		},
	})
}

func (s *Service) TypeLoader() *generated.TypeLoader {
	return generated.NewTypeLoader(generated.TypeLoaderConfig{
		MaxBatch: s.batch,
		Wait:     s.wait,
		Fetch: func(ctx context.Context, keys []uint) ([]*skillz.Type, []error) {
			var errors = make([]error, len(keys))
			var results = make([]*skillz.Type, len(keys))

			for i, k := range keys {
				result, err := s.universe.Type(ctx, k)
				if err != nil {
					errors[i] = err
					continue
				}

				results[i] = result
			}

			return results, errors
		},
	})
}

func (s *Service) GroupLoader() *generated.GroupLoader {
	return generated.NewGroupLoader(generated.GroupLoaderConfig{
		MaxBatch: s.batch,
		Wait:     s.wait,
		Fetch: func(ctx context.Context, keys []uint) ([]*skillz.Group, []error) {
			var errors = make([]error, len(keys))
			var results = make([]*skillz.Group, len(keys))

			for i, k := range keys {
				result, err := s.universe.Group(ctx, k)
				if err != nil {
					errors[i] = err
					continue
				}

				results[i] = result
			}

			return results, errors
		},
	})
}

func (s *Service) CategoryLoader() *generated.CategoryLoader {
	return generated.NewCategoryLoader(generated.CategoryLoaderConfig{
		MaxBatch: s.batch,
		Wait:     s.wait,
		Fetch: func(ctx context.Context, keys []uint) ([]*skillz.Category, []error) {
			var errors = make([]error, len(keys))
			var results = make([]*skillz.Category, len(keys))

			for i, k := range keys {
				result, err := s.universe.Category(ctx, k)
				if err != nil {
					errors[i] = err
					continue
				}

				results[i] = result
			}

			return results, errors
		},
	})
}
