package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/booking-api-service/internal/dto"
	"github.com/booking-api-service/internal/service"
	"github.com/booking-api-service/pkg/logger"
	"github.com/go-chi/chi/v5"
)

type BookingHandler struct {
	bookingService service.BookingService
}

func NewBookingHandler(bookingService service.BookingService) *BookingHandler {
	return &BookingHandler{
		bookingService: bookingService,
	}
}

// CreateBooking creates a new booking
func (h *BookingHandler) CreateBooking(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()

	var req dto.CreateBookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	booking, err := h.bookingService.CreateBooking(r.Context(), &req)
	if err != nil {
		log.Error("failed to create booking")
		if err.Error() == "slot is already booked" {
			respondError(w, http.StatusConflict, err.Error())
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := dto.BookingResponse{
		ID:        booking.ID,
		UserID:    booking.UserID,
		CoachID:   booking.CoachID,
		StartTime: booking.StartTime.Format("2006-01-02T15:04:05Z07:00"),
		EndTime:   booking.EndTime.Format("2006-01-02T15:04:05Z07:00"),
		Status:    booking.Status,
		CreatedAt: booking.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: booking.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	respondJSON(w, http.StatusCreated, response)
}

// GetUserBookings retrieves all bookings for a user
func (h *BookingHandler) GetUserBookings(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()
	userIDStr := chi.URLParam(r, "user_id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid user_id")
		return
	}

	bookings, err := h.bookingService.GetUserBookings(r.Context(), userID)
	if err != nil {
		log.Error("failed to get user bookings")
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	responses := make([]dto.BookingResponse, len(bookings))
	for i, b := range bookings {
		responses[i] = dto.BookingResponse{
			ID:        b.ID,
			UserID:    b.UserID,
			CoachID:   b.CoachID,
			StartTime: b.StartTime.Format("2006-01-02T15:04:05Z07:00"),
			EndTime:   b.EndTime.Format("2006-01-02T15:04:05Z07:00"),
			Status:    b.Status,
			CreatedAt: b.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: b.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	response := dto.BookingsListResponse{
		Bookings: responses,
		Count:    len(responses),
	}

	respondJSON(w, http.StatusOK, response)
}

// ModifyBooking modifies an existing booking
func (h *BookingHandler) ModifyBooking(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()
	bookingIDStr := chi.URLParam(r, "id")
	bookingID, err := strconv.ParseInt(bookingIDStr, 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid booking id")
		return
	}

	var req dto.ModifyBookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	booking, err := h.bookingService.ModifyBooking(r.Context(), bookingID, &req)
	if err != nil {
		log.Error("failed to modify booking")
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := dto.BookingResponse{
		ID:        booking.ID,
		UserID:    booking.UserID,
		CoachID:   booking.CoachID,
		StartTime: booking.StartTime.Format("2006-01-02T15:04:05Z07:00"),
		EndTime:   booking.EndTime.Format("2006-01-02T15:04:05Z07:00"),
		Status:    booking.Status,
		CreatedAt: booking.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: booking.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	respondJSON(w, http.StatusOK, response)
}

// CancelBooking cancels an existing booking
func (h *BookingHandler) CancelBooking(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()
	bookingIDStr := chi.URLParam(r, "id")
	bookingID, err := strconv.ParseInt(bookingIDStr, 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid booking id")
		return
	}

	booking, err := h.bookingService.CancelBooking(r.Context(), bookingID)
	if err != nil {
		log.Error("failed to cancel booking")
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := dto.BookingResponse{
		ID:        booking.ID,
		UserID:    booking.UserID,
		CoachID:   booking.CoachID,
		StartTime: booking.StartTime.Format("2006-01-02T15:04:05Z07:00"),
		EndTime:   booking.EndTime.Format("2006-01-02T15:04:05Z07:00"),
		Status:    booking.Status,
		CreatedAt: booking.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: booking.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	respondJSON(w, http.StatusOK, response)
}
