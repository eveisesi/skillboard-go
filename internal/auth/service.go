package auth

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/eveisesi/skillz"
	"github.com/eveisesi/skillz/internal/cache"
	"github.com/go-redis/redis/v8"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
	"golang.org/x/oauth2"
)

type API interface {
	InitializeAttempt(ctx context.Context) (*skillz.AuthAttempt, error)
	AuthAttempt(ctx context.Context, state string) (*skillz.AuthAttempt, error)
	UpdateAuthAttempt(ctx context.Context, attempt *skillz.AuthAttempt) error
	AuthorizationURI(ctx context.Context, state string) string
	BearerForCode(ctx context.Context, code string) (*oauth2.Token, error)
	ParseAndVerifyToken(ctx context.Context, t string) (jwt.Token, error)

	// ValidateToken(ctx context.Context, user *skillz.User) (*skillz.User, error)
}

type Service struct {
	oauth    *oauth2.Config
	client   *http.Client
	cache    cache.AuthAPI
	endpoint string // JWKS Endpoint
}

func New(client *http.Client, oauth *oauth2.Config, cache cache.AuthAPI, jwksEndpoint string) *Service {
	return &Service{
		client: client,
		oauth:  oauth,

		cache: cache,

		endpoint: jwksEndpoint,
	}
}

// have func(ctx context.Context, state string) (*github.com/eveisesi/skillz.AuthAttempt, error),
// want func(ctx context.Context, attempt *github.com/eveisesi/skillz.AuthAttempt) error)

func (s *Service) InitializeAttempt(ctx context.Context) (*skillz.AuthAttempt, error) {

	h := hmac.New(sha256.New, nil)
	_, _ = h.Write([]byte(time.Now().Format(time.RFC3339Nano)))
	b := h.Sum(nil)

	attempt := &skillz.AuthAttempt{
		Status: skillz.PendingAuthStatus,
		State:  fmt.Sprintf("%x", string(b)),
	}

	return attempt, s.cache.CreateAuthAttempt(ctx, attempt)

}

func (s *Service) AuthAttempt(ctx context.Context, state string) (*skillz.AuthAttempt, error) {

	attempt, err := s.cache.AuthAttempt(ctx, state)
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, fmt.Errorf("failed to fetch attempt with state of %s: %w", attempt.State, err)
	}

	if err == nil {
		return attempt, nil
	}

	attempt.Status = skillz.InvalidAuthStatus

	return attempt, nil

}

func (s *Service) UpdateAuthAttempt(ctx context.Context, attempt *skillz.AuthAttempt) error {

	// attempt.State = hash

	return s.cache.CreateAuthAttempt(ctx, attempt)

}

func (s *Service) AuthorizationURI(ctx context.Context, state string) string {
	return s.oauth.AuthCodeURL(state)
}

// func (s *Service) ValidateToken(ctx context.Context, user *skillz.User) error {

// 	ctx = context.WithValue(ctx, oauth2.HTTPClient, s.client)

// 	token := &oauth2.Token{
// 		AccessToken:  user.AccessToken,
// 		RefreshToken: user.RefreshToken,
// 		Expiry:       user.Expires,
// 	}

// 	tokenSource := s.oauth.TokenSource(ctx, token)
// 	newToken, err := tokenSource.Token()
// 	if err != nil {
// 		return err
// 	}

// 	if user.AccessToken != newToken.AccessToken {
// 		user.AccessToken = newToken.AccessToken
// 		user.Expires = newToken.Expiry
// 		user.RefreshToken = newToken.RefreshToken
// 	}

// 	return nil

// }

func (s *Service) BearerForCode(ctx context.Context, code string) (*oauth2.Token, error) {
	return s.oauth.Exchange(ctx, code)
}

func (s *Service) ParseAndVerifyToken(ctx context.Context, t string) (jwt.Token, error) {

	set, err := s.getSet()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch jwks: %w", err)
	}

	token, err := jwt.ParseString(t, jwt.WithKeySet(set))
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	// err = jwt.Validate(token, jwt.WithIssuer("login.eveonline.com"), jwt.WithClaimValue("azp", "27a0d315019c4d15bf909abefe67282b"))
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to validate token: %w", err)
	// }

	return token, nil

}

func (s *Service) getSet() (jwk.Set, error) {

	ctx := context.Background()

	b, err := s.cache.JSONWebKeySet(ctx)
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, fmt.Errorf("unexpected error occured querying redis for jwks: %w", err)
	}

	if b == nil {
		res, err := s.client.Get(s.endpoint)
		if err != nil {
			return nil, fmt.Errorf("unable to retrieve jwks from sso: %w", err)
		}

		if res.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("unexpected status code recieved while fetching jwks. %d", res.StatusCode)
		}

		buf, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read jwk response body: %w", err)
		}

		err = s.cache.SaveJSONWebKeySet(ctx, buf)
		if err != nil {
			return nil, fmt.Errorf("failed to save jwks to cache layer: %w", err)
		}

		b = buf
	}

	return jwk.Parse(b)

}
