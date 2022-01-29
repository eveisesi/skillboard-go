SHELL := /bin/bash

deps:
	docker compose up -d redis mysql

processor:
	go mod tidy
	go run ./cmd/skillz/*.go processor

server:
	go mod tidy
	go run ./cmd/skillz/*.go server

buffalo:
	go mod tidy
	go run ./cmd/skillz/*.go buffalo

test:
	go mod tidy
	go run ./cmd/skillz/*.go test

dup:
	aws-vault exec skillboard -- chamber exec skillboard/production -- docker compose up -d

dpull:
	aws-vault exec skillboard -- chamber exec skillboard/production -- docker compose pull

ddown:
	aws-vault exec skillboard -- chamber exec skillboard/production -- docker compose down

ddownv:
	aws-vault exec skillboard -- chamber exec skillboard/production -- docker compose down -v

dlogsf:
	aws-vault exec skillboard -- chamber exec skillboard/production -- docker compose logs -f server cron processor