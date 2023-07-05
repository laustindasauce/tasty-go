package tasty

import (
	"fmt"
	"net/http"
)

// Returns a set of outright futures given an array of one or more symbols.
func (c *Client) GetFutures(query FuturesQuery) ([]Future, error) {
	path := "/instruments/futures"

	type instrumentResponse struct {
		Data struct {
			Futures []Future `json:"items"`
		} `json:"data"`
	}

	instrumentRes := new(instrumentResponse)

	err := c.request(http.MethodGet, path, query, nil, instrumentRes)
	if err != nil {
		return []Future{}, err
	}

	return instrumentRes.Data.Futures, nil
}

// Returns an outright future given a symbol.
func (c *Client) GetFuture(symbol string) (Future, error) {
	path := fmt.Sprintf("/instruments/futures/%s", symbol)

	type instrumentResponse struct {
		Future Future `json:"data"`
	}

	instrumentRes := new(instrumentResponse)

	err := c.request(http.MethodGet, path, nil, nil, instrumentRes)
	if err != nil {
		return Future{}, err
	}

	return instrumentRes.Future, nil
}

// Returns metadata for all supported future option products.
func (c *Client) GetFutureOptionProducts() ([]FutureOptionProduct, error) {
	path := "/instruments/future-option-products"

	type instrumentResponse struct {
		Data struct {
			FutureOptionProducts []FutureOptionProduct `json:"items"`
		} `json:"data"`
	}

	instrumentRes := new(instrumentResponse)

	err := c.request(http.MethodGet, path, nil, nil, instrumentRes)
	if err != nil {
		return []FutureOptionProduct{}, err
	}

	return instrumentRes.Data.FutureOptionProducts, nil
}

// Get a future option product by exchange and root symbol.
func (c *Client) GetFutureOptionProduct(exchange, rootSymbol string) (FutureOptionProduct, error) {
	path := fmt.Sprintf("/instruments/future-option-products/%s/%s", exchange, rootSymbol)

	type instrumentResponse struct {
		FutureOptionProduct FutureOptionProduct `json:"data"`
	}

	instrumentRes := new(instrumentResponse)

	err := c.request(http.MethodGet, path, nil, nil, instrumentRes)
	if err != nil {
		return FutureOptionProduct{}, err
	}

	return instrumentRes.FutureOptionProduct, nil
}

// Returns a set of future option(s) given an array of one or more symbols.
// Uses TW symbology: [./ESZ9 EW4U9 190927P2975].
func (c *Client) GetFutureOptions(query FutureOptionsQuery) ([]FutureOption, error) {
	path := "/instruments/future-options"

	type instrumentResponse struct {
		Data struct {
			FutureOptions []FutureOption `json:"items"`
		} `json:"data"`
	}

	instrumentRes := new(instrumentResponse)

	err := c.request(http.MethodGet, path, query, nil, instrumentRes)
	if err != nil {
		return []FutureOption{}, err
	}

	return instrumentRes.Data.FutureOptions, nil
}

// Returns a future option given a symbol. Uses TW symbology: ./ESZ9 EW4U9 190927P2975.
func (c *Client) GetFutureOption(symbol string) (FutureOption, error) {
	path := fmt.Sprintf("/instruments/future-options/%s", symbol)

	type instrumentResponse struct {
		FutureOption FutureOption `json:"data"`
		Context      string       `json:"context"`
	}

	instrumentRes := new(instrumentResponse)

	err := c.request(http.MethodGet, path, nil, nil, instrumentRes)
	if err != nil {
		return FutureOption{}, err
	}

	return instrumentRes.FutureOption, nil
}

// Returns metadata for all supported futures products.
func (c *Client) GetFutureProducts() ([]FutureProduct, error) {
	path := "/instruments/future-products"

	type instrumentResponse struct {
		Data struct {
			FutureProducts []FutureProduct `json:"items"`
		} `json:"data"`
	}

	instrumentRes := new(instrumentResponse)

	err := c.request(http.MethodGet, path, nil, nil, instrumentRes)
	if err != nil {
		return []FutureProduct{}, err
	}

	return instrumentRes.Data.FutureProducts, nil
}

// Get future product from exchange and product code.
func (c *Client) GetFutureProduct(exchange Exchange, productCode string) (FutureProduct, error) {
	path := fmt.Sprintf("/instruments/future-products/%s/%s", exchange, productCode)

	type instrumentResponse struct {
		FutureProduct FutureProduct `json:"data"`
	}

	instrumentRes := new(instrumentResponse)

	err := c.request(http.MethodGet, path, nil, nil, instrumentRes)
	if err != nil {
		return FutureProduct{}, err
	}

	return instrumentRes.FutureProduct, nil
}
