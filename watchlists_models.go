package tasty

import (
	"time"
)

type RemovedWatchlist struct {
	ID               *int             `json:"id"`
	UserID           *int             `json:"user_id"`
	Name             string           `json:"name"`
	CreatedAt        time.Time        `json:"created_at"`
	UpdatedAt        time.Time        `json:"updated_at"`
	WatchlistEntries []WatchlistEntry `json:"watchlist_entries"`
	GroupName        string           `json:"group_name"`
	OrderIndex       int              `json:"order_index"`
	CmsID            *string          `json:"cms_id"`
}

type Watchlist struct {
	ID                  *int             `json:"id"`
	Name                string           `json:"name"`
	WatchlistEntries    []WatchlistEntry `json:"watchlist-entries"`
	GroupName           string           `json:"group-name"`
	OrderIndex          int              `json:"order-index"`
	CmsID               *string          `json:"cms-id"`
	WatchlistEntryCount *int             `json:"watchlist-entry-count"`
}

type NewWatchlist struct {
	Name             string           `json:"name"`
	WatchlistEntries []WatchlistEntry `json:"watchlist-entries"`
	GroupName        string           `json:"group-name,omitempty"`
	OrderIndex       int              `json:"order-index,omitempty"`
}

type WatchlistEntry struct {
	Symbol         string         `json:"symbol"`
	InstrumentType InstrumentType `json:"instrument-type"`
}

type PublicWatchlist struct {
	ID                  *int                   `json:"id"`
	Name                string                 `json:"name"`
	WatchlistEntries    []PublicWatchlistEntry `json:"watchlist-entries"`
	GroupName           string                 `json:"group-name"`
	OrderIndex          int                    `json:"order-index"`
	WatchlistEntryCount *int                   `json:"watchlist-entry-count"`
}

// Something weird here in the api where instrument_type instead of instrument-type.
type PublicWatchlistEntry struct {
	Symbol         string         `json:"symbol"`
	InstrumentType InstrumentType `json:"instrument_type"`
}

type PairsWatchlist struct {
	Name           string          `json:"name"`
	PairsEquations []PairsEquation `json:"pairs-equations"`
	OrderIndex     int             `json:"order-index,omitempty"`
}

type PairsEquation struct {
	LeftAction    string `json:"left-action"`
	LeftSymbol    string `json:"left-symbol"`
	LeftQuantity  int    `json:"left-quantity"`
	RightAction   string `json:"right-action"`
	RightSymbol   string `json:"right-symbol"`
	RightQuantity int    `json:"right-quantity"`
}
