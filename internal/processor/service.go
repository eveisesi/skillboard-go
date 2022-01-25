package processor

import (
	"context"
	"time"

	"github.com/eveisesi/skillz"
	"github.com/eveisesi/skillz/internal"
	"github.com/eveisesi/skillz/internal/cache"
	"github.com/eveisesi/skillz/internal/user/v2"
	"github.com/go-redis/redis/v8"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Service struct {
	logger *logrus.Logger
	redis  *redis.Client
	cache  *cache.Service
	user   user.API

	scopes skillz.ScopeProcessors
}

func New(logger *logrus.Logger, redisClient *redis.Client, user user.API, scopes skillz.ScopeProcessors) *Service {
	return &Service{
		logger: logger,
		redis:  redisClient,
		cache:  cache.New(redisClient),
		user:   user,
		scopes: scopes,
	}
}

func (s *Service) Run() error {

	s.logger.Info("Processor has started....")

	for {
		var entry = logrus.NewEntry(s.logger)
		var ctx context.Context = context.Background()
		result, err := s.redis.BZPopMin(ctx, 0, internal.UpdateQueue).Result()
		if err != nil && !errors.Is(err, redis.Nil) {
			return err
		}

		userIDStr, ok := result.Member.(string)
		if !ok {
			s.logger.Errorf("unexpected value for member, expected string, got %T", result.Member)
			continue
		}

		entry = entry.WithField("userIDStr", userIDStr)

		userID, err := uuid.FromString(userIDStr)
		if err != nil {
			entry.WithError(err).Errorf("invalid uuid found on update queue")
			continue
		}

		_, err = s.redis.ZAdd(ctx, internal.UpdatingQueue, &result.Z).Result()
		if err != nil {
			return err
		}

		err = s.processUser(context.Background(), userID)
		if err != nil {
			entry.WithError(err).Error("failed to process user id")
			continue
		}

		_, err = s.redis.ZRem(ctx, internal.UpdatingQueue, result.Member).Result()
		if err != nil {
			return err
		}
	}

}

func (s *Service) processUser(ctx context.Context, userID uuid.UUID) error {

	user, err := s.user.User(ctx, userID)
	if err != nil {
		return errors.Wrap(err, "failed to fetch user from data store")
	}

	err = s.user.ValidateCurrentToken(ctx, user)
	if err != nil {
		return err
	}

	err = s.cache.ResetUserCache(ctx, user)
	if err != nil {
		return err
	}

	user.IsProcessing = true
	err = s.user.UpdateUser(ctx, user)
	if err != nil {
		s.logger.WithError(err).Error("failed to update user")
	}

	defer func() {
		user.IsNew = false
		user.LastProcessed.SetValid(time.Now())
		user.IsProcessing = false
		err = s.user.UpdateUser(ctx, user)
		if err != nil {
			s.logger.WithError(err).Error("failed to update user")
		}
	}()

	for _, processor := range s.scopes {
		err = processor.Process(ctx, user)
		if err != nil {
			err = errors.Wrap(err, "processor failed to process user")
			user.Disabled = true
			user.DisabledReason.SetValid(err.Error())
			user.DisabledTimestamp.SetValid(time.Now())
			return err
		}
	}

	return nil
}
