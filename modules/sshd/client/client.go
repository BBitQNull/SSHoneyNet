package sshd_client

import (
	"context"
	"errors"

	pb "github.com/BBitQNull/SSHoneyNet/pb/cmdparser"
	pbdis "github.com/BBitQNull/SSHoneyNet/pb/dispatcher"
	"github.com/go-kit/kit/endpoint"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"
)

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

type RawRequest struct {
	Cmd string
}

func encodeCmdParserRequest(ctx context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*pb.CmdParserRequest)
	if !ok {
		return nil, errors.New("encodeCmdParserRequest error")
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

type RawResponse struct {
	Result  string
	ErrCode int32
	ErrMsg  string
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
		return nil, errors.New("decodeCmdDispatchResponse error")
	}
	return resp, nil
}
