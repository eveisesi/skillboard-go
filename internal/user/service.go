package user

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/eveisesi/skillz"
	"github.com/eveisesi/skillz/internal/auth"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/pkg/errors"
)

type Service struct {
	auth auth.API
}

func New(auth auth.API) *Service {
	return &Service{
		auth: auth,
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

}

func (s *Service) MemberFromToken(ctx context.Context, token jwt.Token) (*skillz.User, error) {

	subject := token.Subject()
	if subject == "" {
		return nil, errors.New("token subject is empty, expected parsable value")
	}

	memberID, err := memberIDFromSubject(subject)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse character id from token subject")
	}

}

func memberIDFromSubject(sub string) (uint, error) {

	parts := strings.Split(sub, ":")

	if len(parts) != 3 {
		return 0, fmt.Errorf("invalid sub format")
	}

	id, err := strconv.ParseUint(parts[2], 10, 32)

	return uint(id), err

}
