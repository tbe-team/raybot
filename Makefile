########################
# Code generation
########################
.PHONY: gen-openapi
gen-openapi:
	set -eux

	npx --yes @redocly/cli bundle ./api/openapi/openapi.yml --output api/openapi/gen/openapi.yml --ext yml
	go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@v2.4.1 \
		-config internal/controller/http/oas/gen/oapi-codegen.yml \
		api/openapi/gen/openapi.yml

.PHONY: gen-sqlc
gen-sqlc:
	go run github.com/sqlc-dev/sqlc/cmd/sqlc@v1.28.0 generate --file internal/db/sqlc/sqlc.yml

.PHONY: gen-mock
gen-mock:
	go run github.com/vektra/mockery/v2@v2.53.1 --config .mockery.yml

.PHONY: gen-all
gen-all: gen-openapi gen-mock gen-sqlc

#########################
# Database
#########################
GOOSE_DRIVER=sqlite3
GOOSE_DBSTRING="file:./.raybot/data/raybot.db"
GOOSE_MIGRATION_DIR=internal/db/migration

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
	go build -o bin/raybot cmd/raybot/main.go

.PHONY: build-ui
build-ui:
	make -C ui build

.PHONY: build-arm64
build-arm64:
	GOOS=linux GOARCH=arm64 go build -o bin/raybot-arm64 cmd/raybot/main.go

.PHONY: build-windows
build-windows:
	GOOS=windows GOARCH=amd64 go build -o bin/raybot.exe cmd/raybot/main.go

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
