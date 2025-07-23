package transport

import (
	"context"
	"errors"
	"log"

	"github.com/BBitQNull/SSHoneyNet/modules/commandparser/endpoint"
	pb "github.com/BBitQNull/SSHoneyNet/pb/cmdparser"
	"github.com/BBitQNull/SSHoneyNet/pkg/utils/convert"
	grpctransport "github.com/go-kit/kit/transport/grpc"
)

type grpcServer struct {
	pb.UnimplementedCmdParserServer
	CmdParser grpctransport.Handler
}

func decodeGRPCCmdParserRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req, ok := grpcReq.(*pb.CmdParserRequest)
	if !ok {
		log.Fatal("failed to assert")
		return nil, errors.New("failed to assert")
	}
	return endpoint.CmdParserRequest{Cmd: req.Cmd}, nil
}

func encodeGRPCCmdParserResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(endpoint.CmdParserResponse)
	if !ok {
		log.Fatal("failed to assert")
		return nil, errors.New("failed to assert")
	}
	return &pb.CmdParserResponse{Ast: convert.ConvertScript(resp.Ast)}, nil
}
