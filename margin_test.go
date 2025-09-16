package tasty //nolint:testpackage // testing private field

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func TestGetMarginRequirements(t *testing.T) {
	setup()
	defer teardown()

	accountNumber := "5YZ55555"

	mux.HandleFunc(fmt.Sprintf("/margin/accounts/%s/requirements", accountNumber), func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, marginReqResp)
	})

	resp, httpResp, err := client.GetMarginRequirements(accountNumber)
	require.Nil(t, err)
	require.NotNil(t, httpResp)

	require.Equal(t, accountNumber, resp.AccountNumber)
	require.Equal(t, "Total", resp.Description)
	require.Equal(t, "IRA Margin", resp.MarginCalculationType)
	require.Equal(t, "Defined Risk Spreads", resp.OptionLevel)
	require.Equal(t, decimal.NewFromFloat(432457.79), resp.MarginRequirement)
	require.Equal(t, Debit, resp.MarginRequirementEffect)
	require.Equal(t, decimal.NewFromFloat(432457.79), resp.InitialRequirement)
	require.Equal(t, Debit, resp.InitialRequirementEffect)
	require.Equal(t, decimal.NewFromFloat(432457.79), resp.MaintenanceRequirement)
	require.Equal(t, Debit, resp.MaintenanceRequirementEffect)
	require.Equal(t, decimal.NewFromFloat(452578.552), resp.MarginEquity)
	require.Equal(t, Credit, resp.MarginEquityEffect)
	require.Equal(t, decimal.NewFromFloat(12900.762), resp.OptionBuyingPower)
	require.Equal(t, Credit, resp.OptionBuyingPowerEffect)
	require.Equal(t, decimal.NewFromFloat(432457.79), resp.RegTMarginRequirement)
	require.Equal(t, Debit, resp.RegTMarginRequirementEffect)
	require.Equal(t, decimal.NewFromFloat(12900.762), resp.RegTOptionBuyingPower)
	require.Equal(t, Credit, resp.RegTOptionBuyingPowerEffect)
	require.Equal(t, 1686135853860, resp.LastStateTimestamp)

	// Groups
	groups := resp.MarginGroups

	require.Equal(t, 2, len(groups))

	amd := groups[0]

	require.Equal(t, "AMD", amd.Description)
	require.Equal(t, "AMD", amd.Code)
	require.Equal(t, "AMD", amd.UnderlyingSymbol)
	require.Equal(t, "Equity", amd.UnderlyingType)
	require.Equal(t, decimal.NewFromFloat(0.3), amd.ExpectedPriceRangeUpPercent)
	require.Equal(t, decimal.NewFromFloat(0.3), amd.ExpectedPriceRangeDownPercent)
	require.Equal(t, decimal.NewFromFloat(1.01), amd.PointOfNoReturnPercent)
	require.Equal(t, "IRA Margin", amd.MarginCalculationType)
	require.Equal(t, decimal.NewFromFloat(39653.49), amd.MarginRequirement)
	require.Equal(t, Debit, amd.MarginRequirementEffect)
	require.Equal(t, decimal.NewFromFloat(39653.49), amd.InitialRequirement)
	require.Equal(t, Debit, amd.InitialRequirementEffect)
	require.Equal(t, decimal.NewFromFloat(39653.49), amd.MaintenanceRequirement)
	require.Equal(t, Debit, amd.MaintenanceRequirementEffect)
	require.Equal(t, decimal.NewFromFloat(39653.49), amd.BuyingPower)
	require.Equal(t, Credit, amd.BuyingPowerEffect)
	require.True(t, amd.PriceIncreasePercent.Equal(decimal.NewFromInt(1)))
	require.True(t, amd.PriceDecreasePercent.Equal(decimal.NewFromInt(-1)))

	// AMD Holdings
	holdings := amd.Holdings

	require.Equal(t, 1, len(holdings))

	holding := holdings[0]

	require.Equal(t, "LONG_UNDERLYING", holding.Description)
	require.Equal(t, decimal.NewFromFloat(39653.49), holding.MarginRequirement)
	require.Equal(t, Debit, holding.MarginRequirementEffect)
	require.Equal(t, decimal.NewFromFloat(39653.49), holding.InitialRequirement)
	require.Equal(t, Debit, holding.InitialRequirementEffect)
	require.False(t, holding.IncludesWorkingOrder)
	require.Equal(t, decimal.NewFromFloat(39653.49), holding.BuyingPower)
	require.Equal(t, Credit, holding.BuyingPowerEffect)

	// AMD Holding Position Entries
	entries := holding.PositionEntries

	require.Equal(t, 1, len(entries))

	entry := entries[0]

	require.Equal(t, "AMD", entry.InstrumentSymbol)
	require.Equal(t, EquityIT, entry.InstrumentType)
	require.Equal(t, decimal.NewFromFloat(336.446), entry.Quantity)
	require.Zero(t, entry.AverageOpenPrice)
	require.Equal(t, decimal.NewFromFloat(117.86), entry.ClosePrice)
	require.Zero(t, entry.FixingPrice)

	rivn := groups[1]

	require.Equal(t, "RIVN", rivn.Description)
	require.Equal(t, "RIVN", rivn.Code)
	require.Equal(t, "RIVN", rivn.UnderlyingSymbol)
	require.Equal(t, "Equity", rivn.UnderlyingType)
	require.Equal(t, decimal.NewFromFloat(0.5), rivn.ExpectedPriceRangeUpPercent)
	require.Equal(t, decimal.NewFromFloat(0.5), rivn.ExpectedPriceRangeDownPercent)
	require.Equal(t, decimal.NewFromFloat(1.01), rivn.PointOfNoReturnPercent)
	require.Equal(t, "IRA Margin", rivn.MarginCalculationType)
	require.True(t, rivn.MarginRequirement.Equal(decimal.NewFromInt(56000)))
	require.Equal(t, Debit, rivn.MarginRequirementEffect)
	require.True(t, rivn.InitialRequirement.Equal(decimal.NewFromInt(56000)))
	require.Equal(t, Debit, rivn.InitialRequirementEffect)
	require.True(t, rivn.MaintenanceRequirement.Equal(decimal.NewFromInt(56000)))
	require.Equal(t, Debit, rivn.MaintenanceRequirementEffect)
	require.True(t, rivn.BuyingPower.Equal(decimal.NewFromInt(54440)))
	require.Equal(t, Credit, rivn.BuyingPowerEffect)
	require.True(t, rivn.PriceIncreasePercent.Equal(decimal.NewFromInt(1)))
	require.True(t, rivn.PriceDecreasePercent.Equal(decimal.NewFromInt(-1)))

	// RIVN Holdings
	holdings = rivn.Holdings

	require.Equal(t, 1, len(holdings))

	holding = holdings[0]

	require.Equal(t, "NAKED_OPTION", holding.Description)
	require.True(t, holding.MarginRequirement.Equal(decimal.NewFromInt(56000)))
	require.Equal(t, Debit, holding.MarginRequirementEffect)
	require.True(t, holding.InitialRequirement.Equal(decimal.NewFromInt(56000)))
	require.Equal(t, Debit, holding.InitialRequirementEffect)
	require.False(t, holding.IncludesWorkingOrder)
	require.True(t, holding.BuyingPower.Equal(decimal.NewFromInt(54440)))
	require.Equal(t, Credit, holding.BuyingPowerEffect)

	// RIVN Holding Position Entries
	entries = holding.PositionEntries

	require.Equal(t, 1, len(entries))

	entry = entries[0]

	require.Equal(t, "RIVN  230609P00014000", entry.InstrumentSymbol)
	require.Equal(t, EquityOptionIT, entry.InstrumentType)
	require.True(t, entry.Quantity.Equal(decimal.NewFromInt(-40)))
	require.Zero(t, entry.AverageOpenPrice)
	require.Equal(t, decimal.NewFromFloat(0.32), entry.ClosePrice)
	require.Zero(t, entry.FixingPrice)
	require.True(t, entry.StrikePrice.Equal(decimal.NewFromInt(14)))
	require.Equal(t, Put, entry.OptionType)
	require.True(t, entry.DeliverableQuantity.Equal(decimal.NewFromInt(4000)))
	require.Equal(t, "2023-06-09", entry.ExpirationDate)
}

func TestGetMarginRequirementsError(t *testing.T) {
	setup()
	defer teardown()

	accountNumber := "5YZ55555"

	mux.HandleFunc(fmt.Sprintf("/margin/accounts/%s/requirements", accountNumber), func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(401)
		fmt.Fprint(writer, tastyUnauthorizedError)
	})

	_, _, err := client.GetMarginRequirements(accountNumber)
	expectedUnauthorized(t, err)
}

func TestGetEffectiveMarginRequirements(t *testing.T) {
	setup()
	defer teardown()

	accountNumber := "5YZ55555"
	underlyingSymbol := "AAPL"

	mux.HandleFunc(fmt.Sprintf("/accounts/%s/margin-requirements/%s/effective", accountNumber, underlyingSymbol), func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, effectiveMarginRequirementsResp)
	})

	resp, httpResp, err := client.GetEffectiveMarginRequirements(accountNumber, underlyingSymbol)
	require.Nil(t, err)
	require.NotNil(t, httpResp)

	require.Equal(t, underlyingSymbol, resp.UnderlyingSymbol)
	require.Equal(t, decimal.NewFromFloat(.5), resp.LongEquityInitial)
	require.Equal(t, decimal.NewFromFloat(.5), resp.ShortEquityInitial)
	require.Equal(t, decimal.NewFromFloat(.25), resp.LongEquityMaintenance)
	require.Equal(t, decimal.NewFromFloat(.3), resp.ShortEquityMaintenance)
	require.Equal(t, decimal.NewFromFloat(.2), resp.NakedOptionStandard)
	require.Equal(t, decimal.NewFromFloat(.1), resp.NakedOptionMinimum)
	require.True(t, decimal.NewFromInt(250).Equal(resp.NakedOptionFloor))
}

func TestGetEffectiveMarginRequirementsError(t *testing.T) {
	setup()
	defer teardown()

	accountNumber := "5YZ55555"
	underlyingSymbol := "AAPL"

	mux.HandleFunc(fmt.Sprintf("/accounts/%s/margin-requirements/%s/effective", accountNumber, underlyingSymbol), func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(401)
		fmt.Fprint(writer, tastyUnauthorizedError)
	})

	_, _, err := client.GetEffectiveMarginRequirements(accountNumber, underlyingSymbol)
	expectedUnauthorized(t, err)
}

const marginReqResp = `{
  "data": {
    "account-number": "5YZ55555",
    "description": "Total",
    "margin-calculation-type": "IRA Margin",
    "option-level": "Defined Risk Spreads",
    "margin-requirement": "432457.79",
    "margin-requirement-effect": "Debit",
    "initial-requirement": "432457.79",
    "initial-requirement-effect": "Debit",
    "maintenance-requirement": "432457.79",
    "maintenance-requirement-effect": "Debit",
    "margin-equity": "452578.552",
    "margin-equity-effect": "Credit",
    "option-buying-power": "12900.762",
    "option-buying-power-effect": "Credit",
    "reg-t-margin-requirement": "432457.79",
    "reg-t-margin-requirement-effect": "Debit",
    "reg-t-option-buying-power": "12900.762",
    "reg-t-option-buying-power-effect": "Credit",
    "maintenance-excess": "12900.762",
    "maintenance-excess-effect": "Credit",
    "groups": [
      {
        "description": "AMD",
        "code": "AMD",
        "underlying-symbol": "AMD",
        "underlying-type": "Equity",
        "expected-price-range-up-percent": "0.3",
        "expected-price-range-down-percent": "0.3",
        "point-of-no-return-percent": "1.01",
        "margin-calculation-type": "IRA Margin",
        "margin-requirement": "39653.49",
        "margin-requirement-effect": "Debit",
        "initial-requirement": "39653.49",
        "initial-requirement-effect": "Debit",
        "maintenance-requirement": "39653.49",
        "maintenance-requirement-effect": "Debit",
        "buying-power": "39653.49",
        "buying-power-effect": "Credit",
        "groups": [
          {
            "description": "LONG_UNDERLYING",
            "margin-requirement": "39653.49",
            "margin-requirement-effect": "Debit",
            "initial-requirement": "39653.49",
            "initial-requirement-effect": "Debit",
            "maintenance-requirement": "39653.49",
            "maintenance-requirement-effect": "Debit",
            "includes-working-order": false,
            "buying-power": "39653.49",
            "buying-power-effect": "Credit",
            "position-entries": [
              {
                "instrument-symbol": "AMD",
                "instrument-type": "Equity",
                "quantity": "336.446",
                "average-open-price": "NaN",
                "close-price": "117.86",
                "fixing-price": "NaN"
              }
            ]
          }
        ],
        "price-increase-percent": "1.0",
        "price-decrease-percent": "-1.0"
      },
      {
        "description": "RIVN",
        "code": "RIVN",
        "underlying-symbol": "RIVN",
        "underlying-type": "Equity",
        "expected-price-range-up-percent": "0.5",
        "expected-price-range-down-percent": "0.5",
        "point-of-no-return-percent": "1.01",
        "margin-calculation-type": "IRA Margin",
        "margin-requirement": "56000.0",
        "margin-requirement-effect": "Debit",
        "initial-requirement": "56000.0",
        "initial-requirement-effect": "Debit",
        "maintenance-requirement": "56000.0",
        "maintenance-requirement-effect": "Debit",
        "buying-power": "54440.0",
        "buying-power-effect": "Credit",
        "groups": [
          {
            "description": "NAKED_OPTION",
            "margin-requirement": "56000.0",
            "margin-requirement-effect": "Debit",
            "initial-requirement": "56000.0",
            "initial-requirement-effect": "Debit",
            "maintenance-requirement": "56000.0",
            "maintenance-requirement-effect": "Debit",
            "includes-working-order": false,
            "buying-power": "54440.0",
            "buying-power-effect": "Credit",
            "position-entries": [
              {
                "instrument-symbol": "RIVN  230609P00014000",
                "instrument-type": "Equity Option",
                "quantity": "-40.0",
                "average-open-price": "NaN",
                "close-price": "0.32",
                "fixing-price": "NaN",
                "strike-price": "14.0",
                "option-type": "P",
                "deliverable-quantity": "4000.0",
                "expiration-date": "2023-06-09"
              }
            ]
          }
        ],
        "price-increase-percent": "1.0",
        "price-decrease-percent": "-1.0"
      }
    ],
    "last-state-timestamp": 1686135853860
  }
}`

const effectiveMarginRequirementsResp = `{
    "data": {
        "underlying-symbol": "AAPL",
        "long-equity-initial": "0.5",
        "short-equity-initial": "0.5",
        "long-equity-maintenance": "0.25",
        "short-equity-maintenance": "0.3",
        "naked-option-standard": "0.2",
        "naked-option-minimum": "0.1",
        "naked-option-floor": "250.0"
    },
    "context": "/accounts/5YZ55555/margin-requirements/AAPL/effective"
}`
