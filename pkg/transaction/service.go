package transaction

import (
	"context"
	"github.com/google/uuid"
	"github.com/mdshahjahanmiah/banking-ledger/model"
	"github.com/mdshahjahanmiah/banking-ledger/pkg/broker"
	"github.com/mdshahjahanmiah/banking-ledger/pkg/config"
	"github.com/mdshahjahanmiah/banking-ledger/pkg/db"
	eError "github.com/mdshahjahanmiah/explore-go/error"
	"github.com/mdshahjahanmiah/explore-go/logging"
	"github.com/mdshahjahanmiah/explore-go/repository"
	"github.com/pkg/errors"
	"net/http"
	"time"
)

var (
	ErrDuplicateTransaction   = errors.New("duplicate transaction")
	ErrAccountNotFound        = errors.New("account not found")
	ErrAccountNotActive       = errors.New("account not active")
	ErrInvalidAmount          = errors.New("invalid amount")
	ErrInsufficientFunds      = errors.New("insufficient funds")
	ErrInvalidTransactionType = errors.New("invalid transaction type")
)

const (
	TransactionTypeDeposit    = "deposit"
	TransactionTypeWithdrawal = "withdrawal"

	TransactionStatusPending   = "pending"
	TransactionStatusCompleted = "completed"
	TransactionStatusFailed    = "failed"
)

type Service interface {
	CreateTransaction(ctx context.Context, input model.Transaction) (model.Transaction, error)
	ProcessTransaction(ctx context.Context, txn model.Transaction) error
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
	if input.ReferenceID == "" {
		input.ReferenceID = uuid.NewString()
	} else if !model.IsValidUUID(input.ReferenceID) {
		return model.Transaction{}, eError.NewServiceError(
			errors.New("invalid reference_id format"), "reference_id must be a valid UUID", "INVALID_REFERENCE_ID", http.StatusBadRequest)
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

	// Validate transaction
	if err := txn.Validate(); err != nil {
		return model.Transaction{}, err
	}

	// Publish to Kafka (or other broker)
	if err := s.producer.PublishTransaction(txn); err != nil {
		s.logger.Error("failed to publish transaction", "reference_id", txn.ReferenceID, "error", err)
		return model.Transaction{}, err
	}

	s.logger.Info("transaction queued successfully", "reference_id", txn.ReferenceID, "amount", txn.Amount, "currency", txn.Currency)
	return txn, nil
}

func (s *service) ProcessTransaction(ctx context.Context, txn model.Transaction) error {
	// Validate transaction
	if err := txn.Validate(); err != nil {
		return err
	}

	// Process transaction in the store
	if err := s.store.ProcessTransaction(ctx, txn); err != nil {
		return err
	}

	s.logger.Info("transaction processed successfully", "reference_id", txn.ReferenceID, "amount", txn.Amount, "currency", txn.Currency)
	return nil

}

func (s *service) GetTransactions(accountID string) ([]model.Transaction, error) {
	// Get transactions from mongodb
	return s.repo.FindByField("accountid", accountID)
}
