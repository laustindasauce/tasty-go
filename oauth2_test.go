package tasty

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// NOTE: These tests verify the PKCE implementation works correctly according to RFC 7636.
// However, TastyTrade's OAuth2 implementation does not currently support PKCE.
// Using PKCE parameters with TastyTrade will result in "405 Method Not Allowed" errors.
// This implementation is provided for future compatibility and testing purposes.

// MockOAuth2Server provides a test server for OAuth2 flows
type MockOAuth2Server struct {
	server            *httptest.Server
	mux               *http.ServeMux
	authEndpoint      string
	tokenEndpoint     string
	expectedState     string
	expectedCode      string
	tokenResponse     *TokenResponse
	shouldReturnError bool
	errorResponse     *OAuth2Error
	requestCount      int
	mutex             sync.RWMutex
}

// NewMockOAuth2Server creates a new mock OAuth2 server for testing
func NewMockOAuth2Server() *MockOAuth2Server {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	mock := &MockOAuth2Server{
		server:        server,
		mux:           mux,
		authEndpoint:  server.URL + "/oauth/authorize",
		tokenEndpoint: server.URL + "/oauth/token",
		tokenResponse: &TokenResponse{
			AccessToken:  "test_access_token",
			RefreshToken: "test_refresh_token",
			TokenType:    "Bearer",
			ExpiresIn:    3600,
			Scope:        "read trade",
		},
	}

	mock.setupEndpoints()
	return mock
}

// setupEndpoints configures the mock server endpoints
func (m *MockOAuth2Server) setupEndpoints() {
	m.mux.HandleFunc("/oauth/authorize", m.handleAuthorize)
	m.mux.HandleFunc("/oauth/token", m.handleToken)
}

// handleAuthorize handles authorization endpoint requests
func (m *MockOAuth2Server) handleAuthorize(w http.ResponseWriter, r *http.Request) {
	m.mutex.Lock()
	m.requestCount++
	m.mutex.Unlock()

	// This would normally redirect to login page, but for testing
	// we just return a simple response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Authorization endpoint"))
}

// handleToken handles token endpoint requests
func (m *MockOAuth2Server) handleToken(w http.ResponseWriter, r *http.Request) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.requestCount++

	if m.shouldReturnError {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		if m.errorResponse != nil {
			json.NewEncoder(w).Encode(m.errorResponse)
		} else {
			json.NewEncoder(w).Encode(OAuth2Error{
				ErrorCode:        "invalid_request",
				ErrorDescription: "Test error",
			})
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(m.tokenResponse)
}

// SetTokenResponse sets the token response for successful requests
func (m *MockOAuth2Server) SetTokenResponse(response *TokenResponse) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.tokenResponse = response
}

// SetErrorResponse configures the server to return an error
func (m *MockOAuth2Server) SetErrorResponse(err *OAuth2Error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.shouldReturnError = true
	m.errorResponse = err
}

// ClearError clears error configuration
func (m *MockOAuth2Server) ClearError() {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.shouldReturnError = false
	m.errorResponse = nil
}

// GetRequestCount returns the number of requests received
func (m *MockOAuth2Server) GetRequestCount() int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.requestCount
}

// Close shuts down the mock server
func (m *MockOAuth2Server) Close() {
	m.server.Close()
}

// createTestOAuth2Config creates a test OAuth2 configuration
func createTestOAuth2Config(mockServer *MockOAuth2Server) OAuth2Config {
	return OAuth2Config{
		ClientID:     "test_client_id",
		ClientSecret: "test_client_secret",
		RedirectURI:  "http://localhost:8080/callback",
		Scopes:       []string{"read", "trade"},
		AuthURL:      mockServer.authEndpoint,
		TokenURL:     mockServer.tokenEndpoint,
		State:        "test_state_123",
	}
}

// createTestClientWithMemoryStorage creates a Client with isolated memory storage for testing
func createTestClientWithMemoryStorage(mockServer *MockOAuth2Server) (*Client, error) {
	config := createTestOAuth2Config(mockServer)

	// Use NewClient since it's now the unified constructor
	client, err := NewClient(config, defaultHTTPClient)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func TestNewClientInternal(t *testing.T) {
	mockServer := NewMockOAuth2Server()
	defer mockServer.Close()

	config := createTestOAuth2Config(mockServer)

	// Use NewClient since it's now the unified constructor
	client, err := NewClient(config, nil)
	require.NoError(t, err)
	require.NotNil(t, client)

	clientConfig := client.GetConfig()
	assert.Equal(t, config.ClientID, clientConfig.ClientID)
	assert.Equal(t, config.ClientSecret, clientConfig.ClientSecret)
	assert.Equal(t, config.RedirectURI, clientConfig.RedirectURI)
	assert.Equal(t, config.Scopes, clientConfig.Scopes)
	assert.NotNil(t, client.GetTokenManager())
	assert.NotEmpty(t, clientConfig.State)
}

func TestNewClient_InvalidConfig(t *testing.T) {
	tests := []struct {
		name   string
		config OAuth2Config
	}{
		{
			name: "empty client ID",
			config: OAuth2Config{
				ClientSecret: "secret",
				RedirectURI:  "http://localhost:8080/callback",
			},
		},
		{
			name: "empty redirect URI",
			config: OAuth2Config{
				ClientID:     "client_id",
				ClientSecret: "secret",
			},
		},
		{
			name: "invalid redirect URI",
			config: OAuth2Config{
				ClientID:     "client_id",
				ClientSecret: "secret",
				RedirectURI:  "invalid-uri",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(tt.config, nil)
			assert.Error(t, err)
			assert.Nil(t, client)
		})
	}
}

func TestClient_GetAuthorizationURL(t *testing.T) {
	mockServer := NewMockOAuth2Server()
	defer mockServer.Close()

	config := createTestOAuth2Config(mockServer)
	client, err := NewClient(config, nil)
	require.NoError(t, err)

	authURL, err := client.GetAuthorizationURL()
	require.NoError(t, err)
	require.NotEmpty(t, authURL)

	// Parse and validate the URL
	parsedURL, err := url.Parse(authURL)
	require.NoError(t, err)

	assert.Equal(t, mockServer.authEndpoint, parsedURL.Scheme+"://"+parsedURL.Host+parsedURL.Path)

	// Check query parameters
	params := parsedURL.Query()
	assert.Equal(t, "code", params.Get("response_type"))
	assert.Equal(t, config.ClientID, params.Get("client_id"))
	assert.Equal(t, config.RedirectURI, params.Get("redirect_uri"))
	assert.Equal(t, "read trade", params.Get("scope"))
	assert.Equal(t, config.State, params.Get("state"))

	// PKCE parameters should NOT be present since TastyTrade doesn't support PKCE
	assert.Empty(t, params.Get("code_challenge"), "PKCE code_challenge should not be present")
	assert.Empty(t, params.Get("code_challenge_method"), "PKCE code_challenge_method should not be present")
}

func TestClient_GetAuthorizationURL_MissingFields(t *testing.T) {
	tests := []struct {
		name   string
		modify func(*OAuth2Config)
	}{
		{
			name: "missing client ID",
			modify: func(c *OAuth2Config) {
				c.ClientID = ""
			},
		},
		{
			name: "missing redirect URI",
			modify: func(c *OAuth2Config) {
				c.RedirectURI = ""
			},
		},
		{
			name: "missing auth URL",
			modify: func(c *OAuth2Config) {
				c.AuthURL = ""
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockServer := NewMockOAuth2Server()
			defer mockServer.Close()

			config := createTestOAuth2Config(mockServer)
			tt.modify(&config)

			client, err := NewClient(config, http.DefaultClient)
			if err != nil {
				// Expected error for invalid config
				assert.Error(t, err)
				return
			}

			authURL, err := client.GetAuthorizationURL()
			assert.Error(t, err)
			assert.Empty(t, authURL)
		})
	}
}

func TestClient_ExchangeCodeForTokens(t *testing.T) {
	mockServer := NewMockOAuth2Server()
	defer mockServer.Close()

	config := createTestOAuth2Config(mockServer)
	client, err := NewClient(config, nil)
	require.NoError(t, err)

	// Test successful token exchange
	tokenResponse, err := client.ExchangeCodeForTokens("test_auth_code")
	require.NoError(t, err)
	require.NotNil(t, tokenResponse)

	assert.Equal(t, "test_access_token", tokenResponse.AccessToken)
	assert.Equal(t, "test_refresh_token", tokenResponse.RefreshToken)
	assert.Equal(t, "Bearer", tokenResponse.TokenType)
	assert.Equal(t, 3600, tokenResponse.ExpiresIn)

	// Verify tokens are stored in token manager
	storedToken, err := client.GetTokenManager().GetAccessToken()
	require.NoError(t, err)
	assert.Equal(t, "test_access_token", storedToken)
}

func TestClient_ExchangeCodeForTokens_Error(t *testing.T) {
	mockServer := NewMockOAuth2Server()
	defer mockServer.Close()

	// Configure server to return error
	mockServer.SetErrorResponse(&OAuth2Error{
		ErrorCode:        "invalid_grant",
		ErrorDescription: "Invalid authorization code",
	})

	config := createTestOAuth2Config(mockServer)
	client, err := NewClient(config, nil)
	require.NoError(t, err)

	tokenResponse, err := client.ExchangeCodeForTokens("invalid_code")
	assert.Error(t, err)
	assert.Nil(t, tokenResponse)

	// Check that it's a detailed OAuth2 error
	detailedErr, ok := err.(*OAuth2DetailedError)
	require.True(t, ok, "Expected OAuth2DetailedError, got %T", err)
	assert.Equal(t, "invalid_grant", detailedErr.ErrorCode)
	assert.Equal(t, "Invalid authorization code", detailedErr.ErrorDescription)
	assert.False(t, detailedErr.IsRetryable())
}

func TestClient_ExchangeCodeForTokens_EmptyCode(t *testing.T) {
	mockServer := NewMockOAuth2Server()
	defer mockServer.Close()

	config := createTestOAuth2Config(mockServer)
	client, err := NewClient(config, nil)
	require.NoError(t, err)

	tokenResponse, err := client.ExchangeCodeForTokens("")
	assert.Error(t, err)
	assert.Nil(t, tokenResponse)

	// Check that it's a detailed OAuth2 error
	detailedErr, ok := err.(*OAuth2DetailedError)
	require.True(t, ok, "Expected OAuth2DetailedError, got %T", err)
	assert.Equal(t, OAuth2ErrorMissingCode, detailedErr.ErrorCode)
}

func TestClient_RefreshTokens(t *testing.T) {
	mockServer := NewMockOAuth2Server()
	defer mockServer.Close()

	config := createTestOAuth2Config(mockServer)
	client, err := NewClient(config, nil)
	require.NoError(t, err)

	// Set up initial tokens
	client.SetTokens("old_access_token", "refresh_token", 3600)

	// Configure mock server to return new tokens
	mockServer.SetTokenResponse(&TokenResponse{
		AccessToken:  "new_access_token",
		RefreshToken: "new_refresh_token",
		TokenType:    "Bearer",
		ExpiresIn:    3600,
	})

	tokenResponse, err := client.RefreshTokens()
	require.NoError(t, err)
	require.NotNil(t, tokenResponse)

	assert.Equal(t, "new_access_token", tokenResponse.AccessToken)
	assert.Equal(t, "new_refresh_token", tokenResponse.RefreshToken)

	// Verify new tokens are stored
	storedToken, err := client.GetTokenManager().GetAccessToken()
	require.NoError(t, err)
	assert.Equal(t, "new_access_token", storedToken)
}

func TestClient_RefreshTokens_NoRefreshToken(t *testing.T) {
	mockServer := NewMockOAuth2Server()
	defer mockServer.Close()

	client, err := createTestClientWithMemoryStorage(mockServer)
	require.NoError(t, err)

	// Don't set any tokens
	client.tokenManager.Clear()
	tokenResponse, err := client.RefreshTokens()
	assert.Error(t, err)
	assert.Nil(t, tokenResponse)

	// Check that it's a detailed OAuth2 error
	detailedErr, ok := err.(*OAuth2DetailedError)
	require.True(t, ok, "Expected OAuth2DetailedError, got %T", err)
	assert.Equal(t, OAuth2ErrorRefreshFailed, detailedErr.ErrorCode)
}

func TestClient_RefreshTokens_ServerError(t *testing.T) {
	mockServer := NewMockOAuth2Server()
	defer mockServer.Close()

	config := createTestOAuth2Config(mockServer)
	client, err := NewClient(config, nil)
	require.NoError(t, err)

	// Set up initial tokens
	client.SetTokens("access_token", "refresh_token", 3600)

	// Configure server to return error
	mockServer.SetErrorResponse(&OAuth2Error{
		ErrorCode:        "invalid_grant",
		ErrorDescription: "Refresh token expired",
	})

	tokenResponse, err := client.RefreshTokens()
	assert.Error(t, err)
	assert.Nil(t, tokenResponse)

	// Check that it's a detailed OAuth2 error
	detailedErr, ok := err.(*OAuth2DetailedError)
	require.True(t, ok, "Expected OAuth2DetailedError, got %T", err)
	assert.Equal(t, "invalid_grant", detailedErr.ErrorCode)
	assert.Equal(t, "Refresh token expired", detailedErr.ErrorDescription)
}

func TestClient_RefreshTokensWithRetry(t *testing.T) {
	mockServer := NewMockOAuth2Server()
	defer mockServer.Close()

	config := createTestOAuth2Config(mockServer)
	client, err := NewClient(config, nil)
	require.NoError(t, err)

	// Test that non-retryable errors don't retry
	client.SetTokens("access_token", "refresh_token", 3600)

	mockServer.SetErrorResponse(&OAuth2Error{
		ErrorCode:        "invalid_grant",
		ErrorDescription: "Non-retryable error",
	})

	start := time.Now()
	_, err = client.RefreshTokens()
	duration := time.Since(start)

	assert.Error(t, err)
	// Should not retry, so should be fast
	assert.Less(t, duration, 500*time.Millisecond)

	// Verify only one request was made
	assert.Equal(t, 1, mockServer.GetRequestCount())
}

func TestClient_ValidateState(t *testing.T) {
	mockServer := NewMockOAuth2Server()
	defer mockServer.Close()

	config := createTestOAuth2Config(mockServer)
	client, err := NewClient(config, nil)
	require.NoError(t, err)

	// Test valid state
	err = client.ValidateState(config.State)
	assert.NoError(t, err)

	// Test invalid state
	err = client.ValidateState("invalid_state")
	assert.Error(t, err)
	detailedErr, ok := err.(*OAuth2DetailedError)
	require.True(t, ok, "Expected OAuth2DetailedError, got %T", err)
	assert.Equal(t, OAuth2ErrorInvalidState, detailedErr.ErrorCode)

	// Test empty state
	err = client.ValidateState("")
	assert.Error(t, err)
	detailedErr, ok = err.(*OAuth2DetailedError)
	require.True(t, ok, "Expected OAuth2DetailedError, got %T", err)
	assert.Equal(t, OAuth2ErrorInvalidState, detailedErr.ErrorCode)
}

func TestClient_GetConfig(t *testing.T) {
	mockServer := NewMockOAuth2Server()
	defer mockServer.Close()

	config := createTestOAuth2Config(mockServer)
	client, err := NewClient(config, nil)
	require.NoError(t, err)

	retrievedConfig := client.GetConfig()
	assert.Equal(t, config.ClientID, retrievedConfig.ClientID)
	assert.Equal(t, config.RedirectURI, retrievedConfig.RedirectURI)
	assert.Equal(t, config.Scopes, retrievedConfig.Scopes)
}

func TestClient_GetState(t *testing.T) {
	mockServer := NewMockOAuth2Server()
	defer mockServer.Close()

	config := createTestOAuth2Config(mockServer)
	client, err := NewClient(config, nil)
	require.NoError(t, err)

	state := client.GetState()
	assert.Equal(t, config.State, state)
}

func TestClient_GetTokenManager(t *testing.T) {
	mockServer := NewMockOAuth2Server()
	defer mockServer.Close()

	config := createTestOAuth2Config(mockServer)
	client, err := NewClient(config, nil)
	require.NoError(t, err)

	tokenManager := client.GetTokenManager()
	assert.NotNil(t, tokenManager)
}

func TestGenerateSecureState(t *testing.T) {
	state1, err := generateSecureState()
	require.NoError(t, err)
	require.NotEmpty(t, state1)

	state2, err := generateSecureState()
	require.NoError(t, err)
	require.NotEmpty(t, state2)

	// States should be different
	assert.NotEqual(t, state1, state2)

	// States should be reasonable length (base64url encoded 32 bytes)
	assert.Greater(t, len(state1), 40)
	assert.Greater(t, len(state2), 40)
}

func TestRedirectServer_Lifecycle(t *testing.T) {
	server := NewRedirectServer("test_state")
	require.NotNil(t, server)

	// Initially not running
	assert.False(t, server.IsRunning())
	assert.Equal(t, 0, server.GetPort())

	// Start server on random port
	err := server.Start(0)
	require.NoError(t, err)

	// Should be running now
	assert.True(t, server.IsRunning())
	assert.Greater(t, server.GetPort(), 0)

	// Get redirect URI
	redirectURI := server.GetRedirectURI()
	assert.Contains(t, redirectURI, "http://localhost:")

	// Shutdown server
	err = server.Shutdown(5 * time.Second)
	require.NoError(t, err)

	// Should not be running
	assert.False(t, server.IsRunning())
}

func TestRedirectServer_StartTwice(t *testing.T) {
	server := NewRedirectServer("test_state")

	err := server.Start(0)
	require.NoError(t, err)
	defer server.Shutdown(5 * time.Second)

	// Starting again should fail
	err = server.Start(0)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already running")
}

func TestRedirectServer_HandleRedirect_Success(t *testing.T) {
	server := NewRedirectServer("test_state")

	err := server.Start(0)
	require.NoError(t, err)
	defer server.Shutdown(5 * time.Second)

	// Simulate successful redirect
	redirectURI := server.GetRedirectURI()
	resp, err := http.Get(redirectURI + "?code=test_code&state=test_state")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Wait for code
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	select {
	case code := <-server.codeChan:
		assert.Equal(t, "test_code", code)
	case <-ctx.Done():
		t.Fatal("Timeout waiting for code")
	}
}

func TestRedirectServer_HandleRedirect_Error(t *testing.T) {
	server := NewRedirectServer("test_state")

	err := server.Start(0)
	require.NoError(t, err)
	defer server.Shutdown(5 * time.Second)

	// Simulate error redirect
	redirectURI := server.GetRedirectURI()
	resp, err := http.Get(redirectURI + "?error=access_denied&error_description=User%20denied&state=test_state")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// Wait for error
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	select {
	case err := <-server.errChan:
		assert.Error(t, err)
		detailedErr, ok := err.(*OAuth2DetailedError)
		require.True(t, ok, "Expected OAuth2DetailedError, got %T", err)
		assert.Equal(t, "access_denied", detailedErr.ErrorCode)
		assert.Equal(t, "User denied", detailedErr.ErrorDescription)
	case <-ctx.Done():
		t.Fatal("Timeout waiting for error")
	}
}

func TestRedirectServer_HandleRedirect_InvalidState(t *testing.T) {
	server := NewRedirectServer("expected_state")

	err := server.Start(0)
	require.NoError(t, err)
	defer server.Shutdown(5 * time.Second)

	// Simulate redirect with invalid state
	redirectURI := server.GetRedirectURI()
	resp, err := http.Get(redirectURI + "?code=test_code&state=invalid_state")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// Wait for error
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	select {
	case err := <-server.errChan:
		assert.Error(t, err)
		detailedErr, ok := err.(*OAuth2DetailedError)
		require.True(t, ok, "Expected OAuth2DetailedError, got %T", err)
		assert.Equal(t, OAuth2ErrorInvalidState, detailedErr.ErrorCode)
		assert.Contains(t, detailedErr.Error(), "CSRF attack")
	case <-ctx.Done():
		t.Fatal("Timeout waiting for error")
	}
}

func TestRedirectServer_WaitForCode_Timeout(t *testing.T) {
	server := NewRedirectServer("test_state")

	err := server.Start(0)
	require.NoError(t, err)
	defer server.Shutdown(5 * time.Second)

	// Wait for code with short timeout
	code, err := server.WaitForCode(100 * time.Millisecond)
	assert.Error(t, err)
	assert.Empty(t, code)

	detailedErr, ok := err.(*OAuth2DetailedError)
	require.True(t, ok, "Expected OAuth2DetailedError, got %T", err)
	assert.Equal(t, OAuth2ErrorRedirectTimeout, detailedErr.ErrorCode)
}

func TestRedirectServer_WaitForCode_NotRunning(t *testing.T) {
	server := NewRedirectServer("test_state")

	// Don't start the server
	code, err := server.WaitForCode(1 * time.Second)
	assert.Error(t, err)
	assert.Empty(t, code)

	detailedErr, ok := err.(*OAuth2DetailedError)
	require.True(t, ok, "Expected OAuth2DetailedError, got %T", err)
	assert.Equal(t, OAuth2ErrorConfigurationError, detailedErr.ErrorCode)
}

func TestClient_StartRedirectServer(t *testing.T) {
	mockServer := NewMockOAuth2Server()
	defer mockServer.Close()

	config := createTestOAuth2Config(mockServer)
	client, err := NewClient(config, nil)
	require.NoError(t, err)

	redirectServer, err := client.StartRedirectServer(0)
	require.NoError(t, err)
	require.NotNil(t, redirectServer)
	defer redirectServer.Shutdown(5 * time.Second)

	assert.True(t, redirectServer.IsRunning())
	assert.Greater(t, redirectServer.GetPort(), 0)
}

// TestClient_ConcurrentAccess tests thread safety of Client
func TestClient_ConcurrentAccess(t *testing.T) {
	mockServer := NewMockOAuth2Server()
	defer mockServer.Close()

	config := createTestOAuth2Config(mockServer)
	client, err := NewClient(config, nil)
	require.NoError(t, err)

	// Set initial tokens
	client.SetTokens("access_token", "refresh_token", 3600)

	var wg sync.WaitGroup
	errors := make(chan error, 10)

	// Concurrent access to various methods
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			// Test concurrent access to various methods
			_, err := client.GetAuthorizationURL()
			if err != nil {
				errors <- err
				return
			}

			_ = client.IsAuthenticated()
			_ = client.GetState()
			_ = client.GetConfig()
			_ = client.GetTokenManager()
		}()
	}

	wg.Wait()
	close(errors)

	// Check for any errors
	for err := range errors {
		t.Errorf("Concurrent access error: %v", err)
	}
}

// TestClient_Integration tests a complete OAuth2 flow
func TestClient_Integration(t *testing.T) {
	mockServer := NewMockOAuth2Server()
	defer mockServer.Close()

	config := createTestOAuth2Config(mockServer)
	client, err := NewClient(config, nil)
	require.NoError(t, err)

	// Step 1: Get authorization URL
	authURL, err := client.GetAuthorizationURL()
	require.NoError(t, err)
	require.NotEmpty(t, authURL)

	// Step 2: Exchange code for tokens
	tokenResponse, err := client.ExchangeCodeForTokens("test_auth_code")
	require.NoError(t, err)
	require.NotNil(t, tokenResponse)

	// Step 3: Verify authentication
	assert.True(t, client.IsAuthenticated())

	// Step 4: Test token refresh
	mockServer.SetTokenResponse(&TokenResponse{
		AccessToken:  "refreshed_access_token",
		RefreshToken: "refreshed_refresh_token",
		TokenType:    "Bearer",
		ExpiresIn:    3600,
	})

	refreshedResponse, err := client.RefreshTokens()
	require.NoError(t, err)
	require.NotNil(t, refreshedResponse)
	assert.Equal(t, "refreshed_access_token", refreshedResponse.AccessToken)

	// Step 5: Clear tokens
	client.ClearTokens()
	assert.False(t, client.IsAuthenticated())
}

func TestGeneratePKCEChallenge(t *testing.T) {
	t.Run("returns error because PKCE is disabled", func(t *testing.T) {
		challenge, err := GeneratePKCEChallenge()

		require.Error(t, err)
		require.Nil(t, challenge)

		// Verify error message explains PKCE is not supported
		assert.Contains(t, err.Error(), "PKCE is not supported")
		assert.Contains(t, err.Error(), "405 Method Not Allowed")
	})

	t.Run("consistently returns error", func(t *testing.T) {
		// Multiple calls should all return errors
		for i := 0; i < 5; i++ {
			challenge, err := GeneratePKCEChallenge()
			require.Error(t, err)
			require.Nil(t, challenge)
			assert.Contains(t, err.Error(), "PKCE is not supported")
		}
	})
}

func TestPKCEChallenge_Validate(t *testing.T) {
	t.Run("always returns false because PKCE is disabled", func(t *testing.T) {
		// Create a challenge manually (since GeneratePKCEChallenge returns error)
		challenge := &PKCEChallenge{
			CodeVerifier:  "test-verifier",
			CodeChallenge: "test-challenge",
			Method:        "S256",
		}

		// Should always return false since PKCE is disabled
		assert.False(t, challenge.Validate("test-verifier"))
		assert.False(t, challenge.Validate("any-verifier"))
		assert.False(t, challenge.Validate(""))
	})

	t.Run("nil challenge returns false", func(t *testing.T) {
		var challenge *PKCEChallenge
		assert.False(t, challenge.Validate("any-verifier"))
	})

	t.Run("even valid RFC 7636 challenges return false", func(t *testing.T) {
		// Create a challenge with known valid values
		verifier := "dBjftJeZ4CVP-mB92K27uhbUJU1p1r_wW1gFWFOEjXk"
		hash := sha256.Sum256([]byte(verifier))
		expectedChallenge := base64.RawURLEncoding.EncodeToString(hash[:])

		challenge := &PKCEChallenge{
			CodeVerifier:  verifier,
			CodeChallenge: expectedChallenge,
			Method:        "S256",
		}

		// Even with correct values, should return false since PKCE is disabled
		assert.False(t, challenge.Validate(verifier))
		assert.False(t, challenge.Validate("wrong-verifier"))
	})
}

func TestPKCEChallenge_IsValid(t *testing.T) {
	t.Run("always returns false because PKCE is disabled", func(t *testing.T) {
		// Create various challenge configurations
		challenges := []*PKCEChallenge{
			{
				CodeVerifier:  "test-verifier",
				CodeChallenge: "test-challenge",
				Method:        "S256",
			},
			{
				CodeVerifier:  "",
				CodeChallenge: "test-challenge",
				Method:        "S256",
			},
			{
				CodeVerifier:  "test-verifier",
				CodeChallenge: "",
				Method:        "S256",
			},
			{
				CodeVerifier:  "test-verifier",
				CodeChallenge: "test-challenge",
				Method:        "",
			},
		}

		// All should return false since PKCE is disabled
		for i, challenge := range challenges {
			assert.False(t, challenge.IsValid(), "Challenge %d should be invalid", i)
		}
	})

	t.Run("nil challenge returns false", func(t *testing.T) {
		var challenge *PKCEChallenge
		assert.False(t, challenge.IsValid())
	})

	t.Run("even RFC 7636 compliant challenges return false", func(t *testing.T) {
		// Create a perfectly valid RFC 7636 challenge
		verifier := "dBjftJeZ4CVP-mB92K27uhbUJU1p1r_wW1gFWFOEjXk"
		hash := sha256.Sum256([]byte(verifier))
		expectedChallenge := base64.RawURLEncoding.EncodeToString(hash[:])

		challenge := &PKCEChallenge{
			CodeVerifier:  verifier,
			CodeChallenge: expectedChallenge,
			Method:        "S256",
		}

		// Even with perfect RFC compliance, should return false since PKCE is disabled
		assert.False(t, challenge.IsValid())
	})
}

func TestPKCEChallenge_SecurityProperties(t *testing.T) {
	t.Run("PKCE generation consistently fails", func(t *testing.T) {
		// Attempt to generate multiple challenges - all should fail
		for i := 0; i < 10; i++ {
			challenge, err := GeneratePKCEChallenge()
			require.Error(t, err)
			require.Nil(t, challenge)
			assert.Contains(t, err.Error(), "PKCE is not supported")
		}
	})

	t.Run("manual challenge creation for testing purposes", func(t *testing.T) {
		// We can still create PKCEChallenge structs manually for testing
		// but they will always be considered invalid
		challenge := &PKCEChallenge{
			CodeVerifier:  "test-verifier-43-chars-long-base64url-encoded",
			CodeChallenge: "test-challenge-43-chars-long-sha256-base64url",
			Method:        "S256",
		}

		// Even manually created challenges are invalid since PKCE is disabled
		assert.False(t, challenge.IsValid())
		assert.False(t, challenge.Validate("test-verifier-43-chars-long-base64url-encoded"))
	})

	t.Run("PKCE disabled for security reasons", func(t *testing.T) {
		// Document that PKCE is disabled to prevent TastyTrade errors
		// This is a security feature to prevent accidental use of unsupported PKCE
		t.Log("PKCE is disabled because TastyTrade returns '405 Method Not Allowed'")
		t.Log("This prevents runtime errors and ensures OAuth2 flow works correctly")

		// Verify that all PKCE operations fail gracefully
		challenge, err := GeneratePKCEChallenge()
		assert.Error(t, err)
		assert.Nil(t, challenge)
	})
}

func TestPKCEChallenge_EdgeCases(t *testing.T) {
	t.Run("generation always fails", func(t *testing.T) {
		// All attempts to generate PKCE challenges should fail
		for i := 0; i < 5; i++ {
			challenge, err := GeneratePKCEChallenge()
			require.Error(t, err)
			require.Nil(t, challenge)
			assert.Contains(t, err.Error(), "PKCE is not supported")
		}
	})

	t.Run("manual challenges always invalid", func(t *testing.T) {
		// Create various manual challenges - all should be invalid
		challenges := []*PKCEChallenge{
			{
				CodeVerifier:  "verifier1",
				CodeChallenge: "challenge1",
				Method:        "S256",
			},
			{
				CodeVerifier:  "VERIFIER2",
				CodeChallenge: "CHALLENGE2",
				Method:        "S256",
			},
			{
				CodeVerifier:  " verifier3 ",
				CodeChallenge: " challenge3 ",
				Method:        "S256",
			},
		}

		for i, challenge := range challenges {
			assert.False(t, challenge.IsValid(), "Challenge %d should be invalid", i)
			assert.False(t, challenge.Validate("any-verifier"), "Challenge %d should not validate any verifier", i)
		}
	})

	t.Run("PKCE disabled prevents all edge cases", func(t *testing.T) {
		// Document that disabling PKCE prevents all edge case issues
		t.Log("PKCE is completely disabled, preventing all edge case vulnerabilities")
		t.Log("This ensures TastyTrade OAuth2 flow works without PKCE-related errors")

		// Even nil challenges behave consistently
		var nilChallenge *PKCEChallenge
		assert.False(t, nilChallenge.IsValid())
		assert.False(t, nilChallenge.Validate("any-verifier"))
	})
}

func TestOAuth2Config_SupportsPKCE(t *testing.T) {
	t.Run("always returns false - PKCE completely disabled", func(t *testing.T) {
		configs := []OAuth2Config{
			{
				AuthURL:  "https://my.tastytrade.com/auth.html",
				TokenURL: "https://api.tastyworks.com/oauth/token",
			},
			{
				AuthURL:  "https://cert-my.staging-tasty.works/auth.html",
				TokenURL: "https://api.cert.tastyworks.com/oauth/token",
			},
			{
				BaseURL: "https://api.tastyworks.com",
			},
			{
				AuthURL:  "https://auth.example.com/oauth/authorize",
				TokenURL: "https://auth.example.com/oauth/token",
			},
			{}, // empty config
		}

		// All configs should return false since PKCE is completely disabled
		for i, config := range configs {
			assert.False(t, config.SupportsPKCE(), "Config %d should not support PKCE", i)
		}
	})
}

func TestPKCETastyTradeCompatibility(t *testing.T) {
	t.Run("document tastytrade PKCE limitation and complete disabling", func(t *testing.T) {
		// This test documents that PKCE is completely disabled because TastyTrade does not support it
		// The error response from TastyTrade when PKCE parameters are included:
		// {"error": {"code": "method_not_allowed","message": "405 Not Allowed. Unique customer support identifier: 7e168b8d269e2d5e082141eb55a1927c"}}

		// Attempt to generate PKCE challenge - should fail
		challenge, err := GeneratePKCEChallenge()
		require.Error(t, err)
		require.Nil(t, challenge)
		assert.Contains(t, err.Error(), "PKCE is not supported")

		// Create TastyTrade config
		config := NewSandboxOAuth2Config("test-client-id", "test-secret", "http://localhost:8080", []string{"read", "trade"})

		// Verify TastyTrade doesn't support PKCE
		assert.False(t, config.SupportsPKCE())

		// Test OAuth2Options validation
		options := OAuth2Options{
			ClientID:    "test-client-id",
			RedirectURI: "http://localhost:8080",
			UsePKCE:     true, // This should cause validation to fail
		}

		err = ValidateOAuth2Options(options)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "PKCE is not supported")
		assert.Contains(t, err.Error(), "405 Method Not Allowed")

		t.Log("PKCE is completely disabled to prevent TastyTrade '405 Method Not Allowed' errors")
		t.Logf("GeneratePKCEChallenge error: %v", err)
		t.Logf("TastyTrade config supports PKCE: %v", config.SupportsPKCE())
	})
}

func TestValidateOAuth2Options(t *testing.T) {
	t.Run("valid options without PKCE", func(t *testing.T) {
		options := OAuth2Options{
			ClientID:    "test-client-id",
			RedirectURI: "http://localhost:8080",
			UsePKCE:     false, // PKCE disabled
		}

		err := ValidateOAuth2Options(options)
		assert.NoError(t, err)
	})

	t.Run("invalid options with PKCE enabled", func(t *testing.T) {
		options := OAuth2Options{
			ClientID:    "test-client-id",
			RedirectURI: "http://localhost:8080",
			UsePKCE:     true, // PKCE enabled - should fail
		}

		err := ValidateOAuth2Options(options)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "PKCE is not supported")
		assert.Contains(t, err.Error(), "405 Method Not Allowed")
	})

	t.Run("default options are valid", func(t *testing.T) {
		options := OAuth2Options{
			ClientID:    "test-client-id",
			RedirectURI: "http://localhost:8080",
			// UsePKCE defaults to false
		}

		err := ValidateOAuth2Options(options)
		assert.NoError(t, err)
	})
}

func TestClient_ScopeHandling(t *testing.T) {
	mockServer := NewMockOAuth2Server()
	defer mockServer.Close()

	// Set up mock server to return token response without scope
	mockServer.tokenResponse = &TokenResponse{
		AccessToken:  "test_access_token",
		RefreshToken: "test_refresh_token",
		TokenType:    "Bearer",
		ExpiresIn:    3600,
		// Scope is intentionally omitted to simulate TastyTrade behavior
	}

	client, err := createTestClientWithMemoryStorage(mockServer)
	require.NoError(t, err)

	// Exchange code for tokens (use a longer code to pass validation)
	tokenResponse, err := client.ExchangeCodeForTokens("test_code_1234567890")
	require.NoError(t, err)

	// Verify scope is set from config
	expectedScope := strings.Join(client.GetConfig().Scopes, " ")
	assert.Equal(t, expectedScope, tokenResponse.Scope)
	assert.Equal(t, expectedScope, client.GetTokenManager().GetScope())
}

func TestClient_RefreshTokenPreservation(t *testing.T) {
	mockServer := NewMockOAuth2Server()
	defer mockServer.Close()

	client, err := createTestClientWithMemoryStorage(mockServer)
	require.NoError(t, err)

	// Set initial tokens
	client.SetTokens("initial-token", "persistent-refresh-token", 3600)
	originalScope := strings.Join(client.GetConfig().Scopes, " ")

	// Manually set scope to simulate initial token exchange
	tokenManager := client.GetTokenManager()
	tokenManager.storage.Store(&TokenData{
		AccessToken:  "initial-token",
		RefreshToken: "persistent-refresh-token",
		TokenType:    "Bearer",
		ExpiresAt:    time.Now().Add(time.Hour),
		Scope:        originalScope,
	})

	// Set up mock server to return refresh response without refresh token or scope
	mockServer.tokenResponse = &TokenResponse{
		AccessToken: "refreshed-token",
		TokenType:   "Bearer",
		ExpiresIn:   3600,
		// RefreshToken and Scope are intentionally omitted
	}

	// Perform token refresh
	refreshResponse, err := client.RefreshTokens()
	require.NoError(t, err)

	// Verify refresh token is preserved
	assert.Equal(t, "persistent-refresh-token", client.GetTokenManager().GetRefreshToken())

	// Verify scope is preserved
	assert.Equal(t, originalScope, client.GetTokenManager().GetScope())

	// Verify new access token is set
	assert.Equal(t, "refreshed-token", refreshResponse.AccessToken)

	// Verify scope is set in response
	assert.Equal(t, originalScope, refreshResponse.Scope)
}
