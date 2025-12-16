package gotropipay

import "context"

// AddTropicardAccountRequest represents the payload to link a Tropicard
type AddTropicardAccountRequest struct {
	TropicardNumber string `json:"tropicardNumber"`
	Pin             string `json:"pin"`
}

// CryptoAddress represents a deposit address for a specific network and currency
type CryptoAddress struct {
	Address  string `json:"address"`
	Network  string `json:"network"`
	Currency string `json:"currency"`
}

// GetCryptoAddressResponse represents the response for self-charge crypto addresses
type GetCryptoAddressResponse struct {
	FeePercent int             `json:"feePercent"` // Fee as percentage (e.g., 300 = 3.00%)
	FeeFixed   int             `json:"feeFixed"`   // Fixed fee in cents
	Accounts   []CryptoAddress `json:"accounts"`
}

// AddTropicardAccount links a Tropicard to the user's account.
// It returns a generic map.
// You can likely expect an "id" field in the response to use with other Account endpoints.
func (c *Client) AddTropicardAccount(ctx context.Context, req AddTropicardAccountRequest) (map[string]interface{}, error) {
	var resp map[string]interface{}
	err := c.Request(ctx, "POST", "/accounts/", req, &resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// GetCryptoAddressForSelfCharge retrieves cryptocurrency addresses for depositing funds into a specific account.
func (c *Client) GetCryptoAddressForSelfCharge(ctx context.Context, accountID string) (*GetCryptoAddressResponse, error) {
	var resp GetCryptoAddressResponse
	// /accounts/{accountId}/selfcharge/crypto
	path := "/accounts/" + accountID + "/selfcharge/crypto"
	err := c.Request(ctx, "GET", path, nil, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
