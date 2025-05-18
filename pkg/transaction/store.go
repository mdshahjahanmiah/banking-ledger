package transaction

import (
	"context"
	"database/sql"
	"github.com/mdshahjahanmiah/banking-ledger/model"
	"github.com/mdshahjahanmiah/banking-ledger/pkg/db"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"time"
)

type Store interface {
	ProcessTransaction(ctx context.Context, txn model.Transaction) error
}

type store struct {
	db *db.DB
}

func NewStore(db *db.DB) *store {
	return &store{db: db}
}

func (s *store) ProcessTransaction(ctx context.Context, txn model.Transaction) error {
	tx, err := s.db.DB.BeginTx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "failed to begin transaction")
	}
	defer tx.Rollback()

	// Idempotency Checking for existing transaction
	var existingID string
	err = tx.QueryRowContext(ctx,
		`SELECT id FROM transactions WHERE reference_id = $1 AND currency = $2`, txn.ReferenceID, txn.Currency,
	).Scan(&existingID)

	if err == nil {
		return ErrDuplicateTransaction
	} else if !errors.Is(err, sql.ErrNoRows) {
		return errors.Wrap(err, "failed to check existing transactions")
	}

	// Get account details with locking
	var account model.Account
	err = tx.QueryRowContext(ctx,
		`SELECT id, balance, currency, status FROM accounts WHERE id = $1 FOR UPDATE`,
		txn.AccountID,
	).Scan(&account.ID, &account.Balance, &account.Currency, &account.Status)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrAccountNotFound
		}
		return errors.Wrap(err, "failed to get account details")
	}

	if account.Status != model.AccountStatusActive {
		return ErrAccountNotActive
	}

	// Validating transaction amount
	if txn.Amount.LessThanOrEqual(decimal.Zero) {
		return ErrInvalidAmount
	}

	// Calculating new balance
	newBalance := account.Balance
	switch txn.Type {
	case TransactionTypeDeposit:
		newBalance = newBalance.Add(txn.Amount.Unwrap())
	case TransactionTypeWithdrawal:
		if account.Balance.LessThan(txn.Amount.Unwrap()) {
			return ErrInsufficientFunds
		}
		newBalance = newBalance.Sub(txn.Amount.Unwrap())
	default:
		return ErrInvalidTransactionType
	}

	// Updating account balance
	_, err = tx.ExecContext(ctx, `UPDATE accounts SET balance = $1, updated_at = NOW() WHERE id = $2`, newBalance, txn.AccountID)
	if err != nil {
		return errors.Wrap(err, "failed to update account balance")
	}

	// Creating transaction record
	_, err = tx.ExecContext(ctx,
		`INSERT INTO transactions 
		(id, account_id, amount, type, reference_id, currency, status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		txn.ID, txn.AccountID, txn.Amount, txn.Type,
		txn.ReferenceID, txn.Currency, TransactionStatusCompleted, time.Now().UTC(),
	)
	if err != nil {
		return errors.Wrap(err, "failed to create transaction record")
	}

	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "transaction commit failed")
	}
	return nil
}
