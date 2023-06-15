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

	resp, err := client.SymbolSearch(symbolSearch)
	require.Nil(t, err)

	require.Equal(t, 2, len(resp))

	aapl := resp[1]

	require.Equal(t, "AAPL", aapl.Symbol)
	require.Equal(t, "Apple Inc. - Common Stock", aapl.Description)
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
