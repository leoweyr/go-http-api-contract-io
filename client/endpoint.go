package client

import (
	"context"
	"net/http"
)

// Endpoint is a bound outbound endpoint over a shared HTTP client, it creates and executes requests to a resolved URL.
type Endpoint struct {
	httpClient *http.Client
	method     string
	url        string
}

func newEndpoint(httpClient *http.Client, method string, url string) *Endpoint {
	var endpoint *Endpoint = &Endpoint{
		httpClient: httpClient,
		method:     method,
		url:        url,
	}

	return endpoint
}

// NewRequest creates an empty-bodied request bound to this endpoint's method and URL.
func (endpoint *Endpoint) NewRequest(requestContext context.Context) (*http.Request, error) {
	return http.NewRequestWithContext(requestContext, endpoint.method, endpoint.url, nil)
}

// Do executes the request over the shared HTTP client and returns the response.
func (endpoint *Endpoint) Do(httpRequest *http.Request) (*http.Response, error) {
	return endpoint.httpClient.Do(httpRequest)
}
