package bookings

import (
	"time"

	"gorm.io/gorm"
)

type Booking struct {
	ID        string `gorm:"type:uuid;primaryKey"`
	UserID    string `gorm:"type:uuid;not null"`
	EventID   string `gorm:"type:uuid;not null"`
	Seats     int    `gorm:"not null"`
	Status    string `gorm:"type:varchar(10);default:'booked'"`
	Role      string `gorm:"type:varchar(10);default:'user'"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
