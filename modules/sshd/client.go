package sshd

import (
	"github.com/go-kit/kit/endpoint"
	"google.golang.org/grpc"
)

func MakeCmdParserEndpoint(conn *grpc.ClientConn) endpoint.Endpoint {
	return nil
}
