package dispatcher

import (
	"context"

	"github.com/BBitQNull/SSHoneyNet/core/commandparser"
)

type CmdEcho struct {
	CmdResult string
	ErrCode   int32
	ErrMsg    string
}

type CmdDispatcherService interface {
	CmdDispatcher(ctx context.Context, ir commandparser.Script) (CmdEcho, error)
}
