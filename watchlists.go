package tasty

import (
	"fmt"
	"net/http"
	"net/url"
)

// Returns a list of all watchlists for the given account.
func (c *Client) GetMyWatchlists() ([]Watchlist, *http.Response, error) {
	path := "/watchlists"

	type watchlistResponse struct {
		Data struct {
			Watchlists []Watchlist `json:"items"`
		} `json:"data"`
	}

	watchlistsRes := new(watchlistResponse)

	resp, err := c.request(http.MethodGet, path, nil, nil, watchlistsRes)
	if err != nil {
		return []Watchlist{}, resp, err
	}

	return watchlistsRes.Data.Watchlists, resp, nil
}

// Returns a requested account watchlist.
func (c *Client) GetMyWatchlist(name string) (Watchlist, *http.Response, error) {
	path := fmt.Sprintf("/watchlists/%s", url.PathEscape(name))

	type watchlistResponse struct {
		Watchlist Watchlist `json:"data"`
	}

	watchlistsRes := new(watchlistResponse)

	// Must be customRequest in an instance where name has / within
	resp, err := c.customRequest(http.MethodGet, path, nil, nil, watchlistsRes)
	if err != nil {
		return Watchlist{}, resp, err
	}

	return watchlistsRes.Watchlist, resp, nil
}

// Create an account watchlist.
func (c *Client) CreateWatchlist(watchlist NewWatchlist) (Watchlist, *http.Response, error) {
	path := "/watchlists"

	type watchlistResponse struct {
		Watchlist Watchlist `json:"data"`
	}

	watchlistsRes := new(watchlistResponse)

	resp, err := c.request(http.MethodPost, path, nil, watchlist, watchlistsRes)
	if err != nil {
		return Watchlist{}, resp, err
	}

	return watchlistsRes.Watchlist, resp, nil
}

// Replace all properties of an account watchlist.
func (c *Client) EditWatchlist(name string, watchlist NewWatchlist) (Watchlist, *http.Response, error) {
	path := fmt.Sprintf("/watchlists/%s", url.PathEscape(name))

	type watchlistResponse struct {
		Watchlist Watchlist `json:"data"`
	}

	watchlistsRes := new(watchlistResponse)

	// Must be customRequest in an instance where name has / within
	resp, err := c.customRequest(http.MethodPut, path, nil, watchlist, watchlistsRes)
	if err != nil {
		return Watchlist{}, resp, err
	}

	return watchlistsRes.Watchlist, resp, nil
}

// Delete a watchlist for the given account.
func (c *Client) DeleteWatchlist(name string) (RemovedWatchlist, *http.Response, error) {
	path := fmt.Sprintf("/watchlists/%s", url.PathEscape(name))

	removedWatchlist := new(RemovedWatchlist)

	// Must be customRequest in an instance where name has / within
	resp, err := c.customRequest(http.MethodDelete, path, nil, nil, removedWatchlist)
	if err != nil {
		return RemovedWatchlist{}, resp, err
	}

	return *removedWatchlist, resp, nil
}

// Returns a list of all tastytrade pairs watchlists.
func (c *Client) GetPairsWatchlists() ([]PairsWatchlist, *http.Response, error) {
	path := "/pairs-watchlists"

	type watchlistResponse struct {
		Data struct {
			Watchlists []PairsWatchlist `json:"items"`
		} `json:"data"`
	}

	watchlistsRes := new(watchlistResponse)

	resp, err := c.request(http.MethodGet, path, nil, nil, watchlistsRes)
	if err != nil {
		return []PairsWatchlist{}, resp, err
	}

	return watchlistsRes.Data.Watchlists, resp, nil
}

// Returns a requested tastytrade pairs watchlist.
func (c *Client) GetPairsWatchlist(name string) (PairsWatchlist, *http.Response, error) {
	path := fmt.Sprintf("/pairs-watchlists/%s", url.PathEscape(name))

	type watchlistResponse struct {
		Watchlist PairsWatchlist `json:"data"`
	}

	watchlistsRes := new(watchlistResponse)

	// Must be customRequest in an instance where name has / within
	resp, err := c.customRequest(http.MethodGet, path, nil, nil, watchlistsRes)
	if err != nil {
		return PairsWatchlist{}, resp, err
	}

	return watchlistsRes.Watchlist, resp, nil
}

// Returns a list of all tastytrade watchlists.
func (c *Client) GetPublicWatchlists(countsOnly bool) ([]PublicWatchlist, *http.Response, error) {
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

	resp, err := c.request(http.MethodGet, path, query, nil, watchlistsRes)
	if err != nil {
		return []PublicWatchlist{}, resp, err
	}

	return watchlistsRes.Data.Watchlists, resp, nil
}

// Returns a requested tastytrade watchlist.
func (c *Client) GetPublicWatchlist(name string) (Watchlist, *http.Response, error) {
	path := fmt.Sprintf("/public-watchlists/%s", url.PathEscape(name))

	type watchlistResponse struct {
		Watchlist Watchlist `json:"data"`
	}

	watchlistsRes := new(watchlistResponse)

	// Must be customRequest in an instance where name has / within
	resp, err := c.customRequest(http.MethodGet, path, nil, nil, watchlistsRes)
	if err != nil {
		return Watchlist{}, resp, err
	}

	return watchlistsRes.Watchlist, resp, nil
}
