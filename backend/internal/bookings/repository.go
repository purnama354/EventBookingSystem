package bookings

import (
	"gorm.io/gorm"
)

type BookingRepository interface {
	Create(booking *Booking) error
	GetByID(id string) (*Booking, error)
	GetByUserID(userID string) ([]Booking, error)
	GetByEventID(eventID string) ([]Booking, error)
	Update(booking *Booking) error
	Delete(id string) error
}

type BookingRepositoryImpl struct {
	DB *gorm.DB
}

func NewBookingRepository(db *gorm.DB) BookingRepository {
	return &BookingRepositoryImpl{DB: db}
}

func (r *BookingRepositoryImpl) Create(booking *Booking) error {
	return r.DB.Create(booking).Error
}

func (r *BookingRepositoryImpl) GetByID(id string) (*Booking, error) {
	var booking Booking
	err := r.DB.First(&booking, "id = ?", id).Error
	return &booking, err
}

func (r *BookingRepositoryImpl) GetByUserID(userID string) ([]Booking, error) {
	var bookings []Booking
	err := r.DB.Where("user_id = ?", userID).Find(&bookings).Error
	return bookings, err
}

func (r *BookingRepositoryImpl) GetByEventID(eventID string) ([]Booking, error) {
	var bookings []Booking
	err := r.DB.Where("event_id = ?", eventID).Find(&bookings).Error
	return bookings, err
}

func (r *BookingRepositoryImpl) Update(booking *Booking) error {
	return r.DB.Save(booking).Error
}

func (r *BookingRepositoryImpl) Delete(id string) error {
	return r.DB.Delete(&Booking{}, "id = ?", id).Error
}
