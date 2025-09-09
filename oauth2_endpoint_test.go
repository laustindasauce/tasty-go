package tasty

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestOAuth2EndpointConstants verifies that the OAuth2 endpoint constants are correctly defined
func TestOAuth2EndpointConstants(t *testing.T) {
	// Test production endpoints
	assert.Equal(t, "https://my.tastytrade.com/auth.html", oauth2ProductionAuthURL)
	assert.Equal(t, "https://api.tastyworks.com/oauth/token", oauth2ProductionTokenURL)
	
	// Test sandbox endpoints
	assert.Equal(t, "https://cert-my.staging-tasty.works/auth.html", oauth2SandboxAuthURL)
	assert.Equal(t, "https://api.cert.tastyworks.com/oauth/token", oauth2SandboxTokenURL)
	
	// Verify endpoints are different between environments
	assert.NotEqual(t, oauth2ProductionAuthURL, oauth2SandboxAuthURL)
	assert.NotEqual(t, oauth2ProductionTokenURL, oauth2SandboxTokenURL)
}

// TestNewProductionOAuth2Config tests the production OAuth2 config helper
func TestNewProductionOAuth2Config(t *testing.T) {
	config := NewProductionOAuth2Config("client_id", "client_secret", "http://localhost:8080/callback", []string{"read", "trade"})
	
	assert.Equal(t, "client_id", config.ClientID)
	assert.Equal(t, "client_secret", config.ClientSecret)
	assert.Equal(t, "http://localhost:8080/callback", config.RedirectURI)
	assert.Equal(t, []string{"read", "trade"}, config.Scopes)
	assert.Equal(t, apiBaseURL, config.BaseURL)
	assert.Equal(t, oauth2ProductionAuthURL, config.AuthURL)
	assert.Equal(t, oauth2ProductionTokenURL, config.TokenURL)
	
	// Test environment detection
	assert.True(t, config.IsProduction())
	assert.False(t, config.IsSandbox())
	assert.Equal(t, "production", config.GetEnvironment())
}

// TestNewSandboxOAuth2Config tests the sandbox OAuth2 config helper
func TestNewSandboxOAuth2Config(t *testing.T) {
	config := NewSandboxOAuth2Config("client_id", "client_secret", "http://localhost:8080/callback", []string{"read", "trade"})
	
	assert.Equal(t, "client_id", config.ClientID)
	assert.Equal(t, "client_secret", config.ClientSecret)
	assert.Equal(t, "http://localhost:8080/callback", config.RedirectURI)
	assert.Equal(t, []string{"read", "trade"}, config.Scopes)
	assert.Equal(t, apiCertBaseURL, config.BaseURL)
	assert.Equal(t, oauth2SandboxAuthURL, config.AuthURL)
	assert.Equal(t, oauth2SandboxTokenURL, config.TokenURL)
	
	// Test environment detection
	assert.False(t, config.IsProduction())
	assert.True(t, config.IsSandbox())
	assert.Equal(t, "sandbox", config.GetEnvironment())
}

// TestNewProductionOAuth2Config_DefaultScopes tests default scope handling
func TestNewProductionOAuth2Config_DefaultScopes(t *testing.T) {
	config := NewProductionOAuth2Config("client_id", "client_secret", "http://localhost:8080/callback", nil)
	
	assert.Equal(t, []string{defaultScope}, config.Scopes)
}

// TestNewSandboxOAuth2Config_DefaultScopes tests default scope handling
func TestNewSandboxOAuth2Config_DefaultScopes(t *testing.T) {
	config := NewSandboxOAuth2Config("client_id", "client_secret", "http://localhost:8080/callback", nil)
	
	assert.Equal(t, []string{defaultScope}, config.Scopes)
}

// TestOAuth2Config_SetEnvironment tests environment switching
func TestOAuth2Config_SetEnvironment(t *testing.T) {
	config := OAuth2Config{
		ClientID:    "client_id",
		RedirectURI: "http://localhost:8080/callback",
	}
	
	// Test setting production environment
	err := config.SetEnvironment("production")
	require.NoError(t, err)
	assert.Equal(t, apiBaseURL, config.BaseURL)
	assert.Equal(t, oauth2ProductionAuthURL, config.AuthURL)
	assert.Equal(t, oauth2ProductionTokenURL, config.TokenURL)
	assert.True(t, config.IsProduction())
	
	// Test setting sandbox environment
	err = config.SetEnvironment("sandbox")
	require.NoError(t, err)
	assert.Equal(t, apiCertBaseURL, config.BaseURL)
	assert.Equal(t, oauth2SandboxAuthURL, config.AuthURL)
	assert.Equal(t, oauth2SandboxTokenURL, config.TokenURL)
	assert.True(t, config.IsSandbox())
	
	// Test alternative names
	err = config.SetEnvironment("prod")
	require.NoError(t, err)
	assert.True(t, config.IsProduction())
	
	err = config.SetEnvironment("cert")
	require.NoError(t, err)
	assert.True(t, config.IsSandbox())
	
	// Test invalid environment
	err = config.SetEnvironment("invalid")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid environment")
}

// TestOAuth2Config_ValidateEnvironmentConsistency tests environment consistency validation
func TestOAuth2Config_ValidateEnvironmentConsistency(t *testing.T) {
	tests := []struct {
		name      string
		config    OAuth2Config
		wantError bool
		errorMsg  string
	}{
		{
			name: "consistent production config",
			config: OAuth2Config{
				ClientID:    "client_id",
				RedirectURI: "http://localhost:8080/callback",
				BaseURL:     apiBaseURL,
				AuthURL:     oauth2ProductionAuthURL,
				TokenURL:    oauth2ProductionTokenURL,
			},
			wantError: false,
		},
		{
			name: "consistent sandbox config",
			config: OAuth2Config{
				ClientID:    "client_id",
				RedirectURI: "http://localhost:8080/callback",
				BaseURL:     apiCertBaseURL,
				AuthURL:     oauth2SandboxAuthURL,
				TokenURL:    oauth2SandboxTokenURL,
			},
			wantError: false,
		},
		{
			name: "mixed production and sandbox URLs",
			config: OAuth2Config{
				ClientID:    "client_id",
				RedirectURI: "http://localhost:8080/callback",
				BaseURL:     apiBaseURL,
				AuthURL:     oauth2SandboxAuthURL,
				TokenURL:    oauth2ProductionTokenURL,
			},
			wantError: true,
			errorMsg:  "configuration mixes production and sandbox endpoints",
		},
		{
			name: "empty config (valid)",
			config: OAuth2Config{
				ClientID:    "client_id",
				RedirectURI: "http://localhost:8080/callback",
			},
			wantError: false,
		},
		{
			name: "partial production config (valid)",
			config: OAuth2Config{
				ClientID:    "client_id",
				RedirectURI: "http://localhost:8080/callback",
				AuthURL:     oauth2ProductionAuthURL,
			},
			wantError: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.wantError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestValidateEndpointURL tests endpoint URL validation
func TestValidateEndpointURL(t *testing.T) {
	tests := []struct {
		name         string
		url          string
		endpointType string
		wantError    bool
		errorMsg     string
	}{
		{
			name:         "valid production auth URL",
			url:          oauth2ProductionAuthURL,
			endpointType: "authorization",
			wantError:    false,
		},
		{
			name:         "valid sandbox token URL",
			url:          oauth2SandboxTokenURL,
			endpointType: "token",
			wantError:    false,
		},
		{
			name:         "valid localhost URL",
			url:          "http://localhost:8080/oauth/token",
			endpointType: "token",
			wantError:    false,
		},
		{
			name:         "valid test domain URL",
			url:          "https://api.test.com/oauth/authorize",
			endpointType: "authorization",
			wantError:    false,
		},
		{
			name:         "empty URL",
			url:          "",
			endpointType: "authorization",
			wantError:    true,
			errorMsg:     "authorization URL cannot be empty",
		},
		{
			name:         "invalid URL format",
			url:          "not-a-url",
			endpointType: "token",
			wantError:    true,
			errorMsg:     "token URL must include a scheme",
		},
		{
			name:         "missing scheme",
			url:          "api.tastyworks.com/oauth/authorize",
			endpointType: "authorization",
			wantError:    true,
			errorMsg:     "authorization URL must include a scheme",
		},
		{
			name:         "missing host",
			url:          "https:///oauth/token",
			endpointType: "token",
			wantError:    true,
			errorMsg:     "token URL must include a host",
		},
		{
			name:         "TastyTrade domain with HTTP (invalid)",
			url:          "http://api.tastyworks.com/oauth/authorize",
			endpointType: "authorization",
			wantError:    true,
			errorMsg:     "authorization URL must use HTTPS",
		},
		{
			name:         "non-TastyTrade HTTPS domain (invalid for production)",
			url:          "https://evil.com/oauth/token",
			endpointType: "token",
			wantError:    true,
			errorMsg:     "token URL must be a valid TastyTrade endpoint or test domain",
		},
		{
			name:         "HTTP with non-localhost (invalid)",
			url:          "http://malicious.com/oauth/authorize",
			endpointType: "authorization",
			wantError:    true,
			errorMsg:     "authorization URL with HTTP scheme is only allowed for localhost or test domains",
		},
		{
			name:         "invalid scheme",
			url:          "ftp://malicious.com/oauth/token",
			endpointType: "token",
			wantError:    true,
			errorMsg:     "token URL must use HTTP or HTTPS scheme",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateEndpointURL(tt.url, tt.endpointType)
			if tt.wantError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestValidateRedirectURI tests redirect URI validation
func TestValidateRedirectURI(t *testing.T) {
	tests := []struct {
		name        string
		redirectURI string
		wantError   bool
		errorMsg    string
	}{
		{
			name:        "valid HTTPS URI",
			redirectURI: "https://example.com/callback",
			wantError:   false,
		},
		{
			name:        "valid localhost HTTP URI",
			redirectURI: "http://localhost:8080/callback",
			wantError:   false,
		},
		{
			name:        "valid 127.0.0.1 HTTP URI",
			redirectURI: "http://127.0.0.1:3000/callback",
			wantError:   false,
		},
		{
			name:        "valid custom scheme",
			redirectURI: "myapp://oauth/callback",
			wantError:   false,
		},
		{
			name:        "empty URI",
			redirectURI: "",
			wantError:   true,
			errorMsg:    "redirect URI cannot be empty",
		},
		{
			name:        "invalid URI format",
			redirectURI: "not-a-uri",
			wantError:   true,
			errorMsg:    "redirect URI must include a scheme",
		},
		{
			name:        "HTTP with non-localhost",
			redirectURI: "http://example.com/callback",
			wantError:   true,
			errorMsg:    "HTTP redirect URIs are only allowed for localhost",
		},
		{
			name:        "missing scheme",
			redirectURI: "example.com/callback",
			wantError:   true,
			errorMsg:    "redirect URI must include a scheme",
		},
		{
			name:        "short custom scheme",
			redirectURI: "ab://callback",
			wantError:   true,
			errorMsg:    "custom URI schemes must be at least 3 characters long",
		},
		{
			name:        "scheme only",
			redirectURI: "https://",
			wantError:   true,
			errorMsg:    "redirect URI must include a host or path",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateRedirectURI(tt.redirectURI)
			if tt.wantError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestIsLocalhost tests localhost detection
func TestIsLocalhost(t *testing.T) {
	tests := []struct {
		name     string
		host     string
		expected bool
	}{
		{"localhost", "localhost", true},
		{"localhost with port", "localhost:8080", true},
		{"127.0.0.1", "127.0.0.1", true},
		{"127.0.0.1 with port", "127.0.0.1:3000", true},
		{"example.com", "example.com", false},
		{"api.tastyworks.com", "api.tastyworks.com", false},
		{"empty", "", false},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isLocalhost(tt.host)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestIsTestDomain tests test domain detection
func TestIsTestDomain(t *testing.T) {
	tests := []struct {
		name     string
		host     string
		expected bool
	}{
		{"api.test.com", "api.test.com", true},
		{"test.example.com", "test.example.com", true},
		{"example.com", "example.com", true},
		{"httpbin.org", "httpbin.org", true},
		{"subdomain.test", "subdomain.test", true},
		{"something.test", "something.test", true},
		{"api.tastyworks.com", "api.tastyworks.com", false},
		{"google.com", "google.com", false},
		{"empty", "", false},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isTestDomain(tt.host)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestOAuth2Config_GetEnvironment tests environment detection edge cases
func TestOAuth2Config_GetEnvironment(t *testing.T) {
	tests := []struct {
		name     string
		config   OAuth2Config
		expected string
	}{
		{
			name: "production by base URL",
			config: OAuth2Config{
				BaseURL: apiBaseURL,
			},
			expected: "production",
		},
		{
			name: "production by auth URL",
			config: OAuth2Config{
				AuthURL: oauth2ProductionAuthURL,
			},
			expected: "production",
		},
		{
			name: "production by token URL",
			config: OAuth2Config{
				TokenURL: oauth2ProductionTokenURL,
			},
			expected: "production",
		},
		{
			name: "sandbox by base URL",
			config: OAuth2Config{
				BaseURL: apiCertBaseURL,
			},
			expected: "sandbox",
		},
		{
			name: "sandbox by auth URL",
			config: OAuth2Config{
				AuthURL: oauth2SandboxAuthURL,
			},
			expected: "sandbox",
		},
		{
			name: "sandbox by token URL",
			config: OAuth2Config{
				TokenURL: oauth2SandboxTokenURL,
			},
			expected: "sandbox",
		},
		{
			name: "unknown environment",
			config: OAuth2Config{
				BaseURL: "https://unknown.com",
			},
			expected: "unknown",
		},
		{
			name:     "empty config",
			config:   OAuth2Config{},
			expected: "unknown",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.config.GetEnvironment()
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestOAuth2ClientConstructors_EndpointConfiguration tests that client constructors properly configure endpoints
func TestOAuth2ClientConstructors_EndpointConfiguration(t *testing.T) {
	t.Run("NewOAuth2Client sets production endpoints", func(t *testing.T) {
		config := OAuth2Config{
			ClientID:    "test_client_id",
			RedirectURI: "http://localhost:8080/callback",
			Scopes:      []string{"read", "trade"},
		}
		
		client, err := NewOAuth2Client(config, nil)
		require.NoError(t, err)
		require.NotNil(t, client)
		
		// Verify the client is configured for production
		assert.Equal(t, apiBaseURL, client.baseURL)
		assert.Equal(t, apiBaseHost, client.baseHost)
		assert.Equal(t, streamerBaseURL, client.websocket)
		
		// Verify OAuth2 client has production endpoints
		oauth2Config := client.oauth2Client.GetConfig()
		assert.Equal(t, oauth2ProductionAuthURL, oauth2Config.AuthURL)
		assert.Equal(t, oauth2ProductionTokenURL, oauth2Config.TokenURL)
		assert.Equal(t, apiBaseURL, oauth2Config.BaseURL)
	})
	
	t.Run("NewCertOAuth2Client sets sandbox endpoints", func(t *testing.T) {
		config := OAuth2Config{
			ClientID:    "test_client_id",
			RedirectURI: "http://localhost:8080/callback",
			Scopes:      []string{"read", "trade"},
		}
		
		client, err := NewCertOAuth2Client(config, nil)
		require.NoError(t, err)
		require.NotNil(t, client)
		
		// Verify the client is configured for sandbox
		assert.Equal(t, apiCertBaseURL, client.baseURL)
		assert.Equal(t, apiCertBaseHost, client.baseHost)
		assert.Equal(t, streamerCertBaseURL, client.websocket)
		
		// Verify OAuth2 client has sandbox endpoints
		oauth2Config := client.oauth2Client.GetConfig()
		assert.Equal(t, oauth2SandboxAuthURL, oauth2Config.AuthURL)
		assert.Equal(t, oauth2SandboxTokenURL, oauth2Config.TokenURL)
		assert.Equal(t, apiCertBaseURL, oauth2Config.BaseURL)
	})
}

// TestOAuth2ClientConstructors_EndpointValidation tests endpoint validation in constructors
func TestOAuth2ClientConstructors_EndpointValidation(t *testing.T) {
	t.Run("NewOAuth2Client rejects sandbox endpoints", func(t *testing.T) {
		tests := []struct {
			name    string
			config  OAuth2Config
			wantErr string
		}{
			{
				name: "sandbox auth URL",
				config: OAuth2Config{
					ClientID:    "test_client_id",
					RedirectURI: "http://localhost:8080/callback",
					AuthURL:     oauth2SandboxAuthURL,
				},
				wantErr: "NewOAuth2Client requires production authorization URL",
			},
			{
				name: "sandbox token URL",
				config: OAuth2Config{
					ClientID:    "test_client_id",
					RedirectURI: "http://localhost:8080/callback",
					TokenURL:    oauth2SandboxTokenURL,
				},
				wantErr: "NewOAuth2Client requires production token URL",
			},
			{
				name: "sandbox base URL",
				config: OAuth2Config{
					ClientID:    "test_client_id",
					RedirectURI: "http://localhost:8080/callback",
					BaseURL:     apiCertBaseURL,
				},
				wantErr: "use NewCertOAuth2Client for sandbox environment",
			},
		}
		
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				client, err := NewOAuth2Client(tt.config, nil)
				assert.Error(t, err)
				assert.Nil(t, client)
				assert.Contains(t, err.Error(), tt.wantErr)
			})
		}
	})
	
	t.Run("NewCertOAuth2Client rejects production endpoints", func(t *testing.T) {
		tests := []struct {
			name    string
			config  OAuth2Config
			wantErr string
		}{
			{
				name: "production auth URL",
				config: OAuth2Config{
					ClientID:    "test_client_id",
					RedirectURI: "http://localhost:8080/callback",
					AuthURL:     oauth2ProductionAuthURL,
				},
				wantErr: "NewCertOAuth2Client requires sandbox authorization URL",
			},
			{
				name: "production token URL",
				config: OAuth2Config{
					ClientID:    "test_client_id",
					RedirectURI: "http://localhost:8080/callback",
					TokenURL:    oauth2ProductionTokenURL,
				},
				wantErr: "NewCertOAuth2Client requires sandbox token URL",
			},
			{
				name: "production base URL",
				config: OAuth2Config{
					ClientID:    "test_client_id",
					RedirectURI: "http://localhost:8080/callback",
					BaseURL:     apiBaseURL,
				},
				wantErr: "use NewOAuth2Client for production environment",
			},
		}
		
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				client, err := NewCertOAuth2Client(tt.config, nil)
				assert.Error(t, err)
				assert.Nil(t, client)
				assert.Contains(t, err.Error(), tt.wantErr)
			})
		}
	})
}