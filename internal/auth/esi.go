package auth

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/eveisesi/skillz"
	"github.com/go-redis/redis/v8"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

type esi interface {
	InitializeAttempt(ctx context.Context) (*skillz.AuthAttempt, error)
	DeleteAuthAttempt(ctx context.Context, attempt *skillz.AuthAttempt) error
	AuthAttempt(ctx context.Context, state string) (*skillz.AuthAttempt, error)
	AuthorizationURI(ctx context.Context, state string) string
	ValidateESITokenForUser(ctx context.Context, user *skillz.User) (bool, *oauth2.Token, error)
	BearerForESICode(ctx context.Context, code string) (*oauth2.Token, error)
	ParseAndVerifyESIToken(ctx context.Context, t string) (jwt.Token, error)
}

func (s *Service) InitializeAttempt(ctx context.Context) (*skillz.AuthAttempt, error) {

	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(time.Now().Format(time.RFC3339Nano))))

	attempt := &skillz.AuthAttempt{
		Status: skillz.PendingAuthStatus,
		State:  hash,
	}

	return attempt, s.cache.CreateAuthAttempt(ctx, attempt)

}

func (s *Service) AuthAttempt(ctx context.Context, state string) (*skillz.AuthAttempt, error) {

	attempt, err := s.cache.AuthAttempt(ctx, state)
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, fmt.Errorf("failed to fetch attempt with state of %s: %w", attempt.State, err)
	}
	if attempt == nil {
		return nil, errors.New("attempt is invalid or has expired. please try again")
	}
	if err == nil {
		return attempt, nil
	}

	attempt.Status = skillz.InvalidAuthStatus

	return attempt, nil

}

func (s *Service) DeleteAuthAttempt(ctx context.Context, attempt *skillz.AuthAttempt) error {
	return s.cache.DeleteAuthAttempt(ctx, attempt)
}

func (s *Service) AuthorizationURI(ctx context.Context, state string) string {
	return s.esiAuth.oauthConfig.AuthCodeURL(state)
}

func isExpiredJWTError(err error) bool {
	return strings.Contains(err.Error(), "exp not satisfied")
}

// ValidateESITokenForUser takes a user token and, if expired, refreshed it.
func (s *Service) ValidateESITokenForUser(ctx context.Context, user *skillz.User) (bool, *oauth2.Token, error) {

	_, err := s.ParseAndVerifyESIToken(ctx, user.AccessToken)
	if err != nil && !isExpiredJWTError(err) {
		return false, nil, err
	}

	ctx = context.WithValue(ctx, oauth2.HTTPClient, s.client)

	token := &oauth2.Token{
		AccessToken:  user.AccessToken,
		RefreshToken: user.RefreshToken,
		Expiry:       user.Expires,
	}

	tokenSource := s.esiAuth.oauthConfig.TokenSource(ctx, token)
	newToken, err := tokenSource.Token()
	if err != nil {
		return false, nil, err
	}

	return user.AccessToken != newToken.AccessToken, newToken, nil

}

func (s *Service) BearerForESICode(ctx context.Context, code string) (*oauth2.Token, error) {
	return s.esiAuth.oauthConfig.Exchange(ctx, code)
}

func (s *Service) ParseAndVerifyESIToken(ctx context.Context, t string) (jwt.Token, error) {

	token, err := jwt.ParseString(t, jwt.WithKeySet(s.esiAuth.jwks))
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	err = jwt.Validate(token, jwt.WithIssuer("login.eveonline.com"), jwt.WithClaimValue("azp", s.esiAuth.oauthConfig.ClientID))
	if err != nil {
		return token, fmt.Errorf("failed to validate token: %w", err)
	}

	return token, nil

}

func (s *Service) initializeESIJWKSet(keyEndpoint *url.URL) error {

	res, err := s.client.Get(keyEndpoint.String())
	if err != nil {
		return fmt.Errorf("unable to retrieve jwks from sso: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code recieved while fetching jwks. %d", res.StatusCode)
	}

	buf, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("failed to read jwk response body: %w", err)
	}

	set, err := jwk.Parse(buf)
	if err != nil {
		fmt.Println(buf)
		return fmt.Errorf("failed to parse jwk retrieved from ESI SSO: %w", err)
	}

	s.esiAuth.jwks = set
	return nil
}
