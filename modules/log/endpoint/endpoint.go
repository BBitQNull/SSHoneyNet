package log_endpoint

import (
	"context"
	"fmt"
	"time"

	"github.com/BBitQNull/SSHoneyNet/core/log"
	"github.com/go-kit/kit/endpoint"
)

type WriteLogRequest struct {
	LogEntry log.LogEntry
}

type WriteLogResponse struct{}

type GetSinceLogRequest struct {
	Timestamp string `json:"timestamp"`
}

type GetSinceLogResponse struct {
	LogOutput []log.LogEntry
}

func MakeGetLogEndpoint(svc log.LogService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(GetSinceLogRequest)
		if !ok {
			return nil, fmt.Errorf("invalid request type in MakeGetLogEndpoint, got %T", request)
		}
		t, err := time.Parse(time.RFC3339, req.Timestamp)
		if err != nil {
			return nil, fmt.Errorf("invalid timestamp: %v", err)
		}
		v, err := svc.GetLogSince(ctx, t)
		if err != nil {
			return nil, err
		}
		return GetSinceLogResponse{LogOutput: v}, nil
	}
}

func MakeWriteLogEndpoint(svc log.LogService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(WriteLogRequest)
		if !ok {
			return nil, fmt.Errorf("invalid request type in MakeWriteLogEndpoint, got %T", request)
		}
		err = svc.WriteLog(ctx, req.LogEntry)
		if err != nil {
			return nil, err
		}
		return WriteLogResponse{}, nil
	}
}
