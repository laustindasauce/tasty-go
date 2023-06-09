package tasty

import (
	"fmt"
	"net/http"

	"github.com/austinbspencer/tasty-go/constants"
	"github.com/austinbspencer/tasty-go/models"
)

func (c *Client) GetMyAccounts() ([]models.Account, *Error) {
	if c.Session.SessionToken == nil {
		return []models.Account{}, &Error{Message: "Session is invalid: Session Token cannot be nil."}
	}
	reqURL := fmt.Sprintf("%s/customers/me/accounts", c.baseURL)

	type accountResponse struct {
		Data struct {
			Items []struct {
				Account models.Account `json:"account"`
			} `json:"items"`
		} `json:"data"`
	}

	accountsRes := new(accountResponse)

	header := http.Header{}
	header.Add("Authorization", *c.Session.SessionToken)

	err := c.get(reqURL, header, nil, accountsRes)
	if err != nil {
		return []models.Account{}, err
	}

	var accounts []models.Account

	for _, acct := range accountsRes.Data.Items {
		accounts = append(accounts, acct.Account)
	}

	return accounts, nil
}

func (c *Client) GetAccountTradingStatus(accountNumber string) (models.AccountTradingStatus, *Error) {
	if c.Session.SessionToken == nil {
		return models.AccountTradingStatus{}, &Error{Message: "Session is invalid: Session Token cannot be nil."}
	}
	reqURL := fmt.Sprintf("%s/accounts/%s/trading-status", c.baseURL, accountNumber)

	type tradingStatusRes struct {
		AccountTradingStatus models.AccountTradingStatus `json:"data"`
	}
	accountsRes := new(tradingStatusRes)

	header := http.Header{}
	header.Add("Authorization", *c.Session.SessionToken)

	err := c.get(reqURL, header, nil, accountsRes)
	if err != nil {
		return models.AccountTradingStatus{}, err
	}

	return accountsRes.AccountTradingStatus, nil
}

func (c *Client) GetAccountBalances(accountNumber string) (models.AccountBalance, *Error) {
	if c.Session.SessionToken == nil {
		return models.AccountBalance{}, &Error{Message: "Session is invalid: Session Token cannot be nil."}
	}
	reqURL := fmt.Sprintf("%s/accounts/%s/balances", c.baseURL, accountNumber)

	type accountBalanceRes struct {
		AccountBalance models.AccountBalance `json:"data"`
	}
	accountsRes := new(accountBalanceRes)

	header := http.Header{}
	header.Add("Authorization", *c.Session.SessionToken)

	err := c.get(reqURL, header, nil, accountsRes)
	if err != nil {
		return models.AccountBalance{}, err
	}

	return accountsRes.AccountBalance, nil
}

func (c *Client) GetAccountPositions(accountNumber string, query models.AccountPositionQuery) ([]models.AccountPosition, *Error) {
	if c.Session.SessionToken == nil {
		return []models.AccountPosition{}, &Error{Message: "Session is invalid: Session Token cannot be nil."}
	}
	reqURL := fmt.Sprintf("%s/accounts/%s/positions", c.baseURL, accountNumber)

	type accountResponse struct {
		Data struct {
			AccountPositions []models.AccountPosition `json:"items"`
		} `json:"data"`
	}

	accountsRes := new(accountResponse)

	header := http.Header{}
	header.Add("Authorization", *c.Session.SessionToken)

	err := c.get(reqURL, header, query, accountsRes)
	if err != nil {
		return []models.AccountPosition{}, err
	}

	return accountsRes.Data.AccountPositions, nil
}

func (c *Client) GetAccountBalanceSnapshots(accountNumber string, query models.AccountBalanceSnapshotsQuery) ([]models.AccountBalanceSnapshots, *Error) {
	if c.Session.SessionToken == nil {
		return []models.AccountBalanceSnapshots{}, &Error{Message: "Session is invalid: Session Token cannot be nil."}
	}

	if query.TimeOfDay == nil {
		query.TimeOfDay = constants.EndOfDay
	}

	reqURL := fmt.Sprintf("%s/accounts/%s/balance-snapshots", c.baseURL, accountNumber)

	type accountResponse struct {
		Data struct {
			AccountBalanceSnapshots []models.AccountBalanceSnapshots `json:"items"`
		} `json:"data"`
	}

	accountsRes := new(accountResponse)

	header := http.Header{}
	header.Add("Authorization", *c.Session.SessionToken)

	err := c.get(reqURL, header, query, accountsRes)
	if err != nil {
		return []models.AccountBalanceSnapshots{}, err
	}

	return accountsRes.Data.AccountBalanceSnapshots, nil
}

func (c *Client) GetAccountNetLiqHistory(accountNumber string, query models.HistoricLiquidityQuery) ([]models.NetLiqOHLC, *Error) {
	if c.Session.SessionToken == nil {
		return []models.NetLiqOHLC{}, &Error{Message: "Session is invalid: Session Token cannot be nil."}
	}

	reqURL := fmt.Sprintf("%s/accounts/%s/net-liq/history", c.baseURL, accountNumber)

	type accountResponse struct {
		Data struct {
			HistoricLiquidity []models.NetLiqOHLC `json:"items"`
		} `json:"data"`
	}

	accountsRes := new(accountResponse)

	header := http.Header{}
	header.Add("Authorization", *c.Session.SessionToken)

	err := c.get(reqURL, header, query, accountsRes)
	if err != nil {
		return []models.NetLiqOHLC{}, err
	}

	return accountsRes.Data.HistoricLiquidity, nil
}
