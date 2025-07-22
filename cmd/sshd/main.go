package main

import (
	"fmt"
	"net"

	sshd_service "github.com/BBitQNull/SSHoneyNet/modules/sshd/service"
	sshd_transport "github.com/BBitQNull/SSHoneyNet/modules/sshd/transport"
	"github.com/BBitQNull/SSHoneyNet/pb"
	"github.com/BBitQNull/SSHoneyNet/pkg/model"
	"google.golang.org/grpc"
)

func main() {
	svc := sshd_service.NewSSHDService(model.NewEchoRegistry())
	svc.StartSSHServer()
	gs := sshd_transport.NewGRPCSSHDServer(svc)
	listener, err := net.Listen("tcp", ":8972")
	if err != nil {
		fmt.Printf("failed to listen: %v", err)
		return
	}
	s := grpc.NewServer()
	pb.RegisterEchoCmdServer(s, gs)
	if err = s.Serve(listener); err != nil {
		fmt.Printf("failed to serve: %v", err)
		return
	}
}
