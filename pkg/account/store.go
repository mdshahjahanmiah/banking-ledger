package account

import (
	"context"
	"github.com/mdshahjahanmiah/banking-ledger/model"
	"github.com/mdshahjahanmiah/banking-ledger/pkg/db"
	"github.com/pkg/errors"
)

var (
	ErrAccountNotFound = errors.New("account not found")
	ErrInvalidAccount  = errors.New("invalid account")
	ErrInvalidStatus   = errors.New("invalid account status")
)

type Store interface {
	Insert(ctx context.Context, a *model.Account) error
}

type store struct {
	db *db.DB
}

func NewStore(db *db.DB) *store {
	return &store{db: db}
}

func (s *store) Insert(ctx context.Context, a *model.Account) error {
	if err := a.Validate(); err != nil {
		return errors.Wrap(err, "account validation failed")
	}

	err := s.db.DB.QueryRowContext(ctx,
		`INSERT INTO accounts (id, user_id, balance, currency, status)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING created_at, updated_at`,
		a.ID, a.UserID, a.Balance, a.Currency, a.Status,
	).Scan(&a.CreatedAt, &a.UpdatedAt)

	if err != nil {
		return errors.Wrap(err, "failed to create account")
	}

	return nil
}
