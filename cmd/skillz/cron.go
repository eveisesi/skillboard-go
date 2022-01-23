package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/eveisesi/skillz"
	"github.com/eveisesi/skillz/internal/alliance"
	"github.com/eveisesi/skillz/internal/auth"
	"github.com/eveisesi/skillz/internal/cache"
	"github.com/eveisesi/skillz/internal/character"
	"github.com/eveisesi/skillz/internal/corporation"
	"github.com/eveisesi/skillz/internal/esi"
	"github.com/eveisesi/skillz/internal/etag"
	"github.com/eveisesi/skillz/internal/mysql"
	"github.com/eveisesi/skillz/internal/skill"
	"github.com/eveisesi/skillz/internal/universe"
	"github.com/eveisesi/skillz/internal/user"
	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
	"github.com/urfave/cli/v2"
)

func init() {
	commands = append(commands, &cli.Command{
		Name:        "cron",
		Description: "Start the Cron Scheduler",
		Action:      cronCommand,
	})
}

func cronCommand(_ *cli.Context) error {
	cache := cache.New(redisClient)
	allianceRepo := mysql.NewAllianceRepository(mysqlClient)
	characterRepo := mysql.NewCharacterRepository(mysqlClient)
	corporationRepo := mysql.NewCorporationRepository(mysqlClient)
	etagRepo := mysql.NewETagRepository(mysqlClient)
	userRepo := mysql.NewUserRepository(mysqlClient)
	skillsRepo := mysql.NewSkillRepository(mysqlClient)
	universeRepo := mysql.NewUniverseRepository(mysqlClient)

	auth := auth.New(
		skillz.EnvironmentFromString(cfg.Environment),
		httpClient(),
		cache,
		oauth2Config(),
		cfg.Eve.JWKSURI,
	)
	etag := etag.New(cache, etagRepo)
	esi := esi.New(httpClient(), redisClient, logger, etag)
	character := character.New(logger, cache, esi, etag, characterRepo)
	corporation := corporation.New(logger, cache, esi, etag, corporationRepo)
	alliance := alliance.New(logger, cache, esi, etag, allianceRepo)
	universe := universe.New(logger, cache, esi, universeRepo)
	skill := skill.New(logger, cache, esi, universe, skillsRepo)
	user := user.New(redisClient, logger, cache, auth, alliance, character, corporation, skill, userRepo)

	cron := cron.New()

	entryID, err := cron.AddFunc("@every 3h", func() {

		var ctx = context.Background()

		logger.Info("executing process updateable users cron")

		err := user.ProcessUpdatableUsers(ctx)
		if err != nil {
			logger.WithError(err).Fatal("failed to update processable users")
		}
	})
	if err != nil {
		logger.WithError(err).Error("failed to add user update job to cron scheduler")
		return errors.Wrap(err, "failed to add user update job to cron scheduler")
	}

	logger.WithField("entryID", entryID).Info("successfully address cron job to scheduler")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	cron.Start()
	logger.Info("cron running....")
	<-sc
	logger.Info("stopping cron....")

	ctx := cron.Stop()

	logger.Info("cron stopped, waiting for jobs to finish....")
	<-ctx.Done()
	logger.Info("cron stopped successfully, exiting process....")
	return nil
}
