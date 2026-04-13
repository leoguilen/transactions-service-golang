package main

import (
	"log"
	"net/http"
	"os"

	"github.com/leoguilen/transactions/internal/adapters/db/postgres"
	"github.com/leoguilen/transactions/internal/adapters/handlers"
	"github.com/leoguilen/transactions/internal/core/services"
)

func init() {
	if _, exists := os.LookupEnv("DB_CONNECTION_STRING"); !exists {
		log.Fatal("DB_CONNECTION_STRING environment variable is required")
	}
}

func main() {
	postgresConnStr := os.Getenv("DB_CONNECTION_STRING")

	accountRepo := postgres.NewAccountRepository(postgresConnStr)
	transactionRepo := postgres.NewTransactionRepository(postgresConnStr)

	accountService := services.NewAccountService(accountRepo)
	transactionService := services.NewTransactionService(accountRepo, transactionRepo)

	httpHandler := handlers.NewHttpHandler(accountService, transactionService)

	mux := http.NewServeMux()
	httpHandler.RegisterRoutes(mux)

	addr := ":" + os.Getenv("HTTP_PORT")
	log.Printf("Starting HTTP server on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
