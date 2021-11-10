package auth

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/eveisesi/skillz"
	"github.com/eveisesi/skillz/internal/cache"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
	"golang.org/x/oauth2"
)

type API interface {
	InitializeAttempt(ctx context.Context) (*skillz.AuthAttempt, error)
	AuthAttempt(ctx context.Context, hash string) (*skillz.AuthAttempt, error)
	UpdateAuthAttempt(ctx context.Context, hash string, attempt *skillz.AuthAttempt) (*skillz.AuthAttempt, error)
	AuthorizationURI(ctx context.Context, state string, scopes []string) string
	ValidateToken(ctx context.Context, member *skillz.User) (*skillz.User, error)
	BearerForCode(ctx context.Context, code string) (*oauth2.Token, error)
	ParseAndVerifyToken(ctx context.Context, t string) (jwt.Token, error)
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

func (s *Service) InitializeAttempt(ctx context.Context) (*skillz.AuthAttempt, error) {

	h := hmac.New(sha256.New, nil)
	_, _ = h.Write([]byte(time.Now().Format(time.RFC3339Nano)))
	b := h.Sum(nil)

	attempt := &skillz.AuthAttempt{
		Status: skillz.PendingAuthStatus,
		State:  fmt.Sprintf("%x", string(b)),
	}

	return s.cache.CreateAuthAttempt(ctx, attempt)

}

func (s *Service) AuthAttempt(ctx context.Context, hash string) (*skillz.AuthAttempt, error) {

	attempt, err := s.cache.AuthAttempt(ctx, hash)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch attempt with hash of %s: %w", hash, err)
	}

	if attempt == nil {
		attempt = new(skillz.AuthAttempt)
		attempt.Status = skillz.InvalidAuthStatus
	}

	return attempt, nil

}

func (s *Service) UpdateAuthAttempt(ctx context.Context, hash string, attempt *skillz.AuthAttempt) (*skillz.AuthAttempt, error) {

	attempt.State = hash

	attempt, err := s.cache.CreateAuthAttempt(ctx, attempt)
	if err != nil {
		return nil, err
	}

	return attempt, nil

}

func (s *Service) AuthorizationURI(ctx context.Context, state string, scopes []string) string {
	strScopes := ""
	if len(scopes) > 0 {
		strScopes = strings.Join(scopes, " ")
	}

	return s.oauth.AuthCodeURL(state, oauth2.SetAuthURLParam("scope", strScopes))
}

func (s *Service) ValidateToken(ctx context.Context, member *skillz.User) (*skillz.User, error) {

	ctx = context.WithValue(ctx, oauth2.HTTPClient, s.client)

	token := new(oauth2.Token)
	token.AccessToken = member.AccessToken
	token.RefreshToken = member.RefreshToken.String
	token.Expiry = member.Expires.Time

	tokenSource := s.oauth.TokenSource(ctx, token)
	newToken, err := tokenSource.Token()
	if err != nil {
		return nil, err
	}

	if member.AccessToken.String != newToken.AccessToken {
		member.AccessToken.SetValid(newToken.AccessToken)
		member.Expires.SetValid(newToken.Expiry)
		member.RefreshToken.SetValid(newToken.RefreshToken)
	}

	return member, nil

}

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
	if err != nil {
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
