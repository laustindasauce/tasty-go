package tasty

import (
	"fmt"
	"net/http"
	"net/url"
)

// Retrieve a set of cryptocurrencies given an array of one or more symbols.
func (c *Client) GetCryptocurrencies(symbols []string) ([]CryptocurrencyInfo, *Error) {
	if c.Session.SessionToken == nil {
		return []CryptocurrencyInfo{}, &Error{Message: "Session is invalid: Session Token cannot be nil."}
	}

	path := "/instruments/cryptocurrencies"

	type instrumentResponse struct {
		Data struct {
			Cryptocurrencies []CryptocurrencyInfo `json:"items"`
		} `json:"data"`
	}

	type symbolsQuery struct {
		// Symbols is the list of symbols
		Symbols []string `url:"symbol[]"`
	}

	query := symbolsQuery{Symbols: symbols}

	instrumentRes := new(instrumentResponse)

	header := http.Header{}
	header.Add("Authorization", *c.Session.SessionToken)

	err := c.request(http.MethodGet, path, header, query, nil, instrumentRes)
	if err != nil {
		return []CryptocurrencyInfo{}, err
	}

	return instrumentRes.Data.Cryptocurrencies, nil
}

// Retrieve a cryptocurrency given a symbol.
func (c *Client) GetCryptocurrency(symbol Cryptocurrency) (CryptocurrencyInfo, *Error) {
	if c.Session.SessionToken == nil {
		return CryptocurrencyInfo{}, &Error{Message: "Session is invalid: Session Token cannot be nil."}
	}

	symbol += Cryptocurrency(url.PathEscape("/USD"))

	path := fmt.Sprintf("/instruments/cryptocurrencies/%s", symbol)

	type instrumentResponse struct {
		Crypto CryptocurrencyInfo `json:"data"`
	}

	instrumentRes := new(instrumentResponse)

	header := http.Header{}
	header.Add("Authorization", *c.Session.SessionToken)

	err := c.customRequest(http.MethodGet, path, header, nil, nil, instrumentRes)
	if err != nil {
		return CryptocurrencyInfo{}, err
	}

	return instrumentRes.Crypto, nil
}

// Returns all active equities in a paginated fashion.
func (c *Client) GetActiveEquities(query ActiveEquitiesQuery) ([]Equity, Pagination, *Error) {
	if c.Session.SessionToken == nil {
		return []Equity{}, Pagination{}, &Error{Message: "Session is invalid: Session Token cannot be nil."}
	}

	path := "/instruments/equities/active"

	type instrumentResponse struct {
		Data struct {
			ActiveEquities []Equity `json:"items"`
		} `json:"data"`
		Pagination Pagination `json:"pagination"`
	}

	instrumentRes := new(instrumentResponse)

	header := http.Header{}
	header.Add("Authorization", *c.Session.SessionToken)

	err := c.request(http.MethodGet, path, header, query, nil, instrumentRes)
	if err != nil {
		return []Equity{}, Pagination{}, err
	}

	return instrumentRes.Data.ActiveEquities, instrumentRes.Pagination, nil
}

// Returns a set of equity definitions given an array of one or more symbols.
func (c *Client) GetEquities(query EquitiesQuery) ([]Equity, *Error) {
	if c.Session.SessionToken == nil {
		return []Equity{}, &Error{Message: "Session is invalid: Session Token cannot be nil."}
	}

	path := "/instruments/equities"

	type instrumentResponse struct {
		Data struct {
			Equities []Equity `json:"items"`
		} `json:"data"`
	}

	instrumentRes := new(instrumentResponse)

	header := http.Header{}
	header.Add("Authorization", *c.Session.SessionToken)

	err := c.request(http.MethodGet, path, header, query, nil, instrumentRes)
	if err != nil {
		return []Equity{}, err
	}

	return instrumentRes.Data.Equities, nil
}

// Returns a single equity definition for the provided symbol.
func (c *Client) GetEquity(symbol string) (Equity, *Error) {
	if c.Session.SessionToken == nil {
		return Equity{}, &Error{Message: "Session is invalid: Session Token cannot be nil."}
	}

	// url escape required for instances where "/" exists in symbol i.e. BRK/B
	path := fmt.Sprintf("/instruments/equities/%s", url.PathEscape(symbol))

	type instrumentResponse struct {
		Equity Equity `json:"data"`
	}

	instrumentRes := new(instrumentResponse)

	header := http.Header{}
	header.Add("Authorization", *c.Session.SessionToken)

	// customRequest required for instances where "/" exists in symbol i.e. BRK/B
	err := c.customRequest(http.MethodGet, path, header, nil, nil, instrumentRes)
	if err != nil {
		return Equity{}, err
	}

	return instrumentRes.Equity, nil
}

// Returns a set of equity options given one or more symbols.
func (c *Client) GetEquityOptions(query EquityOptionsQuery) ([]EquityOption, *Error) {
	if c.Session.SessionToken == nil {
		return []EquityOption{}, &Error{Message: "Session is invalid: Session Token cannot be nil."}
	}

	path := "/instruments/equity-options"

	type instrumentResponse struct {
		Data struct {
			EquityOptions []EquityOption `json:"items"`
		} `json:"data"`
	}

	instrumentRes := new(instrumentResponse)

	header := http.Header{}
	header.Add("Authorization", *c.Session.SessionToken)

	err := c.request(http.MethodGet, path, header, query, nil, instrumentRes)
	if err != nil {
		return []EquityOption{}, err
	}

	return instrumentRes.Data.EquityOptions, nil
}

// Returns a set of equity options given one or more symbols.
func (c *Client) GetEquityOption(sym EquityOptionsSymbology, active bool) (EquityOption, *Error) {
	if c.Session.SessionToken == nil {
		return EquityOption{}, &Error{Message: "Session is invalid: Session Token cannot be nil."}
	}

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

	header := http.Header{}
	header.Add("Authorization", *c.Session.SessionToken)

	err := c.request(http.MethodGet, path, header, query, nil, instrumentRes)
	if err != nil {
		return EquityOption{}, err
	}

	return instrumentRes.EquityOption, nil
}

// Returns a set of outright futures given an array of one or more symbols.
func (c *Client) GetFutures(query FuturesQuery) ([]Future, *Error) {
	if c.Session.SessionToken == nil {
		return []Future{}, &Error{Message: "Session is invalid: Session Token cannot be nil."}
	}

	path := "/instruments/futures"

	type instrumentResponse struct {
		Data struct {
			Futures []Future `json:"items"`
		} `json:"data"`
	}

	instrumentRes := new(instrumentResponse)

	header := http.Header{}
	header.Add("Authorization", *c.Session.SessionToken)

	err := c.request(http.MethodGet, path, header, query, nil, instrumentRes)
	if err != nil {
		return []Future{}, err
	}

	return instrumentRes.Data.Futures, nil
}

// Returns an outright future given a symbol.
func (c *Client) GetFuture(symbol string) (Future, *Error) {
	if c.Session.SessionToken == nil {
		return Future{}, &Error{Message: "Session is invalid: Session Token cannot be nil."}
	}

	path := fmt.Sprintf("/instruments/futures/%s", symbol)

	type instrumentResponse struct {
		Future Future `json:"data"`
	}

	instrumentRes := new(instrumentResponse)

	header := http.Header{}
	header.Add("Authorization", *c.Session.SessionToken)

	err := c.request(http.MethodGet, path, header, nil, nil, instrumentRes)
	if err != nil {
		return Future{}, err
	}

	return instrumentRes.Future, nil
}

// Returns metadata for all supported future option products.
func (c *Client) GetFutureOptionProducts() ([]FutureOptionProduct, *Error) {
	if c.Session.SessionToken == nil {
		return []FutureOptionProduct{}, &Error{Message: "Session is invalid: Session Token cannot be nil."}
	}

	path := "/instruments/future-option-products"

	type instrumentResponse struct {
		Data struct {
			FutureOptionProducts []FutureOptionProduct `json:"items"`
		} `json:"data"`
	}

	instrumentRes := new(instrumentResponse)

	header := http.Header{}
	header.Add("Authorization", *c.Session.SessionToken)

	err := c.request(http.MethodGet, path, header, nil, nil, instrumentRes)
	if err != nil {
		return []FutureOptionProduct{}, err
	}

	return instrumentRes.Data.FutureOptionProducts, nil
}

// Get a future option product by exchange and root symbol.
func (c *Client) GetFutureOptionProduct(exchange, rootSymbol string) (FutureOptionProduct, *Error) {
	if c.Session.SessionToken == nil {
		return FutureOptionProduct{}, &Error{Message: "Session is invalid: Session Token cannot be nil."}
	}

	path := fmt.Sprintf("/instruments/future-option-products/%s/%s", exchange, rootSymbol)

	type instrumentResponse struct {
		FutureOptionProduct FutureOptionProduct `json:"data"`
	}

	instrumentRes := new(instrumentResponse)

	header := http.Header{}
	header.Add("Authorization", *c.Session.SessionToken)

	err := c.request(http.MethodGet, path, header, nil, nil, instrumentRes)
	if err != nil {
		return FutureOptionProduct{}, err
	}

	return instrumentRes.FutureOptionProduct, nil
}

// Returns a set of future option(s) given an array of one or more symbols.
// Uses TW symbology: [./ESZ9 EW4U9 190927P2975].
func (c *Client) GetFutureOptions(query FutureOptionsQuery) ([]FutureOption, *Error) {
	if c.Session.SessionToken == nil {
		return []FutureOption{}, &Error{Message: "Session is invalid: Session Token cannot be nil."}
	}

	path := "/instruments/future-options"

	type instrumentResponse struct {
		Data struct {
			FutureOptions []FutureOption `json:"items"`
		} `json:"data"`
	}

	instrumentRes := new(instrumentResponse)

	header := http.Header{}
	header.Add("Authorization", *c.Session.SessionToken)

	err := c.request(http.MethodGet, path, header, query, nil, instrumentRes)
	if err != nil {
		return []FutureOption{}, err
	}

	return instrumentRes.Data.FutureOptions, nil
}

// Returns a future option given a symbol. Uses TW symbology: ./ESZ9 EW4U9 190927P2975.
func (c *Client) GetFutureOption(symbol string) (FutureOption, *Error) {
	if c.Session.SessionToken == nil {
		return FutureOption{}, &Error{Message: "Session is invalid: Session Token cannot be nil."}
	}

	path := fmt.Sprintf("/instruments/future-options/%s", symbol)

	type instrumentResponse struct {
		FutureOption FutureOption `json:"data"`
		Context      string       `json:"context"`
	}

	instrumentRes := new(instrumentResponse)

	header := http.Header{}
	header.Add("Authorization", *c.Session.SessionToken)

	err := c.request(http.MethodGet, path, header, nil, nil, instrumentRes)
	if err != nil {
		return FutureOption{}, err
	}

	return instrumentRes.FutureOption, nil
}

// Returns metadata for all supported futures products.
func (c *Client) GetFutureProducts() ([]FutureProduct, *Error) {
	if c.Session.SessionToken == nil {
		return []FutureProduct{}, &Error{Message: "Session is invalid: Session Token cannot be nil."}
	}

	path := "/instruments/future-products"

	type instrumentResponse struct {
		Data struct {
			FutureProducts []FutureProduct `json:"items"`
		} `json:"data"`
	}

	instrumentRes := new(instrumentResponse)

	header := http.Header{}
	header.Add("Authorization", *c.Session.SessionToken)

	err := c.request(http.MethodGet, path, header, nil, nil, instrumentRes)
	if err != nil {
		return []FutureProduct{}, err
	}

	return instrumentRes.Data.FutureProducts, nil
}

// Get future product from exchange and product code.
func (c *Client) GetFutureProduct(exchange Exchange, productCode string) (FutureProduct, *Error) {
	if c.Session.SessionToken == nil {
		return FutureProduct{}, &Error{Message: "Session is invalid: Session Token cannot be nil."}
	}

	path := fmt.Sprintf("/instruments/future-products/%s/%s", exchange, productCode)

	type instrumentResponse struct {
		FutureProduct FutureProduct `json:"data"`
	}

	instrumentRes := new(instrumentResponse)

	header := http.Header{}
	header.Add("Authorization", *c.Session.SessionToken)

	err := c.request(http.MethodGet, path, header, nil, nil, instrumentRes)
	if err != nil {
		return FutureProduct{}, err
	}

	return instrumentRes.FutureProduct, nil
}

// Retrieve all quantity decimal precisions.
func (c *Client) GetQuantityDecimalPrecisions() ([]QuantityDecimalPrecision, *Error) {
	if c.Session.SessionToken == nil {
		return []QuantityDecimalPrecision{}, &Error{Message: "Session is invalid: Session Token cannot be nil."}
	}

	path := "/instruments/quantity-decimal-precisions"

	type instrumentResponse struct {
		Data struct {
			QuantityDecimalPrecisions []QuantityDecimalPrecision `json:"items"`
		} `json:"data"`
	}

	instrumentRes := new(instrumentResponse)

	header := http.Header{}
	header.Add("Authorization", *c.Session.SessionToken)

	err := c.request(http.MethodGet, path, header, nil, nil, instrumentRes)
	if err != nil {
		return []QuantityDecimalPrecision{}, err
	}

	return instrumentRes.Data.QuantityDecimalPrecisions, nil
}

// Returns a set of warrant definitions that can be filtered by parameters.
func (c *Client) GetWarrants(symbols []string) ([]Warrant, *Error) {
	if c.Session.SessionToken == nil {
		return []Warrant{}, &Error{Message: "Session is invalid: Session Token cannot be nil."}
	}

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

	header := http.Header{}
	header.Add("Authorization", *c.Session.SessionToken)

	err := c.request(http.MethodGet, path, header, query, nil, instrumentRes)
	if err != nil {
		return []Warrant{}, err
	}

	return instrumentRes.Data.Warrants, nil
}

// Returns a single warrant definition for the provided symbol.
func (c *Client) GetWarrant(symbol string) (Warrant, *Error) {
	if c.Session.SessionToken == nil {
		return Warrant{}, &Error{Message: "Session is invalid: Session Token cannot be nil."}
	}

	path := fmt.Sprintf("/instruments/warrants/%s", symbol)

	type instrumentResponse struct {
		Warrant Warrant `json:"data"`
	}

	instrumentRes := new(instrumentResponse)

	header := http.Header{}
	header.Add("Authorization", *c.Session.SessionToken)

	err := c.request(http.MethodGet, path, header, nil, nil, instrumentRes)
	if err != nil {
		return Warrant{}, err
	}

	return instrumentRes.Warrant, nil
}
