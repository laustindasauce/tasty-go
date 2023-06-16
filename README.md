# tasty-go

[![Go Reference](https://pkg.go.dev/badge/github.com/austinbspencer/tasty-go.svg)](https://pkg.go.dev/github.com/austinbspencer/tasty-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/austinbspencer/tasty-go)](https://goreportcard.com/report/github.com/austinbspencer/tasty-go)

This library provides unofficial Go clients for [TastyTrade API](https://tastytrade.com).

> You will need to opt into TastyTrade's API [here](https://developer.tastytrade.com)

## TO-DO

- Separate instruments.go into multiple files
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

</details>
```
