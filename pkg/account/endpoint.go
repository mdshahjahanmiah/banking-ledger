package account

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	eError "github.com/mdshahjahanmiah/explore-go/error"
	"github.com/shopspring/decimal"
	"net/http"
	"strings"
	"time"
)

type AccountResponse struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	Balance   string `json:"balance"` // Decimal as string
	Currency  string `json:"currency"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"` // ISO8601 timestamp
	UpdatedAt string `json:"updated_at"` // ISO8601 timestamp
}

func makePostAccountEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		// Type assertion with validation
		req, ok := request.(AccountRequest)
		if !ok {
			return nil, eError.NewServiceError(
				nil,
				"Invalid request type",
				"invalid_request_type",
				http.StatusBadRequest,
			)
		}

		createReq := CreateAccountRequest{
			UserID:   req.UserID,
			Currency: strings.ToUpper(req.Currency), // Ensure uppercase
			Balance:  decimal.NewFromFloat(req.InitialBalance),
		}

		account, err := s.CreateAccount(ctx, createReq)
		if err != nil {
			return nil, err
		}

		return AccountResponse{
			ID:        account.ID,
			UserID:    account.UserID,
			Balance:   account.Balance.String(), // Return as string for precision
			Currency:  account.Currency,
			Status:    string(account.Status),
			CreatedAt: account.CreatedAt.Format(time.RFC3339),
			UpdatedAt: account.UpdatedAt.Format(time.RFC3339),
		}, nil
	}
}
