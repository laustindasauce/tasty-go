package tasty

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// OAuth2Config holds the configuration for OAuth2 authentication
type OAuth2Config struct {
	// The OAuth2 client ID provided by TastyTrade
	ClientID string `json:"client_id"`
	// The OAuth2 client secret provided by TastyTrade
	ClientSecret string `json:"client_secret"`
	// The redirect URI registered with TastyTrade for OAuth2 flow
	RedirectURI string `json:"redirect_uri"`
	// The requested OAuth2 scopes
	Scopes []string `json:"scopes"`
	// The base URL for OAuth2 endpoints (production or sandbox)
	BaseURL string `json:"base_url"`
	// State parameter for CSRF protection
	State string `json:"state"`
	// The OAuth2 authorization endpoint URL
	AuthURL string `json:"auth_url"`
	// The OAuth2 token endpoint URL
	TokenURL string `json:"token_url"`
}

// TokenResponse represents the response from OAuth2 token exchange
type TokenResponse struct {
	// The access token for API authentication
	AccessToken string `json:"access_token"`
	// The refresh token for obtaining new access tokens
	RefreshToken string `json:"refresh_token"`
	// The type of token (typically "Bearer")
	TokenType string `json:"token_type"`
	// The lifetime in seconds of the access token
	ExpiresIn int `json:"expires_in"`
	// Optional ID token for OpenID Connect
	IDToken string `json:"id_token,omitempty"`
	// The scope of the access token
	Scope string `json:"scope,omitempty"`
}

// OAuth2Error represents an OAuth2 error response
type OAuth2Error struct {
	// The error code as defined by OAuth2 specification
	ErrorCode string `json:"error"`
	// Human-readable description of the error
	ErrorDescription string `json:"error_description,omitempty"`
	// URI identifying a human-readable web page with error information
	ErrorURI string `json:"error_uri,omitempty"`
	// State parameter that was included in the authorization request
	State string `json:"state,omitempty"`
}

// Error implements the error interface for OAuth2Error
func (e OAuth2Error) Error() string {
	if e.ErrorDescription != "" {
		return fmt.Sprintf("OAuth2 error: %s - %s", e.ErrorCode, e.ErrorDescription)
	}
	return fmt.Sprintf("OAuth2 error: %s", e.ErrorCode)
}

// OAuth2Options provides configuration options for OAuth2 client setup
type OAuth2Options struct {
	// The OAuth2 client ID provided by TastyTrade
	ClientID string `json:"client_id"`
	// The OAuth2 client secret provided by TastyTrade
	ClientSecret string `json:"client_secret"`
	// The redirect URI registered with TastyTrade for OAuth2 flow
	RedirectURI string `json:"redirect_uri"`
	// The requested OAuth2 scopes
	Scopes []string `json:"scopes"`
	// DEPRECATED: PKCE is not supported by TastyTrade and will always be disabled.
	// This field is kept for API compatibility but has no effect.
	// Setting this to true will result in an error during client creation.
	UsePKCE bool `json:"use_pkce"`
	// Custom state parameter for CSRF protection (auto-generated if empty)
	State string `json:"state"`
	// Custom timeout for OAuth2 operations
	Timeout time.Duration `json:"timeout"`
}

// PKCEChallenge represents a PKCE challenge/verifier pair for OAuth2 security
// 
// IMPORTANT: TastyTrade's OAuth2 implementation does not currently support PKCE.
// Using PKCE parameters in authorization requests will result in "405 Method Not Allowed" errors.
// This implementation is provided for future compatibility and for use with other OAuth2 providers
// that do support PKCE.
type PKCEChallenge struct {
	// The code verifier - a cryptographically random string
	CodeVerifier string `json:"code_verifier"`
	// The code challenge - SHA256 hash of the code verifier, base64url encoded
	CodeChallenge string `json:"code_challenge"`
	// The method used to derive the challenge from the verifier (always "S256")
	Method string `json:"code_challenge_method"`
}

// GeneratePKCEChallenge creates a new PKCE challenge/verifier pair using SHA256
// and cryptographically secure random generation as required by RFC 7636
//
// DEPRECATED: This function is disabled because TastyTrade does not support PKCE.
// Calling this function will return an error. The implementation is kept for
// future compatibility if TastyTrade adds PKCE support.
func GeneratePKCEChallenge() (*PKCEChallenge, error) {
	return nil, fmt.Errorf("PKCE is not supported by TastyTrade OAuth2 implementation - using PKCE parameters results in '405 Method Not Allowed' errors")
}

// Validate verifies that the provided verifier matches this PKCE challenge
// This method should be used to validate the code verifier during token exchange
//
// DEPRECATED: Always returns false because PKCE is not supported by TastyTrade.
func (p *PKCEChallenge) Validate(verifier string) bool {
	// Always return false since PKCE is disabled
	return false
}

// IsValid checks if the PKCE challenge is properly formed and valid
//
// DEPRECATED: Always returns false because PKCE is not supported by TastyTrade.
func (p *PKCEChallenge) IsValid() bool {
	// Always return false since PKCE is disabled
	return false
}

// TokenStorage defines the interface for token storage backends
type TokenStorage interface {
	// Store saves token data to the storage backend
	Store(data *TokenData) error
	// Load retrieves token data from the storage backend
	Load() (*TokenData, error)
	// Clear removes all token data from the storage backend
	Clear() error
}

// TokenData represents the token information stored by TokenStorage implementations
type TokenData struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	TokenType    string    `json:"token_type"`
	ExpiresAt    time.Time `json:"expires_at"`
	Scope        string    `json:"scope"`
}

// MemoryTokenStorage provides in-memory token storage (default implementation)
type MemoryTokenStorage struct {
	mutex sync.RWMutex
	data  *TokenData
}

// NewMemoryTokenStorage creates a new in-memory token storage
func NewMemoryTokenStorage() *MemoryTokenStorage {
	return &MemoryTokenStorage{}
}

// Store saves token data in memory
func (m *MemoryTokenStorage) Store(data *TokenData) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	// Create a copy to avoid external modifications
	m.data = &TokenData{
		AccessToken:  data.AccessToken,
		RefreshToken: data.RefreshToken,
		TokenType:    data.TokenType,
		ExpiresAt:    data.ExpiresAt,
		Scope:        data.Scope,
	}
	return nil
}

// Load retrieves token data from memory
func (m *MemoryTokenStorage) Load() (*TokenData, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	
	if m.data == nil {
		return nil, nil
	}
	
	// Return a copy to avoid external modifications
	return &TokenData{
		AccessToken:  m.data.AccessToken,
		RefreshToken: m.data.RefreshToken,
		TokenType:    m.data.TokenType,
		ExpiresAt:    m.data.ExpiresAt,
		Scope:        m.data.Scope,
	}, nil
}

// Clear removes token data from memory with secure overwriting
func (m *MemoryTokenStorage) Clear() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	if m.data != nil {
		// Securely overwrite sensitive data
		if m.data.AccessToken != "" {
			tokenBytes := []byte(m.data.AccessToken)
			for i := range tokenBytes {
				tokenBytes[i] = 0
			}
		}
		if m.data.RefreshToken != "" {
			tokenBytes := []byte(m.data.RefreshToken)
			for i := range tokenBytes {
				tokenBytes[i] = 0
			}
		}
		m.data = nil
	}
	return nil
}

// FileTokenStorage provides file-based token storage with encryption
type FileTokenStorage struct {
	mutex    sync.RWMutex
	filePath string
}

// NewFileTokenStorage creates a new file-based token storage
func NewFileTokenStorage(filePath string) *FileTokenStorage {
	return &FileTokenStorage{
		filePath: filePath,
	}
}

// Store saves token data to file
func (f *FileTokenStorage) Store(data *TokenData) error {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	
	// Marshal token data to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal token data: %w", err)
	}
	
	// Write to file with restricted permissions (0600 - owner read/write only)
	err = os.WriteFile(f.filePath, jsonData, 0600)
	if err != nil {
		return fmt.Errorf("failed to write token file: %w", err)
	}
	
	return nil
}

// Load retrieves token data from file
func (f *FileTokenStorage) Load() (*TokenData, error) {
	f.mutex.RLock()
	defer f.mutex.RUnlock()
	
	// Check if file exists
	if _, err := os.Stat(f.filePath); os.IsNotExist(err) {
		return nil, nil // No token file exists
	}
	
	// Read file contents
	jsonData, err := os.ReadFile(f.filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read token file: %w", err)
	}
	
	// Unmarshal JSON data
	var data TokenData
	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal token data: %w", err)
	}
	
	return &data, nil
}

// Clear removes the token file
func (f *FileTokenStorage) Clear() error {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	
	// Check if file exists
	if _, err := os.Stat(f.filePath); os.IsNotExist(err) {
		return nil // File doesn't exist, nothing to clear
	}
	
	// Remove the file
	err := os.Remove(f.filePath)
	if err != nil {
		return fmt.Errorf("failed to remove token file: %w", err)
	}
	
	return nil
}

// TokenManager provides thread-safe storage and lifecycle management for OAuth2 tokens
type TokenManager struct {
	// Thread-safe access to token operations
	mutex sync.RWMutex
	
	// Storage backend for token persistence
	storage TokenStorage
	
	// Configuration for token refresh
	refreshCallback func() (*TokenResponse, error)
}

// NewTokenManager creates a new TokenManager instance with file-based storage
// Tokens are stored in ~/.tasty-go/tokens.json by default
func NewTokenManager() *TokenManager {
	return NewFileTokenManager(getDefaultTokenPath())
}

// NewTokenManagerWithStorage creates a new TokenManager instance with custom storage
func NewTokenManagerWithStorage(storage TokenStorage) *TokenManager {
	return &TokenManager{
		storage: storage,
	}
}

// NewFileTokenManager creates a new TokenManager instance with file-based storage
func NewFileTokenManager(filePath string) *TokenManager {
	return &TokenManager{
		storage: NewFileTokenStorage(filePath),
	}
}

// NewMemoryTokenManager creates a new TokenManager instance with in-memory storage
// This is useful for testing or when you don't want persistent token storage
func NewMemoryTokenManager() *TokenManager {
	return &TokenManager{
		storage: NewMemoryTokenStorage(),
	}
}

// getDefaultTokenPath returns the default path for storing OAuth2 tokens
func getDefaultTokenPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// Fallback to current directory if home directory is not available
		return ".tasty-go-tokens.json"
	}
	
	// Create .tasty-go directory in user's home directory
	configDir := filepath.Join(homeDir, ".tasty-go")
	os.MkdirAll(configDir, 0700) // Create directory with secure permissions
	
	return filepath.Join(configDir, "tokens.json")
}

// SetTokens stores new OAuth2 tokens with thread-safe access
// The expiresIn parameter is in seconds from now
func (tm *TokenManager) SetTokens(accessToken, refreshToken string, expiresIn int) {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()
	
	data := &TokenData{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer", // Default token type
		ExpiresAt:    time.Now().Add(time.Duration(expiresIn) * time.Second),
	}
	
	tm.storage.Store(data)
}

// SetTokensFromResponse stores tokens from a TokenResponse with thread-safe access
func (tm *TokenManager) SetTokensFromResponse(response *TokenResponse) {
	if response == nil {
		return
	}
	
	tm.mutex.Lock()
	defer tm.mutex.Unlock()
	
	tokenType := response.TokenType
	if tokenType == "" {
		tokenType = "Bearer" // Default token type
	}
	
	data := &TokenData{
		AccessToken:  response.AccessToken,
		RefreshToken: response.RefreshToken,
		TokenType:    tokenType,
		ExpiresAt:    time.Now().Add(time.Duration(response.ExpiresIn) * time.Second),
		Scope:        response.Scope,
	}
	
	tm.storage.Store(data)
}

// GetAccessToken returns the current access token if valid, or attempts refresh
// Returns an error if no valid token is available and refresh fails
func (tm *TokenManager) GetAccessToken() (string, error) {
	tm.mutex.RLock()
	
	// Load current token data
	data, err := tm.storage.Load()
	if err != nil {
		tm.mutex.RUnlock()
		return "", fmt.Errorf("failed to load token data: %w", err)
	}
	
	// Check if current token is still valid (with 30 second buffer)
	if data != nil && data.AccessToken != "" && time.Now().Add(30*time.Second).Before(data.ExpiresAt) {
		token := data.AccessToken
		tm.mutex.RUnlock()
		return token, nil
	}
	
	// Token is expired or about to expire, need to refresh
	var refreshToken string
	if data != nil {
		refreshToken = data.RefreshToken
	}
	callback := tm.refreshCallback
	tm.mutex.RUnlock()
	
	// Attempt automatic refresh if callback is available
	if callback != nil && refreshToken != "" {
		if err := tm.refreshTokens(); err != nil {
			return "", fmt.Errorf("failed to refresh access token: %w", err)
		}
		
		// Return the new token after successful refresh
		tm.mutex.RLock()
		data, err := tm.storage.Load()
		tm.mutex.RUnlock()
		if err != nil {
			return "", fmt.Errorf("failed to load refreshed token: %w", err)
		}
		if data != nil {
			return data.AccessToken, nil
		}
	}
	
	return "", fmt.Errorf("access token expired and no refresh mechanism available")
}

// GetRefreshToken returns the current refresh token with thread-safe access
func (tm *TokenManager) GetRefreshToken() string {
	tm.mutex.RLock()
	defer tm.mutex.RUnlock()
	
	data, err := tm.storage.Load()
	if err != nil || data == nil {
		return ""
	}
	return data.RefreshToken
}

// GetTokenType returns the token type (typically "Bearer")
func (tm *TokenManager) GetTokenType() string {
	tm.mutex.RLock()
	defer tm.mutex.RUnlock()
	
	data, err := tm.storage.Load()
	if err != nil || data == nil {
		return "Bearer" // Default token type
	}
	if data.TokenType == "" {
		return "Bearer" // Default token type
	}
	return data.TokenType
}

// GetScope returns the token scope
func (tm *TokenManager) GetScope() string {
	tm.mutex.RLock()
	defer tm.mutex.RUnlock()
	
	data, err := tm.storage.Load()
	if err != nil || data == nil {
		return ""
	}
	return data.Scope
}

// IsExpired checks if the access token is expired or about to expire (within 30 seconds)
func (tm *TokenManager) IsExpired() bool {
	tm.mutex.RLock()
	defer tm.mutex.RUnlock()
	
	data, err := tm.storage.Load()
	if err != nil || data == nil || data.AccessToken == "" {
		return true
	}
	
	// Consider token expired if it expires within 30 seconds
	return time.Now().Add(30 * time.Second).After(data.ExpiresAt)
}

// HasValidToken checks if there's a valid access token available
func (tm *TokenManager) HasValidToken() bool {
	return !tm.IsExpired()
}

// HasRefreshToken checks if a refresh token is available
func (tm *TokenManager) HasRefreshToken() bool {
	tm.mutex.RLock()
	defer tm.mutex.RUnlock()
	
	data, err := tm.storage.Load()
	if err != nil || data == nil {
		return false
	}
	return data.RefreshToken != ""
}

// GetExpiresAt returns the expiration time of the current access token
func (tm *TokenManager) GetExpiresAt() time.Time {
	tm.mutex.RLock()
	defer tm.mutex.RUnlock()
	
	data, err := tm.storage.Load()
	if err != nil || data == nil {
		return time.Time{}
	}
	return data.ExpiresAt
}

// GetTimeUntilExpiry returns the duration until the token expires
func (tm *TokenManager) GetTimeUntilExpiry() time.Duration {
	tm.mutex.RLock()
	defer tm.mutex.RUnlock()
	
	data, err := tm.storage.Load()
	if err != nil || data == nil || data.AccessToken == "" {
		return 0
	}
	
	remaining := data.ExpiresAt.Sub(time.Now())
	if remaining < 0 {
		return 0
	}
	return remaining
}

// SetRefreshCallback sets the callback function for automatic token refresh
func (tm *TokenManager) SetRefreshCallback(callback func() (*TokenResponse, error)) {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()
	tm.refreshCallback = callback
}

// refreshTokens performs the actual token refresh using the callback
func (tm *TokenManager) refreshTokens() error {
	tm.mutex.RLock()
	callback := tm.refreshCallback
	tm.mutex.RUnlock()
	
	if callback == nil {
		return fmt.Errorf("no refresh callback configured")
	}
	
	response, err := callback()
	if err != nil {
		return fmt.Errorf("refresh callback failed: %w", err)
	}
	
	if response == nil {
		return fmt.Errorf("refresh callback returned nil response")
	}
	
	// Update tokens with new response
	tm.SetTokensFromResponse(response)
	return nil
}

// Clear securely clears all stored token data
func (tm *TokenManager) Clear() {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()
	
	// Clear tokens from storage backend
	tm.storage.Clear()
	tm.refreshCallback = nil
}

// Clone creates a copy of the TokenManager with the same token data
// This is useful for creating independent instances while preserving token state
// Note: Clones use in-memory storage to avoid file conflicts
func (tm *TokenManager) Clone() *TokenManager {
	tm.mutex.RLock()
	defer tm.mutex.RUnlock()
	
	// Load current data
	data, err := tm.storage.Load()
	if err != nil || data == nil {
		// Return empty clone if no data or error
		return NewMemoryTokenManager()
	}
	
	// Create clone with memory storage and copy data
	clone := NewMemoryTokenManager()
	
	// Store the data in the clone
	clone.storage.Store(data)
	
	return clone
}

// NewProductionOAuth2Config creates an OAuth2Config for the production environment
func NewProductionOAuth2Config(clientID, clientSecret, redirectURI string, scopes []string) OAuth2Config {
	if len(scopes) == 0 {
		scopes = []string{defaultScope}
	}
	
	return OAuth2Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURI:  redirectURI,
		Scopes:       scopes,
		BaseURL:      apiBaseURL,
		AuthURL:      oauth2ProductionAuthURL,
		TokenURL:     oauth2ProductionTokenURL,
	}
}

// NewSandboxOAuth2Config creates an OAuth2Config for the sandbox environment
func NewSandboxOAuth2Config(clientID, clientSecret, redirectURI string, scopes []string) OAuth2Config {
	if len(scopes) == 0 {
		scopes = []string{defaultScope}
	}
	
	return OAuth2Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURI:  redirectURI,
		Scopes:       scopes,
		BaseURL:      apiCertBaseURL,
		AuthURL:      oauth2SandboxAuthURL,
		TokenURL:     oauth2SandboxTokenURL,
	}
}

// IsProduction returns true if this configuration is for the production environment
func (c OAuth2Config) IsProduction() bool {
	return c.BaseURL == apiBaseURL || c.AuthURL == oauth2ProductionAuthURL || c.TokenURL == oauth2ProductionTokenURL
}

// IsSandbox returns true if this configuration is for the sandbox environment
func (c OAuth2Config) IsSandbox() bool {
	return c.BaseURL == apiCertBaseURL || c.AuthURL == oauth2SandboxAuthURL || c.TokenURL == oauth2SandboxTokenURL
}

// GetEnvironment returns a string representation of the environment ("production" or "sandbox")
func (c OAuth2Config) GetEnvironment() string {
	if c.IsProduction() {
		return "production"
	} else if c.IsSandbox() {
		return "sandbox"
	}
	return "unknown"
}

// Validate performs comprehensive validation of the OAuth2 configuration
func (c OAuth2Config) Validate() error {
	// Validate required fields
	if c.ClientID == "" {
		return fmt.Errorf("OAuth2 client ID is required")
	}
	
	if c.RedirectURI == "" {
		return fmt.Errorf("OAuth2 redirect URI is required")
	}
	
	// Validate redirect URI format
	if err := validateRedirectURI(c.RedirectURI); err != nil {
		return fmt.Errorf("invalid redirect URI: %w", err)
	}
	
	// Validate endpoint URLs if provided
	if c.AuthURL != "" {
		if err := validateEndpointURL(c.AuthURL, "authorization"); err != nil {
			return fmt.Errorf("invalid authorization URL: %w", err)
		}
	}
	
	if c.TokenURL != "" {
		if err := validateEndpointURL(c.TokenURL, "token"); err != nil {
			return fmt.Errorf("invalid token URL: %w", err)
		}
	}
	
	if c.BaseURL != "" {
		if err := validateEndpointURL(c.BaseURL, "base"); err != nil {
			return fmt.Errorf("invalid base URL: %w", err)
		}
	}
	
	// Validate environment consistency
	if err := c.validateEnvironmentConsistency(); err != nil {
		return fmt.Errorf("environment configuration inconsistency: %w", err)
	}
	
	return nil
}

// validateEnvironmentConsistency ensures all URLs are consistent with the same environment
func (c OAuth2Config) validateEnvironmentConsistency() error {
	// If no URLs are set, configuration is valid
	if c.BaseURL == "" && c.AuthURL == "" && c.TokenURL == "" {
		return nil
	}
	
	// Check if configuration mixes production and sandbox URLs
	hasProduction := false
	hasSandbox := false
	
	if c.BaseURL != "" {
		if c.BaseURL == apiBaseURL {
			hasProduction = true
		} else if c.BaseURL == apiCertBaseURL {
			hasSandbox = true
		}
	}
	
	if c.AuthURL != "" {
		if c.AuthURL == oauth2ProductionAuthURL {
			hasProduction = true
		} else if c.AuthURL == oauth2SandboxAuthURL {
			hasSandbox = true
		}
	}
	
	if c.TokenURL != "" {
		if c.TokenURL == oauth2ProductionTokenURL {
			hasProduction = true
		} else if c.TokenURL == oauth2SandboxTokenURL {
			hasSandbox = true
		}
	}
	
	// Error if mixing environments
	if hasProduction && hasSandbox {
		return fmt.Errorf("configuration mixes production and sandbox endpoints")
	}
	
	return nil
}

// SetEnvironment configures the OAuth2Config for the specified environment
func (c *OAuth2Config) SetEnvironment(environment string) error {
	switch environment {
	case "production", "prod":
		c.BaseURL = apiBaseURL
		c.AuthURL = oauth2ProductionAuthURL
		c.TokenURL = oauth2ProductionTokenURL
	case "sandbox", "cert":
		c.BaseURL = apiCertBaseURL
		c.AuthURL = oauth2SandboxAuthURL
		c.TokenURL = oauth2SandboxTokenURL
	default:
		return fmt.Errorf("invalid environment: %s (must be 'production' or 'sandbox')", environment)
	}
	return nil
}

// validateRedirectURI validates the format and security of a redirect URI
func validateRedirectURI(redirectURI string) error {
	if redirectURI == "" {
		return fmt.Errorf("redirect URI cannot be empty")
	}
	
	// Parse the URI
	parsedURL, err := url.Parse(redirectURI)
	if err != nil {
		return fmt.Errorf("invalid URI format: %w", err)
	}
	
	// Check scheme
	switch parsedURL.Scheme {
	case "http":
		// HTTP is only allowed for localhost for development
		if parsedURL.Hostname() != "localhost" && parsedURL.Hostname() != "127.0.0.1" {
			return fmt.Errorf("HTTP redirect URIs are only allowed for localhost")
		}
	case "https":
		// HTTPS is always allowed
	case "":
		return fmt.Errorf("redirect URI must include a scheme (http:// or https://)")
	default:
		// Custom schemes are allowed for mobile apps
		if len(parsedURL.Scheme) < 3 {
			return fmt.Errorf("custom URI schemes must be at least 3 characters long")
		}
	}
	
	// Validate that URI is not just a scheme
	if parsedURL.Host == "" && parsedURL.Path == "" {
		return fmt.Errorf("redirect URI must include a host or path")
	}
	
	return nil
}

// validateEndpointURL validates OAuth2 endpoint URLs
func validateEndpointURL(endpointURL, endpointType string) error {
	if endpointURL == "" {
		return fmt.Errorf("%s URL cannot be empty", endpointType)
	}
	
	// Parse the URL
	parsedURL, err := url.Parse(endpointURL)
	if err != nil {
		return fmt.Errorf("invalid %s URL format: %w", endpointType, err)
	}
	
	// Must have a scheme
	if parsedURL.Scheme == "" {
		return fmt.Errorf("%s URL must include a scheme", endpointType)
	}
	
	// Must have a host
	if parsedURL.Host == "" {
		return fmt.Errorf("%s URL must include a host", endpointType)
	}
	
	// Check if this is a TastyTrade production domain
	tastyTradeHosts := []string{
		"api.tastyworks.com",
		"api.cert.tastyworks.com",
		"my.tastytrade.com",
		"cert-my.staging-tasty.works",
	}
	
	isTastyTradeDomain := false
	for _, host := range tastyTradeHosts {
		if parsedURL.Host == host {
			isTastyTradeDomain = true
			break
		}
	}
	
	// TastyTrade domains must always use HTTPS
	if isTastyTradeDomain && parsedURL.Scheme != "https" {
		return fmt.Errorf("%s URL must use HTTPS", endpointType)
	}
	
	// For production use, require HTTPS and valid TastyTrade endpoints
	// Allow HTTP and test domains for development/testing
	if parsedURL.Scheme == "https" {
		// HTTPS URLs should be valid TastyTrade endpoints for production use or test domains
		if !isTastyTradeDomain && !isTestDomain(parsedURL.Host) && !isLocalhost(parsedURL.Host) {
			return fmt.Errorf("%s URL must be a valid TastyTrade endpoint or test domain", endpointType)
		}
	} else if parsedURL.Scheme == "http" {
		// HTTP is only allowed for localhost and test domains (not TastyTrade domains)
		if !isLocalhost(parsedURL.Host) && !isTestDomain(parsedURL.Host) {
			return fmt.Errorf("%s URL with HTTP scheme is only allowed for localhost or test domains", endpointType)
		}
	} else {
		return fmt.Errorf("%s URL must use HTTP or HTTPS scheme", endpointType)
	}
	
	return nil
}

// isLocalhost checks if the host is localhost or 127.0.0.1
func isLocalhost(host string) bool {
	return host == "localhost" || host == "127.0.0.1" || 
		   strings.HasPrefix(host, "localhost:") || strings.HasPrefix(host, "127.0.0.1:")
}

// isTestDomain checks if the host is a test domain (for testing purposes)
func isTestDomain(host string) bool {
	testDomains := []string{
		"api.test.com",
		"test.example.com",
		"example.com",
		"httpbin.org",
	}
	
	for _, testDomain := range testDomains {
		if host == testDomain || strings.HasSuffix(host, "."+testDomain) {
			return true
		}
	}
	
	// Allow any .test TLD for testing
	return strings.HasSuffix(host, ".test")
}

// SupportsPKCE returns whether the OAuth2 configuration supports PKCE
// 
// DEPRECATED: Always returns false because PKCE is completely disabled.
// TastyTrade does not support PKCE and this library is specifically for TastyTrade.
func (c OAuth2Config) SupportsPKCE() bool {
	// PKCE is completely disabled for this TastyTrade-specific library
	return false
}