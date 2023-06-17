# tasty-go

[![Go Reference](https://pkg.go.dev/badge/github.com/austinbspencer/tasty-go.svg)](https://pkg.go.dev/github.com/austinbspencer/tasty-go)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/austinbspencer/tasty-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/austinbspencer/tasty-go)](https://goreportcard.com/report/github.com/austinbspencer/tasty-go)
[![codecov](https://codecov.io/gh/austinbspencer/tasty-go/branch/main/graph/badge.svg?token=ZVVJF2RFQO)](https://codecov.io/gh/austinbspencer/tasty-go)

This library provides unofficial Go clients for [TastyTrade API](https://tastytrade.com).

> You will need to opt into TastyTrade's API [here](https://developer.tastytrade.com)

## TO-DO

- Untested endpoints
  - Margin requirements dry-run
  - Order reconfirm

## Installation

```
go get github.com/austinbspencer.com/tasty-go
```

## Example Usage

Simple usage to get you started.

```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/austinbspencer/tasty-go"
	_ "github.com/joho/godotenv/autoload" // load .env file automatically
)

var (
	hClient = http.Client{Timeout: time.Duration(30) * time.Second}
	client  *tasty.Client
)

var certCreds = tasty.LoginInfo{Login: os.Getenv("certUsername"), Password: os.Getenv("certPassword")}

func main() {
	client, _ = tasty.NewCertClient(&hClient)
	_, err := client.CreateSession(certCreds, nil)
	if err != nil {
		log.Fatal(err)
	}

	accounts, err := client.GetMyAccounts()
	if err != nil {
		log.Fatal(err)
	}

	balances, err := client.GetAccountBalances(accounts[0].AccountNumber)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(balances.CashBalance)
}

```

## Basic Usage:

<details>
<summary>Auth Patterns (Token, session lifetime)</summary>

Check out TastyTrade's [documentation](https://developer.tastytrade.com/#auth-patterns-token-session-lifetime)

- Create / validate / create from remember token

```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/austinbspencer/tasty-go"
	_ "github.com/joho/godotenv/autoload" // load .env file automatically
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

	_, err = client.ValidateSession()
	if err != nil {
		_, err = client.
			CreateSession(tasty.LoginInfo{
				Login:    client.Session.User.Email,
				Password: *client.Session.RememberToken,
			}, nil)
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("Session is valid")

	// Destroy the session
	err = client.DestroySession()
	if err != nil {
		log.Fatal(err)
	}
}

```

</details>

<details>
<summary>User Management</summary>

Check out TastyTrade's [documentation](https://developer.tastytrade.com/#user-management)

> Password Reset

```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/austinbspencer/tasty-go"
	_ "github.com/joho/godotenv/autoload" // load .env file automatically
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

	err = client.RequestPasswordResetEmail(client.Session.User.Email)
	if err != nil {
		log.Fatal(err)
	}

	// You will get an email with a reset link after the above request
	// This link will have a token in the query
	// https://developer.tastytrade.com/password/reset/?token=this-is-your-token

	// Attach the token along with new password in change request
	// Password change will invalidate all current sessions
	err = client.ChangePassword(tasty.PasswordReset{
		Password:             "newPassword",
		PasswordConfirmation: "newPassword",
		ResetPasswordToken:   "this-is-your-token",
	})
	if err != nil {
		log.Fatal(err)
	}
}

```

</details>

<details>
<summary>Customer Account Information</summary>

Check out TastyTrade's [documentation](https://developer.tastytrade.com/#customer-account-information)

```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/austinbspencer/tasty-go"
	_ "github.com/joho/godotenv/autoload" // load .env file automatically
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

	accounts, err := client.GetMyAccounts()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("I have access to %d accounts!", len(accounts))
}

```

</details>

<details>
<summary>Streaming Market Data</summary>

**This requires using the DXFeed Streamer which isn't supported by TastyTrade or this unofficial TastyTrade API wrapper.**

Check out TastyTrade's [documentation](https://developer.tastytrade.com/#streaming-market-data)

```go
package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/austinbspencer/tasty-go"
	_ "github.com/joho/godotenv/autoload" // load .env file automatically
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

<details>
<summary>Position Retrieval</summary>

View all current account positions

Check out TastyTrade's [documentation](https://developer.tastytrade.com/#position-retrieval)

```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/austinbspencer/tasty-go"
	_ "github.com/joho/godotenv/autoload" // load .env file automatically
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

	positions, err := client.GetAccountPositions(accountNumber, tasty.AccountPositionQuery{})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("You have %d positions on your account!", len(positions))
}

```

</details>

<details>
<summary>Balance Retrieval</summary>

Check out TastyTrade's [documentation](https://developer.tastytrade.com/#balance-retrieval)

```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/austinbspencer/tasty-go"
	_ "github.com/joho/godotenv/autoload" // load .env file automatically
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

	balances, err := client.GetAccountBalances(accountNumber)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Your account %s has a cash balance of %f.", balances.AccountNumber, balances.CashBalance)
}

```

</details>

<details>
<summary>Public Watchlists</summary>

Check out TastyTrade's [documentation](https://developer.tastytrade.com/#watchlists)

```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/austinbspencer/tasty-go"
	_ "github.com/joho/godotenv/autoload" // load .env file automatically
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

Check out TastyTrade's [documentation](https://developer.tastytrade.com/#instruments) and [Open API Spec](https://developer.tastytrade.com/open-api-spec/instruments/)

> Equity Options

```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/austinbspencer/tasty-go"
	_ "github.com/joho/godotenv/autoload" // load .env file automatically
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

	"github.com/austinbspencer/tasty-go"
	_ "github.com/joho/godotenv/autoload" // load .env file automatically
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
<summary>Order Management</summary>

Check out TastyTrade's [documentation](https://developer.tastytrade.com/#order-management)

> Live Orders

```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/austinbspencer/tasty-go"
	_ "github.com/joho/godotenv/autoload" // load .env file automatically
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

> Dry Run

```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/austinbspencer/tasty-go"
	_ "github.com/joho/godotenv/autoload" // load .env file automatically
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

> Send Order

```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/austinbspencer/tasty-go"
	_ "github.com/joho/godotenv/autoload" // load .env file automatically
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

	symbol := "RIVN"
	quantity := 1
	action1 := tasty.BTC

	symbol1 := tasty.EquityOptionsSymbology{
		Symbol:     symbol,
		OptionType: tasty.Call,
		Strike:     15,
		Expiration: time.Date(2023, 6, 23, 0, 0, 0, 0, time.Local),
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

	resp, orderErr, err := client.SubmitOrder(accountNumber, order)
	if err != nil {
		log.Fatal(err)
	} else if orderErr != nil {
		log.Fatal(orderErr)
	}

	fmt.Printf("Your order with id: %d has a status of %s!", resp.Order.ID, resp.Order.Status)
}

```

> Cancel Replace

```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/austinbspencer/tasty-go"
	_ "github.com/joho/godotenv/autoload" // load .env file automatically
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
<summary>Example Order Requests</summary>

Check out TastyTrade's [documentation](https://developer.tastytrade.com/#example-order-requests)

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

Check out TastyTrade's [documentation](https://developer.tastytrade.com/#example-order-requests)

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

<details>
<summary>Example Order Requests</summary>
All transactions impacting an accounts balances or positions are available at this endpoint.

Check out TastyTrade's [documentation](https://developer.tastytrade.com/#transaction-history)

```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/austinbspencer/tasty-go"
	_ "github.com/joho/godotenv/autoload" // load .env file automatically
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

	transactions, err := client.GetAccountTransactions(accountNumber, tasty.TransactionsQuery{PerPage: 2})
	if err != nil {
		log.Fatal(err)
	}

	latest := transactions[0]

	fmt.Printf("Your latest transaction was a %s of %s!", latest.TransactionType, latest.UnderlyingSymbol)
}

```

</details>
