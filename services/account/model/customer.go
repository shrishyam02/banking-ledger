package model

import (
	"time"

	"github.com/google/uuid"
)

type Customer struct {
	ID          uuid.UUID `gorm:"uuid;default:uuid_generate_v4();primaryKey"`
	Name        string    `gorm:"type:varchar(255);not null"`
	Email       string    `gorm:"type:varchar(255);unique"`
	PhoneNumber string    `gorm:"type:varchar(20);unique"`
	Address     string    `gorm:"type:text"`
	CreatedAt   time.Time `gorm:"type:timestamp with time zone"`
	UpdatedAt   time.Time `gorm:"type:timestamp with time zone"`
}
