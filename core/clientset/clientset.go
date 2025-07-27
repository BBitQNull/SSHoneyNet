package clientset

import (
	fs_client "github.com/BBitQNull/SSHoneyNet/modules/dispatcher/client/fsclient"
	proc_client "github.com/BBitQNull/SSHoneyNet/modules/dispatcher/client/procclient"
	sshd_client "github.com/BBitQNull/SSHoneyNet/modules/sshd/client"
	"google.golang.org/grpc"
)

// dispatcher client 调用 procsystem
// dispatcher client 调用 filesystem

type ClientSetDispatcher struct {
	ProcClient proc_client.ProcManageClient
	FSClient   fs_client.FSManageClient
}

// sshd client 调用 log
// sshd client 调用 cmdparser
// sshd client 调用 dispatcher
type ClientSetSSHD struct {
	SSHDClient *sshd_client.SSHDManageClient
}

func NewClientSetDispatcher(connProc, connFS *grpc.ClientConn) *ClientSetDispatcher {
	return &ClientSetDispatcher{
		ProcClient: proc_client.NewProcManageClient(connProc),
		FSClient:   fs_client.NewFSManageClient(connFS),
	}
}

func NewClientSetSSHD(connLog, connParser, connDispatcher *grpc.ClientConn) *ClientSetSSHD {
	return &ClientSetSSHD{
		SSHDClient: sshd_client.NewSSHDManageClient(connLog, connParser, connDispatcher),
	}
}
