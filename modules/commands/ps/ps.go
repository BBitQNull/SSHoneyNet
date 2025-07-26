package ps

import (
	"context"
	"log"

	"github.com/BBitQNull/SSHoneyNet/core/dispatcher"
	proc_client "github.com/BBitQNull/SSHoneyNet/modules/dispatcher/client"
	"github.com/BBitQNull/SSHoneyNet/modules/sshd/service/handler"
	"github.com/BBitQNull/SSHoneyNet/pkg/model"
	"github.com/BBitQNull/SSHoneyNet/pkg/utils/counterpid"
	"github.com/BBitQNull/SSHoneyNet/pkg/utils/exescript"
)

type PsHandler struct {
	procClient proc_client.ProcManageClient
}

func NewPsHandler(procClient proc_client.ProcManageClient) *PsHandler {
	return &PsHandler{procClient: procClient}
}

func (h *PsHandler) Execute(ctx context.Context, cmd exescript.ExecCommand, sessionID string) (dispatcher.CmdEcho, error) {
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
	_, err := h.procClient.CreateProc(ctx, &proc_client.RawRequest{
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
	resp, err := h.procClient.ListProc(ctx, &proc_client.RawRequest{
		Pid:     psPid,
		Ppid:    shellPid,
		Command: "ps",
	})
	v, ok := resp.(*proc_client.ListProcResponse)
	if !ok {
		log.Printf("assert error ps: unexpected type %T, value: %#v", resp, resp)
		log.Println("assert error ps")
	}
	procs := make([]model.Process, 0, 5)
	for _, item := range v.Processes {
		if item == nil {
			continue
		}
		if item.State == "R" {
			p := model.Process{
				PID:  item.PID,
				TTY:  "pts/0",
				TIME: "00:00:00",
				CMD:  item.Command,
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
