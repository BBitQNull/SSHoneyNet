package model

import (
	"fmt"
	"strings"
)

type Process struct {
	PID  int64
	TTY  string
	TIME string
	CMD  string
}

func FormatPsOutput(processes []Process) string {
	var b strings.Builder
	// 表头
	b.WriteString(fmt.Sprintf("%-8s %-12s %-10s %s\n", "PID", "TTY", "TIME", "CMD"))

	// 每行进程
	for _, p := range processes {
		b.WriteString(fmt.Sprintf("%-8d %-12s %-10s %s\n", p.PID, p.TTY, p.TIME, p.CMD))
	}
	return b.String()
}
