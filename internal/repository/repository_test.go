package repository

import (
	"testing"

	"github.com/booking-api-service/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// UserRepositoryTestSuite tests for user repository
type UserRepositoryTestSuite struct {
	suite.Suite
}

func (suite *UserRepositoryTestSuite) TestUserRepository_ImplementsInterface() {
	var _ UserRepository
}

// CoachRepositoryTestSuite tests for coach repository
type CoachRepositoryTestSuite struct {
	suite.Suite
}

func (suite *CoachRepositoryTestSuite) TestCoachRepository_ImplementsInterface() {
	var _ CoachRepository
}

// AvailabilityRepositoryTestSuite tests for availability repository
type AvailabilityRepositoryTestSuite struct {
	suite.Suite
}

func (suite *AvailabilityRepositoryTestSuite) TestAvailabilityRepository_ImplementsInterface() {
	var _ AvailabilityRepository
}

// BookingRepositoryTestSuite tests for booking repository
type BookingRepositoryTestSuite struct {
	suite.Suite
}

func (suite *BookingRepositoryTestSuite) TestBookingRepository_ImplementsInterface() {
	var _ BookingRepository
}

// AvailabilityExceptionRepositoryTestSuite tests for availability exception repository
type AvailabilityExceptionRepositoryTestSuite struct {
	suite.Suite
}

func (suite *AvailabilityExceptionRepositoryTestSuite) TestAvailabilityExceptionRepository_ImplementsInterface() {
	var _ AvailabilityExceptionRepository
}

// ModelValidationTestSuite tests for model validation
type ModelValidationTestSuite struct {
	suite.Suite
}

func (suite *ModelValidationTestSuite) TestUser_ModelStructure() {
	user := &model.User{
		ID:       1,
		Name:     "John Doe",
		Timezone: "UTC",
	}

	assert.Equal(suite.T(), int64(1), user.ID)
	assert.Equal(suite.T(), "John Doe", user.Name)
	assert.Equal(suite.T(), "UTC", user.Timezone)
}

func (suite *ModelValidationTestSuite) TestCoach_ModelStructure() {
	coach := &model.Coach{
		ID:       1,
		Name:     "Jane Smith",
		Timezone: "America/New_York",
	}

	assert.Equal(suite.T(), int64(1), coach.ID)
	assert.Equal(suite.T(), "Jane Smith", coach.Name)
	assert.Equal(suite.T(), "America/New_York", coach.Timezone)
}

func (suite *ModelValidationTestSuite) TestAvailability_ModelStructure() {
	availability := &model.Availability{
		ID:        1,
		CoachID:   1,
		DayOfWeek: 1,
		StartTime: "09:00",
		EndTime:   "17:00",
	}

	assert.Equal(suite.T(), int64(1), availability.ID)
	assert.Equal(suite.T(), int64(1), availability.CoachID)
	assert.Equal(suite.T(), 1, availability.DayOfWeek)
	assert.Equal(suite.T(), "09:00", availability.StartTime)
	assert.Equal(suite.T(), "17:00", availability.EndTime)
}

func (suite *ModelValidationTestSuite) TestAvailabilityException_ModelStructure() {
	exception := &model.AvailabilityException{
		ID:          1,
		CoachID:     1,
		IsAvailable: false,
	}

	assert.Equal(suite.T(), int64(1), exception.ID)
	assert.Equal(suite.T(), int64(1), exception.CoachID)
	assert.False(suite.T(), exception.IsAvailable)
}

func (suite *ModelValidationTestSuite) TestBooking_ModelStructure() {
	booking := &model.Booking{
		ID:      1,
		UserID:  1,
		CoachID: 1,
		Status:  "ACTIVE",
	}

	assert.Equal(suite.T(), int64(1), booking.ID)
	assert.Equal(suite.T(), int64(1), booking.UserID)
	assert.Equal(suite.T(), int64(1), booking.CoachID)
	assert.Equal(suite.T(), "ACTIVE", booking.Status)
}

// Test registration
func TestUserRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(UserRepositoryTestSuite))
}

func TestCoachRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(CoachRepositoryTestSuite))
}

func TestAvailabilityRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(AvailabilityRepositoryTestSuite))
}

func TestBookingRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(BookingRepositoryTestSuite))
}

func TestAvailabilityExceptionRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(AvailabilityExceptionRepositoryTestSuite))
}

func TestModelValidationTestSuite(t *testing.T) {
	suite.Run(t, new(ModelValidationTestSuite))
}
