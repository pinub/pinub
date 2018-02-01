package http_test

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/pinub/pinub"
	"github.com/pinub/pinub/claim"
)

const linksURL = "/"

func TestHomepage(t *testing.T) {
	t.Parallel()

	userToken := "e12e72dc-ea60-4532-8f23-6c8ccaa26fb7"
	linkIDs := []string{"10", "12"}
	count := 0

	userFn := func(string) (*pinub.User, error) {
		return &pinub.User{}, nil
	}
	userByTokenFn := func(string) (*pinub.User, error) {
		return &pinub.User{Token: userToken}, nil
	}
	refreshTokenFn := func(*pinub.User) error {
		return nil
	}
	linksFn := func(*pinub.User) ([]pinub.Link, error) {
		return []pinub.Link{}, nil
	}
	deleteLinkFn := func(l *pinub.Link, _ *pinub.User) error {
		claim.Equals(t, l.ID, linkIDs[count])
		count = count + 1
		return nil
	}

	t.Run("show unauthorized homepage", func(t *testing.T) {
		handler, client, rec := setUp()
		client.Us.UserFn = userFn

		req := httptest.NewRequest("GET", linksURL, nil)
		handler.ServeHTTP(rec, req)

		claim.Equals(t, 200, rec.Code)
		claim.Equals(t, false, client.Us.UserByTokenInvoked)
		claim.Equals(t, false, client.Ls.LinksInvoked)
	})

	t.Run("show links", func(t *testing.T) {
		handler, client, rec := setUp()
		client.Us.UserByTokenFn = userByTokenFn
		client.Us.RefreshTokenFn = refreshTokenFn
		client.Ls.LinksFn = linksFn

		req := httptest.NewRequest("GET", linksURL, nil)
		req.Header.Add("Cookie", "keks="+userToken)
		handler.ServeHTTP(rec, req)

		claim.Equals(t, 200, rec.Code)
		claim.Equals(t, true, client.Us.UserByTokenInvoked)
		claim.Equals(t, true, client.Ls.LinksInvoked)
		claim.Equals(t, false, client.Ls.DeleteLinkInvoked)
	})

	t.Run("delete links", func(t *testing.T) {
		handler, client, rec := setUp()
		client.Us.UserByTokenFn = userByTokenFn
		client.Us.RefreshTokenFn = refreshTokenFn
		client.Ls.LinksFn = linksFn
		client.Ls.DeleteLinkFn = deleteLinkFn

		req := httptest.NewRequest("GET", linksURL, nil)
		req.Header.Add("Cookie", "keks="+userToken)
		req.Header.Add("Cookie", "deleteMe="+strings.Join(linkIDs, ","))
		handler.ServeHTTP(rec, req)

		claim.Equals(t, 200, rec.Code)
		claim.Equals(t, true, client.Us.UserByTokenInvoked)
		claim.Equals(t, true, client.Ls.LinksInvoked)
		claim.Equals(t, true, client.Ls.DeleteLinkInvoked)
	})
}
