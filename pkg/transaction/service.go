package transaction

import (
	"context"
	"github.com/google/uuid"
	"github.com/mdshahjahanmiah/banking-ledger/model"
	"github.com/mdshahjahanmiah/banking-ledger/pkg/broker"
	"github.com/mdshahjahanmiah/banking-ledger/pkg/config"
	"github.com/mdshahjahanmiah/banking-ledger/pkg/db"
	"github.com/mdshahjahanmiah/banking-ledger/repository"
	"github.com/mdshahjahanmiah/explore-go/logging"
	"github.com/pkg/errors"
	"time"
)

type Service interface {
	CreateTransaction(ctx context.Context, input model.Transaction) (model.Transaction, error)
	GetTransactions(accountID string) ([]model.Transaction, error)
}

type service struct {
	config   config.Config
	logger   *logging.Logger
	store    *store
	repo     *repository.Repository[model.Transaction]
	producer broker.Producer
}

func NewService(config config.Config, logger *logging.Logger, database *db.DB, repo *repository.Repository[model.Transaction], producer broker.Producer) (Service, error) {
	return &service{
		config:   config,
		logger:   logger,
		store:    NewStore(database),
		repo:     repo,
		producer: producer,
	}, nil
}

func (s *service) CreateTransaction(ctx context.Context, input model.Transaction) (model.Transaction, error) {
	if !model.IsValidUUID(input.ReferenceID) {
		return model.Transaction{}, errors.New("invalid reference_id format")
	}

	if input.ReferenceID == "" {
		input.ReferenceID = uuid.NewString()
	}

	txn := model.Transaction{
		ID:          model.NewUUID(),
		AccountID:   input.AccountID,
		Type:        input.Type,
		Amount:      input.Amount,
		Currency:    input.Currency,
		ReferenceID: input.ReferenceID,
		Status:      TransactionStatusPending,
		CreatedAt:   time.Now().UTC(),
	}

	// Publish to Kafka (or other broker)
	if err := s.producer.PublishTransaction(txn); err != nil {
		s.logger.Error("failed to publish transaction", "reference_id", txn.ReferenceID, "error", err)
		return model.Transaction{}, err
	}

	s.logger.Info("transaction queued successfully", "reference_id", txn.ReferenceID, "amount", txn.Amount, "currency", txn.Currency)
	return txn, nil
}

func (s *service) GetTransactions(accountID string) ([]model.Transaction, error) {
	return s.repo.FindByField("accountid", accountID)
}
