package tasty //nolint:testpackage // testing private field

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSymbolSearch(t *testing.T) {
	setup()
	defer teardown()

	symbolSearch := "apple"

	mux.HandleFunc(fmt.Sprintf("/symbols/search/%s", symbolSearch), func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, symbolSearchResp)
	})

	resp, httpResp, err := client.SymbolSearch(symbolSearch)
	require.Nil(t, err)
	require.NotNil(t, httpResp)

	require.Equal(t, 2, len(resp))

	aapl := resp[1]

	require.Equal(t, "AAPL", aapl.Symbol)
	require.Equal(t, "Apple Inc. - Common Stock", aapl.Description)
}

func TestSymbolSearchError(t *testing.T) {
	setup()
	defer teardown()

	symbolSearch := "apple"

	mux.HandleFunc(fmt.Sprintf("/symbols/search/%s", symbolSearch), func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(401)
		fmt.Fprint(writer, tastyUnauthorizedError)
	})

	_, httpResp, err := client.SymbolSearch(symbolSearch)
	expectedUnauthorized(t, err)
	require.NotNil(t, httpResp)
}

const symbolSearchResp = `{
    "data": {
        "items": [
            {
                "symbol": "APLE",
                "description": "Apple Hospitality REIT, Inc. Common Shares"
            },
            {
                "symbol": "AAPL",
                "description": "Apple Inc. - Common Stock"
            }
        ]
    }
}`
