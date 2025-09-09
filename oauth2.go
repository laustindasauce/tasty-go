package tasty

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	// Default OAuth2 scopes for TastyTrade API
	defaultScope = "read trade"
	
	// OAuth2 grant types
	grantTypeAuthorizationCode = "authorization_code"
	grantTypeRefreshToken      = "refresh_token"
	
	// OAuth2 response types
	responseTypeCode = "code"
)

// OAuth2Client provides OAuth2 authentication functionality for the TastyTrade API
type OAuth2Client struct {
	config       OAuth2Config
	tokenManager *TokenManager
	httpClient   *http.Client
	pkce         *PKCEChallenge
}

// newOAuth2ClientInternal creates a new OAuth2Client with the provided configuration
func newOAuth2ClientInternal(config OAuth2Config, httpClient *http.Client) (*OAuth2Client, error) {
	if httpClient == nil {
		httpClient = defaultHTTPClient
	}
	
	// Validate the configuration
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid OAuth2 configuration: %w", err)
	}
	
	// Validate that PKCE is not being requested (since TastyTrade doesn't support it)
	if err := validateNoPKCERequested(config); err != nil {
		return nil, err
	}
	
	// Set default endpoints if not provided
	if config.AuthURL == "" || config.TokenURL == "" {
		if config.BaseURL == "" {
			// Default to production if not specified
			config.BaseURL = apiBaseURL
		}
		
		// Set appropriate endpoints based on base URL
		if config.BaseURL == apiCertBaseURL {
			config.AuthURL = oauth2SandboxAuthURL
			config.TokenURL = oauth2SandboxTokenURL
		} else {
			config.AuthURL = oauth2ProductionAuthURL
			config.TokenURL = oauth2ProductionTokenURL
		}
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
	
	// PKCE is not supported by TastyTrade - skip PKCE generation
	// This prevents the "405 Method Not Allowed" error from TastyTrade
	var pkce *PKCEChallenge = nil
	
	tokenManager := NewTokenManager()
	
	client := &OAuth2Client{
		config:       config,
		tokenManager: tokenManager,
		httpClient:   httpClient,
		pkce:         pkce,
	}
	
	// Set up token refresh callback
	tokenManager.SetRefreshCallback(client.refreshTokensInternal)
	
	return client, nil
}

// GetAuthorizationURL generates the OAuth2 authorization URL with state parameters
// Note: PKCE is not used because TastyTrade does not support it
func (c *OAuth2Client) GetAuthorizationURL() (string, error) {
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
	
	// Note: PKCE parameters (code_challenge, code_challenge_method) are intentionally
	// omitted because TastyTrade returns "405 Method Not Allowed" when they are included
	
	authURL.RawQuery = params.Encode()
	
	return authURL.String(), nil
}

// ExchangeCodeForTokens exchanges an authorization code for access and refresh tokens
func (c *OAuth2Client) ExchangeCodeForTokens(code string) (*TokenResponse, error) {
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
	
	// Note: code_verifier is intentionally omitted because TastyTrade doesn't support PKCE
	
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
	
	// Store tokens in token manager
	c.tokenManager.SetTokensFromResponse(&tokenResponse)
	
	return &tokenResponse, nil
}

// RefreshTokens refreshes the access token using the stored refresh token
func (c *OAuth2Client) RefreshTokens() (*TokenResponse, error) {
	refreshToken := c.tokenManager.GetRefreshToken()
	if refreshToken == "" {
		return nil, NewOAuth2Error(OAuth2ErrorRefreshFailed, "no refresh token available")
	}
	
	return c.refreshTokensWithRetry(refreshToken, 3)
}

// refreshTokensWithRetry attempts to refresh tokens with automatic retry logic
func (c *OAuth2Client) refreshTokensWithRetry(refreshToken string, maxRetries int) (*TokenResponse, error) {
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
func (c *OAuth2Client) performTokenRefresh(refreshToken string) (*TokenResponse, error) {
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
	
	// Store new tokens in token manager
	c.tokenManager.SetTokensFromResponse(&tokenResponse)
	
	return &tokenResponse, nil
}

// refreshTokensInternal is the internal callback used by TokenManager
func (c *OAuth2Client) refreshTokensInternal() (*TokenResponse, error) {
	return c.RefreshTokens()
}

// GetTokenManager returns the token manager for accessing token information
func (c *OAuth2Client) GetTokenManager() *TokenManager {
	return c.tokenManager
}

// GetConfig returns a copy of the OAuth2 configuration
func (c *OAuth2Client) GetConfig() OAuth2Config {
	return c.config
}

// GetState returns the current state parameter for CSRF protection
func (c *OAuth2Client) GetState() string {
	return c.config.State
}

// ValidateState validates that the provided state matches the expected state
func (c *OAuth2Client) ValidateState(state string) error {
	return ValidateState(c.config.State, state)
}

// IsAuthenticated checks if the client has valid authentication tokens
func (c *OAuth2Client) IsAuthenticated() bool {
	return c.tokenManager.HasValidToken() || c.tokenManager.HasRefreshToken()
}

// ClearTokens securely clears all stored authentication tokens
func (c *OAuth2Client) ClearTokens() {
	c.tokenManager.Clear()
}

// SetTokens stores OAuth2 tokens directly in the client for "bring your own tokens" usage
// This is the primary method for initializing the client with externally obtained tokens
// The expiresIn parameter is in seconds from now
func (c *OAuth2Client) SetTokens(accessToken, refreshToken string, expiresIn int) {
	c.tokenManager.SetTokens(accessToken, refreshToken, expiresIn)
}

// SetTokensFromResponse stores tokens from a TokenResponse object
// This is useful when you have a complete TokenResponse from an external OAuth2 flow
func (c *OAuth2Client) SetTokensFromResponse(response *TokenResponse) {
	c.tokenManager.SetTokensFromResponse(response)
}

// HasValidToken checks if the client has a valid (non-expired) access token
func (c *OAuth2Client) HasValidToken() bool {
	return c.tokenManager.HasValidToken()
}

// HasRefreshToken checks if the client has a refresh token available
func (c *OAuth2Client) HasRefreshToken() bool {
	return c.tokenManager.HasRefreshToken()
}

// GetTokenExpiration returns the expiration time of the current access token
func (c *OAuth2Client) GetTokenExpiration() time.Time {
	return c.tokenManager.GetExpiresAt()
}

// GetTimeUntilExpiry returns the duration until the current access token expires
func (c *OAuth2Client) GetTimeUntilExpiry() time.Duration {
	return c.tokenManager.GetTimeUntilExpiry()
}

// IsTokenExpired checks if the current access token is expired or about to expire
func (c *OAuth2Client) IsTokenExpired() bool {
	return c.tokenManager.IsExpired()
}

// generateSecureState generates a cryptographically secure state parameter
func generateSecureState() (string, error) {
	bytes := make([]byte, 32) // 32 bytes = 256 bits of entropy
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate secure random bytes: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(bytes), nil
}



// RedirectServer provides HTTP server functionality for handling OAuth2 redirects
type RedirectServer struct {
	server   *http.Server
	listener net.Listener
	codeChan chan string
	errChan  chan error
	state    string
	mutex    sync.RWMutex
	running  bool
}

// NewRedirectServer creates a new RedirectServer instance
func NewRedirectServer(state string) *RedirectServer {
	return &RedirectServer{
		codeChan: make(chan string, 1),
		errChan:  make(chan error, 1),
		state:    state,
	}
}

// Start starts the HTTP server on the specified port with timeout handling
// If port is 0, a random available port will be chosen
func (rs *RedirectServer) Start(port int) error {
	rs.mutex.Lock()
	defer rs.mutex.Unlock()
	
	if rs.running {
		return NewOAuth2Error(OAuth2ErrorConfigurationError, "redirect server is already running")
	}
	
	// Create listener on specified port
	addr := ":" + strconv.Itoa(port)
	if port == 0 {
		addr = ":0" // Let the system choose an available port
	}
	
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return NewOAuth2ErrorWithContext(OAuth2ErrorNetworkError, 
			fmt.Sprintf("failed to create listener on port %d", port), 0, err)
	}
	
	rs.listener = listener
	
	// Create HTTP server with redirect handler
	mux := http.NewServeMux()
	mux.HandleFunc("/", rs.handleRedirect)
	
	rs.server = &http.Server{
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}
	
	rs.running = true
	
	// Start server in goroutine
	go func() {
		if err := rs.server.Serve(rs.listener); err != nil && err != http.ErrServerClosed {
			rs.mutex.Lock()
			rs.running = false
			rs.mutex.Unlock()
			
			select {
			case rs.errChan <- NewOAuth2ErrorWithContext(OAuth2ErrorServerError, 
				"redirect server error", 0, err):
			default:
				// Channel is full, ignore
			}
		}
	}()
	
	return nil
}

// GetPort returns the port the server is listening on
func (rs *RedirectServer) GetPort() int {
	rs.mutex.RLock()
	defer rs.mutex.RUnlock()
	
	if rs.listener == nil {
		return 0
	}
	
	addr := rs.listener.Addr().(*net.TCPAddr)
	return addr.Port
}

// GetRedirectURI returns the full redirect URI for this server
func (rs *RedirectServer) GetRedirectURI() string {
	port := rs.GetPort()
	if port == 0 {
		return ""
	}
	return fmt.Sprintf("http://localhost:%d", port)
}

// WaitForCode waits for an authorization code with timeout protection
func (rs *RedirectServer) WaitForCode(timeout time.Duration) (string, error) {
	rs.mutex.RLock()
	if !rs.running {
		rs.mutex.RUnlock()
		return "", NewOAuth2Error(OAuth2ErrorConfigurationError, "redirect server is not running")
	}
	rs.mutex.RUnlock()
	
	// Create timeout context
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	
	select {
	case code := <-rs.codeChan:
		return code, nil
	case err := <-rs.errChan:
		if detailedErr, ok := err.(*OAuth2DetailedError); ok {
			return "", detailedErr
		}
		return "", NewOAuth2ErrorWithContext(OAuth2ErrorServerError, 
			"server error while waiting for code", 0, err)
	case <-ctx.Done():
		return "", NewOAuth2ErrorWithContext(OAuth2ErrorRedirectTimeout, 
			fmt.Sprintf("timeout waiting for authorization code after %v", timeout), 0, ctx.Err())
	}
}

// handleRedirect handles the OAuth2 redirect callback
func (rs *RedirectServer) handleRedirect(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	query := r.URL.Query()
	
	// Check for error parameter first
	if errorCode := query.Get("error"); errorCode != "" {
		errorDesc := query.Get("error_description")
		errorURI := query.Get("error_uri")
		
		oauthErr := OAuth2Error{
			ErrorCode:        errorCode,
			ErrorDescription: errorDesc,
			ErrorURI:         errorURI,
			State:            query.Get("state"),
		}
		
		// Convert to detailed error
		detailedErr := NewOAuth2ErrorFromStandard(oauthErr)
		
		// Send error response to user
		rs.sendErrorResponse(w, oauthErr)
		
		// Send detailed error to channel
		select {
		case rs.errChan <- detailedErr:
		default:
			// Channel is full, ignore
		}
		return
	}
	
	// Get authorization code
	code := query.Get("code")
	if code == "" {
		detailedErr := NewOAuth2Error(OAuth2ErrorMissingCode, "missing authorization code in redirect")
		rs.sendErrorResponse(w, OAuth2Error{
			ErrorCode:        OAuth2ErrorMissingCode,
			ErrorDescription: "Missing authorization code",
		})
		
		select {
		case rs.errChan <- detailedErr:
		default:
			// Channel is full, ignore
		}
		return
	}
	
	// Validate state parameter for CSRF protection
	receivedState := query.Get("state")
	if err := ValidateState(rs.state, receivedState); err != nil {
		rs.sendErrorResponse(w, OAuth2Error{
			ErrorCode:        OAuth2ErrorInvalidState,
			ErrorDescription: "Invalid state parameter",
		})
		
		select {
		case rs.errChan <- err:
		default:
			// Channel is full, ignore
		}
		return
	}
	
	// Send success response to user
	rs.sendSuccessResponse(w)
	
	// Send code to channel
	select {
	case rs.codeChan <- code:
	default:
		// Channel is full, ignore
	}
}

// sendSuccessResponse sends a success HTML response to the user
func (rs *RedirectServer) sendSuccessResponse(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	
	html := `<!DOCTYPE html>
<html>
<head>
    <title>Authorization Successful</title>
    <style>
        body { font-family: Arial, sans-serif; text-align: center; margin-top: 50px; }
        .success { color: #28a745; }
        .container { max-width: 500px; margin: 0 auto; padding: 20px; }
    </style>
</head>
<body>
    <div class="container">
        <h1 class="success">✓ Authorization Successful</h1>
        <p>You have successfully authorized the application.</p>
        <p>You can now close this window and return to the application.</p>
    </div>
    <script>
        // Auto-close window after 3 seconds if opened in popup
        if (window.opener) {
            setTimeout(function() {
                window.close();
            }, 3000);
        }
    </script>
</body>
</html>`
	
	w.Write([]byte(html))
}

// sendErrorResponse sends an error HTML response to the user
func (rs *RedirectServer) sendErrorResponse(w http.ResponseWriter, oauthErr OAuth2Error) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusBadRequest)
	
	html := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <title>Authorization Failed</title>
    <style>
        body { font-family: Arial, sans-serif; text-align: center; margin-top: 50px; }
        .error { color: #dc3545; }
        .container { max-width: 500px; margin: 0 auto; padding: 20px; }
        .error-details { background: #f8f9fa; padding: 15px; border-radius: 5px; margin-top: 20px; }
    </style>
</head>
<body>
    <div class="container">
        <h1 class="error">✗ Authorization Failed</h1>
        <p>There was an error during the authorization process.</p>
        <div class="error-details">
            <strong>Error:</strong> %s<br>
            %s
        </div>
        <p>Please close this window and try again.</p>
    </div>
    <script>
        // Auto-close window after 5 seconds if opened in popup
        if (window.opener) {
            setTimeout(function() {
                window.close();
            }, 5000);
        }
    </script>
</body>
</html>`, oauthErr.ErrorCode, 
		func() string {
			if oauthErr.ErrorDescription != "" {
				return fmt.Sprintf("<strong>Description:</strong> %s", oauthErr.ErrorDescription)
			}
			return ""
		}())
	
	w.Write([]byte(html))
}

// Shutdown gracefully shuts down the HTTP server with proper resource cleanup
func (rs *RedirectServer) Shutdown(timeout time.Duration) error {
	rs.mutex.Lock()
	defer rs.mutex.Unlock()
	
	if !rs.running {
		return nil // Already shut down
	}
	
	rs.running = false
	
	if rs.server == nil {
		return nil
	}
	
	// Create shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	
	// Attempt graceful shutdown
	if err := rs.server.Shutdown(ctx); err != nil {
		// Force close if graceful shutdown fails
		if rs.listener != nil {
			rs.listener.Close()
		}
		return NewOAuth2ErrorWithContext(OAuth2ErrorServerError, 
			"failed to shutdown server gracefully", 0, err)
	}
	
	// Close channels to prevent goroutine leaks
	close(rs.codeChan)
	close(rs.errChan)
	
	return nil
}

// IsRunning returns whether the redirect server is currently running
func (rs *RedirectServer) IsRunning() bool {
	rs.mutex.RLock()
	defer rs.mutex.RUnlock()
	return rs.running
}

// StartRedirectServer is a convenience method on OAuth2Client to start a redirect server
func (c *OAuth2Client) StartRedirectServer(port int) (*RedirectServer, error) {
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

// validateNoPKCERequested ensures that PKCE is not being requested since TastyTrade doesn't support it
func validateNoPKCERequested(config OAuth2Config) error {
	// For now, we don't have a direct PKCE flag in OAuth2Config, but we can add this validation
	// if needed in the future. The main protection is that GeneratePKCEChallenge() returns an error.
	return nil
}

// ValidateOAuth2Options validates OAuth2Options and ensures PKCE is not requested
func ValidateOAuth2Options(options OAuth2Options) error {
	if options.UsePKCE {
		return fmt.Errorf("PKCE is not supported by TastyTrade OAuth2 implementation - using PKCE parameters results in '405 Method Not Allowed' errors")
	}
	return nil
}