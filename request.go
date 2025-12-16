package gotropipay

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Request executes an HTTP request with authentication
func (c *Client) Request(ctx context.Context, method, path string, body interface{}, result interface{}) error {
	// Get Token
	token, err := c.auth.GetToken()
	if err != nil {
		return fmt.Errorf("failed to get token: %w", err)
	}

	var reqBody io.Reader
	if body != nil {
		jsonBytes, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonBytes)
	}

	// Build full URL
	// Simple concatenation, assuming Request path starts with / or baseURL doesn't end with it.
	// Ideally use path.Join or url.Parse but strict strings are faster if careful.
	fullURL := c.baseURL + path

	req, err := http.NewRequestWithContext(ctx, method, fullURL, reqBody)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Handle specific error codes or generic 4xx/5xx
	if resp.StatusCode >= 400 {
		// Try to parse error body if any
		errResp := new(bytes.Buffer)
		_, _ = io.Copy(errResp, resp.Body)
		return fmt.Errorf("API error: %s (status: %d) - %s", req.URL.String(), resp.StatusCode, errResp.String())
	}

	if result != nil {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return nil
}
