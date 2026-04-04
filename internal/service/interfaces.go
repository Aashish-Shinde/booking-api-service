package service

import (
	"context"

	"github.com/booking-api-service/internal/dto"
	"github.com/booking-api-service/internal/model"
)

// AvailabilityService defines business logic for availability
type AvailabilityService interface {
	SetAvailability(ctx context.Context, coachID int64, req *dto.CreateAvailabilityRequest) error
	AddException(ctx context.Context, coachID int64, req *dto.CreateAvailabilityExceptionRequest) error
	GetAvailableSlots(ctx context.Context, coachID int64, date string, timezone string) (*dto.GetSlotsResponse, error)
}

// BookingService defines business logic for bookings
type BookingService interface {
	CreateBooking(ctx context.Context, req *dto.CreateBookingRequest) (*model.Booking, error)
	GetUserBookings(ctx context.Context, userID int64) ([]*model.Booking, error)
	ModifyBooking(ctx context.Context, bookingID int64, req *dto.ModifyBookingRequest) (*model.Booking, error)
	CancelBooking(ctx context.Context, bookingID int64) (*model.Booking, error)
}
