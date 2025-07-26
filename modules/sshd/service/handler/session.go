package handler

import (
	"context"
	"log"
	"math/rand"
	"regexp"
	"strings"
	"sync"
	"time"

	proc_client "github.com/BBitQNull/SSHoneyNet/modules/dispatcher/client/procclient"
	sshd_client "github.com/BBitQNull/SSHoneyNet/modules/sshd/client"
	pb "github.com/BBitQNull/SSHoneyNet/pb/cmdparser"
	pbdis "github.com/BBitQNull/SSHoneyNet/pb/dispatcher"
	"github.com/chzyer/readline"
	"github.com/gliderlabs/ssh"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	re1to9         = regexp.MustCompile(`^[1-9]$`)
	e10plus        = regexp.MustCompile(`^[1-9]\d+$`)
	rng            = rand.New(rand.NewSource(time.Now().UnixNano()))
	min, max       = 3378, 4679
	SessionPidMap  = make(map[string]int64)
	SessionPidLock sync.RWMutex
)

func getRandom() int64 {
	return int64(min + rand.Intn(max-min))
}

func SessionHandler(procClient proc_client.ProcManageClient) ssh.Handler {
	return func(s ssh.Session) {
		var rl *readline.Instance
		var err error
		rl, err = readline.NewEx(&readline.Config{
			Prompt:          "honeypot> ",
			HistoryLimit:    100,
			InterruptPrompt: "^C",
			EOFPrompt:       "exit",
			Stdin:           s,
			Stdout:          s,
			Stderr:          s,
			// 处理 Ctrl+L、Ctrl+Z
			FuncFilterInputRune: func(r rune) (rune, bool) {
				switch r {
				case 12: // Ctrl+L (ASCII 12)
					//	rl.Clean()
					s.Write([]byte("\033[H\033[2J"))
					//	rl.Refresh()
					return 0, false // 阻止进一步处理
				case 26: // Ctrl+Z (ASCII 26)
					s.Write([]byte("\n[模拟] 当前进程已挂起 (实际上并未挂起)"))
					return 0, false
				}
				return r, true
			},
		})
		if err != nil {
			return
		}
		defer rl.Close()

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
		if err != nil {
			return
		}
		defer conndispatch.Close()
		parserendpoint := sshd_client.MakeCmdParserEndpoint(connparser)
		dispatchendpoint := sshd_client.MakeCmdDispatchEndpoint(conndispatch)

		sessionID := string(s.Context().SessionID())
		// 创建shell进程
		shellResp, err := procClient.CreateProc(context.Background(), &proc_client.RawRequest{
			Command: "/bin/zsh",
			Pid:     getRandom(),
			Ppid:    3377,
		})
		if err != nil {
			log.Printf("Failed to create shell process: %v", err)
			s.Exit(1)
			return
		}
		v, ok := shellResp.(*proc_client.RawProcResponse)
		if !ok {
			log.Printf("shellResp type: %T", shellResp)
			log.Printf("assert fail shellResp")
			return
		}
		SessionPidLock.Lock()
		SessionPidMap[sessionID] = v.Process.PID
		SessionPidLock.Unlock()

		// 交互
		for {
			line, err := rl.Readline()
			if err != nil {
				break
			}
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			if line == "\n" {
				continue
			}
			switch line {
			case "exit":
				return
			case "clear":
				rl.Clean()
				rl.Write([]byte("\033[H\033[2J"))
				rl.Refresh()
				continue
			case "help":
				rl.Write([]byte("还没有帮助^-^" + "\n"))
				continue
			default:
				//命令解析
				parserRespRaw, err := parserendpoint(context.Background(), &pb.CmdParserRequest{
					Cmd: line,
				})
				if err != nil {
					log.Printf("parserRespRaw error: %v", err)
				}
				if parserRespRaw == nil {
					switch {
					case line == "0":
						rl.Write([]byte("zsh: command not found: " + line + "\n"))
						continue
					case re1to9.MatchString(line):
						rl.Write([]byte("cd: no such entry in dir stack\n"))
						continue
					case e10plus.MatchString(line):
						rl.Write([]byte("zsh: command not found: " + line + "\n"))
						continue
					}
					continue
				}
				// 命令调用
				astReq, ok := parserRespRaw.(*pb.CmdParserResponse)
				if !ok {
					log.Printf("parserRespRaw断言失败")
					continue
				}
				dispatchRespRaw, err := dispatchendpoint(context.Background(), &pbdis.DispatcherRequest{
					Ast:       astReq.Ast,
					SessionID: sessionID,
				})
				if err != nil {
					log.Printf("dispatchRespRaw error: %v", err)
				}
				result, ok := dispatchRespRaw.(*pbdis.DispatcherResponse)
				if !ok {
					log.Printf("dispatchRespRaw断言失败")
					continue
				}
				if result.Errcode == 0 {
					rl.Write([]byte(result.Cmdresult + "\n"))
				} else {
					rl.Write([]byte(result.Errmsg + "\n"))
				}
			}
		}
	}
}
