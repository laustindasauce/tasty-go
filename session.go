package tasty

import (
	"net/http"
)

// Create a new user session.
func (c *Client) CreateSession(login LoginInfo, twoFactorCode *string) (Session, *http.Response, error) {
	path := "/sessions"

	type sessionResponse struct {
		Session Session `json:"data"`
	}

	session := new(sessionResponse)

	header := http.Header{}

	if twoFactorCode != nil {
		header.Add("X-Tastyworks-OTP", *twoFactorCode)
	}

	resp, err := c.noAuthRequest(http.MethodPost, path, header, nil, login, session)
	if err != nil {
		return Session{}, resp, err
	}

	c.Session = session.Session

	return session.Session, resp, nil
}

// Validate the user session.
func (c *Client) ValidateSession() (User, *http.Response, error) {
	path := "/sessions/validate"

	type validSessionResponse struct {
		User User `json:"data"`
	}

	user := new(validSessionResponse)

	resp, err := c.request(http.MethodPost, path, nil, nil, user)
	if err != nil {
		return User{}, resp, err
	}

	c.Session.User = user.User

	return user.User, resp, nil
}

// Destroy the user session and invalidate the token.
func (c *Client) DestroySession() (*http.Response, error) {
	path := "/sessions"

	return c.request(http.MethodDelete, path, nil, nil, nil)
}

// Request a password reset email.
func (c *Client) RequestPasswordResetEmail(email string) (*http.Response, error) {
	path := "/password/reset"

	type reset struct {
		Email string `json:"email"`
	}

	resetInfo := new(reset)
	resetInfo.Email = email

	return c.noAuthRequest(http.MethodPost, path, http.Header{}, nil, resetInfo, nil)
}

// Request a password reset email.
func (c *Client) ChangePassword(resetInfo PasswordReset) (*http.Response, error) {
	path := "/password"

	return c.noAuthRequest(http.MethodPost, path, http.Header{}, nil, resetInfo, nil)
}
