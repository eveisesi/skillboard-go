package main

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type config struct {
	MySQL struct {
		Host string `envconfig:"MYSQL_HOST" required:"true"`
		User string `envconfig:"MYSQL_USER" required:"true"`
		Pass string `envconfig:"MYSQL_PASSWORD" required:"true"`
		DB   string `envconfig:"MYSQL_DB" required:"true"`
	}

	Redis struct {
		Host string `envconfig:"REDIS_HOST" required:"true"`
		Pass string `envconfig:"REDIS_PASS" required:"true"`
	}

	Log struct {
		Level string `envconfig:"LOG_LEVEL" required:"true"`
	}

	Eve struct {
		ClientID     string `envconfig:"EVE_CLIENT_ID" required:"true"`
		ClientSecret string `envconfig:"EVE_CLIENT_SECRET" required:"true"`
		JWKSURI      string `envconfig:"EVE_JWKS_URI" required:"true"`
	}

	Auth struct {
		PrivateKey     []byte `envconfig:"AUTH_KEY"`
		TokenExpiryStr string `envconfig:"AUTH_EXPIRY" required:"true"`
		TokenExpiry    time.Duration
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
		panic(fmt.Sprintf("failed to config env: %s", err))
	}

	// TODO (@ddouglas): Remove ME
	if len(cfg.Auth.PrivateKey) == 0 {
		if cfg.Environment == "production" {
			panic("failed to config env: AUTH_KEY is required but has a length of 0")
		}

		cfg.Auth.PrivateKey, err = os.ReadFile("../../.config/.key")
		if err != nil {
			err = fmt.Errorf("failed to config env: AUTH_KEY is required and an error was encountered reading key file: %w", err)
			fmt.Println(err)
			panic(err)
		}
	}

	dur, err := time.ParseDuration(cfg.Auth.TokenExpiryStr)
	if err != nil {
		panic(fmt.Errorf("failed to parse AUTH_EXPIRY: %w", err))
	}

	cfg.Auth.TokenExpiry = dur

}
