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

func (c *Client) GetHistoricDividends(symbol string) ([]models.DividendInfo, *Error) {
	if c.Session.SessionToken == nil {
		return []models.DividendInfo{}, &Error{Message: "Session is invalid: Session Token cannot be nil."}
	}

	reqURL := fmt.Sprintf("%s/market-metrics/historic-corporate-events/dividends/%s", c.baseURL, symbol)

	type marketMetricResponse struct {
		Data struct {
			HistoricDividends []models.DividendInfo `json:"items"`
		} `json:"data"`
	}

	marketMetricsRes := new(marketMetricResponse)

	header := http.Header{}
	header.Add("Authorization", *c.Session.SessionToken)

	err := c.get(reqURL, header, nil, marketMetricsRes)
	if err != nil {
		return []models.DividendInfo{}, err
	}

	return marketMetricsRes.Data.HistoricDividends, nil
}

func (c *Client) GetHistoricEarnings(symbol string, query models.HistoricEarningsQuery) ([]models.EarningsInfo, *Error) {
	if c.Session.SessionToken == nil {
		return []models.EarningsInfo{}, &Error{Message: "Session is invalid: Session Token cannot be nil."}
	}

	reqURL := fmt.Sprintf("%s/market-metrics/historic-corporate-events/earnings-reports/%s", c.baseURL, symbol)

	type marketMetricResponse struct {
		Data struct {
			HistoricEarnings []models.EarningsInfo `json:"items"`
		} `json:"data"`
	}

	marketMetricsRes := new(marketMetricResponse)

	header := http.Header{}
	header.Add("Authorization", *c.Session.SessionToken)

	err := c.get(reqURL, header, query, marketMetricsRes)
	if err != nil {
		return []models.EarningsInfo{}, err
	}

	return marketMetricsRes.Data.HistoricEarnings, nil
}
