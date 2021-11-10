package esi

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/eveisesi/skillz"
	"github.com/eveisesi/skillz/internal/etag"
	"github.com/go-http-utils/headers"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

type API interface {
	characters
}

type Service struct {
	client *http.Client
	redis  *redis.Client

	etag etag.API

	userAgent string
}

// Compile Check
var _ API = new(Service)

const (
	esiHost               = "esi.evetech.net"
	headerTimestampFormat = "Mon, 02 Jan 2006 15:04:05 MST"
)

func New(client *http.Client, redis *redis.Client, etag etag.API, userAgent string) *Service {
	return &Service{
		client, redis, etag, userAgent,
	}
}

type out struct {
	Data    interface{} `json:"data"`
	Headers http.Header `json:"headers"`
	Status  int         `json:"status"`
}

func (s *Service) request(ctx context.Context, method, path string, body io.Reader, expected int, out *out) error {

	var etag *skillz.Etag
	var err error
	var res = new(http.Response)
	if method == http.MethodGet {
		err = s.getResponseCache(ctx, path, out)
		if err == nil {
			return nil
		}

		etag, err = s.etag.Etag(ctx, path)
		if err != nil {
			return errors.Wrap(err, "failed to fetch etag for request")
		}

	}

	uri, _ := url.ParseRequestURI(path)
	uri.Scheme = "https"
	uri.Host = esiHost

	for {
		req, err := http.NewRequestWithContext(ctx, method, uri.String(), body)
		if err != nil {
			return errors.Wrap(err, "failed to build request")
		}

		req.Header.Add(headers.UserAgent, s.userAgent)

		if method == http.MethodGet && etag != nil {
			req.Header.Add(headers.IfNoneMatch, etag.Etag)
		}

		res, err = s.client.Do(req)
		if err != nil {
			return errors.Wrap(err, "failed to execute request")
		}

		if res != nil && res.StatusCode >= http.StatusContinue && res.StatusCode < http.StatusInternalServerError {
			break
		}

		// Sleep for 1/2 second
		time.Sleep(time.Millisecond * 500)

	}

	defer func(requestID string, body io.ReadCloser) {
		err := body.Close()
		if err != nil {
			fmt.Printf("failed to close requst body for %s\n", requestID)
		}
	}(fmt.Sprintf("%s %s", method, path), res.Body)

	out.Status = res.StatusCode
	out.Headers = res.Header

	if out.Status >= http.StatusBadRequest {
		data, err := io.ReadAll(res.Body)
		if err != nil {
			return errors.Wrapf(err, "expected status %d, got %d: unable to parse request body", expected, res.StatusCode)
		}

		return errors.Errorf("expected status %d, got %d: %s", expected, res.StatusCode, string(data))
	}

	var cacheDuration time.Duration
	expiresHeader := res.Header.Get(headers.Expires)
	if expiresHeader != "" {
		expires, err := time.Parse(headerTimestampFormat, expiresHeader)
		if err == nil {
			cacheDuration = time.Until(expires)
		}
		etagHeader := res.Header.Get(headers.ETag)
		if etagHeader != "" {
			etag.Path = path
			etag.Etag = etagHeader
			etag.CachedUntil = expires

			s.etag.InsertEtag(ctx, etag)

		}
	}

	if out.Status == http.StatusNotModified {
		return nil
	}

	err = json.NewDecoder(res.Body).Decode(out.Data)
	if err != nil {
		return errors.Wrap(err, "failed to decode request body to json")
	}

	if method == http.MethodGet {
		_ = s.cacheResponse(ctx, path, cacheDuration, out)
	}

	return nil

}

func (s *Service) cacheResponse(ctx context.Context, path string, duration time.Duration, out *out) error {

	payload, err := json.Marshal(out)
	if err != nil {
		return errors.Wrap(err, "failed to cache response for cache")
	}

	_, err = s.redis.Set(ctx, hash(path), string(payload), duration).Result()
	if err != nil {
		return errors.Wrap(err, "failed to write response body to cache layer")
	}

	return nil

}

func (s *Service) getResponseCache(ctx context.Context, path string, out *out) error {

	b, err := s.redis.Get(ctx, hash(path)).Bytes()
	if err != nil {
		return err
	}

	return json.Unmarshal(b, out)

}

func hash(s string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(s)))
}
