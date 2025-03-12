package events

import (
	"time"

	"github.com/google/uuid"
)

type EventService interface {
	CreateEvent(title, description string, date string, location string, capacity int) (*Event, error)
	GetEventByID(id string) (*Event, error)
	GetAllEvents() ([]Event, error)
	UpdateEvent(event *Event) error
	DeleteEvent(id string) error
}

type EventServiceImpl struct {
	EventRepository EventRepository
}

func NewEventService(eventRepository EventRepository) EventService {
	return &EventServiceImpl{EventRepository: eventRepository}
}

func (s *EventServiceImpl) CreateEvent(title, description string, date string, location string, capacity int) (*Event, error) {
	parsedDate, err := time.Parse(time.RFC3339, date)
	if err != nil {
		return nil, err
	}
	event := &Event{
		ID:          uuid.New().String(),
		Title:       title,
		Description: description,
		Date:        parsedDate,
		Location:    location,
		Capacity:    capacity,
	}

	err = s.EventRepository.Create(event)
	if err != nil {
		return nil, err
	}

	return event, nil
}

func (s *EventServiceImpl) GetEventByID(id string) (*Event, error) {
	return s.EventRepository.GetByID(id)
}

func (s *EventServiceImpl) GetAllEvents() ([]Event, error) {
	return s.EventRepository.GetAll()
}

func (s *EventServiceImpl) UpdateEvent(event *Event) error {
	return s.EventRepository.Update(event)
}

func (s *EventServiceImpl) DeleteEvent(id string) error {
	return s.EventRepository.Delete(id)
}
