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

const registerURL = "/register"

func TestRegister(t *testing.T) {
	t.Parallel()

	userID := "42"
	userToken := "57f4dd05-638c-43c7-80c6-738900ce9b80"
	userEmail := "test@test.test"
	userPassword := "blah"

	noUserFn := func(string) (*pinub.User, error) {
		return nil, errors.New("")
	}
	userFn := func(string) (*pinub.User, error) {
		return &pinub.User{ID: userID}, nil
	}
	createUserFn := func(u *pinub.User) error {
		u.ID = userID
		return nil
	}
	addTokenFn := func(u *pinub.User) error {
		u.Token = userToken
		return nil
	}

	t.Run("get register", func(t *testing.T) {
		handler, _, rec := setUp()
		req := httptest.NewRequest("GET", registerURL, nil)
		handler.ServeHTTP(rec, req)

		claim.Equals(t, 200, rec.Code)
	})

	t.Run("post register", func(t *testing.T) {
		handler, client, rec := setUp()
		client.Us.UserFn = noUserFn
		client.Us.CreateUserFn = createUserFn
		client.Us.AddTokenFn = addTokenFn

		data := url.Values{
			"email":            {userEmail},
			"password":         {userPassword},
			"password_confirm": {userPassword},
		}
		req := httptest.NewRequest("POST", registerURL, strings.NewReader(data.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		handler.ServeHTTP(rec, req)

		claim.Equals(t, 303, rec.Code)
		claim.Equals(t, true, client.Us.UserInvoked)
		claim.Equals(t, true, client.Us.CreateUserInvoked)
		claim.Equals(t, true, client.Us.AddTokenInvoked)

		claim.Assert(t,
			strings.Contains(rec.Header().Get("Set-Cookie"), "keks="+userToken),
			"cookie should be set")
	})

	t.Run("post register invalid email", func(t *testing.T) {
		t.Parallel()

		handler, client, rec := setUp()
		data := url.Values{
			"email":            {"test@test"},
			"password":         {userPassword},
			"password_confirm": {userPassword},
		}
		req := httptest.NewRequest("POST", registerURL, strings.NewReader(data.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		handler.ServeHTTP(rec, req)

		claim.Equals(t, 200, rec.Code)
		claim.Equals(t, false, client.Us.CreateUserInvoked)

		claim.Assert(t,
			strings.Contains(rec.Body.String(), pinub.ErrUserInvalidEmail.Error()),
			"error message about email format should appear")
	})

	t.Run("post register short password", func(t *testing.T) {
		t.Parallel()

		handler, client, rec := setUp()
		data := url.Values{
			"email":            {userEmail},
			"password":         {"bla"},
			"password_confirm": {"bla"},
		}
		req := httptest.NewRequest("POST", registerURL, strings.NewReader(data.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		handler.ServeHTTP(rec, req)

		claim.Equals(t, 200, rec.Code)
		claim.Equals(t, false, client.Us.UserInvoked)

		claim.Assert(t,
			strings.Contains(rec.Body.String(), pinub.ErrUserShortPassword.Error()),
			"error message about short password should appear")
	})

	t.Run("post register different confirm password", func(t *testing.T) {
		t.Parallel()

		handler, client, rec := setUp()
		data := url.Values{
			"email":            {userEmail},
			"password":         {userPassword},
			"password_confirm": {"different password"},
		}
		req := httptest.NewRequest("POST", registerURL, strings.NewReader(data.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		handler.ServeHTTP(rec, req)

		claim.Equals(t, 200, rec.Code)
		claim.Equals(t, false, client.Us.UserInvoked)

		claim.Assert(t,
			strings.Contains(rec.Body.String(), "Passwords do not match"),
			"error message about confirm password should appear")
	})

	t.Run("post register present user", func(t *testing.T) {
		t.Parallel()

		handler, client, rec := setUp()
		client.Us.UserFn = userFn

		data := url.Values{
			"email":            {userEmail},
			"password":         {userPassword},
			"password_confirm": {userPassword},
		}
		req := httptest.NewRequest("POST", registerURL, strings.NewReader(data.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		handler.ServeHTTP(rec, req)

		claim.Equals(t, 200, rec.Code)
		claim.Equals(t, true, client.Us.UserInvoked)

		claim.Assert(t,
			strings.Contains(rec.Body.String(), "User exists already"),
			"error message about existing user should appear")
	})
}
