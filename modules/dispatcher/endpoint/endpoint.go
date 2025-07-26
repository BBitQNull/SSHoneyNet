package endpoint

import (
	"context"
	"errors"

	"github.com/BBitQNull/SSHoneyNet/core/commandparser"
	"github.com/BBitQNull/SSHoneyNet/core/dispatcher"
	"github.com/go-kit/kit/endpoint"
)

type CmdDispatchRequest struct {
	Ast       commandparser.Script
	SessionID string
}

type CmdDispatchResponse struct {
	CmdResult dispatcher.CmdEcho
}

func MakeCmdDispatch(svc dispatcher.CmdDispatcherService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(CmdDispatchRequest)
		if !ok {
			return nil, errors.New("error")
		}
		result, err := svc.CmdDispatcher(ctx, req.Ast, req.SessionID)
		if err != nil {
			return nil, err
		}
		return CmdDispatchResponse{CmdResult: result}, nil
	}
}
