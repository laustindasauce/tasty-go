package tasty

import (
	"fmt"
	"net/http"

	"github.com/austinbspencer/tasty-go/models"
)

// Get authenticated customer
func (c *Client) GetMyCustomerInfo() (models.Customer, *Error) {
	if c.Session.SessionToken == nil {
		return models.Customer{}, &Error{Message: "Session is invalid: Session Token cannot be nil."}
	}
	path := "/customers/me"

	type customerResponse struct {
		Customer models.Customer `json:"data"`
	}

	customersRes := new(customerResponse)

	header := http.Header{}
	header.Add("Authorization", *c.Session.SessionToken)

	err := c.request(http.MethodGet, path, header, nil, nil, customersRes)
	if err != nil {
		return models.Customer{}, err
	}

	return customersRes.Customer, nil
}

// Get a full customer resource.
func (c *Client) GetCustomer(customerID string) (models.Customer, *Error) {
	if c.Session.SessionToken == nil {
		return models.Customer{}, &Error{Message: "Session is invalid: Session Token cannot be nil."}
	}
	path := fmt.Sprintf("/customers/%s", customerID)

	type customerResponse struct {
		Customer models.Customer `json:"data"`
	}

	customersRes := new(customerResponse)

	header := http.Header{}
	header.Add("Authorization", *c.Session.SessionToken)

	err := c.request(http.MethodGet, path, header, nil, nil, customersRes)
	if err != nil {
		return models.Customer{}, err
	}

	return customersRes.Customer, nil
}

// Get a full customer resource.
func (c *Client) GetCustomerAccounts(customerID string) ([]models.Account, *Error) {
	if c.Session.SessionToken == nil {
		return []models.Account{}, &Error{Message: "Session is invalid: Session Token cannot be nil."}
	}
	path := fmt.Sprintf("/customers/%s/accounts", customerID)

	type customerResponse struct {
		Data struct {
			Items []struct {
				Account models.Account `json:"account"`
			} `json:"items"`
		} `json:"data"`
	}

	customersRes := new(customerResponse)

	header := http.Header{}
	header.Add("Authorization", *c.Session.SessionToken)

	err := c.request(http.MethodGet, path, header, nil, nil, customersRes)
	if err != nil {
		return []models.Account{}, err
	}

	var accounts []models.Account

	for _, acct := range customersRes.Data.Items {
		accounts = append(accounts, acct.Account)
	}

	return accounts, nil
}

// Returns the appropriate quote streamer endpoint, level and identification token
// for the current customer to receive market data.
func (c *Client) GetQuoteStreamerTokens() (models.QuoteStreamerTokenAuthResult, *Error) {
	if c.Session.SessionToken == nil {
		return models.QuoteStreamerTokenAuthResult{}, &Error{Message: "Session is invalid: Session Token cannot be nil."}
	}
	path := "/quote-streamer-tokens"

	type customerResponse struct {
		Streamer models.QuoteStreamerTokenAuthResult `json:"data"`
	}

	customersRes := new(customerResponse)

	header := http.Header{}
	header.Add("Authorization", *c.Session.SessionToken)

	err := c.request(http.MethodGet, path, header, nil, nil, customersRes)
	if err != nil {
		return models.QuoteStreamerTokenAuthResult{}, err
	}

	return customersRes.Streamer, nil
}
