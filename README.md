# Go Tropipay SDK

A Go wrapper for the Tropipay Payments API.

## Installation

```bash
go get github.com/yosle/gotropipay
```

## Usage

```go
package main

import (
	"context"
	"fmt"
	"log"
	
	"github.com/yosle/gotropipay"
)

func main() {
	// Initialize Client
	client := gotropipay.NewClient(
		"YOUR_CLIENT_ID",
		"YOUR_CLIENT_SECRET",
		gotropipay.WithEnvironment(gotropipay.SandboxEnv), // or gotropipay.ProductionEnv
	)

	// Context with timeout
	ctx := context.Background()

	// List Payment Cards
	cards, err := client.ListPaymentCards(ctx)
	if err != nil {
		log.Fatalf("Error listing cards: %v", err)
	}

	for _, card := range cards {
		fmt.Printf("Card: %s - %s\n", card.ID, card.Last4)
	}
}
```

## Features

- **Authentication**: Automatic OAuth2 Token retrieval and caching.
- **Environments**: Support for Sandbox and Production.
- **Context support**: All methods support `context.Context`.
- **Resources**:
  - Payment Cards (CRUD)

## License

MIT
