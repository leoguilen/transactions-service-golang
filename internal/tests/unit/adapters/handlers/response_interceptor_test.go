package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/leoguilen/transactions/internal/adapters/handlers"
)

func TestNewResponseInterceptor(t *testing.T) {
	w := httptest.NewRecorder()
	interceptor := handlers.NewResponseInterceptor(w)

	if interceptor.ResponseWriter == nil {
		t.Errorf("expected ResponseWriter to be set")
	}
	if interceptor.GetStatusCode() != http.StatusOK {
		t.Errorf("expected initial status code to be 200, got %d", interceptor.GetStatusCode())
	}
}
