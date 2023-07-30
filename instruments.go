package tasty

import (
	"fmt"
	"net/http"
)

// Retrieve all quantity decimal precisions.
func (c *Client) GetQuantityDecimalPrecisions() ([]QuantityDecimalPrecision, *http.Response, error) {
	path := "/instruments/quantity-decimal-precisions"

	type instrumentResponse struct {
		Data struct {
			QuantityDecimalPrecisions []QuantityDecimalPrecision `json:"items"`
		} `json:"data"`
	}

	instrumentRes := new(instrumentResponse)

	resp, err := c.request(http.MethodGet, path, nil, nil, instrumentRes)
	if err != nil {
		return []QuantityDecimalPrecision{}, resp, err
	}

	return instrumentRes.Data.QuantityDecimalPrecisions, resp, nil
}

// Returns a set of warrant definitions that can be filtered by parameters.
func (c *Client) GetWarrants(symbols []string) ([]Warrant, *http.Response, error) {
	path := "/instruments/warrants"

	type instrumentResponse struct {
		Data struct {
			Warrants []Warrant `json:"items"`
		} `json:"data"`
	}

	type symbolsQuery struct {
		// Symbols is the list of symbols
		Symbols []string `url:"symbol[]"`
	}

	query := symbolsQuery{Symbols: symbols}

	instrumentRes := new(instrumentResponse)

	resp, err := c.request(http.MethodGet, path, query, nil, instrumentRes)
	if err != nil {
		return []Warrant{}, resp, err
	}

	return instrumentRes.Data.Warrants, resp, nil
}

// Returns a single warrant definition for the provided symbol.
func (c *Client) GetWarrant(symbol string) (Warrant, *http.Response, error) {
	path := fmt.Sprintf("/instruments/warrants/%s", symbol)

	type instrumentResponse struct {
		Warrant Warrant `json:"data"`
	}

	instrumentRes := new(instrumentResponse)

	resp, err := c.request(http.MethodGet, path, nil, nil, instrumentRes)
	if err != nil {
		return Warrant{}, resp, err
	}

	return instrumentRes.Warrant, resp, nil
}
