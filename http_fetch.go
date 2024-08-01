package gotil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// base http fetcher
type HTTPFetcher struct {
	Method  string
	Url     string
	Body    JSON
	Headers JSON
}

// fetch http
func (h HTTPFetcher) Fetch(result any) error {
	var err error
	var req *http.Request

	if h.Body == nil {
		req, err = http.NewRequest(h.Method, h.Url, nil)
		if err != nil {
			return err
		}
	} else {
		b, err := json.Marshal(h.Body)
		if err != nil {
			return err
		}
		req, err = http.NewRequest(h.Method, h.Url, bytes.NewBuffer(b))
		if err != nil {
			return err
		}
	}

	for k, v := range h.Headers {
		req.Header.Set(k, fmt.Sprintf("%v", v))
	}

	c := http.Client{}
	res, err := c.Do(req)
	if err != nil {
		return err
	}

	return json.NewDecoder(res.Body).Decode(&result)
}
