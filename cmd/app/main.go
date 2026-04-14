package main

import (
	"log"
	"net/http"
	"os"

	"github.com/leoguilen/transactions/internal/adapters/db/postgres"
	"github.com/leoguilen/transactions/internal/adapters/handlers"
	"github.com/leoguilen/transactions/internal/core/services"
)

// @title Transactions Service API
// @version 1.0.0
// @description Account and Transaction Management Service
// @host localhost:5000
// @BasePath /
// @schemes http https

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

	loggedMux := handlers.LoggingMiddleware(mux)

	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = "5000"
	}
	addr := ":" + port

	log.Printf("Starting HTTP server on %s", addr)
	if err := http.ListenAndServe(addr, loggedMux); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
