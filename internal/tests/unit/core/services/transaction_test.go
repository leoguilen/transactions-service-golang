package services

import (
	"context"
	"testing"

	"github.com/leoguilen/transactions/internal/core/domain"
	"github.com/leoguilen/transactions/internal/core/services"
	"github.com/leoguilen/transactions/internal/tests/unit/mocks"
)

func TestTransactionService_CreateTransaction_Success(t *testing.T) {
	mockAccountRepo := &mocks.MockAccountRepository{
		GetByIDFn: func(ctx context.Context, id int) (*domain.Account, error) {
			return &domain.Account{ID: 1, DocumentNumber: "12345678900"}, nil
		},
	}
	mockTransactionRepo := &mocks.MockTransactionRepository{
		InsertFn: func(ctx context.Context, transaction *domain.Transaction) (*domain.Transaction, error) {
			return &domain.Transaction{ID: 1, AccountID: 1, OperationTypeID: domain.OperationTypePurchase, Amount: -100.0}, nil
		},
	}
	service := services.NewTransactionService(mockAccountRepo, mockTransactionRepo)

	tr, err := service.CreateTransaction(context.Background(), 1, int(domain.OperationTypePurchase), 100.0)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if tr == nil {
		t.Fatalf("expected transaction, got nil")
	}
	if tr.AccountID != 1 {
		t.Errorf("expected AccountID 1, got %d", tr.AccountID)
	}
}

func TestTransactionService_CreateTransaction_AccountNotFound(t *testing.T) {
	mockAccountRepo := &mocks.MockAccountRepository{
		GetByIDFn: func(ctx context.Context, id int) (*domain.Account, error) {
			return nil, nil
		},
	}
	mockTransactionRepo := &mocks.MockTransactionRepository{}
	service := services.NewTransactionService(mockAccountRepo, mockTransactionRepo)

	_, err := service.CreateTransaction(context.Background(), 999, int(domain.OperationTypePurchase), 100.0)

	if err != domain.ErrTransactionAccountNotExists {
		t.Errorf("expected ErrTransactionAccountNotExists, got %v", err)
	}
}

func TestTransactionService_CreateTransaction_AccountRepoError(t *testing.T) {
	mockAccountRepo := &mocks.MockAccountRepository{
		GetByIDFn: func(ctx context.Context, id int) (*domain.Account, error) {
			return nil, domain.ErrTransactionAccountNotExists
		},
	}
	mockTransactionRepo := &mocks.MockTransactionRepository{}
	service := services.NewTransactionService(mockAccountRepo, mockTransactionRepo)

	_, err := service.CreateTransaction(context.Background(), 1, int(domain.OperationTypePurchase), 100.0)

	if err != domain.ErrTransactionAccountNotExists {
		t.Errorf("expected error, got %v", err)
	}
}

func TestTransactionService_CreateTransaction_InvalidTransaction(t *testing.T) {
	mockAccountRepo := &mocks.MockAccountRepository{
		InsertFn: func(ctx context.Context, account *domain.Account) (*domain.Account, error) {
			return &domain.Account{ID: 1, DocumentNumber: "12345678900"}, nil
		},
	}
	mockTransactionRepo := &mocks.MockTransactionRepository{}
	service := services.NewTransactionService(mockAccountRepo, mockTransactionRepo)

	_, err := service.CreateTransaction(context.Background(), 1, 999, 100.0)

	if err == nil {
		t.Errorf("expected error for invalid operation type")
	}
}

func TestTransactionService_CreateTransaction_TransactionRepoInsertError(t *testing.T) {
	mockAccountRepo := &mocks.MockAccountRepository{
		GetByIDFn: func(ctx context.Context, id int) (*domain.Account, error) {
			return &domain.Account{ID: 1, DocumentNumber: "12345678900"}, nil
		},
	}
	mockTransactionRepo := &mocks.MockTransactionRepository{
		InsertFn: func(ctx context.Context, transaction *domain.Transaction) (*domain.Transaction, error) {
			return nil, domain.ErrTransactionAccountNotExists
		},
	}
	service := services.NewTransactionService(mockAccountRepo, mockTransactionRepo)

	_, err := service.CreateTransaction(context.Background(), 1, int(domain.OperationTypePurchase), 100.0)

	if err != domain.ErrTransactionAccountNotExists {
		t.Errorf("expected error from repo insert, got %v", err)
	}
}
