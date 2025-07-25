package endpoint

import (
	"context"
	"errors"
	"log"

	"github.com/BBitQNull/SSHoneyNet/core/commandparser"
	"github.com/go-kit/kit/endpoint"
)

type CmdParserRequest struct {
	Cmd       string
	SessionID string
}

type CmdParserResponse struct {
	Ast *commandparser.Script
}

func MakeCmdParserEndpoint(svc commandparser.CmdParserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(CmdParserRequest)
		if !ok {
			log.Printf("failed to assert:")
			return nil, errors.New("failed to assert")
		}
		ast, err := svc.CommandParser(ctx, req.Cmd, req.SessionID)
		if err != nil {
			log.Printf("failed to cmdparser: %v", err)
			return nil, err
		}
		return CmdParserResponse{Ast: ast}, nil
	}
}
