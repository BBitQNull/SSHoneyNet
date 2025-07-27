package sshd

import "github.com/BBitQNull/SSHoneyNet/core/clientset"

type SSHDService interface {
	StartSSHServer(SSHDClient *clientset.ClientSetSSHD)
}
