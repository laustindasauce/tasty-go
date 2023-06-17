package tasty

import (
	"net/http"
)

type sessionResponse struct {
	Session Session `json:"data"`
}

// Create a new user session.
func (c *Client) CreateSession(login LoginInfo, twoFactorCode *string) (Session, *Error) {
	path := "/sessions"

	session := new(sessionResponse)

	header := http.Header{}

	if twoFactorCode != nil {
		header.Add("X-Tastyworks-OTP", *twoFactorCode)
	}

	err := c.noAuthRequest(http.MethodPost, path, header, nil, login, session)
	if err != nil {
		return Session{}, err
	}

	c.Session = session.Session

	return session.Session, nil
}

// Validate the user session.
func (c *Client) ValidateSession() (Session, *Error) {
	path := "/sessions/validate"

	session := new(sessionResponse)

	err := c.request(http.MethodPost, path, nil, nil, session)
	if err != nil {
		return Session{}, err
	}

	c.Session = session.Session

	return session.Session, nil
}

// Destroy the user session and invalidate the token.
func (c *Client) DestroySession() *Error {
	path := "/sessions"

	return c.request(http.MethodDelete, path, nil, nil, nil)
}

// Request a password reset email.
func (c *Client) RequestPasswordResetEmail(email string) *Error {
	path := "/password/reset"

	type reset struct {
		Email string `json:"email"`
	}

	resetInfo := new(reset)
	resetInfo.Email = email

	return c.noAuthRequest(http.MethodPost, path, http.Header{}, nil, resetInfo, nil)
}

// Request a password reset email.
func (c *Client) ChangePassword(resetInfo PasswordReset) *Error {
	path := "/password"

	return c.noAuthRequest(http.MethodPost, path, http.Header{}, nil, resetInfo, nil)
}
