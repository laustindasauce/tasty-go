package tasty //nolint:testpackage // testing private field

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestGetFutures(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/instruments/futures", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, futuresResp)
	})

	productCode := "ES"

	resp, err := client.GetFutures(FuturesQuery{ProductCode: []string{productCode}})
	require.Nil(t, err)

	require.Equal(t, 1, len(resp))

	future := resp[0]

	require.Equal(t, "/ESM3", future.Symbol)
	require.Equal(t, productCode, future.ProductCode)
	require.Equal(t, StringToFloat32(50), future.ContractSize)
	require.Equal(t, StringToFloat32(.25), future.TickSize)
	require.Equal(t, StringToFloat32(50), future.NotionalMultiplier)
	require.Equal(t, StringToFloat32(0), future.MainFraction)
	require.Equal(t, StringToFloat32(0), future.SubFraction)
	require.Equal(t, StringToFloat32(0.01), future.DisplayFactor)
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
	require.Equal(t, StringToFloat32(50), futureProd.NotionalMultiplier)
	require.Equal(t, StringToFloat32(.25), futureProd.TickSize)
	require.Equal(t, StringToFloat32(.01), futureProd.DisplayFactor)
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

	require.Equal(t, StringToFloat32(.25), future.TickSizes[0].Value)
	require.Equal(t, StringToFloat32(.05), future.OptionTickSizes[0].Value)
	require.Equal(t, StringToFloat32(5), future.OptionTickSizes[0].Threshold)
	require.Equal(t, StringToFloat32(.05), future.SpreadTickSizes[0].Value)
	require.Equal(t, "/ESU3", future.SpreadTickSizes[0].Symbol)
}

func TestGetFuturesError(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/instruments/futures", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(401)
		fmt.Fprint(writer, tastyUnauthorizedError)
	})

	productCode := "ES"

	_, err := client.GetFutures(FuturesQuery{ProductCode: []string{productCode}})
	expectedUnauthorized(t, err)
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
	require.Equal(t, StringToFloat32(50), future.ContractSize)
	require.Equal(t, StringToFloat32(.25), future.TickSize)
	require.Equal(t, StringToFloat32(50), future.NotionalMultiplier)
	require.Equal(t, StringToFloat32(0), future.MainFraction)
	require.Equal(t, StringToFloat32(0), future.SubFraction)
	require.Equal(t, StringToFloat32(0.01), future.DisplayFactor)
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
	require.Equal(t, StringToFloat32(50), futureProd.NotionalMultiplier)
	require.Equal(t, StringToFloat32(.25), futureProd.TickSize)
	require.Equal(t, StringToFloat32(.01), futureProd.DisplayFactor)
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

	require.Equal(t, StringToFloat32(.25), future.TickSizes[0].Value)
	require.Equal(t, StringToFloat32(.05), future.OptionTickSizes[0].Value)
	require.Equal(t, StringToFloat32(5), future.OptionTickSizes[0].Threshold)
	require.Equal(t, StringToFloat32(.05), future.SpreadTickSizes[0].Value)
	require.Equal(t, "/ESU3", future.SpreadTickSizes[0].Symbol)
}

func TestGetFutureError(t *testing.T) {
	setup()
	defer teardown()

	symbol := "ESM3"

	mux.HandleFunc(fmt.Sprintf("/instruments/futures/%s", symbol), func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(401)
		fmt.Fprint(writer, tastyUnauthorizedError)
	})

	_, err := client.GetFuture(symbol)
	expectedUnauthorized(t, err)
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
	require.Equal(t, StringToFloat32(100), fop.ClearingPriceMultiplier)
	require.Equal(t, StringToFloat32(0.01), fop.DisplayFactor)
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
	require.Equal(t, StringToFloat32(1.0), futureProd.NotionalMultiplier)
	require.Equal(t, StringToFloat32(.25), futureProd.TickSize)
	require.Equal(t, StringToFloat32(1.0), futureProd.DisplayFactor)
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

func TestGetFutureOptionProductsError(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/instruments/future-option-products", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(401)
		fmt.Fprint(writer, tastyUnauthorizedError)
	})

	_, err := client.GetFutureOptionProducts()
	expectedUnauthorized(t, err)
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
	require.Equal(t, StringToFloat32(100), fop.ClearingPriceMultiplier)
	require.Equal(t, StringToFloat32(0.01), fop.DisplayFactor)
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
	require.Equal(t, StringToFloat32(1.0), futureProd.NotionalMultiplier)
	require.Equal(t, StringToFloat32(.25), futureProd.TickSize)
	require.Equal(t, StringToFloat32(1.0), futureProd.DisplayFactor)
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

func TestGetFutureOptionProductError(t *testing.T) {
	setup()
	defer teardown()

	rootSymbol := "ECYM"
	exchange := "CME"

	mux.HandleFunc(fmt.Sprintf("/instruments/future-option-products/%s/%s", exchange, rootSymbol), func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(401)
		fmt.Fprint(writer, tastyUnauthorizedError)
	})

	_, err := client.GetFutureOptionProduct(exchange, rootSymbol)
	expectedUnauthorized(t, err)
}

func TestGetFutureOptions(t *testing.T) {
	setup()
	defer teardown()

	future := FutureSymbology{ProductCode: "ES", MonthCode: December, YearDigit: 9}

	expiry := time.Date(2019, 9, 27, 0, 0, 0, 0, time.Local)
	fcc := FutureOptionsSymbology{
		OptionContractCode: "EW4U9",
		FutureContractCode: future.Build(),
		OptionType:         Put,
		Strike:             2975,
		Expiration:         expiry,
	}

	symbol := fcc.Build()

	query := FutureOptionsQuery{
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
	require.Equal(t, StringToFloat32(2975), fo.StrikePrice)
	require.Equal(t, "CME", fo.Exchange)
	require.Equal(t, "EW4U9 P2975", fo.ExchangeSymbol)
	require.Equal(t, Put, fo.OptionType)
	require.Equal(t, "American", fo.ExerciseStyle)
	require.True(t, fo.IsVanilla)
	require.True(t, fo.IsPrimaryDeliverable)
	require.Equal(t, StringToFloat32(1), fo.FuturePriceRatio)
	require.Equal(t, StringToFloat32(1), fo.Multiplier)
	require.Equal(t, StringToFloat32(1), fo.UnderlyingCount)
	require.True(t, fo.IsConfirmed)
	require.Equal(t, StringToFloat32(.5), fo.NotionalValue)
	require.Equal(t, StringToFloat32(.01), fo.DisplayFactor)
	require.Equal(t, "2", fo.SecurityExchange)
	require.Equal(t, "0", fo.SxID)
	require.Equal(t, "Future", fo.SettlementType)
	require.Equal(t, StringToFloat32(1), fo.StrikeFactor)
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
	require.Equal(t, StringToFloat32(1), fop.ClearingPriceMultiplier)
	require.Equal(t, StringToFloat32(.01), fop.DisplayFactor)
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

	future := FutureSymbology{ProductCode: "ES", MonthCode: December, YearDigit: 9}

	expiry := time.Date(2019, 9, 27, 0, 0, 0, 0, time.Local)
	fcc := FutureOptionsSymbology{
		OptionContractCode: "EW4U9",
		FutureContractCode: future.Build(),
		OptionType:         Put,
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
	require.Equal(t, StringToFloat32(2975), fo.StrikePrice)
	require.Equal(t, "CME", fo.Exchange)
	require.Equal(t, "EW4U9 P2975", fo.ExchangeSymbol)
	require.Equal(t, Put, fo.OptionType)
	require.Equal(t, "American", fo.ExerciseStyle)
	require.True(t, fo.IsVanilla)
	require.True(t, fo.IsPrimaryDeliverable)
	require.Equal(t, StringToFloat32(1), fo.FuturePriceRatio)
	require.Equal(t, StringToFloat32(1), fo.Multiplier)
	require.Equal(t, StringToFloat32(1), fo.UnderlyingCount)
	require.True(t, fo.IsConfirmed)
	require.Equal(t, StringToFloat32(.5), fo.NotionalValue)
	require.Equal(t, StringToFloat32(.01), fo.DisplayFactor)
	require.Equal(t, "2", fo.SecurityExchange)
	require.Equal(t, "0", fo.SxID)
	require.Equal(t, "Future", fo.SettlementType)
	require.Equal(t, StringToFloat32(1), fo.StrikeFactor)
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
	require.Equal(t, StringToFloat32(1), fop.ClearingPriceMultiplier)
	require.Equal(t, StringToFloat32(.01), fop.DisplayFactor)
	require.Equal(t, "CME", fop.Exchange)
	require.Equal(t, "Physical", fop.ProductType)
	require.Equal(t, "Weekly", fop.ExpirationType)
	require.Equal(t, 0, fop.SettlementDelayDays)
	require.True(t, fop.IsRollover)
	require.Equal(t, "Equity Index", fop.MarketSector)
}

func TestGetFutureOptionsError(t *testing.T) {
	setup()
	defer teardown()

	future := FutureSymbology{ProductCode: "ES", MonthCode: December, YearDigit: 9}

	expiry := time.Date(2019, 9, 27, 0, 0, 0, 0, time.Local)
	fcc := FutureOptionsSymbology{
		OptionContractCode: "EW4U9",
		FutureContractCode: future.Build(),
		OptionType:         Put,
		Strike:             2975,
		Expiration:         expiry,
	}

	symbol := fcc.Build()

	query := FutureOptionsQuery{
		Symbols: []string{symbol},
	}

	mux.HandleFunc("/instruments/future-options", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(401)
		fmt.Fprint(writer, tastyUnauthorizedError)
	})

	_, err := client.GetFutureOptions(query)
	expectedUnauthorized(t, err)
}

func TestGetFutureProduct(t *testing.T) {
	setup()
	defer teardown()

	exchange := CME
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
	require.Equal(t, StringToFloat32(50), fp.NotionalMultiplier)
	require.Equal(t, StringToFloat32(.25), fp.TickSize)
	require.Equal(t, StringToFloat32(.01), fp.DisplayFactor)
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
	require.Equal(t, StringToFloat32(1), op.ClearingPriceMultiplier)
	require.Equal(t, StringToFloat32(.01), op.DisplayFactor)
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

func TestGetFutureProductError(t *testing.T) {
	setup()
	defer teardown()

	exchange := CME
	code := "ES"

	mux.HandleFunc(fmt.Sprintf("/instruments/future-products/%s/%s", exchange, code), func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(401)
		fmt.Fprint(writer, tastyUnauthorizedError)
	})

	_, err := client.GetFutureProduct(exchange, code)
	expectedUnauthorized(t, err)
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
	require.Equal(t, StringToFloat32(50), fp.NotionalMultiplier)
	require.Equal(t, StringToFloat32(.25), fp.TickSize)
	require.Equal(t, StringToFloat32(.01), fp.DisplayFactor)
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
	require.Equal(t, StringToFloat32(1), op.ClearingPriceMultiplier)
	require.Equal(t, StringToFloat32(.01), op.DisplayFactor)
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

func TestGetFutureProductsError(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/instruments/future-products", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(401)
		fmt.Fprint(writer, tastyUnauthorizedError)
	})

	_, err := client.GetFutureProducts()
	expectedUnauthorized(t, err)
}

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
