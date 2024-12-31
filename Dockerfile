# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install git and build dependencies
RUN apk add --no-cache git build-base

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/golicensemanager cmd/golicensemanager/main.go

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/bin/golicensemanager .
COPY --from=builder /app/.env.example .env

# Expose port
EXPOSE 8080

# Run the application
CMD ["./golicensemanager"]