package tasty //nolint:testpackage // testing private field

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetQuantityDecimalPrecisions(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/instruments/quantity-decimal-precisions", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, quantityDecimalPrecisionsResp)
	})

	resp, err := client.GetQuantityDecimalPrecisions()
	require.Nil(t, err)

	prec := resp[0]

	require.Equal(t, Crypto, prec.InstrumentType)
	require.Equal(t, "AAVE/USD", prec.Symbol)
	require.Equal(t, 8, prec.Value)
	require.Equal(t, 6, prec.MinimumIncrementPrecision)
}

func TestGetQuantityDecimalPrecisionsError(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/instruments/quantity-decimal-precisions", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(401)
		fmt.Fprint(writer, tastyUnauthorizedError)
	})

	_, err := client.GetQuantityDecimalPrecisions()
	expectedUnauthorized(t, err)
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
	require.Equal(t, WarrantIT, war.InstrumentType)
	require.Equal(t, "XNAS", war.ListedMarket)
	require.Equal(t, "Nikola Corporation - Warrant expiring 6/3/2025", war.Description)
	require.False(t, war.IsClosingOnly)
	require.False(t, war.Active)
}

func TestGetWarrantsError(t *testing.T) {
	setup()
	defer teardown()

	symbol := "NKLAW"

	mux.HandleFunc("/instruments/warrants", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(401)
		fmt.Fprint(writer, tastyUnauthorizedError)
	})

	_, err := client.GetWarrants([]string{symbol})
	expectedUnauthorized(t, err)
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
	require.Equal(t, WarrantIT, war.InstrumentType)
	require.Equal(t, "XNAS", war.ListedMarket)
	require.Equal(t, "Nikola Corporation - Warrant expiring 6/3/2025", war.Description)
	require.False(t, war.IsClosingOnly)
	require.False(t, war.Active)
}
func TestGetWarrantError(t *testing.T) {
	setup()
	defer teardown()

	symbol := "NKLAW"

	mux.HandleFunc(fmt.Sprintf("/instruments/warrants/%s", symbol), func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(401)
		fmt.Fprint(writer, tastyUnauthorizedError)
	})

	_, err := client.GetWarrant(symbol)
	expectedUnauthorized(t, err)
}

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
