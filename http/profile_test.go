package http_test

import (
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/pinub/pinub"
)

const profileURL = "/profile"

func TestHTTP_ProfileNoAuthGet(t *testing.T) {
	t.Parallel()
	h, _, r := setUp()

	req, err := http.NewRequest("GET", profileURL, nil)
	ok(t, err)

	h.ServeHTTP(r, req)
	equals(t, r.Code, http.StatusSeeOther)
}

func TestHTTP_ProfileGet(t *testing.T) {
	t.Parallel()
	h, c, r := setUp()

	email := "test@example.com"
	id := "42"
	password := "4242"
	token := "b96562b1-fcfa-40a8-8f14-3f09005f9f43"

	c.Us.UserByTokenFn = func(t string) (*pinub.User, error) {
		return &pinub.User{
			ID:       id,
			Email:    email,
			Password: password,
			Token:    token,
		}, nil
	}

	c.Us.RefreshTokenFn = func(u *pinub.User) error {
		return nil
	}

	cookie := &http.Cookie{
		Name:    "keks",
		Value:   token,
		Path:    "/",
		Expires: time.Now().Add(14 * time.Hour * 24),
	}

	req, err := http.NewRequest("GET", profileURL, nil)
	req.AddCookie(cookie)
	ok(t, err)

	h.ServeHTTP(r, req)
	equals(t, r.Code, http.StatusOK)

	equals(t, true, c.Us.UserByTokenInvoked)
	assert(t, strings.Contains(r.Body.String(), email), "email should be prefilled")
	// cannot test because this runs in a goroutine
	//equals(t, true, c.Us.RefreshTokenInvoked)
}

func TestHTTP_ProfileChangeEmail(t *testing.T) {
	t.Parallel()
	h, c, r := setUp()

	email := "test@example.com"
	email_new := "new@example.com"
	id := "42"
	password := "4242"
	token := "8f9020a9-b4c0-48b5-9951-765cb7a69419"

	c.Us.UserByTokenFn = func(t string) (*pinub.User, error) {
		return &pinub.User{
			ID:       id,
			Email:    email,
			Password: password,
			Token:    token,
		}, nil
	}

	c.Us.RefreshTokenFn = func(u *pinub.User) error {
		return nil
	}

	cookie := &http.Cookie{
		Name:    "keks",
		Value:   token,
		Path:    "/",
		Expires: time.Now().Add(14 * time.Hour * 24),
	}

	data := url.Values{
		"email": {email_new},
	}
	req, err := http.NewRequest("PUT", profileURL, strings.NewReader(data.Encode()))
	req.AddCookie(cookie)
	ok(t, err)

	h.ServeHTTP(r, req)
	equals(t, r.Code, http.StatusOK)

	equals(t, true, c.Us.UserByTokenInvoked)
	// cannot test because this runs in a goroutine
	//equals(t, true, c.Us.RefreshTokenInvoked)
	t.Log(r.Body.String())
	//assert(t, !strings.Contains(r.Body.String(), "error"), "No error should apear")
}
