# Weekly Availability Implementation

## Summary
Updated the booking API to support setting availability for the entire week in a single request, in addition to the existing single-day functionality.

## Changes Made

### 1. **DTO Updates** (`internal/dto/requests.go`)
Added new request types to support batch availability:

```go
// SetWeeklyAvailabilityRequest is the request to set availability for entire week
type SetWeeklyAvailabilityRequest struct {
	Availabilities []DayAvailabilityInput `json:"availabilities"`
}

// DayAvailabilityInput represents availability for a specific day of the week
type DayAvailabilityInput struct {
	DayOfWeek int             `json:"day_of_week"` // 0=Sunday, 1=Monday, ..., 6=Saturday
	Slots     []TimeSlotInput `json:"slots"`
}
```

### 2. **Service Interface Update** (`internal/service/interfaces.go`)
Added new method to `AvailabilityService` interface:

```go
SetWeeklyAvailability(ctx context.Context, coachID int64, req *dto.SetWeeklyAvailabilityRequest) error
```

### 3. **Service Implementation** (`internal/service/availability_service.go`)
Implemented `SetWeeklyAvailability` method with:
- Validation of coach existence
- Validation of day of week (0-6)
- Validation of time formats (30-minute boundaries)
- Automatic clearing of all existing availability for the coach
- Creation of new availability slots for all provided days
- Comprehensive error handling with detailed error messages

### 4. **Handler Update** (`internal/handler/availability_handler.go`)
Added new handler method `SetWeeklyAvailability`:
- Extracts coach_id from URL parameters
- Decodes JSON request body
- Calls service method
- Returns appropriate HTTP responses

### 5. **Routes Update** (`cmd/main.go`)
Added new endpoint:
```
POST /api/v1/coaches/{coach_id}/availability/weekly
```

## API Usage Examples

### Set Availability for Single Day (Existing)
```bash
curl -X POST http://localhost:8080/api/v1/coaches/1/availability \
  -H "Content-Type: application/json" \
  -d '{
    "day_of_week": 1,
    "slots": [
      {"start_time": "09:00", "end_time": "12:00"},
      {"start_time": "14:00", "end_time": "17:00"}
    ]
  }'
```

### Set Availability for Entire Week (New)
```bash
curl -X POST http://localhost:8080/api/v1/coaches/1/availability/weekly \
  -H "Content-Type: application/json" \
  -d '{
    "availabilities": [
      {
        "day_of_week": 1,
        "slots": [
          {"start_time": "09:00", "end_time": "12:00"},
          {"start_time": "14:00", "end_time": "17:00"}
        ]
      },
      {
        "day_of_week": 2,
        "slots": [
          {"start_time": "09:00", "end_time": "12:00"},
          {"start_time": "14:00", "end_time": "18:00"}
        ]
      },
      {
        "day_of_week": 3,
        "slots": [
          {"start_time": "10:00", "end_time": "13:00"}
        ]
      },
      {
        "day_of_week": 4,
        "slots": [
          {"start_time": "09:00", "end_time": "12:00"},
          {"start_time": "14:00", "end_time": "17:00"}
        ]
      },
      {
        "day_of_week": 5,
        "slots": [
          {"start_time": "09:00", "end_time": "12:00"},
          {"start_time": "14:00", "end_time": "17:00"}
        ]
      }
    ]
  }'
```

## Key Features

1. **Backward Compatible**: Original single-day endpoint still works
2. **Atomic Operation**: All days are processed in a single request
3. **Automatic Cleanup**: Old availability is cleared and replaced with new settings
4. **Flexible**: Can set availability for any combination of days
5. **Validation**: Comprehensive validation for time formats and day of week values
6. **Error Handling**: Detailed error messages for debugging

## Day of Week Mapping

- `0` = Sunday
- `1` = Monday
- `2` = Tuesday
- `3` = Wednesday
- `4` = Thursday
- `5` = Friday
- `6` = Saturday

## Time Format

- All times must be in `HH:MM` format
- Times must align to 30-minute boundaries (e.g., 09:00, 09:30, 10:00)
- Start time must be before end time
