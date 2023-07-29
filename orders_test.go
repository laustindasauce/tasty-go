package tasty //nolint:testpackage // testing private field

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func TestSubmitMarketOrderDryRun(t *testing.T) {
	setup()
	defer teardown()

	accountNumber := "5YZ55555"
	symbol := "AAPL"
	quantity := float32(1)
	action := BTO

	mux.HandleFunc(fmt.Sprintf("/accounts/%s/orders/dry-run", accountNumber), func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, orderDryRunResp)
	})

	order := NewOrder{
		TimeInForce: Day,
		OrderType:   Market,
		Legs: []NewOrderLeg{
			{
				InstrumentType: EquityIT,
				Symbol:         symbol,
				Quantity:       quantity,
				Action:         action,
			},
		},
	}

	resp, orderErr, httpResp, err := client.SubmitOrderDryRun(accountNumber, order)
	require.Nil(t, err)
	require.NotNil(t, httpResp)
	require.Nil(t, orderErr)

	o := resp.Order

	require.Equal(t, accountNumber, o.AccountNumber)
	require.Equal(t, Day, o.TimeInForce)
	require.Equal(t, Market, o.OrderType)
	require.Equal(t, 1, o.Size)
	require.Equal(t, symbol, o.UnderlyingSymbol)
	require.Equal(t, EquityIT, o.UnderlyingInstrumentType)
	require.Equal(t, Contingent, o.Status)
	require.Equal(t, "Pending Condition", o.ContingentStatus)
	require.True(t, o.Cancellable)
	require.True(t, o.Editable)
	require.False(t, o.Edited)
	require.Zero(t, o.UpdatedAt, "hello")

	ol := o.Legs[0]

	require.Equal(t, EquityIT, ol.InstrumentType)
	require.Equal(t, symbol, ol.Symbol)
	require.Equal(t, quantity, ol.Quantity)
	require.Equal(t, quantity, ol.RemainingQuantity)
	require.Equal(t, action, ol.Action)
	require.Empty(t, ol.Fills)

	require.Empty(t, resp.Warnings)

	bpe := resp.BuyingPowerEffect

	require.Equal(t, decimal.NewFromFloat(183.08), bpe.ChangeInMarginRequirement)
	require.Equal(t, Debit, bpe.ChangeInMarginRequirementEffect)
	require.Equal(t, decimal.NewFromFloat(183.081), bpe.ChangeInBuyingPower)
	require.Equal(t, Debit, bpe.ChangeInBuyingPowerEffect)
	require.Equal(t, decimal.NewFromFloat(241.62), bpe.CurrentBuyingPower)
	require.Equal(t, Credit, bpe.CurrentBuyingPowerEffect)
	require.Equal(t, decimal.NewFromFloat(58.539), bpe.NewBuyingPower)
	require.Equal(t, Credit, bpe.NewBuyingPowerEffect)
	require.Equal(t, decimal.NewFromFloat(183.08), bpe.IsolatedOrderMarginRequirement)
	require.Equal(t, Debit, bpe.IsolatedOrderMarginRequirementEffect)
	require.False(t, bpe.IsSpread)
	require.Equal(t, decimal.NewFromFloat(183.081), bpe.Impact)
	require.Equal(t, Debit, bpe.Effect)

	fee := resp.FeeCalculation

	require.True(t, fee.RegulatoryFees.Equal(decimal.Zero), "regulatory fees")
	require.Equal(t, None, fee.RegulatoryFeesEffect)
	require.Equal(t, decimal.NewFromFloat(0.001), fee.ClearingFees)
	require.Equal(t, Debit, fee.ClearingFeesEffect)
	require.True(t, fee.Commission.Equal(decimal.Zero))
	require.Equal(t, None, fee.CommissionEffect)
	require.True(t, fee.ProprietaryIndexOptionFees.Equal(decimal.Zero))
	require.Equal(t, None, fee.ProprietaryIndexOptionFeesEffect)
	require.Equal(t, decimal.NewFromFloat(0.001), fee.TotalFees)
	require.Equal(t, Debit, fee.TotalFeesEffect)
}

func TestSubmitMarketOrderDryRunError(t *testing.T) {
	setup()
	defer teardown()

	accountNumber := "5YZ55555"
	symbol := "AAPL"
	quantity := float32(1)
	action := BTO

	mux.HandleFunc(fmt.Sprintf("/accounts/%s/orders/dry-run", accountNumber), func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(401)
		fmt.Fprint(writer, tastyUnauthorizedError)
	})

	order := NewOrder{
		TimeInForce: Day,
		OrderType:   Market,
		Legs: []NewOrderLeg{
			{
				InstrumentType: EquityIT,
				Symbol:         symbol,
				Quantity:       quantity,
				Action:         action,
			},
		},
	}

	_, _, httpResp, err := client.SubmitOrderDryRun(accountNumber, order)
	expectedUnauthorized(t, err)
	require.NotNil(t, httpResp)
}

func TestReconfirmOrderError(t *testing.T) {
	setup()
	defer teardown()

	accountNumber := "5YZ55555"
	orderID := 272985726

	mux.HandleFunc(fmt.Sprintf("/accounts/%s/orders/%d/reconfirm", accountNumber, orderID), func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(422)
		fmt.Fprint(writer, reconfirmResp)
	})

	_, httpResp, err := client.ReconfirmOrder(accountNumber, orderID)
	require.NotNil(t, err)
	require.NotNil(t, httpResp)

	require.Equal(t, "\nError in request 422;\nCode: cannot_reconfirm_order\nMessage: The order could not be reconfirmed.", err.Error())
}

func TestSubmitGTCOrderDryRun(t *testing.T) {
	setup()
	defer teardown()

	accountNumber := "5YZ55555"
	symbol := "GOOGL"
	quantity := float32(1)
	action := STC
	price := float32(124.55)

	mux.HandleFunc(fmt.Sprintf("/accounts/%s/orders/dry-run", accountNumber), func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, orderDryRunGTCResp)
	})

	order := NewOrder{
		TimeInForce: GTC,
		OrderType:   Limit,
		Price:       price,
		PriceEffect: Credit,
		Legs: []NewOrderLeg{
			{
				InstrumentType: EquityIT,
				Symbol:         symbol,
				Quantity:       quantity,
				Action:         action,
			},
		},
	}

	resp, orderErr, httpResp, err := client.SubmitOrderDryRun(accountNumber, order)
	require.Nil(t, err)
	require.NotNil(t, httpResp)
	require.Nil(t, orderErr)

	o := resp.Order

	require.Equal(t, accountNumber, o.AccountNumber)
	require.Equal(t, GTC, o.TimeInForce)
	require.Equal(t, Limit, o.OrderType)
	require.Equal(t, 1, o.Size)
	require.Equal(t, symbol, o.UnderlyingSymbol)
	require.Equal(t, EquityIT, o.UnderlyingInstrumentType)
	require.Equal(t, decimal.NewFromFloat32(price), o.Price)
	require.Equal(t, Credit, o.PriceEffect)
	require.Equal(t, Contingent, o.Status)
	require.Equal(t, "Pending Condition", o.ContingentStatus)
	require.True(t, o.Cancellable)
	require.True(t, o.Editable)
	require.False(t, o.Edited)
	require.Zero(t, o.UpdatedAt)

	ol := o.Legs[0]

	require.Equal(t, EquityIT, ol.InstrumentType)
	require.Equal(t, symbol, ol.Symbol)
	require.Equal(t, quantity, ol.Quantity)
	require.Equal(t, quantity, ol.RemainingQuantity)
	require.Equal(t, action, ol.Action)
	require.Empty(t, ol.Fills)

	require.Empty(t, resp.Warnings)

	bpe := resp.BuyingPowerEffect

	require.Equal(t, decimal.NewFromFloat(123.965), bpe.ChangeInMarginRequirement)
	require.Equal(t, Credit, bpe.ChangeInMarginRequirementEffect)
	require.Equal(t, decimal.NewFromFloat(124.538855), bpe.ChangeInBuyingPower)
	require.Equal(t, Credit, bpe.ChangeInBuyingPowerEffect)
	require.Equal(t, decimal.NewFromFloat(241.62), bpe.CurrentBuyingPower)
	require.Equal(t, Credit, bpe.CurrentBuyingPowerEffect)
	require.Equal(t, decimal.NewFromFloat(366.158855), bpe.NewBuyingPower)
	require.Equal(t, Credit, bpe.NewBuyingPowerEffect)
	require.Equal(t, decimal.NewFromFloat(123.965), bpe.IsolatedOrderMarginRequirement)
	require.Equal(t, Debit, bpe.IsolatedOrderMarginRequirementEffect)
	require.False(t, bpe.IsSpread)
	require.Equal(t, decimal.NewFromFloat(124.538855), bpe.Impact)
	require.Equal(t, Credit, bpe.Effect)
}

func TestSubmitGTCOrderDryRunError(t *testing.T) {
	setup()
	defer teardown()

	accountNumber := "5YZ55555"
	symbol := "GOOGL"
	quantity := float32(1)
	action := STC
	price := float32(124.55)

	mux.HandleFunc(fmt.Sprintf("/accounts/%s/orders/dry-run", accountNumber), func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(401)
		fmt.Fprint(writer, tastyUnauthorizedError)
	})

	order := NewOrder{
		TimeInForce: GTC,
		OrderType:   Limit,
		Price:       price,
		PriceEffect: Credit,
		Legs: []NewOrderLeg{
			{
				InstrumentType: EquityIT,
				Symbol:         symbol,
				Quantity:       quantity,
				Action:         action,
			},
		},
	}

	_, _, httpResp, err := client.SubmitOrderDryRun(accountNumber, order)
	expectedUnauthorized(t, err)
	require.NotNil(t, httpResp)
}

func TestSubmitErrorOrderDryRun(t *testing.T) {
	setup()
	defer teardown()

	accountNumber := "5YZ55555"
	symbol := "AAPL"
	quantity := float32(10)
	action := BTO

	mux.HandleFunc(fmt.Sprintf("/accounts/%s/orders/dry-run", accountNumber), func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, orderErrorDryRunResp)
	})

	order := NewOrder{
		TimeInForce: Day,
		OrderType:   Market,
		Legs: []NewOrderLeg{
			{
				InstrumentType: EquityIT,
				Symbol:         symbol,
				Quantity:       quantity,
				Action:         action,
			},
		},
	}

	resp, orderErr, httpResp, err := client.SubmitOrderDryRun(accountNumber, order)
	require.Nil(t, err)
	require.NotNil(t, httpResp)
	require.NotNil(t, orderErr)

	bpe := resp.BuyingPowerEffect

	require.Equal(t, decimal.NewFromFloat(1828.8), bpe.ChangeInMarginRequirement)
	require.Equal(t, Debit, bpe.ChangeInMarginRequirementEffect)
	require.Equal(t, decimal.NewFromFloat(1829.008), bpe.ChangeInBuyingPower)
	require.Equal(t, Debit, bpe.ChangeInBuyingPowerEffect)
	require.Equal(t, decimal.NewFromFloat(241.62), bpe.CurrentBuyingPower)
	require.Equal(t, Credit, bpe.CurrentBuyingPowerEffect)
	require.Equal(t, decimal.NewFromFloat(1587.388), bpe.NewBuyingPower)
	require.Equal(t, Debit, bpe.NewBuyingPowerEffect)
	require.Equal(t, decimal.NewFromFloat(1828.8), bpe.IsolatedOrderMarginRequirement)
	require.Equal(t, Debit, bpe.IsolatedOrderMarginRequirementEffect)
	require.False(t, bpe.IsSpread)
	require.Equal(t, decimal.NewFromFloat(1829.008), bpe.Impact)
	require.Equal(t, Debit, bpe.Effect)

	require.Equal(t, "preflight_check_failure", orderErr.Code)
	require.Equal(t, "One or more preflight checks failed", orderErr.Message)
	require.Equal(t, "margin_check_failed", orderErr.Errors[0].Code)
	require.Equal(t, "Account does not have sufficient buying power available for this order.", orderErr.Errors[0].Message)
}

func TestSubmitErrorOrderDryRunError(t *testing.T) {
	setup()
	defer teardown()

	accountNumber := "5YZ55555"
	symbol := "AAPL"
	quantity := float32(10)
	action := BTO

	mux.HandleFunc(fmt.Sprintf("/accounts/%s/orders/dry-run", accountNumber), func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(401)
		fmt.Fprint(writer, tastyUnauthorizedError)
	})

	order := NewOrder{
		TimeInForce: Day,
		OrderType:   Market,
		Legs: []NewOrderLeg{
			{
				InstrumentType: EquityIT,
				Symbol:         symbol,
				Quantity:       quantity,
				Action:         action,
			},
		},
	}

	_, _, httpResp, err := client.SubmitOrderDryRun(accountNumber, order)
	expectedUnauthorized(t, err)
	require.NotNil(t, httpResp)
}

func TestSubmitOrder(t *testing.T) {
	setup()
	defer teardown()

	accountNumber := "5YZ55555"

	mux.HandleFunc(fmt.Sprintf("/accounts/%s/orders", accountNumber), func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, orderResp)
	})

	symbol := "RIVN"
	quantity := float32(1)
	action1 := BTC

	symbol1 := EquityOptionsSymbology{
		Symbol:     symbol,
		OptionType: Call,
		Strike:     15,
		Expiration: time.Date(2023, 6, 23, 0, 0, 0, 0, time.Local),
	}

	order := NewOrder{
		TimeInForce: GTC,
		OrderType:   Limit,
		PriceEffect: Debit,
		Price:       0.04,
		Legs: []NewOrderLeg{
			{
				InstrumentType: EquityOptionIT,
				Symbol:         symbol1.Build(),
				Quantity:       quantity,
				Action:         action1,
			},
		},
		Rules: NewOrderRules{Conditions: []NewOrderCondition{
			{
				Action:         Route,
				Symbol:         symbol,
				InstrumentType: "Equity",
				Indicator:      Last,
				Comparator:     LTE,
				Threshold:      0.01,
			},
		}},
	}

	resp, orderErr, httpResp, err := client.SubmitOrder(accountNumber, order)
	require.Nil(t, err)
	require.NotNil(t, httpResp)
	require.Nil(t, orderErr)

	o := resp.Order

	require.Equal(t, accountNumber, o.AccountNumber)
	require.Equal(t, GTC, o.TimeInForce)
	require.Equal(t, Limit, o.OrderType)
	require.Equal(t, 1, o.Size)
	require.Equal(t, symbol, o.UnderlyingSymbol)
	require.Equal(t, EquityIT, o.UnderlyingInstrumentType)
	require.Equal(t, decimal.NewFromFloat(0.04), o.Price)
	require.Equal(t, Debit, o.PriceEffect)
	require.Equal(t, Contingent, o.Status)
	require.Equal(t, "Pending Condition", o.ContingentStatus)
	require.True(t, o.Cancellable)
	require.True(t, o.Editable)
	require.False(t, o.Edited)
	require.Equal(t, "2023-06-13T20:25:43.22Z", o.ReceivedAt.Format(time.RFC3339Nano))
	require.Equal(t, 1686687943220, o.UpdatedAt)

	ol := o.Legs[0]

	require.Equal(t, EquityOptionIT, ol.InstrumentType)
	require.Equal(t, symbol1.Build(), ol.Symbol)
	require.Equal(t, quantity, ol.Quantity)
	require.Equal(t, quantity, ol.RemainingQuantity)
	require.Equal(t, action1, ol.Action)
	require.Empty(t, ol.Fills)

	oc := o.Rules.Conditions[0]

	require.Equal(t, 207561, oc.ID)
	require.Equal(t, Route, oc.Action)
	require.Equal(t, symbol, oc.Symbol)
	require.Equal(t, EquityIT, oc.InstrumentType)
	require.Equal(t, Last, oc.Indicator)
	require.Equal(t, LTE, oc.Comparator)
	require.Equal(t, decimal.NewFromFloat(0.01), oc.Threshold)
	require.False(t, oc.IsThresholdBasedOnNotional)

	pc := oc.PriceComponents[0]

	require.Equal(t, symbol, pc.Symbol)
	require.Equal(t, EquityIT, pc.InstrumentType)
	require.Equal(t, quantity, pc.Quantity)
	require.Equal(t, Long, pc.QuantityDirection)

	warn := resp.Warnings[0]

	require.Equal(t, "tif_next_valid_sesssion", warn.Code)
	require.Equal(t, "Your order will begin working during next valid session.", warn.Message)

	bpe := resp.BuyingPowerEffect

	require.True(t, bpe.ChangeInMarginRequirement.Equal(decimal.Zero))
	require.Equal(t, None, bpe.ChangeInMarginRequirementEffect)
	require.Equal(t, decimal.NewFromFloat(4.13), bpe.ChangeInBuyingPower)
	require.Equal(t, Debit, bpe.ChangeInBuyingPowerEffect)
	require.Equal(t, decimal.NewFromFloat(241.62), bpe.CurrentBuyingPower)
	require.Equal(t, Credit, bpe.CurrentBuyingPowerEffect)
	require.Equal(t, decimal.NewFromFloat(237.49), bpe.NewBuyingPower)
	require.Equal(t, Credit, bpe.NewBuyingPowerEffect)
	require.True(t, bpe.IsolatedOrderMarginRequirement.Equal(decimal.Zero))
	require.Equal(t, None, bpe.IsolatedOrderMarginRequirementEffect)
	require.False(t, bpe.IsSpread)
	require.Equal(t, decimal.NewFromFloat(4.13), bpe.Impact)
	require.Equal(t, Debit, bpe.Effect)

	fee := resp.FeeCalculation

	require.Equal(t, decimal.NewFromFloat(.03), fee.RegulatoryFees)
	require.Equal(t, Debit, fee.RegulatoryFeesEffect)
	require.Equal(t, decimal.NewFromFloat(0.1), fee.ClearingFees)
	require.Equal(t, Debit, fee.ClearingFeesEffect)
	require.True(t, fee.Commission.Equal(decimal.Zero))
	require.Equal(t, None, fee.CommissionEffect)
	require.True(t, fee.ProprietaryIndexOptionFees.Equal(decimal.Zero))
	require.Equal(t, None, fee.ProprietaryIndexOptionFeesEffect)
	require.Equal(t, decimal.NewFromFloat(0.13), fee.TotalFees)
	require.Equal(t, Debit, fee.TotalFeesEffect)
}

func TestSubmitOrderError(t *testing.T) {
	setup()
	defer teardown()

	accountNumber := "5YZ55555"

	mux.HandleFunc(fmt.Sprintf("/accounts/%s/orders", accountNumber), func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(401)
		fmt.Fprint(writer, tastyUnauthorizedError)
	})

	symbol := "RIVN"
	quantity := float32(1)
	action1 := BTC

	symbol1 := EquityOptionsSymbology{
		Symbol:     symbol,
		OptionType: Call,
		Strike:     15,
		Expiration: time.Date(2023, 6, 23, 0, 0, 0, 0, time.Local),
	}

	order := NewOrder{
		TimeInForce: GTC,
		OrderType:   Limit,
		PriceEffect: Debit,
		Price:       0.04,
		Legs: []NewOrderLeg{
			{
				InstrumentType: EquityOptionIT,
				Symbol:         symbol1.Build(),
				Quantity:       quantity,
				Action:         action1,
			},
		},
		Rules: NewOrderRules{Conditions: []NewOrderCondition{
			{
				Action:         Route,
				Symbol:         symbol,
				InstrumentType: "Equity",
				Indicator:      Last,
				Comparator:     LTE,
				Threshold:      0.01,
			},
		}},
	}

	_, _, httpResp, err := client.SubmitOrder(accountNumber, order)
	expectedUnauthorized(t, err)
	require.NotNil(t, httpResp)
}

func TestGetAccountLiveOrders(t *testing.T) {
	setup()
	defer teardown()

	accountNumber := "5YZ55555"

	mux.HandleFunc(fmt.Sprintf("/accounts/%s/orders/live", accountNumber), func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, liveOrdersResp)
	})

	resp, httpResp, err := client.GetAccountLiveOrders(accountNumber)
	require.Nil(t, err)
	require.NotNil(t, httpResp)

	o := resp[0]

	require.Equal(t, accountNumber, o.AccountNumber)
	require.Equal(t, GTC, o.TimeInForce)
	require.Equal(t, Limit, o.OrderType)
	require.Equal(t, 1, o.Size)
	require.Equal(t, "RIVN", o.UnderlyingSymbol)
	require.Equal(t, EquityIT, o.UnderlyingInstrumentType)
	require.Equal(t, decimal.NewFromFloat(0.04), o.Price)
	require.Equal(t, Debit, o.PriceEffect)
	require.Equal(t, Contingent, o.Status)
	require.Equal(t, "Pending Condition", o.ContingentStatus)
	require.True(t, o.Cancellable)
	require.True(t, o.Editable)
	require.False(t, o.Edited)
	require.Equal(t, "2023-06-13T20:25:43.22Z", o.ReceivedAt.Format(time.RFC3339Nano))
	require.Equal(t, 1686687943220, o.UpdatedAt)

	ol := o.Legs[0]

	require.Equal(t, EquityOptionIT, ol.InstrumentType)
	require.Equal(t, "RIVN  230623C00015000", ol.Symbol)
	require.Equal(t, float32(1), ol.Quantity)
	require.Equal(t, float32(1), ol.RemainingQuantity)
	require.Equal(t, BTC, ol.Action)
	require.Empty(t, ol.Fills)

	oc := o.Rules.Conditions[0]

	require.Equal(t, 207561, oc.ID)
	require.Equal(t, Route, oc.Action)
	require.Equal(t, "RIVN", oc.Symbol)
	require.Equal(t, EquityIT, oc.InstrumentType)
	require.Equal(t, Last, oc.Indicator)
	require.Equal(t, LTE, oc.Comparator)
	require.Equal(t, decimal.NewFromFloat(0.01), oc.Threshold)
	require.False(t, oc.IsThresholdBasedOnNotional)

	pc := oc.PriceComponents[0]

	require.Equal(t, "RIVN", pc.Symbol)
	require.Equal(t, EquityIT, pc.InstrumentType)
	require.Equal(t, float32(1), pc.Quantity)
	require.Equal(t, Long, pc.QuantityDirection)
}

func TestGetAccountLiveOrdersError(t *testing.T) {
	setup()
	defer teardown()

	accountNumber := "5YZ55555"

	mux.HandleFunc(fmt.Sprintf("/accounts/%s/orders/live", accountNumber), func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(401)
		fmt.Fprint(writer, tastyUnauthorizedError)
	})

	_, httpResp, err := client.GetAccountLiveOrders(accountNumber)
	expectedUnauthorized(t, err)
	require.NotNil(t, httpResp)
}

func TestGetAccountOrders(t *testing.T) {
	setup()
	defer teardown()

	accountNumber := "5WV48989"

	mux.HandleFunc(fmt.Sprintf("/accounts/%s/orders", accountNumber), func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, accountOrdersResp)
	})

	resp, pagination, httpResp, err := client.GetAccountOrders(accountNumber, OrdersQuery{PerPage: 2})
	require.Nil(t, err)
	require.NotNil(t, httpResp)

	require.Equal(t, 2, len(resp))

	require.Equal(t, 2, pagination.PerPage)
}

func TestGetAccountOrdersError(t *testing.T) {
	setup()
	defer teardown()

	accountNumber := "5WV48989"

	mux.HandleFunc(fmt.Sprintf("/accounts/%s/orders", accountNumber), func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(401)
		fmt.Fprint(writer, tastyUnauthorizedError)
	})

	_, _, httpResp, err := client.GetAccountOrders(accountNumber, OrdersQuery{PerPage: 2})
	expectedUnauthorized(t, err)
	require.NotNil(t, httpResp)
}

func TestSubmitOrderECRDryRun(t *testing.T) {
	setup()
	defer teardown()

	accountNumber := "5WV48989"
	orderID := 68675

	mux.HandleFunc(fmt.Sprintf("/accounts/%s/orders/%d/dry-run", accountNumber, orderID), func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, orderECRDryRunResp)
	})

	orderECR := NewOrderECR{
		TimeInForce: Day,
		Price:       185.45,
		OrderType:   Limit,
		PriceEffect: Debit,
	}

	resp, httpResp, err := client.SubmitOrderECRDryRun(accountNumber, orderID, orderECR)
	require.Nil(t, err)
	require.NotNil(t, httpResp)

	o := resp.Order

	require.Equal(t, accountNumber, o.AccountNumber)
	require.Equal(t, Day, o.TimeInForce)
	require.Equal(t, Limit, o.OrderType)
	require.Equal(t, 1, o.Size)
	require.Equal(t, "AAPL", o.UnderlyingSymbol)
	require.Equal(t, EquityIT, o.UnderlyingInstrumentType)
	require.Equal(t, decimal.NewFromFloat(185.45), o.Price)
	require.Equal(t, Debit, o.PriceEffect)
	require.Equal(t, Contingent, o.Status)
	require.Equal(t, "Pending Order", o.ContingentStatus)
	require.True(t, o.Cancellable)
	require.True(t, o.Editable)
	require.False(t, o.Edited)
	require.Zero(t, o.UpdatedAt)
}

func TestSubmitOrderECRDryRunError(t *testing.T) {
	setup()
	defer teardown()

	accountNumber := "5WV48989"
	orderID := 68675

	mux.HandleFunc(fmt.Sprintf("/accounts/%s/orders/%d/dry-run", accountNumber, orderID), func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(401)
		fmt.Fprint(writer, tastyUnauthorizedError)
	})

	orderECR := NewOrderECR{
		TimeInForce: Day,
		Price:       185.45,
		OrderType:   Limit,
		PriceEffect: Debit,
	}

	_, httpResp, err := client.SubmitOrderECRDryRun(accountNumber, orderID, orderECR)
	expectedUnauthorized(t, err)
	require.NotNil(t, httpResp)
}

func TestGetOrder(t *testing.T) {
	setup()
	defer teardown()

	accountNumber := "5WV48989"
	orderID := 68675

	mux.HandleFunc(fmt.Sprintf("/accounts/%s/orders/%d", accountNumber, orderID), func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, getOrderResp)
	})

	o, httpResp, err := client.GetOrder(accountNumber, orderID)
	require.Nil(t, err)
	require.NotNil(t, httpResp)

	require.Equal(t, accountNumber, o.AccountNumber)
	require.Equal(t, Day, o.TimeInForce)
	require.Equal(t, Limit, o.OrderType)
	require.Equal(t, 1, o.Size)
	require.Equal(t, "AAPL", o.UnderlyingSymbol)
	require.Equal(t, EquityIT, o.UnderlyingInstrumentType)
	require.Equal(t, decimal.NewFromFloat(124.55), o.Price)
	require.Equal(t, Debit, o.PriceEffect)
	require.Equal(t, Contingent, o.Status)
	require.Equal(t, "Pending Condition", o.ContingentStatus)
	require.True(t, o.Cancellable)
	require.True(t, o.Editable)
	require.False(t, o.Edited)
	require.Equal(t, 1686698526525, o.UpdatedAt)
}

func TestGetOrderError(t *testing.T) {
	setup()
	defer teardown()

	accountNumber := "5WV48989"
	orderID := 68675

	mux.HandleFunc(fmt.Sprintf("/accounts/%s/orders/%d", accountNumber, orderID), func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(401)
		fmt.Fprint(writer, tastyUnauthorizedError)
	})

	_, httpResp, err := client.GetOrder(accountNumber, orderID)
	expectedUnauthorized(t, err)
	require.NotNil(t, httpResp)
}

func TestCancelOrder(t *testing.T) {
	setup()
	defer teardown()

	accountNumber := "5WV48989"
	orderID := 68677

	mux.HandleFunc(fmt.Sprintf("/accounts/%s/orders/%d", accountNumber, orderID), func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, cancelledOrderResp)
	})

	o, httpResp, err := client.CancelOrder(accountNumber, orderID)
	require.Nil(t, err)
	require.NotNil(t, httpResp)

	require.Equal(t, accountNumber, o.AccountNumber)
	require.Equal(t, Day, o.TimeInForce)
	require.Equal(t, Limit, o.OrderType)
	require.Equal(t, 1, o.Size)
	require.Equal(t, "AAPL", o.UnderlyingSymbol)
	require.Equal(t, EquityIT, o.UnderlyingInstrumentType)
	require.Equal(t, decimal.NewFromFloat(186.99), o.Price)
	require.Equal(t, Debit, o.PriceEffect)
	require.Equal(t, Cancelled, o.Status)
	require.False(t, o.Cancellable)
	require.Equal(t, "2023-06-14T01:02:24.795Z", o.CancelledAt.Format(time.RFC3339Nano))
	require.False(t, o.Editable)
	require.False(t, o.Edited)
	require.Equal(t, "2023-06-14T01:02:24.669Z", o.ReceivedAt.Format(time.RFC3339Nano))
	require.Equal(t, 1686704544800, o.UpdatedAt)
	require.Equal(t, "2023-06-14T01:02:24.794Z", o.TerminalAt.Format(time.RFC3339Nano))
}

func TestCancelOrderError(t *testing.T) {
	setup()
	defer teardown()

	accountNumber := "5WV48989"
	orderID := 68677

	mux.HandleFunc(fmt.Sprintf("/accounts/%s/orders/%d", accountNumber, orderID), func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(401)
		fmt.Fprint(writer, tastyUnauthorizedError)
	})

	_, httpResp, err := client.CancelOrder(accountNumber, orderID)
	expectedUnauthorized(t, err)
	require.NotNil(t, httpResp)
}

func TestReplaceOrder(t *testing.T) {
	setup()
	defer teardown()

	accountNumber := "5WV48989"
	orderID := 68678

	mux.HandleFunc(fmt.Sprintf("/accounts/%s/orders/%d", accountNumber, orderID), func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, replaceOrderResp)
	})

	orderECR := NewOrderECR{
		TimeInForce: Day,
		Price:       185.45,
		OrderType:   Limit,
		PriceEffect: Debit,
		ValueEffect: Debit,
	}

	o, httpResp, err := client.ReplaceOrder(accountNumber, orderID, orderECR)
	require.Nil(t, err)
	require.NotNil(t, httpResp)

	require.Equal(t, accountNumber, o.AccountNumber)
	require.Equal(t, Day, o.TimeInForce)
	require.Equal(t, Limit, o.OrderType)
	require.Equal(t, 1, o.Size)
	require.Equal(t, "AAPL", o.UnderlyingSymbol)
	require.Equal(t, EquityIT, o.UnderlyingInstrumentType)
	require.Equal(t, decimal.NewFromFloat(185.45), o.Price)
	require.Equal(t, Debit, o.PriceEffect)
	require.Equal(t, Contingent, o.Status)
	require.Equal(t, "Pending Condition", o.ContingentStatus)
	require.True(t, o.Cancellable)
	require.True(t, o.Editable)
	require.False(t, o.Edited)
	require.Equal(t, 1686706739960, o.UpdatedAt)
}

func TestReplaceOrderError(t *testing.T) {
	setup()
	defer teardown()

	accountNumber := "5WV48989"
	orderID := 68678

	mux.HandleFunc(fmt.Sprintf("/accounts/%s/orders/%d", accountNumber, orderID), func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(401)
		fmt.Fprint(writer, tastyUnauthorizedError)
	})

	orderECR := NewOrderECR{
		TimeInForce: Day,
		Price:       185.45,
		OrderType:   Limit,
		PriceEffect: Debit,
		ValueEffect: Debit,
	}

	_, httpResp, err := client.ReplaceOrder(accountNumber, orderID, orderECR)
	expectedUnauthorized(t, err)
	require.NotNil(t, httpResp)
}

func TestPatchOrder(t *testing.T) {
	setup()
	defer teardown()

	accountNumber := "5WV48989"
	orderID := 68680

	mux.HandleFunc(fmt.Sprintf("/accounts/%s/orders/%d", accountNumber, orderID), func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, patchOrderResp)
	})

	orderECR := NewOrderECR{
		TimeInForce: Day,
		Price:       187.45,
		OrderType:   Limit,
		PriceEffect: Debit,
		ValueEffect: Debit,
	}

	o, httpResp, err := client.PatchOrder(accountNumber, orderID, orderECR)
	require.Nil(t, err)
	require.NotNil(t, httpResp)

	require.Equal(t, accountNumber, o.AccountNumber)
	require.Equal(t, Day, o.TimeInForce)
	require.Equal(t, Limit, o.OrderType)
	require.Equal(t, 1, o.Size)
	require.Equal(t, "AAPL", o.UnderlyingSymbol)
	require.Equal(t, EquityIT, o.UnderlyingInstrumentType)
	require.Equal(t, decimal.NewFromFloat(187.45), o.Price)
	require.Equal(t, Debit, o.PriceEffect)
	require.Equal(t, Contingent, o.Status)
	require.Equal(t, "Pending Condition", o.ContingentStatus)
	require.True(t, o.Cancellable)
	require.True(t, o.Editable)
	require.False(t, o.Edited)
	require.Equal(t, 1686707204835, o.UpdatedAt)
}

func TestPatchOrderError(t *testing.T) {
	setup()
	defer teardown()

	accountNumber := "5WV48989"
	orderID := 68680

	mux.HandleFunc(fmt.Sprintf("/accounts/%s/orders/%d", accountNumber, orderID), func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(401)
		fmt.Fprint(writer, tastyUnauthorizedError)
	})

	orderECR := NewOrderECR{
		TimeInForce: Day,
		Price:       187.45,
		OrderType:   Limit,
		PriceEffect: Debit,
		ValueEffect: Debit,
	}

	_, httpResp, err := client.PatchOrder(accountNumber, orderID, orderECR)
	expectedUnauthorized(t, err)
	require.NotNil(t, httpResp)
}

func TestGetCustomerLiveOrders(t *testing.T) {
	setup()
	defer teardown()

	accountNumber := "5YZ55555"
	customerID := "me"

	mux.HandleFunc(fmt.Sprintf("/customers/%s/orders/live", customerID), func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, customerLiveOrdersResp)
	})

	resp, httpResp, err := client.GetCustomerLiveOrders(customerID, OrdersQuery{AccountNumbers: []string{accountNumber}})
	require.Nil(t, err)
	require.NotNil(t, httpResp)

	o := resp[0]

	require.Equal(t, accountNumber, o.AccountNumber)
	require.Equal(t, Day, o.TimeInForce)
	require.Equal(t, Limit, o.OrderType)
	require.Equal(t, 1, o.Size)
	require.Equal(t, "RIVN", o.UnderlyingSymbol)
	require.Equal(t, EquityIT, o.UnderlyingInstrumentType)
	require.Equal(t, decimal.NewFromFloat(0.01), o.Price)
	require.Equal(t, Debit, o.PriceEffect)
	require.Equal(t, Filled, o.Status)
	require.False(t, o.Cancellable)
	require.False(t, o.Editable)
	require.False(t, o.Edited)

	ol := o.Legs[0]

	require.Equal(t, EquityOptionIT, ol.InstrumentType)
	require.Equal(t, "RIVN  230623C00015000", ol.Symbol)
	require.Equal(t, float32(1), ol.Quantity)
	require.Equal(t, float32(0), ol.RemainingQuantity)
	require.Equal(t, BTC, ol.Action)

	fi := ol.Fills[0]

	require.Equal(t, "2263911504", fi.ExtGroupFillID)
	require.Equal(t, "90305", fi.ExtExecID)
	require.Equal(t, "3_OPT850090305", fi.FillID)
	require.Equal(t, float32(1), fi.Quantity)
	require.Equal(t, decimal.NewFromFloat(0.01), fi.FillPrice)
	require.Equal(t, "2023-06-23T14:12:04.214Z", fi.FilledAt.Format(time.RFC3339Nano))
	require.Equal(t, "CITADEL_OPTIONS_A", fi.DestinationVenue)
}

func TestGetCustomerLiveOrdersError(t *testing.T) {
	setup()
	defer teardown()

	customerID := "me"

	mux.HandleFunc(fmt.Sprintf("/customers/%s/orders/live", customerID), func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(401)
		fmt.Fprint(writer, customerOrdersErrorResp)
	})

	_, httpResp, err := client.GetCustomerLiveOrders(customerID, OrdersQuery{})
	require.NotNil(t, err)
	require.NotNil(t, httpResp)

	require.Equal(t, "\nError in request 401;\nCode: validation_error\nMessage: Request validation failed", err.Error())
	// require.Equal(t, "Request validation failed", err.Message)
	// require.NotEmpty(t, err.Errors)
	// require.Equal(t, "account-numbers", err.Errors[0].Domain)
	// require.Equal(t, "is missing", err.Errors[0].Reason)
}

func TestGetCustomerOrders(t *testing.T) {
	setup()
	defer teardown()

	accountNumber := "5YZ55555"
	customerID := "me"

	mux.HandleFunc(fmt.Sprintf("/customers/%s/orders", customerID), func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, customerLiveOrdersResp)
	})

	resp, httpResp, err := client.GetCustomerOrders(customerID, OrdersQuery{AccountNumbers: []string{accountNumber}})
	require.Nil(t, err)
	require.NotNil(t, httpResp)

	o := resp[0]

	require.Equal(t, accountNumber, o.AccountNumber)
	require.Equal(t, Day, o.TimeInForce)
	require.Equal(t, Limit, o.OrderType)
	require.Equal(t, 1, o.Size)
	require.Equal(t, "RIVN", o.UnderlyingSymbol)
	require.Equal(t, EquityIT, o.UnderlyingInstrumentType)
	require.Equal(t, decimal.NewFromFloat(0.01), o.Price)
	require.Equal(t, Debit, o.PriceEffect)
	require.Equal(t, Filled, o.Status)
	require.False(t, o.Cancellable)
	require.False(t, o.Editable)
	require.False(t, o.Edited)

	ol := o.Legs[0]

	require.Equal(t, EquityOptionIT, ol.InstrumentType)
	require.Equal(t, "RIVN  230623C00015000", ol.Symbol)
	require.Equal(t, float32(1), ol.Quantity)
	require.Equal(t, float32(0), ol.RemainingQuantity)
	require.Equal(t, BTC, ol.Action)

	fi := ol.Fills[0]

	require.Equal(t, "2263911504", fi.ExtGroupFillID)
	require.Equal(t, "90305", fi.ExtExecID)
	require.Equal(t, "3_OPT850090305", fi.FillID)
	require.Equal(t, float32(1), fi.Quantity)
	require.Equal(t, decimal.NewFromFloat(0.01), fi.FillPrice)
	require.Equal(t, "2023-06-23T14:12:04.214Z", fi.FilledAt.Format(time.RFC3339Nano))
	require.Equal(t, "CITADEL_OPTIONS_A", fi.DestinationVenue)
}

func TestGetCustomerOrdersError(t *testing.T) {
	setup()
	defer teardown()

	customerID := "me"

	mux.HandleFunc(fmt.Sprintf("/customers/%s/orders", customerID), func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(401)
		fmt.Fprint(writer, customerOrdersErrorResp)
	})

	_, httpResp, err := client.GetCustomerOrders(customerID, OrdersQuery{})
	require.NotNil(t, err)
	require.NotNil(t, httpResp)

	require.Equal(t, "\nError in request 401;\nCode: validation_error\nMessage: Request validation failed", err.Error())
	// require.Equal(t, "Request validation failed", err.Message)
	// require.NotEmpty(t, err.Errors)
	// require.Equal(t, "account-numbers", err.Errors[0].Domain)
	// require.Equal(t, "is missing", err.Errors[0].Reason)
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

const accountOrdersResp = `{
  "data": {
    "items": [
      {
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
      {
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
        "status": "Cancelled",
        "cancellable": false,
        "cancelled-at": "2023-06-14T01:46:44.799+00:00",
        "editable": false,
        "edited": true,
        "received-at": "2023-06-14T01:38:59.936+00:00",
        "updated-at": 1686707204813,
        "terminal-at": "2023-06-14T01:46:44.799+00:00",
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
      }
    ]
  },
  "context": "/accounts/5WV48989/orders",
  "pagination": {
    "per-page": 2,
    "page-offset": 0,
    "item-offset": 0,
    "total-items": 8,
    "total-pages": 4,
    "current-item-count": 2,
    "previous-link": null,
    "next-link": null,
    "paging-link-template": null
  }
}`

const reconfirmResp = `{"error":{"code":"cannot_reconfirm_order","message":"The order could not be reconfirmed."}}`

const customerOrdersErrorResp = `{
    "error": {
        "code": "validation_error",
        "message": "Request validation failed",
        "errors": [
            {
                "domain": "account-numbers",
                "reason": "is missing"
            }
        ]
    }
}`

const customerLiveOrdersResp = `{
  "data": {
    "items": [
      {
        "id": 274344092,
        "account-number": "5YZ55555",
        "time-in-force": "Day",
        "order-type": "Limit",
        "size": 1,
        "underlying-symbol": "RIVN",
        "underlying-instrument-type": "Equity",
        "price": "0.01",
        "price-effect": "Debit",
        "status": "Filled",
        "cancellable": false,
        "editable": false,
        "edited": false,
        "ext-exchange-order-number": "60722521974940",
        "ext-client-order-id": "9c0000373a105a289c",
        "ext-global-order-number": 14138,
        "received-at": "2023-06-23T13:52:59.062+00:00",
        "updated-at": 1687529524255,
        "terminal-at": "2023-06-23T14:12:04.250+00:00",
        "legs": [
          {
            "instrument-type": "Equity Option",
            "symbol": "RIVN  230623C00015000",
            "quantity": 1,
            "remaining-quantity": 0,
            "action": "Buy to Close",
            "fills": [
              {
                "ext-group-fill-id": "2263911504",
                "ext-exec-id": "90305",
                "fill-id": "3_OPT850090305",
                "quantity": 1,
                "fill-price": "0.01",
                "filled-at": "2023-06-23T14:12:04.214+00:00",
                "destination-venue": "CITADEL_OPTIONS_A"
              }
            ]
          }
        ]
      }
    ]
  },
  "context": "/customers/me/orders/live"
}`
