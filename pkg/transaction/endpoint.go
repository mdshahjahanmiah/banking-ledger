package transaction

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/mdshahjahanmiah/banking-ledger/model"
	"github.com/mdshahjahanmiah/explore-go/logging"
	"github.com/shopspring/decimal"
	"time"
)

func makeDepositEndpoint(s Service, logger *logging.Logger) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(TransactionRequest)

		txn := model.Transaction{
			ID:          model.NewUUID(),
			AccountID:   req.AccountID,
			Type:        TransactionTypeDeposit,
			Amount:      decimal.NewFromFloat(req.Amount),
			Currency:    req.Currency,
			ReferenceID: req.ReferenceID,
			Status:      TransactionStatusPending,
			CreatedAt:   time.Now().UTC(),
		}

		result, err := s.CreateTransaction(ctx, txn)
		if err != nil {
			logger.Error("deposit failed", "account", req.AccountID, "error", err)
			return nil, err
		}

		logger.Info("deposit queued", "transaction_id", txn.ID, "amount", req.Amount)
		return result, nil
	}
}

func makeWithdrawEndpoint(s Service, logger *logging.Logger) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(TransactionRequest)

		txn := model.Transaction{
			ID:          model.NewUUID(),
			AccountID:   req.AccountID,
			Type:        TransactionTypeWithdrawal,
			Amount:      decimal.NewFromFloat(req.Amount),
			Currency:    req.Currency,
			ReferenceID: req.ReferenceID,
			Status:      TransactionStatusPending,
			CreatedAt:   time.Now().UTC(),
		}

		result, err := s.CreateTransaction(ctx, txn)
		if err != nil {
			logger.Error("withdrawal failed", "account", req.AccountID, "error", err)
			return nil, err
		}

		logger.Info("withdrawal queued", "transaction_id", txn.ID, "amount", req.Amount)
		return result, nil
	}
}

func makeAuditEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(AuditRequest)
		return s.GetTransactions(req.AccountID)
	}
}
