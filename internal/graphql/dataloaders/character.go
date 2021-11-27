package dataloaders

import (
	"context"

	"github.com/eveisesi/skillz"
	"github.com/eveisesi/skillz/internal/graphql/dataloaders/generated"
)

func (s *Service) CharacterLoader() *generated.CharacterLoader {
	return generated.NewCharacterLoader(generated.CharacterLoaderConfig{
		MaxBatch: s.batch,
		Wait:     s.wait,
		Fetch: func(ctx context.Context, keys []uint64) ([]*skillz.Character, []error) {

			var errors = make([]error, len(keys))
			var results = make([]*skillz.Character, len(keys))

			for i, k := range keys {
				result, err := s.character.Character(ctx, k)
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
