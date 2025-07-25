package proc_endpoint

import (
	"context"
	"errors"
	"log"

	"github.com/BBitQNull/SSHoneyNet/core/procsystem"
	"github.com/go-kit/kit/endpoint"
)

type ProcessRequest struct {
	Pid     int64
	Command string
	Ppid    int64
}

type ProcessResponse struct {
	PCB     *procsystem.PCB
	PCBList []*procsystem.PCB
}

func MakeProcessCreateEndpoint(svc procsystem.ProcessManager) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(ProcessRequest)
		if !ok {
			log.Println("assert error MakeProcessEndpoint")
			return nil, errors.New("assert error MakeProcessEndpoint")
		}
		v, err := svc.CreateProcess(ctx, req.Command, req.Pid, req.Ppid)
		if err != nil {
			log.Println("CreateProcess error")
			return nil, errors.New("CreateProcess error")
		}
		return ProcessResponse{
			PCB:     v,
			PCBList: nil,
		}, nil
	}
}

func MakeProcessKillEndpoint(svc procsystem.ProcessManager) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(ProcessRequest)
		if !ok {
			log.Println("assert error MakeProcessKillEndpointt")
			return nil, errors.New("assert error MakeProcessKillEndpoint")
		}
		err := svc.KillProcess(ctx, req.Pid)
		if err != nil {
			log.Println("KillProcess error")
			return nil, errors.New("KillProcess error")
		}
		return ProcessResponse{
			PCB:     nil,
			PCBList: nil,
		}, nil
	}
}

func MakeProcessGetEndpoint(svc procsystem.ProcessManager) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(ProcessRequest)
		if !ok {
			log.Println("assert error MakeProcessGetEndpoint")
			return nil, errors.New("assert error MakeProcessGetEndpoint")
		}
		v, err := svc.GetProcess(ctx, req.Pid)
		if err != nil {
			log.Println("GetProcess error")
			return nil, errors.New("GetProcess error")
		}
		return ProcessResponse{
			PCB:     v,
			PCBList: nil,
		}, nil
	}
}

func MakeProcessListEndpoint(svc procsystem.ProcessManager) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		_, ok := request.(ProcessRequest)
		if !ok {
			log.Println("assert error MakeProcessListEndpoint")
			return nil, errors.New("assert error MakeProcessListEndpoint")
		}
		v, err := svc.ListProcess(ctx)
		if err != nil {
			log.Println("ListProcess error")
			return nil, errors.New("ListProcess error")
		}
		return ProcessResponse{
			PCB:     nil,
			PCBList: v,
		}, nil
	}
}
