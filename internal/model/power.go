package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type PowerInput struct {
	ID   string
	From time.Time
	To   time.Time
}

type PowerOutput struct {
	PowerInput

	Power map[time.Time]decimal.Decimal
}
