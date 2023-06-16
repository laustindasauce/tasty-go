package tasty

import (
	"net/http"

	"github.com/austinbspencer/tasty-go/models"
)

type sessionResponse struct {
	Session models.Session `json:"data"`
}

// Create a new user session.
func (c *Client) CreateSession(login models.LoginInfo, twoFactorCode *string) (models.Session, *Error) {
	path := "/sessions"

	session := new(sessionResponse)

	header := http.Header{}

	if twoFactorCode != nil {
		header.Add("X-Tastyworks-OTP", *twoFactorCode)
	}

	err := c.request(http.MethodPost, path, header, nil, login, session)
	if err == nil {
		c.Session = session.Session
	} else {
		return models.Session{}, err
	}

	return session.Session, nil
}

// Validate the user session.
func (c *Client) ValidateSession() (models.Session, *Error) {
	if c.Session.SessionToken == nil {
		return models.Session{}, &Error{Message: "Session is invalid: Session Token cannot be nil."}
	}
	path := "/sessions/validate"

	session := new(sessionResponse)

	header := http.Header{}
	header.Add("Authorization", *c.Session.SessionToken)

	err := c.request(http.MethodPost, path, header, nil, nil, session)
	if err != nil {
		return models.Session{}, err
	}

	c.Session = session.Session

	return session.Session, nil
}

// Destroy the user session and invalidate the token.
func (c *Client) DestroySession() *Error {
	if c.Session.SessionToken == nil {
		return &Error{Message: "Session is invalid: Session Token cannot be nil."}
	}
	path := "/sessions"

	header := http.Header{}
	header.Add("Authorization", *c.Session.SessionToken)

	return c.request(http.MethodDelete, path, header, nil, nil, nil)
}
