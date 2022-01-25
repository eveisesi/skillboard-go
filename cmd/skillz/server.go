package main

// "github.com/eveisesi/skillz/internal/graphql/engine/dataloaders"

// func init() {
// 	commands = append(commands, &cli.Command{
// 		Name:        "server",
// 		Description: "Starts the GraphQL API",
// 		Action:      serverCommand,
// 	})
// }

// func serverCommand(_ *cli.Context) error {

// 	var env skillz.Environment = skillz.Development
// 	if cfg.Environment == "production" {
// 		env = skillz.Production
// 	}

// 	cache := cache.New(redisClient)

// 	allianceRepo := mysql.NewAllianceRepository(mysqlClient)
// 	characterRepo := mysql.NewCharacterRepository(mysqlClient)
// 	clonesRepo := mysql.NewCloneRepository(mysqlClient)
// 	contactRepo := mysql.NewContactRepository(mysqlClient)
// 	corporationRepo := mysql.NewCorporationRepository(mysqlClient)
// 	etagRepo := mysql.NewETagRepository(mysqlClient)
// 	skillsRepo := mysql.NewSkillRepository(mysqlClient)
// 	universeRepo := mysql.NewUniverseRepository(mysqlClient)
// 	userRepo := mysql.NewUserRepository(mysqlClient)

// 	etag := etag.New(cache, etagRepo)
// 	esi := esi.New(httpClient(), redisClient, logger, etag)
// 	auth := auth.New(
// 		skillz.EnvironmentFromString(cfg.Environment),
// 		httpClient(),
// 		cache,
// 		oauth2Config(),
// 		keyConfig(),
// 		cfg.Auth.TokenKID,
// 		cfg.Auth.TokenDomain,
// 		cfg.Auth.TokenExpiry,
// 		cfg.Eve.JWKSURI,
// 	)
// 	character := character.New(logger, cache, esi, etag, characterRepo)
// 	corporation := corporation.New(logger, cache, esi, etag, corporationRepo)
// 	alliance := alliance.New(logger, cache, esi, etag, allianceRepo)
// 	universe := universe.New(logger, cache, esi, universeRepo)
// 	clone := clone.New(logger, cache, etag, esi, universe, clonesRepo)
// 	contact := contact.New(logger, cache, etag, esi, character, corporation, alliance, contactRepo)
// 	skill := skill.New(logger, cache, esi, universe, skillsRepo)
// 	user := user.New(redisClient, logger, cache, auth, alliance, character, corporation, skill, userRepo)
// 	server := server.New(cfg.Server.Port, env, logger, alliance, auth, character, clone, contact, corporation, skill, user)
// 	errChan := make(chan error, 1)
// 	go func() {
// 		errChan <- server.Run()
// 	}()

// 	sc := make(chan os.Signal, 1)
// 	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

// 	select {
// 	case err := <-errChan:
// 		logger.WithError(err).Error("server error encountered, shutting down")

// 		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)

// 		err = server.Shutdown(ctx)
// 		if err != nil {
// 			logger.WithError(err).Fatal("failed to gracefully shutdown server")
// 		}

// 		cancel()
// 	case <-sc:
// 		logger.Info("gracefully shutting down server")

// 		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)

// 		err := server.Shutdown(ctx)
// 		if err != nil {
// 			logger.WithError(err).Fatal("failed to gracefully shutdown server")
// 		}

// 		cancel()

// 	}

// 	return nil
// }
