package tasty

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOAuth2Client_SetTokens(t *testing.T) {
	config := NewProductionOAuth2Config("test-client-id", "test-client-secret", "http://localhost:8080/callback", []string{"read", "trade"})
	
	// Create client with memory storage for isolated testing
	tokenManager := NewMemoryTokenManager()
	client := &OAuth2Client{
		config:       config,
		tokenManager: tokenManager,
		httpClient:   defaultHTTPClient,
		pkce:         nil,
	}
	tokenManager.SetRefreshCallback(client.refreshTokensInternal)

	accessToken := "test-access-token"
	refreshToken := "test-refresh-token"
	expiresIn := 3600

	// Set tokens
	client.SetTokens(accessToken, refreshToken, expiresIn)

	// Verify tokens are set correctly
	assert.True(t, client.HasValidToken())
	assert.True(t, client.HasRefreshToken())
	assert.False(t, client.IsTokenExpired())

	// Verify token manager has the tokens
	storedAccessToken, err := client.tokenManager.GetAccessToken()
	require.NoError(t, err)
	assert.Equal(t, accessToken, storedAccessToken)
	assert.Equal(t, refreshToken, client.tokenManager.GetRefreshToken())
}

func TestOAuth2Client_SetTokensFromResponse(t *testing.T) {
	config := NewProductionOAuth2Config("test-client-id", "test-client-secret", "http://localhost:8080/callback", []string{"read", "trade"})
	
	// Create client with memory storage for isolated testing
	tokenManager := NewMemoryTokenManager()
	client := &OAuth2Client{
		config:       config,
		tokenManager: tokenManager,
		httpClient:   defaultHTTPClient,
		pkce:         nil,
	}
	tokenManager.SetRefreshCallback(client.refreshTokensInternal)

	tokenResponse := &TokenResponse{
		AccessToken:  "test-access-token",
		RefreshToken: "test-refresh-token",
		TokenType:    "Bearer",
		ExpiresIn:    3600,
		Scope:        "read trade",
	}

	// Set tokens from response
	client.SetTokensFromResponse(tokenResponse)

	// Verify tokens are set correctly
	assert.True(t, client.HasValidToken())
	assert.True(t, client.HasRefreshToken())
	assert.False(t, client.IsTokenExpired())

	// Verify token manager has the tokens
	storedAccessToken, err := client.tokenManager.GetAccessToken()
	require.NoError(t, err)
	assert.Equal(t, tokenResponse.AccessToken, storedAccessToken)
	assert.Equal(t, tokenResponse.RefreshToken, client.tokenManager.GetRefreshToken())
	assert.Equal(t, tokenResponse.TokenType, client.tokenManager.GetTokenType())
	assert.Equal(t, tokenResponse.Scope, client.tokenManager.GetScope())
}

func TestOAuth2Client_SetTokensFromResponse_Nil(t *testing.T) {
	config := NewProductionOAuth2Config("test-client-id", "test-client-secret", "http://localhost:8080/callback", []string{"read", "trade"})
	
	// Create client with memory storage for isolated testing
	tokenManager := NewMemoryTokenManager()
	client := &OAuth2Client{
		config:       config,
		tokenManager: tokenManager,
		httpClient:   defaultHTTPClient,
		pkce:         nil,
	}
	tokenManager.SetRefreshCallback(client.refreshTokensInternal)

	// Should not panic with nil response
	client.SetTokensFromResponse(nil)

	// Should not have valid tokens
	assert.False(t, client.HasValidToken())
	assert.False(t, client.HasRefreshToken())
}

func TestOAuth2Client_HasValidToken(t *testing.T) {
	config := NewProductionOAuth2Config("test-client-id", "test-client-secret", "http://localhost:8080/callback", []string{"read", "trade"})
	
	// Create client with memory storage for isolated testing
	tokenManager := NewMemoryTokenManager()
	client := &OAuth2Client{
		config:       config,
		tokenManager: tokenManager,
		httpClient:   defaultHTTPClient,
		pkce:         nil,
	}
	tokenManager.SetRefreshCallback(client.refreshTokensInternal)

	// Initially should not have valid token
	assert.False(t, client.HasValidToken())

	// Set valid token
	client.SetTokens("access-token", "refresh-token", 3600)
	assert.True(t, client.HasValidToken())

	// Set expired token
	client.SetTokens("expired-token", "refresh-token", -10)
	assert.False(t, client.HasValidToken())

	// Set token expiring soon (within 30 second buffer)
	client.SetTokens("expiring-token", "refresh-token", 10)
	assert.False(t, client.HasValidToken())
}

func TestOAuth2Client_HasRefreshToken(t *testing.T) {
	config := NewProductionOAuth2Config("test-client-id", "test-client-secret", "http://localhost:8080/callback", []string{"read", "trade"})
	
	// Create client with memory storage for isolated testing
	tokenManager := NewMemoryTokenManager()
	client := &OAuth2Client{
		config:       config,
		tokenManager: tokenManager,
		httpClient:   defaultHTTPClient,
		pkce:         nil,
	}
	tokenManager.SetRefreshCallback(client.refreshTokensInternal)

	// Initially should not have refresh token
	assert.False(t, client.HasRefreshToken())

	// Set tokens with refresh token
	client.SetTokens("access-token", "refresh-token", 3600)
	assert.True(t, client.HasRefreshToken())

	// Set tokens without refresh token
	client.SetTokens("access-token", "", 3600)
	assert.False(t, client.HasRefreshToken())
}

func TestOAuth2Client_GetTokenExpiration(t *testing.T) {
	config := NewProductionOAuth2Config("test-client-id", "test-client-secret", "http://localhost:8080/callback", []string{"read", "trade"})
	
	// Create client with memory storage for isolated testing
	tokenManager := NewMemoryTokenManager()
	client := &OAuth2Client{
		config:       config,
		tokenManager: tokenManager,
		httpClient:   defaultHTTPClient,
		pkce:         nil,
	}
	tokenManager.SetRefreshCallback(client.refreshTokensInternal)

	// Set token with 1 hour expiry
	beforeSet := time.Now()
	client.SetTokens("access-token", "refresh-token", 3600)
	afterSet := time.Now()

	expiration := client.GetTokenExpiration()
	expectedExpiration := beforeSet.Add(time.Hour)

	// Should be within reasonable range
	assert.True(t, expiration.After(expectedExpiration.Add(-time.Second)))
	assert.True(t, expiration.Before(afterSet.Add(time.Hour).Add(time.Second)))
}

func TestOAuth2Client_GetTimeUntilExpiry(t *testing.T) {
	config := NewProductionOAuth2Config("test-client-id", "test-client-secret", "http://localhost:8080/callback", []string{"read", "trade"})
	
	// Create client with memory storage for isolated testing
	tokenManager := NewMemoryTokenManager()
	client := &OAuth2Client{
		config:       config,
		tokenManager: tokenManager,
		httpClient:   defaultHTTPClient,
		pkce:         nil,
	}
	tokenManager.SetRefreshCallback(client.refreshTokensInternal)

	// No token should return 0
	assert.Equal(t, time.Duration(0), client.GetTimeUntilExpiry())

	// Set token with 1 hour expiry
	client.SetTokens("access-token", "refresh-token", 3600)
	timeUntilExpiry := client.GetTimeUntilExpiry()

	// Should be close to 1 hour (allowing for small timing differences)
	assert.True(t, timeUntilExpiry > 59*time.Minute)
	assert.True(t, timeUntilExpiry <= time.Hour)

	// Set expired token
	client.SetTokens("expired-token", "refresh-token", -10)
	assert.Equal(t, time.Duration(0), client.GetTimeUntilExpiry())
}

func TestOAuth2Client_IsTokenExpired(t *testing.T) {
	config := NewProductionOAuth2Config("test-client-id", "test-client-secret", "http://localhost:8080/callback", []string{"read", "trade"})
	
	// Create client with memory storage for isolated testing
	tokenManager := NewMemoryTokenManager()
	client := &OAuth2Client{
		config:       config,
		tokenManager: tokenManager,
		httpClient:   defaultHTTPClient,
		pkce:         nil,
	}
	tokenManager.SetRefreshCallback(client.refreshTokensInternal)

	// No token should be considered expired
	assert.True(t, client.IsTokenExpired())

	// Set valid token
	client.SetTokens("access-token", "refresh-token", 3600)
	assert.False(t, client.IsTokenExpired())

	// Set expired token
	client.SetTokens("expired-token", "refresh-token", -10)
	assert.True(t, client.IsTokenExpired())

	// Set token expiring soon (within 30 second buffer)
	client.SetTokens("expiring-token", "refresh-token", 10)
	assert.True(t, client.IsTokenExpired())
}

func TestNewOAuth2ClientWithTokens(t *testing.T) {
	config := NewProductionOAuth2Config("test-client-id", "test-client-secret", "http://localhost:8080/callback", []string{"read", "trade"})
	
	accessToken := "test-access-token"
	refreshToken := "test-refresh-token"
	expiresIn := 3600

	client, err := NewOAuth2ClientWithTokens(config, accessToken, refreshToken, expiresIn, nil)
	require.NoError(t, err)
	require.NotNil(t, client)

	// Verify client is in OAuth2 mode
	assert.Equal(t, AuthModeOAuth2, client.GetAuthMode())
	assert.True(t, client.IsOAuth2Mode())

	// Verify tokens are set
	assert.True(t, client.IsAuthenticated())
	assert.True(t, client.HasValidToken())
	assert.True(t, client.HasRefreshToken())
	assert.False(t, client.IsTokenExpired())

	// Verify token values
	storedAccessToken, err := client.oauth2Client.tokenManager.GetAccessToken()
	require.NoError(t, err)
	assert.Equal(t, accessToken, storedAccessToken)
	assert.Equal(t, refreshToken, client.oauth2Client.tokenManager.GetRefreshToken())
}

func TestNewOAuth2ClientWithTokenResponse(t *testing.T) {
	config := NewProductionOAuth2Config("test-client-id", "test-client-secret", "http://localhost:8080/callback", []string{"read", "trade"})
	
	tokenResponse := &TokenResponse{
		AccessToken:  "test-access-token",
		RefreshToken: "test-refresh-token",
		TokenType:    "Bearer",
		ExpiresIn:    3600,
		Scope:        "read trade",
	}

	client, err := NewOAuth2ClientWithTokenResponse(config, tokenResponse, nil)
	require.NoError(t, err)
	require.NotNil(t, client)

	// Verify client is in OAuth2 mode
	assert.Equal(t, AuthModeOAuth2, client.GetAuthMode())
	assert.True(t, client.IsOAuth2Mode())

	// Verify tokens are set
	assert.True(t, client.IsAuthenticated())
	assert.True(t, client.HasValidToken())
	assert.True(t, client.HasRefreshToken())
	assert.False(t, client.IsTokenExpired())

	// Verify token values
	storedAccessToken, err := client.oauth2Client.tokenManager.GetAccessToken()
	require.NoError(t, err)
	assert.Equal(t, tokenResponse.AccessToken, storedAccessToken)
	assert.Equal(t, tokenResponse.RefreshToken, client.oauth2Client.tokenManager.GetRefreshToken())
	assert.Equal(t, tokenResponse.TokenType, client.oauth2Client.tokenManager.GetTokenType())
	assert.Equal(t, tokenResponse.Scope, client.oauth2Client.tokenManager.GetScope())
}

func TestNewCertOAuth2ClientWithTokens(t *testing.T) {
	config := NewSandboxOAuth2Config("test-client-id", "test-client-secret", "http://localhost:8080/callback", []string{"read", "trade"})
	
	accessToken := "test-access-token"
	refreshToken := "test-refresh-token"
	expiresIn := 3600

	client, err := NewCertOAuth2ClientWithTokens(config, accessToken, refreshToken, expiresIn, nil)
	require.NoError(t, err)
	require.NotNil(t, client)

	// Verify client is in OAuth2 mode and using cert environment
	assert.Equal(t, AuthModeOAuth2, client.GetAuthMode())
	assert.True(t, client.IsOAuth2Mode())
	assert.Equal(t, apiCertBaseURL, client.baseURL)

	// Verify tokens are set
	assert.True(t, client.IsAuthenticated())
	assert.True(t, client.HasValidToken())
	assert.True(t, client.HasRefreshToken())
}

func TestNewCertOAuth2ClientWithTokenResponse(t *testing.T) {
	config := NewSandboxOAuth2Config("test-client-id", "test-client-secret", "http://localhost:8080/callback", []string{"read", "trade"})
	
	tokenResponse := &TokenResponse{
		AccessToken:  "test-access-token",
		RefreshToken: "test-refresh-token",
		TokenType:    "Bearer",
		ExpiresIn:    3600,
		Scope:        "read trade",
	}

	client, err := NewCertOAuth2ClientWithTokenResponse(config, tokenResponse, nil)
	require.NoError(t, err)
	require.NotNil(t, client)

	// Verify client is in OAuth2 mode and using cert environment
	assert.Equal(t, AuthModeOAuth2, client.GetAuthMode())
	assert.True(t, client.IsOAuth2Mode())
	assert.Equal(t, apiCertBaseURL, client.baseURL)

	// Verify tokens are set
	assert.True(t, client.IsAuthenticated())
	assert.True(t, client.HasValidToken())
	assert.True(t, client.HasRefreshToken())
}

func TestClient_SetTokens(t *testing.T) {
	config := NewProductionOAuth2Config("test-client-id", "test-client-secret", "http://localhost:8080/callback", []string{"read", "trade"})
	client, err := NewOAuth2Client(config, nil)
	require.NoError(t, err)

	accessToken := "test-access-token"
	refreshToken := "test-refresh-token"
	expiresIn := 3600

	// Set tokens
	err = client.SetTokens(accessToken, refreshToken, expiresIn)
	require.NoError(t, err)

	// Verify tokens are set
	assert.True(t, client.IsAuthenticated())
	assert.True(t, client.HasValidToken())
	assert.True(t, client.HasRefreshToken())
}

func TestClient_SetTokens_SessionMode(t *testing.T) {
	client := NewClient(nil)

	// Should return error in session mode
	err := client.SetTokens("access-token", "refresh-token", 3600)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "OAuth2 mode")
}

func TestClient_SetTokensFromResponse(t *testing.T) {
	config := NewProductionOAuth2Config("test-client-id", "test-client-secret", "http://localhost:8080/callback", []string{"read", "trade"})
	client, err := NewOAuth2Client(config, nil)
	require.NoError(t, err)

	tokenResponse := &TokenResponse{
		AccessToken:  "test-access-token",
		RefreshToken: "test-refresh-token",
		TokenType:    "Bearer",
		ExpiresIn:    3600,
		Scope:        "read trade",
	}

	// Set tokens from response
	err = client.SetTokensFromResponse(tokenResponse)
	require.NoError(t, err)

	// Verify tokens are set
	assert.True(t, client.IsAuthenticated())
	assert.True(t, client.HasValidToken())
	assert.True(t, client.HasRefreshToken())
}

func TestClient_SetTokensFromResponse_SessionMode(t *testing.T) {
	client := NewClient(nil)

	tokenResponse := &TokenResponse{
		AccessToken:  "test-access-token",
		RefreshToken: "test-refresh-token",
		TokenType:    "Bearer",
		ExpiresIn:    3600,
	}

	// Should return error in session mode
	err := client.SetTokensFromResponse(tokenResponse)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "OAuth2 mode")
}

func TestClient_HasValidToken(t *testing.T) {
	// Test OAuth2 mode with isolated storage
	config := NewProductionOAuth2Config("test-client-id", "test-client-secret", "http://localhost:8080/callback", []string{"read", "trade"})
	
	// Create OAuth2Client with memory storage for isolated testing
	tokenManager := NewMemoryTokenManager()
	oauth2ClientInternal := &OAuth2Client{
		config:       config,
		tokenManager: tokenManager,
		httpClient:   defaultHTTPClient,
		pkce:         nil,
	}
	tokenManager.SetRefreshCallback(oauth2ClientInternal.refreshTokensInternal)
	
	// Create Client with the isolated OAuth2Client
	oauth2Client := &Client{
		httpClient:   defaultHTTPClient,
		baseURL:      config.BaseURL,
		oauth2Client: oauth2ClientInternal,
		authMode:     AuthModeOAuth2,
	}

	// Initially should not have valid token
	assert.False(t, oauth2Client.HasValidToken())

	// Set valid token
	oauth2Client.SetTokens("access-token", "refresh-token", 3600)
	assert.True(t, oauth2Client.HasValidToken())

	// Test session mode
	sessionClient := NewClient(nil)
	assert.False(t, sessionClient.HasValidToken())
}

func TestClient_HasRefreshToken(t *testing.T) {
	// Test OAuth2 mode with isolated storage
	config := NewProductionOAuth2Config("test-client-id", "test-client-secret", "http://localhost:8080/callback", []string{"read", "trade"})
	
	// Create OAuth2Client with memory storage for isolated testing
	tokenManager := NewMemoryTokenManager()
	oauth2ClientInternal := &OAuth2Client{
		config:       config,
		tokenManager: tokenManager,
		httpClient:   defaultHTTPClient,
		pkce:         nil,
	}
	tokenManager.SetRefreshCallback(oauth2ClientInternal.refreshTokensInternal)
	
	// Create Client with the isolated OAuth2Client
	oauth2Client := &Client{
		httpClient:   defaultHTTPClient,
		baseURL:      config.BaseURL,
		oauth2Client: oauth2ClientInternal,
		authMode:     AuthModeOAuth2,
	}

	// Initially should not have refresh token
	assert.False(t, oauth2Client.HasRefreshToken())

	// Set tokens with refresh token
	oauth2Client.SetTokens("access-token", "refresh-token", 3600)
	assert.True(t, oauth2Client.HasRefreshToken())

	// Test session mode
	sessionClient := NewClient(nil)
	assert.False(t, sessionClient.HasRefreshToken())
}

func TestClient_GetTokenExpiration(t *testing.T) {
	// Test OAuth2 mode
	config := NewProductionOAuth2Config("test-client-id", "test-client-secret", "http://localhost:8080/callback", []string{"read", "trade"})
	oauth2Client, err := NewOAuth2Client(config, nil)
	require.NoError(t, err)

	// Set token
	beforeSet := time.Now()
	oauth2Client.SetTokens("access-token", "refresh-token", 3600)
	afterSet := time.Now()

	expiration, err := oauth2Client.GetTokenExpiration()
	require.NoError(t, err)

	expectedExpiration := beforeSet.Add(time.Hour)
	assert.True(t, expiration.After(expectedExpiration.Add(-time.Second)))
	assert.True(t, expiration.Before(afterSet.Add(time.Hour).Add(time.Second)))

	// Test session mode
	sessionClient := NewClient(nil)
	_, err = sessionClient.GetTokenExpiration()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "OAuth2 mode")
}

func TestClient_GetTimeUntilExpiry(t *testing.T) {
	// Test OAuth2 mode
	config := NewProductionOAuth2Config("test-client-id", "test-client-secret", "http://localhost:8080/callback", []string{"read", "trade"})
	oauth2Client, err := NewOAuth2Client(config, nil)
	require.NoError(t, err)

	// Set token with 1 hour expiry
	oauth2Client.SetTokens("access-token", "refresh-token", 3600)

	timeUntilExpiry, err := oauth2Client.GetTimeUntilExpiry()
	require.NoError(t, err)

	// Should be close to 1 hour
	assert.True(t, timeUntilExpiry > 59*time.Minute)
	assert.True(t, timeUntilExpiry <= time.Hour)

	// Test session mode
	sessionClient := NewClient(nil)
	_, err = sessionClient.GetTimeUntilExpiry()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "OAuth2 mode")
}

func TestClient_IsTokenExpired(t *testing.T) {
	// Test OAuth2 mode with isolated storage
	config := NewProductionOAuth2Config("test-client-id", "test-client-secret", "http://localhost:8080/callback", []string{"read", "trade"})
	
	// Create OAuth2Client with memory storage for isolated testing
	tokenManager := NewMemoryTokenManager()
	oauth2ClientInternal := &OAuth2Client{
		config:       config,
		tokenManager: tokenManager,
		httpClient:   defaultHTTPClient,
		pkce:         nil,
	}
	tokenManager.SetRefreshCallback(oauth2ClientInternal.refreshTokensInternal)
	
	// Create Client with the isolated OAuth2Client
	oauth2Client := &Client{
		httpClient:   defaultHTTPClient,
		baseURL:      config.BaseURL,
		oauth2Client: oauth2ClientInternal,
		authMode:     AuthModeOAuth2,
	}

	// Initially should be expired (no token)
	assert.True(t, oauth2Client.IsTokenExpired())

	// Set valid token
	oauth2Client.SetTokens("access-token", "refresh-token", 3600)
	assert.False(t, oauth2Client.IsTokenExpired())

	// Set expired token
	oauth2Client.SetTokens("expired-token", "refresh-token", -10)
	assert.True(t, oauth2Client.IsTokenExpired())

	// Test session mode
	sessionClient := NewClient(nil)
	assert.True(t, sessionClient.IsTokenExpired())
}

func TestClient_TokenManagement_Integration(t *testing.T) {
	config := NewProductionOAuth2Config("test-client-id", "test-client-secret", "http://localhost:8080/callback", []string{"read", "trade"})
	
	// Test creating client with tokens
	client, err := NewOAuth2ClientWithTokens(config, "access-token", "refresh-token", 3600, nil)
	require.NoError(t, err)

	// Verify all token status methods work
	assert.True(t, client.IsAuthenticated())
	assert.True(t, client.HasValidToken())
	assert.True(t, client.HasRefreshToken())
	assert.False(t, client.IsTokenExpired())

	expiration, err := client.GetTokenExpiration()
	require.NoError(t, err)
	assert.True(t, expiration.After(time.Now()))

	timeUntilExpiry, err := client.GetTimeUntilExpiry()
	require.NoError(t, err)
	assert.True(t, timeUntilExpiry > 0)

	// Test updating tokens
	newTokenResponse := &TokenResponse{
		AccessToken:  "new-access-token",
		RefreshToken: "new-refresh-token",
		TokenType:    "Bearer",
		ExpiresIn:    7200,
		Scope:        "read trade",
	}

	err = client.SetTokensFromResponse(newTokenResponse)
	require.NoError(t, err)

	// Verify new tokens are active
	storedToken, err := client.oauth2Client.tokenManager.GetAccessToken()
	require.NoError(t, err)
	assert.Equal(t, "new-access-token", storedToken)

	// Test clearing tokens
	client.ClearAuthentication()
	assert.False(t, client.IsAuthenticated())
	assert.False(t, client.HasValidToken())
	assert.False(t, client.HasRefreshToken())
}