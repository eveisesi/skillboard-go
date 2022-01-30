package main

import (
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
	"github.com/eveisesi/skillz/internal/skill"
	"github.com/eveisesi/skillz/internal/universe"
	"github.com/eveisesi/skillz/internal/user/v2"
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

func buffaloCmd(c *cli.Context) error {

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

	var env skillz.Environment = skillz.Development
	if cfg.Environment == "production" {
		env = skillz.Production
	}

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

	return web.NewService(
		env,
		cfg.SessionName,
		logger,
		auth,
		user,
		renderer(),
	).Start()

}
