# Development commands
default:
    @just --list

# Install development dependencies
setup:
    go mod download
    go mod tidy

# Run the application
run:
    go run cmd/golicensemanager/main.go

# Run tests
test:
    go test -v ./...

# Run tests with coverage
coverage:
    go test -coverprofile=coverage.out ./...
    go tool cover -html=coverage.out

# Run linter
lint:
    golangci-lint run

# Create a new migration
migrate-create name:
    migrate create -ext sql -dir scripts/db/migrations -seq {{name}}

# Run migrations up
migrate-up:
    migrate -database "postgres://postgres:postgres@localhost:5432/licensedb?sslmode=disable" -path scripts/db/migrations up

# Run migrations down
migrate-down:
    migrate -database "postgres://postgres:postgres@localhost:5432/licensedb?sslmode=disable" -path scripts/db/migrations down

# Generate mock data
generate-mocks:
    mockgen -destination test/mocks/repository_mock.go -package mocks github.com/yourusername/golicensemanager/internal/repository Repository

# Start development environment with Docker
docker-dev:
    docker-compose up -d

# Stop development environment
docker-down:
    docker-compose down