package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/leoguilen/transactions/internal/adapters/handlers"
	"github.com/leoguilen/transactions/internal/core/domain"
)

func TestRespondWithError(t *testing.T) {
	tests := []struct {
		name               string
		err                error
		defaultStatusCode  int
		expectedStatusCode int
		expectedMessage    string
		expectedCode       string
	}{
		{
			name:               "account not found error",
			err:                domain.ErrAccountNotFound,
			defaultStatusCode:  http.StatusInternalServerError,
			expectedStatusCode: http.StatusNotFound,
			expectedMessage:    "Account not found",
			expectedCode:       "ACCOUNT_NOT_FOUND",
		},
		{
			name:               "invalid account ID error",
			err:                domain.ErrInvalidAccountID,
			defaultStatusCode:  http.StatusInternalServerError,
			expectedStatusCode: http.StatusBadRequest,
			expectedMessage:    "Invalid account ID",
			expectedCode:       "INVALID_ACCOUNT_ID",
		},
		{
			name:               "account already exists error",
			err:                domain.ErrAccountAlreadyExists,
			defaultStatusCode:  http.StatusInternalServerError,
			expectedStatusCode: http.StatusConflict,
			expectedMessage:    "Account already exists",
			expectedCode:       "ACCOUNT_ALREADY_EXISTS",
		},
		{
			name:               "invalid document number error",
			err:                domain.ErrInvalidAccountDocumentNumber,
			defaultStatusCode:  http.StatusInternalServerError,
			expectedStatusCode: http.StatusBadRequest,
			expectedMessage:    "Invalid document number",
			expectedCode:       "INVALID_DOCUMENT_NUMBER",
		},
		{
			name:               "invalid transaction error",
			err:                domain.ErrInvalidTransaction,
			defaultStatusCode:  http.StatusInternalServerError,
			expectedStatusCode: http.StatusBadRequest,
			expectedMessage:    "Invalid transaction",
			expectedCode:       "INVALID_TRANSACTION",
		},
		{
			name:               "transaction account not exists error",
			err:                domain.ErrTransactionAccountNotExists,
			defaultStatusCode:  http.StatusInternalServerError,
			expectedStatusCode: http.StatusNotFound,
			expectedMessage:    "Account not found",
			expectedCode:       "ACCOUNT_NOT_FOUND",
		},
		{
			name:               "invalid operation type error",
			err:                domain.ErrTransactionOperationTypeInvalid,
			defaultStatusCode:  http.StatusInternalServerError,
			expectedStatusCode: http.StatusBadRequest,
			expectedMessage:    "Invalid operation type",
			expectedCode:       "INVALID_OPERATION_TYPE",
		},
		{
			name:               "invalid amount error",
			err:                domain.ErrTransactionAmountInvalid,
			defaultStatusCode:  http.StatusInternalServerError,
			expectedStatusCode: http.StatusBadRequest,
			expectedMessage:    "Invalid amount",
			expectedCode:       "INVALID_AMOUNT",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			handlers.RespondWithError(w, tt.err, tt.defaultStatusCode, "Failed operation")

			if w.Code != tt.expectedStatusCode {
				t.Errorf("Expected status %d, got %d", tt.expectedStatusCode, w.Code)
			}

			var response handlers.ErrorResponse
			if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
				t.Fatalf("Failed to parse response: %v", err)
			}

			if response.Message != tt.expectedMessage {
				t.Errorf("Expected message %q, got %q", tt.expectedMessage, response.Message)
			}

			if response.Code != tt.expectedCode {
				t.Errorf("Expected code %q, got %q", tt.expectedCode, response.Code)
			}

			contentType := w.Header().Get("Content-Type")
			if contentType != "application/json" {
				t.Errorf("Expected Content-Type application/json, got %s", contentType)
			}
		})
	}
}

func TestRespondWithBadRequest(t *testing.T) {
	w := httptest.NewRecorder()
	handlers.RespondWithBadRequest(w, "Invalid input")

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	var response handlers.ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if response.Code != "BAD_REQUEST" {
		t.Errorf("Expected code BAD_REQUEST, got %s", response.Code)
	}
}

func TestRespondWithNotFound(t *testing.T) {
	w := httptest.NewRecorder()
	handlers.RespondWithNotFound(w, "Resource not found")

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}

	var response handlers.ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if response.Code != "NOT_FOUND" {
		t.Errorf("Expected code NOT_FOUND, got %s", response.Code)
	}
}

func TestRespondWithConflict(t *testing.T) {
	w := httptest.NewRecorder()
	handlers.RespondWithConflict(w, "Resource already exists")

	if w.Code != http.StatusConflict {
		t.Errorf("Expected status %d, got %d", http.StatusConflict, w.Code)
	}

	var response handlers.ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if response.Code != "CONFLICT" {
		t.Errorf("Expected code CONFLICT, got %s", response.Code)
	}
}

func TestRespondWithInternalServerError(t *testing.T) {
	w := httptest.NewRecorder()
	handlers.RespondWithInternalServerError(w, "Internal server error")

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}

	var response handlers.ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if response.Code != "INTERNAL_SERVER_ERROR" {
		t.Errorf("Expected code INTERNAL_SERVER_ERROR, got %s", response.Code)
	}
}

func TestErrorResponseFormat(t *testing.T) {
	w := httptest.NewRecorder()
	handlers.RespondWithError(w, domain.ErrAccountNotFound, http.StatusInternalServerError, "Default message")

	var response handlers.ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if response.Error == "" {
		t.Error("Error field should not be empty")
	}

	if response.Message == "" {
		t.Error("Message field should not be empty")
	}

	if response.Code == "" {
		t.Error("Code field should not be empty")
	}
}

func TestUnmappedErrorUsesDefault(t *testing.T) {
	customErr := domain.ErrInvalidAccount
	w := httptest.NewRecorder()

	handlers.RespondWithError(w, customErr, http.StatusTeapot, "Custom default")

	if w.Code != http.StatusBadRequest {
		t.Errorf("Mapped error should use mapped status, got %d", w.Code)
	}
}
