package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/eveisesi/skillz"
	"github.com/eveisesi/skillz/internal/auth"
	"github.com/eveisesi/skillz/internal/cache"
	"github.com/eveisesi/skillz/internal/character"
	"github.com/eveisesi/skillz/internal/esi"
	"github.com/eveisesi/skillz/internal/etag"
	"github.com/eveisesi/skillz/internal/mysql"
	"github.com/eveisesi/skillz/internal/server"
	"github.com/eveisesi/skillz/internal/user"
	"github.com/eveisesi/skillz/pkg/roundtripper"
	"github.com/urfave/cli/v2"
	"golang.org/x/oauth2"
)

func serverCommand(_ *cli.Context) error {

	var env skillz.Environment = skillz.Development
	if cfg.Environment == "production" {
		env = skillz.Production
	}

	cache := cache.New(redisClient)

	etagRepo := mysql.NewETagRepository(mysqlClient, "etags")
	characterRepo := mysql.NewCharacterRepository(mysqlClient, "characters", "character_corporation_history")
	userRepo := mysql.NewUserRepository(mysqlClient, "users")

	httpClient := &http.Client{
		Transport: roundtripper.UserAgent(cfg.UserAgent, http.DefaultTransport),
		Timeout:   time.Second * 5,
	}

	oauth2Config := &oauth2.Config{
		ClientID:     cfg.Eve.ClientID,
		ClientSecret: cfg.Eve.ClientSecret,
		Scopes: []string{
			"esi-skills.read_skills.v1",
			"esi-skills.read_skillqueue.v1",
			"esi-clones.read_clones.v1",
			"esi-universe.read_structures.v1",
			"esi-characters.read_standings.v1",
			"esi-clones.read_implants.v1",
		},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://login.eveonline.com/v2/oauth/authorize",
			TokenURL: "https://login.eveonline.com/v2/oauth/token",
		},
		RedirectURL: "http://localhost:54400/auth/callback",
	}

	etag := etag.New(cache, etagRepo)
	esi := esi.New(httpClient, redisClient, etag)
	auth := auth.New(httpClient, oauth2Config, cache, cfg.Eve.JWKSURI)
	character := character.New(cache, esi, etag, characterRepo)
	user := user.New(auth, character, userRepo)
	server := server.New(cfg.Server.Port, env, logger, auth, character, user)
	errChan := make(chan error, 1)
	go func() {
		errChan <- server.Run()
	}()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	select {
	case err := <-errChan:
		logger.WithError(err).Error("server error encountered, shutting down")

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)

		err = server.Shutdown(ctx)
		if err != nil {
			logger.WithError(err).Fatal("failed to gracefully shutdown server")
		}

		cancel()
	case <-sc:
		logger.Info("gracefully shutting down server")

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)

		err := server.Shutdown(ctx)
		if err != nil {
			logger.WithError(err).Fatal("failed to gracefully shutdown server")
		}

		cancel()

	}

	return nil
}
