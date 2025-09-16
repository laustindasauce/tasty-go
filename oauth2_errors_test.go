package tasty

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewOAuth2Error tests OAuth2 error creation
func TestNewOAuth2Error(t *testing.T) {
	err := NewOAuth2Error("test_error", "Test error description")
	
	assert.Equal(t, "test_error", err.ErrorCode)
	assert.Equal(t, "Test error description", err.ErrorDescription)
	assert.NotZero(t, err.Timestamp)
	assert.NotEmpty(t, err.TroubleshootingGuide)
	assert.NotEmpty(t, err.SuggestedActions)
}

// TestNewOAuth2ErrorFromStandard tests creating detailed error from standard OAuth2Error
func TestNewOAuth2ErrorFromStandard(t *testing.T) {
	stdErr := OAuth2Error{
		ErrorCode:        "invalid_grant",
		ErrorDescription: "Invalid authorization code",
		ErrorURI:         "https://example.com/error",
		State:            "test_state",
	}
	
	detailedErr := NewOAuth2ErrorFromStandard(stdErr)
	
	assert.Equal(t, stdErr.ErrorCode, detailedErr.ErrorCode)
	assert.Equal(t, stdErr.ErrorDescription, detailedErr.ErrorDescription)
	assert.Equal(t, stdErr.ErrorURI, detailedErr.ErrorURI)
	assert.Equal(t, stdErr.State, detailedErr.State)
	assert.Equal(t, OAuth2ErrorTypeClient, detailedErr.Type)
	assert.Equal(t, OAuth2ErrorSeverityHigh, detailedErr.Severity)
}

// TestNewOAuth2ErrorWithContext tests creating error with additional context
func TestNewOAuth2ErrorWithContext(t *testing.T) {
	originalErr := errors.New("original error")
	
	detailedErr := NewOAuth2ErrorWithContext("network_error", "Network failure", 500, originalErr)
	
	assert.Equal(t, "network_error", detailedErr.ErrorCode)
	assert.Equal(t, "Network failure", detailedErr.ErrorDescription)
	assert.Equal(t, 500, detailedErr.HTTPStatusCode)
	assert.Equal(t, originalErr, detailedErr.OriginalError)
	assert.Equal(t, "original error", detailedErr.InternalMessage)
}

// TestOAuth2DetailedError_Error tests error message formatting
func TestOAuth2DetailedError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      OAuth2DetailedError
		expected string
	}{
		{
			name: "full error with all fields",
			err: OAuth2DetailedError{
				ErrorCode:        "invalid_request",
				ErrorDescription: "Missing parameter",
				HTTPStatusCode:   400,
			},
			expected: "OAuth2 error: invalid_request - Missing parameter - (HTTP 400)",
		},
		{
			name: "error with code and internal message",
			err: OAuth2DetailedError{
				ErrorCode:       "network_error",
				InternalMessage: "Connection timeout",
			},
			expected: "OAuth2 error: network_error - Connection timeout",
		},
		{
			name: "error with only code",
			err: OAuth2DetailedError{
				ErrorCode: "unknown_error",
			},
			expected: "OAuth2 error: unknown_error",
		},
		{
			name:     "empty error",
			err:      OAuth2DetailedError{},
			expected: "OAuth2 error occurred",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.err.Error())
		})
	}
}

// TestOAuth2DetailedError_IsRetryable tests retry logic
func TestOAuth2DetailedError_IsRetryable(t *testing.T) {
	tests := []struct {
		name      string
		err       OAuth2DetailedError
		retryable bool
	}{
		{
			name: "server error - retryable",
			err: OAuth2DetailedError{
				ErrorCode: OAuth2ErrorServerError,
				Type:      OAuth2ErrorTypeServer,
			},
			retryable: true,
		},
		{
			name: "network error - retryable",
			err: OAuth2DetailedError{
				ErrorCode: OAuth2ErrorNetworkError,
				Type:      OAuth2ErrorTypeNetwork,
			},
			retryable: true,
		},
		{
			name: "temporarily unavailable - retryable",
			err: OAuth2DetailedError{
				ErrorCode: OAuth2ErrorTemporarilyUnavailable,
			},
			retryable: true,
		},
		{
			name: "invalid client - not retryable",
			err: OAuth2DetailedError{
				ErrorCode: OAuth2ErrorInvalidClient,
			},
			retryable: false,
		},
		{
			name: "invalid grant - not retryable",
			err: OAuth2DetailedError{
				ErrorCode: OAuth2ErrorInvalidGrant,
			},
			retryable: false,
		},
		{
			name: "security error - not retryable",
			err: OAuth2DetailedError{
				ErrorCode: OAuth2ErrorInvalidState,
				Type:      OAuth2ErrorTypeSecurity,
			},
			retryable: false,
		},
		{
			name: "configuration error - not retryable",
			err: OAuth2DetailedError{
				ErrorCode: OAuth2ErrorConfigurationError,
				Type:      OAuth2ErrorTypeConfiguration,
			},
			retryable: false,
		},
		{
			name: "access denied - not retryable",
			err: OAuth2DetailedError{
				ErrorCode: OAuth2ErrorAccessDenied,
			},
			retryable: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.retryable, tt.err.IsRetryable())
		})
	}
}

// TestOAuth2DetailedError_GetSeverityString tests severity string conversion
func TestOAuth2DetailedError_GetSeverityString(t *testing.T) {
	tests := []struct {
		name     string
		severity OAuth2ErrorSeverity
		expected string
	}{
		{"Low", OAuth2ErrorSeverityLow, "Low"},
		{"Medium", OAuth2ErrorSeverityMedium, "Medium"},
		{"High", OAuth2ErrorSeverityHigh, "High"},
		{"Critical", OAuth2ErrorSeverityCritical, "Critical"},
		{"Unknown", OAuth2ErrorSeverity(999), "Unknown"},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := OAuth2DetailedError{Severity: tt.severity}
			assert.Equal(t, tt.expected, err.GetSeverityString())
		})
	}
}

// TestOAuth2DetailedError_GetTypeString tests type string conversion
func TestOAuth2DetailedError_GetTypeString(t *testing.T) {
	tests := []struct {
		name     string
		errType  OAuth2ErrorType
		expected string
	}{
		{"Client Error", OAuth2ErrorTypeClient, "Client Error"},
		{"Server Error", OAuth2ErrorTypeServer, "Server Error"},
		{"Network Error", OAuth2ErrorTypeNetwork, "Network Error"},
		{"Configuration Error", OAuth2ErrorTypeConfiguration, "Configuration Error"},
		{"Security Error", OAuth2ErrorTypeSecurity, "Security Error"},
		{"Token Error", OAuth2ErrorTypeToken, "Token Error"},
		{"Unknown Error", OAuth2ErrorType(999), "Unknown Error"},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := OAuth2DetailedError{Type: tt.errType}
			assert.Equal(t, tt.expected, err.GetTypeString())
		})
	}
}

// TestOAuth2DetailedError_Unwrap tests error unwrapping
func TestOAuth2DetailedError_Unwrap(t *testing.T) {
	originalErr := errors.New("original error")
	detailedErr := OAuth2DetailedError{OriginalError: originalErr}
	
	assert.Equal(t, originalErr, detailedErr.Unwrap())
	
	// Test with nil original error
	detailedErr2 := OAuth2DetailedError{}
	assert.Nil(t, detailedErr2.Unwrap())
}

// TestClassifyErrorType tests error type classification
func TestClassifyErrorType(t *testing.T) {
	tests := []struct {
		name      string
		errorCode string
		expected  OAuth2ErrorType
	}{
		{"invalid_request", OAuth2ErrorInvalidRequest, OAuth2ErrorTypeClient},
		{"invalid_client", OAuth2ErrorInvalidClient, OAuth2ErrorTypeClient},
		{"invalid_grant", OAuth2ErrorInvalidGrant, OAuth2ErrorTypeClient},
		{"server_error", OAuth2ErrorServerError, OAuth2ErrorTypeServer},
		{"temporarily_unavailable", OAuth2ErrorTemporarilyUnavailable, OAuth2ErrorTypeServer},
		{"network_error", OAuth2ErrorNetworkError, OAuth2ErrorTypeNetwork},
		{"configuration_error", OAuth2ErrorConfigurationError, OAuth2ErrorTypeConfiguration},
		{"invalid_state", OAuth2ErrorInvalidState, OAuth2ErrorTypeSecurity},
		{"pkce_error", OAuth2ErrorPKCEError, OAuth2ErrorTypeSecurity},
		{"token_expired", OAuth2ErrorTokenExpired, OAuth2ErrorTypeToken},
		{"refresh_failed", OAuth2ErrorRefreshFailed, OAuth2ErrorTypeToken},
		{"access_denied", OAuth2ErrorAccessDenied, OAuth2ErrorTypeSecurity},
		{"unknown_error", "unknown_error", OAuth2ErrorTypeUnknown},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := classifyErrorType(tt.errorCode)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestClassifyErrorSeverity tests error severity classification
func TestClassifyErrorSeverity(t *testing.T) {
	tests := []struct {
		name      string
		errorCode string
		expected  OAuth2ErrorSeverity
	}{
		{"invalid_state", OAuth2ErrorInvalidState, OAuth2ErrorSeverityCritical},
		{"pkce_error", OAuth2ErrorPKCEError, OAuth2ErrorSeverityCritical},
		{"invalid_client", OAuth2ErrorInvalidClient, OAuth2ErrorSeverityHigh},
		{"invalid_grant", OAuth2ErrorInvalidGrant, OAuth2ErrorSeverityHigh},
		{"access_denied", OAuth2ErrorAccessDenied, OAuth2ErrorSeverityMedium},
		{"token_expired", OAuth2ErrorTokenExpired, OAuth2ErrorSeverityMedium},
		{"server_error", OAuth2ErrorServerError, OAuth2ErrorSeverityMedium},
		{"unknown_error", "unknown_error", OAuth2ErrorSeverityLow},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := classifyErrorSeverity(tt.errorCode)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestGenerateTroubleshootingGuide tests troubleshooting guide generation
func TestGenerateTroubleshootingGuide(t *testing.T) {
	guide := generateTroubleshootingGuide(OAuth2ErrorInvalidRequest)
	assert.NotEmpty(t, guide)
	assert.Contains(t, guide, "missing a required parameter")
	
	// Test unknown error code
	unknownGuide := generateTroubleshootingGuide("unknown_error")
	assert.NotEmpty(t, unknownGuide)
	assert.Contains(t, unknownGuide, "OAuth2 error occurred")
}

// TestGenerateSuggestedActions tests suggested actions generation
func TestGenerateSuggestedActions(t *testing.T) {
	actions := generateSuggestedActions(OAuth2ErrorInvalidRequest)
	assert.NotEmpty(t, actions)
	assert.Greater(t, len(actions), 0)
	
	// Test unknown error code
	unknownActions := generateSuggestedActions("unknown_error")
	assert.NotEmpty(t, unknownActions)
	assert.Greater(t, len(unknownActions), 0)
}

// TestOAuth2ErrorHandler_NewOAuth2ErrorHandler tests error handler creation
func TestOAuth2ErrorHandler_NewOAuth2ErrorHandler(t *testing.T) {
	handler := NewOAuth2ErrorHandler()
	
	assert.True(t, handler.EnableRetry)
	assert.Equal(t, 3, handler.MaxRetries)
	assert.Equal(t, time.Second, handler.RetryDelay)
	assert.True(t, handler.EnableLogging)
	assert.False(t, handler.LogSensitiveData)
}

// TestOAuth2ErrorHandler_HandleError tests error handling
func TestOAuth2ErrorHandler_HandleError(t *testing.T) {
	handler := NewOAuth2ErrorHandler()
	
	// Test with nil error
	detailedErr, shouldRetry := handler.HandleError(nil)
	assert.Nil(t, detailedErr)
	assert.False(t, shouldRetry)
	
	// Test with detailed OAuth2 error
	originalDetailedErr := NewOAuth2Error(OAuth2ErrorServerError, "Server error")
	detailedErr, shouldRetry = handler.HandleError(originalDetailedErr)
	assert.Equal(t, originalDetailedErr, detailedErr)
	assert.True(t, shouldRetry) // Server errors are retryable
	
	// Test with standard OAuth2 error
	stdErr := OAuth2Error{
		ErrorCode:        "invalid_grant",
		ErrorDescription: "Invalid code",
	}
	detailedErr, shouldRetry = handler.HandleError(stdErr)
	assert.Equal(t, "invalid_grant", detailedErr.ErrorCode)
	assert.False(t, shouldRetry) // Invalid grant is not retryable
	
	// Test with generic error
	genericErr := errors.New("connection timeout")
	detailedErr, shouldRetry = handler.HandleError(genericErr)
	assert.Equal(t, OAuth2ErrorNetworkError, detailedErr.ErrorCode)
	assert.True(t, shouldRetry) // Network errors are retryable
}

// TestOAuth2ErrorHandler_ClassifyGenericError tests generic error classification
func TestOAuth2ErrorHandler_ClassifyGenericError(t *testing.T) {
	handler := NewOAuth2ErrorHandler()
	
	tests := []struct {
		name          string
		err           error
		expectedCode  string
		expectedType  OAuth2ErrorType
	}{
		{
			name:         "network error",
			err:          errors.New("connection refused"),
			expectedCode: OAuth2ErrorNetworkError,
			expectedType: OAuth2ErrorTypeNetwork,
		},
		{
			name:         "timeout error",
			err:          errors.New("timeout occurred"),
			expectedCode: OAuth2ErrorNetworkError,
			expectedType: OAuth2ErrorTypeNetwork,
		},
		{
			name:         "http error",
			err:          errors.New("http status 500"),
			expectedCode: OAuth2ErrorServerError,
			expectedType: OAuth2ErrorTypeServer,
		},
		{
			name:         "configuration error",
			err:          errors.New("invalid config parameter"),
			expectedCode: OAuth2ErrorConfigurationError,
			expectedType: OAuth2ErrorTypeConfiguration,
		},
		{
			name:         "unknown error",
			err:          errors.New("something went wrong"),
			expectedCode: "unknown_error",
			expectedType: OAuth2ErrorTypeUnknown,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			detailedErr := handler.classifyGenericError(tt.err)
			assert.Equal(t, tt.expectedCode, detailedErr.ErrorCode)
			assert.Equal(t, tt.expectedType, detailedErr.Type)
			assert.Equal(t, tt.err, detailedErr.OriginalError)
		})
	}
}

// TestValidateState tests state parameter validation
func TestValidateState(t *testing.T) {
	tests := []struct {
		name     string
		expected string
		received string
		wantErr  bool
		errCode  string
	}{
		{
			name:     "valid state",
			expected: "test_state_123",
			received: "test_state_123",
			wantErr:  false,
		},
		{
			name:     "empty expected state",
			expected: "",
			received: "any_state",
			wantErr:  true,
			errCode:  OAuth2ErrorConfigurationError,
		},
		{
			name:     "empty received state",
			expected: "test_state",
			received: "",
			wantErr:  true,
			errCode:  OAuth2ErrorInvalidState,
		},
		{
			name:     "mismatched state",
			expected: "expected_state",
			received: "received_state",
			wantErr:  true,
			errCode:  OAuth2ErrorInvalidState,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateState(tt.expected, tt.received)
			
			if tt.wantErr {
				require.Error(t, err)
				detailedErr, ok := err.(*OAuth2DetailedError)
				require.True(t, ok)
				assert.Equal(t, tt.errCode, detailedErr.ErrorCode)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestValidateAuthorizationCode tests authorization code validation
func TestValidateAuthorizationCode(t *testing.T) {
	tests := []struct {
		name    string
		code    string
		wantErr bool
		errCode string
	}{
		{
			name:    "valid code",
			code:    "valid_authorization_code_123",
			wantErr: false,
		},
		{
			name:    "empty code",
			code:    "",
			wantErr: true,
			errCode: OAuth2ErrorMissingCode,
		},
		{
			name:    "too short code",
			code:    "short",
			wantErr: true,
			errCode: OAuth2ErrorInvalidGrant,
		},
		{
			name:    "code with invalid characters",
			code:    "invalid@code#with$special%chars",
			wantErr: true,
			errCode: OAuth2ErrorInvalidGrant,
		},
		{
			name:    "valid code with allowed characters",
			code:    "valid-code_123.abc",
			wantErr: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateAuthorizationCode(tt.code)
			
			if tt.wantErr {
				require.Error(t, err)
				detailedErr, ok := err.(*OAuth2DetailedError)
				require.True(t, ok)
				assert.Equal(t, tt.errCode, detailedErr.ErrorCode)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestValidateTokenResponse tests token response validation
func TestValidateTokenResponse(t *testing.T) {
	tests := []struct {
		name     string
		response *TokenResponse
		wantErr  bool
		errCode  string
	}{
		{
			name: "valid response",
			response: &TokenResponse{
				AccessToken:  "valid_access_token",
				RefreshToken: "valid_refresh_token",
				TokenType:    "Bearer",
				ExpiresIn:    3600,
			},
			wantErr: false,
		},
		{
			name:     "nil response",
			response: nil,
			wantErr:  true,
			errCode:  OAuth2ErrorServerError,
		},
		{
			name: "missing access token",
			response: &TokenResponse{
				RefreshToken: "valid_refresh_token",
				TokenType:    "Bearer",
				ExpiresIn:    3600,
			},
			wantErr: true,
			errCode: OAuth2ErrorServerError,
		},
		{
			name: "missing token type",
			response: &TokenResponse{
				AccessToken:  "valid_access_token",
				RefreshToken: "valid_refresh_token",
				ExpiresIn:    3600,
			},
			wantErr: true,
			errCode: OAuth2ErrorServerError,
		},
		{
			name: "invalid expires in",
			response: &TokenResponse{
				AccessToken:  "valid_access_token",
				RefreshToken: "valid_refresh_token",
				TokenType:    "Bearer",
				ExpiresIn:    0,
			},
			wantErr: true,
			errCode: OAuth2ErrorServerError,
		},
		{
			name: "too short access token",
			response: &TokenResponse{
				AccessToken:  "short",
				RefreshToken: "valid_refresh_token",
				TokenType:    "Bearer",
				ExpiresIn:    3600,
			},
			wantErr: true,
			errCode: OAuth2ErrorServerError,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTokenResponse(tt.response)
			
			if tt.wantErr {
				require.Error(t, err)
				detailedErr, ok := err.(*OAuth2DetailedError)
				require.True(t, ok)
				assert.Equal(t, tt.errCode, detailedErr.ErrorCode)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestIsTemporaryError tests temporary error detection
func TestIsTemporaryError(t *testing.T) {
	tests := []struct {
		name      string
		err       error
		temporary bool
	}{
		{
			name:      "nil error",
			err:       nil,
			temporary: false,
		},
		{
			name: "retryable detailed error",
			err: &OAuth2DetailedError{
				ErrorCode: OAuth2ErrorServerError,
				Type:      OAuth2ErrorTypeServer,
			},
			temporary: true,
		},
		{
			name: "non-retryable detailed error",
			err: &OAuth2DetailedError{
				ErrorCode: OAuth2ErrorInvalidClient,
				Type:      OAuth2ErrorTypeClient,
			},
			temporary: false,
		},
		{
			name:      "timeout error",
			err:       errors.New("timeout occurred"),
			temporary: true,
		},
		{
			name:      "connection error",
			err:       errors.New("connection refused"),
			temporary: true,
		},
		{
			name:      "503 error",
			err:       errors.New("HTTP 503 Service Unavailable"),
			temporary: true,
		},
		{
			name:      "generic error",
			err:       errors.New("something went wrong"),
			temporary: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsTemporaryError(tt.err)
			assert.Equal(t, tt.temporary, result)
		})
	}
}

// TestWrapHTTPError tests HTTP error wrapping
func TestWrapHTTPError(t *testing.T) {
	tests := []struct {
		name           string
		statusCode     int
		originalErr    error
		expectedCode   string
		expectedType   OAuth2ErrorType
		expectedSeverity OAuth2ErrorSeverity
	}{
		{
			name:             "400 Bad Request",
			statusCode:       http.StatusBadRequest,
			originalErr:      errors.New("bad request"),
			expectedCode:     OAuth2ErrorInvalidRequest,
			expectedType:     OAuth2ErrorTypeClient,
			expectedSeverity: OAuth2ErrorSeverityMedium,
		},
		{
			name:             "401 Unauthorized",
			statusCode:       http.StatusUnauthorized,
			originalErr:      errors.New("unauthorized"),
			expectedCode:     OAuth2ErrorInvalidClient,
			expectedType:     OAuth2ErrorTypeClient,
			expectedSeverity: OAuth2ErrorSeverityHigh,
		},
		{
			name:             "403 Forbidden",
			statusCode:       http.StatusForbidden,
			originalErr:      errors.New("forbidden"),
			expectedCode:     OAuth2ErrorAccessDenied,
			expectedType:     OAuth2ErrorTypeClient,
			expectedSeverity: OAuth2ErrorSeverityMedium,
		},
		{
			name:             "500 Internal Server Error",
			statusCode:       http.StatusInternalServerError,
			originalErr:      errors.New("server error"),
			expectedCode:     OAuth2ErrorServerError,
			expectedType:     OAuth2ErrorTypeServer,
			expectedSeverity: OAuth2ErrorSeverityMedium,
		},
		{
			name:             "503 Service Unavailable",
			statusCode:       http.StatusServiceUnavailable,
			originalErr:      errors.New("service unavailable"),
			expectedCode:     OAuth2ErrorTemporarilyUnavailable,
			expectedType:     OAuth2ErrorTypeServer,
			expectedSeverity: OAuth2ErrorSeverityLow,
		},
		{
			name:             "418 I'm a teapot (unknown status)",
			statusCode:       http.StatusTeapot,
			originalErr:      errors.New("teapot error"),
			expectedCode:     "http_error",
			expectedType:     OAuth2ErrorTypeUnknown,
			expectedSeverity: OAuth2ErrorSeverityLow,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := &http.Response{
				StatusCode: tt.statusCode,
			}
			
			detailedErr := WrapHTTPError(resp, tt.originalErr)
			
			assert.Equal(t, tt.expectedCode, detailedErr.ErrorCode)
			assert.Equal(t, tt.expectedType, detailedErr.Type)
			assert.Equal(t, tt.expectedSeverity, detailedErr.Severity)
			assert.Equal(t, tt.statusCode, detailedErr.HTTPStatusCode)
			assert.Equal(t, tt.originalErr, detailedErr.OriginalError)
			assert.Contains(t, detailedErr.ErrorDescription, http.StatusText(tt.statusCode))
		})
	}
}

// TestWrapHTTPError_NilResponse tests HTTP error wrapping with nil response
func TestWrapHTTPError_NilResponse(t *testing.T) {
	originalErr := errors.New("network error")
	detailedErr := WrapHTTPError(nil, originalErr)
	
	assert.Equal(t, OAuth2ErrorNetworkError, detailedErr.ErrorCode)
	assert.Equal(t, originalErr, detailedErr.OriginalError)
	assert.Contains(t, detailedErr.ErrorDescription, "HTTP response is nil")
}

// TestPredefinedErrors tests predefined error variables
func TestPredefinedErrors(t *testing.T) {
	tests := []struct {
		name         string
		err          *OAuth2DetailedError
		expectedCode string
	}{
		{"ErrOAuth2InvalidState", ErrOAuth2InvalidState, OAuth2ErrorInvalidState},
		{"ErrOAuth2TokenExpired", ErrOAuth2TokenExpired, OAuth2ErrorTokenExpired},
		{"ErrOAuth2RefreshFailed", ErrOAuth2RefreshFailed, OAuth2ErrorRefreshFailed},
		{"ErrOAuth2InvalidCode", ErrOAuth2InvalidCode, OAuth2ErrorMissingCode},
		{"ErrOAuth2RedirectTimeout", ErrOAuth2RedirectTimeout, OAuth2ErrorRedirectTimeout},
		{"ErrOAuth2InvalidCredentials", ErrOAuth2InvalidCredentials, OAuth2ErrorInvalidCredentials},
		{"ErrOAuth2ConfigurationError", ErrOAuth2ConfigurationError, OAuth2ErrorConfigurationError},
		{"ErrOAuth2NetworkError", ErrOAuth2NetworkError, OAuth2ErrorNetworkError},
		{"ErrOAuth2PKCEError", ErrOAuth2PKCEError, OAuth2ErrorPKCEError},
		{"ErrOAuth2ServerUnavailable", ErrOAuth2ServerUnavailable, OAuth2ErrorServerUnavailable},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expectedCode, tt.err.ErrorCode)
			assert.NotEmpty(t, tt.err.ErrorDescription)
			assert.NotEmpty(t, tt.err.TroubleshootingGuide)
			assert.NotEmpty(t, tt.err.SuggestedActions)
		})
	}
}

// TestOAuth2ErrorIntegration tests error handling integration
func TestOAuth2ErrorIntegration(t *testing.T) {
	// Create a server that returns various error types
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/server_error":
			w.WriteHeader(http.StatusInternalServerError)
		case "/client_error":
			w.WriteHeader(http.StatusBadRequest)
		case "/oauth_error":
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"error":"invalid_grant","error_description":"Invalid code"}`))
		default:
			w.WriteHeader(http.StatusOK)
		}
	}))
	defer server.Close()
	
	handler := NewOAuth2ErrorHandler()
	
	// Test server error (retryable)
	resp, _ := http.Get(server.URL + "/server_error")
	serverErr := WrapHTTPError(resp, nil)
	detailedErr, shouldRetry := handler.HandleError(serverErr)
	assert.True(t, shouldRetry)
	assert.Equal(t, OAuth2ErrorServerError, detailedErr.ErrorCode)
	
	// Test client error (not retryable)
	resp, _ = http.Get(server.URL + "/client_error")
	clientErr := WrapHTTPError(resp, nil)
	detailedErr, shouldRetry = handler.HandleError(clientErr)
	assert.False(t, shouldRetry)
	assert.Equal(t, OAuth2ErrorInvalidRequest, detailedErr.ErrorCode)
}

// TestOAuth2ErrorConcurrency tests error handling under concurrent access
func TestOAuth2ErrorConcurrency(t *testing.T) {
	handler := NewOAuth2ErrorHandler()
	
	const numGoroutines = 50
	const numOperations = 100
	
	var wg sync.WaitGroup
	errorsChan := make(chan error, numGoroutines*numOperations)
	
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				// Create various types of errors
				testErr := NewOAuth2Error("test_error", "Test error")
				_, _ = handler.HandleError(testErr)
				
				genericErr := errors.New("generic error")
				_, _ = handler.HandleError(genericErr)
				
				// Test error classification
				_ = classifyErrorType("test_code")
				_ = classifyErrorSeverity("test_code")
			}
		}(i)
	}
	
	wg.Wait()
	close(errorsChan)
	
	// Check for any errors
	for err := range errorsChan {
		t.Errorf("Concurrent error handling error: %v", err)
	}
}