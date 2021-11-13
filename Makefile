SHELL := /bin/bash

gqlgen:
	gqlgen generate --config .config/gql/gqlgen.yml

processor:
	go mod tidy
	go run ./cmd/skillz/*.go processor