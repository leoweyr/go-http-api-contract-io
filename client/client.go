package client

import (
	"net/http"
)

// Client is the outbound entry to a single service identified by its base URL.
type Client struct {
	httpClient *http.Client
	origin     *Origin
}

// NewClient creates a service entry over the given round tripper that does not follow redirects, bound to a base URL.
func NewClient(roundTripper http.RoundTripper, baseURL string) *Client {
	var httpClient *http.Client = &http.Client{
		Transport: roundTripper,
		CheckRedirect: func(httpRequest *http.Request, viaRequests []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	var client *Client = &Client{
		httpClient: httpClient,
		origin:     newOrigin(httpClient, baseURL),
	}

	return client
}

// GetOrigin exposes the service origin as the sole endpoint resolution point.
func (client *Client) GetOrigin() *Origin {
	return client.origin
}
