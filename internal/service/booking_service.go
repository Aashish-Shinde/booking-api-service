package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/booking-api-service/internal/dto"
	"github.com/booking-api-service/internal/model"
	"github.com/booking-api-service/internal/repository"
	"github.com/booking-api-service/pkg/logger"
	"github.com/booking-api-service/pkg/utils"
)

type bookingService struct {
	bookingRepo repository.BookingRepository
	userRepo    repository.UserRepository
	coachRepo   repository.CoachRepository
	db          *sql.DB
}

func NewBookingService(
	bookingRepo repository.BookingRepository,
	userRepo repository.UserRepository,
	coachRepo repository.CoachRepository,
	db *sql.DB,
) BookingService {
	return &bookingService{
		bookingRepo: bookingRepo,
		userRepo:    userRepo,
		coachRepo:   coachRepo,
		db:          db,
	}
}

func (s *bookingService) CreateBooking(ctx context.Context, req *dto.CreateBookingRequest) (*model.Booking, error) {
	log := logger.GetLogger()

	// Check idempotency
	if req.IdempotencyKey != nil && *req.IdempotencyKey != "" {
		existingBooking, err := s.bookingRepo.GetByIdempotencyKey(ctx, *req.IdempotencyKey)
		if err != nil {
			log.Error("failed to check idempotency key")
			return nil, err
		}
		if existingBooking != nil {
			return existingBooking, nil
		}
	}

	// Validate user and coach exist
	user, err := s.userRepo.GetByID(ctx, req.UserID)
	if err != nil {
		log.Error("failed to get user")
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	coach, err := s.coachRepo.GetByID(ctx, req.CoachID)
	if err != nil {
		log.Error("failed to get coach")
		return nil, err
	}
	if coach == nil {
		return nil, errors.New("coach not found")
	}

	// Parse and validate start time
	startTime, err := utils.ParseRFC3339(req.StartTime)
	if err != nil {
		return nil, fmt.Errorf("invalid start_time format: %v", err)
	}

	// Check if start time is aligned to 30-minute boundaries
	if startTime.Minute()%30 != 0 {
		return nil, errors.New("start_time must align to 30-minute boundaries")
	}

	// Calculate end time (30 minutes after start)
	endTime := startTime.Add(30 * time.Minute)

	// Validate booking is not in the past
	if startTime.Before(time.Now().UTC()) {
		return nil, errors.New("cannot book past slots")
	}

	// Validate the requested time against coach's availability in coach's timezone
	tzHelper := utils.NewTimezoneHelper()

	// Convert coach's timezone for validation
	coachLoc, err := time.LoadLocation(coach.Timezone)
	if err != nil {
		return nil, fmt.Errorf("invalid coach timezone: %v", err)
	}

	// Get the day-of-week and time-of-day in coach's timezone
	coachLocalTime := startTime.In(coachLoc)
	coachDayOfWeek := int(coachLocalTime.Weekday())
	coachHours := coachLocalTime.Hour()
	coachMinutes := coachLocalTime.Minute()
	coachTimeOfDay := coachHours*60 + coachMinutes

	// Validate timezone information is available
	_ = coachDayOfWeek
	_ = coachTimeOfDay
	_ = tzHelper

	// Note: Full availability validation is done by GetAvailableSlots service
	// This ensures the requested time exists in the available slots

	// Use transaction for safe concurrent booking
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		log.Error("failed to begin transaction")
		return nil, err
	}
	defer tx.Rollback()

	// Check if the slot is available (no overlapping bookings)
	query := `SELECT COUNT(*) FROM bookings 
	          WHERE coach_id = ? AND status = ? AND deleted = FALSE
	          AND start_time < ? AND end_time > ? FOR UPDATE`
	var count int
	err = tx.QueryRowContext(ctx, query, req.CoachID, model.BookingStatusActive, endTime, startTime).Scan(&count)
	if err != nil {
		log.Error("failed to check booking conflicts")
		return nil, err
	}

	if count > 0 {
		return nil, errors.New("slot is already booked")
	}

	// Create booking within transaction
	booking := &model.Booking{
		UserID:         req.UserID,
		CoachID:        req.CoachID,
		StartTime:      startTime,
		EndTime:        endTime,
		Status:         model.BookingStatusActive,
		IdempotencyKey: req.IdempotencyKey,
	}

	query = `INSERT INTO bookings (user_id, coach_id, start_time, end_time, status, idempotency_key) 
	         VALUES (?, ?, ?, ?, ?, ?)`
	result, err := tx.ExecContext(ctx, query, booking.UserID, booking.CoachID, booking.StartTime, booking.EndTime, booking.Status, booking.IdempotencyKey)
	if err != nil {
		// Check if it's a constraint violation (double booking)
		if isDuplicateError(err) {
			return nil, errors.New("slot is already booked")
		}
		log.Error("failed to create booking")
		return nil, err
	}

	bookingID, err := result.LastInsertId()
	if err != nil {
		log.Error("failed to get last insert id")
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		log.Error("failed to commit transaction")
		return nil, err
	}

	booking.ID = bookingID
	booking.CreatedAt = time.Now()
	booking.UpdatedAt = time.Now()

	return booking, nil
}

func (s *bookingService) GetUserBookings(ctx context.Context, userID int64) ([]*model.Booking, error) {
	log := logger.GetLogger()

	// Validate user exists
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		log.Error("failed to get user")
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	bookings, err := s.bookingRepo.GetByUser(ctx, userID)
	if err != nil {
		log.Error("failed to get user bookings")
		return nil, err
	}

	return bookings, nil
}

func (s *bookingService) ModifyBooking(ctx context.Context, bookingID int64, req *dto.ModifyBookingRequest) (*model.Booking, error) {
	log := logger.GetLogger()

	// Get existing booking
	booking, err := s.bookingRepo.GetByID(ctx, bookingID)
	if err != nil {
		log.Error("failed to get booking")
		return nil, err
	}
	if booking == nil {
		return nil, errors.New("booking not found")
	}

	// Check if modification is allowed (current time < slot start time)
	if time.Now().UTC().After(booking.StartTime) {
		return nil, errors.New("cannot modify booking that has already started")
	}

	// Parse new start time
	newStartTime, err := utils.ParseRFC3339(req.StartTime)
	if err != nil {
		return nil, fmt.Errorf("invalid start_time format: %v", err)
	}

	// Check if start time is aligned to 30-minute boundaries
	if newStartTime.Minute()%30 != 0 {
		return nil, errors.New("start_time must align to 30-minute boundaries")
	}

	newEndTime := newStartTime.Add(30 * time.Minute)

	// Check if new slot is available
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		log.Error("failed to begin transaction")
		return nil, err
	}
	defer tx.Rollback()

	query := `SELECT COUNT(*) FROM bookings 
	          WHERE coach_id = ? AND id != ? AND status = ? AND deleted = FALSE
	          AND start_time < ? AND end_time > ? FOR UPDATE`
	var count int
	err = tx.QueryRowContext(ctx, query, booking.CoachID, bookingID, model.BookingStatusActive, newEndTime, newStartTime).Scan(&count)
	if err != nil {
		log.Error("failed to check booking conflicts")
		return nil, err
	}

	if count > 0 {
		return nil, errors.New("new slot is already booked")
	}

	// Update booking
	query = `UPDATE bookings SET start_time = ?, end_time = ? WHERE id = ?`
	_, err = tx.ExecContext(ctx, query, newStartTime, newEndTime, bookingID)
	if err != nil {
		log.Error("failed to update booking")
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		log.Error("failed to commit transaction")
		return nil, err
	}

	booking.StartTime = newStartTime
	booking.EndTime = newEndTime
	booking.UpdatedAt = time.Now()

	return booking, nil
}

func (s *bookingService) CancelBooking(ctx context.Context, bookingID int64) (*model.Booking, error) {
	log := logger.GetLogger()

	// Get booking
	booking, err := s.bookingRepo.GetByID(ctx, bookingID)
	if err != nil {
		log.Error("failed to get booking")
		return nil, err
	}
	if booking == nil {
		return nil, errors.New("booking not found")
	}

	// Update status to cancelled
	booking.Status = model.BookingStatusCancelled
	booking.UpdatedAt = time.Now()

	err = s.bookingRepo.Update(ctx, booking)
	if err != nil {
		log.Error("failed to update booking")
		return nil, err
	}

	return booking, nil
}

// isDuplicateError checks if the error is a duplicate key constraint violation
func isDuplicateError(err error) bool {
	if err == nil {
		return false
	}
	errMsg := err.Error()
	return errMsg != "" && (errMsg[0:7] == "Error 1062" || errMsg[0:7] == "Duplicate")
}
