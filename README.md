# Go Tropipay SDK

A comprehensive, idiomatic Go wrapper for the [Tropipay](https://tropipay.com) Payments API. This SDK simplifies the integration of Tropipay's financial services into your Go applications, handling authentication, request signing, and data modeling.

## Features

*   **Robust Authentication**: Automatic OAuth2 token retrieval, caching, and refreshing.
*   **Environment Support**: Easy switching between Sandbox and Production environments.
*   **Context Aware**: All API operations support `context.Context` for cancellation and timeouts.
*   **Comprehensive Coverage**:
    *   **Users**: specific profile management, security codes, 2FA configuration.
    *   **Payment Cards (Links)**: Create, list, delete, and manage payment links/cards.
    *   **Accounts**: Link Tropicards and retrieve crypto deposit addresses.
    *   **Beneficiaries (Deposit Accounts)**: Manage recipients for transfers.
    *   **Movements**: Full transaction history with advanced filtering (REST & GraphQL support).

## Installation

```bash
go get github.com/tropipay/gotropipay
```

## Getting Started

Initialize the client with your credentials. We recommend using environment variables to store secrets.

```go
package main

import (
    "context"
    "log"
    "os"

    "github.com/tropipay/gotropipay"
)

func main() {
    clientID := os.Getenv("TROPIPAY_CLIENT_ID")
    clientSecret := os.Getenv("TROPIPAY_CLIENT_SECRET")

    // Initialize Client in Sandbox mode
    client := gotropipay.NewClient(
        clientID,
        clientSecret,
        gotropipay.WithEnvironment(gotropipay.SandboxEnv), // Remove for Production
    )

    // Verify connection by getting user profile
    user, err := client.GetUserProfile(context.Background())
    if err != nil {
        log.Fatalf("Failed to get profile: %v", err)
    }

    log.Printf("Connected as: %s %s (Balance: %d cents)", user.Name, user.Surname, user.Balance)
}
```

## Usage Examples

### 1. Creating a Payment Link (PaymentCard)

Create a payment URL to share with a customer.

```go
ctx := context.Background()

req := gotropipay.CreatePaymentCardRequest{
    Reference:   "ORDER-1234",
    Concept:     "Product Purchase",
    Amount:      1500, // 15.00 EUR
    Currency:    "EUR",
    Description: "Payment for Order #1234",
    SingleUse:   true,
}

card, err := client.CreatePaymentCard(ctx, req)
if err != nil {
    log.Printf("Error creating link: %v", err)
    return
}

fmt.Printf("Payment Link Created: %s\n", card.PaymentURL)
```

### 2. Listing Movements with Filters

Retrieve transaction history with specific criteria.

```go
filter := &gotropipay.MovementFilter{
    State:     []string{"completed"},
    Currency:  "EUR",
    AmountGte: 1000,
}

// Get the first 20 completed EUR transactions over 10.00
resp, err := client.ListMovements(ctx, 20, 0, filter)
if err != nil {
    log.Printf("Error fetching movements: %v", err)
    return
}

for _, m := range resp.Items {
    fmt.Printf("[%s] %d %s - Ref: %s\n", m.CreatedAt, m.Amount, m.Currency, m.Reference)
}
```

### 3. Managing Beneficiaries

Add a new bank account for transfers.

```go
newBeneficiary := gotropipay.CreateDepositAccountRequest{
    FirstName:            "Jane",
    LastName:             "Doe",
    AccountNumber:        "ES9121000418450200051332",
    CountryDestinationID: 1, // Spain
    Type:                 0,
    Alias:                "Jane Primary",
}

data, err := client.CreateDepositAccount(ctx, newBeneficiary)
if err != nil {
    log.Println("Could not create beneficiary:", err)
    return
}
fmt.Printf("Beneficiary Created ID: %d\n", data.ID)
```

### 4. Advanced: GraphQL Movement Search

For complex queries involving nested fields.

```go
// Search specifically using the GraphQL endpoint
gqlResp, err := client.SearchMovements(ctx, filter, 10, 0)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Total matches found: %d\n", gqlResp.TotalCount)
```

## Best Practices

### Context and Timeouts
Always use `context.Context` to manage request lifecycles. This is crucial for production applications to handle network latency or cancellations gracefully.

```go
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

client.GetUserProfile(ctx)
```

### Error Handling
The SDK attempts to return meaningful errors. Check for standard Go errors or wrapped API error messages.

### Security
*   **Never hardcode credentials.** Use environment variables or a secure vault.
*   **Token Management:** The SDK handles token refresh automatically. You do not need to manually manage the `Bearer` token.
*   **Sandboxing:** Always develop and test against `gotropipay.SandboxEnv` before switching to `ProductionEnv`.

## License

MIT
