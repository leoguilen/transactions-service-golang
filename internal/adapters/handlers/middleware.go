package handlers

import (
	"net/http"
	"time"

	"github.com/leoguilen/transactions/internal/logger"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logRequest(r)
		interceptor := NewResponseInterceptor(w)
		startTime := time.Now()
		next.ServeHTTP(interceptor, r)
		duration := time.Since(startTime)
		logResponse(interceptor, duration)
	})
}

func logRequest(r *http.Request) {
	entry := logger.NewLogEntry("http.request")
	entry.Method = r.Method
	entry.Path = r.RequestURI
	entry.Print()
}

func logResponse(interceptor *ResponseInterceptor, duration time.Duration) {
	entry := logger.NewLogEntry("http.response")
	entry.StatusCode = interceptor.GetStatusCode()
	entry.DurationMs = duration.Milliseconds()
	entry.SetLogLevel(entry.StatusCode)
	entry.Print()
}
