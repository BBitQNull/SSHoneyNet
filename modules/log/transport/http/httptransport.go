package log_httptransport

import (
	"context"
	"encoding/json"
	"net/http"

	log_endpoint "github.com/BBitQNull/SSHoneyNet/modules/log/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
)

func decodeGetSinceLogRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request log_endpoint.GetSinceLogRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func encodeLogResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(response)
}

func decodeReadAllLogRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request log_endpoint.ReadAllLogRequest
	return request, nil
}

func NewHTTPHandler(endpoints log_endpoint.Endpoints) http.Handler {
	mux := http.NewServeMux()

	// 绑定不同的路径到不同的 handler
	mux.Handle("/logs/since", httptransport.NewServer(
		endpoints.GetLogEndpoint,
		decodeGetSinceLogRequest,
		encodeLogResponse,
	))

	mux.Handle("/logs/all", httptransport.NewServer(
		endpoints.ReadAllLogEndpoint,
		decodeReadAllLogRequest,
		encodeLogResponse,
	))

	return mux
}
