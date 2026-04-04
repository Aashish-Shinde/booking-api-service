.PHONY: build run test migrate-up migrate-down clean help

BINARY_NAME=booking-api-service
DOCKER_IMAGE=booking-api:latest

help:
	@echo "Makefile targets:"
	@echo "  build         - Build the application"
	@echo "  run           - Run the application"
	@echo "  test          - Run tests"
	@echo "  migrate-up    - Run database migrations"
	@echo "  migrate-down  - Rollback database migrations"
	@echo "  clean         - Clean build artifacts"

build:
	@echo "Building application..."
	go build -o bin/$(BINARY_NAME) cmd/main.go

run: build
	@echo "Running application..."
	DATABASE_URL="root:root@tcp(localhost:3306)/booking_api?parseTime=true" ./bin/$(BINARY_NAME)

test:
	@echo "Running tests..."
	go test -v -cover ./...

migrate-up:
	@echo "Running migrations..."
	migrate -path db/migrations -database "mysql://root:root@tcp(localhost:3306)/booking_api" up

migrate-down:
	@echo "Rolling back migrations..."
	migrate -path db/migrations -database "mysql://root:root@tcp(localhost:3306)/booking_api" down

migrate-force:
	@echo "Force migration version..."
	migrate -path db/migrations -database "mysql://root:root@tcp(localhost:3306)/booking_api" force 1

clean:
	@echo "Cleaning..."
	rm -rf bin/

deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

fmt:
	@echo "Formatting code..."
	go fmt ./...

lint:
	@echo "Linting code..."
	golint ./...

vet:
	@echo "Running go vet..."
	go vet ./...
