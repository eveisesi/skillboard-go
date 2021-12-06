package main

import (
	"os"

	"github.com/eveisesi/skillz/internal/mysql"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var (
	cfg         *config
	logger      *logrus.Logger
	app         *cli.App
	redisClient *redis.Client
	mysqlClient mysql.QueryExecContext
)

func init() {
	app = cli.NewApp()
	app.Name = "Skillz CLI"
	app.HelpName = "skillz"
	app.Usage = "A CLI for Eve Online Skillboards"
	app.UsageText = "skillz --help"
	app.Commands = []*cli.Command{
		{
			Name:        "server",
			Description: "Starts the GraphQL API",
			Action:      serverCommand,
		},
		{
			Name:        "processor",
			Description: "Start the Job Processor",
			Action:      processorCommand,
		},
		{
			Name:        "cron",
			Description: "Start the Cron Scheduler",
			Action:      cronCommand,
		},
		{
			Name:        "importTypes",
			Description: "Run an types importer",
			Action:      importTypes,
			Flags: []cli.Flag{
				&cli.UintFlag{
					Name:  "categoryID",
					Usage: "ID for the Category that we need to import groups and types for",
				},
			},
		},
		{
			Name:        "importMap",
			Description: "Run an map importer",
			Action:      importMap,
		},
		{
			Name:   "test",
			Action: testCommand,
		},
	}

	buildConfig()
	buildLogger()
	buildMySQL()
	buildRedis()
}

func main() {

	err := app.Run(os.Args)
	if err != nil {
		logger.WithError(err).Fatal("failed to initialize CLI")
	}

}
