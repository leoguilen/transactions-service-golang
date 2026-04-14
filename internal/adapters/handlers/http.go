package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/leoguilen/transactions/internal/core/ports"
)

type HttpHandler struct {
	accountService     ports.AccountService
	transactionService ports.TransactionService
}

type AccountResponse struct {
	ID             int    `json:"id"`
	DocumentNumber string `json:"document_number"`
	CreatedAt      string `json:"created_at"`
}

type TransactionResponse struct {
	ID              int     `json:"id"`
	AccountID       int     `json:"account_id"`
	OperationTypeID int     `json:"operation_type_id"`
	Amount          float64 `json:"amount"`
	EventDate       string  `json:"event_date"`
}

type CreateAccountRequest struct {
	DocumentNumber string `json:"document_number"`
}

type CreateTransactionRequest struct {
	AccountID       int     `json:"account_id"`
	OperationTypeID int     `json:"operation_type_id"`
	Amount          float64 `json:"amount"`
}

func NewHttpHandler(accountService ports.AccountService, transactionService ports.TransactionService) *HttpHandler {
	return &HttpHandler{
		accountService:     accountService,
		transactionService: transactionService,
	}
}

func (h *HttpHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /accounts/{id}", h.GetAccountByID)
	mux.HandleFunc("POST /accounts", h.CreateAccount)
	mux.HandleFunc("POST /transactions", h.CreateTransaction)

	// Swagger/OpenAPI routes
	mux.HandleFunc("GET /swagger/index.html", ServeSwaggerUI)
	mux.HandleFunc("GET /swagger/", ServeSwaggerUI)
	mux.HandleFunc("GET /swagger.json", ServeSwaggerUI)
}

// GetAccountByID retrieves account information by ID
// @Summary Get Account by ID
// @Description Retrieve account details by account ID
// @Tags accounts
// @Accept json
// @Produce json
// @Param id path int true "Account ID"
// @Success 200 {object} AccountResponse "Account details"
// @Failure 400 {object} ErrorResponse "Invalid account ID"
// @Failure 404 {object} ErrorResponse "Account not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /accounts/{id} [get]
func (h *HttpHandler) GetAccountByID(w http.ResponseWriter, r *http.Request) {
	accountID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		RespondWithBadRequest(w, "Invalid account ID format")
		return
	}

	account, err := h.accountService.GetAccountByID(r.Context(), accountID)

	if err != nil {
		RespondWithError(w, err, http.StatusInternalServerError, "Failed to retrieve account")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&AccountResponse{
		ID:             account.ID,
		DocumentNumber: account.DocumentNumber,
		CreatedAt:      account.CreatedAt.Format(time.RFC3339),
	})
}

// CreateAccount creates a new account
// @Summary Create Account
// @Description Create a new account with a document number
// @Tags accounts
// @Accept json
// @Produce json
// @Param request body CreateAccountRequest true "Account creation request"
// @Success 201 {object} AccountResponse "Account created"
// @Failure 400 {object} ErrorResponse "Invalid document number"
// @Failure 409 {object} ErrorResponse "Account already exists"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /accounts [post]
func (h *HttpHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	var req CreateAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondWithBadRequest(w, "Invalid request body")
		return
	}

	account, err := h.accountService.CreateAccount(r.Context(), req.DocumentNumber)

	if err != nil {
		RespondWithError(w, err, http.StatusInternalServerError, "Failed to create account")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Location", "/accounts/"+strconv.Itoa(account.ID))
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(&AccountResponse{
		ID:             account.ID,
		DocumentNumber: account.DocumentNumber,
		CreatedAt:      account.CreatedAt.Format(time.RFC3339),
	})
}

// CreateTransaction creates a new transaction
// @Summary Create Transaction
// @Description Create a new transaction for an account
// @Tags transactions
// @Accept json
// @Produce json
// @Param request body CreateTransactionRequest true "Transaction creation request"
// @Success 201 {object} TransactionResponse "Transaction created"
// @Failure 400 {object} ErrorResponse "Invalid transaction data"
// @Failure 404 {object} ErrorResponse "Account not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /transactions [post]
func (h *HttpHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var req CreateTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondWithBadRequest(w, "Invalid request body")
		return
	}

	transaction, err := h.transactionService.CreateTransaction(r.Context(), req.AccountID, req.OperationTypeID, req.Amount)

	if err != nil {
		RespondWithError(w, err, http.StatusInternalServerError, "Failed to create transaction")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(&TransactionResponse{
		ID:              transaction.ID,
		AccountID:       transaction.AccountID,
		OperationTypeID: int(transaction.OperationTypeID),
		Amount:          transaction.Amount,
		EventDate:       transaction.EventDate.Format(time.RFC3339),
	})
}
