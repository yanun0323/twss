package model

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func Test_parseSymbol(t *testing.T) {
	assert.Equal(t, "-", parseSymbol("<p style= color:green>-</p>"))
	assert.Equal(t, "+", parseSymbol("<p style= color:red>+</p>"))
	assert.Equal(t, " ", parseSymbol("<p> </p>"))
}

func Test_calculatePercentage(t *testing.T) {
	assert.Equal(t, "-0.5", calculatePercentage("-0.0005", "0.1"))
	assert.Equal(t, "0", calculatePercentage("0", "0.1"))
	assert.Equal(t, "1.2", calculatePercentage("0.0012", "0.1"))
}

func Test_fixPrice(t *testing.T) {
	good := func(expected, input float64) {
		assert.True(t, d(expected).Equal(fixPrice(d(input))))
	}
	good(9.040, 9.044)
	good(9.060, 9.055)
	good(10.00, 10.044)
	good(10.05, 10.055)

	good(49.00, 49.044)
	good(49.05, 49.055)
	good(50.40, 50.444)
	good(50.60, 50.555)

	good(99.40, 99.44)
	good(99.60, 99.55)
	good(100.0, 100.44)
	good(100.5, 100.55)

	good(499.0, 499.44)
	good(499.5, 499.55)
	good(504.0, 504.4)
	good(506.0, 505.5)

	good(994.0, 994.4)
	good(996.0, 995.5)
	good(1000, 1004.4)
	good(1005, 1005.5)
}

func Test_calculateLimit(t *testing.T) {
	good := func(close, grade float64) {
		assert.True(t, calculateLimit(d(close).String(), d(grade).String()))
	}
	bad := func(close, grade float64) {
		assert.False(t, calculateLimit(d(close).String(), d(grade).String()))
	}

	good(21.05, 1.90)
	good(11.30, 1.00)
	good(10.90, 0.95)
	good(10.95, 0.99)
	good(8.97, -0.99)

	bad(21.05, 1.85)
	bad(11.30, 0.95)
	bad(10.90, 0.90)

	assert.Equal(t, int32(-3), d(0.003).Exponent())
}

func d(f float64) decimal.Decimal {
	return decimal.NewFromFloat(f)
}
