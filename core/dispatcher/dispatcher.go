package dispatcher

import (
	"context"

	"github.com/BBitQNull/SSHoneyNet/pkg/utils/exescript"
)

type CmdEcho struct {
	CmdResult string
	ErrCode   int16
	ErrMsg    string
}

type CmdDispatcherService interface {
	CmdDispatcher(ctx context.Context, ir exescript.ExecScript) (CmdEcho, error)
}
