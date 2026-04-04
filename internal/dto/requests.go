package dto

// CreateAvailabilityRequest is the request to set availability for a coach
type CreateAvailabilityRequest struct {
	DayOfWeek int             `json:"day_of_week"`
	Slots     []TimeSlotInput `json:"slots"`
}

// TimeSlotInput represents a time slot
type TimeSlotInput struct {
	StartTime string `json:"start_time"` // HH:MM format
	EndTime   string `json:"end_time"`   // HH:MM format
}

// CreateAvailabilityExceptionRequest is the request to add an exception
type CreateAvailabilityExceptionRequest struct {
	Date        string  `json:"date"` // YYYY-MM-DD format
	IsAvailable bool    `json:"is_available"`
	StartTime   *string `json:"start_time"` // Optional, HH:MM format
	EndTime     *string `json:"end_time"`   // Optional, HH:MM format
}

// GetSlotsRequest is the query parameters for getting available slots
type GetSlotsRequest struct {
	Date     string `form:"date"`     // YYYY-MM-DD format
	Timezone string `form:"timezone"` // e.g., "Asia/Kolkata"
}

// CreateBookingRequest is the request to create a booking
type CreateBookingRequest struct {
	UserID         int64   `json:"user_id"`
	CoachID        int64   `json:"coach_id"`
	StartTime      string  `json:"start_time"`      // RFC3339 format
	IdempotencyKey *string `json:"idempotency_key"` // Optional
}

// ModifyBookingRequest is the request to modify a booking
type ModifyBookingRequest struct {
	StartTime string `json:"start_time"` // RFC3339 format
}

// AvailabilitySlot represents an available slot
type AvailabilitySlot struct {
	StartTime string `json:"start_time"` // RFC3339 format
	EndTime   string `json:"end_time"`   // RFC3339 format
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
	EndTime   string `json:"end_time"`
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
