package events

import (
	"encoding/json"
	"eventBookingSystem/internal/types"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type EventHandler struct {
	EventService EventService
}

func NewEventHandler(eventService EventService) *EventHandler {
	return &EventHandler{EventService: eventService}
}

func (h *EventHandler) ListEvents(w http.ResponseWriter, r *http.Request) {
	events, err := h.EventService.GetAllEvents()
	if err != nil {
		types.SendError(w, http.StatusInternalServerError, "RETRIEVAL_FAILED", "Failed to get events", nil)
		return
	}

	types.SendSuccess(w, http.StatusOK, "Events retrieved successfully", events)
}

func (h *EventHandler) HandleEvents(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if r.URL.Query().Has("search") || r.URL.Query().Has("startDate") ||
			r.URL.Query().Has("endDate") || r.URL.Query().Has("minSeats") {
			h.SearchEvents(w, r)
			return
		}
		// Extract event ID from the URL path
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) > 3 {
			eventID := parts[3]
			if eventID != "" {
				h.GetEventDetails(w, r)
				return
			}
		}
		h.ListEvents(w, r)
	case http.MethodPost:
		h.CreateEvent(w, r)
	case http.MethodPut:
		h.UpdateEvent(w, r)
	case http.MethodDelete:
		h.DeleteEvent(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *EventHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Date        string `json:"date"`
		Location    string `json:"location"`
		Capacity    int    `json:"capacity"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		types.SendError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body", nil)
		return
	}

	// Validation errors map
	validationErrors := make(map[string]string)

	// Title validation
	if title := strings.TrimSpace(req.Title); title == "" {
		validationErrors["title"] = "Title is required"
	} else if len(title) < 3 {
		validationErrors["title"] = "Title must be at least 3 characters long"
	}

	// Date validation
	if date := strings.TrimSpace(req.Date); date == "" {
		validationErrors["date"] = "Date is required"
	} else if _, err := time.Parse(time.RFC3339, date); err != nil {
		validationErrors["date"] = "Date must be in RFC3339 format (e.g., 2024-03-15T14:00:00Z)"
	} else {
		// Check if date is in the future
		parsedDate, _ := time.Parse(time.RFC3339, date)
		if parsedDate.Before(time.Now()) {
			validationErrors["date"] = "Event date must be in the future"
		}
	}

	// Location validation
	if location := strings.TrimSpace(req.Location); location == "" {
		validationErrors["location"] = "Location is required"
	} else if len(location) < 3 {
		validationErrors["location"] = "Location must be at least 3 characters long"
	}

	// Capacity validation
	if req.Capacity <= 0 {
		validationErrors["capacity"] = "Capacity must be a positive integer"
	} else if req.Capacity > 10000 { // Example maximum capacity
		validationErrors["capacity"] = "Capacity cannot exceed 10000"
	}

	// Description validation (optional field)
	if desc := strings.TrimSpace(req.Description); len(desc) > 1000 {
		validationErrors["description"] = "Description cannot exceed 1000 characters"
	}

	// If there are validation errors, return them
	if len(validationErrors) > 0 {
		types.SendError(w, http.StatusBadRequest, "VALIDATION_FAILED", "Validation failed", validationErrors)
		return
	}

	event, err := h.EventService.CreateEvent(
		strings.TrimSpace(req.Title),
		strings.TrimSpace(req.Description),
		req.Date,
		strings.TrimSpace(req.Location),
		req.Capacity,
	)
	if err != nil {
		types.SendError(w, http.StatusInternalServerError, "EVENT_CREATION_FAILED", "Failed to create event", nil)
		return
	}

	types.SendSuccess(w, http.StatusCreated, "Event created successfully", event)
}

func (h *EventHandler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 4 {
		types.SendError(w, http.StatusBadRequest, "INVALID_URL", "Invalid URL format", nil)
		return
	}

	eventID := parts[3]
	if eventID == "" {
		types.SendError(w, http.StatusBadRequest, "MISSING_ID", "Event ID is required", nil)
		return
	}

	var req struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Date        string `json:"date"`
		Location    string `json:"location"`
		Capacity    int    `json:"capacity"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		types.SendError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body", nil)
		return
	}

	// Validation errors map
	validationErrors := make(map[string]string)
	var parsedDate time.Time

	// Title validation
	if title := strings.TrimSpace(req.Title); title == "" {
		validationErrors["title"] = "Title is required"
	} else if len(title) < 3 {
		validationErrors["title"] = "Title must be at least 3 characters long"
	}

	// Date validation
	if date := strings.TrimSpace(req.Date); date == "" {
		validationErrors["date"] = "Date is required"
	} else if t, err := time.Parse(time.RFC3339, date); err != nil {
		validationErrors["date"] = "Date must be in RFC3339 format (e.g., 2024-03-15T14:00:00Z)"
	} else {
		parsedDate = t
		if t.Before(time.Now()) {
			validationErrors["date"] = "Event date must be in the future"
		}
	}

	// Location validation
	if location := strings.TrimSpace(req.Location); location == "" {
		validationErrors["location"] = "Location is required"
	} else if len(location) < 3 {
		validationErrors["location"] = "Location must be at least 3 characters long"
	}

	// Capacity validation
	if req.Capacity <= 0 {
		validationErrors["capacity"] = "Capacity must be a positive integer"
	} else if req.Capacity > 10000 {
		validationErrors["capacity"] = "Capacity cannot exceed 10000"
	}

	// Description validation (optional field)
	if desc := strings.TrimSpace(req.Description); len(desc) > 1000 {
		validationErrors["description"] = "Description cannot exceed 1000 characters"
	}

	// If there are validation errors, return them
	if len(validationErrors) > 0 {
		types.SendError(w, http.StatusBadRequest, "VALIDATION_FAILED", "Validation failed", validationErrors)
		return
	}

	// Get existing event
	existingEvent, err := h.EventService.GetEventByID(eventID)
	if err != nil {
		types.SendError(w, http.StatusNotFound, "EVENT_NOT_FOUND", "Event not found", nil)
		return
	}

	// Update event fields
	existingEvent.Title = strings.TrimSpace(req.Title)
	existingEvent.Description = strings.TrimSpace(req.Description)
	existingEvent.Date = parsedDate
	existingEvent.Location = strings.TrimSpace(req.Location)
	existingEvent.Capacity = req.Capacity

	// Save updated event
	if err := h.EventService.UpdateEvent(existingEvent); err != nil {
		types.SendError(w, http.StatusInternalServerError, "UPDATE_FAILED", "Failed to update event", nil)
		return
	}

	types.SendSuccess(w, http.StatusOK, "Event updated successfully", existingEvent)
}

func (h *EventHandler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 4 {
		types.SendError(w, http.StatusBadRequest, "INVALID_URL", "Invalid URL format", nil)
		return
	}

	eventID := parts[3]
	if eventID == "" {
		types.SendError(w, http.StatusBadRequest, "MISSING_ID", "Event ID is required", nil)
		return
	}

	// Verify event exists before deletion
	_, err := h.EventService.GetEventByID(eventID)
	if err != nil {
		types.SendError(w, http.StatusNotFound, "EVENT_NOT_FOUND", "Event not found", nil)
		return
	}

	err = h.EventService.DeleteEvent(eventID)
	if err != nil {
		types.SendError(w, http.StatusInternalServerError, "DELETION_FAILED", "Failed to delete event", nil)
		return
	}

	types.SendSuccess(w, http.StatusOK, "Event deleted successfully", nil)
}

func (h *EventHandler) GetEventDetails(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 4 {
		types.SendError(w, http.StatusBadRequest, "INVALID_URL", "Invalid URL format", nil)
		return
	}

	eventID := parts[3]
	if eventID == "" {
		types.SendError(w, http.StatusBadRequest, "MISSING_ID", "Event ID is required", nil)
		return
	}

	event, err := h.EventService.GetEventByID(eventID)
	if err != nil {
		types.SendError(w, http.StatusNotFound, "EVENT_NOT_FOUND", "Event not found", nil)
		return
	}

	types.SendSuccess(w, http.StatusOK, "Event retrieved successfully", event)
}

func (h *EventHandler) SearchEvents(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	search := query.Get("search")

	validationErrors := make(map[string]string)
	var startDate, endDate *time.Time

	if s := query.Get("startDate"); s != "" {
		t, err := time.Parse(time.RFC3339, s)
		if err != nil {
			validationErrors["startDate"] = "Invalid date format, use RFC3339 format"
		} else {
			startDate = &t
		}
	}

	if e := query.Get("endDate"); e != "" {
		t, err := time.Parse(time.RFC3339, e)
		if err != nil {
			validationErrors["endDate"] = "Invalid date format, use RFC3339 format"
		} else {
			endDate = &t
		}
	}

	var minSeats *int
	if s := query.Get("minSeats"); s != "" {
		seats, err := strconv.Atoi(s)
		if err != nil {
			validationErrors["minSeats"] = "Must be a valid integer"
		} else if seats < 0 {
			validationErrors["minSeats"] = "Must be a positive integer"
		} else {
			minSeats = &seats
		}
	}

	if len(validationErrors) > 0 {
		types.SendError(w, http.StatusBadRequest, "VALIDATION_FAILED", "Invalid search parameters", validationErrors)
		return
	}

	events, err := h.EventService.SearchEvents(search, startDate, endDate, minSeats)
	if err != nil {
		types.SendError(w, http.StatusInternalServerError, "SEARCH_FAILED", "Failed to search events", nil)
		return
	}

	types.SendSuccess(w, http.StatusOK, "Events search completed", events)
}
