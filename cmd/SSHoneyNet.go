package main

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/BBitQNull/SSHoneyNet/core/clientset"
	parser_service "github.com/BBitQNull/SSHoneyNet/modules/commandparser/service"
	parser_transport "github.com/BBitQNull/SSHoneyNet/modules/commandparser/transport"
	"github.com/BBitQNull/SSHoneyNet/modules/commands/ls"
	"github.com/BBitQNull/SSHoneyNet/modules/commands/ps"
	"github.com/BBitQNull/SSHoneyNet/modules/commands/uname"
	dispatch_service "github.com/BBitQNull/SSHoneyNet/modules/dispatcher/service"
	dispatch_transport "github.com/BBitQNull/SSHoneyNet/modules/dispatcher/transport"
	fs_service "github.com/BBitQNull/SSHoneyNet/modules/filesystem/service"
	fs_transport "github.com/BBitQNull/SSHoneyNet/modules/filesystem/transport"
	log_endpoint "github.com/BBitQNull/SSHoneyNet/modules/log/endpoint"
	log_service "github.com/BBitQNull/SSHoneyNet/modules/log/service"
	log_transport "github.com/BBitQNull/SSHoneyNet/modules/log/transport/grpc"
	log_httptransport "github.com/BBitQNull/SSHoneyNet/modules/log/transport/http"
	process_service "github.com/BBitQNull/SSHoneyNet/modules/procsystem/service"
	proc_transport "github.com/BBitQNull/SSHoneyNet/modules/procsystem/transport"
	sshd_service "github.com/BBitQNull/SSHoneyNet/modules/sshd/service"
	parser_Pb "github.com/BBitQNull/SSHoneyNet/pb/cmdparser"
	dispatch_Pb "github.com/BBitQNull/SSHoneyNet/pb/dispatcher"
	fs_Pb "github.com/BBitQNull/SSHoneyNet/pb/filesystem"
	log_Pb "github.com/BBitQNull/SSHoneyNet/pb/log"
	proc_Pb "github.com/BBitQNull/SSHoneyNet/pb/procsystem"
	"github.com/BBitQNull/SSHoneyNet/pkg/utils/pathconfig"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// 获取path
	jsonPath := pathconfig.GetPath(pathconfig.JSON)
	logPath := pathconfig.GetPath(pathconfig.LOG)
	fmt.Println("[DEBUG] logPath:", logPath)
	fmt.Println("[DEBUG] jsonPath:", jsonPath)

	// server
	// proc 启动
	procSvc := process_service.NewProcessServer()
	procGs := proc_transport.NewGRPCProcServer(procSvc)
	procListener, err := net.Listen("tcp", ":9001")
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
	fs, err := fs_service.NewFileSystem(jsonPath)
	if err != nil {
		log.Fatal(err)
	}
	fsSvc := fs_service.NewFSService(fs)
	fsGs := fs_transport.NewFSServer(fsSvc)
	fsListener, err := net.Listen("tcp", ":9002")
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
	parserListener, err := net.Listen("tcp", ":9003")
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
	// log grpc 启动
	writer := log_service.NewFileLogWriter(logPath)
	logSvc := log_service.NewLogServer(writer, writer)
	logGs := log_transport.NewWriteLogServer(logSvc)

	logListener, err := net.Listen("tcp", ":9005")
	if err != nil {
		return
	}
	logS := grpc.NewServer()
	log_Pb.RegisterLogServiceServer(logS, logGs)
	go func() {
		err = logS.Serve(logListener)
		if err != nil {
			return
		}
	}()
	// client
	/*
	 dispatcher client 调用 procsystem
	 dispatcher client 调用 filesystem
	 sshd client 调用 Cmdparser
	 sshd client 调用 Log
	*/
	connProc, err := grpc.NewClient(
		"127.0.0.1:9001",
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
		"127.0.0.1:9002",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Println("dispatcher connFs client start error")
	}
	defer connFs.Close()
	if connFs == nil {
		log.Fatal("Failed to connect to fs service")
	}

	connParser, err := grpc.NewClient(
		"127.0.0.1:9003",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Println("connParser client start error")
	}
	defer connParser.Close()
	if connParser == nil {
		log.Fatal("Failed to connect to sshd service")
	}

	connLog, err := grpc.NewClient(
		"127.0.0.1:9005",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Println("connLog client start error")
	}
	defer connLog.Close()
	if connLog == nil {
		log.Fatal("Failed to connect to sshd service")
	}

	clientsCmd := clientset.NewClientSetDispatcher(connProc, connFs)

	// dispatcher初始化
	dispatchSvc := dispatch_service.NewDispatcherServer(clientsCmd)
	// 命令注册
	// uname
	unameHandler := uname.NewUnameHandler(clientsCmd.ProcClient, clientsCmd.FSClient)
	dispatchSvc.RegisterCmd("uname", unameHandler)
	// ps
	psHandler := ps.NewPsHandler(clientsCmd.ProcClient, clientsCmd.FSClient)
	dispatchSvc.RegisterCmd("ps", psHandler)
	// ls
	lsHandler := ls.NewLsHandler(clientsCmd.ProcClient, clientsCmd.FSClient)
	dispatchSvc.RegisterCmd("ls", lsHandler)

	// dispatcher 启动
	dispatchGs := dispatch_transport.NewCmdDispatcherServer(dispatchSvc)
	dispatchListener, err := net.Listen("tcp", ":9004")
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
	// sshd client 调用 Dispatcher
	connDispatcher, err := grpc.NewClient(
		"127.0.0.1:9004",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Println("connDispatcher client start error")
	}
	defer connDispatcher.Close()
	if connDispatcher == nil {
		log.Fatal("Failed to connect to sshd service")
	}

	clientsLog := clientset.NewClientSetSSHD(connLog, connParser, connDispatcher)
	// sshd启动
	go func() {
		sshd_service.NewSSHDService(clientsCmd.ProcClient).StartSSHServer(clientsLog)

	}()
	// log http启动
	logEndpoints := log_endpoint.Endpoints{
		GetLogEndpoint:     log_endpoint.MakeGetLogEndpoint(logSvc),
		ReadAllLogEndpoint: log_endpoint.MakeReadAllLogEndpoint(logSvc),
	}
	httpHandler := log_httptransport.NewHTTPHandler(logEndpoints)
	go func() {
		log.Println("HTTP server is running...")
		if err := http.ListenAndServe(":8080", httpHandler); err != nil {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	select {}
}
