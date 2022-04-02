package user

import (
	"context"
	"crypto/sha1"
	"database/sql"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/eveisesi/skillz"
	"github.com/eveisesi/skillz/internal"
	"github.com/eveisesi/skillz/internal/alliance"
	"github.com/eveisesi/skillz/internal/auth"
	"github.com/eveisesi/skillz/internal/cache"
	"github.com/eveisesi/skillz/internal/character"
	"github.com/eveisesi/skillz/internal/clone"
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
	LoadUserAll(ctx context.Context, id string) (*skillz.User, error)

	UserFromToken(ctx context.Context, token jwt.Token) (*skillz.User, error)
	ValidateCurrentToken(ctx context.Context, user *skillz.User) error

	User(ctx context.Context, id string, rels ...UserRel) (*skillz.User, error)
	RefreshUser(ctx context.Context, user *skillz.User) error
	UserByCharacterID(ctx context.Context, characterID uint64) (*skillz.User, error)
	UserByUUID(ctx context.Context, id uuid.UUID) (*skillz.User, error)
	SearchUsers(ctx context.Context, q string) ([]*skillz.UserSearchResult, error)
	CreateUser(ctx context.Context, user *skillz.User) error
	DeleteUser(ctx context.Context, user *skillz.User) error

	Recent(ctx context.Context) ([]*skillz.User, []*skillz.User, error)
	// NewUsersBySP(ctx context.Context) ([]*skillz.UserWithSkillMeta, error)
	ResetUserCache(ctx context.Context, user *skillz.User) error
	ProcessUpdatableUsers(ctx context.Context) error

	UserSettings(ctx context.Context, id string) (*skillz.UserSettings, error)
	CreateUserSettings(ctx context.Context, userID string, settings *skillz.UserSettings) error
}

type Service struct {
	redis  *redis.Client
	logger *logrus.Logger
	cache  cache.UserAPI

	auth        auth.API
	character   character.API
	corporation corporation.API
	alliance    alliance.API

	clones clone.API
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
	clones clone.API,
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
		clones:         clones,
		UserRepository: user,
	}
}

type UserRel uint

var allRels = []UserRel{
	UserCharacterRel, UserAttributesRel,
	UserSkillsRel, UserFlyableRel,
	UserSkillQueueRel, UserSkillMetaRel,
}

func (r UserRel) Valid() bool {
	for _, rel := range allRels {
		if r == rel {
			return true
		}
	}

	return false
}

const (
	UserCharacterRel UserRel = iota
	UserAttributesRel
	UserSkillsRel
	UserFlyableRel
	UserSkillQueueRel
	UserSkillMetaRel
)

func (s *Service) User(ctx context.Context, id string, rels ...UserRel) (*skillz.User, error) {

	var mx = new(sync.Mutex)
	var wg = new(sync.WaitGroup)

	user, err := s.UserRepository.User(ctx, id)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, errors.Wrap(err, "unexpected error encountered fetching user")
	}

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}

	user.Settings, err = s.UserSettings(ctx, id)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, errors.Wrap(err, "unexpected error encountered fetching user settings")
	}

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}

	entry := s.logger.
		WithField("userID", user.ID)

	for _, rel := range rels {
		if !rel.Valid() {
			continue
		}

		wg.Add(1)
		switch rel {
		case UserCharacterRel:
			go s.LoadCharacter(ctx, user, entry, mx, wg)
		case UserAttributesRel:
			go s.LoadAttributes(ctx, user, entry, mx, wg)
		case UserSkillsRel:
			go s.LoadSkillGrouped(ctx, user, entry, mx, wg)
		case UserFlyableRel:
			go s.LoadFlyable(ctx, user, entry, mx, wg)
		case UserSkillQueueRel:
			go s.LoadSkillQueue(ctx, user, entry, mx, wg)
		case UserSkillMetaRel:
			go s.LoadSkillMeta(ctx, user, entry, mx, wg)
		}
	}

	wg.Wait()

	return user, nil

}

func (s *Service) Recent(ctx context.Context) ([]*skillz.User, []*skillz.User, error) {

	var users []*skillz.User
	var err error

	users, err = s.UserRepository.NewUsersBySP(ctx)
	if err != nil {
		return nil, nil, err
	}

	var wg = new(sync.WaitGroup)
	var mx = new(sync.Mutex)
	for i, user := range users {
		wg.Add(1)
		go func(ctx context.Context, i int, user *skillz.User, wg *sync.WaitGroup) {
			defer wg.Done()
			u, err := s.User(ctx, user.ID, UserCharacterRel, UserSkillQueueRel, UserSkillMetaRel)
			if err != nil {
				s.logger.WithError(err).Error("failed to load character")
				return
			}

			mx.Lock()
			defer mx.Unlock()
			users[i] = u

		}(context.Background(), i, user, wg)

	}

	wg.Wait()

	var populated = make([]*skillz.User, 0, len(users))
	for _, user := range users {
		if user.Character == nil {
			continue
		}

		populated = append(populated, user)
	}

	highlighted := append(make([]*skillz.User, 0, len(populated)), populated...)

	sort.Slice(highlighted, func(i, j int) bool {
		return highlighted[i].Meta.TotalSP > highlighted[j].Meta.TotalSP
	})

	if len(highlighted) > 6 {
		highlighted = highlighted[0:6]
	}

	return highlighted, populated, nil

}

var ErrUserNotFound = errors.New("user does not exist")

func (s *Service) LoadUserAll(ctx context.Context, id string) (*skillz.User, error) {

	var mx = new(sync.Mutex)
	var wg = new(sync.WaitGroup)

	user, err := s.UserRepository.User(ctx, id)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, errors.Wrap(err, "unexpected error encountered fetch user")
	}

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}

	user.Settings, err = s.UserSettings(ctx, user.ID)
	if err != nil {
		return nil, errors.Wrap(err, "unexpected error encountered fetch user")
	}

	entry := s.logger.
		WithField("userID", user.ID)

	wg.Add(1)
	go s.LoadCharacter(ctx, user, entry, mx, wg)

	wg.Add(1)
	go s.LoadSkillMeta(ctx, user, entry, mx, wg)

	if !user.Settings.HideSkills {
		wg.Add(1)
		go s.LoadSkillGrouped(ctx, user, entry, mx, wg)
	}

	if !user.Settings.HideFlyable {
		wg.Add(1)
		go s.LoadFlyable(ctx, user, entry, mx, wg)
	}

	if !user.Settings.HideQueue {
		wg.Add(1)
		go s.LoadSkillQueue(ctx, user, entry, mx, wg)
	}

	if !user.Settings.HideAttributes {
		wg.Add(1)
		go s.LoadAttributes(ctx, user, entry, mx, wg)
	}

	if !user.Settings.HideImplants {
		wg.Add(1)
		go s.LoadImplants(ctx, user, entry, mx, wg)

	}

	wg.Wait()

	return user, nil

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
		s.logger.WithFields(logrus.Fields{
			"claims":        claims,
			"character_id":  user.CharacterID,
			"userOwnerHash": user.OwnerHash,
		}).Error("character owner hash mismatch. please contact support")
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

	err = s.UserRepository.CreateUser(ctx, user)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create user in data store")
	}

	err = s.auth.DeleteAuthAttempt(ctx, attempt)
	if err != nil {
		return nil, errors.Wrap(err, "failed to remove auth attempt from cache")
	}

	err = s.redis.ZAdd(ctx, internal.UpdateQueue, &redis.Z{Score: float64(time.Now().Unix()), Member: user.ID}).Err()
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
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, errors.Wrap(err, "unexpected error encountered fetch user settings")
	}

	if user.Settings == nil {
		err = s.CreateUserSettings(ctx, user.ID, &skillz.UserSettings{
			Visibility: skillz.VisibilityPrivate,
		})
		if err != nil {
			return nil, errors.Wrap(err, "unexpected error encountered whilst attempting to create user settings")
		}
	}

	return user, nil

}

func (s *Service) RefreshUser(ctx context.Context, user *skillz.User) error {

	err := s.redis.ZAdd(ctx, internal.UpdateQueue, &redis.Z{Score: float64(time.Now().Unix()), Member: user.ID}).Err()
	if err != nil {
		return errors.Wrap(err, "failed to push user id to processing queue")
	}

	return nil

}

func (s *Service) UserFromToken(ctx context.Context, token jwt.Token) (*skillz.User, error) {

	subject := token.Subject()
	if subject == "" {
		return nil, errors.New("token subject is empty, expected parsable value")
	}

	characterID, err := memberIDFromSubject(subject)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse character id from token subject")
	}

	claims := token.PrivateClaims()
	if _, ok := claims["owner"]; !ok {
		return nil, errors.New("invalid token, owner claim is missing")
	}

	ownerHash := claims["owner"].(string)

	userID := hashedUserID(ownerHash, strconv.FormatUint(characterID, 10))

	user, err := s.UserRepository.UserByCharacterID(ctx, characterID)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		user = &skillz.User{
			ID:          userID,
			CharacterID: characterID,
			IsNew:       true,
		}
		err = nil
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

func hashedUserID(ownerHash, characterID string) string {
	s := []string{ownerHash, characterID}
	return fmt.Sprintf("%x", sha1.Sum([]byte(strings.Join(s, ":"))))

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

	return s.UserRepository.CreateUser(ctx, user)

}

func memberIDFromSubject(sub string) (uint64, error) {

	parts := strings.Split(sub, ":")

	if len(parts) != 3 {
		return 0, fmt.Errorf("invalid sub format")
	}

	return strconv.ParseUint(parts[2], 10, 64)

}

func (s *Service) ProcessUpdatableUsers(ctx context.Context) error {

	users, err := s.UserRepository.UsersSortedByProcessedAtLimit(ctx)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return errors.Wrap(err, "failed to query users table for processable users")
	}

	s.logger.WithField("user_count", len(users)).Info("updatable user count")

	for _, user := range users {
		err = s.redis.ZAdd(ctx, internal.UpdateQueue, &redis.Z{Score: float64(time.Now().Unix()), Member: user.ID}).Err()
		if err != nil {
			return errors.Wrap(err, "failed to push user id to processing queue")
		}
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

func (s *Service) UserSettings(ctx context.Context, id string) (*skillz.UserSettings, error) {

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

	defer func(ctx context.Context, id string, settings *skillz.UserSettings) {
		err = s.cache.SetUserSettings(ctx, id, settings, time.Hour)
		if err != nil {
			s.logger.WithError(err).WithField("user_id", id).Error("failed to cache user settings")
		}

	}(context.Background(), id, settings)

	return settings, nil

}

func (s *Service) CreateUserSettings(ctx context.Context, id string, settings *skillz.UserSettings) error {

	settings.UserID = id
	if settings.VisibilityToken == "" {
		settings.VisibilityToken = fmt.Sprintf("%x", sha1.Sum([]byte(time.Now().Format(time.RFC3339Nano))))
	}

	err := s.UserRepository.CreateUserSettings(ctx, settings)
	if err != nil {
		return errors.Wrap(err, "failed to cache user settings")
	}

	defer func() {
		err = s.cache.SetUserSettings(ctx, settings.UserID, settings, time.Hour)
		if err != nil {
			s.logger.WithError(err).WithField("user_id", settings.UserID).Error("failed to cache user settings")
		}
	}()

	return nil

}

func (s *Service) ResetUserCache(ctx context.Context, user *skillz.User) error {
	return s.cache.ResetUserCache(ctx, user)
}
