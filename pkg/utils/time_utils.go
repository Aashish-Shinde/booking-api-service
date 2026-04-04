package utils

import (
	"fmt"
	"time"
)

// TimeStringToTime converts HH:MM format to time.Time (with today's date)
func TimeStringToTime(timeStr string) (time.Time, error) {
	return time.Parse("15:04", timeStr)
}

// DateStringToTime converts YYYY-MM-DD format to time.Time
func DateStringToTime(dateStr string) (time.Time, error) {
	return time.Parse("2006-01-02", dateStr)
}

// GetDayOfWeek returns 0-6 for a given date (0=Sunday)
func GetDayOfWeek(date time.Time) int {
	return int(date.Weekday())
}

// TimeToHHMM converts time.Time to HH:MM format
func TimeToHHMM(t time.Time) string {
	return t.Format("15:04")
}

// DateToYYYYMMDD converts time.Time to YYYY-MM-DD format
func DateToYYYYMMDD(t time.Time) string {
	return t.Format("2006-01-02")
}

// ParseRFC3339 parses RFC3339 timestamp
func ParseRFC3339(timestamp string) (time.Time, error) {
	return time.Parse(time.RFC3339, timestamp)
}

// AddMinutes adds minutes to a time
func AddMinutes(t time.Time, minutes int) time.Time {
	return t.Add(time.Duration(minutes) * time.Minute)
}

// TimeStrToMinutes converts HH:MM to minutes since midnight
func TimeStrToMinutes(timeStr string) (int, error) {
	t, err := TimeStringToTime(timeStr)
	if err != nil {
		return 0, err
	}
	return t.Hour()*60 + t.Minute(), nil
}

// MinutesToTimeStr converts minutes since midnight to HH:MM
func MinutesToTimeStr(minutes int) string {
	hours := minutes / 60
	mins := minutes % 60
	return fmt.Sprintf("%02d:%02d", hours, mins)
}

// CheckTimeAlignment checks if a time is aligned to 30-minute boundaries
func CheckTimeAlignment(timeStr string) bool {
	minutes, err := TimeStrToMinutes(timeStr)
	if err != nil {
		return false
	}
	return minutes%30 == 0
}
