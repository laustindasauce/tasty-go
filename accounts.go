package tasty

import (
	"fmt"
	"net/http"

	"github.com/austinbspencer/tasty-go/constants"
	"github.com/austinbspencer/tasty-go/models"
	"github.com/austinbspencer/tasty-go/queries"
)

// Get the accounts for the authenticated client.
func (c *Client) GetMyAccounts() ([]models.Account, *Error) {
	if c.Session.SessionToken == nil {
		return []models.Account{}, &Error{Message: "Session is invalid: Session Token cannot be nil."}
	}
	path := "/customers/me/accounts"

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

	err := c.request(http.MethodGet, path, header, nil, nil, accountsRes)
	if err != nil {
		return []models.Account{}, err
	}

	var accounts []models.Account

	for _, acct := range accountsRes.Data.Items {
		accounts = append(accounts, acct.Account)
	}

	return accounts, nil
}

// Returns current trading status for an account.
func (c *Client) GetAccountTradingStatus(accountNumber string) (models.AccountTradingStatus, *Error) {
	if c.Session.SessionToken == nil {
		return models.AccountTradingStatus{}, &Error{Message: "Session is invalid: Session Token cannot be nil."}
	}
	path := fmt.Sprintf("/accounts/%s/trading-status", accountNumber)

	type tradingStatusRes struct {
		AccountTradingStatus models.AccountTradingStatus `json:"data"`
	}
	accountsRes := new(tradingStatusRes)

	header := http.Header{}
	header.Add("Authorization", *c.Session.SessionToken)

	err := c.request(http.MethodGet, path, header, nil, nil, accountsRes)
	if err != nil {
		return models.AccountTradingStatus{}, err
	}

	return accountsRes.AccountTradingStatus, nil
}

// Returns the current balance values for an account.
func (c *Client) GetAccountBalances(accountNumber string) (models.AccountBalance, *Error) {
	if c.Session.SessionToken == nil {
		return models.AccountBalance{}, &Error{Message: "Session is invalid: Session Token cannot be nil."}
	}
	path := fmt.Sprintf("/accounts/%s/balances", accountNumber)

	type accountBalanceRes struct {
		AccountBalance models.AccountBalance `json:"data"`
	}
	accountsRes := new(accountBalanceRes)

	header := http.Header{}
	header.Add("Authorization", *c.Session.SessionToken)

	err := c.request(http.MethodGet, path, header, nil, nil, accountsRes)
	if err != nil {
		return models.AccountBalance{}, err
	}

	return accountsRes.AccountBalance, nil
}

// Returns a list of the account's positions.
// Can be filtered by symbol, underlying_symbol.
func (c *Client) GetAccountPositions(accountNumber string, query queries.AccountPosition) ([]models.AccountPosition, *Error) {
	if c.Session.SessionToken == nil {
		return []models.AccountPosition{}, &Error{Message: "Session is invalid: Session Token cannot be nil."}
	}
	path := fmt.Sprintf("/accounts/%s/positions", accountNumber)

	type accountResponse struct {
		Data struct {
			AccountPositions []models.AccountPosition `json:"items"`
		} `json:"data"`
	}

	accountsRes := new(accountResponse)

	header := http.Header{}
	header.Add("Authorization", *c.Session.SessionToken)

	err := c.request(http.MethodGet, path, header, query, nil, accountsRes)
	if err != nil {
		return []models.AccountPosition{}, err
	}

	return accountsRes.Data.AccountPositions, nil
}

// Returns most recent snapshot and current balance for an account.
func (c *Client) GetAccountBalanceSnapshots(accountNumber string, query queries.AccountBalanceSnapshots) ([]models.AccountBalanceSnapshots, *Error) {
	if c.Session.SessionToken == nil {
		return []models.AccountBalanceSnapshots{}, &Error{Message: "Session is invalid: Session Token cannot be nil."}
	}

	// Default to EOD
	if query.TimeOfDay == "" {
		query.TimeOfDay = constants.EndOfDay
	}

	path := fmt.Sprintf("/accounts/%s/balance-snapshots", accountNumber)

	type accountResponse struct {
		Data struct {
			AccountBalanceSnapshots []models.AccountBalanceSnapshots `json:"items"`
		} `json:"data"`
	}

	accountsRes := new(accountResponse)

	header := http.Header{}
	header.Add("Authorization", *c.Session.SessionToken)

	err := c.request(http.MethodGet, path, header, query, nil, accountsRes)
	if err != nil {
		return []models.AccountBalanceSnapshots{}, err
	}

	return accountsRes.Data.AccountBalanceSnapshots, nil
}

// Returns a list of account net liquidating value snapshots.
func (c *Client) GetAccountNetLiqHistory(accountNumber string, query queries.HistoricLiquidity) ([]models.NetLiqOHLC, *Error) {
	if c.Session.SessionToken == nil {
		return []models.NetLiqOHLC{}, &Error{Message: "Session is invalid: Session Token cannot be nil."}
	}

	path := fmt.Sprintf("/accounts/%s/net-liq/history", accountNumber)

	type accountResponse struct {
		Data struct {
			HistoricLiquidity []models.NetLiqOHLC `json:"items"`
		} `json:"data"`
	}

	accountsRes := new(accountResponse)

	header := http.Header{}
	header.Add("Authorization", *c.Session.SessionToken)

	err := c.request(http.MethodGet, path, header, query, nil, accountsRes)
	if err != nil {
		return []models.NetLiqOHLC{}, err
	}

	return accountsRes.Data.HistoricLiquidity, nil
}

// Get the position limit.
func (c *Client) GetAccountPositionLimit(accountNumber string) (models.PositionLimit, *Error) {
	if c.Session.SessionToken == nil {
		return models.PositionLimit{}, &Error{Message: "Session is invalid: Session Token cannot be nil."}
	}

	path := fmt.Sprintf("/accounts/%s/position-limit", accountNumber)

	type accountResponse struct {
		PositionLimit models.PositionLimit `json:"data"`
	}

	accountsRes := new(accountResponse)

	header := http.Header{}
	header.Add("Authorization", *c.Session.SessionToken)

	err := c.request(http.MethodGet, path, header, nil, nil, accountsRes)
	if err != nil {
		return models.PositionLimit{}, err
	}

	return accountsRes.PositionLimit, nil
}
