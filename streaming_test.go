package tasty //nolint:testpackage // testing private field

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetQuoteStreamerTokens(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api-quote-tokens", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, quoteStreamerTokensResp)
	})

	resp, httpResp, err := client.GetQuoteStreamerTokens()
	require.Nil(t, err)
	require.NotNil(t, httpResp)

	require.Equal(t, "example-token-here", resp.Token)
	require.Equal(t, "wss://tasty-openapi-ws.dxfeed.com/realtime", resp.DXLinkURL)
	require.Equal(t, "api", resp.Level)
}

func TestGetQuoteStreamerTokensError(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api-quote-tokens", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(401)
		fmt.Fprint(writer, tastyUnauthorizedError)
	})

	_, httpResp, err := client.GetQuoteStreamerTokens()
	expectedUnauthorized(t, err)
	require.NotNil(t, httpResp)
}

const quoteStreamerTokensResp = `{
  "data": {
    "token": "example-token-here",
    "dxlink-url": "wss://tasty-openapi-ws.dxfeed.com/realtime",
    "level": "api"
  },
  "context": "/api-quote-tokens"
}`
