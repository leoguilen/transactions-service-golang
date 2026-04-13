package http

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/leoguilen/transactions/internal/adapters/db/postgres"
	"github.com/leoguilen/transactions/internal/adapters/handlers"
	"github.com/leoguilen/transactions/internal/core/services"
	"github.com/leoguilen/transactions/internal/tests/integration/utils"
)

var ConnStr string

func TestMain(m *testing.M) {
	ctx := context.Background()

	var teardown func()
	var err error

	ConnStr, teardown, err = utils.SetupPostgresContainer(ctx)
	if err != nil {
		log.Fatalf("failed to setup PostgreSQL container: %v", err)
	}
	defer teardown()

	code := m.Run()
	os.Exit(code)
}

func NewHttpHandler() *handlers.HttpHandler {
	accountRepo := postgres.NewAccountRepository(ConnStr)
	transactionRepo := postgres.NewTransactionRepository(ConnStr)

	accountService := services.NewAccountService(accountRepo)
	transactionService := services.NewTransactionService(accountRepo, transactionRepo)

	return handlers.NewHttpHandler(accountService, transactionService)
}
