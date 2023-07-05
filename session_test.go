package tasty //nolint:testpackage // testing private field

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateSession(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/sessions", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, sessionResp)
	})

	resp, err := client.CreateSession(LoginInfo{Login: "default", Password: "Password"}, nil)
	require.Nil(t, err)

	require.Equal(t, "default@gmail.com", resp.User.Email)
	require.Equal(t, "default", resp.User.Username)
	require.Equal(t, "U0001563674", resp.User.ExternalID)
	require.Nil(t, resp.User.ID)
	require.NotNil(t, resp.SessionToken)
	require.Equal(t, "example-session-token+C", *resp.SessionToken)
	require.Nil(t, resp.RememberToken)
}

func TestCreateTwoFactorSession(t *testing.T) {
	setup()
	defer teardown()

	twoFaCode := "2-fa-code"

	mux.HandleFunc("/sessions", func(writer http.ResponseWriter, request *http.Request) {
		require.Equal(t, twoFaCode, request.Header.Get("X-Tastyworks-OTP"))
		fmt.Fprint(writer, sessionResp)
	})

	resp, err := client.CreateSession(LoginInfo{Login: "default", Password: "Password"}, &twoFaCode)
	require.Nil(t, err)

	require.Equal(t, "default@gmail.com", resp.User.Email)
	require.Equal(t, "default", resp.User.Username)
	require.Equal(t, "U0001563674", resp.User.ExternalID)
	require.NotNil(t, resp.SessionToken)
	require.Equal(t, "example-session-token+C", *resp.SessionToken)
	require.Nil(t, resp.RememberToken)
}

func TestCreateSessionError(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/sessions", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(401)
		fmt.Fprint(writer, tastyInvalidCredentialsError)
	})

	_, err := client.CreateSession(LoginInfo{Login: "default", Password: "Password"}, nil)
	expectedInvalidCredentials(t, err)
}

func TestValidateSession(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/sessions/validate", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, sessionValidateResp)
	})

	resp, err := client.ValidateSession()
	require.Nil(t, err)

	require.Equal(t, "default@gmail.com", resp.Email)
	require.Equal(t, "default", resp.Username)
	require.Equal(t, "U0001563674", resp.ExternalID)
	require.Equal(t, 123456, *resp.ID)
}

func TestValidateSessionError(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/sessions/validate", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(401)
		fmt.Fprint(writer, tastyInvalidSessionError)
	})

	_, err := client.ValidateSession()
	expectedInvalidSession(t, err)
}

func TestDestroySession(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/sessions", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, sessionResp)
	})

	err := client.DestroySession()
	require.Nil(t, err)
}

func TestDestroySessionError(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/sessions", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(401)
		fmt.Fprint(writer, tastyInvalidSessionError)
	})

	err := client.DestroySession()
	expectedInvalidSession(t, err)
}

func TestRequestPasswordResetEmail(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/password/reset", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(200)
	})

	err := client.RequestPasswordResetEmail("some-email@domain.com")
	require.Nil(t, err)
}

func TestRequestPasswordResetEmailError(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/password/reset", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(400)
		fmt.Fprint(writer, passwordChangeRequestErrorResp)
	})

	err := client.RequestPasswordResetEmail("")
	require.NotNil(t, err)

	require.Equal(t, "\nError in request 400;\nCode: validation_error\nMessage: Request validation failed", err.Error())
}

func TestChangePassword(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/password", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(200)
	})

	err := client.ChangePassword(
		PasswordReset{
			Password:             "newPassword",
			PasswordConfirmation: "newPassword",
			ResetPasswordToken:   "test-token",
		})
	require.Nil(t, err)
}

func TestChangePasswordError(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/password", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(400)
		fmt.Fprint(writer, passwordResetErrorResp)
	})

	err := client.ChangePassword(
		PasswordReset{
			Password:             "newPassword",
			PasswordConfirmation: "newPassword",
		})
	require.NotNil(t, err)

	require.Equal(t, "\nError in request 400;\nCode: validation_error\nMessage: Request validation failed", err.Error())
}

const sessionResp = `{
  "data": {
    "user": {
      "email": "default@gmail.com",
      "username": "default",
      "external-id": "U0001563674"
    },
    "session-token": "example-session-token+C"
  },
  "context": "/sessions"
}`

const sessionValidateResp = `{
  "data": {
    "email": "default@gmail.com",
    "username": "default",
    "external-id": "U0001563674",
    "id": 123456
  },
  "context": "/sessions/validate"
}`

const passwordChangeRequestErrorResp = `{
    "error": {
        "code": "validation_error",
        "message": "Request validation failed",
        "errors": [
            {
                "domain": "email",
                "reason": "is empty"
            }
        ]
    }
}`

const passwordResetErrorResp = `{
    "error": {
        "code": "validation_error",
        "message": "Request validation failed",
        "errors": [
            {
                "domain": "current-password",
                "reason": "are missing, exactly one parameter must be provided"
            },
            {
                "domain": "reset-password-token",
                "reason": "are missing, exactly one parameter must be provided"
            }
        ]
    }
}`
