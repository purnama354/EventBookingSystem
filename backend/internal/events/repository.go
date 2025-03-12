package events

import (
	"gorm.io/gorm"
)

type EventRepository interface {
	Create(event *Event) error
	GetByID(id string) (*Event, error)
	GetAll() ([]Event, error)
	Update(event *Event) error
	Delete(id string) error
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
