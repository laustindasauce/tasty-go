package tasty

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// Test data structures for OAuth2 request testing
type testResponse struct {
	Message string `json:"message"`
	Data    string `json:"data"`
}

type testParams struct {
	Filter string `url:"filter,omitempty"`
	Limit  int    `url:"limit,omitempty"`
}

// setupOAuth2Test creates a test environment for OAuth2 request testing
func setupOAuth2Test(t *testing.T) (*Client, *httptest.Server, *http.ServeMux) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	// Create OAuth2 client using internal function to bypass endpoint validation
	config := OAuth2Config{
		ClientID:     "test_client_id",
		ClientSecret: "test_client_secret",
		RedirectURI:  "http://localhost:8080/callback",
		Scopes:       []string{"read", "trade"},
		AuthURL:      server.URL + "/oauth/authorize",
		TokenURL:     server.URL + "/oauth/token",
		BaseURL:      server.URL,
	}

	client, err := NewClient(config, http.DefaultClient)
	require.NoError(t, err)

	// Override the base URL for testing
	client.baseURL = server.URL
	client.baseHost = strings.Split(server.URL, "/")[2]

	// Set up valid tokens in the token manager
	client.SetTokens("valid-access-token", "valid-refresh-token", 3600)

	return client, server, mux
}

// TestOAuth2Request_Success tests successful OAuth2 authenticated requests
func TestOAuth2Request_Success(t *testing.T) {
	client, server, mux := setupOAuth2Test(t)
	defer server.Close()

	// Set up test endpoint
	mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		// Verify Bearer token is present
		auth := r.Header.Get("Authorization")
		require.Equal(t, "Bearer valid-access-token", auth)

		// Verify content type
		require.Empty(t, r.Header.Get("Content-Type"))

		// Return success response
		response := testResponse{
			Message: "success",
			Data:    "test data",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	var result testResponse
	resp, err := client.request("GET", "/test", nil, nil, &result)

	require.Nil(t, err)
	require.NotNil(t, resp)
	require.Equal(t, "success", result.Message)
	require.Equal(t, "test data", result.Data)
}

// TestOAuth2Request_WithParams tests OAuth2 requests with query parameters
func TestOAuth2Request_WithParams(t *testing.T) {
	client, server, mux := setupOAuth2Test(t)
	defer server.Close()

	mux.HandleFunc("/test-params", func(w http.ResponseWriter, r *http.Request) {
		// Verify Bearer token
		auth := r.Header.Get("Authorization")
		require.Equal(t, "Bearer valid-access-token", auth)

		// Verify query parameters
		require.Equal(t, "active", r.URL.Query().Get("filter"))
		require.Equal(t, "10", r.URL.Query().Get("limit"))

		response := testResponse{Message: "success with params"}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	params := testParams{Filter: "active", Limit: 10}
	var result testResponse
	resp, err := client.request("GET", "/test-params", params, nil, &result)

	require.Nil(t, err)
	require.NotNil(t, resp)
	require.Equal(t, "success with params", result.Message)
}

// TestOAuth2Request_WithPayload tests OAuth2 requests with JSON payload
func TestOAuth2Request_WithPayload(t *testing.T) {
	client, server, mux := setupOAuth2Test(t)
	defer server.Close()

	mux.HandleFunc("/test-payload", func(w http.ResponseWriter, r *http.Request) {
		// Verify Bearer token
		auth := r.Header.Get("Authorization")
		require.Equal(t, "Bearer valid-access-token", auth)

		// Verify method
		require.Equal(t, "POST", r.Method)

		// Parse and verify payload
		var payload testResponse
		err := json.NewDecoder(r.Body).Decode(&payload)
		require.NoError(t, err)
		require.Equal(t, "test payload", payload.Message)

		response := testResponse{Message: "payload received"}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	payload := testResponse{Message: "test payload"}
	var result testResponse
	resp, err := client.request("POST", "/test-payload", nil, payload, &result)

	require.Nil(t, err)
	require.NotNil(t, resp)
	require.Equal(t, "payload received", result.Message)
}

// TestOAuth2Request_NoContent tests OAuth2 requests that return 204 No Content
func TestOAuth2Request_NoContent(t *testing.T) {
	client, server, mux := setupOAuth2Test(t)
	defer server.Close()

	mux.HandleFunc("/no-content", func(w http.ResponseWriter, r *http.Request) {
		// Verify Bearer token
		auth := r.Header.Get("Authorization")
		require.Equal(t, "Bearer valid-access-token", auth)

		w.WriteHeader(http.StatusNoContent)
	})

	resp, err := client.request("DELETE", "/no-content", nil, nil, nil)

	require.Nil(t, err)
	require.NotNil(t, resp)
	require.Equal(t, http.StatusNoContent, resp.StatusCode)
}

// TestOAuth2Request_TokenRefresh tests automatic token refresh on 401
func TestOAuth2Request_TokenRefresh(t *testing.T) {
	client, server, mux := setupOAuth2Test(t)
	defer server.Close()

	requestCount := 0

	// Set up token endpoint for refresh
	mux.HandleFunc("/oauth/token", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "POST", r.Method)
		require.Equal(t, "application/x-www-form-urlencoded", r.Header.Get("Content-Type"))

		// Return new tokens
		response := TokenResponse{
			AccessToken:  "new-access-token",
			RefreshToken: "new-refresh-token",
			TokenType:    "Bearer",
			ExpiresIn:    3600,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	// Set up test endpoint that returns 401 on first request, success on second
	mux.HandleFunc("/test-refresh", func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		auth := r.Header.Get("Authorization")

		if requestCount == 1 {
			// First request with expired token - return 401
			require.Equal(t, "Bearer valid-access-token", auth)
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "invalid_token",
			})
			return
		}

		// Second request with refreshed token - return success
		require.Equal(t, "Bearer new-access-token", auth)
		response := testResponse{Message: "success after refresh"}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	var result testResponse
	resp, err := client.request("GET", "/test-refresh", nil, nil, &result)

	require.Nil(t, err)
	require.NotNil(t, resp)
	require.Equal(t, "success after refresh", result.Message)
	require.Equal(t, 2, requestCount) // Should have made 2 requests

	// Verify token was updated
	token, tokenErr := client.GetTokenManager().GetAccessToken()
	require.NoError(t, tokenErr)
	require.Equal(t, "new-access-token", token)
}

// TestOAuth2Request_TokenRefreshFails tests behavior when token refresh fails
func TestOAuth2Request_TokenRefreshFails(t *testing.T) {
	client, server, mux := setupOAuth2Test(t)
	defer server.Close()

	// Set up token endpoint to return error
	mux.HandleFunc("/oauth/token", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error":             "invalid_grant",
			"error_description": "Refresh token expired",
		})
	})

	// Set up test endpoint that returns 401
	mux.HandleFunc("/test-refresh-fail", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "invalid_token",
		})
	})

	var result testResponse
	resp, err := client.request("GET", "/test-refresh-fail", nil, nil, &result)

	require.NotNil(t, err)
	require.NotNil(t, resp)
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

// TestCustomOAuth2Request_Success tests successful custom OAuth2 requests
func TestCustomOAuth2Request_Success(t *testing.T) {
	client, server, mux := setupOAuth2Test(t)
	defer server.Close()

	mux.HandleFunc("/custom/test", func(w http.ResponseWriter, r *http.Request) {
		// Verify Bearer token
		auth := r.Header.Get("Authorization")
		require.Equal(t, "Bearer valid-access-token", auth)

		response := testResponse{Message: "custom success"}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	var result testResponse
	resp, err := client.customRequest("GET", "/custom/test", nil, nil, &result)

	require.Nil(t, err)
	require.NotNil(t, resp)
	require.Equal(t, "custom success", result.Message)
}

// TestCustomOAuth2Request_TokenRefresh tests custom request with token refresh
func TestCustomOAuth2Request_TokenRefresh(t *testing.T) {
	client, server, mux := setupOAuth2Test(t)
	defer server.Close()

	requestCount := 0

	// Set up token endpoint for refresh
	mux.HandleFunc("/oauth/token", func(w http.ResponseWriter, r *http.Request) {
		response := TokenResponse{
			AccessToken:  "refreshed-access-token",
			RefreshToken: "refreshed-refresh-token",
			TokenType:    "Bearer",
			ExpiresIn:    3600,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	// Set up custom endpoint that returns 401 on first request
	mux.HandleFunc("/custom/refresh", func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		auth := r.Header.Get("Authorization")

		if requestCount == 1 {
			require.Equal(t, "Bearer valid-access-token", auth)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		require.Equal(t, "Bearer refreshed-access-token", auth)
		response := testResponse{Message: "custom refresh success"}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	var result testResponse
	resp, err := client.customRequest("GET", "/custom/refresh", nil, nil, &result)

	require.Nil(t, err)
	require.NotNil(t, resp)
	require.Equal(t, "custom refresh success", result.Message)
	require.Equal(t, 2, requestCount)
}

// TestOAuth2Request_ErrorResponses tests OAuth2 requests with various error responses
func TestOAuth2Request_ErrorResponses(t *testing.T) {
	client, server, mux := setupOAuth2Test(t)
	defer server.Close()

	errorCodes := []int{400, 403, 404, 415, 422, 500}

	for _, code := range errorCodes {
		path := fmt.Sprintf("/error-%d", code)
		mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
			// Verify Bearer token
			auth := r.Header.Get("Authorization")
			require.Equal(t, "Bearer valid-access-token", auth)

			w.WriteHeader(code)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error": map[string]interface{}{
					"code":    fmt.Sprintf("error_%d", code),
					"message": fmt.Sprintf("Test error %d", code),
				},
			})
		})

		var result testResponse
		resp, err := client.request("GET", path, nil, nil, &result)

		require.NotNil(t, err, "Expected error for status code %d", code)
		require.NotNil(t, resp, "Expected response for status code %d", code)
		require.Equal(t, code, resp.StatusCode, "Expected status code %d", code)
	}
}

// TestOAuth2Request_InvalidPayload tests OAuth2 requests with invalid JSON payload
func TestOAuth2Request_InvalidPayload(t *testing.T) {
	client, server, _ := setupOAuth2Test(t)
	defer server.Close()

	// Test with invalid payload that can't be marshaled
	invalid := math.Inf(1)
	resp, err := client.request("POST", "/test", nil, invalid, nil)

	require.NotNil(t, err)
	require.Nil(t, resp)
	require.Contains(t, err.Message, "Client Side Error")
}

// TestOAuth2Request_InvalidParams tests OAuth2 requests with invalid query parameters
func TestOAuth2Request_InvalidParams(t *testing.T) {
	client, server, _ := setupOAuth2Test(t)
	defer server.Close()

	// Test with invalid params that can't be encoded
	invalid := math.Inf(1)
	resp, err := client.request("GET", "/test", invalid, nil, nil)

	require.NotNil(t, err)
	require.Nil(t, resp)
	require.Contains(t, err.Message, "Client Side Error")
}

// TestOAuth2Request_InvalidResult tests OAuth2 requests with invalid result type
func TestOAuth2Request_InvalidResult(t *testing.T) {
	client, server, mux := setupOAuth2Test(t)
	defer server.Close()

	mux.HandleFunc("/invalid-result", func(w http.ResponseWriter, r *http.Request) {
		response := testResponse{Message: "test"}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	// Test with invalid result type that can't be decoded
	invalid := math.Inf(1)
	resp, err := client.request("GET", "/invalid-result", nil, nil, &invalid)

	require.NotNil(t, err)
	require.NotNil(t, resp)
	require.Contains(t, err.Message, "Client Side Error")
}

// TestOAuth2Request_NoOAuth2Client tests requests when OAuth2 client is not initialized
func TestOAuth2Request_NoTokens(t *testing.T) {
	// Create client without tokens
	config := OAuth2Config{
		ClientID:     "test_client_id",
		ClientSecret: "test_client_secret",
		RedirectURI:  "http://localhost:8080/callback",
		Scopes:       []string{"read", "trade"},
	}
	client, err := NewClient(config, http.DefaultClient)
	require.NoError(t, err)

	resp, err := client.request("GET", "/test", nil, nil, nil)

	require.NotNil(t, err)
	require.Nil(t, resp)
	// require.Equal(t, "oauth2_token_error", err.Code)
	// require.Contains(t, err.Message, "Failed to get access token")
}

// TestOAuth2Request_TokenError tests requests when token retrieval fails
func TestOAuth2Request_TokenError(t *testing.T) {
	client, server, _ := setupOAuth2Test(t)
	defer server.Close()

	// Clear tokens to simulate token error
	client.ClearTokens()

	resp, err := client.request("GET", "/test", nil, nil, nil)

	require.NotNil(t, err)
	require.Nil(t, resp)
	require.Equal(t, "oauth2_token_error", err.Code)
	require.Contains(t, err.Message, "Failed to get access token")
}

// TestOAuth2Request_NetworkError tests OAuth2 requests with network errors
func TestOAuth2Request_NetworkError(t *testing.T) {
	client, server, _ := setupOAuth2Test(t)
	server.Close() // Close server to simulate network error

	resp, err := client.request("GET", "/test", nil, nil, nil)

	require.NotNil(t, err)
	require.Nil(t, resp)
	require.Contains(t, err.Message, "Client Side Error")
}

// TestOAuth2Request_InvalidURL tests OAuth2 requests with invalid base URL
func TestOAuth2Request_InvalidURL(t *testing.T) {
	client, server, _ := setupOAuth2Test(t)
	defer server.Close()

	// Set invalid base URL
	client.baseURL = "invalid-url"

	resp, err := client.request("GET", "/test", nil, nil, nil)

	require.NotNil(t, err)
	require.Nil(t, resp)
	require.Contains(t, err.Message, "Client Side Error")
}
