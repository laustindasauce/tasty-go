# tasty-go

[![Go Reference](https://pkg.go.dev/badge/github.com/laustindasauce/tasty-go.svg)](https://pkg.go.dev/github.com/laustindasauce/tasty-go)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/austinbspencer/tasty-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/laustindasauce/tasty-go)](https://goreportcard.com/report/github.com/laustindasauce/tasty-go)
[![codecov](https://codecov.io/gh/austinbspencer/tasty-go/branch/main/graph/badge.svg?token=ZVVJF2RFQO)](https://codecov.io/gh/austinbspencer/tasty-go)

This library provides `unofficial` Go clients for [tastytrade API](https://tastytrade.com).

> **Important:** TastyTrade has migrated to OAuth2 authentication. Session-based authentication is deprecated and will be discontinued on December 1st, 2024. Please migrate to OAuth2 authentication as shown in the examples below.

> You will need to opt into tastytrade's API [here](https://developer.tastytrade.com)

## tastytrade

[tastytrade](https://tastytrade.com/about-us/) pioneered options trading technology for retail traders.

[Create your account](https://start.tastytrade.com/#/login?referralCode=MS53QAT6DS) if you don't already have one to begin trading with tastytrade.

## Dependencies

There are very few direct dependencies for this lightweight API wrapper.

- [decimal](https://github.com/shopspring/decimal)
- [go-querystring](https://github.com/google/go-querystring)
- [testify](https://github.com/stretchr/testify) `for testing`

## Untested endpoints

- Order reconfirm
  - tastytrade API support has informed me that this endpoint is for Equity Offering orders only.

## Installation

```
go get github.com/austinbspencer.com/tasty-go
```

## OAuth2 Authentication Setup

TastyTrade now uses OAuth2 for authentication. The **recommended approach** is to handle OAuth2 authorization in your own application and use this library with pre-existing tokens ("bring your own tokens").

### 1. Register Your Application

First, register your application with TastyTrade to get your OAuth2 credentials:

- Visit [TastyTrade Developer Portal](https://developer.tastytrade.com)
- Create a new application
- Note your `Client ID` and `Client Secret`
- Configure your redirect URI (e.g., `http://localhost:8080` for development)

### 2. Recommended: "Bring Your Own Tokens" Usage

The primary usage pattern is to obtain OAuth2 tokens through your own authorization flow and initialize the client with those tokens:

```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/laustindasauce/tasty-go"
)

func main() {
	// Create OAuth2 configuration
	config := tasty.NewProductionOAuth2Config(
		os.Getenv("TASTY_CLIENT_ID"),
		os.Getenv("TASTY_CLIENT_SECRET"),
		"http://localhost:8080/callback",
		[]string{"read", "trade"},
	)

	// Option 1: Create client with individual token parameters
	// (tokens obtained from your external OAuth2 flow)
	client, err := tasty.NewOAuth2ClientWithTokens(
		config,
		"your-access-token-from-external-flow",
		"your-refresh-token-from-external-flow",
		3600, // expires in 1 hour
		nil,  // use default HTTP client
	)
	if err != nil {
		log.Fatal(err)
	}

	// Option 2: Create client with TokenResponse object
	tokenResponse := &tasty.TokenResponse{
		AccessToken:  "your-access-token",
		RefreshToken: "your-refresh-token",
		TokenType:    "Bearer",
		ExpiresIn:    3600,
		Scope:        "read trade",
	}

	client2, err := tasty.NewOAuth2ClientWithTokenResponse(config, tokenResponse, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Option 3: Set tokens after client creation
	client3, err := tasty.NewOAuth2Client(config, nil)
	if err != nil {
		log.Fatal(err)
	}

	err = client3.SetTokens("access-token", "refresh-token", 3600)
	if err != nil {
		log.Fatal(err)
	}

	// Now use the client for API calls - tokens refresh automatically
	accounts, err := client.GetMyAccounts()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Found %d accounts\n", len(accounts))

	// Check token status
	fmt.Printf("Is authenticated: %v\n", client.IsAuthenticated())
	fmt.Printf("Has valid token: %v\n", client.HasValidToken())
	fmt.Printf("Token expires in: %v\n", client.GetTimeUntilExpiry())
}
```

### 3. Alternative: Built-in OAuth2 Flow (Optional)

```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/laustindasauce/tasty-go"
)

func main() {
	// Configure OAuth2 for sandbox environment
	config := tasty.OAuth2Config{
		ClientID:     os.Getenv("TASTY_CLIENT_ID"),
		ClientSecret: os.Getenv("TASTY_CLIENT_SECRET"),
		RedirectURI:  "http://localhost:8080",
		Scopes:       []string{"read", "trade"},
	}

	// Create OAuth2 client for sandbox
	httpClient := &http.Client{Timeout: 30 * time.Second}
	client, err := tasty.NewOAuth2Client(config, httpClient)
	if err != nil {
		log.Fatal(err)
	}

	// Check if we already have valid tokens
	if client.HasValidToken() {
		fmt.Println("✓ Found existing valid tokens, skipping authentication...")

		// Test API call with existing tokens
		accounts, _, err := client.GetMyAccounts()
		if err != nil {
			fmt.Printf("Existing tokens invalid, need to re-authenticate: %v\n", err)
		} else {
			fmt.Println("✓ Existing tokens work! Making API call...")
			balances, _, err := client.GetAccountBalances(accounts[0].AccountNumber)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Cash balance: $%s\n", balances.CashBalance.String())
			fmt.Println("Authentication not needed - using saved tokens.")
			return
		}
	}

	// Get authorization URL
	authURL, err := client.GetAuthorizationURL()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Visit this URL to authorize: %s\n", authURL)

	// Start built-in redirect server
	server, err := client.StartRedirectServer(8080)
	if err != nil {
		log.Fatal(err)
	}
	defer server.Shutdown(5 * time.Second)

	// Wait for authorization code
	code, err := server.WaitForCode(5 * time.Minute)
	if err != nil {
		log.Fatal(err)
	}

	// Exchange code for tokens
	tokens, err := client.ExchangeCodeForTokens(code)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Access token obtained: %s...\n", tokens.AccessToken[:20])

	// Debug: Show where tokens are stored
	homeDir, _ := os.UserHomeDir()
	tokenPath := fmt.Sprintf("%s/.tasty-go/tokens.json", homeDir)
	fmt.Printf("Tokens stored at: %s\n", tokenPath)

	// Check if token file exists
	if _, err := os.Stat(tokenPath); err == nil {
		fmt.Println("✓ Token file created successfully!")
	} else {
		fmt.Printf("✗ Token file not found: %v\n", err)
	}

	// Now you can make API calls
	accounts, _, err := client.GetMyAccounts()
	if err != nil {
		log.Fatal(err)
	}

	balances, _, err := client.GetAccountBalances(accounts[0].AccountNumber)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Cash balance: $%s\n", balances.CashBalance.String())
}
```

### 4. Production vs Sandbox

For production, use the production constructors:

```go
// Production with tokens
config := tasty.NewProductionOAuth2Config(clientID, clientSecret, redirectURI, scopes)
client, err := tasty.NewOAuth2ClientWithTokens(config, accessToken, refreshToken, expiresIn, nil)

// Sandbox with tokens
config := tasty.NewSandboxOAuth2Config(clientID, clientSecret, redirectURI, scopes)
client, err := tasty.NewCertOAuth2ClientWithTokens(config, accessToken, refreshToken, expiresIn, nil)
```

### 5. Manual Token Exchange (Without Built-in Server)

If you prefer to handle the redirect yourself:

```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/laustindasauce/tasty-go"
)

func main() {
	config := tasty.OAuth2Config{
		ClientID:     os.Getenv("TASTY_CLIENT_ID"),
		ClientSecret: os.Getenv("TASTY_CLIENT_SECRET"),
		RedirectURI:  "https://yourapp.com/callback",
		Scopes:       []string{"read", "trade"},
	}

	httpClient := &http.Client{Timeout: 30 * time.Second}
	client, err := tasty.NewOAuth2Client(config, httpClient)
	if err != nil {
		log.Fatal(err)
	}

	// Get authorization URL
	authURL, err := client.GetAuthorizationURL()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Visit this URL: %s\n", authURL)
	fmt.Print("Enter the authorization code: ")

	var code string
	fmt.Scanln(&code)

	// Validate state parameter (important for security)
	// You should extract this from your callback URL
	var receivedState string
	fmt.Print("Enter the state parameter: ")
	fmt.Scanln(&receivedState)

	if err := client.ValidateState(receivedState); err != nil {
		log.Fatal("Invalid state parameter:", err)
	}

	// Exchange code for tokens
	tokens, err := client.ExchangeCodeForTokens(code)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Successfully authenticated! Token expires in %d seconds\n", tokens.ExpiresIn)

	// Make API calls - tokens are automatically refreshed as needed
	accounts, err := client.GetMyAccounts()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Found %d accounts\n", len(accounts))
}
```

### 6. Token Management and Status

The library provides comprehensive token management methods:

```go
// Check authentication status
isAuth := client.IsAuthenticated()
hasValidToken := client.HasValidToken()
hasRefreshToken := client.HasRefreshToken()
isExpired := client.IsTokenExpired()

// Get token timing information
expiration, err := client.GetTokenExpiration()
timeUntilExpiry, err := client.GetTimeUntilExpiry()

// Update tokens at runtime
err = client.SetTokens("new-access-token", "new-refresh-token", 3600)
err = client.SetTokensFromResponse(newTokenResponse)

// Clear authentication
client.ClearAuthentication()
}
```

## Migration Guide: Session to OAuth2

If you're migrating from session-based authentication, here are the key changes:

### Before (Session-based - Deprecated)

```go
// OLD - Session-based authentication (deprecated)
client, _ := tasty.NewCertClient(&hClient)
creds := tasty.LoginInfo{
    Login:    os.Getenv("username"),
    Password: os.Getenv("password"),
}
_, err := client.CreateSession(creds, nil)
if err != nil {
    log.Fatal(err)
}
```

### After (OAuth2)

```go
// NEW - OAuth2 authentication
config := tasty.OAuth2Config{
    ClientID:     os.Getenv("TASTY_CLIENT_ID"),
    ClientSecret: os.Getenv("TASTY_CLIENT_SECRET"),
    RedirectURI:  "http://localhost:8080",
}
client, err := tasty.NewCertOAuth2Client(config, &hClient)
if err != nil {
    log.Fatal(err)
}

// Handle OAuth2 flow (see examples above)
```

### Key Differences

1. **Authentication Method**: OAuth2 uses authorization codes and tokens instead of username/password
2. **Client Creation**: Use `NewOAuth2Client()` or `NewCertOAuth2Client()` instead of `NewClient()` or `NewCertClient()`
3. **Configuration**: OAuth2 requires client credentials from TastyTrade developer portal
4. **Token Management**: Tokens are automatically refreshed - no manual session management needed
5. **Security**: OAuth2 provides better security with PKCE and state parameters

### Environment Variables

Update your environment variables:

```bash
# Old session-based variables (remove these)
# export certUsername="your_username"
# export certPassword="your_password"

# New OAuth2 variables
export TASTY_CLIENT_ID="your_client_id"
export TASTY_CLIENT_SECRET="your_client_secret"
```

### Common Migration Patterns

#### Pattern 1: Simple API Calls

**Before:**

```go
client, _ := tasty.NewCertClient(&hClient)
client.CreateSession(creds, nil)
accounts, err := client.GetMyAccounts()
```

**After:**

```go
client, _ := tasty.NewCertOAuth2Client(config, &hClient)
// Complete OAuth2 flow (see examples above)
accounts, err := client.GetMyAccounts() // Same API call!
```

#### Pattern 2: Long-running Applications

**Before:**

```go
// Session validation and refresh
_, err := client.ValidateSession()
if err != nil {
    client.CreateSession(creds, nil)
}
```

**After:**

```go
// OAuth2 tokens are automatically refreshed
// No manual validation needed!
accounts, err := client.GetMyAccounts()
// Token refresh happens automatically if needed
```

## Basic API Usage

Check out tastytrade's [documentation](https://developer.tastytrade.com/basic-api-usage/)

<details>
<summary>OAuth2 Token Management</summary>

> OAuth2 tokens are automatically managed - no manual validation needed!

```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/laustindasauce/tasty-go"
)

func main() {
	config := tasty.OAuth2Config{
		ClientID:     os.Getenv("TASTY_CLIENT_ID"),
		ClientSecret: os.Getenv("TASTY_CLIENT_SECRET"),
		RedirectURI:  "http://localhost:8080",
	}

	httpClient := &http.Client{Timeout: 30 * time.Second}
	client, err := tasty.NewCertOAuth2Client(config, httpClient)
	if err != nil {
		log.Fatal(err)
	}

	// Complete OAuth2 flow (see main examples above)
	// ... authorization flow code ...

	// Check authentication status
	if client.IsAuthenticated() {
		fmt.Println("Client is authenticated")
	}

	// Get token information
	tokenManager := client.GetOAuth2Client().GetTokenManager()
	if !tokenManager.IsExpired() {
		timeLeft := tokenManager.GetTimeUntilExpiry()
		fmt.Printf("Token expires in: %v\n", timeLeft)
	}

	// Tokens are automatically refreshed when making API calls
	accounts, err := client.GetMyAccounts()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Successfully retrieved %d accounts\n", len(accounts))

	// Clear authentication when done (optional)
	client.ClearAuthentication()
	fmt.Println("Authentication cleared")
}
```

</details>

<details>
<summary>User Management</summary>

> [docs](https://developer.tastytrade.com/basic-api-usage/#user-management)

> Password Reset (OAuth2)

```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/laustindasauce/tasty-go"
)

func main() {
	config := tasty.OAuth2Config{
		ClientID:     os.Getenv("TASTY_CLIENT_ID"),
		ClientSecret: os.Getenv("TASTY_CLIENT_SECRET"),
		RedirectURI:  "http://localhost:8080",
	}

	httpClient := &http.Client{Timeout: 30 * time.Second}
	client, err := tasty.NewCertOAuth2Client(config, httpClient)
	if err != nil {
		log.Fatal(err)
	}

	// Complete OAuth2 authentication first
	// ... OAuth2 flow code (see main examples) ...

	// Get user information
	customer, err := client.GetMyCustomerInfo()
	if err != nil {
		log.Fatal(err)
	}

	// Request password reset email
	err = client.RequestPasswordResetEmail(customer.Email)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Password reset email sent!")

	// You will get an email with a reset link after the above request
	// This link will have a token in the query
	// https://developer.tastytrade.com/password/reset/?token=this-is-your-token

	// Attach the token along with new password in change request
	// Password change will invalidate all current OAuth2 tokens
	err = client.ChangePassword(tasty.PasswordReset{
		Password:             "newPassword",
		PasswordConfirmation: "newPassword",
		ResetPasswordToken:   "this-is-your-token",
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Password changed successfully!")
	// Note: You'll need to re-authenticate after password change
}
```

</details>

<details>
<summary>Customer Account Information</summary>

> [docs](https://developer.tastytrade.com/basic-api-usage/#customer-account-information)

```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/laustindasauce/tasty-go"
)

func main() {
	config := tasty.OAuth2Config{
		ClientID:     os.Getenv("TASTY_CLIENT_ID"),
		ClientSecret: os.Getenv("TASTY_CLIENT_SECRET"),
		RedirectURI:  "http://localhost:8080",
	}

	httpClient := &http.Client{Timeout: 30 * time.Second}
	client, err := tasty.NewCertOAuth2Client(config, httpClient)
	if err != nil {
		log.Fatal(err)
	}

	// Complete OAuth2 authentication first
	// ... OAuth2 flow code (see main examples) ...

	accounts, err := client.GetMyAccounts()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("I have access to %d accounts!\n", len(accounts))

	// Get detailed customer information
	customer, err := client.GetMyCustomerInfo()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Customer: %s %s\n", customer.FirstName, customer.LastName)
	fmt.Printf("Email: %s\n", customer.Email)
}
```

</details>

<details>
<summary>Account Positions</summary>

View all current account positions

> [docs](https://developer.tastytrade.com/basic-api-usage/#account-positions)

```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/laustindasauce/tasty-go"
)

func main() {
	config := tasty.OAuth2Config{
		ClientID:     os.Getenv("TASTY_CLIENT_ID"),
		ClientSecret: os.Getenv("TASTY_CLIENT_SECRET"),
		RedirectURI:  "http://localhost:8080",
	}

	httpClient := &http.Client{Timeout: 30 * time.Second}
	client, err := tasty.NewCertOAuth2Client(config, httpClient)
	if err != nil {
		log.Fatal(err)
	}

	// Complete OAuth2 authentication first
	// ... OAuth2 flow code (see main examples) ...

	// Get accounts
	accounts, err := client.GetMyAccounts()
	if err != nil {
		log.Fatal(err)
	}

	accountNumber := accounts[0].AccountNumber

	positions, err := client.GetAccountPositions(accountNumber, tasty.AccountPositionQuery{})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("You have %d positions on your account!\n", len(positions))

	// Display position details
	for _, position := range positions {
		fmt.Printf("Symbol: %s, Quantity: %.2f, Market Value: $%.2f\n",
			position.Symbol, position.Quantity, position.MarketValue)
	}
}
```

</details>

<details>
<summary>Account Balances</summary>

> [docs](https://developer.tastytrade.com/basic-api-usage/#account-balances)

```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/laustindasauce/tasty-go"
)

func main() {
	config := tasty.OAuth2Config{
		ClientID:     os.Getenv("TASTY_CLIENT_ID"),
		ClientSecret: os.Getenv("TASTY_CLIENT_SECRET"),
		RedirectURI:  "http://localhost:8080",
	}

	httpClient := &http.Client{Timeout: 30 * time.Second}
	client, err := tasty.NewCertOAuth2Client(config, httpClient)
	if err != nil {
		log.Fatal(err)
	}

	// Complete OAuth2 authentication first
	// ... OAuth2 flow code (see main examples) ...

	// Get accounts
	accounts, err := client.GetMyAccounts()
	if err != nil {
		log.Fatal(err)
	}

	accountNumber := accounts[0].AccountNumber

	balances, err := client.GetAccountBalances(accountNumber)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Account %s balances:\n", balances.AccountNumber)
	fmt.Printf("  Cash Balance: $%.2f\n", balances.CashBalance)
	fmt.Printf("  Net Liquidating Value: $%.2f\n", balances.NetLiquidatingValue)
	fmt.Printf("  Buying Power: $%.2f\n", balances.BuyingPower)
}
```

</details>

<details>
<summary>Watchlists</summary>

> [docs](https://developer.tastytrade.com/basic-api-usage/#watchlists)

> Public Watchlists

```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/laustindasauce/tasty-go"
)

var (
	hClient   = http.Client{Timeout: time.Duration(30) * time.Second}
	certCreds = tasty.LoginInfo{
		Login:      os.Getenv("certUsername"),
		Password:   os.Getenv("certPassword"),
		RememberMe: true,
	}
)

const accountNumber = "5WV48989"

func main() {
	client, _ := tasty.NewCertClient(&hClient)
	_, err := client.CreateSession(certCreds, nil)
	if err != nil {
		log.Fatal(err)
	}

	countsOnly := false

	watchlists, err := client.GetPublicWatchlists(countsOnly)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("There are %d public watchlists!", len(watchlists))
}

```

</details>

<details>
<summary>Instruments</summary>

> [docs](https://developer.tastytrade.com/basic-api-usage/#instruments) and [Open API Spec](https://developer.tastytrade.com/open-api-spec/instruments/)

> Equity Options

```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/laustindasauce/tasty-go"
)

var (
	hClient   = http.Client{Timeout: time.Duration(30) * time.Second}
	certCreds = tasty.LoginInfo{
		Login:      os.Getenv("certUsername"),
		Password:   os.Getenv("certPassword"),
		RememberMe: true,
	}
)

const accountNumber = "5WV48989"

func main() {
	client, _ := tasty.NewCertClient(&hClient)
	_, err := client.CreateSession(certCreds, nil)
	if err != nil {
		log.Fatal(err)
	}

	eoSymbol := tasty.EquityOptionsSymbology{
		Symbol:     "AMD",
		OptionType: tasty.Call,
		Strike:     180,
		Expiration: time.Date(2023, 06, 23, 0, 0, 0, 0, time.UTC),
	}

	equityOptions, err := client.GetEquityOptions(tasty.EquityOptionsQuery{Symbols: []string{eoSymbol.Build()}})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Your equity option with underlying symbol: %s", equityOptions[0].UnderlyingSymbol)
}

```

> Future Options

```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/laustindasauce/tasty-go"
)

var (
	hClient   = http.Client{Timeout: time.Duration(30) * time.Second}
	certCreds = tasty.LoginInfo{
		Login:      os.Getenv("certUsername"),
		Password:   os.Getenv("certPassword"),
		RememberMe: true,
	}
)

const accountNumber = "5WV48989"

func main() {
	client, _ := tasty.NewCertClient(&hClient)
	_, err := client.CreateSession(certCreds, nil)
	if err != nil {
		log.Fatal(err)
	}

	future := tasty.FutureSymbology{ProductCode: "ES", MonthCode: tasty.December, YearDigit: 9}

	expiry := time.Date(2019, 9, 27, 0, 0, 0, 0, time.Local)
	fcc := tasty.FutureOptionsSymbology{
		OptionContractCode: "EW4U9",
		FutureContractCode: future.Build(),
		OptionType:         tasty.Put,
		Strike:             2975,
		Expiration:         expiry,
	}

	query := tasty.FutureOptionsQuery{
		Symbols: []string{fcc.Build()},
	}

	futureOptions, err := client.GetFutureOptions(query)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Your future option with underlying symbol: %s", futureOptions[0].UnderlyingSymbol)
}

```

</details>

<details>
<summary>Transaction History</summary>
All transactions impacting an accounts balances or positions are available at this endpoint.

> [docs](https://developer.tastytrade.com/basic-api-usage/#transaction-history)

```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/laustindasauce/tasty-go"
)

var (
	hClient   = http.Client{Timeout: time.Duration(30) * time.Second}
	certCreds = tasty.LoginInfo{
		Login:      os.Getenv("certUsername"),
		Password:   os.Getenv("certPassword"),
		RememberMe: true,
	}
)

const accountNumber = "5WV48989"

func main() {
	client, _ := tasty.NewCertClient(&hClient)
	_, err := client.CreateSession(certCreds, nil)
	if err != nil {
		log.Fatal(err)
	}

	transactions, _, err := client.GetAccountTransactions(accountNumber, tasty.TransactionsQuery{PerPage: 2})
	if err != nil {
		log.Fatal(err)
	}

	latest := transactions[0]

	fmt.Printf("Your latest transaction was a %s of %s!", latest.TransactionType, latest.UnderlyingSymbol)
}

```

> With Pagination Handling

```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/laustindasauce/tasty-go"
)

var (
	hClient   = http.Client{Timeout: time.Duration(30) * time.Second}
	certCreds = tasty.LoginInfo{
		Login:      os.Getenv("certUsername"),
		Password:   os.Getenv("certPassword"),
		RememberMe: true,
	}
)

const accountNumber = "5WV48989"

func main() {
	client, _ := tasty.NewCertClient(&hClient)
	_, err := client.CreateSession(certCreds, nil)
	if err != nil {
		log.Fatal(err)
	}

	query := tasty.TransactionsQuery{PerPage: 25}

	transactions, pagination, err := client.GetAccountTransactions(accountNumber, query)
	if err != nil {
		log.Fatal(err)
	}

	for pagination.PageOffset < (pagination.TotalPages - 1) {
		query.PageOffset += 1
		moreTransactions, newPagination, err := client.GetAccountTransactions(accountNumber, query)
		if err != nil {
			log.Fatal(err)
		}

		transactions = append(transactions, moreTransactions...)
		pagination = newPagination
	}

	latest := transactions[0]

	fmt.Printf("Your latest transaction was a %s of %s!", latest.TransactionType, latest.UnderlyingSymbol)
}

```

</details>

## Order Management

Check out tastytrade's [documentation](https://developer.tastytrade.com/order-management/)

<details>
<summary>Search Orders</summary>

> [docs](https://developer.tastytrade.com/order-management/#search-orders)

```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/laustindasauce/tasty-go"
)

func main() {
	config := tasty.OAuth2Config{
		ClientID:     os.Getenv("TASTY_CLIENT_ID"),
		ClientSecret: os.Getenv("TASTY_CLIENT_SECRET"),
		RedirectURI:  "http://localhost:8080",
	}

	httpClient := &http.Client{Timeout: 30 * time.Second}
	client, err := tasty.NewCertOAuth2Client(config, httpClient)
	if err != nil {
		log.Fatal(err)
	}

	// Complete OAuth2 authentication first
	// ... OAuth2 flow code (see main examples) ...

	// Get accounts
	accounts, err := client.GetMyAccounts()
	if err != nil {
		log.Fatal(err)
	}

	accountNumber := accounts[0].AccountNumber

	// Query for narrowing search of orders
	query := tasty.OrdersQuery{Status: []tasty.OrderStatus{tasty.Filled}}

	orders, _, err := client.GetAccountOrders(accountNumber, query)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Your account has %d filled orders!\n", len(orders))

	// Display order details
	for _, order := range orders {
		fmt.Printf("Order ID: %d, Status: %s, Symbol: %s\n",
			order.ID, order.Status, order.Legs[0].Symbol)
	}
}
```

</details>

<details>
<summary>Search Orders</summary>

> [docs](https://developer.tastytrade.com/order-management/#live-orders)

```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/laustindasauce/tasty-go"
)

var (
	hClient   = http.Client{Timeout: time.Duration(30) * time.Second}
	certCreds = tasty.LoginInfo{
		Login:      os.Getenv("certUsername"),
		Password:   os.Getenv("certPassword"),
		RememberMe: true,
	}
)

const accountNumber = "5WV48989"

func main() {
	client, _ := tasty.NewCertClient(&hClient)
	_, err := client.CreateSession(certCreds, nil)
	if err != nil {
		log.Fatal(err)
	}

	liveOrders, err := client.GetAccountLiveOrders(accountNumber)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Your account has %d live orders!", len(liveOrders))
}

```

</details>

<details>
<summary>Order Dry Run</summary>

> [docs](https://developer.tastytrade.com/order-management/#order-dry-run)

```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/laustindasauce/tasty-go"
)

var (
	hClient   = http.Client{Timeout: time.Duration(30) * time.Second}
	certCreds = tasty.LoginInfo{
		Login:      os.Getenv("certUsername"),
		Password:   os.Getenv("certPassword"),
		RememberMe: true,
	}
)

const accountNumber = "5WV48989"

func main() {
	client, _ := tasty.NewCertClient(&hClient)
	_, err := client.CreateSession(certCreds, nil)
	if err != nil {
		log.Fatal(err)
	}

	symbol := "AMD"
	quantity := 1
	action := tasty.BTO

	order := tasty.NewOrder{
		TimeInForce: tasty.Day,
		OrderType:   tasty.Market,
		Legs: []tasty.NewOrderLeg{
			{
				InstrumentType: tasty.EquityIT,
				Symbol:         symbol,
				Quantity:       quantity,
				Action:         action,
			},
		},
	}

	resp, orderErr, err := client.SubmitOrderDryRun(accountNumber, order)
	if err != nil {
		log.Fatal(err)
	} else if orderErr != nil {
		log.Fatal(orderErr)
	}

	fmt.Printf("Your dry run order status is %s!", resp.Order.Status)
}

```

</details>

<details>
<summary>Submit Order</summary>

> [docs](https://developer.tastytrade.com/order-management/#submit-order)

```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/laustindasauce/tasty-go"
)

func main() {
	config := tasty.OAuth2Config{
		ClientID:     os.Getenv("TASTY_CLIENT_ID"),
		ClientSecret: os.Getenv("TASTY_CLIENT_SECRET"),
		RedirectURI:  "http://localhost:8080",
	}

	httpClient := &http.Client{Timeout: 30 * time.Second}
	client, err := tasty.NewCertOAuth2Client(config, httpClient)
	if err != nil {
		log.Fatal(err)
	}

	// Complete OAuth2 authentication first
	// ... OAuth2 flow code (see main examples) ...

	// Get accounts
	accounts, err := client.GetMyAccounts()
	if err != nil {
		log.Fatal(err)
	}

	accountNumber := accounts[0].AccountNumber

	symbol := "RIVN"
	quantity := 1
	action1 := tasty.BTC

	// Create option symbol for expiration date
	expirationDate := time.Now().AddDate(0, 1, 0) // 1 month from now
	symbol1 := tasty.EquityOptionsSymbology{
		Symbol:     symbol,
		OptionType: tasty.Call,
		Strike:     15,
		Expiration: expirationDate,
	}

	order := tasty.NewOrder{
		TimeInForce: tasty.GTC,
		OrderType:   tasty.Limit,
		PriceEffect: tasty.Debit,
		Price:       0.04,
		Legs: []tasty.NewOrderLeg{
			{
				InstrumentType: tasty.EquityOptionIT,
				Symbol:         symbol1.Build(),
				Quantity:       quantity,
				Action:         action1,
			},
		},
		Rules: tasty.NewOrderRules{Conditions: []tasty.NewOrderCondition{
			{
				Action:         tasty.Route,
				Symbol:         symbol,
				InstrumentType: "Equity",
				Indicator:      tasty.Last,
				Comparator:     tasty.LTE,
				Threshold:      0.01,
			},
		}},
	}

	// Submit order dry run first (recommended)
	dryRunResp, orderErr, err := client.SubmitOrderDryRun(accountNumber, order)
	if err != nil {
		log.Fatal(err)
	} else if orderErr != nil {
		log.Fatal("Dry run failed:", orderErr)
	}

	fmt.Printf("Dry run successful! Estimated cost: $%.2f\n", dryRunResp.Order.Price)

	// Submit actual order
	resp, orderErr, err := client.SubmitOrder(accountNumber, order)
	if err != nil {
		log.Fatal(err)
	} else if orderErr != nil {
		log.Fatal("Order submission failed:", orderErr)
	}

	fmt.Printf("Order submitted successfully!\n")
	fmt.Printf("Order ID: %d\n", resp.Order.ID)
	fmt.Printf("Status: %s\n", resp.Order.Status)
	fmt.Printf("Symbol: %s\n", resp.Order.Legs[0].Symbol)
}
```

</details>

<details>
<summary>Cancel Order</summary>

> [docs](https://developer.tastytrade.com/order-management/#cancel-order)

```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/laustindasauce/tasty-go"
)

var (
	hClient   = http.Client{Timeout: time.Duration(30) * time.Second}
	certCreds = tasty.LoginInfo{
		Login:      os.Getenv("certUsername"),
		Password:   os.Getenv("certPassword"),
		RememberMe: true,
	}
)

const accountNumber = "5WV48989"
const orderID = 123456

func main() {
	client, _ := tasty.NewCertClient(&hClient)
	_, err := client.CreateSession(certCreds, nil)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := client.CancelOrder(accountNumber, orderID); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Order has been cancelled!")
}

```

</details>

<details>
<summary>Cancel Replace</summary>

> [docs](https://developer.tastytrade.com/order-management/#cancel-replace)

```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/laustindasauce/tasty-go"
)

var (
	hClient   = http.Client{Timeout: time.Duration(30) * time.Second}
	certCreds = tasty.LoginInfo{
		Login:      os.Getenv("certUsername"),
		Password:   os.Getenv("certPassword"),
		RememberMe: true,
	}
)

const accountNumber = "5WV48989"

func main() {
	client, _ := tasty.NewCertClient(&hClient)
	_, err := client.CreateSession(certCreds, nil)
	if err != nil {
		log.Fatal(err)
	}

	orderID := 68678

	orderECR := tasty.NewOrderECR{
		TimeInForce: tasty.Day,
		Price:       185.45,
		OrderType:   tasty.Limit,
		PriceEffect: tasty.Debit,
		ValueEffect: tasty.Debit,
	}

	newOrder, err := client.ReplaceOrder(accountNumber, orderID, orderECR)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Your order was replaced with order with id: %d has a status of %s!", newOrder.ID, newOrder.Status)
}

```

</details>

<details>
<summary>Examples</summary>

> [docs](https://developer.tastytrade.com/order-management/#example-order-requests)

> Market Order

```go
order := tasty.NewOrder{
	TimeInForce: tasty.Day,
	OrderType:   tasty.Market,
	Legs: []tasty.NewOrderLeg{
		{
			InstrumentType: tasty.EquityIT,
			Symbol: "AMD",
			Quantity: 1,
			Action: tasty.BTO,
		},
	},
}
```

> GTC Closing Order

```go
order := tasty.NewOrder{
	TimeInForce: tasty.GTC,
	Price: 150.25,
	PriceEffect: tasty.Credit,
	OrderType:   tasty.Limit,
	Legs: []tasty.NewOrderLeg{
		{
			InstrumentType: tasty.EquityIT,
			Symbol: "AMD",
			Quantity: 1,
			Action: tasty.STC,
		},
	},
}
```

> Short Futures Limit Order

```go
order := tasty.NewOrder{
	TimeInForce: tasty.Day,
	Price: 90.03,
	PriceEffect: tasty.Credit,
	OrderType:   tasty.Limit,
	Legs: []tasty.NewOrderLeg{
		{
			InstrumentType: tasty.FutureIT,
			Symbol: "/CLZ2",
			Quantity: 1,
			Action: tasty.STO,
		},
	},
}
```

> Bear Call Spread

```go
eoSymbolShort := tasty.EquityOptionsSymbology{
	Symbol:     "AMD",
	OptionType: tasty.Call,
	Strike:     185,
	Expiration: time.Date(2023, 06, 23, 0, 0, 0, 0, time.UTC),
}

eoSymbolLong := tasty.EquityOptionsSymbology{
	Symbol:     "AMD",
	OptionType: tasty.Call,
	Strike:     187.5,
	Expiration: time.Date(2023, 06, 23, 0, 0, 0, 0, time.UTC),
}

order := tasty.NewOrder{
	TimeInForce: tasty.Day,
	Price:       0.85,
	PriceEffect: tasty.Credit,
	OrderType:   tasty.Limit,
	Legs: []tasty.NewOrderLeg{
		{
			InstrumentType: tasty.EquityOptionIT,
			Symbol:         eoSymbolShort.Build(),
			Quantity:       1,
			Action:         tasty.STO,
		},
		{
			InstrumentType: tasty.EquityOptionIT,
			Symbol:         eoSymbolLong.Build(),
			Quantity:       1,
			Action:         tasty.BTO,
		},
	},
}
```

> GTD Order

```go
order := tasty.NewOrder{
	TimeInForce: tasty.GTD,
	GtcDate:     "2023-06-23",
	Price:       0.85,
	PriceEffect: tasty.Credit,
	OrderType:   tasty.Limit,
	Legs: []tasty.NewOrderLeg{
		{
			InstrumentType: tasty.EquityIT,
			Symbol:         "AMD",
			Quantity:       1,
			Action:         tasty.BTO,
		},
	},
}
```

> Stop Limit Order

```go
order := tasty.NewOrder{
	TimeInForce: tasty.Day,
	Price:       180.0,
	PriceEffect: tasty.Debit,
	OrderType:   tasty.Limit,
	StopTrigger: 180.0,
	Legs: []tasty.NewOrderLeg{
		{
			InstrumentType: tasty.EquityIT,
			Symbol:         "AMD",
			Quantity:       1,
			Action:         tasty.BTO,
		},
	},
}
```

> Notional Cryptocurrency Order

```go
order := tasty.NewOrder{
	TimeInForce: tasty.GTC,
	OrderType:   tasty.NotionalMarket,
	Value:       10.0,
	ValueEffect: tasty.Debit,
	Legs: []tasty.NewOrderLeg{
		{
			InstrumentType: tasty.Crypto,
			Symbol:         string(tasty.Bitcoin),
			Action:         tasty.BTO,
		},
	},
}
```

</details>

<details>
<summary>Example Order Requests</summary>

**Tastytrade only supports fractional trading of certain equity products.**

- To determine if an equity can be fractionally traded, fetch the equity instrument and check the is-fractional-quantity-eligible field

Check out tastytrade's [documentation](https://developer.tastytrade.com/order-management/#example-order-requests)

> Fractional Quantity Order

```go
// Fractional orders must have a minimum monetary value of $5.
// Buy orders for 0.5 shares of a $1 stock will be rejected.
order := tasty.NewOrder{
	TimeInForce: tasty.Day,
	OrderType:   tasty.Market,
	Legs: []tasty.NewOrderLeg{
		{
			InstrumentType: tasty.EquityIT,
			Symbol:         "AMD",
			Quantity:       0.5,
			Action:         tasty.BTO,
		},
	},
}
```

> Notional Amount Order

```go
// To buy $10 of AMD stock, submit a Notional Market order with a value
// instead of a price. Omit the quantity field from the legs:
order := tasty.NewOrder{
	TimeInForce: tasty.Day,
	OrderType:   tasty.NotionalMarket,
	Value: 10.0,
	ValueEffect: tasty.Debit,
	Legs: []tasty.NewOrderLeg{
		{
			InstrumentType: tasty.EquityIT,
			Symbol:         "AMD",
			Action:         tasty.BTO,
		},
	},
}
```

</details>

## Streaming Market Data

Check out tastytrade's [documentation](https://developer.tastytrade.com/streaming-market-data/)

<details>
<summary>Get a Streamer Token</summary>

**This requires using the DXFeed Streamer which isn't supported by tastytrade or this unofficial tastytrade API wrapper.**

Check out tastytrade's [documentation](https://developer.tastytrade.com/streaming-market-data)

```go
package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/laustindasauce/tasty-go"
)

var (
	hClient   = http.Client{Timeout: time.Duration(30) * time.Second}
	certCreds = tasty.LoginInfo{
		Login:      os.Getenv("certUsername"),
		Password:   os.Getenv("certPassword"),
		RememberMe: true,
	}
)

const accountNumber = "5WV48989"

func main() {
	client, _ := tasty.NewCertClient(&hClient)
	_, err := client.CreateSession(certCreds, nil)
	if err != nil {
		log.Fatal(err)
	}

	dxFeedData, err := client.GetQuoteStreamerTokens()
	if err != nil {
		log.Fatal(err)
	}

	// Do something with the streamer data
}

```

</details>

## Streaming Account Data

Check out tastytrade's [documentation](https://developer.tastytrade.com/streaming-account-data/)

<details>
<summary>Simple Websocket Account Streamer</summary>

**This is an oversimplified websocket connection example for streaming account data**

```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/laustindasauce/tasty-go"
	"golang.org/x/net/websocket"
)

var (
	hClient   = http.Client{Timeout: time.Duration(30) * time.Second}
	certCreds = tasty.LoginInfo{
		Login:      os.Getenv("certUsername"),
		Password:   os.Getenv("certPassword"),
		RememberMe: true,
	}
)

const accountNumber = "5WV48989"

func main() {
	client := tasty.NewCertClient(&hClient)
	_, _, err := client.CreateSession(certCreds, nil)
	if err != nil {
		log.Fatal(err)
	}

	protocol := ""
	origin := "http://localhost:8080"

	// Open Websocket connection
	ws, err := websocket.Dial(client.GetWebsocketURL(), protocol, origin)
	if err != nil {
		log.Fatal(err)
	}

	incomingMessages := make(chan string)
	go readClientMessages(ws, incomingMessages)

	// Send connect message
	response := new(tasty.WebsocketMessage)
	response.Action = "connect"
	response.Value = []string{accountNumber}
	response.AuthToken = *client.Session.SessionToken
	err = websocket.JSON.Send(ws, response)
	if err != nil {
		fmt.Printf("Send failed: %s\n", err.Error())
		os.Exit(1)
	}

	// Subscribe to notifications
	// Add notification subscription message here
	// All available -> https://developer.tastytrade.com/streaming-account-data/#available-actions

	// Await responses and send heartbeats
	i := 0
	for {
		select {
		case <-time.After(time.Duration(time.Second * 15)):
			// Send heartbeat every 15 seconds to keep connection alive
			fmt.Println("sending heartbeat")
			i++
			response := new(tasty.WebsocketMessage)
			response.Action = "heartbeat"
			response.AuthToken = *client.Session.SessionToken
			err = websocket.JSON.Send(ws, response)
			if err != nil {
				fmt.Printf("Send failed: %s\n", err.Error())
				os.Exit(1)
			}
		case message := <-incomingMessages:
			fmt.Println(`Message Received:`, message)
		}
	}
}

func readClientMessages(ws *websocket.Conn, incomingMessages chan string) {
	for {
		var message string
		err := websocket.Message.Receive(ws, &message)
		if err != nil {
			fmt.Printf("Error::: %s\n", err.Error())
			return
		}
		incomingMessages <- message
	}
}

```

</details>

## Testing

Nearly 100% code coverage testing.

> Run all tests

```bash
go test .
```

> Run all tests with code coverage information

```bash
go test -race -covermode=atomic -coverprofile=coverage.out -v .
```

## OAuth2 Troubleshooting

### Common Issues and Solutions

#### 1. "Invalid client credentials" Error

**Problem:** Your client ID or client secret is incorrect.

**Solution:**

- Verify your credentials in the TastyTrade developer portal
- Ensure you're using the correct environment (production vs sandbox)
- Check that your environment variables are set correctly

```bash
echo $TASTY_CLIENT_ID
echo $TASTY_CLIENT_SECRET
```

#### 2. "Invalid redirect URI" Error

**Problem:** The redirect URI doesn't match what's registered with TastyTrade.

**Solution:**

- Ensure the redirect URI in your code exactly matches the one registered in the developer portal
- For development, use `http://localhost:8080` (HTTP is allowed for localhost)
- For production, use HTTPS URLs only

#### 3. "Invalid state parameter" Error

**Problem:** State parameter mismatch, which could indicate a CSRF attack or implementation error.

**Solution:**

- Ensure you're properly validating the state parameter
- Don't manually modify the state parameter
- Make sure the state from the authorization URL matches the one in the callback

```go
// Always validate state parameter
if err := client.ValidateState(receivedState); err != nil {
    log.Fatal("Invalid state parameter:", err)
}
```

#### 4. "Token expired" Error

**Problem:** Access token has expired and refresh failed.

**Solution:**

- Tokens are automatically refreshed - this usually indicates a refresh token issue
- Re-authenticate the user through the OAuth2 flow
- Check that your refresh token hasn't been revoked

```go
// Check if client is still authenticated
if !client.IsAuthenticated() {
    // Need to re-authenticate
    // ... perform OAuth2 flow again ...
}
```

#### 5. "Authorization code expired" Error

**Problem:** Too much time passed between getting the authorization code and exchanging it for tokens.

**Solution:**

- Exchange the authorization code for tokens immediately after receiving it
- Authorization codes typically expire within 10 minutes
- Don't store authorization codes - exchange them right away

#### 6. Network/Connection Issues

**Problem:** Network timeouts or connection errors during OAuth2 flow.

**Solution:**

- Increase HTTP client timeout
- Implement retry logic for network errors
- Check your internet connection and firewall settings

```go
// Increase timeout for OAuth2 operations
httpClient := &http.Client{
    Timeout: 60 * time.Second, // Increased timeout
}
```

#### 7. "Server temporarily unavailable" Error

**Problem:** TastyTrade servers are experiencing issues.

**Solution:**

- Wait and retry after a few minutes
- Check TastyTrade's status page for known issues
- Implement exponential backoff for retries

### Environment-Specific Issues

#### Sandbox vs Production

Make sure you're using the correct client constructor:

```go
// For sandbox/testing
client, err := tasty.NewCertOAuth2Client(config, httpClient)

// For production
client, err := tasty.NewOAuth2Client(config, httpClient)
```

#### HTTPS Requirements

- Production OAuth2 endpoints require HTTPS
- Redirect URIs must use HTTPS in production (except localhost for development)
- Ensure your callback server uses HTTPS in production

### Debugging Tips

#### Enable Detailed Logging

```go
// Add detailed error logging
if err != nil {
    if oauthErr, ok := err.(*tasty.OAuth2DetailedError); ok {
        log.Printf("OAuth2 Error: %s", oauthErr.Error())
        log.Printf("Error Type: %s", oauthErr.GetTypeString())
        log.Printf("Severity: %s", oauthErr.GetSeverityString())
        if oauthErr.InternalMessage != "" {
            log.Printf("Internal: %s", oauthErr.InternalMessage)
        }
    } else {
        log.Printf("General Error: %s", err.Error())
    }
}
```

#### Check Token Status

```go
tokenManager := client.GetOAuth2Client().GetTokenManager()
fmt.Printf("Token expired: %v\n", tokenManager.IsExpired())
fmt.Printf("Has refresh token: %v\n", tokenManager.HasRefreshToken())
fmt.Printf("Time until expiry: %v\n", tokenManager.GetTimeUntilExpiry())
```

#### Validate Configuration

```go
config := tasty.OAuth2Config{
    ClientID:     os.Getenv("TASTY_CLIENT_ID"),
    ClientSecret: os.Getenv("TASTY_CLIENT_SECRET"),
    RedirectURI:  "http://localhost:8080",
}

if err := config.Validate(); err != nil {
    log.Fatal("Invalid configuration:", err)
}
```

### Getting Help

If you're still experiencing issues:

1. Check the [TastyTrade Developer Documentation](https://developer.tastytrade.com)
2. Review the OAuth2 specification: [RFC 6749](https://tools.ietf.org/html/rfc6749)
3. Open an issue on this repository with:
   - Your Go version
   - The exact error message
   - A minimal code example (without credentials)
   - Whether you're using sandbox or production

## Contributing

Please consider opening an [issue](https://github.com/laustindasauce/tasty-go/issues) if you notice any bugs or areas of possible improvement. You can also fork this repo and open a pull request with your own changes. Be sure that all changes have adequate testing in a similar fashion to the rest of the repository.
