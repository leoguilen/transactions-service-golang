package postgres

import (
	"context"
	"testing"

	"github.com/leoguilen/transactions/internal/adapters/db/postgres"
	"github.com/leoguilen/transactions/internal/core/domain"
	"github.com/leoguilen/transactions/internal/core/ports"
)

func newAccountRepo() ports.AccountRepository {
	return postgres.NewAccountRepository(ConnStr)
}

func TestAccountRepository_Insert_Success(t *testing.T) {
	repo := newAccountRepo()
	account := &domain.Account{DocumentNumber: "12345678901"}

	result, err := repo.Insert(context.Background(), account)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result == nil || result.ID == 0 {
		t.Fatal("expected account with ID, got nil or zero ID")
	}
	if result.DocumentNumber != "12345678901" {
		t.Errorf("expected DocumentNumber 12345678901, got %s", result.DocumentNumber)
	}
}

func TestAccountRepository_GetByID_Success(t *testing.T) {
	repo := newAccountRepo()
	account := &domain.Account{DocumentNumber: "98765432100"}
	inserted, err := repo.Insert(context.Background(), account)
	if err != nil {
		t.Fatalf("setup failed: could not insert account: %v", err)
	}

	result, err := repo.GetByID(context.Background(), inserted.ID)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result == nil || result.ID != inserted.ID {
		t.Errorf("expected account ID %d, got %v", inserted.ID, result)
	}
}

func TestAccountRepository_GetByID_NotFound(t *testing.T) {
	repo := newAccountRepo()
	result, err := repo.GetByID(context.Background(), 99999)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result != nil {
		t.Fatal("expected nil for nonexistent account")
	}
}
