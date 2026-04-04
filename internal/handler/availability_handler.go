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

type AvailabilityHandler struct {
	availabilityService service.AvailabilityService
}

func NewAvailabilityHandler(availabilityService service.AvailabilityService) *AvailabilityHandler {
	return &AvailabilityHandler{
		availabilityService: availabilityService,
	}
}

// SetAvailability sets weekly availability for a coach
func (h *AvailabilityHandler) SetAvailability(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()
	coachIDStr := chi.URLParam(r, "coach_id")
	coachID, err := strconv.ParseInt(coachIDStr, 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid coach_id")
		return
	}

	var req dto.CreateAvailabilityRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	err = h.availabilityService.SetAvailability(r.Context(), coachID, &req)
	if err != nil {
		log.Error("failed to set availability")
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondSuccess(w, http.StatusOK, "availability set successfully")
}

// AddException adds a one-time exception to availability
func (h *AvailabilityHandler) AddException(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()
	coachIDStr := chi.URLParam(r, "coach_id")
	coachID, err := strconv.ParseInt(coachIDStr, 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid coach_id")
		return
	}

	var req dto.CreateAvailabilityExceptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	err = h.availabilityService.AddException(r.Context(), coachID, &req)
	if err != nil {
		log.Error("failed to add exception")
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondSuccess(w, http.StatusOK, "exception added successfully")
}

// GetSlots gets available slots for a coach on a specific date
func (h *AvailabilityHandler) GetSlots(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()
	coachIDStr := chi.URLParam(r, "coach_id")
	coachID, err := strconv.ParseInt(coachIDStr, 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid coach_id")
		return
	}

	date := r.URL.Query().Get("date")
	timezone := r.URL.Query().Get("timezone")

	if date == "" || timezone == "" {
		respondError(w, http.StatusBadRequest, "date and timezone query parameters are required")
		return
	}

	slots, err := h.availabilityService.GetAvailableSlots(r.Context(), coachID, date, timezone)
	if err != nil {
		log.Error("failed to get slots")
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, slots)
}

// Helper functions
func respondError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(dto.ErrorResponse{Error: message})
}

func respondSuccess(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(dto.SuccessResponse{Message: message})
}

func respondJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
