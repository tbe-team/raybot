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
