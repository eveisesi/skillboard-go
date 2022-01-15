package auth

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/eveisesi/skillz"
	"github.com/eveisesi/skillz/internal/cache"
	"github.com/gofrs/uuid"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

type API interface {
	esi
	user
	GetPublicJWKSet() string
	GetJWKSURI() string
}

type esiAuth struct {
	oauthConfig *oauth2.Config
	jwks        jwk.Set
}

type userAuth struct {
	rsaKey        *rsa.PrivateKey
	privateJWKSet jwk.Set
	publicJWKSet  jwk.Set
	tokenKid      uuid.UUID
	tokenExpiry   time.Duration
	tokenDomain   string
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
	tokenKid uuid.UUID,
	tokenDomain string,
	tokenExpiry time.Duration,
	esiAuthJWKSEndpoint *url.URL,
) *Service {
	s := &Service{
		env:    env,
		client: client,
		cache:  cache,
		esiAuth: esiAuth{
			oauthConfig: esiOAuth,
		},
		userAuth: userAuth{
			rsaKey:      userPrivateKey,
			tokenKid:    tokenKid,
			tokenExpiry: tokenExpiry,
			tokenDomain: tokenDomain,
		},
	}

	err := s.initializeESIJWKSet(esiAuthJWKSEndpoint)
	if err != nil {
		panic(fmt.Errorf("internal.Auth: %w", err))
	}

	err = s.initializeJWTSet()
	if err != nil {
		panic(fmt.Errorf("internal.Auth: %w", err))
	}

	return s
}

func (s *Service) GetJWKSURI() string {
	return fmt.Sprintf("%s/.well-known/jwks", s.userAuth.tokenDomain)
}

func (s *Service) GetPublicJWKSet() string {
	data, _ := json.Marshal(s.userAuth.publicJWKSet)
	return string(data)
}

func (s *Service) initializeJWTSet() error {

	privateKey, err := jwk.New(s.userAuth.rsaKey)
	if err != nil {
		return errors.Wrap(err, "failed to initalize jwk using provided rsa private key")
	}

	if _, ok := privateKey.(jwk.RSAPrivateKey); !ok {
		return errors.Wrapf(err, "unexpected type for jwk. expected jwk.RSAPrivateKey, got %T", privateKey)
	}

	err = privateKey.Set(jwk.AlgorithmKey, jwa.RS256)
	if err != nil {
		return errors.Wrap(err, "failed to set alg key on jwk")
	}

	err = privateKey.Set(jwk.KeyIDKey, s.userAuth.tokenKid.String())
	if err != nil {
		return errors.Wrap(err, "failed to set kid on jwk")
	}

	privateSet := jwk.NewSet()
	privateSet.Add(privateKey.(jwk.RSAPrivateKey))

	s.userAuth.privateJWKSet = privateSet

	publicKey, err := jwk.New(s.userAuth.rsaKey.PublicKey)
	if err != nil {
		return errors.Wrap(err, "failed to initalize jwk using provided rsa public key")
	}

	if _, ok := publicKey.(jwk.RSAPublicKey); !ok {
		return errors.Wrapf(err, "unexpected type for jwk. expected jwk.RSAPublicKey, got %T", publicKey)
	}

	err = publicKey.Set(jwk.AlgorithmKey, jwa.RS256)
	if err != nil {
		return errors.Wrap(err, "failed to set alg key on jwk")
	}

	err = publicKey.Set(jwk.KeyIDKey, s.userAuth.tokenKid.String())
	if err != nil {
		return errors.Wrap(err, "failed to set kid on jwk")
	}

	publicSet := jwk.NewSet()
	publicSet.Add(publicKey.(jwk.RSAPublicKey))

	s.userAuth.publicJWKSet = publicSet

	return nil
}
