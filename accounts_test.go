package tasty //nolint:testpackage // testing private field

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestGetMyAccounts(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/customers/me/accounts", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, myAccountsResp)
	})

	resp, err := client.GetMyAccounts()
	require.Nil(t, err)

	require.Equal(t, 3, len(resp))

	roth := resp[0]

	require.Equal(t, "5YZ55555", roth.AccountNumber)
	require.Equal(t, "A1d557b2a-e5f1-483a-9798-13923403f442", roth.ExternalID)
	require.Equal(t, "2022-10-27T20:49:52.79Z", roth.OpenedAt.Format(time.RFC3339Nano))
	require.Equal(t, "Roth IRA", roth.Nickname)
	require.Equal(t, "Roth IRA", roth.AccountTypeName)
	require.False(t, roth.DayTraderStatus)
	require.False(t, roth.IsClosed)
	require.False(t, roth.IsFirmError)
	require.False(t, roth.IsFirmProprietary)
	require.False(t, roth.IsFuturesApproved)
	require.False(t, roth.IsTestDrive)
	require.Equal(t, "Cash", roth.MarginOrCash)
	require.False(t, roth.IsForeign)
	require.Equal(t, "2022-11-04", roth.FundingDate)
	require.Equal(t, "GROWTH", roth.InvestmentObjective)
	require.Equal(t, "Defined Risk Spreads", roth.SuitableOptionsLevel)
	require.Equal(t, "2022-10-27T20:49:52.793Z", roth.CreatedAt.Format(time.RFC3339Nano))
}

func TestGetMyAccountsError(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/customers/me/accounts", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(401)
		fmt.Fprint(writer, tastyUnauthorizedError)
	})

	_, err := client.GetMyAccounts()
	expectedUnauthorized(t, err)
}

func TestGetAccountTradingStatus(t *testing.T) {
	setup()
	defer teardown()

	accountNumber := "5YZ55555"

	mux.HandleFunc(fmt.Sprintf("/accounts/%s/trading-status", accountNumber), func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, accountTradingStatusResp)
	})

	resp, err := client.GetAccountTradingStatus(accountNumber)
	require.Nil(t, err)

	require.Equal(t, "5YZ55555", resp.AccountNumber)
	require.Equal(t, 0, resp.DayTradeCount)
	require.Equal(t, "IRA Margin", resp.EquitiesMarginCalculationType)
	require.Equal(t, "default", resp.FeeScheduleName)
	require.Equal(t, StringToFloat32(0.0), resp.FuturesMarginRateMultiplier)
	require.False(t, resp.HasIntradayEquitiesMargin)
	require.Equal(t, 447096, resp.ID)
	require.False(t, resp.IsAggregatedAtClearing)
	require.False(t, resp.IsClosed)
	require.False(t, resp.IsClosingOnly)
	require.False(t, resp.IsCryptocurrencyClosingOnly)
	require.False(t, resp.IsCryptocurrencyEnabled)
	require.False(t, resp.IsFrozen)
	require.True(t, resp.IsFullEquityMarginRequired)
	require.False(t, resp.IsFuturesClosingOnly)
	require.False(t, resp.IsFuturesIntraDayEnabled)
	require.False(t, resp.IsFuturesEnabled)
	require.False(t, resp.IsInDayTradeEquityMaintenanceCall)
	require.False(t, resp.IsInMarginCall)
	require.False(t, resp.IsPatternDayTrader)
	require.False(t, resp.IsPortfolioMarginEnabled)
	require.False(t, resp.IsRiskReducingOnly)
	require.False(t, resp.IsSmallNotionalFuturesIntraDayEnabled)
	require.True(t, resp.IsRollTheDayForwardEnabled)
	require.True(t, resp.AreFarOtmNetOptionsRestricted)
	require.Equal(t, "Defined Risk Spreads", resp.OptionsLevel)
	require.False(t, resp.ShortCallsEnabled)
	require.Equal(t, StringToFloat32(0.0), resp.SmallNotionalFuturesMarginRateMultiplier)
	require.False(t, resp.IsEquityOfferingEnabled)
	require.False(t, resp.IsEquityOfferingClosingOnly)
	require.Equal(t, "2022-10-27T20:49:52.928Z", resp.EnhancedFraudSafeguardsEnabledAt.Format(time.RFC3339Nano))
	require.Equal(t, "2023-05-28T20:44:40.32Z", resp.UpdatedAt.Format(time.RFC3339Nano))
}

func TestGetAccountTradingStatusError(t *testing.T) {
	setup()
	defer teardown()

	accountNumber := "5YZ55555"

	mux.HandleFunc(fmt.Sprintf("/accounts/%s/trading-status", accountNumber), func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(401)
		fmt.Fprint(writer, tastyUnauthorizedError)
	})

	_, err := client.GetAccountTradingStatus(accountNumber)
	expectedUnauthorized(t, err)
}

func TestGetAccountBalances(t *testing.T) {
	setup()
	defer teardown()

	accountNumber := "5YZ55555"

	mux.HandleFunc(fmt.Sprintf("/accounts/%s/balances", accountNumber), func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, accountBalancesResp)
	})

	resp, err := client.GetAccountBalances(accountNumber)
	require.Nil(t, err)

	require.Equal(t, "5YZ55555", resp.AccountNumber)
	require.Equal(t, StringToFloat32(51600.762), resp.CashBalance)
	require.Equal(t, StringToFloat32(281983.415), resp.LongEquityValue)
	require.Equal(t, StringToFloat32(0.0), resp.ShortEquityValue)
	require.Equal(t, StringToFloat32(0.0), resp.LongDerivativeValue)
	require.Equal(t, StringToFloat32(82680.5), resp.ShortDerivativeValue)
	require.Equal(t, StringToFloat32(0.0), resp.LongFuturesValue)
	require.Equal(t, StringToFloat32(0.0), resp.ShortFuturesValue)
	require.Equal(t, StringToFloat32(0.0), resp.LongFuturesDerivativeValue)
	require.Equal(t, StringToFloat32(0.0), resp.ShortFuturesDerivativeValue)
	require.Equal(t, StringToFloat32(0.0), resp.LongMargineableValue)
	require.Equal(t, StringToFloat32(0.0), resp.ShortMargineableValue)
	require.Equal(t, StringToFloat32(452284.177), resp.MarginEquity)
	require.Equal(t, StringToFloat32(20078.762), resp.EquityBuyingPower)
	require.Equal(t, StringToFloat32(20078.762), resp.DerivativeBuyingPower)
	require.Equal(t, StringToFloat32(0.0), resp.DayTradingBuyingPower)
	require.Equal(t, StringToFloat32(0.0), resp.FuturesMarginRequirement)
	require.Equal(t, StringToFloat32(0.0), resp.AvailableTradingFunds)
	require.Equal(t, StringToFloat32(432279.234), resp.MaintenanceRequirement)
	require.Equal(t, StringToFloat32(0.0), resp.MaintenanceCallValue)
	require.Equal(t, StringToFloat32(0.0), resp.RegTCallValue)
	require.Equal(t, StringToFloat32(0.0), resp.DayTradingCallValue)
	require.Equal(t, StringToFloat32(0.0), resp.DayEquityCallValue)
	require.Equal(t, StringToFloat32(543557.677), resp.NetLiquidatingValue)
	require.Equal(t, StringToFloat32(20078.76), resp.CashAvailableToWithdraw)
	require.Equal(t, StringToFloat32(0.0), resp.EquityOfferingMarginRequirement)
	require.Equal(t, StringToFloat32(0.0), resp.LongBondValue)
	require.Equal(t, StringToFloat32(0.0), resp.BondMarginRequirement)
	require.Equal(t, "2023-06-08", resp.SnapshotDate)
	require.Equal(t, StringToFloat32(432279.2338), resp.RegTMarginRequirement)
	require.Equal(t, StringToFloat32(0.0), resp.FuturesOvernightMarginRequirement)
	require.Equal(t, StringToFloat32(0.0), resp.FuturesIntradayMarginRequirement)
	require.Equal(t, StringToFloat32(20078.762), resp.MaintenanceExcess)
	require.Equal(t, StringToFloat32(0.0), resp.PendingMarginInterest)
	require.Equal(t, StringToFloat32(20078.76), resp.EffectiveCryptocurrencyBuyingPower)
	require.Equal(t, "2023-06-08T16:30:18.889Z", resp.UpdatedAt.Format(time.RFC3339Nano))
}

func TestGetAccountBalancesError(t *testing.T) {
	setup()
	defer teardown()

	accountNumber := "5YZ55555"

	mux.HandleFunc(fmt.Sprintf("/accounts/%s/balances", accountNumber), func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(401)
		fmt.Fprint(writer, tastyUnauthorizedError)
	})

	_, err := client.GetAccountBalances(accountNumber)
	expectedUnauthorized(t, err)
}

func TestGetAccountPositions(t *testing.T) {
	setup()
	defer teardown()

	accountNumber := "5YZ55555"
	query := AccountPositionQuery{UnderlyingSymbol: []string{"RIVN"}}

	mux.HandleFunc(fmt.Sprintf("/accounts/%s/positions", accountNumber), func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, accountPositionsResp)
	})

	resp, err := client.GetAccountPositions(accountNumber, query)
	require.Nil(t, err)

	rivn := resp[0]

	require.Equal(t, "5YZ55555", rivn.AccountNumber)
	require.Equal(t, "RIVN  230609P00014000", rivn.Symbol)
	require.Equal(t, EquityOptionIT, rivn.InstrumentType)
	require.Equal(t, "RIVN", rivn.UnderlyingSymbol)
	require.Equal(t, 40, rivn.Quantity)
	require.Equal(t, Short, rivn.QuantityDirection)
	require.Equal(t, StringToFloat32(0.41), rivn.ClosePrice)
	require.Equal(t, StringToFloat32(0.79), rivn.AverageOpenPrice)
	require.Equal(t, StringToFloat32(0.79), rivn.AverageYearlyMarketClosePrice)
	require.Equal(t, StringToFloat32(0.41), rivn.AverageDailyMarketClosePrice)
	require.Equal(t, 100, rivn.Multiplier)
	require.Equal(t, Debit, rivn.CostEffect)
	require.False(t, rivn.IsSuppressed)
	require.False(t, rivn.IsFrozen)
	require.Equal(t, 0, rivn.RestrictedQuantity)
	require.Equal(t, "2023-06-09T20:00:00Z", rivn.ExpiresAt.Format(time.RFC3339))
	require.Equal(t, StringToFloat32(0.0), rivn.RealizedDayGain)
	require.Equal(t, None, rivn.RealizedDayGainEffect)
	require.Equal(t, "2023-05-24", rivn.RealizedDayGainDate)
	require.Equal(t, StringToFloat32(0.0), rivn.RealizedToday)
	require.Equal(t, None, rivn.RealizedTodayEffect)
	require.Equal(t, "2023-05-24", rivn.RealizedTodayDate)
	require.Equal(t, "2023-05-24T17:17:57.615Z", rivn.CreatedAt.Format(time.RFC3339Nano))
	require.Equal(t, "2023-05-24T17:17:58.632Z", rivn.UpdatedAt.Format(time.RFC3339Nano))
}

func TestGetAccountPositionsError(t *testing.T) {
	setup()
	defer teardown()

	accountNumber := "5YZ55555"
	query := AccountPositionQuery{UnderlyingSymbol: []string{"RIVN"}}

	mux.HandleFunc(fmt.Sprintf("/accounts/%s/positions", accountNumber), func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(401)
		fmt.Fprint(writer, tastyUnauthorizedError)
	})

	_, err := client.GetAccountPositions(accountNumber, query)
	expectedUnauthorized(t, err)
}

func TestGetAccountBalanceSnapshots(t *testing.T) {
	setup()
	defer teardown()

	accountNumber := "5YZ55555"
	query := AccountBalanceSnapshotsQuery{SnapshotDate: time.Now().AddDate(0, -2, 0)}

	mux.HandleFunc(fmt.Sprintf("/accounts/%s/balance-snapshots", accountNumber), func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, balanceSnapshotsResp)
	})

	resp, err := client.GetAccountBalanceSnapshots(accountNumber, query)
	require.Nil(t, err)

	snap := resp[0]

	require.Equal(t, "5YZ55555", snap.AccountNumber)
	require.Equal(t, StringToFloat32(51600.762), snap.CashBalance)
	require.Equal(t, StringToFloat32(281983.965), snap.LongEquityValue)
	require.Equal(t, StringToFloat32(0.0), snap.ShortEquityValue)
	require.Equal(t, StringToFloat32(0.0), snap.LongDerivativeValue)
	require.Equal(t, StringToFloat32(82680.0), snap.ShortDerivativeValue)
	require.Equal(t, StringToFloat32(0.0), snap.LongFuturesValue)
	require.Equal(t, StringToFloat32(0.0), snap.ShortFuturesValue)
	require.Equal(t, StringToFloat32(0.0), snap.LongFuturesDerivativeValue)
	require.Equal(t, StringToFloat32(0.0), snap.ShortFuturesDerivativeValue)
	require.Equal(t, StringToFloat32(0.0), snap.LongMargineableValue)
	require.Equal(t, StringToFloat32(0.0), snap.ShortMargineableValue)
	require.Equal(t, StringToFloat32(452284.727), snap.MarginEquity)
	require.Equal(t, StringToFloat32(20078.762), snap.EquityBuyingPower)
	require.Equal(t, StringToFloat32(20078.762), snap.DerivativeBuyingPower)
	require.Equal(t, StringToFloat32(0.0), snap.DayTradingBuyingPower)
	require.Equal(t, StringToFloat32(0.0), snap.FuturesMarginRequirement)
	require.Equal(t, StringToFloat32(0.0), snap.AvailableTradingFunds)
	require.Equal(t, StringToFloat32(432279.047), snap.MaintenanceRequirement)
	require.Equal(t, StringToFloat32(0.0), snap.MaintenanceCallValue)
	require.Equal(t, StringToFloat32(0.0), snap.RegTCallValue)
	require.Equal(t, StringToFloat32(0.0), snap.DayTradingCallValue)
	require.Equal(t, StringToFloat32(0.0), snap.DayEquityCallValue)
	require.Equal(t, StringToFloat32(4544.727), snap.NetLiquidatingValue)
	require.Equal(t, StringToFloat32(20078.76), snap.CashAvailableToWithdraw)
	require.Equal(t, StringToFloat32(20078.76), snap.DayTradeExcess)
	require.Equal(t, StringToFloat32(0.0), snap.PendingCash)
	require.Equal(t, None, snap.PendingCashEffect)
	require.Equal(t, StringToFloat32(0.0), snap.LongCryptocurrencyValue)
	require.Equal(t, StringToFloat32(0.0), snap.ShortCryptocurrencyValue)
	require.Equal(t, StringToFloat32(0.0), snap.CryptocurrencyMarginRequirement)
	require.Equal(t, StringToFloat32(0.0), snap.UnsettledCryptocurrencyFiatAmount)
	require.Equal(t, None, snap.UnsettledCryptocurrencyFiatEffect)
	require.Equal(t, StringToFloat32(20078.76), snap.ClosedLoopAvailableBalance)
	require.Equal(t, StringToFloat32(0.0), snap.EquityOfferingMarginRequirement)
	require.Equal(t, StringToFloat32(0.0), snap.LongBondValue)
	require.Equal(t, StringToFloat32(0.0), snap.BondMarginRequirement)
	require.Equal(t, "2023-06-08", snap.SnapshotDate)
	require.Equal(t, StringToFloat32(432279.0465), snap.RegTMarginRequirement)
	require.Equal(t, StringToFloat32(0.0), snap.FuturesOvernightMarginRequirement)
	require.Equal(t, StringToFloat32(0.0), snap.FuturesIntradayMarginRequirement)
	require.Equal(t, StringToFloat32(20078.762), snap.MaintenanceExcess)
	require.Equal(t, StringToFloat32(0.0), snap.PendingMarginInterest)
	require.Equal(t, StringToFloat32(20078.76), snap.EffectiveCryptocurrencyBuyingPower)
	require.Equal(t, "2023-06-08T18:37:39.568Z", snap.UpdatedAt.Format(time.RFC3339Nano))
}

func TestGetAccountBalanceSnapshotsError(t *testing.T) {
	setup()
	defer teardown()

	accountNumber := "5YZ55555"
	query := AccountBalanceSnapshotsQuery{SnapshotDate: time.Now().AddDate(0, -2, 0)}

	mux.HandleFunc(fmt.Sprintf("/accounts/%s/balance-snapshots", accountNumber), func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(401)
		fmt.Fprint(writer, tastyUnauthorizedError)
	})

	_, err := client.GetAccountBalanceSnapshots(accountNumber, query)
	expectedUnauthorized(t, err)
}

func TestGetAccountNetLiqHistory(t *testing.T) {
	setup()
	defer teardown()

	accountNumber := "5YZ55555"
	query := HistoricLiquidityQuery{TimeBack: OneDay}

	mux.HandleFunc(fmt.Sprintf("/accounts/%s/net-liq/history", accountNumber), func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, netLiqHistoryResp)
	})

	resp, err := client.GetAccountNetLiqHistory(accountNumber, query)
	require.Nil(t, err)

	require.Equal(t, 3, len(resp))

	liq := resp[0]

	require.Equal(t, StringToFloat32(4498.667), liq.Open)
	require.Equal(t, StringToFloat32(4498.667), liq.High)
	require.Equal(t, StringToFloat32(4498.667), liq.Low)
	require.Equal(t, StringToFloat32(4498.667), liq.Close)
	require.Equal(t, StringToFloat32(0.0), liq.PendingCashOpen)
	require.Equal(t, StringToFloat32(0.0), liq.PendingCashHigh)
	require.Equal(t, StringToFloat32(0.0), liq.PendingCashLow)
	require.Equal(t, StringToFloat32(0.0), liq.PendingCashClose)
	require.Equal(t, StringToFloat32(4498.667), liq.TotalOpen)
	require.Equal(t, StringToFloat32(4498.667), liq.TotalHigh)
	require.Equal(t, StringToFloat32(4498.667), liq.TotalLow)
	require.Equal(t, StringToFloat32(4498.667), liq.TotalClose)
	require.Equal(t, "2023-06-08 13:30:00+00", liq.Time)
}

func TestGetAccountNetLiqHistoryError(t *testing.T) {
	setup()
	defer teardown()

	accountNumber := "5YZ55555"
	query := HistoricLiquidityQuery{TimeBack: OneDay}

	mux.HandleFunc(fmt.Sprintf("/accounts/%s/net-liq/history", accountNumber), func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(401)
		fmt.Fprint(writer, tastyUnauthorizedError)
	})

	_, err := client.GetAccountNetLiqHistory(accountNumber, query)
	expectedUnauthorized(t, err)
}

func TestGetAccountPositionLimit(t *testing.T) {
	setup()
	defer teardown()

	accountNumber := "5YZ55555"

	mux.HandleFunc(fmt.Sprintf("/accounts/%s/position-limit", accountNumber), func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, positionLimitResp)
	})

	resp, err := client.GetAccountPositionLimit(accountNumber)
	require.Nil(t, err)

	require.Equal(t, accountNumber, resp.AccountNumber)
	require.Equal(t, 50000, resp.EquityOrderSize)
	require.Equal(t, 10000, resp.EquityOptionOrderSize)
	require.Equal(t, 2500, resp.FutureOrderSize)
	require.Equal(t, 2500, resp.FutureOptionOrderSize)
	require.Equal(t, 500, resp.UnderlyingOpeningOrderLimit)
	require.Equal(t, 500000, resp.EquityPositionSize)
	require.Equal(t, 20000, resp.EquityOptionPositionSize)
	require.Equal(t, 5000, resp.FuturePositionSize)
	require.Equal(t, 5000, resp.FutureOptionPositionSize)
}

func TestGetAccountPositionLimitError(t *testing.T) {
	setup()
	defer teardown()

	accountNumber := "5YZ55555"

	mux.HandleFunc(fmt.Sprintf("/accounts/%s/position-limit", accountNumber), func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(401)
		fmt.Fprint(writer, tastyUnauthorizedError)
	})

	_, err := client.GetAccountPositionLimit(accountNumber)
	expectedUnauthorized(t, err)
}

const myAccountsResp = `{
  "data": {
    "items": [
      {
        "account": {
          "account-number": "5YZ55555",
          "external-id": "A1d557b2a-e5f1-483a-9798-13923403f442",
          "opened-at": "2022-10-27T20:49:52.790+00:00",
          "nickname": "Roth IRA",
          "account-type-name": "Roth IRA",
          "day-trader-status": false,
          "is-closed": false,
          "is-firm-error": false,
          "is-firm-proprietary": false,
          "is-futures-approved": false,
          "is-test-drive": false,
          "margin-or-cash": "Cash",
          "is-foreign": false,
          "funding-date": "2022-11-04",
          "investment-objective": "GROWTH",
          "suitable-options-level": "Defined Risk Spreads",
          "created-at": "2022-10-27T20:49:52.793+00:00"
        },
        "authority-level": "owner"
      },
      {
        "account": {
          "account-number": "5WW55555",
          "external-id": "A0002465882",
          "opened-at": "2021-03-26T18:31:41.070+00:00",
          "nickname": "Margin acct",
          "account-type-name": "Individual",
          "day-trader-status": false,
          "is-closed": false,
          "is-firm-error": false,
          "is-firm-proprietary": false,
          "is-futures-approved": false,
          "is-test-drive": false,
          "margin-or-cash": "Margin",
          "is-foreign": false,
          "funding-date": "2021-04-07",
          "investment-objective": "GROWTH",
          "suitable-options-level": "No Restrictions",
          "created-at": "2021-03-26T18:31:41.078+00:00"
        },
        "authority-level": "owner"
      },
      {
        "account": {
          "account-number": "5WZ55555",
          "external-id": "A0001236586",
          "opened-at": "2019-12-27T18:12:03.420+00:00",
          "nickname": "Individual",
          "account-type-name": "Individual",
          "day-trader-status": false,
          "is-closed": false,
          "is-firm-error": false,
          "is-firm-proprietary": false,
          "is-futures-approved": false,
          "is-test-drive": false,
          "margin-or-cash": "Cash",
          "is-foreign": false,
          "funding-date": "2021-03-82680",
          "investment-objective": "GROWTH",
          "suitable-options-level": "Covered And Cash Secured",
          "created-at": "2019-12-27T18:12:03.424+00:00"
        },
        "authority-level": "owner"
      }
    ]
  },
  "context": "/customers/me/accounts"
}`

const accountTradingStatusResp = `{
  "data": {
    "account-number": "5YZ55555",
    "day-trade-count": 0,
    "equities-margin-calculation-type": "IRA Margin",
    "fee-schedule-name": "default",
    "futures-margin-rate-multiplier": "0.0",
    "has-intraday-equities-margin": false,
    "id": 447096,
    "is-aggregated-at-clearing": false,
    "is-closed": false,
    "is-closing-only": false,
    "is-cryptocurrency-closing-only": false,
    "is-cryptocurrency-enabled": false,
    "is-frozen": false,
    "is-full-equity-margin-required": true,
    "is-futures-closing-only": false,
    "is-futures-intra-day-enabled": false,
    "is-futures-enabled": false,
    "is-in-day-trade-equity-maintenance-call": false,
    "is-in-margin-call": false,
    "is-pattern-day-trader": false,
    "is-portfolio-margin-enabled": false,
    "is-risk-reducing-only": false,
    "is-small-notional-futures-intra-day-enabled": false,
    "is-roll-the-day-forward-enabled": true,
    "are-far-otm-net-options-restricted": true,
    "options-level": "Defined Risk Spreads",
    "short-calls-enabled": false,
    "small-notional-futures-margin-rate-multiplier": "0.0",
    "is-equity-offering-enabled": false,
    "is-equity-offering-closing-only": false,
    "enhanced-fraud-safeguards-enabled-at": "2022-10-27T20:49:52.928+00:00",
    "updated-at": "2023-05-28T20:44:40.320+00:00"
  },
  "context": "/accounts/5YZ55555/trading-status"
}`

const accountBalancesResp = `{
  "data": {
    "account-number": "5YZ55555",
    "cash-balance": "51600.762",
    "long-equity-value": "281983.415",
    "short-equity-value": "0.0",
    "long-derivative-value": "0.0",
    "short-derivative-value": "82680.5",
    "long-futures-value": "0.0",
    "short-futures-value": "0.0",
    "long-futures-derivative-value": "0.0",
    "short-futures-derivative-value": "0.0",
    "long-margineable-value": "0.0",
    "short-margineable-value": "0.0",
    "margin-equity": "452284.177",
    "equity-buying-power": "20078.762",
    "derivative-buying-power": "20078.762",
    "day-trading-buying-power": "0.0",
    "futures-margin-requirement": "0.0",
    "available-trading-funds": "0.0",
    "maintenance-requirement": "432279.234",
    "maintenance-call-value": "0.0",
    "reg-t-call-value": "0.0",
    "day-trading-call-value": "0.0",
    "day-equity-call-value": "0.0",
    "net-liquidating-value": "543557.677",
    "cash-available-to-withdraw": "20078.76",
    "day-trade-excess": "20078.76",
    "pending-cash": "0.0",
    "pending-cash-effect": "None",
    "long-cryptocurrency-value": "0.0",
    "short-cryptocurrency-value": "0.0",
    "cryptocurrency-margin-requirement": "0.0",
    "unsettled-cryptocurrency-fiat-amount": "0.0",
    "unsettled-cryptocurrency-fiat-effect": "None",
    "closed-loop-available-balance": "20078.76",
    "equity-offering-margin-requirement": "0.0",
    "long-bond-value": "0.0",
    "bond-margin-requirement": "0.0",
    "snapshot-date": "2023-06-08",
    "reg-t-margin-requirement": "432279.2338",
    "futures-overnight-margin-requirement": "0.0",
    "futures-intraday-margin-requirement": "0.0",
    "maintenance-excess": "20078.762",
    "pending-margin-interest": "0.0",
    "effective-cryptocurrency-buying-power": "20078.76",
    "updated-at": "2023-06-08T16:30:18.889+00:00"
  },
  "context": "/accounts/5YZ55555/balances"
}`

const accountPositionsResp = `{
  "data": {
    "items": [
      {
        "account-number": "5YZ55555",
        "symbol": "RIVN  230609P00014000",
        "instrument-type": "Equity Option",
        "underlying-symbol": "RIVN",
        "quantity": 40,
        "quantity-direction": "Short",
        "close-price": "0.41",
        "average-open-price": "0.79",
        "average-yearly-market-close-price": "0.79",
        "average-daily-market-close-price": "0.41",
        "multiplier": 100,
        "cost-effect": "Debit",
        "is-suppressed": false,
        "is-frozen": false,
        "restricted-quantity": 0,
        "expires-at": "2023-06-09T20:00:00.000+00:00",
        "realized-day-gain": "0.0",
        "realized-day-gain-effect": "None",
        "realized-day-gain-date": "2023-05-24",
        "realized-today": "0.0",
        "realized-today-effect": "None",
        "realized-today-date": "2023-05-24",
        "created-at": "2023-05-24T17:17:57.615+00:00",
        "updated-at": "2023-05-24T17:17:58.632+00:00"
      }
    ]
  },
  "context": "/accounts/5YZ55555/positions"
}`

const balanceSnapshotsResp = `{
  "data": {
    "items": [
      {
        "account-number": "5YZ55555",
        "cash-balance": "51600.762",
        "long-equity-value": "281983.965",
        "short-equity-value": "0.0",
        "long-derivative-value": "0.0",
        "short-derivative-value": "82680.0",
        "long-futures-value": "0.0",
        "short-futures-value": "0.0",
        "long-futures-derivative-value": "0.0",
        "short-futures-derivative-value": "0.0",
        "long-margineable-value": "0.0",
        "short-margineable-value": "0.0",
        "margin-equity": "452284.727",
        "equity-buying-power": "20078.762",
        "derivative-buying-power": "20078.762",
        "day-trading-buying-power": "0.0",
        "futures-margin-requirement": "0.0",
        "available-trading-funds": "0.0",
        "maintenance-requirement": "432279.047",
        "maintenance-call-value": "0.0",
        "reg-t-call-value": "0.0",
        "day-trading-call-value": "0.0",
        "day-equity-call-value": "0.0",
        "net-liquidating-value": "4544.727",
        "cash-available-to-withdraw": "20078.76",
        "day-trade-excess": "20078.76",
        "pending-cash": "0.0",
        "pending-cash-effect": "None",
        "long-cryptocurrency-value": "0.0",
        "short-cryptocurrency-value": "0.0",
        "cryptocurrency-margin-requirement": "0.0",
        "unsettled-cryptocurrency-fiat-amount": "0.0",
        "unsettled-cryptocurrency-fiat-effect": "None",
        "closed-loop-available-balance": "20078.76",
        "equity-offering-margin-requirement": "0.0",
        "long-bond-value": "0.0",
        "bond-margin-requirement": "0.0",
        "snapshot-date": "2023-06-08",
        "reg-t-margin-requirement": "432279.0465",
        "futures-overnight-margin-requirement": "0.0",
        "futures-intraday-margin-requirement": "0.0",
        "maintenance-excess": "20078.762",
        "pending-margin-interest": "0.0",
        "effective-cryptocurrency-buying-power": "20078.76",
        "updated-at": "2023-06-08T18:37:39.568+00:00"
      }
    ]
  },
  "context": "/accounts/5YZ55555/balance-snapshots"
}`

const netLiqHistoryResp = `{
  "data": {
    "items": [
      {
        "open": "4498.667",
        "high": "4498.667",
        "low": "4498.667",
        "close": "4498.667",
        "pending-cash-open": "0.0",
        "pending-cash-high": "0.0",
        "pending-cash-low": "0.0",
        "pending-cash-close": "0.0",
        "total-open": "4498.667",
        "total-high": "4498.667",
        "total-low": "4498.667",
        "total-close": "4498.667",
        "time": "2023-06-08 13:30:00+00"
      },
      {
        "open": "4498.712",
        "high": "4498.712",
        "low": "4498.712",
        "close": "4498.712",
        "pending-cash-open": "0.0",
        "pending-cash-high": "0.0",
        "pending-cash-low": "0.0",
        "pending-cash-close": "0.0",
        "total-open": "4498.712",
        "total-high": "4498.712",
        "total-low": "4498.712",
        "total-close": "4498.712",
        "time": "2023-06-08 13:35:00+00"
      },
      {
        "open": "4507.5383",
        "high": "4507.5383",
        "low": "4507.5383",
        "close": "4507.5383",
        "pending-cash-open": "0.0",
        "pending-cash-high": "0.0",
        "pending-cash-low": "0.0",
        "pending-cash-close": "0.0",
        "total-open": "4507.5383",
        "total-high": "4507.5383",
        "total-low": "4507.5383",
        "total-close": "4507.5383",
        "time": "2023-06-08 13:40:00+00"
      }
    ]
  }
}`

const positionLimitResp = `{
  "data": {
      "account-number": "5YZ55555",
      "equity-order-size": 50000,
      "equity-option-order-size": 10000,
      "future-order-size": 2500,
      "future-option-order-size": 2500,
      "underlying-opening-order-limit": 500,
      "equity-position-size": 500000,
      "equity-option-position-size": 20000,
      "future-position-size": 5000,
      "future-option-position-size": 5000
  },
  "context": "/accounts/5YZ55555/position-limit"
}`
