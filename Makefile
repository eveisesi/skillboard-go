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
	docker compose up -d

dpull:
	docker compose pull

ddown:
	docker compose down

ddownv:
	docker compose down -v

dlogsf:
	docker compose logs -f server cron processor