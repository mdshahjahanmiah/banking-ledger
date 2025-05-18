package model

import (
	"encoding/json"

	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

// Decimal wraps decimal.Decimal to support BSON and JSON marshaling.
type Decimal struct {
	decimal.Decimal
}

// MarshalBSONValue stores Decimal as a string in MongoDB.
func (d Decimal) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return bson.MarshalValue(d.String())
}

// UnmarshalBSONValue reads Decimal stored as string from MongoDB.
func (d *Decimal) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	var s string
	if err := bson.UnmarshalValue(t, data, &s); err != nil {
		return err
	}
	dec, err := decimal.NewFromString(s)
	if err != nil {
		return err
	}
	d.Decimal = dec
	return nil
}

// MarshalJSON returns Decimal as a JSON string.
func (d Decimal) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

// UnmarshalJSON parses Decimal from a JSON string.
func (d *Decimal) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	dec, err := decimal.NewFromString(s)
	if err != nil {
		return err
	}
	d.Decimal = dec
	return nil
}

func (d Decimal) Unwrap() decimal.Decimal {
	return d.Decimal
}
