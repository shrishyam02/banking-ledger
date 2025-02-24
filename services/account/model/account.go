package model

import (
	"time"

	"github.com/google/uuid"
)

type Account struct {
	ID            uuid.UUID `gorm:"uuid;default:uuid_generate_v4();primaryKey"`
	AccountNumber string    `gorm:"type:varchar(20);not null;uniqueIndex"`
	AccountType   string    `gorm:"type:varchar(20);not null"`
	Status        string    `gorm:"type:varchar(20);not null"`
	Balance       float64   `gorm:"type:decimal(10,2)"`
	CreatedAt     time.Time `gorm:"type:timestamp with time zone"`
	UpdatedAt     time.Time `gorm:"type:timestamp with time zone"`
	CustomerID    uuid.UUID `gorm:"-"`
	Customer      Customer  `gorm:"foreignKey:CustomerID"`
}
