package tasty //nolint:testpackage // testing private field

import (
	"fmt"
	"net/http"
	"testing"

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
	require.Equal(t, Crypto, btc.InstrumentType)
	require.Equal(t, "Bitcoin", btc.ShortDescription)
	require.Equal(t, "Bitcoin to USD", btc.Description)
	require.False(t, btc.IsClosingOnly)
	require.True(t, btc.Active)
	require.Equal(t, StringToFloat32(0.01), btc.TickSize)
	require.Equal(t, "BTC/USD:CXTALP", btc.StreamerSymbol)

	venueSymbol := btc.DestinationVenueSymbols[0]

	require.Equal(t, 71, venueSymbol.ID)
	require.Equal(t, "BTC", venueSymbol.Symbol)
	require.Equal(t, "CITADEL_CRYPTOCURRENCY", venueSymbol.DestinationVenue)
	require.Equal(t, 8, venueSymbol.MaxQuantityPrecision)
	require.Equal(t, 8, venueSymbol.MaxPricePrecision)
	require.True(t, venueSymbol.Routable)
}

func TestGetCryptocurrenciesError(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/instruments/cryptocurrencies", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(401)
		fmt.Fprint(writer, tastyUnauthorizedError)
	})

	_, err := client.GetCryptocurrencies([]string{"BTC/USD", "ETH/USD"})
	expectedUnauthorized(t, err)
}

func TestGetCryptocurrency(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/instruments/cryptocurrencies/BTC/USD", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, cryptoResp)
	})

	btc, err := client.GetCryptocurrency(Bitcoin)
	require.Nil(t, err)

	require.Equal(t, 1, btc.ID)
	require.Equal(t, "BTC/USD", btc.Symbol)
	require.Equal(t, Crypto, btc.InstrumentType)
	require.Equal(t, "Bitcoin", btc.ShortDescription)
	require.Equal(t, "Bitcoin to USD", btc.Description)
	require.False(t, btc.IsClosingOnly)
	require.True(t, btc.Active)
	require.Equal(t, StringToFloat32(0.01), btc.TickSize)
	require.Equal(t, "BTC/USD:CXTALP", btc.StreamerSymbol)

	venueSymbol := btc.DestinationVenueSymbols[0]

	require.Equal(t, 71, venueSymbol.ID)
	require.Equal(t, "BTC", venueSymbol.Symbol)
	require.Equal(t, "CITADEL_CRYPTOCURRENCY", venueSymbol.DestinationVenue)
	require.Equal(t, 8, venueSymbol.MaxQuantityPrecision)
	require.Equal(t, 8, venueSymbol.MaxPricePrecision)
	require.True(t, venueSymbol.Routable)
}

func TestGetCryptocurrencyError(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/instruments/cryptocurrencies/BTC/USD", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(401)
		fmt.Fprint(writer, tastyUnauthorizedError)
	})

	_, err := client.GetCryptocurrency(Bitcoin)
	expectedUnauthorized(t, err)
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

const cryptoResp = `{
    "data": {
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
            },
            {
                "id": 10,
                "symbol": "BTCUSD",
                "destination-venue": "DV Chain",
                "max-quantity-precision": 8,
                "max-price-precision": 2,
                "routable": true
            },
            {
                "id": 11,
                "symbol": "XBT",
                "destination-venue": "Jane Street",
                "max-quantity-precision": 8,
                "max-price-precision": 7,
                "routable": true
            },
            {
                "id": 12,
                "symbol": "BTC",
                "destination-venue": "Zero Hash",
                "routable": false
            },
            {
                "id": 75,
                "symbol": "BTC/USD",
                "destination-venue": "CBOE_DIGITAL",
                "max-quantity-precision": 6,
                "max-price-precision": 1,
                "routable": true
            }
        ]
    }
}`
