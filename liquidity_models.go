package tasty

import (
	"time"

	"github.com/shopspring/decimal"
)

// NetLiqOHLC Represents the open, high, low, close values for a time interval.
// The close value is the value at time. Also includes pending cash values and
// a sum of the two for convenience.
type NetLiqOHLC struct {
	Open             decimal.Decimal `json:"open"`
	High             decimal.Decimal `json:"high"`
	Low              decimal.Decimal `json:"low"`
	Close            decimal.Decimal `json:"close"`
	PendingCashOpen  decimal.Decimal `json:"pending-cash-open"`
	PendingCashHigh  decimal.Decimal `json:"pending-cash-high"`
	PendingCashLow   decimal.Decimal `json:"pending-cash-low"`
	PendingCashClose decimal.Decimal `json:"pending-cash-close"`
	TotalOpen        decimal.Decimal `json:"total-open"`
	TotalHigh        decimal.Decimal `json:"total-high"`
	TotalLow         decimal.Decimal `json:"total-low"`
	TotalClose       decimal.Decimal `json:"total-close"`
	Time             string          `json:"time"`
}

type HistoricLiquidityQuery struct {
	// If given, will return data for a specific period of time with a pre-defined
	// time interval. Passing 1d will return the previous day of data in 5 minute
	// intervals. This param is required if start-time is not given. 1d - If
	// equities market is open, this will return data starting from market open in
	// 5 minute intervals. If market is closed, will return data from previous
	// market open.
	TimeBack TimeBack `url:"time-back,omitempty"`
	// StartTime is The start point for this query. This param is required is
	// time-back is not given. If given, will take precedence over time-back.
	StartTime time.Time `layout:"2006-01-02T15:04:05Z" url:"start-time,omitempty"`
}
