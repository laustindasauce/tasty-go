package tasty //nolint:testpackage // testing private field

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var (
	watchlist = NewWatchlist{
		Name: "testing with / - odd chars",
		WatchlistEntries: []WatchlistEntry{
			{
				Symbol:         "AAPL",
				InstrumentType: EquityIT,
			},
		},
		GroupName:  "test",
		OrderIndex: 1,
	}

	newWatchlist = NewWatchlist{
		Name: "edited watchlist",
		WatchlistEntries: []WatchlistEntry{
			{
				Symbol:         "TSLA",
				InstrumentType: EquityIT,
			},
			{
				Symbol:         "AMZN",
				InstrumentType: EquityIT,
			},
		},
		GroupName:  "test",
		OrderIndex: 1,
	}
)

func TestGetMyWatchlists(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/watchlists", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, getWatchlistsResp)
	})

	resp, httpResp, err := client.GetMyWatchlists()
	require.Nil(t, err)
	require.NotNil(t, httpResp)
	require.Equal(t, 1, len(resp))

	w := resp[0]

	require.Equal(t, watchlist.Name, w.Name)
	require.Equal(t, len(watchlist.WatchlistEntries), len(w.WatchlistEntries))
	require.Equal(t, watchlist.WatchlistEntries[0].Symbol, w.WatchlistEntries[0].Symbol)
	require.Equal(t, watchlist.WatchlistEntries[0].InstrumentType, w.WatchlistEntries[0].InstrumentType)
	require.Equal(t, watchlist.GroupName, w.GroupName)
	require.Equal(t, watchlist.OrderIndex, w.OrderIndex)
}

func TestGetMyWatchlistsError(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/watchlists", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(401)
		fmt.Fprint(writer, tastyUnauthorizedError)
	})

	_, httpResp, err := client.GetMyWatchlists()
	expectedUnauthorized(t, err)
	require.NotNil(t, httpResp)
}

func TestGetMyWatchlist(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(fmt.Sprintf("/watchlists/%s", url.PathEscape(watchlist.Name)), func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, getWatchlistResp)
	})

	w, httpResp, err := client.GetMyWatchlist(watchlist.Name)
	require.Nil(t, err)
	require.NotNil(t, httpResp)

	require.Equal(t, watchlist.Name, w.Name)
	require.Equal(t, len(watchlist.WatchlistEntries), len(w.WatchlistEntries))
	require.Equal(t, watchlist.WatchlistEntries[0].Symbol, w.WatchlistEntries[0].Symbol)
	require.Equal(t, watchlist.WatchlistEntries[0].InstrumentType, w.WatchlistEntries[0].InstrumentType)
	require.Equal(t, watchlist.GroupName, w.GroupName)
	require.Equal(t, watchlist.OrderIndex, w.OrderIndex)
}

func TestGetMyWatchlistError(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(fmt.Sprintf("/watchlists/%s", url.PathEscape(watchlist.Name)), func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(401)
		fmt.Fprint(writer, tastyUnauthorizedError)
	})

	_, httpResp, err := client.GetMyWatchlist(watchlist.Name)
	expectedUnauthorized(t, err)
	require.NotNil(t, httpResp)
}

func TestCreateWatchlist(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/watchlists", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, createWatchlistResp)
	})

	w, httpResp, err := client.CreateWatchlist(watchlist)
	require.Nil(t, err)
	require.NotNil(t, httpResp)

	require.Equal(t, watchlist.Name, w.Name)
	require.Equal(t, len(watchlist.WatchlistEntries), len(w.WatchlistEntries))
	require.Equal(t, watchlist.WatchlistEntries[0].Symbol, w.WatchlistEntries[0].Symbol)
	require.Equal(t, watchlist.WatchlistEntries[0].InstrumentType, w.WatchlistEntries[0].InstrumentType)
	require.Equal(t, watchlist.GroupName, w.GroupName)
	require.Equal(t, watchlist.OrderIndex, w.OrderIndex)
}

func TestCreateWatchlistError(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/watchlists", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(401)
		fmt.Fprint(writer, tastyUnauthorizedError)
	})

	_, httpResp, err := client.CreateWatchlist(watchlist)
	expectedUnauthorized(t, err)
	require.NotNil(t, httpResp)
}

func TestEditWatchlist(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(fmt.Sprintf("/watchlists/%s", url.PathEscape(watchlist.Name)), func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, editedWatchlistResp)
	})

	w, httpResp, err := client.EditWatchlist(watchlist.Name, newWatchlist)
	require.Nil(t, err)
	require.NotNil(t, httpResp)

	require.Equal(t, newWatchlist.Name, w.Name)
	require.Equal(t, len(newWatchlist.WatchlistEntries), len(w.WatchlistEntries))
	require.Equal(t, newWatchlist.WatchlistEntries[0].Symbol, w.WatchlistEntries[0].Symbol)
	require.Equal(t, newWatchlist.WatchlistEntries[0].InstrumentType, w.WatchlistEntries[0].InstrumentType)
	require.Equal(t, newWatchlist.GroupName, w.GroupName)
	require.Equal(t, newWatchlist.OrderIndex, w.OrderIndex)
}

func TestEditWatchlistError(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(fmt.Sprintf("/watchlists/%s", url.PathEscape(watchlist.Name)), func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(401)
		fmt.Fprint(writer, tastyUnauthorizedError)
	})

	_, httpResp, err := client.EditWatchlist(watchlist.Name, newWatchlist)
	expectedUnauthorized(t, err)
	require.NotNil(t, httpResp)
}

func TestDeleteWatchlist(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(fmt.Sprintf("/watchlists/%s", url.PathEscape(newWatchlist.Name)), func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, deleteWatchlistResp)
	})

	w, httpResp, err := client.DeleteWatchlist(newWatchlist.Name)
	require.Nil(t, err)
	require.NotNil(t, httpResp)

	require.NotNil(t, w.ID)
	require.Equal(t, 162, *w.ID)
	require.Equal(t, newWatchlist.Name, w.Name)
	require.NotNil(t, w.UserID)
	require.Equal(t, 564, *w.UserID)
	require.Equal(t, "2023-06-14T16:34:05.909Z", w.CreatedAt.Format(time.RFC3339Nano))
	require.Equal(t, "2023-06-14T17:24:07.48Z", w.UpdatedAt.Format(time.RFC3339Nano))
	require.Equal(t, len(newWatchlist.WatchlistEntries), len(w.WatchlistEntries))
	require.Equal(t, newWatchlist.WatchlistEntries[0].Symbol, w.WatchlistEntries[0].Symbol)
	require.Equal(t, newWatchlist.WatchlistEntries[0].InstrumentType, w.WatchlistEntries[0].InstrumentType)
	require.Nil(t, w.CmsID)
	require.Equal(t, newWatchlist.GroupName, w.GroupName)
	require.Equal(t, newWatchlist.OrderIndex, w.OrderIndex)
}

func TestDeleteWatchlistError(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(fmt.Sprintf("/watchlists/%s", url.PathEscape(newWatchlist.Name)), func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(401)
		fmt.Fprint(writer, tastyUnauthorizedError)
	})

	_, httpResp, err := client.DeleteWatchlist(newWatchlist.Name)
	expectedUnauthorized(t, err)
	require.NotNil(t, httpResp)
}

func TestGetPairsWatchlists(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/pairs-watchlists", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, getPairsWatchlistsResp)
	})

	resp, httpResp, err := client.GetPairsWatchlists()
	require.Nil(t, err)
	require.NotNil(t, httpResp)
	require.Equal(t, 1, len(resp))

	w := resp[0]
	pe := w.PairsEquations[0]

	require.Equal(t, "Stocks", w.Name)
	require.Equal(t, 2, len(w.PairsEquations))
	require.Equal(t, "Buy", pe.LeftAction)
	require.Equal(t, "/ZN", pe.LeftSymbol)
	require.Equal(t, 2, pe.LeftQuantity)
	require.Equal(t, "Sell", pe.RightAction)
	require.Equal(t, "/ZB", pe.RightSymbol)
	require.Equal(t, 1, pe.RightQuantity)
	require.Equal(t, 4, w.OrderIndex)
}

func TestGetPairsWatchlistsError(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/pairs-watchlists", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(401)
		fmt.Fprint(writer, tastyUnauthorizedError)
	})

	_, httpResp, err := client.GetPairsWatchlists()
	expectedUnauthorized(t, err)
	require.NotNil(t, httpResp)
}

func TestGetPairsWatchlist(t *testing.T) {
	setup()
	defer teardown()

	name := "Stocks"

	mux.HandleFunc(fmt.Sprintf("/pairs-watchlists/%s", name), func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, getPairsWatchlistResp)
	})

	w, httpResp, err := client.GetPairsWatchlist(name)
	require.Nil(t, err)
	require.NotNil(t, httpResp)

	pe := w.PairsEquations[0]

	require.Equal(t, name, w.Name)
	require.Equal(t, 2, len(w.PairsEquations))
	require.Equal(t, "Buy", pe.LeftAction)
	require.Equal(t, "/ZN", pe.LeftSymbol)
	require.Equal(t, 2, pe.LeftQuantity)
	require.Equal(t, "Sell", pe.RightAction)
	require.Equal(t, "/ZB", pe.RightSymbol)
	require.Equal(t, 1, pe.RightQuantity)
	require.Equal(t, 4, w.OrderIndex)
}

func TestGetPairsWatchlistError(t *testing.T) {
	setup()
	defer teardown()

	name := "Stocks"

	mux.HandleFunc(fmt.Sprintf("/pairs-watchlists/%s", name), func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(401)
		fmt.Fprint(writer, tastyUnauthorizedError)
	})

	_, httpResp, err := client.GetPairsWatchlist(name)
	expectedUnauthorized(t, err)
	require.NotNil(t, httpResp)
}

func TestGetPublicWatchlists(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/public-watchlists", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, getPublicWatchlistsResp)
	})

	countsOnly := false

	resp, httpResp, err := client.GetPublicWatchlists(countsOnly)
	require.Nil(t, err)
	require.NotNil(t, httpResp)
	require.Equal(t, 1, len(resp))

	w := resp[0]
	pe := w.WatchlistEntries[0]

	require.Equal(t, "All Earnings", w.Name)
	require.Equal(t, 5, len(w.WatchlistEntries))
	require.Equal(t, "GIS", pe.Symbol)
	require.Equal(t, EquityIT, pe.InstrumentType)
	require.Equal(t, "Earnings", w.GroupName)
	require.Equal(t, 100, w.OrderIndex)
}

func TestGetPublicWatchlistsError(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/public-watchlists", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(401)
		fmt.Fprint(writer, tastyUnauthorizedError)
	})

	countsOnly := false

	_, httpResp, err := client.GetPublicWatchlists(countsOnly)
	expectedUnauthorized(t, err)
	require.NotNil(t, httpResp)
}

func TestGetPublicWatchlistsCountsError(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/public-watchlists", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(401)
		fmt.Fprint(writer, tastyUnauthorizedError)
	})

	countsOnly := true

	_, httpResp, err := client.GetPublicWatchlists(countsOnly)
	expectedUnauthorized(t, err)
	require.NotNil(t, httpResp)
}

func TestGetPublicWatchlist(t *testing.T) {
	setup()
	defer teardown()

	name := "High Options Volume"

	mux.HandleFunc(fmt.Sprintf("/public-watchlists/%s", url.PathEscape(name)), func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, getPublicWatchlistResp)
	})

	w, httpResp, err := client.GetPublicWatchlist(name)
	require.Nil(t, err)
	require.NotNil(t, httpResp)

	pe := w.WatchlistEntries[0]

	require.Equal(t, name, w.Name)
	require.Equal(t, 2, len(w.WatchlistEntries))
	require.Equal(t, "SPY", pe.Symbol)
	require.Equal(t, EquityIT, pe.InstrumentType)
	require.Equal(t, "Liquidity", w.GroupName)
	require.Equal(t, 100, w.OrderIndex)
}

func TestGetPublicWatchlistError(t *testing.T) {
	setup()
	defer teardown()

	name := "High Options Volume"

	mux.HandleFunc(fmt.Sprintf("/public-watchlists/%s", url.PathEscape(name)), func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(401)
		fmt.Fprint(writer, tastyUnauthorizedError)
	})

	_, httpResp, err := client.GetPublicWatchlist(name)
	expectedUnauthorized(t, err)
	require.NotNil(t, httpResp)
}

const getWatchlistsResp = `{
  "data": {
    "items": [
      {
        "name": "testing with / - odd chars",
		"user-id": 12345,
        "watchlist-entries": [
          { "symbol": "AAPL", "instrument-type": "Equity" }
        ],
        "group-name": "test",
        "order-index": 1
      }
    ]
  },
  "context": "/watchlists"
}`

const getWatchlistResp = `{
  "data": {
    "name": "testing with / - odd chars",
	"user-id": 12345,
    "watchlist-entries": [{ "symbol": "AAPL", "instrument-type": "Equity" }],
    "group-name": "test",
    "order-index": 1
  }
}`

const createWatchlistResp = `{
  "data": {
    "name": "testing with / - odd chars",
    "watchlist-entries": [{ "symbol": "AAPL", "instrument-type": "Equity" }],
    "group-name": "test",
    "order-index": 1
  },
  "context": "/watchlists"
}`

const editedWatchlistResp = `{
  "data": {
    "name": "edited watchlist",
    "watchlist-entries": [
      { "symbol": "TSLA", "instrument-type": "Equity" },
      { "symbol": "AMZN", "instrument-type": "Equity" }
    ],
    "group-name": "test",
    "order-index": 1
  }
}`

const deleteWatchlistResp = `{
  "id": 162,
  "name": "edited watchlist",
  "user_id": 564,
  "created_at": "2023-06-14T16:34:05.909Z",
  "updated_at": "2023-06-14T17:24:07.480Z",
  "watchlist_entries": [
    { "symbol": "TSLA", "instrument-type": "Equity" },
    { "symbol": "AMZN", "instrument-type": "Equity" }
  ],
  "cms_id": null,
  "group_name": "test",
  "order_index": 1
}`

const getPairsWatchlistsResp = `{
  "data": {
    "items": [
      {
        "name": "Stocks",
		"order-index": 4,
		"pairs-equations": [
			{
				"left-action": "Buy",
				"left-symbol": "/ZN",
				"left-quantity": 2,
				"right-action": "Sell",
				"right-symbol": "/ZB",
				"right-quantity": 1
			},
			{
				"left-action": "Buy",
				"left-symbol": "/GC",
				"left-quantity": 1,
				"right-action": "Sell",
				"right-symbol": "/SI",
				"right-quantity": 2
			}
		],
		"created-at": "2020-11-17T22:17:27.424Z",
		"updated-at": "2020-11-17T22:17:27.424Z"
      }
    ]
  }
}`

const getPairsWatchlistResp = `{
  "data": {
	"name": "Stocks",
	"order-index": 4,
    "pairs-equations": [
      {
        "left-action": "Buy",
        "left-symbol": "/ZN",
        "left-quantity": 2,
        "right-action": "Sell",
        "right-symbol": "/ZB",
        "right-quantity": 1
      },
      {
        "left-action": "Buy",
        "left-symbol": "/GC",
        "left-quantity": 1,
        "right-action": "Sell",
        "right-symbol": "/SI",
        "right-quantity": 2
      }
    ],
	"created-at": "2020-11-17T22:17:27.424Z",
	"updated-at": "2020-11-17T22:17:27.424Z"
  }
}`

const getPublicWatchlistsResp = `{
  "data": {
    "items": [
		{
			"id": "216375",
			"name": "All Earnings",
			"group-name": "Earnings",
			"user-id": -1,
			"order-index": 100,
			"watchlist-entries": [
				{
					"symbol": "GIS",
					"instrument-type": "Equity"
				},
				{
					"symbol": "LEN",
					"instrument-type": "Equity"
				},
				{
					"symbol": "DRI",
					"instrument-type": "Equity"
				},
				{
					"symbol": "FDS",
					"instrument-type": "Equity"
				},
				{
					"symbol": "FDX",
					"instrument-type": "Equity"
				}
			]
		}
    ]
  }
}`

const getPublicWatchlistResp = `{
  "data": {
    "id": "69136",
	"name": "High Options Volume",
	"group-name": "Liquidity",
	"user-id": -1,
	"order-index": 100,
    "watchlist-entries": [
      {
        "symbol": "SPY",
        "instrument-type": "Equity"
      },
      {
        "symbol": "QQQ",
        "instrument-type": "Equity"
      }
    ],
    "group-name": "Liquidity",
    "order-index": 100
  }
}`
