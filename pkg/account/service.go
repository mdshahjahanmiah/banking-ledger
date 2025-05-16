package account

import (
	"context"
	"github.com/google/uuid"
	"github.com/mdshahjahanmiah/banking-ledger/model"
	"github.com/mdshahjahanmiah/banking-ledger/pkg/config"
	"github.com/mdshahjahanmiah/banking-ledger/pkg/db"
	"github.com/mdshahjahanmiah/explore-go/logging"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type Service interface {
	CreateAccount(ctx context.Context, input CreateAccountRequest) (*model.Account, error)
}

type service struct {
	config config.Config
	logger *logging.Logger
	store  Store
}

type CreateAccountRequest struct {
	UserID   string          `json:"user_id" validate:"required"`
	Currency string          `json:"currency" validate:"required,len=3"`
	Balance  decimal.Decimal `json:"balance" validate:"gte=0"`
}

func NewService(config config.Config, logger *logging.Logger, database *db.DB) Service {
	return &service{
		config: config,
		logger: logger,
		store:  NewStore(database),
	}
}

func (s *service) CreateAccount(ctx context.Context, req CreateAccountRequest) (*model.Account, error) {
	account := &model.Account{
		ID:       uuid.NewString(),
		UserID:   req.UserID,
		Balance:  req.Balance,
		Currency: req.Currency,
		Status:   model.AccountStatusActive,
	}

	// Validate account
	if err := account.Validate(); err != nil {
		s.logger.Error("account validation failed", "error", err)
		return nil, errors.Wrap(ErrInvalidAccount, err.Error())
	}

	// Store in database
	if err := s.store.Insert(ctx, account); err != nil {
		s.logger.Error("failed to create account", "error", err)
		return nil, errors.Wrap(err, "failed to create account")
	}

	s.logger.Info("account created successfully", "account_id", account.ID)
	return account, nil
}
