package main

import (
	"github.com/eveisesi/skillz"
	"github.com/eveisesi/skillz/internal/alliance"
	"github.com/eveisesi/skillz/internal/auth"
	"github.com/eveisesi/skillz/internal/cache"
	"github.com/eveisesi/skillz/internal/character"
	"github.com/eveisesi/skillz/internal/clone"
	"github.com/eveisesi/skillz/internal/contact"
	"github.com/eveisesi/skillz/internal/corporation"
	"github.com/eveisesi/skillz/internal/esi"
	"github.com/eveisesi/skillz/internal/etag"
	"github.com/eveisesi/skillz/internal/mysql"
	"github.com/eveisesi/skillz/internal/processor"
	"github.com/eveisesi/skillz/internal/skill"
	"github.com/eveisesi/skillz/internal/universe"
	"github.com/eveisesi/skillz/internal/user"
	"github.com/urfave/cli/v2"
)

func processorCommand(c *cli.Context) error {

	etagRepo := mysql.NewETagRepository(mysqlClient)

	cache := cache.New(redisClient)
	etag := etag.New(cache, etagRepo)
	esi := esi.New(httpClient(), redisClient, logger, etag)

	allianceRepo := mysql.NewAllianceRepository(mysqlClient)
	corporationRepo := mysql.NewCorporationRepository(mysqlClient)
	characterRepo := mysql.NewCharacterRepository(mysqlClient)
	contactRepo := mysql.NewContactRepository(mysqlClient)
	cloneRepo := mysql.NewCloneRepository(mysqlClient)
	skillsRepo := mysql.NewSkillRepository(mysqlClient)
	userRepo := mysql.NewUserRepository(mysqlClient)
	universeRepo := mysql.NewUniverseRepository(mysqlClient)

	alliance := alliance.New(cache, esi, etag, allianceRepo)
	corporation := corporation.New(cache, esi, etag, corporationRepo)
	character := character.New(cache, esi, etag, characterRepo)

	auth := auth.New(
		httpClient(),
		cache,
		oauth2Config(),
		keyConfig(),
		cfg.Eve.JWKSURI,
		cfg.Auth.TokenExpiry,
	)
	universe := universe.New(cache, esi, universeRepo)
	clone := clone.New(logger, cache, etag, esi, universe, cloneRepo)
	skills := skill.New(logger, cache, esi, universe, skillsRepo)
	contact := contact.New(logger, cache, etag, esi, character, corporation, alliance, contactRepo)

	user := user.New(redisClient, logger, cache, auth, alliance, character, corporation, skills, userRepo)

	return processor.New(logger, redisClient, user, skillz.ScopeProcessors{
		clone,
		skills,
		contact,
	}).Run()

}
