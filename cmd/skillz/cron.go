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
	"github.com/eveisesi/skillz/internal/processor"
	"github.com/eveisesi/skillz/internal/skill"
	"github.com/eveisesi/skillz/internal/universe"
	"github.com/eveisesi/skillz/internal/user/v2"
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
	cache := cache.New(redisClient, cfg.Redis.DisableCache == 1)
	allianceRepo := mysql.NewAllianceRepository(mysqlClient)
	characterRepo := mysql.NewCharacterRepository(mysqlClient)
	corporationRepo := mysql.NewCorporationRepository(mysqlClient)
	etagRepo := mysql.NewETagRepository(mysqlClient)
	userRepo := mysql.NewUserRepository(mysqlClient)
	cloneRepo := mysql.NewCloneRepository(mysqlClient)
	skillsRepo := mysql.NewSkillRepository(mysqlClient)
	universeRepo := mysql.NewUniverseRepository(mysqlClient)

	auth := auth.New(
		skillz.EnvironmentFromString(cfg.Environment),
		httpClient(),
		cache,
		oauth2Config(),
	)
	etag := etag.New(cache, etagRepo)
	esi := esi.New(httpClient(), redisClient, logger, etag)
	character := character.New(logger, cache, esi, etag, characterRepo)
	corporation := corporation.New(logger, cache, esi, etag, corporationRepo)
	alliance := alliance.New(logger, cache, esi, etag, allianceRepo)
	universe := universe.New(logger, cache, esi, universeRepo)
	clone := clone.New(logger, cache, etag, esi, universe, cloneRepo)
	skills := skill.New(logger, cache, esi, universe, skillsRepo)
	user := user.New(redisClient, logger, cache, auth, alliance, character, corporation, skills, clone, userRepo)

	cron := cron.New()

	processor := processor.New(logger, redisClient, nr, user, skillz.ScopeProcessors{
		clone,
		skills,
	})

	entryID, err := cron.AddFunc("0 */3 * * *", func() {

		var ctx = context.Background()

		logger.Info("executing process updateable users cron")

		users, err := user.ProcessUpdatableUsers(ctx)
		if err != nil {
			logger.WithError(err).Fatal("failed to update processable users")
		}

		logger.WithField("count", len(users)).Info("updateable users")

		for _, user := range users {
			err = processor.ProcessUser(ctx, user)
			if err != nil {
				logger.WithError(err).Error("failed to process user id")
				time.Sleep(time.Second * 3)
			}
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
