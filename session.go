package tasty

import (
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
	path := "/sessions"

	session := new(sessionResponse)

	header := http.Header{}
	if login.TwoFactorCode != nil {
		header.Add("X-Tastyworks-OTP", *login.TwoFactorCode)
	}

	err := c.request(http.MethodPost, path, header, nil, login, session)
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
	path := "/sessions/validate"

	session := new(sessionResponse)

	header := http.Header{}
	header.Add("Authorization", *c.Session.SessionToken)

	err := c.request(http.MethodPost, path, header, nil, nil, session)
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
	path := "/sessions"

	header := http.Header{}
	header.Add("Authorization", *c.Session.SessionToken)

	return c.request(http.MethodDelete, path, header, nil, nil, nil)
}
