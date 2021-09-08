package httpclientmock

import (
	"bytes"
	"io"
	"net/http"
)

type Response struct {
	StatusCode int
	Headers    map[string]string
	Body       []byte
	Err        error
}

// Make makes a response or an error according to the passed spec
func (res *Response) Make() (*http.Response, error) {
	if res.Err != nil {
		return nil, res.Err
	}

	head := http.Header{}
	if res.Headers != nil && len(res.Headers) > 0 {
		for k, v := range res.Headers {
			head.Set(k, v)
		}
	}

	if res.Body == nil {
		res.Body = []byte("")
	}

	return &http.Response{
		StatusCode: res.StatusCode,
		Header:     head,
		Body:       io.NopCloser(bytes.NewBuffer(res.Body)),
	}, nil
}
