package service

import (
	"testing"
	"time"

	"github.com/booking-api-service/internal/dto"
	"github.com/booking-api-service/internal/model"
)

func TestSlotGeneration(t *testing.T) {
	// Test slot generation logic
	service := &availabilityService{}

	date := time.Date(2026, time.April, 7, 0, 0, 0, 0, time.UTC)
	loc := time.UTC

	slots := service.generateSlots(date, "10:00", "13:00", loc)

	// Should generate 6 slots (10:00-10:30, 10:30-11:00, 11:00-11:30, 11:30-12:00, 12:00-12:30, 12:30-13:00)
	expectedCount := 6
	if len(slots) != expectedCount {
		t.Errorf("expected %d slots, got %d", expectedCount, len(slots))
	}

	// Check first slot
	if slots[0].StartTime != "2026-04-07T10:00:00Z" {
		t.Errorf("expected first slot start time to be 2026-04-07T10:00:00Z, got %s", slots[0].StartTime)
	}

	// Check last slot
	if slots[len(slots)-1].EndTime != "2026-04-07T13:00:00Z" {
		t.Errorf("expected last slot end time to be 2026-04-07T13:00:00Z, got %s", slots[len(slots)-1].EndTime)
	}
}

func TestRemoveBookedSlots(t *testing.T) {
	service := &availabilityService{}

	slots := []dto.AvailabilitySlot{
		{StartTime: "2026-04-07T10:00:00Z", EndTime: "2026-04-07T10:30:00Z"},
		{StartTime: "2026-04-07T10:30:00Z", EndTime: "2026-04-07T11:00:00Z"},
		{StartTime: "2026-04-07T11:00:00Z", EndTime: "2026-04-07T11:30:00Z"},
		{StartTime: "2026-04-07T11:30:00Z", EndTime: "2026-04-07T12:00:00Z"},
	}

	// Create a booking for the second slot
	bookings := []*model.Booking{
		{
			StartTime: time.Date(2026, time.April, 7, 10, 30, 0, 0, time.UTC),
			EndTime:   time.Date(2026, time.April, 7, 11, 0, 0, 0, time.UTC),
		},
	}

	result := service.removeBookedSlots(slots, bookings, time.UTC)

	// Should have 3 available slots (removed the booked one)
	expectedCount := 3
	if len(result) != expectedCount {
		t.Errorf("expected %d available slots, got %d", expectedCount, len(result))
	}

	// Verify the booked slot was removed
	for _, slot := range result {
		if slot.StartTime == "2026-04-07T10:30:00Z" {
			t.Error("booked slot was not removed")
		}
	}
}

// Note: These require the service package to export dto and model
// For actual integration tests, you would use a test database
