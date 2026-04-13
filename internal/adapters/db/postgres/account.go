package postgres

import (
	"context"
	"database/sql"

	"github.com/leoguilen/transactions/internal/core/domain"
	"github.com/leoguilen/transactions/internal/core/ports"

	_ "github.com/lib/pq"
)

const (
	SelectAccountByIDQuery = `SELECT Id, DocumentNumber, CreatedAt FROM Accounts WHERE Id = $1`
	InsertAccountQuery     = `INSERT INTO Accounts (DocumentNumber) VALUES ($1) RETURNING Id, DocumentNumber, CreatedAt`
)

type AccountRepository struct {
	db *sql.DB
}

func NewAccountRepository(connStr string) ports.AccountRepository {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	return &AccountRepository{db: db}
}

func (a *AccountRepository) GetByID(ctx context.Context, id int) (*domain.Account, error) {
	conn, err := a.db.Conn(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	row := conn.QueryRowContext(ctx, SelectAccountByIDQuery, id)

	var account domain.Account
	if err := row.Scan(&account.ID, &account.DocumentNumber, &account.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &account, nil
}

func (a *AccountRepository) Insert(ctx context.Context, account *domain.Account) (*domain.Account, error) {
	conn, err := a.db.Conn(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	row := conn.QueryRowContext(ctx, InsertAccountQuery, account.DocumentNumber)

	var newAccount domain.Account
	if err := row.Scan(&newAccount.ID, &newAccount.DocumentNumber, &newAccount.CreatedAt); err != nil {
		return nil, err
	}

	return &newAccount, nil
}
