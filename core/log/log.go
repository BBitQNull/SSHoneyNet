package log

import (
	"context"
	"time"
)

type LogEntry struct {
	Timestamp time.Time `json:"timestamp"`
	SessionID string    `json:"session_id"`
	UserInput string    `json:"user_input"`
	Output    string    `json:"output"`
	IP        string    `json:"ip,omitempty"`
}

type LogService interface {
	WriteLog(ctx context.Context, entry LogEntry) error
	GetLog(ctx context.Context) ([]LogEntry, error)
	GetLogSince(ctx context.Context, t time.Time) ([]LogEntry, error)
}

type LogWriter interface {
	Write(entry LogEntry) error
	ReadAll() ([]LogEntry, error)
}

type StreamLogReader interface {
	ReadSince(t time.Time) ([]LogEntry, error)
}
