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
	Balance   string `json:"balance"`
	Currency  string `json:"currency"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func makePostAccountEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
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
			Balance:   account.Balance.String(),
			Currency:  account.Currency,
			Status:    string(account.Status),
			CreatedAt: account.CreatedAt.Format(time.RFC3339),
			UpdatedAt: account.UpdatedAt.Format(time.RFC3339),
		}, nil
	}
}
