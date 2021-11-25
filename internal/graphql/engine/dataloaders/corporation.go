package dataloaders

import (
	"context"

	"github.com/eveisesi/skillz"
	"github.com/eveisesi/skillz/internal/graphql/engine/dataloaders/generated"
)

func (s *Service) CorporationLoader() *generated.CorporationLoader {
	return generated.NewCorporationLoader(generated.CorporationLoaderConfig{
		MaxBatch: s.batch,
		Wait:     s.wait,
		Fetch: func(ctx context.Context, keys []uint) ([]*skillz.Corporation, []error) {

			var errors = make([]error, len(keys))
			var results = make([]*skillz.Corporation, len(keys))

			for i, k := range keys {
				result, err := s.corporation.Corporation(ctx, k)
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
