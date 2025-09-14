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
	client = NewClient(http.DefaultClient)
	client.Session = Session{
		SessionToken: &testToken,
	}
	client.baseURL = server.URL
	// Required for customRequest method
	client.baseHost = strings.Split(server.URL, "/")[2]
}

func teardown() {
	server.Close()
}

func TestTastyCertSession(t *testing.T) {
	c := NewCertClient(nil)

	require.NotNil(t, c.httpClient)
	require.Equal(t, apiCertBaseURL, c.baseURL)
	require.Equal(t, apiCertBaseHost, c.baseHost)
	require.Equal(t, streamerCertBaseURL, c.websocket)
	require.Equal(t, streamerCertBaseURL, c.GetWebsocketURL())

	cWithHTTP := NewCertClient(&http.Client{Timeout: time.Duration(30) * time.Second})

	require.NotNil(t, cWithHTTP.httpClient)
	require.Equal(t, time.Duration(30)*time.Second, cWithHTTP.httpClient.Timeout)
}

func TestTastySession(t *testing.T) {
	c := NewClient(nil)

	require.NotNil(t, c.httpClient)
	require.Equal(t, apiBaseURL, c.baseURL)
	require.Equal(t, apiBaseHost, c.baseHost)
	require.Equal(t, streamerBaseURL, c.websocket)
	require.Equal(t, streamerBaseURL, c.GetWebsocketURL())

	cWithHTTP := NewClient(&http.Client{Timeout: time.Duration(30) * time.Second})

	require.NotNil(t, cWithHTTP.httpClient)
	require.Equal(t, time.Duration(30)*time.Second, cWithHTTP.httpClient.Timeout)
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
	c := NewCertClient(&http.Client{Timeout: time.Duration(30) * time.Second})
	c.Session.SessionToken = &testToken

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
	c := NewCertClient(&http.Client{Timeout: time.Duration(30) * time.Second})

	httpResp, tastyError := c.request(http.MethodGet, "/no-auth", nil, nil, nil)
	require.NotNil(t, tastyError)
	require.Nil(t, httpResp)

	require.Equal(t,
		"\nError in request 0;\nCode: invalid_session\nMessage: Session is invalid: Session Token cannot be nil.",
		tastyError.Error())

	c.Session.SessionToken = &testToken

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

func TestNoAuthRequest(t *testing.T) {
	c := NewCertClient(&http.Client{Timeout: time.Duration(30) * time.Second})

	// Test invalid payload
	invalid := math.Inf(1)
	httpResp, tastyError := c.noAuthRequest(http.MethodGet, "/test", nil, nil, invalid, nil)
	require.NotNil(t, tastyError)
	require.Nil(t, httpResp)

	require.Equal(t,
		"\nError in request 0;\nCode: \nMessage: Client Side Error: json: unsupported value: +Inf",
		tastyError.Error())

	// Test invalid query
	httpResp, tastyError = c.noAuthRequest(http.MethodGet, "/test", nil, invalid, nil, nil)
	require.NotNil(t, tastyError)
	require.Nil(t, httpResp)

	require.Equal(t,
		"\nError in request 0;\nCode: \nMessage: Client Side Error: <nil>",
		tastyError.Error())

	// Test invalid method
	httpResp, tastyError = c.noAuthRequest(http.MethodGet+"/sdfl/", "/test", nil, nil, nil, nil)
	require.NotNil(t, tastyError)
	require.Nil(t, httpResp)

	require.Equal(t,
		"\nError in request 0;\nCode: \nMessage: Client Side Error: net/http: invalid method \"GET/sdfl/\"",
		tastyError.Error())

	// Test invalid URL
	c.baseURL = "invalid"
	httpResp, tastyError = c.noAuthRequest(http.MethodGet, "/test", nil, nil, nil, nil)
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
	setup()
	defer teardown()

	mux.HandleFunc("/no-content", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusNoContent)
	})

	httpResp, err := client.noAuthRequest(http.MethodGet, "/no-content", nil, nil, nil, nil)
	require.Nil(t, err)
	require.NotNil(t, httpResp)
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
	setup()
	defer teardown()

	for _, errCode := range errorStatusCodes {
		path := fmt.Sprintf("/error/%d", errCode)
		mux.HandleFunc(path, func(writer http.ResponseWriter, request *http.Request) {
			writer.WriteHeader(errCode)
		})

		httpResp, err := client.noAuthRequest(http.MethodGet, path, nil, nil, nil, nil)
		require.NotNil(t, err)
		require.NotNil(t, httpResp)

		require.Equal(t, errCode, err.StatusCode)
	}
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
	setup()
	defer teardown()

	mux.HandleFunc("/invalid", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, map[string]string{"test-key": "value"})
	})

	httpResp, err := client.noAuthRequest(http.MethodGet, "/invalid", nil, nil, nil, math.Inf(1))
	require.NotNil(t, err)
	require.NotNil(t, httpResp)
}

func TestCustomRequestMissingCredentials(t *testing.T) {
	c := NewClient(&http.Client{Timeout: time.Duration(30) * time.Second})

	httpResp, tastyErr := c.customRequest(http.MethodGet, "/invalid", nil, nil, nil)
	require.NotNil(t, tastyErr)
	require.Nil(t, httpResp)

	require.Equal(t,
		"\nError in request 0;\nCode: invalid_session\nMessage: Session is invalid: Session Token cannot be nil.",
		tastyErr.Error())
}

func TestRequestMissingCredentials(t *testing.T) {
	c := NewClient(&http.Client{Timeout: time.Duration(30) * time.Second})

	httpResp, tastyErr := c.customRequest(http.MethodGet, "/invalid", nil, nil, nil)
	require.NotNil(t, tastyErr)
	require.Nil(t, httpResp)

	require.Equal(t,
		"\nError in request 0;\nCode: invalid_session\nMessage: Session is invalid: Session Token cannot be nil.",
		tastyErr.Error())
}

func TestNoAuthRequestWithParams(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/with-params", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusNoContent)
		require.Equal(t, "AAPL", request.URL.Query().Get("symbol"))
	})

	httpResp, err := client.noAuthRequest(http.MethodGet, "/with-params", nil, AccountPositionQuery{Symbol: "AAPL"}, nil, nil)
	require.Nil(t, err)
	require.NotNil(t, httpResp)
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

func TestNewOAuth2Client(t *testing.T) {
	config := OAuth2Config{
		ClientID:     "test_client_id",
		ClientSecret: "test_client_secret",
		RedirectURI:  "http://localhost:8080/callback",
		Scopes:       []string{"read", "trade"},
	}

	client, err := NewOAuth2Client(config, nil)
	require.NoError(t, err)
	require.NotNil(t, client)
	require.Equal(t, AuthModeOAuth2, client.authMode)
	require.Equal(t, apiBaseURL, client.baseURL)
	require.Equal(t, apiBaseHost, client.baseHost)
	require.Equal(t, streamerBaseURL, client.websocket)
	require.NotNil(t, client.oauth2Client)
	require.True(t, client.IsOAuth2Mode())
	require.False(t, client.IsSessionMode())

	// Test with custom HTTP client
	customClient := &http.Client{Timeout: time.Duration(60) * time.Second}
	client2, err := NewOAuth2Client(config, customClient)
	require.NoError(t, err)
	require.Equal(t, customClient, client2.httpClient)
}

func TestNewOAuth2Client_WithEndpoints(t *testing.T) {
	config := OAuth2Config{
		ClientID:     "test_client_id",
		ClientSecret: "test_client_secret",
		RedirectURI:  "http://localhost:8080/callback",
		Scopes:       []string{"read", "trade"},
		BaseURL:      apiBaseURL,
		AuthURL:      oauth2ProductionAuthURL,
		TokenURL:     oauth2ProductionTokenURL,
	}

	client, err := NewOAuth2Client(config, nil)
	require.NoError(t, err)
	require.NotNil(t, client)
	require.Equal(t, AuthModeOAuth2, client.authMode)
}

func TestNewOAuth2Client_InvalidEndpoints(t *testing.T) {
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
			wantErr: "NewOAuth2Client requires production authorization URL",
		},
		{
			name: "sandbox token URL with production constructor",
			config: OAuth2Config{
				ClientID:     "test_client_id",
				ClientSecret: "test_client_secret",
				RedirectURI:  "http://localhost:8080/callback",
				TokenURL:     oauth2SandboxTokenURL,
			},
			wantErr: "NewOAuth2Client requires production token URL",
		},
		{
			name: "cert base URL with production constructor",
			config: OAuth2Config{
				ClientID:     "test_client_id",
				ClientSecret: "test_client_secret",
				RedirectURI:  "http://localhost:8080/callback",
				BaseURL:      apiCertBaseURL,
			},
			wantErr: "use NewCertOAuth2Client for sandbox environment",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewOAuth2Client(tt.config, nil)
			require.Error(t, err)
			require.Nil(t, client)
			require.Contains(t, err.Error(), tt.wantErr)
		})
	}
}

func TestNewCertOAuth2Client(t *testing.T) {
	config := OAuth2Config{
		ClientID:     "test_client_id",
		ClientSecret: "test_client_secret",
		RedirectURI:  "http://localhost:8080/callback",
		Scopes:       []string{"read", "trade"},
	}

	client, err := NewCertOAuth2Client(config, nil)
	require.NoError(t, err)
	require.NotNil(t, client)
	require.Equal(t, AuthModeOAuth2, client.authMode)
	require.Equal(t, apiCertBaseURL, client.baseURL)
	require.Equal(t, apiCertBaseHost, client.baseHost)
	require.Equal(t, streamerCertBaseURL, client.websocket)
	require.NotNil(t, client.oauth2Client)
	require.True(t, client.IsOAuth2Mode())
	require.False(t, client.IsSessionMode())

	// Test with custom HTTP client
	customClient := &http.Client{Timeout: time.Duration(60) * time.Second}
	client2, err := NewCertOAuth2Client(config, customClient)
	require.NoError(t, err)
	require.Equal(t, customClient, client2.httpClient)
}

func TestNewCertOAuth2Client_WithEndpoints(t *testing.T) {
	config := OAuth2Config{
		ClientID:     "test_client_id",
		ClientSecret: "test_client_secret",
		RedirectURI:  "http://localhost:8080/callback",
		Scopes:       []string{"read", "trade"},
		BaseURL:      apiCertBaseURL,
		AuthURL:      oauth2SandboxAuthURL,
		TokenURL:     oauth2SandboxTokenURL,
	}

	client, err := NewCertOAuth2Client(config, nil)
	require.NoError(t, err)
	require.NotNil(t, client)
	require.Equal(t, AuthModeOAuth2, client.authMode)
}

func TestNewCertOAuth2Client_InvalidEndpoints(t *testing.T) {
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
			wantErr: "NewCertOAuth2Client requires sandbox authorization URL",
		},
		{
			name: "production token URL with cert constructor",
			config: OAuth2Config{
				ClientID:     "test_client_id",
				ClientSecret: "test_client_secret",
				RedirectURI:  "http://localhost:8080/callback",
				TokenURL:     oauth2ProductionTokenURL,
			},
			wantErr: "NewCertOAuth2Client requires sandbox token URL",
		},
		{
			name: "production base URL with cert constructor",
			config: OAuth2Config{
				ClientID:     "test_client_id",
				ClientSecret: "test_client_secret",
				RedirectURI:  "http://localhost:8080/callback",
				BaseURL:      apiBaseURL,
			},
			wantErr: "use NewOAuth2Client for production environment",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewCertOAuth2Client(tt.config, nil)
			require.Error(t, err)
			require.Nil(t, client)
			require.Contains(t, err.Error(), tt.wantErr)
		})
	}
}

func TestClient_AuthMode(t *testing.T) {
	// Test session mode
	sessionClient := NewClient(nil)
	require.Equal(t, AuthModeSession, sessionClient.GetAuthMode())
	require.True(t, sessionClient.IsSessionMode())
	require.False(t, sessionClient.IsOAuth2Mode())

	// Test OAuth2 mode
	config := OAuth2Config{
		ClientID:     "test_client_id",
		ClientSecret: "test_client_secret",
		RedirectURI:  "http://localhost:8080/callback",
		Scopes:       []string{"read", "trade"},
	}
	oauth2Client, err := NewOAuth2Client(config, nil)
	require.NoError(t, err)
	require.Equal(t, AuthModeOAuth2, oauth2Client.GetAuthMode())
	require.True(t, oauth2Client.IsOAuth2Mode())
	require.False(t, oauth2Client.IsSessionMode())
}

func TestClient_OAuth2Methods_SessionMode(t *testing.T) {
	client := NewClient(nil)

	// Test OAuth2 methods fail in session mode
	_, err := client.GetAuthorizationURL()
	require.Error(t, err)
	require.Contains(t, err.Error(), "authorization URL is only available in OAuth2 mode")

	_, err = client.ExchangeCodeForTokens("test_code")
	require.Error(t, err)
	require.Contains(t, err.Error(), "token exchange is only available in OAuth2 mode")

	_, err = client.RefreshTokens()
	require.Error(t, err)
	require.Contains(t, err.Error(), "token refresh is only available in OAuth2 mode")

	_, err = client.StartRedirectServer(8080)
	require.Error(t, err)
	require.Contains(t, err.Error(), "redirect server is only available in OAuth2 mode")

	err = client.ValidateState("test_state")
	require.Error(t, err)
	require.Contains(t, err.Error(), "state validation is only available in OAuth2 mode")

	require.Nil(t, client.GetOAuth2Client())
}

func TestClient_OAuth2Methods_OAuth2Mode(t *testing.T) {
	config := OAuth2Config{
		ClientID:     "test_client_id",
		ClientSecret: "test_client_secret",
		RedirectURI:  "http://localhost:8080/callback",
		Scopes:       []string{"read", "trade"},
		State:        "test_state",
	}
	client, err := NewOAuth2Client(config, nil)
	require.NoError(t, err)

	// Test OAuth2 methods work in OAuth2 mode
	authURL, err := client.GetAuthorizationURL()
	require.NoError(t, err)
	require.Contains(t, authURL, "client_id=test_client_id")
	require.Contains(t, authURL, "state=test_state")

	err = client.ValidateState("test_state")
	require.NoError(t, err)

	err = client.ValidateState("wrong_state")
	require.Error(t, err)

	require.NotNil(t, client.GetOAuth2Client())
}

func TestClient_IsAuthenticated(t *testing.T) {
	// Test session mode
	sessionClient := NewClient(nil)
	require.False(t, sessionClient.IsAuthenticated())

	sessionClient.Session.SessionToken = &testToken
	require.True(t, sessionClient.IsAuthenticated())

	// Test OAuth2 mode with isolated storage
	config := OAuth2Config{
		ClientID:     "test_client_id",
		ClientSecret: "test_client_secret",
		RedirectURI:  "http://localhost:8080/callback",
		Scopes:       []string{"read", "trade"},
		BaseURL:      apiBaseURL,
		AuthURL:      oauth2ProductionAuthURL,
		TokenURL:     oauth2ProductionTokenURL,
	}

	// Create OAuth2Client with memory storage for isolated testing
	tokenManager := NewMemoryTokenManager()
	oauth2ClientInternal := &OAuth2Client{
		config:       config,
		tokenManager: tokenManager,
		httpClient:   defaultHTTPClient,
		pkce:         nil,
	}
	tokenManager.SetRefreshCallback(oauth2ClientInternal.refreshTokensInternal)

	// Create Client with the isolated OAuth2Client
	oauth2Client := &Client{
		httpClient:   defaultHTTPClient,
		baseURL:      config.BaseURL,
		oauth2Client: oauth2ClientInternal,
		authMode:     AuthModeOAuth2,
	}

	require.False(t, oauth2Client.IsAuthenticated())

	// Set tokens to make it authenticated
	oauth2Client.oauth2Client.tokenManager.SetTokens("access_token", "refresh_token", 3600)
	require.True(t, oauth2Client.IsAuthenticated())
}

func TestClient_ClearAuthentication(t *testing.T) {
	// Test session mode
	sessionClient := NewClient(nil)
	sessionClient.Session.SessionToken = &testToken
	require.True(t, sessionClient.IsAuthenticated())

	sessionClient.ClearAuthentication()
	require.False(t, sessionClient.IsAuthenticated())
	require.Equal(t, Session{}, sessionClient.Session)

	// Test OAuth2 mode
	config := OAuth2Config{
		ClientID:     "test_client_id",
		ClientSecret: "test_client_secret",
		RedirectURI:  "http://localhost:8080/callback",
		Scopes:       []string{"read", "trade"},
	}
	oauth2Client, err := NewOAuth2Client(config, nil)
	require.NoError(t, err)

	oauth2Client.oauth2Client.tokenManager.SetTokens("access_token", "refresh_token", 3600)
	require.True(t, oauth2Client.IsAuthenticated())

	oauth2Client.ClearAuthentication()
	require.False(t, oauth2Client.IsAuthenticated())
}

func TestAuthMode_String(t *testing.T) {
	require.Equal(t, "session", AuthModeSession.String())
	require.Equal(t, "oauth2", AuthModeOAuth2.String())
	require.Equal(t, "unknown", AuthMode(999).String())
}

func TestClient_BackwardCompatibility(t *testing.T) {
	// Test that existing session-based constructors still work
	sessionClient := NewClient(nil)
	require.NotNil(t, sessionClient)
	require.Equal(t, AuthModeSession, sessionClient.authMode)

	certClient := NewCertClient(nil)
	require.NotNil(t, certClient)
	require.Equal(t, AuthModeSession, certClient.authMode)

	// Test that session-based methods still work
	sessionClient.Session.SessionToken = &testToken
	require.True(t, sessionClient.IsAuthenticated())
}
