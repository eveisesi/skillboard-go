package main

import (
	"github.com/eveisesi/skillz"
	"github.com/eveisesi/skillz/internal/auth"
	"github.com/eveisesi/skillz/internal/cache"
	"github.com/eveisesi/skillz/internal/user/v2"
	"github.com/eveisesi/skillz/internal/web"
	"github.com/urfave/cli/v2"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func init() {
	commands = append(commands, &cli.Command{
		Name:        "buffalo",
		Description: "Starts a Buffalo Web App to Serve content",
		Action:      buffaloCmd,
	})
}

func buffaloCmd(_ *cli.Context) error {

	var env skillz.Environment = skillz.Development
	if cfg.Environment == "production" {
		env = skillz.Production
	}

	boil.SetDB(dbConn)

	cache := cache.New(redisClient)

	auth := auth.New(
		skillz.EnvironmentFromString(cfg.Environment),
		httpClient(),
		cache,
		oauth2Config(),
		keyConfig(),
		cfg.Auth.TokenKID,
		cfg.Auth.TokenDomain,
		cfg.Auth.TokenExpiry,
		cfg.Eve.JWKSURI,
	)

	user := user.NewService(redisClient, logger, cache, auth)

	web.NewService(
		env,
		cfg.SessionName,
		logger,
		auth,
		user,
		renderer(),
	).Start()

	return nil

}
