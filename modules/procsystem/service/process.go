package process_service

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/BBitQNull/SSHoneyNet/core/procsystem"
)

type ProcessTable struct {
	mu      sync.RWMutex
	process map[int64]*(procsystem.PCB)
	nextPID int64
}

func (pt *ProcessTable) CreateProcess(ctx context.Context, name string, pid int64, ppid int64) (*procsystem.PCB, error) {
	pt.mu.Lock()
	defer pt.mu.Unlock()
	pt.process[pid] = &procsystem.PCB{
		PID:       pid,
		Command:   name,
		TGID:      pid,
		PPID:      ppid,
		State:     "S",
		CreatedAt: time.Now(),
		ExitedAt:  nil,
	}
	return pt.process[pid], nil
}

// 等待后台协程清理
func (pt *ProcessTable) KillProcess(ctx context.Context, pid int64) error {
	pt.mu.Lock()
	defer pt.mu.Unlock()
	proc, ok := pt.process[pid]
	if !ok {
		return errors.New("process not found")
	}
	now := time.Now()
	proc.State = "Z"
	proc.ExitedAt = &now
	return nil
}

func (pt *ProcessTable) GetProcess(ctx context.Context, pid int64) (*procsystem.PCB, error) {
	pt.mu.Lock()
	defer pt.mu.Unlock()
	proc, ok := pt.process[pid]
	if !ok {
		return nil, errors.New("process not found")
	}
	return proc, nil
}

func (pt *ProcessTable) ListProcess(ctx context.Context) ([]*procsystem.PCB, error) {
	values := make([]*procsystem.PCB, 0, len(pt.process))
	for _, v := range pt.process {
		values = append(values, v)
	}
	return values, nil
}

func (pt *ProcessTable) CleanupZumbies(ttl time.Duration) {
	pt.mu.Lock()
	defer pt.mu.Unlock()
	now := time.Now()
	for pid, proc := range pt.process {
		if proc.State == "Z" && proc.ExitedAt != nil {
			if now.Sub(*proc.ExitedAt) > ttl {
				delete(pt.process, pid)
				fmt.Printf("[Reaper] Cleaned up process &d\n", pid)
			}
		}
	}
}
func StartReaper(pm *ProcessTable, interval time.Duration, ttl time.Duration) {
	go func() {
		for {
			time.Sleep(interval)
			pm.CleanupZumbies(ttl)
		}
	}()
}
