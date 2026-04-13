package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/leoguilen/transactions/internal/adapters/handlers"
)

func TestCreateAccount_Success(t *testing.T) {
	handler := NewHttpHandler()

	req := httptest.NewRequest("POST", "/accounts",
		strings.NewReader(`{"document_number":"98765432101"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreateAccount(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, w.Code)
	}

	var resp handlers.AccountResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.DocumentNumber != "98765432101" {
		t.Errorf("unexpected document number: %s", resp.DocumentNumber)
	}
}

func TestCreateAccount_InvalidJSON(t *testing.T) {
	handler := NewHttpHandler()

	req := httptest.NewRequest("POST", "/accounts",
		strings.NewReader(`{invalid json}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreateAccount(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestCreateAccount_MissingDocumentNumber(t *testing.T) {
	handler := NewHttpHandler()

	req := httptest.NewRequest("POST", "/accounts",
		strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreateAccount(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestCreateAccount_InvalidDocumentNumber(t *testing.T) {
	handler := NewHttpHandler()

	req := httptest.NewRequest("POST", "/accounts",
		strings.NewReader(`{"document_number":"123"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreateAccount(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestCreateAccount_DuplicateDocumentNumber(t *testing.T) {
	handler := NewHttpHandler()

	req := httptest.NewRequest("POST", "/accounts",
		strings.NewReader(`{"document_number":"11122233344"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreateAccount(w, req)

	req2 := httptest.NewRequest("POST", "/accounts",
		strings.NewReader(`{"document_number":"11122233344"}`))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()

	handler.CreateAccount(w2, req2)

	if w2.Code != http.StatusConflict {
		t.Errorf("expected status %d, got %d", http.StatusConflict, w2.Code)
	}
}

func TestCreateAccount_EmptyDocumentNumber(t *testing.T) {
	handler := NewHttpHandler()

	req := httptest.NewRequest("POST", "/accounts",
		strings.NewReader(`{"document_number":""}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreateAccount(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}
