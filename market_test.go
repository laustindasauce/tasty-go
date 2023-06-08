package tasty

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/austinbspencer/tasty-go/models"
	"github.com/stretchr/testify/require"
)

func TestGetMarketMetrics(t *testing.T) {
	setup()
	defer teardown()

	query := models.MarketMetricsQuery{Symbols: []string{"AAPL", "TSLA"}}

	mux.HandleFunc("/market-metrics", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, marketMetricsResp)
	})

	resp, err := client.GetMarketMetrics(query)
	require.Nil(t, err)

	require.Equal(t, 2, len(resp))

	aapl := resp[0]

	require.Equal(t, models.StringToFloat32(0.224186663), aapl.ImpliedVolatilityIndex)
	require.Equal(t, models.StringToFloat32(-0.014691909), aapl.ImpliedVolatilityIndex5DayChange)
	require.Equal(t, models.StringToFloat32(0.02231313), aapl.ImpliedVolatilityIndexRank)
	require.Equal(t, models.StringToFloat32(0.02231313), aapl.TosImpliedVolatilityIndexRank)
	require.Equal(t, models.StringToFloat32(0.037374278), aapl.TwImpliedVolatilityIndexRank)
	require.Equal(t, time.Time(time.Date(2023, time.June, 7, 19, 55, 22, 151000000, time.UTC)), aapl.TosImpliedVolatilityIndexRankUpdatedAt)
	require.Equal(t, "tos", aapl.ImpliedVolatilityIndexRankSource)
	require.Equal(t, models.StringToFloat32(0.031249801), aapl.ImpliedVolatilityPercentile)
	require.Equal(t, time.Time(time.Date(2023, time.June, 7, 19, 55, 22, 110000000, time.UTC)), aapl.ImpliedVolatilityUpdatedAt)
	require.Equal(t, models.StringToFloat32(899091.420242437), aapl.LiquidityValue)
	require.Equal(t, models.StringToFloat32(1.516319123), aapl.LiquidityRank)
	require.Equal(t, 4, aapl.LiquidityRating)
	require.Equal(t, time.Time(time.Date(2023, time.June, 7, 19, 55, 23, 376000000, time.UTC)), aapl.UpdatedAt)
	require.Equal(t, models.StringToFloat32(1.265634436), aapl.Beta)
	require.Equal(t, time.Time(time.Date(2023, time.June, 4, 17, 0, 35, 805000000, time.UTC)), aapl.BetaUpdatedAt)
	require.Equal(t, models.StringToFloat32(0.79), aapl.CorrSpy3month)
	require.Equal(t, models.StringToFloat32(0.24), aapl.DividendRatePerShare)
	require.Equal(t, models.StringToFloat32(0.0), aapl.AnnualDividendPerShare)
	require.Equal(t, models.StringToFloat32(0.006396146), aapl.DividendYield)
	require.Equal(t, "2023-05-12", aapl.DividendExDate)
	require.Equal(t, "2022-08-05", aapl.DividendNextDate)
	require.Equal(t, "2023-05-18", aapl.DividendPayDate)
	require.Equal(t, time.Time(time.Date(2023, time.May, 8, 0, 16, 36, 676000000, time.UTC)), aapl.DividendUpdatedAt)
	require.Equal(t, "XNAS", aapl.ListedMarket)
	require.Equal(t, "Easy To Borrow", aapl.Lendability)
	require.Equal(t, models.StringToFloat32(0.0), aapl.BorrowRate)
	require.Equal(t, 2846108626900, aapl.MarketCap)
	require.Equal(t, models.StringToFloat32(21.87), aapl.ImpliedVolatility30Day)
	require.Equal(t, models.StringToFloat32(20.02), aapl.HistoricalVolatility30Day)
	require.Equal(t, models.StringToFloat32(19.55), aapl.HistoricalVolatility60Day)
	require.Equal(t, models.StringToFloat32(20.24), aapl.HistoricalVolatility90Day)
	require.Equal(t, models.StringToFloat32(1.85), aapl.IvHv30DayDifference)
	require.Equal(t, models.StringToFloat32(0.0), aapl.PriceEarningsRatio)
	require.Equal(t, models.StringToFloat32(0.0), aapl.EarningsPerShare)

	// Option Expiration IVs
	optExpiryIVs := aapl.OptionExpirationImpliedVolatilities

	require.Equal(t, 3, len(optExpiryIVs))
	require.Equal(t, "2023-06-09", optExpiryIVs[0].ExpirationDate)
	require.Equal(t, "Standard", optExpiryIVs[0].OptionChainType)
	require.Equal(t, "PM", optExpiryIVs[0].SettlementType)
	require.Equal(t, models.StringToFloat32(0.288187496), optExpiryIVs[0].ImpliedVolatility)

	// Liquidity Running State
	liquidityRunningState := aapl.LiquidityRunningState

	require.Equal(t, models.StringToFloat32(38347762.009590656), liquidityRunningState.Sum)
	require.Equal(t, 48, liquidityRunningState.Count)
	require.Equal(t, time.Time(time.Date(2023, time.June, 3, 10, 0, 10, 347000000, time.UTC)), liquidityRunningState.StartedAt)
	require.Equal(t, time.Time(time.Date(2023, time.June, 7, 19, 55, 22, 110000000, time.UTC)), liquidityRunningState.UpdatedAt)

	// Earnings
	earnings := aapl.Earnings

	require.True(t, earnings.Visible)
	require.Equal(t, "2023-07-27", earnings.ExpectedReportDate)
	require.False(t, earnings.Estimated)
	require.Zero(t, earnings.LateFlag)
	require.Equal(t, "2023-06-01", earnings.QuarterEndDate)
	require.Equal(t, models.StringToFloat32(1.2), earnings.ActualEPS)
	require.Equal(t, models.StringToFloat32(1.18), earnings.ConsensusEstimate)
	require.Equal(t, time.Date(2023, time.June, 7, 11, 1, 36, 77000000, time.UTC), earnings.UpdatedAt)
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
