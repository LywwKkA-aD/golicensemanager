#!/bin/bash

# Project name
PROJECT_NAME="golicensemanager"

# Create main project directory

# Create root level files
touch README.md
touch .gitignore
touch go.mod
touch go.sum
touch Makefile
touch justfile
touch .env.example
touch docker-compose.yml
touch Dockerfile

# Create main directories and their subdirectories
mkdir -p {cmd,internal,pkg,api,scripts,deployments,test,docs}

# cmd structure
mkdir -p cmd/$PROJECT_NAME
touch cmd/$PROJECT_NAME/main.go
touch cmd/$PROJECT_NAME/.gitkeep

# internal structure
mkdir -p internal/{app,config,middleware,models,repository,service,utils}
touch internal/{app,config,middleware,models,repository,service,utils}/.gitkeep

# Create specific internal subdirectories
mkdir -p internal/app/handlers
touch internal/app/handlers/.gitkeep

mkdir -p internal/repository/postgres
touch internal/repository/postgres/.gitkeep

# pkg structure (shared libraries)
touch pkg/.gitkeep

# api structure
mkdir -p api/{http,proto}
touch api/http/.gitkeep
touch api/proto/.gitkeep

# Create OpenAPI/Swagger documentation directory
mkdir -p api/http/swagger
touch api/http/swagger/.gitkeep

# scripts structure
mkdir -p scripts/{db,dev,ci}
touch scripts/db/.gitkeep
touch scripts/dev/.gitkeep
touch scripts/ci/.gitkeep

# Create database migration directory
mkdir -p scripts/db/migrations
touch scripts/db/migrations/.gitkeep

# deployments structure
mkdir -p deployments/{docker,k8s}
touch deployments/docker/.gitkeep

# test structure
mkdir -p test/{integration,mocks,fixtures}
touch test/{integration,mocks,fixtures}/.gitkeep

# docs structure
mkdir -p docs/{swagger,architecture,api,development}
touch docs/{swagger,architecture,api,development}/.gitkeep

# Create .gitignore content
cat > .gitignore << 'EOL'
# Binaries
bin/
*.exe
*.exe~
*.dll
*.so
*.dylib

# Test binary, built with `go test -c`
*.test

# Output of the go coverage tool
*.out

# Dependency directories
vendor/

# IDE specific files
.idea/
.vscode/
*.swp
*.swo

# Environment files
.env
.env.local

# OS specific
.DS_Store
Thumbs.db

# Build output
dist/
EOL

# Make all scripts executable
chmod +x scripts/**/*.sh 2>/dev/null || true

echo "Project structure created successfully!"