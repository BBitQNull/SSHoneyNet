package sshd_test

import (
	"context"
	"fmt"
	"net"
	"testing"

	sshd_service "github.com/BBitQNull/SSHoneyNet/modules/sshd/service"
	sshd_transport "github.com/BBitQNull/SSHoneyNet/modules/sshd/transport"
	"github.com/BBitQNull/SSHoneyNet/pb"
	"github.com/BBitQNull/SSHoneyNet/pkg/model"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func init() {
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

func SSHDTest(t *testing.T) {
	conn, err := grpc.NewClient(
		"127.0.0.1:8972",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return
	}
	defer conn.Close()
	c := pb.NewEchoCmdClient(conn)
	resp, err := c.EchoCommand(context.Background(), &pb.EchoCmdRequest{
		Result: &pb.CmdResult{
			Output:   "hello sshd!",
			Exitcode: 0,
			Errmsg:   "no error!",
			Log:      "",
		},
	})
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, true, resp)
}
