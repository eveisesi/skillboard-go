package dataloaders

import (
	"context"

	"github.com/eveisesi/skillz"
	"github.com/eveisesi/skillz/internal/graphql/engine/dataloaders/generated"
)

func (s *Service) CloneLoader() *generated.CloneLoader {
	return generated.NewCloneLoader(generated.CloneLoaderConfig{
		MaxBatch: s.batch,
		Wait:     s.wait,
		Fetch: func(ctx context.Context, keys []*skillz.User) ([]*skillz.CharacterCloneMeta, []error) {
			var errors = make([]error, len(keys))
			var results = make([]*skillz.CharacterCloneMeta, len(keys))

			for i, k := range keys {
				result, err := s.clone.Clones(ctx, k)
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

func (s *Service) ImplantLoader() *generated.ImplantLoader {
	return generated.NewImplantLoader(generated.ImplantLoaderConfig{
		MaxBatch: s.batch,
		Wait:     s.wait,
		Fetch: func(ctx context.Context, keys []*skillz.User) ([][]*skillz.CharacterImplant, []error) {
			var errors = make([]error, len(keys))
			var results = make([][]*skillz.CharacterImplant, len(keys))

			for i, k := range keys {
				result, err := s.clone.Implants(ctx, k)
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
