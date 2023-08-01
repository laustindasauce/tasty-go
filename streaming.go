package tasty

import "net/http"

// Returns the appropriate API quote streamer endpoint, level and identification token
// for the current customer to receive market data.
func (c *Client) GetQuoteStreamerTokens() (QuoteStreamerTokenAuthResult, *http.Response, error) {
	path := "/api-quote-tokens"

	type customerResponse struct {
		Streamer QuoteStreamerTokenAuthResult `json:"data"`
	}

	customersRes := new(customerResponse)

	resp, err := c.request(http.MethodGet, path, nil, nil, customersRes)
	if err != nil {
		return QuoteStreamerTokenAuthResult{}, resp, err
	}

	return customersRes.Streamer, resp, nil
}
