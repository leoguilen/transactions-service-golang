package services

import (
	"context"
	"errors"
	"testing"

	"github.com/leoguilen/transactions/internal/core/domain"
	"github.com/leoguilen/transactions/internal/core/services"
	"github.com/leoguilen/transactions/internal/tests/unit/mocks"
)

func TestAccountService_CreateAccount_Success(t *testing.T) {
	docNumber := "12345678901"
	expectedAccount := &domain.Account{ID: 1, DocumentNumber: docNumber}

	mock := &mocks.MockAccountRepository{
		InsertFn: func(ctx context.Context, account *domain.Account) (*domain.Account, error) {
			return expectedAccount, nil
		},
	}

	service := services.NewAccountService(mock)
	account, err := service.CreateAccount(context.Background(), docNumber)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if account == nil {
		t.Fatal("expected account, got nil")
	}
	if account.ID != expectedAccount.ID {
		t.Errorf("expected ID %d, got %d", expectedAccount.ID, account.ID)
	}
	if account.DocumentNumber != docNumber {
		t.Errorf("expected DocumentNumber %s, got %s", docNumber, account.DocumentNumber)
	}
}

func TestAccountService_CreateAccount_InvalidDocumentNumber(t *testing.T) {
	mock := &mocks.MockAccountRepository{}
	service := services.NewAccountService(mock)

	account, err := service.CreateAccount(context.Background(), "123")

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if account != nil {
		t.Fatal("expected nil account, got value")
	}
	if err != domain.ErrInvalidAccountDocumentNumber {
		t.Errorf("expected ErrInvalidAccountDocumentNumber, got %v", err)
	}
}

func TestAccountService_CreateAccount_UniqueConstraintError(t *testing.T) {
	docNumber := "12345678901"

	mock := &mocks.MockAccountRepository{
		InsertFn: func(ctx context.Context, account *domain.Account) (*domain.Account, error) {
			return nil, errors.New("unique constraint violation")
		},
	}

	service := services.NewAccountService(mock)
	account, err := service.CreateAccount(context.Background(), docNumber)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if account != nil {
		t.Fatal("expected nil account, got value")
	}
	if err != domain.ErrAccountAlreadyExists {
		t.Errorf("expected ErrAccountAlreadyExists, got %v", err)
	}
}

func TestAccountService_CreateAccount_RepositoryError(t *testing.T) {
	docNumber := "12345678901"
	expectedErr := errors.New("database error")

	mock := &mocks.MockAccountRepository{
		InsertFn: func(ctx context.Context, account *domain.Account) (*domain.Account, error) {
			return nil, expectedErr
		},
	}

	service := services.NewAccountService(mock)
	account, err := service.CreateAccount(context.Background(), docNumber)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if account != nil {
		t.Fatal("expected nil account, got value")
	}
	if err != expectedErr {
		t.Errorf("expected %v, got %v", expectedErr, err)
	}
}

func TestAccountService_GetAccountByID_Success(t *testing.T) {
	accountID := 1
	expectedAccount := &domain.Account{ID: accountID, DocumentNumber: "12345678901"}

	mock := &mocks.MockAccountRepository{
		GetByIDFn: func(ctx context.Context, id int) (*domain.Account, error) {
			if id == accountID {
				return expectedAccount, nil
			}
			return nil, nil
		},
	}

	service := services.NewAccountService(mock)
	account, err := service.GetAccountByID(context.Background(), accountID)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if account == nil {
		t.Fatal("expected account, got nil")
	}
	if account.ID != accountID {
		t.Errorf("expected ID %d, got %d", accountID, account.ID)
	}
}

func TestAccountService_GetAccountByID_InvalidID_Zero(t *testing.T) {
	mock := &mocks.MockAccountRepository{}
	service := services.NewAccountService(mock)

	account, err := service.GetAccountByID(context.Background(), 0)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if account != nil {
		t.Fatal("expected nil account, got value")
	}
	if err != domain.ErrInvalidAccountID {
		t.Errorf("expected ErrInvalidAccountID, got %v", err)
	}
}

func TestAccountService_GetAccountByID_InvalidID_Negative(t *testing.T) {
	mock := &mocks.MockAccountRepository{}
	service := services.NewAccountService(mock)

	account, err := service.GetAccountByID(context.Background(), -1)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if account != nil {
		t.Fatal("expected nil account, got value")
	}
	if err != domain.ErrInvalidAccountID {
		t.Errorf("expected ErrInvalidAccountID, got %v", err)
	}
}

func TestAccountService_GetAccountByID_NotFound(t *testing.T) {
	mock := &mocks.MockAccountRepository{
		GetByIDFn: func(ctx context.Context, id int) (*domain.Account, error) {
			return nil, nil
		},
	}

	service := services.NewAccountService(mock)
	account, err := service.GetAccountByID(context.Background(), 999)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if account != nil {
		t.Fatal("expected nil account, got value")
	}
	if err != domain.ErrAccountNotFound {
		t.Errorf("expected ErrAccountNotFound, got %v", err)
	}
}

func TestAccountService_GetAccountByID_RepositoryError(t *testing.T) {
	expectedErr := errors.New("database error")

	mock := &mocks.MockAccountRepository{
		GetByIDFn: func(ctx context.Context, id int) (*domain.Account, error) {
			return nil, expectedErr
		},
	}

	service := services.NewAccountService(mock)
	account, err := service.GetAccountByID(context.Background(), 1)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if account != nil {
		t.Fatal("expected nil account, got value")
	}
	if err != expectedErr {
		t.Errorf("expected %v, got %v", expectedErr, err)
	}
}
