# Booking API - Postman Collection Guide

## Overview

Complete Postman collection for testing the Booking API Service with **80+ ready-to-use requests** utilizing the comprehensive seed data.

**File:** `Booking_API_Postman_Collection.json`

---

## Setup Instructions

### 1. Import Collection into Postman

1. Open Postman
2. Click **Import** → Select **Booking_API_Postman_Collection.json**
3. Collection will be imported with all organized folders and requests

### 2. Set Environment Variables

The collection uses a variable `{{baseUrl}}` which defaults to `http://localhost:8080`

**To modify:**
- Open the collection → Click **Variables** tab
- Change `baseUrl` value to your target environment
- Examples:
  - Development: `http://localhost:8080`
  - Production: `https://api.booking-service.com`

### 3. Start Your API Server

```bash
cd /path/to/booking-api-service
make run
# or
go run cmd/main.go
```

The server should be running on `http://localhost:8080`

---

## Collection Structure

### 📋 Folder Organization

#### 1. **Health & Status**
Basic health check to verify API is running.

| Request | Method | Purpose |
|---------|--------|---------|
| Health Check | GET | Verify service is running |

---

#### 2. **Availability - Single Day**
Set availability for individual days (legacy/single-use).

| Request | Coach | Day | Hours |
|---------|-------|-----|-------|
| Set Coach Alice (ID:1) - Monday 9-5 | Alice (1) | Monday | 09:00-17:00 |
| Set Coach Bob (ID:2) - Wednesday with lunch break | Bob (2) | Wednesday | 10:00-12:30, 13:30-18:00 |
| Set Coach Diana (ID:4) - Thursday early hours | Diana (4) | Thursday | 07:00-15:00 |

---

#### 3. **Availability - Weekly Setup**
Set availability for entire week in single request (new feature).

| Request | Coach | Pattern | Days Covered |
|---------|-------|---------|--------------|
| Set Weekly Availability - Coach Alice (Full Week Mon-Fri) | Alice (1) | 9-5 Standard | Mon-Fri |
| Set Weekly Availability - Coach Edward (Flexible Schedule) | Edward (5) | Split shifts | Tue-Sat with breaks |
| Set Weekly Availability - Coach Fiona (All Days) | Fiona (6) | Maximum availability | Every day |

**Key Feature:** These requests demonstrate the new weekly availability endpoint that allows setting multiple days in one request.

---

#### 4. **Availability - Exceptions**
Handle special cases like days off, special availability, partial days.

| Request | Coach | Date | Type | Details |
|---------|-------|------|------|---------|
| Add Exception - Coach Alice Day Off (Full Day) | Alice (1) | 2026-04-12 | Day off | is_available: false |
| Add Exception - Coach Edward Special Evening Session | Edward (5) | 2026-04-10 | Special hours | 19:00-21:00 |
| Add Exception - Coach Fiona Partial Day Availability | Fiona (6) | 2026-04-13 | Partial day | 10:00-13:00 only |

---

#### 5. **Available Slots**
Query available slots for coaches across different timezones.

| Request | Coach | Date | Timezone | Purpose |
|---------|-------|------|----------|---------|
| Get Slots - Coach Alice (NY Timezone) | Alice (1) | 2026-04-06 | America/New_York | Standard query |
| Get Slots - Coach Bob (London Timezone) | Bob (2) | 2026-04-07 | Europe/London | Extended hours |
| Get Slots - Coach Edward (Berlin Timezone) | Edward (5) | 2026-04-08 | Europe/Berlin | Multiple slots |
| Get Slots - Coach Fiona (Melbourne Timezone) | Fiona (6) | 2026-04-06 | Australia/Melbourne | Max availability |
| Get Slots - Coach George (Chicago Timezone) | George (7) | 2026-04-09 | America/Chicago | Limited slots |

**Test Timezone Conversions:** Each request uses the coach's native timezone to test accurate slot generation across timezones.

---

#### 6. **Bookings**
Create bookings and retrieve booking history.

| Request | User | Coach | Date-Time | Purpose |
|---------|------|-------|-----------|---------|
| Create Booking - User John with Coach Alice | John (1) | Alice (1) | 2026-04-06 09:30 | Standard booking |
| Create Booking - User Emma with Coach Diana | Emma (6) | Diana (4) | 2026-04-07 11:00 | Early hours |
| Create Booking - User Alex with Coach Edward | Alex (7) | Edward (5) | 2026-04-08 14:30 | Multiple slots |
| Create Booking - User Lisa with Coach Fiona | Lisa (8) | Fiona (6) | 2026-04-06 14:30 | Max availability |
| Get Bookings - User John | John (1) | - | - | View all bookings |
| Get Bookings - User Emma | Emma (6) | - | - | View all bookings |
| Get Bookings - User Alex | Alex (7) | - | - | View all bookings |

---

#### 7. **Bookings - Modify & Cancel**
Modify existing bookings and cancel them.

| Request | Booking ID | Action | New Time |
|---------|------------|--------|----------|
| Modify Booking - Change time (Booking ID:1) | 1 | Modify | 2026-04-06 10:00 |
| Modify Booking - Change time (Booking ID:5) | 5 | Modify | 2026-04-10 16:00 |
| Cancel Booking - Booking ID:15 | 15 | Cancel | - |
| Cancel Booking - Booking ID:16 | 16 | Cancel | - |

---

#### 8. **Error Scenarios & Edge Cases**
Test error handling and validation.

| Request | Scenario | Expected Error |
|---------|----------|-----------------|
| ERROR - Invalid Coach ID | Non-existent coach (999) | 404 Coach not found |
| ERROR - Invalid Day of Week | Day of week = 7 (invalid) | 400 Invalid day |
| ERROR - Invalid Time Format | Time not on 30-min boundary | 400 Time alignment error |
| ERROR - Invalid User ID for Booking | Non-existent user (999) | 404 User not found |
| ERROR - Start time after end time | End before start | 400 Invalid time range |

---

## Seed Data Reference

### Users (10 Total)
```
1. John Doe (America/New_York)
2. Jane Smith (Europe/London)
3. Mike Johnson (Asia/Tokyo)
4. Sarah Williams (Australia/Sydney)
5. David Brown (America/Los_Angeles)
6. Emma Davis (America/Chicago)
7. Alex Martinez (Europe/Paris)
8. Lisa Anderson (Asia/Singapore)
9. Tom Wilson (Canada/Toronto)
10. Rachel Green (America/Denver)
```

### Coaches (8 Total)
```
1. Coach Alice (America/New_York) - Mon-Fri 9-5
2. Coach Bob (Europe/London) - Tue-Sat 10-18 (with lunch break)
3. Coach Charlie (Asia/Tokyo) - Mon/Wed/Fri 8-16
4. Coach Diana (America/Los_Angeles) - Mon-Thu 7-15
5. Coach Edward (Europe/Berlin) - Tue-Sat flexible
6. Coach Fiona (Australia/Melbourne) - Daily 8-18
7. Coach George (America/Chicago) - Wed-Thu limited
8. Coach Hannah (Asia/Singapore) - Daily multiple slots
```

### Existing Bookings (17 Total)
- **IDs 1-10:** ACTIVE bookings (test retrieval and modification)
- **IDs 11-14:** COMPLETED bookings (test status filtering)
- **IDs 15-17:** CANCELLED bookings (test cancellation)

---

## Testing Workflows

### Workflow 1: Complete Booking Flow
1. ✅ **Health Check** - Verify API running
2. ✅ **Set Weekly Availability** - Coach Alice (Mon-Fri)
3. ✅ **Get Slots** - Coach Alice on April 6
4. ✅ **Create Booking** - User John with Coach Alice
5. ✅ **Get Bookings** - User John (verify booking created)
6. ✅ **Modify Booking** - Change time
7. ✅ **Get Bookings** - User John (verify modification)

### Workflow 2: Timezone Testing
1. ✅ **Get Slots** - Coach Alice (New York)
2. ✅ **Get Slots** - Coach Bob (London)
3. ✅ **Get Slots** - Coach Edward (Berlin)
4. ✅ **Get Slots** - Coach Fiona (Melbourne)
5. Compare returned times to verify timezone conversions

### Workflow 3: Exception Handling
1. ✅ **Add Exception** - Coach Alice day off
2. ✅ **Get Slots** - Coach Alice (should have no slots)
3. ✅ **Add Exception** - Coach Edward special session
4. ✅ **Get Slots** - Coach Edward (should show evening slots)

### Workflow 4: Weekly Availability
1. ✅ **Set Weekly Availability** - Full week setup
2. ✅ **Get Slots** - Multiple days to verify
3. ✅ **Verify** - All 5 days have correct availability
4. ✅ **Modify** - Change one day
5. ✅ **Verify** - Changes persisted

### Workflow 5: Error Testing
1. ✅ Run all **Error Scenarios** requests
2. ✅ Verify proper error messages
3. ✅ Verify correct HTTP status codes

---

## Request Details & Usage

### Headers Used
- `Content-Type: application/json` - For POST/PUT requests
- `X-Request-ID` - Unique request identifier for logging/tracing

### Variables in Requests
- `{{baseUrl}}` - Base URL (defaults to http://localhost:8080)
- Coach IDs: 1-8 (from seed data)
- User IDs: 1-10 (from seed data)
- Booking IDs: 1-17 (from seed data)

### Response Examples

**Successful Availability Set:**
```json
{
  "message": "weekly availability set successfully"
}
```

**Available Slots Response:**
```json
{
  "slots": [
    {
      "start_time": "2026-04-06T09:00:00Z",
      "end_time": "2026-04-06T09:30:00Z"
    },
    {
      "start_time": "2026-04-06T09:30:00Z",
      "end_time": "2026-04-06T10:00:00Z"
    }
  ]
}
```

**Booking Created:**
```json
{
  "id": 18,
  "user_id": 1,
  "coach_id": 1,
  "start_time": "2026-04-06T09:30:00Z",
  "end_time": "2026-04-06T10:00:00Z",
  "status": "ACTIVE",
  "created_at": "2026-04-05T10:15:30Z",
  "updated_at": "2026-04-05T10:15:30Z"
}
```

---

## Tips & Best Practices

### 1. **Before Running Requests**
- Ensure database is populated with seed data
- Server is running on correct port
- `baseUrl` variable is set correctly

### 2. **Request Ordering**
- Health Check first to verify connectivity
- Set availability before getting slots
- Create bookings only on available slots
- Test errors last

### 3. **Timezone Testing**
- Each coach has different timezone
- When getting slots, use coach's timezone
- Verify returned times account for timezone offsets

### 4. **Idempotency Keys**
- Each booking request includes unique `idempotency_key`
- Allows safe retries without duplicate bookings
- Modify keys if re-running same request

### 5. **Response Validation**
- Check HTTP status codes (200, 201, 400, 404, 500)
- Validate response structure matches schema
- Verify data returned matches request parameters

---

## Common Issues & Solutions

### Issue: "Coach not found"
- **Cause:** Coach ID doesn't exist (valid range: 1-8)
- **Solution:** Use IDs 1-8 from seed data

### Issue: "times must align to 30-minute boundaries"
- **Cause:** Time like "09:15" or "10:45" used
- **Solution:** Use only :00 or :30 (e.g., 09:00, 09:30, 10:00)

### Issue: "No available slots"
- **Cause:** No availability set or exceptions override
- **Solution:** Set availability first, check exceptions

### Issue: "Slot already booked"
- **Cause:** Another booking exists at same time
- **Solution:** Get slots first, book from available list

### Issue: Connection refused
- **Cause:** API server not running
- **Solution:** Start server with `make run` or `go run cmd/main.go`

---

## Integration with CI/CD

The collection can be used with **Newman** (Postman CLI) for automated testing:

```bash
# Install Newman
npm install -g newman

# Run collection
newman run Booking_API_Postman_Collection.json \
  --environment environment.json \
  --reporters cli,json \
  --reporter-json-export test-results.json
```

---

## Additional Resources

- **Swagger File:** `docs/swagger.yaml`
- **Seed Data Docs:** `SEED_DATA_DOCUMENTATION.md`
- **Weekly Availability Docs:** `WEEKLY_AVAILABILITY_IMPLEMENTATION.md`
- **API Service:** Running on http://localhost:8080

---

## Version History

- **v1.0** - Initial collection with all endpoints and comprehensive test data
