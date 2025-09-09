package tasty

import (
	"fmt"
	"log"
	"os"
)

// MigrationHelper provides utilities to help developers migrate from session-based to OAuth2 authentication
type MigrationHelper struct {
	// Whether to log deprecation warnings (default: true)
	LogWarnings bool
	// Custom logger for deprecation warnings (default: standard log)
	Logger *log.Logger
}

// NewMigrationHelper creates a new migration helper with default settings
func NewMigrationHelper() *MigrationHelper {
	return &MigrationHelper{
		LogWarnings: true,
		Logger:      log.New(os.Stderr, "[TASTY-GO DEPRECATION] ", log.LstdFlags),
	}
}

// LogDeprecationWarning logs a deprecation warning for session-based methods
func (m *MigrationHelper) LogDeprecationWarning(method, replacement string) {
	if m.LogWarnings && m.Logger != nil {
		m.Logger.Printf("WARNING: %s is deprecated and will be removed in a future version. Use %s instead.", method, replacement)
	}
}

// ConvertSessionConfigToOAuth2 helps convert session-based configuration to OAuth2
func ConvertSessionConfigToOAuth2(clientID, clientSecret, redirectURI string, isProduction bool) OAuth2Config {
	var config OAuth2Config
	
	if isProduction {
		config = NewProductionOAuth2Config(clientID, clientSecret, redirectURI, nil)
	} else {
		config = NewSandboxOAuth2Config(clientID, clientSecret, redirectURI, nil)
	}
	
	return config
}

// MigrationGuide provides step-by-step migration guidance
type MigrationGuide struct {
	FromSessionBased bool
	ToOAuth2         bool
	Environment      string // "production" or "sandbox"
}

// GetMigrationSteps returns a list of migration steps for the user
func (mg *MigrationGuide) GetMigrationSteps() []string {
	if !mg.FromSessionBased || !mg.ToOAuth2 {
		return []string{"No migration needed"}
	}
	
	steps := []string{
		"1. Register your application with TastyTrade to get OAuth2 credentials (client_id, client_secret)",
		"2. Set up a redirect URI for your application (e.g., http://localhost:8080/callback for development)",
		"3. Replace session-based client creation:",
	}
	
	if mg.Environment == "production" {
		steps = append(steps, []string{
			"   OLD: client := tasty.NewClient(httpClient)",
			"   NEW: config := tasty.NewProductionOAuth2Config(clientID, clientSecret, redirectURI, nil)",
			"        client, err := tasty.NewOAuth2Client(config, httpClient)",
		}...)
	} else {
		steps = append(steps, []string{
			"   OLD: client := tasty.NewCertClient(httpClient)",
			"   NEW: config := tasty.NewSandboxOAuth2Config(clientID, clientSecret, redirectURI, nil)",
			"        client, err := tasty.NewCertOAuth2Client(config, httpClient)",
		}...)
	}
	
	steps = append(steps, []string{
		"4. Replace session creation with OAuth2 flow:",
		"   OLD: session, _, err := client.CreateSession(loginInfo, nil)",
		"   NEW: authURL, err := client.GetAuthorizationURL()",
		"        // Direct user to authURL, then:",
		"        tokens, err := client.ExchangeCodeForTokens(authCode)",
		"5. Remove session validation and destruction calls (handled automatically by OAuth2)",
		"6. Test your application with the new OAuth2 flow",
		"7. Update your documentation and deployment scripts",
	}...)
	
	return steps
}

// PrintMigrationGuide prints the migration guide to stdout
func (mg *MigrationGuide) PrintMigrationGuide() {
	fmt.Println("=== TastyTrade Go SDK Migration Guide ===")
	fmt.Printf("Migrating from session-based to OAuth2 authentication (%s environment)\n\n", mg.Environment)
	
	steps := mg.GetMigrationSteps()
	for _, step := range steps {
		fmt.Println(step)
	}
	
	fmt.Println("\nFor more information, visit: https://developer.tastytrade.com/")
	fmt.Println("For OAuth2 documentation, see the README.md file")
}

// ValidateOAuth2Migration checks if OAuth2 configuration is valid for migration
func ValidateOAuth2Migration(config OAuth2Config) error {
	if err := config.Validate(); err != nil {
		return fmt.Errorf("OAuth2 configuration validation failed: %w", err)
	}
	
	// Additional migration-specific validations
	if config.ClientID == "" {
		return fmt.Errorf("OAuth2 client ID is required for migration")
	}
	
	if config.ClientSecret == "" {
		return fmt.Errorf("OAuth2 client secret is required for migration")
	}
	
	if config.RedirectURI == "" {
		return fmt.Errorf("OAuth2 redirect URI is required for migration")
	}
	
	return nil
}

// GetMigrationExamples returns code examples for common migration scenarios
func GetMigrationExamples() map[string]string {
	examples := make(map[string]string)
	
	examples["session_to_oauth2_production"] = `
// Before (Session-based - Deprecated)
client := tasty.NewClient(nil)
session, _, err := client.CreateSession(tasty.LoginInfo{
    Login:    "your_username",
    Password: "your_password",
}, nil)

// After (OAuth2)
config := tasty.NewProductionOAuth2Config(
    "your_client_id",
    "your_client_secret", 
    "http://localhost:8080/callback",
    nil,
)
client, err := tasty.NewOAuth2Client(config, nil)
authURL, err := client.GetAuthorizationURL()
// Direct user to authURL, then:
tokens, err := client.ExchangeCodeForTokens(authCode)
`
	
	examples["session_to_oauth2_sandbox"] = `
// Before (Session-based - Deprecated)
client := tasty.NewCertClient(nil)
session, _, err := client.CreateSession(tasty.LoginInfo{
    Login:    "your_username",
    Password: "your_password",
}, nil)

// After (OAuth2)
config := tasty.NewSandboxOAuth2Config(
    "your_client_id",
    "your_client_secret",
    "http://localhost:8080/callback", 
    nil,
)
client, err := tasty.NewCertOAuth2Client(config, nil)
authURL, err := client.GetAuthorizationURL()
// Direct user to authURL, then:
tokens, err := client.ExchangeCodeForTokens(authCode)
`
	
	examples["oauth2_with_redirect_server"] = `
// OAuth2 with built-in redirect server
config := tasty.NewProductionOAuth2Config(clientID, clientSecret, "http://localhost:8080/callback", nil)
client, err := tasty.NewOAuth2Client(config, nil)

// Start redirect server
server, err := client.StartRedirectServer(8080)
if err != nil {
    log.Fatal(err)
}
defer server.Shutdown()

// Get authorization URL and direct user to it
authURL, err := client.GetAuthorizationURL()
fmt.Printf("Please visit: %s\n", authURL)

// Wait for authorization code
code, err := server.WaitForCode(5 * time.Minute)
if err != nil {
    log.Fatal(err)
}

// Exchange code for tokens
tokens, err := client.ExchangeCodeForTokens(code)
`
	
	return examples
}

// CheckSessionDeprecation checks if the client is using deprecated session-based authentication
func CheckSessionDeprecation(client *Client) bool {
	return client.authMode == AuthModeSession
}

// SuggestOAuth2Alternative suggests OAuth2 alternatives for session-based operations
func SuggestOAuth2Alternative(operation string) string {
	alternatives := map[string]string{
		"CreateSession":   "Use OAuth2 authorization flow with GetAuthorizationURL() and ExchangeCodeForTokens()",
		"ValidateSession": "OAuth2 tokens are validated automatically during API requests",
		"DestroySession":  "Use ClearAuthentication() to clear OAuth2 tokens",
		"NewClient":       "Use NewOAuth2Client() with OAuth2Config",
		"NewCertClient":   "Use NewCertOAuth2Client() with OAuth2Config",
	}
	
	if alternative, exists := alternatives[operation]; exists {
		return alternative
	}
	
	return "Consider migrating to OAuth2 authentication for better security and future compatibility"
}

// Global migration helper instance
var defaultMigrationHelper = NewMigrationHelper()

// LogSessionDeprecation logs a deprecation warning for session-based methods
func LogSessionDeprecation(method, replacement string) {
	defaultMigrationHelper.LogDeprecationWarning(method, replacement)
}

// SetDeprecationLogging enables or disables deprecation warning logging
func SetDeprecationLogging(enabled bool) {
	defaultMigrationHelper.LogWarnings = enabled
}

// SetDeprecationLogger sets a custom logger for deprecation warnings
func SetDeprecationLogger(logger *log.Logger) {
	defaultMigrationHelper.Logger = logger
}