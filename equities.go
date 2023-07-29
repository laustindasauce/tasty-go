package tasty

import (
	"fmt"
	"net/http"
	"net/url"
)

// Returns all active equities in a paginated fashion.
func (c *Client) GetActiveEquities(query ActiveEquitiesQuery) ([]Equity, Pagination, *http.Response, error) {
	path := "/instruments/equities/active"

	type instrumentResponse struct {
		Data struct {
			ActiveEquities []Equity `json:"items"`
		} `json:"data"`
		Pagination Pagination `json:"pagination"`
	}

	instrumentRes := new(instrumentResponse)

	resp, err := c.request(http.MethodGet, path, query, nil, instrumentRes)
	if err != nil {
		return []Equity{}, Pagination{}, resp, err
	}

	return instrumentRes.Data.ActiveEquities, instrumentRes.Pagination, resp, nil
}

// Returns a set of equity definitions given an array of one or more symbols.
func (c *Client) GetEquities(query EquitiesQuery) ([]Equity, *http.Response, error) {
	path := "/instruments/equities"

	type instrumentResponse struct {
		Data struct {
			Equities []Equity `json:"items"`
		} `json:"data"`
	}

	instrumentRes := new(instrumentResponse)

	resp, err := c.request(http.MethodGet, path, query, nil, instrumentRes)
	if err != nil {
		return []Equity{}, resp, err
	}

	return instrumentRes.Data.Equities, resp, nil
}

// Returns a single equity definition for the provided symbol.
func (c *Client) GetEquity(symbol string) (Equity, *http.Response, error) {
	// url escape required for instances where "/" exists in symbol i.e. BRK/B
	path := fmt.Sprintf("/instruments/equities/%s", url.PathEscape(symbol))

	type instrumentResponse struct {
		Equity Equity `json:"data"`
	}

	instrumentRes := new(instrumentResponse)

	// customRequest required for instances where "/" exists in symbol i.e. BRK/B
	resp, err := c.customRequest(http.MethodGet, path, nil, nil, instrumentRes)
	if err != nil {
		return Equity{}, resp, err
	}

	return instrumentRes.Equity, resp, nil
}

// Returns a set of equity options given one or more symbols.
func (c *Client) GetEquityOptions(query EquityOptionsQuery) ([]EquityOption, *http.Response, error) {
	path := "/instruments/equity-options"

	type instrumentResponse struct {
		Data struct {
			EquityOptions []EquityOption `json:"items"`
		} `json:"data"`
	}

	instrumentRes := new(instrumentResponse)

	resp, err := c.request(http.MethodGet, path, query, nil, instrumentRes)
	if err != nil {
		return []EquityOption{}, resp, err
	}

	return instrumentRes.Data.EquityOptions, resp, nil
}

// Returns a set of equity options given one or more symbols.
func (c *Client) GetEquityOption(sym EquityOptionsSymbology, active bool) (EquityOption, *http.Response, error) {
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

	resp, err := c.request(http.MethodGet, path, query, nil, instrumentRes)
	if err != nil {
		return EquityOption{}, resp, err
	}

	return instrumentRes.EquityOption, resp, nil
}
