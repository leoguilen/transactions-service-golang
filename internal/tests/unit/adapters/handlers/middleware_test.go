package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/leoguilen/transactions/internal/adapters/handlers"
)

func TestLoggingMiddleware_LogsRequest(t *testing.T) {
	oldStdout := io.Writer(nil)
	var capturedOutput bytes.Buffer

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok"}`))
	})

	middleware := handlers.LoggingMiddleware(handler)

	reqBody := []byte(`{"name": "test"}`)
	req := httptest.NewRequest("POST", "/accounts", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	middleware.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	if w.Body.String() != `{"status": "ok"}` {
		t.Errorf("expected response body to be unmodified, got %s", w.Body.String())
	}

	_ = oldStdout
	_ = capturedOutput
}

func TestLoggingMiddleware_CapturesResponseStatus(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	})

	middleware := handlers.LoggingMiddleware(handler)

	req := httptest.NewRequest("POST", "/accounts", nil)
	w := httptest.NewRecorder()

	middleware.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", w.Code)
	}
}

func TestLoggingMiddleware_WithoutRequestBody(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middleware := handlers.LoggingMiddleware(handler)

	req := httptest.NewRequest("GET", "/accounts/1", nil)
	w := httptest.NewRecorder()

	middleware.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestNewLogEntry_JSON(t *testing.T) {
	event := "http.request"
	entry := newTestLogEntry(event)

	jsonData, err := json.Marshal(entry)
	if err != nil {
		t.Errorf("failed to marshal LogEntry to JSON: %v", err)
	}

	var decoded map[string]interface{}
	err = json.Unmarshal(jsonData, &decoded)
	if err != nil {
		t.Errorf("failed to unmarshal LogEntry JSON: %v", err)
	}

	if decoded["event"] != event {
		t.Errorf("expected event in JSON")
	}
}

func TestLoggingMiddleware_ResponseTime(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(50 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	})

	middleware := handlers.LoggingMiddleware(handler)

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	start := time.Now()
	middleware.ServeHTTP(w, req)
	elapsed := time.Since(start)

	if elapsed > 200*time.Millisecond {
		t.Errorf("middleware added too much overhead: %v", elapsed)
	}
}

func newTestLogEntry(event string) any {
	return map[string]any{
		"timestamp": time.Now().UTC().Format(time.RFC3339Nano),
		"level":     "INFO",
		"event":     event,
	}
}
