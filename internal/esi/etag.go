package esi

import (
	"context"

	"github.com/eveisesi/skillz"
	"github.com/pkg/errors"
)

func (s *Service) Etag(ctx context.Context, endpointID EndpointID, params *Params) (*skillz.Etag, error) {

	etagID, err := Resolvers[endpointID](params)
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate etag ID")
	}

	return s.etag.Etag(ctx, etagID)

}
