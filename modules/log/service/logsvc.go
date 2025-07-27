package log_service

import (
	"bufio"
	"context"
	"encoding/json"
	"os"
	"sync"
	"time"

	"github.com/BBitQNull/SSHoneyNet/core/log"
)

type LogServer struct {
	writer          log.LogWriter
	StreamLogReader log.StreamLogReader
}

type FileLogWriter struct {
	mu       sync.RWMutex
	filePath string
}

func NewFileLogWriter(path string) *FileLogWriter {
	return &FileLogWriter{
		filePath: path,
	}
}

func NewLogServer(writer log.LogWriter, streamLogReader log.StreamLogReader) *LogServer {
	return &LogServer{
		writer:          writer,
		StreamLogReader: streamLogReader,
	}
}

func (s *LogServer) WriteLog(ctx context.Context, entry log.LogEntry) error {
	return s.writer.Write(entry)
}

func (s *LogServer) GetLog(ctx context.Context) ([]log.LogEntry, error) {
	return s.writer.ReadAll()
}

func (s *LogServer) GetLogSince(ctx context.Context, t time.Time) ([]log.LogEntry, error) {
	return s.StreamLogReader.ReadSince(t)
}

func (w *FileLogWriter) Write(entry log.LogEntry) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	file, err := os.OpenFile(w.filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	data, err := json.Marshal(entry)
	if err != nil {
		return err
	}
	_, err = file.Write(append(data, '\n'))
	return err
}

func (w *FileLogWriter) ReadAll() ([]log.LogEntry, error) {
	w.mu.RLock()
	defer w.mu.RUnlock()
	file, err := os.OpenFile(w.filePath, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var logs []log.LogEntry
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var entry log.LogEntry
		if err := json.Unmarshal(scanner.Bytes(), &entry); err == nil {
			logs = append(logs, entry)
		} else {
			continue
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return logs, nil
}

func (w *FileLogWriter) ReadSince(t time.Time) ([]log.LogEntry, error) {
	w.mu.RLock()
	defer w.mu.RUnlock()

	file, err := os.Open(w.filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var results []log.LogEntry
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		var entry log.LogEntry
		if err := json.Unmarshal(line, &entry); err != nil {
			continue
		}
		if entry.Timestamp.After(t) {
			results = append(results, entry)
		}
	}
	return results, scanner.Err()
}
