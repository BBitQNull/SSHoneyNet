package log_httptransport

import (
	"context"
	"encoding/json"
	"net/http"

	log_endpoint "github.com/BBitQNull/SSHoneyNet/modules/log/endpoint"
)

func decodeGetSinceLogRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request log_endpoint.GetSinceLogRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func encodeGetSinceLogResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(response)
}
