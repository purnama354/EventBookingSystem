package bookings

import (
	"encoding/json"
	"eventBookingSystem/internal/middleware"
	"eventBookingSystem/internal/types"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

type BookingHandler struct {
	BookingService BookingService
}

func NewBookingHandler(bookingService BookingService) *BookingHandler {
	return &BookingHandler{BookingService: bookingService}
}

func (h *BookingHandler) CreateBooking(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		types.SendError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed", nil)
		return
	}

	var req struct {
		EventID string `json:"eventID"`
		Seats   int    `json:"seats"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		types.SendError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body", nil)
		return
	}

	// Input validation
	validationErrors := make(map[string]string)

	if _, err := uuid.Parse(req.EventID); err != nil {
		validationErrors["eventID"] = "Invalid event ID format"
	}

	if req.Seats <= 0 {
		validationErrors["seats"] = "Seats must be a positive integer"
	}

	if len(validationErrors) > 0 {
		types.SendError(w, http.StatusBadRequest, "VALIDATION_FAILED", "Validation failed", validationErrors)
		return
	}

	// Get the user ID from the request context
	userID := r.Context().Value(middleware.UserIDKey).(string)

	booking, err := h.BookingService.CreateBooking(userID, req.EventID, req.Seats)
	if err != nil {
		types.SendError(w, http.StatusInternalServerError, "BOOKING_FAILED", "Failed to create booking", nil)
		return
	}

	types.SendSuccess(w, http.StatusCreated, "Booking created successfully", booking)
}

func (h *BookingHandler) GetBookingByID(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 4 {
		types.SendError(w, http.StatusBadRequest, "INVALID_URL", "Invalid URL format", nil)
		return
	}

	bookingID := parts[3]
	if bookingID == "" {
		types.SendError(w, http.StatusBadRequest, "MISSING_ID", "Booking ID is required", nil)
		return
	}

	if _, err := uuid.Parse(bookingID); err != nil {
		types.SendError(w, http.StatusBadRequest, "INVALID_ID", "Invalid booking ID format", nil)
		return
	}

	booking, err := h.BookingService.GetBookingByID(bookingID)
	if err != nil {
		types.SendError(w, http.StatusNotFound, "BOOKING_NOT_FOUND", "Booking not found", nil)
		return
	}

	types.SendSuccess(w, http.StatusOK, "Booking retrieved successfully", booking)
}

func (h *BookingHandler) GetBookingsByUserID(w http.ResponseWriter, r *http.Request) {
	// Get the user ID from the request context
	userID := r.Context().Value(middleware.UserIDKey).(string)

	// Input validation
	if _, err := uuid.Parse(userID); err != nil {
		types.SendError(w, http.StatusBadRequest, "INVALID_USER_ID", "Invalid user ID format", nil)
		return
	}

	bookings, err := h.BookingService.GetBookingsByUserID(userID)
	if err != nil {
		types.SendError(w, http.StatusInternalServerError, "RETRIEVAL_FAILED", "Failed to get bookings", nil)
		return
	}

	types.SendSuccess(w, http.StatusOK, "Bookings retrieved successfully", bookings)
}

func (h *BookingHandler) CancelBooking(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 4 {
		types.SendError(w, http.StatusBadRequest, "INVALID_URL", "Invalid URL format", nil)
		return
	}

	bookingID := parts[3]
	if bookingID == "" {
		types.SendError(w, http.StatusBadRequest, "MISSING_ID", "Booking ID is required", nil)
		return
	}

	if _, err := uuid.Parse(bookingID); err != nil {
		types.SendError(w, http.StatusBadRequest, "INVALID_ID", "Invalid booking ID format", nil)
		return
	}

	err := h.BookingService.CancelBooking(bookingID)
	if err != nil {
		types.SendError(w, http.StatusInternalServerError, "CANCELLATION_FAILED", "Failed to cancel booking", nil)
		return
	}

	types.SendSuccess(w, http.StatusOK, "Booking cancelled successfully", nil)
}

func (h *BookingHandler) HandleBookings(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.CreateBooking(w, r)
	case http.MethodGet:
		parts := strings.Split(r.URL.Path, "/")
		fmt.Println(parts)
		if len(parts) > 3 {
			if parts[2] == "users" {
				h.GetBookingsByUserID(w, r)
				return
			}
			bookingID := parts[3]
			if bookingID != "" {
				h.GetBookingByID(w, r)
				return
			}
		}
		types.SendError(w, http.StatusBadRequest, "INVALID_URL", "Invalid URL", nil)
		return
	case http.MethodDelete:
		h.CancelBooking(w, r)
	default:
		types.SendError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed", nil)
	}
}
