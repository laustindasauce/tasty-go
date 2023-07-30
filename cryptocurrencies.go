package tasty

import (
	"fmt"
	"net/http"
	"net/url"
)

// Retrieve a set of cryptocurrencies given an array of one or more symbols.
func (c *Client) GetCryptocurrencies(symbols []string) ([]CryptocurrencyInfo, *http.Response, error) {
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

	resp, err := c.request(http.MethodGet, path, query, nil, instrumentRes)
	if err != nil {
		return []CryptocurrencyInfo{}, resp, err
	}

	return instrumentRes.Data.Cryptocurrencies, resp, nil
}

// Retrieve a cryptocurrency given a symbol.
func (c *Client) GetCryptocurrency(symbol Cryptocurrency) (CryptocurrencyInfo, *http.Response, error) {
	symbolString := url.PathEscape(string(symbol))

	path := fmt.Sprintf("/instruments/cryptocurrencies/%s", symbolString)

	type instrumentResponse struct {
		Crypto CryptocurrencyInfo `json:"data"`
	}

	instrumentRes := new(instrumentResponse)

	resp, err := c.customRequest(http.MethodGet, path, nil, nil, instrumentRes)
	if err != nil {
		return CryptocurrencyInfo{}, resp, err
	}

	return instrumentRes.Crypto, resp, nil
}
