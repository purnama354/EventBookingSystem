package users

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           string `gorm:"type:uuid;primaryKey"`
	Username     string `gorm:"uniqueIndex;not null"`
	Email        string `gorm:"uniqueIndex;not null"`
	PasswordHash string `gorm:"not null"`
	Role         string `gorm:"type:varchar(10);default:'user'"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}
