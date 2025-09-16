package tasty

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

// OAuth2 error constants as defined by RFC 6749
const (
	// Authorization endpoint errors
	OAuth2ErrorInvalidRequest          = "invalid_request"
	OAuth2ErrorUnauthorizedClient      = "unauthorized_client"
	OAuth2ErrorAccessDenied            = "access_denied"
	OAuth2ErrorUnsupportedResponseType = "unsupported_response_type"
	OAuth2ErrorInvalidScope            = "invalid_scope"
	OAuth2ErrorServerError             = "server_error"
	OAuth2ErrorTemporarilyUnavailable  = "temporarily_unavailable"

	// Token endpoint errors
	OAuth2ErrorInvalidClient           = "invalid_client"
	OAuth2ErrorInvalidGrant            = "invalid_grant"
	OAuth2ErrorUnsupportedGrantType    = "unsupported_grant_type"

	// Custom error codes for internal use
	OAuth2ErrorTokenExpired            = "token_expired"
	OAuth2ErrorRefreshFailed           = "refresh_failed"
	OAuth2ErrorInvalidState            = "invalid_state"
	OAuth2ErrorMissingCode             = "missing_code"
	OAuth2ErrorRedirectTimeout         = "redirect_timeout"
	OAuth2ErrorInvalidCredentials      = "invalid_credentials"
	OAuth2ErrorConfigurationError      = "configuration_error"
	OAuth2ErrorNetworkError            = "network_error"
	OAuth2ErrorPKCEError               = "pkce_error"
	OAuth2ErrorServerUnavailable       = "server_unavailable"
)

// OAuth2ErrorType represents different categories of OAuth2 errors
type OAuth2ErrorType int

const (
	OAuth2ErrorTypeUnknown OAuth2ErrorType = iota
	OAuth2ErrorTypeClient                   // Client-side errors (4xx)
	OAuth2ErrorTypeServer                   // Server-side errors (5xx)
	OAuth2ErrorTypeNetwork                  // Network connectivity errors
	OAuth2ErrorTypeConfiguration            // Configuration/setup errors
	OAuth2ErrorTypeSecurity                 // Security-related errors (CSRF, invalid state)
	OAuth2ErrorTypeToken                    // Token-related errors (expired, invalid)
)

// OAuth2ErrorSeverity indicates the severity level of an OAuth2 error
type OAuth2ErrorSeverity int

const (
	OAuth2ErrorSeverityLow OAuth2ErrorSeverity = iota
	OAuth2ErrorSeverityMedium
	OAuth2ErrorSeverityHigh
	OAuth2ErrorSeverityCritical
)

// OAuth2DetailedError provides comprehensive error information for OAuth2 operations
type OAuth2DetailedError struct {
	// Standard OAuth2 error fields
	ErrorCode        string `json:"error"`
	ErrorDescription string `json:"error_description,omitempty"`
	ErrorURI         string `json:"error_uri,omitempty"`
	State            string `json:"state,omitempty"`

	// Additional context fields
	Type             OAuth2ErrorType     `json:"type"`
	Severity         OAuth2ErrorSeverity `json:"severity"`
	HTTPStatusCode   int                 `json:"http_status_code,omitempty"`
	Timestamp        time.Time           `json:"timestamp"`
	RequestID        string              `json:"request_id,omitempty"`
	
	// Troubleshooting information
	TroubleshootingGuide string   `json:"troubleshooting_guide,omitempty"`
	SuggestedActions     []string `json:"suggested_actions,omitempty"`
	
	// Internal context
	InternalMessage string `json:"internal_message,omitempty"`
	OriginalError   error  `json:"-"` // Don't serialize the original error
}

// Error implements the error interface
func (e OAuth2DetailedError) Error() string {
	var parts []string
	
	if e.ErrorCode != "" {
		parts = append(parts, fmt.Sprintf("OAuth2 error: %s", e.ErrorCode))
	}
	
	if e.ErrorDescription != "" {
		parts = append(parts, e.ErrorDescription)
	} else if e.InternalMessage != "" {
		parts = append(parts, e.InternalMessage)
	}
	
	if e.HTTPStatusCode > 0 {
		parts = append(parts, fmt.Sprintf("(HTTP %d)", e.HTTPStatusCode))
	}
	
	if len(parts) == 0 {
		return "OAuth2 error occurred"
	}
	
	return strings.Join(parts, " - ")
}

// IsRetryable determines if the error is retryable based on error type and code
func (e OAuth2DetailedError) IsRetryable() bool {
	// Never retry security-related errors
	if e.Type == OAuth2ErrorTypeSecurity {
		return false
	}
	
	// Never retry client configuration errors
	if e.Type == OAuth2ErrorTypeConfiguration {
		return false
	}
	
	// Check specific error codes that should not be retried
	nonRetryableErrors := []string{
		OAuth2ErrorInvalidClient,
		OAuth2ErrorInvalidGrant,
		OAuth2ErrorUnauthorizedClient,
		OAuth2ErrorUnsupportedGrantType,
		OAuth2ErrorInvalidScope,
		OAuth2ErrorAccessDenied,
		OAuth2ErrorUnsupportedResponseType,
		OAuth2ErrorInvalidState,
		OAuth2ErrorInvalidCredentials,
	}
	
	for _, nonRetryable := range nonRetryableErrors {
		if e.ErrorCode == nonRetryable {
			return false
		}
	}
	
	// Retry server errors and network errors
	return e.Type == OAuth2ErrorTypeServer || e.Type == OAuth2ErrorTypeNetwork ||
		   e.ErrorCode == OAuth2ErrorServerError || e.ErrorCode == OAuth2ErrorTemporarilyUnavailable
}

// GetSeverityString returns a human-readable severity level
func (e OAuth2DetailedError) GetSeverityString() string {
	switch e.Severity {
	case OAuth2ErrorSeverityLow:
		return "Low"
	case OAuth2ErrorSeverityMedium:
		return "Medium"
	case OAuth2ErrorSeverityHigh:
		return "High"
	case OAuth2ErrorSeverityCritical:
		return "Critical"
	default:
		return "Unknown"
	}
}

// GetTypeString returns a human-readable error type
func (e OAuth2DetailedError) GetTypeString() string {
	switch e.Type {
	case OAuth2ErrorTypeClient:
		return "Client Error"
	case OAuth2ErrorTypeServer:
		return "Server Error"
	case OAuth2ErrorTypeNetwork:
		return "Network Error"
	case OAuth2ErrorTypeConfiguration:
		return "Configuration Error"
	case OAuth2ErrorTypeSecurity:
		return "Security Error"
	case OAuth2ErrorTypeToken:
		return "Token Error"
	default:
		return "Unknown Error"
	}
}

// Unwrap returns the original error for error unwrapping
func (e OAuth2DetailedError) Unwrap() error {
	return e.OriginalError
}

// NewOAuth2Error creates a new OAuth2DetailedError with basic information
func NewOAuth2Error(code, description string) *OAuth2DetailedError {
	detailedErr := &OAuth2DetailedError{
		ErrorCode:        code,
		ErrorDescription: description,
		Type:             classifyErrorType(code),
		Severity:         classifyErrorSeverity(code),
		Timestamp:        time.Now(),
	}
	
	detailedErr.TroubleshootingGuide = generateTroubleshootingGuide(code)
	detailedErr.SuggestedActions = generateSuggestedActions(code)
	
	return detailedErr
}

// NewOAuth2ErrorFromStandard creates a detailed error from a standard OAuth2Error
func NewOAuth2ErrorFromStandard(stdErr OAuth2Error) *OAuth2DetailedError {
	detailedErr := &OAuth2DetailedError{
		ErrorCode:        stdErr.ErrorCode,
		ErrorDescription: stdErr.ErrorDescription,
		ErrorURI:         stdErr.ErrorURI,
		State:            stdErr.State,
		Type:             classifyErrorType(stdErr.ErrorCode),
		Severity:         classifyErrorSeverity(stdErr.ErrorCode),
		Timestamp:        time.Now(),
	}
	
	detailedErr.TroubleshootingGuide = generateTroubleshootingGuide(detailedErr.ErrorCode)
	detailedErr.SuggestedActions = generateSuggestedActions(detailedErr.ErrorCode)
	
	return detailedErr
}

// NewOAuth2ErrorWithContext creates a detailed error with additional context
func NewOAuth2ErrorWithContext(code, description string, httpStatus int, originalErr error) *OAuth2DetailedError {
	detailedErr := &OAuth2DetailedError{
		ErrorCode:        code,
		ErrorDescription: description,
		HTTPStatusCode:   httpStatus,
		Type:             classifyErrorType(code),
		Severity:         classifyErrorSeverity(code),
		Timestamp:        time.Now(),
		OriginalError:    originalErr,
	}
	
	// Add internal message if original error exists
	if originalErr != nil {
		detailedErr.InternalMessage = originalErr.Error()
	}
	
	detailedErr.TroubleshootingGuide = generateTroubleshootingGuide(code)
	detailedErr.SuggestedActions = generateSuggestedActions(code)
	
	return detailedErr
}

// classifyErrorType determines the error type based on the error code
func classifyErrorType(errorCode string) OAuth2ErrorType {
	switch errorCode {
	case OAuth2ErrorInvalidRequest, OAuth2ErrorUnauthorizedClient, OAuth2ErrorInvalidClient,
		 OAuth2ErrorInvalidGrant, OAuth2ErrorUnsupportedGrantType, OAuth2ErrorInvalidScope,
		 OAuth2ErrorUnsupportedResponseType:
		return OAuth2ErrorTypeClient
		
	case OAuth2ErrorServerError, OAuth2ErrorTemporarilyUnavailable:
		return OAuth2ErrorTypeServer
		
	case OAuth2ErrorNetworkError, OAuth2ErrorServerUnavailable:
		return OAuth2ErrorTypeNetwork
		
	case OAuth2ErrorConfigurationError, OAuth2ErrorInvalidCredentials:
		return OAuth2ErrorTypeConfiguration
		
	case OAuth2ErrorInvalidState, OAuth2ErrorPKCEError:
		return OAuth2ErrorTypeSecurity
		
	case OAuth2ErrorTokenExpired, OAuth2ErrorRefreshFailed:
		return OAuth2ErrorTypeToken
		
	case OAuth2ErrorAccessDenied:
		return OAuth2ErrorTypeSecurity
		
	default:
		return OAuth2ErrorTypeUnknown
	}
}

// classifyErrorSeverity determines the error severity based on the error code
func classifyErrorSeverity(errorCode string) OAuth2ErrorSeverity {
	switch errorCode {
	case OAuth2ErrorInvalidState, OAuth2ErrorPKCEError:
		return OAuth2ErrorSeverityCritical
		
	case OAuth2ErrorInvalidClient, OAuth2ErrorInvalidGrant, OAuth2ErrorUnauthorizedClient,
		 OAuth2ErrorInvalidCredentials:
		return OAuth2ErrorSeverityHigh
		
	case OAuth2ErrorAccessDenied, OAuth2ErrorTokenExpired, OAuth2ErrorRefreshFailed:
		return OAuth2ErrorSeverityMedium
		
	case OAuth2ErrorServerError, OAuth2ErrorTemporarilyUnavailable, OAuth2ErrorNetworkError:
		return OAuth2ErrorSeverityMedium
		
	default:
		return OAuth2ErrorSeverityLow
	}
}

// generateTroubleshootingGuide provides detailed troubleshooting information
func generateTroubleshootingGuide(errorCode string) string {
	guides := map[string]string{
		OAuth2ErrorInvalidRequest: "The request is missing a required parameter, includes an invalid parameter value, includes a parameter more than once, or is otherwise malformed. Check your authorization URL parameters and ensure all required fields are present and properly formatted.",
		
		OAuth2ErrorUnauthorizedClient: "The client is not authorized to request an authorization code using this method. Verify that your client ID is correct and that your application is properly registered with TastyTrade.",
		
		OAuth2ErrorAccessDenied: "The resource owner or authorization server denied the request. This typically occurs when the user cancels the authorization process. You may want to retry the authorization flow or check if the user has the necessary permissions.",
		
		OAuth2ErrorUnsupportedResponseType: "The authorization server does not support obtaining an authorization code using this method. Ensure you're using 'code' as the response_type parameter.",
		
		OAuth2ErrorInvalidScope: "The requested scope is invalid, unknown, or malformed. Check that you're requesting valid scopes for the TastyTrade API (typically 'read' and/or 'trade').",
		
		OAuth2ErrorServerError: "The authorization server encountered an unexpected condition that prevented it from fulfilling the request. This is typically a temporary issue - try again after a short delay.",
		
		OAuth2ErrorTemporarilyUnavailable: "The authorization server is currently unable to handle the request due to a temporary overloading or maintenance. Try again after a short delay.",
		
		OAuth2ErrorInvalidClient: "Client authentication failed (e.g., unknown client, no client authentication included, or unsupported authentication method). Verify your client ID and client secret are correct.",
		
		OAuth2ErrorInvalidGrant: "The provided authorization grant (e.g., authorization code, resource owner credentials) or refresh token is invalid, expired, revoked, does not match the redirection URI used in the authorization request, or was issued to another client. You may need to restart the authorization flow.",
		
		OAuth2ErrorUnsupportedGrantType: "The authorization grant type is not supported by the authorization server. Ensure you're using 'authorization_code' for initial token exchange or 'refresh_token' for token refresh.",
		
		OAuth2ErrorTokenExpired: "The access token has expired and needs to be refreshed. Use the refresh token to obtain a new access token, or restart the authorization flow if no refresh token is available.",
		
		OAuth2ErrorRefreshFailed: "Failed to refresh the access token using the refresh token. The refresh token may be expired or invalid. You may need to restart the authorization flow.",
		
		OAuth2ErrorInvalidState: "The state parameter in the callback does not match the expected value. This could indicate a CSRF attack or a problem with your application's state management. Ensure you're properly generating and validating state parameters.",
		
		OAuth2ErrorMissingCode: "The authorization callback is missing the required authorization code parameter. This may indicate an error in the authorization flow or a malformed callback URL.",
		
		OAuth2ErrorRedirectTimeout: "Timed out waiting for the authorization callback. The user may not have completed the authorization process, or there may be an issue with the redirect server. Try increasing the timeout or check your redirect URI configuration.",
		
		OAuth2ErrorInvalidCredentials: "The provided OAuth2 credentials (client ID, client secret, or redirect URI) are invalid or not properly configured. Verify your credentials with TastyTrade and ensure they match your application registration.",
		
		OAuth2ErrorConfigurationError: "There is an error in your OAuth2 configuration. Check that all required fields are properly set and that URLs are valid and accessible.",
		
		OAuth2ErrorNetworkError: "A network error occurred while communicating with the OAuth2 server. Check your internet connection and ensure that TastyTrade's OAuth2 endpoints are accessible.",
		
		OAuth2ErrorPKCEError: "An error occurred with PKCE (Proof Key for Code Exchange) validation. This is a security feature that helps protect against authorization code interception attacks. Ensure your PKCE implementation is correct.",
		
		OAuth2ErrorServerUnavailable: "The OAuth2 server is currently unavailable. This may be due to maintenance or high load. Try again after a short delay.",
	}
	
	if guide, exists := guides[errorCode]; exists {
		return guide
	}
	
	return "An OAuth2 error occurred. Please check your configuration and try again. If the problem persists, contact TastyTrade support."
}

// generateSuggestedActions provides actionable steps to resolve the error
func generateSuggestedActions(errorCode string) []string {
	actions := map[string][]string{
		OAuth2ErrorInvalidRequest: {
			"Verify all required parameters are included in the authorization URL",
			"Check parameter formatting and encoding",
			"Ensure no parameters are duplicated",
			"Review the OAuth2 specification for required parameters",
		},
		
		OAuth2ErrorUnauthorizedClient: {
			"Verify your client ID is correct",
			"Check that your application is registered with TastyTrade",
			"Ensure your application has the necessary permissions",
			"Contact TastyTrade support if the issue persists",
		},
		
		OAuth2ErrorAccessDenied: {
			"Retry the authorization flow",
			"Check if the user has necessary permissions",
			"Verify the requested scopes are appropriate",
			"Ensure the user completes the authorization process",
		},
		
		OAuth2ErrorUnsupportedResponseType: {
			"Use 'code' as the response_type parameter",
			"Check the OAuth2 specification for supported response types",
			"Verify your authorization URL is properly formatted",
		},
		
		OAuth2ErrorInvalidScope: {
			"Use valid scopes: 'read', 'trade', or 'read trade'",
			"Check TastyTrade API documentation for available scopes",
			"Remove any invalid or unsupported scopes from your request",
		},
		
		OAuth2ErrorServerError: {
			"Wait a few minutes and try again",
			"Check TastyTrade status page for known issues",
			"Implement exponential backoff for retries",
			"Contact TastyTrade support if the issue persists",
		},
		
		OAuth2ErrorTemporarilyUnavailable: {
			"Wait and retry after a short delay",
			"Implement exponential backoff",
			"Check TastyTrade status page",
			"Try again during off-peak hours",
		},
		
		OAuth2ErrorInvalidClient: {
			"Verify your client ID and client secret",
			"Check that credentials match your TastyTrade application registration",
			"Ensure you're using the correct environment (production vs sandbox)",
			"Regenerate credentials if necessary",
		},
		
		OAuth2ErrorInvalidGrant: {
			"Restart the authorization flow to get a new authorization code",
			"Check that the authorization code hasn't expired",
			"Verify the redirect URI matches exactly",
			"Ensure the code hasn't been used already",
		},
		
		OAuth2ErrorUnsupportedGrantType: {
			"Use 'authorization_code' for initial token exchange",
			"Use 'refresh_token' for token refresh",
			"Check the OAuth2 specification for supported grant types",
		},
		
		OAuth2ErrorTokenExpired: {
			"Use the refresh token to get a new access token",
			"Restart the authorization flow if no refresh token is available",
			"Implement automatic token refresh in your application",
		},
		
		OAuth2ErrorRefreshFailed: {
			"Restart the authorization flow to get new tokens",
			"Check that the refresh token hasn't expired",
			"Verify your client credentials are correct",
		},
		
		OAuth2ErrorInvalidState: {
			"Ensure state parameter is properly generated and stored",
			"Verify state validation logic in your application",
			"Check for potential CSRF attacks",
			"Restart the authorization flow with a new state parameter",
		},
		
		OAuth2ErrorMissingCode: {
			"Check your redirect URI configuration",
			"Verify the authorization callback is working correctly",
			"Ensure the user completes the authorization process",
			"Check for errors in the authorization flow",
		},
		
		OAuth2ErrorRedirectTimeout: {
			"Increase the timeout duration",
			"Check that your redirect server is running and accessible",
			"Verify the redirect URI is correct",
			"Ensure the user completes authorization within the timeout period",
		},
		
		OAuth2ErrorInvalidCredentials: {
			"Verify client ID, client secret, and redirect URI",
			"Check your TastyTrade application registration",
			"Ensure you're using the correct environment",
			"Regenerate credentials if necessary",
		},
		
		OAuth2ErrorConfigurationError: {
			"Review all OAuth2 configuration parameters",
			"Validate URLs and ensure they're accessible",
			"Check environment-specific settings",
			"Verify required fields are not empty",
		},
		
		OAuth2ErrorNetworkError: {
			"Check your internet connection",
			"Verify TastyTrade endpoints are accessible",
			"Check for firewall or proxy issues",
			"Try again after a short delay",
		},
		
		OAuth2ErrorPKCEError: {
			"Verify PKCE challenge generation is correct",
			"Check code verifier and challenge matching",
			"Ensure SHA256 hashing is implemented correctly",
			"Review PKCE implementation against RFC 7636",
		},
		
		OAuth2ErrorServerUnavailable: {
			"Wait and retry after a delay",
			"Check TastyTrade status page",
			"Implement exponential backoff",
			"Try again later if maintenance is ongoing",
		},
	}
	
	if actionList, exists := actions[errorCode]; exists {
		return actionList
	}
	
	return []string{
		"Check your OAuth2 configuration",
		"Verify your credentials are correct",
		"Try again after a short delay",
		"Contact TastyTrade support if the issue persists",
	}
}

// Predefined error variables for common OAuth2 errors
var (
	ErrOAuth2InvalidState = NewOAuth2Error(OAuth2ErrorInvalidState, 
		"Invalid state parameter: potential CSRF attack detected")
	
	ErrOAuth2TokenExpired = NewOAuth2Error(OAuth2ErrorTokenExpired, 
		"Access token has expired and needs to be refreshed")
	
	ErrOAuth2RefreshFailed = NewOAuth2Error(OAuth2ErrorRefreshFailed, 
		"Failed to refresh access token using refresh token")
	
	ErrOAuth2InvalidCode = NewOAuth2Error(OAuth2ErrorMissingCode, 
		"Authorization code is missing or invalid")
	
	ErrOAuth2RedirectTimeout = NewOAuth2Error(OAuth2ErrorRedirectTimeout, 
		"Timeout waiting for authorization callback")
	
	ErrOAuth2InvalidCredentials = NewOAuth2Error(OAuth2ErrorInvalidCredentials, 
		"Invalid OAuth2 client credentials")
	
	ErrOAuth2ConfigurationError = NewOAuth2Error(OAuth2ErrorConfigurationError, 
		"OAuth2 configuration error")
	
	ErrOAuth2NetworkError = NewOAuth2Error(OAuth2ErrorNetworkError, 
		"Network error during OAuth2 operation")
	
	ErrOAuth2PKCEError = NewOAuth2Error(OAuth2ErrorPKCEError, 
		"PKCE validation error")
	
	ErrOAuth2ServerUnavailable = NewOAuth2Error(OAuth2ErrorServerUnavailable, 
		"OAuth2 server is currently unavailable")
)

// OAuth2ErrorHandler provides centralized error handling for OAuth2 operations
type OAuth2ErrorHandler struct {
	// Configuration for error handling behavior
	EnableRetry     bool
	MaxRetries      int
	RetryDelay      time.Duration
	EnableLogging   bool
	LogSensitiveData bool
}

// NewOAuth2ErrorHandler creates a new error handler with default settings
func NewOAuth2ErrorHandler() *OAuth2ErrorHandler {
	return &OAuth2ErrorHandler{
		EnableRetry:      true,
		MaxRetries:       3,
		RetryDelay:       time.Second,
		EnableLogging:    true,
		LogSensitiveData: false,
	}
}

// HandleError processes an OAuth2 error and returns appropriate action
func (h *OAuth2ErrorHandler) HandleError(err error) (*OAuth2DetailedError, bool) {
	if err == nil {
		return nil, false
	}
	
	// Check if it's already a detailed OAuth2 error
	if detailedErr, ok := err.(*OAuth2DetailedError); ok {
		return detailedErr, detailedErr.IsRetryable() && h.EnableRetry
	}
	
	// Check if it's a standard OAuth2 error
	if oauth2Err, ok := err.(OAuth2Error); ok {
		detailedErr := NewOAuth2ErrorFromStandard(oauth2Err)
		return detailedErr, detailedErr.IsRetryable() && h.EnableRetry
	}
	
	// Handle other error types
	detailedErr := h.classifyGenericError(err)
	return detailedErr, detailedErr.IsRetryable() && h.EnableRetry
}

// classifyGenericError converts generic errors to OAuth2DetailedError
func (h *OAuth2ErrorHandler) classifyGenericError(err error) *OAuth2DetailedError {
	errMsg := err.Error()
	
	// Network-related errors
	if strings.Contains(errMsg, "connection") || strings.Contains(errMsg, "timeout") ||
	   strings.Contains(errMsg, "network") || strings.Contains(errMsg, "dial") {
		return NewOAuth2ErrorWithContext(OAuth2ErrorNetworkError, 
			"Network error during OAuth2 operation", 0, err)
	}
	
	// HTTP-related errors
	if strings.Contains(errMsg, "http") || strings.Contains(errMsg, "status") {
		return NewOAuth2ErrorWithContext(OAuth2ErrorServerError, 
			"HTTP error during OAuth2 operation", 0, err)
	}
	
	// Configuration-related errors
	if strings.Contains(errMsg, "config") || strings.Contains(errMsg, "invalid") ||
	   strings.Contains(errMsg, "missing") || strings.Contains(errMsg, "required") {
		return NewOAuth2ErrorWithContext(OAuth2ErrorConfigurationError, 
			"Configuration error in OAuth2 setup", 0, err)
	}
	
	// Default to unknown error
	return NewOAuth2ErrorWithContext("unknown_error", 
		"An unknown error occurred during OAuth2 operation", 0, err)
}

// ValidateState performs comprehensive state parameter validation
func ValidateState(expected, received string) error {
	if expected == "" {
		return NewOAuth2Error(OAuth2ErrorConfigurationError, 
			"Expected state parameter is empty - this indicates a configuration error")
	}
	
	if received == "" {
		return NewOAuth2Error(OAuth2ErrorInvalidState, 
			"State parameter is missing from the callback - this could indicate a CSRF attack or malformed callback")
	}
	
	if expected != received {
		return NewOAuth2ErrorWithContext(OAuth2ErrorInvalidState, 
			fmt.Sprintf("State parameter mismatch: expected '%s' but received '%s' - this indicates a potential CSRF attack", 
				expected, received), 0, nil)
	}
	
	return nil
}

// ValidateAuthorizationCode performs validation of authorization codes
func ValidateAuthorizationCode(code string) error {
	if code == "" {
		return NewOAuth2Error(OAuth2ErrorMissingCode, 
			"Authorization code is missing from the callback")
	}
	
	// Basic format validation (authorization codes are typically alphanumeric)
	if len(code) < 10 {
		return NewOAuth2Error(OAuth2ErrorInvalidGrant, 
			"Authorization code appears to be too short - it may be invalid")
	}
	
	// Check for suspicious characters that might indicate tampering
	for _, char := range code {
		if !((char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || 
			 (char >= '0' && char <= '9') || char == '-' || char == '_' || char == '.') {
			return NewOAuth2Error(OAuth2ErrorInvalidGrant, 
				"Authorization code contains invalid characters")
		}
	}
	
	return nil
}

// ValidateTokenResponse performs validation of token responses
func ValidateTokenResponse(response *TokenResponse) error {
	if response == nil {
		return NewOAuth2Error(OAuth2ErrorServerError, 
			"Token response is nil - server may have returned an empty response")
	}
	
	if response.AccessToken == "" {
		return NewOAuth2Error(OAuth2ErrorServerError, 
			"Access token is missing from token response")
	}
	
	if response.TokenType == "" {
		return NewOAuth2Error(OAuth2ErrorServerError, 
			"Token type is missing from token response")
	}
	
	if response.ExpiresIn <= 0 {
		return NewOAuth2Error(OAuth2ErrorServerError, 
			"Token expiration time is invalid or missing")
	}
	
	// Validate token format (basic checks)
	if len(response.AccessToken) < 10 {
		return NewOAuth2Error(OAuth2ErrorServerError, 
			"Access token appears to be too short - it may be invalid")
	}
	
	return nil
}

// IsTemporaryError determines if an error is temporary and should be retried
func IsTemporaryError(err error) bool {
	if err == nil {
		return false
	}
	
	// Check if it's a detailed OAuth2 error
	if detailedErr, ok := err.(*OAuth2DetailedError); ok {
		return detailedErr.IsRetryable()
	}
	
	// Check error message for temporary indicators
	errMsg := strings.ToLower(err.Error())
	temporaryIndicators := []string{
		"timeout", "temporary", "unavailable", "server_error", 
		"connection", "network", "503", "502", "504",
	}
	
	for _, indicator := range temporaryIndicators {
		if strings.Contains(errMsg, indicator) {
			return true
		}
	}
	
	return false
}

// WrapHTTPError converts HTTP errors to OAuth2DetailedError
func WrapHTTPError(resp *http.Response, originalErr error) *OAuth2DetailedError {
	if resp == nil {
		return NewOAuth2ErrorWithContext(OAuth2ErrorNetworkError, 
			"HTTP response is nil", 0, originalErr)
	}
	
	var errorCode string
	var errorType OAuth2ErrorType
	var severity OAuth2ErrorSeverity
	
	switch resp.StatusCode {
	case http.StatusBadRequest:
		errorCode = OAuth2ErrorInvalidRequest
		errorType = OAuth2ErrorTypeClient
		severity = OAuth2ErrorSeverityMedium
	case http.StatusUnauthorized:
		errorCode = OAuth2ErrorInvalidClient
		errorType = OAuth2ErrorTypeClient
		severity = OAuth2ErrorSeverityHigh
	case http.StatusForbidden:
		errorCode = OAuth2ErrorAccessDenied
		errorType = OAuth2ErrorTypeClient
		severity = OAuth2ErrorSeverityMedium
	case http.StatusInternalServerError:
		errorCode = OAuth2ErrorServerError
		errorType = OAuth2ErrorTypeServer
		severity = OAuth2ErrorSeverityMedium
	case http.StatusBadGateway, http.StatusServiceUnavailable, http.StatusGatewayTimeout:
		errorCode = OAuth2ErrorTemporarilyUnavailable
		errorType = OAuth2ErrorTypeServer
		severity = OAuth2ErrorSeverityLow
	default:
		errorCode = "http_error"
		errorType = OAuth2ErrorTypeUnknown
		severity = OAuth2ErrorSeverityLow
	}
	
	description := fmt.Sprintf("HTTP %d: %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	
	detailedErr := &OAuth2DetailedError{
		ErrorCode:        errorCode,
		ErrorDescription: description,
		HTTPStatusCode:   resp.StatusCode,
		Type:             errorType,
		Severity:         severity,
		Timestamp:        time.Now(),
		OriginalError:    originalErr,
	}
	
	detailedErr.TroubleshootingGuide = generateTroubleshootingGuide(errorCode)
	detailedErr.SuggestedActions = generateSuggestedActions(errorCode)
	
	return detailedErr
}