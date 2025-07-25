package main

import sshd_service "github.com/BBitQNull/SSHoneyNet/modules/sshd/service"

func main() {
	// sshd client 调用 parser&dispatcher
	sshd_service.NewSSHDService().StartSSHServer()
	// dispatcher client 调用 procsystem
}
