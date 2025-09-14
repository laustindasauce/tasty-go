package tasty //nolint:testpackage // testing private field

import (
	"fmt"
	"math"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var (
	mux       *http.ServeMux
	server    *httptest.Server
	client    *Client
	testToken = "fake-access-token+C"
)

func setup() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	// Create OAuth2 client for testing
	config := OAuth2Config{
		ClientID:     "test_client_id",
		ClientSecret: "test_client_secret",
		RedirectURI:  "http://localhost:8080/callback",
		Scopes:       []string{"read", "trade"},
	}
	var err error
	client, err = NewClient(config, http.DefaultClient)
	if err != nil {
		panic(err)
	}

	// Set test tokens
	client.SetTokens(testToken, "refresh_token", 3600)

	client.baseURL = server.URL
	// Required for customRequest method
	client.baseHost = strings.Split(server.URL, "/")[2]
}

func teardown() {
	server.Close()
}

func TestDecodeError(t *testing.T) {
	r, err := http.NewRequest(http.MethodGet, apiBaseURL+"/customers/me/accounts", nil)
	require.NoError(t, err)

	r.Header.Add("Content-Type", "application/json")
	httpClient := http.Client{Timeout: time.Duration(30) * time.Second}
	resp, err := httpClient.Do(r)
	require.NoError(t, err)

	defer resp.Body.Close()

	require.True(t, containsInt(errorStatusCodes, resp.StatusCode))
	tastyErr := decodeError(resp)
	require.NotNil(t, tastyErr)

	require.Equal(t, 401, tastyErr.StatusCode)
	require.Equal(t, "token_invalid", tastyErr.Code)
	require.Equal(t, "This token is invalid or has expired", tastyErr.Message)
	require.Empty(t, tastyErr.Errors)

	require.Equal(
		t,
		"\nError in request 401;\nCode: token_invalid\nMessage: This token is invalid or has expired",
		tastyErr.Error(),
	)
}

func TestCustomRequest(t *testing.T) {
	config := OAuth2Config{
		ClientID:     "test_client_id",
		ClientSecret: "test_client_secret",
		RedirectURI:  "http://localhost:8080/callback",
		Scopes:       []string{"read", "trade"},
	}
	c, err := NewCertClient(config, &http.Client{Timeout: time.Duration(30) * time.Second})
	require.NoError(t, err)
	c.SetTokens(testToken, "refresh_token", 3600)

	// Test invalid payload
	invalid := math.Inf(1)
	httpResp, tastyError := c.customRequest(http.MethodGet, "/test", nil, invalid, nil)
	require.NotNil(t, tastyError)
	require.Nil(t, httpResp, "payload error")

	require.Equal(t,
		"\nError in request 0;\nCode: \nMessage: Client Side Error: json: unsupported value: +Inf",
		tastyError.Error())

	// Test invalid query
	httpResp, tastyError = c.customRequest(http.MethodGet, "/test", invalid, nil, nil)
	require.NotNil(t, tastyError)
	require.Nil(t, httpResp, "invalid query")

	require.Equal(t,
		"\nError in request 0;\nCode: \nMessage: Client Side Error: <nil>",
		tastyError.Error())

	// Test invalid method
	httpResp, tastyError = c.customRequest(http.MethodGet+"/sdfl/", "/test", nil, nil, nil)
	require.NotNil(t, tastyError)
	require.Nil(t, httpResp, "invalid method")

	require.Equal(t,
		"\nError in request 0;\nCode: \nMessage: Client Side Error: Get/sdfl/ \"https://api.cert.tastyworks.com/test\": net/http: invalid method \"GET/sdfl/\"",
		tastyError.Error())
}

func TestRequest(t *testing.T) {
	config := OAuth2Config{
		ClientID:     "test_client_id",
		ClientSecret: "test_client_secret",
		RedirectURI:  "http://localhost:8080/callback",
		Scopes:       []string{"read", "trade"},
	}
	c, err := NewCertClient(config, &http.Client{Timeout: time.Duration(30) * time.Second})
	require.NoError(t, err)

	httpResp, tastyError := c.request(http.MethodGet, "/no-auth", nil, nil, nil)
	require.NotNil(t, tastyError)
	require.Nil(t, httpResp)

	require.Equal(t,
		"\nError in request 0;\nCode: oauth2_token_error\nMessage: Failed to get access token: no access token available",
		tastyError.Error())

	c.SetTokens(testToken, "refresh_token", 3600)

	// Test invalid payload
	invalid := math.Inf(1)
	httpResp, tastyError = c.request(http.MethodGet, "/test", nil, invalid, nil)
	require.NotNil(t, tastyError)
	require.Nil(t, httpResp)

	require.Equal(t,
		"\nError in request 0;\nCode: \nMessage: Client Side Error: json: unsupported value: +Inf",
		tastyError.Error())

	// Test invalid query
	httpResp, tastyError = c.request(http.MethodGet, "/test", invalid, nil, nil)
	require.NotNil(t, tastyError)
	require.Nil(t, httpResp)

	require.Equal(t,
		"\nError in request 0;\nCode: \nMessage: Client Side Error: <nil>",
		tastyError.Error())

	// Test invalid method
	httpResp, tastyError = c.request(http.MethodGet+"/sdfl/", "/test", nil, nil, nil)
	require.NotNil(t, tastyError)
	require.Nil(t, httpResp)

	require.Equal(t,
		"\nError in request 0;\nCode: \nMessage: Client Side Error: net/http: invalid method \"GET/sdfl/\"",
		tastyError.Error())

	// Test invalid URL
	c.baseURL = "invalid"
	httpResp, tastyError = c.request(http.MethodGet, "/test", nil, nil, nil)
	require.NotNil(t, tastyError)
	require.Nil(t, httpResp)

	require.Equal(t,
		"\nError in request 0;\nCode: \nMessage: Client Side Error: Get \"invalid/test\": unsupported protocol scheme \"\"",
		tastyError.Error())
}

func TestCustomRequestNoContent(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/no-content", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusNoContent)
	})

	httpResp, err := client.customRequest(http.MethodGet, "/no-content", nil, nil, nil)
	require.Nil(t, err)
	require.NotNil(t, httpResp)
}

func TestRequestNoContent(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/no-content", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusNoContent)
	})

	httpResp, err := client.request(http.MethodGet, "/no-content", nil, nil, nil)
	require.Nil(t, err)
	require.NotNil(t, httpResp)
}

func TestNoAuthRequestNoContent(t *testing.T) {
	// This test is no longer applicable since all requests now require OAuth2 authentication
	t.Skip("NoAuthRequest method removed - all requests now use OAuth2 authentication")
}

func TestCustomRequestErrorResponses(t *testing.T) {
	setup()
	defer teardown()

	for _, errCode := range errorStatusCodes {
		path := fmt.Sprintf("/error/%d", errCode)
		mux.HandleFunc(path, func(writer http.ResponseWriter, request *http.Request) {
			writer.WriteHeader(errCode)
		})

		httpResp, err := client.customRequest(http.MethodGet, path, nil, nil, nil)
		require.NotNil(t, err)
		require.NotNil(t, httpResp)

		require.Equal(t, errCode, err.StatusCode)
	}
}

func TestRequestErrorResponses(t *testing.T) {
	setup()
	defer teardown()

	for _, errCode := range errorStatusCodes {
		path := fmt.Sprintf("/error/%d", errCode)
		mux.HandleFunc(path, func(writer http.ResponseWriter, request *http.Request) {
			writer.WriteHeader(errCode)
		})

		httpResp, err := client.request(http.MethodGet, path, nil, nil, nil)
		require.NotNil(t, err)
		require.NotNil(t, httpResp)

		require.Equal(t, errCode, err.StatusCode)
	}
}

func TestNoAuthRequestErrorResponses(t *testing.T) {
	// This test is no longer applicable since all requests now require OAuth2 authentication
	t.Skip("NoAuthRequest method removed - all requests now use OAuth2 authentication")
}

func TestCustomRequestInvalidResult(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/invalid", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, map[string]string{"test-key": "value"})
	})

	httpResp, err := client.customRequest(http.MethodGet, "/invalid", nil, nil, math.Inf(1))
	require.NotNil(t, err)
	require.NotNil(t, httpResp)
}

func TestRequestInvalidResult(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/invalid", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, map[string]string{"test-key": "value"})
	})

	httpResp, err := client.request(http.MethodGet, "/invalid", nil, nil, math.Inf(1))
	require.NotNil(t, err)
	require.NotNil(t, httpResp)
}

func TestNoAuthRequestInvalidResult(t *testing.T) {
	// This test is no longer applicable since all requests now require OAuth2 authentication
	t.Skip("NoAuthRequest method removed - all requests now use OAuth2 authentication")
}

func TestCustomRequestMissingCredentials(t *testing.T) {
	config := OAuth2Config{
		ClientID:     "test_client_id",
		ClientSecret: "test_client_secret",
		RedirectURI:  "http://localhost:8080/callback",
		Scopes:       []string{"read", "trade"},
	}
	c, err := NewClient(config, &http.Client{Timeout: time.Duration(30) * time.Second})
	require.NoError(t, err)

	httpResp, tastyErr := c.customRequest(http.MethodGet, "/invalid", nil, nil, nil)
	require.NotNil(t, tastyErr)
	require.Nil(t, httpResp)

	require.Equal(t,
		"\nError in request 0;\nCode: oauth2_token_error\nMessage: Failed to get access token: no access token available",
		tastyErr.Error())
}

func TestRequestMissingCredentials(t *testing.T) {
	config := OAuth2Config{
		ClientID:     "test_client_id",
		ClientSecret: "test_client_secret",
		RedirectURI:  "http://localhost:8080/callback",
		Scopes:       []string{"read", "trade"},
	}
	c, err := NewClient(config, &http.Client{Timeout: time.Duration(30) * time.Second})
	require.NoError(t, err)

	httpResp, tastyErr := c.request(http.MethodGet, "/invalid", nil, nil, nil)
	require.NotNil(t, tastyErr)
	require.Nil(t, httpResp)

	require.Equal(t,
		"\nError in request 0;\nCode: oauth2_token_error\nMessage: Failed to get access token: no access token available",
		tastyErr.Error())
}

func TestNoAuthRequestWithParams(t *testing.T) {
	// This test is no longer applicable since all requests now require OAuth2 authentication
	t.Skip("NoAuthRequest method removed - all requests now use OAuth2 authentication")
}

const tastyUnauthorizedError = `{
    "error": {
        "code": "unauthorized",
        "message": "Unauthorized. Unique customer support identifier: test"
    }
}`

const tastyInvalidCredentialsError = `{
    "error": {
        "code": "invalid_credentials",
        "message": "Invalid login, please check your username and password. Unique customer support identifier: test-id"
    }
}`

const tastyInvalidSessionError = `{
    "error": {
        "code": "invalid_session",
        "message": "Session user not present. Unique customer support identifier: test-id"
    }
}`

func expectedUnauthorized(t *testing.T, err error) {
	require.NotNil(t, err)

	require.Equal(t, "\nError in request 401;\nCode: unauthorized\nMessage: Unauthorized. Unique customer support identifier: test", err.Error())
}

func expectedInvalidCredentials(t *testing.T, err error) {
	require.NotNil(t, err)

	require.Equal(t,
		"\nError in request 401;\nCode: invalid_credentials\nMessage: Invalid login, please check your username and password. Unique customer support identifier: test-id",
		err.Error())
}

func expectedInvalidSession(t *testing.T, err error) {
	require.NotNil(t, err)

	require.Equal(t,
		"\nError in request 401;\nCode: invalid_session\nMessage: Session user not present. Unique customer support identifier: test-id",
		err.Error())
}

// OAuth2 Client Integration Tests

func TestNewClient(t *testing.T) {
	config := OAuth2Config{
		ClientID:     "test_client_id",
		ClientSecret: "test_client_secret",
		RedirectURI:  "http://localhost:8080/callback",
		Scopes:       []string{"read", "trade"},
	}

	client, err := NewClient(config, nil)
	require.NoError(t, err)
	require.NotNil(t, client)
	require.Equal(t, apiBaseURL, client.baseURL)
	require.Equal(t, apiBaseHost, client.baseHost)
	require.Equal(t, streamerBaseURL, client.websocket)

	// Test with custom HTTP client
	customClient := &http.Client{Timeout: time.Duration(60) * time.Second}
	client2, err := NewClient(config, customClient)
	require.NoError(t, err)
	require.Equal(t, customClient, client2.httpClient)
}

func TestNewClient_WithEndpoints(t *testing.T) {
	config := OAuth2Config{
		ClientID:     "test_client_id",
		ClientSecret: "test_client_secret",
		RedirectURI:  "http://localhost:8080/callback",
		Scopes:       []string{"read", "trade"},
		BaseURL:      apiBaseURL,
		AuthURL:      oauth2ProductionAuthURL,
		TokenURL:     oauth2ProductionTokenURL,
	}

	client, err := NewClient(config, nil)
	require.NoError(t, err)
	require.NotNil(t, client)
}

func TestNewClient_InvalidEndpoints(t *testing.T) {
	tests := []struct {
		name    string
		config  OAuth2Config
		wantErr string
	}{
		{
			name: "sandbox auth URL with production constructor",
			config: OAuth2Config{
				ClientID:     "test_client_id",
				ClientSecret: "test_client_secret",
				RedirectURI:  "http://localhost:8080/callback",
				AuthURL:      oauth2SandboxAuthURL,
			},
			wantErr: "NewClient requires production authorization URL",
		},
		{
			name: "sandbox token URL with production constructor",
			config: OAuth2Config{
				ClientID:     "test_client_id",
				ClientSecret: "test_client_secret",
				RedirectURI:  "http://localhost:8080/callback",
				TokenURL:     oauth2SandboxTokenURL,
			},
			wantErr: "NewClient requires production token URL",
		},
		{
			name: "cert base URL with production constructor",
			config: OAuth2Config{
				ClientID:     "test_client_id",
				ClientSecret: "test_client_secret",
				RedirectURI:  "http://localhost:8080/callback",
				BaseURL:      apiCertBaseURL,
			},
			wantErr: "use NewCertClient for sandbox environment",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(tt.config, nil)
			require.Error(t, err)
			require.Nil(t, client)
			require.Contains(t, err.Error(), tt.wantErr)
		})
	}
}

func TestNewCertClient(t *testing.T) {
	config := OAuth2Config{
		ClientID:     "test_client_id",
		ClientSecret: "test_client_secret",
		RedirectURI:  "http://localhost:8080/callback",
		Scopes:       []string{"read", "trade"},
	}

	client, err := NewCertClient(config, nil)
	require.NoError(t, err)
	require.NotNil(t, client)
	require.Equal(t, apiCertBaseURL, client.baseURL)
	require.Equal(t, apiCertBaseHost, client.baseHost)
	require.Equal(t, streamerCertBaseURL, client.websocket)

	// Test with custom HTTP client
	customClient := &http.Client{Timeout: time.Duration(60) * time.Second}
	client2, err := NewCertClient(config, customClient)
	require.NoError(t, err)
	require.Equal(t, customClient, client2.httpClient)
}

func TestNewCertClient_WithEndpoints(t *testing.T) {
	config := OAuth2Config{
		ClientID:     "test_client_id",
		ClientSecret: "test_client_secret",
		RedirectURI:  "http://localhost:8080/callback",
		Scopes:       []string{"read", "trade"},
		BaseURL:      apiCertBaseURL,
		AuthURL:      oauth2SandboxAuthURL,
		TokenURL:     oauth2SandboxTokenURL,
	}

	client, err := NewCertClient(config, nil)
	require.NoError(t, err)
	require.NotNil(t, client)
}

func TestNewCertClient_InvalidEndpoints(t *testing.T) {
	tests := []struct {
		name    string
		config  OAuth2Config
		wantErr string
	}{
		{
			name: "production auth URL with cert constructor",
			config: OAuth2Config{
				ClientID:     "test_client_id",
				ClientSecret: "test_client_secret",
				RedirectURI:  "http://localhost:8080/callback",
				AuthURL:      oauth2ProductionAuthURL,
			},
			wantErr: "NewCertClient requires sandbox authorization URL",
		},
		{
			name: "production token URL with cert constructor",
			config: OAuth2Config{
				ClientID:     "test_client_id",
				ClientSecret: "test_client_secret",
				RedirectURI:  "http://localhost:8080/callback",
				TokenURL:     oauth2ProductionTokenURL,
			},
			wantErr: "NewCertClient requires sandbox token URL",
		},
		{
			name: "production base URL with cert constructor",
			config: OAuth2Config{
				ClientID:     "test_client_id",
				ClientSecret: "test_client_secret",
				RedirectURI:  "http://localhost:8080/callback",
				BaseURL:      apiBaseURL,
			},
			wantErr: "use NewClient for production environment",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewCertClient(tt.config, nil)
			require.Error(t, err)
			require.Nil(t, client)
			require.Contains(t, err.Error(), tt.wantErr)
		})
	}
}

func TestClient_AuthMode(t *testing.T) {
	// Test OAuth2 mode (now the only mode)
	config := OAuth2Config{
		ClientID:     "test_client_id",
		ClientSecret: "test_client_secret",
		RedirectURI:  "http://localhost:8080/callback",
		Scopes:       []string{"read", "trade"},
	}
	client, err := NewClient(config, nil)
	require.NoError(t, err)

	// All clients are now OAuth2-only
	require.True(t, client.IsAuthenticated() || !client.HasValidToken()) // Either has tokens or doesn't
}

func TestClient_OAuth2Methods_InvalidConfig(t *testing.T) {
	// Test with invalid config
	config := OAuth2Config{
		// Missing required fields
	}
	_, err := NewClient(config, nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid OAuth2 configuration")
}

func TestClient_OAuth2Methods(t *testing.T) {
	config := OAuth2Config{
		ClientID:     "test_client_id",
		ClientSecret: "test_client_secret",
		RedirectURI:  "http://localhost:8080/callback",
		Scopes:       []string{"read", "trade"},
		State:        "test_state",
	}
	client, err := NewClient(config, nil)
	require.NoError(t, err)

	// Test OAuth2 methods work
	authURL, err := client.GetAuthorizationURL()
	require.NoError(t, err)
	require.Contains(t, authURL, "client_id=test_client_id")
	require.Contains(t, authURL, "state=test_state")

	err = client.ValidateState("test_state")
	require.NoError(t, err)

	err = client.ValidateState("wrong_state")
	require.Error(t, err)
}

func TestClient_IsAuthenticated(t *testing.T) {
	config := OAuth2Config{
		ClientID:     "test_client_id",
		ClientSecret: "test_client_secret",
		RedirectURI:  "http://localhost:8080/callback",
		Scopes:       []string{"read", "trade"},
	}

	client, err := NewClient(config, nil)
	client.ClearTokens()
	require.NoError(t, err)
	require.False(t, client.IsAuthenticated())

	// Set tokens to make it authenticated
	client.SetTokens("access_token", "refresh_token", 3600)
	require.True(t, client.IsAuthenticated())
}

func TestClient_ClearAuthentication(t *testing.T) {
	config := OAuth2Config{
		ClientID:     "test_client_id",
		ClientSecret: "test_client_secret",
		RedirectURI:  "http://localhost:8080/callback",
		Scopes:       []string{"read", "trade"},
	}
	client, err := NewClient(config, nil)
	require.NoError(t, err)

	client.SetTokens("access_token", "refresh_token", 3600)
	require.True(t, client.IsAuthenticated())

	client.ClearTokens()
	require.False(t, client.IsAuthenticated())
}

func TestClient_TokenMethods(t *testing.T) {
	config := OAuth2Config{
		ClientID:     "test_client_id",
		ClientSecret: "test_client_secret",
		RedirectURI:  "http://localhost:8080/callback",
		Scopes:       []string{"read", "trade"},
	}
	client, err := NewClient(config, nil)
	require.NoError(t, err)

	// Test token methods
	require.False(t, client.HasValidToken())
	require.False(t, client.HasRefreshToken())
	require.True(t, client.IsTokenExpired())

	client.SetTokens("access_token", "refresh_token", 3600)
	require.True(t, client.HasValidToken())
	require.True(t, client.HasRefreshToken())
	require.False(t, client.IsTokenExpired())
}

func TestClient_WithTokens(t *testing.T) {
	config := OAuth2Config{
		ClientID:     "test_client_id",
		ClientSecret: "test_client_secret",
		RedirectURI:  "http://localhost:8080/callback",
		Scopes:       []string{"read", "trade"},
	}

	// Test NewClientWithTokens
	client, err := NewClientWithTokens(config, "access_token", "refresh_token", 3600, nil)
	require.NoError(t, err)
	require.NotNil(t, client)
	require.True(t, client.IsAuthenticated())

	// Test NewCertClientWithTokens
	certClient, err := NewCertClientWithTokens(config, "access_token", "refresh_token", 3600, nil)
	require.NoError(t, err)
	require.NotNil(t, certClient)
	require.True(t, certClient.IsAuthenticated())
}
