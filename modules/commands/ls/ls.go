package ls

import (
	"context"

	"github.com/BBitQNull/SSHoneyNet/core/dispatcher"
	fs_client "github.com/BBitQNull/SSHoneyNet/modules/dispatcher/client/fsclient"
	"github.com/BBitQNull/SSHoneyNet/pkg/utils/exescript"
)

type LsHandler struct {
	procClient fs_client.FSManageClient
}

func NewLsHandler(procClient fs_client.FSManageClient) *LsHandler {
	return &LsHandler{procClient: procClient}
}

func (h *LsHandler) Execute(ctx context.Context, cmd exescript.ExecCommand, sessionID string) (dispatcher.CmdEcho, error) {

}
