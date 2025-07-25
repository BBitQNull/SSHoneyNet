package procsystem

import (
	"context"
	"time"
)

type PCB struct {
	PID       int64
	PPID      int64
	TGID      int64
	Command   string
	State     string
	User      string
	CreatedAt time.Time
	ExitedAt  *time.Time
}

type ProcessManager interface {
	CreateProcess(ctx context.Context, name string, pid int64, ppid int64) (*PCB, error)
	KillProcess(ctx context.Context, pid int64) error
	GetProcess(ctx context.Context, pid int64) (*PCB, error)
	ListProcess(ctx context.Context) ([]*PCB, error)
	CleanupZumbies(ttl time.Duration)
}

var stateDescriptions = map[string]string{
	"R": "Running",
	"S": "Sleeping",
	"D": "Uninterruptible Sleep",
	"T": "Stopped",
	"Z": "Zombie",
	"X": "Dead",
}
