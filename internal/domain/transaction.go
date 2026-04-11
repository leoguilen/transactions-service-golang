package domain

import "time"

// TransactionType represents the type of transaction
type TransactionType string

const (
	TransactionTypePurchase            TransactionType = "purchase"
	TransactionTypeWithdrawal          TransactionType = "withdrawal"
	TransactionTypeCreditVoucher       TransactionType = "credit_voucher"
	TransactionTypePurchaseInstallment TransactionType = "purchase_installments"
)

// Transaction represents a financial operation linked to an account
type Transaction struct {
	TransactionID   string    `json:"transaction_id" db:"transaction_id"`
	AccountID       string    `json:"account_id" db:"account_id"`
	OperationTypeID int       `json:"operation_type_id" db:"operation_type_id"`
	Amount          string    `json:"amount" db:"amount"` // use string to preserve NUMERIC precision
	EventDate       time.Time `json:"event_date" db:"event_date"`
	Description     string    `json:"description,omitempty" db:"description"`
	Metadata        string    `json:"metadata,omitempty" db:"metadata"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
}
