package ports

import (
	"context"

	"github.com/leoguilen/transactions/internal/core/domain"
)

type TransactionService interface {
	CreateTransaction(ctx context.Context, accountID, operationTypeID int, amount float64) (*domain.Transaction, error)
}

type TransactionRepository interface {
	Insert(ctx context.Context, transaction *domain.Transaction) (*domain.Transaction, error)
}
