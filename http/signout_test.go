package http_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/pinub/pinub"
	"github.com/pinub/pinub/claim"
)

const signoutURL = "/signout"

func TestSignout(t *testing.T) {
	t.Parallel()

	userToken := "2c316f20-ed46-4dab-9f76-96d4ea3a4bc7"

	userByTokenFn := func(token string) (*pinub.User, error) {
		return &pinub.User{ID: "", Email: "", Password: "", Token: token}, nil
	}
	userRefreshTokenFn := func(*pinub.User) error {
		return nil
	}

	cookie := &http.Cookie{
		Name:    "keks",
		Value:   userToken,
		Path:    "/",
		Expires: time.Now().Add(14 * time.Hour * 24),
	}

	t.Run("signout", func(t *testing.T) {
		handler, client, rec := setUp()
		client.Us.UserByTokenFn = userByTokenFn
		client.Us.RefreshTokenFn = userRefreshTokenFn

		req := httptest.NewRequest("GET", signoutURL, nil)
		req.AddCookie(cookie)
		handler.ServeHTTP(rec, req)

		claim.Equals(t, 303, rec.Code)
		claim.Equals(t, true, client.Us.UserByTokenInvoked)
		claim.Equals(t, false, client.Us.RefreshTokenInvoked)
	})
}
