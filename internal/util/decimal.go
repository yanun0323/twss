package util

import (
	"regexp"

	"github.com/shopspring/decimal"
)

var _regexpDecimal = regexp.MustCompile(`[\,]`)

func Decimal(s string) decimal.Decimal {
	s = _regexpDecimal.ReplaceAllString(s, "")
	d, err := decimal.NewFromString(s)
	if err != nil {
		return decimal.Zero
	}
	return d
}
