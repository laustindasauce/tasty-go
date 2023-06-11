package tasty

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/austinbspencer/tasty-go/constants"
	"github.com/austinbspencer/tasty-go/models"
	"github.com/austinbspencer/tasty-go/utils"
	"github.com/stretchr/testify/require"
)

func TestGetCryptocurrencies(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/instruments/cryptocurrencies", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, cryptocurrenciesResp)
	})

	resp, err := client.GetCryptocurrencies([]string{"BTC/USD", "ETH/USD"})
	require.Nil(t, err)

	require.Equal(t, 2, len(resp))

	btc := resp[0]

	require.Equal(t, 1, btc.ID)
	require.Equal(t, "BTC/USD", btc.Symbol)
	require.Equal(t, "Cryptocurrency", btc.InstrumentType)
	require.Equal(t, "Bitcoin", btc.ShortDescription)
	require.Equal(t, "Bitcoin to USD", btc.Description)
	require.False(t, btc.IsClosingOnly)
	require.True(t, btc.Active)
	require.Equal(t, models.StringToFloat32(0.01), btc.TickSize)
	require.Equal(t, "BTC/USD:CXTALP", btc.StreamerSymbol)

	venueSymbol := btc.DestinationVenueSymbols[0]

	require.Equal(t, 71, venueSymbol.ID)
	require.Equal(t, "BTC", venueSymbol.Symbol)
	require.Equal(t, "CITADEL_CRYPTOCURRENCY", venueSymbol.DestinationVenue)
	require.Equal(t, 8, venueSymbol.MaxQuantityPrecision)
	require.Equal(t, 8, venueSymbol.MaxPricePrecision)
	require.True(t, venueSymbol.Routable)
}

func TestGetActiveEquities(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/instruments/equities/active", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, activeEquitiesResp)
	})

	resp, pagination, err := client.GetActiveEquities(models.ActiveEquitiesQuery{PerPage: 4})
	require.Nil(t, err)

	require.Equal(t, 4, len(resp))

	equity := resp[0]

	require.Equal(t, 15625396, equity.ID)
	require.Equal(t, "YMM", equity.Symbol)
	require.Equal(t, "Equity", equity.InstrumentType)
	require.Equal(t, "35969L108", equity.Cusip)
	require.Equal(t, "Full Truck Alliance Co. Ltd. American Depositary Shares (each representing 20 Class A Ordinary Shares)", equity.ShortDescription)
	require.False(t, equity.IsIndex)
	require.Equal(t, "XNYS", equity.ListedMarket)
	require.Equal(t, "Full Truck Alliance Co. Ltd. American Depositary Shares (each representing 20 Class A Ordinary Shares)", equity.Description)
	require.Equal(t, "Easy To Borrow", equity.Lendability)
	require.Equal(t, models.StringToFloat32(0.0), equity.BorrowRate)
	require.Equal(t, "Equity", equity.MarketTimeInstrumentCollection)
	require.False(t, equity.IsClosingOnly)
	require.False(t, equity.IsOptionsClosingOnly)
	require.True(t, equity.Active)
	require.False(t, equity.IsFractionalQuantityEligible)
	require.False(t, equity.IsIlliquid)
	require.False(t, equity.IsEtf)
	require.Equal(t, "YMM", equity.StreamerSymbol)

	tickSize := equity.TickSizes[0]

	require.Equal(t, models.StringToFloat32(0.0001), tickSize.Value)
	require.Equal(t, models.StringToFloat32(1.0), tickSize.Threshold)

	optionTickSize := equity.OptionTickSizes[0]

	require.Equal(t, models.StringToFloat32(0.05), optionTickSize.Value)
	require.Equal(t, models.StringToFloat32(3.0), optionTickSize.Threshold)

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

func TestGetEquities(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/instruments/equities", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, equitiesResp)
	})

	resp, err := client.GetEquities(models.EquitiesQuery{Symbols: []string{"AAPL", "TSLA"}})
	require.Nil(t, err)

	require.Equal(t, 2, len(resp))

	equity := resp[0]

	require.Equal(t, 726, equity.ID)
	require.Equal(t, "AAPL", equity.Symbol)
	require.Equal(t, "Equity", equity.InstrumentType)
	require.Equal(t, "037833100", equity.Cusip)
	require.Equal(t, "APPLE INC", equity.ShortDescription)
	require.False(t, equity.IsIndex)
	require.Equal(t, "XNAS", equity.ListedMarket)
	require.Equal(t, "APPLE INC", equity.Description)
	require.Equal(t, "Easy To Borrow", equity.Lendability)
	require.Equal(t, models.StringToFloat32(0.0), equity.BorrowRate)
	require.Equal(t, "Equity", equity.MarketTimeInstrumentCollection)
	require.False(t, equity.IsClosingOnly)
	require.False(t, equity.IsOptionsClosingOnly)
	require.True(t, equity.Active)
	require.True(t, equity.IsFractionalQuantityEligible)
	require.False(t, equity.IsIlliquid)
	require.False(t, equity.IsEtf)
	require.Equal(t, "AAPL", equity.StreamerSymbol)

	tickSize := equity.TickSizes[0]

	require.Equal(t, models.StringToFloat32(0.0001), tickSize.Value)
	require.Equal(t, models.StringToFloat32(1.0), tickSize.Threshold)

	optionTickSize := equity.OptionTickSizes[0]

	require.Equal(t, models.StringToFloat32(0.01), optionTickSize.Value)
	require.Equal(t, models.StringToFloat32(3.0), optionTickSize.Threshold)
}

func TestGetEquity(t *testing.T) {
	setup()
	defer teardown()

	symbol := "AAPL"

	mux.HandleFunc(fmt.Sprintf("/instruments/equities/%s", symbol), func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, equityResp)
	})

	equity, err := client.GetEquity(symbol)
	require.Nil(t, err)

	require.Equal(t, 726, equity.ID)
	require.Equal(t, symbol, equity.Symbol)
	require.Equal(t, "Equity", equity.InstrumentType)
	require.Equal(t, "037833100", equity.Cusip)
	require.Equal(t, "APPLE INC", equity.ShortDescription)
	require.False(t, equity.IsIndex)
	require.Equal(t, "XNAS", equity.ListedMarket)
	require.Equal(t, "APPLE INC", equity.Description)
	require.Equal(t, "Easy To Borrow", equity.Lendability)
	require.Equal(t, models.StringToFloat32(0.0), equity.BorrowRate)
	require.Equal(t, "Equity", equity.MarketTimeInstrumentCollection)
	require.False(t, equity.IsClosingOnly)
	require.False(t, equity.IsOptionsClosingOnly)
	require.True(t, equity.Active)
	require.True(t, equity.IsFractionalQuantityEligible)
	require.False(t, equity.IsIlliquid)
	require.False(t, equity.IsEtf)
	require.Equal(t, symbol, equity.StreamerSymbol)

	tickSize := equity.TickSizes[0]

	require.Equal(t, models.StringToFloat32(0.0001), tickSize.Value)
	require.Equal(t, models.StringToFloat32(1.0), tickSize.Threshold)

	optionTickSize := equity.OptionTickSizes[0]

	require.Equal(t, models.StringToFloat32(0.01), optionTickSize.Value)
	require.Equal(t, models.StringToFloat32(3.0), optionTickSize.Threshold)
}

func TestGetEquityOptions(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/instruments/equity-options", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, equityOptionsResp)
	})

	symbol := "AAPL"
	optionType := constants.Call

	sym := utils.OCCSymbology{
		Symbol:     symbol,
		Strike:     185,
		OptionType: optionType,
		Expiration: time.Date(2023, 6, 16, 0, 0, 0, 0, time.UTC),
	}
	occSymbol := sym.GetOCCSymbology()

	resp, err := client.GetEquityOptions(models.EquityOptionsQuery{Symbols: []string{occSymbol}})
	require.Nil(t, err)

	require.Equal(t, 1, len(resp))

	equity := resp[0]

	require.Equal(t, occSymbol, equity.Symbol)
	require.Equal(t, "Equity Option", equity.InstrumentType)
	require.True(t, equity.Active)
	require.Equal(t, models.StringToFloat32(185.0), equity.StrikePrice)
	require.Equal(t, symbol, equity.RootSymbol)
	require.Equal(t, symbol, equity.UnderlyingSymbol)
	require.Equal(t, "2023-06-16", equity.ExpirationDate)
	require.Equal(t, "American", equity.ExerciseStyle)
	require.Equal(t, 100, equity.SharesPerContract)
	require.Equal(t, string(optionType), equity.OptionType)
	require.Equal(t, "Standard", equity.OptionChainType)
	require.Equal(t, "Regular", equity.ExpirationType)
	require.Equal(t, "PM", equity.SettlementType)
	require.Equal(t, "2023-06-16T20:00:00Z", equity.StopsTradingAt.Format(time.RFC3339))
	require.Equal(t, "Equity Option", equity.MarketTimeInstrumentCollection)
	require.Equal(t, 6, equity.DaysToExpiration)
	require.Equal(t, "2023-06-16T20:00:00Z", equity.ExpiresAt.Format(time.RFC3339))
	require.False(t, equity.IsClosingOnly)
	require.Equal(t, ".AAPL230616C185", equity.StreamerSymbol)
}

func TestGetEquityOption(t *testing.T) {
	setup()
	defer teardown()

	symbol := "AAPL"
	optionType := constants.Call

	sym := utils.OCCSymbology{
		Symbol:     symbol,
		Strike:     185,
		OptionType: optionType,
		Expiration: time.Date(2023, 6, 16, 0, 0, 0, 0, time.UTC),
	}
	occSymbol := sym.GetOCCSymbology()

	mux.HandleFunc(fmt.Sprintf("/instruments/equity-options/%s", occSymbol), func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, equityOptionResp)
	})

	equity, err := client.GetEquityOption(sym, true)
	require.Nil(t, err)

	require.Equal(t, occSymbol, equity.Symbol)
	require.Equal(t, "Equity Option", equity.InstrumentType)
	require.True(t, equity.Active)
	require.Equal(t, models.StringToFloat32(185.0), equity.StrikePrice)
	require.Equal(t, symbol, equity.RootSymbol)
	require.Equal(t, symbol, equity.UnderlyingSymbol)
	require.Equal(t, "2023-06-16", equity.ExpirationDate)
	require.Equal(t, "American", equity.ExerciseStyle)
	require.Equal(t, 100, equity.SharesPerContract)
	require.Equal(t, string(optionType), equity.OptionType)
	require.Equal(t, "Standard", equity.OptionChainType)
	require.Equal(t, "Regular", equity.ExpirationType)
	require.Equal(t, "PM", equity.SettlementType)
	require.Equal(t, "2023-06-16T20:00:00Z", equity.StopsTradingAt.Format(time.RFC3339))
	require.Equal(t, "Equity Option", equity.MarketTimeInstrumentCollection)
	require.Equal(t, 6, equity.DaysToExpiration)
	require.Equal(t, "2023-06-16T20:00:00Z", equity.ExpiresAt.Format(time.RFC3339))
	require.False(t, equity.IsClosingOnly)
	require.Equal(t, ".AAPL230616C185", equity.StreamerSymbol)
}

func TestGetFutures(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/instruments/futures", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, futuresResp)
	})

	productCode := "ES"

	resp, err := client.GetFutures(models.FuturesQuery{ProductCode: []string{productCode}})
	require.Nil(t, err)

	require.Equal(t, 1, len(resp))

	future := resp[0]

	require.Equal(t, "/ESM3", future.Symbol)
	require.Equal(t, productCode, future.ProductCode)
	require.Equal(t, models.StringToFloat32(50), future.ContractSize)
	require.Equal(t, models.StringToFloat32(.25), future.TickSize)
	require.Equal(t, models.StringToFloat32(50), future.NotionalMultiplier)
	require.Equal(t, models.StringToFloat32(0), future.MainFraction)
	require.Equal(t, models.StringToFloat32(0), future.SubFraction)
	require.Equal(t, models.StringToFloat32(0.01), future.DisplayFactor)
	require.Equal(t, "2023-06-16", future.LastTradeDate)
	require.Equal(t, "2023-06-16", future.ExpirationDate)
	require.Equal(t, "2023-06-16", future.ClosingOnlyDate)
	require.True(t, future.Active)
	require.True(t, future.ActiveMonth)
	require.False(t, future.NextActiveMonth)
	require.False(t, future.IsClosingOnly)
	require.Equal(t, "2023-06-16T13:30:00Z", future.StopsTradingAt.Format(time.RFC3339Nano))
	require.Equal(t, "2023-06-16T13:30:00Z", future.ExpiresAt.Format(time.RFC3339Nano))
	require.Equal(t, "CME_ES", future.ProductGroup)
	require.Equal(t, "CME", future.Exchange)
	require.Equal(t, "/ESU3", future.RollTargetSymbol)
	require.Equal(t, "XCME", future.StreamerExchangeCode)
	require.Equal(t, "/ESM23:XCME", future.StreamerSymbol)
	require.True(t, future.BackMonthFirstCalendarSymbol)
	require.True(t, future.IsTradeable)

	futureETF := future.FutureETFEquivalent

	require.Equal(t, "SPY", futureETF.Symbol)
	require.Equal(t, 501, futureETF.ShareQuantity)

	futureProd := future.FutureProduct

	require.Equal(t, "/ES", futureProd.RootSymbol)
	require.Equal(t, "ES", futureProd.Code)
	require.Equal(t, "E Mini S&P", futureProd.Description)
	require.Equal(t, "ES", futureProd.ClearingCode)
	require.Equal(t, "16", futureProd.ClearingExchangeCode)
	require.Equal(t, "ES", futureProd.ClearportCode)
	require.Equal(t, "ES", futureProd.LegacyCode)
	require.Equal(t, "CME", futureProd.Exchange)
	require.Equal(t, "CME", futureProd.LegacyExchangeCode)
	require.Equal(t, "Financial", futureProd.ProductType)

	require.Equal(t, 4, len(futureProd.ListedMonths))
	require.Equal(t, 4, len(futureProd.ActiveMonths))
	require.Equal(t, models.StringToFloat32(50), futureProd.NotionalMultiplier)
	require.Equal(t, models.StringToFloat32(.25), futureProd.TickSize)
	require.Equal(t, models.StringToFloat32(.01), futureProd.DisplayFactor)
	require.Equal(t, "XCME", futureProd.StreamerExchangeCode)
	require.False(t, futureProd.SmallNotional)
	require.True(t, futureProd.BackMonthFirstCalendarSymbol)
	require.False(t, futureProd.FirstNotice)
	require.True(t, futureProd.CashSettled)
	require.Equal(t, "ES", futureProd.SecurityGroup)
	require.Equal(t, "Equity Index", futureProd.MarketSector)

	roll := futureProd.Roll

	require.Equal(t, "equity_index", roll.Name)
	require.Equal(t, 3, roll.ActiveCount)
	require.True(t, roll.CashSettled)
	require.Equal(t, 4, roll.BusinessDaysOffset)
	require.False(t, roll.FirstNotice)

	require.Equal(t, models.StringToFloat32(.25), future.TickSizes[0].Value)
	require.Equal(t, models.StringToFloat32(.05), future.OptionTickSizes[0].Value)
	require.Equal(t, models.StringToFloat32(5), future.OptionTickSizes[0].Threshold)
	require.Equal(t, models.StringToFloat32(.05), future.SpreadTickSizes[0].Value)
	require.Equal(t, "/ESU3", future.SpreadTickSizes[0].Symbol)
}

func TestGetFuture(t *testing.T) {
	setup()
	defer teardown()

	symbol := "ESM3"

	mux.HandleFunc(fmt.Sprintf("/instruments/futures/%s", symbol), func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, futureResp)
	})

	productCode := "ES"

	future, err := client.GetFuture(symbol)
	require.Nil(t, err)

	require.Equal(t, "/ESM3", future.Symbol)
	require.Equal(t, productCode, future.ProductCode)
	require.Equal(t, models.StringToFloat32(50), future.ContractSize)
	require.Equal(t, models.StringToFloat32(.25), future.TickSize)
	require.Equal(t, models.StringToFloat32(50), future.NotionalMultiplier)
	require.Equal(t, models.StringToFloat32(0), future.MainFraction)
	require.Equal(t, models.StringToFloat32(0), future.SubFraction)
	require.Equal(t, models.StringToFloat32(0.01), future.DisplayFactor)
	require.Equal(t, "2023-06-16", future.LastTradeDate)
	require.Equal(t, "2023-06-16", future.ExpirationDate)
	require.Equal(t, "2023-06-16", future.ClosingOnlyDate)
	require.True(t, future.Active)
	require.True(t, future.ActiveMonth)
	require.False(t, future.NextActiveMonth)
	require.False(t, future.IsClosingOnly)
	require.Equal(t, "2023-06-16T13:30:00Z", future.StopsTradingAt.Format(time.RFC3339Nano))
	require.Equal(t, "2023-06-16T13:30:00Z", future.ExpiresAt.Format(time.RFC3339Nano))
	require.Equal(t, "CME_ES", future.ProductGroup)
	require.Equal(t, "CME", future.Exchange)
	require.Equal(t, "/ESU3", future.RollTargetSymbol)
	require.Equal(t, "XCME", future.StreamerExchangeCode)
	require.Equal(t, "/ESM23:XCME", future.StreamerSymbol)
	require.True(t, future.BackMonthFirstCalendarSymbol)
	require.True(t, future.IsTradeable)

	futureETF := future.FutureETFEquivalent

	require.Equal(t, "SPY", futureETF.Symbol)
	require.Equal(t, 501, futureETF.ShareQuantity)

	futureProd := future.FutureProduct

	require.Equal(t, "/ES", futureProd.RootSymbol)
	require.Equal(t, "ES", futureProd.Code)
	require.Equal(t, "E Mini S&P", futureProd.Description)
	require.Equal(t, "ES", futureProd.ClearingCode)
	require.Equal(t, "16", futureProd.ClearingExchangeCode)
	require.Equal(t, "ES", futureProd.ClearportCode)
	require.Equal(t, "ES", futureProd.LegacyCode)
	require.Equal(t, "CME", futureProd.Exchange)
	require.Equal(t, "CME", futureProd.LegacyExchangeCode)
	require.Equal(t, "Financial", futureProd.ProductType)

	require.Equal(t, 4, len(futureProd.ListedMonths))
	require.Equal(t, 4, len(futureProd.ActiveMonths))
	require.Equal(t, models.StringToFloat32(50), futureProd.NotionalMultiplier)
	require.Equal(t, models.StringToFloat32(.25), futureProd.TickSize)
	require.Equal(t, models.StringToFloat32(.01), futureProd.DisplayFactor)
	require.Equal(t, "XCME", futureProd.StreamerExchangeCode)
	require.False(t, futureProd.SmallNotional)
	require.True(t, futureProd.BackMonthFirstCalendarSymbol)
	require.False(t, futureProd.FirstNotice)
	require.True(t, futureProd.CashSettled)
	require.Equal(t, "ES", futureProd.SecurityGroup)
	require.Equal(t, "Equity Index", futureProd.MarketSector)

	roll := futureProd.Roll

	require.Equal(t, "equity_index", roll.Name)
	require.Equal(t, 3, roll.ActiveCount)
	require.True(t, roll.CashSettled)
	require.Equal(t, 4, roll.BusinessDaysOffset)
	require.False(t, roll.FirstNotice)

	require.Equal(t, models.StringToFloat32(.25), future.TickSizes[0].Value)
	require.Equal(t, models.StringToFloat32(.05), future.OptionTickSizes[0].Value)
	require.Equal(t, models.StringToFloat32(5), future.OptionTickSizes[0].Threshold)
	require.Equal(t, models.StringToFloat32(.05), future.SpreadTickSizes[0].Value)
	require.Equal(t, "/ESU3", future.SpreadTickSizes[0].Symbol)
}

const cryptocurrenciesResp = `{
  "data": {
    "items": [
      {
        "id": 1,
        "symbol": "BTC/USD",
        "instrument-type": "Cryptocurrency",
        "short-description": "Bitcoin",
        "description": "Bitcoin to USD",
        "is-closing-only": false,
        "active": true,
        "tick-size": "0.01",
        "streamer-symbol": "BTC/USD:CXTALP",
        "destination-venue-symbols": [
          {
            "id": 71,
            "symbol": "BTC",
            "destination-venue": "CITADEL_CRYPTOCURRENCY",
            "max-quantity-precision": 8,
            "max-price-precision": 8,
            "routable": true
          },
          {
            "id": 80,
            "symbol": "BTC/USD",
            "destination-venue": "CBOE_DIGITAL_CRYPTOCURRENCY",
            "max-quantity-precision": 6,
            "max-price-precision": 1,
            "routable": true
          }
        ]
      },
      {
        "id": 3,
        "symbol": "ETH/USD",
        "instrument-type": "Cryptocurrency",
        "short-description": "Ethereum",
        "description": "Ethereum to USD",
        "is-closing-only": false,
        "active": true,
        "tick-size": "0.01",
        "streamer-symbol": "ETH/USD:CXTALP"
      }
    ]
  },
  "context": "/instruments/cryptocurrencies"
}`

const activeEquitiesResp = `{
    "data": {
        "items": [
            {
                "id": 15625396,
                "symbol": "YMM",
                "instrument-type": "Equity",
                "cusip": "35969L108",
                "short-description": "Full Truck Alliance Co. Ltd. American Depositary Shares (each representing 20 Class A Ordinary Shares)",
                "is-index": false,
                "listed-market": "XNYS",
                "description": "Full Truck Alliance Co. Ltd. American Depositary Shares (each representing 20 Class A Ordinary Shares)",
                "lendability": "Easy To Borrow",
                "borrow-rate": "0.0",
                "market-time-instrument-collection": "Equity",
                "is-closing-only": false,
                "is-options-closing-only": false,
                "active": true,
                "is-fractional-quantity-eligible": false,
                "is-illiquid": false,
                "is-etf": false,
                "streamer-symbol": "YMM",
                "tick-sizes": [
                    {
                        "value": "0.0001",
                        "threshold": "1.0"
                    },
                    {
                        "value": "0.01"
                    }
                ],
                "option-tick-sizes": [
                    {
                        "value": "0.05",
                        "threshold": "3.0"
                    },
                    {
                        "value": "0.1"
                    }
                ]
            },
            {
                "id": 15578315,
                "symbol": "AMAM",
                "instrument-type": "Equity",
                "cusip": "02290A102",
                "short-description": "Ambrx Biopharma Inc. American Depositary Shares (each representing seven Ordinary Shares)",
                "is-index": false,
                "listed-market": "XNAS",
                "description": "Ambrx Biopharma Inc. American Depositary Shares (each representing seven Ordinary Shares)",
                "lendability": "Locate Required",
                "borrow-rate": "24.5013",
                "market-time-instrument-collection": "Equity",
                "is-closing-only": false,
                "is-options-closing-only": false,
                "active": true,
                "is-fractional-quantity-eligible": false,
                "is-illiquid": false,
                "is-etf": false,
                "streamer-symbol": "AMAM",
                "tick-sizes": [
                    {
                        "value": "0.0001",
                        "threshold": "1.0"
                    },
                    {
                        "value": "0.01"
                    }
                ],
                "option-tick-sizes": [
                    {
                        "value": "0.05",
                        "threshold": "3.0"
                    },
                    {
                        "value": "0.1"
                    }
                ]
            },
            {
                "id": 696,
                "symbol": "AFGE",
                "instrument-type": "Equity",
                "cusip": "025932864",
                "short-description": "AMERICAN FINANC",
                "is-index": false,
                "listed-market": "XNYS",
                "description": "AMERICAN FINANCIAL GROUP INC SUBORDINATED DEBENTURES DUE 2054",
                "lendability": "Locate Required",
                "borrow-rate": "10.6765",
                "market-time-instrument-collection": "Equity",
                "is-closing-only": false,
                "is-options-closing-only": false,
                "active": true,
                "is-fractional-quantity-eligible": false,
                "is-illiquid": false,
                "is-etf": false,
                "streamer-symbol": "AFGE",
                "tick-sizes": [
                    {
                        "value": "0.0001",
                        "threshold": "1.0"
                    },
                    {
                        "value": "0.01"
                    }
                ]
            },
            {
                "id": 24247603,
                "symbol": "RMIF",
                "instrument-type": "Equity",
                "short-description": "ETF Series Solutions LHA Risk-Managed Income ETF",
                "is-index": false,
                "listed-market": "BATS",
                "description": "ETF Series Solutions LHA Risk-Managed Income ETF",
                "lendability": "Preborrow",
                "market-time-instrument-collection": "Equity",
                "is-closing-only": false,
                "is-options-closing-only": false,
                "active": true,
                "is-fractional-quantity-eligible": false,
                "is-illiquid": false,
                "is-etf": false,
                "streamer-symbol": "RMIF",
                "tick-sizes": [
                    {
                        "value": "0.0001",
                        "threshold": "1.0"
                    },
                    {
                        "value": "0.01"
                    }
                ]
            }
        ]
    },
    "context": "/instruments/equities/active",
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

const equitiesResp = `{
    "data": {
        "items": [
            {
                "id": 726,
                "symbol": "AAPL",
                "instrument-type": "Equity",
                "cusip": "037833100",
                "short-description": "APPLE INC",
                "is-index": false,
                "listed-market": "XNAS",
                "description": "APPLE INC",
                "lendability": "Easy To Borrow",
                "borrow-rate": "0.0",
                "market-time-instrument-collection": "Equity",
                "is-closing-only": false,
                "is-options-closing-only": false,
                "active": true,
                "is-fractional-quantity-eligible": true,
                "is-illiquid": false,
                "is-etf": false,
                "streamer-symbol": "AAPL",
                "tick-sizes": [
                    {
                        "value": "0.0001",
                        "threshold": "1.0"
                    },
                    {
                        "value": "0.01"
                    }
                ],
                "option-tick-sizes": [
                    {
                        "value": "0.01",
                        "threshold": "3.0"
                    },
                    {
                        "value": "0.05"
                    }
                ]
            },
            {
                "id": 13754,
                "symbol": "TSLA",
                "instrument-type": "Equity",
                "cusip": "88160R101",
                "short-description": "TESLA MOTORS IN",
                "is-index": false,
                "listed-market": "XNAS",
                "description": "TESLA MOTORS INC",
                "lendability": "Easy To Borrow",
                "borrow-rate": "0.0",
                "market-time-instrument-collection": "Equity",
                "is-closing-only": false,
                "is-options-closing-only": false,
                "active": true,
                "is-fractional-quantity-eligible": true,
                "is-illiquid": false,
                "is-etf": false,
                "streamer-symbol": "TSLA",
                "tick-sizes": [
                    {
                        "value": "0.0001",
                        "threshold": "1.0"
                    },
                    {
                        "value": "0.01"
                    }
                ],
                "option-tick-sizes": [
                    {
                        "value": "0.01",
                        "threshold": "3.0"
                    },
                    {
                        "value": "0.05"
                    }
                ]
            }
        ]
    },
    "context": "/instruments/equities"
}`

const equityResp = `{
    "data": {
        "id": 726,
        "symbol": "AAPL",
        "instrument-type": "Equity",
        "cusip": "037833100",
        "short-description": "APPLE INC",
        "is-index": false,
        "listed-market": "XNAS",
        "description": "APPLE INC",
        "lendability": "Easy To Borrow",
        "borrow-rate": "0.0",
        "market-time-instrument-collection": "Equity",
        "is-closing-only": false,
        "is-options-closing-only": false,
        "active": true,
        "is-fractional-quantity-eligible": true,
        "is-illiquid": false,
        "is-etf": false,
        "streamer-symbol": "AAPL",
        "tick-sizes": [
            {
                "value": "0.0001",
                "threshold": "1.0"
            },
            {
                "value": "0.01"
            }
        ],
        "option-tick-sizes": [
            {
                "value": "0.01",
                "threshold": "3.0"
            },
            {
                "value": "0.05"
            }
        ]
    },
    "context": "/instruments/equities/AAPL"
}`

const equityOptionsResp = `{
    "data": {
        "items": [
            {
                "symbol": "AAPL  230616C00185000",
                "instrument-type": "Equity Option",
                "active": true,
                "strike-price": "185.0",
                "root-symbol": "AAPL",
                "underlying-symbol": "AAPL",
                "expiration-date": "2023-06-16",
                "exercise-style": "American",
                "shares-per-contract": 100,
                "option-type": "C",
                "option-chain-type": "Standard",
                "expiration-type": "Regular",
                "settlement-type": "PM",
                "stops-trading-at": "2023-06-16T20:00:00.000+00:00",
                "market-time-instrument-collection": "Equity Option",
                "days-to-expiration": 6,
                "expires-at": "2023-06-16T20:00:00.000+00:00",
                "is-closing-only": false,
                "streamer-symbol": ".AAPL230616C185"
            }
        ]
    },
    "context": "/instruments/equity-options"
}`

const equityOptionResp = `{
    "data": {
        "symbol": "AAPL  230616C00185000",
        "instrument-type": "Equity Option",
        "active": true,
        "strike-price": "185.0",
        "root-symbol": "AAPL",
        "underlying-symbol": "AAPL",
        "expiration-date": "2023-06-16",
        "exercise-style": "American",
        "shares-per-contract": 100,
        "option-type": "C",
        "option-chain-type": "Standard",
        "expiration-type": "Regular",
        "settlement-type": "PM",
        "stops-trading-at": "2023-06-16T20:00:00.000+00:00",
        "market-time-instrument-collection": "Equity Option",
        "days-to-expiration": 6,
        "expires-at": "2023-06-16T20:00:00.000+00:00",
        "is-closing-only": false,
        "streamer-symbol": ".AAPL230616C185"
    },
    "context": "/instruments/equity-options/AAPL%20%20230616C00185000"
}`

const futuresResp = `{
  "data": {
    "items": [
      {
        "symbol": "/ESM3",
        "product-code": "ES",
        "contract-size": "50.0",
        "tick-size": "0.25",
        "notional-multiplier": "50.0",
        "main-fraction": "0.0",
        "sub-fraction": "0.0",
        "display-factor": "0.01",
        "last-trade-date": "2023-06-16",
        "expiration-date": "2023-06-16",
        "closing-only-date": "2023-06-16",
        "active": true,
        "active-month": true,
        "next-active-month": false,
        "is-closing-only": false,
        "stops-trading-at": "2023-06-16T13:30:00.000+00:00",
        "expires-at": "2023-06-16T13:30:00.000+00:00",
        "product-group": "CME_ES",
        "exchange": "CME",
        "roll-target-symbol": "/ESU3",
        "streamer-exchange-code": "XCME",
        "streamer-symbol": "/ESM23:XCME",
        "back-month-first-calendar-symbol": true,
        "is-tradeable": true,
        "future-etf-equivalent": {
          "symbol": "SPY",
          "share-quantity": 501
        },
        "future-product": {
          "root-symbol": "/ES",
          "code": "ES",
          "description": "E Mini S&P",
          "clearing-code": "ES",
          "clearing-exchange-code": "16",
          "clearport-code": "ES",
          "legacy-code": "ES",
          "exchange": "CME",
          "legacy-exchange-code": "CME",
          "product-type": "Financial",
          "listed-months": ["H", "M", "U", "Z"],
          "active-months": ["H", "M", "U", "Z"],
          "notional-multiplier": "50.0",
          "tick-size": "0.25",
          "display-factor": "0.01",
          "streamer-exchange-code": "XCME",
          "small-notional": false,
          "back-month-first-calendar-symbol": true,
          "first-notice": false,
          "cash-settled": true,
          "security-group": "ES",
          "market-sector": "Equity Index",
          "roll": {
            "name": "equity_index",
            "active-count": 3,
            "cash-settled": true,
            "business-days-offset": 4,
            "first-notice": false
          }
        },
        "tick-sizes": [
          {
            "value": "0.25"
          }
        ],
        "option-tick-sizes": [
          {
            "value": "0.05",
            "threshold": "5.0"
          },
          {
            "value": "0.25"
          }
        ],
        "spread-tick-sizes": [
          {
            "value": "0.05",
            "symbol": "/ESU3"
          },
          {
            "value": "0.05",
            "symbol": "/ESZ3"
          }
        ]
      }
    ]
  },
  "context": "/instruments/futures"
}`

const futureResp = `{
    "data": {
        "symbol": "/ESM3",
        "product-code": "ES",
        "contract-size": "50.0",
        "tick-size": "0.25",
        "notional-multiplier": "50.0",
        "main-fraction": "0.0",
        "sub-fraction": "0.0",
        "display-factor": "0.01",
        "last-trade-date": "2023-06-16",
        "expiration-date": "2023-06-16",
        "closing-only-date": "2023-06-16",
        "active": true,
        "active-month": true,
        "next-active-month": false,
        "is-closing-only": false,
        "stops-trading-at": "2023-06-16T13:30:00.000+00:00",
        "expires-at": "2023-06-16T13:30:00.000+00:00",
        "product-group": "CME_ES",
        "exchange": "CME",
        "roll-target-symbol": "/ESU3",
        "streamer-exchange-code": "XCME",
        "streamer-symbol": "/ESM23:XCME",
        "back-month-first-calendar-symbol": true,
        "is-tradeable": true,
        "future-etf-equivalent": {
            "symbol": "SPY",
            "share-quantity": 501
        },
        "future-product": {
            "root-symbol": "/ES",
            "code": "ES",
            "description": "E Mini S&P",
            "clearing-code": "ES",
            "clearing-exchange-code": "16",
            "clearport-code": "ES",
            "legacy-code": "ES",
            "exchange": "CME",
            "legacy-exchange-code": "CME",
            "product-type": "Financial",
            "listed-months": [
                "H",
                "M",
                "U",
                "Z"
            ],
            "active-months": [
                "H",
                "M",
                "U",
                "Z"
            ],
            "notional-multiplier": "50.0",
            "tick-size": "0.25",
            "display-factor": "0.01",
            "streamer-exchange-code": "XCME",
            "small-notional": false,
            "back-month-first-calendar-symbol": true,
            "first-notice": false,
            "cash-settled": true,
            "security-group": "ES",
            "market-sector": "Equity Index",
            "roll": {
                "name": "equity_index",
                "active-count": 3,
                "cash-settled": true,
                "business-days-offset": 4,
                "first-notice": false
            }
        },
        "tick-sizes": [
            {
                "value": "0.25"
            }
        ],
        "option-tick-sizes": [
            {
                "value": "0.05",
                "threshold": "5.0"
            },
            {
                "value": "0.25"
            }
        ],
        "spread-tick-sizes": [
            {
                "value": "0.05",
                "symbol": "/ESU3"
            },
            {
                "value": "0.05",
                "symbol": "/ESZ3"
            }
        ]
    },
    "context": "/instruments/futures/ESM3"
}`
