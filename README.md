# tasty-go

[![Go Reference](https://pkg.go.dev/badge/github.com/austinbspencer/tasty-go.svg)](https://pkg.go.dev/github.com/austinbspencer/tasty-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/austinbspencer/tasty-go)](https://goreportcard.com/report/github.com/austinbspencer/tasty-go)

This library provides unofficial Go clients for [TastyTrade API](https://tastytrade.com).

## TO-DO

- Separate instruments.go into multiple files
- Reconfirm order is untested as-well-as customers endpoints

## Installation

```
go get github.com/austinbspencer.com/tasty-go
```

## Example Usage

Check tests for other example usages.

```go
import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/austinbspencer/tasty-go"
)

var (
	hClient http.Client = http.Client{Timeout: time.Duration(30) * time.Second}
)

func main() {
	client, err := tasty.NewCertClient(&hClient)
	_, err := Client.CreateSession(tasty.LoginInfo{Login: "username", Password: "password"})
	if err != nil {
		log.Fatal(err)
	}

	res, err := Client.GetMyAccounts()
		if err != nil {
			log.Fatal(err)
		}

	balances, err := Client.GetAccountBalances()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(balances.CashBalance)
}

```
