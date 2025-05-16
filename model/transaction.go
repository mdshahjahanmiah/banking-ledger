package model

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"time"
)

type Transaction struct {
	ID          string          `json:"id"`
	AccountID   string          `json:"account_id"`
	Type        string          `json:"type"`
	Amount      decimal.Decimal `json:"amount"`
	Currency    string          `json:"currency"`
	ReferenceID string          `json:"reference_id"`
	Status      string          `json:"status"`
	CreatedAt   time.Time       `json:"created_at"`
}

func NewUUID() string {
	return uuid.New().String()
}

func IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}
