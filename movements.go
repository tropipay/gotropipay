package gotropipay

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

// MovementState represents the state of a movement
type MovementState string

const (
	MovementStatePending   MovementState = "pending"
	MovementStateCompleted MovementState = "completed"
	MovementStateFailed    MovementState = "failed"
	MovementStateCancelled MovementState = "cancelled"
)

// Movement represents a transaction or movement record
type Movement struct {
	ID            interface{} `json:"id"` // Can be int or string depending on the endpoint (REST vs GraphQL)
	Amount        int64       `json:"amount"`
	Currency      string      `json:"currency"`
	State         string      `json:"state"` // using string instead of MovementState to be flexible with casing
	Reference     string      `json:"reference"`
	CreatedAt     string      `json:"createdAt"`
	CompletedAt   string      `json:"completedAt"`
	BalanceBefore int64       `json:"balanceBefore"`
	BalanceAfter  int64       `json:"balanceAfter"`
	Recipient     *User       `json:"recipient,omitempty"` // Populated in GraphQL
	Sender        *User       `json:"sender,omitempty"`    // Populated in GraphQL
	Account       interface{} `json:"account,omitempty"`   // Populated in GraphQL generally
}

// MovementFilter represents the filter criteria for listing movements
type MovementFilter struct {
	State         []string `json:"state,omitempty"`
	Currency      string   `json:"currency,omitempty"`
	AmountGte     int64    `json:"amountGte,omitempty"`
	AmountLte     int64    `json:"amountLte,omitempty"`
	CreatedAtFrom string   `json:"createdAtFrom,omitempty"`
	CreatedAtTo   string   `json:"createdAtTo,omitempty"`
	Reference     string   `json:"reference,omitempty"`
	AccountID     string   `json:"accountId,omitempty"` // For GraphQL filter
}

// ListMovementsResponse is the response structure for listing movements
type ListMovementsResponse struct {
	Items      []Movement `json:"items"`
	TotalCount int        `json:"totalCount"`
	HasMore    bool       `json:"hasMore"`
}

// GraphQLRequest structures
type graphQLRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables,omitempty"`
}

// REST Endpoints

// ListMovements retrieves a list of movements for the authenticated user
func (c *Client) ListMovements(ctx context.Context, limit, offset int, filter *MovementFilter) (*ListMovementsResponse, error) {
	return c.listMovementsCommon(ctx, "/movements/", limit, offset, filter)
}

// ListAccountMovements retrieves movements for a specific account
func (c *Client) ListAccountMovements(ctx context.Context, accountID string, limit, offset int, filter *MovementFilter) (*ListMovementsResponse, error) {
	path := fmt.Sprintf("/accounts/%s/movements", accountID)
	return c.listMovementsCommon(ctx, path, limit, offset, filter)
}

func (c *Client) listMovementsCommon(ctx context.Context, path string, limit, offset int, filter *MovementFilter) (*ListMovementsResponse, error) {
	params := url.Values{}
	if limit > 0 {
		params.Add("limit", strconv.Itoa(limit))
	}
	if offset > 0 {
		params.Add("offset", strconv.Itoa(offset))
	}
	if filter != nil {
		filterJSON, err := json.Marshal(filter)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal filter: %w", err)
		}
		params.Add("query", string(filterJSON))
	}

	if len(params) > 0 {
		path += "?" + params.Encode()
	}

	var resp ListMovementsResponse
	err := c.Request(ctx, "GET", path, nil, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// GraphQL Business Endpoint

// SearchMovements performs an advanced search using the GraphQL endpoint
// SearchMovements performs an advanced search using the GraphQL endpoint
func (c *Client) SearchMovements(ctx context.Context, filter *MovementFilter, limit, offset int) (*ListMovementsResponse, error) {
	query := `query GetMovements($filter: MovementFilter, $pagination: PaginationInput) { 
		movements(filter: $filter, pagination: $pagination) { 
			items { 
				id 
				reference
				concept
				state
				createdAt 
				completedAt 
				amount {
					value
					currency
				}
				sender
				recipient
				movementDetail {
					senderData {
						name
						email
					}
					recipientData {
						name
						account
					}
				}
			} 
			totalCount 
		} 
	}`

	vars := map[string]interface{}{
		"filter": filter,
		"pagination": map[string]int{
			"limit":  limit,
			"offset": offset,
		},
	}

	req := graphQLRequest{
		Query:     query,
		Variables: vars,
	}

	// Internal structs for GraphQL response
	type gqlAmount struct {
		Value    int64  `json:"value"`
		Currency string `json:"currency"`
	}
	type gqlPersonaData struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	type gqlRecipientData struct {
		Name    string `json:"name"`
		Account string `json:"account"`
	}
	type gqlMovementDetail struct {
		SenderData    gqlPersonaData   `json:"senderData"`
		RecipientData gqlRecipientData `json:"recipientData"`
	}
	type gqlMovement struct {
		ID             interface{}       `json:"id"`
		Reference      string            `json:"reference"`
		Concept        string            `json:"concept"`
		State          string            `json:"state"`
		CreatedAt      string            `json:"createdAt"`
		CompletedAt    string            `json:"completedAt"`
		Amount         gqlAmount         `json:"amount"`
		Sender         string            `json:"sender"`
		Recipient      string            `json:"recipient"`
		MovementDetail gqlMovementDetail `json:"movementDetail"`
	}

	var gqlResp struct {
		Data struct {
			Movements struct {
				Items      []gqlMovement `json:"items"`
				TotalCount int           `json:"totalCount"`
			} `json:"movements"`
		} `json:"data"`
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors"`
	}

	err := c.Request(ctx, "POST", "/movements/business", req, &gqlResp)
	if err != nil {
		return nil, err
	}

	if len(gqlResp.Errors) > 0 {
		return nil, fmt.Errorf("graphql error: %s", gqlResp.Errors[0].Message)
	}

	// Map back to standard Movement struct
	var movements []Movement
	for _, item := range gqlResp.Data.Movements.Items {
		// Construct User objects from details if available, otherwise just use names
		sender := &User{Name: item.Sender}
		if item.MovementDetail.SenderData.Name != "" {
			sender.Name = item.MovementDetail.SenderData.Name
			sender.Email = item.MovementDetail.SenderData.Email
		}

		recipient := &User{Name: item.Recipient}
		if item.MovementDetail.RecipientData.Name != "" {
			recipient.Name = item.MovementDetail.RecipientData.Name
			// item.MovementDetail.RecipientData.Account contains account number usually
		}

		movements = append(movements, Movement{
			ID:          item.ID,
			Amount:      item.Amount.Value,
			Currency:    item.Amount.Currency,
			State:       item.State,
			Reference:   item.Reference,
			CreatedAt:   item.CreatedAt,
			CompletedAt: item.CompletedAt,
			// Balance fields are not in the standard GQL item response unless added, skipping for now
			Sender:    sender,
			Recipient: recipient,
		})
	}

	return &ListMovementsResponse{
		Items:      movements,
		TotalCount: gqlResp.Data.Movements.TotalCount,
		HasMore:    len(movements) < gqlResp.Data.Movements.TotalCount, // Approximate
	}, nil
}
