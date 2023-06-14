package tasty

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/austinbspencer/tasty-go/models"
)

// Returns a list of all watchlists for the given account
func (c *Client) GetMyWatchlists() ([]models.Watchlist, *Error) {
	if c.Session.SessionToken == nil {
		return []models.Watchlist{}, &Error{Message: "Session is invalid: Session Token cannot be nil."}
	}

	path := "/watchlists"

	type watchlistResponse struct {
		Data struct {
			Watchlists []models.Watchlist `json:"items"`
		} `json:"data"`
	}

	watchlistsRes := new(watchlistResponse)

	header := http.Header{}
	header.Add("Authorization", *c.Session.SessionToken)

	err := c.request(http.MethodGet, path, header, nil, nil, watchlistsRes)
	if err != nil {
		return []models.Watchlist{}, err
	}

	return watchlistsRes.Data.Watchlists, nil
}

// Returns a requested account watchlist
func (c *Client) GetMyWatchlist(name string) (models.Watchlist, *Error) {
	if c.Session.SessionToken == nil {
		return models.Watchlist{}, &Error{Message: "Session is invalid: Session Token cannot be nil."}
	}

	path := fmt.Sprintf("/watchlists/%s", url.PathEscape(name))

	type watchlistResponse struct {
		Watchlist models.Watchlist `json:"data"`
	}

	watchlistsRes := new(watchlistResponse)

	header := http.Header{}
	header.Add("Authorization", *c.Session.SessionToken)

	// Must be customRequest in an instance where name has / within
	err := c.customRequest(http.MethodGet, path, header, nil, nil, watchlistsRes)
	if err != nil {
		return models.Watchlist{}, err
	}

	return watchlistsRes.Watchlist, nil
}

// Create an account watchlist
func (c *Client) CreateWatchlist(watchlist models.NewWatchlist) (models.Watchlist, *Error) {
	if c.Session.SessionToken == nil {
		return models.Watchlist{}, &Error{Message: "Session is invalid: Session Token cannot be nil."}
	}

	path := "/watchlists"

	type watchlistResponse struct {
		Watchlist models.Watchlist `json:"data"`
	}

	watchlistsRes := new(watchlistResponse)

	header := http.Header{}
	header.Add("Authorization", *c.Session.SessionToken)

	err := c.request(http.MethodPost, path, header, nil, watchlist, watchlistsRes)
	if err != nil {
		return models.Watchlist{}, err
	}

	return watchlistsRes.Watchlist, nil
}

// Replace all properties of an account watchlist
func (c *Client) EditWatchlist(name string, watchlist models.NewWatchlist) (models.Watchlist, *Error) {
	if c.Session.SessionToken == nil {
		return models.Watchlist{}, &Error{Message: "Session is invalid: Session Token cannot be nil."}
	}

	path := fmt.Sprintf("/watchlists/%s", url.PathEscape(name))

	type watchlistResponse struct {
		Watchlist models.Watchlist `json:"data"`
	}

	watchlistsRes := new(watchlistResponse)

	header := http.Header{}
	header.Add("Authorization", *c.Session.SessionToken)

	// Must be customRequest in an instance where name has / within
	err := c.customRequest(http.MethodPut, path, header, nil, watchlist, watchlistsRes)
	if err != nil {
		return models.Watchlist{}, err
	}

	return watchlistsRes.Watchlist, nil
}

// Delete a watchlist for the given account
func (c *Client) DeleteWatchlist(name string) (models.RemovedWatchlist, *Error) {
	if c.Session.SessionToken == nil {
		return models.RemovedWatchlist{}, &Error{Message: "Session is invalid: Session Token cannot be nil."}
	}

	path := fmt.Sprintf("/watchlists/%s", url.PathEscape(name))

	header := http.Header{}
	header.Add("Authorization", *c.Session.SessionToken)

	removedWatchlist := new(models.RemovedWatchlist)

	// Must be customRequest in an instance where name has / within
	err := c.customRequest(http.MethodDelete, path, header, nil, nil, removedWatchlist)
	if err != nil {
		return models.RemovedWatchlist{}, err
	}

	return *removedWatchlist, nil
}

// Returns a list of all tastyworks pairs watchlists
func (c *Client) GetPairsWatchlists() ([]models.PairsWatchlist, *Error) {
	if c.Session.SessionToken == nil {
		return []models.PairsWatchlist{}, &Error{Message: "Session is invalid: Session Token cannot be nil."}
	}

	path := "/pairs-watchlists"

	type watchlistResponse struct {
		Data struct {
			Watchlists []models.PairsWatchlist `json:"items"`
		} `json:"data"`
	}

	watchlistsRes := new(watchlistResponse)

	header := http.Header{}
	header.Add("Authorization", *c.Session.SessionToken)

	err := c.request(http.MethodGet, path, header, nil, nil, watchlistsRes)
	if err != nil {
		return []models.PairsWatchlist{}, err
	}

	return watchlistsRes.Data.Watchlists, nil
}

// Returns a requested tastyworks pairs watchlist
func (c *Client) GetPairsWatchlist(name string) (models.PairsWatchlist, *Error) {
	if c.Session.SessionToken == nil {
		return models.PairsWatchlist{}, &Error{Message: "Session is invalid: Session Token cannot be nil."}
	}

	path := fmt.Sprintf("/pairs-watchlists/%s", url.PathEscape(name))

	type watchlistResponse struct {
		Watchlist models.PairsWatchlist `json:"data"`
	}

	watchlistsRes := new(watchlistResponse)

	header := http.Header{}
	header.Add("Authorization", *c.Session.SessionToken)

	// Must be customRequest in an instance where name has / within
	err := c.customRequest(http.MethodGet, path, header, nil, nil, watchlistsRes)
	if err != nil {
		return models.PairsWatchlist{}, err
	}

	return watchlistsRes.Watchlist, nil
}

// Returns a list of all tastyworks watchlists
func (c *Client) GetPublicWatchlists(countsOnly bool) ([]models.PublicWatchlist, *Error) {
	if c.Session.SessionToken == nil {
		return []models.PublicWatchlist{}, &Error{Message: "Session is invalid: Session Token cannot be nil."}
	}

	path := "/public-watchlists"

	type watchlistResponse struct {
		Data struct {
			Watchlists []models.PublicWatchlist `json:"items"`
		} `json:"data"`
	}

	watchlistsRes := new(watchlistResponse)

	header := http.Header{}
	header.Add("Authorization", *c.Session.SessionToken)

	type watchlistsQuery struct {
		CountsOnly bool `url:"counts-only"`
	}

	query := watchlistsQuery{CountsOnly: countsOnly}

	err := c.request(http.MethodGet, path, header, query, nil, watchlistsRes)
	if err != nil {
		return []models.PublicWatchlist{}, err
	}

	return watchlistsRes.Data.Watchlists, nil
}

// Returns a requested tastyworks watchlist
func (c *Client) GetPublicWatchlist(name string) (models.Watchlist, *Error) {
	if c.Session.SessionToken == nil {
		return models.Watchlist{}, &Error{Message: "Session is invalid: Session Token cannot be nil."}
	}

	path := fmt.Sprintf("/public-watchlists/%s", url.PathEscape(name))

	type watchlistResponse struct {
		Watchlist models.Watchlist `json:"data"`
	}

	watchlistsRes := new(watchlistResponse)

	header := http.Header{}
	header.Add("Authorization", *c.Session.SessionToken)

	// Must be customRequest in an instance where name has / within
	err := c.customRequest(http.MethodGet, path, header, nil, nil, watchlistsRes)
	if err != nil {
		return models.Watchlist{}, err
	}

	return watchlistsRes.Watchlist, nil
}
