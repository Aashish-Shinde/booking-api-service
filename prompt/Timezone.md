User: how do we find out which day  of week on particular date as to check availiablilty we use date to check coach availibility

GitHub Copilot: Read [](file:///home/ashish/Desktop/personal/booking-api-service/pkg/utils/timezone.go#1-1), lines 1 to 50

Read [](file:///home/ashish/Desktop/personal/booking-api-service/internal/service/availability_service.go#1-1), lines 1 to 100

Read [](file:///home/ashish/Desktop/personal/booking-api-service/internal/model/models.go)

Now let me check the time_utils.go file to see what utilities already exist: 

Read [](file:///home/ashish/Desktop/personal/booking-api-service/pkg/utils/time_utils.go)

Perfect! I can see the solution. You already have a **`GetDayOfWeek()`** function in your time_utils.go file that does exactly what you need. Let me show you how it works and how to use it in the availability checking:

## Solution: Finding Day of Week for a Date

You have a utility function already available in time_utils.go:

```go
// GetDayOfWeek returns 0-6 for a given date (0=Sunday)
func GetDayOfWeek(date time.Time) int {
	return int(date.Weekday())
}
```

### How to use it to check coach availability:

1. **Parse the date string** (e.g., "2026-04-05"):
```go
date, err := utils.DateStringToTime("2026-04-05")
```

2. **Get the day of week** (returns 0-6, where 0=Sunday, 1=Monday, etc.):
```go
dayOfWeek := utils.GetDayOfWeek(date)
```

3. **Query availability** using the dayOfWeek:
Your models already show that `Availability` has a `DayOfWeek` field (0-6) that matches this format:
```go
// Use dayOfWeek to query availability_repository
availabilitySlots := s.availabilityRepo.GetByCoachAndDay(ctx, coachID, dayOfWeek)
```

### Complete example:
```go
// Given a booking request with a date string
dateStr := "2026-04-05"
date, _ := utils.DateStringToTime(dateStr)
dayOfWeek := utils.GetDayOfWeek(date)  // Returns 3 (Wednesday)

// Now use dayOfWeek to check coach's recurring availability
availabilitySlots, _ := s.availabilityRepo.GetByCoachAndDay(ctx, coachID, dayOfWeek)
```

This approach lets you check the coach's **recurring weekly availability** (set once for each day of the week) against any specific date you want to book.

User: give me list of timezone we can use to set to the user,in chat

GitHub Copilot: Here's a comprehensive list of commonly used timezones you can set for users:

## **Americas**
- `America/New_York` - Eastern Time
- `America/Chicago` - Central Time
- `America/Denver` - Mountain Time
- `America/Los_Angeles` - Pacific Time
- `America/Anchorage` - Alaska Time
- `Pacific/Honolulu` - Hawaii Time
- `America/Toronto` - Eastern Canada
- `America/Vancouver` - Pacific Canada
- `America/Mexico_City` - Mexico City
- `America/Sao_Paulo` - Brazil
- `America/Buenos_Aires` - Argentina

## **Europe**
- `Europe/London` - GMT/BST
- `Europe/Paris` - CET/CEST
- `Europe/Berlin` - CET/CEST
- `Europe/Madrid` - CET/CEST
- `Europe/Rome` - CET/CEST
- `Europe/Amsterdam` - CET/CEST
- `Europe/Stockholm` - CET/CEST
- `Europe/Zurich` - CET/CEST
- `Europe/Vienna` - CET/CEST
- `Europe/Athens` - EET/EEST
- `Europe/Istanbul` - EET
- `Europe/Moscow` - MSK

## **Asia-Pacific**
- `Asia/Kolkata` - India Standard Time
- `Asia/Bangkok` - Indochina Time
- `Asia/Singapore` - Singapore Time
- `Asia/Hong_Kong` - Hong Kong Time
- `Asia/Shanghai` - China Standard Time
- `Asia/Tokyo` - Japan Standard Time
- `Asia/Seoul` - Korea Standard Time
- `Asia/Manila` - Philippines Time
- `Asia/Jakarta` - Western Indonesia Time
- `Australia/Sydney` - Australian Eastern Time
- `Australia/Melbourne` - Australian Eastern Time
- `Australia/Brisbane` - Australian Eastern Time
- `Australia/Perth` - Australian Western Time
- `Pacific/Auckland` - New Zealand Time
- `Pacific/Fiji` - Fiji Time

## **Middle East & Africa**
- `Asia/Dubai` - Gulf Standard Time
- `Asia/Karachi` - Pakistan Standard Time
- `Africa/Cairo` - Egypt Standard Time
- `Africa/Johannesburg` - South Africa Standard Time
- `Africa/Lagos` - West Africa Time
- `Africa/Nairobi` - East Africa Time

## **UTC**
- `UTC` - Coordinated Universal Time

You can validate these in your application using Go's `time.LoadLocation()` function. To set it for a user, just store one of these timezone strings in your `User` model's `Timezone` field, like:

```go
user.Timezone = "America/New_York"  // or any from the list above
```