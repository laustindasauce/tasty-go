# tasty-go

[![Go Reference](https://pkg.go.dev/badge/github.com/austinbspencer/tasty-go.svg)](https://pkg.go.dev/github.com/austinbspencer/tasty-go)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/austinbspencer/tasty-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/austinbspencer/tasty-go)](https://goreportcard.com/report/github.com/austinbspencer/tasty-go)
[![codecov](https://codecov.io/gh/austinbspencer/tasty-go/branch/main/graph/badge.svg?token=ZVVJF2RFQO)](https://codecov.io/gh/austinbspencer/tasty-go)

This library provides `unofficial` Go clients for [tastytrade API](https://tastytrade.com).

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

## Basic API Usage

Check out tastytrade's [documentation](https://developer.tastytrade.com/basic-api-usage/)

<details>
<summary>Auth Patterns (Token, session lifetime)</summary>

> [docs](https://developer.tastytrade.com/#auth-patterns-token-session-lifetime)

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

> [docs](https://developer.tastytrade.com/basic-api-usage/#user-management)

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

> [docs](https://developer.tastytrade.com/basic-api-usage/#customer-account-information)

```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/austinbspencer/tasty-go"
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

	"github.com/austinbspencer/tasty-go"
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

	"github.com/austinbspencer/tasty-go"
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

	"github.com/austinbspencer/tasty-go"
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

	"github.com/austinbspencer/tasty-go"
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

	"github.com/austinbspencer/tasty-go"
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

	"github.com/austinbspencer/tasty-go"
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

	"github.com/austinbspencer/tasty-go"
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

	// Query for narrowing search of orders
	query := tasty.OrdersQuery{Status: []tasty.OrderStatus{tasty.Filled}}

	orders, _, err := client.GetAccountOrders(accountNumber, query)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Your account has %d live orders!", len(orders))
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

	"github.com/austinbspencer/tasty-go"
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

	"github.com/austinbspencer/tasty-go"
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

	"github.com/austinbspencer/tasty-go"
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

	"github.com/austinbspencer/tasty-go"
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

	"github.com/austinbspencer/tasty-go"
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

	"github.com/austinbspencer/tasty-go"
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

	"github.com/austinbspencer/tasty-go"
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

## Contributing

Please consider opening an [issue](https://github.com/austinbspencer/tasty-go/issues) if you notice any bugs or areas of possible improvement. You can also fork this repo and open a pull request with your own changes. Be sure that all changes have adequate testing in a similar fashion to the rest of the repository.
