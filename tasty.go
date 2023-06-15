package tasty

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/austinbspencer/tasty-go/utils"
	"github.com/google/go-querystring/query"
)

const (
	apiBaseURL      = "https://api.tastyworks.com"
	apiBaseHost     = "api.tastyworks.com"
	apiCertBaseUrl  = "https://api.cert.tastyworks.com"
	apiCertBaseHost = "api.cert.tastyworks.com"
)

var (
	errorStatusCodes = []int{400, 401, 403, 404, 415, 422, 500}
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

type TastyError struct {
	Domain string `json:"domain"`
	Reason string `json:"reason"`
}

// Error represents an error returned by the TastyWorks API.
type Error struct {
	// Simple code error string
	Code string `json:"code"`
	// A short description of the error.
	Message string `json:"message"`
	// Slice of errors
	Errors []TastyError `json:"errors"`
	// The HTTP status code.
	StatusCode int `json:"error,omitempty"`
}

// Error ...
func (e Error) Error() string {
	return fmt.Sprintf("\nError in request %d;\nCode: %s\nMessage: %s", e.StatusCode, e.Code, e.Message)
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
		e.Message = fmt.Sprintf("TastyWorks: unexpected HTTP %d: %s (empty error)", resp.StatusCode, err.Error())
		e.StatusCode = resp.StatusCode
		return e
	}

	errRes.Error.StatusCode = resp.StatusCode

	e = &errRes.Error

	return e
}

// customRequest handles any requests for the client with unique paths
func (c *Client) customRequest(method, path string, header http.Header, params, payload, result any) *Error {
	r := new(http.Request)

	r.Method = method

	r.URL = &url.URL{
		Scheme: strings.Split(c.baseURL, ":")[0],
		Host:   c.baseHost,
		Opaque: fmt.Sprintf("//%s%s", c.baseHost, path),
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return &Error{Message: fmt.Sprintf("Client Side Error: %v", err)}
	}

	r.Body = io.NopCloser(bytes.NewBuffer(body))

	r.Header = header

	r.Header.Add("Content-Type", "application/json")

	if params != nil {
		queryString, err := query.Values(params)
		if err != nil {
			return &Error{Message: fmt.Sprintf("Client Side Error: %v", err)}
		}
		r.URL.RawQuery = queryString.Encode()
	}

	resp, err := c.httpClient.Do(r)
	if err != nil {
		return &Error{Message: fmt.Sprintf("Client Side Error: %v", err)}
	}

	// body, err = ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	return &Error{Message: fmt.Sprintf("Client Side Error: %v", err)}
	// }

	// fmt.Println(string(body))

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		return nil
	}
	if utils.ContainsInt(errorStatusCodes, resp.StatusCode) {
		return c.decodeError(resp)
	}

	if result != nil {
		err = json.NewDecoder(resp.Body).Decode(result)
		if err != nil {
			return &Error{Message: fmt.Sprintf("Client Side Error: %v", err)}
		}
	}

	return nil
}

// request handles any requests for the client
func (c *Client) request(method, path string, header http.Header, params, payload, result any) *Error {
	body, err := json.Marshal(payload)
	if err != nil {
		return &Error{Message: fmt.Sprintf("Client Side Error: %v", err)}
	}

	fullURL := c.baseURL + path

	r, err := http.NewRequest(method, fullURL, bytes.NewBuffer(body))
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
		r.URL.RawQuery = queryString.Encode()
	}

	resp, err := c.httpClient.Do(r)
	if err != nil {
		return &Error{Message: fmt.Sprintf("Client Side Error: %v", err)}
	}

	// body, err = ioutil.ReadAll(resp.Body)

	// if err != nil {
	// 	return &Error{Message: fmt.Sprintf("Client Side Error: %v", err)}
	// }

	// fmt.Println(string(body))

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		return nil
	}
	if utils.ContainsInt(errorStatusCodes, resp.StatusCode) {
		return c.decodeError(resp)
	}

	if result != nil {
		err = json.NewDecoder(resp.Body).Decode(result)
		if err != nil {
			return &Error{Message: fmt.Sprintf("Client Side Error: %v", err)}
		}
	}

	return nil
}
