.PHONY: help run-local build test lint test-api clean

help:
	@echo "Available targets:"
	@echo "  make run-local   - Run directly without building"
	@echo "  make build       - Build the binary"
	@echo "  make test        - Run Go tests"
	@echo "  make lint        - Run linter"
	@echo "  make test-api    - Run API integration tests (requires server running)"
	@echo "  make clean       - Clean build artifacts"

run-local:
	@echo "Starting server on :8080..."
	@go run ./cmd/api

build:
	@echo "Building binary..."
	@go build -o bin/api ./cmd/api
	@echo "Binary created at bin/api"

test:
	@go test -v ./...

lint:
	golangci-lint run ./...

test-api:
	@echo "Running API integration tests..."
	@echo "Note: Make sure the server is running (make run-local)"
	@./test-api.sh

clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/
	@echo "Done!"
