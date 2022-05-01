package main

import (
	"net/url"

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
		Host         string `envconfig:"REDIS_HOST" required:"true"`
		Pass         string `envconfig:"REDIS_PASS" required:"true"`
		DisableCache uint   `envconfig:"DISABLE_CACHE" required:"true"`
	}

	Log struct {
		Level string `envconfig:"LOG_LEVEL" required:"true"`
	}

	Eve struct {
		ClientID           string `envconfig:"EVE_CLIENT_ID" required:"true"`
		ClientSecret       string `envconfig:"EVE_CLIENT_SECRET" required:"true"`
		CallbackURIStr     string `envconfig:"EVE_CALLBACK_URI" required:"true"`
		CallbackURI        *url.URL
		InitializeUniverse uint `envconfig:"INITIALIZE_UNIVERSE" required:"true"`
	}

	SessionName string `envconfig:"SESSION_NAME" default:"__skillboard_session"`

	Environment string `envconfig:"ENVIRONMENT" required:"true"`

	UserAgent string `envconfig:"USER_AGENT" required:"true"`
}

func buildConfig() {
	_ = godotenv.Load("app.env")

	cfg = new(config)
	err := envconfig.Process("", cfg)
	if err != nil {
		panic(errors.Wrap(err, "failed to config env"))
	}

	callbackURI, err := url.Parse(cfg.Eve.CallbackURIStr)
	if err != nil {
		panic(errors.Wrap(err, "failed to parse EVE_CALLBACK_URI as a valid URI"))
	}

	cfg.Eve.CallbackURI = callbackURI

}
