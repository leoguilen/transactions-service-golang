package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/leoguilen/transactions/internal/core/domain"
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
}

func (h *HttpHandler) GetAccountByID(w http.ResponseWriter, r *http.Request) {
	accountID, err := strconv.Atoi(r.PathValue("id"))

	account, err := h.accountService.GetAccountByID(r.Context(), accountID)
	log.Printf("GetAccountByID: accountID=%d, err=%v", accountID, err)
	switch err {
	case nil:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(&AccountResponse{
			ID:             account.ID,
			DocumentNumber: account.DocumentNumber,
			CreatedAt:      account.CreatedAt.Format(time.RFC3339),
		})
	case domain.ErrAccountNotFound:
		http.Error(w, "account not found", http.StatusNotFound)
	case domain.ErrInvalidAccountID:
		http.Error(w, "invalid account ID", http.StatusBadRequest)
	default:
		http.Error(w, "internal server error: "+err.Error(), http.StatusInternalServerError)
	}
}

func (h *HttpHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	var req CreateAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	account, err := h.accountService.CreateAccount(r.Context(), req.DocumentNumber)

	switch err {
	case nil:
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Location", "/accounts/"+strconv.Itoa(account.ID))
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(&AccountResponse{
			ID:             account.ID,
			DocumentNumber: account.DocumentNumber,
			CreatedAt:      account.CreatedAt.Format(time.RFC3339),
		})
	case domain.ErrInvalidAccount, domain.ErrInvalidAccountDocumentNumber:
		http.Error(w, "invalid account document number", http.StatusBadRequest)
	case domain.ErrAccountAlreadyExists:
		http.Error(w, "account already exists", http.StatusConflict)
	default:
		http.Error(w, "internal server error: "+err.Error(), http.StatusInternalServerError)
	}
}

func (h *HttpHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var req CreateTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	transaction, err := h.transactionService.CreateTransaction(r.Context(), req.AccountID, req.OperationTypeID, req.Amount)

	switch err {
	case nil:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(&TransactionResponse{
			ID:              transaction.ID,
			AccountID:       transaction.AccountID,
			OperationTypeID: int(transaction.OperationTypeID),
			Amount:          transaction.Amount,
			EventDate:       transaction.EventDate.Format(time.RFC3339),
		})
	case domain.ErrInvalidTransaction, domain.ErrTransactionOperationTypeInvalid, domain.ErrTransactionAmountInvalid:
		log.Printf("%s", err.Error())
		http.Error(w, "invalid data: "+err.Error(), http.StatusBadRequest)
	case domain.ErrTransactionAccountNotExists:
		http.Error(w, "transaction account not exists", http.StatusBadRequest)
	default:
		http.Error(w, "internal server error: "+err.Error(), http.StatusInternalServerError)
	}
}
