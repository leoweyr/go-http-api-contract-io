package client

import "net/http"

// RoundTripperFunc adapts a bare function to the http.RoundTripper interface.
type RoundTripperFunc func(*http.Request) (*http.Response, error)

// RoundTrip calls the wrapped function.
func (function RoundTripperFunc) RoundTrip(httpRequest *http.Request) (*http.Response, error) {
	return function(httpRequest)
}
