package service

import (
	"context"

	"github.com/BBitQNull/SSHoneyNet/core/sshd"
	"github.com/BBitQNull/SSHoneyNet/modules/sshd/service/server"
	"github.com/BBitQNull/SSHoneyNet/pkg/model"
	"github.com/BBitQNull/SSHoneyNet/pkg/utils/auth"
)

type sshdService struct{}

func NewSSHDService() sshd.SSHDService {
	return &sshdService{}
}

/*
	func (s *sshdService) DeliverCommand(ctx context.Context) (string, error) {
		return "stub", nil
	}
*/
func (s *sshdService) EchoCommand(ctx context.Context) (model.CmdResult, error) {
	return model.CmdResult{}, nil
}

func (s *sshdService) StartSSHServer() {
	authSvc := &auth.SimpleAuthService{}
	server.StartServer(authSvc)
}
