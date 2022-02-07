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
	UserSettings(ctx context.Context, id string) (*skillz.UserSettings, error)
	SetUserSettings(ctx context.Context, id string, settings *skillz.UserSettings, expires time.Duration) error

	SearchUsers(ctx context.Context, q string) ([]*skillz.UserSearchResult, error)
	SetSearchUsersResults(ctx context.Context, q string, users []*skillz.UserSearchResult, expires time.Duration) error
	NewUsersBySP(ctx context.Context) ([]*skillz.User, error)
	SetNewUsersBySP(ctx context.Context, users []*skillz.User, expires time.Duration) error

	BustNewUsersBySP(ctx context.Context) error
	ResetUserCache(ctx context.Context, user *skillz.User) error
}

const (
	userSearchKeyPrefix   = "user::search"
	userSettingsKeyPrefix = "user::settings"
	usersNewBySPPrefix    = "users::new-by-sp"
	recentUsersPrefix     = "users::recent"
)

func (s *Service) RecentUsers(ctx context.Context) (*skillz.RecentUsers, error) {
	if s.disabled {
		return nil, nil
	}
	key := generateKey(recentUsersPrefix)
	result, err := s.redis.Get(ctx, key).Bytes()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, errors.Wrapf(err, errorFFormat, userAPI, "SearchUsers", "failed to fetch results from cache")
	}

	if errors.Is(err, redis.Nil) {
		return nil, nil
	}

	var out = new(skillz.RecentUsers)
	err = json.Unmarshal(result, out)
	return out, errors.Wrapf(err, errorFFormat, userAPI, "UserSettings", "failed to decode json to structure")

}

func (s *Service) ResetUserCache(ctx context.Context, user *skillz.User) error {

	keys := []string{
		generateKey(characterSkillMetaKeyPrefix, strconv.FormatUint(user.CharacterID, 10)),
		generateKey(characterSkillsKeyPrefix, strconv.FormatUint(user.CharacterID, 10)),
		generateKey(characterSkillsGroupedKeyPrefix, strconv.FormatUint(user.CharacterID, 10)),
		generateKey(characterFlyableKeyPrefix, strconv.FormatUint(user.CharacterID, 10)),
		generateKey(characterSkillQueueKeySummaryPrefix, strconv.FormatUint(user.CharacterID, 10)),
		generateKey(characterAttributesKeyPrefix, strconv.FormatUint(user.CharacterID, 10)),
	}

	for _, key := range keys {
		err := s.redis.Del(ctx, key).Err()
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) SearchUsers(ctx context.Context, q string) ([]*skillz.UserSearchResult, error) {
	if s.disabled {
		return nil, nil
	}
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
	if s.disabled {
		return nil
	}
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

func (s *Service) NewUsersBySP(ctx context.Context) ([]*skillz.User, error) {
	if s.disabled {
		return nil, nil
	}
	var users = make([]*skillz.User, 0)

	key := generateKey(usersNewBySPPrefix)
	results, err := s.redis.SMembers(ctx, key).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return users, errors.Wrapf(err, errorFFormat, userAPI, "NewUsersBySP", "failed to fetch results from cache")
	}

	if errors.Is(err, redis.Nil) || len(results) == 0 {
		return users, nil
	}

	users = make([]*skillz.User, 0, len(results))
	for _, result := range results {
		var user = new(skillz.User)
		err = json.Unmarshal([]byte(result), user)
		if err != nil {
			return users, errors.Wrapf(err, errorFFormat, userAPI, "NewUsersBySP", "failed to decode json to structure")
		}

		users = append(users, user)
	}

	return users, nil

}

func (s *Service) SetNewUsersBySP(ctx context.Context, users []*skillz.User, expires time.Duration) error {
	if s.disabled {
		return nil
	}
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

	key := generateKey(usersNewBySPPrefix)
	err := s.redis.SAdd(ctx, key, members...).Err()
	if err != nil {
		return errors.Wrapf(err, errorFFormat, userAPI, "SetNewUsersBySP", "failed to write cache")
	}

	err = s.redis.Expire(ctx, key, expires).Err()

	return errors.Wrapf(err, errorFFormat, userAPI, "SetNewUsersBySP", "failed to set expiry on set")

}

func (s *Service) BustNewUsersBySP(ctx context.Context) error {
	return s.redis.Del(ctx, generateKey(usersNewBySPPrefix)).Err()
}

func (s *Service) UserSettings(ctx context.Context, id string) (*skillz.UserSettings, error) {
	if s.disabled {
		return nil, nil
	}
	key := generateKey(userSettingsKeyPrefix, id)
	result, err := s.redis.Get(ctx, key).Bytes()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, errors.Wrapf(err, errorFFormat, userAPI, "UserSettings", "failed to fetch results from cache")
	}

	if errors.Is(err, redis.Nil) {
		return nil, nil
	}

	var settings = new(skillz.UserSettings)
	err = json.Unmarshal(result, settings)
	return settings, errors.Wrapf(err, errorFFormat, userAPI, "UserSettings", "failed to decode json to structure")

}

func (s *Service) SetUserSettings(ctx context.Context, id string, settings *skillz.UserSettings, expires time.Duration) error {
	if s.disabled {
		return nil
	}
	key := generateKey(userSettingsKeyPrefix, id)
	data, err := json.Marshal(settings)
	if err != nil {
		return errors.Wrapf(err, errorFFormat, userAPI, "SetType", "failed to encode structure as json")
	}

	err = s.redis.Set(ctx, key, data, time.Hour).Err()
	return errors.Wrapf(err, errorFFormat, userAPI, "SetType", "failed to write cache")

}
