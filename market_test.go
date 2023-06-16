package tasty //nolint:testpackage // testing private field

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestGetMarketMetrics(t *testing.T) {
	setup()
	defer teardown()

	symbols := []string{"AAPL", "TSLA"}

	mux.HandleFunc("/market-metrics", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, marketMetricsResp)
	})

	resp, err := client.GetMarketMetrics(symbols)
	require.Nil(t, err)

	require.Equal(t, 2, len(resp))

	aapl := resp[0]

	require.Equal(t, StringToFloat32(0.224186663), aapl.ImpliedVolatilityIndex)
	require.Equal(t, StringToFloat32(-0.014691909), aapl.ImpliedVolatilityIndex5DayChange)
	require.Equal(t, StringToFloat32(0.02231313), aapl.ImpliedVolatilityIndexRank)
	require.Equal(t, StringToFloat32(0.02231313), aapl.TosImpliedVolatilityIndexRank)
	require.Equal(t, StringToFloat32(0.037374278), aapl.TwImpliedVolatilityIndexRank)
	require.Equal(t, time.Date(2023, time.June, 7, 19, 55, 22, 151000000, time.UTC), aapl.TosImpliedVolatilityIndexRankUpdatedAt)
	require.Equal(t, "tos", aapl.ImpliedVolatilityIndexRankSource)
	require.Equal(t, StringToFloat32(0.031249801), aapl.ImpliedVolatilityPercentile)
	require.Equal(t, time.Date(2023, time.June, 7, 19, 55, 22, 110000000, time.UTC), aapl.ImpliedVolatilityUpdatedAt)
	require.Equal(t, StringToFloat32(899091.420242437), aapl.LiquidityValue)
	require.Equal(t, StringToFloat32(1.516319123), aapl.LiquidityRank)
	require.Equal(t, 4, aapl.LiquidityRating)
	require.Equal(t, time.Date(2023, time.June, 7, 19, 55, 23, 376000000, time.UTC), aapl.UpdatedAt)
	require.Equal(t, StringToFloat32(1.265634436), aapl.Beta)
	require.Equal(t, time.Date(2023, time.June, 4, 17, 0, 35, 805000000, time.UTC), aapl.BetaUpdatedAt)
	require.Equal(t, StringToFloat32(0.79), aapl.CorrSpy3month)
	require.Equal(t, StringToFloat32(0.24), aapl.DividendRatePerShare)
	require.Equal(t, StringToFloat32(0.0), aapl.AnnualDividendPerShare)
	require.Equal(t, StringToFloat32(0.006396146), aapl.DividendYield)
	require.Equal(t, "2023-05-12", aapl.DividendExDate)
	require.Equal(t, "2022-08-05", aapl.DividendNextDate)
	require.Equal(t, "2023-05-18", aapl.DividendPayDate)
	require.Equal(t, time.Date(2023, time.May, 8, 0, 16, 36, 676000000, time.UTC), aapl.DividendUpdatedAt)
	require.Equal(t, "XNAS", aapl.ListedMarket)
	require.Equal(t, "Easy To Borrow", aapl.Lendability)
	require.Equal(t, StringToFloat32(0.0), aapl.BorrowRate)
	require.Equal(t, 2846108626900, aapl.MarketCap)
	require.Equal(t, StringToFloat32(21.87), aapl.ImpliedVolatility30Day)
	require.Equal(t, StringToFloat32(20.02), aapl.HistoricalVolatility30Day)
	require.Equal(t, StringToFloat32(19.55), aapl.HistoricalVolatility60Day)
	require.Equal(t, StringToFloat32(20.24), aapl.HistoricalVolatility90Day)
	require.Equal(t, StringToFloat32(1.85), aapl.IvHv30DayDifference)
	require.Equal(t, StringToFloat32(0.0), aapl.PriceEarningsRatio)
	require.Equal(t, StringToFloat32(0.0), aapl.EarningsPerShare)

	// Option Expiration IVs
	optExpiryIVs := aapl.OptionExpirationImpliedVolatilities

	require.Equal(t, 3, len(optExpiryIVs))
	require.Equal(t, "2023-06-09", optExpiryIVs[0].ExpirationDate)
	require.Equal(t, "Standard", optExpiryIVs[0].OptionChainType)
	require.Equal(t, "PM", optExpiryIVs[0].SettlementType)
	require.Equal(t, StringToFloat32(0.288187496), optExpiryIVs[0].ImpliedVolatility)

	// Liquidity Running State
	liquidityRunningState := aapl.LiquidityRunningState

	require.Equal(t, StringToFloat32(38347762.009590656), liquidityRunningState.Sum)
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
	require.Equal(t, StringToFloat32(1.2), earnings.ActualEPS)
	require.Equal(t, StringToFloat32(1.18), earnings.ConsensusEstimate)
	require.Equal(t, time.Date(2023, time.June, 7, 11, 1, 36, 77000000, time.UTC), earnings.UpdatedAt)
}

func TestGetHistoricDividends(t *testing.T) {
	setup()
	defer teardown()

	symbol := "AAPL"

	mux.HandleFunc(fmt.Sprintf("/market-metrics/historic-corporate-events/dividends/%s", symbol), func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, historicDividendsResp)
	})

	resp, err := client.GetHistoricDividends(symbol)
	require.Nil(t, err)

	require.Equal(t, 39, len(resp))
	require.Equal(t, "2023-05-12", resp[0].OccurredDate)
	require.Equal(t, StringToFloat32(0.24), resp[0].Amount)
}

func TestGetHistoricEarnings(t *testing.T) {
	setup()
	defer teardown()

	startDate := time.Now().AddDate(-1, 0, 0)
	symbol := "AAPL"

	mux.HandleFunc(fmt.Sprintf("/market-metrics/historic-corporate-events/earnings-reports/%s", symbol), func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, historicEarningsResp)
	})

	resp, err := client.GetHistoricEarnings(symbol, startDate)
	require.Nil(t, err)

	require.Equal(t, 4, len(resp))
	require.Equal(t, "2022-06-30", resp[0].OccurredDate)
	require.Equal(t, StringToFloat32(1.2), resp[0].Eps)
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
