package account

import (
	"context"
	"github.com/lib/pq"
	"github.com/mdshahjahanmiah/banking-ledger/model"
	"github.com/mdshahjahanmiah/banking-ledger/pkg/db"
	eError "github.com/mdshahjahanmiah/explore-go/error"
	"github.com/pkg/errors"
	"net/http"
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
		// PostgreSQL unique violation
		if pgErr, ok := err.(*pq.Error); ok {
			if pgErr.Code == "23505" {
				return eError.NewServiceError(err, ErrDuplicateAccountMsg, ErrDuplicateAccountCode, http.StatusConflict)
			}
		}

		// Other unexpected errors
		return eError.NewServiceError(err, ErrInternalServerMsg, ErrInternalServerCode, http.StatusInternalServerError)
	}

	return nil
}
