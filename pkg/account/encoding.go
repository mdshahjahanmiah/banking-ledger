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

func decodeCreateAccountRequest(ctx context.Context, request *http.Request) (interface{}, error) {
	decoder := json.NewDecoder(request.Body)
	decoder.DisallowUnknownFields()

	var accountRequest AccountRequest
	err := decoder.Decode(&accountRequest)
	if err != nil {
		slog.Error("decode account create request", "err", err)
		return nil, eError.NewServiceError(err, "decode account create request", "payload", http.StatusBadRequest)
	}

	if accountRequest.UserID == "" || accountRequest.InitialBalance < 0 || accountRequest.Currency == "" {
		slog.Warn("invalid account create request", "request", accountRequest)
		return nil, eError.NewServiceError(nil, "invalid account data", "validation", http.StatusBadRequest)
	}

	return accountRequest, nil
}
