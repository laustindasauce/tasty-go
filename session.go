package tasty

import (
	"net/http"
)

// Create a new user session.
func (c *Client) CreateSession(login LoginInfo, twoFactorCode *string) (Session, error) {
	path := "/sessions"

	type sessionResponse struct {
		Session Session `json:"data"`
	}

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
func (c *Client) ValidateSession() (User, error) {
	path := "/sessions/validate"

	type validSessionResponse struct {
		User User `json:"data"`
	}

	user := new(validSessionResponse)

	err := c.request(http.MethodPost, path, nil, nil, user)
	if err != nil {
		return User{}, err
	}

	c.Session.User = user.User

	return user.User, nil
}

// Destroy the user session and invalidate the token.
func (c *Client) DestroySession() error {
	path := "/sessions"

	return c.request(http.MethodDelete, path, nil, nil, nil)
}

// Request a password reset email.
func (c *Client) RequestPasswordResetEmail(email string) error {
	path := "/password/reset"

	type reset struct {
		Email string `json:"email"`
	}

	resetInfo := new(reset)
	resetInfo.Email = email

	return c.noAuthRequest(http.MethodPost, path, http.Header{}, nil, resetInfo, nil)
}

// Request a password reset email.
func (c *Client) ChangePassword(resetInfo PasswordReset) error {
	path := "/password"

	return c.noAuthRequest(http.MethodPost, path, http.Header{}, nil, resetInfo, nil)
}
