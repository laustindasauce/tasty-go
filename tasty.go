package tasty

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/google/go-querystring/query"
)

const (
	apiBaseURL          = "https://api.tastyworks.com"
	apiBaseHost         = "api.tastyworks.com"
	apiCertBaseURL      = "https://api.cert.tastyworks.com"
	apiCertBaseHost     = "api.cert.tastyworks.com"
	streamerBaseURL     = "wss://streamer.tastyworks.com"
	streamerCertBaseURL = "wss://streamer.cert.tastyworks.com"
)

var (
	defaultHTTPClient = &http.Client{Timeout: time.Duration(30) * time.Second}
	errorStatusCodes  = []int{400, 401, 403, 404, 415, 422, 500}
)

// Client for the tasty api wrapper with OAuth2 authentication.
type Client struct {
	// HTTP client for making requests
	httpClient *http.Client

	// API endpoints
	baseURL   string
	baseHost  string
	websocket string

	// OAuth2 configuration and token management
	config       OAuth2Config
	tokenManager *TokenManager
	pkce         *PKCEChallenge
}

// NewClient creates a new Tasty Client using OAuth2 authentication for production environment.
func NewClient(config OAuth2Config, httpClient *http.Client) (*Client, error) {
	if httpClient == nil {
		httpClient = defaultHTTPClient
	}

	// Ensure production environment configuration
	if config.BaseURL == "" && config.AuthURL == "" && config.TokenURL == "" {
		// Set production endpoints if none are specified
		config.BaseURL = apiBaseURL
		config.AuthURL = oauth2ProductionAuthURL
		config.TokenURL = oauth2ProductionTokenURL
	} else if config.BaseURL != "" && (config.AuthURL == "" || config.TokenURL == "") {
		// Set endpoints based on base URL
		if config.BaseURL == apiBaseURL {
			config.AuthURL = oauth2ProductionAuthURL
			config.TokenURL = oauth2ProductionTokenURL
		} else if config.BaseURL == apiCertBaseURL {
			return nil, fmt.Errorf("use NewCertClient for sandbox environment")
		}
	}

	// Validate that we're using production endpoints
	if config.AuthURL != "" && config.AuthURL != oauth2ProductionAuthURL {
		return nil, fmt.Errorf("NewClient requires production authorization URL, got: %s", config.AuthURL)
	}
	if config.TokenURL != "" && config.TokenURL != oauth2ProductionTokenURL {
		return nil, fmt.Errorf("NewClient requires production token URL, got: %s", config.TokenURL)
	}

	// Validate the configuration
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid OAuth2 configuration: %w", err)
	}

	// Set default scopes if not provided
	if len(config.Scopes) == 0 {
		config.Scopes = []string{defaultScope}
	}

	// Generate state parameter if not provided
	if config.State == "" {
		state, err := generateSecureState()
		if err != nil {
			return nil, fmt.Errorf("failed to generate state parameter: %w", err)
		}
		config.State = state
	}

	tokenManager := NewTokenManager()

	c := &Client{
		httpClient:   httpClient,
		baseURL:      apiBaseURL,
		baseHost:     apiBaseHost,
		websocket:    streamerBaseURL,
		config:       config,
		tokenManager: tokenManager,
		pkce:         nil, // PKCE not supported by TastyTrade
	}

	// Set up token refresh callback
	tokenManager.SetRefreshCallback(c.refreshTokensInternal)

	return c, nil
}

// NewClientWithTokens creates a new Tasty Client using OAuth2 authentication for production environment
// and initializes it with the provided tokens. This is the recommended method for "bring your own tokens" usage.
func NewClientWithTokens(config OAuth2Config, accessToken, refreshToken string, expiresIn int, httpClient *http.Client) (*Client, error) {
	client, err := NewClient(config, httpClient)
	if err != nil {
		return nil, err
	}

	// Set the provided tokens
	client.SetTokens(accessToken, refreshToken, expiresIn)

	return client, nil
}

// NewClientWithTokenResponse creates a new Tasty Client using OAuth2 authentication for production environment
// and initializes it with tokens from a TokenResponse. This is useful when you have a complete TokenResponse
// from an external OAuth2 flow.
func NewClientWithTokenResponse(config OAuth2Config, tokenResponse *TokenResponse, httpClient *http.Client) (*Client, error) {
	client, err := NewClient(config, httpClient)
	if err != nil {
		return nil, err
	}

	// Set tokens from the response
	client.SetTokensFromResponse(tokenResponse)

	return client, nil
}

// NewCertClient creates a new Tasty Cert Client using OAuth2 authentication for sandbox environment.
func NewCertClient(config OAuth2Config, httpClient *http.Client) (*Client, error) {
	if httpClient == nil {
		httpClient = defaultHTTPClient
	}

	// Ensure sandbox environment configuration
	if config.BaseURL == "" && config.AuthURL == "" && config.TokenURL == "" {
		// Set sandbox endpoints if none are specified
		config.BaseURL = apiCertBaseURL
		config.AuthURL = oauth2SandboxAuthURL
		config.TokenURL = oauth2SandboxTokenURL
	} else if config.BaseURL != "" && (config.AuthURL == "" || config.TokenURL == "") {
		// Set endpoints based on base URL
		if config.BaseURL == apiCertBaseURL {
			config.AuthURL = oauth2SandboxAuthURL
			config.TokenURL = oauth2SandboxTokenURL
		} else if config.BaseURL == apiBaseURL {
			return nil, fmt.Errorf("use NewClient for production environment")
		}
	}

	// Validate that we're using sandbox endpoints
	if config.AuthURL != "" && config.AuthURL != oauth2SandboxAuthURL {
		return nil, fmt.Errorf("NewCertClient requires sandbox authorization URL, got: %s", config.AuthURL)
	}
	if config.TokenURL != "" && config.TokenURL != oauth2SandboxTokenURL {
		return nil, fmt.Errorf("NewCertClient requires sandbox token URL, got: %s", config.TokenURL)
	}

	// Validate the configuration
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid OAuth2 configuration: %w", err)
	}

	// Set default scopes if not provided
	if len(config.Scopes) == 0 {
		config.Scopes = []string{defaultScope}
	}

	// Generate state parameter if not provided
	if config.State == "" {
		state, err := generateSecureState()
		if err != nil {
			return nil, fmt.Errorf("failed to generate state parameter: %w", err)
		}
		config.State = state
	}

	tokenManager := NewTokenManager()

	c := &Client{
		httpClient:   httpClient,
		baseURL:      apiCertBaseURL,
		baseHost:     apiCertBaseHost,
		websocket:    streamerCertBaseURL,
		config:       config,
		tokenManager: tokenManager,
		pkce:         nil, // PKCE not supported by TastyTrade
	}

	// Set up token refresh callback
	tokenManager.SetRefreshCallback(c.refreshTokensInternal)

	return c, nil
}

// NewCertClientWithTokens creates a new Tasty Cert Client using OAuth2 authentication for sandbox environment
// and initializes it with the provided tokens. This is the recommended method for "bring your own tokens" usage.
func NewCertClientWithTokens(config OAuth2Config, accessToken, refreshToken string, expiresIn int, httpClient *http.Client) (*Client, error) {
	client, err := NewCertClient(config, httpClient)
	if err != nil {
		return nil, err
	}

	// Set the provided tokens
	client.SetTokens(accessToken, refreshToken, expiresIn)

	return client, nil
}

// NewCertClientWithTokenResponse creates a new Tasty Cert Client using OAuth2 authentication for sandbox environment
// and initializes it with tokens from a TokenResponse. This is useful when you have a complete TokenResponse
// from an external OAuth2 flow.
func NewCertClientWithTokenResponse(config OAuth2Config, tokenResponse *TokenResponse, httpClient *http.Client) (*Client, error) {
	client, err := NewCertClient(config, httpClient)
	if err != nil {
		return nil, err
	}

	// Set tokens from the response
	client.SetTokensFromResponse(tokenResponse)

	return client, nil
}

// Getter for the tastytrade account streaming websocket url.
func (c Client) GetWebsocketURL() string {
	return c.websocket
}

// GetAuthorizationURL generates the OAuth2 authorization URL with state parameters
func (c *Client) GetAuthorizationURL() (string, error) {
	if c.config.ClientID == "" {
		return "", NewOAuth2Error(OAuth2ErrorConfigurationError, "client ID is required")
	}
	if c.config.RedirectURI == "" {
		return "", NewOAuth2Error(OAuth2ErrorConfigurationError, "redirect URI is required")
	}
	if c.config.AuthURL == "" {
		return "", NewOAuth2Error(OAuth2ErrorConfigurationError, "authorization URL not configured")
	}

	// Build authorization URL
	authURL, err := url.Parse(c.config.AuthURL)
	if err != nil {
		return "", NewOAuth2ErrorWithContext(OAuth2ErrorConfigurationError,
			"failed to parse authorization URL", 0, err)
	}

	// Add query parameters (without PKCE since TastyTrade doesn't support it)
	params := url.Values{}
	params.Set("response_type", responseTypeCode)
	params.Set("client_id", c.config.ClientID)
	params.Set("redirect_uri", c.config.RedirectURI)
	params.Set("scope", strings.Join(c.config.Scopes, " "))
	params.Set("state", c.config.State)

	authURL.RawQuery = params.Encode()

	return authURL.String(), nil
}

// ExchangeCodeForTokens exchanges an authorization code for access and refresh tokens
func (c *Client) ExchangeCodeForTokens(code string) (*TokenResponse, error) {
	// Validate authorization code
	if err := ValidateAuthorizationCode(code); err != nil {
		return nil, err
	}

	// Prepare token exchange request
	if c.config.TokenURL == "" {
		return nil, NewOAuth2Error(OAuth2ErrorConfigurationError, "token URL not configured")
	}
	tokenURL := c.config.TokenURL

	// Build form data (without PKCE since TastyTrade doesn't support it)
	data := url.Values{}
	data.Set("grant_type", grantTypeAuthorizationCode)
	data.Set("client_id", c.config.ClientID)
	data.Set("code", code)
	data.Set("redirect_uri", c.config.RedirectURI)

	// Add client secret if provided (for confidential clients)
	if c.config.ClientSecret != "" {
		data.Set("client_secret", c.config.ClientSecret)
	}

	// Make token exchange request
	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, NewOAuth2ErrorWithContext(OAuth2ErrorNetworkError,
			"failed to create token request", 0, err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, NewOAuth2ErrorWithContext(OAuth2ErrorNetworkError,
			"token exchange request failed", 0, err)
	}
	defer resp.Body.Close()

	// Handle error responses
	if resp.StatusCode != http.StatusOK {
		var oauthErr OAuth2Error
		if err := json.NewDecoder(resp.Body).Decode(&oauthErr); err != nil {
			// Failed to parse OAuth2 error, create generic HTTP error
			return nil, WrapHTTPError(resp, err)
		}
		// Convert standard OAuth2 error to detailed error
		return nil, NewOAuth2ErrorFromStandard(oauthErr)
	}

	// Parse successful response
	var tokenResponse TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return nil, NewOAuth2ErrorWithContext(OAuth2ErrorServerError,
			"failed to parse token response", resp.StatusCode, err)
	}

	// Validate response
	if err := ValidateTokenResponse(&tokenResponse); err != nil {
		return nil, err
	}

	// Set scope from config if not provided in response
	if tokenResponse.Scope == "" && len(c.config.Scopes) > 0 {
		tokenResponse.Scope = strings.Join(c.config.Scopes, " ")
	}

	// Store tokens in token manager
	c.tokenManager.SetTokensFromResponse(&tokenResponse)

	return &tokenResponse, nil
}

// RefreshTokens refreshes the access token using the stored refresh token
func (c *Client) RefreshTokens() (*TokenResponse, error) {
	refreshToken := c.tokenManager.GetRefreshToken()
	if refreshToken == "" {
		return nil, NewOAuth2Error(OAuth2ErrorRefreshFailed, "no refresh token available")
	}

	return c.refreshTokensWithRetry(refreshToken, 3)
}

// refreshTokensWithRetry attempts to refresh tokens with automatic retry logic
func (c *Client) refreshTokensWithRetry(refreshToken string, maxRetries int) (*TokenResponse, error) {
	var lastErr error
	errorHandler := NewOAuth2ErrorHandler()

	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			// Exponential backoff: wait 1s, 2s, 4s between retries
			backoff := time.Duration(1<<uint(attempt-1)) * time.Second
			time.Sleep(backoff)
		}

		tokenResponse, err := c.performTokenRefresh(refreshToken)
		if err == nil {
			return tokenResponse, nil
		}

		lastErr = err

		// Use error handler to determine if we should retry
		detailedErr, shouldRetry := errorHandler.HandleError(err)
		if !shouldRetry {
			return nil, detailedErr
		}
	}

	// Create a comprehensive error for the final failure
	if detailedErr, ok := lastErr.(*OAuth2DetailedError); ok {
		detailedErr.InternalMessage = fmt.Sprintf("token refresh failed after %d attempts", maxRetries)
		return nil, detailedErr
	}

	return nil, NewOAuth2ErrorWithContext(OAuth2ErrorRefreshFailed,
		fmt.Sprintf("token refresh failed after %d attempts", maxRetries), 0, lastErr)
}

// performTokenRefresh performs the actual token refresh request
func (c *Client) performTokenRefresh(refreshToken string) (*TokenResponse, error) {
	if c.config.TokenURL == "" {
		return nil, NewOAuth2Error(OAuth2ErrorConfigurationError, "token URL not configured")
	}
	tokenURL := c.config.TokenURL

	// Build form data
	data := url.Values{}
	data.Set("grant_type", grantTypeRefreshToken)
	data.Set("client_id", c.config.ClientID)
	data.Set("refresh_token", refreshToken)

	// Add client secret if provided (for confidential clients)
	if c.config.ClientSecret != "" {
		data.Set("client_secret", c.config.ClientSecret)
	}

	// Make refresh request
	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, NewOAuth2ErrorWithContext(OAuth2ErrorNetworkError,
			"failed to create refresh request", 0, err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, NewOAuth2ErrorWithContext(OAuth2ErrorNetworkError,
			"refresh request failed", 0, err)
	}
	defer resp.Body.Close()

	// Handle error responses
	if resp.StatusCode != http.StatusOK {
		var oauthErr OAuth2Error
		if err := json.NewDecoder(resp.Body).Decode(&oauthErr); err != nil {
			// Failed to parse OAuth2 error, create generic HTTP error
			return nil, WrapHTTPError(resp, err)
		}
		// Convert standard OAuth2 error to detailed error
		return nil, NewOAuth2ErrorFromStandard(oauthErr)
	}

	// Parse successful response
	var tokenResponse TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return nil, NewOAuth2ErrorWithContext(OAuth2ErrorServerError,
			"failed to parse refresh response", resp.StatusCode, err)
	}

	// Validate response
	if err := ValidateTokenResponse(&tokenResponse); err != nil {
		return nil, err
	}

	// Get existing token data to preserve refresh token and scope if not provided in response
	existingRefreshToken := c.tokenManager.GetRefreshToken()
	existingScope := c.tokenManager.GetScope()

	// Set scope from config or preserve existing scope if not provided in response (common for refresh responses)
	if tokenResponse.Scope == "" {
		if len(c.config.Scopes) > 0 {
			tokenResponse.Scope = strings.Join(c.config.Scopes, " ")
		} else if existingScope != "" {
			tokenResponse.Scope = existingScope
		}
	}

	// Preserve existing refresh token if not provided in response (common for TastyTrade)
	if tokenResponse.RefreshToken == "" && existingRefreshToken != "" {
		tokenResponse.RefreshToken = existingRefreshToken
	}

	// Store new tokens in token manager
	c.tokenManager.SetTokensFromResponse(&tokenResponse)

	return &tokenResponse, nil
}

// refreshTokensInternal is the internal callback used by TokenManager
func (c *Client) refreshTokensInternal() (*TokenResponse, error) {
	return c.RefreshTokens()
}

// GetTokenManager returns the token manager for accessing token information
func (c *Client) GetTokenManager() *TokenManager {
	return c.tokenManager
}

// GetConfig returns a copy of the OAuth2 configuration
func (c *Client) GetConfig() OAuth2Config {
	return c.config
}

// GetState returns the current state parameter for CSRF protection
func (c *Client) GetState() string {
	return c.config.State
}

// ValidateState validates that the provided state matches the expected state
func (c *Client) ValidateState(state string) error {
	return ValidateState(c.config.State, state)
}

// IsAuthenticated checks if the client has valid authentication tokens
func (c *Client) IsAuthenticated() bool {
	return c.tokenManager.HasValidToken() || c.tokenManager.HasRefreshToken()
}

// ClearTokens securely clears all stored authentication tokens
func (c *Client) ClearTokens() {
	c.tokenManager.Clear()
}

// SetTokens stores OAuth2 tokens directly in the client for "bring your own tokens" usage
func (c *Client) SetTokens(accessToken, refreshToken string, expiresIn int) {
	c.tokenManager.SetTokens(accessToken, refreshToken, expiresIn)
}

// SetTokensFromResponse stores tokens from a TokenResponse object
func (c *Client) SetTokensFromResponse(response *TokenResponse) {
	c.tokenManager.SetTokensFromResponse(response)
}

// HasValidToken checks if the client has a valid (non-expired) access token
func (c *Client) HasValidToken() bool {
	return c.tokenManager.HasValidToken()
}

// HasRefreshToken checks if the client has a refresh token available
func (c *Client) HasRefreshToken() bool {
	return c.tokenManager.HasRefreshToken()
}

// GetTokenExpiration returns the expiration time of the current access token
func (c *Client) GetTokenExpiration() time.Time {
	return c.tokenManager.GetExpiresAt()
}

// GetTimeUntilExpiry returns the duration until the current access token expires
func (c *Client) GetTimeUntilExpiry() time.Duration {
	return c.tokenManager.GetTimeUntilExpiry()
}

// IsTokenExpired checks if the current access token is expired or about to expire
func (c *Client) IsTokenExpired() bool {
	return c.tokenManager.IsExpired()
}

// StartRedirectServer starts an HTTP server for OAuth2 redirects
func (c *Client) StartRedirectServer(port int) (*RedirectServer, error) {
	server := NewRedirectServer(c.config.State)

	if err := server.Start(port); err != nil {
		if detailedErr, ok := err.(*OAuth2DetailedError); ok {
			return nil, detailedErr
		}
		return nil, NewOAuth2ErrorWithContext(OAuth2ErrorConfigurationError,
			"failed to start redirect server", 0, err)
	}

	return server, nil
}

// Error reasoning given by tastytrade.
type ErrorResponse struct {
	Domain string `json:"domain"`
	Reason string `json:"reason"`
}

// Error represents an error returned by the tastytrade API.
type Error struct {
	// Simple code error string
	Code string `json:"code"`
	// A short description of the error.
	Message string `json:"message"`
	// Slice of errors
	Errors []ErrorResponse `json:"errors"`
	// The HTTP status code.
	StatusCode int `json:"error,omitempty"`
}

// Error ...
func (e Error) Error() string {
	return fmt.Sprintf("\nError in request %d;\nCode: %s\nMessage: %s", e.StatusCode, e.Code, e.Message)
}

// decodeError decodes an Error from response status code based off
// the developer docs in tastytrade -> https://developer.tastytrade.com/#error-codes
func decodeError(resp *http.Response) *Error {
	e := new(Error)

	type errorRes struct {
		Error Error `json:"error"`
	}

	errRes := new(errorRes)

	err := json.NewDecoder(resp.Body).Decode(errRes)
	if err != nil {
		e.Message = fmt.Sprintf("tastytrade: unexpected HTTP %d: %s (empty error)", resp.StatusCode, err.Error())
		e.StatusCode = resp.StatusCode
		return e
	}

	errRes.Error.StatusCode = resp.StatusCode

	e = &errRes.Error

	return e
}

// customRequest handles any requests for the client with unique paths using OAuth2 authentication.
func (c *Client) customRequest(method, path string, params, payload, result any) (*http.Response, *Error) {
	return c.customOAuthRequest(method, path, params, payload, result)
}

// customOAuthRequest handles requests with unique paths using OAuth2 authentication.
func (c *Client) customOAuthRequest(method, path string, params, payload, result any) (*http.Response, *Error) {
	// Get access token (automatically refreshes if needed)
	accessToken, err := c.tokenManager.GetAccessToken()
	if err != nil {
		return nil, &Error{Code: "oauth2_token_error", Message: fmt.Sprintf("Failed to get access token: %v", err)}
	}

	fullURL := c.baseURL + path

	var bodyReader io.Reader
	if payload != nil && (method == http.MethodPost || method == http.MethodPut || method == http.MethodPatch) {
		body, err := json.Marshal(payload)
		if err != nil {
			return nil, &Error{Message: fmt.Sprintf("Client Side Error: %v", err)}
		}
		bodyReader = bytes.NewBuffer(body)
	}

	r, err := http.NewRequest(method, fullURL, bodyReader)
	if err != nil {
		return nil, &Error{Message: fmt.Sprintf("Client Side Error: %v", err)}
	}

	r.Header = http.Header{}
	r.Header.Add("Authorization", "Bearer "+accessToken)
	// Add Content-Type only if a body is present
	if bodyReader != nil {
		r.Header.Add("Content-Type", "application/json")
	}

	if params != nil {
		queryString, queryErr := query.Values(params)
		if queryErr != nil {
			return nil, &Error{Message: fmt.Sprintf("Client Side Error: %v", err)}
		}
		r.URL.RawQuery = queryString.Encode()
	}

	resp, err := c.httpClient.Do(r)
	if err != nil {
		return nil, &Error{Message: fmt.Sprintf("Client Side Error: %v", err)}
	}

	defer resp.Body.Close()

	// Handle 401 Unauthorized - attempt token refresh and retry once
	if resp.StatusCode == http.StatusUnauthorized {
		if _, refreshErr := c.RefreshTokens(); refreshErr == nil {
			if newAccessToken, tokenErr := c.tokenManager.GetAccessToken(); tokenErr == nil {
				return c.retryCustomOAuthRequest(method, path, params, payload, result, newAccessToken)
			}
		}
	}

	if resp.StatusCode == http.StatusNoContent {
		return resp, nil
	}
	if containsInt(errorStatusCodes, resp.StatusCode) {
		return resp, decodeError(resp)
	}

	if result != nil {
		err = json.NewDecoder(resp.Body).Decode(result)
		if err != nil {
			return resp, &Error{Message: fmt.Sprintf("Client Side Error: %v", err)}
		}
	}

	return resp, nil
}

// retryCustomOAuthRequest retries a custom OAuth2 request with a new access token
func (c *Client) retryCustomOAuthRequest(method, path string, params, payload, result any, accessToken string) (*http.Response, *Error) {
	fullURL := c.baseURL + path

	var bodyReader io.Reader
	if payload != nil && (method == http.MethodPost || method == http.MethodPut || method == http.MethodPatch) {
		body, err := json.Marshal(payload)
		if err != nil {
			return nil, &Error{Message: fmt.Sprintf("Client Side Error: %v", err)}
		}
		bodyReader = bytes.NewBuffer(body)
	}

	r, err := http.NewRequest(method, fullURL, bodyReader)
	if err != nil {
		return nil, &Error{Message: fmt.Sprintf("Client Side Error: %v", err)}
	}

	r.Header = http.Header{}
	r.Header.Add("Authorization", "Bearer "+accessToken)
	// Add Content-Type only if a body is present
	if bodyReader != nil {
		r.Header.Add("Content-Type", "application/json")
	}

	if params != nil {
		queryString, queryErr := query.Values(params)
		if queryErr != nil {
			return nil, &Error{Message: fmt.Sprintf("Client Side Error: %v", err)}
		}
		r.URL.RawQuery = queryString.Encode()
	}

	resp, err := c.httpClient.Do(r)
	if err != nil {
		return nil, &Error{Message: fmt.Sprintf("Client Side Error: %v", err)}
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		return resp, nil
	}
	if containsInt(errorStatusCodes, resp.StatusCode) {
		return resp, decodeError(resp)
	}

	if result != nil {
		err = json.NewDecoder(resp.Body).Decode(result)
		if err != nil {
			return resp, &Error{Message: fmt.Sprintf("Client Side Error: %v", err)}
		}
	}

	return resp, nil
}

// request handles any requests for the client using OAuth2 authentication.
func (c *Client) request(method, path string, params, payload, result any) (*http.Response, *Error) {
	return c.oauthRequest(method, path, params, payload, result)
}

// oauthRequest handles requests using OAuth2 authentication with automatic token refresh.
func (c *Client) oauthRequest(method, path string, params, payload, result any) (*http.Response, *Error) {
	accessToken, err := c.tokenManager.GetAccessToken()
	if err != nil {
		return nil, &Error{Code: "oauth2_token_error", Message: fmt.Sprintf("Failed to get access token: %v", err)}
	}

	// Initialize an empty buffer for the request body
	var bodyReader io.Reader

	// Condition to set the body and Content-Type header only for methods that need a body
	if payload != nil && (method == http.MethodPost || method == http.MethodPut || method == http.MethodPatch) {
		body, err := json.Marshal(payload)
		if err != nil {
			return nil, &Error{Message: fmt.Sprintf("Client Side Error: %v", err)}
		}
		bodyReader = bytes.NewBuffer(body)
	}

	fullURL := c.baseURL + path

	r, err := http.NewRequest(method, fullURL, bodyReader)
	if err != nil {
		return nil, &Error{Message: fmt.Sprintf("Client Side Error: %v", err)}
	}

	r.Header = http.Header{}
	r.Header.Add("Authorization", "Bearer "+accessToken)

	// Condition to add Content-Type header
	if bodyReader != nil {
		r.Header.Add("Content-Type", "application/json")
	}

	if params != nil {
		queryString, queryErr := query.Values(params)
		if queryErr != nil {
			return nil, &Error{Message: fmt.Sprintf("Client Side Error: %v", err)}
		}
		r.URL.RawQuery = queryString.Encode()
	}

	// ----------------------------------------
	// Start of new logging code for the request
	// ----------------------------------------
	// requestDump, err := httputil.DumpRequestOut(r, true)
	// if err != nil {
	// 	return nil, &Error{Message: fmt.Sprintf("Client Side Error: failed to dump request: %v", err)}
	// }
	// fmt.Println("--- Request Dump ---")
	// fmt.Printf("%s", requestDump)
	// fmt.Println("--------------------")
	// ----------------------------------------
	// End of new logging code
	// ----------------------------------------

	resp, err := c.httpClient.Do(r)
	if err != nil {
		return nil, &Error{Message: fmt.Sprintf("Client Side Error: %v", err)}
	}

	defer resp.Body.Close()

	// ----------------------------------------
	// Start of new logging code for the response
	// ----------------------------------------
	// bodyBytes, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	return resp, &Error{Message: fmt.Sprintf("Client Side Error: failed to read response body: %v", err)}
	// }

	// // Re-create the response body so it can be decoded later
	// resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	// fmt.Println("--- Response Body Dump ---")
	// fmt.Printf("%s\n", string(bodyBytes))
	// fmt.Println("--------------------------")
	// ----------------------------------------
	// End of new logging code
	// ----------------------------------------

	if resp.StatusCode == http.StatusUnauthorized {
		if _, refreshErr := c.RefreshTokens(); refreshErr == nil {
			if newAccessToken, tokenErr := c.tokenManager.GetAccessToken(); tokenErr == nil {
				return c.retryOAuthRequest(method, path, params, payload, result, newAccessToken)
			}
		}
	}

	if resp.StatusCode == http.StatusNoContent {
		return resp, nil
	}
	if containsInt(errorStatusCodes, resp.StatusCode) {
		return resp, decodeError(resp)
	}

	if result != nil {
		err = json.NewDecoder(resp.Body).Decode(result)
		if err != nil {
			return resp, &Error{Message: fmt.Sprintf("Client Side Error: %v", err)}
		}
	}

	return resp, nil
}

// retryOAuthRequest retries an OAuth2 request with a new access token
func (c *Client) retryOAuthRequest(method, path string, params, payload, result any, accessToken string) (*http.Response, *Error) {
	// Initialize an empty reader for the request body.
	var bodyReader io.Reader

	// Only create a body for methods that require a payload.
	if payload != nil && (method == http.MethodPost || method == http.MethodPut || method == http.MethodPatch) {
		body, err := json.Marshal(payload)
		if err != nil {
			return nil, &Error{Message: fmt.Sprintf("Client Side Error: %v", err)}
		}
		bodyReader = bytes.NewBuffer(body)
	}

	fullURL := c.baseURL + path

	r, err := http.NewRequest(method, fullURL, bodyReader)
	if err != nil {
		return nil, &Error{Message: fmt.Sprintf("Client Side Error: %v", err)}
	}

	r.Header = http.Header{}
	r.Header.Add("Authorization", "Bearer "+accessToken)

	// Add Content-Type header only when a body is present.
	if bodyReader != nil {
		r.Header.Add("Content-Type", "application/json")
	}

	if params != nil {
		queryString, queryErr := query.Values(params)
		if queryErr != nil {
			return nil, &Error{Message: fmt.Sprintf("Client Side Error: %v", err)}
		}
		r.URL.RawQuery = queryString.Encode()
	}

	resp, err := c.httpClient.Do(r)
	if err != nil {
		return nil, &Error{Message: fmt.Sprintf("Client Side Error: %v", err)}
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		return resp, nil
	}
	if containsInt(errorStatusCodes, resp.StatusCode) {
		return resp, decodeError(resp)
	}

	if result != nil {
		err = json.NewDecoder(resp.Body).Decode(result)
		if err != nil {
			return resp, &Error{Message: fmt.Sprintf("Client Side Error: %v", err)}
		}
	}

	return resp, nil
}

// generateSecureState generates a cryptographically secure state parameter
func generateSecureState() (string, error) {
	bytes := make([]byte, 32) // 32 bytes = 256 bits of entropy
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate secure random bytes: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(bytes), nil
}
