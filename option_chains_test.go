package tasty

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/austinbspencer/tasty-go/constants"
	"github.com/austinbspencer/tasty-go/models"
	"github.com/stretchr/testify/require"
)

func TestGetFuturesOptionChains(t *testing.T) {
	setup()
	defer teardown()

	productCode := "ES"

	mux.HandleFunc(fmt.Sprintf("/futures-option-chains/%s", productCode), func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, futuresOptionChainsResp)
	})

	resp, err := client.GetFuturesOptionChains(productCode)
	require.Nil(t, err)

	require.Equal(t, 1, len(resp))

	fo := resp[0]

	require.Equal(t, "./ESU3 EW4N3 230728C3990", fo.Symbol)
	require.Equal(t, "/ESU3", fo.UnderlyingSymbol)
	require.Equal(t, "ES", fo.ProductCode)
	require.Equal(t, "2023-07-28", fo.ExpirationDate)
	require.Equal(t, "/ES", fo.RootSymbol)
	require.Equal(t, "EW4", fo.OptionRootSymbol)
	require.Equal(t, models.StringToFloat32(3990), fo.StrikePrice)
	require.Equal(t, "CME", fo.Exchange)
	require.Equal(t, "EW4N3 C3990", fo.ExchangeSymbol)
	require.Equal(t, string(constants.Call), fo.OptionType)
	require.Equal(t, "American", fo.ExerciseStyle)
	require.True(t, fo.IsVanilla)
	require.True(t, fo.IsPrimaryDeliverable)
	require.Equal(t, models.StringToFloat32(1), fo.FuturePriceRatio)
	require.Equal(t, models.StringToFloat32(1), fo.Multiplier)
	require.Equal(t, models.StringToFloat32(1), fo.UnderlyingCount)
	require.True(t, fo.IsConfirmed)
	require.Equal(t, models.StringToFloat32(.5), fo.NotionalValue)
	require.Equal(t, models.StringToFloat32(.01), fo.DisplayFactor)
	require.Equal(t, "2", fo.SecurityExchange)
	require.Equal(t, "0", fo.SxID)
	require.Equal(t, "Future", fo.SettlementType)
	require.Equal(t, models.StringToFloat32(1), fo.StrikeFactor)
	require.Equal(t, "2023-07-28", fo.MaturityDate)
	require.True(t, fo.IsExercisableWeekly)
	require.Equal(t, "0", fo.LastTradeTime)
	require.Equal(t, 47, fo.DaysToExpiration)
	require.False(t, fo.IsClosingOnly)
	require.True(t, fo.Active)
	require.Equal(t, "2023-07-28T20:00:00Z", fo.StopsTradingAt.Format(time.RFC3339))
	require.Equal(t, "2023-07-28T20:00:00Z", fo.ExpiresAt.Format(time.RFC3339))

	fop := fo.FutureOptionProduct

	require.Equal(t, "EW4", fop.RootSymbol)
	require.False(t, fop.CashSettled)
	require.Equal(t, "EW4", fop.Code)
	require.Equal(t, "EW4", fop.LegacyCode)
	require.Equal(t, "EW4", fop.ClearportCode)
	require.Equal(t, "W4", fop.ClearingCode)
	require.Equal(t, "9C", fop.ClearingExchangeCode)
	require.Equal(t, models.StringToFloat32(1), fop.ClearingPriceMultiplier)
	require.Equal(t, models.StringToFloat32(.01), fop.DisplayFactor)
	require.Equal(t, "CME", fop.Exchange)
	require.Equal(t, "Physical", fop.ProductType)
	require.Equal(t, "Weekly", fop.ExpirationType)
	require.Equal(t, 0, fop.SettlementDelayDays)
	require.True(t, fop.IsRollover)
	require.Equal(t, "Equity Index", fop.MarketSector)
}

func TestGetNestedFuturesOptionChains(t *testing.T) {
	setup()
	defer teardown()

	productCode := "ES"

	mux.HandleFunc(fmt.Sprintf("/futures-option-chains/%s/nested", productCode), func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, futuresOptionChainsNested)
	})

	resp, err := client.GetNestedFuturesOptionChains(productCode)
	require.Nil(t, err)

	nestedFuture := resp.Futures[0]

	require.Equal(t, "/ESM3", nestedFuture.Symbol)
	require.Equal(t, "/ES", nestedFuture.RootSymbol)
	require.Equal(t, "2023-06-16", nestedFuture.ExpirationDate)
	require.Equal(t, 5, nestedFuture.DaysToExpiration)
	require.True(t, nestedFuture.ActiveMonth)
	require.False(t, nestedFuture.NextActiveMonth)
	require.Equal(t, "2023-06-16T13:30:00Z", nestedFuture.StopsTradingAt.Format(time.RFC3339))
	require.Equal(t, "2023-06-16T13:30:00Z", nestedFuture.ExpiresAt.Format(time.RFC3339))

	oc := resp.OptionChains[0]

	require.Equal(t, "/ES", oc.UnderlyingSymbol)
	require.Equal(t, "/ES", oc.RootSymbol)
	require.Equal(t, "American", oc.ExerciseStyle)

	exp := oc.Expirations[0]

	require.Equal(t, "/ESU3", exp.UnderlyingSymbol)
	require.Equal(t, "/ES", exp.RootSymbol)
	require.Equal(t, "EW4", exp.OptionRootSymbol)
	require.Equal(t, "EW4N3", exp.OptionContractSymbol)
	require.Equal(t, "EW4", exp.Asset)
	require.Equal(t, "2023-07-28", exp.ExpirationDate)
	require.Equal(t, 47, exp.DaysToExpiration)
	require.Equal(t, "Weekly", exp.ExpirationType)
	require.Equal(t, "PM", exp.SettlementType)
	require.Equal(t, models.StringToFloat32(0.5), exp.NotionalValue)
	require.Equal(t, models.StringToFloat32(0.01), exp.DisplayFactor)
	require.Equal(t, models.StringToFloat32(1), exp.StrikeFactor)
	require.Equal(t, "2023-07-28T20:00:00Z", exp.StopsTradingAt.Format(time.RFC3339))
	require.Equal(t, "2023-07-28T20:00:00Z", exp.ExpiresAt.Format(time.RFC3339))

	tick := exp.TickSizes[0]

	require.Equal(t, models.StringToFloat32(0.05), tick.Value)
	require.Equal(t, models.StringToFloat32(5), tick.Threshold)

	strike := exp.Strikes[0]

	require.Equal(t, models.StringToFloat32(3990), strike.StrikePrice)
	require.Equal(t, "./ESU3 EW4N3 230728C3990", strike.Call)
	require.Equal(t, "./EW4N23C3990:XCME", strike.CallStreamerSymbol)
	require.Equal(t, "./ESU3 EW4N3 230728P3990", strike.Put)
	require.Equal(t, "./EW4N23P3990:XCME", strike.PutStreamerSymbol)
}

const futuresOptionChainsResp = `{
  "data": {
    "items": [
      {
        "symbol": "./ESU3 EW4N3 230728C3990",
        "underlying-symbol": "/ESU3",
        "product-code": "ES",
        "expiration-date": "2023-07-28",
        "root-symbol": "/ES",
        "option-root-symbol": "EW4",
        "strike-price": "3990.0",
        "exchange": "CME",
        "exchange-symbol": "EW4N3 C3990",
        "streamer-symbol": "./EW4N23C3990:XCME",
        "option-type": "C",
        "exercise-style": "American",
        "is-vanilla": true,
        "is-primary-deliverable": true,
        "future-price-ratio": "1.0",
        "multiplier": "1.0",
        "underlying-count": "1.0",
        "is-confirmed": true,
        "notional-value": "0.5",
        "display-factor": "0.01",
        "security-exchange": "2",
        "sx-id": "0",
        "settlement-type": "Future",
        "strike-factor": "1.0",
        "maturity-date": "2023-07-28",
        "is-exercisable-weekly": true,
        "last-trade-time": "0",
        "days-to-expiration": 47,
        "is-closing-only": false,
        "active": true,
        "stops-trading-at": "2023-07-28T20:00:00.000+00:00",
        "expires-at": "2023-07-28T20:00:00.000+00:00",
        "future-option-product": {
          "root-symbol": "EW4",
          "cash-settled": false,
          "code": "EW4",
          "legacy-code": "EW4",
          "clearport-code": "EW4",
          "clearing-code": "W4",
          "clearing-exchange-code": "9C",
          "clearing-price-multiplier": "1.0",
          "display-factor": "0.01",
          "exchange": "CME",
          "product-type": "Physical",
          "expiration-type": "Weekly",
          "settlement-delay-days": 0,
          "is-rollover": true,
          "market-sector": "Equity Index"
        }
      }
    ]
  },
  "context": "/futures-option-chains/ES"
}`

const futuresOptionChainsNested = `{
  "data": {
    "futures": [
      {
        "symbol": "/ESM3",
        "root-symbol": "/ES",
        "expiration-date": "2023-06-16",
        "days-to-expiration": 5,
        "active-month": true,
        "next-active-month": false,
        "stops-trading-at": "2023-06-16T13:30:00.000+00:00",
        "expires-at": "2023-06-16T13:30:00.000+00:00"
      },
      {
        "symbol": "/ESU3",
        "root-symbol": "/ES",
        "expiration-date": "2023-09-15",
        "days-to-expiration": 96,
        "active-month": false,
        "next-active-month": true,
        "stops-trading-at": "2023-09-15T13:30:00.000+00:00",
        "expires-at": "2023-09-15T13:30:00.000+00:00"
      },
      {
        "symbol": "/ESZ3",
        "root-symbol": "/ES",
        "expiration-date": "2023-12-15",
        "days-to-expiration": 187,
        "active-month": false,
        "next-active-month": false,
        "stops-trading-at": "2023-12-15T14:30:00.000+00:00",
        "expires-at": "2023-12-15T14:30:00.000+00:00"
      },
      {
        "symbol": "/ESZ7",
        "root-symbol": "/ES",
        "expiration-date": "2027-12-17",
        "days-to-expiration": 1650,
        "active-month": false,
        "next-active-month": false,
        "stops-trading-at": "2027-12-17T14:30:00.000+00:00",
        "expires-at": "2027-12-17T14:30:00.000+00:00"
      }
    ],
    "option-chains": [
      {
        "underlying-symbol": "/ES",
        "root-symbol": "/ES",
        "exercise-style": "American",
        "expirations": [
          {
            "underlying-symbol": "/ESU3",
            "root-symbol": "/ES",
            "option-root-symbol": "EW4",
            "option-contract-symbol": "EW4N3",
            "asset": "EW4",
            "expiration-date": "2023-07-28",
            "days-to-expiration": 47,
            "expiration-type": "Weekly",
            "settlement-type": "PM",
            "notional-value": "0.5",
            "display-factor": "0.01",
            "strike-factor": "1.0",
            "stops-trading-at": "2023-07-28T20:00:00.000+00:00",
            "expires-at": "2023-07-28T20:00:00.000+00:00",
            "tick-sizes": [
              {
                "value": "0.05",
                "threshold": "5.0"
              },
              {
                "value": "0.25"
              }
            ],
            "strikes": [
              {
                "strike-price": "3990.0",
                "call": "./ESU3 EW4N3 230728C3990",
                "call-streamer-symbol": "./EW4N23C3990:XCME",
                "put": "./ESU3 EW4N3 230728P3990",
                "put-streamer-symbol": "./EW4N23P3990:XCME"
              },
              {
                "strike-price": "4530.0",
                "call": "./ESU3 EW4N3 230728C4530",
                "call-streamer-symbol": "./EW4N23C4530:XCME",
                "put": "./ESU3 EW4N3 230728P4530",
                "put-streamer-symbol": "./EW4N23P4530:XCME"
              },
              {
                "strike-price": "4610.0",
                "call": "./ESU3 EW4N3 230728C4610",
                "call-streamer-symbol": "./EW4N23C4610:XCME",
                "put": "./ESU3 EW4N3 230728P4610",
                "put-streamer-symbol": "./EW4N23P4610:XCME"
              }
            ]
          }
        ]
      }
    ]
  },
  "context": "/futures-option-chains/ES/nested"
}`
