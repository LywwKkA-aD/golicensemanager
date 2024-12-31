BINARY_NAME=golicensemanager
BUILD_DIR=bin

.PHONY: all build clean test coverage deps lint run

all: clean build

build:
	@echo "Building..."
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) cmd/$(BINARY_NAME)/main.go

clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)

test:
	@echo "Running tests..."
	@go test -v ./...

coverage:
	@echo "Running tests with coverage..."
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out

deps:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy

lint:
	@echo "Running linter..."
	@golangci-lint run

run:
	@echo "Running application..."
	@go run cmd/$(BINARY_NAME)/main.go