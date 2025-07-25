package main

import (
	"log"
	"net"

	parser_service "github.com/BBitQNull/SSHoneyNet/modules/commandparser/service"
	parser_transport "github.com/BBitQNull/SSHoneyNet/modules/commandparser/transport"
	_ "github.com/BBitQNull/SSHoneyNet/modules/commands/uname"
	dispatch_service "github.com/BBitQNull/SSHoneyNet/modules/dispatcher/service"
	dispatch_transport "github.com/BBitQNull/SSHoneyNet/modules/dispatcher/transport"
	process_service "github.com/BBitQNull/SSHoneyNet/modules/procsystem/service"
	proc_transport "github.com/BBitQNull/SSHoneyNet/modules/procsystem/transport"
	parser_Pb "github.com/BBitQNull/SSHoneyNet/pb/cmdparser"
	dispatch_Pb "github.com/BBitQNull/SSHoneyNet/pb/dispatcher"
	proc_Pb "github.com/BBitQNull/SSHoneyNet/pb/procsystem"
	"google.golang.org/grpc"
)

func main() {
	// 命令解析
	parserSvc, err := parser_service.NewCmdParserService()
	if err != nil {
		log.Printf("NewCmdParserService error: %v", err)
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

	// 命令调用
	dispatchSvc := dispatch_service.NewDispatcherServer()
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

	// 进程系统
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
	select {}
}
