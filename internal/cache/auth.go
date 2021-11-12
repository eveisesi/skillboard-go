package cache

import (
	"context"
	"encoding/json"
	"fmt"
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
		return nil, err
	}

	if len(result) == 0 {
		return nil, nil
	}

	return result, nil

}

func (s *Service) SaveJSONWebKeySet(ctx context.Context, jwks []byte) error {
	key := generateKey(keyAuthJWKS)
	_, err := s.redis.Set(ctx, key, jwks, time.Hour*6).Result()
	return err
}

func (s *Service) AuthAttempt(ctx context.Context, state string) (*skillz.AuthAttempt, error) {

	key := generateKey(keyAuthAttempt, state)

	result, err := s.redis.Get(ctx, key).Bytes()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}

	if len(result) == 0 {
		return nil, nil
	}

	var attempt = new(skillz.AuthAttempt)
	err = json.Unmarshal(result, attempt)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal data onto result struct")
	}

	return attempt, nil

}

func (s *Service) CreateAuthAttempt(ctx context.Context, attempt *skillz.AuthAttempt) error {

	if attempt.State == "" {
		return fmt.Errorf("empty state provided")
	}

	b, err := json.Marshal(attempt)
	if err != nil {
		return fmt.Errorf("failed to cache auth attempt: %w", err)
	}

	key := generateKey(keyAuthAttempt, attempt.State)

	_, err = s.redis.Set(ctx, key, b, time.Minute*5).Result()
	if err != nil {
		return errors.Wrap(err, "failed to create auth attempt")
	}

	return nil
}
