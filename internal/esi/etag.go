package esi

import (
	"context"

	"github.com/eveisesi/skillz"
	"github.com/pkg/errors"
)

type etags interface {
	Etag(ctx context.Context, endpointID EndpointID, params *Params) (string, *skillz.Etag, error)
}

func (s *Service) Etag(ctx context.Context, endpointID EndpointID, params *Params) (string, *skillz.Etag, error) {

	etagID, err := Resolvers[endpointID](params)
	if err != nil {
		return "", nil, errors.Wrap(err, "failed to generate etag ID")
	}

	etag, err := s.etag.Etag(ctx, etagID)
	return etagID, etag, err

}
