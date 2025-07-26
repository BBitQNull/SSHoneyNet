package main

import (
	"log"
	"net"

	"github.com/BBitQNull/SSHoneyNet/core/clientset"
	parser_service "github.com/BBitQNull/SSHoneyNet/modules/commandparser/service"
	parser_transport "github.com/BBitQNull/SSHoneyNet/modules/commandparser/transport"
	"github.com/BBitQNull/SSHoneyNet/modules/commands/ps"
	"github.com/BBitQNull/SSHoneyNet/modules/commands/uname"
	dispatch_service "github.com/BBitQNull/SSHoneyNet/modules/dispatcher/service"
	dispatch_transport "github.com/BBitQNull/SSHoneyNet/modules/dispatcher/transport"
	fs_service "github.com/BBitQNull/SSHoneyNet/modules/filesystem/service"
	fs_transport "github.com/BBitQNull/SSHoneyNet/modules/filesystem/transport"
	process_service "github.com/BBitQNull/SSHoneyNet/modules/procsystem/service"
	proc_transport "github.com/BBitQNull/SSHoneyNet/modules/procsystem/transport"
	sshd_service "github.com/BBitQNull/SSHoneyNet/modules/sshd/service"
	parser_Pb "github.com/BBitQNull/SSHoneyNet/pb/cmdparser"
	dispatch_Pb "github.com/BBitQNull/SSHoneyNet/pb/dispatcher"
	fs_Pb "github.com/BBitQNull/SSHoneyNet/pb/filesystem"
	proc_Pb "github.com/BBitQNull/SSHoneyNet/pb/procsystem"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// proc 启动
	procSvc := process_service.NewProcessServer()
	procGs := proc_transport.NewGRPCProcServer(procSvc)
	procListener, err := net.Listen("tcp", ":9003")
	if err != nil {
		return
	}
	procS := grpc.NewServer()
	proc_Pb.RegisterProcManageServer(procS, procGs)
	go func() {
		err = procS.Serve(procListener)
		if err != nil {
			return
		}
	}()
	// fs 启动
	fsSvc := fs_service.NewFSService(fs_service.NewFileSystem())
	fsGs := fs_transport.NewFSServer(fsSvc)
	fsListener, err := net.Listen("tcp", ":9004")
	if err != nil {
		return
	}
	fsS := grpc.NewServer()
	fs_Pb.RegisterFileManageServer(fsS, fsGs)
	go func() {
		err = fsS.Serve(fsListener)
		if err != nil {
			return
		}
	}()
	// parser 启动
	parserSvc, err := parser_service.NewCmdParserService()
	if err != nil {
		log.Fatalf("NewCmdParserService error: %v", err)
	}
	parserGs := parser_transport.NewGRPCmdParserServer(parserSvc)
	parserListener, err := net.Listen("tcp", ":9001")
	if err != nil {
		return
	}
	parserS := grpc.NewServer()
	parser_Pb.RegisterCmdParserServer(parserS, parserGs)
	go func() {
		err = parserS.Serve(parserListener)
		if err != nil {
			return
		}
	}()

	// dispatcher连接proc
	// dispatcher client 调用 procsystem
	// dispatcher client 调用 filesystem
	connProc, err := grpc.NewClient(
		"127.0.0.1:9003",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Println("dispatcher connProc client start error")
	}
	defer connProc.Close()
	if connProc == nil {
		log.Fatal("Failed to connect to proc service")
	}

	connFs, err := grpc.NewClient(
		"127.0.0.1:9004",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Println("dispatcher connFs client start error")
	}
	defer connFs.Close()
	if connFs == nil {
		log.Fatal("Failed to connect to fs service")
	}
	clients := clientset.NewClientSet(connProc, connFs)

	// dispatcher初始化
	dispatchSvc := dispatch_service.NewDispatcherServer(clients)
	// 命令注册
	// uname
	unameHandler := uname.NewUnameHandler(clients.ProcClient)
	dispatchSvc.RegisterCmd("uname", unameHandler)
	// ps
	psHandler := ps.NewPsHandler(clients.ProcClient)
	dispatchSvc.RegisterCmd("ps", psHandler)

	// dispatcher 启动
	dispatchGs := dispatch_transport.NewCmdDispatcherServer(dispatchSvc)
	dispatchListener, err := net.Listen("tcp", ":9002")
	if err != nil {
		return
	}
	dispatchS := grpc.NewServer()
	dispatch_Pb.RegisterCmdEchoServer(dispatchS, dispatchGs)
	go func() {
		err = dispatchS.Serve(dispatchListener)
		if err != nil {
			return
		}
	}()

	// sshd启动
	sshd_service.NewSSHDService(clients.ProcClient).StartSSHServer()

	select {}
}
