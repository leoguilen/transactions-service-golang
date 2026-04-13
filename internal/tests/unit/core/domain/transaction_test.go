package domain

import (
	"testing"
	"time"

	"github.com/leoguilen/transactions/internal/core/domain"
)

func TestNewTransaction_ValidInputs(t *testing.T) {
	tests := []struct {
		name            string
		accountID       int
		operationTypeID int
		amount          float64
		wantAmount      float64
	}{
		{
			name:            "Purchase negative amount",
			accountID:       1,
			operationTypeID: int(domain.OperationTypePurchase),
			amount:          100.0,
			wantAmount:      -100.0,
		},
		{
			name:            "Withdrawal negative amount",
			accountID:       2,
			operationTypeID: int(domain.OperationTypeWithdrawal),
			amount:          50.0,
			wantAmount:      -50.0,
		},
		{
			name:            "PurchaseInstallment negative amount",
			accountID:       3,
			operationTypeID: int(domain.OperationTypePurchaseInstallment),
			amount:          200.0,
			wantAmount:      -200.0,
		},
		{
			name:            "CreditVoucher positive amount",
			accountID:       4,
			operationTypeID: int(domain.OperationTypeCreditVoucher),
			amount:          75.0,
			wantAmount:      75.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr, err := domain.NewTransaction(tt.accountID, tt.operationTypeID, tt.amount)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if tr.AccountID != tt.accountID {
				t.Errorf("expected AccountID %d, got %d", tt.accountID, tr.AccountID)
			}
			if tr.OperationTypeID != domain.OperationType(tt.operationTypeID) {
				t.Errorf("expected OperationTypeID %d, got %d", tt.operationTypeID, tr.OperationTypeID)
			}
			if tr.Amount != tt.wantAmount {
				t.Errorf("expected Amount %v, got %v", tt.wantAmount, tr.Amount)
			}
			if time.Since(tr.EventDate) > time.Second {
				t.Errorf("EventDate not set to now")
			}
		})
	}
}

func TestNewTransaction_InvalidOperationType(t *testing.T) {
	invalidTypes := []int{0, 5, -1, 100}
	for _, opType := range invalidTypes {
		t.Run("InvalidOperationType", func(t *testing.T) {
			_, err := domain.NewTransaction(1, opType, 10.0)
			if err == nil {
				t.Errorf("expected error for invalid operation type %d", opType)
			}
		})
	}
}

func TestNewTransaction_InvalidAmount(t *testing.T) {
	invalidAmounts := []float64{0, -10.0}
	for _, amt := range invalidAmounts {
		t.Run("InvalidAmount", func(t *testing.T) {
			_, err := domain.NewTransaction(1, int(domain.OperationTypePurchase), amt)
			if err == nil {
				t.Errorf("expected error for invalid amount %v", amt)
			}
		})
	}
}
