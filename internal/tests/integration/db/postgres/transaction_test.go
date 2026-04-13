package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/leoguilen/transactions/internal/adapters/db/postgres"
	"github.com/leoguilen/transactions/internal/core/domain"
	"github.com/leoguilen/transactions/internal/core/ports"
)

func newTransactionRepo() ports.TransactionRepository {
	return postgres.NewTransactionRepository(ConnStr)
}

func TestTransactionRepository_Insert_Success(t *testing.T) {
	repo := newTransactionRepo()

	transaction := &domain.Transaction{
		AccountID:       1,
		OperationTypeID: domain.OperationTypePurchase,
		Amount:          -100.0,
		EventDate:       time.Now(),
	}

	result, err := repo.Insert(context.Background(), transaction)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result == nil {
		t.Fatalf("expected transaction, got nil")
	}
	if result.ID == 0 {
		t.Errorf("expected non-zero ID, got %d", result.ID)
	}
	if result.AccountID != transaction.AccountID {
		t.Errorf("expected AccountID %d, got %d", transaction.AccountID, result.AccountID)
	}
}

func TestTransactionRepository_Insert_ConnectionError(t *testing.T) {
	repo := newTransactionRepo()

	transaction := &domain.Transaction{
		AccountID:       1,
		OperationTypeID: domain.OperationTypePurchase,
		Amount:          -100.0,
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := repo.Insert(ctx, transaction)

	if err == nil {
		t.Errorf("expected error for cancelled context, got nil")
	}
}
