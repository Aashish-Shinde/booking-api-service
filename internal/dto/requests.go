package dto

import (
	"errors"
	"fmt"
	"strings"
)

// CreateAvailabilityRequest is the request to set availability for a coach (single day)
type CreateAvailabilityRequest struct {
	DayOfWeek int             `json:"day_of_week"`
	Slots     []TimeSlotInput `json:"slots"`
}

// SetWeeklyAvailabilityRequest is the request to set availability for entire week
type SetWeeklyAvailabilityRequest struct {
	Availabilities []DayAvailabilityInput `json:"availabilities"`
}

// DayAvailabilityInput represents availability for a specific day of the week
type DayAvailabilityInput struct {
	DayOfWeek int             `json:"day_of_week"` // 0=Sunday, 1=Monday, ..., 6=Saturday
	Slots     []TimeSlotInput `json:"slots"`
}

// TimeSlotInput represents a time slot (must be exactly 30 minutes)
type TimeSlotInput struct {
	StartTime string `json:"start_time"` // HH:MM format, must align to :00 or :30
	EndTime   string `json:"end_time"`   // HH:MM format, must be exactly 30 min after start
}

// Validate ensures the time slot is valid (HH:MM format, 30-minute boundaries, exactly 30 min duration)
func (t *TimeSlotInput) Validate() error {
	// Validate format
	if err := validateTimeFormat(t.StartTime); err != nil {
		return fmt.Errorf("invalid start_time: %v", err)
	}
	if err := validateTimeFormat(t.EndTime); err != nil {
		return fmt.Errorf("invalid end_time: %v", err)
	}

	// Parse times
	startMinutes := parseTimeToMinutes(t.StartTime)
	endMinutes := parseTimeToMinutes(t.EndTime)

	if startMinutes == -1 || endMinutes == -1 {
		return errors.New("invalid time format")
	}

	// Check 30-minute boundary alignment
	if startMinutes%30 != 0 {
		return errors.New("start_time must align to 30-minute boundaries (:00 or :30)")
	}
	if endMinutes%30 != 0 {
		return errors.New("end_time must align to 30-minute boundaries (:00 or :30)")
	}

	// Check duration is exactly 30 minutes
	if endMinutes-startMinutes != 30 {
		return fmt.Errorf("slot duration must be exactly 30 minutes, got %d minutes", endMinutes-startMinutes)
	}

	// Check start is before end
	if startMinutes >= endMinutes {
		return errors.New("start_time must be before end_time")
	}

	return nil
}

// CreateAvailabilityExceptionRequest is the request to add an exception
type CreateAvailabilityExceptionRequest struct {
	Date        string  `json:"date"` // YYYY-MM-DD format
	IsAvailable bool    `json:"is_available"`
	StartTime   *string `json:"start_time"` // Optional, HH:MM format, must align to :00 or :30
	EndTime     *string `json:"end_time"`   // Optional, HH:MM format, must be exactly 30 min after start
}

// Validate ensures the exception is valid
func (e *CreateAvailabilityExceptionRequest) Validate() error {
	// If only one of start/end time is provided, that's an error
	if (e.StartTime != nil && e.EndTime == nil) || (e.StartTime == nil && e.EndTime != nil) {
		return errors.New("both start_time and end_time must be provided together, or neither")
	}

	// If both are provided, validate them
	if e.StartTime != nil && e.EndTime != nil {
		startMinutes := parseTimeToMinutes(*e.StartTime)
		endMinutes := parseTimeToMinutes(*e.EndTime)

		if startMinutes == -1 || endMinutes == -1 {
			return errors.New("invalid time format")
		}

		// Check 30-minute boundary alignment
		if startMinutes%30 != 0 {
			return errors.New("start_time must align to 30-minute boundaries (:00 or :30)")
		}
		if endMinutes%30 != 0 {
			return errors.New("end_time must align to 30-minute boundaries (:00 or :30)")
		}

		// Check duration is exactly 30 minutes
		if endMinutes-startMinutes != 30 {
			return fmt.Errorf("slot duration must be exactly 30 minutes, got %d minutes", endMinutes-startMinutes)
		}

		if startMinutes >= endMinutes {
			return errors.New("start_time must be before end_time")
		}
	}

	return nil
}

// GetSlotsRequest is the query parameters for getting available slots
type GetSlotsRequest struct {
	Date     string `form:"date"`     // YYYY-MM-DD format
	Timezone string `form:"timezone"` // e.g., "Asia/Kolkata"
}

// CreateBookingRequest is the request to create a booking
// Note: Duration is always exactly 30 minutes, enforced by service
type CreateBookingRequest struct {
	UserID         int64   `json:"user_id"`
	CoachID        int64   `json:"coach_id"`
	StartTime      string  `json:"start_time"`      // RFC3339 format, must align to 30-minute boundaries
	IdempotencyKey *string `json:"idempotency_key"` // Optional
}

// ModifyBookingRequest is the request to modify a booking
// Note: Duration is always exactly 30 minutes, enforced by service
type ModifyBookingRequest struct {
	StartTime string `json:"start_time"` // RFC3339 format, must align to 30-minute boundaries
}

// AvailabilitySlot represents an available slot (always 30 minutes)
type AvailabilitySlot struct {
	StartTime string `json:"start_time"` // RFC3339 format
	EndTime   string `json:"end_time"`   // RFC3339 format (always 30 min after start_time)
}

// GetSlotsResponse is the response with available slots
type GetSlotsResponse struct {
	Slots []AvailabilitySlot `json:"slots"`
}

// BookingResponse is the response for booking operations
type BookingResponse struct {
	ID        int64  `json:"id"`
	UserID    int64  `json:"user_id"`
	CoachID   int64  `json:"coach_id"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"` // Always exactly 30 minutes after start_time
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// BookingsListResponse is the response for listing bookings
type BookingsListResponse struct {
	Bookings []BookingResponse `json:"bookings"`
	Count    int               `json:"count"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

// SuccessResponse represents a generic success response
type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Helper functions

// validateTimeFormat checks if time is in HH:MM format
func validateTimeFormat(t string) error {
	if len(t) != 5 {
		return errors.New("time must be in HH:MM format")
	}
	if t[2] != ':' {
		return errors.New("time must be in HH:MM format")
	}
	if !isDigit(t[0]) || !isDigit(t[1]) || !isDigit(t[3]) || !isDigit(t[4]) {
		return errors.New("time must contain only digits and colon")
	}
	return nil
}

// parseTimeToMinutes converts HH:MM to minutes since midnight
// Returns -1 if parsing fails
func parseTimeToMinutes(timeStr string) int {
	parts := strings.Split(timeStr, ":")
	if len(parts) != 2 {
		return -1
	}
	var hours, minutes int
	_, err := fmt.Sscanf(timeStr, "%d:%d", &hours, &minutes)
	if err != nil || hours < 0 || hours > 23 || minutes < 0 || minutes > 59 {
		return -1
	}
	return hours*60 + minutes
}

// isDigit checks if a character is a digit
func isDigit(b byte) bool {
	return b >= '0' && b <= '9'
}
