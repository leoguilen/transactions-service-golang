package repository

import "github.com/leoguilen/transactions/internal/domain"

// AccountRepository defines persistence operations for accounts (port)
type AccountRepository interface {
	Create(acc *domain.Account) (*domain.Account, error)
	GetByID(id string) (*domain.Account, error)
}
