package client

import (
	"net/http"
)

// Origin is the base URL of a service, it resolves paths into sendable endpoints over a shared HTTP client.
type Origin struct {
	httpClient *http.Client
	baseURL    string
}

func newOrigin(httpClient *http.Client, baseURL string) *Origin {
	var origin *Origin = &Origin{
		httpClient: httpClient,
		baseURL:    baseURL,
	}

	return origin
}

// Endpoint binds an HTTP method and path to a sendable endpoint against the origin base URL.
func (origin *Origin) Endpoint(method string, path string) *Endpoint {
	var endpoint *Endpoint = newEndpoint(origin.httpClient, method, origin.baseURL+path)

	return endpoint
}
