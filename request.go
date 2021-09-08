package httpclientmock

import (
	"io"
	"net/http"
	"testing"
)

type Request struct {
	Method  string
	Url     string
	Headers map[string]string
	Body    []byte
}

// Test tests whether the request conforms to the spec
func (r *Request) Test(t *testing.T, req *http.Request) {
	if r.Method != "" {
		AssertEqual(t, r.Method, req.Method, `Unexpected method "%s" for request. Expected: %s`, req.Method, r.Method)
	}

	if r.Url != "" {
		AssertEqual(t, r.Url, req.URL.String(), `Unexpected url "%s" for request. Expected: %s`, req.URL.String(), r.Url)
	}

	if req.Header != nil && len(r.Headers) > 0 {
		for k, v := range r.Headers {
			hVal := req.Header.Get(k)
			if hVal == "" {
				Errorf(t, `Expected header "%s" is not present in the request`, k)
			} else {
				AssertEqual(t, v, hVal, `Unexpected value "%s" for header "%s". Expected: %s`, hVal, k, v)
			}
		}
	}

	if r.Body != nil && len(r.Body) > 0 {
		b, err := io.ReadAll(req.Body)
		if err != nil {
			panic(err)
		}
		AssertEqual(t, r.Body, b, "Request body is not the same as expected")
	}
}
