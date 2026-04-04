# Architecture & Design Documentation

## System Overview

The Booking API Service is a production-grade REST API for managing appointment booking with availability management. It follows clean architecture principles with clear separation of concerns.

## Architecture Layers

### 1. Handler Layer (`internal/handler/`)

**Responsibility**: HTTP request/response handling

- **AvailabilityHandler**: Manages availability-related endpoints
  - Set weekly availability
  - Add/override exceptions
  - Get available slots

- **BookingHandler**: Manages booking-related endpoints
  - Create bookings
  - Get user bookings
  - Modify bookings
  - Cancel bookings

**Key Points**:
- Input validation at HTTP level
- Request parsing and serialization
- Response formatting
- HTTP status code management

### 2. Service Layer (`internal/service/`)

**Responsibility**: Business logic and orchestration

- **AvailabilityService**: 
  - Validates coach existence
  - Manages weekly availability (CRUD)
  - Handles exceptions and overrides
  - Generates available slots dynamically
  - Applies timezone conversions

- **BookingService**:
  - Validates user/coach existence
  - Creates bookings with concurrency safety
  - Handles idempotency via idempotency_key
  - Validates time boundaries and slot availability
  - Manages booking modifications and cancellations

**Key Points**:
- No direct HTTP knowledge
- Implements business rules
- Uses repository interfaces (testable)
- Handles complex validations

### 3. Repository Layer (`internal/repository/`)

**Responsibility**: Data access and persistence

- **UserRepository**: User CRUD operations
- **CoachRepository**: Coach CRUD operations
- **AvailabilityRepository**: Weekly availability management
- **AvailabilityExceptionRepository**: One-time overrides
- **BookingRepository**: Booking operations

**Key Points**:
- Database-agnostic interfaces
- Query execution
- Soft deletes support
- No business logic

### 4. Model Layer (`internal/model/`)

**Responsibility**: Domain entities

- User, Coach, Availability, AvailabilityException, Booking
- Constants (BookingStatus, DayOfWeek)
- Database schema representation

### 5. DTO Layer (`internal/dto/`)

**Responsibility**: Request/Response contracts

- Input validation structures
- API response objects
- Separates internal models from API contracts

## Concurrency & Safety

### Double Booking Prevention

**Strategy**: Layered defense

1. **Application Layer**: Slot availability check
   ```go
   SELECT COUNT(*) FROM bookings 
   WHERE coach_id = ? AND status = ? 
   AND start_time < ? AND end_time > ? FOR UPDATE
   ```

2. **Database Level**: UNIQUE constraint
   ```sql
   UNIQUE KEY unique_coach_start_time (coach_id, start_time)
   ```

### Transaction Safety

```go
tx, err := s.db.BeginTx(ctx, nil)
defer tx.Rollback()

// Check slot availability with row lock
// Insert booking
// Commit transaction
```

**How it works**:
- `FOR UPDATE` locks the row during check
- Atomic insert within transaction
- UNIQUE constraint catches race conditions
- Returns conflict error if booking exists

### Idempotency

```go
if idempotencyKey != nil {
    existingBooking, _ := s.bookingRepo.GetByIdempotencyKey(ctx, *key)
    if existingBooking != nil {
        return existingBooking  // Idempotent response
    }
}
```

## Slot Generation Algorithm

### Input
- Coach ID
- Date (YYYY-MM-DD)
- Timezone

### Process

```
1. Get day of week from date
2. Check for exceptions
   - If full day disabled: return empty
   - If partial override: use custom times
   - If no exception: use weekly availability
3. Generate 30-minute slots
   - For each availability window
   - Create slots at 30-min boundaries
4. Get all active bookings for that day
5. Remove booked slots from available slots
6. Convert to requested timezone
7. Return available slots
```

### Example

**Weekly Availability** (Monday): 10:00-13:00, 16:00-18:00

**Generated Slots**:
- 10:00-10:30, 10:30-11:00, 11:00-11:30, 11:30-12:00, 12:00-12:30, 12:30-13:00
- 16:00-16:30, 16:30-17:00, 17:00-17:30, 17:30-18:00

**With 1 Booking** (11:00-11:30):
- 10:00-10:30, 10:30-11:00, ~~11:00-11:30~~, 11:30-12:00, 12:00-12:30, 12:30-13:00
- 16:00-16:30, 16:30-17:00, 17:00-17:30, 17:30-18:00

## Timezone Handling

### Storage
- All times stored in **UTC** in database
- Each user/coach has a `timezone` field (IANA format, e.g., "Asia/Kolkata")

### API Boundary
- Requests: Accept RFC3339 format (can be any timezone)
- Responses: Return UTC timestamps
- Slot requests: Accept timezone query param

### Example
```
Coach timezone: Asia/Kolkata (UTC+5:30)
Request: GET /coaches/1/slots?date=2026-04-10&timezone=Asia/Kolkata

Internal Processing:
1. Parse date in requested timezone
2. Calculate day of week
3. Generate slots in UTC
4. Convert display to requested timezone (optional client-side)
5. Return UTC timestamps
```

## Validation Rules

### Time Format
- Must use HH:MM format for daily times
- Must use RFC3339 for bookings
- Must align to 30-minute boundaries

### Business Rules
- `start_time < end_time`
- Cannot book past slots
- Cannot modify ongoing bookings
- Cannot create overlapping bookings

### Soft Deletes
- Records marked with `deleted = TRUE`
- Queries exclude deleted records
- Maintains data integrity

## Error Handling

### HTTP Status Codes
- **400**: Invalid input (bad format, invalid data)
- **409**: Conflict (slot already booked, double booking)
- **500**: Server error (database, validation)

### Error Types
```
1. Validation Errors: Return 400
2. Not Found: Return 400 ("resource not found")
3. Conflict: Return 409 (double booking)
4. Database: Return 500
```

## Performance Considerations

### Indexes
```sql
-- Coaches
INDEX idx_deleted (deleted)

-- Availability
INDEX idx_coach_id (coach_id)
INDEX idx_deleted (deleted)

-- Bookings
INDEX idx_user_id (user_id)
INDEX idx_coach_id (coach_id)
INDEX idx_start_time (start_time)
INDEX idx_status (status)
INDEX idx_deleted (deleted)

-- Unique constraints
UNIQUE KEY unique_coach_start_time (coach_id, start_time)
```

### Query Optimization
- Filtered by coach_id (indexed)
- Range queries on start_time (indexed)
- Status filtering (indexed)

### Slot Generation
- Generated dynamically (not cached)
- Fetches only necessary data
- Filters before returning

## Logging

### Framework: Zap (structured logging)

```go
log.Error("failed to create booking")
log.Info("starting server on :8080")
```

### Log Levels
- ERROR: Critical failures
- INFO: Important operations
- DEBUG: Detailed debugging (in dev mode)

## Testing Strategy

### Unit Tests
- Time utility functions
- Validation logic
- Timezone conversions

### Integration Tests (future)
- End-to-end booking flow
- Concurrency scenarios
- Database transactions
- Slot generation with real data

### Test Tools
- Standard `testing` package
- Table-driven tests
- Mocked repositories (dependency injection)

## Extension Points

### Adding New Features
1. Create DTOs in `internal/dto/`
2. Add repository methods if needed
3. Implement service logic
4. Create handlers
5. Add routes in `main.go`

### Database Migration
- Create `.up.sql` and `.down.sql` files
- Use `golang-migrate` to apply
- Maintain backward compatibility

### Changing Models
- Update repository queries
- Update service logic
- Add migration if schema changes
- Update tests

## Security Considerations

### Input Validation
- All user inputs validated
- SQL injection prevention (parameterized queries)
- Time format validation

### Database Constraints
- UNIQUE constraints prevent duplicates
- Foreign keys maintain referential integrity
- Soft deletes prevent accidental loss

### Future Enhancements
- Rate limiting
- Authentication/authorization
- API key management
- Request signing
- Audit logging

## Monitoring & Observability

### Current
- Structured logging with Zap
- Error logging with context

### Future
- Metrics (Prometheus)
- Distributed tracing (Jaeger)
- Health checks
- Performance monitoring
