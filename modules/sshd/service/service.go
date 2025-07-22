package sshd_service

import (
	"context"
	"errors"

	"github.com/BBitQNull/SSHoneyNet/core/sshd"
	"github.com/BBitQNull/SSHoneyNet/modules/sshd/service/server"
	"github.com/BBitQNull/SSHoneyNet/pkg/model"
	"github.com/BBitQNull/SSHoneyNet/pkg/utils/auth"
)

type sshdService struct {
	echo *model.EchoRegistry
}

func NewSSHDService(echo *model.EchoRegistry) sshd.SSHDService {
	return &sshdService{
		echo: echo,
	}
}

// 处理逻辑待完善
func (s *sshdService) EchoCommand(ctx context.Context, result model.CmdResult) (bool, error) {
	if result.NextPrompt == true {
		s.echo.Send("1", result.Output)
		return true, nil
	}
	return false, errors.New("nextpromt is false")
}

func (s *sshdService) StartSSHServer() {
	authSvc := &auth.SimpleAuthService{}
	server.StartServer(authSvc, s.echo)
}
