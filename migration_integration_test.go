package tasty

import (
	"bytes"
	"log"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMigrationIntegration_FullWorkflow(t *testing.T) {
	// Capture deprecation warnings
	var buf bytes.Buffer
	originalLogger := defaultMigrationHelper.Logger
	defaultMigrationHelper.Logger = log.New(&buf, "[INTEGRATION] ", 0)
	defer func() {
		defaultMigrationHelper.Logger = originalLogger
	}()

	t.Run("SessionToOAuth2Migration", func(t *testing.T) {
		buf.Reset()
		
		// Step 1: Create session-based client (deprecated)
		sessionClient := NewClient(nil)
		require.NotNil(t, sessionClient)
		require.True(t, CheckSessionDeprecation(sessionClient))
		
		// Verify deprecation warning was logged
		output := buf.String()
		require.Contains(t, output, "NewClient")
		require.Contains(t, output, "deprecated")
		
		// Step 2: Check migration status
		status := sessionClient.GetMigrationStatus()
		require.True(t, status["needs_migration"].(bool))
		require.True(t, status["is_session"].(bool))
		require.False(t, status["is_oauth2"].(bool))
		
		// Step 3: Create OAuth2 configuration
		oauth2Config := ConvertSessionConfigToOAuth2(
			"test_client_id",
			"test_client_secret", 
			"http://localhost:8080/callback",
			true, // production
		)
		
		// Step 4: Validate OAuth2 configuration
		err := ValidateOAuth2Migration(oauth2Config)
		require.NoError(t, err)
		
		// Step 5: Migrate to OAuth2
		err = sessionClient.MigrateToOAuth2(oauth2Config)
		require.NoError(t, err)
		
		// Step 6: Verify migration was successful
		require.Equal(t, AuthModeOAuth2, sessionClient.authMode)
		require.NotNil(t, sessionClient.oauth2Client)
		require.False(t, CheckSessionDeprecation(sessionClient))
		
		// Step 7: Check post-migration status
		status = sessionClient.GetMigrationStatus()
		require.False(t, status["needs_migration"].(bool))
		require.False(t, status["is_session"].(bool))
		require.True(t, status["is_oauth2"].(bool))
		require.Equal(t, "production", status["oauth2_environment"])
	})

	t.Run("MigrationGuideGeneration", func(t *testing.T) {
		// Test production migration guide
		prodGuide := &MigrationGuide{
			FromSessionBased: true,
			ToOAuth2:         true,
			Environment:      "production",
		}
		
		steps := prodGuide.GetMigrationSteps()
		require.NotEmpty(t, steps)
		require.Contains(t, steps[0], "Register your application")
		require.Contains(t, steps[2], "Replace session-based client creation")
		
		// Test sandbox migration guide
		sandboxGuide := &MigrationGuide{
			FromSessionBased: true,
			ToOAuth2:         true,
			Environment:      "sandbox",
		}
		
		sandboxSteps := sandboxGuide.GetMigrationSteps()
		require.NotEmpty(t, sandboxSteps)
		// Check that the steps contain the expected content
		stepsText := ""
		for _, step := range sandboxSteps {
			stepsText += step + " "
		}
		require.Contains(t, stepsText, "NewCertOAuth2Client")
		require.Contains(t, stepsText, "NewSandboxOAuth2Config")
	})

	t.Run("MigrationExamples", func(t *testing.T) {
		examples := GetMigrationExamples()
		
		// Verify all expected examples exist
		require.Contains(t, examples, "session_to_oauth2_production")
		require.Contains(t, examples, "session_to_oauth2_sandbox")
		require.Contains(t, examples, "oauth2_with_redirect_server")
		
		// Verify examples contain proper migration patterns
		prodExample := examples["session_to_oauth2_production"]
		require.Contains(t, prodExample, "// Before (Session-based - Deprecated)")
		require.Contains(t, prodExample, "// After (OAuth2)")
		require.Contains(t, prodExample, "CreateSession")
		require.Contains(t, prodExample, "GetAuthorizationURL")
		require.Contains(t, prodExample, "ExchangeCodeForTokens")
	})

	t.Run("OAuth2AlternativeSuggestions", func(t *testing.T) {
		// Test suggestions for common session methods
		testCases := map[string]string{
			"CreateSession":   "OAuth2 authorization flow",
			"ValidateSession": "validated automatically",
			"DestroySession":  "ClearAuthentication",
			"NewClient":       "NewOAuth2Client",
			"NewCertClient":   "NewCertOAuth2Client",
		}
		
		for method, expectedContent := range testCases {
			suggestion := SuggestOAuth2Alternative(method)
			require.Contains(t, suggestion, expectedContent, "Method: %s", method)
		}
	})

	t.Run("DeprecationWarningConfiguration", func(t *testing.T) {
		// Test disabling deprecation warnings
		SetDeprecationLogging(false)
		buf.Reset()
		
		_ = NewClient(nil) // Should not log warning
		output := buf.String()
		require.Empty(t, output)
		
		// Re-enable warnings
		SetDeprecationLogging(true)
		buf.Reset()
		
		_ = NewClient(nil) // Should log warning
		output = buf.String()
		require.Contains(t, output, "deprecated")
	})
}

func TestMigrationIntegration_ErrorHandling(t *testing.T) {
	t.Run("InvalidOAuth2Configuration", func(t *testing.T) {
		client := NewClient(nil)
		
		// Test migration with invalid config
		invalidConfig := OAuth2Config{
			ClientID: "test",
			// Missing required fields
		}
		
		err := client.MigrateToOAuth2(invalidConfig)
		require.Error(t, err)
		require.Contains(t, err.Error(), "validation failed")
		
		// Verify client is still in session mode
		require.Equal(t, AuthModeSession, client.authMode)
	})

	t.Run("FallbackMechanism", func(t *testing.T) {
		client := NewClient(nil)
		
		// Test fallback with invalid OAuth2 and invalid session
		invalidOAuth2Config := &OAuth2Config{
			ClientID: "invalid",
		}
		invalidSessionLogin := &LoginInfo{
			Login:    "",
			Password: "",
		}
		
		err := client.TryOAuth2FallbackToSession(invalidOAuth2Config, invalidSessionLogin, nil)
		require.Error(t, err)
		// The error could be either authentication failed or invalid credentials
		require.True(t, err.Error() != "", "Expected an error message")
	})

	t.Run("ConfigurationValidation", func(t *testing.T) {
		testCases := []struct {
			name   string
			config OAuth2Config
			hasErr bool
		}{
			{
				name: "valid_config",
				config: OAuth2Config{
					ClientID:     "test_client",
					ClientSecret: "test_secret",
					RedirectURI:  "http://localhost:8080/callback",
				},
				hasErr: false,
			},
			{
				name: "missing_client_id",
				config: OAuth2Config{
					ClientSecret: "test_secret",
					RedirectURI:  "http://localhost:8080/callback",
				},
				hasErr: true,
			},
			{
				name: "missing_client_secret",
				config: OAuth2Config{
					ClientID:    "test_client",
					RedirectURI: "http://localhost:8080/callback",
				},
				hasErr: true,
			},
			{
				name: "missing_redirect_uri",
				config: OAuth2Config{
					ClientID:     "test_client",
					ClientSecret: "test_secret",
				},
				hasErr: true,
			},
		}
		
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				err := ValidateOAuth2Migration(tc.config)
				if tc.hasErr {
					require.Error(t, err)
				} else {
					require.NoError(t, err)
				}
			})
		}
	})
}

func TestMigrationIntegration_BackwardCompatibility(t *testing.T) {
	t.Run("ExistingSessionCodeStillWorks", func(t *testing.T) {
		// Simulate existing session-based code
		client := NewClient(nil)
		
		// Set up a mock session token (as existing code might do)
		testToken := "mock-session-token"
		client.Session.SessionToken = &testToken
		
		// Verify session-based functionality still works
		require.True(t, client.IsAuthenticated())
		require.True(t, client.IsSessionMode())
		require.False(t, client.IsOAuth2Mode())
		
		// Clear authentication should work
		client.ClearAuthentication()
		require.False(t, client.IsAuthenticated())
	})

	t.Run("OAuth2CodeWorks", func(t *testing.T) {
		// Test that new OAuth2 code works
		config := NewProductionOAuth2Config(
			"test_client_id",
			"test_client_secret",
			"http://localhost:8080/callback",
			nil,
		)
		
		client, err := NewOAuth2Client(config, nil)
		require.NoError(t, err)
		require.NotNil(t, client)
		
		// Verify OAuth2 functionality
		require.True(t, client.IsOAuth2Mode())
		require.False(t, client.IsSessionMode())
		
		// OAuth2 methods should be available
		authURL, err := client.GetAuthorizationURL()
		require.NoError(t, err)
		require.NotEmpty(t, authURL)
		
		// Clear authentication should work
		client.ClearAuthentication()
		require.False(t, client.IsAuthenticated())
	})

	t.Run("BothModesCanCoexist", func(t *testing.T) {
		// Create both types of clients
		sessionClient := NewClient(nil)
		
		config := NewProductionOAuth2Config(
			"test_client_id",
			"test_client_secret",
			"http://localhost:8080/callback",
			nil,
		)
		oauth2Client, err := NewOAuth2Client(config, nil)
		require.NoError(t, err)
		
		// Verify they have different auth modes
		require.Equal(t, AuthModeSession, sessionClient.authMode)
		require.Equal(t, AuthModeOAuth2, oauth2Client.authMode)
		
		// Verify they work independently
		require.True(t, sessionClient.IsSessionMode())
		require.False(t, sessionClient.IsOAuth2Mode())
		
		require.False(t, oauth2Client.IsSessionMode())
		require.True(t, oauth2Client.IsOAuth2Mode())
	})
}