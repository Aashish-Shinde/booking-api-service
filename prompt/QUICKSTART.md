# Quick Start Guide

## Prerequisites

- Go 1.24+
- MySQL 8.0+
- Docker & Docker Compose (optional, but recommended)

## 5-Minute Setup

### Step 1: Start MySQL Database

```bash
# Using Docker Compose (easiest)
docker-compose up -d

# Or use existing MySQL:
# Create database: CREATE DATABASE booking_api;
```

### Step 2: Run Migrations

```bash
# Download migrate tool first (one-time)
go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Run migrations
migrate -path db/migrations -database "mysql://root:root@tcp(localhost:3306)/booking_api" up
```

### Step 3: Build & Run

```bash
make run

# Or manually:
go build -o bin/booking-api-service cmd/main.go
./bin/booking-api-service
```

Server will start on `http://localhost:8080`

## Test the API

### 1. Create a Coach

```bash
# First, you need to create users and coaches in the database
# You can do this via SQL or by extending the API

# Example via SQL:
mysql -u root -p booking_api
INSERT INTO coaches (name, timezone) VALUES ('John Doe', 'UTC');
INSERT INTO users (name, timezone) VALUES ('Jane Doe', 'UTC');
```

### 2. Set Availability

```bash
curl -X POST http://localhost:8080/coaches/1/availability \
  -H "Content-Type: application/json" \
  -d '{
    "day_of_week": 1,
    "slots": [
      {"start_time": "10:00", "end_time": "13:00"},
      {"start_time": "16:00", "end_time": "18:00"}
    ]
  }'
```

### 3. Get Available Slots

```bash
curl "http://localhost:8080/coaches/1/slots?date=2026-04-13&timezone=UTC"
```

### 4. Book a Slot

```bash
curl -X POST http://localhost:8080/bookings \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 1,
    "coach_id": 1,
    "start_time": "2026-04-13T10:00:00Z"
  }'
```

### 5. View Bookings

```bash
curl http://localhost:8080/users/1/bookings
```

### 6. Cancel Booking

```bash
curl -X DELETE http://localhost:8080/bookings/1
```

## Useful Commands

```bash
# Build
make build

# Run tests
make test

# Run with development settings
DATABASE_URL="root:root@tcp(localhost:3306)/booking_api" make run

# Format code
make fmt

# Lint code
make lint

# Check for errors
make vet

# Clean build artifacts
make clean
```

## Environment Variables

Create `.env` file:

```env
DATABASE_URL=root:root@tcp(localhost:3306)/booking_api?parseTime=true
PORT=8080
ENVIRONMENT=development
LOG_LEVEL=info
```

## Database Access

```bash
# Connect to MySQL
mysql -u root -p booking_api

# View users
SELECT * FROM users;

# View coaches
SELECT * FROM coaches;

# View bookings
SELECT * FROM bookings;

# View availability
SELECT * FROM availability;
```

## Troubleshooting

### Database Connection Error
```
Error: failed to connect to database
Solution: 
1. Check MySQL is running: docker-compose ps
2. Verify credentials in DATABASE_URL
3. Check database exists: CREATE DATABASE booking_api;
```

### Port Already in Use
```
Error: listen tcp :8080: bind: address already in use
Solution: 
1. Kill existing process: lsof -i :8080 | kill -9 PID
2. Or use different port: PORT=8081 make run
```

### Migration Errors
```
Error: migration already exists
Solution:
make migrate-force
make migrate-up
```

### Build Errors
```
Error: package not found
Solution:
go mod download
go mod tidy
make build
```

## Project Structure

```
.
├── cmd/
│   └── main.go                 # Application entry point
├── internal/
│   ├── handler/                # HTTP handlers
│   ├── service/                # Business logic
│   ├── repository/             # Data access
│   ├── model/                  # Domain models
│   ├── dto/                    # Request/response objects
│   └── middleware/             # HTTP middleware
├── pkg/
│   ├── logger/                 # Logging utilities
│   └── utils/                  # Helper functions
├── db/
│   └── migrations/             # SQL migrations
├── bin/                        # Build output
├── Makefile                    # Build commands
├── docker-compose.yml          # Database setup
├── README.md                   # Main documentation
├── ARCHITECTURE.md             # Architecture details
└── go.mod / go.sum            # Dependency management
```

## Next Steps

1. **Understand the architecture**: Read [ARCHITECTURE.md](ARCHITECTURE.md)
2. **Explore the code**: Check `internal/` directory
3. **Run tests**: `make test`
4. **Add features**: Follow the clean architecture pattern
5. **Deploy**: Use Makefile targets for building

## API Documentation

See [README.md](README.md) for full API endpoint documentation.

## Support

For issues or questions:
1. Check the logs: Server logs all errors and operations
2. Review the README and ARCHITECTURE files
3. Check test files for examples
4. Verify database setup with `make migrate-up`

## Example Workflow

```bash
# 1. Start database
docker-compose up -d

# 2. Run migrations
migrate -path db/migrations \
  -database "mysql://root:root@tcp(localhost:3306)/booking_api" up

# 3. Build and run
make run

# 4. In another terminal, test the API
# Create coach and user (via SQL or API)
# Set availability
# Book slots
# View bookings
# Cancel if needed
```

Happy booking! 🚀
