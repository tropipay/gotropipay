package gotropipay_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/tropipay/gotropipay"
)

func newTestClient(t *testing.T) *gotropipay.Client {
	t.Helper()
	// Load .env from project root (assuming test runs from root or package dir)
	_ = godotenv.Load(".env") // Ignore error, env vars might be set in system

	clientID := os.Getenv("TROPIPAY_CLIENT_ID")
	clientSecret := os.Getenv("TROPIPAY_CLIENT_SECRET")

	if clientID == "" || clientSecret == "" {
		t.Skip("Skipping integration test: TROPIPAY_CLIENT_ID or TROPIPAY_CLIENT_SECRET not set")
	}

	return gotropipay.NewClient(clientID, clientSecret, gotropipay.WithEnvironment(gotropipay.SandboxEnv))
}

func TestAuthentication(t *testing.T) {
	client := newTestClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// ListPaymentCards uses helper Request logic which does automatic authentication
	cards, err := client.ListPaymentCards(ctx)
	if err != nil {
		t.Fatalf("Authentication or Request failed: %v", err)
	}
	t.Logf("Authentication successful. Found %d cards", len(cards))
}

func TestListPaymentCards(t *testing.T) {
	client := newTestClient(t)
	ctx := context.Background()

	cards, err := client.ListPaymentCards(ctx)
	if err != nil {
		t.Fatalf("Failed to list payment cards: %v", err)
	}

	for _, card := range cards {
		t.Logf("Card: %s - %s", card.Concept, card.ShortURL)
	}
}

func TestGetUserProfile(t *testing.T) {
	client := newTestClient(t)
	ctx := context.Background()

	user, err := client.GetUserProfile(ctx)
	if err != nil {
		t.Fatalf("Failed to get user profile: %v", err)
	}

	t.Logf("User Profile: %s %s (Email: %s)", user.Name, user.Surname, user.Email)
}
