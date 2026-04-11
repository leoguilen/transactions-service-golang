package domain

import "time"

// Account represents a cardholder account
type Account struct {
	AccountID      string    `json:"account_id" db:"account_id"`
	DocumentNumber string    `json:"document_number" db:"document_number"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}
