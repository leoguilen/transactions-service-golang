package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/leoguilen/transactions/internal/adapters/handlers"
)

func TestCreateTransaction_Success(t *testing.T) {
	handler := NewHttpHandler()

	createAccountReq := httptest.NewRequest("POST", "/accounts",
		strings.NewReader(`{"document_number":"00011122233"}`))
	createAccountReq.Header.Set("Content-Type", "application/json")
	createAccountW := httptest.NewRecorder()
	handler.CreateAccount(createAccountW, createAccountReq)

	if createAccountW.Code != http.StatusCreated {
		t.Fatalf("failed to create account, status: %d, body: %s", createAccountW.Code, createAccountW.Body.String())
	}

	var createdAccount handlers.AccountResponse
	if err := json.NewDecoder(createAccountW.Body).Decode(&createdAccount); err != nil {
		t.Fatalf("failed to decode create account response: %v", err)
	}

	createTransactionReq := httptest.NewRequest("POST", "/transactions",
		strings.NewReader(`{"account_id":`+strconv.Itoa(createdAccount.ID)+`, "operation_type_id":1, "amount":100.50}`))
	createTransactionReq.Header.Set("Content-Type", "application/json")
	createTransactionW := httptest.NewRecorder()
	handler.CreateTransaction(createTransactionW, createTransactionReq)

	if createTransactionW.Code != http.StatusCreated {
		t.Fatalf("failed to create transaction, status: %d, body: %s", createTransactionW.Code, createTransactionW.Body.String())
	}

	var createdTransaction handlers.TransactionResponse
	if err := json.NewDecoder(createTransactionW.Body).Decode(&createdTransaction); err != nil {
		t.Fatalf("failed to decode create transaction response: %v", err)
	}

	if createdTransaction.ID == 0 {
		t.Errorf("unexpected transaction data: %+v", createdTransaction)
	}
}

func TestCreateTransaction_InvalidOperationType(t *testing.T) {
	handler := NewHttpHandler()

	createAccountReq := httptest.NewRequest("POST", "/accounts",
		strings.NewReader(`{"document_number":"00011122244"}`))
	createAccountReq.Header.Set("Content-Type", "application/json")
	createAccountW := httptest.NewRecorder()
	handler.CreateAccount(createAccountW, createAccountReq)

	var createdAccount handlers.AccountResponse
	json.NewDecoder(createAccountW.Body).Decode(&createdAccount)

	createTransactionReq := httptest.NewRequest("POST", "/transactions",
		strings.NewReader(`{"account_id":`+strconv.Itoa(createdAccount.ID)+`, "operation_type_id":999, "amount":100.50}`))
	createTransactionReq.Header.Set("Content-Type", "application/json")
	createTransactionW := httptest.NewRecorder()
	handler.CreateTransaction(createTransactionW, createTransactionReq)

	if createTransactionW.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d, body: %s", http.StatusBadRequest, createTransactionW.Code, createTransactionW.Body.String())
	}
}

func TestCreateTransaction_InvalidAmount(t *testing.T) {
	handler := NewHttpHandler()

	createAccountReq := httptest.NewRequest("POST", "/accounts",
		strings.NewReader(`{"document_number":"00011122255"}`))
	createAccountReq.Header.Set("Content-Type", "application/json")
	createAccountW := httptest.NewRecorder()
	handler.CreateAccount(createAccountW, createAccountReq)

	var createdAccount handlers.AccountResponse
	json.NewDecoder(createAccountW.Body).Decode(&createdAccount)

	createTransactionReq := httptest.NewRequest("POST", "/transactions",
		strings.NewReader(`{"account_id":`+strconv.Itoa(createdAccount.ID)+`, "operation_type_id":1, "amount":-50.00}`))
	createTransactionReq.Header.Set("Content-Type", "application/json")
	createTransactionW := httptest.NewRecorder()
	handler.CreateTransaction(createTransactionW, createTransactionReq)

	if createTransactionW.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d, body: %s", http.StatusBadRequest, createTransactionW.Code, createTransactionW.Body.String())
	}
}

func TestCreateTransaction_MalformedJSON(t *testing.T) {
	handler := NewHttpHandler()

	createTransactionReq := httptest.NewRequest("POST", "/transactions",
		strings.NewReader(`{invalid json}`))
	createTransactionReq.Header.Set("Content-Type", "application/json")
	createTransactionW := httptest.NewRecorder()
	handler.CreateTransaction(createTransactionW, createTransactionReq)

	if createTransactionW.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d, body: %s", http.StatusBadRequest, createTransactionW.Code, createTransactionW.Body.String())
	}
}

func TestCreateTransaction_MissingRequiredFields(t *testing.T) {
	handler := NewHttpHandler()

	createTransactionReq := httptest.NewRequest("POST", "/transactions",
		strings.NewReader(`{"account_id":1}`))
	createTransactionReq.Header.Set("Content-Type", "application/json")
	createTransactionW := httptest.NewRecorder()
	handler.CreateTransaction(createTransactionW, createTransactionReq)

	if createTransactionW.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d, body: %s", http.StatusBadRequest, createTransactionW.Code, createTransactionW.Body.String())
	}
}

func TestCreateTransaction_ZeroAmount(t *testing.T) {
	handler := NewHttpHandler()

	createAccountReq := httptest.NewRequest("POST", "/accounts",
		strings.NewReader(`{"document_number":"00011122266"}`))
	createAccountReq.Header.Set("Content-Type", "application/json")
	createAccountW := httptest.NewRecorder()
	handler.CreateAccount(createAccountW, createAccountReq)

	var createdAccount handlers.AccountResponse
	json.NewDecoder(createAccountW.Body).Decode(&createdAccount)

	createTransactionReq := httptest.NewRequest("POST", "/transactions",
		strings.NewReader(`{"account_id":`+strconv.Itoa(createdAccount.ID)+`, "operation_type_id":1, "amount":0}`))
	createTransactionReq.Header.Set("Content-Type", "application/json")
	createTransactionW := httptest.NewRecorder()
	handler.CreateTransaction(createTransactionW, createTransactionReq)

	if createTransactionW.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d, body: %s", http.StatusBadRequest, createTransactionW.Code, createTransactionW.Body.String())
	}
}
