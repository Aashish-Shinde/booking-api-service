BINARY_NAME=booking-api-service

help:
	@echo "Makefile targets:"
	@echo "  build         - Build the application"
	@echo "  run           - Run the application"
	@echo "  test          - Run tests"
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

clean:
	@echo "Cleaning..."
	rm -rf bin/

deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

