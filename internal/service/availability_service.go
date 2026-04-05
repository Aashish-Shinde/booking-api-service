package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/booking-api-service/internal/dto"
	"github.com/booking-api-service/internal/model"
	"github.com/booking-api-service/internal/repository"
	"github.com/booking-api-service/pkg/logger"
	"github.com/booking-api-service/pkg/utils"
)

type availabilityService struct {
	availabilityRepo          repository.AvailabilityRepository
	availabilityExceptionRepo repository.AvailabilityExceptionRepository
	bookingRepo               repository.BookingRepository
	coachRepo                 repository.CoachRepository
}

func NewAvailabilityService(
	availabilityRepo repository.AvailabilityRepository,
	availabilityExceptionRepo repository.AvailabilityExceptionRepository,
	bookingRepo repository.BookingRepository,
	coachRepo repository.CoachRepository,
) AvailabilityService {
	return &availabilityService{
		availabilityRepo:          availabilityRepo,
		availabilityExceptionRepo: availabilityExceptionRepo,
		bookingRepo:               bookingRepo,
		coachRepo:                 coachRepo,
	}
}

func (s *availabilityService) SetAvailability(ctx context.Context, coachID int64, req *dto.CreateAvailabilityRequest) error {
	log := logger.GetLogger()

	// Validate coach exists
	coach, err := s.coachRepo.GetByID(ctx, coachID)
	if err != nil {
		log.Error("failed to get coach")
		return err
	}
	if coach == nil {
		return errors.New("coach not found")
	}

	// Validate day of week
	if req.DayOfWeek < 0 || req.DayOfWeek > 6 {
		return errors.New("invalid day of week")
	}

	// Clear existing availability for this day
	err = s.availabilityRepo.DeleteByCoachAndDay(ctx, coachID, req.DayOfWeek)
	if err != nil {
		log.Error("failed to delete existing availability")
		return err
	}

	// Validate and create new availability slots
	for _, slot := range req.Slots {
		// Validate time format and alignment
		if !utils.CheckTimeAlignment(slot.StartTime) || !utils.CheckTimeAlignment(slot.EndTime) {
			return errors.New("times must align to 30-minute boundaries")
		}

		startMinutes, _ := utils.TimeStrToMinutes(slot.StartTime)
		endMinutes, _ := utils.TimeStrToMinutes(slot.EndTime)
		if startMinutes >= endMinutes {
			return errors.New("start_time must be before end_time")
		}

		availability := &model.Availability{
			CoachID:   coachID,
			DayOfWeek: req.DayOfWeek,
			StartTime: slot.StartTime,
			EndTime:   slot.EndTime,
		}

		_, err := s.availabilityRepo.Create(ctx, availability)
		if err != nil {
			log.Error("failed to create availability")
			return err
		}
	}

	return nil
}

func (s *availabilityService) AddException(ctx context.Context, coachID int64, req *dto.CreateAvailabilityExceptionRequest) error {
	log := logger.GetLogger()

	// Validate coach exists
	coach, err := s.coachRepo.GetByID(ctx, coachID)
	if err != nil {
		log.Error("failed to get coach")
		return err
	}
	if coach == nil {
		return errors.New("coach not found")
	}

	// Parse date
	date, err := utils.DateStringToTime(req.Date)
	if err != nil {
		return errors.New("invalid date format")
	}

	// If times are provided, validate them
	if req.StartTime != nil && req.EndTime != nil {
		if !utils.CheckTimeAlignment(*req.StartTime) || !utils.CheckTimeAlignment(*req.EndTime) {
			return errors.New("times must align to 30-minute boundaries")
		}
		startMinutes, _ := utils.TimeStrToMinutes(*req.StartTime)
		endMinutes, _ := utils.TimeStrToMinutes(*req.EndTime)
		if startMinutes >= endMinutes {
			return errors.New("start_time must be before end_time")
		}
	}

	exception := &model.AvailabilityException{
		CoachID:     coachID,
		Date:        date,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		IsAvailable: req.IsAvailable,
	}

	_, err = s.availabilityExceptionRepo.Create(ctx, exception)
	if err != nil {
		log.Error("failed to create exception")
		return err
	}

	return nil
}

func (s *availabilityService) GetAvailableSlots(ctx context.Context, coachID int64, date string, timezone string) (*dto.GetSlotsResponse, error) {
	log := logger.GetLogger()

	// Validate coach exists
	coach, err := s.coachRepo.GetByID(ctx, coachID)
	if err != nil {
		log.Error("failed to get coach")
		return nil, err
	}
	if coach == nil {
		return nil, errors.New("coach not found")
	}

	// Timezone helper for proper time conversion
	tzHelper := utils.NewTimezoneHelper()

	// Get day range in UTC for the user's requested date in their timezone
	startUTC, endUTC, err := tzHelper.GetDayRangeInUTC(date, timezone)
	if err != nil {
		return nil, fmt.Errorf("invalid date or timezone: %v", err)
	}

	// Check for exceptions in this UTC range
	exception, err := s.availabilityExceptionRepo.GetByCoachAndDate(ctx, coachID, date)
	if err != nil {
		log.Error("failed to get exception")
		return nil, err
	}

	// If there's an exception and day is not available, return empty slots
	if exception != nil && !exception.IsAvailable {
		return &dto.GetSlotsResponse{Slots: []dto.AvailabilitySlot{}}, nil
	}

	// Get slots from exception or regular availability
	var slots []dto.AvailabilitySlot

	if exception != nil && exception.IsAvailable && exception.StartTime != nil && exception.EndTime != nil {
		// Use exception times - generate slots for the user's requested date in their timezone
		slots = s.generateSlotsForTimezone(startUTC, endUTC, *exception.StartTime, *exception.EndTime, timezone, coach.Timezone)
	} else if exception == nil {
		// Use regular availability - get all availability rules for this coach
		availabilities, err := s.availabilityRepo.GetByCoach(ctx, coachID)
		if err != nil {
			log.Error("failed to get availabilities")
			return nil, err
		}

		// Filter and convert availabilities to slots in UTC
		for _, avail := range availabilities {
			slotList := s.generateSlotsForTimezone(startUTC, endUTC, avail.StartTime, avail.EndTime, timezone, coach.Timezone)
			slots = append(slots, slotList...)
		}
	}

	// Remove booked slots
	bookings, err := s.getBookingsForTimeRange(ctx, coachID, startUTC, endUTC)
	if err != nil {
		log.Error("failed to get bookings")
		return nil, err
	}

	slots = s.removeBookedSlots(slots, bookings, nil)

	return &dto.GetSlotsResponse{Slots: slots}, nil
}

func (s *availabilityService) generateSlots(date time.Time, startTime string, endTime string, loc *time.Location) []dto.AvailabilitySlot {
	startMinutes, _ := utils.TimeStrToMinutes(startTime)
	endMinutes, _ := utils.TimeStrToMinutes(endTime)

	var slots []dto.AvailabilitySlot

	for currentMinutes := startMinutes; currentMinutes+30 <= endMinutes; currentMinutes += 30 {
		slotStart := time.Date(date.Year(), date.Month(), date.Day(),
			currentMinutes/60, currentMinutes%60, 0, 0, loc)
		slotEnd := slotStart.Add(30 * time.Minute)

		slots = append(slots, dto.AvailabilitySlot{
			StartTime: slotStart.UTC().Format(time.RFC3339),
			EndTime:   slotEnd.UTC().Format(time.RFC3339),
		})
	}

	return slots
}

// generateSlotsForTimezone generates 30-minute slots considering timezone differences between user and coach
func (s *availabilityService) generateSlotsForTimezone(startUTC, endUTC time.Time, startTimeStr, endTimeStr string, userTimezone, coachTimezone string) []dto.AvailabilitySlot {
	var slots []dto.AvailabilitySlot

	// Load both timezones
	coachLoc, _ := time.LoadLocation(coachTimezone)
	_ = userTimezone // userTimezone available for future enhancements

	// Parse availability hours in coach's timezone
	startMinutes, _ := utils.TimeStrToMinutes(startTimeStr)
	endMinutes, _ := utils.TimeStrToMinutes(endTimeStr)

	// Iterate through the UTC time range and find overlaps
	currentUTC := startUTC
	for currentUTC.Before(endUTC) {
		// Convert current UTC time to coach's timezone
		coachLocalTime := currentUTC.In(coachLoc)
		coachMinutes := coachLocalTime.Hour()*60 + coachLocalTime.Minute()

		// Check if this time falls within coach's working hours
		if coachMinutes >= startMinutes && coachMinutes+30 <= endMinutes {
			// Create a 30-minute slot
			slotEndUTC := currentUTC.Add(30 * time.Minute)

			// Only include if both start and end are within the UTC range
			if slotEndUTC.Before(endUTC) || slotEndUTC.Equal(endUTC) {
				slots = append(slots, dto.AvailabilitySlot{
					StartTime: currentUTC.Format(time.RFC3339),
					EndTime:   slotEndUTC.Format(time.RFC3339),
				})
			}
		}

		currentUTC = currentUTC.Add(30 * time.Minute)
	}

	return slots
}

func (s *availabilityService) getBookingsForDate(ctx context.Context, coachID int64, date time.Time) ([]*model.Booking, error) {
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	endOfDay := startOfDay.Add(24 * time.Hour)

	return s.bookingRepo.GetByCoachAndDateRange(ctx, coachID, startOfDay.Format(time.RFC3339), endOfDay.Format(time.RFC3339))
}

// getBookingsForTimeRange gets all bookings for a coach within a UTC time range
func (s *availabilityService) getBookingsForTimeRange(ctx context.Context, coachID int64, startUTC, endUTC time.Time) ([]*model.Booking, error) {
	return s.bookingRepo.GetByCoachAndDateRange(ctx, coachID, startUTC.Format(time.RFC3339), endUTC.Format(time.RFC3339))
}

func (s *availabilityService) removeBookedSlots(slots []dto.AvailabilitySlot, bookings []*model.Booking, loc *time.Location) []dto.AvailabilitySlot {
	var availableSlots []dto.AvailabilitySlot

	for _, slot := range slots {
		slotStart, _ := utils.ParseRFC3339(slot.StartTime)
		slotEnd, _ := utils.ParseRFC3339(slot.EndTime)

		isBooked := false
		for _, booking := range bookings {
			// Check if slot overlaps with booking
			if slotStart.Before(booking.EndTime) && slotEnd.After(booking.StartTime) {
				isBooked = true
				break
			}
		}

		if !isBooked {
			availableSlots = append(availableSlots, slot)
		}
	}

	return availableSlots
}
