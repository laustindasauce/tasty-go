package queries

import (
	"time"

	"github.com/austinbspencer/tasty-go/constants"
)

type HistoricLiquidity struct {
	// If given, will return data for a specific period of time with a pre-defined
	// time interval. Passing 1d will return the previous day of data in 5 minute
	// intervals. This param is required if start-time is not given. 1d - If
	// equities market is open, this will return data starting from market open in
	// 5 minute intervals. If market is closed, will return data from previous
	// market open.
	TimeBack constants.TimeBack `url:"time-back,omitempty"`
	// StartTime is The start point for this query. This param is required is
	// time-back is not given. If given, will take precedence over time-back.
	StartTime time.Time `layout:"2006-01-02T15:04:05Z" url:"start-time,omitempty"`
}
