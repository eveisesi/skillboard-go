package main

import (
	"net/url"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

type config struct {
	MySQL struct {
		Host string `envconfig:"MYSQL_HOST" required:"true"`
		User string `envconfig:"MYSQL_USER" required:"true"`
		Pass string `envconfig:"MYSQL_PASSWORD" required:"true"`
		DB   string `envconfig:"MYSQL_DATABASE" required:"true"`
	}

	Redis struct {
		Host string `envconfig:"REDIS_HOST" required:"true"`
		Pass string `envconfig:"REDIS_PASS" required:"true"`
	}

	Log struct {
		Level string `envconfig:"LOG_LEVEL" required:"true"`
	}

	Eve struct {
		ClientID       string `envconfig:"EVE_CLIENT_ID" required:"true"`
		ClientSecret   string `envconfig:"EVE_CLIENT_SECRET" required:"true"`
		CallbackURIStr string `envconfig:"EVE_CALLBACK_URI" required:"true"`
		CallbackURI    *url.URL
		JWKSURIStr     string `envconfig:"EVE_JWKS_URI" required:"true"`
		JWKSURI        *url.URL
	}

	Auth struct {
		PrivateKey      []byte `envconfig:"AUTH_KEY"`
		CookieExpiryStr string `envconfig:"AUTH_EXPIRY" required:"true"`
		CookieExpiry    time.Duration
		CookieURI       string `envconfig:"AUTH_COOKIE_URI" required:"true"`
	}

	Server struct {
		Port uint `envconfig:"SERVER_PORT" required:"true"`
	}

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

	if len(cfg.Auth.PrivateKey) == 0 {
		if cfg.Environment == "production" {
			panic(errors.New("failed to config env: AUTH_KEY is required but has a length of 0"))
		}

		cfg.Auth.PrivateKey, err = os.ReadFile("../../.config/.key")
		if err != nil {
			panic(errors.Wrap(err, "failed to config env: AUTH_KEY is required and an error was encountered reading key file"))
		}
	}

	dur, err := time.ParseDuration(cfg.Auth.CookieExpiryStr)
	if err != nil {
		panic(errors.Wrap(err, "failed to parse AUTH_EXPIRY"))
	}

	cfg.Auth.CookieExpiry = dur

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
