package main

import sshd_service "github.com/BBitQNull/SSHoneyNet/modules/sshd/service"

func main() {
	sshd_service.NewSSHDService().StartSSHServer()
}
