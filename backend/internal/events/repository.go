package events

import (
	"time"

	"gorm.io/gorm"
)

type EventFilter struct {
	Search    string     // Search in title and location
	StartDate *time.Time // Filter events after this date
	EndDate   *time.Time // Filter events before this date
	MinSeats  *int       // Filter events with at least this many seats
}

type EventRepository interface {
	Create(event *Event) error
	GetByID(id string) (*Event, error)
	GetAll() ([]Event, error)
	Update(event *Event) error
	Delete(id string) error
	Search(filter EventFilter) ([]Event, error)
}

type EventRepositoryImpl struct {
	DB *gorm.DB
}

func NewEventRepository(db *gorm.DB) EventRepository {
	return &EventRepositoryImpl{DB: db}
}

func (r *EventRepositoryImpl) Create(event *Event) error {
	return r.DB.Create(event).Error
}

func (r *EventRepositoryImpl) GetByID(id string) (*Event, error) {
	var event Event
	err := r.DB.First(&event, "id = ?", id).Error
	return &event, err
}

func (r *EventRepositoryImpl) GetAll() ([]Event, error) {
	var events []Event
	err := r.DB.Find(&events).Error
	return events, err
}

func (r *EventRepositoryImpl) Update(event *Event) error {
	return r.DB.Save(event).Error
}

func (r *EventRepositoryImpl) Delete(id string) error {
	return r.DB.Delete(&Event{}, "id = ?", id).Error
}

func (r *EventRepositoryImpl) Search(filter EventFilter) ([]Event, error) {
	var events []Event
	query := r.DB.Model(&Event{})

	if filter.Search != "" {
		search := "%" + filter.Search + "%"
		query = query.Where("title ILIKE ? OR location ILIKE ?", search, search)
	}

	if filter.StartDate != nil {
		query = query.Where("date >= ?", filter.StartDate)
	}

	if filter.EndDate != nil {
		query = query.Where("date <= ?", filter.EndDate)
	}

	if filter.MinSeats != nil {
		query = query.Where("capacity >= ?", *filter.MinSeats)
	}

	err := query.Find(&events).Error
	return events, err
}
