package repository

import (
	"context"

	"github.com/booking-api-service/internal/model"
)

// UserRepository defines operations for users
type UserRepository interface {
	Create(ctx context.Context, user *model.User) (int64, error)
	GetByID(ctx context.Context, id int64) (*model.User, error)
	GetAll(ctx context.Context) ([]*model.User, error)
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id int64) error
}

// CoachRepository defines operations for coaches
type CoachRepository interface {
	Create(ctx context.Context, coach *model.Coach) (int64, error)
	GetByID(ctx context.Context, id int64) (*model.Coach, error)
	GetAll(ctx context.Context) ([]*model.Coach, error)
	Update(ctx context.Context, coach *model.Coach) error
	Delete(ctx context.Context, id int64) error
}

// AvailabilityRepository defines operations for availability
type AvailabilityRepository interface {
	Create(ctx context.Context, availability *model.Availability) (int64, error)
	GetByCoachAndDay(ctx context.Context, coachID int64, dayOfWeek int) ([]*model.Availability, error)
	GetByCoach(ctx context.Context, coachID int64) ([]*model.Availability, error)
	Update(ctx context.Context, availability *model.Availability) error
	Delete(ctx context.Context, id int64) error
	DeleteByCoachAndDay(ctx context.Context, coachID int64, dayOfWeek int) error
}

// AvailabilityExceptionRepository defines operations for availability exceptions
type AvailabilityExceptionRepository interface {
	Create(ctx context.Context, exception *model.AvailabilityException) (int64, error)
	GetByCoachAndDate(ctx context.Context, coachID int64, date string) (*model.AvailabilityException, error)
	GetByCoach(ctx context.Context, coachID int64) ([]*model.AvailabilityException, error)
	Update(ctx context.Context, exception *model.AvailabilityException) error
	Delete(ctx context.Context, id int64) error
}

// BookingRepository defines operations for bookings
type BookingRepository interface {
	Create(ctx context.Context, booking *model.Booking) (int64, error)
	GetByID(ctx context.Context, id int64) (*model.Booking, error)
	GetByUser(ctx context.Context, userID int64) ([]*model.Booking, error)
	GetByCoachAndDateRange(ctx context.Context, coachID int64, startTime string, endTime string) ([]*model.Booking, error)
	Update(ctx context.Context, booking *model.Booking) error
	Delete(ctx context.Context, id int64) error
	GetByIdempotencyKey(ctx context.Context, key string) (*model.Booking, error)
}
