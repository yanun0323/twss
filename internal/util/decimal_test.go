package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegexp(t *testing.T) {
	good := func(expected, input string) {
		assert.Equal(t, expected, _regexpDecimal.ReplaceAllString(input, ""))
	}

	good("1000000", "1,000,000")
}

func TestDecimal(t *testing.T) {
	good := func(s string) {
		assert.NotEqual(t, "0", Decimal(s).String())
	}

	good("1,000,000")
	good("1,000.000")
	good("+200")
	good("-200")
}
