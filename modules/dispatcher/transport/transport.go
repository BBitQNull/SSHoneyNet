package transport

import (
	"github.com/BBitQNull/SSHoneyNet/core/dispatcher"
	"github.com/BBitQNull/SSHoneyNet/modules/dispatcher/endpoint"
	pb "github.com/BBitQNull/SSHoneyNet/pb/dispatcher"
	"github.com/BBitQNull/SSHoneyNet/pkg/utils/convert"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"golang.org/x/net/context"
)

type grpcServer struct {
	pb.UnimplementedCmdEchoServer
	cmdecho grpctransport.Handler
}

func decodeGRPCCmdDispatcher(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.DispatcherRequest)
	return endpoint.CmdDispatchRequest{Ast: *convert.ConvertScriptFormpb(req.Ast)}, nil
}

func encodeGRPCCmdDispatcher(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(endpoint.CmdDispatchResponse)
	return &pb.DispatcherResponse{
		Cmdresult: resp.CmdResult.CmdResult,
		Errcode:   resp.CmdResult.ErrCode,
		Errmsg:    resp.CmdResult.ErrMsg,
	}, nil
}

func NewCmdDispatcherServer(svc dispatcher.CmdDispatcherService) pb.CmdEchoServer {
	return &grpcServer{
		cmdecho: grpctransport.NewServer(
			endpoint.MakeCmdDispatch(svc),
			decodeGRPCCmdDispatcher,
			encodeGRPCCmdDispatcher,
		),
	}
}

func (s *grpcServer) CmdEcho(ctx context.Context, req *pb.DispatcherRequest) (*pb.DispatcherResponse, error) {
	_, rep, err := s.cmdecho.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.DispatcherResponse), nil
}
