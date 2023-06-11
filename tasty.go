package tasty

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/google/go-querystring/query"
)

const (
	apiBaseURL      = "https://api.tastyworks.com"
	apiBaseHost     = "api.tastyworks.com"
	apiCertBaseUrl  = "https://api.cert.tastyworks.com"
	apiCertBaseHost = "api.cert.tastyworks.com"
)

var (
	errorStatusCodes = []int{400, 401, 403, 404, 415, 500}
)

type Client struct {
	httpClient *http.Client
	baseURL    string
	baseHost   string
	Session    Session
}

// NewClient creates a new Tasty Client.
func NewClient(httpClient *http.Client) (*Client, error) {
	c := &Client{
		httpClient: httpClient,
		baseURL:    apiBaseURL,
		baseHost:   apiBaseHost,
	}

	return c, nil
}

// NewCertClient creates a new Tasty Cert Client.
func NewCertClient(httpClient *http.Client) (*Client, error) {
	c := &Client{
		httpClient: httpClient,
		baseURL:    apiCertBaseUrl,
		baseHost:   apiCertBaseHost,
	}

	return c, nil
}

// Error represents an error returned by the TastyWorks API.
type Error struct {
	// A short description of the error.
	Message string `json:"message"`
	// The HTTP status code.
	Err int `json:"error,omitempty"`
}

// Error ...
func (e Error) Error() string {
	return e.Message
}

// decodeError decodes an Error from response status code based off
// the developer docs in TastyWorks -> https://developer.tastytrade.com/#error-codes
func (c *Client) decodeError(resp *http.Response) *Error {
	e := new(Error)

	type errorRes struct {
		Error Error `json:"error"`
	}

	errRes := new(errorRes)

	err := json.NewDecoder(resp.Body).Decode(errRes)
	if err != nil {
		fmt.Println(err)
		e.Message = fmt.Sprintf("TastyWorks: unexpected HTTP %d: %s (empty error)", resp.StatusCode, err.Error())
		e.Err = resp.StatusCode
		return e
	}

	e = &errRes.Error

	return e
}

// customGet handles the get requests for the client with unique paths
func (c *Client) customGet(path string, header http.Header, params, result any) *Error {
	for {
		r := new(http.Request)

		r.Method = "GET"

		r.URL = &url.URL{
			Scheme: strings.Split(c.baseURL, ":")[0],
			Host:   c.baseHost,
			Opaque: path,
		}

		r.Header = header

		r.Header.Add("Content-Type", "application/json")

		if params != nil {
			queryString, err := query.Values(params)
			if err != nil {
				return &Error{Message: fmt.Sprintf("Client Side Error: %v", err)}
			}
			// fmt.Println(queryString.Encode())
			r.URL.RawQuery = queryString.Encode()
		}

		resp, err := c.httpClient.Do(r)
		if err != nil {
			return &Error{Message: fmt.Sprintf("Client Side Error: %v", err)}
		}

		// body, err := ioutil.ReadAll(resp.Body)
		// if err != nil {
		// 	return &Error{Message: fmt.Sprintf("Client Side Error: %v", err)}
		// }

		// fmt.Println(string(body))

		defer resp.Body.Close()

		if resp.StatusCode == http.StatusNoContent {
			return nil
		}
		if ContainsInt(errorStatusCodes, resp.StatusCode) {
			return c.decodeError(resp)
		}

		if result != nil {
			err = json.NewDecoder(resp.Body).Decode(result)
			if err != nil {
				return &Error{Message: fmt.Sprintf("Client Side Error: %v", err)}
			}
		}

		break
	}

	return nil
}

// get handles the get requests for the client
func (c *Client) get(url string, header http.Header, params, result any) *Error {
	for {
		r, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return &Error{Message: fmt.Sprintf("Client Side Error: %v", err)}
		}

		r.Header = header

		r.Header.Add("Content-Type", "application/json")

		if params != nil {
			queryString, err := query.Values(params)
			if err != nil {
				return &Error{Message: fmt.Sprintf("Client Side Error: %v", err)}
			}
			// fmt.Println(queryString.Encode())
			r.URL.RawQuery = queryString.Encode()
		}

		fmt.Println("Path: ", r.URL.Path)
		resp, err := c.httpClient.Do(r)
		if err != nil {
			return &Error{Message: fmt.Sprintf("Client Side Error: %v", err)}
		}

		// body, err := ioutil.ReadAll(resp.Body)
		// if err != nil {
		// 	return &Error{Message: fmt.Sprintf("Client Side Error: %v", err)}
		// }

		// fmt.Println(string(body))

		defer resp.Body.Close()

		if resp.StatusCode == http.StatusNoContent {
			return nil
		}
		if ContainsInt(errorStatusCodes, resp.StatusCode) {
			return c.decodeError(resp)
		}

		if result != nil {
			err = json.NewDecoder(resp.Body).Decode(result)
			if err != nil {
				return &Error{Message: fmt.Sprintf("Client Side Error: %v", err)}
			}
		}

		break
	}

	return nil
}

// post handles the post requests for the client
func (c *Client) post(url string, header http.Header, payload, result any) *Error {
	for {
		body, err := json.Marshal(payload)

		r, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
		if err != nil {
			return &Error{Message: fmt.Sprintf("Client Side Error: %v", err)}
		}

		r.Header = header

		r.Header.Add("Content-Type", "application/json")

		resp, err := c.httpClient.Do(r)
		if err != nil {
			return &Error{Message: fmt.Sprintf("Client Side Error: %v", err)}
		}

		// body, err = ioutil.ReadAll(resp.Body)

		// if err != nil {
		// 	return  &Error{Message: fmt.Sprintf("Client Side Error: %v", err)}
		// }

		// fmt.Println(string(body))

		defer resp.Body.Close()

		fmt.Println(resp.StatusCode)
		if resp.StatusCode == http.StatusNoContent {
			return nil
		}
		if ContainsInt(errorStatusCodes, resp.StatusCode) {
			return c.decodeError(resp)
		}

		if result != nil {
			err = json.NewDecoder(resp.Body).Decode(result)
			if err != nil {
				return &Error{Message: fmt.Sprintf("Client Side Error: %v", err)}
			}
		}

		break
	}

	return nil
}

// delete handles the delete requests for the client
func (c *Client) delete(url string, header http.Header, result any) *Error {
	for {
		r, err := http.NewRequest("DELETE", url, nil)
		if err != nil {
			return &Error{Message: fmt.Sprintf("Client Side Error: %v", err)}
		}

		r.Header = header

		r.Header.Add("Content-Type", "application/json")

		resp, err := c.httpClient.Do(r)
		if err != nil {
			return &Error{Message: fmt.Sprintf("Client Side Error: %v", err)}
		}

		// body, err := ioutil.ReadAll(resp.Body)

		// if err != nil {
		// 	return err
		// }

		// fmt.Println(string(body))

		defer resp.Body.Close()

		if resp.StatusCode == http.StatusNoContent {
			return nil
		}
		if ContainsInt(errorStatusCodes, resp.StatusCode) {
			return c.decodeError(resp)
		}

		if result != nil {
			err = json.NewDecoder(resp.Body).Decode(result)
			if err != nil {
				return &Error{Message: fmt.Sprintf("Client Side Error: %v", err)}
			}
		}

		break
	}

	return nil
}
