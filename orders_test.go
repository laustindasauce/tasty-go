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

func TestSubmitMarketOrderDryRun(t *testing.T) {
	setup()
	defer teardown()

	accountNumber := "5YZ55555"
	symbol := "AAPL"
	quantity := 1
	action := constants.BTO

	mux.HandleFunc(fmt.Sprintf("/accounts/%s/orders/dry-run", accountNumber), func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, orderDryRunResp)
	})

	order := models.NewOrder{
		TimeInForce: constants.Day,
		OrderType:   constants.Market,
		Legs: []models.NewOrderLeg{
			{
				InstrumentType: constants.Equity,
				Symbol:         symbol,
				Quantity:       quantity,
				Action:         action,
			},
		},
	}

	resp, orderErr, err := client.SubmitOrderDryRun(accountNumber, order)
	require.Nil(t, err)
	require.Nil(t, orderErr)

	o := resp.Order

	require.Equal(t, accountNumber, o.AccountNumber)
	require.Equal(t, constants.Day, o.TimeInForce)
	require.Equal(t, constants.Market, o.OrderType)
	require.Equal(t, 1, o.Size)
	require.Equal(t, symbol, o.UnderlyingSymbol)
	require.Equal(t, constants.Equity, o.UnderlyingInstrumentType)
	require.Equal(t, constants.Contingent, o.Status)
	require.Equal(t, "Pending Condition", o.ContingentStatus)
	require.True(t, o.Cancellable)
	require.True(t, o.Editable)
	require.False(t, o.Edited)
	require.Zero(t, o.UpdatedAt)

	ol := o.Legs[0]

	require.Equal(t, constants.Equity, ol.InstrumentType)
	require.Equal(t, symbol, ol.Symbol)
	require.Equal(t, quantity, ol.Quantity)
	require.Equal(t, quantity, ol.RemainingQuantity)
	require.Equal(t, action, ol.Action)
	require.Empty(t, ol.Fills)

	require.Empty(t, resp.Warnings)

	bpe := resp.BuyingPowerEffect

	require.Equal(t, models.StringToFloat32(183.08), bpe.ChangeInMarginRequirement)
	require.Equal(t, constants.Debit, bpe.ChangeInMarginRequirementEffect)
	require.Equal(t, models.StringToFloat32(183.081), bpe.ChangeInBuyingPower)
	require.Equal(t, constants.Debit, bpe.ChangeInBuyingPowerEffect)
	require.Equal(t, models.StringToFloat32(241.62), bpe.CurrentBuyingPower)
	require.Equal(t, constants.Credit, bpe.CurrentBuyingPowerEffect)
	require.Equal(t, models.StringToFloat32(58.539), bpe.NewBuyingPower)
	require.Equal(t, constants.Credit, bpe.NewBuyingPowerEffect)
	require.Equal(t, models.StringToFloat32(183.08), bpe.IsolatedOrderMarginRequirement)
	require.Equal(t, constants.Debit, bpe.IsolatedOrderMarginRequirementEffect)
	require.False(t, bpe.IsSpread)
	require.Equal(t, models.StringToFloat32(183.081), bpe.Impact)
	require.Equal(t, constants.Debit, bpe.Effect)

	fee := resp.FeeCalculation

	require.Zero(t, fee.RegulatoryFees)
	require.Equal(t, constants.None, fee.RegulatoryFeesEffect)
	require.Equal(t, models.StringToFloat32(0.001), fee.ClearingFees)
	require.Equal(t, constants.Debit, fee.ClearingFeesEffect)
	require.Equal(t, models.StringToFloat32(0.0), fee.Commission)
	require.Equal(t, constants.None, fee.CommissionEffect)
	require.Equal(t, models.StringToFloat32(0.0), fee.ProprietaryIndexOptionFees)
	require.Equal(t, constants.None, fee.ProprietaryIndexOptionFeesEffect)
	require.Equal(t, models.StringToFloat32(0.001), fee.TotalFees)
	require.Equal(t, constants.Debit, fee.TotalFeesEffect)
}

func TestSubmitGTCOrderDryRun(t *testing.T) {
	setup()
	defer teardown()

	accountNumber := "5YZ55555"
	symbol := "GOOGL"
	quantity := 1
	action := constants.STC
	price := float32(124.55)

	mux.HandleFunc(fmt.Sprintf("/accounts/%s/orders/dry-run", accountNumber), func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, orderDryRunGTCResp)
	})

	order := models.NewOrder{
		TimeInForce: constants.GTC,
		OrderType:   constants.Limit,
		Price:       price,
		PriceEffect: constants.Credit,
		Legs: []models.NewOrderLeg{
			{
				InstrumentType: constants.Equity,
				Symbol:         symbol,
				Quantity:       quantity,
				Action:         action,
			},
		},
	}

	resp, orderErr, err := client.SubmitOrderDryRun(accountNumber, order)
	require.Nil(t, err)
	require.Nil(t, orderErr)

	o := resp.Order

	require.Equal(t, accountNumber, o.AccountNumber)
	require.Equal(t, constants.GTC, o.TimeInForce)
	require.Equal(t, constants.Limit, o.OrderType)
	require.Equal(t, 1, o.Size)
	require.Equal(t, symbol, o.UnderlyingSymbol)
	require.Equal(t, constants.Equity, o.UnderlyingInstrumentType)
	require.Equal(t, models.StringToFloat32(price), o.Price)
	require.Equal(t, constants.Credit, o.PriceEffect)
	require.Equal(t, constants.Contingent, o.Status)
	require.Equal(t, "Pending Condition", o.ContingentStatus)
	require.True(t, o.Cancellable)
	require.True(t, o.Editable)
	require.False(t, o.Edited)
	require.Zero(t, o.UpdatedAt)

	ol := o.Legs[0]

	require.Equal(t, constants.Equity, ol.InstrumentType)
	require.Equal(t, symbol, ol.Symbol)
	require.Equal(t, quantity, ol.Quantity)
	require.Equal(t, quantity, ol.RemainingQuantity)
	require.Equal(t, action, ol.Action)
	require.Empty(t, ol.Fills)

	require.Empty(t, resp.Warnings)

	bpe := resp.BuyingPowerEffect

	require.Equal(t, models.StringToFloat32(123.965), bpe.ChangeInMarginRequirement)
	require.Equal(t, constants.Credit, bpe.ChangeInMarginRequirementEffect)
	require.Equal(t, models.StringToFloat32(124.538855), bpe.ChangeInBuyingPower)
	require.Equal(t, constants.Credit, bpe.ChangeInBuyingPowerEffect)
	require.Equal(t, models.StringToFloat32(241.62), bpe.CurrentBuyingPower)
	require.Equal(t, constants.Credit, bpe.CurrentBuyingPowerEffect)
	require.Equal(t, models.StringToFloat32(366.158855), bpe.NewBuyingPower)
	require.Equal(t, constants.Credit, bpe.NewBuyingPowerEffect)
	require.Equal(t, models.StringToFloat32(123.965), bpe.IsolatedOrderMarginRequirement)
	require.Equal(t, constants.Debit, bpe.IsolatedOrderMarginRequirementEffect)
	require.False(t, bpe.IsSpread)
	require.Equal(t, models.StringToFloat32(124.538855), bpe.Impact)
	require.Equal(t, constants.Credit, bpe.Effect)
}

func TestSubmitErrorOrderDryRun(t *testing.T) {
	setup()
	defer teardown()

	accountNumber := "5YZ55555"
	symbol := "AAPL"
	quantity := 10
	action := constants.BTO

	mux.HandleFunc(fmt.Sprintf("/accounts/%s/orders/dry-run", accountNumber), func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, orderErrorDryRunResp)
	})

	order := models.NewOrder{
		TimeInForce: constants.Day,
		OrderType:   constants.Market,
		Legs: []models.NewOrderLeg{
			{
				InstrumentType: constants.Equity,
				Symbol:         symbol,
				Quantity:       quantity,
				Action:         action,
			},
		},
	}

	resp, orderErr, err := client.SubmitOrderDryRun(accountNumber, order)
	require.Nil(t, err)
	require.NotNil(t, orderErr)

	bpe := resp.BuyingPowerEffect

	require.Equal(t, models.StringToFloat32(1828.8), bpe.ChangeInMarginRequirement)
	require.Equal(t, constants.Debit, bpe.ChangeInMarginRequirementEffect)
	require.Equal(t, models.StringToFloat32(1829.008), bpe.ChangeInBuyingPower)
	require.Equal(t, constants.Debit, bpe.ChangeInBuyingPowerEffect)
	require.Equal(t, models.StringToFloat32(241.62), bpe.CurrentBuyingPower)
	require.Equal(t, constants.Credit, bpe.CurrentBuyingPowerEffect)
	require.Equal(t, models.StringToFloat32(1587.388), bpe.NewBuyingPower)
	require.Equal(t, constants.Debit, bpe.NewBuyingPowerEffect)
	require.Equal(t, models.StringToFloat32(1828.8), bpe.IsolatedOrderMarginRequirement)
	require.Equal(t, constants.Debit, bpe.IsolatedOrderMarginRequirementEffect)
	require.False(t, bpe.IsSpread)
	require.Equal(t, models.StringToFloat32(1829.008), bpe.Impact)
	require.Equal(t, constants.Debit, bpe.Effect)

	require.Equal(t, "preflight_check_failure", orderErr.Code)
	require.Equal(t, "One or more preflight checks failed", orderErr.Message)
	require.Equal(t, "margin_check_failed", orderErr.Errors[0].Code)
	require.Equal(t, "Account does not have sufficient buying power available for this order.", orderErr.Errors[0].Message)
}

func TestSubmitOrder(t *testing.T) {
	setup()
	defer teardown()

	accountNumber := "5YZ55555"

	mux.HandleFunc(fmt.Sprintf("/accounts/%s/orders", accountNumber), func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, orderResp)
	})

	symbol := "RIVN"
	quantity := 1
	action1 := constants.BTC

	symbol1 := utils.EquityOptionsSymbology{
		Symbol:     symbol,
		OptionType: constants.Call,
		Strike:     15,
		Expiration: time.Date(2023, 6, 23, 0, 0, 0, 0, time.Local),
	}

	order := models.NewOrder{
		TimeInForce: constants.GTC,
		OrderType:   constants.Limit,
		PriceEffect: constants.Debit,
		Price:       0.04,
		Legs: []models.NewOrderLeg{
			{
				InstrumentType: constants.EquityOption,
				Symbol:         symbol1.Build(),
				Quantity:       quantity,
				Action:         action1,
			},
		},
		Rules: models.NewOrderRules{Conditions: []models.NewOrderCondition{
			{
				Action:         constants.Route,
				Symbol:         symbol,
				InstrumentType: "Equity",
				Indicator:      constants.Last,
				Comparator:     constants.LTE,
				Threshold:      0.01,
			},
		}},
	}

	resp, orderErr, err := client.SubmitOrder(accountNumber, order)
	require.Nil(t, err)
	require.Nil(t, orderErr)

	o := resp.Order

	require.Equal(t, accountNumber, o.AccountNumber)
	require.Equal(t, constants.GTC, o.TimeInForce)
	require.Equal(t, constants.Limit, o.OrderType)
	require.Equal(t, 1, o.Size)
	require.Equal(t, symbol, o.UnderlyingSymbol)
	require.Equal(t, constants.Equity, o.UnderlyingInstrumentType)
	require.Equal(t, models.StringToFloat32(0.04), o.Price)
	require.Equal(t, constants.Debit, o.PriceEffect)
	require.Equal(t, constants.Contingent, o.Status)
	require.Equal(t, "Pending Condition", o.ContingentStatus)
	require.True(t, o.Cancellable)
	require.True(t, o.Editable)
	require.False(t, o.Edited)
	require.Equal(t, "2023-06-13T20:25:43.22Z", o.ReceivedAt.Format(time.RFC3339Nano))
	require.Equal(t, 1686687943220, o.UpdatedAt)

	ol := o.Legs[0]

	require.Equal(t, constants.EquityOption, ol.InstrumentType)
	require.Equal(t, symbol1.Build(), ol.Symbol)
	require.Equal(t, quantity, ol.Quantity)
	require.Equal(t, quantity, ol.RemainingQuantity)
	require.Equal(t, action1, ol.Action)
	require.Empty(t, ol.Fills)

	oc := o.Rules.Conditions[0]

	require.Equal(t, 207561, oc.ID)
	require.Equal(t, constants.Route, oc.Action)
	require.Equal(t, symbol, oc.Symbol)
	require.Equal(t, constants.Equity, oc.InstrumentType)
	require.Equal(t, constants.Last, oc.Indicator)
	require.Equal(t, constants.LTE, oc.Comparator)
	require.Equal(t, models.StringToFloat32(0.01), oc.Threshold)
	require.False(t, oc.IsThresholdBasedOnNotional)

	pc := oc.PriceComponents[0]

	require.Equal(t, symbol, pc.Symbol)
	require.Equal(t, constants.Equity, pc.InstrumentType)
	require.Equal(t, quantity, pc.Quantity)
	require.Equal(t, constants.Long, pc.QuantityDirection)

	warn := resp.Warnings[0]

	require.Equal(t, "tif_next_valid_sesssion", warn.Code)
	require.Equal(t, "Your order will begin working during next valid session.", warn.Message)

	bpe := resp.BuyingPowerEffect

	require.Equal(t, models.StringToFloat32(0.0), bpe.ChangeInMarginRequirement)
	require.Equal(t, constants.None, bpe.ChangeInMarginRequirementEffect)
	require.Equal(t, models.StringToFloat32(4.13), bpe.ChangeInBuyingPower)
	require.Equal(t, constants.Debit, bpe.ChangeInBuyingPowerEffect)
	require.Equal(t, models.StringToFloat32(241.62), bpe.CurrentBuyingPower)
	require.Equal(t, constants.Credit, bpe.CurrentBuyingPowerEffect)
	require.Equal(t, models.StringToFloat32(237.49), bpe.NewBuyingPower)
	require.Equal(t, constants.Credit, bpe.NewBuyingPowerEffect)
	require.Equal(t, models.StringToFloat32(0.0), bpe.IsolatedOrderMarginRequirement)
	require.Equal(t, constants.None, bpe.IsolatedOrderMarginRequirementEffect)
	require.False(t, bpe.IsSpread)
	require.Equal(t, models.StringToFloat32(4.13), bpe.Impact)
	require.Equal(t, constants.Debit, bpe.Effect)

	fee := resp.FeeCalculation

	require.Equal(t, models.StringToFloat32(.03), fee.RegulatoryFees)
	require.Equal(t, constants.Debit, fee.RegulatoryFeesEffect)
	require.Equal(t, models.StringToFloat32(0.1), fee.ClearingFees)
	require.Equal(t, constants.Debit, fee.ClearingFeesEffect)
	require.Equal(t, models.StringToFloat32(0.0), fee.Commission)
	require.Equal(t, constants.None, fee.CommissionEffect)
	require.Equal(t, models.StringToFloat32(0.0), fee.ProprietaryIndexOptionFees)
	require.Equal(t, constants.None, fee.ProprietaryIndexOptionFeesEffect)
	require.Equal(t, models.StringToFloat32(0.13), fee.TotalFees)
	require.Equal(t, constants.Debit, fee.TotalFeesEffect)
}

func TestGetAccountLiveOrders(t *testing.T) {
	setup()
	defer teardown()

	accountNumber := "5YZ55555"

	mux.HandleFunc(fmt.Sprintf("/accounts/%s/orders/live", accountNumber), func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, liveOrdersResp)
	})

	resp, err := client.GetAccountLiveOrders(accountNumber)
	require.Nil(t, err)

	o := resp[0]

	require.Equal(t, accountNumber, o.AccountNumber)
	require.Equal(t, constants.GTC, o.TimeInForce)
	require.Equal(t, constants.Limit, o.OrderType)
	require.Equal(t, 1, o.Size)
	require.Equal(t, "RIVN", o.UnderlyingSymbol)
	require.Equal(t, constants.Equity, o.UnderlyingInstrumentType)
	require.Equal(t, models.StringToFloat32(0.04), o.Price)
	require.Equal(t, constants.Debit, o.PriceEffect)
	require.Equal(t, constants.Contingent, o.Status)
	require.Equal(t, "Pending Condition", o.ContingentStatus)
	require.True(t, o.Cancellable)
	require.True(t, o.Editable)
	require.False(t, o.Edited)
	require.Equal(t, "2023-06-13T20:25:43.22Z", o.ReceivedAt.Format(time.RFC3339Nano))
	require.Equal(t, 1686687943220, o.UpdatedAt)

	ol := o.Legs[0]

	require.Equal(t, constants.EquityOption, ol.InstrumentType)
	require.Equal(t, "RIVN  230623C00015000", ol.Symbol)
	require.Equal(t, 1, ol.Quantity)
	require.Equal(t, 1, ol.RemainingQuantity)
	require.Equal(t, constants.BTC, ol.Action)
	require.Empty(t, ol.Fills)

	oc := o.Rules.Conditions[0]

	require.Equal(t, 207561, oc.ID)
	require.Equal(t, constants.Route, oc.Action)
	require.Equal(t, "RIVN", oc.Symbol)
	require.Equal(t, constants.Equity, oc.InstrumentType)
	require.Equal(t, constants.Last, oc.Indicator)
	require.Equal(t, constants.LTE, oc.Comparator)
	require.Equal(t, models.StringToFloat32(0.01), oc.Threshold)
	require.False(t, oc.IsThresholdBasedOnNotional)

	pc := oc.PriceComponents[0]

	require.Equal(t, "RIVN", pc.Symbol)
	require.Equal(t, constants.Equity, pc.InstrumentType)
	require.Equal(t, 1, pc.Quantity)
	require.Equal(t, constants.Long, pc.QuantityDirection)
}

func TestSubmitOrderECRDryRun(t *testing.T) {
	setup()
	defer teardown()

	accountNumber := "5WV48989"
	orderID := 68675

	mux.HandleFunc(fmt.Sprintf("/accounts/%s/orders/%d/dry-run", accountNumber, orderID), func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, orderECRDryRunResp)
	})

	orderECR := models.NewOrderECR{
		TimeInForce: constants.Day,
		Price:       185.45,
		OrderType:   constants.Limit,
		PriceEffect: constants.Debit,
	}

	resp, err := client.SubmitOrderECRDryRun(accountNumber, orderID, orderECR)
	require.Nil(t, err)

	o := resp.Order

	require.Equal(t, accountNumber, o.AccountNumber)
	require.Equal(t, constants.Day, o.TimeInForce)
	require.Equal(t, constants.Limit, o.OrderType)
	require.Equal(t, 1, o.Size)
	require.Equal(t, "AAPL", o.UnderlyingSymbol)
	require.Equal(t, constants.Equity, o.UnderlyingInstrumentType)
	require.Equal(t, models.StringToFloat32(185.45), o.Price)
	require.Equal(t, constants.Debit, o.PriceEffect)
	require.Equal(t, constants.Contingent, o.Status)
	require.Equal(t, "Pending Order", o.ContingentStatus)
	require.True(t, o.Cancellable)
	require.True(t, o.Editable)
	require.False(t, o.Edited)
	require.Zero(t, o.UpdatedAt)
}

func TestGetOrder(t *testing.T) {
	setup()
	defer teardown()

	accountNumber := "5WV48989"
	orderID := 68675

	mux.HandleFunc(fmt.Sprintf("/accounts/%s/orders/%d", accountNumber, orderID), func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, getOrderResp)
	})

	o, err := client.GetOrder(accountNumber, orderID)
	require.Nil(t, err)

	require.Equal(t, accountNumber, o.AccountNumber)
	require.Equal(t, constants.Day, o.TimeInForce)
	require.Equal(t, constants.Limit, o.OrderType)
	require.Equal(t, 1, o.Size)
	require.Equal(t, "AAPL", o.UnderlyingSymbol)
	require.Equal(t, constants.Equity, o.UnderlyingInstrumentType)
	require.Equal(t, models.StringToFloat32(124.55), o.Price)
	require.Equal(t, constants.Debit, o.PriceEffect)
	require.Equal(t, constants.Contingent, o.Status)
	require.Equal(t, "Pending Condition", o.ContingentStatus)
	require.True(t, o.Cancellable)
	require.True(t, o.Editable)
	require.False(t, o.Edited)
	require.Equal(t, 1686698526525, o.UpdatedAt)
}

func TestCancelOrder(t *testing.T) {
	setup()
	defer teardown()

	accountNumber := "5WV48989"
	orderID := 68677

	mux.HandleFunc(fmt.Sprintf("/accounts/%s/orders/%d", accountNumber, orderID), func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, cancelledOrderResp)
	})

	o, err := client.CancelOrder(accountNumber, orderID)
	require.Nil(t, err)

	require.Equal(t, accountNumber, o.AccountNumber)
	require.Equal(t, constants.Day, o.TimeInForce)
	require.Equal(t, constants.Limit, o.OrderType)
	require.Equal(t, 1, o.Size)
	require.Equal(t, "AAPL", o.UnderlyingSymbol)
	require.Equal(t, constants.Equity, o.UnderlyingInstrumentType)
	require.Equal(t, models.StringToFloat32(186.99), o.Price)
	require.Equal(t, constants.Debit, o.PriceEffect)
	require.Equal(t, constants.Cancelled, o.Status)
	require.False(t, o.Cancellable)
	require.Equal(t, "2023-06-14T01:02:24.795Z", o.CancelledAt.Format(time.RFC3339Nano))
	require.False(t, o.Editable)
	require.False(t, o.Edited)
	require.Equal(t, "2023-06-14T01:02:24.669Z", o.ReceivedAt.Format(time.RFC3339Nano))
	require.Equal(t, 1686704544800, o.UpdatedAt)
	require.Equal(t, "2023-06-14T01:02:24.794Z", o.TerminalAt.Format(time.RFC3339Nano))
}

func TestReplaceOrder(t *testing.T) {
	setup()
	defer teardown()

	accountNumber := "5WV48989"
	orderID := 68678

	mux.HandleFunc(fmt.Sprintf("/accounts/%s/orders/%d", accountNumber, orderID), func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, replaceOrderResp)
	})

	orderECR := models.NewOrderECR{
		TimeInForce: constants.Day,
		Price:       185.45,
		OrderType:   constants.Limit,
		PriceEffect: constants.Debit,
		ValueEffect: constants.Debit,
	}

	o, err := client.ReplaceOrder(accountNumber, orderID, orderECR)
	require.Nil(t, err)

	require.Equal(t, accountNumber, o.AccountNumber)
	require.Equal(t, constants.Day, o.TimeInForce)
	require.Equal(t, constants.Limit, o.OrderType)
	require.Equal(t, 1, o.Size)
	require.Equal(t, "AAPL", o.UnderlyingSymbol)
	require.Equal(t, constants.Equity, o.UnderlyingInstrumentType)
	require.Equal(t, models.StringToFloat32(185.45), o.Price)
	require.Equal(t, constants.Debit, o.PriceEffect)
	require.Equal(t, constants.Contingent, o.Status)
	require.Equal(t, "Pending Condition", o.ContingentStatus)
	require.True(t, o.Cancellable)
	require.True(t, o.Editable)
	require.False(t, o.Edited)
	require.Equal(t, 1686706739960, o.UpdatedAt)
}

func TestPatchOrder(t *testing.T) {
	setup()
	defer teardown()

	accountNumber := "5WV48989"
	orderID := 68680

	mux.HandleFunc(fmt.Sprintf("/accounts/%s/orders/%d", accountNumber, orderID), func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, patchOrderResp)
	})

	orderECR := models.NewOrderECR{
		TimeInForce: constants.Day,
		Price:       187.45,
		OrderType:   constants.Limit,
		PriceEffect: constants.Debit,
		ValueEffect: constants.Debit,
	}

	o, err := client.PatchOrder(accountNumber, orderID, orderECR)
	require.Nil(t, err)

	require.Equal(t, accountNumber, o.AccountNumber)
	require.Equal(t, constants.Day, o.TimeInForce)
	require.Equal(t, constants.Limit, o.OrderType)
	require.Equal(t, 1, o.Size)
	require.Equal(t, "AAPL", o.UnderlyingSymbol)
	require.Equal(t, constants.Equity, o.UnderlyingInstrumentType)
	require.Equal(t, models.StringToFloat32(187.45), o.Price)
	require.Equal(t, constants.Debit, o.PriceEffect)
	require.Equal(t, constants.Contingent, o.Status)
	require.Equal(t, "Pending Condition", o.ContingentStatus)
	require.True(t, o.Cancellable)
	require.True(t, o.Editable)
	require.False(t, o.Edited)
	require.Equal(t, 1686707204835, o.UpdatedAt)
}

const orderDryRunResp = `{
  "data": {
    "order": {
      "account-number": "5YZ55555",
      "time-in-force": "Day",
      "order-type": "Market",
      "size": 1,
      "underlying-symbol": "AAPL",
      "underlying-instrument-type": "Equity",
      "status": "Contingent",
      "contingent-status": "Pending Condition",
      "cancellable": true,
      "editable": true,
      "edited": false,
      "updated-at": 0,
      "legs": [
        {
          "instrument-type": "Equity",
          "symbol": "AAPL",
          "quantity": 1,
          "remaining-quantity": 1,
          "action": "Buy to Open",
          "fills": []
        }
      ],
      "rules": { "conditions": [] }
    },
    "warnings": [],
    "buying-power-effect": {
      "change-in-margin-requirement": "183.08",
      "change-in-margin-requirement-effect": "Debit",
      "change-in-buying-power": "183.081",
      "change-in-buying-power-effect": "Debit",
      "current-buying-power": "241.62",
      "current-buying-power-effect": "Credit",
      "new-buying-power": "58.539",
      "new-buying-power-effect": "Credit",
      "isolated-order-margin-requirement": "183.08",
      "isolated-order-margin-requirement-effect": "Debit",
      "is-spread": false,
      "impact": "183.081",
      "effect": "Debit"
    },
    "fee-calculation": {
      "regulatory-fees": "0.0",
      "regulatory-fees-effect": "None",
      "clearing-fees": "0.001",
      "clearing-fees-effect": "Debit",
      "commission": "0.0",
      "commission-effect": "None",
      "proprietary-index-option-fees": "0.0",
      "proprietary-index-option-fees-effect": "None",
      "total-fees": "0.001",
      "total-fees-effect": "Debit"
    }
  },
  "context": "/accounts/5YZ55555/orders/dry-run"
}`

const orderDryRunGTCResp = `{
  "data": {
    "order": {
      "account-number": "5YZ55555",
      "time-in-force": "GTC",
      "order-type": "Limit",
      "size": 1,
      "underlying-symbol": "GOOGL",
      "underlying-instrument-type": "Equity",
      "price": "124.55",
      "price-effect": "Credit",
      "status": "Contingent",
      "contingent-status": "Pending Condition",
      "cancellable": true,
      "editable": true,
      "edited": false,
      "updated-at": 0,
      "legs": [
        {
          "instrument-type": "Equity",
          "symbol": "GOOGL",
          "quantity": 1,
          "remaining-quantity": 1,
          "action": "Sell to Close",
          "fills": []
        }
      ],
      "rules": { "conditions": [] }
    },
    "warnings": [],
    "buying-power-effect": {
      "change-in-margin-requirement": "123.965",
      "change-in-margin-requirement-effect": "Credit",
      "change-in-buying-power": "124.538855",
      "change-in-buying-power-effect": "Credit",
      "current-buying-power": "241.62",
      "current-buying-power-effect": "Credit",
      "new-buying-power": "366.158855",
      "new-buying-power-effect": "Credit",
      "isolated-order-margin-requirement": "123.965",
      "isolated-order-margin-requirement-effect": "Debit",
      "is-spread": false,
      "impact": "124.538855",
      "effect": "Credit"
    },
    "fee-calculation": {
      "regulatory-fees": "0.010145",
      "regulatory-fees-effect": "Debit",
      "clearing-fees": "0.001",
      "clearing-fees-effect": "Debit",
      "commission": "0.0",
      "commission-effect": "None",
      "proprietary-index-option-fees": "0.0",
      "proprietary-index-option-fees-effect": "None",
      "total-fees": "0.011145",
      "total-fees-effect": "Debit"
    }
  },
  "context": "/accounts/5YZ55555/orders/dry-run"
}`

const orderErrorDryRunResp = `{
  "data": {
    "buying-power-effect": {
      "change-in-margin-requirement": "1828.8",
      "change-in-margin-requirement-effect": "Debit",
      "change-in-buying-power": "1829.008",
      "change-in-buying-power-effect": "Debit",
      "current-buying-power": "241.62",
      "current-buying-power-effect": "Credit",
      "new-buying-power": "1587.388",
      "new-buying-power-effect": "Debit",
      "isolated-order-margin-requirement": "1828.8",
      "isolated-order-margin-requirement-effect": "Debit",
      "is-spread": false,
      "impact": "1829.008",
      "effect": "Debit"
    }
  },
  "error": {
    "code": "preflight_check_failure",
    "message": "One or more preflight checks failed",
    "errors": [
      {
        "code": "margin_check_failed",
        "message": "Account does not have sufficient buying power available for this order."
      }
    ]
  }
}`

const orderResp = `{
  "data": {
    "order": {
      "id": 272985726,
      "account-number": "5YZ55555",
      "time-in-force": "GTC",
      "order-type": "Limit",
      "size": 1,
      "underlying-symbol": "RIVN",
      "underlying-instrument-type": "Equity",
      "price": "0.04",
      "price-effect": "Debit",
      "status": "Contingent",
      "contingent-status": "Pending Condition",
      "cancellable": true,
      "editable": true,
      "edited": false,
      "received-at": "2023-06-13T20:25:43.220+00:00",
      "updated-at": 1686687943220,
      "legs": [
        {
          "instrument-type": "Equity Option",
          "symbol": "RIVN  230623C00015000",
          "quantity": 1,
          "remaining-quantity": 1,
          "action": "Buy to Close",
          "fills": []
        }
      ],
      "rules": {
        "conditions": [
          {
            "id": 207561,
            "action": "route",
            "symbol": "RIVN",
            "instrument-type": "Equity",
            "indicator": "last",
            "comparator": "lte",
            "threshold": "0.01",
            "is-threshold-based-on-notional": false,
            "price-components": [
              {
                "symbol": "RIVN",
                "instrument-type": "Equity",
                "quantity": 1,
                "quantity-direction": "Long"
              }
            ]
          }
        ]
      }
    },
    "warnings": [
      {
        "code": "tif_next_valid_sesssion",
        "message": "Your order will begin working during next valid session."
      }
    ],
    "buying-power-effect": {
      "change-in-margin-requirement": "0.0",
      "change-in-margin-requirement-effect": "None",
      "change-in-buying-power": "4.13",
      "change-in-buying-power-effect": "Debit",
      "current-buying-power": "241.62",
      "current-buying-power-effect": "Credit",
      "new-buying-power": "237.49",
      "new-buying-power-effect": "Credit",
      "isolated-order-margin-requirement": "0.0",
      "isolated-order-margin-requirement-effect": "None",
      "is-spread": false,
      "impact": "4.13",
      "effect": "Debit"
    },
    "fee-calculation": {
      "regulatory-fees": "0.03",
      "regulatory-fees-effect": "Debit",
      "clearing-fees": "0.1",
      "clearing-fees-effect": "Debit",
      "commission": "0.0",
      "commission-effect": "None",
      "proprietary-index-option-fees": "0.0",
      "proprietary-index-option-fees-effect": "None",
      "total-fees": "0.13",
      "total-fees-effect": "Debit"
    }
  },
  "context": "/accounts/5YZ55555/orders"
}`

const liveOrdersResp = `{
  "data": {
    "items": [
      {
        "id": 272985726,
        "account-number": "5YZ55555",
        "time-in-force": "GTC",
        "order-type": "Limit",
        "size": 1,
        "underlying-symbol": "RIVN",
        "underlying-instrument-type": "Equity",
        "price": "0.04",
        "price-effect": "Debit",
        "status": "Contingent",
        "contingent-status": "Pending Condition",
        "cancellable": true,
        "editable": true,
        "edited": false,
        "received-at": "2023-06-13T20:25:43.220+00:00",
        "updated-at": 1686687943220,
        "legs": [
          {
            "instrument-type": "Equity Option",
            "symbol": "RIVN  230623C00015000",
            "quantity": 1,
            "remaining-quantity": 1,
            "action": "Buy to Close",
            "fills": []
          }
        ],
        "rules": {
          "conditions": [
            {
              "id": 207561,
              "action": "route",
              "symbol": "RIVN",
              "instrument-type": "Equity",
              "indicator": "last",
              "comparator": "lte",
              "threshold": "0.01",
              "is-threshold-based-on-notional": false,
              "price-components": [
                {
                  "symbol": "RIVN",
                  "instrument-type": "Equity",
                  "quantity": 1,
                  "quantity-direction": "Long"
                }
              ]
            }
          ]
        }
      }
    ]
  },
  "context": "/accounts/5YZ55555/orders/live"
}`

const orderECRDryRunResp = `{
  "data": {
    "order": {
      "account-number": "5WV48989",
      "time-in-force": "Day",
      "order-type": "Limit",
      "size": 1,
      "underlying-symbol": "AAPL",
      "underlying-instrument-type": "Equity",
      "price": "185.45",
      "price-effect": "Debit",
      "value-effect": "",
      "status": "Contingent",
      "contingent-status": "Pending Order",
      "cancellable": true,
      "editable": true,
      "edited": false,
      "updated-at": 0,
      "legs": [
        {
          "instrument-type": "Equity",
          "symbol": "AAPL",
          "quantity": 1,
          "remaining-quantity": 1,
          "action": "Buy to Open",
          "fills": []
        }
      ],
      "rules": {
        "conditions": [
          {
            "action": "route",
            "symbol": "AAPL",
            "instrument-type": "Equity",
            "indicator": "last",
            "comparator": "lte",
            "threshold": "0.01",
            "is-threshold-based-on-notional": false,
            "price-components": [
              {
                "symbol": "AAPL",
                "instrument-type": "Equity",
                "quantity": 1,
                "quantity-direction": "Long"
              }
            ]
          }
        ]
      }
    },
    "warnings": [
      {
        "code": "tif_next_valid_sesssion",
        "message": "Your order will begin working during next valid session."
      }
    ],
    "buying-power-effect": {
      "change-in-margin-requirement": "92.725",
      "change-in-margin-requirement-effect": "Debit",
      "change-in-buying-power": "92.726",
      "change-in-buying-power-effect": "Debit",
      "current-buying-power": "1000.0",
      "current-buying-power-effect": "Credit",
      "new-buying-power": "907.274",
      "new-buying-power-effect": "Credit",
      "isolated-order-margin-requirement": "92.725",
      "isolated-order-margin-requirement-effect": "Debit",
      "is-spread": false,
      "impact": "92.726",
      "effect": "Debit"
    },
    "fee-calculation": {
      "regulatory-fees": "0.0",
      "regulatory-fees-effect": "None",
      "clearing-fees": "0.001",
      "clearing-fees-effect": "Debit",
      "commission": "0.0",
      "commission-effect": "None",
      "proprietary-index-option-fees": "0.0",
      "proprietary-index-option-fees-effect": "None",
      "total-fees": "0.001",
      "total-fees-effect": "Debit"
    }
  },
  "context": "/accounts/5WV48989/orders/68675/dry-run"
}`

const getOrderResp = `{
  "data": {
    "id": 68675,
    "account-number": "5WV48989",
    "time-in-force": "Day",
    "order-type": "Limit",
    "size": 1,
    "underlying-symbol": "AAPL",
    "underlying-instrument-type": "Equity",
    "price": "124.55",
    "price-effect": "Debit",
    "status": "Contingent",
    "contingent-status": "Pending Condition",
    "cancellable": true,
    "editable": true,
    "edited": false,
    "received-at": "2023-06-13T23:22:06.525+00:00",
    "updated-at": 1686698526525,
    "legs": [
      {
        "instrument-type": "Equity",
        "symbol": "AAPL",
        "quantity": 1,
        "remaining-quantity": 1,
        "action": "Buy to Open",
        "fills": []
      }
    ],
    "rules": {
      "conditions": [
        {
          "id": 281,
          "action": "route",
          "symbol": "AAPL",
          "instrument-type": "Equity",
          "indicator": "last",
          "comparator": "lte",
          "threshold": "0.01",
          "is-threshold-based-on-notional": false,
          "price-components": [
            {
              "symbol": "AAPL",
              "instrument-type": "Equity",
              "quantity": 1,
              "quantity-direction": "Long"
            }
          ]
        }
      ]
    }
  },
  "context": "/accounts/5WV48989/orders/68675"
}`

const cancelledOrderResp = `{
  "data": {
    "id": 68677,
    "account-number": "5WV48989",
    "time-in-force": "Day",
    "order-type": "Limit",
    "size": 1,
    "underlying-symbol": "AAPL",
    "underlying-instrument-type": "Equity",
    "price": "186.99",
    "price-effect": "Debit",
    "status": "Cancelled",
    "cancellable": false,
    "cancelled-at": "2023-06-14T01:02:24.795+00:00",
    "editable": false,
    "edited": false,
    "received-at": "2023-06-14T01:02:24.669+00:00",
    "updated-at": 1686704544800,
    "terminal-at": "2023-06-14T01:02:24.794+00:00",
    "legs": [
      {
        "instrument-type": "Equity",
        "symbol": "AAPL",
        "quantity": 1,
        "remaining-quantity": 1,
        "action": "Buy to Open",
        "fills": []
      }
    ],
    "rules": {
      "conditions": [
        {
          "id": 283,
          "action": "route",
          "symbol": "AAPL",
          "instrument-type": "Equity",
          "indicator": "last",
          "comparator": "lte",
          "threshold": "0.01",
          "is-threshold-based-on-notional": false,
          "price-components": [
            {
              "symbol": "AAPL",
              "instrument-type": "Equity",
              "quantity": 1,
              "quantity-direction": "Long"
            }
          ]
        }
      ]
    }
  },
  "context": "/accounts/5WV48989/orders/68677"
}`

const replaceOrderResp = `{
  "data": {
    "id": 68680,
    "account-number": "5WV48989",
    "time-in-force": "Day",
    "order-type": "Limit",
    "size": 1,
    "underlying-symbol": "AAPL",
    "underlying-instrument-type": "Equity",
    "price": "185.45",
    "price-effect": "Debit",
    "value-effect": "Debit",
    "status": "Contingent",
    "contingent-status": "Pending Condition",
    "cancellable": true,
    "editable": true,
    "edited": false,
    "received-at": "2023-06-14T01:38:59.936+00:00",
    "updated-at": 1686706739960,
    "legs": [
      {
        "instrument-type": "Equity",
        "symbol": "AAPL",
        "quantity": 1,
        "remaining-quantity": 1,
        "action": "Buy to Open",
        "fills": []
      }
    ],
    "rules": {
      "conditions": [
        {
          "id": 286,
          "action": "route",
          "symbol": "AAPL",
          "instrument-type": "Equity",
          "indicator": "last",
          "comparator": "lte",
          "threshold": "0.01",
          "is-threshold-based-on-notional": false,
          "price-components": [
            {
              "symbol": "AAPL",
              "instrument-type": "Equity",
              "quantity": 1,
              "quantity-direction": "Long"
            }
          ]
        }
      ]
    }
  },
  "context": "/accounts/5WV48989/orders/68678"
}`

const patchOrderResp = `{
  "data": {
    "id": 68681,
    "account-number": "5WV48989",
    "time-in-force": "Day",
    "order-type": "Limit",
    "size": 1,
    "underlying-symbol": "AAPL",
    "underlying-instrument-type": "Equity",
    "price": "187.45",
    "price-effect": "Debit",
    "value-effect": "Debit",
    "status": "Contingent",
    "contingent-status": "Pending Condition",
    "cancellable": true,
    "editable": true,
    "edited": false,
    "received-at": "2023-06-14T01:46:44.803+00:00",
    "updated-at": 1686707204835,
    "legs": [
      {
        "instrument-type": "Equity",
        "symbol": "AAPL",
        "quantity": 1,
        "remaining-quantity": 1,
        "action": "Buy to Open",
        "fills": []
      }
    ],
    "rules": {
      "conditions": [
        {
          "id": 287,
          "action": "route",
          "symbol": "AAPL",
          "instrument-type": "Equity",
          "indicator": "last",
          "comparator": "lte",
          "threshold": "0.01",
          "is-threshold-based-on-notional": false,
          "price-components": [
            {
              "symbol": "AAPL",
              "instrument-type": "Equity",
              "quantity": 1,
              "quantity-direction": "Long"
            }
          ]
        }
      ]
    }
  },
  "context": "/accounts/5WV48989/orders/68680"
}`
