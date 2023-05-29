package tasty

import (
	"fmt"
	"net/http"
)

type LoginInfo struct {
	Login         string `json:"login"`
	Password      string `json:"password"`
	TwoFactorCode *string
	RememberMe    bool
}

type User struct {
	Email      string `json:"email"`
	Username   string `json:"username"`
	ExternalID string `json:"external-id"`
}

type Session struct {
	User          User    `json:"user"`
	SessionToken  *string `json:"session-token"`
	RememberToken *string `json:"remember-token"`
}

type sessionResponse struct {
	Session Session `json:"data"`
	Context string  `json:"context"`
}

func (c *Client) CreateSession(login LoginInfo) (Session, *Error) {
	reqURL := fmt.Sprintf("%s/sessions", c.baseURL)

	session := new(sessionResponse)

	header := http.Header{}
	if login.TwoFactorCode != nil {
		header.Add("X-Tastyworks-OTP", *login.TwoFactorCode)
	}

	err := c.post(reqURL, header, login, session)
	if err == nil {
		(*c).Session = session.Session
	} else {
		return Session{}, err
	}

	return session.Session, err
}

func (c *Client) ValidateSession() (Session, *Error) {
	if c.Session.SessionToken == nil {
		return Session{}, &Error{Message: "Session is invalid: Session Token cannot be nil."}
	}
	reqURL := fmt.Sprintf("%s/sessions/validate", c.baseURL)

	session := new(sessionResponse)

	header := http.Header{}
	header.Add("Authorization", *c.Session.SessionToken)

	err := c.post(reqURL, header, nil, session)
	if err == nil {
		(*c).Session = session.Session
	} else {
		return Session{}, err
	}

	return session.Session, nil
}

func (c *Client) DestroySession() *Error {
	if c.Session.SessionToken == nil {
		return &Error{Message: "Session is invalid: Session Token cannot be nil."}
	}
	reqURL := fmt.Sprintf("%s/sessions", c.baseURL)

	header := http.Header{}
	header.Add("Authorization", *c.Session.SessionToken)

	return c.delete(reqURL, header, nil)
}
