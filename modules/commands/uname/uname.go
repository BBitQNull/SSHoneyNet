package uname

import (
	"context"
	"fmt"

	"github.com/BBitQNull/SSHoneyNet/core/dispatcher"
	proc_client "github.com/BBitQNull/SSHoneyNet/modules/dispatcher/client"
	"github.com/BBitQNull/SSHoneyNet/pkg/utils/exescript"
)

type UnameHandler struct {
	procClient proc_client.ProcManageClient
}

const (
	UNAME_A = "Linux myhostname 5.15.0-78-generic #85-Ubuntu SMP Fri Jul 7 15:25:09 UTC 2023 x86_64 x86_64 x86_64 GNU/Linux"
)

func NewUnameHandler(procClient proc_client.ProcManageClient) *UnameHandler {
	return &UnameHandler{procClient: procClient}
}

func (h *UnameHandler) Execute(ctx context.Context, cmd exescript.ExecCommand, sessionID string) (dispatcher.CmdEcho, error) {
	if cmd.Name == "uname" {
		var result dispatcher.CmdEcho
		fmt.Println("uname flags:")
		for k := range cmd.Flags {
			fmt.Printf("flag key raw bytes: [% x], as string: [%s]\n", []byte(k), k)
		}
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
			CmdResult: UNAME_A,
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
