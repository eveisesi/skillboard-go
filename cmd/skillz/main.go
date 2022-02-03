package main

import (
	"os"

	"github.com/eveisesi/skillz/internal/mysql"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var (
	cfg         *config
	logger      *logrus.Logger
	app         *cli.App
	redisClient *redis.Client
	dbConn      *sqlx.DB
	mysqlClient mysql.QueryExecContext
	commands    []*cli.Command

	nr *newrelic.Application
)

func main() {

	app = cli.NewApp()
	app.Name = "Skillz CLI"
	app.HelpName = "skillz"
	app.Usage = "CLI for Skillboard.Evie"
	app.UsageText = "skillz --help"
	app.Commands = commands

	buildConfig()
	buildLogger()
	buildMySQL()
	buildRedis()
	buildNewRelic()

	err := app.Run(os.Args)
	if err != nil {
		logger.WithError(err).Fatal("failed to initialize CLI")
	}

}
