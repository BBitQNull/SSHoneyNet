package sshd_service

import (
	"github.com/BBitQNull/SSHoneyNet/core/sshd"
	sshserver "github.com/BBitQNull/SSHoneyNet/modules/sshd/service/sshserver"
	"github.com/BBitQNull/SSHoneyNet/pkg/utils/auth"
)

type sshdService struct{}

func NewSSHDService() sshd.SSHDService {
	return &sshdService{}
}

func (s *sshdService) StartSSHServer() {
	authSvc := &auth.SimpleAuthService{}
	sshserver.StartServer(authSvc)
}
