package sshd_client

import (
	"context"
	"errors"
	"fmt"

	"github.com/BBitQNull/SSHoneyNet/core/log"
	pb "github.com/BBitQNull/SSHoneyNet/pb/cmdparser"
	pbdis "github.com/BBitQNull/SSHoneyNet/pb/dispatcher"
	log_Pb "github.com/BBitQNull/SSHoneyNet/pb/log"
	"github.com/BBitQNull/SSHoneyNet/pkg/utils/convert"
	"github.com/go-kit/kit/endpoint"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"
)

type RawCmdParserRequest struct {
	Cmd string
}

type RawCmdParserResponse struct {
	Result  string
	ErrCode int32
	ErrMsg  string
}

type RawWriteLogRequest struct {
	LogEntry log.LogEntry
}

type RawWriteLogResponse struct{}

type SSHDManageClient struct {
	WriteLog      endpoint.Endpoint
	CommandParser endpoint.Endpoint
	Dispatcher    endpoint.Endpoint
}

func NewSSHDManageClient(connLog, connParser, connDispatcher *grpc.ClientConn) *SSHDManageClient {
	return &SSHDManageClient{
		WriteLog:      MakeWriteLogEndpoint(connLog),
		CommandParser: MakeCmdParserEndpoint(connParser),
		Dispatcher:    MakeCmdDispatchEndpoint(connDispatcher),
	}
}

func MakeCmdParserEndpoint(conn *grpc.ClientConn) endpoint.Endpoint {
	return grpctransport.NewClient(
		conn,
		"pb.CmdParser",
		"CommandParser",
		encodeCmdParserRequest,
		decodeCmdParserResponse,
		pb.CmdParserResponse{},
	).Endpoint()
}

func MakeCmdDispatchEndpoint(conn *grpc.ClientConn) endpoint.Endpoint {
	return grpctransport.NewClient(
		conn,
		"pb.CmdEcho",
		"Dispatcher",
		encodeCmdDispatchRequest,
		decodeCmdDispatchResponse,
		pbdis.DispatcherResponse{},
	).Endpoint()
}

func MakeWriteLogEndpoint(conn *grpc.ClientConn) endpoint.Endpoint {
	return grpctransport.NewClient(
		conn,
		"pb.LogService",
		"WriteLog",
		encodeWriteLogRequest,
		decodeWriteLogResponse,
		log_Pb.WriteLogResponse{},
	).Endpoint()
}

func encodeCmdParserRequest(ctx context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*RawCmdParserRequest)
	if !ok {
		return nil, fmt.Errorf("encodeCmdParserRequest: expected *RawRequest but got %T", request)
	}
	return &pb.CmdParserRequest{
		Cmd: req.Cmd,
	}, nil
}

func decodeCmdParserResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(*pb.CmdParserResponse)
	if !ok {
		return nil, errors.New("decodeCmdParserResponse error")
	}
	if resp.Ast == nil {
		return nil, errors.New("decodeCmdParserResponse error")
	}
	return &pb.CmdParserResponse{Ast: resp.Ast}, nil
}

func encodeCmdDispatchRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*pbdis.DispatcherRequest)
	if !ok {
		return nil, errors.New("encodeCmdDispatchRequest error")
	}
	return req, nil
}

func decodeCmdDispatchResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(*pbdis.DispatcherResponse)
	if !ok {
		return nil, fmt.Errorf("decodeCmdDispatchResponse: expected *RawRequest but got %T", response)
	}
	return RawCmdParserResponse{
		Result:  resp.Cmdresult,
		ErrCode: resp.Errcode,
		ErrMsg:  resp.Errmsg,
	}, nil
}

func encodeWriteLogRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*RawWriteLogRequest)
	if !ok {
		return nil, fmt.Errorf("encodeWriteLogRequest: expected *RawRequest but got %T", request)
	}
	return &log_Pb.WriteLogRequest{Entry: convert.ConvertLogEntryToPb(req.LogEntry)}, nil
}

func decodeWriteLogResponse(_ context.Context, response interface{}) (interface{}, error) {
	_, ok := response.(*log_Pb.WriteLogResponse)
	if !ok {
		return nil, fmt.Errorf("decodeWriteLogResponse: expected *RawRequest but got %T", response)
	}
	return &RawWriteLogResponse{}, nil
}
