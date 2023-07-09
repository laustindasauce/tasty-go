package tasty //nolint:testpackage // testing private field

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func TestGetActiveEquities(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/instruments/equities/active", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, activeEquitiesResp)
	})

	resp, pagination, err := client.GetActiveEquities(ActiveEquitiesQuery{PerPage: 4})
	require.Nil(t, err)

	require.Equal(t, 4, len(resp))

	equity := resp[0]

	require.Equal(t, 15625396, equity.ID)
	require.Equal(t, "YMM", equity.Symbol)
	require.Equal(t, EquityIT, equity.InstrumentType)
	require.Equal(t, "35969L108", equity.Cusip)
	require.Equal(t, "Full Truck Alliance Co. Ltd. American Depositary Shares (each representing 20 Class A Ordinary Shares)", equity.ShortDescription)
	require.False(t, equity.IsIndex)
	require.Equal(t, "XNYS", equity.ListedMarket)
	require.Equal(t, "Full Truck Alliance Co. Ltd. American Depositary Shares (each representing 20 Class A Ordinary Shares)", equity.Description)
	require.Equal(t, "Easy To Borrow", equity.Lendability)
	require.True(t, decimal.Zero.Equal(equity.BorrowRate))
	require.Equal(t, "Equity", equity.MarketTimeInstrumentCollection)
	require.False(t, equity.IsClosingOnly)
	require.False(t, equity.IsOptionsClosingOnly)
	require.True(t, equity.Active)
	require.False(t, equity.IsFractionalQuantityEligible)
	require.False(t, equity.IsIlliquid)
	require.False(t, equity.IsEtf)
	require.Equal(t, "YMM", equity.StreamerSymbol)

	tickSize := equity.TickSizes[0]

	require.True(t, tickSize.Value.Equal(decimal.NewFromFloat(0.0001)))
	require.True(t, decimal.NewFromInt(1).Equal(tickSize.Threshold))

	optionTickSize := equity.OptionTickSizes[0]

	require.True(t, optionTickSize.Value.Equal(decimal.NewFromFloat(0.05)))
	require.True(t, decimal.NewFromInt(3).Equal(optionTickSize.Threshold))

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

func TestGetActiveEquitiesError(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/instruments/equities/active", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(401)
		fmt.Fprint(writer, tastyUnauthorizedError)
	})

	_, _, err := client.GetActiveEquities(ActiveEquitiesQuery{PerPage: 4})
	expectedUnauthorized(t, err)
}
func TestGetEquities(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/instruments/equities", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, equitiesResp)
	})

	resp, err := client.GetEquities(EquitiesQuery{Symbols: []string{"AAPL", "TSLA"}})
	require.Nil(t, err)

	require.Equal(t, 2, len(resp))

	equity := resp[0]

	require.Equal(t, 726, equity.ID)
	require.Equal(t, "AAPL", equity.Symbol)
	require.Equal(t, EquityIT, equity.InstrumentType)
	require.Equal(t, "037833100", equity.Cusip)
	require.Equal(t, "APPLE INC", equity.ShortDescription)
	require.False(t, equity.IsIndex)
	require.Equal(t, "XNAS", equity.ListedMarket)
	require.Equal(t, "APPLE INC", equity.Description)
	require.Equal(t, "Easy To Borrow", equity.Lendability)
	require.True(t, decimal.Zero.Equal(equity.BorrowRate))
	require.Equal(t, "Equity", equity.MarketTimeInstrumentCollection)
	require.False(t, equity.IsClosingOnly)
	require.False(t, equity.IsOptionsClosingOnly)
	require.True(t, equity.Active)
	require.True(t, equity.IsFractionalQuantityEligible)
	require.False(t, equity.IsIlliquid)
	require.False(t, equity.IsEtf)
	require.Equal(t, "AAPL", equity.StreamerSymbol)

	tickSize := equity.TickSizes[0]

	require.True(t, tickSize.Value.Equal(decimal.NewFromFloat(0.0001)))
	require.True(t, tickSize.Threshold.Equal(decimal.NewFromInt(1)))

	optionTickSize := equity.OptionTickSizes[0]

	require.Equal(t, decimal.NewFromFloat(0.01), optionTickSize.Value)
	require.True(t, optionTickSize.Threshold.Equal(decimal.NewFromInt(3)))
}

func TestGetEquitiesError(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/instruments/equities", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(401)
		fmt.Fprint(writer, tastyUnauthorizedError)
	})

	_, err := client.GetEquities(EquitiesQuery{Symbols: []string{"AAPL", "TSLA"}})
	expectedUnauthorized(t, err)
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
	require.Equal(t, EquityIT, equity.InstrumentType)
	require.Equal(t, "037833100", equity.Cusip)
	require.Equal(t, "APPLE INC", equity.ShortDescription)
	require.False(t, equity.IsIndex)
	require.Equal(t, "XNAS", equity.ListedMarket)
	require.Equal(t, "APPLE INC", equity.Description)
	require.Equal(t, "Easy To Borrow", equity.Lendability)
	require.True(t, decimal.Zero.Equal(equity.BorrowRate))
	require.Equal(t, "Equity", equity.MarketTimeInstrumentCollection)
	require.False(t, equity.IsClosingOnly)
	require.False(t, equity.IsOptionsClosingOnly)
	require.True(t, equity.Active)
	require.True(t, equity.IsFractionalQuantityEligible)
	require.False(t, equity.IsIlliquid)
	require.False(t, equity.IsEtf)
	require.Equal(t, symbol, equity.StreamerSymbol)

	tickSize := equity.TickSizes[0]

	require.True(t, tickSize.Value.Equal(decimal.NewFromFloat(0.0001)))
	require.True(t, tickSize.Threshold.Equal(decimal.NewFromInt(1)))

	optionTickSize := equity.OptionTickSizes[0]

	require.Equal(t, decimal.NewFromFloat(0.01), optionTickSize.Value)
	require.True(t, optionTickSize.Threshold.Equal(decimal.NewFromInt(3)))
}

func TestGetEquityError(t *testing.T) {
	setup()
	defer teardown()

	symbol := "AAPL"

	mux.HandleFunc(fmt.Sprintf("/instruments/equities/%s", symbol), func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(401)
		fmt.Fprint(writer, tastyUnauthorizedError)
	})

	_, err := client.GetEquity(symbol)
	expectedUnauthorized(t, err)
}

func TestGetEquityOptions(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/instruments/equity-options", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, equityOptionsResp)
	})

	symbol := "AAPL"
	optionType := Call

	sym := EquityOptionsSymbology{
		Symbol:     symbol,
		Strike:     185,
		OptionType: optionType,
		Expiration: time.Date(2023, 6, 16, 0, 0, 0, 0, time.UTC),
	}
	occSymbol := sym.Build()

	resp, err := client.GetEquityOptions(EquityOptionsQuery{Symbols: []string{occSymbol}})
	require.Nil(t, err)

	require.Equal(t, 1, len(resp))

	equity := resp[0]

	require.Equal(t, occSymbol, equity.Symbol)
	require.Equal(t, EquityOptionIT, equity.InstrumentType)
	require.True(t, equity.Active)
	require.True(t, equity.StrikePrice.Equal(decimal.NewFromInt(185)))
	require.Equal(t, symbol, equity.RootSymbol)
	require.Equal(t, symbol, equity.UnderlyingSymbol)
	require.Equal(t, "2023-06-16", equity.ExpirationDate)
	require.Equal(t, "American", equity.ExerciseStyle)
	require.Equal(t, 100, equity.SharesPerContract)
	require.Equal(t, optionType, equity.OptionType)
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

func TestGetEquityOptionsError(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/instruments/equity-options", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(401)
		fmt.Fprint(writer, tastyUnauthorizedError)
	})

	symbol := "AAPL"
	optionType := Call

	sym := EquityOptionsSymbology{
		Symbol:     symbol,
		Strike:     185,
		OptionType: optionType,
		Expiration: time.Date(2023, 6, 16, 0, 0, 0, 0, time.UTC),
	}
	occSymbol := sym.Build()

	_, err := client.GetEquityOptions(EquityOptionsQuery{Symbols: []string{occSymbol}})
	expectedUnauthorized(t, err)
}

func TestGetEquityOption(t *testing.T) {
	setup()
	defer teardown()

	symbol := "AAPL"
	optionType := Call

	sym := EquityOptionsSymbology{
		Symbol:     symbol,
		Strike:     185,
		OptionType: optionType,
		Expiration: time.Date(2023, 6, 16, 0, 0, 0, 0, time.UTC),
	}
	occSymbol := sym.Build()

	mux.HandleFunc(fmt.Sprintf("/instruments/equity-options/%s", occSymbol), func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, equityOptionResp)
	})

	equity, err := client.GetEquityOption(sym, true)
	require.Nil(t, err)

	require.Equal(t, occSymbol, equity.Symbol)
	require.Equal(t, EquityOptionIT, equity.InstrumentType)
	require.True(t, equity.Active)
	require.True(t, equity.StrikePrice.Equal(decimal.NewFromInt(185)))
	require.Equal(t, symbol, equity.RootSymbol)
	require.Equal(t, symbol, equity.UnderlyingSymbol)
	require.Equal(t, "2023-06-16", equity.ExpirationDate)
	require.Equal(t, "American", equity.ExerciseStyle)
	require.Equal(t, 100, equity.SharesPerContract)
	require.Equal(t, optionType, equity.OptionType)
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

func TestGetEquityOptionError(t *testing.T) {
	setup()
	defer teardown()

	symbol := "AAPL"
	optionType := Call

	sym := EquityOptionsSymbology{
		Symbol:     symbol,
		Strike:     185,
		OptionType: optionType,
		Expiration: time.Date(2023, 6, 16, 0, 0, 0, 0, time.UTC),
	}
	occSymbol := sym.Build()

	mux.HandleFunc(fmt.Sprintf("/instruments/equity-options/%s", occSymbol), func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(401)
		fmt.Fprint(writer, tastyUnauthorizedError)
	})

	_, err := client.GetEquityOption(sym, true)
	expectedUnauthorized(t, err)
}

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
