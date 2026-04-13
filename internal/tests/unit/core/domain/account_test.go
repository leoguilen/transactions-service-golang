package domain

import (
	"testing"

	"github.com/leoguilen/transactions/internal/core/domain"
)

func TestNewAccount_ValidDocumentNumber11Digits(t *testing.T) {
	documentNumber := "12345678901"
	account, err := domain.NewAccount(documentNumber)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if account == nil {
		t.Fatal("expected account, got nil")
	}
	if account.DocumentNumber != documentNumber {
		t.Errorf("expected DocumentNumber %s, got %s", documentNumber, account.DocumentNumber)
	}
	if account.ID != 0 {
		t.Errorf("expected ID 0, got %d", account.ID)
	}
}

func TestNewAccount_ValidDocumentNumber14Digits(t *testing.T) {
	documentNumber := "12345678901234"
	account, err := domain.NewAccount(documentNumber)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if account == nil {
		t.Fatal("expected account, got nil")
	}
	if account.DocumentNumber != documentNumber {
		t.Errorf("expected DocumentNumber %s, got %s", documentNumber, account.DocumentNumber)
	}
}

func TestNewAccount_InvalidDocumentNumberTooShort(t *testing.T) {
	documentNumber := "1234567890"
	account, err := domain.NewAccount(documentNumber)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if account != nil {
		t.Fatal("expected nil account, got value")
	}
	if err != domain.ErrInvalidAccountDocumentNumber {
		t.Errorf("expected ErrInvalidAccountDocumentNumber, got %v", err)
	}
}

func TestNewAccount_InvalidDocumentNumberTooLong(t *testing.T) {
	documentNumber := "123456789012345"
	account, err := domain.NewAccount(documentNumber)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if account != nil {
		t.Fatal("expected nil account, got value")
	}
	if err != domain.ErrInvalidAccountDocumentNumber {
		t.Errorf("expected ErrInvalidAccountDocumentNumber, got %v", err)
	}
}

func TestNewAccount_InvalidDocumentNumberEmpty(t *testing.T) {
	documentNumber := ""
	account, err := domain.NewAccount(documentNumber)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if account != nil {
		t.Fatal("expected nil account, got value")
	}
	if err != domain.ErrInvalidAccountDocumentNumber {
		t.Errorf("expected ErrInvalidAccountDocumentNumber, got %v", err)
	}
}
