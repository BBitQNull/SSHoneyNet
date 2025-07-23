package handler

import (
	"strings"

	"github.com/chzyer/readline"
	"github.com/gliderlabs/ssh"
)

func SessionHandler() ssh.Handler {
	return func(s ssh.Session) {
		rl, err := readline.NewEx(&readline.Config{
			Prompt:          "honeypot> ",
			HistoryLimit:    100,
			InterruptPrompt: "^C",
			EOFPrompt:       "exit",
			Stdin:           s,
			Stdout:          s,
			Stderr:          s,
		})
		if err != nil {
			return
		}
		defer rl.Close()

		//		sessionID := string(s.Context().SessionID())

		for {
			line, err := rl.Readline()
			if err != nil {
				break
			}
			line = strings.TrimSpace(line)
			switch line {
			case "exit":
				return
			default:
				//命令解析和调用
			}
		}
	}
}
