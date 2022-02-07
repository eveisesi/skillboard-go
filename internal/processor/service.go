package processor

import (
	"context"
	"time"

	"github.com/eveisesi/skillz"
	"github.com/eveisesi/skillz/internal"
	"github.com/eveisesi/skillz/internal/user/v2"
	"github.com/go-redis/redis/v8"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Service struct {
	logger   *logrus.Logger
	redis    *redis.Client
	newrelic *newrelic.Application

	user user.API

	processors skillz.ScopeProcessors
}

func New(logger *logrus.Logger, redisClient *redis.Client, newrelic *newrelic.Application, user user.API, processors skillz.ScopeProcessors) *Service {
	return &Service{
		logger:   logger,
		redis:    redisClient,
		newrelic: newrelic,

		user: user,

		processors: processors,
	}
}

func (s *Service) downtime() bool {
	var now = time.Now()
	var startDT = time.Date(now.Year(), now.Month(), now.Day(), 10, 55, 0, 0, time.UTC)
	var endDT = time.Date(now.Year(), now.Month(), now.Day(), 11, 30, 0, 0, time.UTC)

	return now.After(startDT) && now.Before(endDT)
}

func (s *Service) Run() error {

	s.logger.Info("Processor has started....")

	for {

		if s.downtime() {
			s.logger.Info("sleeping for downtime")
			time.Sleep(time.Minute)
			continue
		}

		var entry = logrus.NewEntry(s.logger)
		var ctx context.Context = context.Background()
		result, err := s.redis.BZPopMin(ctx, 0, internal.UpdateQueue).Result()
		if err != nil && !errors.Is(err, redis.Nil) {
			return err
		}

		userID, ok := result.Member.(string)
		if !ok {
			s.logger.Errorf("unexpected value for member, expected string, got %T", result.Member)
			continue
		}

		entry = entry.WithField("userID", userID)

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

func (s *Service) processUser(ctx context.Context, userID string) error {

	txn := s.newrelic.StartTransaction("ProcessUser")
	txn.AddAttribute("userID", userID)
	defer txn.End()

	ctx = newrelic.NewContext(ctx, txn)

	user, err := s.user.User(ctx, userID)
	if err != nil {
		err = errors.Wrap(err, "failed to fetch user from data store")
		txn.NoticeError(err)
		return err
	}

	defer func() {
		user.IsNew = false
		user.LastProcessed.SetValid(time.Now())
		user.IsProcessing = false
		err = s.user.CreateUser(ctx, user)
		if err != nil {
			txn.NoticeError(errors.Wrap(err, "failed to update user"))
			s.logger.WithError(err).Error("failed to update user")
		}
	}()

	err = s.user.ValidateCurrentToken(ctx, user)
	if err != nil {
		err = errors.Wrap(err, "failed to validate token")
		txn.NoticeError(err)
		user.Disabled = true
		user.DisabledReason.SetValid(err.Error())
		user.DisabledTimestamp.SetValid(time.Now())
		return err
	}

	err = s.user.ResetUserCache(ctx, user)
	if err != nil {
		err = errors.Wrap(err, "failed to invalidate user cache")
		txn.NoticeError(err)
		return err
	}

	user.IsProcessing = true
	err = s.user.CreateUser(ctx, user)
	if err != nil {
		err = errors.Wrap(err, "failed to update user")
		txn.NoticeError(err)
		s.logger.WithError(err).Error("failed to update user")
		return err
	}

	for _, processor := range s.processors {
		err = processor.Process(ctx, user)
		if err != nil {
			err = errors.Wrap(err, "processor failed to process user")
			txn.NoticeError(err)
			user.Disabled = true
			user.DisabledReason.SetValid(err.Error())
			user.DisabledTimestamp.SetValid(time.Now())
			return err
		}
	}

	return nil
}
