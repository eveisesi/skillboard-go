package user

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/eveisesi/skillz"
	"github.com/eveisesi/skillz/internal"
	"github.com/eveisesi/skillz/internal/auth"
	"github.com/eveisesi/skillz/internal/boiler"
	"github.com/eveisesi/skillz/internal/cache"
	ierrors "github.com/eveisesi/skillz/internal/errors"
	"github.com/go-redis/redis/v8"
	"github.com/gobuffalo/buffalo"
	"github.com/gofrs/uuid"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type API interface {
	Login(ctx context.Context, state, code string) (*boiler.User, error)
	User(ctx context.Context, id uuid.UUID) (*boiler.User, error)

	// UserFromToken(ctx context.Context, token jwt.Token) (*skillz.User, error)
	// ValidateCurrentToken(ctx context.Context, user *skillz.User) error

	// User(ctx context.Context, id uuid.UUID) (*skillz.User, error)
	// RefreshUser(ctx context.Context, user *skillz.User) error
	// UserByCharacterID(ctx context.Context, characterID uint64) (*skillz.User, error)
	// SearchUsers(ctx context.Context, q string) ([]*skillz.UserSearchResult, error)
	// UpdateUser(ctx context.Context, user *skillz.User) error

	// NewUsersBySP(ctx context.Context) ([]*skillz.UserWithSkillMeta, error)

	// ProcessUpdatableUsers(ctx context.Context) error

	// UserSettings(ctx context.Context, id uuid.UUID) (*skillz.UserSettings, error)
	// CreateUserSettings(ctx context.Context, userID uuid.UUID, settings *skillz.UserSettings) error
}

type service struct {
	redis  *redis.Client
	logger *logrus.Logger
	cache  cache.UserAPI

	auth auth.API
}

func NewService(redis *redis.Client, logger *logrus.Logger, cache cache.UserAPI, auth auth.API) API {

	return &service{
		redis, logger, cache, auth,
	}

}

func (s *service) Login(ctx context.Context, state, code string) (*boiler.User, error) {

	attempt, err := s.auth.AuthAttempt(ctx, state)
	if err != nil {
		return nil, ierrors.NewBuffaloHTTPError(http.StatusBadRequest, fmt.Errorf("invalid or missing state, please try again"))
	}

	if attempt == nil {
		return nil, buffalo.HTTPError{Status: http.StatusBadRequest, Cause: fmt.Errorf("invalid or missing state, please try again")}
	}

	if attempt != nil && attempt.Status == skillz.InvalidAuthStatus {
		return nil, buffalo.HTTPError{Status: http.StatusBadRequest, Cause: fmt.Errorf("state is no longer valid, please attempt to log in again. You have five minutes to complete the authentication flow")}
	}

	bearer, err := s.auth.BearerForESICode(ctx, code)
	if err != nil {
		return nil, ierrors.NewBuffaloHTTPError(http.StatusBadRequest, fmt.Errorf("failed to exchange authorization code for token: %w", err))
	}

	token, err := s.auth.ParseAndVerifyESIToken(ctx, bearer.AccessToken)
	if err != nil {
		return nil, ierrors.NewBuffaloHTTPError(http.StatusBadRequest, fmt.Errorf("failed to verify token: %w", err))
	}

	user, err := s.UserFromToken(ctx, token)
	if err != nil {
		return nil, ierrors.NewBuffaloHTTPError(http.StatusBadRequest, fmt.Errorf("failed to fetch/provision user for the provided token: %w", err))
	}

	claims := token.PrivateClaims()
	if _, ok := claims["owner"]; !ok {
		return nil, ierrors.NewBuffaloHTTPError(http.StatusBadRequest, fmt.Errorf("invalid token, owner claims is missing"))
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

	scopes, _ := json.Marshal(scp)

	user.Scopes = scopes
	user.AccessToken = bearer.AccessToken
	user.RefreshToken = bearer.RefreshToken
	user.Expires = bearer.Expiry
	user.LastLogin = time.Now()

	switch user.IsNew {
	case true:

		err = user.InsertG(ctx, boil.Infer())
		if err != nil {
			return nil, errors.Wrap(err, "failed to insert user")
		}

		settings := &boiler.UserSetting{
			UserID: user.ID,
		}

		err = settings.InsertG(ctx, boil.Infer())
		if err != nil {
			return nil, errors.Wrap(err, "failed to create user settings for new user")
		}

		user.R.UserSetting = settings

		defer func(ctx context.Context) {
			err = s.cache.BustNewUsersBySP(ctx)
			if err != nil {
				s.logger.WithError(err).Error("failed to bust new users by sp cache")
			}
		}(context.Background())
	case false:
		_, err = user.UpdateG(ctx, boil.Infer())
		if err != nil {
			return nil, errors.Wrap(err, "failed to update user")
		}
	}

	err = s.redis.ZAdd(ctx, internal.UpdateQueue, &redis.Z{Score: float64(time.Now().Unix()), Member: user.ID.String()}).Err()
	if err != nil {
		return nil, errors.Wrap(err, "failed to push user id to processing queue")
	}

	return user, err

}

func (s *service) User(ctx context.Context, id uuid.UUID) (*boiler.User, error) {

	// qms := append(make([]qm.QueryMod, 0, 2), )
	// for _, rel := range rels {
	// 	qms = append(qms, qm.Load(rel))
	// }

	initial, err := boiler.Users(
		qm.Load(boiler.UserRels.UserSetting), boiler.UserWhere.ID.EQ(id),
	).OneG(ctx)

	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch user by ID")
	}

	var strSlcScopes = make([]string, 0)
	err = json.Unmarshal(initial.Scopes, &strSlcScopes)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode user scopes")
	}

	var modifiers = append(
		make([]qm.QueryMod, 0, len(strSlcScopes)+2),
		qm.Load(boiler.UserRels.UserSetting),
		boiler.UserWhere.ID.EQ(id),
	)
	for _, scope := range strSlcScopes {
		skScope := skillz.Scope(scope)
		switch skScope {
		case skillz.ReadSkillsV1:
			modifiers = append(
				modifiers,
				qm.Load(boiler.UserRels.CharacterCharacterSkillMetum),
				qm.Load(boiler.UserRels.CharacterCharacterSkills),
				qm.Load(boiler.UserRels.CharacterCharacterAttribute),
				qm.Load(boiler.UserRels.CharacterCharacterFlyableShips),
			)
		case skillz.ReadSkillQueueV1:
			modifiers = append(
				modifiers,
				// qm.Load(boiler.)
			)
		}
	}

	user, err := boiler.Users(qms...).OneG(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user")
	}

	return user, nil

}

func (s *service) UserFromToken(ctx context.Context, token jwt.Token) (*boiler.User, error) {
	subject := token.Subject()
	if subject == "" {
		return nil, ierrors.NewBuffaloHTTPError(http.StatusBadRequest, fmt.Errorf("token subject is empty, expected parsable value"))
	}

	characterID, err := memberIDFromSubject(subject)
	if err != nil {
		return nil, ierrors.NewBuffaloHTTPError(http.StatusBadRequest, fmt.Errorf("failed to parse character id from token subject: %w", err))
	}

	user, err := boiler.Users(
		boiler.UserWhere.CharacterID.EQ(characterID),
		qm.Load(boiler.UserRels.UserSetting),
	).OneG(ctx)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, ierrors.NewBuffaloHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to query user: %w", err))
	}

	if errors.Is(err, sql.ErrNoRows) {
		user = &boiler.User{
			ID:          uuid.Must(uuid.NewV4()),
			CharacterID: characterID,
			IsNew:       true,
		}
	}

	return user, nil

}

func memberIDFromSubject(sub string) (uint64, error) {

	parts := strings.Split(sub, ":")

	if len(parts) != 3 {
		return 0, fmt.Errorf("invalid sub format")
	}

	return strconv.ParseUint(parts[2], 10, 64)

}
