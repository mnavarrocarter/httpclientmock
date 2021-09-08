package httpclientmock

import "net/http"

// A TestHttpFunc represents a function that takes a request and return a response or error
//
// It implements http.RoundTripper as well as the same signature for the Do function in http.Client
type TestHttpFunc func(req *http.Request) (*http.Response, error)

func (tc TestHttpFunc) Do(req *http.Request) (*http.Response, error) {
	return tc(req)
}

func (tc TestHttpFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return tc(req)
}
