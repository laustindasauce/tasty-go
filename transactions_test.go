package tasty //nolint:testpackage // testing private field

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestGetAccountTransactions(t *testing.T) {
	setup()
	defer teardown()

	accountNumber := "5YZ55555"
	query := TransactionsQuery{PerPage: 3}

	mux.HandleFunc(fmt.Sprintf("/accounts/%s/transactions", accountNumber), func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, transactionsResp)
	})

	resp, pagination, err := client.GetAccountTransactions(accountNumber, query)
	require.Nil(t, err)

	require.Equal(t, 3, len(resp))

	tr := resp[0]

	require.Equal(t, 2455333333, tr.ID)
	require.Equal(t, accountNumber, tr.AccountNumber)
	require.Equal(t, "RIVN  230623C00015000", tr.Symbol)
	require.Equal(t, EquityOptionIT, tr.InstrumentType)
	require.Equal(t, "RIVN", tr.UnderlyingSymbol)
	require.Equal(t, "Trade", tr.TransactionType)
	require.Equal(t, STO, tr.TransactionSubType)
	require.Equal(t, "Sold 40 RIVN 06/23/23 Call 15.00 @ 0.47", tr.Description)
	require.Equal(t, STO, tr.Action)
	require.Equal(t, StringToFloat32(40), tr.Quantity)
	require.Equal(t, StringToFloat32(.47), tr.Price)
	require.Equal(t, "2023-06-12T13:53:08.199Z", tr.ExecutedAt.Format(time.RFC3339Nano))
	require.Equal(t, "2023-06-12", tr.TransactionDate)
	require.Equal(t, StringToFloat32(1880), tr.Value)
	require.Equal(t, Credit, tr.ValueEffect)
	require.Equal(t, StringToFloat32(0.042), tr.RegulatoryFees)
	require.Equal(t, Debit, tr.RegulatoryFeesEffect)
	require.Equal(t, StringToFloat32(0.4), tr.ClearingFees)
	require.Equal(t, Debit, tr.ClearingFeesEffect)
	require.Equal(t, StringToFloat32(1878.858), tr.NetValue)
	require.Equal(t, Credit, tr.NetValueEffect)
	require.Equal(t, StringToFloat32(40), tr.Commission)
	require.Equal(t, Debit, tr.CommissionEffect)
	require.Equal(t, StringToFloat32(0), tr.ProprietaryIndexOptionFees)
	require.Equal(t, None, tr.ProprietaryIndexOptionFeesEffect)
	require.True(t, tr.IsEstimatedFee)
	require.Equal(t, "6309579568813", tr.ExtExchangeOrderNumber)
	require.Equal(t, 1469, tr.ExtGlobalOrderNumber)
	require.Equal(t, "0", tr.ExtGroupID)
	require.Equal(t, "613761478", tr.ExtGroupFillID)
	require.Equal(t, "6075255", tr.ExtExecID)
	require.Equal(t, "4_A5F-1DY-3M7P-3", tr.ExecID)
	require.Equal(t, "D", tr.Exchange)
	require.Equal(t, 272610989, tr.OrderID)
	require.Equal(t, "", tr.ExchangeAffiliationIDentifier)
	require.Equal(t, 1, tr.LegCount)
	require.Equal(t, "WOLVERINE_OPTIONS_A", tr.DestinationVenue)

	require.Equal(t, 3, pagination.PerPage)
	require.Equal(t, 0, pagination.PageOffset)
	require.Equal(t, 0, pagination.ItemOffset)
	require.Equal(t, 71, pagination.TotalItems)
	require.Equal(t, 24, pagination.TotalPages)
	require.Equal(t, 3, pagination.CurrentItemCount)
	require.Nil(t, pagination.PreviousLink)
	require.Nil(t, pagination.NextLink)
	require.Nil(t, pagination.PagingLinkTemplate)
}

func TestGetAccountTransactionsError(t *testing.T) {
	setup()
	defer teardown()

	accountNumber := "5YZ55555"
	query := TransactionsQuery{PerPage: 3}

	mux.HandleFunc(fmt.Sprintf("/accounts/%s/transactions", accountNumber), func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(401)
		fmt.Fprint(writer, tastyUnauthorizedError)
	})

	_, _, err := client.GetAccountTransactions(accountNumber, query)
	expectedUnauthorized(t, err)
}

func TestGetAccountTransaction(t *testing.T) {
	setup()
	defer teardown()

	accountNumber := "5YZ55555"
	id := 2455333333

	mux.HandleFunc(fmt.Sprintf("/accounts/%s/transactions/%d", accountNumber, id), func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, transactionResp)
	})

	tr, err := client.GetAccountTransaction(accountNumber, id)
	require.Nil(t, err)

	require.Equal(t, 2455333333, tr.ID)
	require.Equal(t, accountNumber, tr.AccountNumber)
	require.Equal(t, "RIVN  230623C00015000", tr.Symbol)
	require.Equal(t, EquityOptionIT, tr.InstrumentType)
	require.Equal(t, "RIVN", tr.UnderlyingSymbol)
	require.Equal(t, "Trade", tr.TransactionType)
	require.Equal(t, STO, tr.TransactionSubType)
	require.Equal(t, "Sold 40 RIVN 06/23/23 Call 15.00 @ 0.47", tr.Description)
	require.Equal(t, STO, tr.Action)
	require.Equal(t, StringToFloat32(40), tr.Quantity)
	require.Equal(t, StringToFloat32(.47), tr.Price)
	require.Equal(t, "2023-06-12T13:53:08.199Z", tr.ExecutedAt.Format(time.RFC3339Nano))
	require.Equal(t, "2023-06-12", tr.TransactionDate)
	require.Equal(t, StringToFloat32(1880), tr.Value)
	require.Equal(t, Credit, tr.ValueEffect)
	require.Equal(t, StringToFloat32(0.042), tr.RegulatoryFees)
	require.Equal(t, Debit, tr.RegulatoryFeesEffect)
	require.Equal(t, StringToFloat32(0.4), tr.ClearingFees)
	require.Equal(t, Debit, tr.ClearingFeesEffect)
	require.Equal(t, StringToFloat32(1878.858), tr.NetValue)
	require.Equal(t, Credit, tr.NetValueEffect)
	require.Equal(t, StringToFloat32(40), tr.Commission)
	require.Equal(t, Debit, tr.CommissionEffect)
	require.Equal(t, StringToFloat32(0), tr.ProprietaryIndexOptionFees)
	require.Equal(t, None, tr.ProprietaryIndexOptionFeesEffect)
	require.True(t, tr.IsEstimatedFee)
	require.Equal(t, "6309579568813", tr.ExtExchangeOrderNumber)
	require.Equal(t, 1469, tr.ExtGlobalOrderNumber)
	require.Equal(t, "0", tr.ExtGroupID)
	require.Equal(t, "613761478", tr.ExtGroupFillID)
	require.Equal(t, "6075255", tr.ExtExecID)
	require.Equal(t, "4_A5F-1DY-3M7P-3", tr.ExecID)
	require.Equal(t, "D", tr.Exchange)
	require.Equal(t, 272610989, tr.OrderID)
	require.Equal(t, "", tr.ExchangeAffiliationIDentifier)
	require.Equal(t, 1, tr.LegCount)
	require.Equal(t, "WOLVERINE_OPTIONS_A", tr.DestinationVenue)
}

func TestGetAccountTransactionError(t *testing.T) {
	setup()
	defer teardown()

	accountNumber := "5YZ55555"
	id := 2455333333

	mux.HandleFunc(fmt.Sprintf("/accounts/%s/transactions/%d", accountNumber, id), func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(401)
		fmt.Fprint(writer, tastyUnauthorizedError)
	})

	_, err := client.GetAccountTransaction(accountNumber, id)
	expectedUnauthorized(t, err)
}

func TestGetAccountTransactionFees(t *testing.T) {
	setup()
	defer teardown()

	accountNumber := "5YZ55555"

	mux.HandleFunc(fmt.Sprintf("/accounts/%s/transactions/total-fees", accountNumber), func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, transactionFeesResp)
	})

	fees, err := client.GetAccountTransactionFees(accountNumber, nil)
	require.Nil(t, err)

	require.Equal(t, StringToFloat32(0), fees.TotalFees)
	require.Equal(t, None, fees.TotalFeesEffect)
}

func TestGetAccountTransactionFeesWithParams(t *testing.T) {
	setup()
	defer teardown()

	accountNumber := "5YZ55555"
	queryDate := time.Date(2023, 6, 16, 0, 0, 0, 0, time.Local)

	mux.HandleFunc(fmt.Sprintf("/accounts/%s/transactions/total-fees", accountNumber), func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, transactionFeesResp)
		require.Equal(t, "2023-06-16", request.URL.Query().Get("date"))
	})

	fees, err := client.GetAccountTransactionFees(accountNumber, &queryDate)
	require.Nil(t, err)

	require.Equal(t, StringToFloat32(0), fees.TotalFees)
	require.Equal(t, None, fees.TotalFeesEffect)
}

func TestGetAccountTransactionFeesError(t *testing.T) {
	setup()
	defer teardown()

	accountNumber := "5YZ55555"

	mux.HandleFunc(fmt.Sprintf("/accounts/%s/transactions/total-fees", accountNumber), func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(401)
		fmt.Fprint(writer, tastyUnauthorizedError)
	})

	_, err := client.GetAccountTransactionFees(accountNumber, nil)
	expectedUnauthorized(t, err)
}

const transactionsResp = `{
  "data": {
    "items": [
      {
        "id": 2455333333,
        "account-number": "5YZ55555",
        "symbol": "RIVN  230623C00015000",
        "instrument-type": "Equity Option",
        "underlying-symbol": "RIVN",
        "transaction-type": "Trade",
        "transaction-sub-type": "Sell to Open",
        "description": "Sold 40 RIVN 06/23/23 Call 15.00 @ 0.47",
        "action": "Sell to Open",
        "quantity": "40.0",
        "price": "0.47",
        "executed-at": "2023-06-12T13:53:08.199+00:00",
        "transaction-date": "2023-06-12",
        "value": "1880.0",
        "value-effect": "Credit",
        "regulatory-fees": "0.042",
        "regulatory-fees-effect": "Debit",
        "clearing-fees": "0.40",
        "clearing-fees-effect": "Debit",
        "net-value": "1878.858",
        "net-value-effect": "Credit",
        "commission": "40.0",
        "commission-effect": "Debit",
        "proprietary-index-option-fees": "0.0",
        "proprietary-index-option-fees-effect": "None",
        "is-estimated-fee": true,
        "ext-exchange-order-number": "6309579568813",
        "ext-global-order-number": 1469,
        "ext-group-id": "0",
        "ext-group-fill-id": "613761478",
        "ext-exec-id": "6075255",
        "exec-id": "4_A5F-1DY-3M7P-3",
        "exchange": "D",
        "order-id": 272610989,
        "exchange-affiliation-identifier": "",
        "leg-count": 1,
        "destination-venue": "WOLVERINE_OPTIONS_A"
      },
      {
        "id": 2454865989,
        "account-number": "5YZ55555",
        "symbol": "RIVN  230609P00014000",
        "instrument-type": "Equity Option",
        "underlying-symbol": "RIVN",
        "transaction-type": "Receive Deliver",
        "transaction-sub-type": "Assignment",
        "description": "Removal of option due to assignment",
        "quantity": "40.0",
        "executed-at": "2023-06-09T21:00:00.000+00:00",
        "transaction-date": "2023-06-09",
        "value": "0.0",
        "value-effect": "None",
        "net-value": "0.0",
        "net-value-effect": "None",
        "is-estimated-fee": true
      },
      {
        "id": 245478958,
        "account-number": "5YZ55555",
        "symbol": "RIVN",
        "instrument-type": "Equity",
        "underlying-symbol": "RIVN",
        "transaction-type": "Receive Deliver",
        "transaction-sub-type": "Buy to Open",
        "description": "Buy to Open 4000 RIVN @ 14.00",
        "action": "Buy to Open",
        "quantity": "4000.0",
        "price": "14.0",
        "executed-at": "2023-06-09T21:00:00.000+00:00",
        "transaction-date": "2023-06-09",
        "value": "56000.0",
        "value-effect": "Debit",
        "clearing-fees": "200.0",
        "clearing-fees-effect": "Debit",
        "net-value": "56200.0",
        "net-value-effect": "Debit",
        "is-estimated-fee": true
      }
    ]
  },
  "context": "/accounts/5YZ55555/transactions",
  "pagination": {
    "per-page": 3,
    "page-offset": 0,
    "item-offset": 0,
    "total-items": 71,
    "total-pages": 24,
    "current-item-count": 3,
    "previous-link": null,
    "next-link": null,
    "paging-link-template": null
  }
}`

const transactionResp = `{
    "data": {
        "id": 2455333333,
        "account-number": "5YZ55555",
        "symbol": "RIVN  230623C00015000",
        "instrument-type": "Equity Option",
        "underlying-symbol": "RIVN",
        "transaction-type": "Trade",
        "transaction-sub-type": "Sell to Open",
        "description": "Sold 40 RIVN 06/23/23 Call 15.00 @ 0.47",
        "action": "Sell to Open",
        "quantity": "40.0",
        "price": "0.47",
        "executed-at": "2023-06-12T13:53:08.199+00:00",
        "transaction-date": "2023-06-12",
        "value": "1880.0",
        "value-effect": "Credit",
        "regulatory-fees": "0.042",
        "regulatory-fees-effect": "Debit",
        "clearing-fees": "0.40",
        "clearing-fees-effect": "Debit",
        "net-value": "1878.858",
        "net-value-effect": "Credit",
        "commission": "40.0",
        "commission-effect": "Debit",
        "proprietary-index-option-fees": "0.0",
        "proprietary-index-option-fees-effect": "None",
        "is-estimated-fee": true,
        "ext-exchange-order-number": "6309579568813",
        "ext-global-order-number": 1469,
        "ext-group-id": "0",
        "ext-group-fill-id": "613761478",
        "ext-exec-id": "6075255",
        "exec-id": "4_A5F-1DY-3M7P-3",
        "exchange": "D",
        "order-id": 272610989,
        "exchange-affiliation-identifier": "",
        "leg-count": 1,
        "destination-venue": "WOLVERINE_OPTIONS_A"
    },
    "context": "/accounts/5YZ55555/transactions/245564891"
}`

const transactionFeesResp = `{
    "data": {
        "total-fees": "0.0",
        "total-fees-effect": "None"
    },
    "context": "/accounts/5YZ55555/transactions/total-fees"
}`
