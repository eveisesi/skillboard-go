package main

import (
	"github.com/eveisesi/skillz"
	"golang.org/x/oauth2"
)

func oauth2Config() *oauth2.Config {

	scopes := make([]string, 0, len(skillz.AllScopes))
	for _, scope := range skillz.AllScopes {
		scopes = append(scopes, scope.String())
	}

	return &oauth2.Config{
		ClientID:     cfg.Eve.ClientID,
		ClientSecret: cfg.Eve.ClientSecret,
		Scopes:       scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://login.eveonline.com/v2/oauth/authorize",
			TokenURL: "https://login.eveonline.com/v2/oauth/token",
		},
		RedirectURL: "http://localhost:54400/auth/callback",
	}
}
