package queries

import (
	"time"

	"github.com/austinbspencer/tasty-go/constants"
)

// Orders is the query for account orders.
type Orders struct {
	// Default value 10
	PerPage int `url:"per-page,omitempty"`
	// Default value 0
	PageOffset int `url:"page-offset,omitempty"`
	// The start date of orders to query.
	StartDate time.Time `layout:"2006-01-02" url:"start-date,omitempty"`
	// The end date of orders to query.
	EndDate time.Time `layout:"2006-01-02" url:"end-date,omitempty"`
	// The Underlying Symbol. The Ticker Symbol FB or
	// TW Future Symbol with out date component /M6E or
	// the full TW Future Symbol /ESU9
	UnderlyingSymbol string `url:"underlying-symbol,omitempty"`
	// Status of the order
	Status []constants.OrderStatus `url:"status[],omitempty"`
	// The full TW Future Symbol /ESZ9 or
	// /NGZ19 if two year digit are appropriate
	FuturesSymbol string `url:"futures-symbol"`
	// Underlying instrument type i.e. constants.InstrumentType
	UnderlyingInstrumentType constants.InstrumentType `url:"underlying-instrument-type,omitempty"`
	// The order to sort results in. Defaults to Desc, Accepts Desc or Asc.
	Sort constants.SortOrder `url:"sort,omitempty"`
	// DateTime start range for filtering transactions in full date-time
	StartAt time.Time `layout:"2006-01-02T15:04:05Z" url:"start-at,omitempty"`
	// DateTime end range for filtering transactions in full date-time
	EndAt time.Time `layout:"2006-01-02T15:04:05Z" url:"end-at,omitempty"`
}
