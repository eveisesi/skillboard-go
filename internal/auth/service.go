package auth

import (
	"fmt"
	"net/http"

	"github.com/eveisesi/skillz"
	"github.com/eveisesi/skillz/internal/cache"
	"github.com/lestrrat-go/jwx/jwk"
	"golang.org/x/oauth2"
)

type API interface {
	esi
}

type esiAuth struct {
	oauthConfig *oauth2.Config
	jwks        jwk.Set
}

type Service struct {
	env     skillz.Environment
	esiAuth esiAuth
	client  *http.Client
	cache   cache.AuthAPI
}

const jwksURIStr = "https://login.eveonline.com/oauth/jwks"

func New(
	env skillz.Environment,
	client *http.Client,
	cache cache.AuthAPI,
	esiOAuth *oauth2.Config,
) *Service {
	s := &Service{
		env:    env,
		client: client,
		cache:  cache,
		esiAuth: esiAuth{
			oauthConfig: esiOAuth,
		},
	}

	err := s.initializeESIJWKSet()
	if err != nil {
		panic(fmt.Errorf("internal.Auth: %w", err))
	}

	return s
}
