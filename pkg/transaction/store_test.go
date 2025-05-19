package transaction_test

import (
	"context"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/mdshahjahanmiah/banking-ledger/model"
	"github.com/mdshahjahanmiah/banking-ledger/pkg/db"
	"github.com/mdshahjahanmiah/banking-ledger/pkg/transaction"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProcessTransaction_Success(t *testing.T) {
	// Create mock db and sqlmock
	sqlDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer sqlDB.Close()

	// Wrap sqlDB with your db.DB
	mockDB := &db.DB{DB: sqlDB}
	store := transaction.NewStore(mockDB)

	ctx := context.Background()

	// Transaction to test
	txn := model.Transaction{
		ID:          "txn1",
		AccountID:   "acc1",
		ReferenceID: "ref1",
		Currency:    "USD",
		Amount:      model.Decimal{Decimal: decimal.NewFromFloat(10)},
		Type:        transaction.TransactionTypeDeposit,
	}

	// Begin transaction expectation
	mock.ExpectBegin()

	// Expect select to check existing transaction (no rows)
	mock.ExpectQuery(`SELECT id FROM transactions WHERE reference_id = \$1 AND currency = \$2`).
		WithArgs(txn.ReferenceID, txn.Currency).
		WillReturnError(sql.ErrNoRows)

	// Expect select for account details with FOR UPDATE
	rows := sqlmock.NewRows([]string{"id", "balance", "currency", "status"}).
		AddRow("acc1", decimal.NewFromFloat(200).String(), "USD", model.AccountStatusActive)
	mock.ExpectQuery(`SELECT id, balance, currency, status FROM accounts WHERE id = \$1 FOR UPDATE`).
		WithArgs(txn.AccountID).
		WillReturnRows(rows)

	// Expect update account balance
	mock.ExpectExec(`UPDATE accounts SET balance = \$1, updated_at = NOW\(\) WHERE id = \$2`).
		WithArgs(sqlmock.AnyArg(), txn.AccountID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Expect insert transaction
	mock.ExpectExec(`INSERT INTO transactions \(id, account_id, amount, type, reference_id, currency, status, created_at\) VALUES \(\$1, \$2, \$3, \$4, \$5, \$6, \$7, \$8\)`).
		WithArgs(txn.ID, txn.AccountID, txn.Amount, txn.Type, txn.ReferenceID, txn.Currency, transaction.TransactionStatusCompleted, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Expect commit
	mock.ExpectCommit()

	err = store.ProcessTransaction(ctx, txn)
	assert.NoError(t, err)

	// ensure all expectations met
	assert.NoError(t, mock.ExpectationsWereMet())
}
