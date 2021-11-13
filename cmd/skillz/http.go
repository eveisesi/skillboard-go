package main

import (
	"net/http"
	"time"

	"github.com/eveisesi/skillz/pkg/roundtripper"
)

func httpClient() *http.Client {
	return &http.Client{
		Transport: roundtripper.UserAgent(cfg.UserAgent, http.DefaultTransport),
		Timeout:   time.Second * 5,
	}
}
