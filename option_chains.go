package tasty

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/austinbspencer/tasty-go/models"
)

// Returns a futures option chain given a futures product code, i.e. ES.
func (c *Client) GetFuturesOptionChains(productCode string) ([]models.FutureOption, *Error) {
	if c.Session.SessionToken == nil {
		return []models.FutureOption{}, &Error{Message: "Session is invalid: Session Token cannot be nil."}
	}

	reqURL := fmt.Sprintf("/futures-option-chains/%s", productCode)

	type instrumentResponse struct {
		Data struct {
			FutureOptions []models.FutureOption `json:"items"`
		} `json:"data"`
	}

	instrumentRes := new(instrumentResponse)

	header := http.Header{}
	header.Add("Authorization", *c.Session.SessionToken)

	err := c.request(http.MethodGet, reqURL, header, nil, nil, instrumentRes)
	if err != nil {
		return []models.FutureOption{}, err
	}

	return instrumentRes.Data.FutureOptions, nil
}

// Returns a futures option chain given a futures product code in a nested form to minimize
// redundant processing.
func (c *Client) GetNestedFuturesOptionChains(productCode string) (models.NestedFuturesOptionChains, *Error) {
	if c.Session.SessionToken == nil {
		return models.NestedFuturesOptionChains{}, &Error{Message: "Session is invalid: Session Token cannot be nil."}
	}

	reqURL := fmt.Sprintf("/futures-option-chains/%s/nested", productCode)

	type instrumentResponse struct {
		Chains models.NestedFuturesOptionChains `json:"data"`
	}

	instrumentRes := new(instrumentResponse)

	header := http.Header{}
	header.Add("Authorization", *c.Session.SessionToken)

	err := c.request(http.MethodGet, reqURL, header, nil, nil, instrumentRes)
	if err != nil {
		return models.NestedFuturesOptionChains{}, err
	}

	return instrumentRes.Chains, nil
}

// Returns an option chain given an underlying symbol, i.e. AAPL.
func (c *Client) GetEquityOptionChains(symbol string) ([]models.EquityOption, *Error) {
	if c.Session.SessionToken == nil {
		return []models.EquityOption{}, &Error{Message: "Session is invalid: Session Token cannot be nil."}
	}

	// url escape required for instances where "/" exists in symbol i.e. BRK/B
	symbol = url.PathEscape(symbol)

	path := fmt.Sprintf("/option-chains/%s", symbol)

	type instrumentResponse struct {
		Data struct {
			EquityOptions []models.EquityOption `json:"items"`
		} `json:"data"`
	}

	instrumentRes := new(instrumentResponse)

	header := http.Header{}
	header.Add("Authorization", *c.Session.SessionToken)

	// customRequest required for instances where "/" exists in symbol i.e. BRK/B
	err := c.customRequest(http.MethodGet, path, header, nil, nil, instrumentRes)
	if err != nil {
		return []models.EquityOption{}, err
	}

	return instrumentRes.Data.EquityOptions, nil
}

// Returns an option chain given an underlying symbol,
// i.e. AAPL in a nested form to minimize redundant processing.
func (c *Client) GetNestedEquityOptionChains(symbol string) ([]models.NestedOptionChains, *Error) {
	if c.Session.SessionToken == nil {
		return []models.NestedOptionChains{}, &Error{Message: "Session is invalid: Session Token cannot be nil."}
	}

	// url escape required for instances where "/" exists in symbol i.e. BRK/B
	symbol = url.PathEscape(symbol)

	path := fmt.Sprintf("/option-chains/%s/nested", symbol)

	type instrumentResponse struct {
		Data struct {
			Chains []models.NestedOptionChains `json:"items"`
		} `json:"data"`
	}

	instrumentRes := new(instrumentResponse)

	header := http.Header{}
	header.Add("Authorization", *c.Session.SessionToken)

	// customRequest required for instances where "/" exists in symbol i.e. BRK/B
	err := c.customRequest(http.MethodGet, path, header, nil, nil, instrumentRes)
	if err != nil {
		return []models.NestedOptionChains{}, err
	}

	return instrumentRes.Data.Chains, nil
}

// Returns an option chain given an underlying symbol,
// i.e. AAPL in a compact form to minimize content size.
func (c *Client) GetCompactEquityOptionChains(symbol string) ([]models.CompactOptionChains, *Error) {
	if c.Session.SessionToken == nil {
		return []models.CompactOptionChains{}, &Error{Message: "Session is invalid: Session Token cannot be nil."}
	}

	// url escape required for instances where "/" exists in symbol i.e. BRK/B
	symbol = url.PathEscape(symbol)

	path := fmt.Sprintf("/option-chains/%s/compact", symbol)

	type instrumentResponse struct {
		Data struct {
			Chains []models.CompactOptionChains `json:"items"`
		} `json:"data"`
	}

	instrumentRes := new(instrumentResponse)

	header := http.Header{}
	header.Add("Authorization", *c.Session.SessionToken)

	// customRequest required for instances where "/" exists in symbol i.e. BRK/B
	err := c.customRequest(http.MethodGet, path, header, nil, nil, instrumentRes)
	if err != nil {
		return []models.CompactOptionChains{}, err
	}

	return instrumentRes.Data.Chains, nil
}
