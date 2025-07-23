package uname

import (
	"errors"

	"github.com/BBitQNull/SSHoneyNet/core/dispatcher"
	"github.com/BBitQNull/SSHoneyNet/pkg/utils/exescript"
)

type UnameHandler struct{}

const (
	UNAME_A = "Linux myhostname 5.15.0-78-generic #85-Ubuntu SMP Fri Jul 7 15:25:09 UTC 2023 x86_64 x86_64 x86_64 GNU/Linux"
)

func (h *UnameHandler) Execute(cmd exescript.ExecCommand) (dispatcher.CmdEcho, error) {
	if cmd.Name == "uname" {
		var result dispatcher.CmdEcho
		for _, flag := range cmd.Flags {
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
		return result, errors.New("illegal option")
	}
	return dispatcher.CmdEcho{
		CmdResult: "",
		ErrCode:   1,
		ErrMsg:    "command not found: " + cmd.Name,
	}, errors.New("command not found")
}
