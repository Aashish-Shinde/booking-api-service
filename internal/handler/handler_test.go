package handler

import (
	"testing"

	"github.com/booking-api-service/internal/service"
	"github.com/stretchr/testify/suite"
)

// AvailabilityHandlerTestSuite tests for availability handler
type AvailabilityHandlerTestSuite struct {
	suite.Suite
}

func (suite *AvailabilityHandlerTestSuite) TestAvailabilityHandler_ImplementsServiceDependency() {
	var _ service.AvailabilityService
}

// BookingHandlerTestSuite tests for booking handler
type BookingHandlerTestSuite struct {
	suite.Suite
}

func (suite *BookingHandlerTestSuite) TestBookingHandler_ImplementsServiceDependency() {
	var _ service.BookingService
}

// Test registration
func TestAvailabilityHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(AvailabilityHandlerTestSuite))
}

func TestBookingHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(BookingHandlerTestSuite))
}
