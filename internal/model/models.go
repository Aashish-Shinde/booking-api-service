package model

import "time"

// User represents a user in the system
type User struct {
	ID        int64
	Name      string
	Timezone  string
	CreatedAt time.Time
	UpdatedAt time.Time
	Deleted   bool
}

// Coach represents a coach in the system
type Coach struct {
	ID        int64
	Name      string
	Timezone  string
	CreatedAt time.Time
	UpdatedAt time.Time
	Deleted   bool
}

// Availability represents weekly recurring availability
type Availability struct {
	ID        int64
	CoachID   int64
	DayOfWeek int    // 0=Sunday, 1=Monday, ..., 6=Saturday
	StartTime string // HH:MM format
	EndTime   string // HH:MM format
	CreatedAt time.Time
	UpdatedAt time.Time
	Deleted   bool
}

// AvailabilityException represents one-time overrides (holidays/leaves)
type AvailabilityException struct {
	ID          int64
	CoachID     int64
	Date        time.Time
	StartTime   *string // HH:MM format, nil if full day disabled
	EndTime     *string // HH:MM format, nil if full day disabled
	IsAvailable bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Deleted     bool
}

// Booking represents a booking
type Booking struct {
	ID             int64
	UserID         int64
	CoachID        int64
	StartTime      time.Time // UTC
	EndTime        time.Time // UTC
	Status         string    // ACTIVE, CANCELLED, COMPLETED
	IdempotencyKey *string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Deleted        bool
}

// BookingStatus constants
const (
	BookingStatusActive    = "ACTIVE"
	BookingStatusCancelled = "CANCELLED"
	BookingStatusCompleted = "COMPLETED"
)
