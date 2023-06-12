package tasty

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/austinbspencer/tasty-go/constants"
	"github.com/austinbspencer/tasty-go/models"
	"github.com/austinbspencer/tasty-go/utils"
)

func (c *Client) GetCryptocurrencies(symbols []string) ([]models.Cryptocurrency, *Error) {
	if c.Session.SessionToken == nil {
		return []models.Cryptocurrency{}, &Error{Message: "Session is invalid: Session Token cannot be nil."}
	}

	reqURL := fmt.Sprintf("%s/instruments/cryptocurrencies", c.baseURL)

	type instrumentResponse struct {
		Data struct {
			Cryptocurrencies []models.Cryptocurrency `json:"items"`
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

	err := c.get(reqURL, header, query, instrumentRes)
	if err != nil {
		return []models.Cryptocurrency{}, err
	}

	return instrumentRes.Data.Cryptocurrencies, nil
}

func (c *Client) GetCryptocurrency(symbol constants.Cryptocurrency) (models.Cryptocurrency, *Error) {
	if c.Session.SessionToken == nil {
		return models.Cryptocurrency{}, &Error{Message: "Session is invalid: Session Token cannot be nil."}
	}

	symbol += constants.Cryptocurrency(url.QueryEscape("/USD"))

	path := fmt.Sprintf("//%s/instruments/cryptocurrencies/%s", c.baseHost, symbol)

	type instrumentResponse struct {
		Crypto models.Cryptocurrency `json:"data"`
	}

	instrumentRes := new(instrumentResponse)

	header := http.Header{}
	header.Add("Authorization", *c.Session.SessionToken)

	err := c.customGet(path, header, nil, instrumentRes)
	if err != nil {
		return models.Cryptocurrency{}, err
	}

	return instrumentRes.Crypto, nil
}

// GetActiveEquities returns all active equities in a paginated fashion
func (c *Client) GetActiveEquities(query models.ActiveEquitiesQuery) ([]models.Equity, models.Pagination, *Error) {
	if c.Session.SessionToken == nil {
		return []models.Equity{}, models.Pagination{}, &Error{Message: "Session is invalid: Session Token cannot be nil."}
	}

	reqURL := fmt.Sprintf("%s/instruments/equities/active", c.baseURL)

	type instrumentResponse struct {
		Data struct {
			ActiveEquities []models.Equity `json:"items"`
		} `json:"data"`
		Pagination models.Pagination `json:"pagination"`
	}

	instrumentRes := new(instrumentResponse)

	header := http.Header{}
	header.Add("Authorization", *c.Session.SessionToken)

	err := c.get(reqURL, header, query, instrumentRes)
	if err != nil {
		return []models.Equity{}, models.Pagination{}, err
	}

	return instrumentRes.Data.ActiveEquities, instrumentRes.Pagination, nil
}

// GetEquities returns a set of equity definitions given an array of one or more symbols
func (c *Client) GetEquities(query models.EquitiesQuery) ([]models.Equity, *Error) {
	if c.Session.SessionToken == nil {
		return []models.Equity{}, &Error{Message: "Session is invalid: Session Token cannot be nil."}
	}

	reqURL := fmt.Sprintf("%s/instruments/equities", c.baseURL)

	type instrumentResponse struct {
		Data struct {
			Equities []models.Equity `json:"items"`
		} `json:"data"`
	}

	instrumentRes := new(instrumentResponse)

	header := http.Header{}
	header.Add("Authorization", *c.Session.SessionToken)

	err := c.get(reqURL, header, query, instrumentRes)
	if err != nil {
		return []models.Equity{}, err
	}

	return instrumentRes.Data.Equities, nil
}

// GetEquities returns a single equity definition for the provided symbol
func (c *Client) GetEquity(symbol string) (models.Equity, *Error) {
	if c.Session.SessionToken == nil {
		return models.Equity{}, &Error{Message: "Session is invalid: Session Token cannot be nil."}
	}

	reqURL := fmt.Sprintf("%s/instruments/equities/%s", c.baseURL, symbol)

	type instrumentResponse struct {
		Equity models.Equity `json:"data"`
	}

	instrumentRes := new(instrumentResponse)

	header := http.Header{}
	header.Add("Authorization", *c.Session.SessionToken)

	err := c.get(reqURL, header, nil, instrumentRes)
	if err != nil {
		return models.Equity{}, err
	}

	return instrumentRes.Equity, nil
}

// GetEquityOptions returns a set of equity options given one or more symbols
func (c *Client) GetEquityOptions(query models.EquityOptionsQuery) ([]models.EquityOption, *Error) {
	if c.Session.SessionToken == nil {
		return []models.EquityOption{}, &Error{Message: "Session is invalid: Session Token cannot be nil."}
	}

	reqURL := fmt.Sprintf("%s/instruments/equity-options", c.baseURL)

	type instrumentResponse struct {
		Data struct {
			EquityOptions []models.EquityOption `json:"items"`
		} `json:"data"`
	}

	instrumentRes := new(instrumentResponse)

	header := http.Header{}
	header.Add("Authorization", *c.Session.SessionToken)

	err := c.get(reqURL, header, query, instrumentRes)
	if err != nil {
		return []models.EquityOption{}, err
	}

	return instrumentRes.Data.EquityOptions, nil
}

// GetEquityOption returns a set of equity options given one or more symbols
func (c *Client) GetEquityOption(sym utils.EquityOptionsSymbology, active bool) (models.EquityOption, *Error) {
	if c.Session.SessionToken == nil {
		return models.EquityOption{}, &Error{Message: "Session is invalid: Session Token cannot be nil."}
	}

	occSymbol := sym.Build()

	reqURL := fmt.Sprintf("%s/instruments/equity-options/%s", c.baseURL, occSymbol)

	type instrumentResponse struct {
		EquityOption models.EquityOption `json:"data"`
	}

	type activeQuery struct {
		// Whether an option is available for trading with the broker.
		Active bool `url:"active"`
	}

	query := activeQuery{Active: active}

	instrumentRes := new(instrumentResponse)

	header := http.Header{}
	header.Add("Authorization", *c.Session.SessionToken)

	err := c.get(reqURL, header, query, instrumentRes)
	if err != nil {
		return models.EquityOption{}, err
	}

	return instrumentRes.EquityOption, nil
}

// GetFutures returns a set of outright futures given an array of one or more symbols.
func (c *Client) GetFutures(query models.FuturesQuery) ([]models.Future, *Error) {
	if c.Session.SessionToken == nil {
		return []models.Future{}, &Error{Message: "Session is invalid: Session Token cannot be nil."}
	}

	reqURL := fmt.Sprintf("%s/instruments/futures", c.baseURL)

	type instrumentResponse struct {
		Data struct {
			Futures []models.Future `json:"items"`
		} `json:"data"`
	}

	instrumentRes := new(instrumentResponse)

	header := http.Header{}
	header.Add("Authorization", *c.Session.SessionToken)

	err := c.get(reqURL, header, query, instrumentRes)
	if err != nil {
		return []models.Future{}, err
	}

	return instrumentRes.Data.Futures, nil
}

// GetFuture returns an outright future given a symbol.
func (c *Client) GetFuture(symbol string) (models.Future, *Error) {
	if c.Session.SessionToken == nil {
		return models.Future{}, &Error{Message: "Session is invalid: Session Token cannot be nil."}
	}

	reqURL := fmt.Sprintf("%s/instruments/futures/%s", c.baseURL, symbol)

	type instrumentResponse struct {
		Future models.Future `json:"data"`
	}

	instrumentRes := new(instrumentResponse)

	header := http.Header{}
	header.Add("Authorization", *c.Session.SessionToken)

	err := c.get(reqURL, header, nil, instrumentRes)
	if err != nil {
		return models.Future{}, err
	}

	return instrumentRes.Future, nil
}

// GetFutures returns metadata for all supported future option products
func (c *Client) GetFutureOptionProducts() ([]models.FutureOptionProduct, *Error) {
	if c.Session.SessionToken == nil {
		return []models.FutureOptionProduct{}, &Error{Message: "Session is invalid: Session Token cannot be nil."}
	}

	reqURL := fmt.Sprintf("%s/instruments/future-option-products", c.baseURL)

	type instrumentResponse struct {
		Data struct {
			FutureOptionProducts []models.FutureOptionProduct `json:"items"`
		} `json:"data"`
	}

	instrumentRes := new(instrumentResponse)

	header := http.Header{}
	header.Add("Authorization", *c.Session.SessionToken)

	err := c.get(reqURL, header, nil, instrumentRes)
	if err != nil {
		return []models.FutureOptionProduct{}, err
	}

	return instrumentRes.Data.FutureOptionProducts, nil
}

// GetFutures Get a future option product by exchange and root symbol
func (c *Client) GetFutureOptionProduct(exchange, rootSymbol string) (models.FutureOptionProduct, *Error) {
	if c.Session.SessionToken == nil {
		return models.FutureOptionProduct{}, &Error{Message: "Session is invalid: Session Token cannot be nil."}
	}

	reqURL := fmt.Sprintf("%s/instruments/future-option-products/%s/%s", c.baseURL, exchange, rootSymbol)

	type instrumentResponse struct {
		FutureOptionProduct models.FutureOptionProduct `json:"data"`
	}

	instrumentRes := new(instrumentResponse)

	header := http.Header{}
	header.Add("Authorization", *c.Session.SessionToken)

	err := c.get(reqURL, header, nil, instrumentRes)
	if err != nil {
		return models.FutureOptionProduct{}, err
	}

	return instrumentRes.FutureOptionProduct, nil
}

// GetFutures returns metadata for all supported future option products
func (c *Client) GetFutureOptions(query models.FutureOptionsQuery) ([]models.FutureOption, *Error) {
	if c.Session.SessionToken == nil {
		return []models.FutureOption{}, &Error{Message: "Session is invalid: Session Token cannot be nil."}
	}

	reqURL := fmt.Sprintf("%s/instruments/future-options", c.baseURL)

	type instrumentResponse struct {
		Data struct {
			FutureOptions []models.FutureOption `json:"items"`
		} `json:"data"`
	}

	instrumentRes := new(instrumentResponse)

	header := http.Header{}
	header.Add("Authorization", *c.Session.SessionToken)

	err := c.get(reqURL, header, query, instrumentRes)
	if err != nil {
		fmt.Println(err)
		return []models.FutureOption{}, err
	}

	return instrumentRes.Data.FutureOptions, nil
}

// Returns a future option given a symbol. Uses TW symbology: ./ESZ9 EW4U9 190927P2975
func (c *Client) GetFutureOption(symbol string) (models.FutureOption, *Error) {
	if c.Session.SessionToken == nil {
		return models.FutureOption{}, &Error{Message: "Session is invalid: Session Token cannot be nil."}
	}

	reqURL := fmt.Sprintf("%s/instruments/future-options/%s", c.baseURL, symbol)

	type instrumentResponse struct {
		FutureOption models.FutureOption `json:"data"`
		Context      string              `json:"context"`
	}

	instrumentRes := new(instrumentResponse)

	header := http.Header{}
	header.Add("Authorization", *c.Session.SessionToken)

	err := c.get(reqURL, header, nil, instrumentRes)
	if err != nil {
		return models.FutureOption{}, err
	}

	fmt.Println(instrumentRes.Context)

	return instrumentRes.FutureOption, nil
}

// GetFutureProducts returns metadata for all supported futures products
func (c *Client) GetFutureProducts() ([]models.FutureProduct, *Error) {
	if c.Session.SessionToken == nil {
		return []models.FutureProduct{}, &Error{Message: "Session is invalid: Session Token cannot be nil."}
	}

	reqURL := fmt.Sprintf("%s/instruments/future-products", c.baseURL)

	type instrumentResponse struct {
		Data struct {
			FutureProducts []models.FutureProduct `json:"items"`
		} `json:"data"`
	}

	instrumentRes := new(instrumentResponse)

	header := http.Header{}
	header.Add("Authorization", *c.Session.SessionToken)

	err := c.get(reqURL, header, nil, instrumentRes)
	if err != nil {
		fmt.Println(err)
		return []models.FutureProduct{}, err
	}

	return instrumentRes.Data.FutureProducts, nil
}

// GetFutureProduct Get future product from exchange and product code
func (c *Client) GetFutureProduct(exchange constants.Exchange, productCode string) (models.FutureProduct, *Error) {
	if c.Session.SessionToken == nil {
		return models.FutureProduct{}, &Error{Message: "Session is invalid: Session Token cannot be nil."}
	}

	reqURL := fmt.Sprintf("%s/instruments/future-products/%s/%s", c.baseURL, exchange, productCode)

	type instrumentResponse struct {
		FutureProduct models.FutureProduct `json:"data"`
	}

	instrumentRes := new(instrumentResponse)

	header := http.Header{}
	header.Add("Authorization", *c.Session.SessionToken)

	err := c.get(reqURL, header, nil, instrumentRes)
	if err != nil {
		return models.FutureProduct{}, err
	}

	return instrumentRes.FutureProduct, nil
}

// GetQuantityDecimalPrecisions Retrieve all quantity decimal precisions.
func (c *Client) GetQuantityDecimalPrecisions() ([]models.QuantityDecimalPrecision, *Error) {
	if c.Session.SessionToken == nil {
		return []models.QuantityDecimalPrecision{}, &Error{Message: "Session is invalid: Session Token cannot be nil."}
	}

	reqURL := fmt.Sprintf("%s/instruments/quantity-decimal-precisions", c.baseURL)

	type instrumentResponse struct {
		Data struct {
			QuantityDecimalPrecisions []models.QuantityDecimalPrecision `json:"items"`
		} `json:"data"`
	}

	instrumentRes := new(instrumentResponse)

	header := http.Header{}
	header.Add("Authorization", *c.Session.SessionToken)

	err := c.get(reqURL, header, nil, instrumentRes)
	if err != nil {
		return []models.QuantityDecimalPrecision{}, err
	}

	return instrumentRes.Data.QuantityDecimalPrecisions, nil
}

// GetWarrants Returns a set of warrant definitions that can be filtered by parameters
func (c *Client) GetWarrants(symbols []string) ([]models.Warrant, *Error) {
	if c.Session.SessionToken == nil {
		return []models.Warrant{}, &Error{Message: "Session is invalid: Session Token cannot be nil."}
	}

	reqURL := fmt.Sprintf("%s/instruments/warrants", c.baseURL)

	type instrumentResponse struct {
		Data struct {
			Warrants []models.Warrant `json:"items"`
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

	err := c.get(reqURL, header, query, instrumentRes)
	if err != nil {
		return []models.Warrant{}, err
	}

	return instrumentRes.Data.Warrants, nil
}

// GetWarrant Returns a single warrant definition for the provided symbol
func (c *Client) GetWarrant(symbol string) (models.Warrant, *Error) {
	if c.Session.SessionToken == nil {
		return models.Warrant{}, &Error{Message: "Session is invalid: Session Token cannot be nil."}
	}

	reqURL := fmt.Sprintf("%s/instruments/warrants/%s", c.baseURL, symbol)

	type instrumentResponse struct {
		Warrant models.Warrant `json:"data"`
	}

	instrumentRes := new(instrumentResponse)

	header := http.Header{}
	header.Add("Authorization", *c.Session.SessionToken)

	err := c.get(reqURL, header, nil, instrumentRes)
	if err != nil {
		return models.Warrant{}, err
	}

	return instrumentRes.Warrant, nil
}