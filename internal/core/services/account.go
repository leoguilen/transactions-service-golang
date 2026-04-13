package services

import (
	"context"
	"strings"

	"github.com/leoguilen/transactions/internal/core/domain"
	"github.com/leoguilen/transactions/internal/core/ports"
)

type AccountService struct {
	accountRepo ports.AccountRepository
}

func NewAccountService(accountRepo ports.AccountRepository) ports.AccountService {
	return &AccountService{accountRepo: accountRepo}
}

func (a *AccountService) CreateAccount(ctx context.Context, docNumber string) (*domain.Account, error) {
	account, err := domain.NewAccount(docNumber)
	if err != nil {
		return nil, err
	}

	insertedAccount, err := a.accountRepo.Insert(ctx, account)
	if err != nil {
		if strings.Contains(err.Error(), "unique constraint") {
			return nil, domain.ErrAccountAlreadyExists
		}
		return nil, err
	}

	return insertedAccount, nil
}

func (a *AccountService) GetAccountByID(ctx context.Context, id int) (*domain.Account, error) {
	if id <= 0 {
		return nil, domain.ErrInvalidAccountID
	}

	account, err := a.accountRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if account == nil {
		return nil, domain.ErrAccountNotFound
	}

	return account, nil
}
