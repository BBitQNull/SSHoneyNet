package ps

import (
	"context"
	"log"

	"github.com/BBitQNull/SSHoneyNet/core/dispatcher"
	proc_client "github.com/BBitQNull/SSHoneyNet/modules/dispatcher/client"
	"github.com/BBitQNull/SSHoneyNet/modules/sshd/service/handler"
	proc_Pb "github.com/BBitQNull/SSHoneyNet/pb/procsystem"
	"github.com/BBitQNull/SSHoneyNet/pkg/model"
	"github.com/BBitQNull/SSHoneyNet/pkg/utils/counterpid"
	"github.com/BBitQNull/SSHoneyNet/pkg/utils/exescript"
	"google.golang.org/grpc/metadata"
)

type PsHandler struct {
	procClient proc_client.ProcManageClient
}

func NewPsHandler(procClient proc_client.ProcManageClient) *PsHandler {
	return &PsHandler{procClient: procClient}
}

func (h *PsHandler) Execute(ctx context.Context, cmd exescript.ExecCommand) (dispatcher.CmdEcho, error) {
	// 从 metadata 中获取 session-id
	md, ok := metadata.FromIncomingContext(ctx)
	log.Printf("CmdHandler.Execute metadata: %+v, ok=%v", md, ok)
	if !ok {
		return dispatcher.CmdEcho{
			CmdResult: "missing metadata",
			ErrCode:   0,
			ErrMsg:    "",
		}, nil
	}

	sessionIDs := md.Get("session-id")
	if len(sessionIDs) == 0 {
		return dispatcher.CmdEcho{
			CmdResult: "session-id not found",
			ErrCode:   0,
			ErrMsg:    "",
		}, nil
	}
	sessionID := sessionIDs[0]

	// 查询对应 shell pid
	handler.SessionPidLock.RLock()
	shellPid, ok := handler.SessionPidMap[sessionID]
	handler.SessionPidLock.RUnlock()
	if !ok {
		return dispatcher.CmdEcho{
			CmdResult: "shell pid not found for session",
			ErrCode:   0,
			ErrMsg:    "",
		}, nil
	}
	psPid := counterpid.GeneratePidFromShell(shellPid)
	procResp, err := h.procClient.CreateProc(ctx, &proc_client.RawRequest{
		Command: "ps",
		Pid:     psPid,
		Ppid:    shellPid,
	})
	if err != nil {
		log.Panicln("CreateProc ps error")
	}
	defer func() {
		_, err := h.procClient.KillProc(ctx, &proc_client.RawRequest{
			Pid: psPid,
		})
		if err != nil {
			log.Println("kill ps error")
		}
	}()
	v, ok := procResp.(*proc_Pb.ProcResponse)
	if !ok {
		log.Println("assert error ps")
	}
	procs := make([]model.Process, 0, 5)
	for _, item := range v.Pcblist {
		if item.State == "R" {
			p := model.Process{
				PID:  item.Pid,
				TTY:  "pts/0",
				TIME: "00:00:00",
				CMD:  "bash",
			}
			procs = append(procs, p)
		}
	}
	return dispatcher.CmdEcho{
		CmdResult: model.FormatPsOutput(procs),
		ErrCode:   0,
		ErrMsg:    "",
	}, nil
}
