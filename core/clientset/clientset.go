package clientset

import (
	fs_client "github.com/BBitQNull/SSHoneyNet/modules/dispatcher/client/fsclient"
	proc_client "github.com/BBitQNull/SSHoneyNet/modules/dispatcher/client/procclient"
	"google.golang.org/grpc"
)

// dispatcher client 调用 procsystem
// dispatcher client 调用 filesystem

type ClientSet struct {
	ProcClient proc_client.ProcManageClient
	FSClient   fs_client.FSManageClient
}

func NewClientSet(connProc, connFS *grpc.ClientConn) *ClientSet {
	return &ClientSet{
		ProcClient: proc_client.NewProcManageClient(connProc),
		FSClient:   fs_client.NewFSManageClient(connFS),
	}
}
