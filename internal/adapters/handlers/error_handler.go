package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/leoguilen/transactions/internal/core/domain"
)

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    string `json:"code,omitempty"`
}

type HTTPErrorCode struct {
	StatusCode int
	Message    string
	Code       string
}

var errorMapping = map[error]HTTPErrorCode{
	domain.ErrInvalidAccount: {
		StatusCode: http.StatusBadRequest,
		Message:    "Invalid account",
		Code:       "INVALID_ACCOUNT",
	},
	domain.ErrInvalidAccountID: {
		StatusCode: http.StatusBadRequest,
		Message:    "Invalid account ID",
		Code:       "INVALID_ACCOUNT_ID",
	},
	domain.ErrAccountNotFound: {
		StatusCode: http.StatusNotFound,
		Message:    "Account not found",
		Code:       "ACCOUNT_NOT_FOUND",
	},
	domain.ErrAccountAlreadyExists: {
		StatusCode: http.StatusConflict,
		Message:    "Account already exists",
		Code:       "ACCOUNT_ALREADY_EXISTS",
	},
	domain.ErrInvalidAccountDocumentNumber: {
		StatusCode: http.StatusBadRequest,
		Message:    "Invalid document number",
		Code:       "INVALID_DOCUMENT_NUMBER",
	},
	domain.ErrInvalidTransaction: {
		StatusCode: http.StatusBadRequest,
		Message:    "Invalid transaction",
		Code:       "INVALID_TRANSACTION",
	},
	domain.ErrTransactionAccountNotExists: {
		StatusCode: http.StatusNotFound,
		Message:    "Account not found",
		Code:       "ACCOUNT_NOT_FOUND",
	},
	domain.ErrTransactionAccountInvalid: {
		StatusCode: http.StatusBadRequest,
		Message:    "Invalid account",
		Code:       "INVALID_ACCOUNT",
	},
	domain.ErrTransactionOperationTypeInvalid: {
		StatusCode: http.StatusBadRequest,
		Message:    "Invalid operation type",
		Code:       "INVALID_OPERATION_TYPE",
	},
	domain.ErrTransactionAmountInvalid: {
		StatusCode: http.StatusBadRequest,
		Message:    "Invalid amount",
		Code:       "INVALID_AMOUNT",
	},
}

func RespondWithError(w http.ResponseWriter, err error, defaultStatusCode int, defaultMessage string) {
	statusCode := defaultStatusCode
	message := defaultMessage
	code := ""

	if httpErr, exists := errorMapping[err]; exists {
		statusCode = httpErr.StatusCode
		message = httpErr.Message
		code = httpErr.Code
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := ErrorResponse{
		Error:   http.StatusText(statusCode),
		Message: message,
		Code:    code,
	}

	json.NewEncoder(w).Encode(response)
}

func RespondWithErrorf(w http.ResponseWriter, statusCode int, errorCode string, format string, args ...interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := ErrorResponse{
		Error:   http.StatusText(statusCode),
		Message: format,
		Code:    errorCode,
	}

	json.NewEncoder(w).Encode(response)
}

func RespondWithBadRequest(w http.ResponseWriter, message string) {
	RespondWithErrorf(w, http.StatusBadRequest, "BAD_REQUEST", message)
}

func RespondWithNotFound(w http.ResponseWriter, message string) {
	RespondWithErrorf(w, http.StatusNotFound, "NOT_FOUND", message)
}

func RespondWithConflict(w http.ResponseWriter, message string) {
	RespondWithErrorf(w, http.StatusConflict, "CONFLICT", message)
}

func RespondWithInternalServerError(w http.ResponseWriter, message string) {
	RespondWithErrorf(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", message)
}
