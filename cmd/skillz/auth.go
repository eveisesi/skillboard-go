package main

import (
	"crypto/x509"
	"encoding/pem"

	"github.com/eveisesi/skillz"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwk"
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
		RedirectURL: "http://localhost:54400/login",
	}
}

func keyConfig() jwk.RSAPrivateKey {

	pemData, _ := pem.Decode(cfg.Auth.PrivateKey)
	if pemData == nil {
		logger.Fatal("pem.Decode of Private Key failed")
	}

	if pemData.Type != "RSA PRIVATE KEY" {
		logger.Fatalf("Expected RSA PRIVATE KEY, got %s", pemData.Type)
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(pemData.Bytes)
	if err != nil {
		logger.WithError(err).Fatal("failed to parse decode pem")
	}

	privateJWK := jwk.NewRSAPrivateKey()
	err = privateJWK.FromRaw(privateKey)
	if err != nil {
		logger.WithError(err).Fatal("failed to initialize jwk with private key")
	}

	err = privateJWK.Set(jwk.AlgorithmKey, jwa.RS256)
	if err != nil {
		logger.WithError(err).Fatal("failed to set JWK Algo")
	}

	err = privateJWK.Set(jwk.KeyIDKey, cfg.Auth.KeyID)
	if err != nil {
		logger.WithError(err).Fatal("failed to set JWK KeyID")
	}

	return privateJWK

}
