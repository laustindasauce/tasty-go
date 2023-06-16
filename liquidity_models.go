package tasty

import (
	"time"
)

// NetLiqOHLC Represents the open, high, low, close values for a time interval.
// The close value is the value at time. Also includes pending cash values and
// a sum of the two for convenience.
type NetLiqOHLC struct {
	Open             StringToFloat32 `json:"open"`
	High             StringToFloat32 `json:"high"`
	Low              StringToFloat32 `json:"low"`
	Close            StringToFloat32 `json:"close"`
	PendingCashOpen  StringToFloat32 `json:"pending-cash-open"`
	PendingCashHigh  StringToFloat32 `json:"pending-cash-high"`
	PendingCashLow   StringToFloat32 `json:"pending-cash-low"`
	PendingCashClose StringToFloat32 `json:"pending-cash-close"`
	TotalOpen        StringToFloat32 `json:"total-open"`
	TotalHigh        StringToFloat32 `json:"total-high"`
	TotalLow         StringToFloat32 `json:"total-low"`
	TotalClose       StringToFloat32 `json:"total-close"`
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
