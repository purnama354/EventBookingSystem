package events

import (
	"time"

	"gorm.io/gorm"
)

type Event struct {
	ID          string    `gorm:"type:uuid;primaryKey"`
	Title       string    `gorm:"not null"`
	Description string    `gorm:"type:text"`
	Date        time.Time `gorm:"not null"`
	Location    string    `gorm:"not null"`
	Capacity    int       `gorm:"not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}
