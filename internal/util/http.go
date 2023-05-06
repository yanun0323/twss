package util

import (
	"fmt"
	"io"
	"net/http"

	"github.com/pkg/errors"
)

func SendRequest(method, url string, fn func(*http.Request), reqBody io.Reader) ([]byte, error) {
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, errors.Wrap(err, "new request")
	}

	if fn != nil {
		fn(req)
	}

	for k, v := range req.Form {
		fmt.Printf("Form K: %s, V: %+v", k, v)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.Request != nil {
		for k, v := range res.Request.Form {
			fmt.Printf("Form K: %s, V: %+v", k, v)
		}
	}
	fmt.Println("Status Code:", res.StatusCode)
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, errors.New("bad request")
	}

	buf, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return buf, nil
}
