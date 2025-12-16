package gotropipay

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// TokenResponse represents the JSON response from the login endpoint
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
}

type authenticator struct {
	clientID     string
	clientSecret string
	baseURL      string
	httpClient   *http.Client

	mu          sync.Mutex
	accessToken string
	expiresAt   time.Time
}

func newAuthenticator(clientID, clientSecret, baseURL string, client *http.Client) *authenticator {
	return &authenticator{
		clientID:     clientID,
		clientSecret: clientSecret,
		baseURL:      baseURL,
		httpClient:   client,
	}
}

// GetToken returns a valid access token, refreshing it if necessary
func (a *authenticator) GetToken() (string, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	// Check if token is valid (with 10-second buffer)
	if a.accessToken != "" && time.Now().Add(10*time.Second).Before(a.expiresAt) {
		return a.accessToken, nil
	}

	return a.refreshToken()
}

func (a *authenticator) refreshToken() (string, error) {
	// Payload for login
	payload := map[string]string{
		"grant_type":    "client_credentials",
		"client_id":     a.clientID,
		"client_secret": a.clientSecret,
	}

	jsonBody, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	// Actually, most logical is: POST <baseURL>/access/token or similar.

	req, err := http.NewRequest("POST", a.baseURL+"/access/token", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var body bytes.Buffer
		_, _ = body.ReadFrom(resp.Body)
		return "", fmt.Errorf("authentication failed: status %d, body: %s", resp.StatusCode, body.String())
	}

	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", err
	}

	a.accessToken = tokenResp.AccessToken
	a.expiresAt = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)

	return a.accessToken, nil
}
