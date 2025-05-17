package model

import (
	"fmt"
	"github.com/shopspring/decimal"
	"time"
)

type AccountStatus string

const (
	AccountStatusActive    AccountStatus = "active"
	AccountStatusSuspended AccountStatus = "suspended"
	AccountStatusClosed    AccountStatus = "closed"
)

type Account struct {
	ID        string          `json:"id"`
	UserID    string          `json:"user_id"`
	Balance   decimal.Decimal `json:"balance"`  // Using decimal for precise monetary values
	Currency  string          `json:"currency"` // ISO 4217 currency code
	Status    AccountStatus   `json:"status"`   // Enumerated type for safety
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

func (a *Account) Validate() error {
	switch a.Status {
	case AccountStatusActive, AccountStatusSuspended, AccountStatusClosed:
	default:
		return fmt.Errorf("invalid account status: %s", a.Status)
	}

	if len(a.Currency) != 3 {
		return fmt.Errorf("currency must be 3-letter ISO code")
	}

	if a.Balance.IsNegative() {
		return fmt.Errorf("account balance cannot be negative")
	}

	return nil
}
