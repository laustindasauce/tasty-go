package tasty

import (
	"fmt"
	"net/http"
)

// Fetch current margin/capital requirements report for an account.
func (c *Client) GetMarginRequirements(accountNumber string) (MarginRequirements, error) {
	path := fmt.Sprintf("/margin/accounts/%s/requirements", accountNumber)

	type marginResponse struct {
		MarginRequirements MarginRequirements `json:"data"`
	}

	marginRes := new(marginResponse)

	err := c.request(http.MethodGet, path, nil, nil, marginRes)
	if err != nil {
		return MarginRequirements{}, err
	}

	return marginRes.MarginRequirements, nil
}

// Estimate margin requirements for an order given an account
// Need more understanding on the expected payload and response
// https://developer.tastytrade.com/open-api-spec/margin-requirements
// When sending request body in format shown in the docs I am receiving
// internal server error response.
func (c *Client) MarginRequirementsDryRun(accountNumber string, order NewOrder) (any, error) {
	path := fmt.Sprintf("/margin/accounts/%s/dry-run", accountNumber)

	type marginResponse struct {
		Response any `json:"data"`
	}

	marginRes := new(marginResponse)

	err := c.request(http.MethodPost, path, nil, order, marginRes)
	if err != nil {
		return nil, err
	}

	return marginRes.Response, nil
}

// Get effective margin requirements for account.
func (c *Client) GetEffectiveMarginRequirements(accountNumber, underlyingSymbol string) (EffectiveMarginRequirements, error) {
	path := fmt.Sprintf("/accounts/%s/margin-requirements/%s/effective", accountNumber, underlyingSymbol)

	type marginResponse struct {
		EffectiveMarginRequirements EffectiveMarginRequirements `json:"data"`
	}

	marginRes := new(marginResponse)

	err := c.request(http.MethodGet, path, nil, nil, marginRes)
	if err != nil {
		return EffectiveMarginRequirements{}, err
	}

	return marginRes.EffectiveMarginRequirements, nil
}
