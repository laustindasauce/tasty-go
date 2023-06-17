package tasty

type LoginInfo struct {
	// The user name or email of the user.
	Login string `json:"login"`
	// The password for the user's account
	Password string `json:"password"`
	// If the session should be extended for longer than normal
	// via remember token. Defaults to false.
	RememberMe bool `json:"remember-me"`
	// Valid for 28 days
	// Allows skipping for 2 factor with in its window.
	RememberToken string `json:"remember-token"`
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

type PasswordReset struct {
	Password             string `json:"password"`
	PasswordConfirmation string `json:"password-confirmation"`
	ResetPasswordToken   string `json:"reset-password-token"`
}
