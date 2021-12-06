package esi

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/eveisesi/skillz"
	"github.com/go-http-utils/headers"
	"github.com/pkg/errors"
)

type modifiers interface {
	AddIfNoneMatchHeader(ctx context.Context, etag string) ModifierFunc
	AddAuthorizationHeader(ctx context.Context, token string) ModifierFunc
	CacheEtag(ctx context.Context, hash string, expiration *time.Time) ModifierFunc
	BaseCharacterModifiers(ctx context.Context, user *skillz.User, etagID string, etag *skillz.Etag) []ModifierFunc
}

type ModifierFunc func(req *http.Request, res *http.Response) error

func (s *Service) AddIfNoneMatchHeader(ctx context.Context, etag string) ModifierFunc {
	return func(req *http.Request, res *http.Response) error {
		if req == nil {
			return nil
		}

		if etag == "" {
			return nil
		}

		req.Header.Set(headers.IfNoneMatch, etag)
		return nil
	}
}

func (s *Service) AddAuthorizationHeader(ctx context.Context, token string) ModifierFunc {
	return func(req *http.Request, res *http.Response) error {
		if req == nil {
			return nil
		}

		if token == "" {
			return nil
		}

		req.Header.Set(headers.Authorization, fmt.Sprintf("Bearer %s", token))
		return nil
	}
}

func (s *Service) CacheEtag(ctx context.Context, hash string, expiration *time.Time) ModifierFunc {
	return func(req *http.Request, res *http.Response) error {
		if res == nil {
			return nil
		}

		duration := time.Now().Add(time.Hour)
		if header := res.Header.Get(headers.Expires); header != "" {

			d, err := time.Parse(headerTimestampFormat, header)
			if err != nil {
				return errors.Wrap(err, "failed to parse ESI Expires header")
			}

			if expiration != nil && expiration.Unix() > d.Unix() {
				d = *expiration
			}

			duration = d
		}

		if header := res.Header.Get(headers.ETag); header != "" {
			etag := &skillz.Etag{
				Path:        hash,
				Etag:        header,
				CachedUntil: duration,
			}

			err := s.etag.InsertEtag(ctx, etag)
			if err != nil {
				return errors.Wrap(err, "failed to write etag to data store")
			}
		}

		return nil
	}
}

func (s *Service) BaseCharacterModifiers(ctx context.Context, user *skillz.User, etagID string, etag *skillz.Etag) []ModifierFunc {
	mods := append(make([]ModifierFunc, 0, 3), s.CacheEtag(ctx, etagID, nil), s.AddAuthorizationHeader(ctx, user.AccessToken))
	if etag != nil && etag.Etag != "" {
		mods = append(mods, s.AddIfNoneMatchHeader(ctx, etag.Etag))
	}

	return mods
}
