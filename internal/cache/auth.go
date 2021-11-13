package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/eveisesi/skillz"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

type AuthAPI interface {
	JSONWebKeySet(ctx context.Context) ([]byte, error)
	SaveJSONWebKeySet(ctx context.Context, jwks []byte) error
	AuthAttempt(ctx context.Context, hash string) (*skillz.AuthAttempt, error)
	CreateAuthAttempt(ctx context.Context, attempt *skillz.AuthAttempt) error
}

const keyAuthAttempt = "auth::attempt"
const keyAuthJWKS = "auth::jwks"

func (s *Service) JSONWebKeySet(ctx context.Context) ([]byte, error) {
	key := generateKey(keyAuthJWKS)

	result, err := s.redis.Get(ctx, key).Bytes()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, errors.Wrapf(err, errorFFormat, authAPI, "JSONWebKeySet", "failed to fetch results from cache")
	}

	if errors.Is(err, redis.Nil) {
		return nil, nil
	}

	return result, nil

}

func (s *Service) SaveJSONWebKeySet(ctx context.Context, jwks []byte) error {
	key := generateKey(keyAuthJWKS)
	err := s.redis.Set(ctx, key, jwks, time.Hour*6).Err()
	return errors.Wrapf(err, errorFFormat, authAPI, "SaveJSONWebKeySet", "failed to write cache")
}

func (s *Service) AuthAttempt(ctx context.Context, state string) (*skillz.AuthAttempt, error) {

	key := generateKey(keyAuthAttempt, state)
	result, err := s.redis.Get(ctx, key).Bytes()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, errors.Wrapf(err, errorFFormat, authAPI, "AuthAttempt", "failed to fetch results from cache")
	}

	if errors.Is(err, redis.Nil) {
		return nil, nil
	}

	var attempt = new(skillz.AuthAttempt)
	err = json.Unmarshal(result, attempt)
	return attempt, errors.Wrapf(err, errorFFormat, authAPI, "AuthAttempt", "failed to decode json to structure")

}

func (s *Service) CreateAuthAttempt(ctx context.Context, attempt *skillz.AuthAttempt) error {

	b, err := json.Marshal(attempt)
	if err != nil {
		return errors.Wrapf(err, errorFFormat, authAPI, "CreateAuthAttempt", "failed to encode struct as json")
	}

	key := generateKey(keyAuthAttempt, attempt.State)
	err = s.redis.Set(ctx, key, b, time.Minute*5).Err()
	return errors.Wrapf(err, errorFFormat, authAPI, "CreateAuthAttempt", "failed to write cache")

}
