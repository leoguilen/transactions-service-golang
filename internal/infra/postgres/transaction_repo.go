package postgres

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/leoguilen/transactions/internal/domain"
)

// PostgresTransactionRepository implements transaction persistence using Postgres
type PostgresTransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *PostgresTransactionRepository {
	return &PostgresTransactionRepository{db: db}
}

func (r *PostgresTransactionRepository) Create(txn *domain.Transaction) (*domain.Transaction, error) {
	if txn.TransactionID == "" {
		// allow DB to generate UUID if desired; require caller to set ID externally if needed
	}
	if txn.CreatedAt.IsZero() {
		txn.CreatedAt = time.Now().UTC()
	}
	query := `INSERT INTO transactions (transaction_id, account_id, operation_type_id, amount, event_date, description, metadata, created_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`
	_, err := r.db.Exec(query, txn.TransactionID, txn.AccountID, txn.OperationTypeID, txn.Amount, txn.EventDate, nullableString(txn.Description), nullableString(txn.Metadata), txn.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("insert transaction: %w", err)
	}
	return txn, nil
}

func (r *PostgresTransactionRepository) ListByAccount(accountID string, limit, offset int) ([]*domain.Transaction, error) {
	query := `SELECT transaction_id, account_id, operation_type_id, amount, event_date, description, metadata, created_at FROM transactions WHERE account_id = $1 ORDER BY event_date DESC LIMIT $2 OFFSET $3`
	rows, err := r.db.Query(query, accountID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("query transactions: %w", err)
	}
	defer rows.Close()
	res := make([]*domain.Transaction, 0)
	for rows.Next() {
		t := &domain.Transaction{}
		if err := rows.Scan(&t.TransactionID, &t.AccountID, &t.OperationTypeID, &t.Amount, &t.EventDate, &t.Description, &t.Metadata, &t.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan transaction: %w", err)
		}
		res = append(res, t)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return res, nil
}

func (r *PostgresTransactionRepository) SumByAccount(accountID string) (string, error) {
	query := `SELECT COALESCE(SUM(amount)::text, '0') FROM transactions WHERE account_id = $1`
	var sum string
	if err := r.db.QueryRow(query, accountID).Scan(&sum); err != nil {
		return "", fmt.Errorf("sum by account: %w", err)
	}
	return sum, nil
}

func nullableString(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}
