package util

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetRequest(t *testing.T) {
	url := "https://www.google.com"
	buf, err := SendRequest(http.MethodGet, url, nil, nil)
	require.Nil(t, err)
	assert.NotEmpty(t, buf)
}
