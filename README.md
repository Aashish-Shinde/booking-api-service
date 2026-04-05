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

### Step 5: Create Database Manually

Before running the application, create the database:

```bash
mysql -h localhost -u root -p -e "CREATE DATABASE IF NOT EXISTS booking_api CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
```

### Step 6: Run Database Schema Migrations Manually

The database schema must be applied **manually** using the SQL migration files. The migrations are located in `db/migrations/`:

**Run the migration files in order:**

```bash
# 1. Create initial database structure
mysql -h localhost -u root -p booking_api < db/migrations/000_v1.0.0_create_database.up.sql

# 2. Create tables and constraints
mysql -h localhost -u root -p booking_api < db/migrations/001_v1.0.0_init_schema.up.sql

# 3. Seed test data (coaches, users, availability)
mysql -h localhost -u root -p booking_api < db/migrations/002_v1.0.0_seed_data.up.sql
```

**Verify the migrations were applied:**

```bash
mysql -h localhost -u root -p booking_api -e "SHOW TABLES;"
```

You should see these tables:
- `users` - 10 test users in different timezones
- `coaches` - 8 test coaches with availability
- `availability` - Weekly recurring availability windows
- `availability_exceptions` - One-time overrides (days off, special hours)
- `bookings` - Booking records

**If You Need to Rollback (Undo) Migrations:**

```bash
# Undo in reverse order
mysql -h localhost -u root -p booking_api < db/migrations/002_v1.0.0_seed_data.down.sql
mysql -h localhost -u root -p booking_api < db/migrations/001_v1.0.0_init_schema.down.sql
mysql -h localhost -u root -p booking_api < db/migrations/000_v1.0.0_create_database.down.sql
```

**What Each Migration File Does:**

The migration files are plain SQL scripts that set up your database:

- **000_v1.0.0_create_database.up.sql** - Creates the initial database structure
- **001_v1.0.0_init_schema.up.sql** - Creates all tables with:
  - Foreign key constraints
  - UNIQUE constraints (prevents double-booking)
  - Indexes for query performance
  - Soft delete support (deleted flag)
- **002_v1.0.0_seed_data.up.sql** - Inserts test data:
  - 10 test users in different timezones (Tokyo, London, New York, Paris, etc.)
  - 8 test coaches with different availability patterns
  - Sample availability windows and exceptions
  - Sample bookings in different statuses

After all migrations are applied, you'll have a fully configured database ready for testing!

### Step 7: Run the Application

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

## 🔑 Quick Start: Core Concepts

### What Are Slots?

**30-minute time windows** when a coach is available to book. They're generated automatically from availability windows.

**Example:**
```
Coach Available: 9:00 AM - 12:00 PM
↓
Generated Slots:
  9:00 - 9:30 AM   ✓ Available
  9:30 - 10:00 AM  ✓ Available
  10:00 - 10:30 AM ✓ Available
  10:30 - 11:00 AM ✓ Available
  11:00 - 11:30 AM ✓ Available
  11:30 - 12:00 PM ✓ Available
```

### What Is Availability?

Availability is when a coach can accept bookings. It comes in two forms:

**1. Weekly Recurring Availability**
```json
{
  "day_of_week": 1,        // Monday
  "start_time": "09:00",   // Coach's local timezone
  "end_time": "17:00"      // Coach's local timezone
}
```

**2. One-Time Exceptions**
```json
{
  "date": "2026-04-10",
  "is_available": false     // Day off
}

OR

{
  "date": "2026-04-12",
  "is_available": true,
  "start_time": "10:00",
  "end_time": "14:00"       // Special hours on this day
}
```

### Booking a Slot: The Complete Flow

**1. User requests available slots:**
```bash
GET /api/availability/coaches/1/slots?date=2026-04-06&userID=1
```

**2. System performs UTC overlap matching:**
```
User's Date (April 6 in Tokyo):     April 5, 3 PM UTC → April 6, 3 PM UTC
Coach's Availability (Monday 9-5 NY): April 7, 1 PM UTC → April 7, 9 PM UTC

Result: Check what PART of the user's date range overlaps with coach's availability
```

**3. Generate 30-minute slots from overlap:**
```json
{
  "slots": [
    { "start_time": "2026-04-06T14:00:00Z", "end_time": "2026-04-06T14:30:00Z" },
    { "start_time": "2026-04-06T14:30:00Z", "end_time": "2026-04-06T15:00:00Z" }
  ]
}
```

**4. User books one of these slots:**
```bash
POST /api/bookings
{
  "user_id": 1,
  "coach_id": 1,
  "start_time": "2026-04-06T14:00:00Z"
}
```

**5. System validates:**
- ✓ Time is in RFC3339 UTC format
- ✓ Time aligns to 30-minute boundary
- ✓ Coach hasn't already been booked for this slot
- ✓ Time is in the future

**6. Booking created:**
```json
{
  "id": 1,
  "user_id": 1,
  "coach_id": 1,
  "start_time": "2026-04-06T14:00:00Z",
  "end_time": "2026-04-06T14:30:00Z",
  "status": "ACTIVE"
}
```

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

The system supports full timezone awareness with a **UTC overlap strategy**:

- **Storage**: All times stored in UTC in the database
- **Coach Timezone**: Coach's availability is set in their local timezone
- **User Timezone**: Slots are converted to user's timezone in API response
- **Automatic Conversion**: Handled transparently at API boundaries

### Deep Dive: Time Conversion Algorithm

The system uses a **UTC overlap matching** approach to handle timezone-aware bookings across different timezones:

#### Step 1: User's Date → UTC Range

When a user requests slots for a date in their timezone:

```
User: April 13, 2026 in Asia/Tokyo (UTC+9)
↓
Start: April 12, 2026 3:00 PM UTC (April 13, 12:00 AM Tokyo)
End:   April 13, 2026 3:00 PM UTC (April 14, 12:00 AM Tokyo)
```

The system converts the user's requested date **boundaries** to UTC. This means checking availability from midnight to midnight in the user's timezone, which translates to a different UTC range.

#### Step 2: Coach's Availability → UTC Times

Coach's weekly availability (e.g., "Monday 9 AM - 5 PM in America/New_York") is stored as local times. For each day within the UTC range:

```
Coach Availability: Monday 9:00 AM - 5:00 PM (America/New_York, UTC-4 EDT)
↓
For April 14 (Monday) in UTC:
  → Monday in New York
  → 9:00 AM EDT = 1:00 PM UTC
  → 5:00 PM EDT = 9:00 PM UTC
  → Available UTC slots: April 14, 1:00 PM - 9:00 PM UTC
```

The system converts the coach's local availability times to UTC.

#### Step 3: Find Overlapping Times

Compare the user's UTC date range with the coach's UTC availability:

```
User's UTC Range:        April 12, 3 PM — April 13, 3 PM UTC
Coach's UTC Availability: April 14, 1 PM — April 14, 9 PM UTC

Overlap: None! Coach not available during user's date in UTC.

---

User's UTC Range:        April 13, 3 PM — April 14, 3 PM UTC  
Coach's UTC Availability: April 14, 1 PM — April 14, 9 PM UTC

Overlap: April 14, 1 PM — April 14, 3 PM UTC → 2 hours available
```

Only the **overlapping times** are returned as available slots.

#### Step 4: Generate 30-Minute Slots

From the overlapping UTC times, create 30-minute slots:

```
UTC Overlap: April 14, 1:00 PM - 3:00 PM UTC
↓
Slots:
  1. April 14, 1:00 PM - 1:30 PM UTC
  2. April 14, 1:30 PM - 2:00 PM UTC
  3. April 14, 2:00 PM - 2:30 PM UTC
  4. April 14, 2:30 PM - 3:00 PM UTC
```

### Example: Real-World Timezone Scenario

**Setup:**
- **Coach**: Alice in America/New_York (UTC-4 EDT on April 14)
  - Available: Monday-Friday, 9:00 AM - 5:00 PM EDT
- **User**: Bob in Asia/Tokyo (UTC+9 JST)
  - Requests: Available slots for April 13, 2026

**Calculation:**

1. **User's Date to UTC:**
   ```
   Bob's April 13 in Tokyo (9 AM Sunday → 9 AM Monday in UTC)
   = April 12, midnight JST → April 13, midnight JST (UTC)
   = April 12, 3:00 PM UTC → April 13, 3:00 PM UTC
   ```

2. **Coach's Availability on Mondays:**
   ```
   Alice's Monday 9 AM - 5 PM EDT
   = April 14, 9 AM EDT → April 14, 5 PM EDT
   = April 14, 1:00 PM UTC → April 14, 9:00 PM UTC
   ```

3. **Check Overlap:**
   ```
   User's UTC Range:      April 12, 3 PM — April 13, 3 PM UTC
   Coach's Availability:  April 14, 1 PM — April 14, 9 PM UTC
   
   Result: NO OVERLAP (Coach's Monday is AFTER user's date range)
   ```

4. **User Gets:**
   ```
   Empty slots array (Alice not available during Bob's April 13)
   ```

**But if Bob requests April 14:**

1. **User's Date to UTC:**
   ```
   Bob's April 14 in Tokyo
   = April 13, 3:00 PM UTC → April 14, 3:00 PM UTC
   ```

2. **Check Overlap:**
   ```
   User's UTC Range:      April 13, 3 PM — April 14, 3 PM UTC
   Coach's Availability:  April 14, 1 PM — April 14, 9 PM UTC
   
   Overlap: April 14, 1 PM — April 14, 3 PM UTC
   ```

3. **Generate Slots:**
   ```
   April 14, 1:00 PM UTC - 1:30 PM UTC
   April 14, 1:30 PM UTC - 2:00 PM UTC
   April 14, 2:00 PM UTC - 2:30 PM UTC
   April 14, 2:30 PM UTC - 3:00 PM UTC
   
   In Bob's timezone (Tokyo):
   April 14, 10:00 PM JST - April 14, 10:30 PM JST
   April 14, 10:30 PM JST - April 15, 12:00 AM JST
   April 15, 12:00 AM JST - April 15, 12:30 AM JST
   April 15, 12:30 AM JST - April 15, 1:00 AM JST
   ```

**Key Insight:** Alice (in New York) is available Monday afternoon, but for Bob (in Tokyo) requesting Monday, Alice's Monday is already Tuesday in Tokyo! This is why day boundaries shift across timezones.

### Slot Generation Details

**30-Minute Slot System:**

Slots are generated dynamically from availability windows:

1. **Input:** Coach's availability window (e.g., 9:00 AM - 5:00 PM)
2. **Process:** Divide into contiguous 30-minute intervals
3. **Output:** All 30-minute slots that don't conflict with bookings

**Example:**

```
Availability: 9:00 AM - 12:00 PM (3 hours)
↓
Generated Slots:
  9:00 AM - 9:30 AM   (30 min)
  9:30 AM - 10:00 AM  (30 min)
  10:00 AM - 10:30 AM (30 min)
  10:30 AM - 11:00 AM (30 min)
  11:00 AM - 11:30 AM (30 min)
  11:30 AM - 12:00 PM (30 min)
```

**Slot Exclusions:**

- Slots already booked (UNIQUE constraint on coach_id + start_time)
- Slots outside availability windows
- Slots in the past
- Slots on days with exceptions (days off)

### Important Considerations

#### Day Boundary Shifts

When coaches and users are in very different timezones, the user's requested date may span multiple days for the coach:

```
User requests "April 13" from UTC+12 (West)
= Coach sees April 12 to April 14 in UTC-12 (East)
```

#### DST (Daylight Saving Time)

Go automatically handles DST transitions when using IANA timezone identifiers:

```go
loc, _ := time.LoadLocation("America/New_York")
// During EDT (April): UTC-4
// During EST (January): UTC-5
// Transition handled automatically
```

#### Always Use RFC3339 Format

When creating bookings, always use UTC times in RFC3339 format:

```json
{
  "start_time": "2026-04-14T14:00:00Z"
}
```

The `Z` indicates UTC. Never use local times or fixed offsets.

#### Supported Timezones

The system supports all IANA timezone identifiers. Common examples:

- **North America**: America/New_York, America/Chicago, America/Los_Angeles, America/Denver, Canada/Toronto
- **Europe**: Europe/London, Europe/Paris, Europe/Berlin, Europe/Madrid
- **Asia**: Asia/Tokyo, Asia/Shanghai, Asia/Hong_Kong, Asia/Singapore, Asia/Bangkok, Asia/Kolkata
- **Australia**: Australia/Sydney, Australia/Melbourne, Australia/Brisbane
- **UTC**: UTC

### Testing Timezone Logic

Use the test data provided in the seed:

**Coach Alice** (America/New_York, 9 AM - 5 PM):
- Test with User 1 (America/New_York) - Same timezone
- Test with User 3 (Asia/Tokyo) - 13-14 hour difference

**Coach Fiona** (Australia/Melbourne, 8 AM - 6 PM):
- Test with User 7 (Europe/Paris) - 8-9 hour difference

**Test Query:**
```bash
curl -X GET \
  "http://localhost:8080/api/availability/coaches/1/slots?date=2026-04-06&userID=3" \
  -H "Content-Type: application/json"
```

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

### Design Decisions & Rationale

This section documents the key architectural decisions made during the implementation. For detailed discussion and reasoning, see the documentation files in the `/prompt` directory.

#### 1. **UTC Overlap Strategy for Timezone Handling**

**Decision:** Use UTC-based date range overlap matching instead of time offset calculations.

**Rationale:**
- Handles DST (Daylight Saving Time) transitions automatically
- Avoids manual timezone arithmetic errors
- Uses Go's native timezone support via IANA identifiers
- Correctly handles day boundary shifts across timezones

**Alternative Considered:**
- Fixed UTC offset calculations (e.g., UTC-5) - Would break during DST transitions
- Server-side timezone assumption - Doesn't support cross-timezone bookings

**Implementation:**
- Coach stores availability in their local timezone (HH:MM format)
- User requests dates in their local timezone
- System converts both to UTC ranges and finds overlaps
- Generated slots returned in UTC, displayed in user's timezone on client

**Reference:** See `prompt/TIMEZONE_IMPLEMENTATION.md` and `prompt/TIMEZONE_QUICK_REFERENCE.md`

#### 2. **30-Minute Fixed-Duration Slots**

**Decision:** All bookings are fixed 30-minute slots, not flexible durations.

**Rationale:**
- Simplifies availability calculation
- Prevents fragmented availability
- Easier to detect double-booking conflicts
- UNIQUE constraint on (coach_id, start_time) guarantees single booking per slot

**Implementation:**
- `end_time = start_time + 30 minutes` (always)
- Slots generated dynamically from availability windows
- Each availability window divided into contiguous 30-min intervals

#### 3. **Clean Architecture (5-Layer)**

**Decision:** Implement handler → service → repository layered architecture with interfaces.

**Rationale:**
- Clear separation of concerns
- Testable business logic (interfaces allow mocking)
- Reusable service layer
- Easy to add authentication/logging middleware
- Follows Go best practices

**Layers:**
```
HTTP Request
    ↓
[Handler] - HTTP parsing, validation, error→HTTP status conversion
    ↓
[Service] - Business logic, timezone conversion, booking validation
    ↓
[Repository] - Database access, query execution
    ↓
[Model/DTO] - Domain entities and API contracts
    ↓
Database
```

**Reference:** See `prompt/ARCHITECTURE.md`

#### 4. **Database Constraints for Data Integrity**

**Decision:** Enforce business rules at database level with constraints.

**Rationale:**
- UNIQUE(coach_id, start_time) prevents double-booking
- Foreign keys maintain referential integrity
- Soft deletes (deleted flag) support logical deletion without data loss
- Constraints work even if application has bugs

**Constraints:**
```sql
PRIMARY KEY (id)
UNIQUE (coach_id, start_time)           -- Prevent double-booking
FOREIGN KEY (coach_id) REFERENCES coaches(id)
FOREIGN KEY (user_id) REFERENCES users(id)
```

#### 5. **Idempotent Booking Creation**

**Decision:** Support optional `idempotency_key` in booking creation for safe retries.

**Rationale:**
- Handles network failures gracefully
- Clients can retry without creating duplicate bookings
- Essential for production systems with unreliable networks

**How It Works:**
```json
POST /api/bookings
{
  "user_id": 1,
  "coach_id": 1,
  "start_time": "2026-04-06T14:00:00Z",
  "idempotency_key": "booking-20260405-001"
}

// First request → Creates booking, returns 201
// Retry with same key → Returns existing booking, returns 200
// Different key → Creates new booking
```

#### 6. **Transaction-Based Concurrency Safety**

**Decision:** Use database transactions with row-level locking (FOR UPDATE) for booking creation.

**Rationale:**
- Prevents race conditions in concurrent booking
- ACID guarantees at database level
- Multiple layers of protection (transaction + constraint)

**Layers:**
1. **Application Level:** Check availability before attempting booking
2. **Row Level:** FOR UPDATE lock prevents concurrent modifications
3. **Constraint Level:** UNIQUE constraint catches any conflicts

#### 7. **Manual Migration Setup**

**Decision:** Require users to run migrations manually using golang-migrate CLI.

**Rationale:**
- Explicit control over database schema changes
- Clear visibility into migration state
- Safer for production deployments
- Avoids auto-migration surprises

**Setup:**
```bash
# User must install and run manually
go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
migrate -path db/migrations -database "mysql://..." up
```

**Alternative Not Used:**
- Auto-migration on app startup - Can cause data loss if migrations conflict
- Embedded migrations - Less control, harder to debug

#### 8. **No Docker Compose Dependency**

**Decision:** Remove Docker from core setup; support native MySQL or user-provided instance.

**Rationale:**
- Reduces complexity for developers who have MySQL installed
- Deployment doesn't depend on Docker setup
- Better for production environments
- Clear separation: User manages database infrastructure

**Setup:**
```bash
# User manages MySQL independently
# Application just connects via connection string
DATABASE_URL=root:root@tcp(localhost:3306)/booking_api?parseTime=true
```

#### 9. **Structured Logging with Go's Standard log/slog**

**Decision:** Use Go 1.24's built-in `log/slog` package instead of external logging library.

**Rationale:**
- No external dependency for logging
- Structured logging built-in to standard library
- Performance and simplicity
- Future-proof (maintained by Go team)

**Implementation:**
```go
log := logger.GetLogger()
log.Error("failed to create booking", slog.String("error", err.Error()))
```

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

**Error:** `error: Table not found` or `Unknown table 'bookings'`

**Solution:**

Make sure you ran all the migration SQL files:

```bash
# Run migrations in order
mysql -h localhost -u root -p booking_api < db/migrations/000_v1.0.0_create_database.up.sql
mysql -h localhost -u root -p booking_api < db/migrations/001_v1.0.0_init_schema.up.sql
mysql -h localhost -u root -p booking_api < db/migrations/002_v1.0.0_seed_data.up.sql

# Verify tables exist
mysql -h localhost -u root -p booking_api -e "SHOW TABLES;"
```

If you need to reset and re-run:

```bash
# Rollback all migrations (in reverse order)
mysql -h localhost -u root -p booking_api < db/migrations/002_v1.0.0_seed_data.down.sql
mysql -h localhost -u root -p booking_api < db/migrations/001_v1.0.0_init_schema.down.sql
mysql -h localhost -u root -p booking_api < db/migrations/000_v1.0.0_create_database.down.sql

# Then re-apply them
mysql -h localhost -u root -p booking_api < db/migrations/000_v1.0.0_create_database.up.sql
mysql -h localhost -u root -p booking_api < db/migrations/001_v1.0.0_init_schema.up.sql
mysql -h localhost -u root -p booking_api < db/migrations/002_v1.0.0_seed_data.up.sql
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
- Check existing documentation in `/prompt` directory for detailed discussions:
  - `ARCHITECTURE.md` - Detailed architecture and layer explanations
  - `IMPLEMENTATION_SUMMARY.md` - Complete implementation overview
  - `TIMEZONE_IMPLEMENTATION.md` - Deep-dive on timezone handling
  - `TIMEZONE_QUICK_REFERENCE.md` - Quick timezone reference
  - `QUICKSTART.md` - 5-minute setup guide
- Review test cases for usage examples

---

**Last Updated**: April 5, 2026
