package tasty

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// AdvancedMockOAuth2Server provides more sophisticated testing capabilities
type AdvancedMockOAuth2Server struct {
	*MockOAuth2Server
	requestDelay    time.Duration
	failureRate     float64
	requestCounter  int64
	concurrentLimit int
	currentRequests int64
}

// NewAdvancedMockOAuth2Server creates an advanced mock server
func NewAdvancedMockOAuth2Server() *AdvancedMockOAuth2Server {
	base := NewMockOAuth2Server()
	return &AdvancedMockOAuth2Server{
		MockOAuth2Server: base,
		concurrentLimit:  100, // Default limit
	}
}

// SetRequestDelay adds artificial delay to requests
func (m *AdvancedMockOAuth2Server) SetRequestDelay(delay time.Duration) {
	m.requestDelay = delay
}

// SetFailureRate sets the probability of request failure (0.0 to 1.0)
func (m *AdvancedMockOAuth2Server) SetFailureRate(rate float64) {
	m.failureRate = rate
}

// SetConcurrentLimit sets the maximum concurrent requests
func (m *AdvancedMockOAuth2Server) SetConcurrentLimit(limit int) {
	m.concurrentLimit = limit
}

// simulateLoad simulates server load conditions
func (m *AdvancedMockOAuth2Server) simulateLoad() bool {
	// Check concurrent request limit
	current := atomic.AddInt64(&m.currentRequests, 1)
	defer atomic.AddInt64(&m.currentRequests, -1)
	
	if int(current) > m.concurrentLimit {
		return false // Reject due to load
	}
	
	// Add artificial delay
	if m.requestDelay > 0 {
		time.Sleep(m.requestDelay)
	}
	
	// Simulate random failures
	if m.failureRate > 0 {
		counter := atomic.AddInt64(&m.requestCounter, 1)
		if float64(counter%100)/100.0 < m.failureRate {
			return false // Simulate failure
		}
	}
	
	return true
}

// TestOAuth2Client_FullIntegrationFlow tests complete OAuth2 flow
func TestOAuth2Client_FullIntegrationFlow(t *testing.T) {
	mockServer := NewMockOAuth2Server()
	defer mockServer.Close()
	
	config := createTestOAuth2Config(mockServer)
	client, err := newOAuth2ClientInternal(config, nil)
	require.NoError(t, err)
	
	// Step 1: Generate authorization URL
	authURL, err := client.GetAuthorizationURL()
	require.NoError(t, err)
	assert.Contains(t, authURL, config.ClientID)
	assert.Contains(t, authURL, "redirect_uri=http%3A%2F%2Flocalhost%3A8080%2Fcallback") // URL encoded
	assert.Contains(t, authURL, config.State)
	
	// Step 2: Start redirect server
	redirectServer, err := client.StartRedirectServer(0)
	require.NoError(t, err)
	defer redirectServer.Shutdown(5 * time.Second)
	
	port := redirectServer.GetPort()
	assert.Greater(t, port, 0)
	
	// Step 3: Simulate authorization callback
	go func() {
		time.Sleep(100 * time.Millisecond) // Small delay to ensure server is ready
		callbackURL := fmt.Sprintf("http://localhost:%d?code=test_auth_code&state=%s", 
			port, config.State)
		http.Get(callbackURL)
	}()
	
	// Step 4: Wait for authorization code
	code, err := redirectServer.WaitForCode(2 * time.Second)
	require.NoError(t, err)
	assert.Equal(t, "test_auth_code", code)
	
	// Step 5: Exchange code for tokens
	tokenResponse, err := client.ExchangeCodeForTokens(code)
	require.NoError(t, err)
	assert.NotEmpty(t, tokenResponse.AccessToken)
	assert.NotEmpty(t, tokenResponse.RefreshToken)
	
	// Step 6: Verify authentication
	assert.True(t, client.IsAuthenticated())
	
	// Step 7: Test token refresh
	mockServer.SetTokenResponse(&TokenResponse{
		AccessToken:  "new_access_token",
		RefreshToken: "new_refresh_token",
		TokenType:    "Bearer",
		ExpiresIn:    3600,
	})
	
	refreshResponse, err := client.RefreshTokens()
	require.NoError(t, err)
	assert.Equal(t, "new_access_token", refreshResponse.AccessToken)
	
	// Step 8: Clear authentication
	client.ClearTokens()
	assert.False(t, client.IsAuthenticated())
}

// TestOAuth2Client_ConcurrentOperations tests concurrent OAuth2 operations
func TestOAuth2Client_ConcurrentOperations(t *testing.T) {
	mockServer := NewAdvancedMockOAuth2Server()
	defer mockServer.Close()
	
	config := createTestOAuth2Config(mockServer.MockOAuth2Server)
	client, err := newOAuth2ClientInternal(config, nil)
	require.NoError(t, err)
	
	// Set initial tokens
	client.tokenManager.SetTokens("access_token", "refresh_token", 3600)
	
	const numGoroutines = 20
	const numOperations = 50
	
	var wg sync.WaitGroup
	errors := make(chan error, numGoroutines*numOperations)
	
	// Test concurrent authorization URL generation
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				authURL, err := client.GetAuthorizationURL()
				if err != nil {
					errors <- fmt.Errorf("goroutine %d, operation %d: %w", id, j, err)
					return
				}
				if authURL == "" {
					errors <- fmt.Errorf("goroutine %d, operation %d: empty auth URL", id, j)
					return
				}
			}
		}(i)
	}
	
	// Test concurrent state validation
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				err := client.ValidateState(client.GetState())
				if err != nil {
					errors <- fmt.Errorf("state validation goroutine %d, operation %d: %w", id, j, err)
					return
				}
			}
		}(i)
	}
	
	// Test concurrent authentication checks
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				_ = client.IsAuthenticated()
				_ = client.GetConfig()
				_ = client.GetState()
			}
		}(i)
	}
	
	wg.Wait()
	close(errors)
	
	// Check for errors
	errorCount := 0
	for err := range errors {
		t.Errorf("Concurrent operation error: %v", err)
		errorCount++
	}
	
	assert.Equal(t, 0, errorCount, "Expected no errors in concurrent operations")
}

// TestOAuth2Client_LoadTesting tests OAuth2 client under high load
func TestOAuth2Client_LoadTesting(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping load test in short mode")
	}
	
	mockServer := NewAdvancedMockOAuth2Server()
	defer mockServer.Close()
	
	// Configure server for load testing
	mockServer.SetRequestDelay(1 * time.Millisecond)
	mockServer.SetFailureRate(0.05) // 5% failure rate
	mockServer.SetConcurrentLimit(50)
	
	config := createTestOAuth2Config(mockServer.MockOAuth2Server)
	
	const numClients = 10
	const operationsPerClient = 100
	
	type TestResult struct {
		ClientID     int
		Successes    int
		Failures     int
		Duration     time.Duration
		AvgLatency   time.Duration
	}
	
	var wg sync.WaitGroup
	results := make(chan TestResult, numClients)
	
	start := time.Now()
	
	for i := 0; i < numClients; i++ {
		wg.Add(1)
		go func(clientID int) {
			defer wg.Done()
			
			client, err := newOAuth2ClientInternal(config, nil)
			if err != nil {
				results <- TestResult{ClientID: clientID, Failures: operationsPerClient}
				return
			}
			
			client.tokenManager.SetTokens("access_token", "refresh_token", 3600)
			
			successes := 0
			failures := 0
			totalLatency := time.Duration(0)
			
			for j := 0; j < operationsPerClient; j++ {
				opStart := time.Now()
				
				// Perform various operations
				_, err1 := client.GetAuthorizationURL()
				err2 := client.ValidateState(client.GetState())
				_ = client.IsAuthenticated()
				
				latency := time.Since(opStart)
				totalLatency += latency
				
				if err1 != nil || err2 != nil {
					failures++
				} else {
					successes++
				}
			}
			
			avgLatency := totalLatency / time.Duration(operationsPerClient)
			results <- TestResult{
				ClientID:   clientID,
				Successes:  successes,
				Failures:   failures,
				Duration:   time.Since(start),
				AvgLatency: avgLatency,
			}
		}(i)
	}
	
	wg.Wait()
	close(results)
	
	totalDuration := time.Since(start)
	totalSuccesses := 0
	totalFailures := 0
	maxLatency := time.Duration(0)
	
	for result := range results {
		totalSuccesses += result.Successes
		totalFailures += result.Failures
		if result.AvgLatency > maxLatency {
			maxLatency = result.AvgLatency
		}
		
		t.Logf("Client %d: %d successes, %d failures, avg latency: %v", 
			result.ClientID, result.Successes, result.Failures, result.AvgLatency)
	}
	
	totalOperations := totalSuccesses + totalFailures
	successRate := float64(totalSuccesses) / float64(totalOperations) * 100
	
	t.Logf("Load test results:")
	t.Logf("  Total operations: %d", totalOperations)
	t.Logf("  Success rate: %.2f%%", successRate)
	t.Logf("  Total duration: %v", totalDuration)
	t.Logf("  Max average latency: %v", maxLatency)
	t.Logf("  Operations per second: %.2f", float64(totalOperations)/totalDuration.Seconds())
	
	// Assertions for acceptable performance
	assert.Greater(t, successRate, 90.0, "Success rate should be above 90%")
	assert.Less(t, maxLatency, 100*time.Millisecond, "Average latency should be under 100ms")
	assert.Less(t, totalDuration, 30*time.Second, "Total test duration should be under 30 seconds")
}

// TestOAuth2Client_ErrorRecovery tests error recovery scenarios
func TestOAuth2Client_ErrorRecovery(t *testing.T) {
	// Create a server that fails initially, then recovers
	failureCount := int64(0)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count := atomic.AddInt64(&failureCount, 1)
		
		if r.URL.Path == "/oauth/token" {
			if count <= 3 {
				// Fail first 3 requests
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(OAuth2Error{
					ErrorCode:        OAuth2ErrorServerError,
					ErrorDescription: "Temporary server error",
				})
				return
			}
			
			// Succeed after 3 failures
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(TokenResponse{
				AccessToken:  "recovered_access_token",
				RefreshToken: "recovered_refresh_token",
				TokenType:    "Bearer",
				ExpiresIn:    3600,
			})
			return
		}
		
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()
	
	config := OAuth2Config{
		ClientID:     "test_client",
		ClientSecret: "test_secret",
		RedirectURI:  "http://localhost:8080/callback",
		TokenURL:     server.URL + "/oauth/token",
		AuthURL:      server.URL + "/oauth/authorize",
	}
	
	client, err := newOAuth2ClientInternal(config, nil)
	require.NoError(t, err)
	
	// Set initial tokens
	client.tokenManager.SetTokens("old_token", "refresh_token", 3600)
	
	// Attempt refresh - should eventually succeed after retries
	response, err := client.refreshTokensWithRetry("refresh_token", 5)
	require.NoError(t, err)
	require.NotNil(t, response)
	
	assert.Equal(t, "recovered_access_token", response.AccessToken)
	assert.Greater(t, atomic.LoadInt64(&failureCount), int64(3))
}

// TestOAuth2Client_MemoryLeakPrevention tests memory leak prevention
func TestOAuth2Client_MemoryLeakPrevention(t *testing.T) {
	mockServer := NewMockOAuth2Server()
	defer mockServer.Close()
	
	config := createTestOAuth2Config(mockServer)
	
	// Create and destroy many clients rapidly
	for i := 0; i < 1000; i++ {
		client, err := newOAuth2ClientInternal(config, nil)
		require.NoError(t, err)
		
		// Use the client briefly
		client.tokenManager.SetTokens(fmt.Sprintf("token_%d", i), 
			fmt.Sprintf("refresh_%d", i), 3600)
		
		_, err = client.GetAuthorizationURL()
		require.NoError(t, err)
		
		assert.True(t, client.IsAuthenticated())
		
		// Clear resources
		client.ClearTokens()
		assert.False(t, client.IsAuthenticated())
		
		// Verify token manager is cleared
		tm := client.GetTokenManager()
		assert.False(t, tm.HasValidToken())
		assert.False(t, tm.HasRefreshToken())
	}
}

// TestRedirectServer_HighConcurrency tests redirect server under high concurrency
func TestRedirectServer_HighConcurrency(t *testing.T) {
	server := NewRedirectServer("test_state")
	
	err := server.Start(0)
	require.NoError(t, err)
	defer server.Shutdown(5 * time.Second)
	
	const numRequests = 100
	var wg sync.WaitGroup
	results := make(chan error, numRequests)
	
	baseURL := server.GetRedirectURI()
	
	// Send many concurrent requests
	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			
			// Alternate between success and error scenarios
			var url string
			if id%2 == 0 {
				url = fmt.Sprintf("%s?code=test_code_%d&state=test_state", baseURL, id)
			} else {
				url = fmt.Sprintf("%s?error=access_denied&state=test_state", baseURL)
			}
			
			resp, err := http.Get(url)
			if err != nil {
				results <- err
				return
			}
			defer resp.Body.Close()
			
			// Check that server responds appropriately
			if id%2 == 0 && resp.StatusCode != http.StatusOK {
				results <- fmt.Errorf("expected 200 for success case, got %d", resp.StatusCode)
			} else if id%2 == 1 && resp.StatusCode != http.StatusBadRequest {
				results <- fmt.Errorf("expected 400 for error case, got %d", resp.StatusCode)
			}
		}(i)
	}
	
	wg.Wait()
	close(results)
	
	// Check for errors
	errorCount := 0
	for err := range results {
		t.Errorf("Concurrent redirect server error: %v", err)
		errorCount++
	}
	
	assert.Equal(t, 0, errorCount, "Expected no errors in concurrent redirect server operations")
}

// TestOAuth2Client_TimeoutHandling tests timeout handling
func TestOAuth2Client_TimeoutHandling(t *testing.T) {
	// Create a server that responds very slowly
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second) // Longer than client timeout
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(TokenResponse{
			AccessToken:  "slow_token",
			RefreshToken: "slow_refresh",
			TokenType:    "Bearer",
			ExpiresIn:    3600,
		})
	}))
	defer server.Close()
	
	// Create client with short timeout
	httpClient := &http.Client{
		Timeout: 500 * time.Millisecond,
	}
	
	config := OAuth2Config{
		ClientID:     "test_client",
		ClientSecret: "test_secret",
		RedirectURI:  "http://localhost:8080/callback",
		TokenURL:     server.URL + "/token",
		AuthURL:      server.URL + "/auth",
	}
	
	client, err := newOAuth2ClientInternal(config, httpClient)
	require.NoError(t, err)
	
	// Attempt operations that should timeout or fail due to network issues
	_, err = client.ExchangeCodeForTokens("valid_authorization_code_123")
	assert.Error(t, err)
	// The error might be timeout or network error depending on timing
	assert.True(t, strings.Contains(err.Error(), "timeout") || strings.Contains(err.Error(), "network") || strings.Contains(err.Error(), "context deadline exceeded"))
	
	client.tokenManager.SetTokens("token", "refresh", 3600)
	_, err = client.RefreshTokens()
	assert.Error(t, err)
	// The error might be timeout or network error depending on timing
	assert.True(t, strings.Contains(err.Error(), "timeout") || strings.Contains(err.Error(), "network") || strings.Contains(err.Error(), "context deadline exceeded"))
}

// TestOAuth2Client_StateManagement tests state parameter management
func TestOAuth2Client_StateManagement(t *testing.T) {
	mockServer := NewMockOAuth2Server()
	defer mockServer.Close()
	
	config := createTestOAuth2Config(mockServer)
	client, err := newOAuth2ClientInternal(config, nil)
	require.NoError(t, err)
	
	// Test state consistency across multiple authorization URLs
	state1 := client.GetState()
	authURL1, err := client.GetAuthorizationURL()
	require.NoError(t, err)
	
	state2 := client.GetState()
	authURL2, err := client.GetAuthorizationURL()
	require.NoError(t, err)
	
	// State should remain consistent
	assert.Equal(t, state1, state2)
	assert.Contains(t, authURL1, state1)
	assert.Contains(t, authURL2, state1)
	
	// Test state validation
	assert.NoError(t, client.ValidateState(state1))
	assert.Error(t, client.ValidateState("wrong_state"))
	assert.Error(t, client.ValidateState(""))
}

// TestOAuth2Client_ConfigurationImmutability tests configuration immutability
func TestOAuth2Client_ConfigurationImmutability(t *testing.T) {
	mockServer := NewMockOAuth2Server()
	defer mockServer.Close()
	
	originalConfig := createTestOAuth2Config(mockServer)
	client, err := newOAuth2ClientInternal(originalConfig, nil)
	require.NoError(t, err)
	
	// Get configuration copy
	configCopy := client.GetConfig()
	
	// Modify the copy
	configCopy.ClientID = "modified_client_id"
	configCopy.RedirectURI = "https://modified.com/callback"
	
	// Original client configuration should be unchanged
	currentConfig := client.GetConfig()
	assert.Equal(t, originalConfig.ClientID, currentConfig.ClientID)
	assert.Equal(t, originalConfig.RedirectURI, currentConfig.RedirectURI)
	assert.NotEqual(t, "modified_client_id", currentConfig.ClientID)
}

// TestOAuth2Client_ThreadSafety tests thread safety of OAuth2 client
func TestOAuth2Client_ThreadSafety(t *testing.T) {
	mockServer := NewMockOAuth2Server()
	defer mockServer.Close()
	
	config := createTestOAuth2Config(mockServer)
	client, err := newOAuth2ClientInternal(config, nil)
	require.NoError(t, err)
	
	const numGoroutines = 50
	const numOperations = 20
	
	var wg sync.WaitGroup
	errors := make(chan error, numGoroutines*numOperations)
	
	// Test concurrent read operations
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				// These operations should be thread-safe
				_ = client.GetState()
				_ = client.GetConfig()
				_ = client.IsAuthenticated()
				
				_, err := client.GetAuthorizationURL()
				if err != nil {
					errors <- err
					return
				}
			}
		}()
	}
	
	// Test concurrent token operations
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				// Set tokens with unique values
				client.tokenManager.SetTokens(
					fmt.Sprintf("access_%d_%d", id, j),
					fmt.Sprintf("refresh_%d_%d", id, j),
					3600,
				)
				
				// Check authentication status
				_ = client.IsAuthenticated()
				
				// Clear tokens occasionally
				if j%5 == 0 {
					client.ClearTokens()
				}
			}
		}(i)
	}
	
	wg.Wait()
	close(errors)
	
	// Check for errors
	for err := range errors {
		t.Errorf("Thread safety error: %v", err)
	}
}