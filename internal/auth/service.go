// Package auth performs to primary functions
//
// Verifies Authorization State is still validate during an authorization code exchange
// It facilitates that Authoriztion Code Exchange
// Ensures that the Bearer Token that we received back from the IDP is a valid, signed JWT Token.
// Stores that retrived Token in User State
// Generates User Tokens, with the proper claims necessary to properly validate these tokens later on
// Validates User Tokens on inbound requests
package auth

import (
	"fmt"
	"net/http"
	"time"

	"github.com/eveisesi/skillz/internal/cache"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
	"golang.org/x/oauth2"
)

type API interface {
	esi
	user
}

type esiAuth struct {
	oauthConfig *oauth2.Config
	jwks        jwk.Set
}

type userAuth struct {
	jwkKey      jwk.RSAPrivateKey
	jwks        jwk.Set
	parseOpts   []jwt.ParseOption
	tokenIssuer string
	tokenAud    string
	tokenExpiry time.Duration
}

type Service struct {
	userAuth userAuth
	esiAuth  esiAuth
	client   *http.Client
	cache    cache.AuthAPI
}

func New(
	client *http.Client, cache cache.AuthAPI,
	esiOAuth *oauth2.Config, userPrivateKey jwk.RSAPrivateKey,
	esiAuthJWKSEndpoint, userTokenIssuer, userTokenAud string,
	userTokenExpiry time.Duration,
) *Service {
	userPublicKey, _ := userPrivateKey.PublicKey()
	set := jwk.NewSet()
	set.Add(userPublicKey)

	s := &Service{
		client: client,
		cache:  cache,
		esiAuth: esiAuth{
			oauthConfig: esiOAuth,
		},
		userAuth: userAuth{
			jwkKey:      userPrivateKey,
			jwks:        set,
			tokenAud:    userTokenAud,
			tokenIssuer: userTokenIssuer,
			tokenExpiry: userTokenExpiry,
			parseOpts: []jwt.ParseOption{
				jwt.WithIssuer(userTokenIssuer),
				jwt.WithAudience(userTokenAud),
				jwt.WithKeySet(set),
				jwt.WithValidate(true),
			},
		},
	}

	err := s.initializeESIJWKSet(esiAuthJWKSEndpoint)
	if err != nil {
		panic(fmt.Errorf("internal.Auth: %w", err))
	}

	return s
}
