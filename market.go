package tasty

import (
	"fmt"
	"net/http"

	"github.com/austinbspencer/tasty-go/models"
)

func (c *Client) GetMarketMetrics(query models.MarketMetricsQuery) ([]models.MarketMetricVolatility, *Error) {
	if c.Session.SessionToken == nil {
		return []models.MarketMetricVolatility{}, &Error{Message: "Session is invalid: Session Token cannot be nil."}
	}

	reqURL := fmt.Sprintf("%s/market-metrics", c.baseURL)

	type marketMetricResponse struct {
		Data struct {
			MarketMetrics []models.MarketMetricVolatility `json:"items"`
		} `json:"data"`
	}

	marketMetricsRes := new(marketMetricResponse)

	header := http.Header{}
	header.Add("Authorization", *c.Session.SessionToken)

	err := c.get(reqURL, header, query, marketMetricsRes)
	if err != nil {
		return []models.MarketMetricVolatility{}, err
	}

	return marketMetricsRes.Data.MarketMetrics, nil
}
