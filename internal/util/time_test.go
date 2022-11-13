package util

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNextDay(t *testing.T) {
	tm, err := time.Parse("20060102", "20220101")
	require.Nil(t, err)
	assert.Equal(t, "20220101", tm.Format("20060102"))
	NextDay(&tm)
	assert.Equal(t, "20220102", tm.Format("20060102"))
}
