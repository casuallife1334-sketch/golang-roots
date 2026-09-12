SHELL := /bin/sh

include .env
export

export PROJECT_ROOT := $(shell pwd)

env-up:
	@docker compose up -d postgres minio

env-down:
	@docker compose down

env-cleanup:
	@read -p "Очистить все volumes окружения? ОПАСНОСТЬ УТЕРИ ДАННЫХ. [y/N]: " ans; \
	if [ "$$ans" = "y" ]; then \
		docker compose down -v && \
		rm -rf "$(PROJECT_ROOT)/out/pgdata" "$(PROJECT_ROOT)/out/minio" && \
		echo "Файлы окружения очищены"; \
	else \
		echo "Очистка окружения отменена"; \
	fi

migrate-up:
	@make migrate-action action=up

migrate-down:
	@make migrate-action action=down

migrate-action:
	@if [ -z "$(action)" ]; then \
		echo "Отсутствует необходимый параметр 'action'. Пример: 'make migrate-action action=up'"; \
		exit 1; \
	fi; \
	docker compose run --rm postgres-migrate \
		-path /migrations \
		-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@postgres:5432/${POSTGRES_DB}?sslmode=disable \
		"$(action)"

run:
	@LOGGER_FOLDER="$(PROJECT_ROOT)/out/logs" go run ./cmd/genealogy

test:
	@go test ./...

fmt:
	@gofmt -w ./cmd ./internal
