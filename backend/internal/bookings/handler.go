package bookings

import (
	"encoding/json"
	"eventBookingSystem/internal/middleware"
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
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		EventID string `json:"eventID"`
		Seats   int    `json:"seats"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Input validation
	if _, err := uuid.Parse(req.EventID); err != nil {
		http.Error(w, "Invalid event ID", http.StatusBadRequest)
		return
	}

	if req.Seats <= 0 {
		http.Error(w, "Seats must be a positive integer", http.StatusBadRequest)
		return
	}

	// Get the user ID from the request context
	userID := r.Context().Value(middleware.UserIDKey).(string)

	booking, err := h.BookingService.CreateBooking(userID, req.EventID, req.Seats)
	if err != nil {
		http.Error(w, "Failed to create booking", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(booking)
}

func (h *BookingHandler) GetBookingByID(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 4 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	bookingID := parts[3]
	if bookingID == "" {
		http.Error(w, "Booking ID is required", http.StatusBadRequest)
		return
	}

	if _, err := uuid.Parse(bookingID); err != nil {
		http.Error(w, "Invalid booking ID", http.StatusBadRequest)
		return
	}

	booking, err := h.BookingService.GetBookingByID(bookingID)
	if err != nil {
		http.Error(w, "Booking not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(booking)
}

func (h *BookingHandler) GetBookingsByUserID(w http.ResponseWriter, r *http.Request) {
	// Get the user ID from the request context
	userID := r.Context().Value(middleware.UserIDKey).(string)

	// Input validation
	if _, err := uuid.Parse(userID); err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	bookings, err := h.BookingService.GetBookingsByUserID(userID)
	if err != nil {
		http.Error(w, "Failed to get bookings", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bookings)
}

func (h *BookingHandler) CancelBooking(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 4 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	bookingID := parts[3]
	if bookingID == "" {
		http.Error(w, "Booking ID is required", http.StatusBadRequest)
		return
	}

	if _, err := uuid.Parse(bookingID); err != nil {
		http.Error(w, "Invalid booking ID", http.StatusBadRequest)
		return
	}

	err := h.BookingService.CancelBooking(bookingID)
	if err != nil {
		http.Error(w, "Failed to cancel booking", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
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
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	case http.MethodDelete:
		h.CancelBooking(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// func isValidUUID(u string) bool {
// 	_, err := uuid.Parse(u)
// 	return err == nil
// }
