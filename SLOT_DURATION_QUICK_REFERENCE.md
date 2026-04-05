# 30-Minute Slot Enforcement - Quick Reference

## TL;DR

✅ **YES** - Users can book any slot where:
- Start time is `:00` or `:30` (e.g., `09:00`, `09:30`, `10:00`)
- Duration is exactly 30 minutes
- No conflict with existing bookings

❌ **NO** - Users CANNOT book:
- Arbitrary durations (1 hour, 45 min, 15 min, etc.)
- Non-aligned times (`:15`, `:45`, etc.)
- Overlapping time slots

---

## Valid Examples (✅ Can be booked)

```
09:00 - 09:30  ✅ (30 min, aligned)
09:30 - 10:00  ✅ (30 min, aligned)
10:00 - 10:30  ✅ (30 min, aligned)
14:30 - 15:00  ✅ (30 min, aligned)
23:30 - 00:00  ✅ (30 min, aligned, crosses midnight)
```

---

## Invalid Examples (❌ Cannot be booked)

```
09:00 - 10:00  ❌ (60 minutes - too long)
09:00 - 09:15  ❌ (15 minutes - too short)
09:15 - 09:45  ❌ (30 min but NOT aligned)
09:45 - 10:15  ❌ (30 min but NOT aligned)
13:00 - 13:30 (if already booked)  ❌ (conflict)
```

---

## Booking Time Request

```bash
# Valid: Aligned to :00
curl -X POST http://localhost:8080/api/v1/bookings \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 1,
    "coach_id": 1,
    "start_time": "2026-04-06T09:00:00Z"
  }'

# Valid: Aligned to :30
curl -X POST http://localhost:8080/api/v1/bookings \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 1,
    "coach_id": 1,
    "start_time": "2026-04-06T09:30:00Z"
  }'

# Invalid: NOT aligned (will be rejected)
curl -X POST http://localhost:8080/api/v1/bookings \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 1,
    "coach_id": 1,
    "start_time": "2026-04-06T09:15:00Z"
  }'
```

---

## Response Structure

When a booking is created, the system ALWAYS sets:
- **Duration:** Exactly 30 minutes
- **End Time:** Start Time + 30 minutes

```json
{
  "id": 18,
  "user_id": 1,
  "coach_id": 1,
  "start_time": "2026-04-06T09:30:00Z",
  "end_time": "2026-04-06T10:00:00Z",  ← Always start + 30 min
  "status": "ACTIVE",
  "created_at": "2026-04-05T15:30:00Z",
  "updated_at": "2026-04-05T15:30:00Z"
}
```

---

## Error Messages

| Scenario | Error Message | Solution |
|----------|---------------|----------|
| `:15` or `:45` time | `start_time must align to 30-minute boundaries` | Use `:00` or `:30` |
| Wrong duration (availability) | `slot duration must be exactly 30 minutes` | Set end = start + 30 min |
| Slot booked | `slot is already booked` | Choose different time |
| Past time | `cannot book past slots` | Use future date/time |

---

## Key Points

1. **Always 30 Minutes:** Every booking is exactly 30 minutes. The system calculates end time automatically.

2. **Boundary Alignment:** Start times must end in `:00` or `:30`. No exceptions.

3. **No Flexibility:** Cannot book 1 hour, 45 minutes, 15 minutes, etc. Must be 30 minutes.

4. **Timezone Handling:** Times are in RFC3339 format with timezone. API handles conversions.

5. **Conflict Prevention:** Double-booking is impossible. Only available slots can be booked.

---

## Testing Checklist

- [ ] Can book at `:00` boundary ✅
- [ ] Can book at `:30` boundary ✅
- [ ] Cannot book at `:15` boundary ❌
- [ ] Cannot book at `:45` boundary ❌
- [ ] End time is always start + 30 min ✅
- [ ] Cannot book overlapping slots ❌
- [ ] Cannot book past times ❌
- [ ] Availability shows only 30-min slots ✅

---

## Files with 30-Minute Enforcement

1. **`internal/dto/requests.go`**
   - `TimeSlotInput.Validate()` - Validates 30-min duration
   - `CreateAvailabilityExceptionRequest.Validate()` - Validates exceptions

2. **`internal/service/booking_service.go`**
   - `CreateBooking()` - Enforces 30-min alignment
   - `ModifyBooking()` - Enforces 30-min alignment

3. **`internal/service/availability_service.go`**
   - `SetAvailability()` - Enforces 30-min slots
   - `SetWeeklyAvailability()` - Enforces 30-min slots

---

## Postman Collection

Use the included Postman collection to test:
- `Booking_API_Postman_Collection.json`

Pre-built requests demonstrate:
- Valid 30-minute bookings
- Valid availability setup
- Error cases (invalid times, boundaries, etc.)

---

**Summary:** 🎯 All slots are strictly 30 minutes. Users book by providing a start time (must be `:00` or `:30`). System automatically sets end time to 30 minutes later. No exceptions.
