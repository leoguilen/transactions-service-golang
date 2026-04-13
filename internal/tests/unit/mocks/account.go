package mocks

import (
	"context"

	"github.com/leoguilen/transactions/internal/core/domain"
)

type MockAccountService struct {
	GetAccountByIDFn func(ctx context.Context, id int) (*domain.Account, error)
	CreateAccountFn  func(ctx context.Context, docNumber string) (*domain.Account, error)
}

func (m *MockAccountService) GetAccountByID(ctx context.Context, id int) (*domain.Account, error) {
	return m.GetAccountByIDFn(ctx, id)
}

func (m *MockAccountService) CreateAccount(ctx context.Context, docNumber string) (*domain.Account, error) {
	return m.CreateAccountFn(ctx, docNumber)
}

type MockAccountRepository struct {
	InsertFn  func(ctx context.Context, account *domain.Account) (*domain.Account, error)
	GetByIDFn func(ctx context.Context, id int) (*domain.Account, error)
}

func (m *MockAccountRepository) Insert(ctx context.Context, account *domain.Account) (*domain.Account, error) {
	if m.InsertFn != nil {
		return m.InsertFn(ctx, account)
	}
	return account, nil
}

func (m *MockAccountRepository) GetByID(ctx context.Context, id int) (*domain.Account, error) {
	if m.GetByIDFn != nil {
		return m.GetByIDFn(ctx, id)
	}
	return nil, nil
}
