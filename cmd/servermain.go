package main

import (
	"log"
	"net"

	parser_service "github.com/BBitQNull/SSHoneyNet/modules/commandparser/service"
	parser_transport "github.com/BBitQNull/SSHoneyNet/modules/commandparser/transport"
	_ "github.com/BBitQNull/SSHoneyNet/modules/commands/uname"
	dispatch_service "github.com/BBitQNull/SSHoneyNet/modules/dispatcher/service"
	dispatch_transport "github.com/BBitQNull/SSHoneyNet/modules/dispatcher/transport"
	parserPb "github.com/BBitQNull/SSHoneyNet/pb/cmdparser"
	dispatchPb "github.com/BBitQNull/SSHoneyNet/pb/dispatcher"
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
	parserPb.RegisterCmdParserServer(parserS, parserGs)
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
	dispatchPb.RegisterCmdEchoServer(dispatchS, dispatchGs)
	go func() {
		err = dispatchS.Serve(dispatchListener)
		if err != nil {
			return
		}
	}()
	select {}
}
