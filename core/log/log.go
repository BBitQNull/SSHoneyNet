package log

import "context"

type LogEntry struct {
	Timestamp string `json:"timestamp"`
	SessionID string `json:"session_id"`
	UserInput string `json:"user_input"`
	Output    string `json:"output"`
	IP        string `json:"ip,omitempty"`
}

type LogService interface {
	WriteLog(ctx context.Context, entry LogEntry) error
	GetLog(ctx context.Context) ([]LogEntry, error)
	GetLogSince(ctx context.Context, timestamp string) ([]LogEntry, error)
}

type LogWriter interface {
	Write(entry LogEntry) error
	ReadAll() ([]LogEntry, error)
}

type StreamLogReader interface {
	ReadSince(timestamp string) ([]LogEntry, error)
}
