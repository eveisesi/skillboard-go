package user

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/eveisesi/skillz"
	"github.com/eveisesi/skillz/internal/alliance"
	"github.com/eveisesi/skillz/internal/auth"
	"github.com/eveisesi/skillz/internal/character"
	"github.com/eveisesi/skillz/internal/corporation"
	"github.com/gofrs/uuid"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/pkg/errors"
)

type API interface {
	User(ctx context.Context, id uuid.UUID) (*skillz.User, error)
	Login(ctx context.Context, code, state string) error
	UserFromToken(ctx context.Context, token jwt.Token) (*skillz.User, error)
	ValidateToken(ctx context.Context, user *skillz.User) error
}

type Service struct {
	auth        auth.API
	character   character.API
	corporation corporation.API
	alliance    alliance.API

	skillz.UserRepository
}

var _ API = new(Service)

func New(
	auth auth.API,
	alliance alliance.API,
	character character.API,
	corporation corporation.API,
	user skillz.UserRepository,
) *Service {
	return &Service{

		alliance:       alliance,
		auth:           auth,
		character:      character,
		corporation:    corporation,
		UserRepository: user,
	}
}

func (s *Service) Login(ctx context.Context, code, state string) error {

	attempt, err := s.auth.AuthAttempt(ctx, state)
	if err != nil {
		return errors.Wrap(err, "failed to fetch attempt with provided state")
	}

	if attempt != nil && attempt.Status == skillz.InvalidAuthStatus {
		return errors.New("request is no longer valid, please try again")
	}

	bearer, err := s.auth.BearerForCode(ctx, code)
	if err != nil {
		return errors.Wrap(err, "failed to exchange authorization code for token")
	}

	token, err := s.auth.ParseAndVerifyToken(ctx, bearer.AccessToken)
	if err != nil {
		return errors.Wrap(err, "failed to verify token")
	}

	user, err := s.UserFromToken(ctx, token)
	if err != nil {
		return errors.Wrap(err, "failed to fetch/provision user for the provided token")
	}

	claims := token.PrivateClaims()
	if _, ok := claims["owner"]; !ok {
		return errors.New("invalid token, owner claims is missing")
	}

	ownerHash := claims["owner"].(string)

	if user.OwnerHash == "" {
		user.OwnerHash = ownerHash
	} else if user.OwnerHash != ownerHash {
		return errors.New("character owner hash mismatch. please contact support")
	}

	if _, ok := claims["scp"]; !ok {
		return errors.New("invalid token, scope claim is missing")
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
		return errors.New("invalid type for scp claim in token.")
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
	case false:
		err = s.UserRepository.UpdateUser(ctx, user)
	}

	return err

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

	user, err := s.UserRepository.UserByCharacterID(ctx, characterID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	if errors.Is(err, sql.ErrNoRows) {
		user = &skillz.User{
			CharacterID: characterID,
		}
	}

	_, err = s.character.Character(ctx, characterID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch user's character")
	}

	return user, nil

}

func (s *Service) ValidateToken(ctx context.Context, user *skillz.User) error {

	updated, err := s.auth.ValidateToken(ctx, user)
	if err != nil {
		return err
	}

	if updated {
		err = s.UserRepository.UpdateUser(ctx, user)
		if err != nil {
			return err
		}
	}

	return nil

}

func memberIDFromSubject(sub string) (uint64, error) {

	parts := strings.Split(sub, ":")

	if len(parts) != 3 {
		return 0, fmt.Errorf("invalid sub format")
	}

	return strconv.ParseUint(parts[2], 10, 64)

}
