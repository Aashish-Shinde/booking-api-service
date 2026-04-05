# Booking API Service

A production-grade REST API for an appointment booking system built in Go with MySQL. This service manages coach availability, booking slots, and provides timezone-aware availability calculations with concurrency-safe operations.

## 📋 Table of Contents

- [Features](#features)
- [Tech Stack](#tech-stack)
- [Prerequisites](#prerequisites)
- [Setup & Installation](#setup--installation)
- [Configuration](#configuration)
- [API Endpoints](#api-endpoints)
- [Core Concepts](#core-concepts)
- [Timezone Handling](#timezone-handling)
- [Database Schema](#database-schema)
- [Architecture](#architecture)
- [Make Commands](#make-commands)
- [Testing](#testing)
- [Troubleshooting](#troubleshooting)
- [Future Enhancements](#future-enhancements)

## ✨ Features

- ✅ **Coach Availability Management** - Weekly recurring + one-time exceptions
- ✅ **30-Minute Slot System** - Automatic 30-minute slot generation from availability windows
- ✅ **Double Booking Prevention** - Database constraints & transaction-based locking
- ✅ **Concurrent Booking Safety** - Row-level locking with `FOR UPDATE` clause
- ✅ **Timezone Support** - Full timezone awareness for coaches and users
- ✅ **Idempotent Operations** - Safe retry mechanism for booking creation
- ✅ **Booking Lifecycle** - Create, modify, cancel bookings with status tracking
- ✅ **Clean Architecture** - Separation of concerns (Handler → Service → Repository)
- ✅ **Structured Logging** - Production-ready logging with request context
- ✅ **RESTful API** - Chi router with middleware support
- ✅ **Graceful Shutdown** - Proper cleanup of resources

## 🛠️ Tech Stack

- **Language**: Go 1.24
- **Web Framework**: Chi (Lightweight HTTP router)
- **Database**: MySQL 8.0+
- **Database Migrations**: golang-migrate/migrate
- **Logging**: Go's `log/slog` (standard library)
- **Architecture Pattern**: Clean Architecture (Layered)
- **Concurrency**: Database transactions with row-level locking

## 📋 Prerequisites

- **Go** 1.24 or higher
- **MySQL** 8.0 or higher
- **Make** (for build automation)
- **golang-migrate** (for database migrations)

Install golang-migrate:
```bash
go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

## 🚀 Setup & Installation

### Step 1: Clone the Repository

```bash
cd /home/ashish/Desktop/personal/booking-api-service
```

### Step 2: Install Go Dependencies

```bash
go mod download
go mod tidy
```

### Step 3: Start MySQL Database

Ensure you have MySQL 8.0+ running on your system. Update the `.env` file with your MySQL connection details:
```bash
DATABASE_URL=your_user:your_password@tcp(your_host:3306)/booking_api?parseTime=true
```

### Step 4: Create Environment Configuration

Create a `.env` file in the project root:

```bash
# Database Configuration
DATABASE_URL=root:root@tcp(localhost:3306)/booking_api?parseTime=true

# Server Configuration
PORT=8080
ENVIRONMENT=development
LOG_LEVEL=info
```

### Step 5: Run Database Migrations

```bash
make migrate-up
```

Or manually:
```bash
migrate -path db/migrations -database "mysql://root:root@tcp(localhost:3306)/booking_api" up
```

This creates all tables and seeds initial test data.

### Step 6: Run the Application

```bash
make run
```

Or manually:
```bash
go build -o bin/booking-api-service cmd/main.go
./bin/booking-api-service
```

The API will be available at `http://localhost:8080`

### Verify Installation

```bash
curl http://localhost:8080/health
```

## 🔧 Configuration

### Environment Variables

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `DATABASE_URL` | MySQL connection string | - | ✅ |
| `PORT` | Server port | `8080` | ❌ |
| `ENVIRONMENT` | dev/production | `development` | ❌ |
| `LOG_LEVEL` | debug/info/warn/error | `info` | ❌ |

### Database Connection String Format

```
user:password@tcp(host:port)/database?parseTime=true
```

### Example Configurations

**Development**
```
DATABASE_URL=root:root@tcp(localhost:3306)/booking_api?parseTime=true
```

**Production**
```
DATABASE_URL=user:password@tcp(prod-db.example.com:3306)/booking_api_prod?parseTime=true
```

## 📡 API Endpoints

### Base URL
All endpoints are prefixed with the base URL: `http://localhost:8080`

### Availability Management

#### 1. Get Available Slots for a Coach

Retrieves all available 30-minute slots for a coach on a specific date in the user's timezone.

**Request**
```http
GET /api/availability/coaches/{coachID}/slots?date=2026-04-06&userID=1
```

**Query Parameters**
- `date` (required): Date in format `YYYY-MM-DD`
- `userID` (required): User ID to get user's timezone for conversions

**Response (200 OK)**
```json
{
  "coach_id": 1,
  "date": "2026-04-06",
  "timezone": "America/New_York",
  "slots": [
    {
      "start_time": "2026-04-06T14:00:00Z",
      "end_time": "2026-04-06T14:30:00Z"
    },
    {
      "start_time": "2026-04-06T14:30:00Z",
      "end_time": "2026-04-06T15:00:00Z"
    }
  ]
}
```

**Error Responses**
- `400 Bad Request`: Invalid date format or missing parameters
- `404 Not Found`: Coach or user not found
- `500 Internal Server Error`: Database error

#### 2. Get Availability for a Specific Date

Gets the weekly availability pattern for a coach.

**Request**
```http
GET /api/availability/coaches/{coachID}/availability
```

**Response (200 OK)**
```json
{
  "coach_id": 1,
  "availability": [
    {
      "day_of_week": 1,
      "start_time": "09:00",
      "end_time": "17:00"
    },
    {
      "day_of_week": 2,
      "start_time": "09:00",
      "end_time": "17:00"
    }
  ]
}
```

#### 3. Set Weekly Availability

Sets recurring weekly availability for a coach.

**Request**
```http
POST /api/availability/coaches/{coachID}/availability
Content-Type: application/json

{
  "day_of_week": 1,
  "start_time": "09:00",
  "end_time": "17:00"
}
```

**Parameters**
- `day_of_week`: 1 (Monday) to 5 (Friday), 0 = Sunday, 6 = Saturday
- `start_time`: Time in HH:MM format (coach's local timezone)
- `end_time`: Time in HH:MM format (coach's local timezone)

**Response (201 Created)**
```json
{
  "id": 1,
  "coach_id": 1,
  "day_of_week": 1,
  "start_time": "09:00",
  "end_time": "17:00",
  "created_at": "2026-04-05T10:00:00Z"
}
```

#### 4. Add Availability Exception

Creates an exception (override) for a specific date.

**Request**
```http
POST /api/availability/coaches/{coachID}/exceptions
Content-Type: application/json

{
  "date": "2026-04-10",
  "is_available": false
}
```

**With Custom Hours**
```json
{
  "date": "2026-04-10",
  "is_available": true,
  "start_time": "10:00",
  "end_time": "14:00"
}
```

**Response (201 Created)**
```json
{
  "id": 1,
  "coach_id": 1,
  "date": "2026-04-10",
  "is_available": false,
  "start_time": null,
  "end_time": null,
  "created_at": "2026-04-05T10:00:00Z"
}
```

### Booking Management

#### 1. Create a Booking

Creates a new booking for a user with a coach.

**Request**
```http
POST /api/bookings
Content-Type: application/json

{
  "user_id": 1,
  "coach_id": 1,
  "start_time": "2026-04-06T14:00:00Z",
  "idempotency_key": "booking-20260405-001"
}
```

**Parameters**
- `user_id` (required): User ID
- `coach_id` (required): Coach ID
- `start_time` (required): Start time in RFC3339 format (UTC)
- `idempotency_key` (optional): Unique key for idempotent requests

**Response (201 Created)**
```json
{
  "id": 1,
  "user_id": 1,
  "coach_id": 1,
  "start_time": "2026-04-06T14:00:00Z",
  "end_time": "2026-04-06T14:30:00Z",
  "status": "ACTIVE",
  "created_at": "2026-04-05T10:00:00Z",
  "updated_at": "2026-04-05T10:00:00Z"
}
```

**Error Responses**
- `400 Bad Request`: Invalid input or time not aligned to 30-minute boundary
- `409 Conflict`: Slot already booked or time is not available
- `500 Internal Server Error`: Database error

#### 2. Get User Bookings

Retrieves all bookings for a user.

**Request**
```http
GET /api/users/{userID}/bookings?status=ACTIVE
```

**Query Parameters**
- `status` (optional): Filter by status (ACTIVE, COMPLETED, CANCELLED)

**Response (200 OK)**
```json
{
  "bookings": [
    {
      "id": 1,
      "user_id": 1,
      "coach_id": 1,
      "start_time": "2026-04-06T14:00:00Z",
      "end_time": "2026-04-06T14:30:00Z",
      "status": "ACTIVE",
      "created_at": "2026-04-05T10:00:00Z",
      "updated_at": "2026-04-05T10:00:00Z"
    }
  ],
  "count": 1
}
```

#### 3. Get Booking Details

Retrieves details of a specific booking.

**Request**
```http
GET /api/bookings/{bookingID}
```

**Response (200 OK)**
```json
{
  "id": 1,
  "user_id": 1,
  "coach_id": 1,
  "start_time": "2026-04-06T14:00:00Z",
  "end_time": "2026-04-06T14:30:00Z",
  "status": "ACTIVE",
  "created_at": "2026-04-05T10:00:00Z",
  "updated_at": "2026-04-05T10:00:00Z"
}
```

#### 4. Modify a Booking

Updates an existing booking to a new time.

**Request**
```http
PUT /api/bookings/{bookingID}
Content-Type: application/json

{
  "start_time": "2026-04-06T15:00:00Z"
}
```

**Response (200 OK)**
```json
{
  "id": 1,
  "user_id": 1,
  "coach_id": 1,
  "start_time": "2026-04-06T15:00:00Z",
  "end_time": "2026-04-06T15:30:00Z",
  "status": "ACTIVE",
  "created_at": "2026-04-05T10:00:00Z",
  "updated_at": "2026-04-05T10:05:00Z"
}
```

**Constraints**
- Cannot modify bookings in the past
- Cannot modify bookings that have already started
- New time must be available

#### 5. Cancel a Booking

Cancels an existing booking.

**Request**
```http
DELETE /api/bookings/{bookingID}
```

**Response (204 No Content)**

**Error Responses**
- `404 Not Found`: Booking not found
- `400 Bad Request`: Cannot cancel past bookings

## 🔑 Core Concepts

### Slots

Slots are **30-minute time windows** generated dynamically from coach availability. Key points:

- **Automatic Generation**: Slots are computed on-the-fly from availability windows
- **Fixed Duration**: All slots are exactly 30 minutes
- **Timezone-Aware**: Generated in coach's timezone, returned in requested timezone
- **Booking-Aware**: Excludes already-booked slots

**Example:**
If a coach is available 9:00 AM - 12:00 PM, the system generates:
- 9:00 - 9:30
- 9:30 - 10:00
- 10:00 - 10:30
- ... and so on

### Availability Windows

Availability windows define when a coach is available:

- **Weekly Recurring**: Set via `day_of_week` (1=Monday, 7=Sunday)
- **Exceptions**: One-time overrides for specific dates
- **Multiple Windows**: Coach can have multiple non-overlapping windows per day

**Example:**
Coach Alice might have:
- Monday-Friday: 9:00 AM - 5:00 PM (recurring)
- April 10: Day off (exception)
- April 12: 10:00 AM - 2:00 PM (special hours exception)

### Booking Status

Bookings have three states:

1. **ACTIVE**: Current/future booking
2. **COMPLETED**: Past booking that was fulfilled
3. **CANCELLED**: Booking that was cancelled

### Idempotency

The booking creation endpoint supports idempotent requests using `idempotency_key`:

```json
{
  "user_id": 1,
  "coach_id": 1,
  "start_time": "2026-04-06T14:00:00Z",
  "idempotency_key": "booking-20260405-001"
}
```

If you submit the same request twice with the same key, you get the same response without creating duplicate bookings.

## 🌍 Timezone Handling

### Overview

The system supports full timezone awareness:

- **Storage**: All times stored in UTC in the database
- **Coach Timezone**: Coach's availability is set in their local timezone
- **User Timezone**: Slots are converted to user's timezone in API response
- **Automatic Conversion**: Handled transparently at API boundaries

### How It Works

1. **Setting Availability**
   ```
   Coach in America/New_York sets: 9:00 AM - 5:00 PM
   → Stored as UTC equivalent in database
   ```

2. **Getting Slots**
   ```
   User in Asia/Tokyo requests slots for 2026-04-06
   → System loads coach's availability (9:00-17:00 NY time)
   → Converts to coach's UTC equivalent
   → Generates slots
   → Converts slots to user's timezone (Tokyo)
   → Returns to user
   ```

3. **Creating Booking**
   ```
   User provides: 2026-04-06T14:00:00Z (UTC)
   → Checked against coach's availability in coach's timezone
   → If valid, saved to database in UTC
   ```

### Supported Timezones

The system supports all IANA timezone identifiers. Common examples:

- **North America**: America/New_York, America/Chicago, America/Los_Angeles, America/Denver, Canada/Toronto
- **Europe**: Europe/London, Europe/Paris, Europe/Berlin, Europe/Madrid
- **Asia**: Asia/Tokyo, Asia/Shanghai, Asia/Hong_Kong, Asia/Singapore, Asia/Bangkok
- **Australia**: Australia/Sydney, Australia/Melbourne, Australia/Brisbane
- **India**: Asia/Kolkata
- **UTC**: UTC

### Testing Timezone Logic

Use the test data provided in the seed:

**Coach Alice** (America/New_York, 9 AM - 5 PM):
- Test with User 1 (America/New_York) - Same timezone
- Test with User 3 (Asia/Tokyo) - 13-14 hour difference

**Coach Fiona** (Australia/Melbourne, 8 AM - 6 PM):
- Test with User 7 (Europe/Paris) - 8-9 hour difference

## 🗄️ Database Schema

### users
```sql
CREATE TABLE users (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  name VARCHAR(255),
  timezone VARCHAR(50),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted TINYINT DEFAULT 0
)
```

- **id**: Unique user identifier
- **name**: User's name
- **timezone**: User's IANA timezone (e.g., Asia/Tokyo)
- **created_at**: Record creation timestamp (UTC)
- **updated_at**: Last modification timestamp (UTC)
- **deleted**: Soft delete flag (0=active, 1=deleted)

### coaches
```sql
CREATE TABLE coaches (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  name VARCHAR(255),
  timezone VARCHAR(50),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted TINYINT DEFAULT 0
)
```

- Similar structure to users table
- **timezone**: Coach's local timezone for availability scheduling

### availability
```sql
CREATE TABLE availability (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  coach_id BIGINT NOT NULL,
  day_of_week INT (1-7),
  start_time TIME,
  end_time TIME,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted TINYINT DEFAULT 0,
  FOREIGN KEY (coach_id) REFERENCES coaches(id)
)
```

- **coach_id**: Foreign key to coaches table
- **day_of_week**: 1=Monday, 2=Tuesday, ..., 7=Sunday
- **start_time**: Availability start time (HH:MM format, coach's local timezone)
- **end_time**: Availability end time (HH:MM format, coach's local timezone)

### availability_exceptions
```sql
CREATE TABLE availability_exceptions (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  coach_id BIGINT NOT NULL,
  date DATE,
  start_time TIME,
  end_time TIME,
  is_available BOOLEAN,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted TINYINT DEFAULT 0,
  FOREIGN KEY (coach_id) REFERENCES coaches(id)
)
```

- **coach_id**: Foreign key to coaches table
- **date**: Date of the exception
- **is_available**: true (special availability), false (day off)
- **start_time/end_time**: Custom hours (only if is_available=true)

### bookings
```sql
CREATE TABLE bookings (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  user_id BIGINT NOT NULL,
  coach_id BIGINT NOT NULL,
  start_time DATETIME,
  end_time DATETIME,
  status VARCHAR(20),
  idempotency_key VARCHAR(255) UNIQUE,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted TINYINT DEFAULT 0,
  FOREIGN KEY (user_id) REFERENCES users(id),
  FOREIGN KEY (coach_id) REFERENCES coaches(id),
  UNIQUE KEY unique_coach_slot (coach_id, start_time)
)
```

- **user_id**: Foreign key to users table
- **coach_id**: Foreign key to coaches table
- **start_time**: Booking start time (UTC)
- **end_time**: Booking end time (UTC) - Always start_time + 30 minutes
- **status**: ACTIVE, COMPLETED, or CANCELLED
- **idempotency_key**: Unique key for idempotent request handling
- **UNIQUE(coach_id, start_time)**: Prevents double-booking for same coach at same time

## 🏗️ Architecture

### Overview

The application follows **Clean Architecture** principles with clear separation of concerns:

```
HTTP Request
    ↓
[Handler Layer] - HTTP request/response handling
    ↓
[Service Layer] - Business logic and validation
    ↓
[Repository Layer] - Data access (database)
    ↓
Database
```

### Layer Responsibilities

#### Handler Layer (`internal/handler/`)

- Receives HTTP requests
- Validates input data structure
- Calls service methods
- Returns HTTP responses
- Handles error conversion to HTTP status codes

**Files:**
- `availability_handler.go`: Handles availability endpoints
- `booking_handler.go`: Handles booking endpoints
- `handler_test.go`: Handler tests

#### Service Layer (`internal/service/`)

- Implements business logic
- Validates business rules
- Orchestrates repository calls
- Handles timezone conversions
- Manages transactions

**Files:**
- `availability_service.go`: Availability business logic (slot generation, exception handling)
- `booking_service.go`: Booking operations
- `availability_service_test.go`: Slot generation tests
- `service_validation_test.go`: Business validation tests

**Key Methods:**
- `GetSlots()`: Generates available slots for a date
- `CreateAvailability()`: Sets weekly availability
- `AddException()`: Creates one-time overrides
- `CreateBooking()`: Books a slot with conflict checking
- `ModifyBooking()`: Reschedules a booking
- `CancelBooking()`: Cancels a booking

#### Repository Layer (`internal/repository/`)

- Abstracts database access
- Executes SQL queries
- Manages transactions
- Ensures data consistency

**Files:**
- `availability_repository.go`: Availability queries
- `availability_exception_repository.go`: Exception queries
- `booking_repository.go`: Booking queries
- `user_repository.go`: User queries
- `coach_repository.go`: Coach queries
- `interfaces.go`: Repository interfaces
- `repository_test.go`: Repository tests

#### Model Layer (`internal/model/`)

- Defines domain models
- Matches database schema
- Pure data structures

**Files:**
- `models.go`: All model definitions (User, Coach, Booking, etc.)

#### DTO Layer (`internal/dto/`)

- Request/response objects
- Maps between HTTP and domain
- Validation annotations

**Files:**
- `requests.go`: All request DTOs

#### Middleware (`internal/middleware/`)

- Cross-cutting concerns
- Timeout handling
- Logging
- Error recovery

**Files:**
- `timeout.go`: Request timeout middleware

#### Configuration (`pkg/config/`)

- Environment variable loading
- Configuration management
- Database URL building

**Files:**
- `config.go`: Configuration logic

#### Utilities (`pkg/utils/`)

- Helper functions
- Time utilities
- Timezone conversions
- Validation helpers

**Files:**
- `time_utils.go`: Time parsing and formatting
- `timezone.go`: Timezone conversion utilities
- `validator.go`: Input validation

#### Logger (`pkg/logger/`)

- Structured logging
- Consistent log format
- Context propagation

**Files:**
- `logger.go`: Logger setup

### Data Flow Example: Get Available Slots

```
1. HTTP Request (Handler)
   GET /api/availability/coaches/1/slots?date=2026-04-06&userID=1
   
2. Handler Layer (availability_handler.go)
   - Validate query parameters
   - Parse date
   - Call availabilityService.GetSlots()
   
3. Service Layer (availability_service.go)
   - Load coach details (timezone)
   - Load user details (timezone)
   - Get weekly availability for that date
   - Get exceptions for that date
   - Generate 30-minute slots
   - Filter out booked slots
   - Convert to user's timezone
   - Build response
   
4. Repository Layer
   - Coach repository: Get coach by ID
   - Availability repository: Get availability for day_of_week
   - Exception repository: Get exceptions for date
   - Booking repository: Get bookings for date range
   
5. Database Queries
   - SELECT * FROM coaches WHERE id = 1
   - SELECT * FROM availability WHERE coach_id = 1 AND day_of_week = 2
   - SELECT * FROM availability_exceptions WHERE coach_id = 1 AND date = '2026-04-06'
   - SELECT * FROM bookings WHERE coach_id = 1 AND start_time >= ... AND start_time < ...
   
6. Response
   {
     "coach_id": 1,
     "slots": [...]
   }
```

### Concurrency & Safety

#### Double Booking Prevention

The system prevents double bookings through multiple layers:

1. **Database Constraint** (Primary)
   ```sql
   UNIQUE KEY unique_coach_slot (coach_id, start_time)
   ```

2. **Row-Level Locking** (Service Layer)
   ```sql
   SELECT ... FROM bookings WHERE coach_id = ? AND start_time = ? FOR UPDATE
   ```

3. **Transaction Isolation** (Service Layer)
   - Booking creation wrapped in transaction
   - Checks for conflicts before INSERT
   - If conflict exists, rolls back transaction

#### Request Idempotency

Idempotency key mechanism:

```go
// First request
POST /api/bookings
{
  "user_id": 1,
  "coach_id": 1,
  "start_time": "2026-04-06T14:00:00Z",
  "idempotency_key": "booking-001"
}
→ Response: 201 Created, Booking ID 42

// Retry with same key
POST /api/bookings
{
  "user_id": 1,
  "coach_id": 1,
  "start_time": "2026-04-06T14:00:00Z",
  "idempotency_key": "booking-001"
}
→ Response: 200 OK, Booking ID 42 (same booking returned)
```

## 📦 Make Commands

```bash
make build           # Build the application binary
make run             # Build and run the application
make test            # Run all unit tests
make test-verbose    # Run tests with verbose output
make migrate-up      # Run database migrations (up)
make migrate-down    # Rollback database migrations (down)
make migrate-force   # Force migration to specific version
make clean           # Clean build artifacts and binaries
make deps            # Download and verify dependencies
make fmt             # Format code with gofmt
make lint            # Run golangci-lint
make vet             # Run go vet for suspicious constructs
make help            # Show available commands
```

## 🧪 Testing

### Running Tests

```bash
# Run all tests
make test

# Run with verbose output
make test-verbose

# Run specific test file
go test ./internal/service -v

# Run specific test function
go test ./internal/service -run TestGetSlots -v

# Run with coverage
go test -cover ./...

# Run with coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Test Coverage

Tests are included for:

**Service Layer** (`internal/service/`):
- `TestGetSlots`: Slot generation logic
- `TestCreateBooking`: Booking creation with validations
- `TestModifyBooking`: Booking rescheduling
- `TestCancelBooking`: Booking cancellation
- Timezone conversion tests
- Double booking prevention tests

**Repository Layer** (`internal/repository/`):
- Data access operations
- Transaction handling

**Handler Layer** (`internal/handler/`):
- HTTP request/response mapping
- Error handling

### Key Test Scenarios

1. **Slot Generation**
   - Generate slots within business hours
   - Exclude weekends (if not available)
   - Handle exceptions (days off, special hours)
   - Exclude booked slots
   - Timezone conversion

2. **Booking Validation**
   - Cannot book past slots
   - Cannot double book same slot
   - Times must align to 30-minute boundaries
   - User and coach must exist

3. **Concurrency**
   - Multiple concurrent booking requests
   - Race condition handling
   - Transaction isolation

## 🔍 Troubleshooting

### Common Issues and Solutions

#### 1. Database Connection Failed

**Error:** `failed to connect to database`

**Solution:**
```bash
# Check connection string in .env
DATABASE_URL=root:root@tcp(localhost:3306)/booking_api?parseTime=true

# Verify MySQL is accessible
mysql -h localhost -u root -p -e "USE booking_api; SHOW TABLES;"
```

#### 2. Migration Errors

**Error:** `error: Dirty database version`

**Solution:**
```bash
# Force migration to version (check current version first)
migrate -path db/migrations -database "mysql://root:root@tcp(localhost:3306)/booking_api" force 2
```

#### 3. Port Already in Use

**Error:** `listen tcp :8080: bind: address already in use`

**Solution:**
```bash
# Find process using port 8080
lsof -i :8080

# Kill the process
kill -9 <PID>

# Or use different port
PORT=8081 make run
```

#### 4. Slot Generation Returns Empty

**Check:**
1. Coach has availability set for the requested date
2. Date format is correct (YYYY-MM-DD)
3. Coach timezone is valid
4. Check logs for specific errors

**Debug:**
```bash
# Check coach availability
mysql -h localhost -u root -p booking_api
SELECT * FROM availability WHERE coach_id = 1;

# Check for exceptions
SELECT * FROM availability_exceptions WHERE coach_id = 1;

# Check for bookings
SELECT * FROM bookings WHERE coach_id = 1;
```

#### 5. Timezone Conversion Issues

**Verify:**
1. Coach timezone is valid IANA timezone
2. User timezone is valid IANA timezone
3. All times in requests are in RFC3339 format

**Example:**
```bash
curl -X GET \
  "http://localhost:8080/api/availability/coaches/1/slots?date=2026-04-06&userID=3" \
  -H "Content-Type: application/json"
```

## 🚀 Future Enhancements

### Planned Features

- [ ] **Authentication & Authorization**
  - JWT-based authentication
  - Role-based access control (users, coaches, admins)
  - Permission management

- [ ] **Performance Optimization**
  - Caching layer (Redis) for availability slots
  - Database query optimization with indexes
  - Connection pooling optimization

- [ ] **Notifications**
  - Email notifications for bookings
  - SMS alerts
  - Webhook support for external systems
  - Push notifications

- [ ] **Analytics & Reporting**
  - Booking statistics
  - Coach utilization metrics
  - User analytics
  - Revenue reports

- [ ] **Advanced Availability**
  - Recurring exceptions (e.g., every 2 weeks off)
  - Buffer time between bookings
  - Maximum bookings per day limit
  - Lead time requirements

- [ ] **Payment Integration**
  - Stripe integration
  - Refund handling
  - Invoice generation

- [ ] **Admin Dashboard**
  - Booking management UI
  - Coach/user management
  - Exception calendar
  - Analytics dashboard

- [ ] **API Improvements**
  - OpenAPI/Swagger documentation
  - Rate limiting
  - Request validation middleware
  - API versioning

## 📄 License

MIT License - See LICENSE file for details

## 🤝 Contributing

Contributions are welcome! Please follow these steps:

1. Create a feature branch (`git checkout -b feature/amazing-feature`)
2. Commit your changes (`git commit -m 'Add amazing feature'`)
3. Push to the branch (`git push origin feature/amazing-feature`)
4. Open a Pull Request

## 📞 Support

For issues, questions, or suggestions:
- Open an issue on GitHub
- Check existing documentation in `/prompt` directory
- Review test cases for usage examples

---

**Last Updated**: April 5, 2026
