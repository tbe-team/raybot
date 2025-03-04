########################
# Code generation
########################
.PHONY: gen-mock
gen-mock:
	go run github.com/vektra/mockery/v2@v2.50 --config .mockery.yml

#########################
# Build
#########################
.PHONY: build
build:
	go build -o bin/raybot cmd/raybot/main.go

.PHONY: build-arm64
build-arm64:
	GOOS=linux GOARCH=arm64 go build -o bin/raybot-arm64 cmd/raybot/main.go

#########################
# Run
#########################
.PHONY: run
run:
	go run cmd/raybot/main.go raybot-config.yml

#########################
# Testing
#########################
.PHONY: test
test:
	go test -v -short ./...

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
