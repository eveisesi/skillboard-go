package auth

import (
	"context"
	"time"

	"github.com/gofrs/uuid"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jws"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/pkg/errors"
)

type user interface {
	UserToken(ctx context.Context, userID uuid.UUID) (string, error)
	UserIDFromToken(ctx context.Context, raw string) (uuid.UUID, error)
}

const (
	keyUserID = "userID"
)

func (s *Service) UserToken(ctx context.Context, userID uuid.UUID) (string, error) {

	var err error

	token := jwt.New()
	err = token.Set(jwt.AudienceKey, s.userAuth.tokenDomain)
	if err != nil {
		return "", errors.Wrap(err, "failed to set aud key on token")
	}

	err = token.Set(jwt.IssuerKey, s.userAuth.tokenDomain)
	if err != nil {
		return "", errors.Wrap(err, "failed to set iss key on token")
	}

	err = token.Set(jwt.IssuedAtKey, time.Now().Unix())
	if err != nil {
		return "", errors.Wrap(err, "failed to set iat key on token")
	}

	err = token.Set(jwt.ExpirationKey, time.Now().Add(s.userAuth.tokenExpiry).Unix())
	if err != nil {
		return "", errors.Wrap(err, "failed to set exp key on token")
	}

	err = token.Set(keyUserID, userID.String())
	if err != nil {
		return "", errors.Wrap(err, "failed to set custom key userID on token")
	}

	tokenHeaders := jws.NewHeaders()
	err = tokenHeaders.Set(jws.KeyIDKey, s.userAuth.tokenKid.String())
	if err != nil {
		return "", errors.Wrap(err, "failed to set custom key userID on token")
	}

	tokenBytes, err := jwt.Sign(token, jwa.RS256, s.userAuth.rsaKey, jwt.WithHeaders(tokenHeaders))
	if err != nil {
		return "", errors.Wrap(err, "failed to sign token")
	}

	return string(tokenBytes), nil

}

func (s *Service) UserIDFromToken(ctx context.Context, raw string) (uuid.UUID, error) {

	token, err := jwt.Parse(
		[]byte(raw),
		jwt.WithAudience(s.userAuth.tokenDomain),
		jwt.WithIssuer(s.userAuth.tokenDomain),
		jwt.WithValidate(true),
		jwt.WithKeySet(s.userAuth.publicJWKSet),
	)
	if err != nil {
		return uuid.Nil, errors.Wrap(err, "failed to validate jwt")
	}

	userIDStr, ok := token.Get(keyUserID)
	if !ok {
		return uuid.Nil, errors.Wrap(err, "token missing required field in private claims")
	}

	userID, err := uuid.FromString(userIDStr.(string))
	if err != nil {
		return uuid.Nil, errors.Wrap(err, "failed to parse user id from token to valid uuid")
	}

	return userID, nil

}
