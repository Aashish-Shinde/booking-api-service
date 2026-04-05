package utils

import (
	"errors"
	"fmt"
	"time"
)

// TimezoneHelper provides timezone conversion utilities for the booking system
type TimezoneHelper struct{}

// NewTimezoneHelper creates a new timezone helper instance
func NewTimezoneHelper() *TimezoneHelper {
	return &TimezoneHelper{}
}

// ConvertTimeToCoachTimezone converts a user's local time to coach's local time
// This helps determine what time it is for the coach when the user specifies a time
func (th *TimezoneHelper) ConvertTimeToCoachTimezone(userTime time.Time, userTimezone string, coachTimezone string) (time.Time, error) {
	// Load coach timezone
	coachLoc, err := time.LoadLocation(coachTimezone)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid coach timezone: %v", err)
	}

	// Convert user time to coach's timezone
	coachTime := userTime.In(coachLoc)
	return coachTime, nil
}

// ConvertTimeToUserTimezone converts a time to user's local timezone
func (th *TimezoneHelper) ConvertTimeToUserTimezone(utcTime time.Time, userTimezone string) (time.Time, error) {
	userLoc, err := time.LoadLocation(userTimezone)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid user timezone: %v", err)
	}

	userTime := utcTime.In(userLoc)
	return userTime, nil
}

// GetDayRangeInUTC converts a date in user's timezone to UTC range (start and end of day)
// This is used to query the database for availability across day boundaries
func (th *TimezoneHelper) GetDayRangeInUTC(dateStr string, userTimezone string) (time.Time, time.Time, error) {
	// Parse the date string (YYYY-MM-DD)
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid date format: %v", err)
	}

	// Load user timezone
	userLoc, err := time.LoadLocation(userTimezone)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid timezone: %v", err)
	}

	// Create start of day in user's timezone (00:00:00)
	startOfDayLocal := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, userLoc)
	startOfDayUTC := startOfDayLocal.UTC()

	// Create end of day in user's timezone (23:59:59)
	endOfDayLocal := time.Date(date.Year(), date.Month(), date.Day(), 23, 59, 59, 999999999, userLoc)
	endOfDayUTC := endOfDayLocal.UTC()

	return startOfDayUTC, endOfDayUTC, nil
}

// IsTimeInCoachAvailability checks if a given UTC time falls within coach's availability
// considering the coach's timezone and recurring availability
func (th *TimezoneHelper) IsTimeInCoachAvailability(utcTime time.Time, coachTimezone string, dayOfWeek int, startTimeStr string, endTimeStr string) (bool, error) {
	// Load coach timezone
	coachLoc, err := time.LoadLocation(coachTimezone)
	if err != nil {
		return false, fmt.Errorf("invalid coach timezone: %v", err)
	}

	// Convert UTC time to coach's local timezone
	coachTime := utcTime.In(coachLoc)

	// Check if the day of week matches
	// Note: Go's time.Weekday: 0=Sunday, 1=Monday, ..., 6=Saturday
	if int(coachTime.Weekday()) != dayOfWeek {
		return false, nil
	}

	// Parse start and end times (HH:MM format)
	startMinutes, err := TimeStrToMinutes(startTimeStr)
	if err != nil {
		return false, fmt.Errorf("invalid start time: %v", err)
	}

	endMinutes, err := TimeStrToMinutes(endTimeStr)
	if err != nil {
		return false, fmt.Errorf("invalid end time: %v", err)
	}

	// Get current time in coach's timezone as minutes since midnight
	coachTimeMinutes := coachTime.Hour()*60 + coachTime.Minute()

	// Check if time falls within availability window
	return coachTimeMinutes >= startMinutes && coachTimeMinutes < endMinutes, nil
}

// GetDayOfWeekInTimezone returns the day of week for a date in a specific timezone
// This is important because the same UTC time might be different days in different timezones
func (th *TimezoneHelper) GetDayOfWeekInTimezone(utcTime time.Time, timezone string) (int, error) {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return -1, fmt.Errorf("invalid timezone: %v", err)
	}

	localTime := utcTime.In(loc)
	return int(localTime.Weekday()), nil
}

// GetDateInTimezone returns the date (YYYY-MM-DD) for a UTC time in a specific timezone
func (th *TimezoneHelper) GetDateInTimezone(utcTime time.Time, timezone string) (string, error) {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return "", fmt.Errorf("invalid timezone: %v", err)
	}

	localTime := utcTime.In(loc)
	return localTime.Format("2006-01-02"), nil
}

// AdjustTimeToCoachWorkingHours adjusts a requested time to fit within coach's working hours
// If the exact time isn't available, it finds the next available slot
func (th *TimezoneHelper) AdjustTimeToCoachWorkingHours(requestedUTCTime time.Time, coachTimezone string, startTimeStr string, endTimeStr string) (time.Time, error) {
	if coachTimezone == "" {
		return time.Time{}, errors.New("coach timezone is required")
	}

	// Load coach timezone
	coachLoc, err := time.LoadLocation(coachTimezone)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid coach timezone: %v", err)
	}

	// Convert requested UTC time to coach's timezone
	coachTime := requestedUTCTime.In(coachLoc)

	// Parse availability hours
	startMinutes, err := TimeStrToMinutes(startTimeStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid start time: %v", err)
	}

	endMinutes, err := TimeStrToMinutes(endTimeStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid end time: %v", err)
	}

	// Get current time in coach's timezone as minutes since midnight
	currentMinutes := coachTime.Hour()*60 + coachTime.Minute()

	// If current time is before start, adjust to start time
	if currentMinutes < startMinutes {
		adjustedTime := time.Date(coachTime.Year(), coachTime.Month(), coachTime.Day(),
			startMinutes/60, startMinutes%60, 0, 0, coachLoc)
		return adjustedTime.UTC(), nil
	}

	// If current time is within working hours, round up to next 30-minute boundary
	if currentMinutes < endMinutes {
		nextSlot := ((currentMinutes / 30) + 1) * 30
		if nextSlot < endMinutes {
			adjustedTime := time.Date(coachTime.Year(), coachTime.Month(), coachTime.Day(),
				nextSlot/60, nextSlot%60, 0, 0, coachLoc)
			return adjustedTime.UTC(), nil
		}
	}

	// If after working hours, move to next day's start time
	nextDay := coachTime.AddDate(0, 0, 1)
	adjustedTime := time.Date(nextDay.Year(), nextDay.Month(), nextDay.Day(),
		startMinutes/60, startMinutes%60, 0, 0, coachLoc)
	return adjustedTime.UTC(), nil
}

// ValidateTimeInAvailabilityWindow checks if a specific UTC time is available
// considering coach's recurring availability and timezone
func (th *TimezoneHelper) ValidateTimeInAvailabilityWindow(utcTime time.Time, coachTimezone string, availabilityRules []AvailabilityRule) (bool, error) {
	if coachTimezone == "" {
		return false, errors.New("coach timezone is required")
	}

	// Load coach timezone
	coachLoc, err := time.LoadLocation(coachTimezone)
	if err != nil {
		return false, fmt.Errorf("invalid coach timezone: %v", err)
	}

	// Convert UTC time to coach's local timezone
	coachTime := utcTime.In(coachLoc)
	dayOfWeek := int(coachTime.Weekday())
	timeMinutes := coachTime.Hour()*60 + coachTime.Minute()

	// Check each availability rule for this day
	for _, rule := range availabilityRules {
		if rule.DayOfWeek == dayOfWeek {
			startMinutes, _ := TimeStrToMinutes(rule.StartTime)
			endMinutes, _ := TimeStrToMinutes(rule.EndTime)

			if timeMinutes >= startMinutes && timeMinutes < endMinutes {
				return true, nil
			}
		}
	}

	return false, nil
}

// AvailabilityRule represents a coach's recurring availability rule
type AvailabilityRule struct {
	DayOfWeek int    // 0=Sunday, 1=Monday, ..., 6=Saturday
	StartTime string // HH:MM format
	EndTime   string // HH:MM format
}
