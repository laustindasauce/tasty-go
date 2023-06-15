package queries

import (
	"time"

	"github.com/austinbspencer/tasty-go/constants"
)

// Optional filtering available for position endpoint.
type AccountPosition struct {
	// UnderlyingSymbol An array of Underlying symbol(s) for positions
	UnderlyingSymbol []string `url:"underlying-symbol[],omitempty"`
	// Symbol A single symbol. Stock Ticker Symbol AAPL,
	// OCC Option Symbol AAPL 191004P00275000,
	// TW Future Symbol /ESZ9, or
	// TW Future Option Symbol ./ESZ9 EW4U9 190927P2975
	Symbol string `url:"symbol,omitempty"`
	// InstrumentType The type of Instrument
	InstrumentType constants.InstrumentType `url:"instrument-type,omitempty"`
	// IncludeClosedPositions If closed positions should be included in the query
	IncludeClosedPositions bool `url:"include-closed-positions,omitempty"`
	// UnderlyingProductCode The underlying Future's Product code. i.e ES
	UnderlyingProductCode string `url:"underlying-product-code,omitempty"`
	// PartitionKeys Account partition keys
	PartitionKeys []string `url:"partition-keys[],omitempty"`
	// NetPositions Returns net positions grouped by instrument type and symbol
	NetPositions bool `url:"net-positions,omitempty"`
	// IncludeMarks Include current quote mark (note: can decrease performance)
	IncludeMarks bool `url:"include-marks,omitempty"`
}

// Filtering for the account balance snapshots endpoint.
type AccountBalanceSnapshots struct {
	// SnapshotDate The day of the balance snapshot to retrieve
	SnapshotDate time.Time `layout:"2006-01-02" url:"snapshot-date,omitempty"`
	// TimeOfDay The abbreviation for the time of day. Default value: EOD
	// Available values: EOD, BOD.
	TimeOfDay constants.TimeOfDay `url:"time-of-day"`
}
