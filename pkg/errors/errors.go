package errors

import "net/http"

// AppError represents an application error
type AppError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
	Status  int    `json:"-"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	return e.Message
}

// Error codes
const (
	ErrCodeInvalidInput        = "INVALID_INPUT"
	ErrCodeNotFound            = "NOT_FOUND"
	ErrCodeDuplicateBooking    = "DUPLICATE_BOOKING"
	ErrCodeConflict            = "CONFLICT"
	ErrCodeUnauthorized        = "UNAUTHORIZED"
	ErrCodeForbidden           = "FORBIDDEN"
	ErrCodeInternalServer      = "INTERNAL_SERVER_ERROR"
	ErrCodeDatabaseError       = "DATABASE_ERROR"
	ErrCodeInvalidTimezone     = "INVALID_TIMEZONE"
	ErrCodeInvalidTime         = "INVALID_TIME"
	ErrCodeBookingInPast       = "BOOKING_IN_PAST"
	ErrCodeCannotModifyBooking = "CANNOT_MODIFY_BOOKING"
	ErrCodeSlotAlreadyBooked   = "SLOT_ALREADY_BOOKED"
	ErrCodeCoachNotFound       = "COACH_NOT_FOUND"
	ErrCodeUserNotFound        = "USER_NOT_FOUND"
	ErrCodeBookingNotFound     = "BOOKING_NOT_FOUND"
)

// ErrorStatusMap maps error codes to HTTP status codes
var ErrorStatusMap = map[string]int{
	ErrCodeInvalidInput:        http.StatusBadRequest,
	ErrCodeNotFound:            http.StatusNotFound,
	ErrCodeDuplicateBooking:    http.StatusConflict,
	ErrCodeConflict:            http.StatusConflict,
	ErrCodeUnauthorized:        http.StatusUnauthorized,
	ErrCodeForbidden:           http.StatusForbidden,
	ErrCodeInternalServer:      http.StatusInternalServerError,
	ErrCodeDatabaseError:       http.StatusInternalServerError,
	ErrCodeInvalidTimezone:     http.StatusBadRequest,
	ErrCodeInvalidTime:         http.StatusBadRequest,
	ErrCodeBookingInPast:       http.StatusBadRequest,
	ErrCodeCannotModifyBooking: http.StatusBadRequest,
	ErrCodeSlotAlreadyBooked:   http.StatusConflict,
	ErrCodeCoachNotFound:       http.StatusNotFound,
	ErrCodeUserNotFound:        http.StatusNotFound,
	ErrCodeBookingNotFound:     http.StatusNotFound,
}

// NewAppError creates a new application error
func NewAppError(code, message string) *AppError {
	status, exists := ErrorStatusMap[code]
	if !exists {
		status = http.StatusInternalServerError
	}
	return &AppError{
		Code:    code,
		Message: message,
		Status:  status,
	}
}

// NewAppErrorWithDetails creates a new application error with details
func NewAppErrorWithDetails(code, message, details string) *AppError {
	status, exists := ErrorStatusMap[code]
	if !exists {
		status = http.StatusInternalServerError
	}
	return &AppError{
		Code:    code,
		Message: message,
		Details: details,
		Status:  status,
	}
}

// Predefined errors
func ErrorInvalidInput(msg string) *AppError {
	return NewAppError(ErrCodeInvalidInput, msg)
}

func ErrorNotFound(resource string) *AppError {
	return NewAppError(ErrCodeNotFound, resource+" not found")
}

func ErrorSlotAlreadyBooked() *AppError {
	return NewAppError(ErrCodeSlotAlreadyBooked, "slot is already booked")
}

func ErrorCoachNotFound() *AppError {
	return NewAppError(ErrCodeCoachNotFound, "coach not found")
}

func ErrorUserNotFound() *AppError {
	return NewAppError(ErrCodeUserNotFound, "user not found")
}

func ErrorBookingNotFound() *AppError {
	return NewAppError(ErrCodeBookingNotFound, "booking not found")
}

func ErrorInvalidTimezone(tz string) *AppError {
	return NewAppErrorWithDetails(ErrCodeInvalidTimezone, "invalid timezone", tz)
}

func ErrorInvalidTime(details string) *AppError {
	return NewAppErrorWithDetails(ErrCodeInvalidTime, "invalid time format", details)
}

func ErrorBookingInPast() *AppError {
	return NewAppError(ErrCodeBookingInPast, "cannot book past slots")
}

func ErrorCannotModifyBooking() *AppError {
	return NewAppError(ErrCodeCannotModifyBooking, "cannot modify booking that has already started")
}

func ErrorDatabaseError(details string) *AppError {
	return NewAppErrorWithDetails(ErrCodeDatabaseError, "database error", details)
}

func ErrorInternalServer() *AppError {
	return NewAppError(ErrCodeInternalServer, "internal server error")
}
