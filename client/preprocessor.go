package client

import (
	"net/http"
)

// Preprocessor wraps a round tripper to preprocess the outbound request before delegating to the next round tripper.
type Preprocessor interface {
	Call(next http.RoundTripper) http.RoundTripper
}
