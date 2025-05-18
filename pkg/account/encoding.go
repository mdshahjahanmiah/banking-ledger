package account

import (
	"context"
	"encoding/json"
	eError "github.com/mdshahjahanmiah/explore-go/error"
	"log/slog"
	"net/http"
)

type AccountRequest struct {
	UserID         string  `json:"user_id"`
	InitialBalance float64 `json:"initial_balance"`
	Currency       string  `json:"currency"`
}

func decodeCreateAccountRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var req AccountRequest
	if err := decoder.Decode(&req); err != nil {
		slog.Error("failed to decode create account request", "error", err)
		return nil, eError.NewServiceError(err, "failed to decode create account request", "payload", http.StatusBadRequest)
	}

	if req.UserID == "" {
		slog.Warn("user_id is required", "request", req)
		return nil, eError.NewServiceError(nil, "user_id is required", "validation", http.StatusBadRequest)
	}

	if req.Currency == "" {
		slog.Warn("currency is required", "request", req)
		return nil, eError.NewServiceError(nil, "currency is required", "validation", http.StatusBadRequest)
	}

	if req.InitialBalance < 0 {
		slog.Warn("initial_balance must be >= 0", "request", req)
		return nil, eError.NewServiceError(nil, "initial_balance must be >= 0", "validation", http.StatusBadRequest)
	}

	return req, nil
}
