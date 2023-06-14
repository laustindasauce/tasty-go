package tasty

import (
	"fmt"
	"net/http"

	"github.com/austinbspencer/tasty-go/models"
)

// Fetch current margin/capital requirements report for an account
func (c *Client) GetMarginRequirements(accountNumber string) (models.MarginRequirements, *Error) {
	if c.Session.SessionToken == nil {
		return models.MarginRequirements{}, &Error{Message: "Session is invalid: Session Token cannot be nil."}
	}

	path := fmt.Sprintf("/margin/accounts/%s/requirements", accountNumber)

	type accountResponse struct {
		MarginRequirements models.MarginRequirements `json:"data"`
	}

	marginRes := new(accountResponse)

	header := http.Header{}
	header.Add("Authorization", *c.Session.SessionToken)

	err := c.request(http.MethodGet, path, header, nil, nil, marginRes)
	if err != nil {
		return models.MarginRequirements{}, err
	}

	return marginRes.MarginRequirements, nil
}

// Estimate margin requirements for an order given an account
// This is not functional at the moment
// Need more understanding on the expected payload
func (c *Client) marginRequirementsDryRun(accountNumber string, order models.NewOrder) (any, *Error) {
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

// Get effective margin requirements for account
func (c *Client) GetEffectiveMarginRequirements(accountNumber, underlyingSymbol string) (models.EffectiveMarginRequirements, *Error) {
	if c.Session.SessionToken == nil {
		return models.EffectiveMarginRequirements{}, &Error{Message: "Session is invalid: Session Token cannot be nil."}
	}

	path := fmt.Sprintf("/accounts/%s/margin-requirements/%s/effective", accountNumber, underlyingSymbol)

	type accountResponse struct {
		EffectiveMarginRequirements models.EffectiveMarginRequirements `json:"data"`
	}

	marginRes := new(accountResponse)

	header := http.Header{}
	header.Add("Authorization", *c.Session.SessionToken)

	err := c.request(http.MethodGet, path, header, nil, nil, marginRes)
	if err != nil {
		return models.EffectiveMarginRequirements{}, err
	}

	return marginRes.EffectiveMarginRequirements, nil
}
