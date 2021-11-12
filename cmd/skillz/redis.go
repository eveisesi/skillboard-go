package main

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

func buildRedis() {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	redisClient = redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Host,
		Password: cfg.Redis.Pass,
	})

	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		logger.WithError(err).Fatal("failed to ping redis")
	}

}
