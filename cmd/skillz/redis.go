package main

import (
	"context"

	"github.com/go-redis/redis/v8"
)

func buildRedis(ctx context.Context) *redis.Client {
	redis := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Host,
		Password: cfg.Redis.Pass,
	})

	_, err := redis.Ping(ctx).Result()
	if err != nil {
		logger.WithError(err).Fatal("failed to ping redis")
	}
	return redis
}
