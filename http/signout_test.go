package http_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/pinub/pinub"
)

const signoutURL = "/signout"

func TestHTTP_SignoutGet(t *testing.T) {
	t.Parallel()
	h, c, r := setUp()

	token := "2c316f20-ed46-4dab-9f76-96d4ea3a4bc7"

	c.Us.UserByTokenFn = func(t string) (*pinub.User, error) {
		return &pinub.User{ID: "", Email: "", Password: "", Token: token}, nil
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

	req, err := http.NewRequest("GET", signoutURL, nil)
	req.AddCookie(cookie)
	ok(t, err)

	h.ServeHTTP(r, req)
	equals(t, r.Code, http.StatusSeeOther)

	t.Log(r.Result().Cookies())
	equals(t, true, c.Us.UserByTokenInvoked)
	// cannot test because this runs in a goroutine
	//equals(t, true, c.Us.RefreshTokenInvoked)
}
