package gotropipay

import (
	"net/http"
	"time"
)

// Client is the main entry point for the Tropipay API
type Client struct {
	clientID     string
	clientSecret string
	baseURL      string
	httpClient   *http.Client

	// auth holds the authentication state and logic
	auth *authenticator
}

// NewClient creates a new Tropipay API client
func NewClient(clientID, clientSecret string, opts ...Option) *Client {
	c := &Client{
		clientID:     clientID,
		clientSecret: clientSecret,
		baseURL:      string(ProductionEnv), // default is production
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}

	// Apply options
	for _, opt := range opts {
		opt(c)
	}

	// Initialize authenticator
	c.auth = newAuthenticator(clientID, clientSecret, c.baseURL, c.httpClient)

	return c
}
