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

	"github.com/eveisesi/skillz/internal/etag"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

type API interface {
	characters
	corporations
	alliance
}

type Service struct {
	client *http.Client
	redis  *redis.Client

	etag etag.API
}

// Compile Check
var _ API = new(Service)

const (
	esiHost               = "esi.evetech.net"
	headerTimestampFormat = "Mon, 02 Jan 2006 15:04:05 MST"
)

func New(client *http.Client, redis *redis.Client, etag etag.API) *Service {
	return &Service{
		client, redis, etag,
	}
}

type out struct {
	Data    interface{} `json:"data"`
	Headers http.Header `json:"headers"`
	Status  int         `json:"status"`
}

func (s *Service) request(ctx context.Context, method, path string, body io.Reader, expected int, out *out, mods ...ModifierFunc) error {

	var err error
	var res = new(http.Response)

	uri, _ := url.ParseRequestURI(path)
	uri.Scheme = "https"
	uri.Host = esiHost

	for {
		req, err := http.NewRequestWithContext(ctx, method, uri.String(), body)
		if err != nil {
			return errors.Wrap(err, "failed to build request")
		}

		for _, mod := range mods {
			err = mod(req, nil)
			if err != nil {
				return err
			}
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

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return errors.Wrapf(err, "expected status %d, got %d: unable to parse request body", expected, res.StatusCode)
	}

	if out.Status >= http.StatusBadRequest {
		return errors.Errorf("expected status %d, got %d: %s", expected, res.StatusCode, string(data))
	}

	if out.Status == http.StatusNotModified {
		return nil
	}

	for _, mod := range mods {
		err = mod(nil, res)
		if err != nil {
			return err
		}
	}

	err = json.Unmarshal(data, out.Data)
	if err != nil {
		return errors.Wrapf(err, "failed to decode request body to json: %s", string(data))
	}

	return nil

}

func hash(s string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(s)))
}
