package tasty

import (
	"bytes"
	"log"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMigrationHelper(t *testing.T) {
	t.Run("NewMigrationHelper", func(t *testing.T) {
		helper := NewMigrationHelper()
		require.NotNil(t, helper)
		require.True(t, helper.LogWarnings)
		require.NotNil(t, helper.Logger)
	})

	t.Run("LogDeprecationWarning", func(t *testing.T) {
		var buf bytes.Buffer
		helper := &MigrationHelper{
			LogWarnings: true,
			Logger:      log.New(&buf, "[TEST] ", 0),
		}

		helper.LogDeprecationWarning("TestMethod", "NewTestMethod")
		
		output := buf.String()
		require.Contains(t, output, "WARNING")
		require.Contains(t, output, "TestMethod")
		require.Contains(t, output, "deprecated")
		require.Contains(t, output, "NewTestMethod")
	})

	t.Run("LogDeprecationWarning_Disabled", func(t *testing.T) {
		var buf bytes.Buffer
		helper := &MigrationHelper{
			LogWarnings: false,
			Logger:      log.New(&buf, "[TEST] ", 0),
		}

		helper.LogDeprecationWarning("TestMethod", "NewTestMethod")
		
		output := buf.String()
		require.Empty(t, output)
	})
}

func TestConvertSessionConfigToOAuth2(t *testing.T) {
	t.Run("Production", func(t *testing.T) {
		config := ConvertSessionConfigToOAuth2("client123", "secret456", "http://localhost:8080/callback", true)
		
		require.Equal(t, "client123", config.ClientID)
		require.Equal(t, "secret456", config.ClientSecret)
		require.Equal(t, "http://localhost:8080/callback", config.RedirectURI)
		require.True(t, config.IsProduction())
		require.False(t, config.IsSandbox())
	})

	t.Run("Sandbox", func(t *testing.T) {
		config := ConvertSessionConfigToOAuth2("client123", "secret456", "http://localhost:8080/callback", false)
		
		require.Equal(t, "client123", config.ClientID)
		require.Equal(t, "secret456", config.ClientSecret)
		require.Equal(t, "http://localhost:8080/callback", config.RedirectURI)
		require.False(t, config.IsProduction())
		require.True(t, config.IsSandbox())
	})
}

func TestMigrationGuide(t *testing.T) {
	t.Run("GetMigrationSteps_Production", func(t *testing.T) {
		guide := &MigrationGuide{
			FromSessionBased: true,
			ToOAuth2:         true,
			Environment:      "production",
		}

		steps := guide.GetMigrationSteps()
		require.NotEmpty(t, steps)
		require.Contains(t, strings.Join(steps, " "), "NewOAuth2Client")
		require.Contains(t, strings.Join(steps, " "), "NewProductionOAuth2Config")
	})

	t.Run("GetMigrationSteps_Sandbox", func(t *testing.T) {
		guide := &MigrationGuide{
			FromSessionBased: true,
			ToOAuth2:         true,
			Environment:      "sandbox",
		}

		steps := guide.GetMigrationSteps()
		require.NotEmpty(t, steps)
		require.Contains(t, strings.Join(steps, " "), "NewCertOAuth2Client")
		require.Contains(t, strings.Join(steps, " "), "NewSandboxOAuth2Config")
	})

	t.Run("GetMigrationSteps_NoMigration", func(t *testing.T) {
		guide := &MigrationGuide{
			FromSessionBased: false,
			ToOAuth2:         false,
		}

		steps := guide.GetMigrationSteps()
		require.Equal(t, []string{"No migration needed"}, steps)
	})
}

func TestValidateOAuth2Migration(t *testing.T) {
	t.Run("ValidConfig", func(t *testing.T) {
		config := NewProductionOAuth2Config("client123", "secret456", "http://localhost:8080/callback", nil)
		err := ValidateOAuth2Migration(config)
		require.NoError(t, err)
	})

	t.Run("MissingClientID", func(t *testing.T) {
		config := NewProductionOAuth2Config("", "secret456", "http://localhost:8080/callback", nil)
		err := ValidateOAuth2Migration(config)
		require.Error(t, err)
		require.Contains(t, err.Error(), "client ID")
	})

	t.Run("MissingClientSecret", func(t *testing.T) {
		config := OAuth2Config{
			ClientID:    "client123",
			RedirectURI: "http://localhost:8080/callback",
		}
		err := ValidateOAuth2Migration(config)
		require.Error(t, err)
		require.Contains(t, err.Error(), "client secret")
	})

	t.Run("MissingRedirectURI", func(t *testing.T) {
		config := OAuth2Config{
			ClientID:     "client123",
			ClientSecret: "secret456",
		}
		err := ValidateOAuth2Migration(config)
		require.Error(t, err)
		require.Contains(t, err.Error(), "redirect URI")
	})
}

func TestGetMigrationExamples(t *testing.T) {
	examples := GetMigrationExamples()
	
	require.NotEmpty(t, examples)
	require.Contains(t, examples, "session_to_oauth2_production")
	require.Contains(t, examples, "session_to_oauth2_sandbox")
	require.Contains(t, examples, "oauth2_with_redirect_server")
	
	// Check that examples contain expected content
	prodExample := examples["session_to_oauth2_production"]
	require.Contains(t, prodExample, "NewClient")
	require.Contains(t, prodExample, "NewOAuth2Client")
	require.Contains(t, prodExample, "NewProductionOAuth2Config")
	
	sandboxExample := examples["session_to_oauth2_sandbox"]
	require.Contains(t, sandboxExample, "NewCertClient")
	require.Contains(t, sandboxExample, "NewCertOAuth2Client")
	require.Contains(t, sandboxExample, "NewSandboxOAuth2Config")
}

func TestCheckSessionDeprecation(t *testing.T) {
	t.Run("SessionClient", func(t *testing.T) {
		client := NewClient(nil)
		require.True(t, CheckSessionDeprecation(client))
	})

	t.Run("OAuth2Client", func(t *testing.T) {
		config := NewProductionOAuth2Config("client123", "secret456", "http://localhost:8080/callback", nil)
		client, err := NewOAuth2Client(config, nil)
		require.NoError(t, err)
		require.False(t, CheckSessionDeprecation(client))
	})
}

func TestSuggestOAuth2Alternative(t *testing.T) {
	testCases := map[string]string{
		"CreateSession":   "OAuth2 authorization flow",
		"ValidateSession": "validated automatically",
		"DestroySession":  "ClearAuthentication",
		"NewClient":       "NewOAuth2Client",
		"NewCertClient":   "NewCertOAuth2Client",
		"UnknownMethod":   "Consider migrating to OAuth2",
	}

	for method, expectedContent := range testCases {
		suggestion := SuggestOAuth2Alternative(method)
		require.Contains(t, suggestion, expectedContent, "Method: %s", method)
	}
}

func TestDeprecationLogging(t *testing.T) {
	t.Run("GlobalLogging", func(t *testing.T) {
		// Test global deprecation logging functions
		SetDeprecationLogging(false)
		require.False(t, defaultMigrationHelper.LogWarnings)
		
		SetDeprecationLogging(true)
		require.True(t, defaultMigrationHelper.LogWarnings)
		
		// Test custom logger
		var buf bytes.Buffer
		customLogger := log.New(&buf, "[CUSTOM] ", 0)
		SetDeprecationLogger(customLogger)
		require.Equal(t, customLogger, defaultMigrationHelper.Logger)
		
		// Test logging
		LogSessionDeprecation("TestMethod", "NewMethod")
		output := buf.String()
		require.Contains(t, output, "TestMethod")
		require.Contains(t, output, "deprecated")
	})
}

func TestBackwardCompatibility_SessionMethods(t *testing.T) {
	// Capture deprecation warnings
	var buf bytes.Buffer
	originalLogger := defaultMigrationHelper.Logger
	defaultMigrationHelper.Logger = log.New(&buf, "[TEST] ", 0)
	defer func() {
		defaultMigrationHelper.Logger = originalLogger
	}()

	t.Run("NewClient_DeprecationWarning", func(t *testing.T) {
		buf.Reset()
		client := NewClient(nil)
		
		require.NotNil(t, client)
		require.Equal(t, AuthModeSession, client.authMode)
		
		output := buf.String()
		require.Contains(t, output, "NewClient")
		require.Contains(t, output, "deprecated")
	})

	t.Run("NewCertClient_DeprecationWarning", func(t *testing.T) {
		buf.Reset()
		client := NewCertClient(nil)
		
		require.NotNil(t, client)
		require.Equal(t, AuthModeSession, client.authMode)
		
		output := buf.String()
		require.Contains(t, output, "NewCertClient")
		require.Contains(t, output, "deprecated")
	})
}

func TestClient_MigrationMethods(t *testing.T) {
	t.Run("TryOAuth2FallbackToSession", func(t *testing.T) {
		client := NewClient(nil)
		
		// Test with invalid OAuth2 config and valid session login
		oauth2Config := &OAuth2Config{
			ClientID: "invalid",
			// Missing required fields
		}
		sessionLogin := &LoginInfo{
			Login:    "test",
			Password: "test",
		}
		
		// This should fail because we don't have a real server, but it tests the fallback logic
		err := client.TryOAuth2FallbackToSession(oauth2Config, sessionLogin, nil)
		require.Error(t, err) // Expected to fail without real server
	})

	t.Run("MigrateToOAuth2", func(t *testing.T) {
		client := NewClient(nil)
		require.Equal(t, AuthModeSession, client.authMode)
		
		config := NewProductionOAuth2Config("client123", "secret456", "http://localhost:8080/callback", nil)
		err := client.MigrateToOAuth2(config)
		require.NoError(t, err)
		
		require.Equal(t, AuthModeOAuth2, client.authMode)
		require.NotNil(t, client.oauth2Client)
	})

	t.Run("MigrateToOAuth2_InvalidConfig", func(t *testing.T) {
		client := NewClient(nil)
		
		config := OAuth2Config{
			ClientID: "client123",
			// Missing required fields
		}
		err := client.MigrateToOAuth2(config)
		require.Error(t, err)
		require.Contains(t, err.Error(), "validation failed")
	})

	t.Run("GetMigrationStatus", func(t *testing.T) {
		// Test session client
		sessionClient := NewClient(nil)
		status := sessionClient.GetMigrationStatus()
		
		require.Equal(t, "session", status["auth_mode"])
		require.True(t, status["is_session"].(bool))
		require.False(t, status["is_oauth2"].(bool))
		require.True(t, status["needs_migration"].(bool))
		
		// Test OAuth2 client
		config := NewProductionOAuth2Config("client123", "secret456", "http://localhost:8080/callback", nil)
		oauth2Client, err := NewOAuth2Client(config, nil)
		require.NoError(t, err)
		
		status = oauth2Client.GetMigrationStatus()
		require.Equal(t, "oauth2", status["auth_mode"])
		require.False(t, status["is_session"].(bool))
		require.True(t, status["is_oauth2"].(bool))
		require.False(t, status["needs_migration"].(bool))
		require.Equal(t, "production", status["oauth2_environment"])
	})
}

func TestBackwardCompatibility_ExistingCode(t *testing.T) {
	t.Run("SessionBasedWorkflow", func(t *testing.T) {
		// Test that existing session-based code still works
		client := NewClient(nil)
		require.NotNil(t, client)
		require.Equal(t, AuthModeSession, client.authMode)
		
		// Test that session token can be set manually (for testing)
		testToken := "test-session-token"
		client.Session.SessionToken = &testToken
		require.True(t, client.IsAuthenticated())
		require.True(t, client.IsSessionMode())
		require.False(t, client.IsOAuth2Mode())
		
		// Test clearing authentication
		client.ClearAuthentication()
		require.False(t, client.IsAuthenticated())
	})

	t.Run("OAuth2Workflow", func(t *testing.T) {
		// Test that OAuth2 workflow works
		config := NewProductionOAuth2Config("client123", "secret456", "http://localhost:8080/callback", nil)
		client, err := NewOAuth2Client(config, nil)
		require.NoError(t, err)
		require.NotNil(t, client)
		require.Equal(t, AuthModeOAuth2, client.authMode)
		require.True(t, client.IsOAuth2Mode())
		require.False(t, client.IsSessionMode())
		
		// Test OAuth2 methods are available
		authURL, err := client.GetAuthorizationURL()
		require.NoError(t, err)
		require.NotEmpty(t, authURL)
		
		// Test clearing authentication
		client.ClearAuthentication()
		require.False(t, client.IsAuthenticated())
	})
}

func TestBackwardCompatibility_RequestMethods(t *testing.T) {
	t.Run("SessionRequest_NoToken", func(t *testing.T) {
		client := NewClient(nil)
		
		// Test that session request fails without token
		_, err := client.request("GET", "/test", nil, nil, nil)
		require.Error(t, err)
		require.Contains(t, err.Error(), "Session is invalid")
	})

	t.Run("OAuth2Request_NoClient", func(t *testing.T) {
		// Create client in OAuth2 mode but without proper OAuth2 client
		client := &Client{
			authMode: AuthModeOAuth2,
		}
		
		// Test that OAuth2 request fails without OAuth2 client
		_, err := client.request("GET", "/test", nil, nil, nil)
		require.Error(t, err)
		require.Contains(t, err.Error(), "OAuth2 client not initialized")
	})
}