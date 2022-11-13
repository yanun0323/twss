package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetRequest(t *testing.T) {
	url := "https://www.google.com"
	buf, err := GetRequest(url)
	require.Nil(t, err)
	assert.NotEmpty(t, buf)
}
