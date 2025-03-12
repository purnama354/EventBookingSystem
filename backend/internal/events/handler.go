package events

import (
	"encoding/json"
	"net/http"
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
		http.Error(w, "Failed to get events", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}

func (h *EventHandler) HandleEvents(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
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
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Input validation
	if strings.TrimSpace(req.Title) == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(req.Date) == "" {
		http.Error(w, "Date is required", http.StatusBadRequest)
		return
	}

	if _, err := time.Parse(time.RFC3339, req.Date); err != nil {
		http.Error(w, "Invalid date format", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(req.Location) == "" {
		http.Error(w, "Location is required", http.StatusBadRequest)
		return
	}

	if req.Capacity <= 0 {
		http.Error(w, "Capacity must be a positive integer", http.StatusBadRequest)
		return
	}

	event, err := h.EventService.CreateEvent(req.Title, req.Description, req.Date, req.Location, req.Capacity)
	if err != nil {
		http.Error(w, "Failed to create event", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(event)
}

func (h *EventHandler) GetEventDetails(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 4 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	eventID := parts[3]
	if eventID == "" {
		http.Error(w, "Event ID is required", http.StatusBadRequest)
		return
	}

	event, err := h.EventService.GetEventByID(eventID)
	if err != nil {
		http.Error(w, "Event not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(event)
}

func (h *EventHandler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 4 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	eventID := parts[3]
	if eventID == "" {
		http.Error(w, "Event ID is required", http.StatusBadRequest)
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
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Input validation
	if strings.TrimSpace(req.Title) == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(req.Date) == "" {
		http.Error(w, "Date is required", http.StatusBadRequest)
		return
	}

	if _, err := time.Parse(time.RFC3339, req.Date); err != nil {
		http.Error(w, "Invalid date format", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(req.Location) == "" {
		http.Error(w, "Location is required", http.StatusBadRequest)
		return
	}

	if req.Capacity <= 0 {
		http.Error(w, "Capacity must be a positive integer", http.StatusBadRequest)
		return
	}

	existingEvent, err := h.EventService.GetEventByID(eventID)
	if err != nil {
		http.Error(w, "Event not found", http.StatusNotFound)
		return
	}

	parsedDate, err := time.Parse(time.RFC3339, req.Date)
	if err != nil {
		http.Error(w, "Invalid date format", http.StatusBadRequest)
		return
	}

	existingEvent.Title = req.Title
	existingEvent.Description = req.Description
	existingEvent.Date = parsedDate
	existingEvent.Location = req.Location
	existingEvent.Capacity = req.Capacity

	err = h.EventService.UpdateEvent(existingEvent)
	if err != nil {
		http.Error(w, "Failed to update event", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(existingEvent)
}

func (h *EventHandler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 4 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	eventID := parts[3]
	if eventID == "" {
		http.Error(w, "Event ID is required", http.StatusBadRequest)
		return
	}

	err := h.EventService.DeleteEvent(eventID)
	if err != nil {
		http.Error(w, "Failed to delete event", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
