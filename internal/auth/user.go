package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/eveisesi/skillz"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/pkg/errors"
)

type user interface {
	GetPublicUserJWKS() jwk.Set
	ParseAndVerifyUserToken(ctx context.Context, t string) (jwt.Token, error)
	TokenFromUser(ctx context.Context, user *skillz.User) (string, error)
}

func (s *Service) GetPublicUserJWKS() jwk.Set {
	return s.userAuth.jwks
}

func (s *Service) ParseAndVerifyUserToken(ctx context.Context, t string) (jwt.Token, error) {

	token, err := jwt.ParseString(t, s.userAuth.parseOpts...)
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	return token, nil

}

func (s *Service) TokenFromUser(ctx context.Context, user *skillz.User) (string, error) {

	token := jwt.New()
	err := token.Set(jwt.AudienceKey, s.userAuth.tokenAud)
	if err != nil {
		return "", errors.Wrap(err, "failed to set audience on token")
	}

	err = token.Set(jwt.IssuerKey, s.userAuth.tokenIssuer)
	if err != nil {
		return "", errors.Wrap(err, "failed to set issuer on token")
	}

	err = token.Set(jwt.ExpirationKey, time.Now().Add(s.userAuth.tokenExpiry))
	if err != nil {
		return "", errors.Wrap(err, "failed to set expiry on token")
	}

	err = token.Set(jwt.SubjectKey, user.ID.String())
	if err != nil {
		return "", errors.Wrap(err, "failed to set subject on token")
	}

	err = token.Set("owner", user.OwnerHash)
	if err != nil {
		return "", errors.Wrap(err, "failed to set private claim 'owner' on token")
	}

	signed, err := jwt.Sign(token, jwa.RS256, s.userAuth.jwkKey)
	if err != nil {
		return "", errors.Wrap(err, "failed to sign token")
	}

	return string(signed), err

}
