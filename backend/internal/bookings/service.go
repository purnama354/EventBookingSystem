package bookings

import (
	"github.com/google/uuid"
)

type BookingService interface {
	CreateBooking(userID, eventID string, seats int) (*Booking, error)
	GetBookingByID(id string) (*Booking, error)
	GetBookingsByUserID(userID string) ([]Booking, error)
	CancelBooking(id string) error
}

type BookingServiceImpl struct {
	BookingRepository BookingRepository
}

func NewBookingService(bookingRepository BookingRepository) BookingService {
	return &BookingServiceImpl{BookingRepository: bookingRepository}
}

func (s *BookingServiceImpl) CreateBooking(userID, eventID string, seats int) (*Booking, error) {
	booking := &Booking{
		ID:      uuid.New().String(),
		UserID:  userID,
		EventID: eventID,
		Seats:   seats,
		Status:  "booked",
	}

	err := s.BookingRepository.Create(booking)
	if err != nil {
		return nil, err
	}

	return booking, nil
}

func (s *BookingServiceImpl) GetBookingByID(id string) (*Booking, error) {
	return s.BookingRepository.GetByID(id)
}

func (s *BookingServiceImpl) GetBookingsByUserID(userID string) ([]Booking, error) {
	return s.BookingRepository.GetByUserID(userID)
}

func (s *BookingServiceImpl) CancelBooking(id string) error {
	booking, err := s.BookingRepository.GetByID(id)
	if err != nil {
		return err
	}

	booking.Status = "cancelled"
	return s.BookingRepository.Update(booking)
}
