package main

import (
	"golang.org/x/oauth2"
)

func oauth2Config() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     cfg.Eve.ClientID,
		ClientSecret: cfg.Eve.ClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://login.eveonline.com/v2/oauth/authorize",
			TokenURL: "https://login.eveonline.com/v2/oauth/token",
		},
		RedirectURL: cfg.Eve.CallbackURI.String(),
	}
}
