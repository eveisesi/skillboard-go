SHELL := /bin/bash

deps:
	source .env && docker compose up -d redis mysql

processor:
	go mod tidy
	go run ./cmd/skillz/*.go processor

buffalo:
	go mod tidy
	go run ./cmd/skillz/*.go buffalo

test:
	go mod tidy
	go run ./cmd/skillz/*.go test

dup:
	source .env && docker compose up -d

dpull:
	source .env && docker compose pull

ddown:
	source .env && docker compose down

ddownv:
	source .env && docker compose down -v

dlogsf:
	source .env && docker compose logs -f server cron processor