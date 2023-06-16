package tasty

import (
	"fmt"
	"net/http"
	"time"
)

// Returns a paginated list of the account's transactions (as identified by
// the provided authentication token) based on sort param. If no sort is
// passed in, it defaults to descending order.
func (c *Client) GetAccountTransactions(accountNumber string, query TransactionsQuery) ([]Transaction, *Error) {
	if c.Session.SessionToken == nil {
		return []Transaction{}, &Error{Message: "Session is invalid: Session Token cannot be nil."}
	}

	path := fmt.Sprintf("/accounts/%s/transactions", accountNumber)

	type accountResponse struct {
		Data struct {
			Transactions []Transaction `json:"items"`
		} `json:"data"`
	}

	transactionsRes := new(accountResponse)

	header := http.Header{}
	header.Add("Authorization", *c.Session.SessionToken)

	err := c.request(http.MethodGet, path, header, query, nil, transactionsRes)
	if err != nil {
		return []Transaction{}, err
	}

	return transactionsRes.Data.Transactions, nil
}

// Retrieve a transaction by account number and ID.
func (c *Client) GetAccountTransaction(accountNumber string, id int) (Transaction, *Error) {
	if c.Session.SessionToken == nil {
		return Transaction{}, &Error{Message: "Session is invalid: Session Token cannot be nil."}
	}

	path := fmt.Sprintf("/accounts/%s/transactions/%d", accountNumber, id)

	type accountResponse struct {
		Transaction Transaction `json:"data"`
	}

	transactionsRes := new(accountResponse)

	header := http.Header{}
	header.Add("Authorization", *c.Session.SessionToken)

	err := c.request(http.MethodGet, path, header, nil, nil, transactionsRes)
	if err != nil {
		return Transaction{}, err
	}

	return transactionsRes.Transaction, nil
}

// Return the total fees for an account for a given day
// the day will default to today.
func (c *Client) GetAccountTransactionFees(accountNumber string, date *time.Time) (TransactionFees, *Error) {
	if c.Session.SessionToken == nil {
		return TransactionFees{}, &Error{Message: "Session is invalid: Session Token cannot be nil."}
	}

	path := fmt.Sprintf("/accounts/%s/transactions/total-fees", accountNumber)

	type accountResponse struct {
		TransactionFees TransactionFees `json:"data"`
	}

	transactionsRes := new(accountResponse)

	header := http.Header{}
	header.Add("Authorization", *c.Session.SessionToken)

	type dateQuery struct {
		// The date to get fees for, defaults to today
		Date string `url:"date,omitempty"`
	}

	query := dateQuery{}

	if date != nil {
		query.Date = date.Format("2006-01-02")
	}

	err := c.request(http.MethodGet, path, header, query, nil, transactionsRes)
	if err != nil {
		return TransactionFees{}, err
	}

	return transactionsRes.TransactionFees, nil
}
