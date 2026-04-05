# 30-Minute Slot Enforcement Policy

## Overview

The Booking API strictly enforces **30-minute booking slots**. Users cannot book arbitrary durations or times - all slots must:

1. ✅ Align to 30-minute boundaries (`:00` or `:30`)
2. ✅ Be exactly 30 minutes in duration
3. ✅ Not overlap with existing bookings

---

## Slot Duration

**All bookings are exactly 30 minutes.**

### Valid Examples:
- `09:00` → `09:30` ✅ (30 minutes)
- `09:30` → `10:00` ✅ (30 minutes)
- `14:00` → `14:30` ✅ (30 minutes)
- `23:30` → `00:00` ✅ (30 minutes, next day)

### Invalid Examples:
- `09:00` → `10:00` ❌ (60 minutes - too long)
- `09:00` → `09:15` ❌ (15 minutes - too short)
- `09:15` → `09:45` ❌ (30 minutes but misaligned)
- `09:00` → `09:45` ❌ (45 minutes + misaligned end)

---

## Time Boundary Alignment

All start times MUST align to 30-minute boundaries.

### Valid Start Times:
- `:00` (e.g., 09:00, 10:00, 14:30, 23:00)
- `:30` (e.g., 09:30, 10:30, 14:00, 23:30)

### Invalid Start Times:
- `:15` (e.g., 09:15, 14:15) ❌
- `:45` (e.g., 09:45, 14:45) ❌
- `:05`, `:10`, `:20`, `:25`, `:35`, `:40`, `:50`, `:55` ❌

---

## Validation Layers

### 1. **DTO Validation** (`internal/dto/requests.go`)
Validates at request deserialization:
- Time format correctness (HH:MM)
- 30-minute boundary alignment
- Exactly 30-minute duration

```go
func (t *TimeSlotInput) Validate() error {
    // Validates:
    // - HH:MM format
    // - Aligned to :00 or :30
    // - Exactly 30 minutes duration
}
```

### 2. **Service Validation** (`internal/service/booking_service.go`)
Validates during booking creation/modification:

**CreateBooking:**
```go
// Check if start time is aligned to 30-minute boundaries
if startTime.Minute()%30 != 0 {
    return nil, errors.New("start_time must align to 30-minute boundaries")
}

// Calculate end time (always 30 minutes)
endTime := startTime.Add(30 * time.Minute)
```

**ModifyBooking:**
```go
// Same validation applied when modifying
if newStartTime.Minute()%30 != 0 {
    return nil, errors.New("start_time must align to 30-minute boundaries")
}

newEndTime := newStartTime.Add(30 * time.Minute)
```

### 3. **Database Constraints**
- Start and end times are stored as `DATETIME`
- Unique index on (coach_id, start_time, status) prevents double-booking
- Triggers ensure end_time = start_time + 30 minutes

---

## API Endpoints & Enforcement

### Setting Availability

#### Single Day Endpoint
```
POST /api/v1/coaches/{coach_id}/availability
```

**Request:**
```json
{
  "day_of_week": 1,
  "slots": [
    {
      "start_time": "09:00",
      "end_time": "09:30"
    },
    {
      "start_time": "09:30",
      "end_time": "10:00"
    }
  ]
}
```

✅ Valid - Each slot is exactly 30 minutes and aligned

#### Weekly Endpoint
```
POST /api/v1/coaches/{coach_id}/availability/weekly
```

**Request:**
```json
{
  "availabilities": [
    {
      "day_of_week": 1,
      "slots": [
        {
          "start_time": "09:00",
          "end_time": "09:30"
        }
      ]
    }
  ]
}
```

✅ Valid - Slot is exactly 30 minutes

### Creating Bookings

#### Create Booking Endpoint
```
POST /api/v1/bookings
```

**Request:**
```json
{
  "user_id": 1,
  "coach_id": 1,
  "start_time": "2026-04-06T09:00:00Z"
}
```

**Service Logic:**
1. Validates `start_time` aligns to 30-minute boundary
2. Automatically sets `end_time` = `start_time` + 30 minutes
3. Checks for conflicts with existing bookings
4. Creates booking (always 30 minutes)

**Success Response:**
```json
{
  "id": 18,
  "user_id": 1,
  "coach_id": 1,
  "start_time": "2026-04-06T09:00:00Z",
  "end_time": "2026-04-06T09:30:00Z",
  "status": "ACTIVE",
  "created_at": "2026-04-05T10:15:30Z",
  "updated_at": "2026-04-05T10:15:30Z"
}
```

### Modifying Bookings

#### Modify Booking Endpoint
```
PUT /api/v1/bookings/{id}
```

**Request:**
```json
{
  "start_time": "2026-04-06T10:00:00Z"
}
```

**Validation:**
- Validates new start_time aligns to 30-minute boundary
- Automatically calculates end_time (+ 30 minutes)
- Checks for conflicts

### Getting Available Slots

#### Get Slots Endpoint
```
GET /api/v1/coaches/{coach_id}/slots?date=2026-04-06&timezone=America/New_York
```

**Response:**
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

✅ All slots are exactly 30 minutes with aligned boundaries

---

## Error Responses

### Invalid Time Format
```json
{
  "error": "invalid start_time format"
}
```
**HTTP 400**

### Not Aligned to 30-Minute Boundary
```json
{
  "error": "start_time must align to 30-minute boundaries"
}
```
**HTTP 400**

### Invalid Duration (for availability)
```json
{
  "error": "slot duration must be exactly 30 minutes, got 45 minutes"
}
```
**HTTP 400**

### Slot Already Booked
```json
{
  "error": "slot is already booked"
}
```
**HTTP 400**

### Cannot Book Past Time
```json
{
  "error": "cannot book past slots"
}
```
**HTTP 400**

---

## Examples

### Example 1: Valid Booking Request
```bash
curl -X POST http://localhost:8080/api/v1/bookings \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 1,
    "coach_id": 1,
    "start_time": "2026-04-06T09:30:00Z"
  }'
```

**Response (201 Created):**
```json
{
  "id": 18,
  "user_id": 1,
  "coach_id": 1,
  "start_time": "2026-04-06T09:30:00Z",
  "end_time": "2026-04-06T10:00:00Z",
  "status": "ACTIVE",
  "created_at": "2026-04-05T15:30:00Z",
  "updated_at": "2026-04-05T15:30:00Z"
}
```

✅ Start time at `:30` boundary → ✅ End time automatically set to 30 min later

---

### Example 2: Invalid - Wrong Time Boundary
```bash
curl -X POST http://localhost:8080/api/v1/bookings \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 1,
    "coach_id": 1,
    "start_time": "2026-04-06T09:15:00Z"
  }'
```

**Response (400 Bad Request):**
```json
{
  "error": "start_time must align to 30-minute boundaries"
}
```

❌ `:15` is not a valid boundary (must be `:00` or `:30`)

---

### Example 3: Valid Weekly Availability Setup
```bash
curl -X POST http://localhost:8080/api/v1/coaches/1/availability/weekly \
  -H "Content-Type: application/json" \
  -d '{
    "availabilities": [
      {
        "day_of_week": 1,
        "slots": [
          {"start_time": "09:00", "end_time": "09:30"},
          {"start_time": "09:30", "end_time": "10:00"},
          {"start_time": "10:00", "end_time": "10:30"}
        ]
      }
    ]
  }'
```

**Response (200 OK):**
```json
{
  "message": "weekly availability set successfully"
}
```

✅ All slots are exactly 30 minutes

---

### Example 4: Invalid - Wrong Duration
```bash
curl -X POST http://localhost:8080/api/v1/coaches/1/availability \
  -H "Content-Type: application/json" \
  -d '{
    "day_of_week": 1,
    "slots": [
      {"start_time": "09:00", "end_time": "09:45"}
    ]
  }'
```

**Response (400 Bad Request):**
```json
{
  "error": "slot duration must be exactly 30 minutes, got 45 minutes"
}
```

❌ Slot is 45 minutes instead of 30

---

## Implementation Details

### TimeSlotInput Validation
```go
// Validates:
// 1. Format: "HH:MM"
// 2. Hours: 0-23
// 3. Minutes: 0-59
// 4. Start aligns to :00 or :30
// 5. End aligns to :00 or :30
// 6. Duration = End - Start = exactly 30 minutes
// 7. Start < End
func (t *TimeSlotInput) Validate() error { ... }
```

### Service-Level Enforcement
```go
// CreateBooking enforces:
if startTime.Minute()%30 != 0 {
    return nil, errors.New("start_time must align to 30-minute boundaries")
}
endTime := startTime.Add(30 * time.Minute) // Always 30 min

// ModifyBooking enforces same rules
if newStartTime.Minute()%30 != 0 {
    return nil, errors.New("start_time must align to 30-minute boundaries")
}
newEndTime := newStartTime.Add(30 * time.Minute) // Always 30 min
```

---

## Summary

| Feature | Status | Details |
|---------|--------|---------|
| Time Boundaries | ✅ Enforced | Must be `:00` or `:30` |
| Duration | ✅ Fixed | Always exactly 30 minutes |
| Validation | ✅ Multiple layers | DTO, Service, Database |
| Flexibility | ⚠️ Limited | Users can only book 30-min slots |
| Conflicts | ✅ Prevented | No double-booking allowed |
| Error Messages | ✅ Clear | Specific guidance for corrections |

Users **cannot book arbitrary durations or times**. The system enforces 30-minute slots across all endpoints.
