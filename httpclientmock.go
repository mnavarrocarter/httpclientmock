// Package httpclientmock implements a small expectation framework to test code that depends on http.Client
//
// It provides a fluent api to define those expectations and create clients based on them
package httpclientmock

import (
	"net/http"
	"testing"
)

// Mock represents the assertions to be made about an incoming request and the response
// it should return.
type Mock struct {
	Expect *Request
	Return *Response
}

// GetTestFunc returns a TestHttpFunc that can be used as a http.RoundTripper
func (m *Mock) GetTestFunc(t *testing.T) TestHttpFunc {
	return func(req *http.Request) (*http.Response, error) {
		m.Expect.Test(t, req)
		return m.Return.Make()
	}
}

// InjectInClient injects a TestHttpFunc in the pass client
// If client = nil, then http.DefaultClient is used
//
// This function returns another function that will restore the state of client
// to what it was before the injection. You should defer that function
// so the state is restored after your test finish
func (m *Mock) InjectInClient(t *testing.T, client *http.Client) func() {
	if client == nil {
		client = http.DefaultClient
	}

	oldTransport := client.Transport
	client.Transport = m.GetTestFunc(t)

	return func() {
		client.Transport = oldTransport
	}
}

// BuildNewClient returns a http.Client that has a TestHttpFunc as a transport
func (m *Mock) BuildNewClient(t *testing.T) *http.Client {
	return &http.Client{
		Transport: m.GetTestFunc(t),
	}
}
