package util

import (
	"regexp"
	"strconv"

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

func Int(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return i
}
