package ports

import (
	"context"

	"github.com/leoguilen/transactions/internal/core/domain"
)

type AccountService interface {
	CreateAccount(ctx context.Context, docNumber string) (*domain.Account, error)
	GetAccountByID(ctx context.Context, id int) (*domain.Account, error)
}

type AccountRepository interface {
	Insert(ctx context.Context, account *domain.Account) (*domain.Account, error)
	GetByID(ctx context.Context, id int) (*domain.Account, error)
}
