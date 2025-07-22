package handler

import (
	"bufio"
	"io"
	"strings"

	"github.com/BBitQNull/SSHoneyNet/pkg/model"
	"github.com/gliderlabs/ssh"
)

func SessionHandler(echoReg *model.EchoRegistry) ssh.Handler {
	return func(s ssh.Session) {
		sessionID := string(s.Context().SessionID())
		writeChan := make(chan string, 100)
		echoReg.Register(sessionID, writeChan)
		defer echoReg.Unregister(sessionID)

		go func() {
			for msg := range writeChan {
				io.WriteString(s, msg)
			}
		}()
		scanner := bufio.NewScanner(s)
		if scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "exit" {
				return
			}
			// grpc client commandparser

		}
	}
}
