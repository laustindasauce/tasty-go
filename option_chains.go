package tasty

import (
	"fmt"
	"net/http"
	"net/url"
)

// Returns a futures option chain given a futures product code, i.e. ES.
func (c *Client) GetFuturesOptionChains(productCode string) ([]FutureOption, *http.Response, error) {
	path := fmt.Sprintf("/futures-option-chains/%s", productCode)

	type instrumentResponse struct {
		Data struct {
			FutureOptions []FutureOption `json:"items"`
		} `json:"data"`
	}

	instrumentRes := new(instrumentResponse)

	resp, err := c.request(http.MethodGet, path, nil, nil, instrumentRes)
	if err != nil {
		return []FutureOption{}, resp, err
	}

	return instrumentRes.Data.FutureOptions, resp, nil
}

// Returns a futures option chain given a futures product code in a nested form to minimize
// redundant processing.
func (c *Client) GetNestedFuturesOptionChains(productCode string) (NestedFuturesOptionChains, *http.Response, error) {
	path := fmt.Sprintf("/futures-option-chains/%s/nested", productCode)

	type instrumentResponse struct {
		Chains NestedFuturesOptionChains `json:"data"`
	}

	instrumentRes := new(instrumentResponse)

	resp, err := c.request(http.MethodGet, path, nil, nil, instrumentRes)
	if err != nil {
		return NestedFuturesOptionChains{}, resp, err
	}

	return instrumentRes.Chains, resp, nil
}

// Returns an option chain given an underlying symbol, i.e. AAPL.
func (c *Client) GetEquityOptionChains(symbol string) ([]EquityOption, *http.Response, error) {
	// url escape required for instances where "/" exists in symbol i.e. BRK/B
	symbol = url.PathEscape(symbol)

	path := fmt.Sprintf("/option-chains/%s", symbol)

	type instrumentResponse struct {
		Data struct {
			EquityOptions []EquityOption `json:"items"`
		} `json:"data"`
	}

	instrumentRes := new(instrumentResponse)

	// customRequest required for instances where "/" exists in symbol i.e. BRK/B
	resp, err := c.customRequest(http.MethodGet, path, nil, nil, instrumentRes)
	if err != nil {
		return []EquityOption{}, resp, err
	}

	return instrumentRes.Data.EquityOptions, resp, nil
}

// Returns an option chain given an underlying symbol,
// i.e. AAPL in a nested form to minimize redundant processing.
func (c *Client) GetNestedEquityOptionChains(symbol string) ([]NestedOptionChains, *http.Response, error) {
	// url escape required for instances where "/" exists in symbol i.e. BRK/B
	symbol = url.PathEscape(symbol)

	path := fmt.Sprintf("/option-chains/%s/nested", symbol)

	type instrumentResponse struct {
		Data struct {
			Chains []NestedOptionChains `json:"items"`
		} `json:"data"`
	}

	instrumentRes := new(instrumentResponse)

	// customRequest required for instances where "/" exists in symbol i.e. BRK/B
	resp, err := c.customRequest(http.MethodGet, path, nil, nil, instrumentRes)
	if err != nil {
		return []NestedOptionChains{}, resp, err
	}

	return instrumentRes.Data.Chains, resp, nil
}

// Returns an option chain given an underlying symbol,
// i.e. AAPL in a compact form to minimize content size.
func (c *Client) GetCompactEquityOptionChains(symbol string) ([]CompactOptionChains, *http.Response, error) {
	// url escape required for instances where "/" exists in symbol i.e. BRK/B
	symbol = url.PathEscape(symbol)

	path := fmt.Sprintf("/option-chains/%s/compact", symbol)

	type instrumentResponse struct {
		Data struct {
			Chains []CompactOptionChains `json:"items"`
		} `json:"data"`
	}

	instrumentRes := new(instrumentResponse)

	// customRequest required for instances where "/" exists in symbol i.e. BRK/B
	resp, err := c.customRequest(http.MethodGet, path, nil, nil, instrumentRes)
	if err != nil {
		return []CompactOptionChains{}, resp, err
	}

	return instrumentRes.Data.Chains, resp, nil
}
