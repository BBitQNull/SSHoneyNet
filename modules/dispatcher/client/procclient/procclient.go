package proc_client

import (
	"context"
	"fmt"

	"github.com/BBitQNull/SSHoneyNet/core/procsystem"
	proc_Pb "github.com/BBitQNull/SSHoneyNet/pb/procsystem"
	"github.com/BBitQNull/SSHoneyNet/pkg/utils/convert"
	"github.com/go-kit/kit/endpoint"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"
)

type RawRequest struct {
	Pid     int64
	Ppid    int64
	Command string
}

type RawProcResponse struct {
	Process *procsystem.PCB
}

type ListProcResponse struct {
	Processes []*procsystem.PCB
}

type ProcManageClient struct {
	CreateProc endpoint.Endpoint
	KillProc   endpoint.Endpoint
	GetProc    endpoint.Endpoint
	ListProc   endpoint.Endpoint
}

func NewProcManageClient(conn *grpc.ClientConn) ProcManageClient {
	return ProcManageClient{
		CreateProc: MakeProcessCreateEndpoint(conn),
		KillProc:   MakeProcessKillEndpoint(conn),
		GetProc:    MakeProcessGetEndpoint(conn),
		ListProc:   MakeProcessListEndpoint(conn),
	}
}

func MakeProcessCreateEndpoint(conn *grpc.ClientConn) endpoint.Endpoint {
	return grpctransport.NewClient(
		conn,
		"pb.ProcManage",
		"CreateProc",
		encodeProcCreateRequest,
		decodeProcCreateResponse,
		proc_Pb.ProcResponse{},
	).Endpoint()
}

func encodeProcCreateRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*RawRequest)
	if !ok {
		return nil, fmt.Errorf("encodeProcCreateRequest: expected *RawRequest but got %T", request)
	}
	return &proc_Pb.ProcRequest{
		Command: req.Command,
		Pid:     req.Pid,
		Ppid:    req.Ppid,
	}, nil
}

func decodeProcCreateResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(*proc_Pb.ProcResponse)
	if !ok {
		return nil, fmt.Errorf("decodeProcCreateResponse: expected *ProcResponse but got %T", response)
	}
	return &RawProcResponse{
		Process: convert.ConvertPcbFromPb(resp.Pcb),
	}, nil
}

func MakeProcessKillEndpoint(conn *grpc.ClientConn) endpoint.Endpoint {
	return grpctransport.NewClient(
		conn,
		"pb.ProcManage",
		"KillProc",
		encodeProcKillRequest,
		decodeProcKillResponse,
		proc_Pb.ProcResponse{},
	).Endpoint()
}

func encodeProcKillRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*RawRequest)
	if !ok {
		return nil, fmt.Errorf("encodeProcKillRequest: expected *RawRequest but got %T", request)
	}
	return &proc_Pb.ProcRequest{
		Command: req.Command,
		Pid:     req.Pid,
		Ppid:    req.Ppid,
	}, nil
}

func decodeProcKillResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(*proc_Pb.ProcResponse)
	if !ok {
		return nil, fmt.Errorf("decodeProcKillResponse: expected *ProcResponse but got %T", response)
	}
	return &RawProcResponse{
		Process: convert.ConvertPcbFromPb(resp.Pcb),
	}, nil
}

func MakeProcessGetEndpoint(conn *grpc.ClientConn) endpoint.Endpoint {
	return grpctransport.NewClient(
		conn,
		"pb.ProcManage",
		"GetProc",
		encodeProcGetRequest,
		decodeProcGetResponse,
		proc_Pb.ProcResponse{},
	).Endpoint()
}

func encodeProcGetRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*RawRequest)
	if !ok {
		return nil, fmt.Errorf("encodeProcGetRequest: expected *RawRequest but got %T", request)
	}
	return &proc_Pb.ProcRequest{
		Command: req.Command,
		Pid:     req.Pid,
		Ppid:    req.Ppid,
	}, nil
}

func decodeProcGetResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(*proc_Pb.ProcResponse)
	if !ok {
		return nil, fmt.Errorf("decodeProcGetResponse: expected *ProcResponse but got %T", response)
	}
	return &RawProcResponse{
		Process: convert.ConvertPcbFromPb(resp.Pcb),
	}, nil
}

func MakeProcessListEndpoint(conn *grpc.ClientConn) endpoint.Endpoint {
	return grpctransport.NewClient(
		conn,
		"pb.ProcManage",
		"ListProc",
		encodeProcListRequest,
		decodeProcListResponse,
		proc_Pb.ProcResponse{},
	).Endpoint()
}

func encodeProcListRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*RawRequest)
	if !ok {
		return nil, fmt.Errorf("encodeProcListRequest: expected *RawRequest but got %T", request)
	}
	return &proc_Pb.ProcRequest{
		Command: req.Command,
		Pid:     req.Pid,
		Ppid:    req.Ppid,
	}, nil
}

func decodeProcListResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(*proc_Pb.ProcResponse)
	if !ok {
		return nil, fmt.Errorf("decodeProcListResponse: expected *ProcResponse but got %T", response)
	}
	return &ListProcResponse{
		Processes: convert.ConvertPcbListFromPb(resp.Pcblist),
	}, nil
}
