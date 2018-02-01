package http_test

import (
	"errors"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/pinub/pinub"
	"github.com/pinub/pinub/claim"
)

const signinURL = "/signin"

func TestSignin(t *testing.T) {
	t.Parallel()

	userID := "45"
	userToken := "cd6007cd-57f5-47e6-b3b2-e296fe1c4cba"
	userEmail := "test@test.signin"
	userPassword := "signin"

	noUserFn := func(string) (*pinub.User, error) {
		return nil, errors.New("")
	}
	userFn := func(email string) (*pinub.User, error) {
		pwd, _ := pinub.HashPassword(userPassword)
		return &pinub.User{ID: userID, Email: email, Password: pwd}, nil
	}
	userWrongPasswordFn := func(email string) (*pinub.User, error) {
		return &pinub.User{ID: userID, Email: email, Password: "different"}, nil
	}
	addTokenFn := func(u *pinub.User) error {
		u.Token = userToken
		return nil
	}
	addTokenErrorFn := func(*pinub.User) error {
		return errors.New("")
	}

	t.Run("signin get", func(t *testing.T) {
		handler, _, rec := setUp()

		req := httptest.NewRequest("GET", signinURL, nil)
		handler.ServeHTTP(rec, req)

		claim.Equals(t, 200, rec.Code)
	})

	t.Run("signin success", func(t *testing.T) {
		handler, client, rec := setUp()
		client.Us.UserFn = userFn
		client.Us.AddTokenFn = addTokenFn

		data := url.Values{
			"email":    {userEmail},
			"password": {userPassword},
		}
		req := httptest.NewRequest("POST", signinURL, strings.NewReader(data.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		handler.ServeHTTP(rec, req)

		claim.Equals(t, 303, rec.Code)
		claim.Equals(t, true, client.Us.UserInvoked)
		claim.Equals(t, true, client.Us.AddTokenInvoked)

		claim.Assert(t,
			strings.Contains(rec.Header().Get("Set-Cookie"), "keks="+userToken),
			"cookie should be set")
	})

	t.Run("signin short password", func(t *testing.T) {
		handler, client, rec := setUp()

		data := url.Values{
			"email":    {userEmail},
			"password": {"bla"},
		}
		req := httptest.NewRequest("POST", signinURL, strings.NewReader(data.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		handler.ServeHTTP(rec, req)

		claim.Equals(t, 200, rec.Code)
		claim.Equals(t, false, client.Us.UserInvoked)

		claim.Assert(t,
			strings.Contains(rec.Body.String(), pinub.ErrUserShortPassword.Error()),
			"error message about short password should appear")
	})

	t.Run("signin invalid email", func(t *testing.T) {
		handler, client, rec := setUp()

		data := url.Values{
			"email":    {"test@test"},
			"password": {userPassword},
		}
		req := httptest.NewRequest("POST", signinURL, strings.NewReader(data.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		handler.ServeHTTP(rec, req)

		claim.Equals(t, 200, rec.Code)
		claim.Equals(t, false, client.Us.UserInvoked)

		claim.Assert(t,
			strings.Contains(rec.Body.String(), pinub.ErrUserInvalidEmail.Error()),
			"error message about email format should appear")
	})

	t.Run("signin no user", func(t *testing.T) {
		handler, client, rec := setUp()
		client.Us.UserFn = noUserFn

		data := url.Values{
			"email":    {userEmail},
			"password": {userPassword},
		}
		req := httptest.NewRequest("POST", signinURL, strings.NewReader(data.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		handler.ServeHTTP(rec, req)

		claim.Equals(t, 200, rec.Code)
		claim.Equals(t, true, client.Us.UserInvoked)

		claim.Assert(t,
			strings.Contains(rec.Body.String(), pinub.ErrUserPasswordNotCorrect.Error()),
			"error message about email or password should appear")
	})

	t.Run("signin wrong password", func(t *testing.T) {
		handler, client, rec := setUp()
		client.Us.UserFn = userWrongPasswordFn

		data := url.Values{
			"email":    {userEmail},
			"password": {userPassword},
		}
		req := httptest.NewRequest("POST", signinURL, strings.NewReader(data.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		handler.ServeHTTP(rec, req)

		claim.Equals(t, 200, rec.Code)
		claim.Equals(t, true, client.Us.UserInvoked)

		claim.Assert(t,
			strings.Contains(rec.Body.String(), pinub.ErrUserPasswordNotCorrect.Error()),
			"error message about email or password should appear")
	})

	t.Run("signin cannot create token", func(t *testing.T) {
		handler, client, rec := setUp()
		client.Us.UserFn = userFn
		client.Us.AddTokenFn = addTokenErrorFn

		data := url.Values{
			"email":    {userEmail},
			"password": {userPassword},
		}
		req := httptest.NewRequest("POST", signinURL, strings.NewReader(data.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		handler.ServeHTTP(rec, req)

		claim.Equals(t, 200, rec.Code)
		claim.Equals(t, true, client.Us.UserInvoked)
		claim.Equals(t, true, client.Us.AddTokenInvoked)

		t.Log(rec.Body.String())
		claim.Assert(t,
			strings.Contains(rec.Body.String(), "Cannot sign in user"),
			"error message about email or password should appear")
	})
}
