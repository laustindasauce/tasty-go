package tasty

import (
	"fmt"
	"net/http"
	"net/url"
)

// Returns a list of all watchlists for the given account.
func (c *Client) GetMyWatchlists() ([]Watchlist, error) {
	path := "/watchlists"

	type watchlistResponse struct {
		Data struct {
			Watchlists []Watchlist `json:"items"`
		} `json:"data"`
	}

	watchlistsRes := new(watchlistResponse)

	err := c.request(http.MethodGet, path, nil, nil, watchlistsRes)
	if err != nil {
		return []Watchlist{}, err
	}

	return watchlistsRes.Data.Watchlists, nil
}

// Returns a requested account watchlist.
func (c *Client) GetMyWatchlist(name string) (Watchlist, error) {
	path := fmt.Sprintf("/watchlists/%s", url.PathEscape(name))

	type watchlistResponse struct {
		Watchlist Watchlist `json:"data"`
	}

	watchlistsRes := new(watchlistResponse)

	// Must be customRequest in an instance where name has / within
	err := c.customRequest(http.MethodGet, path, nil, nil, watchlistsRes)
	if err != nil {
		return Watchlist{}, err
	}

	return watchlistsRes.Watchlist, nil
}

// Create an account watchlist.
func (c *Client) CreateWatchlist(watchlist NewWatchlist) (Watchlist, error) {
	path := "/watchlists"

	type watchlistResponse struct {
		Watchlist Watchlist `json:"data"`
	}

	watchlistsRes := new(watchlistResponse)

	err := c.request(http.MethodPost, path, nil, watchlist, watchlistsRes)
	if err != nil {
		return Watchlist{}, err
	}

	return watchlistsRes.Watchlist, nil
}

// Replace all properties of an account watchlist.
func (c *Client) EditWatchlist(name string, watchlist NewWatchlist) (Watchlist, error) {
	path := fmt.Sprintf("/watchlists/%s", url.PathEscape(name))

	type watchlistResponse struct {
		Watchlist Watchlist `json:"data"`
	}

	watchlistsRes := new(watchlistResponse)

	// Must be customRequest in an instance where name has / within
	err := c.customRequest(http.MethodPut, path, nil, watchlist, watchlistsRes)
	if err != nil {
		return Watchlist{}, err
	}

	return watchlistsRes.Watchlist, nil
}

// Delete a watchlist for the given account.
func (c *Client) DeleteWatchlist(name string) (RemovedWatchlist, error) {
	path := fmt.Sprintf("/watchlists/%s", url.PathEscape(name))

	removedWatchlist := new(RemovedWatchlist)

	// Must be customRequest in an instance where name has / within
	err := c.customRequest(http.MethodDelete, path, nil, nil, removedWatchlist)
	if err != nil {
		return RemovedWatchlist{}, err
	}

	return *removedWatchlist, nil
}

// Returns a list of all tastyworks pairs watchlists.
func (c *Client) GetPairsWatchlists() ([]PairsWatchlist, error) {
	path := "/pairs-watchlists"

	type watchlistResponse struct {
		Data struct {
			Watchlists []PairsWatchlist `json:"items"`
		} `json:"data"`
	}

	watchlistsRes := new(watchlistResponse)

	err := c.request(http.MethodGet, path, nil, nil, watchlistsRes)
	if err != nil {
		return []PairsWatchlist{}, err
	}

	return watchlistsRes.Data.Watchlists, nil
}

// Returns a requested tastyworks pairs watchlist.
func (c *Client) GetPairsWatchlist(name string) (PairsWatchlist, error) {
	path := fmt.Sprintf("/pairs-watchlists/%s", url.PathEscape(name))

	type watchlistResponse struct {
		Watchlist PairsWatchlist `json:"data"`
	}

	watchlistsRes := new(watchlistResponse)

	// Must be customRequest in an instance where name has / within
	err := c.customRequest(http.MethodGet, path, nil, nil, watchlistsRes)
	if err != nil {
		return PairsWatchlist{}, err
	}

	return watchlistsRes.Watchlist, nil
}

// Returns a list of all tastyworks watchlists.
func (c *Client) GetPublicWatchlists(countsOnly bool) ([]PublicWatchlist, error) {
	path := "/public-watchlists"

	type watchlistResponse struct {
		Data struct {
			Watchlists []PublicWatchlist `json:"items"`
		} `json:"data"`
	}

	watchlistsRes := new(watchlistResponse)

	type watchlistsQuery struct {
		CountsOnly bool `url:"counts-only"`
	}

	query := watchlistsQuery{CountsOnly: countsOnly}

	err := c.request(http.MethodGet, path, query, nil, watchlistsRes)
	if err != nil {
		return []PublicWatchlist{}, err
	}

	return watchlistsRes.Data.Watchlists, nil
}

// Returns a requested tastyworks watchlist.
func (c *Client) GetPublicWatchlist(name string) (Watchlist, error) {
	path := fmt.Sprintf("/public-watchlists/%s", url.PathEscape(name))

	type watchlistResponse struct {
		Watchlist Watchlist `json:"data"`
	}

	watchlistsRes := new(watchlistResponse)

	// Must be customRequest in an instance where name has / within
	err := c.customRequest(http.MethodGet, path, nil, nil, watchlistsRes)
	if err != nil {
		return Watchlist{}, err
	}

	return watchlistsRes.Watchlist, nil
}
