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

### 1. Payment Cards (Paylinks)

Manage payment links for your customers.

```go
// Create a new payment link
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
    log.Fatalf("Error creating link: %v", err)
}
fmt.Printf("Payment Link: %s\n", card.PaymentURL)

// List existing cards
cards, _ := client.ListPaymentCards(ctx)
for _, c := range cards {
    fmt.Printf("Card %s: %s\n", c.Reference, c.ShortURL)
}
```

### 2. User Management

Access profile data and manage security settings.

```go
// Get Profile
user, _ := client.GetUserProfile(ctx)
fmt.Printf("User: %s %s\n", user.Name, user.Surname)

// Send Security Code (SMS/Email)
err := client.SendSecurityCode(ctx, gotropipay.SendSecurityCodeRequest{
    Type: "email",
})
```

### 3. Deposit Accounts (Beneficiaries)

Manage external bank accounts for transfers.

```go
// List beneficiaries
accounts, _ := client.ListDepositAccounts(ctx, 10, 0, "")
for _, acc := range accounts {
    fmt.Printf("Beneficiary: %s %s (%s)\n", acc.FirstName, acc.LastName, acc.AccountNumber)
}

// Validate an account number
valResp, _ := client.ValidateAccountNumber(ctx, gotropipay.ValidateAccountNumberRequest{
    AccountNumber:        "ES9121000418450200051332",
    CountryDestinationID: 1, // Spain
    Currency:             "EUR",
})
fmt.Printf("Is Valid: %v\n", valResp.Valid)
```

### 4. Movements (Transactions)

**Standard List (REST)**

```go
filter := &gotropipay.MovementFilter{
    State:     []string{"completed"},
    Currency:  "EUR",
    AmountGte: 1000,
}

resp, _ := client.ListMovements(ctx, 20, 0, filter)
for _, m := range resp.Items {
    fmt.Printf("Movement: %d %s | Ref: %s\n", m.Amount, m.Currency, m.Reference)
}
```

**Advanced Search (GraphQL)**

Ideal for complex queries, filtering by nested fields, or retrieving detailed sender/recipient info.

```go
gqlFilter := &gotropipay.MovementFilter{
    PaymentMethod: []string{"CARD"}, // Filter by payment method
}

// Fetch advanced details
gqlResp, err := client.SearchMovements(ctx, gqlFilter, 10, 0)
if err != nil {
    log.Fatal(err)
}

for _, m := range gqlResp.Items {
    fmt.Printf("ID: %v, Amount: %d %s, Sender: %s\n", m.ID, m.Amount.Value, m.Amount.Currency, m.Sender.Name)
}
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
