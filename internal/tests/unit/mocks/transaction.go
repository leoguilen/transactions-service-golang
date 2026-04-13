package mocks

import (
	"context"

	"github.com/leoguilen/transactions/internal/core/domain"
)

type MockTransactionService struct {
	CreateTransactionFn func(ctx context.Context, accountID, operationTypeID int, amount float64) (*domain.Transaction, error)
}

func (m *MockTransactionService) CreateTransaction(ctx context.Context, accountID, operationTypeID int, amount float64) (*domain.Transaction, error) {
	return m.CreateTransactionFn(ctx, accountID, operationTypeID, amount)
}

type MockTransactionRepository struct {
	InsertFn func(ctx context.Context, tr *domain.Transaction) (*domain.Transaction, error)
}

func (m *MockTransactionRepository) Insert(ctx context.Context, tr *domain.Transaction) (*domain.Transaction, error) {
	if m.InsertFn != nil {
		return m.InsertFn(ctx, tr)
	}
	return nil, nil
}
