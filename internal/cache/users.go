package cache

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/eveisesi/skillz"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

type UserAPI interface {
	SearchUsers(ctx context.Context, q string) ([]*skillz.UserSearchResult, error)
	SetSearchUsersResults(ctx context.Context, q string, users []*skillz.UserSearchResult, expires time.Duration) error
	NewUsersBySP(ctx context.Context, days int) ([]*skillz.UserWithSkillMeta, error)
	SetNewUsersBySP(ctx context.Context, days int, users []*skillz.UserWithSkillMeta, expires time.Duration) error
}

const (
	userSearchKeyPrefix = "user::search"
	usersNewBySPPrefix  = "users::new-by-sp"
)

func (s *Service) SearchUsers(ctx context.Context, q string) ([]*skillz.UserSearchResult, error) {

	var users = make([]*skillz.UserSearchResult, 0)

	key := generateKey(userSearchKeyPrefix, hash(q))
	results, err := s.redis.SMembers(ctx, key).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return users, errors.Wrapf(err, errorFFormat, userAPI, "SearchUsers", "failed to fetch results from cache")
	}

	if errors.Is(err, redis.Nil) || len(results) == 0 {
		return users, nil
	}

	users = make([]*skillz.UserSearchResult, 0, len(results))
	for _, result := range results {
		var user = new(skillz.UserSearchResult)
		err = json.Unmarshal([]byte(result), user)
		if err != nil {
			return users, errors.Wrapf(err, errorFFormat, userAPI, "SearchUsers", "failed to decode json to structure")
		}

		users = append(users, user)
	}

	return users, nil

}

func (s *Service) SetSearchUsersResults(ctx context.Context, q string, users []*skillz.UserSearchResult, expires time.Duration) error {

	members := make([]interface{}, 0, len(users))
	for _, user := range users {
		data, err := json.Marshal(user)
		if err != nil {
			return errors.Wrapf(err, errorFFormat, userAPI, "SetSearchUsersResults", "failed to encode struct as json")
		}

		members = append(members, string(data))
	}

	if len(members) == 0 {
		return nil
	}

	key := generateKey(userSearchKeyPrefix, hash(q))
	err := s.redis.SAdd(ctx, key, members...).Err()
	if err != nil {
		return errors.Wrapf(err, errorFFormat, userAPI, "SetSearchUsersResults", "failed to write cache")
	}

	err = s.redis.Expire(ctx, key, expires).Err()

	return errors.Wrapf(err, errorFFormat, userAPI, "SetSearchUsersResults", "failed to set expiry on set")

}

func (s *Service) NewUsersBySP(ctx context.Context, days int) ([]*skillz.UserWithSkillMeta, error) {

	var users = make([]*skillz.UserWithSkillMeta, 0)

	key := generateKey(usersNewBySPPrefix, strconv.Itoa(days))
	results, err := s.redis.SMembers(ctx, key).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return users, errors.Wrapf(err, errorFFormat, userAPI, "NewUsersBySP", "failed to fetch results from cache")
	}

	if errors.Is(err, redis.Nil) || len(results) == 0 {
		return users, nil
	}

	users = make([]*skillz.UserWithSkillMeta, 0, len(results))
	for _, result := range results {
		var user = new(skillz.UserWithSkillMeta)
		err = json.Unmarshal([]byte(result), user)
		if err != nil {
			return users, errors.Wrapf(err, errorFFormat, userAPI, "NewUsersBySP", "failed to decode json to structure")
		}

		users = append(users, user)
	}

	return users, nil

}

func (s *Service) SetNewUsersBySP(ctx context.Context, days int, users []*skillz.UserWithSkillMeta, expires time.Duration) error {

	members := make([]interface{}, 0, len(users))
	for _, user := range users {
		data, err := json.Marshal(user)
		if err != nil {
			return errors.Wrapf(err, errorFFormat, userAPI, "SetNewUsersBySP", "failed to encode struct as json")
		}

		members = append(members, string(data))
	}

	if len(members) == 0 {
		return nil
	}

	key := generateKey(usersNewBySPPrefix, strconv.Itoa(days))
	err := s.redis.SAdd(ctx, key, members...).Err()
	if err != nil {
		return errors.Wrapf(err, errorFFormat, userAPI, "SetNewUsersBySP", "failed to write cache")
	}

	err = s.redis.Expire(ctx, key, expires).Err()

	return errors.Wrapf(err, errorFFormat, userAPI, "SetNewUsersBySP", "failed to set expiry on set")

}
