package handler

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/BBitQNull/SSHoneyNet/core/clientset"
	loglog "github.com/BBitQNull/SSHoneyNet/core/log"
	proc_client "github.com/BBitQNull/SSHoneyNet/modules/dispatcher/client/procclient"
	sshd_client "github.com/BBitQNull/SSHoneyNet/modules/sshd/client"
	pb "github.com/BBitQNull/SSHoneyNet/pb/cmdparser"
	pbdis "github.com/BBitQNull/SSHoneyNet/pb/dispatcher"
	"github.com/chzyer/readline"
	"github.com/gliderlabs/ssh"
	"github.com/go-kit/kit/endpoint"
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

func LogUserInteraction(
	ctx context.Context,
	logEndpoint endpoint.Endpoint,
	timestamp time.Time,
	sessionID, userInput, output, ip string,
) {
	if logEndpoint == nil {
		log.Printf("logEndpoint is nil, cannot write log")
		return
	}

	req := &sshd_client.RawWriteLogRequest{
		LogEntry: loglog.LogEntry{
			Timestamp: timestamp,
			SessionID: sessionID,
			UserInput: userInput,
			Output:    output,
			IP:        ip,
		},
	}

	resp, err := logEndpoint(ctx, req)
	if err != nil {
		log.Printf("failed to write log: %v", err)
		return
	}

	if _, ok := resp.(*sshd_client.RawWriteLogResponse); !ok {
		log.Printf("unexpected response type from log endpoint: %T", resp)
	}
}

func SessionHandler(procClient proc_client.ProcManageClient, SSHDClient *clientset.ClientSetSSHD) ssh.Handler {
	return func(s ssh.Session) {
		ParserClient := SSHDClient.SSHDClient.CommandParser
		DispatcherClient := SSHDClient.SSHDClient.Dispatcher
		LogClient := SSHDClient.SSHDClient.WriteLog
		sessionID := string(s.Context().SessionID())
		remoteAddr := s.RemoteAddr()
		cwd := "/" // 当前所在目录，等待“cd”命令更新
		var rl *readline.Instance
		var err error
		rl, err = readline.NewEx(&readline.Config{
			Prompt:          "honeypot#> ",
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
					line := "Ctrl+L"
					output := "\033[H\033[2J"
					LogUserInteraction(context.Background(), LogClient, time.Now(), sessionID, line, "Ctrl+L", remoteAddr.String())
					s.Write([]byte(output))
					return 0, false // 阻止进一步处理
				case 26: // Ctrl+Z (ASCII 26)
					line := "Ctrl+Z"
					output := "\n[模拟] 当前进程已挂起 (实际上并未挂起)"
					LogUserInteraction(context.Background(), LogClient, time.Now(), sessionID, line, "Ctrl+Z", remoteAddr.String())
					s.Write([]byte(output))
					return 0, false
				}
				return r, true
			},
		})
		if err != nil {
			return
		}
		defer rl.Close()

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
			rl.SetPrompt(fmt.Sprintf("honeypot:%s#> ", cwd))
			rl.Refresh()
			line, err := rl.Readline()
			readTime := time.Now()
			if err != nil {
				break
			}
			line = strings.TrimSpace(line)

			if line == "" {
				LogUserInteraction(context.Background(), LogClient, readTime, sessionID, line, "", remoteAddr.String())
				continue
			}
			if line == "\n" {
				LogUserInteraction(context.Background(), LogClient, readTime, sessionID, line, "\n", remoteAddr.String())
				continue
			}
			switch line {
			case "exit":
				LogUserInteraction(context.Background(), LogClient, readTime, sessionID, line, "exit", remoteAddr.String())
				return
			case "clear":
				output := "\033[H\033[2J"
				LogUserInteraction(context.Background(), LogClient, readTime, sessionID, line, "clear", remoteAddr.String())
				rl.Clean()
				rl.Write([]byte(output))
				rl.Refresh()
				continue
			case "help":
				output := "还没有帮助^-^" + "\n"
				LogUserInteraction(context.Background(), LogClient, readTime, sessionID, line, output, remoteAddr.String())
				rl.Write([]byte(output))
				continue
			default:
				//命令解析
				parserRespRaw, err := ParserClient(context.Background(), &sshd_client.RawCmdParserRequest{
					Cmd: line,
				})
				if err != nil {
					log.Printf("parserRespRaw error: %v", err)
				}
				if parserRespRaw == nil {
					switch {
					case line == "0":
						output := "zsh: command not found: " + line + "\n"
						LogUserInteraction(context.Background(), LogClient, readTime, sessionID, line, output, remoteAddr.String())
						rl.Write([]byte(output))
						continue
					case re1to9.MatchString(line):
						output := "cd: no such entry in dir stack\n"
						LogUserInteraction(context.Background(), LogClient, readTime, sessionID, line, output, remoteAddr.String())
						rl.Write([]byte(output))
						continue
					case e10plus.MatchString(line):
						output := "zsh: command not found: " + line + "\n"
						LogUserInteraction(context.Background(), LogClient, readTime, sessionID, line, output, remoteAddr.String())
						rl.Write([]byte(output))
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
				dispatchRespRaw, err := DispatcherClient(context.Background(), &pbdis.DispatcherRequest{
					Ast:       astReq.Ast,
					SessionID: sessionID,
				})
				if err != nil {
					log.Printf("dispatchRespRaw error: %v", err)
				}
				result, ok := dispatchRespRaw.(sshd_client.RawCmdParserResponse)
				if !ok {
					log.Printf("dispatchRespRaw断言失败, want: %T", dispatchRespRaw)
					continue
				}
				if len(result.Result)+len(result.ErrMsg) > 0 {
					if result.ErrCode == 0 {
						output := result.Result
						LogUserInteraction(context.Background(), LogClient, readTime, sessionID, line, output, remoteAddr.String())
						rl.Write([]byte(strings.TrimRight(result.Result, "\n") + "\n"))
					} else {
						output := result.ErrMsg
						LogUserInteraction(context.Background(), LogClient, readTime, sessionID, line, output, remoteAddr.String())
						rl.Write([]byte(strings.TrimRight(result.ErrMsg, "\n") + "\n"))
					}
				}
			}
		}
	}
}
