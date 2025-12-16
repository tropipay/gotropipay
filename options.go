package gotropipay

import (
	"net/http"
	"time"
)

// Environment defines the API environment (Production, Sandbox, etc.)
type Environment string

const (
	// ProductionEnv is the live environment URL
	ProductionEnv Environment = "https://www.tropipay.com/api/v3"
	// SandboxEnv is the test environment URL
	SandboxEnv Environment = "https://sandbox.tropipay.me/api/v3"
)

// Option is a functional option for configuring the Client
type Option func(*Client)

// WithHTTPClient sets a custom HTTP client
func WithHTTPClient(client *http.Client) Option {
	return func(c *Client) {
		if client != nil {
			c.httpClient = client
		}
	}
}

// WithEnvironment sets the Tropipay environment (Sandbox/Production)
func WithEnvironment(env Environment) Option {
	return func(c *Client) {
		c.baseURL = string(env)
	}
}

// WithBaseURL sets a custom base URL (useful for 'Custom' environment or proxies)
func WithBaseURL(url string) Option {
	return func(c *Client) {
		c.baseURL = url
	}
}

// WithTimeout sets the timeout for requests
func WithTimeout(d time.Duration) Option {
	return func(c *Client) {
		c.httpClient.Timeout = d
	}
}
