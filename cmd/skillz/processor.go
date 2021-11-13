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
	"github.com/eveisesi/skillz/internal/processor"
	"github.com/eveisesi/skillz/internal/user"
	"github.com/urfave/cli/v2"
)

func processorCommand(c *cli.Context) error {

	etagRepo := mysql.NewETagRepository(mysqlClient)

	cache := cache.New(redisClient)
	etag := etag.New(cache, etagRepo)
	esi := esi.New(httpClient(), redisClient, etag)

	allianceRepo := mysql.NewAllianceRepository(mysqlClient)
	corporationRepo := mysql.NewCorporationRepository(mysqlClient)
	characterRepo := mysql.NewCharacterRepository(mysqlClient)
	cloneRepo := mysql.NewCloneRepository(mysqlClient)
	userRepo := mysql.NewUserRepository(mysqlClient)

	alliance := alliance.New(cache, esi, etag, allianceRepo)
	corporation := corporation.New(cache, esi, etag, corporationRepo)
	character := character.New(cache, esi, etag, characterRepo)

	auth := auth.New(httpClient(), oauth2Config(), cache, cfg.Eve.JWKSURI)

	user := user.New(redisClient, auth, alliance, character, corporation, userRepo)
	clone := clone.New(cache, etag, esi, cloneRepo)

	return processor.New(logger, redisClient, user, skillz.ScopeProcessors{
		clone,
	}).Run()

}
