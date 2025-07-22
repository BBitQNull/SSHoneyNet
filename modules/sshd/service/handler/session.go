package handler

import (
	"fmt"
	"io"
	"strings"

	"github.com/BBitQNull/SSHoneyNet/pkg/model"
	"github.com/chzyer/readline"
	"github.com/gliderlabs/ssh"
)

func SessionHandler(echoReg *model.EchoRegistry) ssh.Handler {
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

		sessionID := string(s.Context().SessionID())
		writeChan := make(chan string, 100)
		echoReg.Register(sessionID, writeChan)
		defer echoReg.Unregister(sessionID)
		// test
		io.WriteString(s, "Welcome to Honeypot!\n")

		go func() {
			for msg := range writeChan {
				fmt.Fprintln(rl.Stdout(), msg)
			}
		}()
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
				writeChan <- "Yours input: " + line
			}
		}
	}
}
