# Makefile

APP_PORT ?= 8080

.PHONY: run-api
run-api:
	go run ./apps/api

.PHONY: test
test:
	go test ./...

include .env

export

compose = docker compose -f deployments/compose/docker-compose.yml --profile tools

migrate-up:
	$(compose) run --rm migrate -path=/migrations -database=$(DATABASE_URL) up

migrate-down:
	$(compose) run --rm migrate -path=/migrations -database=$(DATABASE_URL) down 1

migrate-force:
	$(compose) run --rm migrate -path=/migrations -database=$(DATABASE_URL) force 1