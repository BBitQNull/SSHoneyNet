package clientset

import (
	proc_client "github.com/BBitQNull/SSHoneyNet/modules/dispatcher/client"
	"google.golang.org/grpc"
)

// dispatcher client 调用 procsystem
// dispatcher client 调用 filesystem

type ClientSet struct {
	ProcClient proc_client.ProcManageClient
	//	FSClient   fsclient.FSManageClient
}

func NewClientSet(connProc *grpc.ClientConn) *ClientSet {
	return &ClientSet{
		ProcClient: proc_client.NewProcManageClient(connProc),
		//	FSClient:   fsclient.NewFSManageClient(connFS),
	}
}
