package utils

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Validator provides validation functions for common patterns
type Validator struct{}

// NewValidator creates a new validator
func NewValidator() *Validator {
	return &Validator{}
}

// ValidateDayOfWeek validates day of week (0-6)
func (v *Validator) ValidateDayOfWeek(day int) error {
	if day < 0 || day > 6 {
		return fmt.Errorf("invalid day of week: must be between 0-6, got %d", day)
	}
	return nil
}

// ValidateTime validates HH:MM format
func (v *Validator) ValidateTime(timeStr string) error {
	if !isValidTimeFormat(timeStr) {
		return fmt.Errorf("invalid time format: %s, expected HH:MM", timeStr)
	}
	return nil
}

// ValidateTimeRange validates start_time < end_time
func (v *Validator) ValidateTimeRange(startTime, endTime string) error {
	if err := v.ValidateTime(startTime); err != nil {
		return err
	}
	if err := v.ValidateTime(endTime); err != nil {
		return err
	}

	startMinutes, _ := TimeStrToMinutes(startTime)
	endMinutes, _ := TimeStrToMinutes(endTime)

	if startMinutes >= endMinutes {
		return fmt.Errorf("start_time (%s) must be before end_time (%s)", startTime, endTime)
	}
	return nil
}

// ValidateDate validates YYYY-MM-DD format
func (v *Validator) ValidateDate(dateStr string) error {
	if !isValidDateFormat(dateStr) {
		return fmt.Errorf("invalid date format: %s, expected YYYY-MM-DD", dateStr)
	}
	return nil
}

// ValidateTimezone validates IANA timezone format
func (v *Validator) ValidateTimezone(tz string) error {
	if tz == "" {
		return fmt.Errorf("timezone cannot be empty")
	}
	if strings.ContainsAny(tz, "\n\r\t") {
		return fmt.Errorf("timezone contains invalid characters")
	}
	// Basic validation - should be a valid IANA timezone
	// Full validation is done when loading location in time package
	if !isValidIANATimezone(tz) {
		return fmt.Errorf("invalid IANA timezone: %s", tz)
	}
	return nil
}

// ValidateID validates a positive integer ID
func (v *Validator) ValidateID(id int64) error {
	if id <= 0 {
		return fmt.Errorf("invalid ID: must be greater than 0, got %d", id)
	}
	return nil
}

// ValidateEmail validates email format (basic)
func (v *Validator) ValidateEmail(email string) error {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return fmt.Errorf("invalid email format: %s", email)
	}
	return nil
}

// ValidateName validates a name (non-empty, reasonable length)
func (v *Validator) ValidateName(name string) error {
	if name = strings.TrimSpace(name); name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	if len(name) > 255 {
		return fmt.Errorf("name too long: maximum 255 characters")
	}
	return nil
}

// ValidateStringLength validates string length
func (v *Validator) ValidateStringLength(s string, minLen, maxLen int) error {
	s = strings.TrimSpace(s)
	if len(s) < minLen {
		return fmt.Errorf("string too short: minimum %d characters", minLen)
	}
	if len(s) > maxLen {
		return fmt.Errorf("string too long: maximum %d characters", maxLen)
	}
	return nil
}

// Helper functions

func isValidTimeFormat(timeStr string) bool {
	if len(timeStr) != 5 {
		return false
	}
	parts := strings.Split(timeStr, ":")
	if len(parts) != 2 {
		return false
	}

	hours, err := strconv.Atoi(parts[0])
	if err != nil || hours < 0 || hours > 23 {
		return false
	}

	minutes, err := strconv.Atoi(parts[1])
	if err != nil || minutes < 0 || minutes > 59 {
		return false
	}

	return true
}

func isValidDateFormat(dateStr string) bool {
	if len(dateStr) != 10 {
		return false
	}
	parts := strings.Split(dateStr, "-")
	if len(parts) != 3 {
		return false
	}

	year, err := strconv.Atoi(parts[0])
	if err != nil || year < 2000 || year > 2100 {
		return false
	}

	month, err := strconv.Atoi(parts[1])
	if err != nil || month < 1 || month > 12 {
		return false
	}

	day, err := strconv.Atoi(parts[2])
	if err != nil || day < 1 || day > 31 {
		return false
	}

	return true
}

func isValidIANATimezone(tz string) bool {
	// Common valid timezones
	validTzs := map[string]bool{
		"UTC":                 true,
		"Asia/Kolkata":        true,
		"Asia/Tokyo":          true,
		"Asia/Shanghai":       true,
		"Europe/London":       true,
		"Europe/Paris":        true,
		"America/New_York":    true,
		"America/Los_Angeles": true,
		"America/Chicago":     true,
		"Australia/Sydney":    true,
	}

	// If in common list, it's valid
	if validTzs[tz] {
		return true
	}

	// Basic check: format should be Region/City or UTC
	if tz == "UTC" || tz == "GMT" {
		return true
	}

	parts := strings.Split(tz, "/")
	if len(parts) == 2 && len(parts[0]) > 0 && len(parts[1]) > 0 {
		return true
	}

	return false
}
