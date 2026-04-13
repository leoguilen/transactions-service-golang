package postgres

import (
	"context"
	"log"
	"os"
	"testing"

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
