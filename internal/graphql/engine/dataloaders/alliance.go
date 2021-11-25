package dataloaders

import (
	"context"

	"github.com/eveisesi/skillz"
	"github.com/eveisesi/skillz/internal/graphql/engine/dataloaders/generated"
)

func (s *Service) AllianceLoader() *generated.AllianceLoader {
	return generated.NewAllianceLoader(generated.AllianceLoaderConfig{
		MaxBatch: s.batch,
		Wait:     s.wait,
		Fetch: func(ctx context.Context, keys []uint) ([]*skillz.Alliance, []error) {

			var errors = make([]error, len(keys))
			var results = make([]*skillz.Alliance, len(keys))

			for i, k := range keys {
				result, err := s.alliance.Alliance(ctx, k)
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
