package main

import (
	"context"

	"github.com/eveisesi/skillz/internal/cache"
	"github.com/eveisesi/skillz/internal/esi"
	"github.com/eveisesi/skillz/internal/etag"
	"github.com/eveisesi/skillz/internal/mysql"
	"github.com/urfave/cli/v2"
)

func init() {
	commands = append(
		commands,
		&cli.Command{
			Name:        "import",
			Description: "Run an types importer",
			Action:      importCmd,
		},
	)
}

func importCmd(c *cli.Context) error {

	universeRepo := mysql.NewUniverseRepository(mysqlClient)
	cache := cache.New(redisClient)
	etagRepo := mysql.NewETagRepository(mysqlClient)

	etag := etag.New(cache, etagRepo)
	esi := esi.New(httpClient(), redisClient, logger, etag)

	var ctx = context.Background()
	for _, categoryID := range []uint{6, 16} {

		entry := logger.WithField("categoryID", categoryID)

		skillCategory, err := esi.GetCategory(ctx, categoryID)
		if err != nil {
			entry.WithError(err).Fatal("failed to fetch category")
		}

		err = universeRepo.CreateCategory(ctx, skillCategory)
		if err != nil {
			entry.WithError(err).Fatal("failed to create category in data store")
		}

		logger.Info("successfully processed category")

		for _, groupID := range skillCategory.Groups {

			entry := entry.WithField("groupID", groupID)

			group, err := esi.GetGroup(ctx, groupID)
			if err != nil {
				entry.WithError(err).Fatal("failed to fetch group from esi")
			}

			err = universeRepo.CreateGroup(ctx, group)
			if err != nil {
				entry.WithError(err).Fatal("failed to create group in data store")
			}

			logger.Info("successfully processed group")

			for _, typeID := range group.TypeIDs {

				entry := entry.WithField("typeID", typeID)

				item, err := esi.GetType(ctx, typeID)
				if err != nil {
					entry.WithError(err).Fatal("failed to fetch type from esi")
				}

				err = universeRepo.CreateType(ctx, item)
				if err != nil {
					entry.WithError(err).Fatal("failed to create type in data store")
				}

				if len(item.Attributes) > 0 {
					err = universeRepo.CreateTypeDogmaAttributes(ctx, item.Attributes)
					if err != nil {
						entry.WithError(err).Fatal("failed to create type attributes in data store")
					}
				}

				logger.Info("successfully processed type")

			}

		}
	}

	return nil

}

// func importMap(_ *cli.Context) error {

// 	universeRepo := mysql.NewUniverseRepository(mysqlClient)
// 	cache := cache.New(redisClient)

// 	etagRepo := mysql.NewETagRepository(mysqlClient)

// 	etag := etag.New(cache, etagRepo)

// 	esi := esi.New(httpClient(), redisClient, logger, etag)

// 	var ctx = context.Background()

// 	regionIDs, err := esi.GetRegions(ctx)
// 	if err != nil {
// 		logger.WithError(err).Fatal("failed to fetch regions from ESI")
// 	}

// 	for _, regionID := range regionIDs {

// 		entry := logger.WithField("regionID", regionID)

// 		region, err := esi.GetRegion(ctx, regionID)
// 		if err != nil {
// 			entry.WithError(err).Fatal("failed to fetch region from ESI")
// 		}

// 		err = universeRepo.CreateRegion(ctx, region)
// 		if err != nil {
// 			entry.WithError(err).Fatal("failed to create region in data store")
// 		}

// 		entry.Info("successfully processed region")

// 		for _, constellationID := range region.ConstellationIDs {
// 			entry := entry.WithField("constellationID", constellationID)

// 			constellation, err := esi.GetConstellation(ctx, constellationID)
// 			if err != nil {
// 				entry.WithError(err).Fatal("failed to fetch constellation from ESI")
// 			}

// 			err = universeRepo.CreateConstellation(ctx, constellation)
// 			if err != nil {
// 				entry.WithError(err).Fatal("failed to create constellation in data store")
// 			}

// 			entry.Info("successfully processed constellation")

// 			for _, systemID := range constellation.SystemIDs {
// 				entry := entry.WithField("solarSystemID", systemID)

// 				solarSystem, err := esi.GetSolarSystem(ctx, systemID)
// 				if err != nil {
// 					entry.WithError(err).Fatal("failed to fetch solarSystem from ESI")
// 				}

// 				err = universeRepo.CreateSolarSystem(ctx, solarSystem)
// 				if err != nil {
// 					entry.WithError(err).Fatal("failed to create solarSystem in data store")
// 				}

// 				entry.Info("successfully processed solar system")

// 			}

// 		}
// 	}

// 	return nil

// }
