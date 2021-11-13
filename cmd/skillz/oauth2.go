package main

import (
	"github.com/eveisesi/skillz"
	"golang.org/x/oauth2"
)

func oauth2Config() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     cfg.Eve.ClientID,
		ClientSecret: cfg.Eve.ClientSecret,
		Scopes: []string{
			skillz.ReadSkillsV1.String(),
			skillz.ReadSkillQueueV1.String(),
			skillz.ReadClonesV1.String(),
			// "esi-universe.read_structures.v1",
			// "esi-characters.read_standings.v1",
			skillz.ReadClonesV1.String(),
		},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://login.eveonline.com/v2/oauth/authorize",
			TokenURL: "https://login.eveonline.com/v2/oauth/token",
		},
		RedirectURL: "http://localhost:54400/auth/callback",
	}
}
