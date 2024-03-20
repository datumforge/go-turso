package turso

import "net/http"

// MockHTTPRequestDoer implements the standard http.Client interface.
type MockHTTPRequestDoer struct {
	Response *http.Response
	Error    error
}

// Do implements the standard http.Client interface for MockHTTPRequestDoer
func (md *MockHTTPRequestDoer) Do(req *http.Request) (*http.Response, error) {
	return md.Response, md.Error
}
