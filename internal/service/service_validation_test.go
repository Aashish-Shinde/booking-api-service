package service

import (
	"testing"
	"time"

	"github.com/booking-api-service/internal/dto"
	"github.com/booking-api-service/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// AvailabilityServiceTestSuite tests for availability service
type AvailabilityServiceTestSuite struct {
	suite.Suite
}

func (suite *AvailabilityServiceTestSuite) TestSlotGeneration() {
	service := &availabilityService{}
	date := time.Date(2026, time.April, 7, 0, 0, 0, 0, time.UTC)
	loc := time.UTC

	slots := service.generateSlots(date, "10:00", "13:00", loc)

	assert.Len(suite.T(), slots, 6)
	assert.Equal(suite.T(), "2026-04-07T10:00:00Z", slots[0].StartTime)
	assert.Equal(suite.T(), "2026-04-07T13:00:00Z", slots[len(slots)-1].EndTime)
}

func (suite *AvailabilityServiceTestSuite) TestRemoveBookedSlots() {
	service := &availabilityService{}
	slots := []dto.AvailabilitySlot{
		{StartTime: "2026-04-07T10:00:00Z", EndTime: "2026-04-07T10:30:00Z"},
		{StartTime: "2026-04-07T10:30:00Z", EndTime: "2026-04-07T11:00:00Z"},
		{StartTime: "2026-04-07T11:00:00Z", EndTime: "2026-04-07T11:30:00Z"},
	}

	bookings := []*model.Booking{
		{
			StartTime: time.Date(2026, time.April, 7, 10, 30, 0, 0, time.UTC),
			EndTime:   time.Date(2026, time.April, 7, 11, 0, 0, 0, time.UTC),
		},
	}

	result := service.removeBookedSlots(slots, bookings, nil)

	assert.Len(suite.T(), result, 2)
	assert.Equal(suite.T(), "2026-04-07T10:00:00Z", result[0].StartTime)
	assert.Equal(suite.T(), "2026-04-07T11:00:00Z", result[1].StartTime)
}

func TestAvailabilityServiceTestSuite(t *testing.T) {
	suite.Run(t, new(AvailabilityServiceTestSuite))
}

// DTOValidationTestSuite tests for DTO validation
type DTOValidationTestSuite struct {
	suite.Suite
}

func (suite *DTOValidationTestSuite) TestTimeSlotInput_ValidateSuccess() {
	slot := &dto.TimeSlotInput{
		StartTime: "09:00",
		EndTime:   "09:30",
	}

	err := slot.Validate()
	assert.NoError(suite.T(), err)
}

func (suite *DTOValidationTestSuite) TestTimeSlotInput_ValidateInvalidFormat() {
	slot := &dto.TimeSlotInput{
		StartTime: "9:00",
		EndTime:   "09:30",
	}

	err := slot.Validate()
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "HH:MM")
}

func (suite *DTOValidationTestSuite) TestTimeSlotInput_ValidateInvalidBoundary() {
	slot := &dto.TimeSlotInput{
		StartTime: "09:15",
		EndTime:   "09:45",
	}

	err := slot.Validate()
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "30-minute boundaries")
}

func (suite *DTOValidationTestSuite) TestTimeSlotInput_ValidateInvalidDuration() {
	slot := &dto.TimeSlotInput{
		StartTime: "09:00",
		EndTime:   "10:00",
	}

	err := slot.Validate()
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "30 minutes")
}

func (suite *DTOValidationTestSuite) TestCreateAvailabilityExceptionRequest_Validate() {
	req := &dto.CreateAvailabilityExceptionRequest{
		Date:        "2026-04-12",
		IsAvailable: false,
	}

	err := req.Validate()
	assert.NoError(suite.T(), err)
}

func (suite *DTOValidationTestSuite) TestCreateAvailabilityExceptionRequest_ValidatePartialTimes() {
	req := &dto.CreateAvailabilityExceptionRequest{
		Date:        "2026-04-12",
		IsAvailable: true,
		StartTime:   stringPtr("10:00"),
	}

	err := req.Validate()
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "together")
}

func stringPtr(s string) *string {
	return &s
}

func TestDTOValidationTestSuite(t *testing.T) {
	suite.Run(t, new(DTOValidationTestSuite))
}
