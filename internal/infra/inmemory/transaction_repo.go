package inmemory

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/leoguilen/transactions/internal/domain"
)

// TransactionRepository is an in-memory implementation for tests
type TransactionRepository struct {
	mu           sync.RWMutex
	transactions map[string]*domain.Transaction
	byAccount    map[string][]*domain.Transaction
}

func NewTransactionRepository() *TransactionRepository {
	return &TransactionRepository{
		transactions: make(map[string]*domain.Transaction),
		byAccount:    make(map[string][]*domain.Transaction),
	}
}

func (r *TransactionRepository) Create(tx *domain.Transaction) (*domain.Transaction, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if tx.TransactionID == "" {
		tx.TransactionID = uuid.New().String()
	}
	if tx.CreatedAt.IsZero() {
		tx.CreatedAt = time.Now().UTC()
	}
	r.transactions[tx.TransactionID] = tx
	r.byAccount[tx.AccountID] = append([]*domain.Transaction{tx}, r.byAccount[tx.AccountID]...)
	return tx, nil
}

func (r *TransactionRepository) ListByAccount(accountID string, limit, offset int) ([]*domain.Transaction, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	list := r.byAccount[accountID]
	if offset >= len(list) {
		return []*domain.Transaction{}, nil
	}
	end := offset + limit
	if end > len(list) {
		end = len(list)
	}
	return list[offset:end], nil
}

func (r *TransactionRepository) SumByAccount(accountID string) (string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	list := r.byAccount[accountID]
	var sum float64
	for _, t := range list {
		// Amount stored as string for numeric precision in domain; simple float conversion for in-memory
		var v float64
		fmt.Sscanf(t.Amount, "%f", &v)
		sum += v
	}
	return fmt.Sprintf("%.2f", sum), nil
}
