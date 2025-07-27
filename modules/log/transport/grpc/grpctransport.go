package log_grpctransport

import (
	"context"
	"fmt"

	"github.com/BBitQNull/SSHoneyNet/core/log"
	log_endpoint "github.com/BBitQNull/SSHoneyNet/modules/log/endpoint"
	log_Pb "github.com/BBitQNull/SSHoneyNet/pb/log"
	"github.com/BBitQNull/SSHoneyNet/pkg/utils/convert"
	grpctransport "github.com/go-kit/kit/transport/grpc"
)

type grpcServer struct {
	log_Pb.UnimplementedLogServiceServer
	writelog grpctransport.Handler
}

func decodeWriteLogRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req, ok := grpcReq.(*log_Pb.WriteLogRequest)
	if !ok {
		return nil, fmt.Errorf("assert error in decodeWriteLogRequest, got %T", grpcReq)
	}
	return log_endpoint.WriteLogRequest{LogEntry: convert.ConvertLogEntryFormPb(req.Entry)}, nil
}

func encodeWriteLogResponse(_ context.Context, response interface{}) (interface{}, error) {
	_, ok := response.(log_endpoint.WriteLogResponse)
	if !ok {
		return nil, fmt.Errorf("assert error in encodeWriteLogResponse")
	}
	return &log_Pb.WriteLogResponse{}, nil
}

func NewWriteLogServer(svc log.LogService) log_Pb.LogServiceServer {
	return &grpcServer{
		writelog: grpctransport.NewServer(
			log_endpoint.MakeWriteLogEndpoint(svc),
			decodeWriteLogRequest,
			encodeWriteLogResponse,
		),
	}
}

func (s *grpcServer) WriteLog(ctx context.Context, req *log_Pb.WriteLogRequest) (*log_Pb.WriteLogResponse, error) {
	_, rep, err := s.writelog.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	resp, ok := rep.(*log_Pb.WriteLogResponse)
	if !ok {
		return nil, fmt.Errorf("unexpected response type %T", rep)
	}
	return resp, nil
}
