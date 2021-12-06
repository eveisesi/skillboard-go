package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/eveisesi/skillz"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

type UserAPI interface {
	SearchUsers(ctx context.Context, q string) ([]*skillz.User, error)
	SetSearchUsersResults(ctx context.Context, q string, users []*skillz.User, expires time.Duration) error
}

const (
	userSearchKeyPrefix = "user::search"
)

func (s *Service) SearchUsers(ctx context.Context, q string) ([]*skillz.User, error) {

	key := generateKey(userSearchKeyPrefix, hash(q))
	results, err := s.redis.SMembers(ctx, key).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, errors.Wrapf(err, errorFFormat, userAPI, "SearchUsers", "failed to fetch results from cache")
	}

	if errors.Is(err, redis.Nil) || len(results) == 0 {
		return nil, nil
	}

	users := make([]*skillz.User, 0, len(results))
	for _, result := range results {
		var user = new(skillz.User)
		err = json.Unmarshal([]byte(result), user)
		if err != nil {
			return users, errors.Wrapf(err, errorFFormat, userAPI, "SearchUsers", "failed to decode json to structure")
		}

		users = append(users, user)
	}

	return users, nil

}

func (s *Service) SetSearchUsersResults(ctx context.Context, q string, users []*skillz.User, expires time.Duration) error {

	members := make([]interface{}, 0, len(users))
	for _, skill := range users {
		data, err := json.Marshal(skill)
		if err != nil {
			return errors.Wrapf(err, errorFFormat, userAPI, "SetSearchUsersResults", "failed to encode struct as json")
		}

		members = append(members, string(data))
	}

	if len(members) == 0 {
		return nil
	}

	hashedQ := hash(q)
	key := generateKey(characterFlyableKeyPrefix, hashedQ)
	err := s.redis.SAdd(ctx, key, members...).Err()
	if err != nil {
		return errors.Wrapf(err, errorFFormat, skillAPI, "SetSearchUsersResults", "failed to write cache")
	}

	err = s.redis.Expire(ctx, key, expires).Err()

	return errors.Wrapf(err, errorFFormat, skillAPI, "SetSearchUsersResults", "failed to set expiry on set")

}
