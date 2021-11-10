package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/eveisesi/skillz"
	"github.com/go-redis/redis"
)

type AuthAPI interface {
	JSONWebKeySet(ctx context.Context) ([]byte, error)
	SaveJSONWebKeySet(ctx context.Context, jwks []byte) error
	AuthAttempt(ctx context.Context, hash string) (*skillz.AuthAttempt, error)
	CreateAuthAttempt(ctx context.Context, attempt *skillz.AuthAttempt) (*skillz.AuthAttempt, error)
}

const keyAuthAttempt = "skillz::auth::attempt::%s"
const keyAuthJWKS = "skillz::auth::jwks"

func (s *Service) JSONWebKeySet(ctx context.Context) ([]byte, error) {

	result, err := s.redis.Get(ctx, keyAuthJWKS).Bytes()
	if err != nil && err != redis.Nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, nil
	}

	return result, nil

}

func (s *Service) SaveJSONWebKeySet(ctx context.Context, jwks []byte) error {

	_, err := s.redis.Set(ctx, keyAuthJWKS, jwks, time.Hour*6).Result()

	return err

}

func (s *Service) AuthAttempt(ctx context.Context, hash string) (*skillz.AuthAttempt, error) {

	var attempt = new(skillz.AuthAttempt)

	result, err := s.redis.Get(ctx, fmt.Sprintf(keyAuthAttempt, hash)).Bytes()
	if err != nil && err != redis.Nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, nil
	}

	err = json.Unmarshal(result, attempt)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal data onto result struct: %w", err)
	}

	return attempt, nil

}

func (s *Service) CreateAuthAttempt(ctx context.Context, attempt *skillz.AuthAttempt) (*skillz.AuthAttempt, error) {

	if attempt.State == "" {
		return nil, fmt.Errorf("empty state provided")
	}

	b, err := json.Marshal(attempt)
	if err != nil {
		return nil, fmt.Errorf("failed to cache auth attempt: %w", err)
	}

	_, err = s.redis.Set(ctx, fmt.Sprintf(keyAuthAttempt, attempt.State), b, time.Minute*5).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to create auth attempt: %w", err)
	}

	return attempt, nil
}
