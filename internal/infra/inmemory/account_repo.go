package inmemory

import (
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/leoguilen/transactions/internal/domain"
)

var (
	ErrAccountNotFound = errors.New("account not found")
)

// AccountRepository is an in-memory implementation for tests
type AccountRepository struct {
	mu       sync.RWMutex
	accounts map[string]*domain.Account
}

func NewAccountRepository() *AccountRepository {
	return &AccountRepository{accounts: make(map[string]*domain.Account)}
}

func (r *AccountRepository) Create(acc *domain.Account) (*domain.Account, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if acc.AccountID == "" {
		acc.AccountID = uuid.New().String()
	}
	acc.CreatedAt = time.Now().UTC()
	r.accounts[acc.AccountID] = acc
	return acc, nil
}

func (r *AccountRepository) GetByID(id string) (*domain.Account, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	acc, ok := r.accounts[id]
	if !ok {
		return nil, ErrAccountNotFound
	}
	return acc, nil
}
