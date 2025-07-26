package uname

import (
	"context"

	"github.com/BBitQNull/SSHoneyNet/core/dispatcher"
	fs_client "github.com/BBitQNull/SSHoneyNet/modules/dispatcher/client/fsclient"
	proc_client "github.com/BBitQNull/SSHoneyNet/modules/dispatcher/client/procclient"
	"github.com/BBitQNull/SSHoneyNet/pkg/utils/exescript"
)

type UnameHandler struct {
	procClient proc_client.ProcManageClient
	fsClient   fs_client.FSManageClient
}

const (
	UNAME_A = "Linux 80a0c3540f07 6.8.0-36-generic #36-Ubuntu SMP PREEMPT_DYNAMIC Mon Jun 10 10:49:14 UTC 2024 x86_64 Linux"
)

func NewUnameHandler(procClient proc_client.ProcManageClient, fsClient fs_client.FSManageClient) *UnameHandler {
	return &UnameHandler{
		procClient: procClient,
		fsClient:   fsClient,
	}
}

func (h *UnameHandler) Execute(ctx context.Context, cmd exescript.ExecCommand, sessionID string) (dispatcher.CmdEcho, error) {
	if cmd.Name == "uname" {
		var result dispatcher.CmdEcho
		if len(cmd.Flags) != 0 {
			for flag := range cmd.Flags {
				switch flag {
				case "a":
					return dispatcher.CmdEcho{
						CmdResult: UNAME_A,
						ErrCode:   0,
						ErrMsg:    "",
					}, nil
				default:
					result = dispatcher.CmdEcho{
						CmdResult: "",
						ErrCode:   1,
						ErrMsg:    "uname: illegal option --" + flag + "\nusage: uname [-amnoprsv]",
					}
				}
			}
			return result, nil
		}
		return dispatcher.CmdEcho{
			CmdResult: "Linux",
			ErrCode:   0,
			ErrMsg:    "",
		}, nil
	}
	return dispatcher.CmdEcho{
		CmdResult: "",
		ErrCode:   1,
		ErrMsg:    "command not found: " + cmd.Name,
	}, nil
}
