package services

import (
	"context"

	"github.com/leoguilen/transactions/internal/core/domain"
	"github.com/leoguilen/transactions/internal/core/ports"
)

type TransactionService struct {
	accountRepo     ports.AccountRepository
	transactionRepo ports.TransactionRepository
}

func NewTransactionService(accountRepo ports.AccountRepository, transactionRepo ports.TransactionRepository) ports.TransactionService {
	return &TransactionService{accountRepo: accountRepo, transactionRepo: transactionRepo}
}

func (t *TransactionService) CreateTransaction(ctx context.Context, accountID, operationTypeID int, amount float64) (*domain.Transaction, error) {
	account, err := t.accountRepo.GetByID(ctx, accountID)
	if err != nil {
		return nil, err
	}
	if account == nil {
		return nil, domain.ErrTransactionAccountNotExists
	}

	transaction, err := domain.NewTransaction(accountID, operationTypeID, amount)
	if err != nil {
		return nil, err
	}

	insertedTransaction, err := t.transactionRepo.Insert(ctx, transaction)
	if err != nil {
		return nil, err
	}

	return insertedTransaction, nil
}
