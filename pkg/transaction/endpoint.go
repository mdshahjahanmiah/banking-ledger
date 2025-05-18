package transaction

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/mdshahjahanmiah/banking-ledger/model"
	"github.com/mdshahjahanmiah/explore-go/logging"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"time"
)

var ErrInvalidRequestType = errors.New("invalid request type")

func makeDepositEndpoint(s Service, logger *logging.Logger) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(TransactionRequest)
		if !ok {
			logger.Error("invalid deposit request type")
			return nil, ErrInvalidRequestType
		}

		txnID := model.NewUUID()
		txn := model.Transaction{
			ID:          txnID,
			AccountID:   req.AccountID,
			Type:        TransactionTypeDeposit,
			Amount:      model.Decimal{Decimal: decimal.NewFromFloat(req.Amount)},
			Currency:    req.Currency,
			ReferenceID: req.ReferenceID,
			Status:      TransactionStatusPending,
			CreatedAt:   time.Now().UTC(),
		}

		result, err := s.CreateTransaction(ctx, txn)
		if err != nil {
			logger.Error("deposit failed", "transaction_id", txnID, "account_id", req.AccountID, "error", err)
			return nil, err
		}

		logger.Info("deposit queued", "transaction_id", txnID, "account_id", req.AccountID, "amount", req.Amount)
		return result, nil
	}
}

func makeWithdrawEndpoint(s Service, logger *logging.Logger) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(TransactionRequest)
		if !ok {
			logger.Error("invalid withdrawal request type")
			return nil, ErrInvalidRequestType
		}

		txnID := model.NewUUID()
		txn := model.Transaction{
			ID:          txnID,
			AccountID:   req.AccountID,
			Type:        TransactionTypeWithdrawal,
			Amount:      model.Decimal{Decimal: decimal.NewFromFloat(req.Amount)},
			Currency:    req.Currency,
			ReferenceID: req.ReferenceID,
			Status:      TransactionStatusPending,
			CreatedAt:   time.Now().UTC(),
		}

		result, err := s.CreateTransaction(ctx, txn)
		if err != nil {
			logger.Error("withdrawal failed", "transaction_id", txnID, "account_id", req.AccountID, "error", err)
			return nil, err
		}

		logger.Info("withdrawal queued", "transaction_id", txnID, "account_id", req.AccountID, "amount", req.Amount)
		return result, nil
	}
}

func makeAuditEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(AuditRequest)
		if !ok {
			return nil, ErrInvalidRequestType
		}
		return s.GetTransactions(req.AccountID)
	}
}
