package tasty //nolint:testpackage // testing private field

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func TestGetMarketMetrics(t *testing.T) {
	setup()
	defer teardown()

	symbols := []string{"AAPL", "TSLA"}

	mux.HandleFunc("/market-metrics", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, marketMetricsResp)
	})

	resp, httpResp, err := client.GetMarketMetrics(symbols)
	require.Nil(t, err)
	require.NotNil(t, httpResp)

	require.Equal(t, 2, len(resp))

	aapl := resp[0]

	require.Equal(t, decimal.NewFromFloat(0.224186663), aapl.ImpliedVolatilityIndex)
	require.Equal(t, decimal.NewFromFloat(-0.014691909), aapl.ImpliedVolatilityIndex5DayChange)
	require.Equal(t, decimal.NewFromFloat(0.02231313), aapl.ImpliedVolatilityIndexRank)
	require.Equal(t, decimal.NewFromFloat(0.02231313), aapl.TosImpliedVolatilityIndexRank)
	require.Equal(t, decimal.NewFromFloat(0.037374278), aapl.TwImpliedVolatilityIndexRank)
	require.Equal(t, time.Date(2023, time.June, 7, 19, 55, 22, 151000000, time.UTC), aapl.TosImpliedVolatilityIndexRankUpdatedAt)
	require.Equal(t, "tos", aapl.ImpliedVolatilityIndexRankSource)
	require.Equal(t, decimal.NewFromFloat(0.031249801), aapl.ImpliedVolatilityPercentile)
	require.Equal(t, time.Date(2023, time.June, 7, 19, 55, 22, 110000000, time.UTC), aapl.ImpliedVolatilityUpdatedAt)
	require.Equal(t, decimal.NewFromFloat(899091.420242437), aapl.LiquidityValue)
	require.Equal(t, decimal.NewFromFloat(1.516319123), aapl.LiquidityRank)
	require.Equal(t, 4, aapl.LiquidityRating)
	require.Equal(t, time.Date(2023, time.June, 7, 19, 55, 23, 376000000, time.UTC), aapl.UpdatedAt)
	require.Equal(t, decimal.NewFromFloat(1.265634436), aapl.Beta)
	require.Equal(t, time.Date(2023, time.June, 4, 17, 0, 35, 805000000, time.UTC), aapl.BetaUpdatedAt)
	require.Equal(t, decimal.NewFromFloat(0.79), aapl.CorrSpy3month)
	require.Equal(t, decimal.NewFromFloat(0.24), aapl.DividendRatePerShare)
	require.True(t, aapl.AnnualDividendPerShare.Equal(decimal.Zero))
	require.Equal(t, decimal.NewFromFloat(0.006396146), aapl.DividendYield)
	require.Equal(t, "2023-05-12", aapl.DividendExDate)
	require.Equal(t, "2022-08-05", aapl.DividendNextDate)
	require.Equal(t, "2023-05-18", aapl.DividendPayDate)
	require.Equal(t, time.Date(2023, time.May, 8, 0, 16, 36, 676000000, time.UTC), aapl.DividendUpdatedAt)
	require.Equal(t, "XNAS", aapl.ListedMarket)
	require.Equal(t, "Easy To Borrow", aapl.Lendability)
	require.True(t, aapl.BorrowRate.Equal(decimal.Zero))
	require.Equal(t, 2846108626900, aapl.MarketCap)
	require.Equal(t, decimal.NewFromFloat(21.87), aapl.ImpliedVolatility30Day)
	require.Equal(t, decimal.NewFromFloat(20.02), aapl.HistoricalVolatility30Day)
	require.Equal(t, decimal.NewFromFloat(19.55), aapl.HistoricalVolatility60Day)
	require.Equal(t, decimal.NewFromFloat(20.24), aapl.HistoricalVolatility90Day)
	require.Equal(t, decimal.NewFromFloat(1.85), aapl.IvHv30DayDifference)
	require.True(t, aapl.PriceEarningsRatio.Equal(decimal.Zero))
	require.True(t, aapl.EarningsPerShare.Equal(decimal.Zero))

	// Option Expiration IVs
	optExpiryIVs := aapl.OptionExpirationImpliedVolatilities

	require.Equal(t, 3, len(optExpiryIVs))
	require.Equal(t, "2023-06-09", optExpiryIVs[0].ExpirationDate)
	require.Equal(t, "Standard", optExpiryIVs[0].OptionChainType)
	require.Equal(t, "PM", optExpiryIVs[0].SettlementType)
	require.Equal(t, decimal.NewFromFloat(0.288187496), optExpiryIVs[0].ImpliedVolatility)

	// Liquidity Running State
	liquidityRunningState := aapl.LiquidityRunningState

	require.Equal(t, decimal.NewFromFloat(38347762.009590656), liquidityRunningState.Sum)
	require.Equal(t, 48, liquidityRunningState.Count)
	require.Equal(t, time.Date(2023, time.June, 3, 10, 0, 10, 347000000, time.UTC), liquidityRunningState.StartedAt)
	require.Equal(t, time.Date(2023, time.June, 7, 19, 55, 22, 110000000, time.UTC), liquidityRunningState.UpdatedAt)

	// Earnings
	earnings := aapl.Earnings

	require.True(t, earnings.Visible)
	require.Equal(t, "2023-07-27", earnings.ExpectedReportDate)
	require.False(t, earnings.Estimated)
	require.Zero(t, earnings.LateFlag)
	require.Equal(t, "2023-06-01", earnings.QuarterEndDate)
	require.Equal(t, decimal.NewFromFloat(1.2), earnings.ActualEPS)
	require.Equal(t, decimal.NewFromFloat(1.18), earnings.ConsensusEstimate)
	require.Equal(t, time.Date(2023, time.June, 7, 11, 1, 36, 77000000, time.UTC), earnings.UpdatedAt)
}

func TestGetMarketMetricsError(t *testing.T) {
	setup()
	defer teardown()

	symbols := []string{"AAPL", "TSLA"}

	mux.HandleFunc("/market-metrics", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(401)
		fmt.Fprint(writer, tastyUnauthorizedError)
	})

	_, _, err := client.GetMarketMetrics(symbols)
	expectedUnauthorized(t, err)
}

func TestGetHistoricDividends(t *testing.T) {
	setup()
	defer teardown()

	symbol := "AAPL"

	mux.HandleFunc(fmt.Sprintf("/market-metrics/historic-corporate-events/dividends/%s", symbol), func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, historicDividendsResp)
	})

	resp, httpResp, err := client.GetHistoricDividends(symbol)
	require.Nil(t, err)
	require.NotNil(t, httpResp)

	require.Equal(t, 39, len(resp))
	require.Equal(t, "2023-05-12", resp[0].OccurredDate)
	require.Equal(t, decimal.NewFromFloat(0.24), resp[0].Amount)
}

func TestGetHistoricDividendsError(t *testing.T) {
	setup()
	defer teardown()

	symbol := "AAPL"

	mux.HandleFunc(fmt.Sprintf("/market-metrics/historic-corporate-events/dividends/%s", symbol), func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(401)
		fmt.Fprint(writer, tastyUnauthorizedError)
	})

	_, _, err := client.GetHistoricDividends(symbol)
	expectedUnauthorized(t, err)
}

func TestGetHistoricEarnings(t *testing.T) {
	setup()
	defer teardown()

	startDate := time.Now().AddDate(-1, 0, 0)
	symbol := "AAPL"

	mux.HandleFunc(fmt.Sprintf("/market-metrics/historic-corporate-events/earnings-reports/%s", symbol), func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, historicEarningsResp)
	})

	resp, httpResp, err := client.GetHistoricEarnings(symbol, startDate)
	require.Nil(t, err)
	require.NotNil(t, httpResp)

	require.Equal(t, 4, len(resp))
	require.Equal(t, "2022-06-30", resp[0].OccurredDate)
	require.Equal(t, decimal.NewFromFloat(1.2), resp[0].Eps)
}

func TestGetHistoricEarningsError(t *testing.T) {
	setup()
	defer teardown()

	startDate := time.Now().AddDate(-1, 0, 0)
	symbol := "AAPL"

	mux.HandleFunc(fmt.Sprintf("/market-metrics/historic-corporate-events/earnings-reports/%s", symbol), func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(401)
		fmt.Fprint(writer, tastyUnauthorizedError)
	})

	_, _, err := client.GetHistoricEarnings(symbol, startDate)
	expectedUnauthorized(t, err)
}

func TestGetMarketDataByType(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/market-data/by-type", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, marketDataByTypeResp)
	})

	query := MarketDataQuery{
		Equity:         []string{"AAPL"},
		Cryptocurrency: []string{"BTC"},
		EquityOption:   []string{"AAPL  250919C00005000"},
		Future:         []string{"ES"},
		FutureOption:   []string{"./ESZ5 EW4V5 251024C5850"},
		Index:          []string{"SPY"},
	}

	resp, pagination, httpResp, err := client.GetMarketDataByType(query)
	require.Nil(t, err)
	require.NotNil(t, httpResp)

	require.Equal(t, 6, len(resp))

	equity := resp[0]

	require.Equal(t, "BTC", equity.Symbol)
	require.Equal(t, EquityIT, equity.InstrumentType)
	require.Equal(t, "2025-09-13T21:33:06.641Z", equity.UpdatedAt.Format(time.RFC3339Nano))
	require.True(t, equity.Bid.Equal(decimal.NewFromFloat(51.4)))
	require.True(t, equity.BidSize.Equal(decimal.NewFromFloat(1.0)))
	require.True(t, equity.Ask.Equal(decimal.NewFromFloat(51.66)))
	require.True(t, equity.AskSize.Equal(decimal.NewFromFloat(4.0)))
	require.True(t, equity.Mid.Equal(decimal.NewFromFloat(51.53)))
	require.True(t, equity.Mark.Equal(decimal.NewFromFloat(51.66)))
	require.True(t, equity.Last.Equal(decimal.NewFromFloat(51.78)))
	require.True(t, equity.LastMkt.Equal(decimal.NewFromFloat(51.78)))
	require.True(t, equity.Beta.Equal(decimal.NewFromFloat(1.719217504)))
	require.True(t, equity.Open.Equal(decimal.NewFromFloat(50.94)))
	require.True(t, equity.DayHighPrice.Equal(decimal.NewFromFloat(51.795)))
	require.True(t, equity.DayLowPrice.Equal(decimal.NewFromFloat(50.87)))
	require.True(t, equity.Close.Equal(decimal.NewFromFloat(51.78)))
	require.Equal(t, "Final", equity.ClosePriceType)
	require.True(t, equity.PrevClose.Equal(decimal.NewFromFloat(50.74)))
	require.Equal(t, "Final", equity.PrevClosePriceType)
	require.Equal(t, "2025-09-12", equity.SummaryDate)
	require.Equal(t, "2025-09-11", equity.PrevCloseDate)
	require.True(t, equity.LowLimitPrice.Equal(decimal.NewFromFloat(46.3)))
	require.True(t, equity.HighLimitPrice.Equal(decimal.NewFromFloat(56.59)))
	require.False(t, equity.IsTradingHalted)
	require.Equal(t, -1, equity.HaltStartTime)
	require.Equal(t, -1, equity.HaltEndTime)
	require.True(t, equity.YearLowPrice.Equal(decimal.NewFromFloat(25.5)))
	require.True(t, equity.YearHighPrice.Equal(decimal.NewFromFloat(54.48)))
	require.True(t, equity.Volume.Equal(decimal.NewFromFloat(995431.0)))

	// Pagination
	require.Equal(t, 4, pagination.PerPage)
	require.Equal(t, 0, pagination.PageOffset)
	require.Equal(t, 0, pagination.ItemOffset)
	require.Equal(t, 12143, pagination.TotalItems)
	require.Equal(t, 3036, pagination.TotalPages)
	require.Equal(t, 4, pagination.CurrentItemCount)
	require.Nil(t, pagination.PreviousLink)
	require.Nil(t, pagination.NextLink)
	require.Nil(t, pagination.PagingLinkTemplate)
}

func TestGetMarketDataByTypeError(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/market-data/by-type", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(401)
		fmt.Fprint(writer, tastyUnauthorizedError)
	})

	_, _, httpResp, err := client.GetMarketDataByType(MarketDataQuery{})
	expectedUnauthorized(t, err)
	require.NotNil(t, httpResp)
}

const marketMetricsResp = `{
  "data": {
    "items": [
      {
        "symbol": "AAPL",
        "implied-volatility-index": "0.224186663",
        "implied-volatility-index-5-day-change": "-0.014691909",
        "implied-volatility-index-rank": "0.02231313",
        "tos-implied-volatility-index-rank": "0.02231313",
        "tw-implied-volatility-index-rank": "0.037374278",
        "tos-implied-volatility-index-rank-updated-at": "2023-06-07T19:55:22.151Z",
        "implied-volatility-index-rank-source": "tos",
        "implied-volatility-percentile": "0.031249801",
        "implied-volatility-updated-at": "2023-06-07T19:55:22.110Z",
        "liquidity-value": "899091.420242437",
        "liquidity-rank": "1.516319123",
        "liquidity-rating": 4,
        "updated-at": "2023-06-07T19:55:23.376Z",
        "option-expiration-implied-volatilities": [
          {
            "expiration-date": "2023-06-09",
            "option-chain-type": "Standard",
            "settlement-type": "PM",
            "implied-volatility": "0.288187496"
          },
          {
            "expiration-date": "2023-06-16",
            "option-chain-type": "Standard",
            "settlement-type": "PM",
            "implied-volatility": "0.253236223"
          },
          {
            "expiration-date": "2023-06-23",
            "option-chain-type": "Standard",
            "settlement-type": "PM",
            "implied-volatility": "0.23260259"
          }
        ],
        "liquidity-running-state": {
          "sum": "38347762.009590656",
          "count": 48,
          "started-at": "2023-06-03T10:00:10.347Z",
          "updated-at": "2023-06-07T19:55:22.110Z"
        },
        "beta": "1.265634436",
        "beta-updated-at": "2023-06-04T17:00:35.805Z",
        "corr-spy-3month": "0.79",
        "dividend-rate-per-share": "0.24",
        "annual-dividend-per-share": "0.0",
        "dividend-yield": "0.006396146",
        "dividend-ex-date": "2023-05-12",
        "dividend-next-date": "2022-08-05",
        "dividend-pay-date": "2023-05-18",
        "dividend-updated-at": "2023-05-08T00:16:36.676Z",
        "earnings": {
          "visible": true,
          "expected-report-date": "2023-07-27",
          "estimated": false,
          "late-flag": 0,
          "quarter-end-date": "2023-06-01",
          "actual-eps": "1.2",
          "consensus-estimate": "1.18",
          "updated-at": "2023-06-07T11:01:36.077Z"
        },
        "listed-market": "XNAS",
        "lendability": "Easy To Borrow",
        "borrow-rate": "0.0",
        "market-cap": 2846108626900,
        "implied-volatility-30-day": "21.87",
        "historical-volatility-30-day": "20.02",
        "historical-volatility-60-day": "19.55",
        "historical-volatility-90-day": "20.24",
        "iv-hv-30-day-difference": "1.85",
        "price-earnings-ratio": "0.0",
        "earnings-per-share": "0.0"
      },
      {
        "symbol": "TSLA",
        "implied-volatility-index": "0.550920788",
        "implied-volatility-index-5-day-change": "0.006567415",
        "implied-volatility-index-rank": "0.12821369",
        "tos-implied-volatility-index-rank": "0.12821369",
        "tw-implied-volatility-index-rank": "0.425296508",
        "tos-implied-volatility-index-rank-updated-at": "2023-06-07T19:46:31.514Z",
        "implied-volatility-index-rank-source": "tos",
        "implied-volatility-percentile": "0.107956139",
        "implied-volatility-updated-at": "2023-06-07T19:46:31.503Z",
        "liquidity-value": "183478.567652843",
        "liquidity-rank": "0.397819778",
        "liquidity-rating": 4,
        "created-at": "2019-01-20T18:02:44.206-06:00",
        "updated-at": "2023-06-07T19:46:32.482Z",
        "option-expiration-implied-volatilities": [
          {
            "expiration-date": "2023-06-09",
            "option-chain-type": "Standard",
            "settlement-type": "PM",
            "implied-volatility": "0.61080228"
          },
          {
            "expiration-date": "2023-06-16",
            "option-chain-type": "Standard",
            "settlement-type": "PM",
            "implied-volatility": "0.562039188"
          },
          {
            "expiration-date": "2023-06-23",
            "option-chain-type": "Standard",
            "settlement-type": "PM",
            "implied-volatility": "0.527234346"
          }
        ],
        "liquidity-running-state": {
          "sum": "7312230.924219939",
          "count": 48,
          "started-at": "2023-06-03T10:00:10.344Z",
          "updated-at": "2023-06-07T19:46:31.503Z"
        },
        "beta": "1.806768437",
        "beta-updated-at": "2023-06-04T17:00:13.353Z",
        "corr-spy-3month": "0.66",
        "dividend-rate-per-share": "0.0",
        "annual-dividend-per-share": "0.0",
        "dividend-yield": "0.0",
        "earnings": {
          "visible": true,
          "expected-report-date": "2023-07-19",
          "estimated": false,
          "late-flag": 0,
          "quarter-end-date": "2023-06-01",
          "actual-eps": "0.65",
          "consensus-estimate": "0.7",
          "updated-at": "2023-06-07T11:02:19.083Z"
        },
        "listed-market": "XNAS",
        "lendability": "Easy To Borrow",
        "borrow-rate": "0.0",
        "market-cap": 678178835284,
        "implied-volatility-30-day": "54.91",
        "historical-volatility-30-day": "36.63",
        "historical-volatility-60-day": "44.98",
        "historical-volatility-90-day": "48.42",
        "iv-hv-30-day-difference": "18.28",
        "price-earnings-ratio": "0.0",
        "earnings-per-share": "0.0"
      }
    ]
  }
}`

const historicDividendsResp = `{
  "data": {
    "items": [
      { "occurred-date": "2023-05-12", "amount": "0.24" },
      { "occurred-date": "2023-02-10", "amount": "0.23" },
      { "occurred-date": "2022-11-04", "amount": "0.23" },
      { "occurred-date": "2022-08-05", "amount": "0.23" },
      { "occurred-date": "2022-05-06", "amount": "0.23" },
      { "occurred-date": "2022-02-04", "amount": "0.22" },
      { "occurred-date": "2021-11-05", "amount": "0.22" },
      { "occurred-date": "2021-08-06", "amount": "0.22" },
      { "occurred-date": "2021-05-07", "amount": "0.22" },
      { "occurred-date": "2021-02-05", "amount": "0.205" },
      { "occurred-date": "2020-11-06", "amount": "0.205" },
      { "occurred-date": "2020-08-07", "amount": "0.205" },
      { "occurred-date": "2020-05-08", "amount": "0.205" },
      { "occurred-date": "2020-02-07", "amount": "0.1925" },
      { "occurred-date": "2019-11-07", "amount": "0.1925" },
      { "occurred-date": "2019-08-09", "amount": "0.1925" },
      { "occurred-date": "2019-05-10", "amount": "0.1925" },
      { "occurred-date": "2019-02-08", "amount": "0.1825" },
      { "occurred-date": "2018-11-08", "amount": "0.1825" },
      { "occurred-date": "2018-08-10", "amount": "0.1825" },
      { "occurred-date": "2018-05-11", "amount": "0.1825" },
      { "occurred-date": "2018-02-09", "amount": "0.1575" },
      { "occurred-date": "2017-11-10", "amount": "0.1575" },
      { "occurred-date": "2017-08-10", "amount": "0.1575" },
      { "occurred-date": "2017-05-11", "amount": "0.1575" },
      { "occurred-date": "2017-02-09", "amount": "0.1425" },
      { "occurred-date": "2016-11-03", "amount": "0.1425" },
      { "occurred-date": "2016-08-04", "amount": "0.1425" },
      { "occurred-date": "2016-05-05", "amount": "0.1425" },
      { "occurred-date": "2016-02-04", "amount": "0.13" },
      { "occurred-date": "2015-11-05", "amount": "0.13" },
      { "occurred-date": "2015-08-06", "amount": "0.13" },
      { "occurred-date": "2015-05-07", "amount": "0.13" },
      { "occurred-date": "2015-02-05", "amount": "0.1175" },
      { "occurred-date": "2014-11-06", "amount": "0.1175" },
      { "occurred-date": "2014-08-07", "amount": "0.1175" },
      { "occurred-date": "2014-05-08", "amount": "0.1175" },
      { "occurred-date": "2014-02-06", "amount": "0.108928571" },
      { "occurred-date": "2013-11-06", "amount": "0.108928571" }
    ]
  }
}`

const historicEarningsResp = `{
  "data": {
    "items": [
      { "occurred-date": "2022-06-30", "eps": "1.2" },
      { "occurred-date": "2022-09-30", "eps": "1.29" },
      { "occurred-date": "2022-12-31", "eps": "1.89" },
      { "occurred-date": "2023-03-31", "eps": "1.53" }
    ]
  }
}`

const marketDataByTypeResp = `{
  "data": {
    "items": [
      {
        "symbol": "BTC",
        "instrument-type": "Equity",
        "updated-at": "2025-09-13T21:33:06.641Z",
        "bid": "51.4",
        "bid-size": "1.0",
        "ask": "51.66",
        "ask-size": "4.0",
        "mid": "51.53",
        "mark": "51.66",
        "last": "51.78",
        "last-mkt": "51.78",
        "beta": "1.719217504",
        "open": "50.94",
        "day-high-price": "51.795",
        "day-low-price": "50.87",
        "close": "51.78",
        "close-price-type": "Final",
        "prev-close": "50.74",
        "prev-close-price-type": "Final",
        "summary-date": "2025-09-12",
        "prev-close-date": "2025-09-11",
        "low-limit-price": "46.3",
        "high-limit-price": "56.59",
        "is-trading-halted": false,
        "halt-start-time": -1,
        "halt-end-time": -1,
        "year-low-price": "25.5",
        "year-high-price": "54.48",
        "volume": "995431.0"
      },
      {
        "symbol": "AAPL",
        "instrument-type": "Equity",
        "updated-at": "2025-09-13T21:35:04.949Z",
        "bid": "233.97",
        "bid-size": "1.0",
        "ask": "234.0",
        "ask-size": "8.0",
        "mid": "233.985",
        "mark": "233.99",
        "last": "233.99",
        "last-ext": "221.8511",
        "last-mkt": "234.07",
        "beta": "1.078586035",
        "dividend-amount": "0.26",
        "dividend-frequency": "4.0",
        "open": "229.22",
        "day-high-price": "234.51",
        "day-low-price": "229.02",
        "close": "234.07",
        "close-price-type": "Final",
        "prev-close": "230.03",
        "prev-close-price-type": "Final",
        "summary-date": "2025-09-12",
        "prev-close-date": "2025-09-11",
        "low-limit-price": "210.46",
        "high-limit-price": "257.23",
        "is-trading-halted": false,
        "trading-halted-reason": "",
        "halt-start-time": -1,
        "halt-end-time": -1,
        "year-low-price": "169.2101",
        "year-high-price": "260.1",
        "volume": "55824216.0"
      },
      {
        "symbol": "AAPL  250919C00005000",
        "instrument-type": "Equity Option",
        "updated-at": "2025-09-13T21:38:31.404Z",
        "bid": "228.7",
        "bid-size": "178.0",
        "ask": "229.5",
        "ask-size": "171.0",
        "mid": "229.1",
        "mark": "229.1",
        "last": "228.97",
        "last-mkt": "228.97",
        "open": "224.35",
        "day-high-price": "229.18",
        "day-low-price": "224.17",
        "close": "228.97",
        "close-price-type": "Regular",
        "prev-close": "225.4",
        "prev-close-price-type": "Regular",
        "summary-date": "2025-09-12",
        "prev-close-date": "2025-09-11",
        "is-trading-halted": false,
        "halt-start-time": -1,
        "halt-end-time": -1,
        "volume": "76.0",
        "volatility": "1.073205822",
        "delta": "1.0",
        "gamma": "0.0",
        "theta": "0.0",
        "rho": "0.000944827",
        "vega": "0.0",
        "theo-price": "228.943960053",
        "dx-mark": "228.943960053",
        "tick-size": "0.01",
        "open-interest": 510
      },
      {
        "symbol": "./ESZ5 EW4V5 251024C5850",
        "instrument-type": "Future Option",
        "updated-at": "2025-09-13T21:35:43.744Z",
        "bid": "747.0",
        "bid-size": "0.0",
        "ask": "754.25",
        "ask-size": "0.0",
        "mid": "750.625",
        "mark": "801.768086004",
        "close": "802.5",
        "close-price-type": "Preliminary",
        "prev-close": "804.75",
        "prev-close-price-type": "Final",
        "summary-date": "2025-09-12",
        "prev-close-date": "2025-09-11",
        "is-trading-halted": false,
        "halt-start-time": -1,
        "halt-end-time": -1,
        "volatility": "0.235349124",
        "delta": "0.943486996",
        "gamma": "0.000196302",
        "theta": "-0.540852829",
        "rho": "6.293186268",
        "vega": "2.347970638",
        "theo-price": "801.768086004",
        "dx-mark": "801.768086004",
        "tick-size": "0.01",
        "open-interest": 0
      },
      {
        "symbol": "SPY",
        "instrument-type": "Equity",
        "updated-at": "2025-09-13T21:35:08.910Z",
        "bid": "657.23",
        "bid-size": "1.0",
        "ask": "657.25",
        "ask-size": "55.0",
        "mid": "657.24",
        "mark": "657.25",
        "last": "657.41",
        "last-mkt": "657.41",
        "beta": "1.007939355",
        "dividend-amount": "1.761117",
        "dividend-frequency": "4.0",
        "open": "657.6",
        "day-high-price": "659.11",
        "day-low-price": "656.9",
        "close": "657.41",
        "close-price-type": "Final",
        "prev-close": "657.63",
        "prev-close-price-type": "Final",
        "summary-date": "2025-09-12",
        "prev-close-date": "2025-09-11",
        "low-limit-price": "591.86",
        "high-limit-price": "723.38",
        "is-trading-halted": false,
        "halt-start-time": -1,
        "halt-end-time": -1,
        "year-low-price": "481.8",
        "year-high-price": "659.11",
        "volume": "72780135.0"
      },
      {
        "symbol": "ES",
        "instrument-type": "Equity",
        "updated-at": "2025-09-13T21:33:44.722Z",
        "bid": "64.99",
        "bid-size": "1.0",
        "ask": "66.24",
        "ask-size": "5.0",
        "mid": "65.615",
        "mark": "65.7",
        "last": "65.7",
        "last-mkt": "65.7",
        "beta": "0.681157943",
        "dividend-amount": "0.7525",
        "dividend-frequency": "4.0",
        "open": "64.79",
        "day-high-price": "65.835",
        "day-low-price": "64.6",
        "close": "65.7",
        "close-price-type": "Final",
        "prev-close": "65.05",
        "prev-close-price-type": "Final",
        "summary-date": "2025-09-12",
        "prev-close-date": "2025-09-11",
        "low-limit-price": "58.89",
        "high-limit-price": "71.98",
        "is-trading-halted": false,
        "halt-start-time": -1,
        "halt-end-time": -1,
        "year-low-price": "52.28",
        "year-high-price": "68.73",
        "volume": "2169037.0"
      }
    ]
  },
  "pagination": {
    "per-page": 4,
    "page-offset": 0,
    "item-offset": 0,
    "total-items": 12143,
    "total-pages": 3036,
    "current-item-count": 4,
    "previous-link": null,
    "next-link": null,
    "paging-link-template": null
  }
}`
