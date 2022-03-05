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
	"github.com/eveisesi/skillz/internal/mysql"
	"github.com/eveisesi/skillz/internal/server"
	"github.com/eveisesi/skillz/internal/skill"
	"github.com/eveisesi/skillz/internal/universe"
	"github.com/eveisesi/skillz/internal/user/v2"
	"github.com/urfave/cli/v2"
)

func init() {
	commands = append(commands, &cli.Command{
		Name:        "server",
		Description: "Starts the HTTP Server",
		Action:      serverCmd,
	})
}

func serverCmd(c *cli.Context) error {

	if cfg.MySQL.RunMigrations == 1 {
		err := migrateUpCommand(c)
		if err != nil {
			return err
		}
	}

	if cfg.Eve.InitializeUniverse == 1 {
		err := importCmd(c)
		if err != nil {
			return err
		}
	}

	// var env skillz.Environment = skillz.Development
	// if cfg.Environment == "production" {
	// 	env = skillz.Production
	// }

	allianceRepo := mysql.NewAllianceRepository(mysqlClient)
	characterRepo := mysql.NewCharacterRepository(mysqlClient)
	corporationRepo := mysql.NewCorporationRepository(mysqlClient)
	etagRepo := mysql.NewETagRepository(mysqlClient)
	cloneRepo := mysql.NewCloneRepository(mysqlClient)
	skillzRepo := mysql.NewSkillRepository(mysqlClient)
	userRepo := mysql.NewUserRepository(mysqlClient)
	universeRepo := mysql.NewUniverseRepository(mysqlClient)

	cache := cache.New(redisClient, cfg.Redis.DisableCache == 1)
	etag := etag.New(cache, etagRepo)
	esi := esi.New(httpClient(), redisClient, logger, etag)
	character := character.New(logger, cache, esi, etag, characterRepo)
	corporation := corporation.New(logger, cache, esi, etag, corporationRepo)
	alliance := alliance.New(logger, cache, esi, etag, allianceRepo)
	universe := universe.New(logger, cache, esi, universeRepo)
	clone := clone.New(logger, cache, etag, esi, universe, cloneRepo)
	skills := skill.New(logger, cache, esi, universe, skillzRepo)

	auth := auth.New(
		skillz.EnvironmentFromString(cfg.Environment),
		httpClient(),
		cache,
		oauth2Config(),
	)

	user := user.New(redisClient, logger, cache, auth, alliance, character, corporation, skills, clone, userRepo)

	srv := server.New(logger, nr, auth, user)

	go func() {
		if err := srv.Start(); err != nil {
			logger.Fatal("failed to start http server")
		}
	}()

	block := make(chan os.Signal, 1)
	signal.Notify(block, os.Interrupt, syscall.SIGTERM)

	<-block

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	srv.GracefullyShutdown(ctx)
	logger.Info("Gracefully shutting down")
	return nil
}
