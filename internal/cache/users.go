package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/eveisesi/skillz"
	"github.com/go-redis/redis/v8"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
)

type UserAPI interface {
	UserSettings(ctx context.Context, id uuid.UUID) (*skillz.UserSettings, error)
	SetUserSettings(ctx context.Context, id uuid.UUID, settings *skillz.UserSettings, expires time.Duration) error

	SearchUsers(ctx context.Context, q string) ([]*skillz.UserSearchResult, error)
	SetSearchUsersResults(ctx context.Context, q string, users []*skillz.UserSearchResult, expires time.Duration) error
	NewUsersBySP(ctx context.Context) ([]*skillz.UserWithSkillMeta, error)
	SetNewUsersBySP(ctx context.Context, users []*skillz.UserWithSkillMeta, expires time.Duration) error
	BustNewUsersBySP(ctx context.Context) error
}

const (
	userSearchKeyPrefix   = "user::search"
	userSettingsKeyPrefix = "user::settings"
	usersNewBySPPrefix    = "users::new-by-sp"
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

func (s *Service) NewUsersBySP(ctx context.Context) ([]*skillz.UserWithSkillMeta, error) {

	var users = make([]*skillz.UserWithSkillMeta, 0)

	key := generateKey(usersNewBySPPrefix)
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

func (s *Service) SetNewUsersBySP(ctx context.Context, users []*skillz.UserWithSkillMeta, expires time.Duration) error {

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
	key := generateKey(usersNewBySPPrefix)
	return s.redis.Del(ctx, key).Err()
}

func (s *Service) UserSettings(ctx context.Context, id uuid.UUID) (*skillz.UserSettings, error) {

	key := generateKey(userSettingsKeyPrefix, id.String())
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

func (s *Service) SetUserSettings(ctx context.Context, id uuid.UUID, settings *skillz.UserSettings, expires time.Duration) error {

	key := generateKey(userSettingsKeyPrefix, id.String())
	data, err := json.Marshal(settings)
	if err != nil {
		return errors.Wrapf(err, errorFFormat, userAPI, "SetType", "failed to encode structure as json")
	}

	err = s.redis.Set(ctx, key, data, time.Hour).Err()
	return errors.Wrapf(err, errorFFormat, userAPI, "SetType", "failed to write cache")

}
