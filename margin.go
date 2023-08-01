package tasty

import (
	"fmt"
	"net/http"
)

// Fetch current margin/capital requirements report for an account.
func (c *Client) GetMarginRequirements(accountNumber string) (MarginRequirements, *http.Response, error) {
	path := fmt.Sprintf("/margin/accounts/%s/requirements", accountNumber)

	type marginResponse struct {
		MarginRequirements MarginRequirements `json:"data"`
	}

	marginRes := new(marginResponse)

	resp, err := c.request(http.MethodGet, path, nil, nil, marginRes)
	if err != nil {
		return MarginRequirements{}, resp, err
	}

	return marginRes.MarginRequirements, resp, nil
}

// Get effective margin requirements for account.
func (c *Client) GetEffectiveMarginRequirements(accountNumber, underlyingSymbol string) (EffectiveMarginRequirements, *http.Response, error) {
	path := fmt.Sprintf("/accounts/%s/margin-requirements/%s/effective", accountNumber, underlyingSymbol)

	type marginResponse struct {
		EffectiveMarginRequirements EffectiveMarginRequirements `json:"data"`
	}

	marginRes := new(marginResponse)

	resp, err := c.request(http.MethodGet, path, nil, nil, marginRes)
	if err != nil {
		return EffectiveMarginRequirements{}, resp, err
	}

	return marginRes.EffectiveMarginRequirements, resp, nil
}

// Publicly accessible, read only margin configuration.
func (c *Client) GetMarginRequirementsPublicConfiguration() (MarginRequirementsGlobalConfiguration,
	*http.Response, error) {
	path := "/margin-requirements-public-configuration"

	type marginResponse struct {
		MarginReq MarginRequirementsGlobalConfiguration `json:"data"`
	}

	marginRes := new(marginResponse)

	resp, err := c.noAuthRequest(http.MethodGet, path, nil, nil, nil, marginRes)
	if err != nil {
		return MarginRequirementsGlobalConfiguration{}, resp, err
	}

	return marginRes.MarginReq, resp, nil
}
