package model

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	ID              uuid.UUID `json:"id"`
	AccountID       uuid.UUID `json:"accountId"`
	Amount          float64   `json:"amount"`
	TransactionType string    `json:"transactionType"` // e.g., "credit", "debit"
	Details         string    `json:"details"`
	AcceptedAt      time.Time `json:"acceptedAt"`
}
