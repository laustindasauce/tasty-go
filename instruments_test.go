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

	sym := utils.EquityOptionsSymbology{
		Symbol:     symbol,
		Strike:     185,
		OptionType: optionType,
		Expiration: time.Date(2023, 6, 16, 0, 0, 0, 0, time.UTC),
	}
	occSymbol := sym.Build()

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

	sym := utils.EquityOptionsSymbology{
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

func TestGetFutureOptionProducts(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/instruments/future-option-products", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, futureOptionsProducts)
	})

	resp, err := client.GetFutureOptionProducts()
	require.Nil(t, err)

	require.Equal(t, 2, len(resp))

	fop := resp[0]

	require.Equal(t, "ECYM", fop.RootSymbol)
	require.True(t, fop.CashSettled)
	require.Equal(t, "ECYM", fop.Code)
	require.Equal(t, "ECYM", fop.LegacyCode)
	require.Equal(t, "ECYM", fop.ClearportCode)
	require.Equal(t, "DB", fop.ClearingCode)
	require.Equal(t, "01", fop.ClearingExchangeCode)
	require.Equal(t, models.StringToFloat32(100), fop.ClearingPriceMultiplier)
	require.Equal(t, models.StringToFloat32(0.01), fop.DisplayFactor)
	require.Equal(t, "CME", fop.Exchange)
	require.Equal(t, "Financial", fop.ProductType)
	require.Equal(t, "Regular", fop.ExpirationType)
	require.Equal(t, 0, fop.SettlementDelayDays)
	require.False(t, fop.IsRollover)
	require.Equal(t, "Event", fop.ProductSubtype)
	require.Equal(t, "Equity Index", fop.MarketSector)

	futureProd := fop.FutureProduct

	require.Equal(t, "/ECYM", futureProd.RootSymbol)
	require.Equal(t, "ECYM", futureProd.Code)
	require.Equal(t, "/YM Event Synthetic Underlying", futureProd.Description)
	require.Equal(t, "DB", futureProd.ClearingCode)
	require.Equal(t, "01", futureProd.ClearingExchangeCode)
	require.Equal(t, "ECYM", futureProd.ClearportCode)
	require.Equal(t, "ECY", futureProd.LegacyCode)
	require.Equal(t, "CME", futureProd.Exchange)
	require.Equal(t, "CBT", futureProd.LegacyExchangeCode)
	require.Equal(t, "Financial", futureProd.ProductType)
	require.Equal(t, 4, len(futureProd.ListedMonths))
	require.Equal(t, 4, len(futureProd.ActiveMonths))
	require.Equal(t, models.StringToFloat32(1.0), futureProd.NotionalMultiplier)
	require.Equal(t, models.StringToFloat32(.25), futureProd.TickSize)
	require.Equal(t, models.StringToFloat32(1.0), futureProd.DisplayFactor)
	require.Equal(t, "XCBT", futureProd.StreamerExchangeCode)
	require.True(t, futureProd.SmallNotional)
	require.True(t, futureProd.BackMonthFirstCalendarSymbol)
	require.False(t, futureProd.FirstNotice)
	require.True(t, futureProd.CashSettled)
	require.Equal(t, "Event", futureProd.ProductSubtype)
	require.Equal(t, "YM", futureProd.TrueUnderlyingCode)
	require.Equal(t, "Equity Index", futureProd.MarketSector)

	roll := futureProd.Roll

	require.Equal(t, "synthetics", roll.Name)
	require.Equal(t, 2, roll.ActiveCount)
	require.True(t, roll.CashSettled)
	require.Equal(t, 0, roll.BusinessDaysOffset)
	require.False(t, roll.FirstNotice)
}

func TestGetFutureOptionProduct(t *testing.T) {
	setup()
	defer teardown()

	rootSymbol := "ECYM"
	exchange := "CME"

	mux.HandleFunc(fmt.Sprintf("/instruments/future-option-products/%s/%s", exchange, rootSymbol), func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, futureOptionProduct)
	})

	fop, err := client.GetFutureOptionProduct(exchange, rootSymbol)
	require.Nil(t, err)

	require.Equal(t, "ECYM", fop.RootSymbol)
	require.True(t, fop.CashSettled)
	require.Equal(t, "ECYM", fop.Code)
	require.Equal(t, "ECYM", fop.LegacyCode)
	require.Equal(t, "ECYM", fop.ClearportCode)
	require.Equal(t, "DB", fop.ClearingCode)
	require.Equal(t, "01", fop.ClearingExchangeCode)
	require.Equal(t, models.StringToFloat32(100), fop.ClearingPriceMultiplier)
	require.Equal(t, models.StringToFloat32(0.01), fop.DisplayFactor)
	require.Equal(t, "CME", fop.Exchange)
	require.Equal(t, "Financial", fop.ProductType)
	require.Equal(t, "Regular", fop.ExpirationType)
	require.Equal(t, 0, fop.SettlementDelayDays)
	require.False(t, fop.IsRollover)
	require.Equal(t, "Event", fop.ProductSubtype)
	require.Equal(t, "Equity Index", fop.MarketSector)

	futureProd := fop.FutureProduct

	require.Equal(t, "/ECYM", futureProd.RootSymbol)
	require.Equal(t, "ECYM", futureProd.Code)
	require.Equal(t, "/YM Event Synthetic Underlying", futureProd.Description)
	require.Equal(t, "DB", futureProd.ClearingCode)
	require.Equal(t, "01", futureProd.ClearingExchangeCode)
	require.Equal(t, "ECYM", futureProd.ClearportCode)
	require.Equal(t, "ECY", futureProd.LegacyCode)
	require.Equal(t, "CME", futureProd.Exchange)
	require.Equal(t, "CBT", futureProd.LegacyExchangeCode)
	require.Equal(t, "Financial", futureProd.ProductType)
	require.Equal(t, 4, len(futureProd.ListedMonths))
	require.Equal(t, 4, len(futureProd.ActiveMonths))
	require.Equal(t, models.StringToFloat32(1.0), futureProd.NotionalMultiplier)
	require.Equal(t, models.StringToFloat32(.25), futureProd.TickSize)
	require.Equal(t, models.StringToFloat32(1.0), futureProd.DisplayFactor)
	require.Equal(t, "XCBT", futureProd.StreamerExchangeCode)
	require.True(t, futureProd.SmallNotional)
	require.True(t, futureProd.BackMonthFirstCalendarSymbol)
	require.False(t, futureProd.FirstNotice)
	require.True(t, futureProd.CashSettled)
	require.Equal(t, "Event", futureProd.ProductSubtype)
	require.Equal(t, "YM", futureProd.TrueUnderlyingCode)
	require.Equal(t, "Equity Index", futureProd.MarketSector)

	roll := futureProd.Roll

	require.Equal(t, "synthetics", roll.Name)
	require.Equal(t, 2, roll.ActiveCount)
	require.True(t, roll.CashSettled)
	require.Equal(t, 0, roll.BusinessDaysOffset)
	require.False(t, roll.FirstNotice)
}

func TestGetFutureOptions(t *testing.T) {
	setup()
	defer teardown()

	future := utils.FutureSymbology{ProductCode: "ES", MonthCode: constants.December, YearDigit: 9}

	expiry := time.Date(2019, 9, 27, 0, 0, 0, 0, time.Local)
	fcc := utils.FutureOptionsSymbology{
		OptionContractCode: "EW4U9",
		FutureContractCode: future.Build(),
		OptionType:         constants.Put,
		Strike:             2975,
		Expiration:         expiry,
	}

	symbol := fcc.Build()

	query := models.FutureOptionsQuery{
		Symbols: []string{symbol},
	}

	mux.HandleFunc("/instruments/future-options", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, futureOptionsResp)
	})

	resp, err := client.GetFutureOptions(query)
	require.Nil(t, err)

	require.Equal(t, 1, len(resp))

	fo := resp[0]

	require.Equal(t, symbol, fo.Symbol)
	require.Equal(t, future.Build(), fo.UnderlyingSymbol)
	require.Equal(t, "ES", fo.ProductCode)
	require.Equal(t, expiry.Format("2006-01-02"), fo.ExpirationDate)
	require.Equal(t, "/ES", fo.RootSymbol)
	require.Equal(t, "EW4", fo.OptionRootSymbol)
	require.Equal(t, models.StringToFloat32(2975), fo.StrikePrice)
	require.Equal(t, "CME", fo.Exchange)
	require.Equal(t, "EW4U9 P2975", fo.ExchangeSymbol)
	require.Equal(t, string(constants.Put), fo.OptionType)
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
	require.Equal(t, expiry.Format("2006-01-02"), fo.MaturityDate)
	require.True(t, fo.IsExercisableWeekly)
	require.Equal(t, "0", fo.LastTradeTime)
	require.Equal(t, -1, fo.DaysToExpiration)
	require.False(t, fo.IsClosingOnly)
	require.False(t, fo.Active)
	require.Equal(t, "2019-09-27T20:00:00Z", fo.StopsTradingAt.Format(time.RFC3339))
	require.Equal(t, "2019-09-27T20:00:00Z", fo.ExpiresAt.Format(time.RFC3339))

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

func TestGetFutureOption(t *testing.T) {
	setup()
	defer teardown()

	future := utils.FutureSymbology{ProductCode: "ES", MonthCode: constants.December, YearDigit: 9}

	expiry := time.Date(2019, 9, 27, 0, 0, 0, 0, time.Local)
	fcc := utils.FutureOptionsSymbology{
		OptionContractCode: "EW4U9",
		FutureContractCode: future.Build(),
		OptionType:         constants.Put,
		Strike:             2975,
		Expiration:         expiry,
	}

	symbol := fcc.Build()

	// Need to hard code the endpoint since mux won't handle the . in the url
	mux.HandleFunc("/instruments/future-options/test", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, futureOptionResp)
	})

	fo, err := client.GetFutureOption("test")
	require.Nil(t, err)

	require.Equal(t, symbol, fo.Symbol)
	require.Equal(t, future.Build(), fo.UnderlyingSymbol)
	require.Equal(t, "ES", fo.ProductCode)
	require.Equal(t, expiry.Format("2006-01-02"), fo.ExpirationDate)
	require.Equal(t, "/ES", fo.RootSymbol)
	require.Equal(t, "EW4", fo.OptionRootSymbol)
	require.Equal(t, models.StringToFloat32(2975), fo.StrikePrice)
	require.Equal(t, "CME", fo.Exchange)
	require.Equal(t, "EW4U9 P2975", fo.ExchangeSymbol)
	require.Equal(t, string(constants.Put), fo.OptionType)
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
	require.Equal(t, expiry.Format("2006-01-02"), fo.MaturityDate)
	require.True(t, fo.IsExercisableWeekly)
	require.Equal(t, "0", fo.LastTradeTime)
	require.Equal(t, -1, fo.DaysToExpiration)
	require.False(t, fo.IsClosingOnly)
	require.False(t, fo.Active)
	require.Equal(t, "2019-09-27T20:00:00Z", fo.StopsTradingAt.Format(time.RFC3339))
	require.Equal(t, "2019-09-27T20:00:00Z", fo.ExpiresAt.Format(time.RFC3339))

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

func TestGetFutureProduct(t *testing.T) {
	setup()
	defer teardown()

	exchange := constants.CME
	code := "ES"

	mux.HandleFunc(fmt.Sprintf("/instruments/future-products/%s/%s", exchange, code), func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, futureProductResp)
	})

	fp, err := client.GetFutureProduct(exchange, code)
	require.Nil(t, err)

	require.Equal(t, "/ES", fp.RootSymbol)
	require.Equal(t, "ES", fp.Code)
	require.Equal(t, "E Mini S&P", fp.Description)
	require.Equal(t, "ES", fp.ClearingCode)
	require.Equal(t, "16", fp.ClearingExchangeCode)
	require.Equal(t, "ES", fp.ClearportCode)
	require.Equal(t, "ES", fp.LegacyCode)
	require.Equal(t, "CME", fp.Exchange)
	require.Equal(t, "CME", fp.LegacyExchangeCode)
	require.Equal(t, "Financial", fp.ProductType)
	require.Equal(t, 4, len(fp.ListedMonths))
	require.Equal(t, 4, len(fp.ActiveMonths))
	require.Equal(t, models.StringToFloat32(50), fp.NotionalMultiplier)
	require.Equal(t, models.StringToFloat32(.25), fp.TickSize)
	require.Equal(t, models.StringToFloat32(.01), fp.DisplayFactor)
	require.Equal(t, "XCME", fp.StreamerExchangeCode)
	require.False(t, fp.SmallNotional)
	require.True(t, fp.BackMonthFirstCalendarSymbol)
	require.False(t, fp.FirstNotice)
	require.True(t, fp.CashSettled)
	require.Equal(t, "ES", fp.SecurityGroup)
	require.Equal(t, "Equity Index", fp.MarketSector)

	op := fp.OptionProducts[1]

	require.Equal(t, "E2B", op.RootSymbol)
	require.False(t, op.CashSettled)
	require.Equal(t, "E2B", op.Code)
	require.Equal(t, "E2B", op.LegacyCode)
	require.Equal(t, "E2B", op.ClearportCode)
	require.Equal(t, "ES", op.ClearingCode)
	require.Equal(t, "9C", op.ClearingExchangeCode)
	require.Equal(t, models.StringToFloat32(1), op.ClearingPriceMultiplier)
	require.Equal(t, models.StringToFloat32(.01), op.DisplayFactor)
	require.Equal(t, "CME", op.Exchange)
	require.Equal(t, "Physical", op.ProductType)
	require.Equal(t, "Weekly", op.ExpirationType)
	require.Equal(t, 0, op.SettlementDelayDays)
	require.False(t, op.IsRollover)
	require.Equal(t, "Equity Index", op.MarketSector)

	roll := fp.Roll

	require.Equal(t, "equity_index", roll.Name)
	require.Equal(t, 3, roll.ActiveCount)
	require.True(t, roll.CashSettled)
	require.Equal(t, 4, roll.BusinessDaysOffset)
	require.False(t, roll.FirstNotice)
}

func TestGetFutureProducts(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/instruments/future-products", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, futureProductsResp)
	})

	resp, err := client.GetFutureProducts()
	require.Nil(t, err)

	require.Equal(t, 1, len(resp))

	fp := resp[0]

	require.Equal(t, "/ES", fp.RootSymbol)
	require.Equal(t, "ES", fp.Code)
	require.Equal(t, "E Mini S&P", fp.Description)
	require.Equal(t, "ES", fp.ClearingCode)
	require.Equal(t, "16", fp.ClearingExchangeCode)
	require.Equal(t, "ES", fp.ClearportCode)
	require.Equal(t, "ES", fp.LegacyCode)
	require.Equal(t, "CME", fp.Exchange)
	require.Equal(t, "CME", fp.LegacyExchangeCode)
	require.Equal(t, "Financial", fp.ProductType)
	require.Equal(t, 4, len(fp.ListedMonths))
	require.Equal(t, 4, len(fp.ActiveMonths))
	require.Equal(t, models.StringToFloat32(50), fp.NotionalMultiplier)
	require.Equal(t, models.StringToFloat32(.25), fp.TickSize)
	require.Equal(t, models.StringToFloat32(.01), fp.DisplayFactor)
	require.Equal(t, "XCME", fp.StreamerExchangeCode)
	require.False(t, fp.SmallNotional)
	require.True(t, fp.BackMonthFirstCalendarSymbol)
	require.False(t, fp.FirstNotice)
	require.True(t, fp.CashSettled)
	require.Equal(t, "ES", fp.SecurityGroup)
	require.Equal(t, "Equity Index", fp.MarketSector)

	op := fp.OptionProducts[1]

	require.Equal(t, "E2B", op.RootSymbol)
	require.False(t, op.CashSettled)
	require.Equal(t, "E2B", op.Code)
	require.Equal(t, "E2B", op.LegacyCode)
	require.Equal(t, "E2B", op.ClearportCode)
	require.Equal(t, "ES", op.ClearingCode)
	require.Equal(t, "9C", op.ClearingExchangeCode)
	require.Equal(t, models.StringToFloat32(1), op.ClearingPriceMultiplier)
	require.Equal(t, models.StringToFloat32(.01), op.DisplayFactor)
	require.Equal(t, "CME", op.Exchange)
	require.Equal(t, "Physical", op.ProductType)
	require.Equal(t, "Weekly", op.ExpirationType)
	require.Equal(t, 0, op.SettlementDelayDays)
	require.False(t, op.IsRollover)
	require.Equal(t, "Equity Index", op.MarketSector)

	roll := fp.Roll

	require.Equal(t, "equity_index", roll.Name)
	require.Equal(t, 3, roll.ActiveCount)
	require.True(t, roll.CashSettled)
	require.Equal(t, 4, roll.BusinessDaysOffset)
	require.False(t, roll.FirstNotice)
}

func TestGetQuantityDecimalPrecisions(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/instruments/quantity-decimal-precisions", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, quantityDecimalPrecisionsResp)
	})

	resp, err := client.GetQuantityDecimalPrecisions()
	require.Nil(t, err)

	prec := resp[0]

	require.Equal(t, "Cryptocurrency", prec.InstrumentType)
	require.Equal(t, "AAVE/USD", prec.Symbol)
	require.Equal(t, 8, prec.Value)
	require.Equal(t, 6, prec.MinimumIncrementPrecision)
}

func TestGetWarrants(t *testing.T) {
	setup()
	defer teardown()

	symbol := "NKLAW"

	mux.HandleFunc("/instruments/warrants", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, warrantsResp)
	})

	resp, err := client.GetWarrants([]string{symbol})
	require.Nil(t, err)

	war := resp[0]

	require.Equal(t, symbol, war.Symbol)
	require.Equal(t, "Warrant", war.InstrumentType)
	require.Equal(t, "XNAS", war.ListedMarket)
	require.Equal(t, "Nikola Corporation - Warrant expiring 6/3/2025", war.Description)
	require.False(t, war.IsClosingOnly)
	require.False(t, war.Active)
}

func TestGetWarrant(t *testing.T) {
	setup()
	defer teardown()

	symbol := "NKLAW"

	mux.HandleFunc(fmt.Sprintf("/instruments/warrants/%s", symbol), func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, warrantResp)
	})

	war, err := client.GetWarrant(symbol)
	require.Nil(t, err)

	require.Equal(t, symbol, war.Symbol)
	require.Equal(t, "Warrant", war.InstrumentType)
	require.Equal(t, "XNAS", war.ListedMarket)
	require.Equal(t, "Nikola Corporation - Warrant expiring 6/3/2025", war.Description)
	require.False(t, war.IsClosingOnly)
	require.False(t, war.Active)
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

const futureOptionsProducts = `{
  "data": {
    "items": [
      {
        "root-symbol": "ECYM",
        "cash-settled": true,
        "code": "ECYM",
        "legacy-code": "ECYM",
        "clearport-code": "ECYM",
        "clearing-code": "DB",
        "clearing-exchange-code": "01",
        "clearing-price-multiplier": "100.0",
        "display-factor": "0.01",
        "exchange": "CME",
        "product-type": "Financial",
        "expiration-type": "Regular",
        "settlement-delay-days": 0,
        "is-rollover": false,
        "product-subtype": "Event",
        "market-sector": "Equity Index",
        "future-product": {
          "root-symbol": "/ECYM",
          "code": "ECYM",
          "description": "/YM Event Synthetic Underlying",
          "clearing-code": "DB",
          "clearing-exchange-code": "01",
          "clearport-code": "ECYM",
          "legacy-code": "ECY",
          "exchange": "CME",
          "legacy-exchange-code": "CBT",
          "product-type": "Financial",
          "listed-months": ["H", "M", "U", "Z"],
          "active-months": ["H", "M", "U", "Z"],
          "notional-multiplier": "1.0",
          "tick-size": "0.25",
          "display-factor": "1.0",
          "streamer-exchange-code": "XCBT",
          "small-notional": true,
          "back-month-first-calendar-symbol": true,
          "first-notice": false,
          "cash-settled": true,
          "product-subtype": "Event",
          "true-underlying-code": "YM",
          "market-sector": "Equity Index",
          "roll": {
            "name": "synthetics",
            "active-count": 2,
            "cash-settled": true,
            "business-days-offset": 0,
            "first-notice": false
          }
        }
      },
      {
        "root-symbol": "Q1B",
        "cash-settled": false,
        "code": "Q1B",
        "legacy-code": "Q1B",
        "clearport-code": "Q1B",
        "clearing-code": "FA",
        "clearing-exchange-code": "9C",
        "clearing-price-multiplier": "1.0",
        "display-factor": "0.01",
        "exchange": "CME",
        "product-type": "Physical",
        "expiration-type": "Weekly",
        "settlement-delay-days": 0,
        "is-rollover": false,
        "market-sector": "Equity Index",
        "future-product": {
          "root-symbol": "/NQ",
          "code": "NQ",
          "description": "E Mini Nasdaq",
          "clearing-code": "NQ",
          "clearing-exchange-code": "16",
          "clearport-code": "NQ",
          "legacy-code": "NQ",
          "exchange": "CME",
          "legacy-exchange-code": "CME",
          "product-type": "Financial",
          "listed-months": ["H", "M", "U", "Z"],
          "active-months": ["H", "M", "U", "Z"],
          "notional-multiplier": "20.0",
          "tick-size": "0.25",
          "display-factor": "0.01",
          "streamer-exchange-code": "XCME",
          "small-notional": false,
          "back-month-first-calendar-symbol": true,
          "first-notice": false,
          "cash-settled": true,
          "security-group": "NQ",
          "market-sector": "Equity Index",
          "roll": {
            "name": "equity_index",
            "active-count": 3,
            "cash-settled": true,
            "business-days-offset": 4,
            "first-notice": false
          }
        }
      }
    ]
  },
  "context": "/instruments/future-option-products"
}`

const futureOptionProduct = `{
  "data": {
    "root-symbol": "ECYM",
    "cash-settled": true,
    "code": "ECYM",
    "legacy-code": "ECYM",
    "clearport-code": "ECYM",
    "clearing-code": "DB",
    "clearing-exchange-code": "01",
    "clearing-price-multiplier": "100.0",
    "display-factor": "0.01",
    "exchange": "CME",
    "product-type": "Financial",
    "expiration-type": "Regular",
    "settlement-delay-days": 0,
    "is-rollover": false,
    "product-subtype": "Event",
    "market-sector": "Equity Index",
    "future-product": {
      "root-symbol": "/ECYM",
      "code": "ECYM",
      "description": "/YM Event Synthetic Underlying",
      "clearing-code": "DB",
      "clearing-exchange-code": "01",
      "clearport-code": "ECYM",
      "legacy-code": "ECY",
      "exchange": "CME",
      "legacy-exchange-code": "CBT",
      "product-type": "Financial",
      "listed-months": ["H", "M", "U", "Z"],
      "active-months": ["H", "M", "U", "Z"],
      "notional-multiplier": "1.0",
      "tick-size": "0.25",
      "display-factor": "1.0",
      "streamer-exchange-code": "XCBT",
      "small-notional": true,
      "back-month-first-calendar-symbol": true,
      "first-notice": false,
      "cash-settled": true,
      "product-subtype": "Event",
      "true-underlying-code": "YM",
      "market-sector": "Equity Index",
      "roll": {
        "name": "synthetics",
        "active-count": 2,
        "cash-settled": true,
        "business-days-offset": 0,
        "first-notice": false
      }
    }
  },
  "context": "/instruments/future-option-products/CME/ECYM"
}`

const futureOptionsResp = `{
    "data": {
        "items": [
            {
                "symbol": "./ESZ9 EW4U9 190927P2975",
                "underlying-symbol": "/ESZ9",
                "product-code": "ES",
                "expiration-date": "2019-09-27",
                "root-symbol": "/ES",
                "option-root-symbol": "EW4",
                "strike-price": "2975.0",
                "exchange": "CME",
                "exchange-symbol": "EW4U9 P2975",
                "option-type": "P",
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
                "maturity-date": "2019-09-27",
                "is-exercisable-weekly": true,
                "last-trade-time": "0",
                "days-to-expiration": -1,
                "is-closing-only": false,
                "active": false,
                "stops-trading-at": "2019-09-27T20:00:00.000+00:00",
                "expires-at": "2019-09-27T20:00:00.000+00:00",
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
    "context": "/instruments/future-options"
}`

const futureOptionResp = `{
    "data": {
        "symbol": "./ESZ9 EW4U9 190927P2975",
        "underlying-symbol": "/ESZ9",
        "product-code": "ES",
        "expiration-date": "2019-09-27",
        "root-symbol": "/ES",
        "option-root-symbol": "EW4",
        "strike-price": "2975.0",
        "exchange": "CME",
        "exchange-symbol": "EW4U9 P2975",
        "option-type": "P",
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
        "maturity-date": "2019-09-27",
        "is-exercisable-weekly": true,
        "last-trade-time": "0",
        "days-to-expiration": -1,
        "is-closing-only": false,
        "active": false,
        "stops-trading-at": "2019-09-27T20:00:00.000+00:00",
        "expires-at": "2019-09-27T20:00:00.000+00:00",
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
    },
    "context": "/instruments/future-options/./ESZ9 W4U9%20190927P2975"
}`

const futureProductsResp = `{
  "data": {
    "items": [
      {
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
        "option-products": [
          {
            "root-symbol": "E1B",
            "cash-settled": false,
            "code": "E1B",
            "legacy-code": "E1B",
            "clearport-code": "E1B",
            "clearing-code": "ER",
            "clearing-exchange-code": "9C",
            "clearing-price-multiplier": "1.0",
            "display-factor": "0.01",
            "exchange": "CME",
            "product-type": "Physical",
            "expiration-type": "Weekly",
            "settlement-delay-days": 0,
            "is-rollover": false,
            "market-sector": "Equity Index"
          },
          {
            "root-symbol": "E2B",
            "cash-settled": false,
            "code": "E2B",
            "legacy-code": "E2B",
            "clearport-code": "E2B",
            "clearing-code": "ES",
            "clearing-exchange-code": "9C",
            "clearing-price-multiplier": "1.0",
            "display-factor": "0.01",
            "exchange": "CME",
            "product-type": "Physical",
            "expiration-type": "Weekly",
            "settlement-delay-days": 0,
            "is-rollover": false,
            "market-sector": "Equity Index"
          }
        ],
        "roll": {
          "name": "equity_index",
          "active-count": 3,
          "cash-settled": true,
          "business-days-offset": 4,
          "first-notice": false
        }
      }
    ]
  },
  "context": "/instruments/future-products"
}`

const futureProductResp = `{
    "data": {
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
        "option-products": [
            {
                "root-symbol": "E1B",
                "cash-settled": false,
                "code": "E1B",
                "legacy-code": "E1B",
                "clearport-code": "E1B",
                "clearing-code": "ER",
                "clearing-exchange-code": "9C",
                "clearing-price-multiplier": "1.0",
                "display-factor": "0.01",
                "exchange": "CME",
                "product-type": "Physical",
                "expiration-type": "Weekly",
                "settlement-delay-days": 0,
                "is-rollover": false,
                "market-sector": "Equity Index"
            },
            {
                "root-symbol": "E2B",
                "cash-settled": false,
                "code": "E2B",
                "legacy-code": "E2B",
                "clearport-code": "E2B",
                "clearing-code": "ES",
                "clearing-exchange-code": "9C",
                "clearing-price-multiplier": "1.0",
                "display-factor": "0.01",
                "exchange": "CME",
                "product-type": "Physical",
                "expiration-type": "Weekly",
                "settlement-delay-days": 0,
                "is-rollover": false,
                "market-sector": "Equity Index"
            },
            {
                "root-symbol": "E3B",
                "cash-settled": false,
                "code": "E3B",
                "legacy-code": "E3B",
                "clearport-code": "E3B",
                "clearing-code": "EU",
                "clearing-exchange-code": "9C",
                "clearing-price-multiplier": "1.0",
                "display-factor": "0.01",
                "exchange": "CME",
                "product-type": "Physical",
                "expiration-type": "Weekly",
                "settlement-delay-days": 0,
                "is-rollover": false,
                "market-sector": "Equity Index"
            },
            {
                "root-symbol": "E4B",
                "cash-settled": false,
                "code": "E4B",
                "legacy-code": "E4B",
                "clearport-code": "E4B",
                "clearing-code": "EV",
                "clearing-exchange-code": "9C",
                "clearing-price-multiplier": "1.0",
                "display-factor": "0.01",
                "exchange": "CME",
                "product-type": "Physical",
                "expiration-type": "Weekly",
                "settlement-delay-days": 0,
                "is-rollover": true,
                "market-sector": "Equity Index"
            },
            {
                "root-symbol": "E5B",
                "cash-settled": false,
                "code": "E5B",
                "legacy-code": "E5B",
                "clearport-code": "E5B",
                "clearing-code": "EW",
                "clearing-exchange-code": "9C",
                "clearing-price-multiplier": "1.0",
                "display-factor": "0.01",
                "exchange": "CME",
                "product-type": "Physical",
                "expiration-type": "Weekly",
                "settlement-delay-days": 0,
                "is-rollover": true,
                "market-sector": "Equity Index"
            },
            {
                "root-symbol": "E1D",
                "cash-settled": false,
                "code": "E1D",
                "legacy-code": "E1D",
                "clearport-code": "E1D",
                "clearing-code": "EX",
                "clearing-exchange-code": "9C",
                "clearing-price-multiplier": "1.0",
                "display-factor": "0.01",
                "exchange": "CME",
                "product-type": "Physical",
                "expiration-type": "Weekly",
                "settlement-delay-days": 0,
                "is-rollover": false,
                "market-sector": "Equity Index"
            },
            {
                "root-symbol": "E3C",
                "cash-settled": false,
                "code": "E3C",
                "legacy-code": "E3C",
                "clearport-code": "E3C",
                "clearing-code": "Z3",
                "clearing-exchange-code": "9C",
                "clearing-price-multiplier": "1.0",
                "display-factor": "0.01",
                "exchange": "CME",
                "product-type": "Physical",
                "expiration-type": "Weekly",
                "settlement-delay-days": 0,
                "is-rollover": false,
                "market-sector": "Equity Index"
            },
            {
                "root-symbol": "ES",
                "cash-settled": false,
                "code": "ES",
                "legacy-code": "ES",
                "clearport-code": "ES",
                "clearing-code": "ES",
                "clearing-exchange-code": "16",
                "clearing-price-multiplier": "1.0",
                "display-factor": "0.01",
                "exchange": "CME",
                "product-type": "Physical",
                "expiration-type": "Regular",
                "settlement-delay-days": 0,
                "is-rollover": false,
                "market-sector": "Equity Index"
            },
            {
                "root-symbol": "EW1",
                "cash-settled": false,
                "code": "EW1",
                "legacy-code": "EW1",
                "clearport-code": "EW1",
                "clearing-code": "W1",
                "clearing-exchange-code": "9C",
                "clearing-price-multiplier": "1.0",
                "display-factor": "0.01",
                "exchange": "CME",
                "product-type": "Physical",
                "expiration-type": "Weekly",
                "settlement-delay-days": 0,
                "is-rollover": false,
                "market-sector": "Equity Index"
            },
            {
                "root-symbol": "EW2",
                "cash-settled": false,
                "code": "EW2",
                "legacy-code": "EW2",
                "clearport-code": "EW2",
                "clearing-code": "W2",
                "clearing-exchange-code": "9C",
                "clearing-price-multiplier": "1.0",
                "display-factor": "0.01",
                "exchange": "CME",
                "product-type": "Physical",
                "expiration-type": "Weekly",
                "settlement-delay-days": 0,
                "is-rollover": false,
                "market-sector": "Equity Index"
            },
            {
                "root-symbol": "EW3",
                "cash-settled": false,
                "code": "EW3",
                "legacy-code": "EW3",
                "clearport-code": "EW3",
                "clearing-code": "W3",
                "clearing-exchange-code": "9C",
                "clearing-price-multiplier": "1.0",
                "display-factor": "0.01",
                "exchange": "CME",
                "product-type": "Physical",
                "expiration-type": "Weekly",
                "settlement-delay-days": 0,
                "is-rollover": false,
                "market-sector": "Equity Index"
            },
            {
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
            },
            {
                "root-symbol": "EW",
                "cash-settled": false,
                "code": "EW",
                "legacy-code": "EW",
                "clearport-code": "EW",
                "clearing-code": "EW",
                "clearing-exchange-code": "16",
                "clearing-price-multiplier": "1.0",
                "display-factor": "0.01",
                "exchange": "CME",
                "product-type": "Physical",
                "expiration-type": "End-Of-Month",
                "settlement-delay-days": 0,
                "is-rollover": true,
                "market-sector": "Equity Index"
            },
            {
                "root-symbol": "E1A",
                "cash-settled": false,
                "code": "E1A",
                "legacy-code": "E1A",
                "clearport-code": "E1A",
                "clearing-code": "A4",
                "clearing-exchange-code": "9C",
                "clearing-price-multiplier": "1.0",
                "display-factor": "0.01",
                "exchange": "CME",
                "product-type": "Physical",
                "expiration-type": "Weekly",
                "settlement-delay-days": 0,
                "is-rollover": false,
                "market-sector": "Equity Index"
            },
            {
                "root-symbol": "E2A",
                "cash-settled": false,
                "code": "E2A",
                "legacy-code": "E2A",
                "clearport-code": "E2A",
                "clearing-code": "A5",
                "clearing-exchange-code": "9C",
                "clearing-price-multiplier": "1.0",
                "display-factor": "0.01",
                "exchange": "CME",
                "product-type": "Physical",
                "expiration-type": "Weekly",
                "settlement-delay-days": 0,
                "is-rollover": false,
                "market-sector": "Equity Index"
            },
            {
                "root-symbol": "E3A",
                "cash-settled": false,
                "code": "E3A",
                "legacy-code": "E3A",
                "clearport-code": "E3A",
                "clearing-code": "A6",
                "clearing-exchange-code": "9C",
                "clearing-price-multiplier": "1.0",
                "display-factor": "0.01",
                "exchange": "CME",
                "product-type": "Physical",
                "expiration-type": "Weekly",
                "settlement-delay-days": 0,
                "is-rollover": false,
                "market-sector": "Equity Index"
            },
            {
                "root-symbol": "E4A",
                "cash-settled": false,
                "code": "E4A",
                "legacy-code": "E4A",
                "clearport-code": "E4A",
                "clearing-code": "A7",
                "clearing-exchange-code": "9C",
                "clearing-price-multiplier": "1.0",
                "display-factor": "0.01",
                "exchange": "CME",
                "product-type": "Physical",
                "expiration-type": "Weekly",
                "settlement-delay-days": 0,
                "is-rollover": true,
                "market-sector": "Equity Index"
            },
            {
                "root-symbol": "E5A",
                "cash-settled": false,
                "code": "E5A",
                "legacy-code": "E5A",
                "clearport-code": "E5A",
                "clearing-code": "A8",
                "clearing-exchange-code": "9C",
                "clearing-price-multiplier": "1.0",
                "display-factor": "0.01",
                "exchange": "CME",
                "product-type": "Physical",
                "expiration-type": "Weekly",
                "settlement-delay-days": 0,
                "is-rollover": true,
                "market-sector": "Equity Index"
            },
            {
                "root-symbol": "E1C",
                "cash-settled": false,
                "code": "E1C",
                "legacy-code": "E1C",
                "clearport-code": "E1C",
                "clearing-code": "Z1",
                "clearing-exchange-code": "9C",
                "clearing-price-multiplier": "1.0",
                "display-factor": "0.01",
                "exchange": "CME",
                "product-type": "Physical",
                "expiration-type": "Weekly",
                "settlement-delay-days": 0,
                "is-rollover": false,
                "market-sector": "Equity Index"
            },
            {
                "root-symbol": "E2C",
                "cash-settled": false,
                "code": "E2C",
                "legacy-code": "E2C",
                "clearport-code": "E2C",
                "clearing-code": "Z2",
                "clearing-exchange-code": "9C",
                "clearing-price-multiplier": "1.0",
                "display-factor": "0.01",
                "exchange": "CME",
                "product-type": "Physical",
                "expiration-type": "Weekly",
                "settlement-delay-days": 0,
                "is-rollover": false,
                "market-sector": "Equity Index"
            },
            {
                "root-symbol": "E4C",
                "cash-settled": false,
                "code": "E4C",
                "legacy-code": "E4C",
                "clearport-code": "E4C",
                "clearing-code": "Z4",
                "clearing-exchange-code": "9C",
                "clearing-price-multiplier": "1.0",
                "display-factor": "0.01",
                "exchange": "CME",
                "product-type": "Physical",
                "expiration-type": "Weekly",
                "settlement-delay-days": 0,
                "is-rollover": true,
                "market-sector": "Equity Index"
            },
            {
                "root-symbol": "E5C",
                "cash-settled": false,
                "code": "E5C",
                "legacy-code": "E5C",
                "clearport-code": "E5C",
                "clearing-code": "Z5",
                "clearing-exchange-code": "9C",
                "clearing-price-multiplier": "1.0",
                "display-factor": "0.01",
                "exchange": "CME",
                "product-type": "Physical",
                "expiration-type": "Weekly",
                "settlement-delay-days": 0,
                "is-rollover": true,
                "market-sector": "Equity Index"
            },
            {
                "root-symbol": "E2D",
                "cash-settled": false,
                "code": "E2D",
                "legacy-code": "E2D",
                "clearport-code": "E2D",
                "clearing-code": "EY",
                "clearing-exchange-code": "9C",
                "clearing-price-multiplier": "1.0",
                "display-factor": "0.01",
                "exchange": "CME",
                "product-type": "Physical",
                "expiration-type": "Weekly",
                "settlement-delay-days": 0,
                "is-rollover": false,
                "market-sector": "Equity Index"
            },
            {
                "root-symbol": "E3D",
                "cash-settled": false,
                "code": "E3D",
                "legacy-code": "E3D",
                "clearport-code": "E3D",
                "clearing-code": "EZ",
                "clearing-exchange-code": "9C",
                "clearing-price-multiplier": "1.0",
                "display-factor": "0.01",
                "exchange": "CME",
                "product-type": "Physical",
                "expiration-type": "Weekly",
                "settlement-delay-days": 0,
                "is-rollover": false,
                "market-sector": "Equity Index"
            },
            {
                "root-symbol": "E4D",
                "cash-settled": false,
                "code": "E4D",
                "legacy-code": "E4D",
                "clearport-code": "E4D",
                "clearing-code": "E1",
                "clearing-exchange-code": "9C",
                "clearing-price-multiplier": "1.0",
                "display-factor": "0.01",
                "exchange": "CME",
                "product-type": "Physical",
                "expiration-type": "Weekly",
                "settlement-delay-days": 0,
                "is-rollover": true,
                "market-sector": "Equity Index"
            },
            {
                "root-symbol": "E5D",
                "cash-settled": false,
                "code": "E5D",
                "legacy-code": "E5D",
                "clearport-code": "E5D",
                "clearing-code": "E2",
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
        ],
        "roll": {
            "name": "equity_index",
            "active-count": 3,
            "cash-settled": true,
            "business-days-offset": 4,
            "first-notice": false
        }
    },
    "context": "/instruments/future-products/CME/ES"
}`

const quantityDecimalPrecisionsResp = `{
  "data": {
    "items": [
      {
        "instrument-type": "Cryptocurrency",
        "symbol": "AAVE/USD",
        "value": 8,
        "minimum-increment-precision": 6
      },
      {
        "instrument-type": "Cryptocurrency",
        "symbol": "ADA/USD",
        "value": 6,
        "minimum-increment-precision": 6
      },
      {
        "instrument-type": "Cryptocurrency",
        "symbol": "BAT/USD",
        "value": 8,
        "minimum-increment-precision": 0
      }
    ]
  },
  "context": "/instruments/quantity-decimal-precisions"
}`

const warrantsResp = `{
  "data": {
    "items": [
      {
        "symbol": "NKLAW",
        "instrument-type": "Warrant",
        "listed-market": "XNAS",
        "description": "Nikola Corporation - Warrant expiring 6/3/2025",
        "is-closing-only": false,
        "active": false
      }
    ]
  },
  "context": "/instruments/warrants"
}`

const warrantResp = `{
    "data": {
        "symbol": "NKLAW",
        "instrument-type": "Warrant",
        "listed-market": "XNAS",
        "description": "Nikola Corporation - Warrant expiring 6/3/2025",
        "is-closing-only": false,
        "active": false
    },
    "context": "/instruments/warrants/NKLAW"
}`
