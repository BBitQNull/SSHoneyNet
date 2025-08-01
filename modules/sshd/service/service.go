package sshd_service

import (
	"github.com/BBitQNull/SSHoneyNet/core/clientset"
	"github.com/BBitQNull/SSHoneyNet/core/sshd"
	proc_client "github.com/BBitQNull/SSHoneyNet/modules/dispatcher/client/procclient"
	sshserver "github.com/BBitQNull/SSHoneyNet/modules/sshd/service/sshserver"
	"github.com/BBitQNull/SSHoneyNet/pkg/utils/auth"
)

type sshdService struct {
	procClient proc_client.ProcManageClient
}

func NewSSHDService(procClient proc_client.ProcManageClient) sshd.SSHDService {
	return &sshdService{procClient: procClient}
}

func (s *sshdService) StartSSHServer(SSHDClient *clientset.ClientSetSSHD) {
	authSvc := &auth.SimpleAuthService{}
	sshserver.StartServer(authSvc, s.procClient, SSHDClient)
}
