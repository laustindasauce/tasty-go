package tasty

import (
	"fmt"
	"net/http"
)

// Get the accounts for the authenticated client.
func (c *Client) GetMyAccounts() ([]Account, *Error) {
	if c.Session.SessionToken == nil {
		return []Account{}, &Error{Message: "Session is invalid: Session Token cannot be nil."}
	}
	path := "/customers/me/accounts"

	type accountResponse struct {
		Data struct {
			Items []struct {
				Account Account `json:"account"`
			} `json:"items"`
		} `json:"data"`
	}

	accountsRes := new(accountResponse)

	header := http.Header{}
	header.Add("Authorization", *c.Session.SessionToken)

	err := c.request(http.MethodGet, path, header, nil, nil, accountsRes)
	if err != nil {
		return []Account{}, err
	}

	var accounts []Account

	for _, acct := range accountsRes.Data.Items {
		accounts = append(accounts, acct.Account)
	}

	return accounts, nil
}

// Returns current trading status for an account.
func (c *Client) GetAccountTradingStatus(accountNumber string) (AccountTradingStatus, *Error) {
	if c.Session.SessionToken == nil {
		return AccountTradingStatus{}, &Error{Message: "Session is invalid: Session Token cannot be nil."}
	}
	path := fmt.Sprintf("/accounts/%s/trading-status", accountNumber)

	type tradingStatusRes struct {
		AccountTradingStatus AccountTradingStatus `json:"data"`
	}
	accountsRes := new(tradingStatusRes)

	header := http.Header{}
	header.Add("Authorization", *c.Session.SessionToken)

	err := c.request(http.MethodGet, path, header, nil, nil, accountsRes)
	if err != nil {
		return AccountTradingStatus{}, err
	}

	return accountsRes.AccountTradingStatus, nil
}

// Returns the current balance values for an account.
func (c *Client) GetAccountBalances(accountNumber string) (AccountBalance, *Error) {
	if c.Session.SessionToken == nil {
		return AccountBalance{}, &Error{Message: "Session is invalid: Session Token cannot be nil."}
	}
	path := fmt.Sprintf("/accounts/%s/balances", accountNumber)

	type accountBalanceRes struct {
		AccountBalance AccountBalance `json:"data"`
	}
	accountsRes := new(accountBalanceRes)

	header := http.Header{}
	header.Add("Authorization", *c.Session.SessionToken)

	err := c.request(http.MethodGet, path, header, nil, nil, accountsRes)
	if err != nil {
		return AccountBalance{}, err
	}

	return accountsRes.AccountBalance, nil
}

// Returns a list of the account's positions.
// Can be filtered by symbol, underlying_symbol.
func (c *Client) GetAccountPositions(accountNumber string, query AccountPositionQuery) ([]AccountPosition, *Error) {
	if c.Session.SessionToken == nil {
		return []AccountPosition{}, &Error{Message: "Session is invalid: Session Token cannot be nil."}
	}
	path := fmt.Sprintf("/accounts/%s/positions", accountNumber)

	type accountResponse struct {
		Data struct {
			AccountPositions []AccountPosition `json:"items"`
		} `json:"data"`
	}

	accountsRes := new(accountResponse)

	header := http.Header{}
	header.Add("Authorization", *c.Session.SessionToken)

	err := c.request(http.MethodGet, path, header, query, nil, accountsRes)
	if err != nil {
		return []AccountPosition{}, err
	}

	return accountsRes.Data.AccountPositions, nil
}

// Returns most recent snapshot and current balance for an account.
func (c *Client) GetAccountBalanceSnapshots(accountNumber string, query AccountBalanceSnapshotsQuery) ([]AccountBalanceSnapshots, *Error) {
	if c.Session.SessionToken == nil {
		return []AccountBalanceSnapshots{}, &Error{Message: "Session is invalid: Session Token cannot be nil."}
	}

	// Default to EOD
	if query.TimeOfDay == "" {
		query.TimeOfDay = EndOfDay
	}

	path := fmt.Sprintf("/accounts/%s/balance-snapshots", accountNumber)

	type accountResponse struct {
		Data struct {
			AccountBalanceSnapshots []AccountBalanceSnapshots `json:"items"`
		} `json:"data"`
	}

	accountsRes := new(accountResponse)

	header := http.Header{}
	header.Add("Authorization", *c.Session.SessionToken)

	err := c.request(http.MethodGet, path, header, query, nil, accountsRes)
	if err != nil {
		return []AccountBalanceSnapshots{}, err
	}

	return accountsRes.Data.AccountBalanceSnapshots, nil
}

// Returns a list of account net liquidating value snapshots.
func (c *Client) GetAccountNetLiqHistory(accountNumber string, query HistoricLiquidityQuery) ([]NetLiqOHLC, *Error) {
	if c.Session.SessionToken == nil {
		return []NetLiqOHLC{}, &Error{Message: "Session is invalid: Session Token cannot be nil."}
	}

	path := fmt.Sprintf("/accounts/%s/net-liq/history", accountNumber)

	type accountResponse struct {
		Data struct {
			HistoricLiquidity []NetLiqOHLC `json:"items"`
		} `json:"data"`
	}

	accountsRes := new(accountResponse)

	header := http.Header{}
	header.Add("Authorization", *c.Session.SessionToken)

	err := c.request(http.MethodGet, path, header, query, nil, accountsRes)
	if err != nil {
		return []NetLiqOHLC{}, err
	}

	return accountsRes.Data.HistoricLiquidity, nil
}

// Get the position limit.
func (c *Client) GetAccountPositionLimit(accountNumber string) (PositionLimit, *Error) {
	if c.Session.SessionToken == nil {
		return PositionLimit{}, &Error{Message: "Session is invalid: Session Token cannot be nil."}
	}

	path := fmt.Sprintf("/accounts/%s/position-limit", accountNumber)

	type accountResponse struct {
		PositionLimit PositionLimit `json:"data"`
	}

	accountsRes := new(accountResponse)

	header := http.Header{}
	header.Add("Authorization", *c.Session.SessionToken)

	err := c.request(http.MethodGet, path, header, nil, nil, accountsRes)
	if err != nil {
		return PositionLimit{}, err
	}

	return accountsRes.PositionLimit, nil
}
