SHELL := /bin/bash

deps:
	aws-vault exec skillboard -- chamber exec skillboard/development -- docker compose up -d redis mysql

processor:
	go mod tidy
	aws-vault exec skillboard -- chamber exec skillboard/development -- go run ./cmd/skillz/*.go processor

server:
	go mod tidy
	aws-vault exec skillboard -- chamber exec skillboard/development -- go run ./cmd/skillz/*.go server

test:
	go mod tidy
	aws-vault exec skillboard -- chamber exec skillboard/development -- go run ./cmd/skillz/*.go test

dup:
	aws-vault exec skillboard -- chamber exec skillboard/production -- docker compose up -d

ddown:
	aws-vault exec skillboard -- chamber exec skillboard/production -- docker compose down

ddownv:
	aws-vault exec skillboard -- chamber exec skillboard/production -- docker compose down -v

dlogsf:
	aws-vault exec skillboard -- chamber exec skillboard/production -- docker compose logs -f