package gotropipay

import (
	"context"
	"fmt"
)

// PaymentCard represents a payment link or card payment order resource
type PaymentCard struct {
	ID                    string      `json:"id"`
	CredentialID          interface{} `json:"credentialId"` // Use interface{} or specific pointer if knwon
	Reference             string      `json:"reference"`
	Concept               string      `json:"concept"`
	Description           string      `json:"description"`
	Amount                int64       `json:"amount"` // Amount in smallest unit (e.g., cents)
	Currency              string      `json:"currency"`
	SingleUse             bool        `json:"singleUse"`
	ReasonID              int         `json:"reasonId"`
	ReasonDes             string      `json:"reasonDes"`
	UserID                string      `json:"userId"`
	QRImage               string      `json:"qrImage"` // Using string, json decoder handles null as "" often or use *string
	ShortURL              string      `json:"shortUrl"`
	State                 int         `json:"state"`
	ExpirationDays        int         `json:"expirationDays"`
	Lang                  string      `json:"lang"`
	URLSuccess            string      `json:"urlSuccess"`
	URLFailed             string      `json:"urlFailed"`
	URLNotification       string      `json:"urlNotification"`
	AccountID             int64       `json:"accountId"`
	ExpirationDate        string      `json:"expirationDate"`
	ServiceDate           string      `json:"serviceDate"`
	HasClient             bool        `json:"hasClient"`
	PaymentURL            string      `json:"paymentUrl"`
	Favorite              bool        `json:"favorite"`
	SaveToken             bool        `json:"saveToken"`
	PaymentCardType       int         `json:"paymentcardType"`
	ImageBase             string      `json:"imageBase"`
	Force3DS              bool        `json:"force3ds"`
	Origin                int         `json:"origin"`
	StrictPostalCodeCheck bool        `json:"strictPostalCodeCheck"`
	StrictAddressCheck    bool        `json:"strictAddressCheck"`
	DestinationCurrency   string      `json:"destinationCurrency"`
	Payment3DS            int         `json:"payment3DS"`
	CreatedAt             string      `json:"createdAt"`
	UpdatedAt             string      `json:"updatedAt"`
}

// CreatePaymentCardRequest represents the payload to create a card
type CreatePaymentCardRequest struct {
	Number      string `json:"number"`
	CVC         string `json:"cvc"`
	HolderName  string `json:"holderName"`
	ExpiryMonth int    `json:"expiryMonth"`
	ExpiryYear  int    `json:"expiryYear"`
}

// CreatePaymentCard adds a new payment card
func (c *Client) CreatePaymentCard(ctx context.Context, req CreatePaymentCardRequest) (*PaymentCard, error) {
	var card PaymentCard
	// Assuming endpoint is /paymentcards
	err := c.Request(ctx, "POST", "/paymentcards", req, &card)
	if err != nil {
		return nil, err
	}
	return &card, nil
}

// GetPaymentCard retrieves a specific payment card
func (c *Client) GetPaymentCard(ctx context.Context, id string) (*PaymentCard, error) {
	var card PaymentCard
	path := fmt.Sprintf("/paymentcards/%s", id)
	err := c.Request(ctx, "GET", path, nil, &card)
	if err != nil {
		return nil, err
	}
	return &card, nil
}

// DeletePaymentCard removes a payment card
func (c *Client) DeletePaymentCard(ctx context.Context, id string) error {
	path := fmt.Sprintf("/paymentcards/%s", id)
	return c.Request(ctx, "DELETE", path, nil, nil)
}

// ListPaymentCards retrieves all payment cards
func (c *Client) ListPaymentCards(ctx context.Context) ([]PaymentCard, error) {
	// Often returns a list or a wrapper. Assuming list for simple example
	var cards []PaymentCard
	err := c.Request(ctx, "GET", "/paymentcards", nil, &cards)
	if err != nil {
		return nil, err
	}
	return cards, nil
}
