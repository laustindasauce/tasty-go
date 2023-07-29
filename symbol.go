package tasty

import (
	"fmt"
	"net/http"
	"net/url"
)

// Returns an array of symbol data.
func (c *Client) SymbolSearch(symbol string) ([]SymbolData, *http.Response, error) {
	// url escape required for instances where "/" exists in symbol i.e. BRK/B
	symbol = url.PathEscape(symbol)

	path := fmt.Sprintf("/symbols/search/%s", symbol)

	type symbolResponse struct {
		Data struct {
			SymbolData []SymbolData `json:"items"`
		} `json:"data"`
	}

	symbolRes := new(symbolResponse)

	// customRequest required for instances where "/" exists in symbol i.e. BRK/B
	resp, err := c.customRequest(http.MethodGet, path, nil, nil, symbolRes)
	if err != nil {
		return []SymbolData{}, resp, err
	}

	return symbolRes.Data.SymbolData, resp, nil
}
