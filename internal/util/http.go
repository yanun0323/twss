package util

import (
	"io"
	"net/http"

	"github.com/pkg/errors"
)

func GetRequest(url string) ([]byte, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, errors.New("Bad request")
	}

	buf, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return buf, nil
}
