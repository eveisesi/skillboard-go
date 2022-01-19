package main

import (
	"github.com/eveisesi/skillz"
	"github.com/eveisesi/skillz/internal/auth"
	"github.com/eveisesi/skillz/internal/cache"
	"github.com/eveisesi/skillz/internal/mysql"
	"github.com/eveisesi/skillz/internal/user"
	"github.com/eveisesi/skillz/internal/web"
	"github.com/urfave/cli/v2"
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

	userRepo := mysql.NewUserRepository(mysqlClient)

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

	user := user.New(redisClient, logger, cache, auth, nil, nil, nil, nil, userRepo)

	return web.NewService(
		env,
		cfg.SessionName,
		logger,
		auth,
		user,
		renderer(),
	).Start()

}
