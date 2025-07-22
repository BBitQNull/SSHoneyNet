package sshd_transport

import (
	"context"
	"errors"

	"github.com/BBitQNull/SSHoneyNet/core/sshd"
	"github.com/BBitQNull/SSHoneyNet/modules/sshd/endpoint"
	"github.com/BBitQNull/SSHoneyNet/pb"
	"github.com/BBitQNull/SSHoneyNet/pkg/model"
	grpctransport "github.com/go-kit/kit/transport/grpc"
)

type grpcSSHDService struct {
	pb.UnimplementedEchoCmdServer
	echoCommand grpctransport.Handler
}

func decodeGRPCEchoCmdRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	if req, ok := grpcReq.(*pb.EchoCmdRequest); ok {
		return endpoint.EchoCommandRequest{Result: model.CmdResult{
			Output:     req.Result.Output,
			ExitCode:   req.Result.Exitcode,
			ErrMsg:     errors.New(req.Result.Errmsg),
			Log:        req.Result.Log,
			NextPrompt: req.Result.Nextprompt,
		}}, nil
	}
	return nil, errors.New("decodeGRPCEchoCmdRequest failed")
}

func encodeGRPCEchoCmdResponse(_ context.Context, response interface{}) (interface{}, error) {
	if resp, ok := response.(endpoint.EchoCommandResponse); ok {
		return &pb.EchoCmdResponse{Iserr: resp.IsErr}, nil
	}
	return nil, errors.New("encodeGRPCEchoCmdResponse failed")
}

func NewGRPCSSHDServer(svc sshd.SSHDService) pb.EchoCmdServer {
	return &grpcSSHDService{
		echoCommand: grpctransport.NewServer(
			endpoint.MakeSSHDEndpoint(svc),
			decodeGRPCEchoCmdRequest,
			encodeGRPCEchoCmdResponse,
		),
	}
}

func (s *grpcSSHDService) EchoCommand(ctx context.Context, req *pb.EchoCmdRequest) (*pb.EchoCmdResponse, error) {
	_, rep, err := s.echoCommand.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.EchoCmdResponse), nil
}
