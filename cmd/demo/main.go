package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/yosle/gotropipay"
)

func main() {
	clientID := "a346e4ff0584c5f0695b24d338874241"
	clientSecret := "ceb6de6e0d76d3aab51518c44cb8ecf6"

	if clientID == "" || clientSecret == "" {
		fmt.Println("Usage: Please set TROPIPAY_CLIENT_ID and TROPIPAY_CLIENT_SECRET environment variables")
		fmt.Println("Example (PowerShell):")
		fmt.Println("$env:TROPIPAY_CLIENT_ID='your_id'; $env:TROPIPAY_CLIENT_SECRET='your_secret'; go run ./cmd/demo/main.go")
		os.Exit(1)
	}

	fmt.Println("Initializing Tropipay Client (Sandbox)...")
	// Initialize the client
	// We use SandboxEnv for the test
	client := gotropipay.NewClient(clientID, clientSecret, gotropipay.WithEnvironment(gotropipay.SandboxEnv))

	ctx := context.Background()

	// 1. Test Authentication (Implicitly tested by the first request, but let's try a simple read)
	fmt.Println("Listing Payment Cards...")
	cards, err := client.ListPaymentCards(ctx)
	if err != nil {
		log.Fatalf("Error listing cards: %v", err)
	}

	fmt.Printf("Successfully retrieved %d cards\n", len(cards))
	for _, card := range cards {
		fmt.Printf("- %s (Amount: %d %s, ShortURL: %s)\n", card.Concept, card.Amount, card.Currency, card.ShortURL)
	}
}
