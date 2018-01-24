package http_test

import (
	"errors"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/pinub/pinub"
)

const signinURL = "/signin"

func TestHTTP_SigninGet(t *testing.T) {
	t.Parallel()
	h, _, r := setUp()

	req, err := http.NewRequest("GET", signinURL, nil)
	ok(t, err)

	h.ServeHTTP(r, req)
	equals(t, r.Code, http.StatusOK)
}

func TestHTTP_SigninPost(t *testing.T) {
	t.Parallel()
	h, c, r := setUp()

	email := "test@test.test"
	password := "blahblah"
	id := "64f1001c-7992-470d-8e91-c64f8bc37326"
	token := "cd6007cd-57f5-47e6-b3b2-e296fe1c4cba"

	c.Us.UserFn = func(e string) (*pinub.User, error) {
		pwd, _ := pinub.HashPassword(password)
		return &pinub.User{ID: id, Email: e, Password: pwd}, nil
	}

	c.Us.AddTokenFn = func(u *pinub.User) error {
		u.Token = token
		return nil
	}

	data := url.Values{"email": {email}, "password": {password}}
	req, err := http.NewRequest("POST", signinURL, strings.NewReader(data.Encode()))
	ok(t, err)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	h.ServeHTTP(r, req)

	equals(t, http.StatusSeeOther, r.Code)
	equals(t, true, c.Us.UserInvoked)
	equals(t, true, c.Us.AddTokenInvoked)

	assert(t,
		strings.Contains(r.Header().Get("Set-Cookie"), "keks="+token),
		"cookie should be set")
}

func TestHTTP_SigninPostShortPassword(t *testing.T) {
	t.Parallel()
	h, c, r := setUp()

	email := "test@test.test"
	password := "bla"

	data := url.Values{"email": {email}, "password": {password}}
	req, err := http.NewRequest("POST", signinURL, strings.NewReader(data.Encode()))
	ok(t, err)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	h.ServeHTTP(r, req)

	equals(t, http.StatusOK, r.Code)
	equals(t, false, c.Us.UserInvoked)

	assert(t,
		strings.Contains(r.Body.String(), pinub.ErrUserShortPassword.Error()),
		"error message about short password should appear")
}

func TestHTTP_SigninPostInvalidEmail(t *testing.T) {
	t.Parallel()
	h, c, r := setUp()

	email := "test@test"
	password := "blah"

	data := url.Values{"email": {email}, "password": {password}}
	req, err := http.NewRequest("POST", signinURL, strings.NewReader(data.Encode()))
	ok(t, err)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	h.ServeHTTP(r, req)

	equals(t, http.StatusOK, r.Code)
	equals(t, false, c.Us.UserInvoked)

	assert(t,
		strings.Contains(r.Body.String(), pinub.ErrUserInvalidEmail.Error()),
		"error message about email format should appear")
}

func TestHTTP_SigninPostNoUser(t *testing.T) {
	t.Parallel()
	h, c, r := setUp()

	email := "test@test.test"
	password := "blahblah"

	c.Us.UserFn = func(e string) (*pinub.User, error) {
		return nil, errors.New("")
	}

	data := url.Values{"email": {email}, "password": {password}}
	req, err := http.NewRequest("POST", signinURL, strings.NewReader(data.Encode()))
	ok(t, err)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	h.ServeHTTP(r, req)

	equals(t, http.StatusOK, r.Code)
	equals(t, true, c.Us.UserInvoked)

	assert(t,
		strings.Contains(r.Body.String(), pinub.ErrUserPasswordNotCorrect.Error()),
		"error message about email or password should appear")
}

func TestHTTP_SigninPostWrongPassword(t *testing.T) {
	t.Parallel()
	h, c, r := setUp()

	email := "test@test.test"
	password := "blahblah"
	id := "64f1001c-7992-470d-8e91-c64f8bc37326"

	c.Us.UserFn = func(e string) (*pinub.User, error) {
		pwd, _ := pinub.HashPassword("different password")
		return &pinub.User{ID: id, Email: e, Password: pwd}, nil
	}

	data := url.Values{"email": {email}, "password": {password}}
	req, err := http.NewRequest("POST", signinURL, strings.NewReader(data.Encode()))
	ok(t, err)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	h.ServeHTTP(r, req)

	equals(t, http.StatusOK, r.Code)
	equals(t, true, c.Us.UserInvoked)

	assert(t,
		strings.Contains(r.Body.String(), pinub.ErrUserPasswordNotCorrect.Error()),
		"error message about email or password should appear")
}

func TestHTTP_SigninPostCannotCreateToken(t *testing.T) {
	t.Parallel()
	h, c, r := setUp()

	email := "test@test.test"
	password := "blahblah"
	id := "64f1001c-7992-470d-8e91-c64f8bc37326"

	c.Us.UserFn = func(e string) (*pinub.User, error) {
		pwd, _ := pinub.HashPassword(password)
		return &pinub.User{ID: id, Email: e, Password: pwd}, nil
	}

	c.Us.AddTokenFn = func(u *pinub.User) error {
		return errors.New("")
	}

	data := url.Values{"email": {email}, "password": {password}}
	req, err := http.NewRequest("POST", signinURL, strings.NewReader(data.Encode()))
	ok(t, err)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	h.ServeHTTP(r, req)

	equals(t, http.StatusOK, r.Code)
	equals(t, true, c.Us.UserInvoked)
	equals(t, true, c.Us.AddTokenInvoked)

	assert(t,
		strings.Contains(r.Body.String(), "Cannot sign in user"),
		"error message that sign in was not successful")
}
