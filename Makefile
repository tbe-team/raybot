########################
# Code generation
########################
.PHONY: gen-openapi
gen-openapi:
	set -eux

	pnpm --package=@redocly/cli@1.34 dlx redocly bundle ./api/openapi/openapi.yml --output api/openapi/gen/openapi.yml --ext yml
	go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@v2.4.1 \
		-config internal/controller/http/oas/gen/oapi-codegen.yml \
		api/openapi/gen/openapi.yml

.PHONY: gen-proto
gen-proto:
	pnpm --package=@bufbuild/buf@1.50.1 dlx buf generate

.PHONY: gen-sqlc
gen-sqlc:
	go run github.com/sqlc-dev/sqlc/cmd/sqlc@v1.28.0 generate --file internal/storage/db/sqlc/sqlc.yml

.PHONY: gen-mock
gen-mock:
	go run github.com/vektra/mockery/v2@v2.53.1 --config .mockery.yml

.PHONY: gen-all
gen-all: gen-openapi gen-proto gen-mock gen-sqlc

#########################
# Database
#########################
GOOSE_DRIVER=sqlite3
GOOSE_DBSTRING="file:./.raybot/data/raybot.db"
GOOSE_MIGRATION_DIR=internal/storage/db/migration

.PHONY: migrate-up
migrate-up:
	GOOSE_DRIVER=$(GOOSE_DRIVER) \
	GOOSE_DBSTRING=$(GOOSE_DBSTRING) \
	GOOSE_MIGRATION_DIR=$(GOOSE_MIGRATION_DIR) \
	go run github.com/pressly/goose/v3/cmd/goose@v3.24.1 up

.PHONY: migrate-down
migrate-down:
	GOOSE_DRIVER=$(GOOSE_DRIVER) \
	GOOSE_DBSTRING=$(GOOSE_DBSTRING) \
	GOOSE_MIGRATION_DIR=$(GOOSE_MIGRATION_DIR) \
	go run github.com/pressly/goose/v3/cmd/goose@v3.24.1 down

.PHONY: migrate-status
migrate-status:
	GOOSE_DRIVER=$(GOOSE_DRIVER) \
	GOOSE_DBSTRING=$(GOOSE_DBSTRING) \
	GOOSE_MIGRATION_DIR=$(GOOSE_MIGRATION_DIR) \
	go run github.com/pressly/goose/v3/cmd/goose@v3.24.1 status

.PHONY: migrate-create
migrate-create:
	GOOSE_DRIVER=$(GOOSE_DRIVER) \
	GOOSE_DBSTRING=$(GOOSE_DBSTRING) \
	GOOSE_MIGRATION_DIR=$(GOOSE_MIGRATION_DIR) \
	go run github.com/pressly/goose/v3/cmd/goose@v3.24.1 create $(name) sql

.PHONY: migrate-reset
migrate-reset:
	GOOSE_DRIVER=$(GOOSE_DRIVER) \
	GOOSE_DBSTRING=$(GOOSE_DBSTRING) \
	GOOSE_MIGRATION_DIR=$(GOOSE_MIGRATION_DIR) \
	go run github.com/pressly/goose/v3/cmd/goose@v3.24.1 reset

#########################
# Build
#########################
.PHONY: build
build:
	CGO_ENABLED=1 \
	GOOS=linux \
	GOARCH=amd64 \
	go build -o bin/raybot cmd/raybot/main.go

.PHONY: build-ui
build-ui:
	make -C ui build

.PHONY: build-arm64
build-arm64:
	CGO_ENABLED=1 \
	GOOS=linux \
	GOARCH=arm64 \
	CC=aarch64-linux-gnu-gcc \
	go build -o bin/raybot-arm64 cmd/raybot/main.go

#########################
# Docker
#########################
.PHONY: docker-build-raybot
docker-build-raybot:
	docker build -t raybot -f docker/raybot.dockerfile .

.PHONY: docker-run-raybot
docker-run-raybot:
	docker run -p 3000:3000 raybot

#########################
# Run
#########################
.PHONY: run
run:
	go run cmd/raybot/main.go

#########################
# Testing
#########################
.PHONY: test
test:
	go test -v ./...

.PHONY: test-cov
test-cov:
	go test -coverprofile=bin/coverage.out ./...
	go tool cover -html=bin/coverage.out -o bin/coverage.html
	@echo "Coverage report saved to bin/coverage.html"

########################
# Lint
########################
.PHONY: lint-go
lint-go:
	golangci-lint run ./... --config .golangci.yml
