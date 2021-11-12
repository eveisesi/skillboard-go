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
