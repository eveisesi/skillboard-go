package user

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/eveisesi/skillz"
	"github.com/eveisesi/skillz/internal"
	"github.com/eveisesi/skillz/internal/alliance"
	"github.com/eveisesi/skillz/internal/auth"
	"github.com/eveisesi/skillz/internal/cache"
	"github.com/eveisesi/skillz/internal/character"
	"github.com/eveisesi/skillz/internal/corporation"
	"github.com/eveisesi/skillz/internal/skill"
	"github.com/go-redis/redis/v8"
	"github.com/gofrs/uuid"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type API interface {
	Login(ctx context.Context, code, state string) (*skillz.User, error)
	UserFromToken(ctx context.Context, token jwt.Token) (*skillz.User, error)
	ValidateCurrentToken(ctx context.Context, user *skillz.User) error

	User(ctx context.Context, id uuid.UUID) (*skillz.User, error)
	RefreshUser(ctx context.Context, user *skillz.User) error
	UserByCharacterID(ctx context.Context, characterID uint64) (*skillz.User, error)
	SearchUsers(ctx context.Context, q string) ([]*skillz.UserSearchResult, error)
	UpdateUser(ctx context.Context, user *skillz.User) error

	NewUsersBySP(ctx context.Context) ([]*skillz.UserWithSkillMeta, error)

	ProcessUpdatableUsers(ctx context.Context) error

	UserSettings(ctx context.Context, id uuid.UUID) (*skillz.UserSettings, error)
	CreateUserSettings(ctx context.Context, userID uuid.UUID, settings *skillz.UserSettings) error
}

type Service struct {
	redis  *redis.Client
	logger *logrus.Logger
	cache  cache.UserAPI

	auth        auth.API
	character   character.API
	corporation corporation.API
	alliance    alliance.API

	skills skill.API

	skillz.UserRepository
}

var _ API = new(Service)

func New(
	redis *redis.Client,
	logger *logrus.Logger,
	cache cache.UserAPI,
	auth auth.API,
	alliance alliance.API,
	character character.API,
	corporation corporation.API,
	skills skill.API,
	user skillz.UserRepository,
) *Service {
	return &Service{
		redis:          redis,
		logger:         logger,
		cache:          cache,
		alliance:       alliance,
		auth:           auth,
		character:      character,
		corporation:    corporation,
		skills:         skills,
		UserRepository: user,
	}
}

func (s *Service) SearchUsers(ctx context.Context, q string) ([]*skillz.UserSearchResult, error) {

	results, err := s.cache.SearchUsers(ctx, q)
	if err != nil {
		return nil, err
	}

	if len(results) > 0 {
		return results, nil
	}

	users, err := s.UserRepository.SearchUsers(ctx, q)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	results = make([]*skillz.UserSearchResult, 0, len(users))
	for _, user := range users {
		info, err := s.character.Character(ctx, user.CharacterID)
		if err != nil {
			s.logger.WithError(err).Error("failed to look up user character")
			continue
		}

		results = append(results, &skillz.UserSearchResult{
			User: user,
			Info: info,
		})
	}

	defer func() {
		err := s.cache.SetSearchUsersResults(ctx, q, results, time.Minute*10)
		if err != nil {
			s.logger.WithError(err).Error("failed to cache search results")
		}
	}()

	return results, errors.Wrap(err, "failed to cache user search results")

}

func (s *Service) Login(ctx context.Context, code, state string) (*skillz.User, error) {

	attempt, err := s.auth.AuthAttempt(ctx, state)
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch attempt with provided state")
	}

	if attempt == nil {
		return nil, errors.New("invalid attempt")
	}

	if attempt != nil && attempt.Status == skillz.InvalidAuthStatus {
		return nil, errors.New("request is no longer valid, please try again")
	}

	bearer, err := s.auth.BearerForESICode(ctx, code)
	if err != nil {
		return nil, errors.Wrap(err, "failed to exchange authorization code for token")
	}

	token, err := s.auth.ParseAndVerifyESIToken(ctx, bearer.AccessToken)
	if err != nil {
		return nil, errors.Wrap(err, "failed to verify token")
	}

	user, err := s.UserFromToken(ctx, token)
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch/provision user for the provided token")
	}

	claims := token.PrivateClaims()
	if _, ok := claims["owner"]; !ok {
		return nil, errors.New("invalid token, owner claims is missing")
	}

	ownerHash := claims["owner"].(string)

	if user.OwnerHash == "" {
		user.OwnerHash = ownerHash
	} else if user.OwnerHash != ownerHash {
		return nil, errors.New("character owner hash mismatch. please contact support")
	}

	scp := make([]skillz.Scope, 0)
	switch a := claims["scp"].(type) {
	case []interface{}:
		for _, v := range a {
			scp = append(scp, skillz.Scope(v.(string)))
		}
	case string:
		scp = append(scp, skillz.Scope(a))
	default:
		return nil, errors.New("invalid type for scp claim in token.")
	}

	user.Scopes = scp
	user.AccessToken = bearer.AccessToken
	user.RefreshToken = bearer.RefreshToken
	user.Expires = bearer.Expiry
	user.LastLogin = time.Now()

	switch user.IsNew {
	case true:
		err = s.UserRepository.CreateUser(ctx, user)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create user in data store")
		}
	case false:
		err = s.UserRepository.UpdateUser(ctx, user)
		if err != nil {
			return nil, errors.Wrap(err, "failed to update user in data store")
		}
	}

	sessionID := internal.SessionIDFromContext(ctx)
	if sessionID != uuid.Nil {
		internal.CacheSet(sessionID, user.ID)
	}

	err = s.auth.DeleteAuthAttempt(ctx, attempt)
	if err != nil {
		return nil, errors.Wrap(err, "failed to remove auth attempt from cache")
	}

	err = s.redis.ZAdd(ctx, internal.UpdateQueue, &redis.Z{Score: float64(time.Now().Unix()), Member: user.ID.String()}).Err()
	if err != nil {
		return nil, errors.Wrap(err, "failed to push user id to processing queue")
	}

	if user.IsNew {

		defer func(ctx context.Context) {
			err = s.cache.BustNewUsersBySP(ctx)
			if err != nil {
				s.logger.WithError(err).Error("failed to bust new users by sp cache")
			}
		}(context.Background())
	}

	user.Settings, err = s.UserSettings(ctx, user.ID)
	// if err != n
	return user, err

}

func (s *Service) RefreshUser(ctx context.Context, user *skillz.User) error {

	err := s.redis.ZAdd(ctx, internal.UpdateQueue, &redis.Z{Score: float64(time.Now().Unix()), Member: user.ID.String()}).Err()
	if err != nil {
		return errors.Wrap(err, "failed to push user id to processing queue")
	}

	return nil

}

func (s *Service) NewUsersBySP(ctx context.Context) ([]*skillz.UserWithSkillMeta, error) {

	users, err := s.cache.NewUsersBySP(ctx)
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}

	if len(users) > 0 {
		return users, nil
	}

	userRecords, err := s.UserRepository.NewUsersBySP(ctx)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, errors.Wrap(err, "failed to fetch users by sp in last seven days")
	}

	for _, record := range userRecords {

		meta, err := s.skills.Meta(ctx, record.CharacterID)
		if err != nil {
			continue
		}

		skills, err := s.skills.Skillz(ctx, record.CharacterID)
		if err != nil {
			continue
		}

		queue, err := s.skills.SkillQueue(ctx, record.CharacterID)
		if err != nil {
			continue
		}

		info, err := s.character.Character(ctx, record.CharacterID)
		if err != nil {
			continue
		}

		users = append(users, &skillz.UserWithSkillMeta{
			User:   record,
			Meta:   meta,
			Skills: skills,
			Queue:  queue,
			Info:   info,
		})

	}

	defer func() {
		err = s.cache.SetNewUsersBySP(ctx, users, time.Minute*5)
		if err != nil {
			s.logger.WithError(err).Error("failed to cache users by sp in last seven days")
		}
	}()

	return users, nil

}

func (s *Service) UserFromToken(ctx context.Context, token jwt.Token) (*skillz.User, error) {

	var user *skillz.User
	var err error

	switch token.Issuer() {
	case "login.eveonline.com":
		subject := token.Subject()
		if subject == "" {
			return nil, errors.New("token subject is empty, expected parsable value")
		}

		characterID, ierr := memberIDFromSubject(subject)
		if ierr != nil {
			return nil, errors.Wrap(err, "failed to parse character id from token subject")
		}

		user, err = s.UserRepository.UserByCharacterID(ctx, characterID)
		if err != nil && errors.Is(err, sql.ErrNoRows) {
			user = &skillz.User{
				ID:          uuid.Must(uuid.NewV4()),
				CharacterID: characterID,
				IsNew:       true,
			}
			err = nil
		}

	default:
		return nil, errors.Errorf("unsupported token issuer")
	}

	if err != nil {
		return nil, errors.Wrap(err, "failed to locate user for this token in our database")
	}

	character, err := s.character.Character(ctx, user.CharacterID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch user's character")
	}

	corporation, err := s.corporation.Corporation(ctx, character.CorporationID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch user's corporation")
	}

	if corporation.AllianceID.Valid {
		_, err := s.alliance.Alliance(ctx, corporation.AllianceID.Uint)
		if err != nil {
			return nil, errors.Wrap(err, "failed to fetch corporation's alliance")
		}
	}

	return user, nil

}

func (s *Service) ValidateCurrentToken(ctx context.Context, user *skillz.User) error {

	updated, token, err := s.auth.ValidateESITokenForUser(ctx, user)
	if err != nil {
		return err
	}
	if !updated {
		return nil
	}

	user.ApplyToken(token)

	return s.UserRepository.UpdateUser(ctx, user)

}

func memberIDFromSubject(sub string) (uint64, error) {

	parts := strings.Split(sub, ":")

	if len(parts) != 3 {
		return 0, fmt.Errorf("invalid sub format")
	}

	return strconv.ParseUint(parts[2], 10, 64)

}

func (s *Service) ProcessUpdatableUsers(ctx context.Context) error {

	users, err := s.UserRepository.UsersSortedByProcessedAtLimit(ctx, 100)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return errors.Wrap(err, "failed to query users table for processable users")
	}

	s.logger.WithField("user_count", len(users)).Info("updatable user count")

	for _, user := range users {
		err = s.redis.ZAdd(ctx, internal.UpdateQueue, &redis.Z{Score: float64(time.Now().Unix()), Member: user.ID.String()}).Err()
		if err != nil {
			return errors.Wrap(err, "failed to push user id to processing queue")
		}

		time.Sleep(time.Second)
	}

	return nil

}

type Permission uint

const (
	PermissionHideQueue Permission = iota
	PermissionHideClones
	PermissionHideStandings
	PermissionHideShips
)

func (s *Service) UserSettings(ctx context.Context, id uuid.UUID) (*skillz.UserSettings, error) {

	settings, err := s.cache.UserSettings(ctx, id)
	if err != nil {
		return nil, err
	}

	if settings != nil {
		return settings, nil
	}

	settings, err = s.UserRepository.UserSettings(ctx, id)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, errors.Wrap(err, "failed to fetch settigns from data store")
	}

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	defer func(ctx context.Context, id uuid.UUID, settings *skillz.UserSettings) {
		err = s.cache.SetUserSettings(ctx, id, settings, time.Hour)
		if err != nil {
			s.logger.WithError(err).WithField("user_id", id.String()).Error("failed to cache user settings")
		}

	}(context.Background(), id, settings)

	return settings, nil

}

func (s *Service) CreateUserSettings(ctx context.Context, id uuid.UUID, settings *skillz.UserSettings) error {

	settings.UserID = id

	err := s.UserRepository.CreateUserSettings(ctx, settings)
	if err != nil {
		return errors.Wrap(err, "failed to cache user settings")
	}

	err = s.cache.SetUserSettings(ctx, settings.UserID, settings, time.Hour)
	if err != nil {
		s.logger.WithError(err).WithField("user_id", settings.UserID.String()).Error("failed to cache user settings")
	}

	return nil

}
