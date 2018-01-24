package http_test

import (
	"errors"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/pinub/pinub"
)

const registerURL = "/register"

func TestHTTP_RegisterGet(t *testing.T) {
	t.Parallel()
	h, _, r := setUp()

	req, err := http.NewRequest("GET", registerURL, nil)
	ok(t, err)

	h.ServeHTTP(r, req)
	equals(t, r.Code, http.StatusOK)
}

func TestHTTP_RegisterPost(t *testing.T) {
	t.Parallel()
	h, c, r := setUp()

	email := "test@test.test"
	password := "blah"
	id := "577edcd9-2d7d-4872-896a-bcd4665d70e6"
	token := "57f4dd05-638c-43c7-80c6-738900ce9b80"

	c.Us.UserFn = func(e string) (*pinub.User, error) {
		return nil, errors.New("")
	}

	c.Us.CreateUserFn = func(u *pinub.User) error {
		u.ID = id
		return nil
	}

	c.Us.AddTokenFn = func(u *pinub.User) error {
		u.Token = token
		return nil
	}

	data := url.Values{
		"email":            {email},
		"password":         {password},
		"password_confirm": {password},
	}
	req, err := http.NewRequest("POST", registerURL, strings.NewReader(data.Encode()))
	ok(t, err)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	h.ServeHTTP(r, req)

	equals(t, http.StatusSeeOther, r.Code)
	equals(t, true, c.Us.UserInvoked)
	equals(t, true, c.Us.CreateUserInvoked)
	equals(t, true, c.Us.AddTokenInvoked)

	assert(t,
		strings.Contains(r.Header().Get("Set-Cookie"), "keks="+token),
		"cookie should be set")

	t.Log(r.Body.String())
}

func TestHTTP_RegsiterPostInvalidEmail(t *testing.T) {
	t.Parallel()
	h, c, r := setUp()

	email := "test@test"
	password := "blah"

	data := url.Values{
		"email":            {email},
		"password":         {password},
		"password_confirm": {password},
	}
	req, err := http.NewRequest("POST", registerURL, strings.NewReader(data.Encode()))
	ok(t, err)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	h.ServeHTTP(r, req)

	equals(t, http.StatusOK, r.Code)
	equals(t, false, c.Us.UserInvoked)

	assert(t,
		strings.Contains(r.Body.String(), pinub.ErrUserInvalidEmail.Error()),
		"error message about email format should appear")
}

func TestHTTP_RegisterPostShortPassword(t *testing.T) {
	t.Parallel()
	h, c, r := setUp()

	email := "test@test.test"
	password := "bla"

	data := url.Values{
		"email":            {email},
		"password":         {password},
		"password_confirm": {password},
	}
	req, err := http.NewRequest("POST", registerURL, strings.NewReader(data.Encode()))
	ok(t, err)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	h.ServeHTTP(r, req)

	equals(t, http.StatusOK, r.Code)
	equals(t, false, c.Us.UserInvoked)

	assert(t,
		strings.Contains(r.Body.String(), pinub.ErrUserShortPassword.Error()),
		"error message about short password should appear")
}

func TestHTTP_RegisterPostDifferentConfirm(t *testing.T) {
	t.Parallel()
	h, c, r := setUp()

	email := "test@test.test"
	password := "blah"

	data := url.Values{
		"email":            {email},
		"password":         {password},
		"password_confirm": {password + "blah"},
	}
	req, err := http.NewRequest("POST", registerURL, strings.NewReader(data.Encode()))
	ok(t, err)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	h.ServeHTTP(r, req)

	equals(t, http.StatusOK, r.Code)
	equals(t, false, c.Us.UserInvoked)

	assert(t,
		strings.Contains(r.Body.String(), "Passwords do not match"),
		"error message about confirm password should appear")
}

func TestHTTP_RegisterPostUserPresent(t *testing.T) {
	t.Parallel()
	h, c, r := setUp()

	email := "test@test.test"
	password := "blah"

	c.Us.UserFn = func(e string) (*pinub.User, error) {
		return &pinub.User{ID: "1", Email: email}, nil
	}

	data := url.Values{
		"email":            {email},
		"password":         {password},
		"password_confirm": {password},
	}
	req, err := http.NewRequest("POST", registerURL, strings.NewReader(data.Encode()))
	ok(t, err)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	h.ServeHTTP(r, req)

	equals(t, http.StatusOK, r.Code)
	equals(t, true, c.Us.UserInvoked)

	t.Log(r.Body.String())
	assert(t,
		strings.Contains(r.Body.String(), "User exists already"),
		"error message about existing user should appear")
}
