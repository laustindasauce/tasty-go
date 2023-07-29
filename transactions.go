package tasty

import (
	"fmt"
	"net/http"
	"time"
)

// Returns a paginated list of the account's transactions (as identified by
// the provided authentication token) based on sort param. If no sort is
// passed in, it defaults to descending order.
func (c *Client) GetAccountTransactions(accountNumber string, query TransactionsQuery) ([]Transaction, Pagination, *http.Response, error) {
	path := fmt.Sprintf("/accounts/%s/transactions", accountNumber)

	type accountResponse struct {
		Data struct {
			Transactions []Transaction `json:"items"`
		} `json:"data"`
		Pagination Pagination `json:"pagination"`
	}

	transactionsRes := new(accountResponse)

	resp, err := c.request(http.MethodGet, path, query, nil, transactionsRes)
	if err != nil {
		return []Transaction{}, Pagination{}, resp, err
	}

	return transactionsRes.Data.Transactions, transactionsRes.Pagination, resp, nil
}

// Retrieve a transaction by account number and ID.
func (c *Client) GetAccountTransaction(accountNumber string, id int) (Transaction, *http.Response, error) {
	path := fmt.Sprintf("/accounts/%s/transactions/%d", accountNumber, id)

	type accountResponse struct {
		Transaction Transaction `json:"data"`
	}

	transactionsRes := new(accountResponse)

	resp, err := c.request(http.MethodGet, path, nil, nil, transactionsRes)
	if err != nil {
		return Transaction{}, resp, err
	}

	return transactionsRes.Transaction, resp, nil
}

// Return the total fees for an account for a given day
// the day will default to today.
func (c *Client) GetAccountTransactionFees(accountNumber string, date *time.Time) (TransactionFees, *http.Response, error) {
	path := fmt.Sprintf("/accounts/%s/transactions/total-fees", accountNumber)

	type accountResponse struct {
		TransactionFees TransactionFees `json:"data"`
	}

	transactionsRes := new(accountResponse)

	type dateQuery struct {
		// The date to get fees for, defaults to today
		Date string `url:"date,omitempty"`
	}

	query := dateQuery{}

	if date != nil {
		query.Date = date.Format("2006-01-02")
	}

	resp, err := c.request(http.MethodGet, path, query, nil, transactionsRes)
	if err != nil {
		return TransactionFees{}, resp, err
	}

	return transactionsRes.TransactionFees, resp, nil
}
