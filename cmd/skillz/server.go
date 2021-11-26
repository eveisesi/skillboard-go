package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/eveisesi/skillz"
	"github.com/eveisesi/skillz/internal/alliance"
	"github.com/eveisesi/skillz/internal/auth"
	"github.com/eveisesi/skillz/internal/cache"
	"github.com/eveisesi/skillz/internal/character"
	"github.com/eveisesi/skillz/internal/clone"
	"github.com/eveisesi/skillz/internal/corporation"
	"github.com/eveisesi/skillz/internal/esi"
	"github.com/eveisesi/skillz/internal/etag"
	"github.com/eveisesi/skillz/internal/graphql"
	"github.com/eveisesi/skillz/internal/universe"

	// "github.com/eveisesi/skillz/internal/graphql/engine/dataloaders"
	"github.com/eveisesi/skillz/internal/mysql"
	"github.com/eveisesi/skillz/internal/server"
	"github.com/eveisesi/skillz/internal/skill"
	"github.com/eveisesi/skillz/internal/user"
	"github.com/urfave/cli/v2"
)

func serverCommand(_ *cli.Context) error {

	var env skillz.Environment = skillz.Development
	if cfg.Environment == "production" {
		env = skillz.Production
	}

	cache := cache.New(redisClient)

	allianceRepo := mysql.NewAllianceRepository(mysqlClient)
	characterRepo := mysql.NewCharacterRepository(mysqlClient)
	clonesRepo := mysql.NewCloneRepository(mysqlClient)
	corporationRepo := mysql.NewCorporationRepository(mysqlClient)
	etagRepo := mysql.NewETagRepository(mysqlClient)
	skillsRepo := mysql.NewSkillRepository(mysqlClient)
	universeRepo := mysql.NewUniverseRepository(mysqlClient)
	userRepo := mysql.NewUserRepository(mysqlClient)

	etag := etag.New(cache, etagRepo)
	esi := esi.New(httpClient(), redisClient, etag)
	auth := auth.New(
		httpClient(),
		cache,
		oauth2Config(),
		keyConfig(),
		cfg.Eve.JWKSURI,
		cfg.Auth.TokenIssuer,
		cfg.Auth.TokenAudience,
		cfg.Auth.TokenExpiry,
	)
	character := character.New(cache, esi, etag, characterRepo)
	corporation := corporation.New(cache, esi, etag, corporationRepo)
	alliance := alliance.New(cache, esi, etag, allianceRepo)
	user := user.New(redisClient, auth, alliance, character, corporation, userRepo)
	clone := clone.New(cache, etag, esi, clonesRepo)
	universe := universe.New(cache, esi, universeRepo)
	skill := skill.New(cache, esi, universe, skillsRepo)
	graphql := graphql.New(alliance, auth, character, clone, corporation, skill, universe, user)
	server := server.New(cfg.Server.Port, env, logger, cache, graphql)
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
