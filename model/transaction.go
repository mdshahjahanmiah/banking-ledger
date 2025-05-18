package model

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"regexp"
	"strings"
	"time"
)

var (
	ErrInvalidTransactionID   = errors.New("invalid transaction ID")
	ErrInvalidAccountID       = errors.New("invalid account ID")
	ErrInvalidReferenceID     = errors.New("invalid reference ID")
	ErrInvalidAmount          = errors.New("amount must be greater than zero")
	ErrInvalidTransactionType = errors.New("transaction type must be deposit or withdrawal")
	ErrInvalidCurrency        = errors.New("invalid currency format")
)

type Transaction struct {
	ID          string    `json:"id"`
	AccountID   string    `json:"account_id"`
	Type        string    `json:"type"`
	Amount      Decimal   `json:"amount" bson:"amount"`
	Currency    string    `json:"currency"`
	ReferenceID string    `json:"reference_id"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

func (t *Transaction) Validate() error {
	if !IsValidUUID(t.ID) {
		return ErrInvalidTransactionID
	}
	if !IsValidUUID(t.AccountID) {
		return ErrInvalidAccountID
	}
	if !IsValidUUID(t.ReferenceID) {
		return ErrInvalidReferenceID
	}
	if t.Amount.LessThanOrEqual(decimal.Zero) {
		return ErrInvalidAmount
	}

	switch strings.ToLower(t.Type) {
	case "deposit", "withdrawal":
	default:
		return ErrInvalidTransactionType
	}

	if !isValidCurrency(t.Currency) {
		return ErrInvalidCurrency
	}

	return nil
}

func isValidCurrency(currency string) bool {
	// Accept 3 uppercase letters (ISO 4217 format)
	match, _ := regexp.MatchString("^[A-Z]{3}$", currency)
	return match
}

func NewUUID() string {
	return uuid.New().String()
}

func IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}
