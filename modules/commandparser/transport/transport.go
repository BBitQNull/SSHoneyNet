package transport

import (
	"context"
	"errors"
	"log"

	"github.com/BBitQNull/SSHoneyNet/core/commandparser"
	"github.com/BBitQNull/SSHoneyNet/modules/commandparser/endpoint"
	pb "github.com/BBitQNull/SSHoneyNet/pb/cmdparser"
	"github.com/BBitQNull/SSHoneyNet/pkg/utils/convert"
	grpctransport "github.com/go-kit/kit/transport/grpc"
)

type grpcServer struct {
	pb.UnimplementedCmdParserServer
	cmdParser grpctransport.Handler
}

func decodeGRPCmdParserRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req, ok := grpcReq.(*pb.CmdParserRequest)
	if !ok {
		log.Printf("failed to assert")
		return nil, errors.New("failed to assert")
	}
	if req.Cmd == "" {
		return nil, errors.New("empty command not allowed")
	}
	return endpoint.CmdParserRequest{Cmd: req.Cmd}, nil
}

func encodeGRPCmdParserResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(endpoint.CmdParserResponse)
	if !ok {
		log.Printf("failed to assert")
		return nil, errors.New("failed to assert")
	}
	if resp.Ast == nil {
		return nil, errors.New("empty CmdParserResponse not allowed")
	}
	return &pb.CmdParserResponse{Ast: convert.ConvertScript(resp.Ast)}, nil
}

func NewGRPCmdParserServer(svc commandparser.CmdParserService) pb.CmdParserServer {
	return &grpcServer{
		cmdParser: grpctransport.NewServer(
			endpoint.MakeCmdParserEndpoint(svc),
			decodeGRPCmdParserRequest,
			encodeGRPCmdParserResponse,
		),
	}
}

func (s *grpcServer) CommandParser(ctx context.Context, req *pb.CmdParserRequest) (*pb.CmdParserResponse, error) {
	_, rep, err := s.cmdParser.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.CmdParserResponse), nil
}
