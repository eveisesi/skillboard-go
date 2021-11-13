package processor

import (
	"context"

	"github.com/eveisesi/skillz"
	"github.com/eveisesi/skillz/internal/user"
	"github.com/go-redis/redis/v8"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Service struct {
	logger *logrus.Logger
	redis  *redis.Client
	user   user.API

	scopes skillz.ScopeProcessors
}

func New(logger *logrus.Logger, redisClient *redis.Client, user user.API, scopes skillz.ScopeProcessors) *Service {
	return &Service{
		logger: logger,
		redis:  redisClient,
		user:   user,
		scopes: scopes,
	}
}

const (
	UpdateQueue   = "skillz::queue::update"
	UpdatingQueue = "skillz::queue::updating"
)

func (s *Service) Run() error {

	s.logger.Info("Processor has started....")

	// for {

	// 	value, err := s.redis.BLMove(ctx, UpdateQueue, UpdatingQueue, "left", "right", time.Second*10).Result()
	// 	// values, err := s.redis.BRPop(ctx, time.Second*10, UpdateQueue).Result()
	// 	if err != nil && !errors.Is(err, redis.Nil) {
	// 		return errors.Wrap(err, "failed to BRPop from List")
	// 	}
	// 	if value == "" {
	// 		continue
	// 	}

	// 	fmt.Println(value)
	// 	fmt.Println("------")
	// 	time.Sleep(time.Second * 30)

	// 	_, err = s.redis.LRem(ctx, UpdatingQueue, 0, value).Result()
	// 	if err != nil {
	// 		return errors.Wrap(err, "failed to BRPop from List")
	// 	}
	// }

	for {
		var entry = logrus.NewEntry(s.logger)
		var ctx context.Context = context.Background()
		result, err := s.redis.BZPopMin(ctx, 0, UpdateQueue).Result()
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

		_, err = s.redis.ZAdd(ctx, UpdatingQueue, &result.Z).Result()
		if err != nil {
			return err
		}

		err = s.processUser(ctx, userID)
		if err != nil {
			entry.WithError(err).Error("failed to process user id")
			continue
		}

		_, err = s.redis.ZRem(ctx, UpdatingQueue, result.Member).Result()
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

	err = s.user.ValidateToken(ctx, user)
	if err != nil {
		return err
	}

processorLoop:
	for _, processor := range s.scopes {
		scopes := processor.Scopes()
		for _, scope := range user.Scopes {
			if scopeInSlcScopes(scope, scopes) {
				err = processor.Process(ctx, user)
				if err != nil {
					return errors.Wrap(err, "processor failed to process user")
				}

				continue processorLoop
			}
		}
	}

	return nil
}

func scopeInSlcScopes(s skillz.Scope, slc []skillz.Scope) bool {
	for _, scope := range slc {
		if s == scope {
			return true
		}
	}

	return false
}