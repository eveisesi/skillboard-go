package main

import (
	"time"

	"github.com/newrelic/go-agent/v3/newrelic"
)

func buildNewRelic() {

	var err error
	nr, err = newrelic.NewApplication(
		newrelic.ConfigFromEnvironment(),
	)
	if err != nil {
		logger.WithError(err).Fatal("failed to initialize NewRelic Application")
	}

	err = nr.WaitForConnection(20 * time.Second)
	if err != nil {
		logger.WithError(err).Fatal("failed to initialize NewRelic Application Connection")
	}

}
