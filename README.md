# Booking API Service

A production-grade REST API for an appointment booking system built in Go with MySQL.

## Features

- ✅ Coach availability management (weekly recurring + exceptions)
- ✅ 30-minute slot booking system
- ✅ Double booking prevention with database constraints
- ✅ Safe concurrent booking with transactions
- ✅ Timezone support
- ✅ Booking modification and cancellation
- ✅ Clean architecture (Handler → Service → Repository)
- ✅ Structured logging with Zap
- ✅ RESTful API with Chi router

## Tech Stack

- **Language**: Go 1.24
- **Framework**: Chi
- **Database**: MySQL 8.0
- **Migrations**: golang-migrate
- **Logging**: Zap
- **Architecture**: Clean Architecture

## Prerequisites

- Go 1.24+
- MySQL 8.0+
- Docker & Docker Compose (optional)

## Setup

### 1. Clone the repository

```bash
cd /home/ashish/Desktop/personal/booking-api-service
```

### 2. Install dependencies

```bash
go mod download
go mod tidy
```

### 3. Start MySQL database

Using Docker Compose (recommended):
```bash
docker-compose up -d
```

Or use an existing MySQL instance and update `.env`

### 4. Run database migrations

```bash
make migrate-up
```

Or manually:
```bash
migrate -path db/migrations -database "mysql://root:root@tcp(localhost:3306)/booking_api" up
```

### 5. Run the application

```bash
make run
```

Or manually:
```bash
go build -o bin/booking-api-service cmd/main.go
./bin/booking-api-service
```

The API will be available at `http://localhost:8080`

## API Endpoints

### Availability Management

#### Set Weekly Availability
```
POST /coaches/{coach_id}/availability

{
  "day_of_week": 1,
  "slots": [
    { "start_time": "10:00", "end_time": "13:00" },
    { "start_time": "16:00", "end_time": "18:00" }
  ]
}
```

#### Add Exception
```
POST /coaches/{coach_id}/exceptions

{
  "date": "2026-04-10",
  "is_available": false
}

// Or with custom times
{
  "date": "2026-04-10",
  "is_available": true,
  "start_time": "09:00",
  "end_time": "12:00"
}
```

#### Get Available Slots
```
GET /coaches/{coach_id}/slots?date=2026-04-10&timezone=Asia/Kolkata

Response:
{
  "slots": [
    {
      "start_time": "2026-04-10T10:00:00Z",
      "end_time": "2026-04-10T10:30:00Z"
    },
    ...
  ]
}
```

### Booking Management

#### Create Booking
```
POST /bookings

{
  "user_id": 1,
  "coach_id": 1,
  "start_time": "2026-04-10T10:00:00Z",
  "idempotency_key": "optional-unique-key"
}

Response:
{
  "id": 1,
  "user_id": 1,
  "coach_id": 1,
  "start_time": "2026-04-10T10:00:00Z",
  "end_time": "2026-04-10T10:30:00Z",
  "status": "ACTIVE",
  "created_at": "2026-04-09T10:00:00Z",
  "updated_at": "2026-04-09T10:00:00Z"
}
```

#### Get User Bookings
```
GET /users/{user_id}/bookings

Response:
{
  "bookings": [...],
  "count": 5
}
```

#### Modify Booking
```
PUT /bookings/{id}

{
  "start_time": "2026-04-10T11:00:00Z"
}
```

#### Cancel Booking
```
DELETE /bookings/{id}
```

## Database Schema

### users
- id, name, timezone, created_at, updated_at, deleted

### coaches
- id, name, timezone, created_at, updated_at, deleted

### availability
- id, coach_id, day_of_week, start_time, end_time, created_at, updated_at, deleted

### availability_exceptions
- id, coach_id, date, start_time, end_time, is_available, created_at, updated_at, deleted

### bookings
- id, user_id, coach_id, start_time (UTC), end_time (UTC), status, idempotency_key, created_at, updated_at, deleted
- **UNIQUE(coach_id, start_time)** - Prevents double booking

## Key Implementation Details

### Concurrency Safety
- Uses database transactions for booking creation
- Row-level locking with `FOR UPDATE` clause
- UNIQUE constraint on (coach_id, start_time)

### Slot Generation
- Generates 30-minute slots dynamically from availability windows
- Applies exceptions (overrides)
- Removes booked slots
- Converts to requested timezone at API layer

### Timezone Handling
- All times stored in UTC in database
- Coach/user have configurable timezone
- Conversion happens at API boundary

### Validation
- Times must align to 30-minute boundaries
- No overlapping availability windows
- Cannot book past slots
- Cannot modify bookings that have started

## Make Targets

```bash
make build         # Build the application
make run           # Build and run
make test          # Run tests
make migrate-up    # Run database migrations
make migrate-down  # Rollback migrations
make clean         # Clean build artifacts
make deps          # Download dependencies
make fmt           # Format code
make lint          # Lint code
make vet           # Run go vet
```

## Project Structure

```
.
├── cmd/
│   └── main.go
├── internal/
│   ├── handler/          # HTTP handlers
│   ├── service/          # Business logic
│   ├── repository/       # Data access
│   ├── model/            # Domain models
│   ├── dto/              # Request/response objects
│   └── middleware/       # HTTP middleware
├── pkg/
│   ├── logger/           # Logging utilities
│   └── utils/            # Helper functions
├── db/
│   └── migrations/       # SQL migrations
├── Makefile
├── docker-compose.yml
├── go.mod
└── go.sum
```

## Error Handling

The API returns appropriate HTTP status codes:

- `200 OK` - Successful operation
- `201 Created` - Resource created
- `400 Bad Request` - Invalid input
- `409 Conflict` - Double booking conflict
- `500 Internal Server Error` - Server error

## Logging

Structured logging using Zap:
- All errors logged with context
- Production-ready format
- Configurable log levels

## Testing

Run tests with:
```bash
make test
```

Tests cover:
- Slot generation logic
- Booking validation
- Concurrency scenarios
- Timezone handling

## Configuration

Create a `.env` file based on `.env.example`:

```
DATABASE_URL=root:root@tcp(localhost:3306)/booking_api?parseTime=true
PORT=8080
ENVIRONMENT=development
LOG_LEVEL=info
```

## Future Enhancements

- [ ] Authentication & authorization
- [ ] Rate limiting
- [ ] Caching layer (Redis)
- [ ] Webhook notifications
- [ ] Analytics & reporting
- [ ] Email notifications
- [ ] Payment integration

## License

MIT
