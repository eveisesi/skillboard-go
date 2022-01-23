package main

import (
	"net/url"
	"time"

	"github.com/gofrs/uuid"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

type config struct {
	MySQL struct {
		Host          string `envconfig:"MYSQL_HOST" required:"true"`
		User          string `envconfig:"MYSQL_USER" required:"true"`
		Pass          string `envconfig:"MYSQL_PASSWORD" required:"true"`
		DB            string `envconfig:"MYSQL_DATABASE" required:"true"`
		RunMigrations uint   `envconfig:"RUN_MIGRATIONS" required:"true"`
	}

	Redis struct {
		Host string `envconfig:"REDIS_HOST" required:"true"`
		Pass string `envconfig:"REDIS_PASS" required:"true"`
	}

	Log struct {
		Level string `envconfig:"LOG_LEVEL" required:"true"`
	}

	Eve struct {
		ClientID           string `envconfig:"EVE_CLIENT_ID" required:"true"`
		ClientSecret       string `envconfig:"EVE_CLIENT_SECRET" required:"true"`
		CallbackURIStr     string `envconfig:"EVE_CALLBACK_URI" required:"true"`
		CallbackURI        *url.URL
		JWKSURIStr         string `envconfig:"EVE_JWKS_URI" required:"true"`
		JWKSURI            *url.URL
		InitializeUniverse uint `envconfig:"INITIALIZE_UNIVERSE" required:"true"`
	}

	Auth struct {
		TokenKIDStr    string `envconfig:"AUTH_TOKEN_KID" required:"true"`
		TokenKID       uuid.UUID
		TokenExpiryStr string `envconfig:"AUTH_EXPIRY" required:"true"`
		TokenExpiry    time.Duration
		TokenDomain    string `envconfig:"AUTH_TOKEN_DOMAIN" required:"true"`
	}

	Server struct {
		Port uint `envconfig:"SERVER_PORT" required:"true"`
	}

	SessionName string `envconfig:"SESSION_NAME" default:"__skillboard_session"`

	Environment string `envconfig:"ENVIRONMENT" required:"true"`

	UserAgent string `envconfig:"USER_AGENT" required:"true"`
}

func buildConfig() {
	_ = godotenv.Load(".config/.env")

	cfg = new(config)
	err := envconfig.Process("", cfg)
	if err != nil {
		panic(errors.Wrap(err, "failed to config env"))
	}

	dur, err := time.ParseDuration(cfg.Auth.TokenExpiryStr)
	if err != nil {
		panic(errors.Wrap(err, "failed to parse AUTH_EXPIRY"))
	}

	cfg.Auth.TokenExpiry = dur

	cfg.Auth.TokenKID = uuid.Must(uuid.FromString(cfg.Auth.TokenKIDStr))

	callbackURI, err := url.Parse(cfg.Eve.CallbackURIStr)
	if err != nil {
		panic(errors.Wrap(err, "failed to parse EVE_CALLBACK_URI as a valid URI"))
	}

	cfg.Eve.CallbackURI = callbackURI

	jwksURI, err := url.Parse(cfg.Eve.JWKSURIStr)
	if err != nil {
		panic(errors.Wrap(err, "failed to parse EVE_JWKS_URI as a valid URI"))
	}

	cfg.Eve.JWKSURI = jwksURI
}
