package user

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/eveisesi/skillz"
	"github.com/eveisesi/skillz/internal/alliance"
	"github.com/eveisesi/skillz/internal/auth"
	"github.com/eveisesi/skillz/internal/character"
	"github.com/eveisesi/skillz/internal/corporation"
	"github.com/go-redis/redis/v8"
	"github.com/gofrs/uuid"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/pkg/errors"
)

type API interface {
	Login(ctx context.Context, code, state string) (*skillz.AuthAttempt, error)
	UserFromToken(ctx context.Context, token jwt.Token) (*skillz.User, error)
	ValidateCurrentToken(ctx context.Context, user *skillz.User) error

	User(ctx context.Context, id uuid.UUID) (*skillz.User, error)
	UpdateUser(ctx context.Context, user *skillz.User) error
}

type Service struct {
	redis *redis.Client

	auth        auth.API
	character   character.API
	corporation corporation.API
	alliance    alliance.API

	skillz.UserRepository
}

var _ API = new(Service)

func New(
	redis *redis.Client,
	auth auth.API,
	alliance alliance.API,
	character character.API,
	corporation corporation.API,
	user skillz.UserRepository,
) *Service {
	return &Service{
		redis:          redis,
		alliance:       alliance,
		auth:           auth,
		character:      character,
		corporation:    corporation,
		UserRepository: user,
	}
}

func (s *Service) Login(ctx context.Context, code, state string) (*skillz.AuthAttempt, error) {

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

	switch user.ID == uuid.Nil {
	case true:
		user.ID = uuid.Must(uuid.NewV4())
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

	userToken, err := s.auth.TokenFromUser(ctx, user)
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate token for user")
	}

	err = s.auth.DeleteAuthAttempt(ctx, attempt)
	if err != nil {
		return nil, errors.Wrap(err, "failed to remove auth attempt from cache")
	}

	attempt.Token.SetValid(userToken)

	return attempt, err

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

	case "http://192.168.1.242:54405":
		userID, ierr := uuid.FromString(token.Subject())
		if ierr != nil {
			return nil, errors.Wrap(err, "failed to parse userID as valid uuid")
		}

		user, err = s.User(ctx, userID)

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
