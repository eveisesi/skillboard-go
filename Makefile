SHELL := /bin/bash

deps:
	aws-vault exec phoenix -- chamber exec phoenix -- docker compose up -d redis mysql

processor:
	go mod tidy
	aws-vault exec phoenix -- chamber exec phoenix -- go run ./cmd/skillz/*.go processor

server:
	go mod tidy
	aws-vault exec phoenix -- chamber exec phoenix -- go run ./cmd/skillz/*.go server

test:
	go mod tidy
	aws-vault exec phoenix -- chamber exec phoenix -- go run ./cmd/skillz/*.go test