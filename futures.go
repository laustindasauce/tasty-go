package tasty

import (
	"fmt"
	"net/http"
)

// Returns a set of outright futures given an array of one or more symbols.
func (c *Client) GetFutures(query FuturesQuery) ([]Future, *http.Response, error) {
	path := "/instruments/futures"

	type instrumentResponse struct {
		Data struct {
			Futures []Future `json:"items"`
		} `json:"data"`
	}

	instrumentRes := new(instrumentResponse)

	resp, err := c.request(http.MethodGet, path, query, nil, instrumentRes)
	if err != nil {
		return []Future{}, resp, err
	}

	return instrumentRes.Data.Futures, resp, nil
}

// Returns an outright future given a symbol.
func (c *Client) GetFuture(symbol string) (Future, *http.Response, error) {
	path := fmt.Sprintf("/instruments/futures/%s", symbol)

	type instrumentResponse struct {
		Future Future `json:"data"`
	}

	instrumentRes := new(instrumentResponse)

	resp, err := c.request(http.MethodGet, path, nil, nil, instrumentRes)
	if err != nil {
		return Future{}, resp, err
	}

	return instrumentRes.Future, resp, nil
}

// Returns metadata for all supported future option products.
func (c *Client) GetFutureOptionProducts() ([]FutureOptionProduct, *http.Response, error) {
	path := "/instruments/future-option-products"

	type instrumentResponse struct {
		Data struct {
			FutureOptionProducts []FutureOptionProduct `json:"items"`
		} `json:"data"`
	}

	instrumentRes := new(instrumentResponse)

	resp, err := c.request(http.MethodGet, path, nil, nil, instrumentRes)
	if err != nil {
		return []FutureOptionProduct{}, resp, err
	}

	return instrumentRes.Data.FutureOptionProducts, resp, nil
}

// Get a future option product by exchange and root symbol.
func (c *Client) GetFutureOptionProduct(exchange, rootSymbol string) (FutureOptionProduct, *http.Response, error) {
	path := fmt.Sprintf("/instruments/future-option-products/%s/%s", exchange, rootSymbol)

	type instrumentResponse struct {
		FutureOptionProduct FutureOptionProduct `json:"data"`
	}

	instrumentRes := new(instrumentResponse)

	resp, err := c.request(http.MethodGet, path, nil, nil, instrumentRes)
	if err != nil {
		return FutureOptionProduct{}, resp, err
	}

	return instrumentRes.FutureOptionProduct, resp, nil
}

// Returns a set of future option(s) given an array of one or more symbols.
// Uses TW symbology: [./ESZ9 EW4U9 190927P2975].
func (c *Client) GetFutureOptions(query FutureOptionsQuery) ([]FutureOption, *http.Response, error) {
	path := "/instruments/future-options"

	type instrumentResponse struct {
		Data struct {
			FutureOptions []FutureOption `json:"items"`
		} `json:"data"`
	}

	instrumentRes := new(instrumentResponse)

	resp, err := c.request(http.MethodGet, path, query, nil, instrumentRes)
	if err != nil {
		return []FutureOption{}, resp, err
	}

	return instrumentRes.Data.FutureOptions, resp, nil
}

// Returns a future option given a symbol. Uses TW symbology: ./ESZ9 EW4U9 190927P2975.
func (c *Client) GetFutureOption(symbol string) (FutureOption, *http.Response, error) {
	path := fmt.Sprintf("/instruments/future-options/%s", symbol)

	type instrumentResponse struct {
		FutureOption FutureOption `json:"data"`
		Context      string       `json:"context"`
	}

	instrumentRes := new(instrumentResponse)

	resp, err := c.request(http.MethodGet, path, nil, nil, instrumentRes)
	if err != nil {
		return FutureOption{}, resp, err
	}

	return instrumentRes.FutureOption, resp, nil
}

// Returns metadata for all supported futures products.
func (c *Client) GetFutureProducts() ([]FutureProduct, *http.Response, error) {
	path := "/instruments/future-products"

	type instrumentResponse struct {
		Data struct {
			FutureProducts []FutureProduct `json:"items"`
		} `json:"data"`
	}

	instrumentRes := new(instrumentResponse)

	resp, err := c.request(http.MethodGet, path, nil, nil, instrumentRes)
	if err != nil {
		return []FutureProduct{}, resp, err
	}

	return instrumentRes.Data.FutureProducts, resp, nil
}

// Get future product from exchange and product code.
func (c *Client) GetFutureProduct(exchange Exchange, productCode string) (FutureProduct, *http.Response, error) {
	path := fmt.Sprintf("/instruments/future-products/%s/%s", exchange, productCode)

	type instrumentResponse struct {
		FutureProduct FutureProduct `json:"data"`
	}

	instrumentRes := new(instrumentResponse)

	resp, err := c.request(http.MethodGet, path, nil, nil, instrumentRes)
	if err != nil {
		return FutureProduct{}, resp, err
	}

	return instrumentRes.FutureProduct, resp, nil
}
