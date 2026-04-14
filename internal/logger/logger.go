package logger

import (
	"encoding/json"
	"fmt"
	"time"
)

type LogLevel string

const (
	LevelInfo  LogLevel = "INFO"
	LevelWarn  LogLevel = "WARN"
	LevelError LogLevel = "ERROR"
)

type LogEntry struct {
	Timestamp  string   `json:"timestamp"`
	Level      LogLevel `json:"level"`
	Event      string   `json:"event"`
	Method     string   `json:"method,omitempty"`
	Path       string   `json:"path,omitempty"`
	StatusCode int      `json:"status_code,omitempty"`
	DurationMs int64    `json:"duration_ms,omitempty"`
	Error      string   `json:"error,omitempty"`
}

func NewLogEntry(event string) *LogEntry {
	return &LogEntry{
		Timestamp: time.Now().UTC().Format(time.RFC3339Nano),
		Level:     LevelInfo,
		Event:     event,
	}
}

func (le *LogEntry) SetLogLevel(statusCode int) {
	switch {
	case statusCode >= 500:
		le.Level = LevelError
	case statusCode >= 400:
		le.Level = LevelWarn
	default:
		le.Level = LevelInfo
	}
}

func (le *LogEntry) Print() {
	data, _ := json.Marshal(le)
	fmt.Println(string(data))
}
