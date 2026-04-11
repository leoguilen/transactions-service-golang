package repository

import "github.com/leoguilen/transactions/internal/domain"

// TransactionRepository defines persistence operations for transactions (port)
type TransactionRepository interface {
	Create(tx *domain.Transaction) (*domain.Transaction, error)
	ListByAccount(accountID string, limit, offset int) ([]*domain.Transaction, error)
	SumByAccount(accountID string) (string, error)
}
