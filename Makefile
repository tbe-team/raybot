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
	go run cmd/raybot/main.go config.yml

#########################
# Testing
#########################
.PHONY: test
test:
	go test -v -short ./...

.PHONY: test-cov
test-cov:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report saved to coverage.html"

########################
# Lint
########################
.PHONY: lint-go
lint-go:
	golangci-lint run ./... --config .golangci.yml
