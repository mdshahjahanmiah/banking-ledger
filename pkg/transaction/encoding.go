package transaction

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	eError "github.com/mdshahjahanmiah/explore-go/error"
	"github.com/pkg/errors"
	"log/slog"
	"net/http"
)

type TransactionRequest struct {
	AccountID   string  `json:"account_id"` // Required
	Amount      float64 `json:"amount"`     // Must be > 0
	Currency    string  `json:"currency"`   // Must match account currency
	ReferenceID string  `json:"reference_id"`
}

type AuditRequest struct {
	AccountID string
}

func decodeDepositRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var req TransactionRequest
	if err := decoder.Decode(&req); err != nil {
		slog.Error("decode deposit request", "err", err)
		return nil, eError.NewServiceError(err, "decode deposit request", "payload", http.StatusBadRequest)
	}

	if req.Amount <= 0 {
		return nil, errors.New("amount must be positive")
	}

	if req.Currency == "" {
		return nil, errors.New("currency is required")
	}

	return req, nil
}

func decodeWithdrawRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var req TransactionRequest
	if err := decoder.Decode(&req); err != nil {
		slog.Error("decode withdraw request", "err", err)
		return nil, eError.NewServiceError(err, "decode withdraw request", "payload", http.StatusBadRequest)
	}

	if req.AccountID == "" || req.Amount <= 0 || req.Currency == "" {
		slog.Warn("invalid withdraw request", "request", req)
		return nil, eError.NewServiceError(nil, "invalid withdraw data", "validation", http.StatusBadRequest)
	}

	return req, nil
}

func decodeAuditRequest(_ context.Context, r *http.Request) (interface{}, error) {
	accountID := chi.URLParam(r, "id")
	if accountID == "" {
		return nil, eError.NewServiceError(nil, "missing account_id in path", "validation", http.StatusBadRequest)
	}

	return AuditRequest{AccountID: accountID}, nil
}
