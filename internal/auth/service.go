package auth

import (
	"crypto/rsa"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/eveisesi/skillz"
	"github.com/eveisesi/skillz/internal/cache"
	"github.com/lestrrat-go/jwx/jwk"
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
	rsaKey       *rsa.PrivateKey
	cookieExpiry time.Duration
	cookieDomain string
}

type Service struct {
	env      skillz.Environment
	userAuth userAuth
	esiAuth  esiAuth
	client   *http.Client
	cache    cache.AuthAPI
}

func New(
	env skillz.Environment,
	client *http.Client,
	cache cache.AuthAPI,
	esiOAuth *oauth2.Config,
	userPrivateKey *rsa.PrivateKey,
	esiAuthJWKSEndpoint *url.URL,
	cookieDomain string,
	cookieExpiry time.Duration,
) *Service {
	s := &Service{
		env:    env,
		client: client,
		cache:  cache,
		esiAuth: esiAuth{
			oauthConfig: esiOAuth,
		},
		userAuth: userAuth{
			rsaKey:       userPrivateKey,
			cookieExpiry: cookieExpiry,
			cookieDomain: cookieDomain,
		},
	}

	err := s.initializeESIJWKSet(esiAuthJWKSEndpoint)
	if err != nil {
		panic(fmt.Errorf("internal.Auth: %w", err))
	}

	return s
}
