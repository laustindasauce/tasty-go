package tasty

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/austinbspencer/tasty-go/models"
)

func (c *Client) GetMarketMetrics(symbols []string) ([]models.MarketMetricVolatility, *Error) {
	if c.Session.SessionToken == nil {
		return []models.MarketMetricVolatility{}, &Error{Message: "Session is invalid: Session Token cannot be nil."}
	}

	path := "/market-metrics"

	type marketMetricResponse struct {
		Data struct {
			MarketMetrics []models.MarketMetricVolatility `json:"items"`
		} `json:"data"`
	}

	marketMetricsRes := new(marketMetricResponse)

	header := http.Header{}
	header.Add("Authorization", *c.Session.SessionToken)

	type marketMetrics struct {
		// Symbols is the list of symbols
		Symbols []string `url:"symbols,comma"`
	}

	query := marketMetrics{Symbols: symbols}

	err := c.request(http.MethodGet, path, header, query, nil, marketMetricsRes)
	if err != nil {
		return []models.MarketMetricVolatility{}, err
	}

	return marketMetricsRes.Data.MarketMetrics, nil
}

func (c *Client) GetHistoricDividends(symbol string) ([]models.DividendInfo, *Error) {
	if c.Session.SessionToken == nil {
		return []models.DividendInfo{}, &Error{Message: "Session is invalid: Session Token cannot be nil."}
	}

	// url escape required for instances where "/" exists in symbol i.e. BRK/B
	path := fmt.Sprintf("/market-metrics/historic-corporate-events/dividends/%s", url.PathEscape(symbol))

	type marketMetricResponse struct {
		Data struct {
			HistoricDividends []models.DividendInfo `json:"items"`
		} `json:"data"`
	}

	marketMetricsRes := new(marketMetricResponse)

	header := http.Header{}
	header.Add("Authorization", *c.Session.SessionToken)

	err := c.customRequest(http.MethodGet, path, header, nil, nil, marketMetricsRes)
	if err != nil {
		return []models.DividendInfo{}, err
	}

	return marketMetricsRes.Data.HistoricDividends, nil
}

func (c *Client) GetHistoricEarnings(symbol string, startDate time.Time) ([]models.EarningsInfo, *Error) {
	if c.Session.SessionToken == nil {
		return []models.EarningsInfo{}, &Error{Message: "Session is invalid: Session Token cannot be nil."}
	}

	// url escape required for instances where "/" exists in symbol i.e. BRK/B
	path := fmt.Sprintf("/market-metrics/historic-corporate-events/earnings-reports/%s", url.PathEscape(symbol))

	type marketMetricResponse struct {
		Data struct {
			HistoricEarnings []models.EarningsInfo `json:"items"`
		} `json:"data"`
	}

	marketMetricsRes := new(marketMetricResponse)

	header := http.Header{}
	header.Add("Authorization", *c.Session.SessionToken)

	type historicEarnings struct {
		// StartDate string
		StartDate time.Time `layout:"2006-01-02" url:"start-date"`
	}

	query := historicEarnings{StartDate: startDate}

	err := c.customRequest(http.MethodGet, path, header, query, nil, marketMetricsRes)
	if err != nil {
		return []models.EarningsInfo{}, err
	}

	return marketMetricsRes.Data.HistoricEarnings, nil
}
