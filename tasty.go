package tasty

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
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

// AuthMode represents the authentication mode for the client
type AuthMode int

const (
	// AuthModeSession uses the legacy session-based authentication (deprecated)
	AuthModeSession AuthMode = iota
	// AuthModeOAuth2 uses OAuth2 authentication
	AuthModeOAuth2
)

// String returns a string representation of the AuthMode
func (am AuthMode) String() string {
	switch am {
	case AuthModeSession:
		return "session"
	case AuthModeOAuth2:
		return "oauth2"
	default:
		return "unknown"
	}
}

// Client for the tasty api wrapper.
type Client struct {
	httpClient *http.Client
	baseURL    string
	baseHost   string
	websocket  string

	// Legacy session support (deprecated)
	Session Session

	// OAuth2 support
	oauth2Client *OAuth2Client
	authMode     AuthMode
}

// NewClient creates a new Tasty Client using session-based authentication (deprecated).
// For new applications, use NewOAuth2Client instead.
func NewClient(httpClient *http.Client) *Client {
	LogSessionDeprecation("NewClient", "NewOAuth2Client() with OAuth2Config")

	if httpClient == nil {
		httpClient = defaultHTTPClient
	}
	c := &Client{
		httpClient: httpClient,
		baseURL:    apiBaseURL,
		baseHost:   apiBaseHost,
		websocket:  streamerBaseURL,
		authMode:   AuthModeSession,
	}

	return c
}

// NewCertClient creates a new Tasty Cert Client using session-based authentication (deprecated).
// For new applications, use NewCertOAuth2Client instead.
func NewCertClient(httpClient *http.Client) *Client {
	LogSessionDeprecation("NewCertClient", "NewCertOAuth2Client() with OAuth2Config")

	if httpClient == nil {
		httpClient = defaultHTTPClient
	}
	c := &Client{
		httpClient: httpClient,
		baseURL:    apiCertBaseURL,
		baseHost:   apiCertBaseHost,
		websocket:  streamerCertBaseURL,
		authMode:   AuthModeSession,
	}

	return c
}

// NewOAuth2Client creates a new Tasty Client using OAuth2 authentication for production environment.
func NewOAuth2Client(config OAuth2Config, httpClient *http.Client) (*Client, error) {
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
			return nil, fmt.Errorf("use NewCertOAuth2Client for sandbox environment")
		}
	}

	// Validate that we're using production endpoints
	if config.AuthURL != "" && config.AuthURL != oauth2ProductionAuthURL {
		return nil, fmt.Errorf("NewOAuth2Client requires production authorization URL, got: %s", config.AuthURL)
	}
	if config.TokenURL != "" && config.TokenURL != oauth2ProductionTokenURL {
		return nil, fmt.Errorf("NewOAuth2Client requires production token URL, got: %s", config.TokenURL)
	}

	// Create OAuth2 client
	oauth2Client, err := newOAuth2ClientInternal(config, httpClient)
	if err != nil {
		return nil, fmt.Errorf("failed to create OAuth2 client: %w", err)
	}

	c := &Client{
		httpClient:   httpClient,
		baseURL:      apiBaseURL,
		baseHost:     apiBaseHost,
		websocket:    streamerBaseURL,
		oauth2Client: oauth2Client,
		authMode:     AuthModeOAuth2,
	}

	return c, nil
}

// NewOAuth2ClientWithTokens creates a new Tasty Client using OAuth2 authentication for production environment
// and initializes it with the provided tokens. This is the recommended method for "bring your own tokens" usage.
func NewOAuth2ClientWithTokens(config OAuth2Config, accessToken, refreshToken string, expiresIn int, httpClient *http.Client) (*Client, error) {
	client, err := NewOAuth2Client(config, httpClient)
	if err != nil {
		return nil, err
	}

	// Set the provided tokens
	client.oauth2Client.SetTokens(accessToken, refreshToken, expiresIn)

	return client, nil
}

// NewOAuth2ClientWithTokenResponse creates a new Tasty Client using OAuth2 authentication for production environment
// and initializes it with tokens from a TokenResponse. This is useful when you have a complete TokenResponse
// from an external OAuth2 flow.
func NewOAuth2ClientWithTokenResponse(config OAuth2Config, tokenResponse *TokenResponse, httpClient *http.Client) (*Client, error) {
	client, err := NewOAuth2Client(config, httpClient)
	if err != nil {
		return nil, err
	}

	// Set tokens from the response
	client.oauth2Client.SetTokensFromResponse(tokenResponse)

	return client, nil
}

// NewCertOAuth2Client creates a new Tasty Cert Client using OAuth2 authentication for sandbox environment.
func NewCertOAuth2Client(config OAuth2Config, httpClient *http.Client) (*Client, error) {
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
			return nil, fmt.Errorf("use NewOAuth2Client for production environment")
		}
	}

	// Validate that we're using sandbox endpoints
	if config.AuthURL != "" && config.AuthURL != oauth2SandboxAuthURL {
		return nil, fmt.Errorf("NewCertOAuth2Client requires sandbox authorization URL, got: %s", config.AuthURL)
	}
	if config.TokenURL != "" && config.TokenURL != oauth2SandboxTokenURL {
		return nil, fmt.Errorf("NewCertOAuth2Client requires sandbox token URL, got: %s", config.TokenURL)
	}

	// Create OAuth2 client
	oauth2Client, err := newOAuth2ClientInternal(config, httpClient)
	if err != nil {
		return nil, fmt.Errorf("failed to create OAuth2 client: %w", err)
	}

	c := &Client{
		httpClient:   httpClient,
		baseURL:      apiCertBaseURL,
		baseHost:     apiCertBaseHost,
		websocket:    streamerCertBaseURL,
		oauth2Client: oauth2Client,
		authMode:     AuthModeOAuth2,
	}

	return c, nil
}

// NewCertOAuth2ClientWithTokens creates a new Tasty Cert Client using OAuth2 authentication for sandbox environment
// and initializes it with the provided tokens. This is the recommended method for "bring your own tokens" usage.
func NewCertOAuth2ClientWithTokens(config OAuth2Config, accessToken, refreshToken string, expiresIn int, httpClient *http.Client) (*Client, error) {
	client, err := NewCertOAuth2Client(config, httpClient)
	if err != nil {
		return nil, err
	}

	// Set the provided tokens
	client.oauth2Client.SetTokens(accessToken, refreshToken, expiresIn)

	return client, nil
}

// NewCertOAuth2ClientWithTokenResponse creates a new Tasty Cert Client using OAuth2 authentication for sandbox environment
// and initializes it with tokens from a TokenResponse. This is useful when you have a complete TokenResponse
// from an external OAuth2 flow.
func NewCertOAuth2ClientWithTokenResponse(config OAuth2Config, tokenResponse *TokenResponse, httpClient *http.Client) (*Client, error) {
	client, err := NewCertOAuth2Client(config, httpClient)
	if err != nil {
		return nil, err
	}

	// Set tokens from the response
	client.oauth2Client.SetTokensFromResponse(tokenResponse)

	return client, nil
}

// Getter for the tastytrade account streaming websocket url.
func (c Client) GetWebsocketURL() string {
	return c.websocket
}

// GetAuthMode returns the current authentication mode
func (c *Client) GetAuthMode() AuthMode {
	return c.authMode
}

// IsOAuth2Mode returns true if the client is using OAuth2 authentication
func (c *Client) IsOAuth2Mode() bool {
	return c.authMode == AuthModeOAuth2
}

// IsSessionMode returns true if the client is using session-based authentication
func (c *Client) IsSessionMode() bool {
	return c.authMode == AuthModeSession
}

// GetOAuth2Client returns the OAuth2 client if available
func (c *Client) GetOAuth2Client() *OAuth2Client {
	return c.oauth2Client
}

// GetAuthorizationURL generates the OAuth2 authorization URL (OAuth2 mode only)
func (c *Client) GetAuthorizationURL() (string, error) {
	if c.authMode != AuthModeOAuth2 {
		return "", fmt.Errorf("authorization URL is only available in OAuth2 mode")
	}
	if c.oauth2Client == nil {
		return "", fmt.Errorf("OAuth2 client not initialized")
	}
	return c.oauth2Client.GetAuthorizationURL()
}

// ExchangeCodeForTokens exchanges an authorization code for tokens (OAuth2 mode only)
func (c *Client) ExchangeCodeForTokens(code string) (*TokenResponse, error) {
	if c.authMode != AuthModeOAuth2 {
		return nil, fmt.Errorf("token exchange is only available in OAuth2 mode")
	}
	if c.oauth2Client == nil {
		return nil, fmt.Errorf("OAuth2 client not initialized")
	}
	return c.oauth2Client.ExchangeCodeForTokens(code)
}

// RefreshTokens refreshes the OAuth2 access token (OAuth2 mode only)
func (c *Client) RefreshTokens() (*TokenResponse, error) {
	if c.authMode != AuthModeOAuth2 {
		return nil, fmt.Errorf("token refresh is only available in OAuth2 mode")
	}
	if c.oauth2Client == nil {
		return nil, fmt.Errorf("OAuth2 client not initialized")
	}
	return c.oauth2Client.RefreshTokens()
}

// StartRedirectServer starts an HTTP server for OAuth2 redirects (OAuth2 mode only)
func (c *Client) StartRedirectServer(port int) (*RedirectServer, error) {
	if c.authMode != AuthModeOAuth2 {
		return nil, fmt.Errorf("redirect server is only available in OAuth2 mode")
	}
	if c.oauth2Client == nil {
		return nil, fmt.Errorf("OAuth2 client not initialized")
	}
	return c.oauth2Client.StartRedirectServer(port)
}

// ValidateState validates the OAuth2 state parameter (OAuth2 mode only)
func (c *Client) ValidateState(state string) error {
	if c.authMode != AuthModeOAuth2 {
		return fmt.Errorf("state validation is only available in OAuth2 mode")
	}
	if c.oauth2Client == nil {
		return fmt.Errorf("OAuth2 client not initialized")
	}
	return c.oauth2Client.ValidateState(state)
}

// IsAuthenticated checks if the client has valid authentication
func (c *Client) IsAuthenticated() bool {
	switch c.authMode {
	case AuthModeOAuth2:
		return c.oauth2Client != nil && c.oauth2Client.IsAuthenticated()
	case AuthModeSession:
		return c.Session.SessionToken != nil
	default:
		return false
	}
}

// ClearAuthentication clears all authentication data
func (c *Client) ClearAuthentication() {
	switch c.authMode {
	case AuthModeOAuth2:
		if c.oauth2Client != nil {
			c.oauth2Client.ClearTokens()
		}
	case AuthModeSession:
		c.Session = Session{}
	}
}

// SetTokens stores OAuth2 tokens directly in the client (OAuth2 mode only)
// This is the primary method for "bring your own tokens" usage
func (c *Client) SetTokens(accessToken, refreshToken string, expiresIn int) error {
	if c.authMode != AuthModeOAuth2 {
		return fmt.Errorf("SetTokens is only available in OAuth2 mode")
	}
	if c.oauth2Client == nil {
		return fmt.Errorf("OAuth2 client not initialized")
	}
	c.oauth2Client.SetTokens(accessToken, refreshToken, expiresIn)
	return nil
}

// SetTokensFromResponse stores tokens from a TokenResponse object (OAuth2 mode only)
func (c *Client) SetTokensFromResponse(response *TokenResponse) error {
	if c.authMode != AuthModeOAuth2 {
		return fmt.Errorf("SetTokensFromResponse is only available in OAuth2 mode")
	}
	if c.oauth2Client == nil {
		return fmt.Errorf("OAuth2 client not initialized")
	}
	c.oauth2Client.SetTokensFromResponse(response)
	return nil
}

// HasValidToken checks if the client has a valid (non-expired) access token (OAuth2 mode only)
func (c *Client) HasValidToken() bool {
	if c.authMode != AuthModeOAuth2 || c.oauth2Client == nil {
		return false
	}
	return c.oauth2Client.HasValidToken()
}

// HasRefreshToken checks if the client has a refresh token available (OAuth2 mode only)
func (c *Client) HasRefreshToken() bool {
	if c.authMode != AuthModeOAuth2 || c.oauth2Client == nil {
		return false
	}
	return c.oauth2Client.HasRefreshToken()
}

// GetTokenExpiration returns the expiration time of the current access token (OAuth2 mode only)
func (c *Client) GetTokenExpiration() (time.Time, error) {
	if c.authMode != AuthModeOAuth2 {
		return time.Time{}, fmt.Errorf("GetTokenExpiration is only available in OAuth2 mode")
	}
	if c.oauth2Client == nil {
		return time.Time{}, fmt.Errorf("OAuth2 client not initialized")
	}
	return c.oauth2Client.GetTokenExpiration(), nil
}

// GetTimeUntilExpiry returns the duration until the current access token expires (OAuth2 mode only)
func (c *Client) GetTimeUntilExpiry() (time.Duration, error) {
	if c.authMode != AuthModeOAuth2 {
		return 0, fmt.Errorf("GetTimeUntilExpiry is only available in OAuth2 mode")
	}
	if c.oauth2Client == nil {
		return 0, fmt.Errorf("OAuth2 client not initialized")
	}
	return c.oauth2Client.GetTimeUntilExpiry(), nil
}

// IsTokenExpired checks if the current access token is expired or about to expire (OAuth2 mode only)
func (c *Client) IsTokenExpired() bool {
	if c.authMode != AuthModeOAuth2 || c.oauth2Client == nil {
		return true
	}
	return c.oauth2Client.IsTokenExpired()
}

// TryOAuth2FallbackToSession attempts OAuth2 authentication first, then falls back to session if needed
// This is a utility method to help with gradual migration
// Deprecated: This method is provided for migration purposes only and will be removed in a future version
func (c *Client) TryOAuth2FallbackToSession(oauth2Config *OAuth2Config, sessionLogin *LoginInfo, twoFactorCode *string) error {
	LogSessionDeprecation("TryOAuth2FallbackToSession", "full OAuth2 authentication flow")

	// If OAuth2 config is provided, try OAuth2 first
	if oauth2Config != nil {
		// Validate OAuth2 config
		if err := ValidateOAuth2Migration(*oauth2Config); err == nil {
			// Try to create OAuth2 client
			var oauth2Client *OAuth2Client
			var err error

			if oauth2Config.IsProduction() {
				oauth2Client, err = newOAuth2ClientInternal(*oauth2Config, c.httpClient)
			} else {
				oauth2Client, err = newOAuth2ClientInternal(*oauth2Config, c.httpClient)
			}

			if err == nil {
				// Successfully created OAuth2 client, switch to OAuth2 mode
				c.oauth2Client = oauth2Client
				c.authMode = AuthModeOAuth2
				return nil
			}
		}
	}

	// OAuth2 failed or not configured, fall back to session
	if sessionLogin != nil {
		_, _, err := c.CreateSession(*sessionLogin, twoFactorCode)
		return err
	}

	return fmt.Errorf("both OAuth2 and session authentication failed or not configured")
}

// MigrateToOAuth2 helps migrate an existing session-based client to OAuth2
// This preserves the existing client instance while switching authentication modes
func (c *Client) MigrateToOAuth2(config OAuth2Config) error {
	// Validate OAuth2 configuration
	if err := ValidateOAuth2Migration(config); err != nil {
		return fmt.Errorf("OAuth2 migration validation failed: %w", err)
	}

	// Create OAuth2 client
	oauth2Client, err := newOAuth2ClientInternal(config, c.httpClient)
	if err != nil {
		return fmt.Errorf("failed to create OAuth2 client during migration: %w", err)
	}

	// Clear existing session data
	c.Session = Session{}

	// Switch to OAuth2
	c.oauth2Client = oauth2Client
	c.authMode = AuthModeOAuth2

	return nil
}

// GetMigrationStatus returns information about the current authentication mode and migration status
func (c *Client) GetMigrationStatus() map[string]interface{} {
	status := make(map[string]interface{})

	status["auth_mode"] = c.authMode.String()
	status["is_oauth2"] = c.IsOAuth2Mode()
	status["is_session"] = c.IsSessionMode()
	status["is_authenticated"] = c.IsAuthenticated()
	status["needs_migration"] = c.IsSessionMode()

	if c.IsOAuth2Mode() && c.oauth2Client != nil {
		status["oauth2_environment"] = c.oauth2Client.config.GetEnvironment()
		status["has_valid_token"] = c.oauth2Client.tokenManager.HasValidToken()
		status["has_refresh_token"] = c.oauth2Client.tokenManager.HasRefreshToken()
	}

	if c.IsSessionMode() {
		status["has_session_token"] = c.Session.SessionToken != nil
	}

	return status
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

// customRequest handles any requests for the client with unique paths and automatic authentication mode detection.
func (c *Client) customRequest(method, path string, params, payload, result any) (*http.Response, *Error) {
	switch c.authMode {
	case AuthModeOAuth2:
		return c.customOAuthRequest(method, path, params, payload, result)
	case AuthModeSession:
		return c.customSessionRequest(method, path, params, payload, result)
	default:
		return nil, &Error{Code: "invalid_auth_mode", Message: "Invalid authentication mode"}
	}
}

// customSessionRequest handles requests with unique paths using session-based authentication (deprecated).
func (c *Client) customSessionRequest(method, path string, params, payload, result any) (*http.Response, *Error) {
	if c.Session.SessionToken == nil {
		return nil, &Error{Code: "invalid_session", Message: "Session is invalid: Session Token cannot be nil."}
	}

	r := new(http.Request)

	r.Method = method

	r.URL = &url.URL{
		Scheme: strings.Split(c.baseURL, ":")[0],
		Host:   c.baseHost,
		Opaque: fmt.Sprintf("//%s%s", c.baseHost, path),
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, &Error{Message: fmt.Sprintf("Client Side Error: %v", err)}
	}

	r.Body = io.NopCloser(bytes.NewBuffer(body))

	r.Header = http.Header{}
	r.Header.Add("Authorization", *c.Session.SessionToken)
	r.Header.Add("Content-Type", "application/json")

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

// customOAuthRequest handles requests with unique paths using OAuth2 authentication.
func (c *Client) customOAuthRequest(method, path string, params, payload, result any) (*http.Response, *Error) {
	if c.oauth2Client == nil {
		return nil, &Error{Code: "invalid_oauth2", Message: "OAuth2 client not initialized"}
	}

	// Get access token (automatically refreshes if needed)
	accessToken, err := c.oauth2Client.GetTokenManager().GetAccessToken()
	if err != nil {
		return nil, &Error{Code: "oauth2_token_error", Message: fmt.Sprintf("Failed to get access token: %v", err)}
	}

	r := new(http.Request)

	r.Method = method

	r.URL = &url.URL{
		Scheme: strings.Split(c.baseURL, ":")[0],
		Host:   c.baseHost,
		Opaque: fmt.Sprintf("//%s%s", c.baseHost, path),
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, &Error{Message: fmt.Sprintf("Client Side Error: %v", err)}
	}

	r.Body = io.NopCloser(bytes.NewBuffer(body))

	r.Header = http.Header{}
	r.Header.Add("Authorization", "Bearer "+accessToken)
	r.Header.Add("Content-Type", "application/json")

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

	fmt.Printf("Retrieve status code: %d\n", resp.StatusCode)

	// Handle 401 Unauthorized - attempt token refresh and retry once
	if resp.StatusCode == http.StatusUnauthorized {
		// Try to refresh the token
		if _, refreshErr := c.oauth2Client.RefreshTokens(); refreshErr == nil {
			// Get the new access token
			if newAccessToken, tokenErr := c.oauth2Client.GetTokenManager().GetAccessToken(); tokenErr == nil {
				// Retry the request with the new token
				return c.retryCustomOAuthRequest(method, path, params, payload, result, newAccessToken)
			}
		}
		// If refresh failed, continue with original error handling
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
	r := new(http.Request)

	r.Method = method

	r.URL = &url.URL{
		Scheme: strings.Split(c.baseURL, ":")[0],
		Host:   c.baseHost,
		Opaque: fmt.Sprintf("//%s%s", c.baseHost, path),
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, &Error{Message: fmt.Sprintf("Client Side Error: %v", err)}
	}

	r.Body = io.NopCloser(bytes.NewBuffer(body))

	r.Header = http.Header{}
	r.Header.Add("Authorization", "Bearer "+accessToken)
	r.Header.Add("Content-Type", "application/json")

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

// request handles any requests for the client with automatic authentication mode detection.
func (c *Client) request(method, path string, params, payload, result any) (*http.Response, *Error) {
	switch c.authMode {
	case AuthModeOAuth2:
		return c.oauthRequest(method, path, params, payload, result)
	case AuthModeSession:
		return c.sessionRequest(method, path, params, payload, result)
	default:
		return nil, &Error{Code: "invalid_auth_mode", Message: "Invalid authentication mode"}
	}
}

// sessionRequest handles requests using session-based authentication (deprecated).
func (c *Client) sessionRequest(method, path string, params, payload, result any) (*http.Response, *Error) {
	if c.Session.SessionToken == nil {
		return nil, &Error{Code: "invalid_session", Message: "Session is invalid: Session Token cannot be nil."}
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, &Error{Message: fmt.Sprintf("Client Side Error: %v", err)}
	}

	fullURL := c.baseURL + path

	r, err := http.NewRequest(method, fullURL, bytes.NewBuffer(body))
	if err != nil {
		return nil, &Error{Message: fmt.Sprintf("Client Side Error: %v", err)}
	}

	r.Header = http.Header{}
	r.Header.Add("Authorization", *c.Session.SessionToken)
	r.Header.Add("Content-Type", "application/json")

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

// oauthRequest handles requests using OAuth2 authentication with automatic token refresh.
func (c *Client) oauthRequest(method, path string, params, payload, result any) (*http.Response, *Error) {
	if c.oauth2Client == nil {
		return nil, &Error{Code: "invalid_oauth2", Message: "OAuth2 client not initialized"}
	}

	accessToken, err := c.oauth2Client.GetTokenManager().GetAccessToken()
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
	requestDump, err := httputil.DumpRequestOut(r, true)
	if err != nil {
		return nil, &Error{Message: fmt.Sprintf("Client Side Error: failed to dump request: %v", err)}
	}
	fmt.Println("--- Request Dump ---")
	fmt.Printf("%s", requestDump)
	fmt.Println("--------------------")
	// ----------------------------------------
	// End of new logging code
	// ----------------------------------------

	resp, err := c.httpClient.Do(r)
	if err != nil {
		return nil, &Error{Message: fmt.Sprintf("Client Side Error: %v", err)}
	}

	defer resp.Body.Close()

	fmt.Printf("Retrieve status code: %d\n", resp.StatusCode)

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
		if _, refreshErr := c.oauth2Client.RefreshTokens(); refreshErr == nil {
			if newAccessToken, tokenErr := c.oauth2Client.GetTokenManager().GetAccessToken(); tokenErr == nil {
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

// noAuthRequest handles any requests for the client without authentication.
func (c *Client) noAuthRequest(method, path string, header http.Header, params, payload, result any) (*http.Response, *Error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, &Error{Message: fmt.Sprintf("Client Side Error: %v", err)}
	}

	fullURL := c.baseURL + path

	r, err := http.NewRequest(method, fullURL, bytes.NewBuffer(body))
	if err != nil {
		return nil, &Error{Message: fmt.Sprintf("Client Side Error: %v", err)}
	}

	if header == nil {
		r.Header = http.Header{}
	} else {
		r.Header = header
	}

	r.Header.Add("Content-Type", "application/json")

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
