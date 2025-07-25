package proc_transport

import (
	"context"
	"errors"
	"log"

	"github.com/BBitQNull/SSHoneyNet/core/procsystem"
	proc_endpoint "github.com/BBitQNull/SSHoneyNet/modules/procsystem/endpoint"
	proc_pb "github.com/BBitQNull/SSHoneyNet/pb/procsystem"
	"github.com/BBitQNull/SSHoneyNet/pkg/utils/convert"
	grpctransport "github.com/go-kit/kit/transport/grpc"
)

type grpcServer struct {
	proc_pb.UnimplementedProcManageServer
	createproc grpctransport.Handler
	killproc   grpctransport.Handler
	getproc    grpctransport.Handler
	listproc   grpctransport.Handler
}

func decodeCreateProcRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req, ok := grpcReq.(*proc_pb.ProcRequest)
	if !ok {
		log.Println("assert error decodeCreateProcRequest")
		return nil, errors.New("assert error decodeCreateProcRequest")
	}
	return proc_endpoint.ProcessRequest{
		Pid:     req.Pid,
		Ppid:    req.Ppid,
		Command: req.Command,
	}, nil
}

func decodeKillProcRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req, ok := grpcReq.(*proc_pb.ProcRequest)
	if !ok {
		log.Println("assert error decodeKillProcRequest")
		return nil, errors.New("assert error decodeKillProcRequest")
	}
	return proc_endpoint.ProcessRequest{
		Pid:     req.Pid,
		Ppid:    req.Ppid,
		Command: req.Command,
	}, nil
}

func decodeGetProcRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req, ok := grpcReq.(*proc_pb.ProcRequest)
	if !ok {
		log.Println("assert error decodeGetProcRequest")
		return nil, errors.New("assert error decodeGetProcRequest")
	}
	return proc_endpoint.ProcessRequest{
		Pid:     req.Pid,
		Ppid:    req.Ppid,
		Command: req.Command,
	}, nil
}

func decodeListProcRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req, ok := grpcReq.(*proc_pb.ProcRequest)
	if !ok {
		log.Println("assert error decodeListProcRequest")
		return nil, errors.New("assert error decodeListProcRequest")
	}
	return proc_endpoint.ProcessRequest{
		Pid:     req.Pid,
		Ppid:    req.Ppid,
		Command: req.Command,
	}, nil
}

func encodeCreateProcResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(proc_endpoint.ProcessResponse)
	if !ok {
		log.Println("assert error encodeCreateProcResponse")
		return nil, errors.New("assert error encodeCreateProcResponse")
	}
	return &proc_pb.ProcResponse{
		Pcb:     convert.ConvertPcbFromEndpoint(resp.PCB),
		Pcblist: convert.ConvertPcbListFromEndpoint(resp.PCBList),
	}, nil
}

func encodeKillProcResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(proc_endpoint.ProcessResponse)
	if !ok {
		log.Println("assert error encodeKillProcResponse")
		return nil, errors.New("assert error encodeKillProcResponse")
	}
	return &proc_pb.ProcResponse{
		Pcb:     convert.ConvertPcbFromEndpoint(resp.PCB),
		Pcblist: convert.ConvertPcbListFromEndpoint(resp.PCBList),
	}, nil
}

func encodeGetProcResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(proc_endpoint.ProcessResponse)
	if !ok {
		log.Println("assert error encodeGetProcResponse")
		return nil, errors.New("assert error encodeGetProcResponse")
	}
	return &proc_pb.ProcResponse{
		Pcb:     convert.ConvertPcbFromEndpoint(resp.PCB),
		Pcblist: convert.ConvertPcbListFromEndpoint(resp.PCBList),
	}, nil
}

func encodeListProcResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(proc_endpoint.ProcessResponse)
	if !ok {
		log.Println("assert error encodeListProcResponse")
		return nil, errors.New("assert error encodeListProcResponse")
	}
	return &proc_pb.ProcResponse{
		Pcb:     convert.ConvertPcbFromEndpoint(resp.PCB),
		Pcblist: convert.ConvertPcbListFromEndpoint(resp.PCBList),
	}, nil
}

func NewGRPCProcServer(svc procsystem.ProcessManager) proc_pb.ProcManageServer {
	return &grpcServer{
		createproc: grpctransport.NewServer(
			proc_endpoint.MakeProcessCreateEndpoint(svc),
			decodeCreateProcRequest,
			encodeCreateProcResponse,
		),
		getproc: grpctransport.NewServer(
			proc_endpoint.MakeProcessGetEndpoint(svc),
			decodeGetProcRequest,
			encodeGetProcResponse,
		),
		killproc: grpctransport.NewServer(
			proc_endpoint.MakeProcessKillEndpoint(svc),
			decodeKillProcRequest,
			encodeKillProcResponse,
		),
		listproc: grpctransport.NewServer(
			proc_endpoint.MakeProcessListEndpoint(svc),
			decodeListProcRequest,
			encodeListProcResponse,
		),
	}
}

func (s *grpcServer) CreateProc(ctx context.Context, req *proc_pb.ProcRequest) (*proc_pb.ProcResponse, error) {
	_, rep, err := s.createproc.ServeGRPC(ctx, req)
	if err != nil {
		log.Println("error CreateProc")
		return nil, errors.New("error CreateProc")
	}
	return rep.(*proc_pb.ProcResponse), nil
}

func (s *grpcServer) KillProc(ctx context.Context, req *proc_pb.ProcRequest) (*proc_pb.ProcResponse, error) {
	_, rep, err := s.killproc.ServeGRPC(ctx, req)
	if err != nil {
		log.Println("error KillProc")
		return nil, errors.New("error KillProc")
	}
	return rep.(*proc_pb.ProcResponse), nil
}

func (s *grpcServer) GetProc(ctx context.Context, req *proc_pb.ProcRequest) (*proc_pb.ProcResponse, error) {
	_, rep, err := s.getproc.ServeGRPC(ctx, req)
	if err != nil {
		log.Println("error GetProc")
		return nil, errors.New("error GetProc")
	}
	return rep.(*proc_pb.ProcResponse), nil
}

func (s *grpcServer) ListProc(ctx context.Context, req *proc_pb.ProcRequest) (*proc_pb.ProcResponse, error) {
	_, rep, err := s.listproc.ServeGRPC(ctx, req)
	if err != nil {
		log.Println("error ListProc")
		return nil, errors.New("error ListProc")
	}
	return rep.(*proc_pb.ProcResponse), nil
}
