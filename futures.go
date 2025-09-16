package tasty

import (
	"fmt"
	"net/http"
)

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
