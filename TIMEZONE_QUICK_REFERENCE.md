# Timezone Implementation - Quick Reference

## TL;DR - How It Works

The system now properly handles bookings across different timezones using a **UTC overlap strategy**:

```
User's Date + User's Timezone → Convert to UTC Range
                                    ↓
Coach's Availability + Coach's Timezone → Convert each time to UTC
                                    ↓
Find Times That Appear in Both Ranges = Available Slots
```

## Key Functions

### 1. Getting Available Slots

```go
// Service method
slots, err := availabilityService.GetAvailableSlots(ctx, coachID, "2024-04-13", "Asia/Kolkata")

// What it does:
// 1. Converts user's date to UTC start/end times in their timezone
// 2. Fetches coach's availability rules  
// 3. For each rule, generates slots that overlap with user's date range
// 4. Removes already-booked slots
// 5. Returns available slots in UTC
```

### 2. Creating a Booking

```go
booking, err := bookingService.CreateBooking(ctx, &dto.CreateBookingRequest{
    UserID:    1,
    CoachID:   5,
    StartTime: "2024-04-13T18:00:00Z",  // RFC3339 format (UTC)
    IdempotencyKey: ptr("unique-key"),
})

// The service validates:
// 1. Coach exists and has timezone set
// 2. Requested time converts properly to coach's timezone
// 3. No conflicts with existing bookings
// 4. Time is in future
```

## Important Things to Know

### ✓ DO
- Always use RFC3339 format for times (UTC)
- Always pass timezone strings (e.g., "America/New_York")
- Use Go's `time.LoadLocation()` for timezone handling
- Store times in database as UTC
- Let Go handle DST automatically

### ✗ DON'T
- Use fixed UTC offsets (UTC-5, UTC+5:30)
- Do manual time arithmetic across timezones
- Assume same day in different timezones
- Forget to handle day boundary shifts

## Example Scenarios

### Scenario 1: Same Timezone
```
Coach: America/Los_Angeles (UTC-7)
User: America/Los_Angeles (UTC-7)
User requests: April 13, 10:00 AM

→ Results in slots around 10:00 AM PDT
```

### Scenario 2: Different Timezones
```
Coach: America/New_York (UTC-4 EDT)
User: Europe/London (UTC+1 BST)

When user requests April 13, 4:00 PM BST:
→ UTC: April 13, 3:00 PM UTC
→ Coach sees: April 13, 11:00 AM EDT
→ Coach availability checked for 11:00 AM EDT times
```

### Scenario 3: Day Shift
```
Coach: America/Los_Angeles (UTC-7)
User: Asia/Tokyo (UTC+9)

When user requests April 13, 8:00 AM JST:
→ UTC: April 12, 11:00 PM UTC (previous day!)
→ Coach sees: April 12, 4:00 PM PDT

So checking "April 13" for user actually checks April 12-13 for coach
```

## Available Timezone Functions

### Timezone Helper Methods

```go
tzHelper := utils.NewTimezoneHelper()

// Convert user's date to UTC range
startUTC, endUTC, _ := tzHelper.GetDayRangeInUTC("2024-04-13", "Asia/Kolkata")

// Get day-of-week accounting for timezone
dayOfWeek := tzHelper.GetDayOfWeekInTimezone(utcTime, "America/New_York")

// Check if time is in coach's working hours
isAvailable := tzHelper.IsTimeInCoachAvailability(utcTime, "America/New_York", availabilityRule)

// Convert between timezones
coachTime, _ := tzHelper.ConvertTimeToCoachTimezone(someTime, "Asia/Kolkata", "America/New_York")
utcTime, _ := tzHelper.ConvertTimeToUserTimezone(someUtcTime, "Asia/Kolkata")
```

## Common Issues & Solutions

### Issue: "Slot not found for available time"
**Cause:** User's date doesn't convert correctly to UTC
**Solution:** Verify timezone string is valid IANA name (e.g., "Asia/Kolkata", not "IST")

### Issue: "DST causes hour difference"
**Cause:** Using fixed UTC offset instead of timezone name
**Solution:** Use `time.LoadLocation("America/New_York")` instead of UTC-5

### Issue: "Day boundary seems wrong"
**Cause:** Not accounting for timezone shifts
**Solution:** Fetch availability for prev/current/next day (service does this automatically)

### Issue: "Slot generation too slow"
**Cause:** Too many availability rules or wide date range
**Solution:** Fetch only 3 consecutive days of availability (already optimized)

## Database Schema

```sql
-- Users table
ALTER TABLE users ADD COLUMN timezone VARCHAR(64) DEFAULT 'UTC';

-- Coaches table
ALTER TABLE coaches ADD COLUMN timezone VARCHAR(64) DEFAULT 'UTC';

-- Bookings table (no changes - already UTC)
-- start_time and end_time are stored as DATETIME in UTC

-- Availability table (no changes)
-- day_of_week (0-6)
-- start_time (HH:MM format, interpreted in coach's timezone)
-- end_time (HH:MM format, interpreted in coach's timezone)
```

## API Changes

### Get Available Slots

```http
GET /api/v1/coaches/:coachID/slots?date=2024-04-13&timezone=Asia/Kolkata
```

**Query Parameters:**
- `date` (required): YYYY-MM-DD format
- `timezone` (required): IANA timezone name

**Response:**
```json
{
  "slots": [
    {
      "start_time": "2024-04-13T18:00:00Z",
      "end_time": "2024-04-13T18:30:00Z"
    }
  ]
}
```

Note: Times are in UTC, but only include slots that fall within user's requested date in their timezone.

### Create Booking

```http
POST /api/v1/bookings
Content-Type: application/json

{
  "user_id": 1,
  "coach_id": 5,
  "start_time": "2024-04-13T18:00:00Z",
  "idempotency_key": "unique-key-123"
}
```

**Requirements:**
- `start_time` must be in RFC3339 format (UTC)
- Must align to 30-minute boundaries
- End time is automatically calculated (start + 30 minutes)

## Testing Tips

1. **Test with real timezones** - Use actual IANA names
2. **Test day boundaries** - Request dates near midnight in different timezones
3. **Test DST transitions** - Request around March 10 and November 3 in US timezones
4. **Test fractional offsets** - Use India (UTC+5:30) and Nepal (UTC+5:45)
5. **Test old bookings** - Ensure existing data still works with UTC times

## Performance

| Operation | Complexity | Notes |
|-----------|-----------|-------|
| Get slots | O(48 * R) | R = number of availability rules, 48 = max slots per day |
| Timezone conversion | O(1) | Built-in Go function |
| Database query | O(log n) | Uses indexed queries on coach_id, start_time |
| Slot filtering | O(b) | b = number of bookings in range |

**Total time:** Sub-millisecond for typical cases

## Debugging

Enable detailed logging:
```go
log := logger.GetLogger()
log.Debug("Converting user date to UTC", "date", "2024-04-13", "timezone", "Asia/Kolkata")
```

Check timezone validity:
```go
loc, err := time.LoadLocation("America/New_York")
if err != nil {
    log.Error("Invalid timezone", "timezone", "America/New_York", "error", err)
}
```

Verify UTC conversion:
```go
startUTC, endUTC, _ := tzHelper.GetDayRangeInUTC("2024-04-13", "Asia/Kolkata")
log.Debug("UTC range", "startUTC", startUTC, "endUTC", endUTC)
```

## Related Documentation

- [TIMEZONE_IMPLEMENTATION.md](./TIMEZONE_IMPLEMENTATION.md) - Complete implementation guide
- [TIMEZONE_SESSION_SUMMARY.md](./TIMEZONE_SESSION_SUMMARY.md) - Session summary
- Go Time Package Docs - https://pkg.go.dev/time (built-in IANA database)

## Quick Links

- **Service:** `internal/service/availability_service.go` - GetAvailableSlots()
- **Utils:** `pkg/utils/timezone.go` - All timezone functions
- **Models:** `internal/model/` - Coach and User models with timezone fields
- **Tests:** (To be created) - Comprehensive timezone test suite

---

**Status:** ✅ Implementation Complete  
**Last Updated:** April 5, 2024  
**Next:** Integration tests and deployment
