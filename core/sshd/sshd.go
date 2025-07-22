package sshd

import (
	"context"

	"github.com/BBitQNull/SSHoneyNet/pkg/model"
)

type SSHDService interface {
	EchoCommand(ctx context.Context, result model.CmdResult) (bool, error)
	StartSSHServer()
}
