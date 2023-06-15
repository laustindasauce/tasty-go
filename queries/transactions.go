package queries

import (
	"time"

	"github.com/austinbspencer/tasty-go/constants"
)

// Transactions is the query for account transactions.
type Transactions struct {
	// Default value 250
	PerPage int `url:"per-page,omitempty"`
	// Default value 0
	PageOffset int `url:"page-offset,omitempty"`
	// The order to sort results in. Defaults to Desc, Accepts Desc or Asc.
	Sort constants.SortOrder `url:"sort,omitempty"`
	// Filter based on transaction_type
	Type string `url:"type,omitempty"`
	// Allows filtering on multiple transaction_types
	Types []string `url:"types[],omitempty"`
	// Filter based on transaction_sub_type
	SubTypes []string `url:"sub-type[],omitempty"`
	// The start date of transactions to query.
	StartDate time.Time `layout:"2006-01-02" url:"start-date,omitempty"`
	// The end date of transactions to query. Defaults to now.
	EndDate time.Time `layout:"2006-01-02" url:"end-date,omitempty"`
	// The type of instrument i.e. constants.InstrumentType
	InstrumentType constants.InstrumentType `url:"instrument-type,omitempty"`
	// The Stock Ticker Symbol AAPL, OCC Option Symbol AAPL 191004P00275000,
	// TW Future Symbol /ESZ9, or TW Future Option Symbol ./ESZ9 EW4U9 190927P2975
	Symbol string `url:"symbol,omitempty"`
	// The Underlying Symbol. The Ticker Symbol FB or
	// TW Future Symbol with out date component /M6E or
	// the full TW Future Symbol /ESU9
	UnderlyingSymbol string `url:"underlying-symbol,omitempty"`
	// The action of the transaction. i.e. constants.OrderAction
	Action constants.OrderAction `url:"action,omitempty"`
	// Account partition key
	PartitionKey string `url:"partition-key,omitempty"`
	// The full TW Future Symbol /ESZ9 or
	// /NGZ19 if two year digit are appropriate
	FuturesSymbol string `url:"futures-symbol"`
	// DateTime start range for filtering transactions in full date-time
	StartAt time.Time `layout:"2006-01-02T15:04:05Z" url:"start-at,omitempty"`
	// DateTime end range for filtering transactions in full date-time
	EndAt time.Time `layout:"2006-01-02T15:04:05Z" url:"end-at,omitempty"`
}
