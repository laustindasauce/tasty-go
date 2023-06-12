package tasty

import (
	"fmt"
	"net/http"

	"github.com/austinbspencer/tasty-go/models"
)

// Returns a futures option chain given a futures product code, i.e. ES
func (c *Client) GetFuturesOptionChains(productCode string) ([]models.FutureOption, *Error) {
	if c.Session.SessionToken == nil {
		return []models.FutureOption{}, &Error{Message: "Session is invalid: Session Token cannot be nil."}
	}

	reqURL := fmt.Sprintf("%s/futures-option-chains/%s", c.baseURL, productCode)

	type instrumentResponse struct {
		Data struct {
			FutureOptions []models.FutureOption `json:"items"`
		} `json:"data"`
	}

	instrumentRes := new(instrumentResponse)

	header := http.Header{}
	header.Add("Authorization", *c.Session.SessionToken)

	err := c.get(reqURL, header, nil, instrumentRes)
	if err != nil {
		return []models.FutureOption{}, err
	}

	return instrumentRes.Data.FutureOptions, nil
}

// Returns a futures option chain given a futures product code in a nested form to minimize
// redundant processing
func (c *Client) GetNestedFuturesOptionChains(productCode string) (models.NestedFuturesOptionChains, *Error) {
	if c.Session.SessionToken == nil {
		return models.NestedFuturesOptionChains{}, &Error{Message: "Session is invalid: Session Token cannot be nil."}
	}

	reqURL := fmt.Sprintf("%s/futures-option-chains/%s/nested", c.baseURL, productCode)

	type instrumentResponse struct {
		Chains models.NestedFuturesOptionChains `json:"data"`
	}

	instrumentRes := new(instrumentResponse)

	header := http.Header{}
	header.Add("Authorization", *c.Session.SessionToken)

	err := c.get(reqURL, header, nil, instrumentRes)
	if err != nil {
		return models.NestedFuturesOptionChains{}, err
	}

	return instrumentRes.Chains, nil
}
