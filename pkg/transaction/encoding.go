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
	AccountID   string  `json:"account_id"`
	Amount      float64 `json:"amount"`
	Currency    string  `json:"currency"`
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
		return nil, eError.NewServiceError(err, "invalid request payload", "INVALID_PAYLOAD", http.StatusBadRequest)
	}

	if req.Amount <= 0 {
		return nil, eError.NewServiceError(
			errors.New("amount must be positive"), "amount must be greater than zero", "INVALID_AMOUNT", http.StatusBadRequest)
	}

	if req.Currency == "" {
		return nil, eError.NewServiceError(
			errors.New("currency is required"), "currency is required", "MISSING_CURRENCY", http.StatusBadRequest)
	}

	if req.AccountID == "" {
		return nil, eError.NewServiceError(
			errors.New("account_id is required"), "account_id is required", "MISSING_ACCOUNT_ID", http.StatusBadRequest)
	}

	return req, nil
}

func decodeWithdrawRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var req TransactionRequest
	if err := decoder.Decode(&req); err != nil {
		slog.Error("decode withdraw request", "err", err)
		return nil, eError.NewServiceError(err, "invalid request payload", "INVALID_PAYLOAD", http.StatusBadRequest)
	}

	if req.AccountID == "" {
		return nil, eError.NewServiceError(
			errors.New("account_id is required"), "account_id is required", "MISSING_ACCOUNT_ID", http.StatusBadRequest)
	}

	if req.Amount <= 0 {
		return nil, eError.NewServiceError(
			errors.New("amount must be positive"), "amount must be greater than zero", "INVALID_AMOUNT", http.StatusBadRequest)
	}

	if req.Currency == "" {
		return nil, eError.NewServiceError(
			errors.New("currency is required"), "currency is required", "MISSING_CURRENCY", http.StatusBadRequest)
	}

	return req, nil
}

func decodeAuditRequest(_ context.Context, r *http.Request) (interface{}, error) {
	accountID := chi.URLParam(r, "id")
	if accountID == "" {
		return nil, eError.NewServiceError(
			errors.New("account_id missing in path"), "missing account_id in path", "MISSING_ACCOUNT_ID", http.StatusBadRequest)
	}

	return AuditRequest{AccountID: accountID}, nil
}
