package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/leoguilen/transactions/internal/adapters/handlers"
	"github.com/leoguilen/transactions/internal/core/domain"
	"github.com/leoguilen/transactions/internal/tests/unit/mocks"
)

func TestGetAccountByID_Success(t *testing.T) {
	mockAcctService := &mocks.MockAccountService{
		GetAccountByIDFn: func(ctx context.Context, id int) (*domain.Account, error) {
			return &domain.Account{
				ID:             1,
				DocumentNumber: "12345678901",
				CreatedAt:      time.Now(),
			}, nil
		},
	}
	handler := handlers.NewHttpHandler(mockAcctService, nil)

	req := httptest.NewRequest("GET", "/accounts/1", nil)
	w := httptest.NewRecorder()

	handler.GetAccountByID(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var resp handlers.AccountResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.ID != 1 || resp.DocumentNumber != "12345678901" {
		t.Errorf("unexpected response data: %+v", resp)
	}
}

func TestGetAccountByID_NotFound(t *testing.T) {
	mockAcctService := &mocks.MockAccountService{
		GetAccountByIDFn: func(ctx context.Context, id int) (*domain.Account, error) {
			return nil, domain.ErrAccountNotFound
		},
	}
	handler := handlers.NewHttpHandler(mockAcctService, nil)

	req := httptest.NewRequest("GET", "/accounts/999", nil)
	w := httptest.NewRecorder()

	handler.GetAccountByID(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestGetAccountByID_InvalidID(t *testing.T) {
	mockAcctService := &mocks.MockAccountService{
		GetAccountByIDFn: func(ctx context.Context, id int) (*domain.Account, error) {
			return nil, domain.ErrInvalidAccountID
		},
	}
	handler := handlers.NewHttpHandler(mockAcctService, nil)

	req := httptest.NewRequest("GET", "/accounts/invalid", nil)
	w := httptest.NewRecorder()

	handler.GetAccountByID(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestCreateAccount_Success(t *testing.T) {
	mockAcctService := &mocks.MockAccountService{
		CreateAccountFn: func(ctx context.Context, docNumber string) (*domain.Account, error) {
			return &domain.Account{
				ID:             1,
				DocumentNumber: docNumber,
				CreatedAt:      time.Now(),
			}, nil
		},
	}
	handler := handlers.NewHttpHandler(mockAcctService, nil)

	body := handlers.CreateAccountRequest{DocumentNumber: "12345678901"}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/accounts", bytes.NewReader(bodyBytes))
	w := httptest.NewRecorder()

	handler.CreateAccount(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, w.Code)
	}

	if w.Header().Get("Location") == "" {
		t.Error("expected Location header to be set")
	}

	var resp handlers.AccountResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.DocumentNumber != "12345678901" {
		t.Errorf("expected document number %s, got %s", "12345678901", resp.DocumentNumber)
	}
}

func TestCreateAccount_InvalidBody(t *testing.T) {
	mockAcctService := &mocks.MockAccountService{}
	handler := handlers.NewHttpHandler(mockAcctService, nil)

	req := httptest.NewRequest("POST", "/accounts", bytes.NewReader([]byte("invalid json")))
	w := httptest.NewRecorder()

	handler.CreateAccount(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestCreateAccount_InvalidDocumentNumber(t *testing.T) {
	mockAcctService := &mocks.MockAccountService{
		CreateAccountFn: func(ctx context.Context, docNumber string) (*domain.Account, error) {
			return nil, domain.ErrInvalidAccountDocumentNumber
		},
	}
	handler := handlers.NewHttpHandler(mockAcctService, nil)

	body := handlers.CreateAccountRequest{DocumentNumber: "invalid"}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/accounts", bytes.NewReader(bodyBytes))
	w := httptest.NewRecorder()

	handler.CreateAccount(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestCreateAccount_AlreadyExists(t *testing.T) {
	mockAcctService := &mocks.MockAccountService{
		CreateAccountFn: func(ctx context.Context, docNumber string) (*domain.Account, error) {
			return nil, domain.ErrAccountAlreadyExists
		},
	}
	handler := handlers.NewHttpHandler(mockAcctService, nil)

	body := handlers.CreateAccountRequest{DocumentNumber: "12345678901"}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/accounts", bytes.NewReader(bodyBytes))
	w := httptest.NewRecorder()

	handler.CreateAccount(w, req)

	if w.Code != http.StatusConflict {
		t.Errorf("expected status %d, got %d", http.StatusConflict, w.Code)
	}
}

func TestCreateTransaction_Success(t *testing.T) {
	mockTxnService := &mocks.MockTransactionService{
		CreateTransactionFn: func(ctx context.Context, accountID, operationTypeID int, amount float64) (*domain.Transaction, error) {
			return &domain.Transaction{
				ID:              1,
				AccountID:       accountID,
				OperationTypeID: domain.OperationType(operationTypeID),
				Amount:          amount,
				EventDate:       time.Now(),
			}, nil
		},
	}
	handler := handlers.NewHttpHandler(nil, mockTxnService)

	body := handlers.CreateTransactionRequest{AccountID: 1, OperationTypeID: 1, Amount: 100.00}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/transactions", bytes.NewReader(bodyBytes))
	w := httptest.NewRecorder()

	handler.CreateTransaction(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, w.Code)
	}

	var resp handlers.TransactionResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Amount != 100.00 {
		t.Errorf("expected amount 100.00, got %f", resp.Amount)
	}
}

func TestCreateTransaction_InvalidBody(t *testing.T) {
	mockTxnService := &mocks.MockTransactionService{}
	handler := handlers.NewHttpHandler(nil, mockTxnService)

	req := httptest.NewRequest("POST", "/transactions", bytes.NewReader([]byte("invalid")))
	w := httptest.NewRecorder()

	handler.CreateTransaction(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestCreateTransaction_InvalidAmount(t *testing.T) {
	mockTxnService := &mocks.MockTransactionService{
		CreateTransactionFn: func(ctx context.Context, accountID, operationTypeID int, amount float64) (*domain.Transaction, error) {
			return nil, domain.ErrTransactionAmountInvalid
		},
	}
	handler := handlers.NewHttpHandler(nil, mockTxnService)

	body := handlers.CreateTransactionRequest{AccountID: 1, OperationTypeID: 1, Amount: -50.00}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/transactions", bytes.NewReader(bodyBytes))
	w := httptest.NewRecorder()

	handler.CreateTransaction(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestCreateTransaction_AccountNotExists(t *testing.T) {
	mockTxnService := &mocks.MockTransactionService{
		CreateTransactionFn: func(ctx context.Context, accountID, operationTypeID int, amount float64) (*domain.Transaction, error) {
			return nil, domain.ErrTransactionAccountNotExists
		},
	}
	handler := handlers.NewHttpHandler(nil, mockTxnService)

	body := handlers.CreateTransactionRequest{AccountID: 999, OperationTypeID: 1, Amount: 100.00}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/transactions", bytes.NewReader(bodyBytes))
	w := httptest.NewRecorder()

	handler.CreateTransaction(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestRegisterRoutes(t *testing.T) {
	mux := http.NewServeMux()
	handler := handlers.NewHttpHandler(nil, nil)

	handler.RegisterRoutes(mux)

	if mux == nil {
		t.Fatal("mux should not be nil after registering routes")
	}
}
