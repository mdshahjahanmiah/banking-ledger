package account_test

import (
	"context"
	database "github.com/mdshahjahanmiah/banking-ledger/pkg/db"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/mdshahjahanmiah/banking-ledger/model"
	"github.com/mdshahjahanmiah/banking-ledger/pkg/account"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestStore_Insert_Success(t *testing.T) {
	// Setup sqlmock
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Wrap *sql.DB into your db.DB struct, assuming it has a field DB *sql.DB
	// (adjust this if your db.DB is different)
	store := account.NewStore(&database.DB{DB: db})

	// Prepare the account to insert
	balance := decimal.NewFromFloat(123.45)
	acc := &model.Account{
		ID:       "acc1",
		UserID:   "user1",
		Balance:  balance,
		Currency: "USD",
		Status:   model.AccountStatusActive,
	}

	// Mock expected DB behavior: QueryRowContext for INSERT returning created_at, updated_at
	rows := sqlmock.NewRows([]string{"created_at", "updated_at"}).
		AddRow(time.Now(), time.Now())

	mock.ExpectQuery(`INSERT INTO accounts .* RETURNING created_at, updated_at`).
		WithArgs(acc.ID, acc.UserID, balance.String(), acc.Currency, acc.Status).
		WillReturnRows(rows)

	// Call the method
	err = store.Insert(context.Background(), acc)
	assert.NoError(t, err)

	// Ensure expectations were met
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
