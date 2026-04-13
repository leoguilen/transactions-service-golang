package postgres

import (
	"context"
	"database/sql"

	"github.com/leoguilen/transactions/internal/core/domain"
	"github.com/leoguilen/transactions/internal/core/ports"

	_ "github.com/lib/pq"
)

const (
	InsertTransactionQuery = `INSERT INTO Transactions (AccountID, OperationTypeID, Amount, EventDate) VALUES ($1, $2, $3, $4) RETURNING Id, AccountID, OperationTypeID, Amount, EventDate`
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(connStr string) ports.TransactionRepository {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	return &TransactionRepository{db: db}
}

func (t *TransactionRepository) Insert(ctx context.Context, transaction *domain.Transaction) (*domain.Transaction, error) {
	conn, err := t.db.Conn(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	row := conn.QueryRowContext(ctx, InsertTransactionQuery, transaction.AccountID, transaction.OperationTypeID, transaction.Amount, transaction.EventDate)

	var newTransaction domain.Transaction
	if err := row.Scan(&newTransaction.ID, &newTransaction.AccountID, &newTransaction.OperationTypeID, &newTransaction.Amount, &newTransaction.EventDate); err != nil {
		return nil, err
	}

	return &newTransaction, nil
}
