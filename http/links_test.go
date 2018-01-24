package http_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/pinub/pinub"
)

const linksURL = "/"

func TestHTTP_IndexGet(t *testing.T) {
	t.Parallel()
	h, c, r := setUp()

	c.Us.UserByTokenFn = func(t string) (*pinub.User, error) {
		return &pinub.User{}, nil
	}

	req, err := http.NewRequest("GET", linksURL, nil)
	ok(t, err)

	h.ServeHTTP(r, req)
	equals(t, http.StatusOK, r.Code)
	equals(t, false, c.Us.UserByTokenInvoked)
}

func TestHTTP_IndexAuthGet(t *testing.T) {
	t.Parallel()
	h, c, r := setUp()

	token := "e12e72dc-ea60-4532-8f23-6c8ccaa26fb7"

	c.Us.UserByTokenFn = func(to string) (*pinub.User, error) {
		equals(t, to, token)
		return &pinub.User{Token: to}, nil
	}
	c.Us.RefreshTokenFn = func(u *pinub.User) error {
		return nil
	}
	c.Ls.LinksFn = func(_ *pinub.User) ([]pinub.Link, error) {
		return []pinub.Link{}, nil
	}

	req, err := http.NewRequest("GET", linksURL, nil)
	req.Header.Add("Cookie", "keks="+token)
	ok(t, err)

	h.ServeHTTP(r, req)
	equals(t, http.StatusOK, r.Code)
	equals(t, true, c.Us.UserByTokenInvoked)
	equals(t, true, c.Ls.LinksInvoked)
	equals(t, false, c.Ls.DeleteLinkInvoked)
}

func TestHTTP_IndexDeleteLinks(t *testing.T) {
	t.Parallel()
	h, c, r := setUp()

	token := "e12e72dc-ea60-4532-8f23-6c8ccaa26fb7"
	linkIDs := []string{"10", "12"}
	count := 0

	c.Us.UserByTokenFn = func(to string) (*pinub.User, error) {
		equals(t, to, token)
		return &pinub.User{Token: to}, nil
	}
	c.Us.RefreshTokenFn = func(u *pinub.User) error {
		return nil
	}
	c.Ls.LinksFn = func(_ *pinub.User) ([]pinub.Link, error) {
		return []pinub.Link{}, nil
	}
	c.Ls.DeleteLinkFn = func(l *pinub.Link, _ *pinub.User) error {
		equals(t, l.ID, linkIDs[count])
		count = count + 1
		return nil
	}

	req, err := http.NewRequest("GET", linksURL, nil)
	req.Header.Add("Cookie", "keks="+token)
	req.Header.Add("Cookie", "deleteMe="+strings.Join(linkIDs, ","))
	ok(t, err)

	h.ServeHTTP(r, req)
	equals(t, http.StatusOK, r.Code)
	equals(t, true, c.Us.UserByTokenInvoked)
	equals(t, true, c.Ls.LinksInvoked)
	equals(t, true, c.Ls.DeleteLinkInvoked)
}
