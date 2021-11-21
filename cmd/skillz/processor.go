package main

// func processorCommand(c *cli.Context) error {

// 	etagRepo := mysql.NewETagRepository(mysqlClient)

// 	cache := cache.New(redisClient)
// 	etag := etag.New(cache, etagRepo)
// 	esi := esi.New(httpClient(), redisClient, etag)

// 	allianceRepo := mysql.NewAllianceRepository(mysqlClient)
// 	corporationRepo := mysql.NewCorporationRepository(mysqlClient)
// 	characterRepo := mysql.NewCharacterRepository(mysqlClient)
// 	cloneRepo := mysql.NewCloneRepository(mysqlClient)
// 	skillsRepo := mysql.NewSkillRepository(mysqlClient)
// 	userRepo := mysql.NewUserRepository(mysqlClient)

// 	alliance := alliance.New(cache, esi, etag, allianceRepo)
// 	corporation := corporation.New(cache, esi, etag, corporationRepo)
// 	character := character.New(cache, esi, etag, characterRepo)

// 	auth := auth.New(httpClient(), oauth2Config(), cache, cfg.Eve.JWKSURI)

// 	user := user.New(redisClient, auth, alliance, character, corporation, userRepo)
// 	clone := clone.New(cache, etag, esi, cloneRepo)
// 	skills := skill.New(cache, esi, skillsRepo)

// 	return processor.New(logger, redisClient, user, skillz.ScopeProcessors{
// 		clone,
// 		skills,
// 	}).Run()

// }
