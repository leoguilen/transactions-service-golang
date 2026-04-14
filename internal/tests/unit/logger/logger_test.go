package logger

import (
	"testing"

	"github.com/leoguilen/transactions/internal/logger"
)

func TestNewLogEntry(t *testing.T) {
	event := "http.request"

	entry := logger.NewLogEntry(event)

	if entry.Event != event {
		t.Errorf("expected event %s, got %s", event, entry.Event)
	}
	if entry.Level != logger.LevelInfo {
		t.Errorf("expected level INFO, got %s", entry.Level)
	}
	if entry.Timestamp == "" {
		t.Errorf("expected timestamp to be set")
	}
}

func TestSetLogLevel_InfoForSuccess(t *testing.T) {
	tests := []struct {
		statusCode int
		level      logger.LogLevel
	}{
		{200, logger.LevelInfo},
		{201, logger.LevelInfo},
		{204, logger.LevelInfo},
		{301, logger.LevelInfo},
		{302, logger.LevelInfo},
	}

	for _, tt := range tests {
		entry := logger.NewLogEntry("test-event")
		entry.SetLogLevel(tt.statusCode)

		if entry.Level != tt.level {
			t.Errorf("for status %d, expected level %s, got %s", tt.statusCode, tt.level, entry.Level)
		}
	}
}

func TestSetLogLevel_WarnForClientError(t *testing.T) {
	tests := []struct {
		statusCode int
	}{
		{400},
		{401},
		{403},
		{404},
		{409},
		{422},
	}

	for _, tt := range tests {
		entry := logger.NewLogEntry("test-event")
		entry.SetLogLevel(tt.statusCode)

		if entry.Level != logger.LevelWarn {
			t.Errorf("for status %d, expected level WARN, got %s", tt.statusCode, entry.Level)
		}
	}
}

func TestSetLogLevel_ErrorForServerError(t *testing.T) {
	tests := []struct {
		statusCode int
	}{
		{500},
		{501},
		{502},
		{503},
		{504},
	}

	for _, tt := range tests {
		entry := logger.NewLogEntry("test-event")
		entry.SetLogLevel(tt.statusCode)

		if entry.Level != logger.LevelError {
			t.Errorf("for status %d, expected level ERROR, got %s", tt.statusCode, entry.Level)
		}
	}
}
