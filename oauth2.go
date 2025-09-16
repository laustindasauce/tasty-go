package tasty

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"
)

const (
	// Default OAuth2 scopes for TastyTrade API
	defaultScope = "read trade"

	// OAuth2 grant types
	grantTypeAuthorizationCode = "authorization_code"
	grantTypeRefreshToken      = "refresh_token"

	// OAuth2 response types
	responseTypeCode = "code"
)

// RedirectServer provides HTTP server functionality for handling OAuth2 redirects
type RedirectServer struct {
	server   *http.Server
	listener net.Listener
	codeChan chan string
	errChan  chan error
	state    string
	mutex    sync.RWMutex
	running  bool
}

// NewRedirectServer creates a new RedirectServer instance
func NewRedirectServer(state string) *RedirectServer {
	return &RedirectServer{
		codeChan: make(chan string, 1),
		errChan:  make(chan error, 1),
		state:    state,
	}
}

// Start starts the HTTP server on the specified port with timeout handling
// If port is 0, a random available port will be chosen
func (rs *RedirectServer) Start(port int) error {
	rs.mutex.Lock()
	defer rs.mutex.Unlock()

	if rs.running {
		return NewOAuth2Error(OAuth2ErrorConfigurationError, "redirect server is already running")
	}

	// Create listener on specified port
	addr := ":" + strconv.Itoa(port)
	if port == 0 {
		addr = ":0" // Let the system choose an available port
	}

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return NewOAuth2ErrorWithContext(OAuth2ErrorNetworkError,
			fmt.Sprintf("failed to create listener on port %d", port), 0, err)
	}

	rs.listener = listener

	// Create HTTP server with redirect handler
	mux := http.NewServeMux()
	mux.HandleFunc("/", rs.handleRedirect)

	rs.server = &http.Server{
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	rs.running = true

	// Start server in goroutine
	go func() {
		if err := rs.server.Serve(rs.listener); err != nil && err != http.ErrServerClosed {
			rs.mutex.Lock()
			rs.running = false
			rs.mutex.Unlock()

			select {
			case rs.errChan <- NewOAuth2ErrorWithContext(OAuth2ErrorServerError,
				"redirect server error", 0, err):
			default:
				// Channel is full, ignore
			}
		}
	}()

	return nil
}

// GetPort returns the port the server is listening on
func (rs *RedirectServer) GetPort() int {
	rs.mutex.RLock()
	defer rs.mutex.RUnlock()

	if rs.listener == nil {
		return 0
	}

	addr := rs.listener.Addr().(*net.TCPAddr)
	return addr.Port
}

// GetRedirectURI returns the full redirect URI for this server
func (rs *RedirectServer) GetRedirectURI() string {
	port := rs.GetPort()
	if port == 0 {
		return ""
	}
	return fmt.Sprintf("http://localhost:%d", port)
}

// WaitForCode waits for an authorization code with timeout protection
func (rs *RedirectServer) WaitForCode(timeout time.Duration) (string, error) {
	rs.mutex.RLock()
	if !rs.running {
		rs.mutex.RUnlock()
		return "", NewOAuth2Error(OAuth2ErrorConfigurationError, "redirect server is not running")
	}
	rs.mutex.RUnlock()

	// Create timeout context
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	select {
	case code := <-rs.codeChan:
		return code, nil
	case err := <-rs.errChan:
		if detailedErr, ok := err.(*OAuth2DetailedError); ok {
			return "", detailedErr
		}
		return "", NewOAuth2ErrorWithContext(OAuth2ErrorServerError,
			"server error while waiting for code", 0, err)
	case <-ctx.Done():
		return "", NewOAuth2ErrorWithContext(OAuth2ErrorRedirectTimeout,
			fmt.Sprintf("timeout waiting for authorization code after %v", timeout), 0, ctx.Err())
	}
}

// handleRedirect handles the OAuth2 redirect callback
func (rs *RedirectServer) handleRedirect(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	query := r.URL.Query()

	// Check for error parameter first
	if errorCode := query.Get("error"); errorCode != "" {
		errorDesc := query.Get("error_description")
		errorURI := query.Get("error_uri")

		oauthErr := OAuth2Error{
			ErrorCode:        errorCode,
			ErrorDescription: errorDesc,
			ErrorURI:         errorURI,
			State:            query.Get("state"),
		}

		// Convert to detailed error
		detailedErr := NewOAuth2ErrorFromStandard(oauthErr)

		// Send error response to user
		rs.sendErrorResponse(w, oauthErr)

		// Send detailed error to channel
		select {
		case rs.errChan <- detailedErr:
		default:
			// Channel is full, ignore
		}
		return
	}

	// Get authorization code
	code := query.Get("code")
	if code == "" {
		detailedErr := NewOAuth2Error(OAuth2ErrorMissingCode, "missing authorization code in redirect")
		rs.sendErrorResponse(w, OAuth2Error{
			ErrorCode:        OAuth2ErrorMissingCode,
			ErrorDescription: "Missing authorization code",
		})

		select {
		case rs.errChan <- detailedErr:
		default:
			// Channel is full, ignore
		}
		return
	}

	// Validate state parameter for CSRF protection
	receivedState := query.Get("state")
	if err := ValidateState(rs.state, receivedState); err != nil {
		rs.sendErrorResponse(w, OAuth2Error{
			ErrorCode:        OAuth2ErrorInvalidState,
			ErrorDescription: "Invalid state parameter",
		})

		select {
		case rs.errChan <- err:
		default:
			// Channel is full, ignore
		}
		return
	}

	// Send success response to user
	rs.sendSuccessResponse(w)

	// Send code to channel
	select {
	case rs.codeChan <- code:
	default:
		// Channel is full, ignore
	}
}

// sendSuccessResponse sends a success HTML response to the user
func (rs *RedirectServer) sendSuccessResponse(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	html := `<!DOCTYPE html>
<html>
<head>
    <title>Authorization Successful</title>
    <style>
        body { font-family: Arial, sans-serif; text-align: center; margin-top: 50px; }
        .success { color: #28a745; }
        .container { max-width: 500px; margin: 0 auto; padding: 20px; }
    </style>
</head>
<body>
    <div class="container">
        <h1 class="success">✓ Authorization Successful</h1>
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

// sendErrorResponse sends an error HTML response to the user
func (rs *RedirectServer) sendErrorResponse(w http.ResponseWriter, oauthErr OAuth2Error) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusBadRequest)

	html := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <title>Authorization Failed</title>
    <style>
        body { font-family: Arial, sans-serif; text-align: center; margin-top: 50px; }
        .error { color: #dc3545; }
        .container { max-width: 500px; margin: 0 auto; padding: 20px; }
        .error-details { background: #f8f9fa; padding: 15px; border-radius: 5px; margin-top: 20px; }
    </style>
</head>
<body>
    <div class="container">
        <h1 class="error">✗ Authorization Failed</h1>
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

// Shutdown gracefully shuts down the HTTP server with proper resource cleanup
func (rs *RedirectServer) Shutdown(timeout time.Duration) error {
	rs.mutex.Lock()
	defer rs.mutex.Unlock()

	if !rs.running {
		return nil // Already shut down
	}

	rs.running = false

	if rs.server == nil {
		return nil
	}

	// Create shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Attempt graceful shutdown
	if err := rs.server.Shutdown(ctx); err != nil {
		// Force close if graceful shutdown fails
		if rs.listener != nil {
			rs.listener.Close()
		}
		return NewOAuth2ErrorWithContext(OAuth2ErrorServerError,
			"failed to shutdown server gracefully", 0, err)
	}

	// Close channels to prevent goroutine leaks
	close(rs.codeChan)
	close(rs.errChan)

	return nil
}

// IsRunning returns whether the redirect server is currently running
func (rs *RedirectServer) IsRunning() bool {
	rs.mutex.RLock()
	defer rs.mutex.RUnlock()
	return rs.running
}

// ValidateOAuth2Options validates OAuth2Options and ensures PKCE is not requested
func ValidateOAuth2Options(options OAuth2Options) error {
	if options.UsePKCE {
		return fmt.Errorf("PKCE is not supported by TastyTrade OAuth2 implementation - using PKCE parameters results in '405 Method Not Allowed' errors")
	}
	return nil
}
