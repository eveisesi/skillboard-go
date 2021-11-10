package main

import (
	"os"
	"runtime"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var (
	cfg    *config
	logger *logrus.Logger
	app    *cli.App
)

func init() {
	app = cli.NewApp()
	app.Name = "Skillboard CLI"
	app.HelpName = "skillboard"
	app.Usage = "A CLI for Eve Online Skillboards"
	app.UsageText = "skillboard --help"
	app.Commands = []*cli.Command{
		{
			Name:   "server",
			Action: serverCommand,
		},
	}

	buildConfig()
	buildLogger()
	runtime.GC()
}

func serverCommand(c *cli.Context) error {
	logger.Info("ping")
	return nil
}

func main() {

	err := app.Run(os.Args)
	if err != nil {
		logger.WithError(err).Fatal("failed to initialize CLI")
	}

}
