.PHONY: help build run run-dev test clean fmt vet

help:
	@echo "Available targets:"
	@echo "  make build    - Build the application"
	@echo "  make run      - Build and run the application"
	@echo "  make run-dev  - Run directly without building"
	@echo "  make test     - Run tests"
	@echo "  make fmt      - Format code"
	@echo "  make vet      - Run go vet"
	@echo "  make clean    - Clean build artifacts"

build:
	@echo "Building go-runtime-demo..."
	@go build -o bin/api ./cmd/api
	@echo "Build complete: bin/api"

run: build
	@echo "Starting server on :8080..."
	@./bin/api

run-dev:
	@echo "Starting server on :8080..."
	@go run ./cmd/api

test:
	@echo "Running tests..."
	@go test -v ./...

fmt:
	@echo "Formatting code..."
	@go fmt ./...
	@echo "Code formatted"

# Run go vet
vet:
	@echo "Running go vet..."
	@go vet ./...
	@echo "Vet complete"

clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/
	@echo "Clean complete"

check: fmt vet test
	@echo "All checks passed"

