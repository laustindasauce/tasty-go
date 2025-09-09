package tasty

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// OAuth2Utils provides utility functions for common OAuth2 operations
type OAuth2Utils struct {
	// Configuration for utility behavior
	EnableDebugLogging bool
	DefaultTimeout     time.Duration
	DefaultPort        int
}

// NewOAuth2Utils creates a new OAuth2Utils instance with default settings
func NewOAuth2Utils() *OAuth2Utils {
	return &OAuth2Utils{
		EnableDebugLogging: false,
		DefaultTimeout:     30 * time.Second,
		DefaultPort:        8080,
	}
}

// OAuth2FlowOptions provides configuration options for different OAuth2 flow variations
type OAuth2FlowOptions struct {
	// Basic configuration
	ClientID     string
	ClientSecret string
	RedirectURI  string
	Scopes       []string
	Environment  string // "production" or "sandbox"
	
	// Flow behavior options
	UseBuiltinServer bool          // Whether to use the built-in redirect server
	ServerPort       int           // Port for the built-in server (0 for random)
	Timeout          time.Duration // Timeout for the authorization flow
	
	// Security options
	UsePKCE      bool   // Whether to use PKCE (recommended)
	CustomState  string // Custom state parameter (auto-generated if empty)
	
	// Advanced options
	CustomAuthURL  string // Custom authorization URL (overrides environment)
	CustomTokenURL string // Custom token URL (overrides environment)
	HTTPClient     *http.Client // Custom HTTP client
}

// DefaultFlowOptions returns OAuth2FlowOptions with sensible defaults
func (u *OAuth2Utils) DefaultFlowOptions() OAuth2FlowOptions {
	return OAuth2FlowOptions{
		Environment:      "production",
		UseBuiltinServer: true,
		ServerPort:       u.DefaultPort,
		Timeout:          u.DefaultTimeout,
		UsePKCE:          true,
		Scopes:           []string{"read", "trade"},
	}
}

// CreateQuickOAuth2Flow sets up a complete OAuth2 flow with minimal configuration
// This is a convenience method for the most common use case
func (u *OAuth2Utils) CreateQuickOAuth2Flow(clientID, clientSecret string, options ...OAuth2FlowOptions) (*OAuth2Client, *RedirectServer, error) {
	// Use default options if none provided
	var opts OAuth2FlowOptions
	if len(options) > 0 {
		opts = options[0]
	} else {
		opts = u.DefaultFlowOptions()
	}
	
	// Set required fields
	opts.ClientID = clientID
	opts.ClientSecret = clientSecret
	
	// Set default redirect URI if not provided
	if opts.RedirectURI == "" {
		if opts.UseBuiltinServer {
			port := opts.ServerPort
			if port == 0 {
				port = u.DefaultPort
			}
			opts.RedirectURI = fmt.Sprintf("http://localhost:%d", port)
		} else {
			return nil, nil, NewOAuth2Error(OAuth2ErrorConfigurationError, 
				"redirect URI is required when not using built-in server")
		}
	}
	
	// Create OAuth2 configuration
	config, err := u.CreateConfigFromOptions(opts)
	if err != nil {
		return nil, nil, err
	}
	
	// Create OAuth2 client
	client, err := newOAuth2ClientInternal(config, opts.HTTPClient)
	if err != nil {
		return nil, nil, err
	}
	
	// Start redirect server if requested
	var server *RedirectServer
	if opts.UseBuiltinServer {
		server, err = client.StartRedirectServer(opts.ServerPort)
		if err != nil {
			return nil, nil, err
		}
		
		// Update redirect URI with actual port if random port was used
		if opts.ServerPort == 0 {
			actualPort := server.GetPort()
			config.RedirectURI = fmt.Sprintf("http://localhost:%d", actualPort)
			// Update client configuration
			client.config.RedirectURI = config.RedirectURI
		}
	}
	
	return client, server, nil
}

// CreateConfigFromOptions creates an OAuth2Config from OAuth2FlowOptions
func (u *OAuth2Utils) CreateConfigFromOptions(opts OAuth2FlowOptions) (OAuth2Config, error) {
	config := OAuth2Config{
		ClientID:     opts.ClientID,
		ClientSecret: opts.ClientSecret,
		RedirectURI:  opts.RedirectURI,
		Scopes:       opts.Scopes,
		State:        opts.CustomState,
	}
	
	// Set environment-specific URLs
	if opts.CustomAuthURL != "" && opts.CustomTokenURL != "" {
		config.AuthURL = opts.CustomAuthURL
		config.TokenURL = opts.CustomTokenURL
		config.BaseURL = extractBaseURL(opts.CustomAuthURL)
	} else {
		switch strings.ToLower(opts.Environment) {
		case "production", "prod":
			config.BaseURL = apiBaseURL
			config.AuthURL = oauth2ProductionAuthURL
			config.TokenURL = oauth2ProductionTokenURL
		case "sandbox", "cert":
			config.BaseURL = apiCertBaseURL
			config.AuthURL = oauth2SandboxAuthURL
			config.TokenURL = oauth2SandboxTokenURL
		default:
			return config, NewOAuth2Error(OAuth2ErrorConfigurationError, 
				fmt.Sprintf("invalid environment: %s (must be 'production' or 'sandbox')", opts.Environment))
		}
	}
	
	// Validate the configuration
	if err := config.Validate(); err != nil {
		return config, err
	}
	
	return config, nil
}

// PerformCompleteFlow executes a complete OAuth2 authorization flow
// This is a high-level convenience method that handles the entire process
func (u *OAuth2Utils) PerformCompleteFlow(ctx context.Context, opts OAuth2FlowOptions) (*OAuth2Client, error) {
	// Create OAuth2 client and server
	client, server, err := u.CreateQuickOAuth2Flow(opts.ClientID, opts.ClientSecret, opts)
	if err != nil {
		return nil, err
	}
	
	// Clean up server when done
	if server != nil {
		defer func() {
			if shutdownErr := server.Shutdown(5 * time.Second); shutdownErr != nil {
				u.logDebug("Failed to shutdown redirect server: %v", shutdownErr)
			}
		}()
	}
	
	// Get authorization URL
	authURL, err := client.GetAuthorizationURL()
	if err != nil {
		return nil, err
	}
	
	u.logDebug("Authorization URL: %s", authURL)
	
	// If using built-in server, wait for callback
	if server != nil {
		u.logDebug("Waiting for authorization callback...")
		
		// Create timeout context
		timeoutCtx, cancel := context.WithTimeout(ctx, opts.Timeout)
		defer cancel()
		
		// Wait for authorization code
		var code string
		select {
		case <-timeoutCtx.Done():
			return nil, NewOAuth2ErrorWithContext(OAuth2ErrorRedirectTimeout, 
				fmt.Sprintf("timeout waiting for authorization after %v", opts.Timeout), 0, timeoutCtx.Err())
		default:
			code, err = server.WaitForCode(opts.Timeout)
			if err != nil {
				return nil, err
			}
		}
		
		u.logDebug("Received authorization code: %s", u.maskSensitiveData(code))
		
		// Exchange code for tokens
		tokenResponse, err := client.ExchangeCodeForTokens(code)
		if err != nil {
			return nil, err
		}
		
		u.logDebug("Token exchange successful, expires in %d seconds", tokenResponse.ExpiresIn)
		
		return client, nil
	}
	
	// If not using built-in server, return client for manual handling
	return client, nil
}

// CreateCustomRedirectHandler creates a custom HTTP handler for OAuth2 redirects
// This is useful when you want to integrate OAuth2 handling into an existing web server
func (u *OAuth2Utils) CreateCustomRedirectHandler(expectedState string, codeChan chan<- string, errChan chan<- error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		
		// Check for error parameter
		if errorCode := query.Get("error"); errorCode != "" {
			errorDesc := query.Get("error_description")
			errorURI := query.Get("error_uri")
			
			oauthErr := OAuth2Error{
				ErrorCode:        errorCode,
				ErrorDescription: errorDesc,
				ErrorURI:         errorURI,
				State:            query.Get("state"),
			}
			
			// Send error response to user
			u.sendCustomErrorResponse(w, oauthErr)
			
			// Send error to channel
			select {
			case errChan <- NewOAuth2ErrorFromStandard(oauthErr):
			default:
				// Channel is full or closed
			}
			return
		}
		
		// Get authorization code
		code := query.Get("code")
		if code == "" {
			err := NewOAuth2Error(OAuth2ErrorMissingCode, "missing authorization code in redirect")
			u.sendCustomErrorResponse(w, OAuth2Error{
				ErrorCode:        OAuth2ErrorMissingCode,
				ErrorDescription: "Missing authorization code",
			})
			
			select {
			case errChan <- err:
			default:
				// Channel is full or closed
			}
			return
		}
		
		// Validate state parameter
		receivedState := query.Get("state")
		if err := ValidateState(expectedState, receivedState); err != nil {
			u.sendCustomErrorResponse(w, OAuth2Error{
				ErrorCode:        OAuth2ErrorInvalidState,
				ErrorDescription: "Invalid state parameter",
			})
			
			select {
			case errChan <- err:
			default:
				// Channel is full or closed
			}
			return
		}
		
		// Send success response
		u.sendCustomSuccessResponse(w)
		
		// Send code to channel
		select {
		case codeChan <- code:
		default:
			// Channel is full or closed
		}
	}
}

// ParseRedirectURL parses an OAuth2 redirect URL and extracts the authorization code or error
func (u *OAuth2Utils) ParseRedirectURL(redirectURL, expectedState string) (string, error) {
	parsedURL, err := url.Parse(redirectURL)
	if err != nil {
		return "", NewOAuth2ErrorWithContext(OAuth2ErrorConfigurationError, 
			"invalid redirect URL format", 0, err)
	}
	
	query := parsedURL.Query()
	
	// Check for error parameter
	if errorCode := query.Get("error"); errorCode != "" {
		errorDesc := query.Get("error_description")
		errorURI := query.Get("error_uri")
		
		oauthErr := OAuth2Error{
			ErrorCode:        errorCode,
			ErrorDescription: errorDesc,
			ErrorURI:         errorURI,
			State:            query.Get("state"),
		}
		
		return "", NewOAuth2ErrorFromStandard(oauthErr)
	}
	
	// Get authorization code
	code := query.Get("code")
	if code == "" {
		return "", NewOAuth2Error(OAuth2ErrorMissingCode, "missing authorization code in redirect URL")
	}
	
	// Validate state parameter
	receivedState := query.Get("state")
	if err := ValidateState(expectedState, receivedState); err != nil {
		return "", err
	}
	
	return code, nil
}

// GenerateSecureState generates a cryptographically secure state parameter
func (u *OAuth2Utils) GenerateSecureState() (string, error) {
	bytes := make([]byte, 32) // 32 bytes = 256 bits of entropy
	if _, err := rand.Read(bytes); err != nil {
		return "", NewOAuth2ErrorWithContext(OAuth2ErrorConfigurationError, 
			"failed to generate secure random bytes", 0, err)
	}
	return base64.RawURLEncoding.EncodeToString(bytes), nil
}

// BuildAuthorizationURL builds an OAuth2 authorization URL with the provided parameters
// This is a lower-level utility for custom OAuth2 implementations
func (u *OAuth2Utils) BuildAuthorizationURL(config OAuth2Config, pkce *PKCEChallenge) (string, error) {
	if config.AuthURL == "" {
		return "", NewOAuth2Error(OAuth2ErrorConfigurationError, "authorization URL not configured")
	}
	
	authURL, err := url.Parse(config.AuthURL)
	if err != nil {
		return "", NewOAuth2ErrorWithContext(OAuth2ErrorConfigurationError, 
			"failed to parse authorization URL", 0, err)
	}
	
	// Build query parameters
	params := url.Values{}
	params.Set("response_type", "code")
	params.Set("client_id", config.ClientID)
	params.Set("redirect_uri", config.RedirectURI)
	params.Set("scope", strings.Join(config.Scopes, " "))
	
	// Add state parameter
	if config.State != "" {
		params.Set("state", config.State)
	}
	
	// Add PKCE parameters if provided
	if pkce != nil {
		params.Set("code_challenge", pkce.CodeChallenge)
		params.Set("code_challenge_method", pkce.Method)
	}
	
	authURL.RawQuery = params.Encode()
	return authURL.String(), nil
}

// ValidateRedirectURI performs comprehensive validation of redirect URIs
func (u *OAuth2Utils) ValidateRedirectURI(redirectURI string, allowInsecure bool) error {
	if redirectURI == "" {
		return NewOAuth2Error(OAuth2ErrorConfigurationError, "redirect URI cannot be empty")
	}
	
	parsedURL, err := url.Parse(redirectURI)
	if err != nil {
		return NewOAuth2ErrorWithContext(OAuth2ErrorConfigurationError, 
			"invalid redirect URI format", 0, err)
	}
	
	// Check scheme
	switch parsedURL.Scheme {
	case "http":
		if !allowInsecure {
			// HTTP is only allowed for localhost in secure mode
			if parsedURL.Hostname() != "localhost" && parsedURL.Hostname() != "127.0.0.1" {
				return NewOAuth2Error(OAuth2ErrorConfigurationError, 
					"HTTP redirect URIs are only allowed for localhost in secure mode")
			}
		}
	case "https":
		// HTTPS is always allowed
	case "":
		return NewOAuth2Error(OAuth2ErrorConfigurationError, 
			"redirect URI must include a scheme (http:// or https://)")
	default:
		// Custom schemes are allowed for mobile apps
		if len(parsedURL.Scheme) < 3 {
			return NewOAuth2Error(OAuth2ErrorConfigurationError, 
				"custom URI schemes must be at least 3 characters long")
		}
	}
	
	// Validate that URI is not just a scheme
	if parsedURL.Host == "" && parsedURL.Path == "" {
		return NewOAuth2Error(OAuth2ErrorConfigurationError, 
			"redirect URI must include a host or path")
	}
	
	return nil
}

// FindAvailablePort finds an available port for the redirect server
func (u *OAuth2Utils) FindAvailablePort(startPort int) (int, error) {
	for port := startPort; port < startPort+100; port++ {
		if u.isPortAvailable(port) {
			return port, nil
		}
	}
	
	return 0, NewOAuth2Error(OAuth2ErrorConfigurationError, 
		fmt.Sprintf("no available ports found starting from %d", startPort))
}

// isPortAvailable checks if a port is available for binding
func (u *OAuth2Utils) isPortAvailable(port int) bool {
	addr := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return false
	}
	listener.Close()
	return true
}

// GetEnvironmentFromURL determines the environment (production/sandbox) from an OAuth2 URL
func (u *OAuth2Utils) GetEnvironmentFromURL(oauthURL string) string {
	if strings.Contains(oauthURL, "api.tastyworks.com") || strings.Contains(oauthURL, "my.tastytrade.com") {
		return "production"
	} else if strings.Contains(oauthURL, "api.cert.tastyworks.com") || strings.Contains(oauthURL, "cert-my.staging-tasty.works") {
		return "sandbox"
	}
	return "unknown"
}

// CreateEnvironmentConfig creates an OAuth2Config for the specified environment
func (u *OAuth2Utils) CreateEnvironmentConfig(environment, clientID, clientSecret, redirectURI string, scopes []string) (OAuth2Config, error) {
	if len(scopes) == 0 {
		scopes = []string{"read", "trade"}
	}
	
	var config OAuth2Config
	switch strings.ToLower(environment) {
	case "production", "prod":
		config = NewProductionOAuth2Config(clientID, clientSecret, redirectURI, scopes)
	case "sandbox", "cert":
		config = NewSandboxOAuth2Config(clientID, clientSecret, redirectURI, scopes)
	default:
		return config, NewOAuth2Error(OAuth2ErrorConfigurationError, 
			fmt.Sprintf("invalid environment: %s (must be 'production' or 'sandbox')", environment))
	}
	
	return config, nil
}

// OAuth2Debugger provides debugging and logging utilities for OAuth2 troubleshooting
type OAuth2Debugger struct {
	EnableLogging    bool
	LogSensitiveData bool
	Logger           *log.Logger
}

// NewOAuth2Debugger creates a new OAuth2Debugger instance
func NewOAuth2Debugger(enableLogging bool) *OAuth2Debugger {
	return &OAuth2Debugger{
		EnableLogging:    enableLogging,
		LogSensitiveData: false,
		Logger:           log.New(os.Stdout, "[OAuth2] ", log.LstdFlags),
	}
}

// LogAuthorizationURL logs the authorization URL (with sensitive data masked)
func (d *OAuth2Debugger) LogAuthorizationURL(authURL string) {
	if !d.EnableLogging {
		return
	}
	
	if d.LogSensitiveData {
		d.Logger.Printf("Authorization URL: %s", authURL)
	} else {
		maskedURL := d.maskURLParameters(authURL, []string{"client_id", "state"})
		d.Logger.Printf("Authorization URL: %s", maskedURL)
	}
}

// LogTokenExchange logs token exchange information
func (d *OAuth2Debugger) LogTokenExchange(code string, success bool, err error) {
	if !d.EnableLogging {
		return
	}
	
	maskedCode := d.maskSensitiveData(code)
	if success {
		d.Logger.Printf("Token exchange successful with code: %s", maskedCode)
	} else {
		d.Logger.Printf("Token exchange failed with code: %s, error: %v", maskedCode, err)
	}
}

// LogTokenRefresh logs token refresh information
func (d *OAuth2Debugger) LogTokenRefresh(success bool, err error) {
	if !d.EnableLogging {
		return
	}
	
	if success {
		d.Logger.Printf("Token refresh successful")
	} else {
		d.Logger.Printf("Token refresh failed: %v", err)
	}
}

// LogRedirectServer logs redirect server events
func (d *OAuth2Debugger) LogRedirectServer(event string, details ...interface{}) {
	if !d.EnableLogging {
		return
	}
	
	if len(details) > 0 {
		d.Logger.Printf("Redirect server %s: %v", event, details)
	} else {
		d.Logger.Printf("Redirect server %s", event)
	}
}

// LogError logs OAuth2 errors with appropriate detail level
func (d *OAuth2Debugger) LogError(err error) {
	if !d.EnableLogging || err == nil {
		return
	}
	
	if detailedErr, ok := err.(*OAuth2DetailedError); ok {
		d.Logger.Printf("OAuth2 Error [%s]: %s", detailedErr.GetTypeString(), detailedErr.Error())
		if d.LogSensitiveData && detailedErr.InternalMessage != "" {
			d.Logger.Printf("Internal details: %s", detailedErr.InternalMessage)
		}
	} else {
		d.Logger.Printf("Error: %v", err)
	}
}

// maskURLParameters masks sensitive parameters in URLs
func (d *OAuth2Debugger) maskURLParameters(rawURL string, sensitiveParams []string) string {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return rawURL // Return original if parsing fails
	}
	
	query := parsedURL.Query()
	for _, param := range sensitiveParams {
		if query.Has(param) {
			query.Set(param, d.maskSensitiveData(query.Get(param)))
		}
	}
	
	parsedURL.RawQuery = query.Encode()
	return parsedURL.String()
}

// maskSensitiveData masks sensitive data for logging
func (d *OAuth2Debugger) maskSensitiveData(data string) string {
	if len(data) <= 8 {
		return strings.Repeat("*", len(data))
	}
	return data[:4] + strings.Repeat("*", len(data)-8) + data[len(data)-4:]
}

// Helper methods for OAuth2Utils

// logDebug logs debug messages if debug logging is enabled
func (u *OAuth2Utils) logDebug(format string, args ...interface{}) {
	if u.EnableDebugLogging {
		log.Printf("[OAuth2Utils] "+format, args...)
	}
}

// maskSensitiveData masks sensitive data for logging
func (u *OAuth2Utils) maskSensitiveData(data string) string {
	if len(data) <= 8 {
		return strings.Repeat("*", len(data))
	}
	return data[:4] + strings.Repeat("*", len(data)-8) + data[len(data)-4:]
}

// sendCustomSuccessResponse sends a success HTML response
func (u *OAuth2Utils) sendCustomSuccessResponse(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	
	html := `<!DOCTYPE html>
<html>
<head>
    <title>Authorization Successful</title>
    <style>
        body { font-family: Arial, sans-serif; text-align: center; margin-top: 50px; background-color: #f8f9fa; }
        .success { color: #28a745; }
        .container { max-width: 500px; margin: 0 auto; padding: 20px; background: white; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
        .icon { font-size: 48px; margin-bottom: 20px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="icon success">✓</div>
        <h1 class="success">Authorization Successful</h1>
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

// sendCustomErrorResponse sends an error HTML response
func (u *OAuth2Utils) sendCustomErrorResponse(w http.ResponseWriter, oauthErr OAuth2Error) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusBadRequest)
	
	html := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <title>Authorization Failed</title>
    <style>
        body { font-family: Arial, sans-serif; text-align: center; margin-top: 50px; background-color: #f8f9fa; }
        .error { color: #dc3545; }
        .container { max-width: 500px; margin: 0 auto; padding: 20px; background: white; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
        .error-details { background: #f8f9fa; padding: 15px; border-radius: 5px; margin-top: 20px; text-align: left; }
        .icon { font-size: 48px; margin-bottom: 20px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="icon error">✗</div>
        <h1 class="error">Authorization Failed</h1>
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

// extractBaseURL extracts the base URL from a full OAuth2 endpoint URL
func extractBaseURL(endpointURL string) string {
	parsedURL, err := url.Parse(endpointURL)
	if err != nil {
		return ""
	}
	
	return fmt.Sprintf("%s://%s", parsedURL.Scheme, parsedURL.Host)
}

// OAuth2ConfigBuilder provides a fluent interface for building OAuth2 configurations
type OAuth2ConfigBuilder struct {
	config OAuth2Config
}

// NewOAuth2ConfigBuilder creates a new OAuth2ConfigBuilder
func NewOAuth2ConfigBuilder() *OAuth2ConfigBuilder {
	return &OAuth2ConfigBuilder{
		config: OAuth2Config{
			Scopes: []string{"read", "trade"}, // Default scopes
		},
	}
}

// ClientID sets the OAuth2 client ID
func (b *OAuth2ConfigBuilder) ClientID(clientID string) *OAuth2ConfigBuilder {
	b.config.ClientID = clientID
	return b
}

// ClientSecret sets the OAuth2 client secret
func (b *OAuth2ConfigBuilder) ClientSecret(clientSecret string) *OAuth2ConfigBuilder {
	b.config.ClientSecret = clientSecret
	return b
}

// RedirectURI sets the OAuth2 redirect URI
func (b *OAuth2ConfigBuilder) RedirectURI(redirectURI string) *OAuth2ConfigBuilder {
	b.config.RedirectURI = redirectURI
	return b
}

// Scopes sets the OAuth2 scopes
func (b *OAuth2ConfigBuilder) Scopes(scopes ...string) *OAuth2ConfigBuilder {
	b.config.Scopes = scopes
	return b
}

// Production configures the builder for production environment
func (b *OAuth2ConfigBuilder) Production() *OAuth2ConfigBuilder {
	b.config.BaseURL = apiBaseURL
	b.config.AuthURL = oauth2ProductionAuthURL
	b.config.TokenURL = oauth2ProductionTokenURL
	return b
}

// Sandbox configures the builder for sandbox environment
func (b *OAuth2ConfigBuilder) Sandbox() *OAuth2ConfigBuilder {
	b.config.BaseURL = apiCertBaseURL
	b.config.AuthURL = oauth2SandboxAuthURL
	b.config.TokenURL = oauth2SandboxTokenURL
	return b
}

// Environment configures the builder for the specified environment
func (b *OAuth2ConfigBuilder) Environment(env string) *OAuth2ConfigBuilder {
	switch strings.ToLower(env) {
	case "production", "prod":
		return b.Production()
	case "sandbox", "cert":
		return b.Sandbox()
	}
	return b
}

// State sets a custom state parameter
func (b *OAuth2ConfigBuilder) State(state string) *OAuth2ConfigBuilder {
	b.config.State = state
	return b
}

// CustomEndpoints sets custom OAuth2 endpoints
func (b *OAuth2ConfigBuilder) CustomEndpoints(authURL, tokenURL string) *OAuth2ConfigBuilder {
	b.config.AuthURL = authURL
	b.config.TokenURL = tokenURL
	b.config.BaseURL = extractBaseURL(authURL)
	return b
}

// Build creates the final OAuth2Config and validates it
func (b *OAuth2ConfigBuilder) Build() (OAuth2Config, error) {
	// Generate state if not provided
	if b.config.State == "" {
		utils := NewOAuth2Utils()
		state, err := utils.GenerateSecureState()
		if err != nil {
			return b.config, err
		}
		b.config.State = state
	}
	
	// Validate the configuration
	if err := b.config.Validate(); err != nil {
		return b.config, err
	}
	
	return b.config, nil
}

// Convenience functions for common OAuth2 operations

// QuickProductionClient creates a production OAuth2 client with minimal configuration
func QuickProductionClient(clientID, clientSecret string) (*Client, error) {
	config := NewProductionOAuth2Config(clientID, clientSecret, "http://localhost:8080", []string{"read", "trade"})
	return NewOAuth2Client(config, nil)
}

// QuickSandboxClient creates a sandbox OAuth2 client with minimal configuration
func QuickSandboxClient(clientID, clientSecret string) (*Client, error) {
	config := NewSandboxOAuth2Config(clientID, clientSecret, "http://localhost:8080", []string{"read", "trade"})
	return NewCertOAuth2Client(config, nil)
}

// QuickAuthorizationFlow performs a complete OAuth2 authorization flow with built-in server
func QuickAuthorizationFlow(clientID, clientSecret, environment string) (*Client, error) {
	utils := NewOAuth2Utils()
	utils.EnableDebugLogging = true
	
	opts := utils.DefaultFlowOptions()
	opts.ClientID = clientID
	opts.ClientSecret = clientSecret
	opts.Environment = environment
	
	oauth2Client, err := utils.PerformCompleteFlow(context.Background(), opts)
	if err != nil {
		return nil, err
	}
	
	// Create a full Client instance with the OAuth2 client
	var client *Client
	if strings.ToLower(environment) == "production" {
		client, err = NewOAuth2Client(oauth2Client.GetConfig(), nil)
	} else {
		client, err = NewCertOAuth2Client(oauth2Client.GetConfig(), nil)
	}
	
	if err != nil {
		return nil, err
	}
	
	// Copy the token manager state
	client.oauth2Client.tokenManager = oauth2Client.tokenManager
	
	return client, nil
}