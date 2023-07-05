package tasty

import (
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// Returns an array of volatility data for given symbols.
func (c *Client) GetMarketMetrics(symbols []string) ([]MarketMetricVolatility, error) {
	path := "/market-metrics"

	type marketMetricResponse struct {
		Data struct {
			MarketMetrics []MarketMetricVolatility `json:"items"`
		} `json:"data"`
	}

	marketMetricsRes := new(marketMetricResponse)

	type marketMetrics struct {
		// Symbols is the list of symbols
		Symbols []string `url:"symbols,comma"`
	}

	query := marketMetrics{Symbols: symbols}

	err := c.request(http.MethodGet, path, query, nil, marketMetricsRes)
	if err != nil {
		return []MarketMetricVolatility{}, err
	}

	return marketMetricsRes.Data.MarketMetrics, nil
}

// Get historical dividend data.
func (c *Client) GetHistoricDividends(symbol string) ([]DividendInfo, error) {
	// url escape required for instances where "/" exists in symbol i.e. BRK/B
	path := fmt.Sprintf("/market-metrics/historic-corporate-events/dividends/%s", url.PathEscape(symbol))

	type marketMetricResponse struct {
		Data struct {
			HistoricDividends []DividendInfo `json:"items"`
		} `json:"data"`
	}

	marketMetricsRes := new(marketMetricResponse)

	err := c.customRequest(http.MethodGet, path, nil, nil, marketMetricsRes)
	if err != nil {
		return []DividendInfo{}, err
	}

	return marketMetricsRes.Data.HistoricDividends, nil
}

// Get historical earnings data.
func (c *Client) GetHistoricEarnings(symbol string, startDate time.Time) ([]EarningsInfo, error) {
	// url escape required for instances where "/" exists in symbol i.e. BRK/B
	path := fmt.Sprintf("/market-metrics/historic-corporate-events/earnings-reports/%s", url.PathEscape(symbol))

	type marketMetricResponse struct {
		Data struct {
			HistoricEarnings []EarningsInfo `json:"items"`
		} `json:"data"`
	}

	marketMetricsRes := new(marketMetricResponse)

	type historicEarnings struct {
		// StartDate string
		StartDate time.Time `layout:"2006-01-02" url:"start-date"`
	}

	query := historicEarnings{StartDate: startDate}

	err := c.customRequest(http.MethodGet, path, query, nil, marketMetricsRes)
	if err != nil {
		return []EarningsInfo{}, err
	}

	return marketMetricsRes.Data.HistoricEarnings, nil
}
