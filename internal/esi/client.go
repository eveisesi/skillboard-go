package esi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
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

func (s *Service) request(ctx context.Context, method, path string, body []byte, expected int, out *out, mods ...ModifierFunc) error {

	uri, _ := url.ParseRequestURI(path)
	uri.Scheme = "https"
	uri.Host = esiHost

	for {

		req, err := http.NewRequest(method, uri.String(), bytes.NewReader(body))
		if err != nil {
			return errors.Wrap(err, "failed to build request")
		}

		for _, mod := range mods {
			err = mod(req, nil)
			if err != nil {
				return err
			}
		}

		res, err := s.client.Do(req)
		if err != nil {
			return errors.Wrap(err, "failed to execute request")
		}

		fmt.Printf("%s %s (%d)\n", method, path, res.StatusCode)

		out.Status = res.StatusCode
		out.Headers = res.Header

		if res.StatusCode >= http.StatusBadRequest {
			if res.StatusCode == http.StatusBadRequest {
				return errors.Wrap(err, "404 Bad Request")
			}

			if res.StatusCode == http.StatusUnauthorized {
				return errors.Wrap(err, "403 Unauthorized")
			}

			if res.StatusCode == http.StatusBadGateway {
				time.Sleep(time.Second)
				continue
			}

			if res.StatusCode == http.StatusTooManyRequests {
				waitStr := out.Headers.Get("x-esi-error-limit-reset")
				if waitStr == "" {
					time.Sleep(time.Second * 10)
					continue
				}

				wait, err := strconv.ParseUint(waitStr, 10, 32)
				if err != nil {
					fmt.Println(errors.Wrap(err, "failed to parse value in expires header to uint"))
					time.Sleep(time.Second * 10)
					continue
				}

				fmt.Printf("http 429, sleeping %d seconds", wait)
				time.Sleep(time.Second * time.Duration(wait))

			}
		}

		data, err := io.ReadAll(res.Body)
		if err != nil {
			return errors.Wrapf(err, "expected status %d, got %d: unable to parse request body", expected, res.StatusCode)
		}

		err = res.Body.Close()
		if err != nil {
			if err.Error() != "context cancelled" {
				fmt.Printf("\n\nfailed to close requst body for %s (%d): %s %s\n\n", fmt.Sprintf("%s %s", method, path), res.StatusCode, err, string(data))
			}
		}

		for _, mod := range mods {
			err = mod(nil, res)
			if err != nil {
				return err
			}
		}

		if out.Status == http.StatusNotModified {
			return nil
		}

		err = json.Unmarshal(data, out.Data)
		if err != nil {
			return errors.Wrapf(err, "failed to decode request body to json: %s", string(data))
		}

		return nil

	}

}

func hash(s string) string {
	// fmt.Sprintf("%x", sha256.Sum256([]byte(s)))
	return s
}
