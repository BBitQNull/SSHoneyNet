package sshd

import (
	"context"

	"github.com/BBitQNull/SSHoneyNet/pkg/model"
)

type SSHDService interface {
	//	DeliverCommand(ctx context.Context) (string, error)
	EchoCommand(ctx context.Context) (model.CmdResult, error)
	StartSSHServer()
}
