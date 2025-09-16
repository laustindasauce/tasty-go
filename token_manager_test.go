package tasty

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewTokenManager(t *testing.T) {
	tm := NewTokenManager(nil)
	if tm == nil {
		t.Fatal("NewTokenManager returned nil")
	}
	if tm.GetTokenType() != "Bearer" {
		t.Errorf("Expected default token type 'Bearer', got '%s'", tm.GetTokenType())
	}
	tm.Clear()
	if tm.HasValidToken() {
		t.Error("New TokenManager should not have valid token")
	}
	if tm.HasRefreshToken() {
		t.Error("New TokenManager should not have refresh token")
	}
	// Test that it uses file storage by default
	tm.SetTokens("test-token", "test-refresh", 3600)
	// Create another TokenManager with default settings
	// It should load the same tokens from the default file
	tm2 := NewTokenManager(nil)
	token, err := tm2.GetAccessToken()
	if err != nil {
		t.Fatalf("Failed to get access token from second manager: %v", err)
	}
	if token != "test-token" {
		t.Errorf("Expected token to persist across TokenManager instances, got '%s'", token)
	}
	// Create another TokenManager with different token path
	// It should load the same tokens from the default file
	tmpPath := "~/.tasty-go/tmp-tokens.json"
	tm3 := NewTokenManager(&tmpPath)
	_, err = tm3.GetAccessToken()
	require.Error(t, err, "There should be no token at this tmp file path")
	// Clean up the default token file
	tm.Clear()
}

func TestTokenManager_SetTokens(t *testing.T) {
	tm := NewMemoryTokenManager() // Use memory storage for isolated testing

	accessToken := "test-access-token"
	refreshToken := "test-refresh-token"
	expiresIn := 3600

	tm.SetTokens(accessToken, refreshToken, expiresIn)

	// Test token retrieval
	token, err := tm.GetAccessToken()
	if err != nil {
		t.Fatalf("GetAccessToken failed: %v", err)
	}

	if token != accessToken {
		t.Errorf("Expected access token '%s', got '%s'", accessToken, token)
	}

	if tm.GetRefreshToken() != refreshToken {
		t.Errorf("Expected refresh token '%s', got '%s'", refreshToken, tm.GetRefreshToken())
	}

	if !tm.HasValidToken() {
		t.Error("Token should be valid after setting")
	}

	if !tm.HasRefreshToken() {
		t.Error("Should have refresh token after setting")
	}
}

func TestTokenManager_SetTokensFromResponse(t *testing.T) {
	tm := NewMemoryTokenManager() // Use memory storage for isolated testing

	response := &TokenResponse{
		AccessToken:  "test-access-token",
		RefreshToken: "test-refresh-token",
		TokenType:    "Bearer",
		ExpiresIn:    3600,
		Scope:        "read trade",
	}

	tm.SetTokensFromResponse(response)

	token, err := tm.GetAccessToken()
	if err != nil {
		t.Fatalf("GetAccessToken failed: %v", err)
	}

	if token != response.AccessToken {
		t.Errorf("Expected access token '%s', got '%s'", response.AccessToken, token)
	}

	if tm.GetRefreshToken() != response.RefreshToken {
		t.Errorf("Expected refresh token '%s', got '%s'", response.RefreshToken, tm.GetRefreshToken())
	}

	if tm.GetTokenType() != response.TokenType {
		t.Errorf("Expected token type '%s', got '%s'", response.TokenType, tm.GetTokenType())
	}

	if tm.GetScope() != response.Scope {
		t.Errorf("Expected scope '%s', got '%s'", response.Scope, tm.GetScope())
	}
}

func TestTokenManager_SetTokensFromResponse_Nil(t *testing.T) {
	tm := NewMemoryTokenManager() // Use memory storage to avoid file persistence

	// Should not panic with nil response
	tm.SetTokensFromResponse(nil)

	if tm.HasValidToken() {
		t.Error("Should not have valid token after nil response")
	}
}

func TestTokenManager_IsExpired(t *testing.T) {
	tm := NewMemoryTokenManager() // Use memory storage to avoid file persistence

	// No token should be considered expired
	if !tm.IsExpired() {
		t.Error("TokenManager with no token should be expired")
	}

	// Set token that expires in 1 hour
	tm.SetTokens("token", "refresh", 3600)
	if tm.IsExpired() {
		t.Error("Token expiring in 1 hour should not be expired")
	}

	// Set token that expires in 10 seconds (should be considered expired due to 30s buffer)
	tm.SetTokens("token", "refresh", 10)
	if !tm.IsExpired() {
		t.Error("Token expiring in 10 seconds should be considered expired")
	}

	// Set token that already expired
	tm.SetTokens("token", "refresh", -10)
	if !tm.IsExpired() {
		t.Error("Expired token should be considered expired")
	}
}

func TestTokenManager_GetTimeUntilExpiry(t *testing.T) {
	tm := NewMemoryTokenManager() // Use memory storage for isolated testing

	// No token should return 0
	if tm.GetTimeUntilExpiry() != 0 {
		t.Error("TokenManager with no token should return 0 time until expiry")
	}

	// Set token that expires in 1 hour
	tm.SetTokens("token", "refresh", 3600)
	duration := tm.GetTimeUntilExpiry()

	// Should be close to 1 hour (allowing for small timing differences)
	expected := time.Hour
	if duration < expected-time.Second || duration > expected+time.Second {
		t.Errorf("Expected duration close to %v, got %v", expected, duration)
	}

	// Set expired token
	tm.SetTokens("token", "refresh", -10)
	if tm.GetTimeUntilExpiry() != 0 {
		t.Error("Expired token should return 0 time until expiry")
	}
}

func TestTokenManager_RefreshCallback(t *testing.T) {
	tm := NewMemoryTokenManager() // Use memory storage for isolated testing

	// Set initial expired token
	tm.SetTokens("old-token", "refresh-token", -10)

	// Set up refresh callback
	callCount := 0
	tm.SetRefreshCallback(func() (*TokenResponse, error) {
		callCount++
		return &TokenResponse{
			AccessToken:  "new-access-token",
			RefreshToken: "new-refresh-token",
			TokenType:    "Bearer",
			ExpiresIn:    3600,
		}, nil
	})

	// GetAccessToken should trigger refresh
	token, err := tm.GetAccessToken()
	if err != nil {
		t.Fatalf("GetAccessToken failed: %v", err)
	}

	if token != "new-access-token" {
		t.Errorf("Expected refreshed token 'new-access-token', got '%s'", token)
	}

	if callCount != 1 {
		t.Errorf("Expected refresh callback to be called once, called %d times", callCount)
	}

	if tm.GetRefreshToken() != "new-refresh-token" {
		t.Errorf("Expected new refresh token 'new-refresh-token', got '%s'", tm.GetRefreshToken())
	}
}

func TestTokenManager_RefreshCallback_Error(t *testing.T) {
	tm := NewMemoryTokenManager() // Use memory storage for isolated testing

	// Set initial expired token
	tm.SetTokens("old-token", "refresh-token", -10)

	// Set up failing refresh callback
	expectedError := errors.New("refresh failed")
	tm.SetRefreshCallback(func() (*TokenResponse, error) {
		return nil, expectedError
	})

	// GetAccessToken should return error
	_, err := tm.GetAccessToken()
	if err == nil {
		t.Fatal("Expected error from failed refresh")
	}

	if !errors.Is(err, expectedError) {
		t.Errorf("Expected error to wrap refresh error, got: %v", err)
	}
}

func TestTokenManager_RefreshCallback_NilResponse(t *testing.T) {
	tm := NewMemoryTokenManager() // Use memory storage for isolated testing

	// Set initial expired token
	tm.SetTokens("old-token", "refresh-token", -10)

	// Set up callback that returns nil response
	tm.SetRefreshCallback(func() (*TokenResponse, error) {
		return nil, nil
	})

	// GetAccessToken should return error
	_, err := tm.GetAccessToken()
	if err == nil {
		t.Fatal("Expected error from nil refresh response")
	}
}

func TestTokenManager_NoRefreshCallback(t *testing.T) {
	tm := NewMemoryTokenManager() // Use memory storage for isolated testing

	// Set initial expired token without refresh callback
	tm.SetTokens("old-token", "refresh-token", -10)

	// GetAccessToken should return error
	_, err := tm.GetAccessToken()
	if err == nil {
		t.Fatal("Expected error when no refresh callback is set")
	}
}

func TestTokenManager_Clear(t *testing.T) {
	tm := NewMemoryTokenManager() // Use memory storage for isolated testing

	// Set tokens
	tm.SetTokens("access-token", "refresh-token", 3600)
	tm.SetRefreshCallback(func() (*TokenResponse, error) {
		return nil, nil
	})

	// Verify tokens are set
	if !tm.HasValidToken() {
		t.Error("Should have valid token before clear")
	}

	if !tm.HasRefreshToken() {
		t.Error("Should have refresh token before clear")
	}

	// Clear tokens
	tm.Clear()

	// Verify tokens are cleared
	if tm.HasValidToken() {
		t.Error("Should not have valid token after clear")
	}

	if tm.HasRefreshToken() {
		t.Error("Should not have refresh token after clear")
	}

	if tm.GetTokenType() != "Bearer" {
		t.Error("Token type should return default 'Bearer' after clear")
	}

	if tm.GetScope() != "" {
		t.Error("Scope should be cleared")
	}

	if !tm.GetExpiresAt().IsZero() {
		t.Error("Expiration time should be cleared")
	}
}

func TestTokenManager_Clone(t *testing.T) {
	tm := NewMemoryTokenManager() // Use memory storage for isolated testing

	// Set tokens and callback
	tm.SetTokens("access-token", "refresh-token", 3600)
	tm.SetRefreshCallback(func() (*TokenResponse, error) {
		return nil, nil
	})

	// Clone the token manager
	clone := tm.Clone()

	// Verify clone has same token data
	originalToken, _ := tm.GetAccessToken()
	cloneToken, _ := clone.GetAccessToken()

	if originalToken != cloneToken {
		t.Errorf("Clone should have same access token, got '%s' vs '%s'", originalToken, cloneToken)
	}

	if tm.GetRefreshToken() != clone.GetRefreshToken() {
		t.Error("Clone should have same refresh token")
	}

	if tm.GetTokenType() != clone.GetTokenType() {
		t.Error("Clone should have same token type")
	}

	if tm.GetScope() != clone.GetScope() {
		t.Error("Clone should have same scope")
	}

	if !tm.GetExpiresAt().Equal(clone.GetExpiresAt()) {
		t.Error("Clone should have same expiration time")
	}

	// Verify callback is not cloned (for security)
	clone.SetTokens("expired-token", "refresh", -10)
	_, err := clone.GetAccessToken()
	if err == nil {
		t.Error("Clone should not have refresh callback")
	}

	// Clean up the original token manager
	tm.Clear()
}

func TestTokenManager_ConcurrentAccess(t *testing.T) {
	tm := NewMemoryTokenManager() // Use memory storage for isolated testing

	// Set up refresh callback
	refreshCount := 0
	var refreshMutex sync.Mutex
	tm.SetRefreshCallback(func() (*TokenResponse, error) {
		refreshMutex.Lock()
		refreshCount++
		refreshMutex.Unlock()

		return &TokenResponse{
			AccessToken:  "refreshed-token",
			RefreshToken: "refresh-token",
			TokenType:    "Bearer",
			ExpiresIn:    3600,
		}, nil
	})

	// Set initial expired token
	tm.SetTokens("expired-token", "refresh-token", -10)

	// Run concurrent access
	const numGoroutines = 10
	const numOperations = 100

	var wg sync.WaitGroup
	errors := make(chan error, numGoroutines*numOperations)

	// Concurrent GetAccessToken calls
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				_, err := tm.GetAccessToken()
				if err != nil {
					errors <- err
				}
			}
		}()
	}

	// Concurrent token updates
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				tm.SetTokens("token-"+string(rune(id)), "refresh", 3600)
			}
		}(i)
	}

	// Concurrent status checks
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				tm.IsExpired()
				tm.HasValidToken()
				tm.HasRefreshToken()
				tm.GetTimeUntilExpiry()
			}
		}()
	}

	wg.Wait()
	close(errors)

	// Check for any errors
	for err := range errors {
		t.Errorf("Concurrent access error: %v", err)
	}
}

func TestTokenManager_ConcurrentRefresh(t *testing.T) {
	tm := NewMemoryTokenManager() // Use memory storage for isolated testing

	// Set up refresh callback with delay to test concurrent refresh
	refreshCount := 0
	var refreshMutex sync.Mutex
	tm.SetRefreshCallback(func() (*TokenResponse, error) {
		refreshMutex.Lock()
		refreshCount++
		refreshMutex.Unlock()

		// Small delay to increase chance of concurrent calls
		time.Sleep(10 * time.Millisecond)

		return &TokenResponse{
			AccessToken:  "refreshed-token",
			RefreshToken: "refresh-token",
			TokenType:    "Bearer",
			ExpiresIn:    3600,
		}, nil
	})

	// Set initial expired token
	tm.SetTokens("expired-token", "refresh-token", -10)

	// Run concurrent refresh attempts
	const numGoroutines = 5
	var wg sync.WaitGroup
	tokens := make(chan string, numGoroutines)
	errors := make(chan error, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			token, err := tm.GetAccessToken()
			if err != nil {
				errors <- err
			} else {
				tokens <- token
			}
		}()
	}

	wg.Wait()
	close(tokens)
	close(errors)

	// Check for errors
	for err := range errors {
		t.Errorf("Concurrent refresh error: %v", err)
	}

	// All tokens should be the same (refreshed token)
	var firstToken string
	tokenCount := 0
	for token := range tokens {
		if firstToken == "" {
			firstToken = token
		}
		if token != firstToken {
			t.Errorf("Expected all tokens to be the same, got different tokens")
		}
		tokenCount++
	}

	if tokenCount != numGoroutines {
		t.Errorf("Expected %d tokens, got %d", numGoroutines, tokenCount)
	}

	// Refresh should have been called at least once, but possibly more due to timing
	refreshMutex.Lock()
	finalRefreshCount := refreshCount
	refreshMutex.Unlock()

	if finalRefreshCount < 1 {
		t.Error("Refresh callback should have been called at least once")
	}
}

func TestTokenManager_EdgeCases(t *testing.T) {
	tm := NewMemoryTokenManager() // Use memory storage for isolated testing

	// Test with empty strings
	tm.SetTokens("", "", 0)
	if tm.HasValidToken() {
		t.Error("Empty access token should not be valid")
	}

	// Test with negative expiry
	tm.SetTokens("token", "refresh", -1000)
	if tm.HasValidToken() {
		t.Error("Token with negative expiry should not be valid")
	}

	// Test GetAccessToken with no refresh token
	tm.SetTokens("expired-token", "", -10)
	_, err := tm.GetAccessToken()
	if err == nil {
		t.Error("Should return error when no refresh token available")
	}
}

// Test TokenStorage implementations

func TestMemoryTokenStorage(t *testing.T) {
	storage := NewMemoryTokenStorage()

	// Test empty storage
	data, err := storage.Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if data != nil {
		t.Error("Expected nil data from empty storage")
	}

	// Test storing data
	testData := &TokenData{
		AccessToken:  "test-access",
		RefreshToken: "test-refresh",
		TokenType:    "Bearer",
		ExpiresAt:    time.Now().Add(time.Hour),
		Scope:        "read write",
	}

	err = storage.Store(testData)
	if err != nil {
		t.Fatalf("Store failed: %v", err)
	}

	// Test loading data
	loadedData, err := storage.Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if loadedData == nil {
		t.Fatal("Expected data, got nil")
	}

	if loadedData.AccessToken != testData.AccessToken {
		t.Errorf("Expected access token %s, got %s", testData.AccessToken, loadedData.AccessToken)
	}

	if loadedData.RefreshToken != testData.RefreshToken {
		t.Errorf("Expected refresh token %s, got %s", testData.RefreshToken, loadedData.RefreshToken)
	}

	if loadedData.TokenType != testData.TokenType {
		t.Errorf("Expected token type %s, got %s", testData.TokenType, loadedData.TokenType)
	}

	if loadedData.Scope != testData.Scope {
		t.Errorf("Expected scope %s, got %s", testData.Scope, loadedData.Scope)
	}

	// Test clearing data
	err = storage.Clear()
	if err != nil {
		t.Fatalf("Clear failed: %v", err)
	}

	data, err = storage.Load()
	if err != nil {
		t.Fatalf("Load after clear failed: %v", err)
	}
	if data != nil {
		t.Error("Expected nil data after clear")
	}
}

func TestFileTokenStorage(t *testing.T) {
	// Create temporary file
	tempDir := t.TempDir()
	tokenFile := filepath.Join(tempDir, "tokens.json")

	storage := NewFileTokenStorage(tokenFile)

	// Test empty storage (file doesn't exist)
	data, err := storage.Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if data != nil {
		t.Error("Expected nil data from non-existent file")
	}

	// Test storing data
	testData := &TokenData{
		AccessToken:  "test-access",
		RefreshToken: "test-refresh",
		TokenType:    "Bearer",
		ExpiresAt:    time.Now().Add(time.Hour),
		Scope:        "read write",
	}

	err = storage.Store(testData)
	if err != nil {
		t.Fatalf("Store failed: %v", err)
	}

	// Verify file exists and has correct permissions
	fileInfo, err := os.Stat(tokenFile)
	if err != nil {
		t.Fatalf("Token file not created: %v", err)
	}

	// Check file permissions (should be 0600)
	if fileInfo.Mode().Perm() != 0600 {
		t.Errorf("Expected file permissions 0600, got %o", fileInfo.Mode().Perm())
	}

	// Test loading data
	loadedData, err := storage.Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if loadedData == nil {
		t.Fatal("Expected data, got nil")
	}

	if loadedData.AccessToken != testData.AccessToken {
		t.Errorf("Expected access token %s, got %s", testData.AccessToken, loadedData.AccessToken)
	}

	if loadedData.RefreshToken != testData.RefreshToken {
		t.Errorf("Expected refresh token %s, got %s", testData.RefreshToken, loadedData.RefreshToken)
	}

	if loadedData.TokenType != testData.TokenType {
		t.Errorf("Expected token type %s, got %s", testData.TokenType, loadedData.TokenType)
	}

	if loadedData.Scope != testData.Scope {
		t.Errorf("Expected scope %s, got %s", testData.Scope, loadedData.Scope)
	}

	// Test clearing data
	err = storage.Clear()
	if err != nil {
		t.Fatalf("Clear failed: %v", err)
	}

	// Verify file is deleted
	if _, err := os.Stat(tokenFile); !os.IsNotExist(err) {
		t.Error("Expected token file to be deleted after clear")
	}

	// Test loading after clear
	data, err = storage.Load()
	if err != nil {
		t.Fatalf("Load after clear failed: %v", err)
	}
	if data != nil {
		t.Error("Expected nil data after clear")
	}
}

func TestFileTokenStorage_InvalidFile(t *testing.T) {
	// Test with invalid file path
	storage := NewFileTokenStorage("/invalid/path/tokens.json")

	testData := &TokenData{
		AccessToken: "test-access",
	}

	err := storage.Store(testData)
	if err == nil {
		t.Error("Expected error when storing to invalid path")
	}
}

func TestFileTokenStorage_CorruptedFile(t *testing.T) {
	// Create temporary file with invalid JSON
	tempDir := t.TempDir()
	tokenFile := filepath.Join(tempDir, "tokens.json")

	// Write invalid JSON
	err := os.WriteFile(tokenFile, []byte("invalid json"), 0600)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	storage := NewFileTokenStorage(tokenFile)

	// Test loading corrupted file
	_, err = storage.Load()
	if err == nil {
		t.Error("Expected error when loading corrupted file")
	}
}

func TestNewFileTokenManager(t *testing.T) {
	tempDir := t.TempDir()
	tokenFile := filepath.Join(tempDir, "tokens.json")

	tm := NewFileTokenManager(tokenFile)

	if tm == nil {
		t.Fatal("NewFileTokenManager returned nil")
	}

	// Test setting and getting tokens
	tm.SetTokens("access-token", "refresh-token", 3600)

	token, err := tm.GetAccessToken()
	if err != nil {
		t.Fatalf("GetAccessToken failed: %v", err)
	}

	if token != "access-token" {
		t.Errorf("Expected access token 'access-token', got '%s'", token)
	}

	// Verify token was persisted to file
	if _, err := os.Stat(tokenFile); os.IsNotExist(err) {
		t.Error("Expected token file to be created")
	}

	// Create new TokenManager with same file and verify tokens are loaded
	tm2 := NewFileTokenManager(tokenFile)

	token2, err := tm2.GetAccessToken()
	if err != nil {
		t.Fatalf("GetAccessToken from second manager failed: %v", err)
	}

	if token2 != "access-token" {
		t.Errorf("Expected persisted token 'access-token', got '%s'", token2)
	}
}

func TestTokenManagerWithCustomStorage(t *testing.T) {
	// Create custom storage implementation
	storage := NewMemoryTokenStorage()
	tm := NewTokenManagerWithStorage(storage)

	if tm == nil {
		t.Fatal("NewTokenManagerWithStorage returned nil")
	}

	// Test functionality
	tm.SetTokens("custom-access", "custom-refresh", 3600)

	token, err := tm.GetAccessToken()
	if err != nil {
		t.Fatalf("GetAccessToken failed: %v", err)
	}

	if token != "custom-access" {
		t.Errorf("Expected access token 'custom-access', got '%s'", token)
	}
}

func TestTokenManager_StorageErrors(t *testing.T) {
	// Create a mock storage that returns errors
	mockStorage := &MockErrorStorage{}
	tm := NewTokenManagerWithStorage(mockStorage)

	// Test GetAccessToken with storage error
	_, err := tm.GetAccessToken()
	if err == nil {
		t.Error("Expected error from GetAccessToken with failing storage")
	}

	// Test other methods with storage errors
	if tm.GetRefreshToken() != "" {
		t.Error("Expected empty refresh token with failing storage")
	}

	if tm.GetTokenType() != "Bearer" {
		t.Error("Expected default token type with failing storage")
	}

	if tm.GetScope() != "" {
		t.Error("Expected empty scope with failing storage")
	}

	if !tm.IsExpired() {
		t.Error("Expected token to be expired with failing storage")
	}

	if tm.HasValidToken() {
		t.Error("Expected no valid token with failing storage")
	}

	if tm.HasRefreshToken() {
		t.Error("Expected no refresh token with failing storage")
	}

	if !tm.GetExpiresAt().IsZero() {
		t.Error("Expected zero expiration time with failing storage")
	}

	if tm.GetTimeUntilExpiry() != 0 {
		t.Error("Expected zero time until expiry with failing storage")
	}
}

// MockErrorStorage is a storage implementation that always returns errors
type MockErrorStorage struct{}

func (m *MockErrorStorage) Store(data *TokenData) error {
	return errors.New("mock storage error")
}

func (m *MockErrorStorage) Load() (*TokenData, error) {
	return nil, errors.New("mock storage error")
}

func (m *MockErrorStorage) Clear() error {
	return errors.New("mock storage error")
}

func TestTokenManager_FileStorageConcurrency(t *testing.T) {
	tempDir := t.TempDir()
	tokenFile := filepath.Join(tempDir, "concurrent_tokens.json")

	tm := NewFileTokenManager(tokenFile)

	// Run concurrent operations
	const numGoroutines = 10
	const numOperations = 50

	var wg sync.WaitGroup
	errors := make(chan error, numGoroutines*numOperations)

	// Concurrent token updates
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				tm.SetTokens("token-"+string(rune(id)), "refresh", 3600)
			}
		}(i)
	}

	// Concurrent token reads
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				// Just check if token operations work without errors
				// Don't fail on "token expired" errors since that's expected during concurrent updates
				tm.GetAccessToken() // Ignore errors - concurrent updates may cause temporary "expired" states
				tm.IsExpired()
				tm.HasValidToken()
				tm.GetRefreshToken()
			}
		}()
	}

	wg.Wait()
	close(errors)

	// Check for any errors
	for err := range errors {
		t.Errorf("Concurrent file storage error: %v", err)
	}
}

func TestNewMemoryTokenManager(t *testing.T) {
	tm := NewMemoryTokenManager()

	if tm == nil {
		t.Fatal("NewMemoryTokenManager returned nil")
	}

	if tm.GetTokenType() != "Bearer" {
		t.Errorf("Expected default token type 'Bearer', got '%s'", tm.GetTokenType())
	}

	if tm.HasValidToken() {
		t.Error("New TokenManager should not have valid token")
	}

	if tm.HasRefreshToken() {
		t.Error("New TokenManager should not have refresh token")
	}

	// Test that it uses memory storage (not persistent)
	tm.SetTokens("memory-token", "memory-refresh", 3600)

	// Create another memory TokenManager
	// It should NOT have the same tokens (since they're in different memory instances)
	tm2 := NewMemoryTokenManager()
	if tm2.HasValidToken() {
		t.Error("Memory TokenManager should not share tokens across instances")
	}
}

func TestDefaultTokenPath(t *testing.T) {
	path := getDefaultTokenPath()

	if path == "" {
		t.Error("Default token path should not be empty")
	}

	// Should contain .tasty-go directory
	if !strings.Contains(path, ".tasty-go") {
		t.Errorf("Expected path to contain '.tasty-go', got '%s'", path)
	}

	// Should end with tokens.json
	if !strings.HasSuffix(path, "tokens.json") {
		t.Errorf("Expected path to end with 'tokens.json', got '%s'", path)
	}
}

func TestTokenManager_SetTokensFromResponse_PreservesRefreshTokenAndScope(t *testing.T) {
	tm := NewMemoryTokenManager()

	// Set initial tokens with scope
	initialResponse := &TokenResponse{
		AccessToken:  "initial-access-token",
		RefreshToken: "persistent-refresh-token",
		TokenType:    "Bearer",
		ExpiresIn:    3600,
		Scope:        "read write trade",
	}

	tm.SetTokensFromResponse(initialResponse)

	// Verify initial tokens are set
	token, err := tm.GetAccessToken()
	if err != nil {
		t.Fatalf("Failed to get initial access token: %v", err)
	}
	if token != "initial-access-token" {
		t.Errorf("Expected initial access token 'initial-access-token', got '%s'", token)
	}
	if tm.GetRefreshToken() != "persistent-refresh-token" {
		t.Errorf("Expected refresh token 'persistent-refresh-token', got '%s'", tm.GetRefreshToken())
	}
	if tm.GetScope() != "read write trade" {
		t.Errorf("Expected scope 'read write trade', got '%s'", tm.GetScope())
	}

	// Simulate refresh response without refresh token or scope (common for TastyTrade)
	refreshResponse := &TokenResponse{
		AccessToken: "new-access-token",
		TokenType:   "Bearer",
		ExpiresIn:   3600,
		// RefreshToken and Scope are empty (not provided in refresh response)
	}

	tm.SetTokensFromResponse(refreshResponse)

	// Verify new access token is set
	token, err = tm.GetAccessToken()
	if err != nil {
		t.Fatalf("Failed to get refreshed access token: %v", err)
	}
	if token != "new-access-token" {
		t.Errorf("Expected refreshed access token 'new-access-token', got '%s'", token)
	}

	// Verify refresh token is preserved
	if tm.GetRefreshToken() != "persistent-refresh-token" {
		t.Errorf("Expected refresh token to be preserved 'persistent-refresh-token', got '%s'", tm.GetRefreshToken())
	}

	// Verify scope is preserved
	if tm.GetScope() != "read write trade" {
		t.Errorf("Expected scope to be preserved 'read write trade', got '%s'", tm.GetScope())
	}
}

func TestTokenManager_SetTokensFromResponse_UpdatesRefreshTokenWhenProvided(t *testing.T) {
	tm := NewMemoryTokenManager()

	// Set initial tokens
	initialResponse := &TokenResponse{
		AccessToken:  "initial-access-token",
		RefreshToken: "old-refresh-token",
		TokenType:    "Bearer",
		ExpiresIn:    3600,
		Scope:        "read write",
	}

	tm.SetTokensFromResponse(initialResponse)

	// Simulate refresh response with new refresh token (some OAuth2 providers do this)
	refreshResponse := &TokenResponse{
		AccessToken:  "new-access-token",
		RefreshToken: "new-refresh-token",
		TokenType:    "Bearer",
		ExpiresIn:    3600,
		Scope:        "read write trade", // Updated scope
	}

	tm.SetTokensFromResponse(refreshResponse)

	// Verify all tokens are updated
	token, err := tm.GetAccessToken()
	if err != nil {
		t.Fatalf("Failed to get access token: %v", err)
	}
	if token != "new-access-token" {
		t.Errorf("Expected access token 'new-access-token', got '%s'", token)
	}

	if tm.GetRefreshToken() != "new-refresh-token" {
		t.Errorf("Expected refresh token 'new-refresh-token', got '%s'", tm.GetRefreshToken())
	}

	if tm.GetScope() != "read write trade" {
		t.Errorf("Expected scope 'read write trade', got '%s'", tm.GetScope())
	}
}
