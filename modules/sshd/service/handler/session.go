package handler

import (
	"context"
	"strings"

	sshd_client "github.com/BBitQNull/SSHoneyNet/modules/sshd/client"
	pb "github.com/BBitQNull/SSHoneyNet/pb/cmdparser"
	pbdis "github.com/BBitQNull/SSHoneyNet/pb/dispatcher"
	"github.com/chzyer/readline"
	"github.com/gliderlabs/ssh"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
		connparser, err := grpc.NewClient(
			"127.0.0.1:9001",
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			return
		}
		defer connparser.Close()
		conndispatch, err := grpc.NewClient(
			"127.0.0.1:9002",
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		defer conndispatch.Close()
		parserendpoint := sshd_client.MakeCmdParserEndpoint(connparser)
		dispatchendpoint := sshd_client.MakeCmdDispatchEndpoint(conndispatch)

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
				//命令解析
				parserRespRaw, err := parserendpoint(context.Background(), &pb.CmdParserRequest{
					Cmd: line,
				})
				if err != nil {
					break
				}
				// 命令调用
				astReq := parserRespRaw.(*pbdis.Script)
				dispatchRespRaw, err := dispatchendpoint(context.Background(), &pbdis.DispatcherRequest{
					Ast: astReq,
				})
				if err != nil {
					break
				}
				result := dispatchRespRaw.(*pbdis.DispatcherResponse)
				if result.Errcode == 0 {
					rl.Write([]byte(result.Cmdresult))
				} else {
					rl.Write([]byte(result.Errmsg))
				}
			}
		}
	}
}
