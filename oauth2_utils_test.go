package tasty

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOAuth2Utils_NewOAuth2Utils(t *testing.T) {
	utils := NewOAuth2Utils()
	
	assert.NotNil(t, utils)
	assert.False(t, utils.EnableDebugLogging)
	assert.Equal(t, 30*time.Second, utils.DefaultTimeout)
	assert.Equal(t, 8080, utils.DefaultPort)
}

func TestOAuth2Utils_DefaultFlowOptions(t *testing.T) {
	utils := NewOAuth2Utils()
	opts := utils.DefaultFlowOptions()
	
	assert.Equal(t, "production", opts.Environment)
	assert.True(t, opts.UseBuiltinServer)
	assert.Equal(t, 8080, opts.ServerPort)
	assert.Equal(t, 30*time.Second, opts.Timeout)
	assert.True(t, opts.UsePKCE)
	assert.Equal(t, []string{"read", "trade"}, opts.Scopes)
}

func TestOAuth2Utils_CreateConfigFromOptions(t *testing.T) {
	utils := NewOAuth2Utils()
	
	tests := []struct {
		name    string
		opts    OAuth2FlowOptions
		wantErr bool
		checks  func(t *testing.T, config OAuth2Config)
	}{
		{
			name: "production environment",
			opts: OAuth2FlowOptions{
				ClientID:     "test_client_id",
				ClientSecret: "test_client_secret",
				RedirectURI:  "http://localhost:8080",
				Environment:  "production",
				Scopes:       []string{"read", "trade"},
			},
			wantErr: false,
			checks: func(t *testing.T, config OAuth2Config) {
				assert.Equal(t, "test_client_id", config.ClientID)
				assert.Equal(t, "test_client_secret", config.ClientSecret)
				assert.Equal(t, "http://localhost:8080", config.RedirectURI)
				assert.Equal(t, apiBaseURL, config.BaseURL)
				assert.Equal(t, oauth2ProductionAuthURL, config.AuthURL)
				assert.Equal(t, oauth2ProductionTokenURL, config.TokenURL)
				assert.Equal(t, []string{"read", "trade"}, config.Scopes)
			},
		},
		{
			name: "sandbox environment",
			opts: OAuth2FlowOptions{
				ClientID:     "test_client_id",
				ClientSecret: "test_client_secret",
				RedirectURI:  "http://localhost:8080",
				Environment:  "sandbox",
				Scopes:       []string{"read"},
			},
			wantErr: false,
			checks: func(t *testing.T, config OAuth2Config) {
				assert.Equal(t, apiCertBaseURL, config.BaseURL)
				assert.Equal(t, oauth2SandboxAuthURL, config.AuthURL)
				assert.Equal(t, oauth2SandboxTokenURL, config.TokenURL)
				assert.Equal(t, []string{"read"}, config.Scopes)
			},
		},
		{
			name: "custom endpoints",
			opts: OAuth2FlowOptions{
				ClientID:       "test_client_id",
				ClientSecret:   "test_client_secret",
				RedirectURI:    "http://localhost:8080",
				CustomAuthURL:  "https://example.com/oauth/authorize",
				CustomTokenURL: "https://example.com/oauth/token",
				Scopes:         []string{"custom"},
			},
			wantErr: false,
			checks: func(t *testing.T, config OAuth2Config) {
				assert.Equal(t, "https://example.com/oauth/authorize", config.AuthURL)
				assert.Equal(t, "https://example.com/oauth/token", config.TokenURL)
				assert.Equal(t, "https://example.com", config.BaseURL)
				assert.Equal(t, []string{"custom"}, config.Scopes)
			},
		},
		{
			name: "invalid environment",
			opts: OAuth2FlowOptions{
				ClientID:     "test_client_id",
				ClientSecret: "test_client_secret",
				RedirectURI:  "http://localhost:8080",
				Environment:  "invalid",
			},
			wantErr: true,
		},
		{
			name: "missing client ID",
			opts: OAuth2FlowOptions{
				ClientSecret: "test_client_secret",
				RedirectURI:  "http://localhost:8080",
				Environment:  "production",
			},
			wantErr: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, err := utils.CreateConfigFromOptions(tt.opts)
			
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			
			require.NoError(t, err)
			if tt.checks != nil {
				tt.checks(t, config)
			}
		})
	}
}

func TestOAuth2Utils_GenerateSecureState(t *testing.T) {
	utils := NewOAuth2Utils()
	
	// Generate multiple states to ensure uniqueness
	states := make(map[string]bool)
	for i := 0; i < 100; i++ {
		state, err := utils.GenerateSecureState()
		require.NoError(t, err)
		assert.NotEmpty(t, state)
		assert.False(t, states[state], "Generated duplicate state: %s", state)
		states[state] = true
		
		// Verify it's base64url encoded
		assert.NotContains(t, state, "+")
		assert.NotContains(t, state, "/")
		assert.NotContains(t, state, "=")
	}
}

func TestOAuth2Utils_BuildAuthorizationURL(t *testing.T) {
	utils := NewOAuth2Utils()
	
	config := OAuth2Config{
		ClientID:    "test_client_id",
		RedirectURI: "http://localhost:8080/callback",
		Scopes:      []string{"read"},
		State:       "test_state",
		AuthURL:     "https://example.com/oauth/authorize",
	}
	
	authURL, err := utils.BuildAuthorizationURL(config, nil)
	require.NoError(t, err)
	
	// Parse and verify the URL
	parsedURL, err := url.Parse(authURL)
	require.NoError(t, err)
	
	query := parsedURL.Query()
	assert.Equal(t, "code", query.Get("response_type"))
	assert.Equal(t, "test_client_id", query.Get("client_id"))
	assert.Equal(t, "http://localhost:8080/callback", query.Get("redirect_uri"))
	assert.Equal(t, "read", query.Get("scope"))
	assert.Equal(t, "test_state", query.Get("state"))
	assert.Empty(t, query.Get("code_challenge"))
	assert.Empty(t, query.Get("code_challenge_method"))
}

func TestOAuth2Utils_ValidateRedirectURI(t *testing.T) {
	utils := NewOAuth2Utils()
	
	tests := []struct {
		name          string
		redirectURI   string
		allowInsecure bool
		wantErr       bool
		errorContains string
	}{
		{
			name:        "valid HTTPS URI",
			redirectURI: "https://example.com/callback",
			wantErr:     false,
		},
		{
			name:        "valid localhost HTTP URI",
			redirectURI: "http://localhost:8080/callback",
			wantErr:     false,
		},
		{
			name:        "valid 127.0.0.1 HTTP URI",
			redirectURI: "http://127.0.0.1:8080/callback",
			wantErr:     false,
		},
		{
			name:          "HTTP URI with allowInsecure",
			redirectURI:   "http://example.com/callback",
			allowInsecure: true,
			wantErr:       false,
		},
		{
			name:          "HTTP URI without allowInsecure",
			redirectURI:   "http://example.com/callback",
			allowInsecure: false,
			wantErr:       true,
			errorContains: "HTTP redirect URIs are only allowed for localhost",
		},
		{
			name:        "valid custom scheme",
			redirectURI: "myapp://callback",
			wantErr:     false,
		},
		{
			name:          "invalid custom scheme (too short)",
			redirectURI:   "ab://callback",
			wantErr:       true,
			errorContains: "custom URI schemes must be at least 3 characters long",
		},
		{
			name:          "empty URI",
			redirectURI:   "",
			wantErr:       true,
			errorContains: "redirect URI cannot be empty",
		},
		{
			name:          "URI without scheme",
			redirectURI:   "example.com/callback",
			wantErr:       true,
			errorContains: "redirect URI must include a scheme",
		},
		{
			name:          "URI with only scheme",
			redirectURI:   "https://",
			wantErr:       true,
			errorContains: "redirect URI must include a host or path",
		},
		{
			name:          "invalid URI format",
			redirectURI:   "ht tp://example.com",
			wantErr:       true,
			errorContains: "invalid redirect URI format",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := utils.ValidateRedirectURI(tt.redirectURI, tt.allowInsecure)
			
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestOAuth2Utils_ParseRedirectURL(t *testing.T) {
	utils := NewOAuth2Utils()
	expectedState := "test_state_123"
	
	tests := []struct {
		name        string
		redirectURL string
		wantCode    string
		wantErr     bool
		errorType   string
	}{
		{
			name:        "valid redirect with code",
			redirectURL: "http://localhost:8080/callback?code=test_code_123&state=test_state_123",
			wantCode:    "test_code_123",
			wantErr:     false,
		},
		{
			name:        "redirect with error",
			redirectURL: "http://localhost:8080/callback?error=access_denied&error_description=User%20denied&state=test_state_123",
			wantErr:     true,
			errorType:   OAuth2ErrorAccessDenied,
		},
		{
			name:        "missing code",
			redirectURL: "http://localhost:8080/callback?state=test_state_123",
			wantErr:     true,
			errorType:   OAuth2ErrorMissingCode,
		},
		{
			name:        "invalid state",
			redirectURL: "http://localhost:8080/callback?code=test_code_123&state=wrong_state",
			wantErr:     true,
			errorType:   OAuth2ErrorInvalidState,
		},
		{
			name:        "missing state",
			redirectURL: "http://localhost:8080/callback?code=test_code_123",
			wantErr:     true,
			errorType:   OAuth2ErrorInvalidState,
		},
		{
			name:        "invalid URL format",
			redirectURL: "ht tp://invalid url",
			wantErr:     true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, err := utils.ParseRedirectURL(tt.redirectURL, expectedState)
			
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errorType != "" {
					if detailedErr, ok := err.(*OAuth2DetailedError); ok {
						assert.Equal(t, tt.errorType, detailedErr.ErrorCode)
					}
				}
				assert.Empty(t, code)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantCode, code)
			}
		})
	}
}

func TestOAuth2Utils_GetEnvironmentFromURL(t *testing.T) {
	utils := NewOAuth2Utils()
	
	tests := []struct {
		name        string
		url         string
		wantEnv     string
	}{
		{
			name:    "production auth URL",
			url:     "https://my.tastytrade.com/auth.html",
			wantEnv: "production",
		},
		{
			name:    "production token URL",
			url:     "https://api.tastyworks.com/oauth/token",
			wantEnv: "production",
		},
		{
			name:    "sandbox auth URL",
			url:     "https://cert-my.staging-tasty.works/auth.html",
			wantEnv: "sandbox",
		},
		{
			name:    "sandbox token URL",
			url:     "https://api.cert.tastyworks.com/oauth/token",
			wantEnv: "sandbox",
		},
		{
			name:    "unknown URL",
			url:     "https://example.com/oauth/authorize",
			wantEnv: "unknown",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := utils.GetEnvironmentFromURL(tt.url)
			assert.Equal(t, tt.wantEnv, env)
		})
	}
}

func TestOAuth2Utils_CreateEnvironmentConfig(t *testing.T) {
	utils := NewOAuth2Utils()
	
	tests := []struct {
		name        string
		environment string
		wantErr     bool
		checkEnv    func(t *testing.T, config OAuth2Config)
	}{
		{
			name:        "production environment",
			environment: "production",
			wantErr:     false,
			checkEnv: func(t *testing.T, config OAuth2Config) {
				assert.True(t, config.IsProduction())
				assert.False(t, config.IsSandbox())
			},
		},
		{
			name:        "prod environment",
			environment: "prod",
			wantErr:     false,
			checkEnv: func(t *testing.T, config OAuth2Config) {
				assert.True(t, config.IsProduction())
			},
		},
		{
			name:        "sandbox environment",
			environment: "sandbox",
			wantErr:     false,
			checkEnv: func(t *testing.T, config OAuth2Config) {
				assert.True(t, config.IsSandbox())
				assert.False(t, config.IsProduction())
			},
		},
		{
			name:        "cert environment",
			environment: "cert",
			wantErr:     false,
			checkEnv: func(t *testing.T, config OAuth2Config) {
				assert.True(t, config.IsSandbox())
			},
		},
		{
			name:        "invalid environment",
			environment: "invalid",
			wantErr:     true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, err := utils.CreateEnvironmentConfig(
				tt.environment,
				"test_client_id",
				"test_client_secret",
				"http://localhost:8080",
				[]string{"read", "trade"},
			)
			
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			
			require.NoError(t, err)
			assert.Equal(t, "test_client_id", config.ClientID)
			assert.Equal(t, "test_client_secret", config.ClientSecret)
			assert.Equal(t, "http://localhost:8080", config.RedirectURI)
			assert.Equal(t, []string{"read", "trade"}, config.Scopes)
			
			if tt.checkEnv != nil {
				tt.checkEnv(t, config)
			}
		})
	}
}

func TestOAuth2Utils_CreateCustomRedirectHandler(t *testing.T) {
	utils := NewOAuth2Utils()
	expectedState := "test_state_123"
	codeChan := make(chan string, 1)
	errChan := make(chan error, 1)
	
	handler := utils.CreateCustomRedirectHandler(expectedState, codeChan, errChan)
	
	tests := []struct {
		name           string
		requestURL     string
		expectCode     string
		expectError    bool
		expectedErrType string
	}{
		{
			name:       "successful authorization",
			requestURL: "/callback?code=test_code_123&state=test_state_123",
			expectCode: "test_code_123",
		},
		{
			name:            "authorization error",
			requestURL:      "/callback?error=access_denied&error_description=User%20denied&state=test_state_123",
			expectError:     true,
			expectedErrType: OAuth2ErrorAccessDenied,
		},
		{
			name:            "missing code",
			requestURL:      "/callback?state=test_state_123",
			expectError:     true,
			expectedErrType: OAuth2ErrorMissingCode,
		},
		{
			name:            "invalid state",
			requestURL:      "/callback?code=test_code_123&state=wrong_state",
			expectError:     true,
			expectedErrType: OAuth2ErrorInvalidState,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear channels
			select {
			case <-codeChan:
			default:
			}
			select {
			case <-errChan:
			default:
			}
			
			// Create test request
			req := httptest.NewRequest("GET", tt.requestURL, nil)
			w := httptest.NewRecorder()
			
			// Call handler
			handler(w, req)
			
			if tt.expectError {
				select {
				case err := <-errChan:
					assert.Error(t, err)
					if tt.expectedErrType != "" {
						if detailedErr, ok := err.(*OAuth2DetailedError); ok {
							assert.Equal(t, tt.expectedErrType, detailedErr.ErrorCode)
						}
					}
				case <-time.After(100 * time.Millisecond):
					t.Fatal("Expected error but none received")
				}
				
				// Should not receive code
				select {
				case code := <-codeChan:
					t.Fatalf("Unexpected code received: %s", code)
				default:
					// Expected
				}
			} else {
				select {
				case code := <-codeChan:
					assert.Equal(t, tt.expectCode, code)
				case <-time.After(100 * time.Millisecond):
					t.Fatal("Expected code but none received")
				}
				
				// Should not receive error
				select {
				case err := <-errChan:
					t.Fatalf("Unexpected error received: %v", err)
				default:
					// Expected
				}
			}
			
			// Check response
			if tt.expectError {
				assert.Equal(t, http.StatusBadRequest, w.Code)
			} else {
				assert.Equal(t, http.StatusOK, w.Code)
			}
			assert.Contains(t, w.Header().Get("Content-Type"), "text/html")
		})
	}
}

func TestOAuth2Debugger(t *testing.T) {
	debugger := NewOAuth2Debugger(true)
	
	assert.NotNil(t, debugger)
	assert.True(t, debugger.EnableLogging)
	assert.False(t, debugger.LogSensitiveData)
	assert.NotNil(t, debugger.Logger)
	
	// Test masking sensitive data
	testData := "1234567890abcdef"
	masked := debugger.maskSensitiveData(testData)
	assert.Equal(t, "1234********cdef", masked)
	
	// Test short data masking
	shortData := "123"
	maskedShort := debugger.maskSensitiveData(shortData)
	assert.Equal(t, "***", maskedShort)
}

func TestOAuth2Debugger_LogMethods(t *testing.T) {
	// Test with logging disabled
	debugger := NewOAuth2Debugger(false)
	
	// These should not panic and should not log anything
	debugger.LogAuthorizationURL("https://example.com/oauth/authorize?client_id=test")
	debugger.LogTokenExchange("test_code", true, nil)
	debugger.LogTokenRefresh(false, fmt.Errorf("test error"))
	debugger.LogRedirectServer("started", 8080)
	debugger.LogError(NewOAuth2Error(OAuth2ErrorInvalidState, "test error"))
	
	// Test with logging enabled
	debugger = NewOAuth2Debugger(true)
	
	// These should not panic
	debugger.LogAuthorizationURL("https://example.com/oauth/authorize?client_id=test&state=secret")
	debugger.LogTokenExchange("test_code_123", true, nil)
	debugger.LogTokenRefresh(false, fmt.Errorf("refresh failed"))
	debugger.LogRedirectServer("started", 8080)
	debugger.LogError(NewOAuth2Error(OAuth2ErrorTokenExpired, "token expired"))
}

func TestOAuth2ConfigBuilder(t *testing.T) {
	builder := NewOAuth2ConfigBuilder()
	assert.NotNil(t, builder)
	assert.Equal(t, []string{"read", "trade"}, builder.config.Scopes)
	
	// Test fluent interface
	config, err := builder.
		ClientID("test_client_id").
		ClientSecret("test_client_secret").
		RedirectURI("http://localhost:8080").
		Scopes("read", "trade", "admin").
		Production().
		State("custom_state").
		Build()
	
	require.NoError(t, err)
	assert.Equal(t, "test_client_id", config.ClientID)
	assert.Equal(t, "test_client_secret", config.ClientSecret)
	assert.Equal(t, "http://localhost:8080", config.RedirectURI)
	assert.Equal(t, []string{"read", "trade", "admin"}, config.Scopes)
	assert.Equal(t, "custom_state", config.State)
	assert.True(t, config.IsProduction())
}

func TestOAuth2ConfigBuilder_Environment(t *testing.T) {
	tests := []struct {
		name        string
		environment string
		checkFunc   func(t *testing.T, config OAuth2Config)
	}{
		{
			name:        "production",
			environment: "production",
			checkFunc: func(t *testing.T, config OAuth2Config) {
				assert.True(t, config.IsProduction())
			},
		},
		{
			name:        "prod",
			environment: "prod",
			checkFunc: func(t *testing.T, config OAuth2Config) {
				assert.True(t, config.IsProduction())
			},
		},
		{
			name:        "sandbox",
			environment: "sandbox",
			checkFunc: func(t *testing.T, config OAuth2Config) {
				assert.True(t, config.IsSandbox())
			},
		},
		{
			name:        "cert",
			environment: "cert",
			checkFunc: func(t *testing.T, config OAuth2Config) {
				assert.True(t, config.IsSandbox())
			},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewOAuth2ConfigBuilder()
			config, err := builder.
				ClientID("test_client_id").
				RedirectURI("http://localhost:8080").
				Environment(tt.environment).
				Build()
			
			require.NoError(t, err)
			tt.checkFunc(t, config)
		})
	}
}

func TestOAuth2ConfigBuilder_CustomEndpoints(t *testing.T) {
	builder := NewOAuth2ConfigBuilder()
	config, err := builder.
		ClientID("test_client_id").
		RedirectURI("http://localhost:8080").
		CustomEndpoints("https://example.com/oauth/authorize", "https://example.com/oauth/token").
		Build()
	
	require.NoError(t, err)
	assert.Equal(t, "https://example.com/oauth/authorize", config.AuthURL)
	assert.Equal(t, "https://example.com/oauth/token", config.TokenURL)
	assert.Equal(t, "https://example.com", config.BaseURL)
}

func TestOAuth2ConfigBuilder_AutoGenerateState(t *testing.T) {
	builder := NewOAuth2ConfigBuilder()
	config, err := builder.
		ClientID("test_client_id").
		RedirectURI("http://localhost:8080").
		Production().
		Build()
	
	require.NoError(t, err)
	assert.NotEmpty(t, config.State)
	
	// Generate another config and ensure states are different
	config2, err := NewOAuth2ConfigBuilder().
		ClientID("test_client_id").
		RedirectURI("http://localhost:8080").
		Production().
		Build()
	
	require.NoError(t, err)
	assert.NotEqual(t, config.State, config2.State)
}

func TestOAuth2ConfigBuilder_ValidationError(t *testing.T) {
	builder := NewOAuth2ConfigBuilder()
	
	// Missing client ID should cause validation error
	_, err := builder.
		RedirectURI("http://localhost:8080").
		Production().
		Build()
	
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "client ID is required")
}

func TestConvenienceFunctions(t *testing.T) {
	// Test QuickProductionClient
	client, err := QuickProductionClient("test_client_id", "test_client_secret")
	require.NoError(t, err)
	assert.NotNil(t, client)
	assert.Equal(t, AuthModeOAuth2, client.authMode)
	
	config := client.oauth2Client.GetConfig()
	assert.True(t, config.IsProduction())
	assert.Equal(t, "test_client_id", config.ClientID)
	assert.Equal(t, "test_client_secret", config.ClientSecret)
	
	// Test QuickSandboxClient
	client, err = QuickSandboxClient("test_client_id", "test_client_secret")
	require.NoError(t, err)
	assert.NotNil(t, client)
	assert.Equal(t, AuthModeOAuth2, client.authMode)
	
	config = client.oauth2Client.GetConfig()
	assert.True(t, config.IsSandbox())
	assert.Equal(t, "test_client_id", config.ClientID)
	assert.Equal(t, "test_client_secret", config.ClientSecret)
}

func TestOAuth2Utils_FindAvailablePort(t *testing.T) {
	utils := NewOAuth2Utils()
	
	// Find an available port starting from 9000
	port, err := utils.FindAvailablePort(9000)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, port, 9000)
	assert.Less(t, port, 9100)
	
	// Verify the port is actually available
	assert.True(t, utils.isPortAvailable(port))
}

func TestOAuth2Utils_CreateQuickOAuth2Flow(t *testing.T) {
	utils := NewOAuth2Utils()
	
	// Test with default options
	client, server, err := utils.CreateQuickOAuth2Flow("test_client_id", "test_client_secret")
	require.NoError(t, err)
	assert.NotNil(t, client)
	assert.NotNil(t, server)
	
	// Clean up
	if server != nil {
		server.Shutdown(time.Second)
	}
	
	// Test with custom options
	opts := OAuth2FlowOptions{
		Environment:      "sandbox",
		UseBuiltinServer: false,
		RedirectURI:      "http://localhost:9999",
		Scopes:           []string{"read"},
	}
	
	client, server, err = utils.CreateQuickOAuth2Flow("test_client_id", "test_client_secret", opts)
	require.NoError(t, err)
	assert.NotNil(t, client)
	assert.Nil(t, server) // Should be nil when UseBuiltinServer is false
	
	config := client.GetConfig()
	assert.True(t, config.IsSandbox())
	assert.Equal(t, "http://localhost:9999", config.RedirectURI)
	assert.Equal(t, []string{"read"}, config.Scopes)
}

func TestOAuth2Utils_CreateQuickOAuth2Flow_Errors(t *testing.T) {
	utils := NewOAuth2Utils()
	
	// Test missing redirect URI when not using built-in server
	opts := OAuth2FlowOptions{
		UseBuiltinServer: false,
		// RedirectURI is missing
	}
	
	_, _, err := utils.CreateQuickOAuth2Flow("test_client_id", "test_client_secret", opts)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "redirect URI is required")
}

// Test helper functions
func TestExtractBaseURL(t *testing.T) {
	tests := []struct {
		name        string
		endpointURL string
		expected    string
	}{
		{
			name:        "HTTPS URL",
			endpointURL: "https://api.example.com/oauth/authorize",
			expected:    "https://api.example.com",
		},
		{
			name:        "HTTP URL with port",
			endpointURL: "http://localhost:8080/oauth/token",
			expected:    "http://localhost:8080",
		},
		{
			name:        "invalid URL",
			endpointURL: "ht tp://invalid",
			expected:    "",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractBaseURL(tt.endpointURL)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Integration test for complete OAuth2 flow simulation
func TestOAuth2Utils_IntegrationFlow(t *testing.T) {
	// Create a mock OAuth2 server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/oauth/authorize":
			// Simulate authorization endpoint
			query := r.URL.Query()
			redirectURI := query.Get("redirect_uri")
			state := query.Get("state")
			
			// Simulate user authorization by redirecting back with code
			callbackURL := fmt.Sprintf("%s?code=test_auth_code&state=%s", redirectURI, state)
			http.Redirect(w, r, callbackURL, http.StatusFound)
			
		case "/oauth/token":
			// Simulate token endpoint
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"access_token": "test_access_token",
				"refresh_token": "test_refresh_token",
				"token_type": "Bearer",
				"expires_in": 3600
			}`))
		}
	}))
	defer mockServer.Close()
	
	utils := NewOAuth2Utils()
	
	// Create OAuth2 flow with custom endpoints pointing to mock server
	opts := OAuth2FlowOptions{
		ClientID:       "test_client_id",
		ClientSecret:   "test_client_secret",
		CustomAuthURL:  mockServer.URL + "/oauth/authorize",
		CustomTokenURL: mockServer.URL + "/oauth/token",
		UseBuiltinServer: true,
		ServerPort:     0, // Use random port
		Timeout:        5 * time.Second,
	}
	
	client, server, err := utils.CreateQuickOAuth2Flow(opts.ClientID, opts.ClientSecret, opts)
	require.NoError(t, err)
	require.NotNil(t, client)
	require.NotNil(t, server)
	
	defer server.Shutdown(time.Second)
	
	// Get authorization URL
	authURL, err := client.GetAuthorizationURL()
	require.NoError(t, err)
	assert.Contains(t, authURL, mockServer.URL)
	
	// Simulate user clicking the authorization URL
	// This would normally open a browser, but we'll simulate the redirect
	httpClient := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// Don't follow redirects, we want to capture the redirect URL
			return http.ErrUseLastResponse
		},
	}
	
	resp, err := httpClient.Get(authURL)
	require.NoError(t, err)
	defer resp.Body.Close()
	
	// The mock server should redirect to our callback URL
	assert.Equal(t, http.StatusFound, resp.StatusCode)
	
	redirectURL := resp.Header.Get("Location")
	assert.NotEmpty(t, redirectURL)
	
	// Parse the authorization code from the redirect URL
	code, err := utils.ParseRedirectURL(redirectURL, client.GetState())
	require.NoError(t, err)
	assert.Equal(t, "test_auth_code", code)
	
	// Exchange code for tokens
	tokenResponse, err := client.ExchangeCodeForTokens(code)
	require.NoError(t, err)
	assert.Equal(t, "test_access_token", tokenResponse.AccessToken)
	assert.Equal(t, "test_refresh_token", tokenResponse.RefreshToken)
	assert.Equal(t, "Bearer", tokenResponse.TokenType)
	assert.Equal(t, 3600, tokenResponse.ExpiresIn)
	
	// Verify token manager has the tokens
	tokenManager := client.GetTokenManager()
	assert.True(t, tokenManager.HasValidToken())
	assert.True(t, tokenManager.HasRefreshToken())
	
	accessToken, err := tokenManager.GetAccessToken()
	require.NoError(t, err)
	assert.Equal(t, "test_access_token", accessToken)
}