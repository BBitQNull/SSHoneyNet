package endpoint

import (
	"context"
	"errors"
	"log"

	"github.com/BBitQNull/SSHoneyNet/core/sshd"
	"github.com/BBitQNull/SSHoneyNet/pkg/model"
	"github.com/go-kit/kit/endpoint"
)

type EchoCommandRequest struct {
	Result model.CmdResult
}

type EchoCommandResponse struct {
	IsErr bool
}

func MakeSSHDEndpoint(svc sshd.SSHDService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(EchoCommandRequest)
		if !ok {
			log.Fatal("failed to assert")
			return EchoCommandResponse{IsErr: false}, errors.New("failed to assert")
		}
		v, err := svc.EchoCommand(ctx, req.Result)
		if err != nil {
			log.Fatal("EchoCommand failed")
			return EchoCommandResponse{IsErr: v}, errors.New("failed to echocmd")
		}
		return EchoCommandResponse{IsErr: v}, nil
	}
}
