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

func TestGetAccountByID_Success(t *testing.T) {
	handler := NewHttpHandler()

	createReq := httptest.NewRequest("POST", "/accounts",
		strings.NewReader(`{"document_number":"12345678901"}`))
	createReq.Header.Set("Content-Type", "application/json")
	createW := httptest.NewRecorder()
	handler.CreateAccount(createW, createReq)

	if createW.Code != http.StatusCreated {
		t.Fatalf("failed to create account, status: %d, body: %s", createW.Code, createW.Body.String())
	}

	var created handlers.AccountResponse
	if err := json.NewDecoder(createW.Body).Decode(&created); err != nil {
		t.Fatalf("failed to decode create response: %v", err)
	}

	accountId := strconv.FormatInt(int64(created.ID), 10)
	req := httptest.NewRequest("GET", "/accounts/"+accountId, nil)
	req.SetPathValue("id", accountId)
	w := httptest.NewRecorder()
	handler.GetAccountByID(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var resp handlers.AccountResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.ID != created.ID || resp.DocumentNumber != "12345678901" {
		t.Errorf("unexpected response data: %+v", resp)
	}
}

func TestGetAccountByID_NotFound(t *testing.T) {
	handler := NewHttpHandler()

	req := httptest.NewRequest("GET", "/accounts/999", nil)
	req.SetPathValue("id", "999")
	w := httptest.NewRecorder()

	handler.GetAccountByID(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestGetAccountByID_InvalidID(t *testing.T) {
	handler := NewHttpHandler()

	req := httptest.NewRequest("GET", "/accounts/invalid", nil)
	req.SetPathValue("id", "invalid")
	w := httptest.NewRecorder()

	handler.GetAccountByID(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}
