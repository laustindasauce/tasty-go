package tasty

import (
	"fmt"
	"net/http"
)

// Fetch current margin/capital requirements report for an account.
func (c *Client) GetMarginRequirements(accountNumber string) (MarginRequirements, *Error) {
	if c.Session.SessionToken == nil {
		return MarginRequirements{}, &Error{Message: "Session is invalid: Session Token cannot be nil."}
	}

	path := fmt.Sprintf("/margin/accounts/%s/requirements", accountNumber)

	type accountResponse struct {
		MarginRequirements MarginRequirements `json:"data"`
	}

	marginRes := new(accountResponse)

	header := http.Header{}
	header.Add("Authorization", *c.Session.SessionToken)

	err := c.request(http.MethodGet, path, header, nil, nil, marginRes)
	if err != nil {
		return MarginRequirements{}, err
	}

	return marginRes.MarginRequirements, nil
}

// Estimate margin requirements for an order given an account
// This is not functional at the moment
// Need more understanding on the expected payload
// https://developer.tastytrade.com/open-api-spec/margin-requirements
func (c *Client) MarginRequirementsDryRun(accountNumber string, order NewOrder) (any, *Error) {
	if c.Session.SessionToken == nil {
		return nil, &Error{Message: "Session is invalid: Session Token cannot be nil."}
	}

	path := fmt.Sprintf("/margin/accounts/%s/requirements", accountNumber)

	type accountResponse struct {
		Response any `json:"data"`
	}

	marginRes := new(accountResponse)

	header := http.Header{}
	header.Add("Authorization", *c.Session.SessionToken)

	err := c.request(http.MethodPost, path, header, nil, order, marginRes)
	if err != nil {
		return nil, err
	}

	return marginRes.Response, nil
}

// Get effective margin requirements for account.
func (c *Client) GetEffectiveMarginRequirements(accountNumber, underlyingSymbol string) (EffectiveMarginRequirements, *Error) {
	if c.Session.SessionToken == nil {
		return EffectiveMarginRequirements{}, &Error{Message: "Session is invalid: Session Token cannot be nil."}
	}

	path := fmt.Sprintf("/accounts/%s/margin-requirements/%s/effective", accountNumber, underlyingSymbol)

	type accountResponse struct {
		EffectiveMarginRequirements EffectiveMarginRequirements `json:"data"`
	}

	marginRes := new(accountResponse)

	header := http.Header{}
	header.Add("Authorization", *c.Session.SessionToken)

	err := c.request(http.MethodGet, path, header, nil, nil, marginRes)
	if err != nil {
		return EffectiveMarginRequirements{}, err
	}

	return marginRes.EffectiveMarginRequirements, nil
}
