package tasty

import (
	"fmt"
	"net/http"
	"net/url"
)

// Returns all active equities in a paginated fashion.
func (c *Client) GetActiveEquities(query ActiveEquitiesQuery) ([]Equity, Pagination, *Error) {
	path := "/instruments/equities/active"

	type instrumentResponse struct {
		Data struct {
			ActiveEquities []Equity `json:"items"`
		} `json:"data"`
		Pagination Pagination `json:"pagination"`
	}

	instrumentRes := new(instrumentResponse)

	err := c.request(http.MethodGet, path, query, nil, instrumentRes)
	if err != nil {
		return []Equity{}, Pagination{}, err
	}

	return instrumentRes.Data.ActiveEquities, instrumentRes.Pagination, nil
}

// Returns a set of equity definitions given an array of one or more symbols.
func (c *Client) GetEquities(query EquitiesQuery) ([]Equity, *Error) {
	path := "/instruments/equities"

	type instrumentResponse struct {
		Data struct {
			Equities []Equity `json:"items"`
		} `json:"data"`
	}

	instrumentRes := new(instrumentResponse)

	err := c.request(http.MethodGet, path, query, nil, instrumentRes)
	if err != nil {
		return []Equity{}, err
	}

	return instrumentRes.Data.Equities, nil
}

// Returns a single equity definition for the provided symbol.
func (c *Client) GetEquity(symbol string) (Equity, *Error) {
	// url escape required for instances where "/" exists in symbol i.e. BRK/B
	path := fmt.Sprintf("/instruments/equities/%s", url.PathEscape(symbol))

	type instrumentResponse struct {
		Equity Equity `json:"data"`
	}

	instrumentRes := new(instrumentResponse)

	// customRequest required for instances where "/" exists in symbol i.e. BRK/B
	err := c.customRequest(http.MethodGet, path, nil, nil, instrumentRes)
	if err != nil {
		return Equity{}, err
	}

	return instrumentRes.Equity, nil
}

// Returns a set of equity options given one or more symbols.
func (c *Client) GetEquityOptions(query EquityOptionsQuery) ([]EquityOption, *Error) {
	path := "/instruments/equity-options"

	type instrumentResponse struct {
		Data struct {
			EquityOptions []EquityOption `json:"items"`
		} `json:"data"`
	}

	instrumentRes := new(instrumentResponse)

	err := c.request(http.MethodGet, path, query, nil, instrumentRes)
	if err != nil {
		return []EquityOption{}, err
	}

	return instrumentRes.Data.EquityOptions, nil
}

// Returns a set of equity options given one or more symbols.
func (c *Client) GetEquityOption(sym EquityOptionsSymbology, active bool) (EquityOption, *Error) {
	occSymbol := sym.Build()

	path := fmt.Sprintf("/instruments/equity-options/%s", occSymbol)

	type instrumentResponse struct {
		EquityOption EquityOption `json:"data"`
	}

	type activeQuery struct {
		// Whether an option is available for trading with the broker.
		Active bool `url:"active"`
	}

	query := activeQuery{Active: active}

	instrumentRes := new(instrumentResponse)

	err := c.request(http.MethodGet, path, query, nil, instrumentRes)
	if err != nil {
		return EquityOption{}, err
	}

	return instrumentRes.EquityOption, nil
}
